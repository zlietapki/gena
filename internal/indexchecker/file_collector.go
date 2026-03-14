package indexchecker

import (
	"path/filepath"

	"github.com/zlietapki/gena/internal/vfs"
)

type fileCollector struct {
	fileMap fileMapType
	dirMap  dirMapType
}

type fileMapType map[string][]fileEntry

type dirMapType map[string][]dirEntry

type fileEntry struct {
	projName string
	file     vfs.File
}

type dirEntry struct {
	projName string
	dir      vfs.Directory
}

func newFileCollector() *fileCollector {
	return &fileCollector{
		fileMap: make(map[string][]fileEntry),
		dirMap:  make(map[string][]dirEntry),
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
		subPath := filepath.Join(prefix, sub.Name)
		c.dirMap[subPath] = append(c.dirMap[subPath], dirEntry{
			projName: projName,
			dir:      sub,
		})
		
		c.collect(projName, sub, subPath)
	}
}
