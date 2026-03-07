package vfs

import "io/fs"

type BlockType int

const (
	BlockTypeOverwrite BlockType = iota
	BlockTypeAdd
	BlockTypeMerge
)

type Block struct {
	Data []byte
	Type BlockType
}

type Blocks map[string]Block

type File struct {
	Name   string
	Mode   fs.FileMode
	Blocks Blocks
}

type Directory struct {
	Name        string
	Mode        fs.FileMode
	Files       []File
	Directories []Directory
}

type Project struct {
	Files       []File
	Directories []Directory
}
