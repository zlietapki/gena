package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/zlietapki/microboiler/internal/merge"
	"github.com/zlietapki/microboiler/internal/vfs"
	"github.com/zlietapki/microboiler/pkg/projects"
)

func main() {
	args, err := getArgs(projects.Names())
	if err != nil {
		printError(err.Error())
	}

	if !pathExistsAndIsDir(args.Output) {
		printError("Output directory does not exist: %s", args.Output)
	}

	var selected []vfs.Directory
	for _, opt := range args.Options {
		proj, err := projects.GetByName(opt)
		if err != nil {
			printError(err.Error())
		}

		selected = append(selected, *proj)
	}

	result := merge.MergeDirs(selected...)
	result.Name = args.ProjectName

	err = createFileSystem(result, args.Output)
	if err != nil {
		printError("Error on write project: %v\n", err)
	}

	resultFolder := filepath.Join(args.Output, args.ProjectName)
	fmt.Printf("Project boilerplate generated %s\n", resultFolder)

	fmt.Printf("Running 'go mod tidy'\n")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = resultFolder
	if out, err := cmd.CombinedOutput(); err != nil {
		printError("Error running go mod tidy: %s %v", out, err)
	}
	fmt.Printf("Done\n")
}

func printError(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}
