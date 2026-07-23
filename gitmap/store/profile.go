package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// LoadProfileConfig reads the profiles.json config file.
func LoadProfileConfig(outputDir string) model.ProfileConfig {
	path := profileConfigPath(outputDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return defaultProfileConfig()
	}

	var cfg model.ProfileConfig

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return defaultProfileConfig()
	}

	return cfg
}

// SaveProfileConfig writes the profiles.json config file.
func SaveProfileConfig(outputDir string, cfg model.ProfileConfig) error {
	path := profileConfigPath(outputDir)

	if err := os.MkdirAll(filepath.Dir(path), constants.DirPermission); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, constants.DirPermission)
}

// ActiveProfileDBFile returns the DB filename for the active profile.
func ActiveProfileDBFile(outputDir string) string {
	cfg := LoadProfileConfig(outputDir)

	return ProfileDBFile(cfg.Active)
}

// ProfileDBFile returns the DB filename for a given profile name.
func ProfileDBFile(name string) string {
	if name == constants.DefaultProfileName {
		return constants.DBFile
	}

	return constants.ProfileDBPrefix + name + ".db"
}

// profileConfigPath returns the full path to profiles.json.
func profileConfigPath(outputDir string) string {
	return filepath.Join(outputDir, constants.DBDir, constants.ProfileConfigFile)
}

// defaultProfileConfig returns the initial config with only "default".
func defaultProfileConfig() model.ProfileConfig {
	return model.ProfileConfig{
		Active:   constants.DefaultProfileName,
		Profiles: []string{constants.DefaultProfileName},
	}
}
