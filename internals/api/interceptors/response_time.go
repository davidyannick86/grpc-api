package interceptors

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ResponseTimeInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Println("ResponseTimeInterceptor: start")
	// Log the start time
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Log the end time
	duration := time.Since(start)

	// Log the response time
	status, _ := status.FromError(err)
	fmt.Printf("Response time for %s: %v, status: %s\n", info.FullMethod, duration, status.Code())

	md := metadata.Pairs(
		"X-Response-Time", duration.String(),
	)

	// Set the metadata in the context
	grpc.SendHeader(ctx, md)

	log.Println("ResponseTimeInterceptor: end")
	// Return the response and error
	return resp, err
}
