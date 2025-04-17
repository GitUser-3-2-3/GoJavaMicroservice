package com.sc.orderservice.controller;

import com.sc.orderservice.model.Order;
import com.sc.orderservice.service.OrderService;
import io.micrometer.observation.annotation.Observed;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/orders")
@RequiredArgsConstructor
@Slf4j
@Observed(name = "orderController")
public class OrderController {

      private final OrderService orderService;

      @GetMapping("/{orderId}")
      public ResponseEntity<Order> getOrderById(@PathVariable String orderId) {
            log.info("Fetching order by id: {}", orderId);
            return ResponseEntity.ok(orderService.getOrderById(orderId));
      }

      @GetMapping("/customer/{customerId}")
      public ResponseEntity<List<Order>> getOrdersByCustomerId(@PathVariable String customerId) {
            log.info("Fetching orders by customer id: {}", customerId);
            return ResponseEntity.ok(orderService.getOrdersByCustomerId(customerId));
      }

      @PostMapping
      public ResponseEntity<Order> createOrder(@Valid @RequestBody Order order) {
            log.info("Creating order for customer: {}", order.getCustomerId());
            Order createdOrder = orderService.createOrder(order);
            return ResponseEntity.status(HttpStatus.CREATED).body(createdOrder);
      }

      @PostMapping("/{orderId}/payment")
      public ResponseEntity<Order> processPayment(@PathVariable String orderId) {
            log.info("Processing payment for order: {}", orderId);
            Order processedOrder = orderService.processPayment(orderId);
            return ResponseEntity.ok(processedOrder);
      }

      @PutMapping("/{id}/status")
      public ResponseEntity<Order> updateOrderStatus(@PathVariable String id,
                                                     @RequestParam Order.OrderStatus status) {
            log.info("Updating order {} status to: {}", id, status);
            Order updatedOrder = orderService.updateOrderStatus(id, status);
            return ResponseEntity.ok(updatedOrder);
      }
}
