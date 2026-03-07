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

	tree, err := buildTree(sourceDir)
	if err != nil {
		panic(err)
	}

	varName := varNameFromPath(outputFile)

	writeGoFile(tree, outputFile, varName)

	fmt.Println("generated:", outputFile)
}
