package result

import (
	"io/fs"

	"github.com/zlietapki/microboiler/pkg/vfs"
)

var FS = vfs.Directory{
	Name: "check",
	Mode: fs.FileMode(0755),
	Files: []vfs.File{
		{
			Name: "file1",
			Mode: fs.FileMode(0644),
			Data: []byte{0x73, 0x6f, 0x6d, 0x65, 0x31},
		},
	},
	Directories: []vfs.Directory{
		{
			Name: "folder1",
			Mode: fs.FileMode(0755),
			Files: nil,
			Directories: nil,
		},
		{
			Name: "folder2",
			Mode: fs.FileMode(0755),
			Files: []vfs.File{
				{
					Name: "file2_1",
					Mode: fs.FileMode(0644),
					Data: []byte{0x73, 0x6f, 0x6d, 0x65, 0x32, 0x5f, 0x31},
				},
				{
					Name: "file2_2",
					Mode: fs.FileMode(0644),
					Data: []byte{0x73, 0x6f, 0x6d, 0x65, 0x32, 0x5f, 0x32},
				},
			},
			Directories: []vfs.Directory{
				{
					Name: "folder3",
					Mode: fs.FileMode(0755),
					Files: []vfs.File{
						{
							Name: "file3_1",
							Mode: fs.FileMode(0644),
							Data: []byte{0x73, 0x6f, 0x6d, 0x65, 0x33, 0x5f, 0x31},
						},
					},
					Directories: nil,
				},
			},
		},
	},
}