// Package config handles loading and merging configuration.
package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// LoadFromFile reads a JSON config file and returns a Config.
//
// Returns the default config when the file does not exist (so users
// can run gitmap with no config at all). When the file DOES exist,
// it is schema-validated before defaults are applied:
//
//  1. ValidateRawConfig checks every required top-level key is
//     present in the raw JSON — catches "I forgot defaultMode"
//     bugs that a typed unmarshal would silently mask with
//     struct defaults.
//  2. ValidateRawShape checks each required key holds the right
//     JSON type (string vs number vs object) — catches
//     `"defaultMode": 42` with a per-key error message instead
//     of letting json.Unmarshal emit its generic "cannot
//     unmarshal number into Go struct field" message.
//  3. The bytes are unmarshaled onto a defaulted Config (so any
//     new optional fields added later don't break old configs).
//  4. ValidateConfig checks the resulting struct for invalid enum
//     values — catches typos like `"defaultMode": "htps"` and
//     explicit empty strings.
//  5. ValidateConfigStruct checks the populated struct for the
//     remaining shape rules: non-empty outputDir, non-negative
//     dashboardRefresh, no-empty entries in excludeDirs, and
//     complete (goos+goarch) release targets.
//
// Any validation step returning an error causes LoadFromFile to
// return that error so the CLI fails fast at startup instead of
// limping along with a partially-broken config.
func LoadFromFile(path string) (model.Config, error) {
	cfg := model.DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {

		return cfg, handleMissingFile(err)
	}
	if err := ValidateRawConfig(data); err != nil {

		return cfg, err
	}
	// Shape check runs BEFORE the typed unmarshal so a wrong-type
	// value (e.g. `"defaultMode": 42`) is reported with the offending
	// key name instead of cascading into json.Unmarshal's generic
	// "cannot unmarshal number into Go struct field ..." message.
	if err := ValidateRawShape(data); err != nil {

		return cfg, err
	}
	cfg, err = parseConfig(data, cfg)
	if err != nil {

		return cfg, err
	}
	if err := ValidateConfig(cfg); err != nil {

		return cfg, err
	}
	// Struct-level checks (non-empty outputDir, non-negative
	// dashboardRefresh, no-empty excludeDirs, complete release
	// targets) run last because they assume a defaulted-then-
	// populated struct, not raw JSON.
	if err := ValidateConfigStruct(cfg); err != nil {

		return cfg, err
	}

	return cfg, nil
}

// handleMissingFile returns nil for missing files, error otherwise.
func handleMissingFile(err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}

// parseConfig unmarshals JSON data into a Config struct.
func parseConfig(data []byte, cfg model.Config) (model.Config, error) {
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// MergeWithFlags applies CLI flag overrides to the loaded config.
// Flags take precedence when they are non-empty.
func MergeWithFlags(cfg model.Config, mode, output, outputDir string) model.Config {
	cfg = applyMode(cfg, mode)
	cfg = applyOutput(cfg, output)
	cfg = applyOutputDir(cfg, outputDir)

	return cfg
}

// applyMode overrides the default mode if the flag is set.
func applyMode(cfg model.Config, mode string) model.Config {
	if len(mode) > 0 {
		cfg.DefaultMode = mode
	}

	return cfg
}

// applyOutput overrides the default output if the flag is set.
func applyOutput(cfg model.Config, output string) model.Config {
	if len(output) > 0 {
		cfg.DefaultOutput = output
	}

	return cfg
}

// applyOutputDir overrides the output directory if the flag is set.
func applyOutputDir(cfg model.Config, outputDir string) model.Config {
	if len(outputDir) > 0 {
		cfg.OutputDir = outputDir
	}

	return cfg
}
