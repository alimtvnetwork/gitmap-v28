package cmdupdate

import (
	"os"
	"path/filepath"
)

func resolveCurrentInstallDir() string {
	selfPath, err := os.Executable()
	if err != nil {
		return ""
	}

	realPath, errEval := filepath.EvalSymlinks(selfPath)
	if errEval == nil {
		selfPath = realPath
	}

	return filepath.Dir(selfPath)
}
