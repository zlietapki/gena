Microboiler (Gena)
==================

Gena is a project boilerplate generator that merges multiple template projects into
a single new project. Instead of maintaining one monolithic template, you compose
smaller, focused templates — gRPC server, REST server, Kafka producer — and combine
them on demand.

Install
-------

```shell
# Project generator (end users)
go install github.com/zlietapki/gena/cmd/gena@latest

# Template manager (template maintainers only)
go install github.com/zlietapki/gena/cmd/indexer@latest
```

Quickstart
----------

```shell
# See available templates
gena list

# Generate a project from one or more templates
gena new -use gena_grpc_server -use gena_rest_server -out /tmp/myproject

# Work with the generated project
cd /tmp/myproject
task run
```

After `gena new`, the following commands run automatically in the output directory:
`task generate`, `go mod tidy`, `go fmt ./...`, `goimports`.

Indexer Reference
-----------------

The `indexer` tool manages the template library. Use it when you want to add,
update, or validate templates.

**Commands**:

```
indexer help
indexer list
indexer add  -name <name> -src <path>
indexer rm   <name>
indexer check
indexer version
```

**`indexer list`** — print all available template names.

```shell
indexer list
```

**`indexer add`** — index a project directory and save it as a named template.
Runs `indexer check` automatically after adding.

```shell
indexer add -name gena_grpc_server -src ../gena_grpc_server/
indexer add -name gena_rest_server -src ../gena_rest_server/
```

Flags:
- `-name` — template name used later with `gena new -use`
- `-src`  — path to the source project directory

**`indexer rm`** — remove a template by name.

```shell
indexer rm gena_grpc_server
```

**`indexer check`** — validate all templates for conflicts (duplicate single blocks,
mismatched merge strategies, inconsistent file modes).

```shell
indexer check
```

**`indexer version`** — print the indexer version.

```shell
indexer version
```

Gena Reference
--------------

The `gena` tool generates new projects by merging templates.

**Commands**:

```
gena help
gena list
gena new  -use <name> [-use <name> ...] -out <path>
gena version
```

**`gena list`** — print all available template names.

```shell
gena list
```

**`gena new`** — generate a project by merging one or more templates into a
directory. The output directory must not already exist.

```shell
gena new -use gena_grpc_server -out /tmp/myproject
gena new -use gena_grpc_server -use gena_rest_server -use gena_kafka_producer -out /tmp/myproject
```

Flags:
- `-use` — template name to include; repeat for each template
- `-out` — output directory path (created by the command)

**`gena version`** — print the gena version.

```shell
gena version
```

Template System
---------------

Templates are stored as YAML indexes inside the binary (compiled via `go:embed`).
Each file in a template is split into named **blocks**. When multiple templates are
merged, blocks with the same name are combined using one of three strategies:

**`single`** — the block must be identical across all templates that define it.
If two templates define the same single block with different content, `indexer check`
reports a conflict and generation is aborted. Use this for things like license headers
or module declarations that must have one canonical value.

**`add`** — content from each template is appended in order. Use this for things
like Makefile targets or README sections where each template contributes its own
independent entries.

**`merge`** — content is merged (deduplicated lines, sorted). Use this for import
lists, dependency declarations, or any content where order does not matter and
duplicates are noise.

Adding a template requires re-indexing (`indexer add`) and rebuilding the binary
(`task build`) to embed the updated templates.

Post-Generation Hooks
---------------------

After `gena new` writes all files, the following commands run automatically inside
the output directory:

1. `task generate` — runs code generation (e.g., protobuf, mocks)
2. `go mod tidy` — resolves Go module dependencies
3. `go fmt ./...` — formats all Go source files
4. `goimports -w -local github.com/zlietapki/boilerplate .` — organizes imports
