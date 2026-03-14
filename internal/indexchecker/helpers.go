package indexchecker

import (
	"github.com/zlietapki/gena/internal/vfs"
	"github.com/zlietapki/gena/pkg/indexes"
)

type blockEntry struct {
	projName string
	block    vfs.Block
}

func getCollector() (*fileCollector, error) {
	projs, err := indexes.GetAll()
	if err != nil {
		return nil, err
	}

	fc := newFileCollector()
	for projName, dir := range projs {
		fc.collect(projName, dir, "")
	}

	return fc, nil
}

func getFileMap() (fileMapType, error) {
	fc, err := getCollector()
	if err != nil {
		return nil, err
	}
	return fc.fileMap, nil
}

func getDirMap() (dirMapType, error) {
	fc, err := getCollector()
	if err != nil {
		return nil, err
	}
	return fc.dirMap, nil
}

func withNewlines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = l + "\n"
	}
	return out
}
