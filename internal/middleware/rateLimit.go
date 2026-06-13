package middleware

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
	"github.com/arnavsaroj/goratelimiter/internal/metrics"
)

func RateLimiterMiddleware(manager *limiter.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			defer func() {
				metrics.RequestDuration.Observe(
					time.Since(start).Seconds(),
				)
			}()

			metrics.RequestsTotal.WithLabelValues("incoming").Inc()

			//extract the ip address of the request however this is not very good as if u use load balancer this wil return the ip of the lb lmao
			//so get x forwaded for port instead of this
			ip := r.Header.Get("X-Forwarded-For")
			if ip != "" {
				ip = strings.Split(ip, ",")[0]
				ip = strings.TrimSpace(ip)
			}
			if ip == "" {

				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err == nil {
					ip = host
				} else {
					ip = r.RemoteAddr
				}
			}

			ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
			defer cancel()

			if !manager.RedisHealthy.Load() {

				localBucket := manager.GetLocalBucket(ip)
				metrics.FallbackRequestsTotal.Inc()

				if !localBucket.Allow() {
					metrics.RequestsTotal.WithLabelValues("blocked").Inc()
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}
				metrics.RequestsTotal.WithLabelValues("allowed").Inc()
				next.ServeHTTP(w, r)
				return
			}

			bucket := manager.GetBucket(ip)

			result, err := bucket.Allow(ctx)

			if err == nil {
				manager.ConsecutiveFailures.Store(0)

				manager.RedisHealthy.Store(true)
			}

			if err != nil {
				metrics.RedisErrorsTotal.Inc()
				failures := manager.ConsecutiveFailures.Add(1)

				if failures >= 3 {
					if manager.RedisHealthy.Load() {
						slog.Info("circuit breaker opened")
					}
					manager.RedisHealthy.Store(false)

				}
				slog.Error("redis_error", "error", err)

				localBucket := manager.GetLocalBucket(ip)
				metrics.FallbackRequestsTotal.Inc()

				slog.Warn("fallback_activated",
					"reason", err.Error(),
				)

				if !localBucket.Allow() {
					metrics.RequestsTotal.WithLabelValues("blocked").Inc()
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}
				metrics.RequestsTotal.WithLabelValues("allowed").Inc()
				next.ServeHTTP(w, r)
				return

			} else if result == false {
				metrics.RequestsTotal.WithLabelValues("blocked").Inc()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			//earlier code
			// if !bucket.Allow(ctx) {
			// 	http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			// 	return
			// }
			metrics.RequestsTotal.WithLabelValues("allowed").Inc()
			next.ServeHTTP(w, r)
		})
	}
}
