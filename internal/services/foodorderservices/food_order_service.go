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
	orders, err := s.queries.ListOrdersByTable(ctx, pgtype.Int4{Int32: tableID, Valid: true})
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
		TableID:           pgtype.Int4{Int32: tableID, Valid: true},
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

func (s *FoodOrderService) ListActiveOrdersByWaiter(ctx context.Context, waiterID int32) ([]repository.ListActiveOrdersByWaiterRow, error) {
	orders, err := s.queries.ListActiveOrdersByWaiter(ctx, waiterID)
	if err != nil {
		return nil, fmt.Errorf("list active orders by waiter: %w", err)
	}
	return orders, nil
}

func (s *FoodOrderService) GetOrderByIdempotencyKey(ctx context.Context, key string) (repository.FoodOrder, error) {
	order, err := s.queries.GetOrderByIdempotencyKey(ctx, helpers.TextFromString(&key))
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("get order by idempotency key: %w", err)
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

// GetActiveOrderForTable returns the currently open order for a table at
// a specific restaurant, if any (order_status not in delivered/cancelled).
// Scoped by both table and restaurant, since the food-hall model allows
// different restaurants to have independent open orders at the same
// physical table simultaneously.
//
// Returns pgx.ErrNoRows (wrapped) if there's no open order — callers
// should check with errors.Is(err, pgx.ErrNoRows) to distinguish "no
// active order" (normal) from a real failure. Never called for takeout
// orders (no table), since takeout orders are always independent and
// must never be merged together.
func (s *FoodOrderService) GetActiveOrderForTable(ctx context.Context, tableID int32, restaurantID int32) (repository.FoodOrder, error) {
	order, err := s.queries.GetActiveOrderForTable(ctx, repository.GetActiveOrderForTableParams{
		TableID:      pgtype.Int4{Int32: tableID, Valid: true},
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
//
// tableID is nil for takeout orders (no table assigned).
func (s *FoodOrderService) CreateOrderWithItemsAtomic(
	// Passes tons of stuff required for order's creation (plus context and pool)
	ctx context.Context,
	pool *pgxpool.Pool,
	restaurantID int32,
	tableID *int32,
	waiterID int32,
	orderNumber int32,
	notes *string,
	idempotencyKey string,
	lines []OrderLine,
) (repository.FoodOrder, error) {

	// Borrows a connection from the PostgreSQL pool and starts a transaction on it, tx is that transaction
	// while err is obviously the error
	tx, err := pool.Begin(ctx)

	// If the transaction couldn't be started, returns the error, marked as "begin transaction"
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("begin transaction: %w", err)
	}

	// If something breaks up anywhere before commit, it undoes everything prior to this point
	// And returns the connection back the pool
	// If commit succeeds, rollback still fires, but basically does nothing.
	defer tx.Rollback(ctx)

	// SQLC generated go code executes the actual query, initiates the struct
	// Passing tx makes every query below run inside this transaction
	// Basically a C# equivalent would be init the class to use the method.
	q := repository.New(tx)

	// First, checks the table ID, if it's null (takeout), the order is saved without a table (NULL in the DB)
	var tableIDParam pgtype.Int4
	if tableID != nil {
		tableIDParam = pgtype.Int4{Int32: *tableID, Valid: true}
	}

	// Creates the food order, with the passed down values
	order, err := q.CreateFoodOrderWithIdempotencyKey(ctx, repository.CreateFoodOrderWithIdempotencyKeyParams{
		RestaurantID:      restaurantID,
		TableID:           tableIDParam,
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

	// For every food line in the order passed to this method
	for _, line := range lines {

		// Looks up the item by ID to check its price, name and description
		item, err := q.GetFoodItem(ctx, line.ItemID)
		if err != nil {
			return repository.FoodOrder{}, fmt.Errorf("get item %d: %w", line.ItemID, err) // Item wasn't gotten - the lookup failed (item missing or DB error)
		}

		// Price is first, because it shouldn't be null, if it's null, it gets turned to 0
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
			ItemNameSnapshot:        helpers.TextFromString(&item.Name), // Copies the item's current name into the snapshot
			ItemDescriptionSnapshot: item.Description,                   // And the description (also a copy)
		})
		// If adding failed
		if err != nil {
			return repository.FoodOrder{}, fmt.Errorf("add item %d to order: %w", line.ItemID, err)
		}
	}

	// Commit saves the order (with its key) and all lines together, all or nothing, making it atomic
	if err := tx.Commit(ctx); err != nil {
		return repository.FoodOrder{}, fmt.Errorf("commit transaction: %w", err)
	}

	// Returns the order and a null for an error
	return order, nil
}

// AppendItemsToOrderAtomic adds item lines to an already-open order as one
// atomic transaction. Idempotency is tracked via the dedicated
// idempotency_keys table (not food_orders), since an append isn't itself
// an order row — this avoids polluting food_orders with marker rows.
//
// Only ever called for table orders — takeout orders have no shared
// identity to append to, so each takeout submission always creates a new
// order via CreateOrderWithItemsAtomic instead.
func (s *FoodOrderService) AppendItemsToOrderAtomic(
	ctx context.Context,
	pool *pgxpool.Pool,
	orderID int32,
	idempotencyKey string,
	lines []OrderLine,
) error {

	// Borrows the connection to the DB and starts the transaction
	tx, err := pool.Begin(ctx)

	// If there was some kind of error, send the error what happened
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// If something breaks up anywhere before commit, it undoes everything prior to this point
	// And returns the connection back the pool
	// If commit succeeds, rollback still fires, but basically does nothing.
	defer tx.Rollback(ctx)

	// SQLC generated go code executes the actual query, initiates the struct
	// Basically a C# equivalent would be init the class to use the method.
	q := repository.New(tx)

	// Checks if Idempotency key was already in use
	// _ means, I don't care about the row itself, only the "error" - if it was actually used
	if _, err := q.CheckIdempotencyKey(ctx, idempotencyKey); err == nil { // err == nil means a row was found, so the key was already used
		return nil // already applied, safe retry
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check idempotency: %w", err) // Query itself failed
	}

	// For every food line in the order passed to this method
	for _, line := range lines {

		// Looks up the item by ID to check its price, name and description
		item, err := q.GetFoodItem(ctx, line.ItemID)
		if err != nil {
			return fmt.Errorf("get item %d: %w", line.ItemID, err)
		} // Item wasn't gotten - the lookup failed (item missing or DB error)

		// Price is first, because it shouldn't be null, if it's null, it gets turned to 0
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
			ItemNameSnapshot:        helpers.TextFromString(&item.Name), // Gets the item's name
			ItemDescriptionSnapshot: item.Description,                   // And the description
		})
		// If adding failed
		if err != nil {
			return fmt.Errorf("add item %d to order: %w", line.ItemID, err)
		}
	}

	// All lines added, now the key is saved as used, so a retry with the same key will be skipped above
	if err := q.RecordIdempotencyKey(ctx, idempotencyKey); err != nil {
		return fmt.Errorf("record idempotency key: %w", err)
	}

	// Commit saves all lines and the key together, all or nothing, making it atomic
	return tx.Commit(ctx)
}

func (s *FoodOrderService) ListTablesWithGlobalOrderStatus(ctx context.Context) ([]repository.ListTablesWithGlobalOrderStatusRow, error) {
	tables, err := s.queries.ListTablesWithGlobalOrderStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tables with global order status: %w", err)
	}
	return tables, nil
}
