package main

import (
	"bytes"
	"fmt"

	"github.com/zlietapki/microboiler/pkg/vfs"
)

func showProj(proj vfs.Project) {
	showFiles(proj.Files, 0)
	showDirs(proj.Dirs, 0)
}

func showFiles(files []vfs.File, ident int) {
	for _, file := range files {
		fmt.Printf("%s-----FILE-----\n", tabs(ident))
		fmt.Printf("%sFile name: %+s\n", tabs(ident), file.Name)
		fmt.Printf("%sFile mode: %+v\n", tabs(ident), file.Mode)
		showBlocks(file.Blocks, ident)
	}
}

func showBlocks(blocks vfs.Blocks, ident int) {
	for key, block := range blocks {
		fmt.Printf("%s-----BLOCK-----\n", tabs(ident))
		fmt.Printf("%sBlock name: %+v\n", tabs(ident), key)
		fmt.Printf("%sBlock type: %+v\n", tabs(ident), blockTypeName(block.Type))
		fmt.Printf("%sBlock data: %s\n", tabs(ident), block.Data)
	}
}

func showDirs(dirs []vfs.Directory, ident int) {
	for _, dir := range dirs {
		fmt.Printf("%s-----DIR-----\n", tabs(ident))
		fmt.Printf("%sDirectory name: %+s\n", tabs(ident), dir.Name)
		fmt.Printf("%sDirectory mode: %+s\n", tabs(ident), dir.Mode)
		showFiles(dir.Files, ident+1)
		showDirs(dir.Dirs, ident+1)
	}
}

func tabs(n int) string {
	return string(bytes.Repeat([]byte("\t"), n))
}

func blockTypeName(t vfs.BlockType) string {
	switch t {
	case vfs.BlockTypeOverwrite:
		return "overwrite"
	case vfs.BlockTypeAdd:
		return "add"
	case vfs.BlockTypeMerge:
		return "merge"
	default:
		return "unknown"
	}
}
