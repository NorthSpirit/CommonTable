package chatservices

import (
	"context"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type ChatMessageService struct {
	queries *repository.Queries
}

func NewChatMessageService(queries *repository.Queries) *ChatMessageService {
	return &ChatMessageService{queries: queries}
}

func (f *ChatMessageService) GetChatMessageByID(ctx context.Context, id int32) (repository.ChatMessage, error) {
	message, err := f.queries.GetChatMessage(ctx, id)
	if err != nil {
		return repository.ChatMessage{}, fmt.Errorf("get chat message: %w", err)
	}

	return message, nil
}

func (f *ChatMessageService) ListMessagesByRestaurant(ctx context.Context, restaurantID int32) ([]repository.ChatMessage, error) {
	messages, err := f.queries.ListMessagesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list messages by restaurant: %w", err)
	}

	return messages, nil
}

func (f *ChatMessageService) ListMessagesByEmployee(ctx context.Context, employeeID int32) ([]repository.ChatMessage, error) {
	messages, err := f.queries.ListMessagesByEmployee(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list messages by employee: %w", err)
	}

	return messages, nil
}

func (f *ChatMessageService) InsertChatMessage(
	ctx context.Context,
	restaurantID int32,
	employeeID int32,
	message string,
	sentAt *time.Time,
) (repository.ChatMessage, error) {

	chatMessage, err := f.queries.CreateChatMessage(ctx, repository.CreateChatMessageParams{
		RestaurantID: restaurantID,
		EmployeeID:   employeeID,
		Message:      message,
		SentAt:       helpers.TimestampFromTimePtr(sentAt),
	})
	if err != nil {
		return repository.ChatMessage{}, fmt.Errorf("create chat message: %w", err)
	}

	return chatMessage, nil
}

func (f *ChatMessageService) DeleteChatMessage(ctx context.Context, id int32) error {
	err := f.queries.DeleteChatMessage(ctx, id)
	if err != nil {
		return fmt.Errorf("delete chat message: %w", err)
	}

	return nil
}

func (f *ChatMessageService) ListMessagesByRestaurantWithSender(ctx context.Context, restaurantID int32) ([]repository.ListMessagesByRestaurantWithSenderRow, error) {
	messages, err := f.queries.ListMessagesByRestaurantWithSender(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list messages by restaurant with sender: %w", err)
	}
	return messages, nil
}

func (f *ChatMessageService) ListMessagesByRestaurantAfterID(ctx context.Context, restaurantID int32, afterID int32) ([]repository.ListMessagesByRestaurantAfterIDRow, error) {
	messages, err := f.queries.ListMessagesByRestaurantAfterID(ctx, repository.ListMessagesByRestaurantAfterIDParams{
		RestaurantID: restaurantID,
		ID:           afterID,
	})
	if err != nil {
		return nil, fmt.Errorf("list messages by restaurant after id: %w", err)
	}
	return messages, nil
}
