package waiterapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/chatservices"

	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	chatMessageService *chatservices.ChatMessageService
}

func NewChatHandler(chatMessageService *chatservices.ChatMessageService) *ChatHandler {
	return &ChatHandler{chatMessageService: chatMessageService}
}

type chatMessageJSON struct {
	ID         int32  `json:"id"`
	SenderName string `json:"sender_name"`
	Message    string `json:"message"`
	SentAt     string `json:"sent_at"`
}

// List handles GET /api/waiter/chat/messages
func (h *ChatHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	messages, err := h.chatMessageService.ListMessagesByRestaurantWithSender(r.Context(), claims.RestaurantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load messages")
		return
	}

	out := make([]chatMessageJSON, 0, len(messages))
	for _, m := range messages {
		out = append(out, chatMessageJSON{
			ID:         m.ID,
			SenderName: m.SenderName,
			Message:    m.Message,
			SentAt:     helpers.FormatSentAt(m.SentAt),
		})
	}

	writeJSON(w, http.StatusOK, out)
}

// Poll handles GET /api/waiter/chat/messages/after/{id}
func (h *ChatHandler) Poll(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	afterID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	messages, err := h.chatMessageService.ListMessagesByRestaurantAfterID(r.Context(), claims.RestaurantID, int32(afterID))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load messages")
		return
	}

	out := make([]chatMessageJSON, 0, len(messages))
	for _, m := range messages {
		out = append(out, chatMessageJSON{
			ID:         m.ID,
			SenderName: m.SenderName,
			Message:    m.Message,
			SentAt:     helpers.FormatSentAt(m.SentAt),
		})
	}

	writeJSON(w, http.StatusOK, out)
}

type sendMessageRequest struct {
	Message string `json:"message"`
}

// Send handles POST /api/waiter/chat/send
func (h *ChatHandler) Send(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	now := time.Now()
	if _, err := h.chatMessageService.InsertChatMessage(r.Context(), claims.RestaurantID, claims.EmployeeID, req.Message, &now); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to send message")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]bool{"success": true})
}
