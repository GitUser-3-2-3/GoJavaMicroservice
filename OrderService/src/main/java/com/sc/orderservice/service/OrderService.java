package com.sc.orderservice.service;

import com.sc.orderservice.exception.OrderNotFoundException;
import com.sc.orderservice.exception.PaymentServiceException;
import com.sc.orderservice.model.Order;
import com.sc.orderservice.model.PaymentRequest;
import com.sc.orderservice.model.PaymentResponse;
import com.sc.orderservice.repository.OrderRepository;
import io.micrometer.observation.annotation.Observed;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.retry.annotation.Backoff;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResourceAccessException;
import org.springframework.web.client.RestTemplate;

import java.util.List;

@Service
@RequiredArgsConstructor
@Slf4j
@Observed(name = "orderService")
public class OrderService {

      private final OrderRepository orderRepository;
      private final RestTemplate restTemplate;

      @Value("${payment.service.url}")
      private String paymentServiceUrl;

      @Transactional(readOnly = true)
      public Order getOrderById(String orderId) {
            return orderRepository.findById(orderId)
                  .orElseThrow(() -> new OrderNotFoundException("Order not found for id: " + orderId));
      }

      @Transactional(readOnly = true)
      public List<Order> getOrdersByCustomerId(String customerId) {
            return orderRepository.findByCustomerId(customerId);
      }

      @Transactional
      public Order createOrder(Order order) {
            log.info("Creating order for customer: {}", order.getCustomerId());
            return orderRepository.save(order);
      }

      @Transactional
      @Retryable(
            retryFor = {ResourceAccessException.class, HttpServerErrorException.class},
            maxAttempts = 5,
            backoff = @Backoff(delay = 1000, multiplier = 2, maxDelay = 10000)
      )
      public Order processPayment(String orderId) {
            log.info("processing payment for order: {}", orderId);

            Order order = getOrderById(orderId);

            if (Order.OrderStatus.PAID.equals(order.getStatus())) {
                  log.info("Order {} already paid", orderId);
                  return order;
            }
            PaymentRequest paymentRequest = PaymentRequest.builder().orderId(order.getOrderId())
                  .customerId(order.getCustomerId()).amount(order.getTotalAmount())
                  .currency("USD")
                  .build();

            try {
                  PaymentResponse response = restTemplate.postForObject(
                        paymentServiceUrl + "/api/payments", paymentRequest, PaymentResponse.class);
                  if (response == null) {
                        throw new PaymentServiceException("Null response from payment service");
                  }
                  log.info("Payment processed with id: {}, status: {}", response.getPaymentId(), response.getStatus());

                  if ("COMPLETED".equals(response.getStatus())) {
                        order.setStatus(Order.OrderStatus.PAID);
                        order.setPaymentId(response.getPaymentId());
                        return orderRepository.save(order);
                  } else {
                        throw new PaymentServiceException("Payment failed: " + response.getErrorMessage());
                  }
            } catch (HttpClientErrorException e) {
                  log.error("Payment service client error: {}", e.getMessage());
                  if (e.getStatusCode() == HttpStatus.BAD_REQUEST) {
                        throw new PaymentServiceException("Invalid payment request: " + e.getResponseBodyAsString());
                  }
                  throw new PaymentServiceException("Payment service error: " + e.getMessage());
            } catch (ResourceAccessException e) {
                  log.error("Cannot connect to payment service: {}", e.getMessage());
                  throw new PaymentServiceException("Payment service unavailable: " + e.getMessage());
            } catch (Exception e) {
                  log.error("Unexpected error processing payment: {}", e.getMessage());
                  throw new PaymentServiceException("Payment processing error: " + e.getMessage());
            }
      }

      @Transactional
      public Order updateOrderStatus(String orderId, Order.OrderStatus status) {
            Order order = getOrderById(orderId);
            order.setStatus(status);
            return orderRepository.save(order);
      }
}
