# Feature Specification: Project Validator CLI

**Feature Branch**: `002-project-validator`
**Created**: 2026-04-02
**Status**: Draft
**Input**: User description: "Необходима единственная команда CLI, которая будет проверять валидность проекта..."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Validate After Any Change (Priority: P1)

A developer has just made changes to Gena (modified a template, refactored merging
logic, updated the generator). They run a single command and within seconds know
whether the project still works correctly. The output lists every validation stage
with a clear PASS or FAIL label, so they immediately know where things broke.

**Why this priority**: The core need — instant feedback after any change. Without
this, a developer must manually run multiple commands and interpret raw output to
know if the project is healthy.

**Independent Test**: Run the validator command on an unmodified project. All stages
show PASS. Then introduce a known conflict between two templates and re-run — the
validator reports FAIL on the affected stage with a description of the conflict.

**Acceptance Scenarios**:

1. **Given** Gena is in a valid state, **When** the developer runs the validator,
   **Then** every stage prints `[PASS]` and the final summary shows success.
2. **Given** Gena has a broken template combination, **When** the developer runs
   the validator, **Then** the stage that failed prints `[FAIL]` with the error
   message, and remaining stages still execute and report their own result.
3. **Given** the validator has completed, **When** the developer reads the output,
   **Then** they can determine which stage failed without looking at any other
   output or running additional commands.

---

### User Story 2 - All Template Combinations Tested (Priority: P2)

The developer wants confidence that not just the "happy path" template combination
works, but every individual template and every multi-template combination works.
The validator automatically discovers all indexed templates and generates every
possible combination, running `gena new` for each.

**Why this priority**: A change to template A might break a combination with template
B without affecting standalone A. Manual testing of all combinations is error-prone
and slow.

**Independent Test**: Index 3 templates, run the validator. Confirm the output shows
results for each individual template (3 stages) and each multi-template combination
(3 pairs + 1 triple = 4 stages), totaling 7 generation stages tested.

**Acceptance Scenarios**:

1. **Given** N templates are indexed, **When** the validator runs,
   **Then** it tests every non-empty combination (2^N − 1 total).
2. **Given** one combination fails, **When** the validator runs,
   **Then** all other combinations still execute and report independently.
3. **Given** the validator completes, **When** the developer reads the summary,
   **Then** they see the count of passed vs failed combinations.

---

### User Story 3 - Zero Configuration (Priority: P3)

The developer runs the validator with no arguments. It discovers templates
automatically, manages temporary output directories itself (creates and cleans up),
and requires no external configuration files or environment setup beyond having
Gena built.

**Why this priority**: A validation tool that requires setup defeats its purpose.
It must be frictionless to run after every change.

**Independent Test**: On a fresh checkout with templates indexed, run the validator
with no arguments. It completes without requiring any flags, config files, or
pre-created directories.

**Acceptance Scenarios**:

1. **Given** no arguments are passed, **When** the validator runs,
   **Then** it discovers templates and runs all combinations automatically.
2. **Given** the validator ran previously, **When** it is run again,
   **Then** there are no leftover temporary directories from the prior run.
3. **Given** the validator exits with an error, **When** the developer inspects
   the filesystem, **Then** temporary output directories have been cleaned up.

---

### Edge Cases

- What if no templates are indexed? The validator should report this clearly and exit.
- What if a temporary output directory already exists before a run (e.g., previous crash)?
  The validator must handle this gracefully — skip or remove, not crash.
- What if a post-generation hook fails (e.g., `go mod tidy` unavailable)? The stage
  should be marked FAIL with the hook name and error in the log.
- What if there are 0 or 1 templates (no multi-template combinations possible)?
  The validator should still run single-template stages without error.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The validator MUST be invocable as a single command with no required
  arguments.
- **FR-002**: The validator MUST discover all currently indexed templates automatically
  without any configuration.
- **FR-003**: The validator MUST test every non-empty combination of indexed templates
  (all single-template runs, all pairs, all triples, etc.).
- **FR-004**: Each combination run MUST be an independent stage — a failure in one
  stage MUST NOT prevent other stages from executing.
- **FR-005**: Each stage MUST produce a log entry with: stage name (template combination),
  outcome (`PASS` or `FAIL`), and on failure, the error message.
- **FR-006**: The validator MUST print a final summary: total stages, passed count,
  failed count, and overall result.
- **FR-007**: The validator MUST clean up all temporary output directories it creates,
  regardless of whether stages pass or fail.
- **FR-008**: The validator MUST exit with a non-zero exit code if any stage fails,
  and exit code 0 if all stages pass.
- **FR-009**: The validator MUST validate that `indexer check` passes before running
  generation stages; if it fails, generation stages are skipped and reported as blocked.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A developer can determine the overall health of the project by reading
  only the final summary line — no need to scroll through full output for a passing run.
- **SC-002**: When any stage fails, the developer can identify the exact failing
  combination and error cause within 30 seconds of reading the output.
- **SC-003**: The validator covers 100% of possible template combinations (2^N − 1
  for N indexed templates) in a single run.
- **SC-004**: The validator leaves no temporary files or directories on disk after
  completion, regardless of outcome.
- **SC-005**: The validator completes a full run (3 templates, 7 combinations)
  in under 2 minutes on a standard developer machine.

## Assumptions

- The validator is invoked from the Gena project root directory.
- All indexed templates are available and can be used to generate projects.
- The developer's machine has all post-generation hook prerequisites installed
  (`task`, `go`, `goimports`); hook failures are reported as stage failures.
- Temporary output directories are created in the system temp directory
  (or a dedicated subfolder) to avoid polluting the project root.
- The validator tests generation correctness only — it does not verify that the
  generated project compiles or runs correctly (that is out of scope for v1).
