package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zlietapki/microboiler/internal/vfs"
)

// usage //go:generate go run ./cmd/fsgen
func main() {
	sourceDir := "./check"
	outputFile := "./generated_fs.go"

	var buf bytes.Buffer

	buf.WriteString("package result\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"io/fs\"\n\n")
	buf.WriteString("\t\"github.com/zlietapki/microboiler/internal/vfs\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("var FS = ")

	tree, err := buildTree(sourceDir)
	if err != nil {
		panic(err)
	}
	writeTree(&buf, tree, 0, true)

	err = os.WriteFile(outputFile, buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("generated:", outputFile)
}

func buildTree(path string) (vfs.Directory, error) {
	info, err := os.Stat(path)
	if err != nil {
		return vfs.Directory{}, err
	}

	dir := vfs.Directory{
		Name: info.Name(),
		Mode: info.Mode(),
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return vfs.Directory{}, err
	}

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			sub, err := buildTree(full)
			if err != nil {
				return vfs.Directory{}, err
			}

			dir.Directories = append(dir.Directories, sub)

			continue
		}

		data, err := os.ReadFile(full)
		if err != nil {
			return vfs.Directory{}, err
		}

		info, _ := entry.Info()

		dir.Files = append(dir.Files, vfs.File{
			Name: entry.Name(),
			Mode: info.Mode(),
			Data: data,
		})
	}

	return dir, nil
}

func writeTree(buf *bytes.Buffer, dir vfs.Directory, indent int, isRoot bool) {
	i := tabs(indent)

	if isRoot {
		fmt.Fprintf(buf, "vfs.Directory{\n")
	} else {
		fmt.Fprintf(buf, "{\n")
	}
	fmt.Fprintf(buf, "%sName: %q,\n", i+"\t", dir.Name)
	fmt.Fprintf(buf, "%sMode: fs.FileMode(%#o),\n", i+"\t", dir.Mode.Perm())

	// files
	if len(dir.Files) > 0 {
		fmt.Fprintf(buf, "%sFiles: []vfs.File{\n", i+"\t")
		for _, f := range dir.Files {
			fmt.Fprintf(buf, "%s{\n", i+"\t\t")
			fmt.Fprintf(buf, "%sName: %q,\n", i+"\t\t\t", f.Name)
			fmt.Fprintf(buf, "%sMode: fs.FileMode(%#o),\n", i+"\t\t\t", f.Mode.Perm())
			fmt.Fprintf(buf, "%sData: %#v,\n", i+"\t\t\t", f.Data)
			fmt.Fprintf(buf, "%s},\n", i+"\t\t")
		}

		fmt.Fprintf(buf, "%s},\n", i+"\t")
	} else {
		fmt.Fprintf(buf, "%sFiles: nil,\n", i+"\t")
	}

	// directories
	if len(dir.Directories) > 0 {
		fmt.Fprintf(buf, "%sDirectories: []vfs.Directory{\n", i+"\t")

		for _, d := range dir.Directories {
			buf.WriteString(i + "\t\t")
			writeTree(buf, d, indent+2, false)
			buf.WriteString(",\n")
		}

		fmt.Fprintf(buf, "%s},\n", i+"\t")
	} else {
		fmt.Fprintf(buf, "%sDirectories: nil,\n", i+"\t")
	}

	fmt.Fprintf(buf, "%s}", i)
}

func tabs(n int) string {
	return string(bytes.Repeat([]byte("\t"), n))
}
