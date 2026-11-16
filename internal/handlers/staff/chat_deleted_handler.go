package staff

import (
	"encoding/json"
	"net/http"

	"CommonTable/internal/services/chatservices"
)

type deletedMessageJSON struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
}

func (h *ChatHandler) Deleted(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	rows, err := h.chatMessageService.ListDeletedMessages(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load deleted messages", http.StatusInternalServerError)
		return
	}

	out := make([]deletedMessageJSON, 0, len(rows))
	for _, row := range rows {
		out = append(out, deletedMessageJSON{ID: row.ID, Message: chatservices.DeletedText(row.DeleteReason.String)})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
