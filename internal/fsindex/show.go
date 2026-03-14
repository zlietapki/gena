package fsindex

import (
	"bytes"
	"fmt"

	"github.com/zlietapki/microboiler/internal/vfs"
)

func showDir(proj vfs.Directory) {
	showFiles(proj.Files, 0)
	showDirs(proj.Dirs, 0)
}

func showFiles(files []vfs.File, ident int) {
	for _, file := range files {
		fmt.Printf("%s-----FILE-----\n", tabs(ident))
		fmt.Printf("%sFile name: %#v\n", tabs(ident), file.Name)
		showBlocks(file.Blocks, ident)
	}
}

func showBlocks(blocks []vfs.Block, ident int) {
	for key, block := range blocks {
		fmt.Printf("%s-----BLOCK-----\n", tabs(ident))
		fmt.Printf("%sBlock name: %#v\n", tabs(ident), key)
		fmt.Printf("%sBlock type: %#v\n", tabs(ident), blockTypeName(block.Type))
		fmt.Printf("%sBlock data: %#v\n", tabs(ident), block.Data)
	}
}

func showDirs(dirs []vfs.Directory, ident int) {
	for _, dir := range dirs {
		fmt.Printf("%s-----DIR-----\n", tabs(ident))
		fmt.Printf("%sDirectory name: %#v\n", tabs(ident), dir.Name)
		showFiles(dir.Files, ident+1)
		showDirs(dir.Dirs, ident+1)
	}
}

func tabs(n int) string {
	return string(bytes.Repeat([]byte("\t"), n))
}

func blockTypeName(t vfs.BlockType) string {
	switch t {
	case vfs.BlockTypeSingle:
		return "single"
	case vfs.BlockTypeAdd:
		return "add"
	case vfs.BlockTypeMerge:
		return "merge"
	default:
		return "unknown"
	}
}
