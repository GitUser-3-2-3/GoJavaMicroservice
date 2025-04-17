package com.sc.orderservice.exception;

import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

      @ExceptionHandler(OrderNotFoundException.class)
      public ResponseEntity<ErrorResponse> handleOrderNotFound(OrderNotFoundException exp) {
            log.warn("Order not found: {}", exp.getMessage());

            return ResponseEntity
                  .status(HttpStatus.NOT_FOUND)
                  .body(new ErrorResponse(exp.getMessage()));
      }

      @ExceptionHandler(PaymentServiceException.class)
      public ResponseEntity<ErrorResponse> handlePaymentServiceException(PaymentServiceException exp) {
            log.error("Payment service exception: {}", exp.getMessage());

            return ResponseEntity
                  .status(HttpStatus.SERVICE_UNAVAILABLE)
                  .body(new ErrorResponse("Payment processing failed:: " + exp.getMessage()));
      }

      @ExceptionHandler(Exception.class)
      public ResponseEntity<ErrorResponse> handleGeneralException(Exception exp) {
            log.error("Unexpected error: {}", exp.getMessage(), exp);

            return ResponseEntity
                  .status(HttpStatus.INTERNAL_SERVER_ERROR)
                  .body(new ErrorResponse("Internal Server Error:: " + exp.getMessage()));
      }

      public record ErrorResponse(String message) {}
}
