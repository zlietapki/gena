package main

import (
	"fmt"
	"path/filepath"

	"github.com/zlietapki/microboiler/internal/vfs"
)

func main() {
	args := getArgs()

	outputFile := filepath.Join("cmd/microboiler", args.NameProject+".yml")

	//vfs.Debug = true
	project, err := vfs.GetDir(args.Src)
	if err != nil {
		panic(err)
	}

	if err := writeYamlFile(project, outputFile); err != nil {
		panic(err)
	}

	if err := vfs.CheckAllFS("cmd/microboiler/"); err != nil {
		panic(err)
	}

	fmt.Println("generated:", outputFile)
}
