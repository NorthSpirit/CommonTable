package helpers

import "github.com/jackc/pgx/v5/pgtype"

func Int4FromInt32(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}
