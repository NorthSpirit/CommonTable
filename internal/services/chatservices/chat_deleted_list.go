package chatservices

import (
	"context"
	"fmt"

	"CommonTable/internal/repository"
)

func (f *ChatMessageService) ListDeletedMessages(ctx context.Context, restaurantID int32) ([]repository.ListDeletedMessagesByRestaurantRow, error) {
	rows, err := f.queries.ListDeletedMessagesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list deleted messages: %w", err)
	}
	return rows, nil
}
