# Research: Update Project Documentation

**Feature**: 001-update-docs
**Date**: 2026-04-02

## Findings

### Decision: Exact command names

**Decision**: `indexer rm` (not `remove`)
**Rationale**: `cmd/indexer/get_args.go` line 42: `case "rm":`; usage string confirms `indexer rm index_name`
**Source**: `cmd/indexer/get_args.go`

---

### Decision: gena new flags

**Decision**: `-use` (repeatable, `arrayFlags`) and `-out` (string)
**Rationale**: `cmd/gena/get_args.go` lines 46-49; both flags required, `-use` can be repeated multiple times
**Source**: `cmd/gena/get_args.go`

---

### Decision: Block types

**Decision**: Three types: `single`, `add`, `merge`
- `single` — block content must be identical across all templates; conflict = error
- `add` — content from each template is appended
- `merge` — content is merged (e.g., import lists, dependency lists)
**Source**: `internal/vfs/vfs.go` lines 13-15

---

### Decision: Post-generation hooks

**Decision**: Four hardcoded commands run after `gena new`:
1. `task generate`
2. `go mod tidy`
3. `go fmt ./...`
4. `goimports -w -local github.com/zlietapki/boilerplate .`

**Rationale**: Hardcoded per Constitution Principle V (Simplicity); no config needed
**Source**: `cmd/gena/main.go` lines 72-75

---

### Decision: indexer add auto-check

**Decision**: `indexer add` automatically runs `check` after indexing
**Rationale**: `cmd/indexer/main.go` line 72: `check()` called inside `addIndex()`
**Source**: `cmd/indexer/main.go`

---

### Decision: Index storage location

**Decision**: Templates stored as YAML in `pkg/indexes/` (compiled into binary via `go:embed`)
**Rationale**: `cmd/indexer/main.go` constant `indexOutput = "pkg/indexes"` (line 18)
**Source**: `cmd/indexer/main.go`

---

### Decision: README scope

**Decision**: Only `README.md` at repo root; no new files
**Rationale**: Spec assumption + Constitution Principle V (Simplicity)

---

### Decision: Remove Debug section

**Decision**: Remove the "Debug" section from README (internal push-all commands)
**Rationale**: Not useful for end users; internal developer workflow; irrelevant to public docs
