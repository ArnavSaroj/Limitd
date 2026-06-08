package main

import (
	"fmt"
	"net/http"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
	"github.com/arnavsaroj/goratelimiter/internal/middleware"
	"github.com/arnavsaroj/goratelimiter/internal/store"
)

func main() {

	rdb := store.NewRedisConnection()
	manager := limiter.NewManager(rdb, 10, 10.0/60.0)

	manager.StartRedisHealthchecker()

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootFunc)

	wrappedMux := middleware.RateLimiterMiddleware(manager)(mux)

	fmt.Print("backend running on port 8080")
	http.ListenAndServe(":8080", wrappedMux)

}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
