package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zlietapki/gena/internal/vfs"
)

func createFileSystem(dir vfs.Directory, path string) error {
	dirPath := filepath.Join(path, dir.Name)
	mode := (fs.FileMode)(dir.Mode)
	if err := os.Mkdir(dirPath, mode); err != nil {
		return err
	}

	for _, file := range dir.Files {
		var sb strings.Builder
		for _, block := range file.Blocks {
			sb.WriteString(strings.Join(block.Data, "\n"))
			sb.WriteString("\n")
		}

		filePath := filepath.Join(dirPath, file.Name)
		mode = (fs.FileMode)(file.Mode)
		err := os.WriteFile(filePath, []byte(sb.String()), mode)
		if err != nil {
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
