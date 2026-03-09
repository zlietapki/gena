package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
)

type Args struct {
	ProjectName string
	Options     arrayFlags
	Output      string
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func getArgs(projectsAvailable map[string]string) (Args, error) {
	var args Args

	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "print version")
	var listProjects bool
	flag.BoolVar(&listProjects, "list", false, "list available projects")

	flag.StringVar(&args.ProjectName, "name", "", "result project name")
	flag.Var(&args.Options, "opt", "selected options")
	flag.StringVar(&args.Output, "output", "", "output path")
	flag.Parse()

	if versionFlag {
		fmt.Println("version v0.0.1")
		os.Exit(0)
	}

	if listProjects {
		for name := range projectsAvailable {
			fmt.Println(name)
		}
		os.Exit(0)
	}

	if args.ProjectName != "" && len(args.Options) > 0 && args.Output != "" {
		return args, nil
	}

	huhProjectName := huh.NewInput().Title("Project name").Value(&args.ProjectName)
	huhOutput := huh.NewInput().Title("Output path").Value(&args.Output)

	//huh.NewOption("gRPC server", "grpc_server"),
	//huh.NewOption("gRPC client", "grpc_client"),
	//	huh.NewOption("REST server", "rest_server"),
	//huh.NewOption("REST client", "web_client"),
	//huh.NewOption("Kafka consumer", "kafka_consumer"),
	//huh.NewOption("Kafka producer", "kafka_producer"),
	//huh.NewOption("Redis", "redis"),
	//huh.NewOption("PostgreSQL", "postgres"),

	var options []huh.Option[string]
	for name := range projectsAvailable {
		options = append(options, huh.NewOption(name, name))
	}

	huhOpts := huh.NewMultiSelect[string]().
		Options(options...).
		Title("Microservice options").
		Value((*[]string)(&args.Options))

	var ready bool
	huhConfirm := huh.NewConfirm().
		Title("Are you sure? ").
		Description("Ready to build").
		Affirmative("Yes!").
		Negative("No.").
		Value(&ready)

	huh.NewForm(huh.NewGroup(huhProjectName, huhOutput, huhOpts, huhConfirm)).Run()
	if !ready {
		os.Exit(0)
	}

	return args, nil
}
