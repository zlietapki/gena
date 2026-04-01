---
name: gena_grpc_server
description: Use this agent to work on the gena_grpc_server template project — a gRPC microservice boilerplate. Use when adding/modifying gRPC handlers, domain interfaces, usecase logic, mock generation, or anything specific to the grpc server template.
---

# gena_grpc_server — gRPC Microservice Template

**Working directory**: `/home/asd/workspace/gena_grpc_server/`

This project is a **template** used by gena (Microboiler) to generate gRPC microservices.
Changes here must be re-indexed into gena: run `task update` from `/home/asd/workspace/gena/`.

## Architecture

Clean Architecture with three layers:
- **domain** (`internal/domain/`) — interfaces (`IUsecase`) and error types; no dependencies
- **usecase** (`internal/usecase/`) — business logic; implements `IUsecase`
- **grpc_handler** (`internal/grpc_handler/`) — gRPC transport; converts proto ↔ domain

## Tech Stack

- Go 1.26.1
- `google.golang.org/grpc` — gRPC framework
- `github.com/zlietapki/microboiler_api_contracts` — protobuf contracts (external repo)
- `github.com/joho/godotenv` + `github.com/kelseyhightower/envconfig` — config from env
- `github.com/gojuno/minimock/v3` — mock generation (`task generate`)

## Commands

```shell
task run       # go run ./cmd/server
task test      # go test ./...
task lint      # go fmt ./... && go vet ./...
task generate  # regenerate mocks via minimock (run after interface changes)
task tools     # install dev tool binaries to bin/
task build     # build binary to bin/server
```

## Key Patterns

- `//go:generate minimock -i IUsecase -o ./mocks/` — mock generation marker
- `var _ domain.IUsecase = (*Usecase)(nil)` — compile-time interface check
- `Depends` struct for constructor dependency injection
- Template merge markers: `// start name:X type:merge` / `// end name:X`
- Config: struct with `envconfig` tags, loaded from `.env` file
- gRPC listen default: `:8888` (`GRPC_LISTEN` env var)

## Template Marker Rules

Do NOT remove `// start name:X` / `// end name:X` markers — they are used by the
gena block-based merging system. Adding new mergeable content: wrap in appropriate
`start/end` markers with the correct block type (`single`, `add`, or `merge`).

## After Any Change

1. `task lint` — verify formatting
2. `task generate` — regenerate mocks if interfaces changed
3. `task test` — verify tests pass
4. `cd /home/asd/workspace/gena && task check` — verify no regressions across all template combinations
