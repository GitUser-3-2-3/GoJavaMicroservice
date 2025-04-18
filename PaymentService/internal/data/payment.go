package data

import (
	"errors"
	"time"
)

const (
	PaymentStatusPending   = "PENDING"
	PaymentStatusCompleted = "COMPLETED"
	PaymentStatusFailed    = "FAILED"
)

var ErrRecordNotFound = errors.New("record not found")

type Payment struct {
	ID          string    `json:"paymentId" db:"id"`
	OrderID     string    `json:"orderId" db:"order_id"`
	CustomerID  string    `json:"customerId" db:"customer_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Currency    string    `json:"currency" db:"currency"`
	Status      string    `json:"status" db:"status"`
	ProcessedAt time.Time `json:"processedAt" db:"processed_at"`
	ErrorMsg    string    `json:"errorMsg,omitempty" db:"error_msg"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type PaymentRequest struct {
	OrderID    string  `json:"orderId" validate:"required"`
	CustomerID string  `json:"customerId" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Currency   string  `json:"currency" validate:"required,len=3"`
}

type PaymentResponse struct {
	PaymentID   string    `json:"paymentId"`
	OrderID     string    `json:"orderId"`
	CustomerID  string    `json:"customerId"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processedAt"`
	ErrorMsg    string    `json:"errorMsg,omitempty"`
}

type ErrorResponse struct {
	Error     string `json:"error"`
	Code      int    `json:"code,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}
