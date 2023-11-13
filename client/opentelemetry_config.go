package main

import (
	"time"

	"google.golang.org/grpc"
)

type OpenTelemetryConfig struct {
	URL         string
	Timeout     time.Duration
	DialOption  []grpc.DialOption
	ServiceName string
}
