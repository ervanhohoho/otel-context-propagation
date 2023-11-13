package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"moul.io/http2curl"
	"net/http"
)

func main() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	config := OpenTelemetryConfig{
		URL:         "54.255.227.155:4317",
		Timeout:     30,
		DialOption:  opts,
		ServiceName: "test-otel1",
	}

	otelService := OpenTelemetryImpl{OpenTelemetryConfig: config}

	shutdown, tracer := otelService.InitiateTracer()
	defer shutdown()

	ctx, span := tracer.Start(context.Background(), "Test")
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/record", nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))
	request.WithContext(ctx)
	command, _ := http2curl.GetCurlCommand(request)
	fmt.Println(request.Header.Get("Traceparent"))
	fmt.Println(command)
	client := http.Client{}

	client.Do(request)
	defer span.End()

}
