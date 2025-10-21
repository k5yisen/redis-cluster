# Redis Cluster Sandbox

A production-ready **6-node Redis 8.2.2 cluster** (3 masters + 3 replicas) built with Docker Compose. This setup is designed for local development and testing, supporting both **macOS** (with port mapping) and **Linux/Ubuntu** (with host networking).

## 🌟 Features

- ✅ **Redis 8.2.2** cluster with automatic failover
- ✅ **3 master + 3 replica** architecture for high availability
- ✅ **Persistent storage** with bind-mounted volumes
- ✅ **Automatic cluster initialization** on first start
- ✅ **Platform-specific configurations** (macOS and Linux)
- ✅ **Go client examples** for cluster interaction
- ✅ **Automated setup script** for Ubuntu/Linux

## 📋 Prerequisites

- **Docker Engine** & **Docker Compose** (v2+)
- 6 free TCP ports: `7001-7006` (Redis ports)
- 6 free TCP ports: `17001-17006` (cluster bus ports)
- **macOS**, **Linux**, or **WSL2**
- *Optional*: **Go 1.19+** (for running Go examples)

## 📁 Project Structure

```
redis-cluster/
├── docker-compose-macos.yml      # macOS config (port mapping)
├── docker-compose-host.yml       # Linux config (host networking)
├── setup-ubuntu.sh               # Automated Ubuntu setup script
├── README.md                     # This file
├── app/                          # Go client examples
│   ├── example-host.go           # Full cluster example
│   ├── test-connection.go        # Connection test
│   └── go.mod                    # Go dependencies
├── config-ubuntu/                # Ubuntu-specific configs
│   ├── redis-7001.conf
│   ├── redis-7002.conf
│   ├── redis-7003.conf
│   ├── redis-7004.conf
│   ├── redis-7005.conf
│   └── redis-7006.conf
└── data/                         # Persistent data (bind-mounted)
    ├── 7001/
    │   ├── redis.conf            # Node configuration
    │   ├── nodes.conf            # Cluster topology (auto-generated)
    │   ├── dump.rdb              # RDB snapshot
    │   └── appendonlydir/        # AOF files
    ├── 7002/
    ├── 7003/
    ├── 7004/
    ├── 7005/
    └── 7006/
```

> **Note**: Runtime artifacts (`nodes.conf`, `dump.rdb`, `appendonlydir`) are git-ignored, only `redis.conf` files are version-controlled.

---

## 🚀 Quick Start

### macOS Setup

**Using Docker Compose:**

```bash
task setup-macos
```

### Ubuntu/Linux Setup

**Option 1: Automated Script** *(Recommended)*

```bash
# Make the script executable
chmod +x setup-ubuntu.sh

# Run the setup script
./setup-ubuntu.sh
```

The script will:
1. ✅ Copy Ubuntu-specific configs to `data/` directories
2. ✅ Clean old cluster state
3. ✅ Stop any existing cluster
4. ✅ Start the cluster with host networking
5. ✅ Wait for cluster formation
6. ✅ Display cluster status

**Option 2: Manual Setup**

```bash
# Copy configs
for port in 7001 7002 7003 7004 7005 7006; do
    cp config-ubuntu/redis-$port.conf data/$port/redis.conf
done

# Clean old cluster state
rm -f data/*/nodes.conf

# Start the cluster
docker compose -f docker-compose-host.yml up -d

# Verify cluster status
docker exec redis-7001 redis-cli -p 7001 cluster info
docker exec redis-7001 redis-cli -p 7001 cluster nodes
```

---

## 🔍 Verifying the Cluster

### Check Cluster Status

```bash
docker exec redis-7001 redis-cli -p 7001 cluster info
```

Expected output should include:
```
cluster_state:ok
cluster_slots_assigned:16384
cluster_known_nodes:6
```

### View Cluster Topology

```bash
docker exec redis-7001 redis-cli -p 7001 cluster nodes
```

You should see 3 masters and 3 replicas with slots distributed across masters (0-16383).

### Test Cluster Operations

