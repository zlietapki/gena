package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zlietapki/microboiler/internal/check_ymls"
	"github.com/zlietapki/microboiler/internal/generate"
	"github.com/zlietapki/microboiler/internal/vfs"
	"gopkg.in/yaml.v3"
)

func main() {
	args := getArgs()

	outputFile := filepath.Join("pkg/projects", args.NameProject+".yml")

	project, err := generate.GetDir(args.Src)
	if err != nil {
		panic(err)
	}

	if err := writeYamlFile(project, outputFile); err != nil {
		panic(err)
	}

	if err := check_ymls.CheckProjects(); err != nil {
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
