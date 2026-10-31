package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

func executePowerShellWithInterceptor(psPath string, cleanArgs []string) error {
	cmd := buildPowerShellCmd(psPath, cleanArgs)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var errBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

	runErr := cmd.Run()
	if runErr == nil {
		return nil
	}

	return normalizePowerShellExitError(runErr, errBuf.String(), cleanArgs)
}

func normalizePowerShellExitError(runErr error, errOutput string, cleanArgs []string) error {
	hasCmdletErr := hasPowerShellCmdletNotFoundError(errOutput)
	if hasCmdletErr {
		return handleCmdletNotFoundFailure(cleanArgs)
	}

	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		return buildPowerShellExitAbortError(exitErr.ExitCode())
	}

	return apperror.WrapSimple(runErr, "execute powershell")
}

func handleCmdletNotFoundFailure(cleanArgs []string) error {
	target := extractPreflightCandidate(cleanArgs)
	opts := PowerShellAdviceOptions{
		RawTarget:   target,
		TargetKind:  TargetKindUnknownCmdlet,
		ContextNote: "Unrecognized PowerShell cmdlet or parser failure",
	}
	renderPowerShellDiagnosticCard(opts)

	return buildPowerShellCmdletAbortError(target)
}

func hasPowerShellCmdletNotFoundError(errOutput string) bool {
	if strings.Contains(errOutput, "is not recognized as a name of a cmdlet") {
		return true
	}

	if strings.Contains(errOutput, "is not recognized as the name of a cmdlet") {
		return true
	}

	if strings.Contains(errOutput, "is not recognized as a function, script file") {
		return true
	}

	return hasPowerShellSyntaxOrPathError(errOutput)
}

func hasPowerShellSyntaxOrPathError(errOutput string) bool {
	if strings.Contains(errOutput, "Cannot find path") {
		return true
	}

	if strings.Contains(errOutput, "ParserError") || strings.Contains(errOutput, "Unexpected token") {
		return true
	}

	return false
}

func buildPowerShellCmdletAbortError(target string) error {
	msg := fmt.Sprintf("powershell cmdlet or command not recognized: %s", target)
	ctx := buildCmdletErrorContext(target)

	return apperror.NewWithDetails("powershell", "E9001", msg, "cmdpowershell", apperror.ErrorTypeAbort, apperror.SeverityWarn, ctx)
}

func buildCmdletErrorContext(target string) map[string]any {
	return map[string]any{
		"reported": true,
		"target":   target,
		"kind":     string(TargetKindUnknownCmdlet),
		"exitCode": int(cliexit.ExitCodeUsageError),
	}
}

func buildPowerShellExitAbortError(code int) error {
	msg := fmt.Sprintf("powershell process exited with code %d", code)
	ctx := map[string]any{"reported": true, "exitCode": code}

	return apperror.NewWithDetails("powershell", "E9001", msg, "cmdpowershell", apperror.ErrorTypeAbort, apperror.SeverityWarn, ctx)
}
