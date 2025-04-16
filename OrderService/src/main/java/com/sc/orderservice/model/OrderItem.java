package com.sc.orderservice.model;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;

@Embeddable
@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class OrderItem {

      @NotEmpty
      private String itemId;

      @NotEmpty
      private String itemName;

      @NotNull
      @Positive
      private BigDecimal price;

      @NotNull
      @Positive
      private Integer quantity;
}
