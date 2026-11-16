package waiterapi

import (
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
		writeError(w, http.StatusInternalServerError, "failed to load deleted messages")
		return
	}

	out := make([]deletedMessageJSON, 0, len(rows))
	for _, row := range rows {
		out = append(out, deletedMessageJSON{ID: row.ID, Message: chatservices.DeletedText(row.DeleteReason.String)})
	}

	writeJSON(w, http.StatusOK, out)
}
