package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

// UserDir returns the absolute path to the user-overlay templates directory,
// e.g. C:\Users\me\.gitmap\templates or /home/me/.gitmap/templates.
//
// The directory is NOT created here; call EnsureUserDir for that.
func UserDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(errTemplateUserDir, err)
	}

	return filepath.Join(home, userTemplatesDirName, userTemplatesSubdir), nil
}

// EnsureUserDir creates the user-overlay templates directory if it is
// missing. It is safe to call repeatedly.
func EnsureUserDir() (string, error) {
	dir, err := UserDir()
	if err != nil {
		return "", err
	}
	if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
		return "", fmt.Errorf(errTemplateMaterialize, dir, mkErr)
	}

	return dir, nil
}
