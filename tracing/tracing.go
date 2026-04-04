package tracing

import (
	"context"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitTracer initializes an OpenTelemetry tracer provider.
// It configures a stdout exporter for demonstration purposes.
// In a production environment, this would typically export to a tracing backend like Jaeger or Zipkin.
// It returns a cleanup function that should be deferred in the main function to ensure proper shutdown.
func InitTracer(serviceName string) func() {
	var exporter trace.SpanExporter
	var err error

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint != "" {
		exporter, err = otlptracehttp.New(context.Background())
		if err != nil {
			log.Fatalf("failed to create OTLP exporter: %v", err)
		}
		log.Println("Initialized OpenTelemetry with OTLP HTTP exporter")
	} else {
		// Create a new stdout exporter.
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(io.Discard), // Discard output for now, will use slog later
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
		if err != nil {
			log.Fatalf("failed to create stdout exporter: %v", err)
		}
		log.Println("Initialized OpenTelemetry with stdout base exporter")
	}

	// Create a resource that describes this application.
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	// Create a new TracerProvider with the exporter and resource.
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter), // Use a batcher for efficiency.
		trace.WithResource(res),
	)
	// Set the global TracerProvider.
	otel.SetTracerProvider(provider)

	// Return a cleanup function to shut down the provider gracefully.
	return func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}
