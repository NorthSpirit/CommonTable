package restaurantservices

import (
	"context"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type EmployeeService struct {
	queries *repository.Queries
}

func NewEmployeeService(queries *repository.Queries) *EmployeeService {
	return &EmployeeService{queries: queries}
}

func (f *EmployeeService) GetEmployeeByID(ctx context.Context, id int32) (repository.Employee, error) {
	employee, err := f.queries.GetEmployee(ctx, id)
	if err != nil {
		return repository.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	return employee, nil
}

func (f *EmployeeService) ListEmployeesByRestaurant(ctx context.Context, restaurantID int32) ([]repository.Employee, error) {
	employees, err := f.queries.ListEmployeesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list employees by restaurant: %w", err)
	}

	return employees, nil
}

func (f *EmployeeService) ListActiveEmployeesByRestaurant(ctx context.Context, restaurantID int32) ([]repository.Employee, error) {
	employees, err := f.queries.ListActiveEmployeesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("list active employees by restaurant: %w", err)
	}

	return employees, nil
}

func (f *EmployeeService) InsertEmployee(
	ctx context.Context,
	restaurantID int32,
	name *string,
	role repository.EmployeeRole,
	pinHash *string,
	hireDate *time.Time,
	email *string,
	phoneNumber *string,
) (repository.Employee, error) {

	employee, err := f.queries.CreateEmployee(ctx, repository.CreateEmployeeParams{
		RestaurantID: restaurantID,
		Name:         helpers.TextFromString(name),
		Role:         role,
		PinHash:      helpers.TextFromString(pinHash),
		HireDate:     helpers.TimestampFromTimePtr(hireDate),
		Email:        helpers.TextFromString(email),
		PhoneNumber:  helpers.TextFromString(phoneNumber),
	})
	if err != nil {
		return repository.Employee{}, fmt.Errorf("create employee: %w", err)
	}

	return employee, nil
}

func (f *EmployeeService) UpdateEmployee(
	ctx context.Context,
	id int32,
	name *string,
	role repository.EmployeeRole,
	email *string,
	phoneNumber *string,
) error {

	err := f.queries.UpdateEmployee(ctx, repository.UpdateEmployeeParams{
		ID:          id,
		Name:        helpers.TextFromString(name),
		Role:        role,
		Email:       helpers.TextFromString(email),
		PhoneNumber: helpers.TextFromString(phoneNumber),
	})
	if err != nil {
		return fmt.Errorf("update employee: %w", err)
	}

	return nil
}

func (f *EmployeeService) UpdateEmployeePinHash(ctx context.Context, id int32, pinHash string) error {
	err := f.queries.UpdateEmployeePinHash(ctx, repository.UpdateEmployeePinHashParams{
		ID:      id,
		PinHash: helpers.TextFromString(&pinHash),
	})
	if err != nil {
		return fmt.Errorf("update employee pin hash: %w", err)
	}

	return nil
}

func (f *EmployeeService) SetEmployeeActive(ctx context.Context, id int32, isActive bool) error {
	err := f.queries.SetEmployeeActive(ctx, repository.SetEmployeeActiveParams{
		ID:       id,
		IsActive: isActive,
	})
	if err != nil {
		return fmt.Errorf("set employee active: %w", err)
	}

	return nil
}

func (f *EmployeeService) UpdateEmployeeLastLogin(ctx context.Context, id int32, lastLogin time.Time) error {
	err := f.queries.UpdateEmployeeLastLogin(ctx, repository.UpdateEmployeeLastLoginParams{
		ID:        id,
		LastLogin: helpers.TimestampFromTimePtr(&lastLogin),
	})
	if err != nil {
		return fmt.Errorf("update employee last login: %w", err)
	}

	return nil
}

func (f *EmployeeService) DeleteEmployee(ctx context.Context, id int32) error {
	err := f.queries.DeleteEmployee(ctx, id)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}

	return nil
}
