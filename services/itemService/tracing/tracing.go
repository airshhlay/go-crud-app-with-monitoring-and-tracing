package tracing

import (
	"io"
	config "itemService/config"
	constants "itemService/constants"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// InitJaeger creates a new tracer.
func InitJaeger(config *config.JaegerConfig, logger *zap.Logger) (opentracing.Tracer, io.Closer, error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: config.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.LogSpans,
			LocalAgentHostPort: config.Host,
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		logger.Fatal(constants.ErrorJaegerInitMsg, zap.Error(err))
	}

	return tracer, closer, err
}
