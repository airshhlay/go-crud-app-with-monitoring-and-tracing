package middleware

import (
	config "gateway/config"
	metrics "gateway/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// increment the total number of requests
		metrics.TotalRequests.With(
			prometheus.Labels{
				"path": c.Request.URL.Path,
			},
		).Inc()
		c.Next()
	}
}
