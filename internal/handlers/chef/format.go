package chef

import "github.com/jackc/pgx/v5/pgtype"

func formatSentAt(t pgtype.Timestamp) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("2006-01-02 15:04")
}
