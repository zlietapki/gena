package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zlietapki/microboiler/internal/merge"
	"github.com/zlietapki/microboiler/internal/vfs"
	"github.com/zlietapki/microboiler/pkg/projects"
	"gopkg.in/yaml.v3"
)

func main() {
	var projectsAvailable = getProjectAvailable()

	args, err := getArgs(projectsAvailable)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var selected []vfs.Directory
	for _, opt := range args.Options {
		yamlData, ok := projectsAvailable[opt]
		if !ok {
			fmt.Printf("Unknown option '%s'\n", opt)
			os.Exit(1)
		}

		var proj vfs.Directory
		if err = yaml.Unmarshal([]byte(yamlData), &proj); err != nil {
			panic(err)
		}

		selected = append(selected, proj)
	}

	result := merge.MergeDirs(selected...)
	result.Name = args.ProjectName

	err = createFileSystem(result, args.Output)
	if err != nil {
		fmt.Printf("Error on write project: %v\n", err)
		if os.IsExist(err) {
			fmt.Printf("Destination folder already exists\n")
		}
		os.Exit(1)
	}

	resultFolder := filepath.Join(args.Output, args.ProjectName)
	fmt.Printf("Project boilerplate generated %s\n", resultFolder)

	fmt.Printf("Running 'go mod tidy'\n")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = resultFolder
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("go mod tidy: %s\n", out)
		os.Exit(1)
	}
	fmt.Printf("Done\n")
}

func getProjectAvailable() map[string]string {
	m := map[string]string{}
	entries, _ := projects.YmlFiles.ReadDir(".")
	for _, e := range entries {
		data, _ := projects.YmlFiles.ReadFile(e.Name())
		key := strings.TrimSuffix(e.Name(), ".yml")
		m[key] = string(data)
	}
	return m
}
