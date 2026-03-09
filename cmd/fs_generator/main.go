package main

import (
	"fmt"
	"os"

	"github.com/zlietapki/microboiler/internal/vfs"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: fs_gen <src_folder> <output_file>")
	}

	sourceDir := os.Args[1]
	outputFile := os.Args[2]
	isGo := os.Args[3]

	//vfs.Debug = true
	project, err := vfs.GetDir(sourceDir)
	if err != nil {
		panic(err)
	}

	if isGo == "go" {
		writeGoFile(project, outputFile, "FS"+project.Name)
	} else {
		if err := writeYamlFile(project, outputFile); err != nil {
			panic(err)
		}
	}

	if err := vfs.CheckAllFS("cmd/microboiler/"); err != nil {
		panic(err)
	}

	fmt.Println("generated:", outputFile)
}
