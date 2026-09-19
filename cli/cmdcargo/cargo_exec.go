package cmdcargo

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
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
	exitErr, isExit := err.(*exec.ExitError)
	if isExit {
		cliexit.Exit(exitErr.ExitCode())
	}

	return err
}

func enrichCargoEnv(cargoBin string) []string {
	binDir := filepath.Dir(cargoBin)
	currentPath := os.Getenv("PATH")
	sep := string(os.PathListSeparator)
	newPath := binDir + sep + currentPath

	env := os.Environ()
	for i, e := range env {
		if len(e) >= 5 && (e[:5] == "PATH=" || e[:5] == "Path=") {
			env[i] = "PATH=" + newPath

			return env
		}
	}

	return append(env, "PATH="+newPath)
}
