package otel

import (
	"context"
	"net/http"

	"github.com/shashtag-ventures/go-common/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	// Register the OTEL trace ID extractor with the logging middleware
	middleware.RegisterTraceIDExtractor(func(ctx context.Context) string {
		if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
			return span.SpanContext().TraceID().String()
		}
		return ""
	})
}

// Middleware creates a new OpenTelemetry middleware.
// It wraps the handler with otelhttp.NewHandler to collect traces.
func Middleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName)
	}
}
