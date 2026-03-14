Microboiler
===========

Install
-------

```shell
go install github.com/zlietapki/gena/cmd/gena@latest
```

Update project templates
------------------------

```shell
go run ./cmd/indexer/ -name microboiler_grpc_server -src ../microboiler_grpc_server/
go run ./cmd/indexer/ -name microboiler_rest_server -src ../microboiler_rest_server/
```

start
-----

```shell
go run ./cmd/gena/
```

Debug
-----

```shell
task check

cd /tmp/some/check
task run
```