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
go run ./cmd/fs_generator/ ../microboiler_grpc_server/ cmd/microboiler/microboiler_grpc_server.yml yml
go run ./cmd/fs_generator/ ../microboiler_rest_server/ cmd/microboiler/microboiler_rest_server.yml yml
```

start
-----

```shell
go run ./cmd/microboiler/
```

Debug
-----

```shell
rm -r /tmp/some/check

go run ./cmd/fs_generator/ -name microboiler_grpc_server -src ../microboiler_grpc_server/
go run ./cmd/fs_generator/ -name microboiler_rest_server -src ../microboiler_rest_server/
go run ./cmd/microboiler/ --skip-options-select

cd /tmp/some/check
go mod tidy
task run
```