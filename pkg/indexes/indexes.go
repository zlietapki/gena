package indexes

import (
	"embed"
	"errors"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/zlietapki/gena/internal/vfs"
)

var ErrNotFound = errors.New("index not found")

//go:embed *.yml
var YmlFiles embed.FS

func GetYmls() map[string]string {
	projects := map[string]string{}

	entries, _ := YmlFiles.ReadDir(".")
	for _, e := range entries {
		data, _ := YmlFiles.ReadFile(e.Name())
		name := strings.TrimSuffix(e.Name(), ".yml")
		projects[name] = string(data)
	}

	return projects
}

func Names() []string {
	var names []string

	entries, _ := YmlFiles.ReadDir(".")
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".yml")
		names = append(names, name)
	}

	return names
}

func GetByName(name string) (*vfs.Directory, error) {
	entries, _ := YmlFiles.ReadDir(".")
	for _, e := range entries {
		data, _ := YmlFiles.ReadFile(e.Name())
		curName := strings.TrimSuffix(e.Name(), ".yml")

		if curName == name {
			var proj vfs.Directory
			err := yaml.Unmarshal(data, &proj)
			if err != nil {
				return nil, err
			}

			return &proj, nil
		}
	}

	return nil, ErrNotFound
}

func GetAll() (map[string]vfs.Directory, error) {
	all := make(map[string]vfs.Directory)

	entries, _ := YmlFiles.ReadDir(".")
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), ".yml")
		data, err := YmlFiles.ReadFile(e.Name())
		if err != nil {
			return nil, err
		}

		var proj vfs.Directory
		err = yaml.Unmarshal(data, &proj)
		if err != nil {
			return nil, err
		}

		all[name] = proj
	}

	return all, nil
}
