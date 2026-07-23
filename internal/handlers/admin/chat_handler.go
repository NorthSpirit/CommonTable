package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/chatservices"
	"CommonTable/internal/services/restaurantservices"
	"CommonTable/internal/templates/admin/adminshared"
	"CommonTable/internal/templates/admin/chat"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ChatHandler struct {
	chatMessageService *chatservices.ChatMessageService
	restaurantService  *restaurantservices.RestaurantService
}

func NewChatHandler(chatMessageService *chatservices.ChatMessageService, restaurantService *restaurantservices.RestaurantService) *ChatHandler {
	return &ChatHandler{chatMessageService: chatMessageService, restaurantService: restaurantService}
}

// Page handles GET /admin/restaurants/{id}/chat
func (h *ChatHandler) Page(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), int32(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	messages, err := h.chatMessageService.ListMessagesByRestaurantWithSender(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "failed to load messages", http.StatusInternalServerError)
		return
	}

	views := make([]adminshared.ChatMessageView, 0, len(messages))
	for _, m := range messages {
		views = append(views, adminshared.ChatMessageView{
			ID:         m.ID,
			SenderName: m.SenderName,
			Message:    m.Message,
			SentAt:     helpers.FormatSentAt(m.SentAt),
		})
	}

	if err := chat.ChatPage(int32(id), restaurant.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

type chatMessageJSON struct {
	ID         int32  `json:"id"`
	SenderName string `json:"sender_name"`
	Message    string `json:"message"`
	SentAt     string `json:"sent_at"`
}

// Poll handles GET /admin/restaurants/{id}/chat/messages?after={id}
func (h *ChatHandler) Poll(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid restaurant id", http.StatusBadRequest)
		return
	}

	afterID, _ := strconv.Atoi(r.URL.Query().Get("after"))

	messages, err := h.chatMessageService.ListMessagesByRestaurantAfterID(r.Context(), int32(id), int32(afterID))
	if err != nil {
		http.Error(w, "failed to load messages", http.StatusInternalServerError)
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

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
