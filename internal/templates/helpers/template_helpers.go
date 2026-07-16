package helpers

import "strconv"

func Itoa(id int32) string {
	return strconv.Itoa(int(id))
}

func FormatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

// WithLang appends the current locale as a query param to an internal path,
// so navigation never silently resets the user back to English.
// English is the default, so it's omitted to keep URLs clean.
func WithLang(path string, locale string) string {
	if locale == "" || locale == "en" {
		return path
	}
	return path + "?lang=" + locale
}
