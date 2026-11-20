package setup

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"CommonTable/internal/dbmigrate"
	"CommonTable/internal/helpers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var identifier = regexp.MustCompile(`^[a-z][a-z0-9_]{0,30}$`)

// Run handles "commontable.exe setup" and returns the process exit code (0 = success).
// The PostgreSQL superuser password is read from COMMONTABLE_SUPERUSER_PASSWORD and the
// admin password from COMMONTABLE_ADMIN_PASSWORD, not from arguments.
func Run(args []string) int {
	fs := flag.NewFlagSet("setup", flag.ContinueOnError)
	host := fs.String("host", "localhost", "PostgreSQL host")
	port := fs.Int("port", 5432, "PostgreSQL port")
	superuser := fs.String("superuser", "postgres", "PostgreSQL login that may create roles and databases")
	dbName := fs.String("db", "commontable", "database to create")
	role := fs.String("role", "commontable", "login (role) to create, owner of the database")
	adminUser := fs.String("admin", "", "username of the CommonTable admin")
	envPath := fs.String("env", ".env", "where to write the settings file")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	superPassword := os.Getenv("COMMONTABLE_SUPERUSER_PASSWORD")
	if superPassword == "" {
		return fail("COMMONTABLE_SUPERUSER_PASSWORD is not set")
	}
	adminPassword := os.Getenv("COMMONTABLE_ADMIN_PASSWORD")
	if *adminUser == "" || adminPassword == "" {
		return fail("the admin username (-admin) and COMMONTABLE_ADMIN_PASSWORD are required")
	}
	if len(adminPassword) < 8 || len(adminPassword) > 72 {
		return fail("the admin password must be 8 to 72 characters long")
	}
	if !identifier.MatchString(*dbName) || !identifier.MatchString(*role) {
		return fail("database and role names may only use lowercase letters, digits and underscores, and must start with a letter")
	}
	if _, err := os.Stat(*envPath); err == nil {
		return fail(fmt.Sprintf("%s already exists. Nothing was changed. Remove it first if you really want a new setup", *envPath))
	}

	adminHash, err := helpers.HashPassword(adminPassword)
	if err != nil {
		return fail("could not hash the admin password: " + err.Error())
	}
	jwtSecret, err := randomBase64(32)
	if err != nil {
		return fail("could not generate the token secret: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	admin, err := pgx.Connect(ctx, connString(*superuser, superPassword, *host, *port, "postgres"))
	if err != nil {
		return fail(describeConnectError(err, *host, *port, *superuser))
	}
	defer admin.Close(context.Background())

	var isSuperuser bool
	if err := admin.QueryRow(ctx, `SELECT rolsuper FROM pg_roles WHERE rolname = current_user`).Scan(&isSuperuser); err != nil {
		return fail("could not check the permissions of " + *superuser + ": " + err.Error())
	}
	if !isSuperuser {
		return fail(fmt.Sprintf("PostgreSQL login %q is not a superuser, so it cannot create the database", *superuser))
	}

	var dbExists, roleExists bool
	err = admin.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1), EXISTS (SELECT 1 FROM pg_roles WHERE rolname = $2)`,
		*dbName, *role,
	).Scan(&dbExists, &roleExists)
	if err != nil {
		return fail("could not check for existing databases: " + err.Error())
	}
	if dbExists {
		return fail(fmt.Sprintf("a database named %q already exists. Nothing was changed. Choose another name with -db", *dbName))
	}
	if roleExists {
		return fail(fmt.Sprintf("a PostgreSQL login named %q already exists. Nothing was changed. Choose another name with -role", *role))
	}

	rolePassword, err := randomHex(24)
	if err != nil {
		return fail("could not generate a password: " + err.Error())
	}

	quotedRole := pgx.Identifier{*role}.Sanitize()
	quotedDB := pgx.Identifier{*dbName}.Sanitize()

	if _, err := admin.Exec(ctx, fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s'", quotedRole, rolePassword)); err != nil {
		return fail("could not create the login: " + err.Error())
	}
	if _, err := admin.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s OWNER %s", quotedDB, quotedRole)); err != nil {
		undo(admin, quotedRole, "")
		return fail("could not create the database: " + err.Error())
	}

	appConn := connString(*role, rolePassword, *host, *port, *dbName)
	if err := dbmigrate.Up(appConn); err != nil {
		undo(admin, quotedRole, quotedDB)
		return fail("could not set up the tables, everything was removed again: " + err.Error())
	}

	if err := writeEnv(*envPath, appConn, *adminUser, adminHash, jwtSecret); err != nil {
		undo(admin, quotedRole, quotedDB)
		return fail("could not write the settings file, everything was removed again: " + err.Error())
	}

	fmt.Printf("setup: database %q and login %q created, tables are ready\n", *dbName, *role)
	fmt.Printf("setup: settings written to %s, admin account %q\n", *envPath, *adminUser)
	return 0
}

func fail(message string) int {
	fmt.Fprintln(os.Stderr, "setup failed: "+message)
	return 1
}

// writeEnv creates the settings file. It fails if the file already exists.
// The hash is in single quotes, because it contains $ characters.
func writeEnv(path string, dbString string, adminUser string, adminHash string, jwtSecret string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	var b strings.Builder
	b.WriteString("GOOSE_DBSTRING=" + dbString + "\n\n")
	b.WriteString("ADMIN_USERNAME=" + adminUser + "\n")
	b.WriteString("ADMIN_PASSWORD_HASH='" + adminHash + "'\n\n")
	b.WriteString("JWT_SECRET=" + jwtSecret + "\n")
	b.WriteString("JWT_TOKEN_TTL_HOURS=12\n\n")
	b.WriteString("TRANSLATIONS_DIR=./locales\n")
	b.WriteString("UPLOADS_DIR=./uploads\n")
	b.WriteString("ALT_LANGUAGE_LABEL=LT\n\n")
	b.WriteString("DEPLOYMENT_MODE=lan\n")
	b.WriteString("PUBLIC_BASE_URL=\n")

	if _, err := file.WriteString(b.String()); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return err
	}
	return nil
}

func connString(user string, password string, host string, port int, database string) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, strconv.Itoa(port)),
		Path:     "/" + database,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

func describeConnectError(err error, host string, port int, user string) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "28P01":
			return fmt.Sprintf("wrong password for PostgreSQL login %q", user)
		case "28000":
			return fmt.Sprintf("PostgreSQL does not accept the login %q, check that it exists", user)
		}
		return "PostgreSQL refused the login: " + pgErr.Message
	}
	return fmt.Sprintf("cannot reach PostgreSQL at %s:%d. Check that it is running and that the host and port are right", host, port)
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func randomBase64(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// undo removes what this run created. An empty quotedDB means the database was not created.
func undo(admin *pgx.Conn, quotedRole string, quotedDB string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if quotedDB != "" {
		_, _ = admin.Exec(ctx, "DROP DATABASE IF EXISTS "+quotedDB)
	}
	_, _ = admin.Exec(ctx, "DROP ROLE IF EXISTS "+quotedRole)
}
