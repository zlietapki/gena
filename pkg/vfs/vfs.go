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

type File struct {
	Name   string
	Mode   fs.FileMode
	Blocks map[string]Block
}

type Directory struct {
	Name        string
	Mode        fs.FileMode
	Files       []File
	Directories []Directory
}
