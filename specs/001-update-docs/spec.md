# Feature Specification: Update Project Documentation

**Feature Branch**: `001-update-docs`
**Created**: 2026-04-02
**Status**: Draft
**Input**: User description: "Обнови документацию проекта в соответствии с его возможностями"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - New User Gets Started (Priority: P1)

A developer discovers Gena (Microboiler) and reads README.md to understand what it does,
how to install it, and how to generate their first project from templates. They can
complete a full workflow — index a template, verify it, generate a project — using
only the README as a guide.

**Why this priority**: First contact with the project. If README does not cover the
basic workflow clearly, new users are lost before they start.

**Independent Test**: A person unfamiliar with the project reads the README and
successfully runs `gena new` to generate a project, without looking at any other docs
or source code.

**Acceptance Scenarios**:

1. **Given** a developer has never used Gena, **When** they read the README,
   **Then** they understand the purpose of the tool in under 2 minutes.
2. **Given** the README install section, **When** a developer follows the steps,
   **Then** both `indexer` and `gena` binaries are available on their system.
3. **Given** the README quickstart, **When** a developer follows it end-to-end,
   **Then** a generated project directory appears at the specified output path.

---

### User Story 2 - Developer Looks Up Command Reference (Priority: P2)

A developer using Gena regularly wants to look up exact flag names and parameters
for any command without running `--help` or reading source code.

**Why this priority**: Reduces friction for existing users; a complete command
reference prevents guesswork.

**Independent Test**: For each command listed in both CLIs, the README provides
the command name, all flags with types, and at least one usage example.

**Acceptance Scenarios**:

1. **Given** the README command reference, **When** a developer looks up `indexer add`,
   **Then** they find the `-name` and `-src` flags with descriptions and an example.
2. **Given** the README command reference, **When** a developer looks up `gena new`,
   **Then** they find the `-use` and `-out` flags with descriptions and an example.
3. **Given** the README, **When** a developer looks for any command that exists
   in the tool, **Then** they find it documented — no command is missing.

---

### User Story 3 - Developer Understands the Template System (Priority: P3)

A developer wants to create their own template and needs to understand the
block-based merging concept before indexing their project.

**Why this priority**: Necessary for template authors, but not required for basic
project generation from existing templates.

**Independent Test**: The README explains what a "block" is, what merge strategies
exist, and how the indexer converts a project into a reusable template — without
requiring the reader to look at YAML files.

**Acceptance Scenarios**:

1. **Given** the README, **When** a developer reads the template system section,
   **Then** they understand the difference between `single`, `add`, and `merge`
   block strategies.
2. **Given** the README, **When** a developer follows the indexing workflow,
   **Then** they can index their own project and verify it with `indexer check`.

---

### Edge Cases

- What if a command documented in README does not exist in the current binary?
- What if a flag name in the README differs from the actual flag (e.g., `rm` vs `remove`)?
- What if the reader uses Windows and the install/usage instructions are Unix-only?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: README MUST contain a clear one-paragraph description of what Gena does
  and the problem it solves.
- **FR-002**: README MUST include an Install section with instructions to obtain
  both `indexer` and `gena` binaries.
- **FR-003**: README MUST document every subcommand of `indexer` (help, list, add,
  remove, check, version) with all accepted flags and at least one usage example each.
- **FR-004**: README MUST document every subcommand of `gena` (help, list, new,
  version) with all accepted flags and at least one usage example each.
- **FR-005**: README MUST include an end-to-end quickstart showing the full workflow:
  index a template → check for conflicts → generate a project.
- **FR-006**: README MUST explain the block-based template merging concept
  (block types: single, add, merge) in plain language.
- **FR-007**: All command examples in README MUST match actual flag names and
  behavior of the current version.
- **FR-008**: README MUST NOT document commands or flags that no longer exist.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer unfamiliar with Gena can complete the quickstart workflow
  (install → index → generate) using only the README, without errors, in under 10 minutes.
- **SC-002**: Every command and flag available in both CLIs is represented in the
  README — 0 undocumented commands.
- **SC-003**: 0 examples in the README reference flags or commands that do not
  exist in the current binary.
- **SC-004**: The block system section is understandable without prior knowledge —
  verifiable by a reviewer unfamiliar with the project being able to explain it back
  in their own words.

## Assumptions

- Documentation scope is limited to `README.md`; no additional doc files are in scope.
- The current behavior of `indexer` and `gena` binaries is the source of truth;
  README is updated to match code, not the other way around.
- Windows-specific instructions are out of scope; Unix/macOS shell examples are sufficient.
- The block system explanation is conceptual only — no YAML schema reference is
  required in this iteration.
