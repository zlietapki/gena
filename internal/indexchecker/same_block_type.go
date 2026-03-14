package indexchecker

import (
	"fmt"
)

func BlocksSameType() error {
	fileMap, err := getFileMap()
	if err != nil {
		return err
	}

	for path, fileEntries := range fileMap {
		if len(fileEntries) < 2 {
			continue
		}

		checkBlockSameTypes(path, fileEntries)
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
				fmt.Printf(`Block type mismatch:
	file=%q block=%q
	%s
	%s
`,
					path, blockName, blockEntries[0].projName, be.projName)
				ok = false
			}
		}
	}

	return ok
}
