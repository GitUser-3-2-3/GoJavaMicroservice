package main

import (
	"PaymentService/internal/data"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type PaymentHandler struct {
	service   *PaymentService
	validator *validator.Validate
}

func NewPaymentHandler(service *PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service,
		validator: validator.New(),
	}
}

func (ph *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	logger := log.With().Str("request_id", requestID).Logger()

	var req data.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode request body")
		respondWithError(w, http.StatusBadRequest, "Invalid request format", requestID)
		return
	}
	if err := ph.validator.Struct(req); err != nil {
		logger.Error().Err(err).Interface("req", req).Msg("Invalid request")
		respondWithError(w, http.StatusBadRequest, "Validation Error: "+err.Error(), requestID)
		return
	}
	logger.Info().Str("order_id", req.OrderID).
		Str("customer_id", req.CustomerID).
		Float64("amount", req.Amount).Msg("Processing payment")

	ctx := r.Context()
	payment, err := ph.service.ProcessPayment(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrPaymentFailed):
			logger.Error().Err(err).Msg("Payment processing failed")
			respondWithError(w, http.StatusBadRequest, err.Error(), requestID)
			return
		default:
			logger.Error().Err(err).Msg("Internal error processing payment")
			respondWithError(w, http.StatusInternalServerError, "Error processing payment", requestID)
			return
		}
	}
	response := data.PaymentResponse{PaymentID: payment.ID,
		OrderID:     payment.OrderID,
		CustomerID:  payment.CustomerID,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		Status:      payment.Status,
		ProcessedAt: payment.ProcessedAt,
		ErrorMsg:    payment.ErrorMsg,
	}
	logger.Info().Str("payment_id", payment.ID).
		Str("status", payment.Status).Msg("Payment processed successfully")

	respondWithJSON(w, http.StatusOK, response)
}

func (ph *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	requestID := getRequestID(r)
	logger := log.With().Str("request_id", requestID).Str("payment_id", id).Logger()

	payment, err := ph.service.GetPaymentById(r.Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("Error retrieving payment")
		respondWithError(w, http.StatusInternalServerError, "Error retrieving payment", requestID)
		return
	}
	if payment == nil {
		logger.Info().Msg("Payment not found")
		respondWithError(w, http.StatusNotFound, "Payment not found", requestID)
		return
	}
	response := data.PaymentResponse{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		CustomerID:  payment.CustomerID,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		Status:      payment.Status,
		ProcessedAt: payment.ProcessedAt,
		ErrorMsg:    payment.ErrorMsg,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (ph *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	requestID := getRequestID(r)
	logger := log.With().Str("request_id", requestID).Str("order_id", orderID).Logger()

	payment, err := ph.service.GetPaymentByOrderId(r.Context(), orderID)
	if err != nil {
		logger.Error().Err(err).Msg("Error retrieving payment")
		respondWithError(w, http.StatusInternalServerError, "Error retrieving payment", requestID)
		return
	}
	if payment == nil {
		logger.Info().Msg("Payment not found for order")
		respondWithError(w, http.StatusNotFound, "Payment not found for order", requestID)
		return
	}

	response := data.PaymentResponse{
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		CustomerID:  payment.CustomerID,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		Status:      payment.Status,
		ProcessedAt: payment.ProcessedAt,
		ErrorMsg:    payment.ErrorMsg,
	}

	respondWithJSON(w, http.StatusOK, response)
}
