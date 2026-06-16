package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arnavsaroj/goratelimiter/internal/handlers"
	"github.com/arnavsaroj/goratelimiter/internal/limiter"
	"github.com/arnavsaroj/goratelimiter/internal/middleware"
	"github.com/arnavsaroj/goratelimiter/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/arnavsaroj/goratelimiter/internal/metrics"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)



	rdb := store.NewRedisConnection()
	manager := limiter.NewManager(rdb, 10, 10.0/60.0)

	manager.StartRedisHealthchecker()

	mux := http.NewServeMux()
	metrics.Init()

	mux.HandleFunc("/", rootFunc)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", handlers.Healthz)

	wrappedMux := middleware.RateLimiterMiddleware(manager)(mux)


	srv := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()
	slog.Info("server-started", "port", 8080)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}

	if err := rdb.Close(); err != nil {
		logger.Error("redis close failed", "error", err)
	}

	logger.Info("server exited gracefully")

}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
