package main

import (
	"bytes"

	"github.com/zlietapki/microboiler/pkg/genfs"
	"github.com/zlietapki/microboiler/pkg/vfs"
)

func main() {
	projects := []vfs.Project{
		genfs.FSfolder1,
		genfs.FSfolder2,
	}
	result := mergeProjects(projects)
	showProj(result)
}

func mergeProjects(projects []vfs.Project) vfs.Project {
	// projects files
	fileGroups := map[string][]vfs.File{}
	for _, p := range projects {
		for _, f := range p.Files {
			fileGroups[f.Name] = append(fileGroups[f.Name], f)
		}
	}

	var mergedFiles []vfs.File
	for _, group := range fileGroups {
		mergedFiles = append(mergedFiles, mergeFiles(group))
	}

	// projects dirs
	dirGroups := map[string][]vfs.Directory{}
	for _, p := range projects {
		for _, d := range p.Directories {
			dirGroups[d.Name] = append(dirGroups[d.Name], d)
		}
	}

	var mergedDirs []vfs.Directory
	for _, group := range dirGroups {
		mergedDirs = append(mergedDirs, mergeDirs(group))
	}

	return vfs.Project{
		Files:       mergedFiles,
		Directories: mergedDirs,
	}
}

func mergeFiles(files []vfs.File) vfs.File {
	blocksList := make([]vfs.Blocks, len(files))
	for i, f := range files {
		blocksList[i] = f.Blocks
	}
	return vfs.File{
		Name:   files[0].Name,
		Mode:   files[0].Mode,
		Blocks: mergeBlocks(blocksList),
	}
}

func mergeBlocks(blocksList []vfs.Blocks) vfs.Blocks {
	keys := map[string]struct{}{}
	for _, blocks := range blocksList {
		for k := range blocks {
			keys[k] = struct{}{}
		}
	}

	result := vfs.Blocks{}
	for k := range keys {
		var collected []vfs.Block
		for _, blocks := range blocksList {
			if b, ok := blocks[k]; ok {
				collected = append(collected, b)
			}
		}
		if len(collected) == 0 {
			continue
		}

		first := collected[0]
		switch first.Type {
		case vfs.BlockTypeOverwrite:
			result[k] = collected[len(collected)-1]
		case vfs.BlockTypeAdd:
			var data []byte
			for _, b := range collected {
				data = append(data, b.Data...)
			}
			result[k] = vfs.Block{
				Type: first.Type,
				Data: data,
			}
		case vfs.BlockTypeMerge:
			seen := map[string]struct{}{}
			var lines [][]byte
			for _, b := range collected {
				for _, line := range bytes.Split(b.Data, []byte("\n")) {
					line = bytes.TrimSpace(line)
					if len(line) == 0 {
						continue
					}
					s := string(line)
					if _, exists := seen[s]; !exists {
						seen[s] = struct{}{}
						lines = append(lines, line)
					}
				}
			}
			result[k] = vfs.Block{Type: first.Type, Data: bytes.Join(lines, []byte("\n"))}
		}
	}

	return result
}

func mergeDirs(dirs []vfs.Directory) vfs.Directory {
	fileGroups := map[string][]vfs.File{}
	for _, d := range dirs {
		for _, f := range d.Files {
			fileGroups[f.Name] = append(fileGroups[f.Name], f)
		}
	}
	var mergedFiles []vfs.File
	for _, group := range fileGroups {
		mergedFiles = append(mergedFiles, mergeFiles(group))
	}

	dirGroups := map[string][]vfs.Directory{}
	for _, d := range dirs {
		for _, sub := range d.Directories {
			dirGroups[sub.Name] = append(dirGroups[sub.Name], sub)
		}
	}
	var mergedDirs []vfs.Directory
	for _, group := range dirGroups {
		mergedDirs = append(mergedDirs, mergeDirs(group))
	}

	return vfs.Directory{
		Name:        dirs[0].Name,
		Mode:        dirs[0].Mode,
		Files:       mergedFiles,
		Directories: mergedDirs,
	}
}
