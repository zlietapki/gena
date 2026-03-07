package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zlietapki/microboiler/internal/vfs"
)

func varNameFromPath(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)

	return "FS" + name[:len(name)-len(ext)]
}

func writeGoFile(project *vfs.Directory, outputFile string, varName string) {
	var buf bytes.Buffer

	buf.WriteString("package genfs\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"github.com/zlietapki/microboiler/pkg/vfs\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("var " + varName + " = ")

	writeTree(&buf, vfs.Directory{Name: project.Name, Files: project.Files, Dirs: project.Dirs}, 0, true)

	err := os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func writeTree(buf *bytes.Buffer, dir vfs.Directory, indent int, isRoot bool) {
	if isRoot {
		fmt.Fprintf(buf, "vfs.Directory{\n")
		fmt.Fprintf(buf, "%sName: %q,\n", tabs(indent+1), dir.Name)
	} else {
		fmt.Fprintf(buf, "{\n")
		fmt.Fprintf(buf, "%sName: %q,\n", tabs(indent+1), dir.Name)
	}

	// files
	if len(dir.Files) > 0 {
		fmt.Fprintf(buf, "%sFiles: []vfs.File{\n", tabs(indent+1))
		for _, f := range dir.Files {
			fmt.Fprintf(buf, "%s{\n", tabs(indent+2))
			fmt.Fprintf(buf, "%sName: %q,\n", tabs(indent+3), f.Name)
			fmt.Fprintf(buf, "%sBlocks: []vfs.Block{\n", tabs(indent+3))
			for _, block := range f.Blocks {
				fmt.Fprintf(buf, "%s{\n", tabs(indent+4))
				fmt.Fprintf(buf, "%sType: %s,\n", tabs(indent+5), blockTypeName(block.Type))
				fmt.Fprintf(buf, "%sData: %#v,\n", tabs(indent+5), block.Data)
				fmt.Fprintf(buf, "%s},\n", tabs(indent+4))
			}
			fmt.Fprintf(buf, "%s},\n", tabs(indent+3))
			fmt.Fprintf(buf, "%s},\n", tabs(indent+2))
		}

		fmt.Fprintf(buf, "%s},\n", tabs(indent+1))
	} else {
		fmt.Fprintf(buf, "%sFiles: nil,\n", tabs(indent+1))
	}

	// directories
	if len(dir.Dirs) > 0 {
		fmt.Fprintf(buf, "%sDirs: []vfs.Directory{\n", tabs(indent+1))

		for _, d := range dir.Dirs {
			buf.WriteString(tabs(indent + 2))
			writeTree(buf, d, indent+2, false)
			buf.WriteString(",\n")
		}

		fmt.Fprintf(buf, "%s},\n", tabs(indent+1))
	} else {
		fmt.Fprintf(buf, "%sDirs: nil,\n", tabs(indent+1))
	}

	fmt.Fprintf(buf, "%s}", tabs(indent))
}

func blockTypeName(t vfs.BlockType) string {
	switch t {
	case vfs.BlockTypeMerge:
		return "vfs.BlockTypeMerge"
	case vfs.BlockTypeAdd:
		return "vfs.BlockTypeAdd"
	default:
		return "vfs.BlockTypeOverwrite"
	}
}

func tabs(n int) string {
	return string(bytes.Repeat([]byte("\t"), n))
}
