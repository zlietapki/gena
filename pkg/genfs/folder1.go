package genfs

import (
	"io/fs"

	"github.com/zlietapki/microboiler/pkg/vfs"
)

var FSfolder1 = vfs.Project{
	Files: []vfs.File{
		{
			Name: "main",
			Mode: fs.FileMode(0644),
			Blocks: vfs.Blocks{
				"1": vfs.Block{
					Type: vfs.BlockTypeOverwrite,
					Data: []byte("folder1 line1 should overwrite\n"),
				},
				"2": vfs.Block{
					Type: vfs.BlockTypeMerge,
					Data: []byte("\tline1\n\tline2\n\tline3\n"),
				},
				"3": vfs.Block{
					Type: vfs.BlockTypeOverwrite,
					Data: []byte("folder1 line3 should overwrite\n"),
				},
			},
		},
	},
	Directories: []vfs.Directory{
		{
			Name: "subFolder1",
			Mode: fs.FileMode(0755),
			Files: []vfs.File{
				{
					Name: "subfile1",
					Mode: fs.FileMode(0644),
					Blocks: vfs.Blocks{
						"hehe": vfs.Block{
							Type: vfs.BlockTypeOverwrite,
							Data: []byte("here it is wins\n"),
						},
					},
				},
			},
			Directories: nil,
		},
	},
}