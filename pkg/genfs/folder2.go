package genfs

import (
	"io/fs"

	"github.com/zlietapki/microboiler/pkg/vfs"
)

var FSfolder2 = vfs.Project{
	Name: "folder2",
	Files: []vfs.File{
		{
			Name: "main",
			Mode: fs.FileMode(0644),
			Blocks: vfs.Blocks{
				"3": vfs.Block{
					Type: vfs.BlockTypeOverwrite,
					Data: []byte("folder2 line 3 should overwrite\n"),
				},
				"1": vfs.Block{
					Type: vfs.BlockTypeOverwrite,
					Data: []byte("folder2 line1 should overwrite\n"),
				},
				"2": vfs.Block{
					Type: vfs.BlockTypeMerge,
					Data: []byte("\tline1\n\tline2\n\tline4\n"),
				},
			},
		},
	},
	Dirs: []vfs.Directory{
		{
			Name: "subFolder1",
			Mode: fs.FileMode(0755),
			Files: []vfs.File{
				{
					Name: "subfile1",
					Mode: fs.FileMode(0644),
					Blocks: vfs.Blocks{
						"hehe": vfs.Block{
							Type: vfs.BlockTypeMerge,
							Data: []byte("here it is hah?\n"),
						},
					},
				},
			},
			Dirs: nil,
		},
	},
}
