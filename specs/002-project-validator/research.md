# Research: Project Validator CLI

**Feature**: 002-project-validator
**Date**: 2026-04-02

## Decisions

### Decision: Entry point

**Decision**: `task check` in `Taskfile.yml` calls a shell script `scripts/check.sh`
**Rationale**: Taskfile is already the project's task runner. A shell script is
simpler than a new Go binary for a sequential command-with-logging workflow.
Consistent with existing `task lint`, `task build` pattern.
**Alternatives considered**:
- New Go binary `cmd/validator/main.go` — more complex, overkill for a dev tool
  that runs sequential shell commands with logging

---

### Decision: Implementation language

**Decision**: Bash shell script (`scripts/check.sh`)
**Rationale**: All validation stages are subprocess invocations (`go build`, `indexer check`,
`gena new`). Bash handles these natively. No need for Go stdlib complexity.
**Alternatives considered**:
- Go binary — would need `os/exec` for every command; adds a binary to compile before
  it can validate compilation

---

### Decision: 11 validation stages (user-confirmed)

**Pre-flight (4)**:
1. `Build` — `go build ./...`
2. `Lint` — `go fmt ./...` + `go vet ./...`
3. `Template update` — `indexer add` for each template from sibling directories
4. `Template integrity` — `indexer check`

**Generation (7)** — all 2³−1 non-empty combinations of 3 templates:
5. single: `gena_grpc_server`
6. single: `gena_rest_server`
7. single: `gena_kafka_producer`
8. pair: `gena_grpc_server` + `gena_rest_server`
9. pair: `gena_grpc_server` + `gena_kafka_producer`
10. pair: `gena_rest_server` + `gena_kafka_producer`
11. triple: `gena_grpc_server` + `gena_rest_server` + `gena_kafka_producer`

---

### Decision: Stage isolation

**Decision**: Each generation stage uses a unique temp directory
(`/tmp/gena-check-<stage>/`). All temp dirs are cleaned up on exit via `trap`.
**Rationale**: Avoids "output directory already exists" errors. Cleanup happens
even on script interruption.

---

### Decision: Template update failure handling

**Decision**: Template update stage fails gracefully if sibling source directories
do not exist (warns but does not abort remaining stages).
**Rationale**: On CI or a fresh checkout, sibling directories may not be present.
The existing indexed templates in `pkg/indexes/` are still valid for generation testing.

---

### Decision: Stage failure behavior

**Decision**: A failing stage prints `[FAIL]` and its error, but does NOT abort
subsequent stages. Only post-flight summary exits with non-zero code if any stage failed.
Exception: if `Build` fails, all generation stages are skipped (can't run unbuilt binaries).
**Rationale**: User wants to see all failures in one run, not just the first.
Early-abort on build failure is practical — there's nothing to test if it won't compile.

---

### Decision: Output format

```
[PASS] Stage 1: Build
[PASS] Stage 2: Lint
[PASS] Stage 3: Template update
[PASS] Stage 4: Template integrity
[PASS] Stage 5: Generate [gena_grpc_server]
[FAIL] Stage 6: Generate [gena_rest_server]
       Error: exit status 1 — ...
...
---
Summary: 10 passed, 1 failed
OVERALL: FAIL
```

---

### Decision: Overwrite existing `task check`

**Decision**: Replace the existing `task check` in `Taskfile.yml` with a call to
`scripts/check.sh`. The `task update` dependency is removed from `check` (update
is now an explicit stage inside the script).
**Source**: User instruction: "Существующая команда task check должна быть перезаписана"
