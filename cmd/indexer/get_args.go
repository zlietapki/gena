package main

import (
	"flag"
	"fmt"
	"os"
)

type flags struct {
	Src         string
	NameProject string
}

func getArgs() flags {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "print version")

	nameProject := ""
	flag.StringVar(&nameProject, "name", "", "name imported project")

	src := ""
	flag.StringVar(&src, "src", "", "source folder")
	flag.Parse()

	if versionFlag {
		fmt.Println("version v0.0.1")
		os.Exit(0)
	}

	return flags{
		Src:         src,
		NameProject: nameProject,
	}
}
