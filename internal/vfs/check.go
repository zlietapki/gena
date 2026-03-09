package vfs

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
)

type fileEntry struct {
	ymlPath string
	file    File
}

type blockEntry struct {
	ymlPath string
	block   Block
}

func CheckAllFS(path string) error {
	ymlPaths, err := filepath.Glob(filepath.Join(path, "*.yml"))
	if err != nil {
		return err
	}

	fileMap := map[string][]fileEntry{}
	for _, ymlPath := range ymlPaths {
		data, err := os.ReadFile(ymlPath)
		if err != nil {
			return err
		}

		var dir Directory
		if err := yaml.Unmarshal(data, &dir); err != nil {
			return err
		}

		collectFiles(dir, "", ymlPath, fileMap)
	}

	for fullPath, entries := range fileMap {
		if len(entries) < 2 {
			continue
		}
		if err := checkSingleBlocks(fullPath, entries); err != nil {
			return err
		}
	}

	return nil
}

func collectFiles(dir Directory, prefix string, ymlPath string, fileMap map[string][]fileEntry) {
	for _, f := range dir.Files {
		fullPath := filepath.Join(prefix, f.Name)

		fileMap[fullPath] = append(fileMap[fullPath], fileEntry{
			ymlPath: ymlPath,
			file:    f,
		})
	}

	for _, sub := range dir.Dirs {
		collectFiles(sub, filepath.Join(prefix, sub.Name), ymlPath, fileMap)
	}
}

func checkSingleBlocks(fullPath string, fileEntries []fileEntry) error {
	blockMap := map[string][]blockEntry{}

	for _, e := range fileEntries {
		for _, b := range e.file.Blocks {
			if b.Type == BlockTypeSingle {
				blockMap[b.Name] = append(blockMap[b.Name], blockEntry{e.ymlPath, b})
			}
		}
	}

	for blockName, blockEntries := range blockMap {
		if len(blockEntries) < 2 {
			continue
		}

		ref := blockEntries[0].block.Data
		for _, be := range blockEntries[1:] {
			if !reflect.DeepEqual(ref, be.block.Data) {
				return fmt.Errorf("conflict: file=%q block=%q: %s vs %s",
					fullPath, blockName, blockEntries[0].ymlPath, be.ymlPath)
			}
		}
	}

	return nil
}
