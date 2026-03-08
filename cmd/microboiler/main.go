package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/zlietapki/microboiler/internal/domain"
	"github.com/zlietapki/microboiler/internal/merge"
	"github.com/zlietapki/microboiler/internal/vfs"
	"gopkg.in/yaml.v3"
)

//go:embed microboiler_grpc_server.yml
var microboilerGrpcServer string

//go:embed microboiler_rest_server.yml
var microboilerRestServer string

var optToYaml = map[string]string{
	"grpc_server": microboilerGrpcServer,
	"rest_server": microboilerRestServer,
}

func main() {
	versionFlag := flag.Bool("version", false, "print version")
	flag.BoolVar(versionFlag, "v", false, "print version")
	skipOptionsSelect := flag.Bool("skip-options-select", false, "skip options selection")
	flag.Parse()

	if *versionFlag {
		fmt.Println("version v0.0.1")
		os.Exit(0)
	}

	var opts domain.SelectedOpts
	var err error
	if *skipOptionsSelect {
		opts = domain.SelectedOpts{
			ProjectName: "check",
			Options:     []string{"grpc_server", "rest_server"},
		}
	} else {
		opts, err = getOpts()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var projects []vfs.Directory
	for _, opt := range opts.Options {
		yamlData, ok := optToYaml[opt]
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
		panic(err)
	}

	fmt.Printf("Project boilerplate generated /tmp/some/%s\n", opts.ProjectName)

	//out, err := yaml.Marshal(result)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Print(string(out))
}

func getOpts() (domain.SelectedOpts, error) {
	opts := domain.SelectedOpts{}

	huhProjectName := huh.NewInput().Title("Project name").Value(&opts.ProjectName)

	hahOpts := huh.NewMultiSelect[string]().
		Options(
			huh.NewOption("gRPC server", "grpc_server"),
			//huh.NewOption("gRPC client", "grpc_client"),
			huh.NewOption("REST server", "rest_server"),
			//huh.NewOption("REST client", "web_client"),
			//huh.NewOption("Kafka consumer", "kafka_consumer"),
			//huh.NewOption("Kafka producer", "kafka_producer"),
			//huh.NewOption("Redis", "redis"),
			//huh.NewOption("PostgreSQL", "postgres"),
		).
		Title("Microservice options").
		Value(&opts.Options)

	var ready bool
	huhConfirm := huh.NewConfirm().
		Title("Are you sure? ").
		Description("Ready to build").
		Affirmative("Yes!").
		Negative("No.").
		Value(&ready)

	huh.NewForm(huh.NewGroup(huhProjectName, hahOpts, huhConfirm)).Run()
	if !ready {
		os.Exit(0)
	}

	return opts, nil
}
