package foodorderservices

import (
	"context"
	"fmt"
	"time"

	"CommonTable/internal/helpers"
	"CommonTable/internal/repository"
)

type PaymentService struct {
	queries *repository.Queries
}

func NewPaymentService(queries *repository.Queries) *PaymentService {
	return &PaymentService{queries: queries}
}

func (f *PaymentService) GetPaymentByID(ctx context.Context, id int32) (repository.Payment, error) {
	payment, err := f.queries.GetPayment(ctx, id)
	if err != nil {
		return repository.Payment{}, fmt.Errorf("get payment: %w", err)
	}

	return payment, nil
}

func (f *PaymentService) ListPaymentsByOrder(ctx context.Context, orderID int32) ([]repository.Payment, error) {
	payments, err := f.queries.ListPaymentsByOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("list payments by order: %w", err)
	}

	return payments, nil
}

func (f *PaymentService) InsertPayment(
	ctx context.Context,
	orderID int32,
	amount float64,
	method repository.PaymentMethod,
	transactionID *string,
	paymentTime *time.Time,
	status repository.PaymentStatus,
) (repository.Payment, error) {

	payment, err := f.queries.CreatePayment(ctx, repository.CreatePaymentParams{
		OrderID:       orderID,
		Amount:        helpers.NumericFromFloat64(amount),
		Method:        method,
		TransactionID: helpers.TextFromString(transactionID),
		PaymentTime:   helpers.TimestampFromTimePtr(paymentTime),
		Status:        status,
	})
	if err != nil {
		return repository.Payment{}, fmt.Errorf("create payment: %w", err)
	}

	return payment, nil
}

func (f *PaymentService) UpdatePaymentStatus(ctx context.Context, id int32, status repository.PaymentStatus) error {
	err := f.queries.UpdatePaymentStatus(ctx, repository.UpdatePaymentStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("update payment status: %w", err)
	}

	return nil
}

func (f *PaymentService) SetPaymentTransactionID(ctx context.Context, id int32, transactionID string) error {
	err := f.queries.SetPaymentTransactionID(ctx, repository.SetPaymentTransactionIDParams{
		ID:            id,
		TransactionID: helpers.TextFromString(&transactionID),
	})
	if err != nil {
		return fmt.Errorf("set payment transaction id: %w", err)
	}

	return nil
}

func (f *PaymentService) DeletePayment(ctx context.Context, id int32) error {
	err := f.queries.DeletePayment(ctx, id)
	if err != nil {
		return fmt.Errorf("delete payment: %w", err)
	}

	return nil
}
