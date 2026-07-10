package helpers

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

// NumericFromFloat64 converts a required float64 into a pgtype.Numeric.
// Use for NOT NULL numeric columns.
func NumericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric

	// pgtype.Numeric.Scan accepts a decimal string; formatting via
	// strconv avoids float imprecision creeping into the DB value.
	if err := n.Scan(strconv.FormatFloat(f, 'f', -1, 64)); err != nil {
		// Scan only fails on malformed input, which strconv.FormatFloat
		// cannot produce, so this path should be unreachable.
		return pgtype.Numeric{Valid: false}
	}

	return n
}

// NumericFromFloat64Ptr converts an optional *float64 into a pgtype.Numeric.
// Use for nullable numeric columns; a nil pointer becomes SQL NULL.
func NumericFromFloat64Ptr(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}

	return NumericFromFloat64(*f)
}

// NumericToFloat64 converts a pgtype.Numeric back into a *float64.
// Returns nil if the value is SQL NULL. Useful when reading rows back
// out of the repository layer for JSON responses.
func NumericToFloat64(n pgtype.Numeric) (*float64, error) {
	if !n.Valid {
		return nil, nil
	}

	f, err := n.Float64Value()
	if err != nil {
		return nil, err
	}

	if !f.Valid {
		return nil, nil
	}

	return &f.Float64, nil
}
