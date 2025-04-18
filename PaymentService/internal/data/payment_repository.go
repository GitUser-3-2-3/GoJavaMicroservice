package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PaymentModel struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentModel {
	return &PaymentModel{db: db}
}
func (mdl *PaymentModel) Create(ctx context.Context, payment *Payment) error {
	query := `INSERT INTO payments (id, order_id, customer_id, amount, currency, 
		    status, processed_at, error_msg, created_at, updated_at) 
		    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	now := time.Now()

	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = now
	}
	if payment.ProcessedAt.IsZero() {
		payment.ProcessedAt = now
	}
	_, err := mdl.db.ExecContext(ctx, query, payment.ID,
		payment.OrderID,
		payment.CustomerID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.ProcessedAt,
		payment.ErrorMsg,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Str("order_id", payment.OrderID).
			Str("customer_id", payment.CustomerID).Msg("Failed to create payment")
		return err
	}
	return nil
}

func (mdl *PaymentModel) GetByPaymentID(ctx context.Context, id string) (*Payment, error) {
	query := `SELECT id, order_id, customer_id, amount, currency, status, 
                processed_at, error_msg, created_at, updated_at 
                FROM payments WHERE id = $1`

	payment := new(Payment)
	err := mdl.db.QueryRowContext(ctx, query, id).Scan(&payment.ID,
		&payment.OrderID,
		&payment.CustomerID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.ProcessedAt,
		&payment.ErrorMsg,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			log.Error().Err(err).Str("id", id).Msg("Failed to get payment")
			return nil, err
		}
	}
	return payment, nil
}

func (mdl *PaymentModel) GetByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	query := `SELECT id, order_id, customer_id, amount, currency, status, 
                processed_at, error_msg, created_at, updated_at 
                FROM payments WHERE order_id = $1`
	payment := new(Payment)

	err := mdl.db.QueryRowContext(ctx, query, orderID).Scan(&payment.ID,
		&payment.OrderID,
		&payment.CustomerID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.ProcessedAt,
		&payment.ErrorMsg,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			log.Error().Err(err).Str("order_id", orderID).Msg("Failed to get payment by order ID")
			return nil, err
		}
	}
	return payment, nil
}

func (mdl *PaymentModel) Update(ctx context.Context, payment *Payment) error {
	query := `UPDATE payments SET status = $1, error_msg = $2, updated_at = $3
                WHERE id = $4`

	payment.UpdatedAt = time.Now()
	args := []any{payment.Status, payment.ErrorMsg, payment.UpdatedAt, payment.ID}

	_, err := mdl.db.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			log.Error().Err(err).Str("id", payment.ID).
				Str("status", payment.Status).Msg("Failed to update payment")
			return err
		}
	}
	return nil
}
