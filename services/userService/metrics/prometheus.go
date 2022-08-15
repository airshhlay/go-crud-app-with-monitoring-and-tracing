package metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// GrpcMetrics collects standard server metrics on the grpc server.
	GrpcMetrics *grpc_prometheus.ServerMetrics
	// Reg is a ServerRegistry.
	Reg *prometheus.Registry
	// TotalRequests counts total requests.
	TotalRequests *prometheus.CounterVec
	// RequestDuration tracks request duration.
	RequestDuration *prometheus.HistogramVec
	// DatabaseOpDuration tracks database op durations.
	DatabaseOpDuration *prometheus.HistogramVec
)

func init() {
	GrpcMetrics = grpc_prometheus.NewServerMetrics()
	Reg = prometheus.NewRegistry()

	// request latency
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Measures the duration taken for each request",
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 2},
		},
		[]string{"service_label", "name", "errorCode"},
	)

	DatabaseOpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "process_database_op_duration_seconds",
			Help:    "Measures the duration taken for a database operation",
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 2},
		},
		[]string{"service_label", "query_type", "query_label", "success"},
	)

	// register collectors
	Reg.MustRegister(GrpcMetrics, RequestDuration, DatabaseOpDuration)
}
