Microboiler
===========

Install
-------

```shell
go install github.com/zlietapki/microboiler/cmd/microboiler@latest
```

Update project templates
------------------------

```shell
go run ./cmd/fs_generator/ -name microboiler_grpc_server -src ../microboiler_grpc_server/
go run ./cmd/fs_generator/ -name microboiler_rest_server -src ../microboiler_rest_server/
```

start
-----

```shell
go run ./cmd/microboiler/
```

Debug
-----

```shell
task check

cd /tmp/some/check
go mod tidy
task run
```