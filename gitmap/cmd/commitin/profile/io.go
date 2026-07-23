package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
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
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &LoadError{Path: path, Reason: "not found", Cause: err}
		}
		return nil, &LoadError{Path: path, Reason: "read failed", Cause: err}
	}
	p, err := Decode(raw)
	if err != nil {
		var le *LoadError
		if errors.As(err, &le) {
			le.Path = path
		}
		return nil, err
	}
	return p, nil
}

// SaveToDisk writes the profile JSON atomically (temp file + rename).
// Refuses overwrite unless allowOverwrite=true.
func SaveToDisk(workspaceRoot string, p *Profile, allowOverwrite bool) error {
	path := ProfilePath(workspaceRoot, p.Name)
	if !allowOverwrite {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("profile %q already exists", p.Name)
		}
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
