package metrics

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// create standard server metrics
	GrpcMetrics *grpc_prometheus.ServerMetrics
	// Create a metrics registry.
	Reg *prometheus.Registry
	// Create a customized counter metric.
	TotalRequests              *prometheus.CounterVec
	RequestLatency             *prometheus.HistogramVec
	DatabaseOpDuration         *prometheus.HistogramVec
	PasswordEncryptionDuration *prometheus.HistogramVec
)

func init() {
	GrpcMetrics = grpc_prometheus.NewServerMetrics()
	Reg = prometheus.NewRegistry()

	// request latency
	RequestLatency = prometheus.NewHistogramVec(
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

	PasswordEncryptionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "process_encryption_op_duration_seconds",
			Help:    "Measures the duration taken to bcrypt a password",
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 2},
		},
		[]string{},
	)

	// register collectors
	Reg.MustRegister(GrpcMetrics, RequestLatency, DatabaseOpDuration, PasswordEncryptionDuration)
}
