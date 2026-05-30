package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method and status code",
		},
		[]string{"method", "status"},
	)

	HttpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_requests_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestsTotal, HttpDuration)
}
