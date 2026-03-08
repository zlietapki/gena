package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zlietapki/microboiler/internal/vfs"

	"github.com/gabriel-vasile/mimetype"
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
		filePath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			sub, err := getDir(filePath)
			if err != nil {
				return nil, err
			}

			dir.Dirs = append(dir.Dirs, *sub)
		} else {
			if !isTextFile(filePath) {
				continue
			}
			file, err := getFile(filePath)
			if err != nil {
				return nil, err
			}

			dir.Files = append(dir.Files, *file)
		}
	}

	return &dir, nil
}

func isTextFile(filePath string) bool {
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return false
	}

	knownMimes := map[string]bool{
		"text/plain; charset=utf-8": true,
		"application/x-executable":  false,
	}

	if val, ok := knownMimes[mtype.String()]; ok {
		return val
	}

	fmt.Println("Unknown MIME type:", mtype.String(), " for file:", filePath)
	os.Exit(1)
	return false
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
	finename := filepath.Base(path)
	ext := filepath.Ext(finename)

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

		if isStartBlock(finename, ext, line) {
			// start block
			inBlock = true
			buf = nil

			blockName = getBlockName(line)
			blockType = getBlockType(line)

		} else if isEndBlock(finename, ext, line) && inBlock {
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

func isStartBlock(finename, ext, line string) bool {
	if ext == ".go" && isRegexp(line, `^//\s?start`) {
		return true
	}

	if ext == ".yml" && isRegexp(line, `^#\s?start`) {
		return true
	}

	if ext == ".md" && isRegexp(line, `^\[//\]: # \(start`) {
		return true
	}

	if finename == "go.mod" && isRegexp(line, `^//\s?start`) {
		return true
	}

	if finename == ".env" && isRegexp(line, `^;\s?start`) {
		return true
	}

	if finename == ".env.example" && isRegexp(line, `^;\s?start`) {
		return true
	}

	if finename == ".gitignore" && isRegexp(line, `^#\s?start`) {
		return true
	}

	if finename == "Dockerfile" && isRegexp(line, `^#\s?start`) {
		return true
	}

	return false
}

func isEndBlock(finename, ext, line string) bool {
	if ext == ".go" && isRegexp(line, `^//\s?end`) {
		return true
	}

	if ext == ".yml" && isRegexp(line, `^#\s?end`) {
		return true
	}

	if ext == ".md" && isRegexp(line, `^\[//\]: # \(end\)`) {
		return true
	}

	if finename == "go.mod" && isRegexp(line, `^//\s?end`) {
		return true
	}

	if finename == ".env" && isRegexp(line, `^;\s?end`) {
		return true
	}

	if finename == ".env.example" && isRegexp(line, `^;\s?end`) {
		return true
	}

	if finename == ".gitignore" && isRegexp(line, `^#\s?end`) {
		return true
	}

	if finename == "Dockerfile" && isRegexp(line, `^#\s?end`) {
		return true
	}

	return false
}

func isRegexp(line string, reg string) bool {
	matched, err := regexp.MatchString(reg, line)
	if err != nil {
		panic(err)
	}

	return matched
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
