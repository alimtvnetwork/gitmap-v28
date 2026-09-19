package cmdcargo

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ResolveCargoBinary locates cargo on the system.
func ResolveCargoBinary() string {
	if p, err := exec.LookPath("cargo"); err == nil {
		return p
	}

	return findCargoInFallbackPaths("cargo")
}

// ResolveRustcBinary locates rustc on the system.
func ResolveRustcBinary() string {
	if p, err := exec.LookPath("rustc"); err == nil {
		return p
	}

	return findCargoInFallbackPaths("rustc")
}

// ResolveRustupBinary locates rustup on the system.
func ResolveRustupBinary() string {
	if p, err := exec.LookPath("rustup"); err == nil {
		return p
	}

	return findCargoInFallbackPaths("rustup")
}

func findCargoInFallbackPaths(bin string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	candidates := buildCargoCandidates(home, bin)
	for _, p := range candidates {
		if isExecutableFile(p) {
			return p
		}
	}

	return ""
}

func buildCargoCandidates(home, bin string) []string {
	var list []string
	if cargoHome := os.Getenv("CARGO_HOME"); cargoHome != "" {
		list = appendExeCandidate(list, filepath.Join(cargoHome, "bin"), bin)
	}

	list = appendExeCandidate(list, filepath.Join(home, ".cargo", "bin"), bin)
	list = append(list, filepath.Join("/usr", "bin", bin), filepath.Join("/usr", "local", "bin", bin))

	return list
}

func appendExeCandidate(list []string, dir, bin string) []string {
	list = append(list, filepath.Join(dir, bin))
	if runtime.GOOS == "windows" {
		list = append(list, filepath.Join(dir, bin+".exe"))
	}

	return list
}

func isExecutableFile(p string) bool {
	info, err := os.Stat(p)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}

	return info.Mode()&0111 != 0
}
