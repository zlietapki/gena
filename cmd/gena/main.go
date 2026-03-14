package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zlietapki/gena/internal/generator"
	"github.com/zlietapki/gena/internal/vfs"
	"github.com/zlietapki/gena/pkg/indexes"
)

func main() {
	args, err := getArgs(indexes.Names())
	if err != nil {
		printError(err.Error())
	}

	if !pathExistsAndIsDir(args.Output) {
		printError("Output directory does not exist: %s", args.Output)
	}

	var selected []vfs.Directory
	for _, opt := range args.Options {
		proj, err := indexes.GetByName(opt)
		if err != nil {
			printError(err.Error())
		}

		selected = append(selected, *proj)
	}

	result := generator.MergeDirs(selected...)
	result.Name = args.ProjectName

	err = writeFiles(result, args.Output)
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
