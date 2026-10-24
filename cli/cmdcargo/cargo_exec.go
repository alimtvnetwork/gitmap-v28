package cmdcargo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/enums"
)

func execCargo(cargoBin string, args []string) error {
	cmd := exec.Command(cargoBin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = enrichCargoEnv(cargoBin)

	runErr := cmd.Run()
	if runErr == nil {
		return nil
	}

	return handleCargoExitError(runErr)
}

func handleCargoExitError(err error) error {
	if exitErr, isExit := err.(*exec.ExitError); isExit {
		return exitErr
	}

	return err
}

func enrichCargoEnv(cargoBin string) []string {
	newPath := buildEnrichedPath(cargoBin)

	return injectPathEnv(os.Environ(), newPath)
}

func buildEnrichedPath(cargoBin string) string {
	binDir := filepath.Dir(cargoBin)
	currentPath := os.Getenv(string(enums.ToolchainEnvPath))
	sep := string(os.PathListSeparator)

	return binDir + sep + currentPath
}

func injectPathEnv(env []string, newPath string) []string {
	targetPrefix := string(enums.ToolchainEnvPath) + "="
	for i, e := range env {
		if isPathEnvEntry(e) {
			env[i] = targetPrefix + newPath

			return env
		}
	}

	return append(env, targetPrefix+newPath)
}

func isPathEnvEntry(e string) bool {
	return len(e) >= 5 && strings.EqualFold(e[:5], "PATH=")
}
