package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Testing Redis Cluster Connection from Host ===\n")

	// Skip single-node test as it's not cluster-aware
	// Use ClusterClient instead for Redis Cluster

	fmt.Println("Testing cluster connection...")

	// Connect to Redis Cluster - use localhost for macOS
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7001",
			"localhost:7002",
			"localhost:7003",
		},
	})
	defer clusterClient.Close()

	pong, err := clusterClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}
	fmt.Printf("✓ Cluster connected! Response: %s\n", pong)

	// Try cluster operations
	err = clusterClient.Set(ctx, "cluster:test", "Cluster works!", 0).Err()
	if err != nil {
		log.Fatalf("Failed to SET in cluster: %v", err)
	}
	fmt.Println("✓ SET cluster:test = Cluster works!")

	val, err := clusterClient.Get(ctx, "cluster:test").Result()
	if err != nil {
		log.Fatalf("Failed to GET from cluster: %v", err)
	}
	fmt.Printf("✓ GET cluster:test = %s\n", val)

	fmt.Println("\n🎉 Cluster connection works from host!")
}
