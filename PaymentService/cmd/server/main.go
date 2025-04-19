package main

import (
	"PaymentService/internal/data"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	cfg, err := Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid log level, defaulting to info")
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	db, err := NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to close connection")
		}
	}(db)
	paymentRepo := data.NewPaymentRepository(db)
	paymentService := NewPaymentService(paymentRepo)
	paymentHandler := NewPaymentHandler(paymentService)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Route("/payments", func(r chi.Router) {
			r.Post("/", paymentHandler.CreatePayment)
			r.Get("/{id}", paymentHandler.GetPayment)
			r.Get("/order/{orderId}", paymentHandler.GetPaymentByOrderID)
		})
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
	})
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSecs) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSecs) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeoutSecs) * time.Second,
	}

	go func() {
		log.Info().Msgf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exited properly")
}
