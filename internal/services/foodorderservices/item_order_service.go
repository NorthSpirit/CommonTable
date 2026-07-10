package foodorderservices

import (
	"context"
	"fmt"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type ItemOrderService struct {
	queries *repository.Queries
}

func NewItemOrderService(queries *repository.Queries) *ItemOrderService {
	return &ItemOrderService{queries: queries}
}

func (s *ItemOrderService) AddItemToOrder(
	ctx context.Context,
	orderID int32,
	itemID int32,
	quantity int32,
	unitPriceSnapshot float64,
	vatSnapshot *float64,
	notes *string,
) (repository.ItemOrder, error) {

	item, err := s.queries.GetFoodItem(ctx, itemID)
	if err != nil {
		return repository.ItemOrder{}, fmt.Errorf("look up food item for snapshot: %w", err)
	}

	itemOrder, err := s.queries.AddItemToOrder(ctx, repository.AddItemToOrderParams{
		OrderID:                 orderID,
		ItemID:                  itemID,
		Quantity:                quantity,
		UnitPriceSnapshot:       helpers.NumericFromFloat64(unitPriceSnapshot),
		VatSnapshot:             helpers.NumericFromFloat64Ptr(vatSnapshot),
		Notes:                   helpers.TextFromString(notes),
		ItemNameSnapshot:        helpers.TextFromString(&item.Name),
		ItemDescriptionSnapshot: item.Description,
	})
	if err != nil {
		return repository.ItemOrder{}, fmt.Errorf("add item to order: %w", err)
	}

	return itemOrder, nil
}

func (s *ItemOrderService) ListItemsForOrder(ctx context.Context, orderID int32) ([]repository.ItemOrder, error) {
	items, err := s.queries.ListItemsForOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("list items for order: %w", err)
	}

	return items, nil
}

func (s *ItemOrderService) UpdateItemOrderQuantity(ctx context.Context, id int32, quantity int32) error {
	err := s.queries.UpdateItemOrderQuantity(ctx, repository.UpdateItemOrderQuantityParams{
		ID:       id,
		Quantity: quantity,
	})
	if err != nil {
		return fmt.Errorf("update item order quantity: %w", err)
	}

	return nil
}

func (s *ItemOrderService) RemoveItemFromOrder(ctx context.Context, id int32) error {
	err := s.queries.RemoveItemFromOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("remove item from order: %w", err)
	}

	return nil
}

func (s *ItemOrderService) ClearOrderItems(ctx context.Context, orderID int32) error {
	err := s.queries.ClearOrderItems(ctx, orderID)
	if err != nil {
		return fmt.Errorf("clear order items: %w", err)
	}

	return nil
}
