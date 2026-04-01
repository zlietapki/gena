#!/usr/bin/env bash
# scripts/check.sh — Project validator: runs 11 validation stages
# Exit 0: all stages passed (warnings allowed)
# Exit 1: one or more stages failed

PASS=0
FAIL=0
SKIP=0
BUILD_FAILED=0

# --- Cleanup -----------------------------------------------------------

cleanup() {
    rm -rf /tmp/gena-check-*/
}
trap cleanup EXIT INT TERM

# --- Helpers -----------------------------------------------------------

run_stage() {
    local num="$1"
    local label="$2"
    local cmd="$3"
    local output
    if output=$(bash -c "$cmd" 2>&1); then
        printf "[PASS] Stage %2d: %s\n" "$num" "$label"
        PASS=$((PASS + 1))
        return 0
    else
        printf "[FAIL] Stage %2d: %s\n" "$num" "$label"
        while IFS= read -r line; do
            printf "       %s\n" "$line"
        done <<< "$output"
        FAIL=$((FAIL + 1))
        return 1
    fi
}

skip_stage() {
    local num="$1"
    local label="$2"
    printf "[SKIP] Stage %2d: %s\n" "$num" "$label"
    SKIP=$((SKIP + 1))
}

print_summary() {
    printf -- "---\n"
    printf "Summary: %d passed, %d failed, %d skipped\n" "$PASS" "$FAIL" "$SKIP"
    if [ "$FAIL" -eq 0 ]; then
        printf "OVERALL: PASS\n"
        exit 0
    else
        printf "OVERALL: FAIL\n"
        exit 1
    fi
}

# --- Stage 1: Build ----------------------------------------------------

if ! run_stage 1 "Build" "go build ./..."; then
    BUILD_FAILED=1
fi

# --- Stage 2: Lint -----------------------------------------------------

run_stage 2 "Lint" "go fmt ./... && go vet ./..."

# --- Stage 3: Template update ------------------------------------------
# Warn-only: missing sibling dirs do not count as failure.

stage3_fail=0
stage3_any=0
for template in gena_grpc_server gena_rest_server gena_kafka_producer; do
    src="../${template}/"
    if [ -d "$src" ]; then
        stage3_any=1
        if ! go run ./cmd/indexer/ add -name "$template" -src "$src" > /dev/null 2>&1; then
            stage3_fail=1
        fi
    fi
done

if [ "$stage3_fail" -eq 1 ]; then
    printf "[FAIL] Stage  3: Template update\n"
    FAIL=$((FAIL + 1))
elif [ "$stage3_any" -eq 0 ]; then
    printf "[WARN] Stage  3: Template update — sibling dirs not found, using existing indexes\n"
else
    printf "[PASS] Stage  3: Template update\n"
    PASS=$((PASS + 1))
fi

# --- Stage 4: Template integrity ---------------------------------------

run_stage 4 "Template integrity" "go run ./cmd/indexer/ check"

# --- Stages 5–11: Generation -------------------------------------------
# If build failed, skip all generation stages.

gen_stage() {
    local num="$1"
    local label="$2"
    shift 2
    local use_flags=""
    for t in "$@"; do
        use_flags="$use_flags -use $t"
    done
    local out="/tmp/gena-check-${num}"
    # Remove leftover dir from previous crashed run
    rm -rf "$out"
    run_stage "$num" "$label" "go run ./cmd/gena/ new$use_flags -out $out"
}

if [ "$BUILD_FAILED" -eq 1 ]; then
    for n in 5 6 7 8 9 10 11; do
        skip_stage "$n" "Generate (build failed)"
    done
else
    gen_stage  5 "Generate [gena_grpc_server]"                                     gena_grpc_server
    gen_stage  6 "Generate [gena_rest_server]"                                     gena_rest_server
    gen_stage  7 "Generate [gena_kafka_producer]"                                  gena_kafka_producer
    gen_stage  8 "Generate [gena_grpc_server gena_rest_server]"                    gena_grpc_server gena_rest_server
    gen_stage  9 "Generate [gena_grpc_server gena_kafka_producer]"                 gena_grpc_server gena_kafka_producer
    gen_stage 10 "Generate [gena_rest_server gena_kafka_producer]"                 gena_rest_server gena_kafka_producer
    gen_stage 11 "Generate [gena_grpc_server gena_rest_server gena_kafka_producer]" gena_grpc_server gena_rest_server gena_kafka_producer
fi

# --- Summary -----------------------------------------------------------

print_summary
