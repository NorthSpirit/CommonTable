package chatservices

import (
	"context"
	"errors"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrMessageNotFound = errors.New("chat message not found or already deleted")

func DeletedText(reason string) string {
	if reason == "" {
		return "[message deleted]"
	}
	return "[message deleted: " + reason + "]"
}

func DisplayMessage(message string, deletedAt pgtype.Timestamp, deleteReason pgtype.Text) string {
	if !deletedAt.Valid {
		return message
	}
	return DeletedText(deleteReason.String)
}

// An empty reason is stored as no reason.
func (f *ChatMessageService) SoftDeleteMessage(ctx context.Context, id int32, restaurantID int32, deletedBy string, reason string) error {
	var reasonPtr *string
	if reason != "" {
		reasonPtr = &reason
	}

	_, err := f.queries.SoftDeleteChatMessage(ctx, repository.SoftDeleteChatMessageParams{
		ID:           id,
		RestaurantID: restaurantID,
		DeletedBy:    helpers.TextFromString(&deletedBy),
		DeleteReason: helpers.TextFromString(reasonPtr),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return fmt.Errorf("soft delete chat message: %w", err)
	}

	return nil
}
