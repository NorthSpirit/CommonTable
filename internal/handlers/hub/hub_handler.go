package hub

import (
	"net/http"

	"CommonTable/internal/templates/hub"
)

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

func (h *HubHandler) ExportDocs(w http.ResponseWriter, r *http.Request) {
	if err := hub.ExportDocsPage().Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}
