package helpers

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// PgUserMessage maps known Postgres errors to a message safe to show the user.
// Returns "" for anything else (caller should still return a 500).
func PgUserMessage(err error) string {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return ""
	}
	switch pgErr.Code {
	case "22001":
		return "One of the values is too long."
	case "22003":
		return "One of the numbers is out of the allowed range."
	case "23505":
		return "A record with this name or number already exists."
	}
	return ""
}
