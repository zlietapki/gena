# Implementation Plan: Project Validator CLI

**Branch**: `002-project-validator` | **Date**: 2026-04-02 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-project-validator/spec.md`

## Summary

Add `scripts/check.sh` — a bash script that runs 11 validation stages with structured
`[PASS]`/`[FAIL]` output and a final summary. Overwrite `task check` in `Taskfile.yml`
to call this script. No new Go binary needed.

## Technical Context

**Language/Version**: Bash (shell script); project is Go 1.25.x
**Primary Dependencies**: `go build`, `go fmt`, `go vet`, `indexer`, `gena` (already in project)
**Storage**: N/A (temp dirs in `/tmp/`, cleaned up via `trap`)
**Testing**: Manual — run `task check` on clean project; expect all 11 stages PASS
**Target Platform**: Linux/macOS (developer workstation)
**Project Type**: Developer tooling (validation script)
**Performance Goals**: Complete 11 stages in under 2 minutes
**Constraints**: No new dependencies; no config files required; cleanup on exit/interrupt
**Scale/Scope**: 3 templates → 11 stages (hardcoded; not auto-generated)

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. CLI-First | ✅ Pass | `task check` is the CLI entry point; script is dev tooling |
| II. Template Integrity | ✅ Pass | Stage 4 (`indexer check`) is a dedicated validation stage |
| III. Block-Based Merging | ✅ Pass | Generation stages exercise the merge path |
| IV. Embedded-Binary Templates | ✅ Pass | Script uses compiled binaries from `bin/` |
| V. Simplicity | ✅ Pass | Bash script, no new Go code, no config files |

## Project Structure

### Documentation (this feature)

```text
specs/002-project-validator/
├── plan.md         # This file
├── research.md     # Phase 0 output
├── quickstart.md   # Phase 1 output
└── tasks.md        # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
scripts/
└── check.sh        # New: validation script (11 stages)

Taskfile.yml        # Modified: overwrite task check
```

**Structure Decision**: Single new file (`scripts/check.sh`) + one task modification.

## Phase 0 Research Summary

See [research.md](research.md).

Key decisions:
- Bash script (not Go binary) — simpler for sequential subprocess + logging
- 11 stages total (user-confirmed): 4 pre-flight + 7 generation
- Failing stage prints error but does NOT abort remaining stages
- Exception: `Build` failure → generation stages skipped
- Temp dirs: `/tmp/gena-check-<N>/` cleaned via `trap`
- `task update` dependency removed from `task check`; template update is now an explicit stage

## Phase 1 Design

### `scripts/check.sh` — Stage Structure

```
Stage  1: Build           — go build ./...
Stage  2: Lint            — go fmt ./... && go vet ./...
Stage  3: Template update — indexer add for each sibling dir (graceful on missing)
Stage  4: Template check  — indexer check
Stage  5: Generate [gena_grpc_server]
Stage  6: Generate [gena_rest_server]
Stage  7: Generate [gena_kafka_producer]
Stage  8: Generate [gena_grpc_server gena_rest_server]
Stage  9: Generate [gena_grpc_server gena_kafka_producer]
Stage 10: Generate [gena_rest_server gena_kafka_producer]
Stage 11: Generate [gena_grpc_server gena_rest_server gena_kafka_producer]
```

### Output Format

```
[PASS] Stage 1:  Build
[PASS] Stage 2:  Lint
[WARN] Stage 3:  Template update (sibling dir not found: ../gena_grpc_server)
[PASS] Stage 4:  Template integrity
[PASS] Stage 5:  Generate [gena_grpc_server]
[FAIL] Stage 6:  Generate [gena_rest_server]
       go run: exit status 1: ...error message...
...
---
Summary: 9 passed, 1 failed, 1 warned
OVERALL: FAIL
```

Exit code 0 if all stages pass (warnings allowed); non-zero if any stage fails.

### `Taskfile.yml` change

Replace current `check` task:

```yaml
# Before:
check:
  deps:
    - update
  cmds:
    - rm -rf /tmp/some/check
    - go run ./cmd/gena/ new -use gena_grpc_server -use gena_rest_server -use gena_kafka_producer -out /tmp/some/

# After:
check:
  desc: Run full project validation (11 stages)
  cmds:
    - bash scripts/check.sh
```

### Contracts

`scripts/check.sh` interface contract:
- Exit 0 — all stages passed (warnings OK)
- Exit 1 — one or more stages failed
- stdout — human-readable stage log + summary
- No required arguments; no config files

## Complexity Tracking

No constitution violations.
