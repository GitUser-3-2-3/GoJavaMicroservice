package main

import (
	"PaymentService/internal/data"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrPaymentFailed = errors.New("payment processing failed")
)

type PaymentModel interface {
	Create(ctx context.Context, payment *data.Payment) error
	Update(ctx context.Context, payment *data.Payment) error
	GetByPaymentID(ctx context.Context, id string) (*data.Payment, error)
	GetByOrderID(ctx context.Context, orderID string) (*data.Payment, error)
}

type PaymentService struct {
	mdl PaymentModel
}

func NewPaymentService(mdl PaymentModel) *PaymentService {
	return &PaymentService{mdl: mdl}
}

func (ps *PaymentService) ProcessPayment(ctx context.Context, req data.PaymentRequest) (*data.Payment, error) {
	logger := log.With().Str("order_id", req.OrderID).
		Str("customer_id", req.CustomerID).Float64("amount", req.Amount).Logger()

	existingPayment, err := ps.mdl.GetByOrderID(ctx, req.OrderID)
	if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
		logger.Error().Err(err).Msg("Error getting existing payment")
		return nil, fmt.Errorf("failed to check existing payment: %w", err)
	}
	if existingPayment != nil && existingPayment.Status == data.PaymentStatusCompleted {
		logger.Info().Str("payment_id", existingPayment.ID).Msg("Payment already processed")
		return existingPayment, nil
	}
	payment := &data.Payment{
		ID:          uuid.New().String(),
		OrderID:     req.OrderID,
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Status:      data.PaymentStatusPending,
		ProcessedAt: time.Now(),
	}
	paymentSuccessful, errMsg := ps.simulatePaymentProcessing(req)

	if paymentSuccessful {
		payment.Status = data.PaymentStatusCompleted
	} else {
		payment.Status = data.PaymentStatusFailed
		payment.ErrorMsg = errMsg
	}
	if existingPayment != nil {
		existingPayment.ErrorMsg = payment.ErrorMsg
		existingPayment.Status = payment.Status
		existingPayment.ProcessedAt = payment.ProcessedAt

		if err := ps.mdl.Update(ctx, existingPayment); err != nil {
			logger.Error().Err(err).Msg("Failed to update payment")
			return nil, fmt.Errorf("failed to update payment: %w", err)
		}
		payment = existingPayment
	} else {
		if err := ps.mdl.Create(ctx, payment); err != nil {
			logger.Error().Err(err).Msg("Failed to create payment")
			return nil, fmt.Errorf("failed to create payment: %w", err)
		}
	}
	if payment.Status == data.PaymentStatusFailed {
		return payment, ErrPaymentFailed
	}
	return payment, nil
}

func (ps *PaymentService) GetPaymentById(ctx context.Context, id string) (*data.Payment, error) {
	payment, err := ps.mdl.GetByPaymentID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to get payment by payment id")
		return nil, err
	}
	return payment, nil
}

func (ps *PaymentService) GetPaymentByOrderId(ctx context.Context, id string) (*data.Payment, error) {
	payment, err := ps.mdl.GetByOrderID(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to get payment by order ID")
		return nil, err
	}
	return payment, nil
}

func (ps *PaymentService) simulatePaymentProcessing(_ data.PaymentRequest) (bool, string) {
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	if rand.Float64() < 0.95 {
		return true, ""
	}
	reasons := []string{"Insufficient funds",
		"Card declined",
		"Payment gateway timeout",
		"Invalid card details",
	}
	return false, reasons[rand.Intn(len(reasons))]
}
