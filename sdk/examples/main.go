package main

import (
	"context"
	"fmt"
	"github.com/Royal17x/flagr/sdk"
	"log"
	"time"
)

func main() {
	client, err := sdk.NewClient(
		"http://localhost:8080",
		"sdk-key-demo",
		sdk.WithTimeout(2*time.Second),
		sdk.WithCacheTTL(30*time.Second),
		sdk.WithDefaultValue(false),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	projectID := "77c00606-0099-4642-83e4-0d03c6f78c36"
	envID := "44e0b5cf-3190-41a8-892f-c407af78eb65"
	ctx := context.Background()

	start := time.Now()
	enabled := client.IsEnabled(ctx, "checkout-v2", projectID, envID)
	fmt.Printf("First call (network): enabled=%v, latency=%v\n", enabled, time.Since(start))

	start = time.Now()
	enabled = client.IsEnabled(ctx, "checkout-v2", projectID, envID)
	fmt.Printf("Second call (cache):  enabled=%v, latency=%v\n", enabled, time.Since(start))

	start = time.Now()
	enabled = client.IsEnabled(ctx, "non-existent-flag", projectID, envID)
	fmt.Printf("Missing flag:         enabled=%v, latency=%v\n", enabled, time.Since(start))

	grpcClient, err := sdk.NewClient(
		"localhost:50051",
		"sdk-key-demo",
		sdk.WithGRPC(),
		sdk.WithTLS(false),
		sdk.WithTimeout(2*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClient.Close()

	start = time.Now()
	enabled = grpcClient.IsEnabled(ctx, "checkout-v2", projectID, envID)
	fmt.Printf("gRPC call:            enabled=%v, latency=%v\n", enabled, time.Since(start))
}
