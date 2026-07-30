package foodorderservices

import (
	"context"
	"errors"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (s *FoodOrderService) GetOrderByIdempotencyKey(ctx context.Context, key string) (repository.FoodOrder, error) {
	order, err := s.queries.GetOrderByIdempotencyKey(ctx, helpers.TextFromString(&key))
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("get order by idempotency key: %w", err)
	}
	return order, nil
}

func (s *FoodOrderService) CreateFoodOrderWithIdempotencyKey(
	ctx context.Context,
	restaurantID int32,
	tableID int32,
	waiterID int32,
	orderCreationDate time.Time,
	orderStatus repository.OrderStatus,
	orderNumber int32,
	notes *string,
	idempotencyKey string,
) (repository.FoodOrder, error) {

	order, err := s.queries.CreateFoodOrderWithIdempotencyKey(ctx, repository.CreateFoodOrderWithIdempotencyKeyParams{
		RestaurantID:      restaurantID,
		TableID:           tableID,
		WaiterID:          waiterID,
		OrderCreationDate: pgtype.Timestamp{Time: orderCreationDate, Valid: true},
		OrderStatus:       orderStatus,
		OrderNumber:       orderNumber,
		Notes:             helpers.TextFromString(notes),
		IdempotencyKey:    helpers.TextFromString(&idempotencyKey),
	})
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("create food order with idempotency key: %w", err)
	}
	return order, nil
}

func (s *FoodOrderService) GetNextOrderNumberForToday(ctx context.Context, restaurantID int32) (int32, error) {
	next, err := s.queries.GetNextOrderNumberForToday(ctx, restaurantID)
	if err != nil {
		return 0, fmt.Errorf("get next order number for today: %w", err)
	}
	return next, nil
}

// GetActiveOrderForTable returns the currently open order for a table, if
// any (order_status not in delivered/cancelled). Returns pgx.ErrNoRows
// (wrapped) if the table has no open order — callers should check with
// errors.Is(err, pgx.ErrNoRows) to distinguish "no active order" (normal)
// from a real failure.
func (s *FoodOrderService) GetActiveOrderForTable(ctx context.Context, tableID int32, restaurantID int32) (repository.FoodOrder, error) {
	order, err := s.queries.GetActiveOrderForTable(ctx, repository.GetActiveOrderForTableParams{
		TableID:      tableID,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("get active order for table: %w", err)
	}
	return order, nil
}

type OrderLine struct {
	ItemID   int32
	Quantity int32
	Notes    *string
}

// CreateOrderWithItemsAtomic inserts the order and every item line as one
// atomic transaction — if any item fails to insert, the whole order rolls
// back rather than reaching the kitchen incomplete. Used by both the
// waiter API and admin's manual order-creation screen.
func (s *FoodOrderService) CreateOrderWithItemsAtomic(
	ctx context.Context,
	pool *pgxpool.Pool,
	restaurantID int32,
	tableID int32,
	waiterID int32,
	orderNumber int32,
	notes *string,
	idempotencyKey string,
	lines []OrderLine,
) (repository.FoodOrder, error) {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	order, err := q.CreateFoodOrderWithIdempotencyKey(ctx, repository.CreateFoodOrderWithIdempotencyKeyParams{
		RestaurantID:      restaurantID,
		TableID:           tableID,
		WaiterID:          waiterID,
		OrderCreationDate: pgtype.Timestamp{Time: time.Now(), Valid: true},
		OrderStatus:       repository.OrderStatusPending,
		OrderNumber:       orderNumber,
		Notes:             helpers.TextFromString(notes),
		IdempotencyKey:    helpers.TextFromString(&idempotencyKey),
	})
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("create order: %w", err)
	}

	for _, line := range lines {
		item, err := q.GetFoodItem(ctx, line.ItemID)
		if err != nil {
			return repository.FoodOrder{}, fmt.Errorf("get item %d: %w", line.ItemID, err)
		}

		price := item.Price
		if !price.Valid {
			price = helpers.NumericFromFloat64(0)
		}

		_, err = q.AddItemToOrder(ctx, repository.AddItemToOrderParams{
			OrderID:                 order.ID,
			ItemID:                  line.ItemID,
			Quantity:                line.Quantity,
			UnitPriceSnapshot:       price,
			VatSnapshot:             pgtype.Numeric{Valid: false},
			Notes:                   helpers.TextFromString(line.Notes),
			ItemNameSnapshot:        helpers.TextFromString(&item.Name),
			ItemDescriptionSnapshot: item.Description,
		})
		if err != nil {
			return repository.FoodOrder{}, fmt.Errorf("add item %d to order: %w", line.ItemID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.FoodOrder{}, fmt.Errorf("commit transaction: %w", err)
	}

	return order, nil
}

// AppendItemsToOrderAtomic adds item lines to an already-open order as one
// atomic transaction. Idempotency is tracked via the dedicated
// idempotency_keys table (not food_orders), since an append isn't itself
// an order row — this avoids polluting food_orders with marker rows.
func (s *FoodOrderService) AppendItemsToOrderAtomic(
	ctx context.Context,
	pool *pgxpool.Pool,
	orderID int32,
	idempotencyKey string,
	lines []OrderLine,
) error {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	if _, err := q.CheckIdempotencyKey(ctx, idempotencyKey); err == nil {
		return nil // already applied, safe retry
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check idempotency: %w", err)
	}

	for _, line := range lines {
		item, err := q.GetFoodItem(ctx, line.ItemID)
		if err != nil {
			return fmt.Errorf("get item %d: %w", line.ItemID, err)
		}

		price := item.Price
		if !price.Valid {
			price = helpers.NumericFromFloat64(0)
		}

		_, err = q.AddItemToOrder(ctx, repository.AddItemToOrderParams{
			OrderID:                 orderID,
			ItemID:                  line.ItemID,
			Quantity:                line.Quantity,
			UnitPriceSnapshot:       price,
			VatSnapshot:             pgtype.Numeric{Valid: false},
			Notes:                   helpers.TextFromString(line.Notes),
			ItemNameSnapshot:        helpers.TextFromString(&item.Name),
			ItemDescriptionSnapshot: item.Description,
		})
		if err != nil {
			return fmt.Errorf("add item %d to order: %w", line.ItemID, err)
		}
	}

	if err := q.RecordIdempotencyKey(ctx, idempotencyKey); err != nil {
		return fmt.Errorf("record idempotency key: %w", err)
	}

	return tx.Commit(ctx)
}
