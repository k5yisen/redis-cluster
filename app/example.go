package main

// This example is designed to run on the HOST machine (not inside Docker)
// It connects to Redis Cluster using localhost:port since docker-compose-host.yml uses host network mode

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	fmt.Println("=== Connecting to Redis Cluster (from Host using localhost) ===\n")

	// Connect to Redis Cluster using localhost
	opt := &redis.ClusterOptions{
		Addrs: []string{
			"localhost:7001",
			"localhost:7002",
			"localhost:7003",
			"localhost:7004",
			"localhost:7005",
			"localhost:7006",
		},
	}

	rdb := redis.NewClusterClient(opt)
	defer rdb.Close()

	// Ping to verify connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis cluster: %v", err)
	}
	fmt.Printf("✓ Connected to Redis cluster: %s\n\n", pong)

	// SET a key
	err = rdb.Set(ctx, "host:example", "Hello from Go on Host!", 0).Err()
	if err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}
	fmt.Println("✓ SET host:example = Hello from Go on Host!")

	// GET the key
	val, err := rdb.Get(ctx, "host:example").Result()
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	fmt.Printf("✓ GET host:example = %s\n\n", val)

	// Multiple operations
	fmt.Println("=== Running multiple operations ===")

	// Counters
	for i := 0; i < 5; i++ {
		count, _ := rdb.Incr(ctx, "counter:visits").Result()
		fmt.Printf("✓ Visit counter: %d\n", count)
	}

	// Hash
	fmt.Println("\n=== Hash operations ===")
	rdb.HSet(ctx, "user:100", map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
		"score": "95",
	})
	user, _ := rdb.HGetAll(ctx, "user:100").Result()
	fmt.Printf("✓ User data: %+v\n", user)

	// List
	fmt.Println("\n=== List operations ===")
	rdb.RPush(ctx, "queue:tasks", "Process payment", "Send email", "Update database")
	tasks, _ := rdb.LRange(ctx, "queue:tasks", 0, -1).Result()
	fmt.Printf("✓ Task queue: %v\n", tasks)

	// Cluster info
	fmt.Println("\n=== Cluster Info ===")
	clusterInfo, _ := rdb.ClusterInfo(ctx).Result()
	fmt.Println(clusterInfo)

	// Show cluster nodes
	fmt.Println("\n=== Cluster Nodes ===")
	clusterNodes, _ := rdb.ClusterNodes(ctx).Result()
	fmt.Println(clusterNodes)

	fmt.Println("\n🎉 All operations completed successfully!")
	fmt.Println("✓ Go client connected to Redis Cluster directly from host using localhost!")
}
