go run ./cmd/fs_generator/ check/folder1 pkg/genfs/folder1.go go
go run ./cmd/fs_generator/ check/folder2 pkg/genfs/folder2.go go
go run ./cmd/merger_go/

go run ./cmd/fs_generator/ check/folder1 cmd/merger_yml/folder1.yml yml
go run ./cmd/fs_generator/ check/folder2 cmd/merger_yml/folder2.yml yml

go run ./cmd/merger_yml/

--real
go run ./cmd/fs_generator/ ../microboiler_grpc_server/ cmd/merger_yml/microboiler_grpc_server.yml yml

go run ./cmd/fs_generator/ check/folder3/ cmd/merger_yml/microboiler_grpc_server.yml yml

# create files int /tmp/some
go run ./cmd/merger_yml/