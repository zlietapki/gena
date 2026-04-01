---

description: "Task list for Project Validator CLI"
---

# Tasks: Project Validator CLI

**Input**: Design documents from `/specs/002-project-validator/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, quickstart.md ✅

**Organization**: Tasks by user story. Two files changed: `scripts/check.sh` (new)
and `Taskfile.yml` (one task overwritten).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1, US2, US3 per spec.md

## Path Conventions

```text
scripts/check.sh     # new validation script
Taskfile.yml         # modified: overwrite task check
```

---

## Phase 1: Setup

**Purpose**: Create file structure.

- [x] T001 Create `scripts/` directory in repo root (if it does not exist) and create empty `scripts/check.sh` with `#!/usr/bin/env bash` header and `set -euo pipefail` disabled (errors handled per-stage manually)

---

## Phase 2: Foundational

**Purpose**: Core helper functions used by all user stories.

- [x] T002 Add `run_stage()` helper function to `scripts/check.sh`: accepts stage number, label, and command; executes command; prints `[PASS] Stage N: label` or `[FAIL] Stage N: label\n       <error>` to stdout; increments global pass/fail counters
- [x] T003 Add `skip_stage()` helper function to `scripts/check.sh`: accepts stage number and label; prints `[SKIP] Stage N: label`; increments skip counter
- [x] T004 Add `trap` cleanup handler in `scripts/check.sh`: on EXIT removes all `/tmp/gena-check-*/` directories created during the run
- [x] T005 Add summary function to `scripts/check.sh`: prints `---`, `Summary: N passed, N failed, N skipped`, `OVERALL: PASS/FAIL`; exits with code 0 if failed==0, else 1

**Checkpoint**: Helper infrastructure complete. Stories can now be implemented independently.

---

## Phase 3: User Story 1 - Validate After Any Change (Priority: P1) 🎯 MVP

**Goal**: Running `task check` shows [PASS]/[FAIL] for every stage with enough detail
to identify failures without running additional commands.

**Independent Test**: Run `task check` on unmodified project → all stages PASS,
exit 0. Introduce a syntax error in `cmd/gena/main.go` → Stage 1 shows [FAIL] with
error message.

### Implementation for User Story 1

- [x] T006 [US1] Add Stage 1 (Build) to `scripts/check.sh`: `run_stage 1 "Build" "go build ./..."` — on failure, set `BUILD_FAILED=1`
- [x] T007 [US1] Add Stage 2 (Lint) to `scripts/check.sh`: run `go fmt ./... 2>&1` and `go vet ./... 2>&1` captured together; pass combined output to `run_stage 2 "Lint"`
- [x] T008 [US1] Add Stage 3 (Template update) to `scripts/check.sh`: for each sibling dir (`../gena_grpc_server/`, `../gena_rest_server/`, `../gena_kafka_producer/`), if dir exists run `go run ./cmd/indexer/ add -name <name> -src <dir>`; if dir missing print `[WARN] Stage 3: Template update — <dir> not found, using existing index` (does not increment fail counter)
- [x] T009 [US1] Add Stage 4 (Template integrity) to `scripts/check.sh`: `run_stage 4 "Template integrity" "go run ./cmd/indexer/ check"`
- [x] T010 [US1] Update `Taskfile.yml`: replace existing `check` task with `desc: Run full project validation (11 stages)` and single command `bash scripts/check.sh`; remove `deps: [update]`
- [x] T011 [US1] Make `scripts/check.sh` executable: `chmod +x scripts/check.sh`

**Checkpoint**: `task check` runs 4 pre-flight stages with PASS/FAIL output.
US1 independently testable at this point.

---

## Phase 4: User Story 2 - All Template Combinations Tested (Priority: P2)

**Goal**: Every non-empty combination of the 3 indexed templates is tested in a
single run. One combination failing does not prevent others from running.

**Independent Test**: With 3 templates indexed, `task check` shows 7 separate
generation stages (5–11). Removing one template and re-running still shows 7
stages but some fail with "index not found" error.

### Implementation for User Story 2

