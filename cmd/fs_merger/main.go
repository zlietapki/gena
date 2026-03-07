package main

import (
	"bytes"
	"fmt"
	"os"

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
	fileNames := map[string][]string{}
	for _, p := range projects {
		for _, f := range p.Files {
			fileGroups[f.Name] = append(fileGroups[f.Name], f)
			fileNames[f.Name] = append(fileNames[f.Name], p.Name)
		}
	}

	var mergedFiles []vfs.File
	for fname, group := range fileGroups {
		mergedFiles = append(mergedFiles, mergeFiles(group, fileNames[fname], ""))
	}

	// projects dirs
	dirGroups := map[string][]vfs.Directory{}
	dirNames := map[string][]string{}
	for _, p := range projects {
		for _, d := range p.Directories {
			dirGroups[d.Name] = append(dirGroups[d.Name], d)
			dirNames[d.Name] = append(dirNames[d.Name], p.Name)
		}
	}

	var mergedDirs []vfs.Directory
	for dname, group := range dirGroups {
		mergedDirs = append(mergedDirs, mergeDirs(group, dirNames[dname], ""))
	}

	return vfs.Project{
		Files:       mergedFiles,
		Directories: mergedDirs,
	}
}

func mergeFiles(files []vfs.File, names []string, path string) vfs.File {
	blocksList := make([]vfs.Blocks, len(files))
	for i, f := range files {
		blocksList[i] = f.Blocks
	}
	return vfs.File{
		Name:   files[0].Name,
		Mode:   files[0].Mode,
		Blocks: mergeBlocks(blocksList, names, path+"/"+files[0].Name),
	}
}

type namedBlock struct {
	name  string
	block vfs.Block
}

func blockTypeName(t vfs.BlockType) string {
	switch t {
	case vfs.BlockTypeOverwrite:
		return "overwrite"
	case vfs.BlockTypeAdd:
		return "add"
	case vfs.BlockTypeMerge:
		return "merge"
	default:
		return "unknown"
	}
}

func mergeBlocks(blocksList []vfs.Blocks, names []string, path string) vfs.Blocks {
	keys := map[string]struct{}{}
	for _, blocks := range blocksList {
		for k := range blocks {
			keys[k] = struct{}{}
		}
	}

	result := vfs.Blocks{}
	for k := range keys {
		var collected []namedBlock
		for i, blocks := range blocksList {
			if b, ok := blocks[k]; ok {
				collected = append(collected, namedBlock{names[i], b})
			}
		}
		if len(collected) == 0 {
			continue
		}

		first := collected[0]
		hasMixed := false
		for _, nb := range collected[1:] {
			if nb.block.Type != first.block.Type {
				hasMixed = true
				break
			}
		}
		if hasMixed {
			fmt.Fprintf(os.Stderr, "warning: mixed block types at %s:%s\n", path, k)
			for _, nb := range collected {
				fmt.Fprintf(os.Stderr, "  %s: %s\n", nb.name, blockTypeName(nb.block.Type))
			}
		}
		switch first.block.Type {
		case vfs.BlockTypeOverwrite:
			result[k] = collected[len(collected)-1].block
		case vfs.BlockTypeAdd:
			var data []byte
			for _, nb := range collected {
				data = append(data, nb.block.Data...)
			}
			result[k] = vfs.Block{
				Type: first.block.Type,
				Data: data,
			}
		case vfs.BlockTypeMerge:
			seen := map[string]struct{}{}
			var lines [][]byte
			for _, nb := range collected {
				for _, line := range bytes.Split(nb.block.Data, []byte("\n")) {
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
			result[k] = vfs.Block{Type: first.block.Type, Data: bytes.Join(lines, []byte("\n"))}
		}
	}

	return result
}

func mergeDirs(dirs []vfs.Directory, names []string, path string) vfs.Directory {
	fileGroups := map[string][]vfs.File{}
	fileSourceNames := map[string][]string{}
	for i, d := range dirs {
		for _, f := range d.Files {
			fileGroups[f.Name] = append(fileGroups[f.Name], f)
			fileSourceNames[f.Name] = append(fileSourceNames[f.Name], names[i])
		}
	}
	var mergedFiles []vfs.File
	for fname, group := range fileGroups {
		mergedFiles = append(mergedFiles, mergeFiles(group, fileSourceNames[fname], path+"/"+dirs[0].Name))
	}

	dirGroups := map[string][]vfs.Directory{}
	dirSourceNames := map[string][]string{}
	for i, d := range dirs {
		for _, sub := range d.Directories {
			dirGroups[sub.Name] = append(dirGroups[sub.Name], sub)
			dirSourceNames[sub.Name] = append(dirSourceNames[sub.Name], names[i])
		}
	}
	var mergedDirs []vfs.Directory
	for dname, group := range dirGroups {
		mergedDirs = append(mergedDirs, mergeDirs(group, dirSourceNames[dname], path+"/"+dirs[0].Name))
	}

	return vfs.Directory{
		Name:        dirs[0].Name,
		Mode:        dirs[0].Mode,
		Files:       mergedFiles,
		Directories: mergedDirs,
	}
}
