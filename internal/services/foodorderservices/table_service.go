package foodorderservices

import (
	"CommonTable/internal/helpers"
	"context"
	"fmt"

	"CommonTable/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type TableService struct {
	queries *repository.Queries
}

func NewTableService(queries *repository.Queries) *TableService {
	return &TableService{queries: queries}
}

func (f *TableService) GetTableByID(ctx context.Context, id int32) (repository.Table, error) {
	table, err := f.queries.GetTable(ctx, id)
	if err != nil {
		return repository.Table{}, fmt.Errorf("get table: %w", err)
	}

	return table, nil
}

func (f *TableService) ListTables(ctx context.Context) ([]repository.Table, error) {
	tables, err := f.queries.ListTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tables: %w", err)
	}

	return tables, nil
}

func (f *TableService) ListTablesByStatus(ctx context.Context, status repository.TableStatus) ([]repository.Table, error) {
	tables, err := f.queries.ListTablesByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("list tables by status: %w", err)
	}

	return tables, nil
}

func (f *TableService) InsertTable(
	ctx context.Context,
	number int32,
	capacity int32,
	section *string,
	status repository.TableStatus,
) (repository.Table, error) {

	table, err := f.queries.CreateTable(ctx, repository.CreateTableParams{
		Number:   number,
		Capacity: capacity,
		Section:  helpers.TextFromString(section),
		Status:   status,
	})
	if err != nil {
		return repository.Table{}, fmt.Errorf("create table: %w", err)
	}

	return table, nil
}

func (f *TableService) UpdateTable(
	ctx context.Context,
	id int32,
	number int32,
	capacity int32,
	section *string,
) error {

	err := f.queries.UpdateTable(ctx, repository.UpdateTableParams{
		ID:       id,
		Number:   number,
		Capacity: capacity,
		Section:  helpers.TextFromString(section),
	})
	if err != nil {
		return fmt.Errorf("update table: %w", err)
	}

	return nil
}

func (f *TableService) SetTableStatus(ctx context.Context, id int32, status repository.TableStatus) error {
	err := f.queries.SetTableStatus(ctx, repository.SetTableStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("set table status: %w", err)
	}

	return nil
}

func (f *TableService) DeleteTable(ctx context.Context, id int32) error {
	err := f.queries.SoftDeleteTable(ctx, id)
	if err != nil {
		return fmt.Errorf("delete table: %w", err)
	}
	return nil
}

func (f *TableService) ListTablesForAdmin(ctx context.Context) ([]repository.Table, error) {
	tables, err := f.queries.ListTablesForAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tables for admin: %w", err)
	}
	return tables, nil
}

func (f *TableService) RestoreTable(ctx context.Context, id int32) error {
	err := f.queries.RestoreTable(ctx, id)
	if err != nil {
		return fmt.Errorf("restore table: %w", err)
	}
	return nil
}

func (s *TableService) ListTablesWithActiveOrderFlag(ctx context.Context, restaurantID int32) ([]repository.ListTablesWithActiveOrderFlagRow, error) {
	tables, err := s.queries.ListTablesWithActiveOrderFlag(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list tables with active order flag: %w", err)
	}
	return tables, nil
}

func (s *TableService) ListTablesWithOrderStatusForRestaurant(ctx context.Context, restaurantID int32) ([]repository.ListTablesWithOrderStatusForRestaurantRow, error) {
	tables, err := s.queries.ListTablesWithOrderStatusForRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list tables with order status for restaurant: %w", err)
	}
	return tables, nil
}

func (s *TableService) ListTablesWithGlobalOrderStatus(ctx context.Context) ([]repository.ListTablesWithGlobalOrderStatusRow, error) {
	tables, err := s.queries.ListTablesWithGlobalOrderStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("list tables with global order status: %w", err)
	}
	return tables, nil
}

func (f *TableService) TableHasActiveOrders(ctx context.Context, tableID int32) (bool, error) {
	has, err := f.queries.TableHasActiveOrders(ctx, pgtype.Int4{Int32: tableID, Valid: true})
	if err != nil {
		return false, fmt.Errorf("check table active orders: %w", err)
	}
	return has, nil
}
