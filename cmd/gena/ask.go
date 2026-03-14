package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"charm.land/huh/v2"
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

func getArgs(projectNames []string) (Args, error) {
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
		for name := range projectNames {
			fmt.Println(name)
		}
		os.Exit(0)
	}

	if args.ProjectName != "" && len(args.Options) > 0 && args.Output != "" {
		return args, nil
	}

	//huh.NewOption("gRPC server", "grpc_server"),
	//huh.NewOption("gRPC client", "grpc_client"),
	//	huh.NewOption("REST server", "rest_server"),
	//huh.NewOption("REST client", "web_client"),
	//huh.NewOption("Kafka consumer", "kafka_consumer"),
	//huh.NewOption("Kafka producer", "kafka_producer"),
	//huh.NewOption("Redis", "redis"),
	//huh.NewOption("PostgreSQL", "postgres"),

	var options []huh.Option[string]
	for _, name := range projectNames {
		options = append(options, huh.NewOption(name, name))
	}

	huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(options...).
				Height(len(options)+1).
				Title("Microservice options").
				Value((*[]string)(&args.Options)).
				Validate(func(s []string) error {
					if len(s) == 0 {
						return errors.New("select at least one option")
					}
					return nil
				}),
			huh.NewInput().Title("Project name").Value(&args.ProjectName).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("required")
					}
					return nil
				}),
			huh.NewInput().Title("Output path").Value(&args.Output).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("required")
					}

					if !pathExistsAndIsDir(s) {
						return errors.New("path does not exist")
					}

					return nil
				}),
		)).Run()

	return args, nil
}

func pathExistsAndIsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
