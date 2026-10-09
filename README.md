# go-kafka-lab

A learning project with hands-on labs for Apache Kafka and clients written in Go. Each lab starts with a prediction,
runs against a local broker and records what actually happened. Everything runs locally, and every topic, key and
payload is synthetic.

The planned labs cover the log and offsets, keys and per-partition ordering, consumer groups and rebalancing, offset
commits and delivery semantics, replay, producer durability (acks, replication, idempotence), idempotent consumers
with the transactional outbox pattern, and retry and dead-letter topics. Each lab compares the behaviour with RabbitMQ
queues and Redis Streams where the two models differ.

## Status

| Component | Status |
|---|---|
| Local stack: single-node Kafka 4.3.1 in KRaft mode and a web UI | done |
| Lab: keys and partitions | done |
| Other labs | planned |

## Stack

- Apache Kafka 4.3.1 (`apache/kafka` image) in KRaft combined mode, with one node acting as both broker and controller
- Go clients built on [franz-go](https://github.com/twmb/franz-go)
- Web UI: [kafbat/kafka-ui](https://github.com/kafbat/kafka-ui)
- Docker Compose and GNU Make, run from a Linux shell such as WSL2

## Run the local stack

```sh
make up      # start Kafka and the web UI, wait for a healthy broker and a running UI
make smoke   # print the KRaft quorum status and the topic list
make down    # remove the containers together with all topics and offsets
```

| Endpoint | Address |
|---|---|
| Bootstrap server for clients on the host | `localhost:9092` |
| Bootstrap server inside the compose network | `kafka:19092` |
| Web UI | http://localhost:8085 |

Published ports bind to `127.0.0.1`. The broker accepts PLAINTEXT connections without authentication and is meant
for local use only.

## Labs

Each lab is a Go program under `labs/`. Run it from the repository root while the local stack is up.

| Lab | Command | Shows |
|---|---|---|
| `labs/keypartition` | `go run ./labs/keypartition` | The partition and offset of each keyed and unkeyed record under the default partitioner |

The key partition lab produces to `lab.products`, which must exist with three partitions:

```sh
docker compose -f deployments/docker-compose.yml exec -T kafka /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 --create --topic lab.products --partitions 3 --replication-factor 1
```

## References

- [Apache Kafka documentation](https://kafka.apache.org/documentation/)
- [apache/kafka Docker image usage](https://github.com/apache/kafka/blob/4.3.1/docker/examples/README.md)
- [franz-go documentation](https://pkg.go.dev/github.com/twmb/franz-go/pkg/kgo)
