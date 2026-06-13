package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total requests processed",
		},
		[]string{"status"},
	)

	RedisErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_errors_total",
			Help: "Total Redis errors",
		},
	)

	FallbackRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "fallback_request_total",
			Help: "Request served when fallback in memory limiter starts working",
		},
	)

	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "requests_duration_seconds",
			Help:    "requests processing latency",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Init() {
	prometheus.MustRegister(
		RequestsTotal,
		RedisErrorsTotal,
		FallbackRequestsTotal,
		RequestDuration,
	)
}
