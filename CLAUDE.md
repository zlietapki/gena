# gena (Microboiler) Development Guidelines

## Role: Orchestrator

**gena is an orchestrator**, not a standalone application. It manages and merges
boilerplate templates from three independent sibling projects:

| Project | Path | Role |
|---------|------|------|
| `gena_grpc_server` | `../gena_grpc_server/` | gRPC microservice template |
| `gena_rest_server` | `../gena_rest_server/` | REST microservice template (Echo, OpenAPI) |
| `gena_kafka_producer` | `../gena_kafka_producer/` | Kafka producer template (franz-go) |

When making changes to gena, consider the impact on all three templates.
Use `/agent:gena_grpc_server`, `/agent:gena_rest_server`, `/agent:gena_kafka_producer`
to work on individual template projects.

## Active Technologies

- Go 1.25.x (no CGO)
- Bash (shell scripts in `scripts/`)

## Project Structure

```text
cmd/indexer/        # indexer CLI — converts projects to YAML template indexes
cmd/gena/           # gena CLI — generates projects by merging templates
internal/vfs/       # virtual file system (Directory, File, Block types)
internal/fsindex/   # IndexDir() — recursively indexes a directory tree
internal/generator/ # MergeDirs() — merges multiple template directories
internal/indexchecker/ # conflict validation across templates
pkg/indexes/        # embedded YAML template indexes (go:embed)
scripts/check.sh    # project validator (11 stages)
```

## Commands

```shell
task build     # build bin/indexer and bin/gena
task update    # re-index templates from sibling workspace dirs
task check     # run full validation: 11 stages (build, lint, templates, all combinations)
task lint      # go fmt ./... && go vet ./...
```

## Code Style

- Go: stdlib `flag`, `gopkg.in/yaml.v3`, no unnecessary abstractions
- No CGO, no frameworks
- Block merge strategies: `single` (identical), `add` (append), `merge` (deduplicate)

## Recent Changes

- 002-project-validator: Added `scripts/check.sh` — 11-stage validator, rewrote `task check`
- 001-update-docs: Updated README.md with full CLI reference and template system docs

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
