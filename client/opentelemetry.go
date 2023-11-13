package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/propagation"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type OpenTelemetry interface {
	InitiateTracer() (func(), trace.Tracer)
}

type OpenTelemetryImpl struct {
	OpenTelemetryConfig OpenTelemetryConfig
}

// This is so that it can be mocked by mockery
type OtelTracer interface {
	trace.Tracer
}

// This is so that it can be mocked by mockery
type OtelSpan interface {
	trace.Span
}

func (o OpenTelemetryImpl) newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	conn, err := grpc.DialContext(ctx, o.OpenTelemetryConfig.URL,
		o.OpenTelemetryConfig.DialOption...,
	)
	if err != nil {
		return nil, fmt.Errorf("gRPC connection error: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("error craeting otlp trace exporter: %w", err)
	}

	return exporter, nil
}

func (o OpenTelemetryImpl) newTraceProvider(exporter sdktrace.SpanExporter, serviceName string) *sdktrace.TracerProvider {
	defaultResource := resource.Default()
	r, err := resource.Merge(
		defaultResource,
		resource.NewWithAttributes(
			defaultResource.SchemaURL(), // This import needs to be the same as resouce version
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to initialize resource: %v", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
}

func (o OpenTelemetryImpl) InitiateTracer() (func(), trace.Tracer) {
	ctx := context.Background()
	exp, err := o.newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tp := o.newTraceProvider(exp, o.OpenTelemetryConfig.ServiceName)

	otel.SetTracerProvider(tp)

	tracer := tp.Tracer(o.OpenTelemetryConfig.ServiceName)

	return func() {
		_ = tp.Shutdown(ctx)
	}, tracer
}
