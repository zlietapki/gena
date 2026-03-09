package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/zlietapki/microboiler/internal/merge"
	"github.com/zlietapki/microboiler/internal/vfs"
	"gopkg.in/yaml.v3"
)

//go:embed microboiler_grpc_server.yml
var microboilerGrpcServer string

//go:embed microboiler_rest_server.yml
var microboilerRestServer string

var projectsAvailable = map[string]string{
	"grpc_server": microboilerGrpcServer,
	"rest_server": microboilerRestServer,
}

func main() {

	var opts SelectedOpts
	var err error

	opts, err = getOpts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var projects []vfs.Directory
	for _, opt := range opts.Options {
		yamlData, ok := projectsAvailable[opt]
		if !ok {
			fmt.Printf("Unknown option '%s'\n", opt)
			os.Exit(1)
		}

		var proj vfs.Directory
		if err = yaml.Unmarshal([]byte(yamlData), &proj); err != nil {
			panic(err)
		}

		projects = append(projects, proj)
	}

	result := merge.Dirs(projects...)
	result.Name = opts.ProjectName

	err = createFileSystem(result, "/tmp/some")
	if err != nil {
		fmt.Printf("Error on write project: %v\n", err)
		if os.IsExist(err) {
			fmt.Printf("Distanation folder already exists\n")
		}
		os.Exit(1)
	}

	fmt.Printf("Project boilerplate generated /tmp/some/%s\n", opts.ProjectName)

	//out, err := yaml.Marshal(result)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Print(string(out))
}
