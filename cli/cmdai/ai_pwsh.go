// Package cmdai — ai_pwsh.go executes PowerShell automation via GitMap and logs to AI history.
package cmdai

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunAiPowershell executes PowerShell commands/scripts and records execution history.
func RunAiPowershell(args []string) *apperror.AppError {
	isEmpty := len(args) == 0
	if isEmpty {
		return apperror.NewValidationError("powershell command or script required")
	}

	repoRoot, rootErr := ResolveRepoRoot()
	hasRootErr := rootErr != nil
	if hasRootErr {
		repoRoot = "."
	}

	return executePowershellStreaming(args, repoRoot)
}

func executePowershellStreaming(args []string, repoRoot string) *apperror.AppError {
	pwshExe := resolvePowershellExe()
	cmdArgs := buildPwshArgs(args)
	cmd := exec.CommandContext(context.Background(), pwshExe, cmdArgs...)
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	start := time.Now()
	runErr := cmd.Run()
	durationMs := int(time.Since(start).Milliseconds())

	recordPwshHistory(pwshExe, args, repoRoot, durationMs, runErr)
	if runErr != nil {
		return wrapProcessError(runErr, pwshExe)
	}

	return nil
}

func resolvePowershellExe() string {
	_, errPwsh := exec.LookPath("pwsh")
	hasPwsh := errPwsh == nil
	if hasPwsh {
		return "pwsh"
	}

	return "powershell"
}

func buildPwshArgs(args []string) []string {
	prefix := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command"}
	commandStr := strings.Join(args, " ")

	return append(prefix, commandStr)
}

func recordPwshHistory(exe string, args []string, repoRoot string, durationMs int, err error) {
	argsJSON, _ := json.Marshal(args)
	exitCode := 0
	errMsg := ""
	isSuccess := err == nil
	if err != nil {
		exitCode = extractExitCode(err)
		errMsg = err.Error()
	}

	_ = store.RecordAiExecution(
		"instruction",
		exe+" "+strings.Join(args, " "),
		string(argsJSON),
		repoRoot,
		"127.0.0.1",
		durationMs,
		exitCode,
		"",
		errMsg,
		isSuccess,
	)
}
