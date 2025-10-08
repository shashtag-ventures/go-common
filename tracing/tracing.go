package tracing

import (
	"context"
	"io"
	"log"

	"go.opentelemetry.io/otel"
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
	// Create a new stdout exporter.
	// For simplicity, output is discarded here, but in a real scenario, it would go to a file or console.
	exporter, err := stdouttrace.New(
		stdouttrace.WithWriter(io.Discard), // Discard output for now, will use slog later
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	// Create a resource that describes this application.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	// Create a new TracerProvider with the exporter and resource.
	provider := trace.NewTracerProvider(
		trace.WithBatcher(exporter), // Use a batcher for efficiency.
		trace.WithResource(resource),
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
