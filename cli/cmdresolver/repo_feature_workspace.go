package cmdresolver

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func ensureDirExists(destPath string) error {
	if _, statErr := os.Stat(destPath); os.IsNotExist(statErr) {
		return os.MkdirAll(destPath, 0755)
	}

	return nil
}

func resolveDefaultWorkspaceRoot() string {
	if root := queryDefaultWorkDirRoot(); len(root) > 0 {
		return root
	}

	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return `D:\work`
	}

	return resolveFallbackParentWorkspace()
}

func queryDefaultWorkDirRoot() string {
	defDB, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer defDB.Close()

	wd, err := defDB.GetDefaultWorkDir()
	if err != nil || wd == nil || len(wd.AbsolutePath) == 0 {
		return ""
	}

	info, errStat := os.Stat(wd.AbsolutePath)
	if errStat == nil && info.IsDir() {
		return wd.AbsolutePath
	}

	return ""
}

func resolveFallbackParentWorkspace() string {
	cwd, err := os.Getwd()
	if err != nil || len(cwd) == 0 {
		return "."
	}

	parent := filepath.Dir(cwd)
	if len(parent) > 0 && parent != cwd {
		return parent
	}

	return cwd
}
