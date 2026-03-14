package indexchecker

import (
	"fmt"
	"reflect"

	"github.com/zlietapki/gena/internal/difflib"
	"github.com/zlietapki/gena/internal/vfs"
)

func SingleBlocksSameContent() error {
	fileMap, err := getFileMap()
	if err != nil {
		return err
	}

	for path, fileEntries := range fileMap {
		if len(fileEntries) < 2 {
			continue
		}

		checkSingleBlocksSameContent(path, fileEntries)
	}

	return nil
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
				fmt.Printf(`Different content in blocks type:single
	file=%q block=%q
	%s
	%s
`,
					path, blockName, blockEntries[0].projName, be.projName)
				diffStr, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
					A:        withNewlines(firstBlockData),
					B:        withNewlines(be.block.Data),
					FromFile: blockEntries[0].projName,
					ToFile:   be.projName,
					Context:  3,
				})
				fmt.Print(diffStr)
				ok = false
			}
		}
	}

	return ok
}
