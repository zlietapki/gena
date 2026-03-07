package main

import (
	_ "embed"
	"fmt"

	"github.com/zlietapki/microboiler/internal/vfs"
	"gopkg.in/yaml.v3"
)

//go:embed folder1.yml
var folder1 string

//go:embed folder2.yml
var folder2 string

func main() {
	var d1, d2 vfs.Directory
	if err := yaml.Unmarshal([]byte(folder1), &d1); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal([]byte(folder2), &d2); err != nil {
		panic(err)
	}

	result := mergeDirs(d1, d2)

	out, err := yaml.Marshal(result)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(out))
}

func mergeDirs(a, b vfs.Directory) vfs.Directory {
	return vfs.Directory{
		Name:  a.Name,
		Files: mergeFileSlices(a.Files, b.Files),
		Dirs:  mergeDirSlices(a.Dirs, b.Dirs),
	}
}

func mergeFileSlices(a, b []vfs.File) []vfs.File {
	index := map[string]vfs.File{}
	order := []string{}
	for _, f := range a {
		if _, ok := index[f.Name]; !ok {
			order = append(order, f.Name)
		}
		index[f.Name] = f
	}
	for _, f := range b {
		if existing, ok := index[f.Name]; ok {
			index[f.Name] = vfs.File{
				Name:   f.Name,
				Blocks: mergeBlockSlices(existing.Blocks, f.Blocks),
			}
		} else {
			order = append(order, f.Name)
			index[f.Name] = f
		}
	}
	result := make([]vfs.File, 0, len(order))
	for _, name := range order {
		result = append(result, index[name])
	}
	return result
}

func mergeDirSlices(a, b []vfs.Directory) []vfs.Directory {
	index := map[string]vfs.Directory{}
	order := []string{}
	for _, d := range a {
		if _, ok := index[d.Name]; !ok {
			order = append(order, d.Name)
		}
		index[d.Name] = d
	}
	for _, d := range b {
		if existing, ok := index[d.Name]; ok {
			index[d.Name] = mergeDirs(existing, d)
		} else {
			order = append(order, d.Name)
			index[d.Name] = d
		}
	}
	result := make([]vfs.Directory, 0, len(order))
	for _, name := range order {
		result = append(result, index[name])
	}
	return result
}

func mergeBlockSlices(a, b []vfs.Block) []vfs.Block {
	index := map[string]vfs.Block{}
	order := []string{}
	for _, blk := range a {
		if _, ok := index[blk.Name]; !ok {
			order = append(order, blk.Name)
		}
		index[blk.Name] = blk
	}
	for _, blk := range b {
		if existing, ok := index[blk.Name]; ok {
			index[blk.Name] = mergeBlock(existing, blk)
		} else {
			order = append(order, blk.Name)
			index[blk.Name] = blk
		}
	}
	result := make([]vfs.Block, 0, len(order))
	for _, name := range order {
		result = append(result, index[name])
	}
	return result
}

func mergeBlock(a, b vfs.Block) vfs.Block {
	switch b.Type {
	case vfs.BlockTypeOverwrite:
		return b
	case vfs.BlockTypeAdd:
		return vfs.Block{Name: a.Name, Type: a.Type, Data: append(a.Data, b.Data...)}
	case vfs.BlockTypeMerge:
		seen := map[string]struct{}{}
		var data []string
		for _, line := range append(a.Data, b.Data...) {
			if _, ok := seen[line]; !ok {
				seen[line] = struct{}{}
				data = append(data, line)
			}
		}
		return vfs.Block{Name: a.Name, Type: a.Type, Data: data}
	}
	return b
}
