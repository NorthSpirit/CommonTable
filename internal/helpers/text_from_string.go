package helpers

import "github.com/jackc/pgx/v5/pgtype"

func TextFromString(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}