- [x] T012 [US2] Add Stage 5 (Generate single: gena_grpc_server) to `scripts/check.sh`: if `BUILD_FAILED`, call `skip_stage 5 "Generate [gena_grpc_server]"`; else create `OUT=/tmp/gena-check-5`, `run_stage 5 "Generate [gena_grpc_server]" "go run ./cmd/gena/ new -use gena_grpc_server -out $OUT"`
- [x] T013 [US2] Add Stage 6 (Generate single: gena_rest_server) to `scripts/check.sh`: same pattern with `gena_rest_server` and `/tmp/gena-check-6`
- [x] T014 [US2] Add Stage 7 (Generate single: gena_kafka_producer) to `scripts/check.sh`: same pattern with `gena_kafka_producer` and `/tmp/gena-check-7`
- [x] T015 [US2] Add Stage 8 (Generate pair: grpc+rest) to `scripts/check.sh`: `run_stage 8 "Generate [gena_grpc_server gena_rest_server]" "go run ./cmd/gena/ new -use gena_grpc_server -use gena_rest_server -out /tmp/gena-check-8"`
- [x] T016 [US2] Add Stage 9 (Generate pair: grpc+kafka) to `scripts/check.sh`: same pattern for grpc+kafka, output `/tmp/gena-check-9`
- [x] T017 [US2] Add Stage 10 (Generate pair: rest+kafka) to `scripts/check.sh`: same pattern for rest+kafka, output `/tmp/gena-check-10`
- [x] T018 [US2] Add Stage 11 (Generate triple: all 3) to `scripts/check.sh`: `run_stage 11 "Generate [gena_grpc_server gena_rest_server gena_kafka_producer]" "go run ./cmd/gena/ new -use gena_grpc_server -use gena_rest_server -use gena_kafka_producer -out /tmp/gena-check-11"`

**Checkpoint**: All 7 generation stages present. US2 independently testable.

---

## Phase 5: User Story 3 - Zero Configuration (Priority: P3)

**Goal**: `task check` requires no arguments, no config, no pre-created directories.
Cleans up temp dirs on exit including on interrupt (Ctrl+C).

**Independent Test**: Run `task check`, interrupt mid-run with Ctrl+C — no
`/tmp/gena-check-*/` directories remain afterward.

### Implementation for User Story 3

- [x] T019 [US3] Verify `trap` in `scripts/check.sh` covers `EXIT`, `INT`, `TERM` signals — cleanup removes `/tmp/gena-check-*/` even on Ctrl+C
- [x] T020 [US3] Verify each generation stage pre-checks that its output directory does not exist and removes it before running (handles leftover from crashed previous run)

**Checkpoint**: Zero-config, self-cleaning execution confirmed.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [x] T021 Run `task check` manually end-to-end; verify output matches `specs/002-project-validator/quickstart.md` example format
- [x] T022 Verify exit code: `task check && echo OK` prints OK on clean project; `echo $?` returns 0; after introducing deliberate failure returns 1
- [x] T023 [P] Update `scripts/check.sh` shebang and add brief header comment describing script purpose and exit codes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Phase 1 (file must exist)
- **US1 (Phase 3)**: Depends on Foundational — needs `run_stage`, `skip_stage`, `trap`
- **US2 (Phase 4)**: Depends on US1 — generation stages reference `BUILD_FAILED` set in Stage 1
- **US3 (Phase 5)**: Depends on Phase 1 (trap already in place from T004) — verification tasks only
- **Polish (Phase 6)**: Depends on all story phases

### Within Each Phase

- T012–T018 (generation stages): each touches `scripts/check.sh` sequentially (same file)
- T006–T010 (US1 implementation): sequential (same file + Taskfile.yml)

### Parallel Opportunities

- T006, T010 can run in parallel (different files: `scripts/check.sh` vs `Taskfile.yml`)
- T012–T018 generation stages are logically independent but edit same file — sequential

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. T001–T005 — setup + helpers
2. T006–T011 — 4 pre-flight stages + Taskfile
3. **STOP and VALIDATE**: `task check` runs and reports 4 stages

### Incremental Delivery

1. Phase 1+2+3 → MVP: pre-flight validation works
2. Phase 4 → All 7 generation combinations tested
3. Phase 5 → Cleanup verified
4. Phase 6 → Polish

---

## Notes

- `scripts/check.sh` does NOT use `set -e` globally — each stage handles its own failure
- `BUILD_FAILED` flag prevents running generation stages against a broken binary
- Template update stage uses `go run` (not compiled binary) for consistency with current `task update`
- Temp dirs pattern: `/tmp/gena-check-<N>/` where N is stage number
