package cmdrestrict

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const restrictionsFileName = "restrictions.json"

var restrictionsDirOverride string

func restrictionsPath() (string, error) {
	if restrictionsDirOverride != "" {
		return filepath.Join(restrictionsDirOverride, restrictionsFileName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gitmap", restrictionsFileName), nil
}

func loadRestrictions() (RestrictionFile, error) {
	path, err := restrictionsPath()
	if err != nil {
		return RestrictionFile{}, err
	}
	raw, errRead := os.ReadFile(path)
	if os.IsNotExist(errRead) {
		return RestrictionFile{}, nil
	}
	if errRead != nil {
		return RestrictionFile{}, errRead
	}
	var file RestrictionFile
	if errJSON := json.Unmarshal(raw, &file); errJSON != nil {
		return RestrictionFile{}, errJSON
	}

	return file, nil
}

func saveRestrictions(file RestrictionFile) error {
	path, err := restrictionsPath()
	if err != nil {
		return err
	}
	if errMk := os.MkdirAll(filepath.Dir(path), 0o755); errMk != nil {
		return errMk
	}
	raw, errJSON := json.MarshalIndent(file, "", "  ")
	if errJSON != nil {
		return errJSON
	}
	raw = append(raw, '\n')

	return os.WriteFile(path, raw, 0o644)
}
