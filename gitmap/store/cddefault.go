package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// LoadCDDefaults reads the cd-defaults.json config file.
func LoadCDDefaults(outputDir string) map[string]string {
	path := cdDefaultsPath(outputDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]string)
	}

	var defaults map[string]string

	err = json.Unmarshal(data, &defaults)
	if err != nil {
		return make(map[string]string)
	}

	return defaults
}

// SaveCDDefaults writes the cd-defaults.json config file.
func SaveCDDefaults(outputDir string, defaults map[string]string) error {
	path := cdDefaultsPath(outputDir)

	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		return err
	}

	data, err := json.MarshalIndent(defaults, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, constants.DirPermission)
}

// cdDefaultsPath returns the full path to cd-defaults.json.
func cdDefaultsPath(outputDir string) string {
	return filepath.Join(outputDir, constants.DBDir, constants.CDDefaultsFile)
}
