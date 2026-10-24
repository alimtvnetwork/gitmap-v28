package cmdcargo

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	cargoconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ResolveCargoBinary locates cargo on the system.
func ResolveCargoBinary() string {
	if p, err := exec.LookPath(cargoconst.BinCargo); err == nil {
		return p
	}

	return findCargoInFallbackPaths(cargoconst.BinCargo)
}

// ResolveRustcBinary locates rustc on the system.
func ResolveRustcBinary() string {
	if p, err := exec.LookPath(cargoconst.BinRustc); err == nil {
		return p
	}

	return findCargoInFallbackPaths(cargoconst.BinRustc)
}

// ResolveRustupBinary locates rustup on the system.
func ResolveRustupBinary() string {
	if p, err := exec.LookPath(cargoconst.BinRustup); err == nil {
		return p
	}

	return findCargoInFallbackPaths(cargoconst.BinRustup)
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
	if cargoHome := os.Getenv(string(enums.ToolchainEnvCargoHome)); cargoHome != "" {
		list = appendExeCandidate(list, filepath.Join(cargoHome, cargoconst.DirBin), bin)
	}

	list = appendExeCandidate(list, filepath.Join(home, cargoconst.DirCargo, cargoconst.DirBin), bin)
	list = append(list, filepath.Join(cargoconst.PathUsrBin, bin), filepath.Join(cargoconst.PathUsrLocalBin, bin))

	return list
}

func appendExeCandidate(list []string, dir, bin string) []string {
	list = append(list, filepath.Join(dir, bin))
	if runtime.GOOS == constants.PlatformWindows {
		list = append(list, filepath.Join(dir, bin+cargoconst.ExtExe))
	}

	return list
}

func isExecutableFile(p string) bool {
	info, err := os.Stat(p)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == constants.PlatformWindows {
		return true
	}

	return info.Mode()&0111 != 0
}
