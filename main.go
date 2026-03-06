package main

import (
	"os"

	"github.com/charmbracelet/huh"
)

func main() {
	var projectName string
	huhProjectName := huh.NewInput().Title("Project name").Value(&projectName)

	var options []string
	hahOpts := huh.NewMultiSelect[string]().
		Options(
			huh.NewOption("Web server", "web_server"),
			huh.NewOption("Web client", "web_client"),
			huh.NewOption("gRPC server", "grpc_server"),
			huh.NewOption("gRPC client", "grpc_client"),
			huh.NewOption("Kafka consumer", "kafka_consumer"),
			huh.NewOption("Kafka producer", "kafka_producer"),
			huh.NewOption("Redis", "redis"),
			huh.NewOption("PostgreSQL", "postgres"),
		).
		Title("Microservice options").
		Value(&options)

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
}
