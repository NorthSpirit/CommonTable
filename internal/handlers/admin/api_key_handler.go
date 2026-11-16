package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/externalorderservices"
	"CommonTable/internal/services/restaurantservices"
	adminrestaurant "CommonTable/internal/templates/admin/restaurant"

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

func (h *ApiKeyHandler) render(w http.ResponseWriter, r *http.Request, restaurantID int32, message string, errorMessage string, newKey string) {
	rest, err := h.restaurantService.GetRestaurantByID(r.Context(), restaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	keys, err := h.keys.ListByRestaurant(r.Context(), restaurantID)
	if err != nil {
		http.Error(w, "failed to load api keys", http.StatusInternalServerError)
		return
	}

	views := make([]adminrestaurant.ApiKeyView, 0, len(keys))
	for _, k := range keys {
		views = append(views, adminrestaurant.ApiKeyView{
			ID:        k.ID,
			Label:     k.Label,
			Prefix:    k.KeyPrefix,
			CreatedAt: helpers.FormatSentAt(k.CreatedAt),
			Active:    !k.RevokedAt.Valid,
		})
	}

	if err := adminrestaurant.ApiKeysPage(restaurantID, rest.Name, message, errorMessage, newKey, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

func restaurantIDFromURL(r *http.Request) (int32, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, false
	}
	return int32(id), true
}

func (h *ApiKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	restaurantID, ok := restaurantIDFromURL(r)
	if !ok {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	h.render(w, r, restaurantID, apiKeyMessages[r.URL.Query().Get("done")], "", "")
}

func (h *ApiKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	restaurantID, ok := restaurantIDFromURL(r)
	if !ok {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	label := strings.TrimSpace(r.FormValue("label"))
	if label == "" {
		h.render(w, r, restaurantID, "", "Label is required.", "")
		return
	}
	if msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Label", Value: label, Max: 255}); msg != "" {
		h.render(w, r, restaurantID, "", msg, "")
		return
	}

	plain, _, err := h.keys.Create(r.Context(), restaurantID, label)
	if err != nil {
		http.Error(w, "failed to create api key", http.StatusInternalServerError)
		return
	}

	h.render(w, r, restaurantID, "", "", plain)
}

func (h *ApiKeyHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	restaurantID, ok := restaurantIDFromURL(r)
	if !ok {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	keyID, err := strconv.Atoi(chi.URLParam(r, "keyID"))
	if err != nil {
		http.Error(w, "invalid key id", http.StatusBadRequest)
		return
	}

	base := "/admin/restaurants/" + strconv.Itoa(int(restaurantID)) + "/api-keys?done="

	if err := h.keys.Revoke(r.Context(), int32(keyID), restaurantID); err != nil {
		if errors.Is(err, externalorderservices.ErrApiKeyNotFound) {
			http.Redirect(w, r, base+"stale", http.StatusSeeOther)
			return
		}
		http.Error(w, "failed to revoke api key", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, base+"revoked", http.StatusSeeOther)
}
