# Redis Cluster Sandbox

This workspace spins up a six-node Redis 7.2 cluster (3 masters + 3 replicas) with Docker Compose. Each node stores its runtime state inside the `data/<port>` folders and reads configuration from a dedicated `redis.conf` file that you can customise.

## Prerequisites

- Docker Engine & Docker Compose (v2+)
- 6 free TCP ports (`7001-7006`) and the corresponding cluster bus ports (`17001-17006`)
- macOS, Linux, or WSL2 (validated on macOS)

## Project Layout

```
redis-cluster/
├── docker-compose.yml          # Defines 6 Redis nodes + init helper
├── README.md                   # This file
└── data/
    ├── 7001/                   # Per-node data volume (bind-mounted)
    │   └── redis.conf          # Node-specific configuration
    ├── 7002/
    │   └── redis.conf
    ├── 7003/
    │   └── redis.conf
    ├── 7004/
    │   └── redis.conf
    ├── 7005/
    │   └── redis.conf
    └── 7006/
        └── redis.conf
```

> Runtime artefacts such as `nodes.conf`, `appendonly.aof`, `appendonlydir`, and `dump.rdb` are ignored via `.gitignore`, so only `redis.conf` files remain version-controlled.

## Usage

### 1. Start the cluster

```bash
# from the repo root
sudo docker compose up -d
```

Compose boots six Redis containers (`redis-7001` … `redis-7006`) and a one-shot helper (`redis-init-cluster`) that creates the cluster topology if it has not been formed yet.

### 2. Verify cluster status

```bash
docker exec redis-7001 redis-cli -p 7001 cluster info
docker exec redis-7001 redis-cli -p 7001 cluster nodes
```

Look for `cluster_state:ok` and 16384 slots assigned across the three master nodes.

### 3. Interact with Redis

```bash
# Write / read a key on any node
docker exec redis-7001 redis-cli -p 7001 SET foo bar
docker exec redis-7003 redis-cli -p 7003 GET foo
```

### 4. Stop the cluster

```bash
sudo docker compose down
```

Use `sudo docker compose down -v` if you want to remove the named `rcnet` network and delete anonymous volumes (not required for the bind-mounted `data/` directories).

## Customising Configuration

Each node reads `/data/redis.conf`, so you can tweak per-node behaviour by editing the corresponding file under `data/<port>/redis.conf`. Typical adjustments include:

```conf
# Enable authentication
requirepass YOUR_PASSWORD
masterauth YOUR_PASSWORD

# Tune memory policy
maxmemory 512mb
maxmemory-policy volatile-lru

# Change persistence strategy
appendonly no
save 900 1
save 300 10
```

After editing a configuration file, restart the affected container to apply the change:

```bash
sudo docker compose restart redis-7004
```

## Rebuilding from Scratch

If you need a clean cluster (for example after major config changes):

```bash
sudo docker compose down
rm -f data/*/nodes.conf data/*/appendonly.aof data/*/dump.rdb
sudo docker compose up -d
```

The init helper will recreate the cluster automatically once all nodes are healthy.

## Troubleshooting

- **`Address already in use`** – ensure no other processes are bound to ports `7001-7006` / `17001-17006` and that old containers/networks are removed (`docker ps -a`, `docker network ls`).
- **Cluster stuck in `fail` state** – run `docker exec redis-7001 redis-cli --cluster fix redis-7001:7001 --cluster-yes` or remove `nodes.conf` files and start fresh.
- **Config changes not applied** – remember to restart the specific container after editing its `redis.conf`.

## Useful Commands

```bash
# Tail logs from a node
docker logs -f redis-7002

# Issue cluster-wide check
docker exec redis-7001 redis-cli --cluster check redis-7001:7001

# Flush data on all nodes (development only)
for port in 7001 7002 7003 7004 7005 7006; do \
  docker exec redis-$port redis-cli -p $port FLUSHALL; \
Done
```

Happy clustering! 🚀

### Run Go Application

```bash
cd app

go run example-host.go

# or

go run test-connection.go
```