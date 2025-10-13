#!/bin/bash

# Ubuntu/Linux Redis Cluster Setup Script
# This script copies Ubuntu-specific config files and starts the cluster

set -e

echo "🚀 Ubuntu Redis Cluster Setup"
echo "================================"
echo ""

# Check if running on Linux
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo "⚠️  Warning: This script is designed for Ubuntu/Linux"
    echo "   You are running on: $OSTYPE"
    echo "   For macOS, use: docker-compose -f docker-compose-macos.yml up -d"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Step 1: Copy config files
echo "📁 Step 1: Copying Ubuntu config files..."
for port in 7001 7002 7003 7004 7005 7006; do
    echo "   Copying config-ubuntu/redis-$port.conf → data/$port/redis.conf"
    cp config-ubuntu/redis-$port.conf data/$port/redis.conf
done
echo "✅ Config files copied"
echo ""

# Step 2: Clean old cluster state
echo "🧹 Step 2: Cleaning old cluster state..."
rm -f data/*/nodes.conf
echo "✅ Old cluster state removed"
echo ""

# Step 3: Stop existing cluster (if any)
echo "🛑 Step 3: Stopping existing cluster..."
docker compose -f docker-compose-host.yml down 2>/dev/null || true
echo "✅ Cluster stopped"
echo ""

# Step 4: Start cluster
echo "🚀 Step 4: Starting Redis Cluster with host networking..."
docker compose -f docker-compose-host.yml up -d
echo "✅ Cluster started"
echo ""

# Step 5: Wait for cluster formation
echo "⏳ Step 5: Waiting for cluster to form (this takes ~15-20 seconds)..."
sleep 5
echo "   Checking cluster status..."

# Wait up to 30 seconds for cluster to be ready
for i in {1..30}; do
    if docker exec redis-7001 redis-cli -p 7001 cluster info 2>/dev/null | grep -q "cluster_state:ok"; then
        echo "✅ Cluster is ready!"
        break
    fi
    echo "   Waiting... ($i/30)"
    sleep 1
done

# Check final status
echo ""
echo "📊 Cluster Status:"
echo "=================="
docker exec redis-7001 redis-cli -p 7001 cluster info | grep -E "cluster_state|cluster_slots|cluster_known_nodes"

echo ""
echo "🎉 Setup Complete!"
echo "=================="
echo ""
echo "Cluster Nodes:"
docker exec redis-7001 redis-cli -p 7001 cluster nodes | grep -E "master|slave" | head -6

echo ""
echo "📝 Next Steps:"
echo "   1. Test connection: redis-cli -p 7001 cluster info"
echo "   2. Run Go example: cd app && go run example-host.go"
echo "   3. View logs: docker logs redis-init-cluster"
echo ""
echo "✅ Your Redis Cluster is ready at localhost:7001-7006"
