# Quickstart: Project Validator

## Prerequisites

- `indexer` and `gena` binaries built: `task build`
- Sibling template directories exist (for Template update stage):
  - `../gena_grpc_server/`
  - `../gena_rest_server/`
  - `../gena_kafka_producer/`

## Run Validator

```shell
task check
```

## Expected Output (all passing)

```
[PASS] Stage  1: Build
[PASS] Stage  2: Lint
[PASS] Stage  3: Template update
[PASS] Stage  4: Template integrity
[PASS] Stage  5: Generate [gena_grpc_server]
[PASS] Stage  6: Generate [gena_rest_server]
[PASS] Stage  7: Generate [gena_kafka_producer]
[PASS] Stage  8: Generate [gena_grpc_server gena_rest_server]
[PASS] Stage  9: Generate [gena_grpc_server gena_kafka_producer]
[PASS] Stage 10: Generate [gena_rest_server gena_kafka_producer]
[PASS] Stage 11: Generate [gena_grpc_server gena_rest_server gena_kafka_producer]
---
Summary: 11 passed, 0 failed
OVERALL: PASS
```

Exit code: 0

## Expected Output (with failures)

```
[PASS] Stage  1: Build
[FAIL] Stage  2: Lint
       go vet: ./internal/vfs/vfs.go:42: error message
[PASS] Stage  3: Template update
...
---
Summary: 10 passed, 1 failed
OVERALL: FAIL
```

Exit code: 1

## When Build Fails

If Stage 1 (Build) fails, generation stages 5–11 are skipped automatically:

```
[FAIL] Stage 1: Build
       ./cmd/gena/main.go:15: syntax error
[PASS] Stage 2: Lint
[PASS] Stage 3: Template update
[PASS] Stage 4: Template integrity
[SKIP] Stage 5–11: Generate (skipped — build failed)
---
Summary: 3 passed, 1 failed, 7 skipped
OVERALL: FAIL
```

## When Sibling Dirs Are Missing (CI / fresh checkout)

Template update stage warns but does not fail:

```
[WARN] Stage 3: Template update
       ../gena_grpc_server/ not found — using existing index
```

The validator continues using whatever templates are already in `pkg/indexes/`.
