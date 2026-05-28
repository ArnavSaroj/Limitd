package middleware

import (
	"net/http"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
)

func RateLimiterMiddleware(manager *limiter.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//extract the ip address of the request however this is not very good as if u use load balancer this wil return the ip of lb lmao 
			//so get x forwaded for port instead of this
			ip := r.Header.Get("X-Forwaded-For")
			if ip == "" {
				ip = r.RemoteAddr

			}

			bucket := manager.GetBucket(ip)

			if !bucket.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
