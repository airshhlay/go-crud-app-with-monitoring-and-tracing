// package tracing

// import (
// 	"context"
// 	config "gateway/config"
// 	"log"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// )

// var (
// 	// todo: remove hardcoding
// 	CollectorURL = "localhost:5000/tracing"
// 	ServiceName  = "gateway"
// )

// func InitTracer(config *config.Config) func(context.Context) error {

// 	secureOption := otlptracegrpc.WithInsecure()

// 	exporter, err := otlptrace.New(
// 		context.Background(),
// 		otlptracegrpc.NewClient(
// 			secureOption,
// 			otlptracegrpc.WithEndpoint(CollectorURL),
// 		),
// 	)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	resources, err := resource.New(
// 		context.Background(),
// 		resource.WithAttributes(
// 			attribute.String(ServiceName),
// 			attribute.String("library.language", "go"),
// 		),
// 	)
// 	if err != nil {
// 		log.Printf("Could not set resources: ", err)
// 	}

// 	otel.SetTracerProvider(
// 		sdktrace.NewTracerProvider(
// 			sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 			sdktrace.WithBatcher(exporter),
// 			sdktrace.WithResource(resources),
// 		),
// 	)
// 	return exporter.Shutdown
// }
