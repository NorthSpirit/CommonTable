package customer

import "net/http"

// GetLocale reads the desired language from the ?lang= query param.
// Only "en" and "alt" are recognized; anything else (including absent)
// falls back to "en".
func GetLocale(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if lang == "alt" {
		return "alt"
	}
	return "en"
}
