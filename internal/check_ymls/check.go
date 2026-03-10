package check_ymls

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/zlietapki/microboiler/internal/vfs"
	"github.com/zlietapki/microboiler/pkg/projects"
)

type fileEntry struct {
	projName string
	file     vfs.File
}

type fileCollector struct {
	fileMap map[string][]fileEntry
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

type blockEntry struct {
	projName string
	block    vfs.Block
}

func CheckProjects() error {
	projects, err := projects.GetAll()
	if err != nil {
		return err
	}

	fc := newFileCollector()
	for projName, dir := range projects {
		fc.collect(projName, dir, "")
	}

	for path, fileEntries := range fc.fileMap {
		if len(fileEntries) < 2 {
			continue
		}

		_ = checkSingleBlocksSameContent(path, fileEntries)
		_ = checkBlockSameTypes(path, fileEntries)
	}

	return nil
}

func checkBlockSameTypes(path string, fileEntries []fileEntry) bool {
	blockMap := map[string][]blockEntry{}

	for _, fileEnt := range fileEntries {
		for _, block := range fileEnt.file.Blocks {
			blockMap[block.Name] = append(blockMap[block.Name], blockEntry{
				projName: fileEnt.projName,
				block:    block,
			})
		}
	}

	ok := true
	for blockName, blockEntries := range blockMap {
		if len(blockEntries) < 2 {
			continue
		}

		ref := blockEntries[0].block.Type
		for _, be := range blockEntries[1:] {
			if ref != be.block.Type {
				fmt.Printf("conflict: file=%q block=%q block type mismatch in projects: %s vs %s",
					path, blockName, blockEntries[0].projName, be.projName)
				ok = false
			}
		}
	}

	return ok
}

func checkSingleBlocksSameContent(path string, fileEntries []fileEntry) bool {
	blockMap := map[string][]blockEntry{}

	for _, fileEnt := range fileEntries {
		for _, block := range fileEnt.file.Blocks {
			if block.Type == vfs.BlockTypeSingle {
				blockMap[block.Name] = append(blockMap[block.Name], blockEntry{
					projName: fileEnt.projName,
					block:    block,
				})
			}
		}
	}

	ok := true
	for blockName, blockEntries := range blockMap {
		if len(blockEntries) < 2 {
			continue
		}

		firstBlockData := blockEntries[0].block.Data
		for _, be := range blockEntries[1:] {
			if !reflect.DeepEqual(firstBlockData, be.block.Data) {
				fmt.Printf("Conflict blocks type:single. Different content: file=%q block=%q in projects %s vs %s",
					path, blockName, blockEntries[0].projName, be.projName)
				ok = false
			}
		}
	}

	return ok
}
