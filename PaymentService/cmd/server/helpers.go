package main

import (
	"PaymentService/internal/data"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func getRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value(middleware.RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

func respondWithError(w http.ResponseWriter, code int, msg, requestID string) {
	respondWithJSON(w, code, data.ErrorResponse{
		Error:     msg,
		Code:      code,
		RequestID: requestID,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "error marshalling JSON response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
