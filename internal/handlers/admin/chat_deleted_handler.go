package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CommonTable/internal/services/chatservices"

	"github.com/go-chi/chi/v5"
)

type deletedMessageJSON struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
}

func (h *ChatHandler) Deleted(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	rows, err := h.chatMessageService.ListDeletedMessages(r.Context(), int32(id))
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