```bash
# Set a key
docker exec redis-7001 redis-cli -p 7001 SET foo bar

# Get from another node (cluster will redirect)
docker exec redis-7003 redis-cli -p 7003 GET foo

# Test counter
docker exec redis-7001 redis-cli -p 7001 INCR counter

# Test hash
docker exec redis-7001 redis-cli -p 7001 HSET user:1 name "Alice" email "alice@example.com"
docker exec redis-7001 redis-cli -p 7001 HGETALL user:1

# Test list
docker exec redis-7001 redis-cli -p 7001 RPUSH queue:tasks "task1" "task2" "task3"
docker exec redis-7001 redis-cli -p 7001 LRANGE queue:tasks 0 -1
```

---

## 💻 Using Go Client Examples

The `app/` directory contains Go examples for connecting to the cluster from your host machine.

### Prerequisites

```bash
cd app
go mod download
```

### Run the Examples

**Full cluster example:**

```bash
cd app
go run example-host.go
```

This will:
- Connect to the Redis cluster via `localhost:7001-7006`
- Perform SET/GET operations
- Test counters, hashes, and lists
- Display cluster info and nodes

**Connection test:**

```bash
cd app
go run test-connection.go
```

### Key Points

- **macOS**: Examples use `localhost:7001-7006` (port mapping)
- **Linux**: Examples use `localhost:7001-7006` (host networking)
- The Go client automatically handles cluster redirections
- Uses `github.com/redis/go-redis/v9` cluster client

---

## 🔧 Configuration

### Per-Node Configuration

Each node reads its configuration from `data/<port>/redis.conf`. Common customizations:

```conf
# Security (add authentication)
requirepass YOUR_PASSWORD
masterauth YOUR_PASSWORD

# Memory Management
maxmemory 512mb
maxmemory-policy volatile-lru

# Persistence Strategy
appendonly yes
appendfsync everysec

# RDB Snapshots
save 900 1     # After 900 sec if at least 1 key changed
save 300 10    # After 300 sec if at least 10 keys changed
save 60 10000  # After 60 sec if at least 10000 keys changed

# Logging
loglevel notice
logfile ""
```

**Apply changes:**

```bash
# Restart a specific node
docker restart redis-7001

# Or restart all nodes
docker compose -f docker-compose-macos.yml restart  # macOS
docker compose -f docker-compose-host.yml restart   # Linux
```

### Network Configuration Differences

**macOS (`docker-compose-macos.yml`):**
- Uses **port mapping** (`-p 7001:7001`)
- Containers communicate via Docker's bridge network
- Host connects via `localhost:7001-7006`
- Nodes announce themselves as `host.docker.internal`

**Linux (`docker-compose-host.yml`):**
- Uses **host networking** (`network_mode: host`)
- Containers directly use host's network stack
- More efficient for cluster communication
- Nodes announce themselves as `127.0.0.1`

---

## 🛠️ Maintenance & Operations

### Viewing Logs

```bash
# View logs from a specific node
docker logs redis-7001

# Follow logs in real-time
docker logs -f redis-7001

# View cluster initialization logs
docker logs redis-init-cluster
```

### Cluster Health Check

```bash
# Check cluster state
docker exec redis-7001 redis-cli --cluster check 127.0.0.1:7001

# Fix cluster issues (if any)
docker exec redis-7001 redis-cli --cluster fix 127.0.0.1:7001 --cluster-yes

# Rebalance cluster slots
docker exec redis-7001 redis-cli --cluster rebalance 127.0.0.1:7001
```

### Monitoring

```bash
# Monitor real-time operations on a node
docker exec redis-7001 redis-cli -p 7001 MONITOR

# Get cluster statistics
docker exec redis-7001 redis-cli -p 7001 INFO stats

# Check memory usage
docker exec redis-7001 redis-cli -p 7001 INFO memory

# View connected clients
docker exec redis-7001 redis-cli -p 7001 CLIENT LIST
```

### Backup & Restore

**Backup:**

```bash
# Trigger RDB snapshot on all nodes
for port in 7001 7002 7003 7004 7005 7006; do
    docker exec redis-$port redis-cli -p $port BGSAVE
done

# Backup data directory
tar -czf redis-backup-$(date +%Y%m%d-%H%M%S).tar.gz data/
```

**Restore:**

