package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: fs_gen <src_folder> <output_file>")
	}

	sourceDir := os.Args[1]
	outputFile := os.Args[2]
	isGo := os.Args[3]

	project, err := getProject(sourceDir)
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

	fmt.Println("generated:", outputFile)
}
