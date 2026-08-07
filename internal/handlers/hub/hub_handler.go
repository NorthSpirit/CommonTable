package hub

import (
	"net/http"

	"CommonTable/internal/templates/hub"
)

// Discover is a minimal endpoint the waiter app's subnet sweep uses to
// confirm "this IP is a CommonTable server" — deliberately tiny and fast,
// since the sweep may probe up to 254 addresses in parallel.
func Discover(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("commontable"))
}

type HubHandler struct{}

func NewHubHandler() *HubHandler {
	return &HubHandler{}
}

func (h *HubHandler) Page(w http.ResponseWriter, r *http.Request) {
	if err := hub.HubPage().Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func (h *HubHandler) ApiDocs(w http.ResponseWriter, r *http.Request) {
	if err := hub.ApiDocsPage().Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}
