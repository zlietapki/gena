package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zlietapki/microboiler/internal/vfs"
	"gopkg.in/yaml.v3"
)

//go:embed microboiler_grpc_server.yml
var folder1 string

//go:embed microboiler_grpc_server.yml
var folder2 string

func main() {
	var d1, d2 vfs.Directory
	if err := yaml.Unmarshal([]byte(folder1), &d1); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal([]byte(folder2), &d2); err != nil {
		panic(err)
	}

	result := mergeDirs(d1, d2)

	//err := createFileSystem(result, "/tmp/some")
	//if err != nil {
	//	panic(err)
	//}

	out, err := yaml.Marshal(result)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(out))
}

func createFileSystem(dir vfs.Directory, path string) error {
	dirPath := filepath.Join(path, dir.Name)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}
	for _, f := range dir.Files {
		var sb strings.Builder
		for _, blk := range f.Blocks {
			sb.WriteString(strings.Join(blk.Data, "\n"))
		}
		if err := os.WriteFile(filepath.Join(dirPath, f.Name), []byte(sb.String()), 0644); err != nil {
			return err
		}
	}
	for _, sub := range dir.Dirs {
		if err := createFileSystem(sub, dirPath); err != nil {
			return err
		}
	}
	return nil
}
