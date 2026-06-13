package main

import (
	
	"log/slog"
	"net/http"
	"os"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
	"github.com/arnavsaroj/goratelimiter/internal/middleware"
	"github.com/arnavsaroj/goratelimiter/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/arnavsaroj/goratelimiter/internal/metrics"
)

func main() {


	logger:=slog.New(slog.NewJSONHandler(os.Stdout,nil))

	slog.SetDefault(logger)

	slog.Info("server-starting")

	slog.Info("server-started","port",8080)


	rdb := store.NewRedisConnection()
	manager := limiter.NewManager(rdb, 10, 10.0/60.0)

	manager.StartRedisHealthchecker()

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootFunc)
	mux.Handle("/metrics",promhttp.Handler())

	wrappedMux := middleware.RateLimiterMiddleware(manager)(mux)

	metrics.Init()

	http.ListenAndServe(":8080", wrappedMux)

}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
