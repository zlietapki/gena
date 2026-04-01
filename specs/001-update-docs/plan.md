# Implementation Plan: Update Project Documentation

**Branch**: `001-update-docs` | **Date**: 2026-04-02 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-update-docs/spec.md`

## Summary

Update `README.md` to accurately reflect actual CLI capabilities of `indexer` and `gena`,
fix command name discrepancies (`remove` → `rm`), add a block-system explanation,
add post-generation hook info, and provide an end-to-end quickstart.

## Technical Context

**Language/Version**: Markdown (documentation only)
**Primary Dependencies**: N/A
**Storage**: `README.md` at repository root
**Testing**: Manual review — follow quickstart.md, verify every command example runs without error
**Target Platform**: GitHub README (rendered Markdown)
**Project Type**: Documentation update
**Performance Goals**: N/A
**Constraints**: Scope limited to `README.md`; no new files in repo
**Scale/Scope**: Single file update

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. CLI-First | ✅ Pass | Doc update only; no code changes |
| II. Template Integrity | ✅ Pass | Not applicable |
| III. Block-Based Merging | ✅ Pass | Will document block types correctly |
| IV. Embedded-Binary Templates | ✅ Pass | Not applicable |
| V. Simplicity | ✅ Pass | Single file, no new abstractions |

## Project Structure

### Documentation (this feature)

```text
specs/001-update-docs/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
README.md                # Only file changed
```

**Structure Decision**: Single-file documentation update. No code changes needed.

## Phase 0 Research Summary

See [research.md](research.md) for full findings.

Key decisions:
- Command `rm` (not `remove`) — confirmed from source
- `gena new` requires both `-use` (repeatable) and `-out` flags
- Post-gen hooks are hardcoded: `task generate`, `go mod tidy`, `go fmt ./...`, `goimports`
- Block types: `single`, `add`, `merge` — from `internal/vfs/vfs.go`
- `indexer add` automatically runs `check` after adding

## Phase 1 Design

### README Sections (new structure)

1. **Overview** — one paragraph: what Gena does, the problem it solves
2. **Install** — `go install` for `gena`; note that `indexer` is for template maintainers
3. **Quickstart** — end-to-end: add template → check → generate
4. **Indexer Reference** — all subcommands with flags and examples
5. **Gena Reference** — all subcommands with flags and examples
6. **Template System** — block types explained in plain language
7. **Post-Generation Hooks** — what runs after `gena new`

### Exact Command Inventory (source of truth)

**indexer**:
```
indexer help
indexer list
indexer add -name <name> -src <path>
indexer rm <name>
indexer check
indexer version
```

**gena**:
```
gena help
gena list
gena new -use <name> [-use <name> ...] -out <path>
gena version
```

### What Changes vs Current README

| Issue | Current | Corrected |
|-------|---------|-----------|
| Remove command name | `remove` (in description) | `rm` |
| Missing `-out` flag description | not described | add description |
| Missing block system section | absent | add |
| Missing post-gen hooks | absent | add |
| Debug section (internal) | present | remove |

## Complexity Tracking

No constitution violations. No complexity justification needed.
