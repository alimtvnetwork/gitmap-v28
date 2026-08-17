package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ProfilePath returns the absolute on-disk path for the given profile
// name inside <workspaceRoot>/.gitmap/commit-in/profiles/.
func ProfilePath(workspaceRoot, name string) string {
	return filepath.Join(workspaceRoot, ".gitmap",
		constants.CommitInDirProfiles,
		name+constants.CommitInProfileFileExt)
}

// LoadFromDisk reads + decodes a profile by name. Missing file →
// LoadError with Reason="not found".
func LoadFromDisk(workspaceRoot, name string) (*Profile, error) {
	path := ProfilePath(workspaceRoot, name)
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, &LoadError{Path: path, Reason: "not found", Cause: err}
	}
	if err != nil {
		return nil, &LoadError{Path: path, Reason: "read failed", Cause: err}
	}
	p, err := Decode(raw)
	var le *LoadError
	if errors.As(err, &le) {
		le.Path = path
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// SaveToDisk writes the profile JSON atomically (temp file + rename).
// Refuses overwrite unless allowOverwrite=true.
func SaveToDisk(workspaceRoot string, p *Profile, allowOverwrite bool) error {
	path := ProfilePath(workspaceRoot, p.Name)
	fileAlreadyExists := isExistingFile(path)
	if allowOverwrite == false && fileAlreadyExists == true {
		return fmt.Errorf("profile %q already exists", p.Name)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir profiles: %w", err)
	}
	out, err := Encode(p)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

// profileExists returns an error if a file already exists at path.
func profileExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("profile %q already exists", filepath.Base(path))
	}
	return nil
}

// isExistingFile returns true when path exists and is accessible.
func isExistingFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
