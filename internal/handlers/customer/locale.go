package customer

import "net/http"

// getLocale reads the desired language from the ?lang= query param.
// Only "en" and "alt" are recognized; anything else (including absent)
// falls back to "en".
func getLocale(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if lang == "alt" {
		return "alt"
	}
	return "en"
}
