go run ./cmd/fs_generator/ check/folder1 pkg/genfs/folder1.go
go run ./cmd/fs_generator/ check/folder2 pkg/genfs/folder2.go

go run ./cmd/fs_generator/ check/folder1 gen_yml/folder1.yml
go run ./cmd/fs_generator/ check/folder2 gen_yml/folder2.yml

go run ./cmd/fs_merger/