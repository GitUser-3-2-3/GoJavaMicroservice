package model

import "time"

const (
	PaymentStatusPending   = "PENDING"
	PaymentStatusCompleted = "COMPLETED"
	PaymentStatusFailed    = "FAILED"
)

type Payment struct {
	ID           string    `json:"paymentId" db:"id"`
	OrderID      string    `json:"orderId" db:"order_id"`
	CustomerID   string    `json:"customerId" db:"customer_id"`
	Amount       int       `json:"amount" db:"amount"`
	Currency     string    `json:"currency" db:"currency"`
	Status       string    `json:"status" db:"status"`
	ProcessedAt  time.Time `json:"processedAt" db:"processed_at"`
	ErrorMessage string    `json:"errorMessage,omitempty" db:"error_message"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type PaymentRequest struct {
	OrderID    string `json:"orderId" validate:"required"`
	CustomerID string `json:"customerId" validate:"required"`
	Amount     int    `json:"amount" validate:"required,gt=0"`
	Currency   string `json:"currency" validate:"required,len=3"`
}

type PaymentResponse struct {
	PaymentID    string    `json:"paymentId"`
	OrderID      string    `json:"orderId"`
	CustomerID   string    `json:"customerId"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	Status       string    `json:"status"`
	ProcessedAt  time.Time `json:"processedAt"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

type ErrorResponse struct {
	Error     string `json:"error"`
	Code      int    `json:"code,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}
