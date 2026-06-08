package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/arnavsaroj/goratelimiter/internal/limiter"
)

func RateLimiterMiddleware(manager *limiter.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Millisecond)
			defer cancel()

			if !manager.RedisHealthy.Load() {


				fmt.Println("redis unhealthy!!")
				fmt.Println("falling back to local bucket in memory bucket implementation")

				localBucket := manager.GetLocalBucket(ip)

				if !localBucket.Allow() {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			bucket := manager.GetBucket(ip)

			result, err := bucket.Allow(ctx)

			if err==nil{
				manager.ConsecutiveFailures.Store(0)

				manager.RedisHealthy.Store(true)
			}

			if err != nil {

				failures := manager.ConsecutiveFailures.Add(1)

				if failures >= 3 {
					if manager.RedisHealthy.Load() {
						fmt.Println("Circuit Breaker opened")

					}
					manager.RedisHealthy.Store(false)

				}

				fmt.Println("redis down!!")
				fmt.Println("falling back to local bucket in memory bucket implementation")

				localBucket := manager.GetLocalBucket(ip)

				if !localBucket.Allow() {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
				return

			} else if result == false {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			//earlier code
			// if !bucket.Allow(ctx) {
			// 	http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			// 	return
			// }
			next.ServeHTTP(w, r)
		})
	}
}
