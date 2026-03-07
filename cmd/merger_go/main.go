package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zlietapki/microboiler/internal/vfs"
	"github.com/zlietapki/microboiler/pkg/genfs"
)

func main() {
	projects := []vfs.Directory{
		genfs.FSfolder1,
		genfs.FSfolder2,
	}
	result := mergeProjects(projects)
	showProj(result)
}

func mergeProjects(projects []vfs.Directory) vfs.Directory {
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
		for _, d := range p.Dirs {
			dirGroups[d.Name] = append(dirGroups[d.Name], d)
			dirNames[d.Name] = append(dirNames[d.Name], p.Name)
		}
	}

	var mergedDirs []vfs.Directory
	for dname, group := range dirGroups {
		mergedDirs = append(mergedDirs, mergeDirs(group, dirNames[dname], ""))
	}

	return vfs.Directory{
		Files: mergedFiles,
		Dirs:  mergedDirs,
	}
}

func mergeFiles(files []vfs.File, names []string, path string) vfs.File {
	blocksList := make([][]vfs.Block, len(files))
	for i, f := range files {
		blocksList[i] = f.Blocks
	}
	return vfs.File{
		Name:   files[0].Name,
		Blocks: mergeBlocks(blocksList, names, path+"/"+files[0].Name),
	}
}

type namedBlock struct {
	name  string
	block vfs.Block
}

func mergeBlocks(blocksList [][]vfs.Block, names []string, path string) []vfs.Block {
	maxLen := 0
	for _, blocks := range blocksList {
		if len(blocks) > maxLen {
			maxLen = len(blocks)
		}
	}

	result := make([]vfs.Block, maxLen)
	for k := 0; k < maxLen; k++ {
		var collected []namedBlock
		for i, blocks := range blocksList {
			if k < len(blocks) {
				collected = append(collected, namedBlock{names[i], blocks[k]})
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
			fmt.Fprintf(os.Stderr, "WARNING: mixed block types at %s:%d\n", path, k)
			for _, nb := range collected {
				fmt.Fprintf(os.Stderr, "  %s: %s\n", nb.name, blockTypeName(nb.block.Type))
			}
		}
		switch first.block.Type {
		case vfs.BlockTypeOverwrite:
			result[k] = collected[len(collected)-1].block
		case vfs.BlockTypeAdd:
			var data []string
			for _, nb := range collected {
				data = append(data, nb.block.Data...)
			}
			result[k] = vfs.Block{
				Type: first.block.Type,
				Data: data,
			}
		case vfs.BlockTypeMerge:
			seen := map[string]struct{}{}
			var lines []string
			for _, nb := range collected {
				for _, line := range nb.block.Data {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					if _, exists := seen[line]; !exists {
						seen[line] = struct{}{}
						lines = append(lines, line)
					}
				}
			}
			result[k] = vfs.Block{Type: first.block.Type, Data: lines}
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
		for _, sub := range d.Dirs {
			dirGroups[sub.Name] = append(dirGroups[sub.Name], sub)
			dirSourceNames[sub.Name] = append(dirSourceNames[sub.Name], names[i])
		}
	}
	var mergedDirs []vfs.Directory
	for dname, group := range dirGroups {
		mergedDirs = append(mergedDirs, mergeDirs(group, dirSourceNames[dname], path+"/"+dirs[0].Name))
	}

	return vfs.Directory{
		Name:  dirs[0].Name,
		Files: mergedFiles,
		Dirs:  mergedDirs,
	}
}
