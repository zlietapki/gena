package genfs

import (
	"github.com/zlietapki/microboiler/pkg/vfs"
)

var FSfolder1 = vfs.Project{
	Name: "folder1",
	Files: []vfs.File{
		{
			Name: "main",
			Blocks: []vfs.Block{
				{
					Type: vfs.BlockTypeOverwrite,
					Data: []string{"folder1 line1 should overwrite"},
				},
				{
					Type: vfs.BlockTypeMerge,
					Data: []string{"\tline1", "\tline2", "\tline3"},
				},
				{
					Type: vfs.BlockTypeOverwrite,
					Data: []string{"folder1 line3 should overwrite"},
				},
			},
		},
	},
	Dirs: []vfs.Directory{
		{
			Name: "subFolder1",
			Files: []vfs.File{
				{
					Name: "subfile1",
					Blocks: []vfs.Block{
						{
							Type: vfs.BlockTypeMerge,
							Data: []string{"here it is wins"},
						},
					},
				},
			},
			Dirs: nil,
		},
	},
}