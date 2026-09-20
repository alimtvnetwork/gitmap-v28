package osclean

import (
	"os"
	"path/filepath"
	"runtime"
)

func resolveHomeDir() string {
	if home := os.Getenv("USERPROFILE"); len(home) > 0 && runtime.GOOS == "windows" {
		return home
	}
	if home := os.Getenv("HOME"); len(home) > 0 {
		return home
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return dir
}

func resolveLocalAppData() string {
	if local := os.Getenv("LOCALAPPDATA"); len(local) > 0 {
		return local
	}
	home := resolveHomeDir()
	if len(home) > 0 && runtime.GOOS == "windows" {
		return filepath.Join(home, "AppData", "Local")
	}

	return ""
}

func filterUniquePaths(paths []string) []string {
	seen := make(map[string]bool)
	var res []string
	for _, p := range paths {
		clean := filepath.Clean(p)
		if len(clean) > 0 && !seen[clean] {
			seen[clean] = true
			res = append(res, clean)
		}
	}

	return res
}

func resolveDevDrives(subpath string) []string {
	var res []string
	if devDir := os.Getenv("DEV_DIR"); len(devDir) > 0 {
		res = append(res, filepath.Join(devDir, subpath))
	}
	if runtime.GOOS == "windows" {
		for _, drive := range []string{`C:\`, `D:\`, `E:\`} {
			res = append(res, filepath.Join(drive, "dev-tool", subpath))
		}
	}

	return res
}
