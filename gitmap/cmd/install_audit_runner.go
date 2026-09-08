package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type commandExecutionResult struct {
	CommandLine, Stdout, Stderr string
	ExitCode                    int
	DurationMs                  int64
	IsSuccess                   bool
	Err                         error
}

func executeCommandWithAudit(args []string, verbose bool) commandExecutionResult {
	cmd, stdoutBuf, stderrBuf := prepareAuditCmd(args, verbose)
	start := time.Now()
	runErr := cmd.Run()
	durationMs := time.Since(start).Milliseconds()

	return buildAuditResult(args, stdoutBuf, stderrBuf, durationMs, runErr)
}

func prepareAuditCmd(args []string, verbose bool) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	cmd := exec.Command(args[0], args[1:]...)
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd.Stdout = selectAuditWriter(&stdoutBuf, os.Stdout, verbose)
	cmd.Stderr = selectAuditWriter(&stderrBuf, os.Stderr, verbose)

	return cmd, &stdoutBuf, &stderrBuf
}

func selectAuditWriter(buf io.Writer, out io.Writer, verbose bool) io.Writer {

	if verbose {

		return io.MultiWriter(out, buf)
	}

	return buf
}

func buildAuditResult(
	args []string, stdoutBuf, stderrBuf *bytes.Buffer,
	durationMs int64, runErr error,
) commandExecutionResult {

	return commandExecutionResult{
		CommandLine: strings.Join(args, " "),
		Stdout:      stdoutBuf.String(),
		Stderr:      stderrBuf.String(),
		ExitCode:    extractRunExitCode(runErr),
		DurationMs:  durationMs,
		IsSuccess:   runErr == nil,
		Err:         runErr,
	}
}

func extractRunExitCode(err error) int {

	if err == nil {

		return 0
	}

	var exitErr *exec.ExitError

	if errors.As(err, &exitErr) {

		return exitErr.ExitCode()
	}

	return 1
}

func recordInstallExecution(
	tool, manager, version string,
	result commandExecutionResult, notes string,
) {
	splitDB, err := store.OpenInstallationSplitDB()

	if err != nil {

		return
	}

	defer splitDB.Close()

	resolvedMgr := resolvePackageManager(manager, tool)
	_ = splitDB.RecordExecution(
		tool, "install", version, resolvedMgr,
		result.DurationMs, result.IsSuccess, result.ExitCode,
		result.Stdout, result.Stderr, result.CommandLine,
		notes, "",
	)
}

func runInstallCommand(args []string, opts installOptions) error {
	result := executeCommandWithAudit(args, opts.Verbose)

	if result.IsSuccess {
		recordInstallExecution(opts.Tool, opts.Manager, opts.Version, result, "install success")
		fmt.Printf("  ✓ %s install command completed successfully.\n", opts.Tool)

		return nil
	}

	recordInstallExecution(opts.Tool, opts.Manager, opts.Version, result, "install failed")
	handleInstallError(args, opts, result)

	return nil
}

func handleInstallError(args []string, opts installOptions, res commandExecutionResult) {
	manager := resolvePackageManager(opts.Manager, opts.Tool)
	errBytes := []byte(res.Stderr + "\n" + res.Stdout)
	logPath := writeInstallErrorLog(opts.Tool, manager, opts.Version, args, errBytes, res.Err)

	printInstallFailureDetails(opts.Tool, manager, opts.Version, args, res.Err, logPath)
	appErr := buildInstallAppError(opts.Tool, manager, args, logPath, res)
	cliexit.HandleError(appErr, res.ExitCode)
}

func buildInstallAppError(
	tool, manager string, args []string,
	logPath string, res commandExecutionResult,
) *apperror.AppError {
	appErr := apperror.WrapWithDetails(
		res.Err, "cmd.installTool", "E9000",
		fmt.Sprintf("tool installation failed for %q via %s", tool, manager),
		"cmd/installtools", apperror.ErrorTypeExecution, apperror.SeverityError,
		map[string]any{
			"tool": tool, "manager": manager,
			"command": strings.Join(args, " "), "log_path": logPath,
			"exit_code": res.ExitCode,
		},
	)
	attachOutputContext(appErr, res.Stdout, res.Stderr)

	return appErr
}

func attachOutputContext(appErr *apperror.AppError, stdout, stderr string) {

	if stdout != "" {
		appErr.WithContext("stdout", strings.TrimSpace(stdout))
	}

	if stderr != "" {
		appErr.WithContext("stderr", strings.TrimSpace(stderr))
	}
}

func runPhaseWithAudit(
	tool, manager, version string,
	args []string, phaseName string, verbose bool,
) error {
	res := executeCommandWithAudit(args, verbose)
	recordInstallExecution(tool, manager, version, res, phaseName)

	if res.IsSuccess {

		return nil
	}

	return apperror.WrapSimple(res.Err, "cmd.runPhase."+phaseName)
}
