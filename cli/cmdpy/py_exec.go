package cmdpy

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdai"
)

func resolveVenvPython() (string, bool) {
	venv := os.Getenv("VIRTUAL_ENV")
	if venv == "" {
		return "", false
	}
	venvBin := filepath.Join(venv, "bin", "python")
	if runtime.GOOS == "windows" {
		venvBin = filepath.Join(venv, "Scripts", "python.exe")
	}
	if _, err := os.Stat(venvBin); err == nil {
		return venvBin, true
	}
	return "", false
}

func resolveRuntimePython() (string, bool) {
	rt, appErr := cmdai.DetectPythonRuntime()
	if appErr == nil && rt.IsAvailable && rt.ExecutablePath != "" {
		return rt.ExecutablePath, true
	}
	return "", false
}

func resolveSystemCandidates() (string, bool) {
	candidates := []string{"python3", "python"}
	if runtime.GOOS == "windows" {
		candidates = append(candidates, "py")
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil && path != "" {
			return path, true
		}
	}
	return "", false
}

func resolvePythonBinary() (string, error) {
	if bin, hasVenv := resolveVenvPython(); hasVenv {
		return bin, nil
	}
	if bin, hasRuntime := resolveRuntimePython(); hasRuntime {
		return bin, nil
	}
	if bin, hasSystem := resolveSystemCandidates(); hasSystem {
		return bin, nil
	}
	return "", apperror.NewWithDetails("cmd.py.resolve", "E1032", "no Python interpreter found on PATH; please install Python 3", "cmdpy", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func buildPyCommand(binary string, args []string, dir string) *exec.Cmd {
	cmd := exec.CommandContext(context.Background(), binary, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func resolveExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}

func executePythonStreaming(binary string, args []string, dir string) (int, error) {
	cmd := buildPyCommand(binary, args, dir)
	start := time.Now()
	runErr := cmd.Run()
	durationMs := time.Since(start).Milliseconds()
	exitCode := resolveExitCode(runErr)
	recordPyTelemetry(binary, args, dir, durationMs, exitCode, runErr)
	return exitCode, runErr
}
