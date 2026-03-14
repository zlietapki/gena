package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zlietapki/microboiler/internal/generator"
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

	result := generator.MergeDirs(selected...)
	result.Name = args.ProjectName

	err = createFileSystem(result, args.Output)
	if err != nil {
		printError("Error on write project: %v\n", err)
	}

	outputDir := filepath.Join(args.Output, args.ProjectName)
	fmt.Printf("Project boilerplate generated %s\n", outputDir)

	// format code
	runCmd(outputDir, "go mod tidy")
	runCmd(outputDir, "go fmt ./...")
	runCmd(outputDir, "goimports -w -local github.com/zlietapki/boilerplate .")
}

func runCmd(currentDir string, command string) {
	args := strings.Split(command, " ")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = currentDir
	if out, err := cmd.CombinedOutput(); err != nil {
		printError("Error running %s: %s %v", command, out, err)
	}
}

func printError(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}
