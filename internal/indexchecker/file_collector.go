package indexchecker

import (
	"path/filepath"

	"github.com/zlietapki/gena/internal/vfs"
)

type fileCollector struct {
	fileMap fileMapType
}

type fileMapType map[string][]fileEntry

type fileEntry struct {
	projName string
	file     vfs.File
}

func newFileCollector() *fileCollector {
	return &fileCollector{
		fileMap: make(map[string][]fileEntry),
	}
}

func (c *fileCollector) collect(projName string, dir vfs.Directory, prefix string) {
	for _, file := range dir.Files {
		fullPath := filepath.Join(prefix, file.Name)

		c.fileMap[fullPath] = append(c.fileMap[fullPath], fileEntry{
			projName: projName,
			file:     file,
		})
	}

	for _, sub := range dir.Dirs {
		c.collect(projName, sub, filepath.Join(prefix, sub.Name))
	}
}
