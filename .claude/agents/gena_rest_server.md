---
name: gena_rest_server
description: Use this agent to work on the gena_rest_server template project — a REST microservice boilerplate with Echo and OpenAPI. Use when modifying REST handlers, domain interfaces, OpenAPI contracts, usecase logic, or anything specific to the REST server template.
---

# gena_rest_server — REST Microservice Template

**Working directory**: `/home/asd/workspace/gena_rest_server/`

This project is a **template** used by gena (Microboiler) to generate REST microservices.
Changes here must be re-indexed into gena: run `task update` from `/home/asd/workspace/gena/`.

## Architecture

Clean Architecture with four layers:
- **domain** (`internal/domain/`) — interfaces (`IUsecase`, `ICounterRepo`), entities; no deps
- **usecase** (`internal/usecase/`) — business logic; implements `IUsecase`
- **rest_handler** (`internal/rest_handler/`) — HTTP transport; implements generated server interface
- **map_repo** (`internal/map_repo/`) — in-memory repository implementation

## Tech Stack

- Go 1.26.1
- `github.com/labstack/echo/v5` — HTTP framework
- `github.com/oapi-codegen/oapi-codegen/v2` — OpenAPI → Go code generation
- `github.com/joho/godotenv` + `github.com/kelseyhightower/envconfig` — config from env
- `github.com/gojuno/minimock/v3` — mock generation

## API Contract

- OpenAPI spec: `api/` directory
- Generated models: `internal/rest_models/` (do NOT edit manually — run `task generate`)
- Endpoint: `GET /api/counter` — returns and increments a counter value

## Commands

```shell
task run       # go run ./cmd/server
task test      # go test ./...
task lint      # go fmt ./... && go vet ./...
task generate  # regenerate mocks (minimock) and REST models (oapi-codegen)
task tools     # install minimock and oapi-codegen to bin/
task build     # build binary to bin/server
```

## Key Patterns

- `//go:generate minimock -i IUsecase -o ./mocks/` — mock generation
- `//go:generate oapi-codegen ...` — REST model generation from OpenAPI spec
- `var _ domain.IUsecase = (*Usecase)(nil)` — compile-time interface check
- `Depends` struct for constructor dependency injection
- Template merge markers: `// start name:X type:merge` / `// end name:X`
- Config: `HTTP_LISTEN` env var (default `:8080`)
- Server: Echo with logging and recovery middleware in `pkg/http_echo_server/`

## Template Marker Rules

Do NOT remove `// start name:X` / `// end name:X` markers. When adding new imports
or dependencies: use the appropriate `add` or `merge` block type markers.

## API Changes

When changing the OpenAPI spec in `api/`:
1. Update the `.yaml` file
2. `task generate` — regenerates `internal/rest_models/`
3. Update handler in `internal/rest_handler/` to match new interface

## After Any Change

1. `task lint` — verify formatting
2. `task generate` — regenerate if interfaces or OpenAPI spec changed
3. `task test` — verify tests pass
4. `cd /home/asd/workspace/gena && task check` — verify no regressions across all template combinations
