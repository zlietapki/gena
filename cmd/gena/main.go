package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/zlietapki/gena/internal/generator"
	"github.com/zlietapki/gena/internal/vfs"
	"github.com/zlietapki/gena/pkg/indexes"
)

const Version = "v1.0.0"

func main() {
	args := getArgs()

	if args.List {
		showIndexList()
		os.Exit(0)
	}

	if args.Version {
		showVersion()
		os.Exit(0)
	}

	if args.New {
		if len(args.Use) > 0 && args.Output != "" {
			newProject(args.Use, args.Output)
			os.Exit(0)
		}

		fmt.Fprintf(os.Stderr, "usage: gena new [-use <index_name>] [-use <index_name>] [-out <path>]\n")
		os.Exit(1)
	}

	usage()
}

func newProject(use []string, output string) {
	if pathExists(output) {
		printError("Output directory already exists: %s", output)
	}

	var idxs []vfs.Directory
	for _, opt := range use {
		proj, err := indexes.GetByName(opt)
		if err != nil {
			printError(err.Error())
		}

		idxs = append(idxs, *proj)
	}

	result := generator.MergeDirs(idxs...)

	if err := os.MkdirAll(output, 0755); err != nil {
		println("11111")
		printError(err.Error())
	}

	err := writeFiles(result, output, true)
	if err != nil {
		printError("Error on write project: %v\n", err)
	}

	fmt.Printf("Project boilerplate generated %s\n", output)

	// post commands
	runCmd(output, "task generate")
	runCmd(output, "go mod tidy")
	runCmd(output, "go fmt ./...")
	runCmd(output, "goimports -w -local github.com/zlietapki/boilerplate .")
}

func showIndexList() {
	idxs := indexes.Names()
	for _, idx := range idxs {
		fmt.Printf("%s\n", idx)
	}
}

func showVersion() {
	fmt.Printf("%s\n", Version)
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
