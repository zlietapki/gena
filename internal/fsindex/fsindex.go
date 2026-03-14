package fsindex

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zlietapki/gena/internal/vfs"
)

var Debug = false

const defaultBlockName = "noblocks"

func IndexDir(path string) (*vfs.Directory, error) {
	mode, err := getMode(path)
	if err != nil {
		return nil, err
	}

	dir := vfs.Directory{
		Name: filepath.Base(path),
		Mode: mode,
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

			sub, err := IndexDir(filePath)
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

			file, err := indexFile(filePath)
			if err != nil {
				return nil, err
			}

			dir.Files = append(dir.Files, *file)
		}
	}

	return &dir, nil
}

func indexFile(path string) (*vfs.File, error) {
	mode, err := getMode(path)
	if err != nil {
		return nil, err
	}

	blocks, err := indexBlocks(path)
	if err != nil {
		return nil, err
	}

	return &vfs.File{
		Name:   filepath.Base(path),
		Mode:   mode,
		Blocks: blocks,
	}, nil
}

func indexBlocks(path string) ([]vfs.Block, error) {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)

	if filename == "go.mod" {
		return getGoModBlocks(path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var blocks []vfs.Block

	var inBlock bool
	var blockName string
	var blockType vfs.BlockType
	var buf []string
	var wholeFile []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		//debug("# Line: %s\n", line)

		wholeFile = append(wholeFile, line)

		if isStartBlock(filename, ext, line) {
			// старт нового блока это неявный конец предыдущего
			if inBlock {
				debug("# Add prev block %q\n", blockName)
				blocks = append(blocks, vfs.Block{
					Name: blockName,
					Type: blockType,
					Data: buf,
				})
			}

			// start block
			inBlock = true
			buf = nil

			blockName, err = getBlockName(line)
			if err != nil {
				return nil, fmt.Errorf("%w in file %s", err, path)
			}
			blockType = getBlockType(line)

			debug("# Block found name:%q type:%q\n", blockName, blockType)

			continue
		}

		if isEndBlock(filename, ext, line) && inBlock {
			debug("# End block found\n")

			blocks = append(blocks, vfs.Block{
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

		blocks = append(blocks, vfs.Block{
			Name: blockName,
			Type: blockType,
			Data: buf,
		})
	}

	// file with no blocks means single block with all file content
	if len(blocks) == 0 {
		debug("No blocks. Add all: %s\n", filename)
		blocks = append(blocks, vfs.Block{
			Name: defaultBlockName,
			Type: vfs.BlockTypeSingle,
			Data: wholeFile,
		})
	}

	return blocks, scanner.Err()
}

func isStartBlock(filename, ext, line string) bool {
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

	if filename == "go.mod" && isRegexp(line, `^//\s?start`) {
		return true
	}

	if filename == ".env" && isRegexp(line, `^#\s?start`) {
		return true
	}

	if filename == ".env.example" && isRegexp(line, `^#\s?start`) {
		return true
	}

	if filename == ".gitignore" && isRegexp(line, `^#\s?start`) {
		return true
	}

	if filename == "Dockerfile" && isRegexp(line, `^#\s?start`) {
		return true
	}

	return false
}

func isEndBlock(filename, ext, line string) bool {
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

	if filename == "go.mod" && isRegexp(line, `^//\s?end`) {
		return true
	}

	if filename == ".env" && isRegexp(line, `^#\s?end`) {
		return true
	}

	if filename == ".env.example" && isRegexp(line, `^#\s?end`) {
		return true
	}

	if filename == ".gitignore" && isRegexp(line, `^#\s?end`) {
		return true
	}

	if filename == "Dockerfile" && isRegexp(line, `^#\s?end`) {
		return true
	}

	return false
}

func getBlockName(line string) (string, error) {
	line = strings.TrimSuffix(line, ")") // for markdown comments like [//]: # (start name:common)

	for _, tok := range strings.Fields(line) {
		if strings.HasPrefix(tok, "name:") {
			return strings.TrimPrefix(tok, "name:"), nil
		}
	}

	return "", errors.New("no block name found")
}

func getBlockType(line string) vfs.BlockType {
	line = strings.TrimSuffix(line, ")") // for markdown comments like [//]: # (start name:common)

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
		return vfs.BlockTypeSingle
	}
}
