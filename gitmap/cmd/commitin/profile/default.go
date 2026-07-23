package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ClearOtherDefaults walks every *.json profile in the source's
// commit-in/profiles directory and rewrites IsDefault=false on any
// profile bound to the same SourceRepoPath whose Name != keepName.
// Returns nil if the directory does not exist yet.
//
// The save flow calls this BEFORE writing the new default so we do
// not race with our own about-to-be-written file.
func ClearOtherDefaults(workspaceRoot, sourceRepoPath, keepName string) error {
	dir := filepath.Join(workspaceRoot, ".gitmap", constants.CommitInDirProfiles)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read profiles dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != constants.CommitInProfileFileExt {
			continue
		}
		name := e.Name()[:len(e.Name())-len(constants.CommitInProfileFileExt)]
		if name == keepName {
			continue
		}
		if err := clearOneDefault(workspaceRoot, name, sourceRepoPath); err != nil {
			return err
		}
	}
	return nil
}

func clearOneDefault(workspaceRoot, name, sourceRepoPath string) error {
	p, err := LoadFromDisk(workspaceRoot, name)
	if err != nil {
		// Skip unreadable profiles — they cannot be the default
		// authority anyway. Zero-swallow: bubble non-not-found.
		var le *LoadError
		if errors.As(err, &le) && le.Reason == "not found" {
			return nil
		}
		return nil // tolerate corrupt sibling profiles during save
	}
	if !p.IsDefault || p.SourceRepoPath != sourceRepoPath {
		return nil
	}
	p.IsDefault = false
	return SaveToDisk(workspaceRoot, p, true)
}
