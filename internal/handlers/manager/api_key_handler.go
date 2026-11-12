package manager

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/externalorderservices"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"

	"github.com/go-chi/chi/v5"
)

var apiKeyMessages = map[string]string{
	"revoked": "API key revoked. It no longer works.",
	"stale":   "That key does not exist or was already revoked.",
}

type ApiKeyHandler struct {
	keys              *externalorderservices.ApiKeyService
	restaurantService *restaurantservices.RestaurantService
}

func NewApiKeyHandler(keys *externalorderservices.ApiKeyService, restaurantService *restaurantservices.RestaurantService) *ApiKeyHandler {
	return &ApiKeyHandler{keys: keys, restaurantService: restaurantService}
}

func (h *ApiKeyHandler) render(w http.ResponseWriter, r *http.Request, message string, errorMessage string, newKey string) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	keys, err := h.keys.ListByRestaurant(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load api keys", http.StatusInternalServerError)
		return
	}

	views := make([]managertmpl.ApiKeyView, 0, len(keys))
	for _, k := range keys {
		views = append(views, managertmpl.ApiKeyView{
			ID:        k.ID,
			Label:     k.Label,
			Prefix:    k.KeyPrefix,
			CreatedAt: helpers.FormatSentAt(k.CreatedAt),
			Active:    !k.RevokedAt.Valid,
		})
	}

	if err := managertmpl.ApiKeysPage(restaurant.Name, claims.Name, message, errorMessage, newKey, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// List handles GET /manager/api-keys
func (h *ApiKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, apiKeyMessages[r.URL.Query().Get("done")], "", "")
}

// Create handles POST /manager/api-keys
// The new key is shown on the page that this request returns, once. It is not stored in plain form.
func (h *ApiKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	label := strings.TrimSpace(r.FormValue("label"))
	if label == "" {
		h.render(w, r, "", "Label is required.", "")
		return
	}
	if msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Label", Value: label, Max: 255}); msg != "" {
		h.render(w, r, "", msg, "")
		return
	}

	plain, _, err := h.keys.Create(r.Context(), claims.RestaurantID, label)
	if err != nil {
		http.Error(w, "failed to create api key", http.StatusInternalServerError)
		return
	}

	h.render(w, r, "", "", plain)
}

// Revoke handles POST /manager/api-keys/{keyID}/revoke
func (h *ApiKeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "keyID"))
	if err != nil {
		http.Error(w, "invalid key id", http.StatusBadRequest)
		return
	}

	if err := h.keys.Revoke(r.Context(), int32(id), claims.RestaurantID); err != nil {
		if errors.Is(err, externalorderservices.ErrApiKeyNotFound) {
			http.Redirect(w, r, "/manager/api-keys?done=stale", http.StatusSeeOther)
			return
		}
		http.Error(w, "failed to revoke api key", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/api-keys?done=revoked", http.StatusSeeOther)
}
