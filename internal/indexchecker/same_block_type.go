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

		err = checkBlockSameTypes(path, fileEntries)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkBlockSameTypes(path string, fileEntries []fileEntry) error {
	blockMap := map[string][]blockEntry{}

	for _, fileEnt := range fileEntries {
		for _, block := range fileEnt.file.Blocks {
			blockMap[block.Name] = append(blockMap[block.Name], blockEntry{
				projName: fileEnt.projName,
				block:    block,
			})
		}
	}

	for blockName, blockEntries := range blockMap {
		if len(blockEntries) < 2 {
			continue
		}

		ref := blockEntries[0].block.Type
		for _, be := range blockEntries[1:] {
			if ref != be.block.Type {
				return fmt.Errorf(`Block type mismatch:
	file=%q block=%q
	%s
	%s
`,
					path, blockName, blockEntries[0].projName, be.projName)
			}
		}
	}

	return nil
}
