package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {

		return ""
	}

	return home
}

func isExecutableCandidate(path string) bool {
	if path == "" {

		return false
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {

		return false
	}
	if runtime.GOOS == "windows" {

		return true
	}

	return info.Mode()&0111 != 0
}

func ensureDirInPath(dir string) {
	if dir == "" {

		return
	}
	currentPath := os.Getenv("PATH")
	sep := string(os.PathListSeparator)
	os.Setenv("PATH", dir+sep+currentPath)
}

func buildUserHomeCandidates(home, binary string) []string {
	return []string{
		filepath.Join(home, ".local", "share", "pnpm", binary),
		filepath.Join(home, ".local", "bin", binary),
		filepath.Join(home, ".local", "agy", "bin", binary),
		filepath.Join(home, ".antigravity", "bin", binary),
		filepath.Join(home, ".agy", "bin", binary),
		filepath.Join(home, "AppData", "Local", "agy", "bin", binary+".exe"),
		filepath.Join(home, ".gemini", "antigravity", "bin", binary+".cmd"),
	}
}

func buildFallbackCandidates(binary string) []string {
	var list []string
	if home := homeDir(); home != "" {
		list = append(list, buildUserHomeCandidates(home, binary)...)
	}
	list = append(list, filepath.Join("/usr", "local", "bin", binary))
	if pnpmHome := os.Getenv("PNPM_HOME"); pnpmHome != "" {
		list = append(list, filepath.Join(pnpmHome, binary))
	}

	return list
}

func findBinaryInFallbackPaths(binary string) string {
	candidates := buildFallbackCandidates(binary)
	for _, p := range candidates {
		if isExecutableCandidate(p) {
			ensureDirInPath(filepath.Dir(p))

			return p
		}
	}

	return ""
}

func resolveToolBinaryPath(binary string) string {
	path, err := exec.LookPath(binary)
	if err == nil {

		return path
	}

	return findBinaryInFallbackPaths(binary)
}
