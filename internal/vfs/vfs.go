package vfs

import "io/fs"

type File struct {
	Name string
	Mode fs.FileMode
	Data []byte
}

type Directory struct {
	Name        string
	Mode        fs.FileMode
	Files       []File
	Directories []Directory
}
