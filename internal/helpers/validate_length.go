package helpers

import (
	"fmt"
	"math"
	"unicode/utf8"
)

type FieldLimit struct {
	Label string
	Value string
	Max   int
}

// CheckLengths returns an error message for the first field that exceeds
// its limit, or "" if all fields fit. Postgres varchar(n) counts characters,
// so this counts runes, not bytes.
// Since function can be reused for different stuff, it's more flexible to pass limit too
func CheckLengths(fields ...FieldLimit) string {
	for _, f := range fields {
		// Because Go's len(string) returns bytes, not chars,
		// meanwhile, PostgreSQL's varchar(n) limits chars, so RuneCountInString must be used.
		if utf8.RuneCountInString(f.Value) > f.Max {
			return fmt.Sprintf("%s must be at most %d characters.", f.Label, f.Max)
		}
	}
	return ""
}

// ValidateSecret checks a PIN/password against bcrypt's 72-byte limit.
// Returns "" if fine. Empty is allowed (callers decide if it's required).
func ValidateSecret(label, value string) string {
	// As mentioned before, len checks for bytes, which is legit for bcrypt.
	if len(value) > 72 {
		return label + " must be at most 72 characters."
	}
	return ""
}

// CheckRange returns an error message if v is NaN/Inf or outside [min, max].
// Works for integers too: CheckRange("Capacity", float64(capacity), 1, 1000).
func CheckRange(label string, v, min, max float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return label + " must be a valid number."
	}
	if v < min || v > max {
		return fmt.Sprintf("%s must be between %g and %g.", label, min, max)
	}
	return ""
}
