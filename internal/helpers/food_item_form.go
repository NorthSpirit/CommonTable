package helpers

import (
	"net/http"
	"strconv"
)

const (
	// MaxFoodItemPrice is a business cap, well inside numeric(10,2).
	MaxFoodItemPrice = 100000.0
	// MaxPrepTimeMins is a business cap for preparation time.
	MaxPrepTimeMins = 1000
)

// ValidateFoodItemForm checks the shared food item fields (name lengths,
// price, prep time) of an already-parsed form. Returns "" if all fine.
// Empty price / prep time are allowed here; callers keep their old handling.
func ValidateFoodItemForm(r *http.Request) string {
	if msg := CheckLengths(
		FieldLimit{Label: "Name", Value: r.FormValue("name"), Max: 255},
		FieldLimit{Label: "Name (alt)", Value: r.FormValue("name_alt"), Max: 255},
	); msg != "" {
		return msg
	}

	if raw := r.FormValue("price"); raw != "" {
		price, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return "Price must be a valid number."
		}
		if msg := CheckRange("Price", price, 0, MaxFoodItemPrice); msg != "" {
			return msg
		}
	}

	if raw := r.FormValue("prep_time_mins"); raw != "" {
		mins, err := strconv.Atoi(raw)
		if err != nil {
			return "Prep time must be a whole number."
		}
		if msg := CheckRange("Prep time", float64(mins), 0, MaxPrepTimeMins); msg != "" {
			return msg
		}
	}

	return ""
}

// SafeFormPrice returns the submitted price, or 0 if it is missing, invalid
// or out of range. Used only to re-fill the form after a validation error.
func SafeFormPrice(r *http.Request) float64 {
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil || CheckRange("Price", price, 0, MaxFoodItemPrice) != "" {
		return 0
	}
	return price
}

// SafeFormPrepTime returns the submitted prep time, or 0 if it is missing,
// invalid or out of range. Used only to re-fill the form after an error.
func SafeFormPrepTime(r *http.Request) int32 {
	mins, err := strconv.Atoi(r.FormValue("prep_time_mins"))
	if err != nil || mins < 0 || mins > MaxPrepTimeMins {
		return 0
	}
	return int32(mins)
}
