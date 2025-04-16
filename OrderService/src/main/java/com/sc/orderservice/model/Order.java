package com.sc.orderservice.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Entity
@Table(name = "orders")
@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class Order {

      @Id
      private String orderId;

      @NotNull
      private String customerId;

      @Enumerated(EnumType.STRING)
      private OrderStatus status;

      @NotNull
      @Positive
      private BigDecimal totalAmount;

      @NotNull
      private LocalDateTime createdAt;

      @ElementCollection
      @CollectionTable(name = "order_items", joinColumns = @JoinColumn(name = "order_id"))
      private List<OrderItem> items = new ArrayList<>();

      private String paymentId;

      @PrePersist
      public void prePersist() {
            if (orderId == null) {
                  orderId = UUID.randomUUID().toString();
            }
            if (createdAt == null) {
                  createdAt = LocalDateTime.now();
            }
            if (status == null) {
                  status = OrderStatus.PENDING;
            }
      }

      public enum OrderStatus {
            CREATED, PENDING, PREPARING, COMPLETED, FULFILLED
      }
}
