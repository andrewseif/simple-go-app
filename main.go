package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer       trace.Tracer
	otlpEndpoint string
)

// List of supported exporters
// https://opentelemetry.io/docs/instrumentation/go/exporters/

// Console Exporter, only for testing
func newConsoleExporter() (oteltrace.SpanExporter, error) {
	return stdouttrace.New()
}

// TracerProvider is an OpenTelemetry TracerProvider.
// It provides Tracers to instrumentation so it can trace operational flow through a system.
func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("simple-go-app"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	ctx := context.Background()

	// For testing to print out traces to the console
	exp, err := newConsoleExporter()

	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer("simple-go-app")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/items", getItems)
	http.ListenAndServe(":8080", r)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HTTP GET /items")
	defer span.End()

	// Simulate Database call to fetch connected devices.
	db(ctx)

	// Add additional delay to simulate HTTP request.
	time.Sleep(1 * time.Second)

	// Return devices
	w.Write([]byte("ok"))
}

func db(ctx context.Context) {
	_, span := tracer.Start(ctx, "SQL SELECT")
	defer span.End()

	// Simulate Database call to SELECT connected devices.
	time.Sleep(2 * time.Second)
}