```bash
# Stop cluster
docker compose -f docker-compose-macos.yml down

# Restore data
tar -xzf redis-backup-YYYYMMDD-HHMMSS.tar.gz

# Start cluster
docker compose -f docker-compose-macos.yml up -d
```

---

## 🔄 Rebuilding from Scratch

If you need to completely reset the cluster:

```bash
# macOS
docker compose -f docker-compose-macos.yml down
rm -f data/*/nodes.conf data/*/dump.rdb
rm -rf data/*/appendonlydir
docker compose -f docker-compose-macos.yml up -d

# Linux
docker compose -f docker-compose-host.yml down
rm -f data/*/nodes.conf data/*/dump.rdb
rm -rf data/*/appendonlydir
docker compose -f docker-compose-host.yml up -d
```

The `init-cluster` service will automatically recreate the cluster topology.

---

## 🐛 Troubleshooting

### Common Issues

**1. Address already in use**

```bash
# Check what's using the ports
lsof -i :7001

# Stop conflicting services or remove old containers
docker ps -a | grep redis
docker rm -f $(docker ps -a -q --filter "name=redis-*")
```

**2. Cluster stuck in `fail` state**

```bash
# Fix cluster
docker exec redis-7001 redis-cli --cluster fix 127.0.0.1:7001 --cluster-yes

# If that doesn't work, rebuild from scratch (see above)
```

**3. Configuration changes not applied**

```bash
# Restart the specific container
docker restart redis-7001

# Verify config was loaded
docker exec redis-7001 redis-cli -p 7001 CONFIG GET maxmemory
```

**4. Node can't join cluster**

```bash
# Check cluster state
docker exec redis-7001 redis-cli -p 7001 cluster nodes

# Reset a failing node
docker exec redis-7001 redis-cli -p 7001 CLUSTER RESET HARD
```

**5. macOS: Go client can't connect**

- Make sure you're using `docker-compose-macos.yml`
- Verify ports are exposed: `docker ps`
- Try connecting to `localhost:7001` not `127.0.0.1:7001`
- Check Go client uses `ClusterClient`, not single-node client

**6. Linux: Permission denied**

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Re-login or use
newgrp docker
```

---

## 📚 Useful Commands

### Cluster Management

```bash
# View cluster information
docker exec redis-7001 redis-cli -p 7001 cluster info

# View cluster nodes and slot distribution
docker exec redis-7001 redis-cli -p 7001 cluster nodes

# Check cluster consistency
docker exec redis-7001 redis-cli --cluster check 127.0.0.1:7001

# Rebalance slots across masters
docker exec redis-7001 redis-cli --cluster rebalance 127.0.0.1:7001
```

### Data Operations

```bash
# Flush all data (development only!)
for port in 7001 7002 7003 7004 7005 7006; do
    docker exec redis-$port redis-cli -p $port FLUSHALL
done

# Count keys on each node
for port in 7001 7002 7003 7004 7005 7006; do
    echo "Port $port:"
    docker exec redis-$port redis-cli -p $port DBSIZE
done

# Get cluster key distribution
docker exec redis-7001 redis-cli --cluster call 127.0.0.1:7001 DBSIZE
```

### Performance Testing

```bash
# Benchmark cluster performance
docker exec redis-7001 redis-cli --cluster call 127.0.0.1:7001 DEBUG SLEEP 0

# Run redis-benchmark against cluster
redis-benchmark -p 7001 -c 50 -n 10000 --cluster

# Memory profiling
docker exec redis-7001 redis-cli -p 7001 MEMORY DOCTOR
```

---

## 🔗 Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Redis Cluster                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Master 1 (7001)  ←──→  Replica 1 (7004)                │
│  Slots: 0-5460                                          │
│                                                         │
│  Master 2 (7002)  ←──→  Replica 2 (7005)                │
│  Slots: 5461-10922                                      │
│                                                         │
│  Master 3 (7003)  ←──→  Replica 3 (7006)                │
│  Slots: 10923-16383                                     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

- **Hash slots**: 16384 total, distributed across 3 masters
- **Replication**: Each master has 1 replica
- **Automatic failover**: If a master fails, its replica is promoted
- **Cluster bus**: Nodes communicate on ports 17001-17006

---

Happy clustering! 🚀
