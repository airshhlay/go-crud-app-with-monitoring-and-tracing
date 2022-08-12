package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalRequests        *prometheus.CounterVec
	RequestLatency       *prometheus.HistogramVec
	ResponseSize         *prometheus.HistogramVec
	AuthenticateDuration *prometheus.HistogramVec
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
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 2},
		},
		[]string{"service_label", "path", "errorCode"},
	)
	prometheus.MustRegister(RequestLatency)

	// response size
	ResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Measures the size of each response in bytes",
			Buckets: prometheus.LinearBuckets(0, 450, 8),
		},
		[]string{"service_label", "path", "errorCode"},
	)
	prometheus.MustRegister(ResponseSize)

	AuthenticateDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "process_authenticate_duration_seconds",
			Help:    "Measures the duration taken to sign and issue a jwt token to the user.",
			Buckets: prometheus.LinearBuckets(0, 0.01, 6),
		},
		[]string{"success"},
	)
	prometheus.MustRegister(AuthenticateDuration)

	return nil
}
