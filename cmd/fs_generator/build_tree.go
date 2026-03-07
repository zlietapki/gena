package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/zlietapki/microboiler/pkg/vfs"
)

func buildTree(path string) (vfs.Directory, error) {
	info, err := os.Stat(path)
	if err != nil {
		return vfs.Directory{}, err
	}

	dir := vfs.Directory{
		Name: info.Name(),
		Mode: info.Mode(),
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return vfs.Directory{}, err
	}

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			sub, err := buildTree(full)
			if err != nil {
				return vfs.Directory{}, err
			}

			dir.Directories = append(dir.Directories, sub)

			continue
		}

		blocks, err := parseBlocks(full)
		if err != nil {
			return vfs.Directory{}, err
		}

		fileMode, _ := entry.Info()

		dir.Files = append(dir.Files, vfs.File{
			Name:   entry.Name(),
			Mode:   fileMode.Mode(),
			Blocks: blocks,
		})
	}

	return dir, nil
}

func parseBlocks(path string) (vfs.Blocks, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	blocks := vfs.Blocks{}
	var (
		inBlock bool
		name    string
		btype   vfs.BlockType
		buf     []string
	)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "// mb ") || strings.HasPrefix(line, "//mb ") {
			inBlock = true
			name = ""
			btype = vfs.BlockTypeOverwrite
			buf = nil

			for _, tok := range strings.Fields(line) {
				if strings.HasPrefix(tok, "name:") {
					name = strings.TrimPrefix(tok, "name:")
				} else if strings.HasPrefix(tok, "type:") {
					if strings.TrimPrefix(tok, "type:") == "merge" {
						btype = vfs.BlockTypeMerge
					}
				}
			}
			continue
		}

		if line == "// mbend" || line == "//mbend" {
			if inBlock && name != "" {
				blocks[name] = vfs.Block{
					Data: []byte(strings.Join(buf, "\n") + "\n"),
					Type: btype,
				}
			}
			inBlock = false
			continue
		}

		if inBlock {
			buf = append(buf, line)
		}
	}

	return blocks, scanner.Err()
}
