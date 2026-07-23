// Package cmd — chromeprofile_process.go: process guard for CPC.
package cmd

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func isChromeRunning(goos string) (bool, error) {
	if goos == constants.OSWindows {
		return isWindowsChromeRunning()
	}
	if goos == constants.OSDarwin {
		return pgrepExact(constants.ChromeProcessMacName)
	}
	return isLinuxChromeRunning()
}

func isWindowsChromeRunning() (bool, error) {
	cmd := exec.Command(constants.ChromeProcessTasklist,
		constants.ChromeProcessTasklistFlag, constants.ChromeProcessTasklistExpr,
		constants.ChromeProcessTasklistNH)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(string(out)), constants.ChromeProcessWindowsImage), nil
}

func isLinuxChromeRunning() (bool, error) {
	var lastErr error
	for _, name := range constants.ChromeProcessLinuxNames {
		ok, err := pgrepExact(name)
		if ok || err == nil {
			return ok, nil
		}
		lastErr = err
	}
	return false, lastErr
}

func pgrepExact(name string) (bool, error) {
	cmd := exec.Command(constants.ChromeProcessPgrep, constants.ChromeProcessPgrepExact, name)
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
