---

description: "Task list for Update Project Documentation"
---

# Tasks: Update Project Documentation

**Input**: Design documents from `/specs/001-update-docs/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, quickstart.md ✅

**Organization**: Tasks grouped by user story. All tasks touch `README.md` —
sequential execution required (same file). Each story phase produces an independently
reviewable increment.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (only where files differ)
- **[Story]**: US1, US2, US3 per spec.md

## Path Conventions

- Single file: `README.md` at repository root

---

## Phase 1: Foundational

**Purpose**: Read current README and prepare for rewrite.

- [x] T001 Read `README.md` in full and note all sections to keep, fix, or remove per `specs/001-update-docs/plan.md` diff table

---

## Phase 2: User Story 1 - New User Gets Started (Priority: P1) 🎯 MVP

**Goal**: README gives a new user everything needed to install and generate a project.

**Independent Test**: Follow only the README from install step to `gena new` completion — no errors, project directory created.

### Implementation for User Story 1

- [x] T002 [US1] Write Overview section in `README.md`: one paragraph explaining what Gena does and the problem it solves (combining multiple boilerplate templates into a new project)
- [x] T003 [US1] Write Install section in `README.md`: `go install` command for `gena`; note that `indexer` is for template maintainers
- [x] T004 [US1] Write Quickstart section in `README.md`: end-to-end flow — `gena list`, then `gena new -use <name> -use <name> -out <path>`, then result description

**Checkpoint**: README can now guide a new user through install and generation independently.

---

## Phase 3: User Story 2 - Complete Command Reference (Priority: P2)

**Goal**: Every command and flag for both CLIs is documented with examples in README.

**Independent Test**: Cross-check each documented command and flag against `cmd/indexer/get_args.go` and `cmd/gena/get_args.go` — 0 discrepancies.

### Implementation for User Story 2

- [x] T005 [US2] Write Indexer section in `README.md` with all 6 subcommands: `help`, `list`, `add -name <name> -src <path>`, `rm <name>`, `check`, `version` — each with a shell example
- [x] T006 [US2] Write Gena section in `README.md` with all 4 subcommands: `help`, `list`, `new -use <name> [-use <name>...] -out <path>`, `version` — each with a shell example; note `-use` is repeatable
- [x] T007 [US2] Verify all command names and flag names in `README.md` match source exactly: `rm` (not `remove`), `-name`/`-src` for add, `-use`/`-out` for new

**Checkpoint**: README is now a complete command reference; can be verified against source code.

---

## Phase 4: User Story 3 - Template System Explanation (Priority: P3)

**Goal**: README explains the block-based template system so a developer can create their own template.

**Independent Test**: A reviewer unfamiliar with the project reads the Template System section and can explain `single`, `add`, `merge` block strategies back in their own words.

### Implementation for User Story 3

- [x] T008 [US3] Write Template System section in `README.md`: explain what a "block" is, describe the three merge strategies (`single` — one instance, conflicts error; `add` — content appended; `merge` — content merged), give a plain-language example of when each applies
- [x] T009 [US3] Write Post-Generation Hooks section in `README.md`: list the 4 hardcoded hooks that run after `gena new` (`task generate`, `go mod tidy`, `go fmt ./...`, `goimports -w -local ...`)

**Checkpoint**: README now covers template authoring; US3 independently testable.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Remove internal/stale content, validate full document.

- [x] T010 Remove the "Debug" section from `README.md` (internal push-all commands not relevant to users)
- [x] T011 Review complete `README.md` against `specs/001-update-docs/quickstart.md` — confirm all examples in README are runnable and consistent with quickstart
- [x] T012 Verify README structure order: Overview → Install → Quickstart → Indexer Reference → Gena Reference → Template System → Post-Generation Hooks

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: No dependencies — start immediately
- **US1 (Phase 2)**: Depends on T001 — then start
- **US2 (Phase 3)**: Depends on Phase 2 completion (same file, sequential)
- **US3 (Phase 4)**: Depends on Phase 3 completion (same file, sequential)
- **Polish (Phase 5)**: Depends on all story phases

### Within Each Phase

- All tasks touch `README.md` — execute sequentially within a phase
- T007 (verification) MUST follow T005 and T006

### Parallel Opportunities

None within phases (single-file edits). Phases themselves are sequential for the same reason.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. T001 — read current state
2. T002–T004 — write Overview + Install + Quickstart
3. **STOP and VALIDATE**: Follow the Install + Quickstart in the README; confirm `gena new` works

### Incremental Delivery

1. Phase 1 + Phase 2 → MVP: new user can get started
2. Phase 3 → Command reference complete
3. Phase 4 → Template authoring documented
4. Phase 5 → Final polish

---

## Notes

- All tasks modify only `README.md`; no code changes
- Source of truth for commands: `cmd/indexer/get_args.go`, `cmd/gena/get_args.go`
- Source of truth for block types: `internal/vfs/vfs.go`
- Post-gen hooks source: `cmd/gena/main.go` lines 72-75
