// Package installer — prompt_update_guard.go handles in-place updates.
package installer

import (
	"os"
	"path/filepath"
)

// HasExistingPrompts returns true if target directory already has a .prompt or prompt directory.
func HasExistingPrompts(targetDir string) bool {
	candidates := []string{
		filepath.Join(targetDir, ".prompts"),
		filepath.Join(targetDir, "prompts"),
		filepath.Join(targetDir, ".prompt-architect"),
	}

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return true
		}
	}

	return false
}
