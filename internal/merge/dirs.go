package merge

import (
	"github.com/zlietapki/microboiler/internal/vfs"
)

func MergeDirs(dirs ...vfs.Directory) vfs.Directory {
	if len(dirs) == 0 {
		return vfs.Directory{}
	}

	result := dirs[0]
	for _, d := range dirs[1:] {
		result = vfs.Directory{
			Name:  result.Name,
			Dirs:  addDirs(result.Dirs, d.Dirs),
			Files: addFiles(result.Files, d.Files),
		}
	}
	return result
}

func addDirs(dirsA []vfs.Directory, dirsB []vfs.Directory) []vfs.Directory {
	dirByName := map[string]vfs.Directory{}
	var orderedNames []string

	for _, dirA := range dirsA {
		orderedNames = append(orderedNames, dirA.Name)
		dirByName[dirA.Name] = dirA
	}

	for _, dirB := range dirsB {
		if dirA, ok := dirByName[dirB.Name]; ok {
			dirByName[dirB.Name] = MergeDirs(dirA, dirB)
		} else {
			orderedNames = append(orderedNames, dirB.Name)
			dirByName[dirB.Name] = dirB
		}
	}

	result := make([]vfs.Directory, 0, len(orderedNames))

	for _, name := range orderedNames {
		result = append(result, dirByName[name])
	}

	return result
}

func addFiles(filesA, filesB []vfs.File) []vfs.File {
	fileByName := map[string]vfs.File{}
	var orderedNames []string

	for _, fileA := range filesA {
		orderedNames = append(orderedNames, fileA.Name)
		fileByName[fileA.Name] = fileA
	}

	for _, fileB := range filesB {
		if fileA, ok := fileByName[fileB.Name]; ok {
			fileByName[fileB.Name] = vfs.File{
				Name:   fileB.Name,
				Blocks: addBlocks(fileA.Blocks, fileB.Blocks),
			}
		} else {
			orderedNames = append(orderedNames, fileB.Name)
			fileByName[fileB.Name] = fileB
		}
	}

	result := make([]vfs.File, 0, len(orderedNames))

	for _, name := range orderedNames {
		result = append(result, fileByName[name])
	}

	return result
}

func addBlocks(blocksA, blocksB []vfs.Block) []vfs.Block {
	blockByName := map[string]vfs.Block{}
	var orderedNames []string

	for _, blockA := range blocksA {
		orderedNames = append(orderedNames, blockA.Name)
		blockByName[blockA.Name] = blockA
	}

	for _, blockB := range blocksB {
		if blockA, ok := blockByName[blockB.Name]; ok {
			blockByName[blockB.Name] = mergeBlocks(blockA, blockB)
		} else {
			orderedNames = append(orderedNames, blockB.Name)
			blockByName[blockB.Name] = blockB
		}
	}

	result := make([]vfs.Block, 0, len(orderedNames))

	for _, name := range orderedNames {
		result = append(result, blockByName[name])
	}

	return result
}

func mergeBlocks(a vfs.Block, b vfs.Block) vfs.Block {
	if a.Type == vfs.BlockTypeSingle {
		return a
	}

	if a.Type == vfs.BlockTypeAdd {
		return vfs.Block{
			Name: a.Name,
			Type: a.Type,
			Data: append(a.Data, b.Data...),
		}
	}

	if a.Type == vfs.BlockTypeMerge {
		return vfs.Block{
			Name: a.Name,
			Type: a.Type,
			Data: mergeLines(a.Data, b.Data),
		}
	}

	return a
}
