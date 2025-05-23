package com.sc.orderservice;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.client.RestTemplateBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.web.client.RestTemplate;

import java.time.Duration;

@SpringBootApplication
public class OrderServiceApplication {

      public static void main(String[] args) {
            SpringApplication.run(OrderServiceApplication.class, args);
      }

      @Bean
      public RestTemplate restTemplate(RestTemplateBuilder builder) {
            return builder.setConnectTimeout(Duration.ofSeconds(3))
                  .setReadTimeout(Duration.ofSeconds(3))
                  .build();
      }
}
