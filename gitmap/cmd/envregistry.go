package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func emptyEnvRegistry() model.EnvRegistry {
	return model.EnvRegistry{
		Variables: []model.EnvVariable{},
		Paths:     []model.EnvPathEntry{},
	}
}

// loadEnvRegistry reads and parses the env-registry.json file.
func loadEnvRegistry() (model.EnvRegistry, *apperror.AppError) {
	path := constants.EnvRegistryFilePath
	data, err := os.ReadFile(path)
	if err != nil {
		return emptyEnvRegistry(), nil
	}

	var registry model.EnvRegistry
	if err = json.Unmarshal(data, &registry); err != nil {
		return model.EnvRegistry{}, apperror.WrapSimple(err, constants.ErrEnvRegistryLoad)
	}

	return registry, nil
}

func writeEnvRegistryFile(path string, registry model.EnvRegistry) error {
	data, err := json.MarshalIndent(registry, "", constants.JSONIndent)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, constants.FilePermission)
}

// saveEnvRegistry writes the env-registry.json file.
func saveEnvRegistry(registry model.EnvRegistry) *apperror.AppError {
	path := constants.EnvRegistryFilePath
	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		return apperror.WrapSimple(err, constants.ErrEnvRegistrySave)
	}
	if err := writeEnvRegistryFile(path, registry); err != nil {
		return apperror.WrapSimple(err, constants.ErrEnvRegistrySave)
	}

	return nil
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

// findEnvVariable returns the variable or an error if not found.
func findEnvVariable(registry model.EnvRegistry, name string) (model.EnvVariable, *apperror.AppError) {
	for _, v := range registry.Variables {
		if v.Name == name {
			return v, nil
		}
	}

	return model.EnvVariable{}, apperror.NewSimple(constants.ErrEnvNotFound, "E9000")
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

// checkEnvPathNotDuplicate returns an error if the path already exists.
func checkEnvPathNotDuplicate(registry model.EnvRegistry, dir string) *apperror.AppError {
	for _, p := range registry.Paths {
		if p.Path == dir {
			return apperror.NewSimple(constants.ErrEnvPathDuplicate, "E9000")
		}
	}

	return nil
}
