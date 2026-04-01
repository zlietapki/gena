# Quickstart: Gena (Microboiler)

## Prerequisites

- Go installed
- `task` (Taskfile runner) installed in the generated project's directory (for post-gen hooks)
- `goimports` installed: `go install golang.org/x/tools/cmd/goimports@latest`

## Step 1: Install

```shell
# Install gena (project generator)
go install github.com/zlietapki/gena/cmd/gena@latest

# Install indexer (template maintainer tool — only needed to manage templates)
go install github.com/zlietapki/gena/cmd/indexer@latest
```

## Step 2: List available templates

```shell
gena list
```

## Step 3: Generate a project

```shell
gena new -use gena_grpc_server -use gena_rest_server -out /tmp/myproject
```

Gena merges the specified templates and writes the project to `/tmp/myproject`.
After writing, it automatically runs: `task generate`, `go mod tidy`, `go fmt ./...`, `goimports`.

## Step 4: Work with your project

```shell
cd /tmp/myproject
task run
```

---

## Managing Templates (for template maintainers)

### Add a template from an existing project

```shell
indexer add -name my_template -src /path/to/my/project
```

This indexes the project and automatically runs `indexer check` to validate
there are no conflicts with existing templates.

### Check for conflicts

```shell
indexer check
```

### List indexed templates

```shell
indexer list
```

### Remove a template

```shell
indexer rm my_template
```

After making template changes, rebuild the binary to embed the updated templates:

```shell
task build
```
