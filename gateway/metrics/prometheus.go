package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalRequests  *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
)

// prometheus handler for the /metrics endpoint
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func Init() error {
	// total requests
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests to the gateway.",
		},
		[]string{"path"},
	)
	prometheus.MustRegister(TotalRequests)

	// request latency
	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Measures the duration taken for each request",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
		},
		[]string{"service_label", "path", "errorCode"},
	)
	prometheus.MustRegister(RequestLatency)

	return nil
}
