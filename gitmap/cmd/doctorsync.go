package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// timeMillisecond is a named duration for readability.
const timeMillisecond = time.Millisecond

// tryCopyWithRetry attempts to copy src to dst with retries.
func tryCopyWithRetry(src, dst string, maxAttempts int, delay time.Duration) bool {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := copyFileOverwrite(src, dst)
		if err == nil {
			return true
		}
		if attempt < maxAttempts {
			fmt.Printf(constants.DoctorRetryFmt, attempt, maxAttempts)
			time.Sleep(delay)
		}
	}

	return false
}

// tryRenameFallback renames the locked target to .old, then copies.
func tryRenameFallback(src, dst string) bool {
	backup := dst + constants.BackupSuffix
	_ = os.Remove(backup)

	err := os.Rename(dst, backup)
	if err != nil {
		return false
	}

	fmt.Println(constants.DoctorRenamedMsg)
	err = copyFileOverwrite(src, dst)
	if err != nil {
		_ = os.Rename(backup, dst)

		return false
	}

	return true
}

// tryKillAndCopy finds stale gitmap processes and terminates them.
func tryKillAndCopy(src, dst string) bool {
	if runtime.GOOS == constants.OSWindows {
		return tryKillWindows(src, dst)
	}

	return false
}

// tryKillWindows kills stale gitmap processes on Windows and retries copy.
func tryKillWindows(src, dst string) bool {
	fmt.Println(constants.DoctorKillingMsg)

	cmd := exec.Command(constants.PSBin, constants.PSNoProfile, constants.PSNonInteractive, constants.PSCommand,
		fmt.Sprintf(`Get-CimInstance Win32_Process -Filter "Name='gitmap.exe'" | `+
			`Where-Object { $_.ExecutablePath -and (Resolve-Path $_.ExecutablePath -ErrorAction SilentlyContinue).Path -eq '%s' -and $_.ProcessId -ne %d } | `+
			`ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue; $_.ProcessId }`,
			dst, os.Getpid()))

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	reportKilledProcesses(string(out))

	return copyFileOverwrite(src, dst) == nil
}

// reportKilledProcesses logs which processes were stopped.
func reportKilledProcesses(output string) {
	killed := strings.TrimSpace(output)
	if len(killed) > 0 {
		fmt.Printf(constants.DoctorKilledFmt, killed)
		time.Sleep(500 * timeMillisecond)
	}
}

// copyFileOverwrite copies src to dst, overwriting dst if it exists.
func copyFileOverwrite(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}
