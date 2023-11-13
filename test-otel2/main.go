package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	config := OpenTelemetryConfig{
		URL:         "54.255.227.155:4317",
		Timeout:     30,
		DialOption:  opts,
		ServiceName: "test-otel2",
	}
	otelService := OpenTelemetryImpl{OpenTelemetryConfig: config}
	shutdown, tracer := otelService.InitiateTracer()
	defer shutdown()

	r := gin.Default()
	r.Use(otelgin.Middleware("otel2"))
	r.GET("/record", func(c *gin.Context) {
		//traceId := c.GetHeader("traceid")
		//spanid := c.GetHeader("spanid")
		//ginCtx := c.Copy()
		//ginCtx.Set("traceID", traceId)
		//ginCtx.Set("spanID", spanid)
		//remoteSpan := trace.SpanFromContext(ginCtx)
		//ctx := trace.ContextWithRemoteSpanContext(
		//	context.Background(),
		//	trace.NewSpanContext(trace.SpanContextConfig{
		//		TraceID:    remoteSpan.SpanContext().TraceID(),
		//		SpanID:     remoteSpan.SpanContext().SpanID(),
		//		TraceFlags: trace.FlagsSampled,
		//	}),
		//)
		fmt.Println(c.Request.Header.Get("Traceparent"))
		//ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		_, span := tracer.Start(c.Request.Context(), "hello from otel 2")
		defer span.End()
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
