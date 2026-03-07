package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/zlietapki/microboiler/internal/vfs"
)

func getDir(path string) (*vfs.Directory, error) {
	dir := vfs.Directory{
		Name: filepath.Base(path),
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			sub, err := getDir(filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}

			dir.Dirs = append(dir.Dirs, *sub)
		} else {
			file, err := getFile(filepath.Join(path, entry.Name()))
			if err != nil {
				return nil, err
			}

			dir.Files = append(dir.Files, *file)
		}
	}

	return &dir, nil
}

func getFile(path string) (*vfs.File, error) {
	blocks, err := getBlocks(path)
	if err != nil {
		return nil, err
	}

	return &vfs.File{
		Name:   filepath.Base(path),
		Blocks: blocks,
	}, nil
}

func getBlocks(path string) ([]vfs.Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	blocks := make([]vfs.Block, 0)

	var inBlock bool
	var blockName string
	var blockType vfs.BlockType
	var buf []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "// start ") {
			inBlock = true
			buf = nil

			blockName = getBlockName(line)
			blockType = getBlockType(line)
		} else if line == "// end" && inBlock {
			blocks = append(blocks, vfs.Block{
				Name: blockName,
				Type: blockType,
				Data: buf,
			})

			inBlock = false
		} else if inBlock {
			buf = append(buf, line)
		}
	}

	return blocks, scanner.Err()
}

func getBlockName(line string) string {
	for _, tok := range strings.Fields(line) {
		if strings.HasPrefix(tok, "name:") {
			return strings.TrimPrefix(tok, "name:")
		}
	}

	panic("No block name in line")
}

func getBlockType(line string) vfs.BlockType {
	blockType := ""
	for _, tok := range strings.Fields(line) {
		if strings.HasPrefix(tok, "type:") {
			blockType = strings.TrimPrefix(tok, "type:")
		}
	}

	switch blockType {
	case "merge":
		return vfs.BlockTypeMerge
	case "add":
		return vfs.BlockTypeAdd
	default:
		return vfs.BlockTypeOverwrite
	}
}
