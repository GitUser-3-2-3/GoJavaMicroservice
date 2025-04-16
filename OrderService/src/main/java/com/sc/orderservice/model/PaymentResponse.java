package com.sc.orderservice.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class PaymentResponse {

      private String paymentId;
      private String orderId;
      private String customerId;
      private BigDecimal amount;
      private String currency;
      private String status;
      private LocalDateTime processedAt;
      private String errorMessage;
}
