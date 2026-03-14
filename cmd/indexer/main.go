package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/zlietapki/gena/internal/fsindex"
	"github.com/zlietapki/gena/internal/indexchecker"
	"github.com/zlietapki/gena/internal/vfs"
)

const indexOutput = "pkg/indexes"

func main() {
	args := getArgs()

	outputFile := filepath.Join(indexOutput, args.NameProject+".yml")

	project, err := fsindex.IndexDir(args.Src)
	if err != nil {
		panic(err)
	}

	if err := writeYamlFile(project, outputFile); err != nil {
		panic(err)
	}

	// post checks
	if err := indexchecker.SingleBlocksSameContent(); err != nil {
		fmt.Printf("ERROR %s", err)
	}

	if err := indexchecker.BlocksSameType(); err != nil {
		fmt.Printf("ERROR %s", err)
	}

	fmt.Println("generated:", outputFile)
}

func writeYamlFile(project *vfs.Directory, outputFile string) error {
	data, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
