package vfs

type BlockType int

const (
	BlockTypeOverwrite BlockType = iota
	BlockTypeAdd
	BlockTypeMerge
)

type Directory struct {
	Name  string
	Files []File
	Dirs  []Directory
}

type File struct {
	Name   string
	Blocks []Block
}

type Block struct {
	Name string
	Type BlockType
	Data []string
}
