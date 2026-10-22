package templates

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// UserDir returns the absolute path to the user-overlay templates directory,
// e.g. C:\Users\me\.gitmap\templates or /home/me/.gitmap/templates.
//
// The directory is NOT created here; call EnsureUserDir for that.
func UserDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "failed to resolve user templates dir")
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
		return "", apperror.WrapSimple(mkErr, "failed to materialize templates to "+dir)
	}

	return dir, nil
}
