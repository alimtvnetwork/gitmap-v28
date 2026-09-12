// Package cmdprompt — prompt_target_sanitizer.go cleans target directory paths.
package cmdprompt

import (
	"path/filepath"
)

func SanitizeTargetDirectory(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return filepath.Clean(dir)
	}

	return abs
}
