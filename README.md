Gena
====

Project generator

Install
-------

```shell
go install github.com/zlietapki/gena/cmd/gena@latest
```

Indexer
-------

Управление шаблонами проектов

```shell
go run ./cmd/indexer/ [command] [-params args]
```

Commands:

* help - помощь
* list - показать список шаблонов
* add - добавить шаблон
  * -name index_name -src template_project_path
* rm index_name - удалить шаблон
* check - проверить конфлиликты загруженных шаблонов
* version - версия indexer

```shell
go run ./cmd/indexer/ help
go run ./cmd/indexer/ list
go run ./cmd/indexer/ add -name gena_grpc_server -src ../gena_grpc_server/
go run ./cmd/indexer/ add -name gena_rest_server -src ../gena_rest_server/
go run ./cmd/indexer/ rm index_name
go run ./cmd/indexer/ check
go run ./cmd/indexer/ version
```

Gena
-----

Генерация проектов

Commands:

* help - помощь
* version - версия gena
* list - показать список шаблонов
* new - генерация проекта
  * -use index_name - использовать шаблон
  * -out - папка назначения

```shell
go run ./cmd/gena/ help
go run ./cmd/gena/ list
go run ./cmd/gena/ new -use gena_grpc_server -use gena_rest_server -out /tmp/gena/project
go run ./cmd/gena/ version
```

Debug
-----

```shell
task check

cd /tmp/some/check
task run

# push all
cd ~/workspace/gena/ && git add . && git commit -m'changes' && git push
cd ~/workspace/gena_kafka_producer/ && git add . && git commit -m'changes' && git push
cd ~/workspace/gena_grpc_server/ && git add . && git commit -m'changes' && git push
cd ~/workspace/gena_rest_server/ && git add . && git commit -m'changes' && git push
```
