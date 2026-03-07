package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: fs_gen <src_folder> <output_file>")
	}

	sourceDir := os.Args[1]
	outputFile := os.Args[2]

	project, err := getProject(sourceDir)
	if err != nil {
		panic(err)
	}

	if err := writeYamlFile(project, outputFile); err != nil {
		panic(err)
	}
	fmt.Println("generated:", outputFile)
}
