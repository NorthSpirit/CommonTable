package helpers

import "github.com/jackc/pgx/v5/pgtype"

func Float8FromFloat64(i *float64) pgtype.Float8 {
	if i == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *i, Valid: true}
}
