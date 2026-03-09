package genfs

import (
	"github.com/zlietapki/microboiler/internal/vfs"
)

var FSfolder1 = vfs.Directory{
	Name: "folder1",
	Files: []vfs.File{
		{
			Name: "main",
			Blocks: []vfs.Block{
				{
					Type: vfs.BlockTypeSingle,
					Data: []string{"folder1 line1 should overwrite"},
				},
				{
					Type: vfs.BlockTypeMerge,
					Data: []string{"\tline1", "\tline2", "\tline3"},
				},
				{
					Type: vfs.BlockTypeSingle,
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
