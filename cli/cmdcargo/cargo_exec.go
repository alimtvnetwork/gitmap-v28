package cmdcargo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/enums"
)

func execCargo(cargoBin string, args []string) *apperror.AppError {
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

func handleCargoExitError(err error) *apperror.AppError {
	if err == nil {
		return nil
	}
	if exitErr, isExit := err.(*exec.ExitError); isExit {
		return apperror.Wrap(exitErr, "handleCargoExitError", map[string]any{"exit_code": exitErr.ExitCode()})
	}

	return apperror.Wrap(err, "handleCargoExitError", nil)
}

func enrichCargoEnv(cargoBin string) []string {
	newPath := buildEnrichedPath(cargoBin)

	return injectPathEnv(os.Environ(), newPath)
}

func buildEnrichedPath(cargoBin string) string {
	binDir := filepath.Dir(cargoBin)
	currentPath := os.Getenv(enums.ToolchainEnvPath.String())
	sep := string(os.PathListSeparator)

	return strings.Join([]string{binDir, currentPath}, sep)
}

func injectPathEnv(env []string, newPath string) []string {
	entry := enums.ToolchainEnvPath.WithValue(newPath)
	for i, e := range env {
		if isPathEnvEntry(e) {
			env[i] = entry

			return env
		}
	}

	return append(env, entry)
}

func isPathEnvEntry(e string) bool {
	prefix := enums.ToolchainEnvPath.AssignmentPrefix()

	return len(e) >= len(prefix) && strings.EqualFold(e[:len(prefix)], prefix)
}
