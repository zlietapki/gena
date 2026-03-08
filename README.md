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
