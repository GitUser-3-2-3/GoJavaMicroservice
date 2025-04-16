package com.sc.orderservice.repository;

import com.sc.orderservice.model.Order;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface OrderRepository extends JpaRepository<Order, String> {
      List<Order> findByCustomerId(String customerId);
}
