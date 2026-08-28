package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// loadEnvRegistry reads and parses the env-registry.json file.
func loadEnvRegistry() model.EnvRegistry {
	path := constants.EnvRegistryFilePath
	data, err := os.ReadFile(path)

	if err != nil {
		return model.EnvRegistry{
			Variables: []model.EnvVariable{},
			Paths:     []model.EnvPathEntry{},
		}
	}

	var registry model.EnvRegistry

	err = json.Unmarshal(data, &registry)
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.New(constants.ErrEnvRegistryLoad, "E9000", nil).Error())
		os.Exit(1)
	}

	return registry
}

// saveEnvRegistry writes the env-registry.json file.
func saveEnvRegistry(registry model.EnvRegistry) {
	path := constants.EnvRegistryFilePath

	err := os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	if err != nil {
		apperror.New(constants.ErrEnvRegistrySave, "E9000", nil)
		return
	}

	data, err := json.MarshalIndent(registry, "", constants.JSONIndent)
	if err != nil {
		apperror.New(constants.ErrEnvRegistrySave, "E9000", nil)
		return
	}

	err = os.WriteFile(path, data, constants.FilePermission)
	if err != nil {
		apperror.New(constants.ErrEnvRegistrySave, "E9000", nil)
		return
	}
}

// upsertEnvVariable adds or updates a variable in the registry.
func upsertEnvVariable(registry model.EnvRegistry, name, value string) model.EnvRegistry {
	for idx, v := range registry.Variables {
		if v.Name == name {
			registry.Variables[idx].Value = value

			return registry
		}
	}

	registry.Variables = append(registry.Variables, model.EnvVariable{Name: name, Value: value})

	return registry
}

// findEnvVariable returns the variable or exits with error.
func findEnvVariable(registry model.EnvRegistry, name string) model.EnvVariable {
	for _, v := range registry.Variables {
		if v.Name == name {
			return v
		}
	}

	fmt.Fprintln(os.Stderr, apperror.New(constants.ErrEnvNotFound, "E9000", nil).Error())
	os.Exit(1)

	return model.EnvVariable{}
}

// removeEnvVariable removes a variable from the registry.
func removeEnvVariable(registry model.EnvRegistry, name string) model.EnvRegistry {
	filtered := make([]model.EnvVariable, 0, len(registry.Variables))

	for _, v := range registry.Variables {
		if v.Name == name {
			continue
		}

		filtered = append(filtered, v)
	}

	registry.Variables = filtered

	return registry
}

// removeEnvPath removes a path entry from the registry.
func removeEnvPath(registry model.EnvRegistry, dir string) model.EnvRegistry {
	filtered := make([]model.EnvPathEntry, 0, len(registry.Paths))

	for _, p := range registry.Paths {
		if p.Path == dir {
			continue
		}

		filtered = append(filtered, p)
	}

	registry.Paths = filtered

	return registry
}

// checkEnvPathNotDuplicate exits if the path already exists.
func checkEnvPathNotDuplicate(registry model.EnvRegistry, dir string) {
	for _, p := range registry.Paths {
		if p.Path == dir {
			apperror.New(constants.ErrEnvPathDuplicate, "E9000", nil)
			return
		}
	}
}
