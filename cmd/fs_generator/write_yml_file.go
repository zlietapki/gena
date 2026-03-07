package main

import (
	"os"

	"github.com/zlietapki/microboiler/pkg/vfs"
	"gopkg.in/yaml.v3"
)

func writeYamlFile(project *vfs.Project, outputFile string) error {
	data, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
