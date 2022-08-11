package metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// create standard server metrics
	GrpcMetrics = grpc_prometheus.NewServerMetrics()
	// Create a metrics registry.
	Reg *prometheus.Registry = prometheus.NewRegistry()
	// Create a customized counter metric.
	TotalRequests  *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
)

func init() {
	// total requests
	TotalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_total_requests",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
	// Reg.MustRegister(GrpcMetrics, TotalRequests)
	TotalRequests.WithLabelValues("test")

	// request latency
	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Measures the duration taken for each request",
			Buckets: prometheus.LinearBuckets(0.01, 0.05, 10),
		},
		[]string{"service_label", "path", "errorCode"},
	)
	// Reg.MustRegister(GrpcMetrics, RequestLatency)
}
