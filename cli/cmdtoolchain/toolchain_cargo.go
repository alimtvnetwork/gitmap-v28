package cmdtoolchain

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CheckHelpFn delegates help checking.
var CheckHelpFn func(command string, args []string)

// RunCargo is the main entry point for gitmap cargo / crg.
func RunCargo(args []string) error {
	checkToolchainHelp("cargo", args)
	if len(args) == 0 {
		return runCargoDefault()
	}

	sub := strings.ToLower(args[0])
	if isStatusSubcommand(sub) {
		return RunRustStatus()
	}

	return dispatchCargo(args)
}

func runCargoDefault() error {
	cargoBin := ResolveCargoBinary()
	if cargoBin == "" {
		return apperror.NewNotFoundError("cargo is not installed or not found in PATH")
	}

	return ExecEnrichedCommand(cargoBin, []string{"--help"})
}

func dispatchCargo(args []string) error {
	cargoBin := ResolveCargoBinary()
	if cargoBin == "" {
		return apperror.NewNotFoundError("cargo is not installed or not found in PATH")
	}

	return ExecEnrichedCommand(cargoBin, args)
}

// ExecEnrichedCommand runs a toolchain executable with enriched PATH and LIBCLANG_PATH.
func ExecEnrichedCommand(binPath string, args []string) error {
	cmd := exec.Command(binPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = EnrichToolchainEnv(binPath)

	return cmd.Run()
}

// ResolveCargoBinary locates cargo on the system.
func ResolveCargoBinary() string {
	return resolveBinaryCandidate("cargo")
}

func resolveBinaryCandidate(name string) string {
	if p, err := exec.LookPath(name); err == nil {
		return p
	}

	return findBinaryInHomeDir(name)
}

func findBinaryInHomeDir(name string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	candidates := buildBinaryCandidates(home, name)
	for _, p := range candidates {
		if isRegularFile(p) {
			return p
		}
	}

	return ""
}

func buildBinaryCandidates(home, name string) []string {
	var list []string
	if cargoHome := os.Getenv("CARGO_HOME"); cargoHome != "" {
		list = appendCandidatePath(list, filepath.Join(cargoHome, "bin"), name)
	}

	list = appendCandidatePath(list, filepath.Join(home, ".cargo", "bin"), name)
	list = append(list, filepath.Join("/usr", "bin", name), filepath.Join("/usr", "local", "bin", name))

	return list
}

func appendCandidatePath(list []string, dir, name string) []string {
	list = append(list, filepath.Join(dir, name))
	if runtime.GOOS == "windows" {
		list = append(list, filepath.Join(dir, name+".exe"))
	}

	return list
}

func isStatusSubcommand(sub string) bool {
	return sub == "status" || sub == "st" || sub == "info"
}

func checkToolchainHelp(cmd string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(cmd, args)
	}
}
