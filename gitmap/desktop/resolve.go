// Package desktop integrates with GitHub Desktop application.
package desktop

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// ResolveCLI returns the absolute path to the GitHub Desktop CLI shim, or
// an empty string when it cannot be located. It probes PATH first (the
// canonical install configures it there), then falls back to the per-user
// install locations the Desktop installer uses on Windows and macOS — those
// directories are NOT on PATH by default on Windows, which is the root cause
// of the silent `gitmap desktop-sync` failure where `exec.LookPath("github")`
// returns ErrNotFound even though GitHub Desktop is installed.
func ResolveCLI() string {
	hit, err := exec.LookPath(constants.GitHubDesktopBin)
	if err == nil {
		return hit
	}
	return resolveFromKnownInstalls()
}

// resolveFromKnownInstalls walks the platform-specific install dirs the
// GitHub Desktop installer writes to and returns the first shim it finds.
func resolveFromKnownInstalls() string {
	for _, candidate := range knownInstallCandidates() {
		info, statErr := os.Stat(candidate)
		if statErr == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}

// knownInstallCandidates returns the platform-specific list of paths to
// probe in priority order. Windows uses a versioned `app-*` directory per
// install, so we sort newest-first by lexical order (matches semver layout).
func knownInstallCandidates() []string {
	if runtime.GOOS == constants.OSWindows {
		return windowsCandidates()
	}
	if runtime.GOOS == "darwin" {
		return darwinCandidates()
	}
	return nil
}

// windowsCandidates lists `%LOCALAPPDATA%\GitHubDesktop\bin\github.bat` and
// any per-version `app-*\bin\github.bat` siblings (newest first).
func windowsCandidates() []string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return nil
	}
	root := filepath.Join(base, "GitHubDesktop")
	out := []string{filepath.Join(root, "bin", "github.bat")}
	entries, err := os.ReadDir(root)
	if err != nil {
		return out
	}
	versionDirs := collectAppDirs(entries)
	sort.Sort(sort.Reverse(sort.StringSlice(versionDirs)))
	for _, name := range versionDirs {
		out = append(out, filepath.Join(root, name, "bin", "github.bat"))
	}
	return out
}

// collectAppDirs filters dir entries to the `app-*` versioned subdirectories
// the Squirrel-based Desktop installer creates.
func collectAppDirs(entries []os.DirEntry) []string {
	out := []string{}
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) > 4 && e.Name()[:4] == "app-" {
			out = append(out, e.Name())
		}
	}
	return out
}

// darwinCandidates returns the macOS bundle path the GitHub Desktop CLI
// installer drops the `github` shim into.
func darwinCandidates() []string {
	return []string{
		"/Applications/GitHub Desktop.app/Contents/Resources/app/static/github",
		filepath.Join(os.Getenv("HOME"), "Applications/GitHub Desktop.app/Contents/Resources/app/static/github"),
	}
}
