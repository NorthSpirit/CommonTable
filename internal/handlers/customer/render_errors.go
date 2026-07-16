package customer

import (
	"net/http"

	"CommonTable/internal/templates/errorpages"
)

func renderNotFound(w http.ResponseWriter, r *http.Request) {
	locale := GetLocale(r)
	w.WriteHeader(http.StatusNotFound)
	_ = errorpages.NotFoundPage(locale).Render(r.Context(), w)
}

func renderServerError(w http.ResponseWriter, r *http.Request) {
	locale := GetLocale(r)
	w.WriteHeader(http.StatusInternalServerError)
	_ = errorpages.ErrorPage(locale).Render(r.Context(), w)
}
