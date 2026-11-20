package dbmigrate

import (
	"database/sql"
	"fmt"

	commontable "CommonTable"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Up applies every migration that is not applied yet to the database in connString.
// The migrations are inside the program (embedded), so no .sql files are needed next to it.
func Up(connString string) error {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(commontable.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
