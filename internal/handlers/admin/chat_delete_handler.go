package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/chatservices"

	"github.com/go-chi/chi/v5"
)

func (h *ChatHandler) Delete(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(chi.URLParam(r, "messageID"))
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	reason := strings.TrimSpace(r.FormValue("reason"))
	if msg := helpers.CheckLengths(helpers.FieldLimit{Label: "Reason", Value: reason, Max: 255}); msg != "" {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := h.chatMessageService.SoftDeleteMessage(r.Context(), int32(messageID), int32(restaurantID), "Admin", reason); err != nil {
		if errors.Is(err, chatservices.ErrMessageNotFound) {
			http.Error(w, "message not found or already deleted", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete message", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/restaurants/"+strconv.Itoa(restaurantID)+"/chat", http.StatusSeeOther)
}
