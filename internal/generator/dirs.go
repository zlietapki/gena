package generator

import (
	"github.com/zlietapki/gena/internal/vfs"
)

func MergeDirs(dirs ...vfs.Directory) vfs.Directory {
	if len(dirs) == 0 {
		return vfs.Directory{}
	}

	result := dirs[0]
	for _, d := range dirs[1:] {
		result = vfs.Directory{
			Name:  result.Name,
			Mode:  result.Mode,
			Dirs:  addDirs(result.Dirs, d.Dirs),
			Files: addFiles(result.Files, d.Files),
		}
	}

	//fmt.Printf("111333 result %v\n", result)
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
	var orderedFileNames []string

	for _, fileA := range filesA {
		orderedFileNames = append(orderedFileNames, fileA.Name)
		fileByName[fileA.Name] = fileA
	}

	for _, fileB := range filesB {
		if fileA, ok := fileByName[fileB.Name]; ok {
			fileByName[fileB.Name] = vfs.File{
				Name:   fileA.Name,
				Mode:   fileA.Mode,
				Blocks: addBlocks(fileA.Blocks, fileB.Blocks),
			}
		} else {
			orderedFileNames = append(orderedFileNames, fileB.Name)
			fileByName[fileB.Name] = fileB
		}
	}

	result := make([]vfs.File, 0, len(orderedFileNames))

	for _, name := range orderedFileNames {
		result = append(result, fileByName[name])
	}

	//fmt.Printf("11111 result: %v\n", result)
	return result
}

func addBlocks(blocksA, blocksB []vfs.Block) []vfs.Block {
	blockByName := map[string]vfs.Block{}
	var orderedBlockNames []string

	for _, blockA := range blocksA {
		orderedBlockNames = append(orderedBlockNames, blockA.Name)
		blockByName[blockA.Name] = blockA
	}

	for _, blockB := range blocksB {
		if blockA, ok := blockByName[blockB.Name]; ok {
			blockByName[blockB.Name] = mergeBlocks(blockA, blockB)

			// FIXME remove
			//fmt.Printf("#### result blockB.Name %s\n", blockB.Name)
			//asd := blockByName[blockB.Name]
			//for _, v := range asd.Data {
			//	fmt.Printf("%s\n", v)
			//}
		} else {
			orderedBlockNames = append(orderedBlockNames, blockB.Name)
			blockByName[blockB.Name] = blockB
		}
	}

	result := make([]vfs.Block, 0, len(orderedBlockNames))

	for _, name := range orderedBlockNames {
		//fmt.Printf("adding block %s\n", name)
		//fmt.Printf("#### result block.Name %s\n", blockByName[name].Data)
		result = append(result, blockByName[name])
	}

	//fmt.Printf("#### result  %s\n", result)
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
