package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zlietapki/microboiler/internal/vfs"
)

const dirMode = 0755
const fileMode = 0644

func createFileSystem(dir vfs.Directory, path string) error {
	dirPath := filepath.Join(path, dir.Name)
	if err := os.Mkdir(dirPath, dirMode); err != nil {
		return err
	}

	for _, f := range dir.Files {
		var sb strings.Builder
		for _, blk := range f.Blocks {
			sb.WriteString(strings.Join(blk.Data, "\n"))
			sb.WriteString("\n")
		}
		if err := os.WriteFile(filepath.Join(dirPath, f.Name), []byte(sb.String()), fileMode); err != nil {
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
