package manager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/services/chatservices"
	"CommonTable/internal/services/restaurantservices"
	managertmpl "CommonTable/internal/templates/manager"
)

type ChatHandler struct {
	chatMessageService *chatservices.ChatMessageService
	restaurantService  *restaurantservices.RestaurantService
}

func NewChatHandler(chatMessageService *chatservices.ChatMessageService, restaurantService *restaurantservices.RestaurantService) *ChatHandler {
	return &ChatHandler{chatMessageService: chatMessageService, restaurantService: restaurantService}
}

// Page handles GET /manager/chat
func (h *ChatHandler) Page(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	restaurant, err := h.restaurantService.GetRestaurantByID(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load restaurant", http.StatusInternalServerError)
		return
	}

	messages, err := h.chatMessageService.ListMessagesByRestaurantWithSender(r.Context(), claims.RestaurantID)
	if err != nil {
		http.Error(w, "failed to load messages", http.StatusInternalServerError)
		return
	}

	views := make([]managertmpl.ChatMessageView, 0, len(messages))
	for _, m := range messages {
		views = append(views, managertmpl.ChatMessageView{
			ID:         m.ID,
			SenderName: m.SenderName,
			Message:    m.Message,
			SentAt:     helpers.FormatSentAt(m.SentAt),
		})
	}

	if err := managertmpl.ChatPage(restaurant.Name, claims.Name, views).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// Send handles POST /manager/chat/send
func (h *ChatHandler) Send(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

	message := r.FormValue("message")
	if message == "" {
		http.Redirect(w, r, "/manager/chat", http.StatusFound)
		return
	}

	now := time.Now()
	if _, err := h.chatMessageService.InsertChatMessage(r.Context(), claims.RestaurantID, claims.EmployeeID, message, &now); err != nil {
		http.Error(w, "failed to send message", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/manager/chat", http.StatusFound)
}

type chatMessageJSON struct {
	ID         int32  `json:"id"`
	SenderName string `json:"sender_name"`
	Message    string `json:"message"`
	SentAt     string `json:"sent_at"`
}

// Poll handles GET /manager/chat/messages?after={id}
func (h *ChatHandler) Poll(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromContext(r.Context())

	afterID, _ := strconv.Atoi(r.URL.Query().Get("after"))

	messages, err := h.chatMessageService.ListMessagesByRestaurantAfterID(r.Context(), claims.RestaurantID, int32(afterID))
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
