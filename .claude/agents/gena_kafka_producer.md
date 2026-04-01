---
name: gena_kafka_producer
description: Use this agent to work on the gena_kafka_producer template project — a Kafka producer microservice boilerplate using franz-go. Use when modifying Kafka client logic, event publishing, domain interfaces, usecase logic, or anything specific to the Kafka producer template.
---

# gena_kafka_producer — Kafka Producer Template

**Working directory**: `/home/asd/workspace/gena_kafka_producer/`

This project is a **template** used by gena (Microboiler) to generate Kafka producer
microservices. Changes here must be re-indexed into gena: run `task update` from
`/home/asd/workspace/gena/`.

## Architecture

Clean Architecture with three layers:
- **domain** (`internal/domain/`) — interfaces (`IEventPublisher`, `IUsecase`), `Event` model
- **usecase** (`internal/usecase/`) — business logic; uses `IEventPublisher` to publish events
- **kafka** (`internal/kafka/`) — Kafka client and producer; implements `IEventPublisher`

## Tech Stack

- Go 1.26.1
- `github.com/twmb/franz-go/pkg/kgo` (v1.20.7) — Kafka client library
- `github.com/joho/godotenv` + `github.com/kelseyhightower/envconfig` — config from env
- `github.com/gojuno/minimock/v3` — mock generation
- **RedPanda** (docker-compose) — local Kafka broker for development

## Infrastructure (local dev)

```shell
# Start RedPanda broker
docker compose up -d

# RedPanda Console UI: http://localhost:8081
# Kafka broker: localhost:9092
```

Config via env vars / `.env`:
- `KAFKA_BROKERS` — comma-separated broker addresses (default: `localhost:9092`)
- `KAFKA_TOPIC` — topic name

## Commands

```shell
task run       # go run ./cmd/server
task test      # go test ./...
task lint      # go fmt ./... && go vet ./...
task generate  # regenerate mocks via minimock
task tools     # install minimock to bin/
```

## Key Patterns

- `IEventPublisher` interface: `Publish(ctx, event) error` — decouples domain from Kafka
- `IUsecase` interface: used for testability via minimock
- `var _ domain.IUsecase = (*Usecase)(nil)` — compile-time interface check
- `Depends` struct for constructor dependency injection
- Async publish with `sync.WaitGroup` in Kafka producer
- Template merge markers: `// start name:X type:merge` / `// end name:X`

## Event Model

`domain.Event` — the core data structure published to Kafka.
Serialized to JSON before publishing. Modify `Event` fields in `internal/domain/`
when the event schema changes.

## Template Marker Rules

Do NOT remove `// start name:X` / `// end name:X` markers. When adding new imports
or Kafka configuration fields: use appropriate block type markers (`add` or `merge`).

## After Any Change

1. `task lint` — verify formatting
2. `task generate` — regenerate mocks if interfaces changed
3. `task test` — verify tests pass
4. `cd /home/asd/workspace/gena && task check` — verify no regressions across all template combinations
