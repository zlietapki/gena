package main

import (
	"flag"
	"fmt"
	"os"
)

type flags struct {
	List    bool
	Add     bool
	Remove  string
	Check   bool
	Version bool
	Name    string
	Src     string
}

func getArgs() flags {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "help":
		usage()
	case "list":
		return flags{
			List: true,
		}
	case "add":
		fs := flag.NewFlagSet("add", flag.ExitOnError)
		name := fs.String("name", "", "index name")
		src := fs.String("src", "", "source folder")
		fs.Parse(os.Args[2:])

		return flags{
			Add:  true,
			Name: *name,
			Src:  *src,
		}
	case "rm":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: indexer rm index_name")
			os.Exit(1)
		}
		return flags{
			Remove: os.Args[2],
		}
	case "check":
		return flags{
			Check: true,
		}
	case "version":
		return flags{
			Version: true,
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
	}

	return flags{}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: indexer [help|list|add|rm|check|version]")
	os.Exit(1)
}
