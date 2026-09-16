package cmdagy

import (
	"os"
	"os/exec"
	"path/filepath"
)

func resolveAntigravityBinary() (string, bool) {
	for _, bin := range []string{"antigravity", "Antigravity", "agy"} {
		if path, err := exec.LookPath(bin); err == nil && path != "" {
			return path, true
		}
	}

	return checkKnownAntigravityPaths()
}

func checkKnownAntigravityPaths() (string, bool) {
	for _, candidate := range getCandidateAntigravityPaths() {
		if candidate != "" {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, true
			}
		}
	}

	return "", false
}

func getCandidateAntigravityPaths() []string {
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	progFiles := os.Getenv("ProgramFiles")

	return []string{
		filepath.Join(localApp, "Programs", "Antigravity", "Antigravity.exe"),
		filepath.Join(progFiles, "Antigravity", "Antigravity.exe"),
		filepath.Join(home, ".local", "share", "antigravity", "Antigravity"),
		"/usr/share/antigravity/Antigravity",
		"/opt/antigravity/Antigravity",
		"/Applications/Antigravity.app/Contents/MacOS/Antigravity",
	}
}
