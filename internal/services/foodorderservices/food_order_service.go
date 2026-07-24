package foodorderservices

import (
	"context"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type FoodOrderService struct {
	queries *repository.Queries
}

func NewFoodOrderService(queries *repository.Queries) *FoodOrderService {
	return &FoodOrderService{queries: queries}
}

func (s *FoodOrderService) GetFoodOrder(ctx context.Context, id int32) (repository.FoodOrder, error) {
	order, err := s.queries.GetFoodOrder(ctx, id)
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("get food order: %w", err)
	}

	return order, nil
}

func (s *FoodOrderService) ListOrdersByRestaurant(ctx context.Context, restaurantID int32) ([]repository.FoodOrder, error) {
	orders, err := s.queries.ListOrdersByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list orders by restaurant: %w", err)
	}

	return orders, nil
}

func (s *FoodOrderService) ListOrdersByTable(ctx context.Context, tableID int32) ([]repository.FoodOrder, error) {
	orders, err := s.queries.ListOrdersByTable(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("list orders by table: %w", err)
	}

	return orders, nil
}

func (s *FoodOrderService) ListOrdersByStatus(
	ctx context.Context,
	restaurantID int32,
	status repository.OrderStatus,
) ([]repository.FoodOrder, error) {

	orders, err := s.queries.ListOrdersByStatus(ctx, repository.ListOrdersByStatusParams{
		RestaurantID: restaurantID,
		OrderStatus:  status,
	})
	if err != nil {
		return nil, fmt.Errorf("list orders by status: %w", err)
	}

	return orders, nil
}

func (s *FoodOrderService) CreateFoodOrder(
	ctx context.Context,
	restaurantID int32,
	tableID int32,
	waiterID int32,
	orderCreationDate time.Time,
	orderStatus repository.OrderStatus,
	orderNumber int32,
	notes *string,
) (repository.FoodOrder, error) {

	order, err := s.queries.CreateFoodOrder(ctx, repository.CreateFoodOrderParams{
		RestaurantID:      restaurantID,
		TableID:           tableID,
		WaiterID:          waiterID,
		OrderCreationDate: pgtype.Timestamp{Time: orderCreationDate, Valid: true},
		OrderStatus:       orderStatus,
		OrderNumber:       orderNumber,
		Notes:             helpers.TextFromString(notes),
	})
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("create food order: %w", err)
	}

	return order, nil
}

func (s *FoodOrderService) UpdateOrderStatus(ctx context.Context, id int32, status repository.OrderStatus) error {
	err := s.queries.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
		ID:          id,
		OrderStatus: status,
	})
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}

func (s *FoodOrderService) UpdateFoodOrderNotes(ctx context.Context, id int32, notes *string) error {
	err := s.queries.UpdateFoodOrderNotes(ctx, repository.UpdateFoodOrderNotesParams{
		ID:    id,
		Notes: helpers.TextFromString(notes),
	})
	if err != nil {
		return fmt.Errorf("update food order notes: %w", err)
	}

	return nil
}

func (s *FoodOrderService) DeleteFoodOrder(ctx context.Context, id int32) error {
	err := s.queries.DeleteFoodOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("delete food order: %w", err)
	}

	return nil
}

func (s *FoodOrderService) ListActiveOrdersByRestaurantForKitchen(ctx context.Context, restaurantID int32) ([]repository.ListActiveOrdersByRestaurantForKitchenRow, error) {
	orders, err := s.queries.ListActiveOrdersByRestaurantForKitchen(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list active orders by restaurant for kitchen: %w", err)
	}
	return orders, nil
}

func (s *FoodOrderService) ListOrdersByRestaurantForAdmin(ctx context.Context, restaurantID int32) ([]repository.ListOrdersByRestaurantForAdminRow, error) {
	orders, err := s.queries.ListOrdersByRestaurantForAdmin(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list orders by restaurant for admin: %w", err)
	}
	return orders, nil
}

func (s *FoodOrderService) GetOrderForAdmin(ctx context.Context, id int32) (repository.GetOrderForAdminRow, error) {
	order, err := s.queries.GetOrderForAdmin(ctx, id)
	if err != nil {
		return repository.GetOrderForAdminRow{}, fmt.Errorf("get order for admin: %w", err)
	}
	return order, nil
}
