<!--
SYNC IMPACT REPORT
==================
Version change: none → 1.0.0 (initial ratification)
Modified principles: N/A (initial)
Added sections: Core Principles, Technology Stack, Development Workflow, Governance
Removed sections: N/A
Templates requiring updates:
  ✅ .specify/templates/plan-template.md — Constitution Check gates align with principles below
  ✅ .specify/templates/spec-template.md — no constitution-specific mandatory sections needed
  ✅ .specify/templates/tasks-template.md — no principle-driven task types added beyond defaults
Follow-up TODOs: none
-->

# Gena (Microboiler) Constitution

## Core Principles

### I. CLI-First

All functionality MUST be exposed exclusively via CLI tools (`indexer`, `gena`).
Interaction protocol: arguments/flags → stdout; errors → stderr.
No interactive TUI required for primary operations; interactive prompts (via `huh`)
are permitted only for optional UX enhancements, never for core workflows.
New features MUST be accessible non-interactively (scriptable).

### II. Template Integrity Before Generation

Templates MUST pass conflict validation (`indexer check`) before any generation attempt.
Conflicts — duplicate blocks with incompatible merge strategies, mismatched file modes,
or contradictory single-block content — MUST be reported with actionable diagnostics.
Generation MUST be aborted on unresolved conflicts; partial output is not acceptable.

### III. Block-Based Merging

Files are decomposed into named blocks with explicit merge strategies:
- `single` — block appears exactly once; conflicts are errors.
- `add` — block content is appended from all templates.
- `merge` — block content is merged (e.g., imports, lists).

Every block's merge strategy MUST be declared in the YAML index.
Undeclared merge behavior is forbidden; ambiguity MUST be resolved at indexing time.

### IV. Embedded-Binary Templates

Template indexes MUST be compiled into the binary via `go:embed`.
YAML is the sole interchange format for template indexes.
No runtime file system dependency for reading templates; all template data is self-contained.
Adding a template requires re-indexing and re-building the binary — this is intentional.

### V. Simplicity

YAGNI: implement only what current use cases require.
Prefer stdlib over third-party libraries. Add dependencies only when stdlib is insufficient.
No organizational-only packages, abstractions for a single use site, or speculative APIs.
Post-generation hooks (e.g., `go mod tidy`, `go fmt`) MUST be hardcoded, not configurable,
until a concrete multi-template need for configurability arises.

## Technology Stack

- **Language**: Go (current: 1.25.x); no CGO.
- **CLI parsing**: stdlib `flag` package.
- **YAML**: `gopkg.in/yaml.v3`.
- **Interactive prompts**: `charm.land/huh/v2` (optional UX only).
- **MIME detection**: `github.com/gabriel-vasile/mimetype` (indexer only).
- **Build & tasks**: Taskfile (`task` runner).
- **Lint**: `go fmt` + `go vet` (mandatory before commit).

New dependencies require explicit justification; stdlib equivalents MUST be ruled out first.

## Development Workflow

- `task update` — re-index all templates from sibling workspace directories.
- `task build` — produces `bin/indexer` and `bin/gena`.
- `task check` — full round-trip: update → generate test project → verify no errors.
- `task lint` — `go fmt ./...` + `go vet ./...`; MUST pass before any commit.
- Generated files (`*.gen.go`) are excluded from version control (`.gitignore`).
- Commit messages: imperative, concise, no period at end.

## Governance

This constitution supersedes all other documented practices.
Amendments require: (1) update this file with version bump, (2) update Sync Impact Report
comment, (3) propagate to affected templates.

Version policy:
- MAJOR: removal or incompatible redefinition of a principle.
- MINOR: new principle or section added.
- PATCH: clarifications, wording fixes.

All spec/plan/task reviews MUST verify compliance with the five Core Principles.
Complexity deviating from Principle V MUST be justified in the plan's Complexity Tracking table.

**Version**: 1.0.0 | **Ratified**: 2026-04-02 | **Last Amended**: 2026-04-02
