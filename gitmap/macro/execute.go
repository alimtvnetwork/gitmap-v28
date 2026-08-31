package macro

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Execute runs all steps in a macro maintaining dynamic directory state and structured reporting.
func Execute(ctx context.Context, m *Macro, opts ExecOptions) error {
	start := time.Now()
	if !isStructuredOutput(opts) {
		printExecutionHeader(m)
	}
	initialDir, _ := os.Getwd()
	dt := NewDirTracker(initialDir)
	rep := NewExecutionReport(m.Name, len(m.Steps), start)
	return runExecuteSteps(ctx, m, opts, dt, rep, start)
}

func isStructuredOutput(opts ExecOptions) bool {
	return opts.JSON || opts.YAML || len(opts.FilePath) > 0
}

func printExecutionHeader(m *Macro) {
	fmt.Printf("  %s▶ Executing Macro: %q (%d steps)%s\n\n",
		constants.ColorCyan, m.Name, len(m.Steps), constants.ColorReset)
}

func runExecuteSteps(ctx context.Context, m *Macro, opts ExecOptions, dt *DirTracker, rep *ExecutionReport, start time.Time) error {
	var lastErr error
	for i, step := range m.Steps {
		stepExec, err := executeSingleStep(ctx, step, i+1, len(m.Steps), opts, dt)
		rep.Steps = append(rep.Steps, stepExec)
		rep.ExecutedSteps++
		if err == nil {
			continue
		}
		rep.FailedSteps++
		lastErr = err
		if !step.ContinueOnError {
			break
		}
	}
	rep.Finalize(time.Now(), lastErr != nil)
	return handleExecutionFinish(m, opts, rep, start, lastErr)
}

func handleExecutionFinish(m *Macro, opts ExecOptions, rep *ExecutionReport, start time.Time, lastErr error) error {
	if isStructuredOutput(opts) {
		_ = HandleReportOutput(rep, opts)
		return lastErr
	}
	if lastErr != nil {
		return lastErr
	}
	return printMacroCompletion(m, start)
}

func printMacroCompletion(m *Macro, start time.Time) error {
	elapsed := time.Since(start)
	fmt.Println()
	fmt.Printf("  %s✔ Macro %q completed (%d steps) · Elapsed: %.1fs%s\n\n",
		constants.ColorGreen, m.Name, len(m.Steps), elapsed.Seconds(), constants.ColorReset)
	return nil
}

func executeSingleStep(ctx context.Context, step MacroStep, idx, total int, opts ExecOptions, dt *DirTracker) (StepExecution, error) {
	if opts.DryRun {
		return executeDryRunStep(step, idx, total, opts, dt), nil
	}
	expandedCmd := ExpandPathAndEnv(step.CommandLine)
	if !isStructuredOutput(opts) {
		fmt.Printf("  [%2d/%d] ➜ %s ... ", idx, total, expandedCmd)
	}
	start := time.Now()
	if isDirChange := dt.ProcessCd(expandedCmd); isDirChange {
		return handleDirChangeStep(step, expandedCmd, dt.CurrentDir, start, opts), nil
	}
	return runStepProcess(ctx, expandedCmd, step, opts, dt, start, idx)
}

func executeDryRunStep(step MacroStep, idx, total int, opts ExecOptions, dt *DirTracker) StepExecution {
	if !isStructuredOutput(opts) {
		fmt.Printf("  [%2d/%d] ➜ (dry-run) %s\n", idx, total, step.CommandLine)
	}
	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    step.CommandLine,
		WorkingDir:     dt.CurrentDir,
		Status:         "dry-run",
		ExitCode:       0,
		ElapsedSeconds: 0,
		Logs:           []string{"dry-run simulation"},
		ErrorLogs:      []string{},
	}
}

func handleDirChangeStep(step MacroStep, expandedCmd, currentDir string, start time.Time, opts ExecOptions) StepExecution {
	elapsed := time.Since(start)
	if !isStructuredOutput(opts) {
		fmt.Printf("%s✔ ok (%.1fs)%s\n", constants.ColorGreen, elapsed.Seconds(), constants.ColorReset)
	}
	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    expandedCmd,
		WorkingDir:     currentDir,
		Status:         "success",
		ExitCode:       0,
		ElapsedSeconds: elapsed.Seconds(),
		Logs:           []string{fmt.Sprintf("Directory changed to %s", currentDir)},
		ErrorLogs:      []string{},
	}
}

func runStepProcess(ctx context.Context, cmdText string, step MacroStep, opts ExecOptions, dt *DirTracker, start time.Time, idx int) (StepExecution, error) {
	targetDir := resolveTargetDir(dt.CurrentDir, step.WorkingDir)
	outBuf, errBuf := &bytes.Buffer{}, &bytes.Buffer{}
	cmd := buildStepCmd(ctx, cmdText, targetDir, opts, outBuf, errBuf)
	err := cmd.Run()
	elapsed := time.Since(start)
	exitCode := resolveExitCode(err)
	logs := splitToLines(outBuf.String())
	errLogs := splitToLines(errBuf.String())
	if err != nil {
		return handleStepFailure(step, cmdText, targetDir, elapsed, exitCode, err, opts, idx, logs, errLogs)
	}
	if !isStructuredOutput(opts) {
		fmt.Printf("%s✔ ok (%.1fs)%s\n", constants.ColorGreen, elapsed.Seconds(), constants.ColorReset)
	}
	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    cmdText,
		WorkingDir:     targetDir,
		Status:         "success",
		ExitCode:       0,
		ElapsedSeconds: elapsed.Seconds(),
		Logs:           logs,
		ErrorLogs:      errLogs,
	}, nil
}

func handleStepFailure(step MacroStep, cmdText, targetDir string, elapsed time.Duration, exitCode int, err error, opts ExecOptions, idx int, logs, errLogs []string) (StepExecution, error) {
	if !isStructuredOutput(opts) {
		printStepFailureMsg(step, elapsed, err, idx)
	}
	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    cmdText,
		WorkingDir:     targetDir,
		Status:         "failed",
		ExitCode:       exitCode,
		ElapsedSeconds: elapsed.Seconds(),
		Logs:           logs,
		Error:          err.Error(),
		ErrorLogs:      errLogs,
	}, err
}

func printStepFailureMsg(step MacroStep, elapsed time.Duration, err error, idx int) {
	fmt.Printf("%s✖ failed (%.1fs)%s\n", constants.ColorRed, elapsed.Seconds(), constants.ColorReset)
	if !step.ContinueOnError {
		fmt.Printf("  %s✖ Step %d failed: %v%s\n", constants.ColorRed, idx, err, constants.ColorReset)
	}
}

func splitToLines(raw string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		text := strings.TrimRight(scanner.Text(), "\r\n")
		if len(text) > 0 {
			lines = append(lines, text)
		}
	}
	if lines == nil {
		return []string{}
	}
	return lines
}

func resolveTargetDir(currentDir, stepDir string) string {
	if len(currentDir) > 0 {
		return currentDir
	}
	return ExpandPathAndEnv(stepDir)
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

func buildStepCmd(ctx context.Context, cmdText, dir string, opts ExecOptions, outBuf, errBuf io.Writer) *exec.Cmd {
	var cmd *exec.Cmd
	if runtime.GOOS == constants.OSWindows {
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", cmdText)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdText)
	}
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	if opts.Verbose {
		cmd.Stdout = io.MultiWriter(os.Stdout, outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, errBuf)
		return cmd
	}
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	return cmd
}
