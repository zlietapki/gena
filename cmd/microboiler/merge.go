package main

import "github.com/zlietapki/microboiler/internal/vfs"

func mergeDirs(dirs ...vfs.Directory) vfs.Directory {
	if len(dirs) == 0 {
		return vfs.Directory{}
	}
	result := dirs[0]
	for _, d := range dirs[1:] {
		result = vfs.Directory{
			Name:  result.Name,
			Files: mergeFileSlices(result.Files, d.Files),
			Dirs:  mergeDirSlices(result.Dirs, d.Dirs),
		}
	}
	return result
}

func mergeFileSlices(slices ...[]vfs.File) []vfs.File {
	index := map[string]vfs.File{}
	order := []string{}
	for _, files := range slices {
		for _, f := range files {
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
	}
	result := make([]vfs.File, 0, len(order))
	for _, name := range order {
		result = append(result, index[name])
	}
	return result
}

func mergeDirSlices(slices ...[]vfs.Directory) []vfs.Directory {
	index := map[string]vfs.Directory{}
	order := []string{}
	for _, dirs := range slices {
		for _, d := range dirs {
			if existing, ok := index[d.Name]; ok {
				index[d.Name] = mergeDirs(existing, d)
			} else {
				order = append(order, d.Name)
				index[d.Name] = d
			}
		}
	}
	result := make([]vfs.Directory, 0, len(order))
	for _, name := range order {
		result = append(result, index[name])
	}
	return result
}

func mergeBlockSlices(slices ...[]vfs.Block) []vfs.Block {
	index := map[string]vfs.Block{}
	order := []string{}
	for _, blocks := range slices {
		for _, blk := range blocks {
			if existing, ok := index[blk.Name]; ok {
				index[blk.Name] = mergeBlock(existing, blk)
			} else {
				order = append(order, blk.Name)
				index[blk.Name] = blk
			}
		}
	}
	result := make([]vfs.Block, 0, len(order))
	for _, name := range order {
		result = append(result, index[name])
	}
	return result
}

func mergeBlock(blocks ...vfs.Block) vfs.Block {
	if len(blocks) == 0 {
		return vfs.Block{}
	}

	result := blocks[0]
	for _, b := range blocks[1:] {
		switch b.Type {
		case vfs.BlockTypeOverwrite:
			result = b
		case vfs.BlockTypeAdd:
			result = vfs.Block{
				Name: result.Name,
				Type: result.Type,
				Data: append(result.Data, b.Data...),
			}
		case vfs.BlockTypeMerge:
			seen := map[string]struct{}{}
			var data []string
			for _, line := range append(result.Data, b.Data...) {
				if _, ok := seen[line]; !ok {
					seen[line] = struct{}{}
					data = append(data, line)
				}
			}
			result = vfs.Block{
				Name: result.Name,
				Type: result.Type,
				Data: data,
			}
		default:
			result = b
		}
	}

	return result
}
