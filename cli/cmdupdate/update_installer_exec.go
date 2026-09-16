package cmdupdate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runRemoteInstaller(scriptPath string) error {
	installDir := resolveCurrentInstallDir()
	cmd := buildRemoteInstallerCmd(scriptPath, installDir)
	cmd.Env = append(os.Environ(), "GITMAP_UPDATING=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = filepath.Dir(scriptPath)

	return cmd.Run()
}

func handleRemoteInstallerError(errRun error) {
	var exitErr *exec.ExitError
	if errors.As(errRun, &exitErr) {
		handleInstallerExitError(exitErr)

		return
	}

	fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteRun, errRun)
}

func handleInstallerExitError(exitErr *exec.ExitError) {
	appErr := apperror.NewWithDetails(
		"cmd.updateremoteinstall.run",
		"E1153",
		fmt.Sprintf("remote installer exited with code %d", exitErr.ExitCode()),
		"cmd.updateremoteinstall",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		map[string]any{"exitCode": exitErr.ExitCode()},
	)
	cliexit.HandleError(appErr, exitErr.ExitCode())
}

func resolveCurrentInstallDir() string {
	selfPath, err := os.Executable()
	if err != nil {
		return ""
	}

	realPath, errEval := filepath.EvalSymlinks(selfPath)
	if errEval == nil {
		selfPath = realPath
	}

	return filepath.Dir(selfPath)
}

func printPostUpdateIdentity() {
	installDir := resolveCurrentInstallDir()
	binName := constants.GitMapBin
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	binPath := filepath.Join(installDir, binName)
	if executeInstalledBinaryIdentity(binPath) {
		return
	}

	printGitmapIdentityBlockLong()
}

func executeInstalledBinaryIdentity(binPath string) bool {
	if _, err := os.Stat(binPath); err != nil {
		return false
	}

	cmd := exec.Command(binPath, "binary")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run() == nil
}
