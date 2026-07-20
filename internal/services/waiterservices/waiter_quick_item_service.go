package waiterservices

import (
	"context"
	"errors"
	"fmt"

	"CommonTable/internal/repository"
)

const maxQuickItems = 15

var ErrQuickItemLimitReached = errors.New("waiter quick item limit reached (max 15)")

type WaiterQuickItemService struct {
	queries *repository.Queries
}

func NewWaiterQuickItemService(queries *repository.Queries) *WaiterQuickItemService {
	return &WaiterQuickItemService{queries: queries}
}

func (s *WaiterQuickItemService) ListQuickItems(ctx context.Context, employeeID int32) ([]repository.ListWaiterQuickItemsRow, error) {
	items, err := s.queries.ListWaiterQuickItems(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("list waiter quick items: %w", err)
	}

	return items, nil
}

// AddQuickItem enforces the 15-item cap before inserting — the database
// itself has no such constraint, so this check must happen here.
func (s *WaiterQuickItemService) AddQuickItem(ctx context.Context, employeeID int32, itemID int32, sortOrder int32) error {
	count, err := s.queries.CountWaiterQuickItems(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("count waiter quick items: %w", err)
	}

	if count >= maxQuickItems {
		return ErrQuickItemLimitReached
	}

	err = s.queries.AddWaiterQuickItem(ctx, repository.AddWaiterQuickItemParams{
		EmployeeID: employeeID,
		ItemID:     itemID,
		SortOrder:  sortOrder,
	})
	if err != nil {
		return fmt.Errorf("add waiter quick item: %w", err)
	}

	return nil
}

func (s *WaiterQuickItemService) RemoveQuickItem(ctx context.Context, employeeID int32, itemID int32) error {
	err := s.queries.RemoveWaiterQuickItem(ctx, repository.RemoveWaiterQuickItemParams{
		EmployeeID: employeeID,
		ItemID:     itemID,
	})
	if err != nil {
		return fmt.Errorf("remove waiter quick item: %w", err)
	}

	return nil
}

func (s *WaiterQuickItemService) ClearQuickItems(ctx context.Context, employeeID int32) error {
	err := s.queries.ClearWaiterQuickItems(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("clear waiter quick items: %w", err)
	}

	return nil
}

func (s *WaiterQuickItemService) UpdateSortOrder(ctx context.Context, employeeID int32, itemID int32, sortOrder int32) error {
	err := s.queries.UpdateWaiterQuickItemSortOrder(ctx, repository.UpdateWaiterQuickItemSortOrderParams{
		EmployeeID: employeeID,
		ItemID:     itemID,
		SortOrder:  sortOrder,
	})
	if err != nil {
		return fmt.Errorf("update waiter quick item sort order: %w", err)
	}

	return nil
}
