package main

import (
	"os"

	"github.com/charmbracelet/huh"
)

type SelectedOpts struct {
	ProjectName string
	Options     []string
}

func getOpts() (SelectedOpts, error) {
	var opts SelectedOpts

	huhProjectName := huh.NewInput().Title("Project name").Value(&opts.ProjectName)

	var options []huh.Option[string]
	for name := range projectsAvailable {
		options = append(options, huh.NewOption(name, name))
	}

	//huh.NewOption("gRPC server", "grpc_server"),
	//huh.NewOption("gRPC client", "grpc_client"),
	//	huh.NewOption("REST server", "rest_server"),
	//huh.NewOption("REST client", "web_client"),
	//huh.NewOption("Kafka consumer", "kafka_consumer"),
	//huh.NewOption("Kafka producer", "kafka_producer"),
	//huh.NewOption("Redis", "redis"),
	//huh.NewOption("PostgreSQL", "postgres"),

	hahOpts := huh.NewMultiSelect[string]().
		Options(options...).
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
