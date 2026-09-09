package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func buildVerifyErrorLog(tool, binary string, paths []string) string {
	args := []string{binary, versionFlag(binary)}
	diagMsg := fmt.Sprintf("binary %q not found in PATH or fallback paths: %v", binary, paths)
	err := fmt.Errorf("%s", diagMsg)
	outBytes := []byte(diagMsg + "\nPATH=" + os.Getenv("PATH"))

	return writeInstallErrorLog(tool, "verifier", "latest", args, outBytes, err)
}

func recordVerificationFailure(tool, binary, logPath string) {
	cmdRes := commandExecutionResult{
		CommandLine: binary + " " + versionFlag(binary),
		Stderr:      fmt.Sprintf("verification failed: binary %s not located", binary),
		ExitCode:    1,
		DurationMs:  0,
		IsSuccess:   false,
		Err:         fmt.Errorf("binary not found: %s", binary),
	}
	recordInstallExecution(tool, "verifier", "latest", cmdRes, "verification-failed")
}

func buildVerifyAppError(tool, binary, logPath string, paths []string) *apperror.AppError {
	msg := fmt.Sprintf("post-install verification failed for %q (binary: %q)", tool, binary)
	ctx := map[string]any{
		"tool":           tool,
		"binary":         binary,
		"searched_paths": paths,
		"log_path":       logPath,
		"path_env":       os.Getenv("PATH"),
	}

	return apperror.NewWithDetails(
		"cmd.verifyInstallation", "E9000", msg,
		"cmd/installverify", apperror.ErrorTypeExecution,
		apperror.SeverityError, ctx,
	)
}

func reportVerificationFailure(tool, binary string) {
	fmt.Fprintf(os.Stderr, constants.ErrInstallVerifyFailed, tool)
	fallbackPaths := buildFallbackCandidates(binary)
	logPath := buildVerifyErrorLog(tool, binary, fallbackPaths)
	recordVerificationFailure(tool, binary, logPath)
	appErr := buildVerifyAppError(tool, binary, logPath, fallbackPaths)
	cliexit.WriteAppErrorReport(os.Stderr, appErr)
}
