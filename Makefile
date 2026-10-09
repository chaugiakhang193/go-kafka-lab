COMPOSE := docker compose -f deployments/docker-compose.yml
KAFKA_BIN := $(COMPOSE) exec -T kafka /opt/kafka/bin

# VCS stamping makes the go command run git in the checkout. Appending keeps
# any GOFLAGS already set in the environment.
export GOFLAGS += -buildvcs=false

.PHONY: help up down ps logs smoke vet

help: ## List available targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-8s %s\n", $$1, $$2}'

up: ## Start Kafka and the web UI, wait for a healthy broker and a running UI
	$(COMPOSE) up -d --wait

down: ## Stop and remove the containers (topics and offsets are discarded)
	$(COMPOSE) down

ps: ## Show container status
	$(COMPOSE) ps

logs: ## Follow container logs
	$(COMPOSE) logs -f

smoke: ## Print the KRaft quorum status and the topic list
	$(KAFKA_BIN)/kafka-metadata-quorum.sh --bootstrap-server localhost:9092 describe --status
	$(KAFKA_BIN)/kafka-topics.sh --bootstrap-server localhost:9092 --list

vet: ## Check formatting and run go vet
	@test -z "$$(gofmt -l .)" || { gofmt -l .; echo "gofmt: files above need formatting"; exit 1; }
	go vet ./...
