package vfs

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

var Debug = false

func GetDir(path string) (*Directory, error) {
	dir := Directory{
		Name: filepath.Base(path),
	}

	debug("ReadDir: %s\n", path)
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		debug("File: %s\n", entry.Name())
		filePath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			if ignoreDir(entry.Name()) {
				debug("Ignore dir %s\n", entry.Name())
				continue
			}

			sub, err := GetDir(filePath)
			if err != nil {
				return nil, err
			}

			dir.Dirs = append(dir.Dirs, *sub)
		} else {
			if !isTextFile(filePath) {
				debug("Ignore binary %s\n", entry.Name())
				continue
			}
			if ignoreFile(entry.Name()) {
				debug("Ignore file %s\n", entry.Name())
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

func getFile(path string) (*File, error) {
	blocks, err := getBlocks(path)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:   filepath.Base(path),
		Blocks: blocks,
	}, nil
}

func getBlocks(path string) ([]Block, error) {
	finename := filepath.Base(path)
	ext := filepath.Ext(finename)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var blocks []Block

	var inBlock bool
	var blockName string
	var blockType BlockType
	var buf []string
	var wholeFile []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		//debug("# Line: %s\n", line)

		wholeFile = append(wholeFile, line)

		if isStartBlock(finename, ext, line) {
			// старт нового блока это неявный конец предыдущего
			if inBlock {
				debug("# Add prev block %q\n", blockName)
				blocks = append(blocks, Block{
					Name: blockName,
					Type: blockType,
					Data: buf,
				})
			}

			// start block
			inBlock = true
			buf = nil

			blockName = getBlockName(line)
			blockType = getBlockType(line)

			debug("# Block %q found\n", blockName)

			continue
		}

		if isEndBlock(finename, ext, line) && inBlock {
			debug("# End block found\n")

			blocks = append(blocks, Block{
				Name: blockName,
				Type: blockType,
				Data: buf,
			})

			inBlock = false
			continue
		}

		if inBlock {
			debug("# In block %q\n", blockName)

			buf = append(buf, line)
			continue
		}
	}

	if inBlock {
		debug("# EOF. Add last block\n")

		blocks = append(blocks, Block{
			Name: blockName,
			Type: blockType,
			Data: buf,
		})
	}

	// file with no blocks means single block with all file content
	if len(blocks) == 0 {
		debug("No blocks. Add all: %s\n", finename)
		blocks = append(blocks, Block{
			Name: "noblocks",
			Type: BlockTypeSingle,
			Data: wholeFile,
		})
	}

	return blocks, scanner.Err()
}

func isStartBlock(finename, ext, line string) bool {
	line = strings.TrimSpace(line)

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
	line = strings.TrimSpace(line)

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

func getBlockName(line string) string {
	for _, tok := range strings.Fields(line) {
		if strings.HasPrefix(tok, "name:") {
			return strings.TrimPrefix(tok, "name:")
		}
	}

	panic("No block name in line")
}

func getBlockType(line string) BlockType {
	blockType := ""
	for _, tok := range strings.Fields(line) {
		if strings.HasPrefix(tok, "type:") {
			blockType = strings.TrimPrefix(tok, "type:")
		}
	}

	switch blockType {
	case "merge":
		return BlockTypeMerge
	case "add":
		return BlockTypeAdd
	default:
		return BlockTypeSingle
	}
}
