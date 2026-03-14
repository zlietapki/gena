package indexchecker

import (
	"github.com/zlietapki/gena/internal/vfs"
	"github.com/zlietapki/gena/pkg/indexes"
)

type blockEntry struct {
	projName string
	block    vfs.Block
}

func SingleBlocksSameContent() error {
	fileMap, err := getFileMap()
	if err != nil {
		return err
	}

	for path, fileEntries := range fileMap {
		if len(fileEntries) < 2 {
			continue
		}

		_ = checkSingleBlocksSameContent(path, fileEntries)
	}

	return nil
}

func BlocksSameType() error {
	fileMap, err := getFileMap()
	if err != nil {
		return err
	}

	for path, fileEntries := range fileMap {
		if len(fileEntries) < 2 {
			continue
		}

		_ = checkBlockSameTypes(path, fileEntries)
	}

	return nil
}

func getFileMap() (fileMapType, error) {
	projs, err := indexes.GetAll()
	if err != nil {
		return nil, err
	}

	fc := newFileCollector()
	for projName, dir := range projs {
		fc.collect(projName, dir, "")
	}

	return fc.fileMap, nil
}

func withNewlines(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = l + "\n"
	}
	return out
}
