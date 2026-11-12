package externalorderservices

import (
	"context"
	"errors"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Limits for orders coming from outside, counted from unapproved_orders.
const (
	rateWindowSeconds  = 60
	maxPerWindow       = 10
	maxPendingPerKey   = 50
	maxPendingPerPhone = 3
)

var (
	ErrRateLimited     = errors.New("too many orders in a short time")
	ErrTooManyPending  = errors.New("too many orders are waiting for approval")
	ErrPhoneLimit      = errors.New("too many orders are waiting for approval for this phone number")
	ErrItemNotFound    = errors.New("item not found")
	ErrItemUnavailable = errors.New("item is not available")
	ErrAlreadyDecided  = errors.New("order was already decided or does not exist")
)

// ItemError tells which item caused ErrItemNotFound or ErrItemUnavailable.
type ItemError struct {
	ItemID int32
	Err    error
}

func (e *ItemError) Error() string { return fmt.Sprintf("item %d: %v", e.ItemID, e.Err) }
func (e *ItemError) Unwrap() error { return e.Err }

type OrderLineInput struct {
	ItemID   int32
	Quantity int32
	Notes    string
}

type SubmitInput struct {
	IdempotencyKey string
	CustomerName   string
	CustomerPhone  string
	Notes          string
	Lines          []OrderLineInput
}

type SubmitResult struct {
	OrderID   int32
	Duplicate bool
}

type ExternalOrderService struct {
	queries *repository.Queries
	pool    *pgxpool.Pool
}

func NewExternalOrderService(queries *repository.Queries, pool *pgxpool.Pool) *ExternalOrderService {
	return &ExternalOrderService{queries: queries, pool: pool}
}

func textOrNull(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (s *ExternalOrderService) ListAvailableItems(ctx context.Context, restaurantID int32) ([]repository.ListAvailableFoodItemsByRestaurantRow, error) {
	items, err := s.queries.ListAvailableFoodItemsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list available items: %w", err)
	}

	return items, nil
}

func (s *ExternalOrderService) Submit(ctx context.Context, key repository.RestaurantApiKey, in SubmitInput) (SubmitResult, error) {
	existing, err := s.queries.GetUnapprovedOrderByIdempotencyKey(ctx, repository.GetUnapprovedOrderByIdempotencyKeyParams{
		ApiKeyID:       key.ID,
		IdempotencyKey: in.IdempotencyKey,
	})
	if err == nil {
		return SubmitResult{OrderID: existing.ID, Duplicate: true}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return SubmitResult{}, fmt.Errorf("check idempotency key: %w", err)
	}

	if err := s.checkLimits(ctx, key, in.CustomerPhone); err != nil {
		return SubmitResult{}, err
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SubmitResult{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	order, err := q.CreateUnapprovedOrder(ctx, repository.CreateUnapprovedOrderParams{
		RestaurantID:   key.RestaurantID,
		ApiKeyID:       key.ID,
		CustomerName:   in.CustomerName,
		CustomerPhone:  in.CustomerPhone,
		Notes:          textOrNull(in.Notes),
		IdempotencyKey: in.IdempotencyKey,
	})
	if err != nil {
		if isUniqueViolation(err) {
			// The same request arrived twice at once, the other one won.
			if again, getErr := s.queries.GetUnapprovedOrderByIdempotencyKey(ctx, repository.GetUnapprovedOrderByIdempotencyKeyParams{
				ApiKeyID:       key.ID,
				IdempotencyKey: in.IdempotencyKey,
			}); getErr == nil {
				return SubmitResult{OrderID: again.ID, Duplicate: true}, nil
			}
		}
		return SubmitResult{}, fmt.Errorf("create unapproved order: %w", err)
	}

	for _, line := range in.Lines {
		item, err := q.GetFoodItemForExternalOrder(ctx, repository.GetFoodItemForExternalOrderParams{
			ID:           line.ItemID,
			RestaurantID: key.RestaurantID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return SubmitResult{}, &ItemError{ItemID: line.ItemID, Err: ErrItemNotFound}
			}
			return SubmitResult{}, fmt.Errorf("get item %d: %w", line.ItemID, err)
		}
		if !item.IsAvailable {
			return SubmitResult{}, &ItemError{ItemID: line.ItemID, Err: ErrItemUnavailable}
		}

		price := item.Price
		if !price.Valid {
			price = helpers.NumericFromFloat64(0)
		}

		err = q.CreateUnapprovedOrderItem(ctx, repository.CreateUnapprovedOrderItemParams{
			UnapprovedOrderID: order.ID,
			ItemID:            line.ItemID,
			Quantity:          line.Quantity,
			UnitPriceSnapshot: price,
			VatSnapshot:       item.Vat,
			ItemNameSnapshot:  item.Name,
			Notes:             textOrNull(line.Notes),
		})
		if err != nil {
			return SubmitResult{}, fmt.Errorf("add item %d: %w", line.ItemID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return SubmitResult{}, fmt.Errorf("commit: %w", err)
	}

	return SubmitResult{OrderID: order.ID}, nil
}

func (s *ExternalOrderService) checkLimits(ctx context.Context, key repository.RestaurantApiKey, phone string) error {
	recent, err := s.queries.CountRecentUnapprovedOrdersByApiKey(ctx, repository.CountRecentUnapprovedOrdersByApiKeyParams{
		ApiKeyID:      key.ID,
		WindowSeconds: rateWindowSeconds,
	})
	if err != nil {
		return fmt.Errorf("count recent orders: %w", err)
	}
	if recent >= maxPerWindow {
		return ErrRateLimited
	}

	pending, err := s.queries.CountPendingUnapprovedOrdersByApiKey(ctx, key.ID)
	if err != nil {
		return fmt.Errorf("count pending orders: %w", err)
	}
	if pending >= maxPendingPerKey {
		return ErrTooManyPending
	}

	pendingForPhone, err := s.queries.CountPendingUnapprovedOrdersByPhone(ctx, repository.CountPendingUnapprovedOrdersByPhoneParams{
		RestaurantID:  key.RestaurantID,
		CustomerPhone: phone,
	})
	if err != nil {
		return fmt.Errorf("count pending orders for phone: %w", err)
	}
	if pendingForPhone >= maxPendingPerPhone {
		return ErrPhoneLimit
	}

	return nil
}

func (s *ExternalOrderService) ListPending(ctx context.Context, restaurantID int32) ([]repository.UnapprovedOrder, error) {
	orders, err := s.queries.ListPendingUnapprovedOrders(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list pending orders: %w", err)
	}

	return orders, nil
}

type PendingOrder struct {
	Order repository.UnapprovedOrder
	Items []repository.UnapprovedOrderItem
	Total float64
}

func (s *ExternalOrderService) ListPendingWithItems(ctx context.Context, restaurantID int32) ([]PendingOrder, error) {
	orders, err := s.ListPending(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	out := make([]PendingOrder, 0, len(orders))
	for _, o := range orders {
		items, err := s.ListItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}

		total := 0.0
		for _, item := range items {
			if p, err := helpers.NumericToFloat64(item.UnitPriceSnapshot); err == nil && p != nil {
				total += *p * float64(item.Quantity)
			}
		}

		out = append(out, PendingOrder{Order: o, Items: items, Total: total})
	}

	return out, nil
}

func (s *ExternalOrderService) ListItems(ctx context.Context, unapprovedOrderID int32) ([]repository.UnapprovedOrderItem, error) {
	items, err := s.queries.ListUnapprovedOrderItems(ctx, unapprovedOrderID)
	if err != nil {
		return nil, fmt.Errorf("list unapproved order items: %w", err)
	}

	return items, nil
}

func (s *ExternalOrderService) Approve(ctx context.Context, restaurantID int32, unapprovedOrderID int32, managerID int32) (repository.FoodOrder, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	pending, err := q.GetUnapprovedOrderForUpdate(ctx, repository.GetUnapprovedOrderForUpdateParams{
		ID:           unapprovedOrderID,
		RestaurantID: restaurantID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.FoodOrder{}, ErrAlreadyDecided
		}
		return repository.FoodOrder{}, fmt.Errorf("get unapproved order: %w", err)
	}
	if pending.Status != "pending" {
		return repository.FoodOrder{}, ErrAlreadyDecided
	}

	items, err := q.ListUnapprovedOrderItems(ctx, pending.ID)
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("list items: %w", err)
	}

	orderNumber, err := q.GetNextOrderNumberForToday(ctx, restaurantID)
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("get order number: %w", err)
	}

	notes := "Online order: " + pending.CustomerName + ", " + pending.CustomerPhone
	if pending.Notes.Valid && pending.Notes.String != "" {
		notes += "\n" + pending.Notes.String
	}

	order, err := q.CreateFoodOrder(ctx, repository.CreateFoodOrderParams{
		RestaurantID:      restaurantID,
		TableID:           pgtype.Int4{},
		WaiterID:          managerID,
		OrderCreationDate: pgtype.Timestamp{Time: time.Now(), Valid: true},
		OrderStatus:       repository.OrderStatusPending,
		OrderNumber:       orderNumber,
		Notes:             textOrNull(notes),
	})
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("create order: %w", err)
	}

	for _, item := range items {
		name := item.ItemNameSnapshot
		_, err := q.AddItemToOrder(ctx, repository.AddItemToOrderParams{
			OrderID:           order.ID,
			ItemID:            item.ItemID,
			Quantity:          item.Quantity,
			UnitPriceSnapshot: item.UnitPriceSnapshot,
			VatSnapshot:       item.VatSnapshot,
			Notes:             item.Notes,
			ItemNameSnapshot:  helpers.TextFromString(&name),
		})
		if err != nil {
			return repository.FoodOrder{}, fmt.Errorf("add item %d: %w", item.ItemID, err)
		}
	}

	rows, err := q.MarkUnapprovedOrderApproved(ctx, repository.MarkUnapprovedOrderApprovedParams{
		ID:              pending.ID,
		ApprovedOrderID: pgtype.Int4{Int32: order.ID, Valid: true},
		DecidedBy:       pgtype.Int4{Int32: managerID, Valid: true},
	})
	if err != nil {
		return repository.FoodOrder{}, fmt.Errorf("mark approved: %w", err)
	}
	if rows == 0 {
		return repository.FoodOrder{}, ErrAlreadyDecided
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.FoodOrder{}, fmt.Errorf("commit: %w", err)
	}

	return order, nil
}

func (s *ExternalOrderService) Reject(ctx context.Context, restaurantID int32, unapprovedOrderID int32, managerID int32, reason string) error {
	// Restaurant check first, so a manager cannot reject another restaurant's order.
	order, err := s.queries.GetUnapprovedOrderForUpdate(ctx, repository.GetUnapprovedOrderForUpdateParams{
		ID:           unapprovedOrderID,
		RestaurantID: restaurantID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAlreadyDecided
		}
		return fmt.Errorf("get unapproved order: %w", err)
	}

	rows, err := s.queries.MarkUnapprovedOrderRejected(ctx, repository.MarkUnapprovedOrderRejectedParams{
		ID:              order.ID,
		RejectionReason: textOrNull(reason),
		DecidedBy:       pgtype.Int4{Int32: managerID, Valid: managerID != 0},
	})
	if err != nil {
		return fmt.Errorf("mark rejected: %w", err)
	}
	if rows == 0 {
		return ErrAlreadyDecided
	}

	return nil
}
