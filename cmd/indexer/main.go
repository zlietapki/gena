package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/zlietapki/gena/internal/fsindex"
	"github.com/zlietapki/gena/internal/indexchecker"
	"github.com/zlietapki/gena/internal/vfs"
	"github.com/zlietapki/gena/pkg/indexes"
)

const (
	indexOutput = "pkg/indexes"
	Version     = "v1.0.0"
)

func main() {
	args := getArgs()

	if args.List {
		showIndexList()
		os.Exit(0)
	}

	if args.Remove != "" {
		removeIndex(args.Remove)
		os.Exit(0)
	}

	if args.Check {
		check()
		os.Exit(0)
	}

	if args.Version {
		showVersion()
		os.Exit(0)
	}

	if args.Add {
		if args.Name != "" && args.Src != "" {
			addIndex(args.Name, args.Src)
			os.Exit(0)
		}

		fmt.Fprintf(os.Stderr, "usage: indexer add [-name <name>] [-src <path>]\n")
		os.Exit(1)
	}

	usage()
}

func addIndex(name string, src string) {
	outputFile := filepath.Join(indexOutput, name+".yml")

	project, err := fsindex.IndexDir(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating index directory: %v\n", err)
		os.Exit(1)
	}

	if err := writeYamlFile(project, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating project: %v\n", err)
		os.Exit(1)
	}

	// post checks
	check()

	fmt.Println("Index added:", name)

	os.Exit(0)
}

func showIndexList() {
	idxs := indexes.Names()
	for _, idx := range idxs {
		fmt.Printf("%s\n", idx)
	}
}

func removeIndex(name string) {
	err := os.Remove(filepath.Join(indexOutput, name+".yml"))
	if err != nil {
		fmt.Printf("ERROR %s", err)
	}
}

func check() {
	ok := true
	if err := indexchecker.SingleBlocksSameContent(); err != nil {
		ok = false
		fmt.Printf("ERROR %s", err)
	}

	if err := indexchecker.BlocksSameType(); err != nil {
		ok = false
		fmt.Printf("ERROR %s", err)
	}

	if err := indexchecker.SameMode(); err != nil {
		ok = false
		fmt.Printf("ERROR %s", err)
	}

	if ok {
		fmt.Printf("Check indexes: OK\n")
	}
}

func showVersion() {
	fmt.Printf("%s\n", Version)
}

func writeYamlFile(project *vfs.Directory, outputFile string) error {
	data, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
