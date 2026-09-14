package macro

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ParseAsyncMacroCommand extracts shell, command, and interval from an async step.
func ParseAsyncMacroCommand(raw string) (bool, AsyncStepOpts) {
	trimmed := strings.TrimSpace(raw)
	isAsyncPrefix := strings.HasPrefix(trimmed, "async ") || trimmed == "async"
	if !isAsyncPrefix {
		return false, AsyncStepOpts{}
	}

	rem := strings.TrimSpace(strings.TrimPrefix(trimmed, "async"))

	return true, parseAsyncTokens(rem)
}

func parseAsyncTokens(rem string) AsyncStepOpts {
	var opts AsyncStepOpts
	rem = extractAsyncShell(rem, &opts)
	rem = extractAsyncInterval(rem, &opts)
	opts.Command = cleanAsyncCommand(rem)

	return opts
}

func extractAsyncShell(rem string, opts *AsyncStepOpts) string {
	shells := []string{"ps", "bash", "shell"}
	for _, s := range shells {
		hasShell := strings.HasPrefix(rem, s+" ") || rem == s
		if hasShell {
			opts.ShellType = s

			return strings.TrimSpace(strings.TrimPrefix(rem, s))
		}
	}

	return rem
}

func extractAsyncInterval(rem string, opts *AsyncStepOpts) string {
	parts := strings.Fields(rem)
	var retained []string
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if p == "-t" && i+1 < len(parts) {
			sec, _ := strconv.Atoi(parts[i+1])
			opts.IntervalSec = sec
			i++
			continue
		}
		if strings.HasPrefix(p, "-t=") {
			sec, _ := strconv.Atoi(strings.TrimPrefix(p, "-t="))
			opts.IntervalSec = sec
			continue
		}
		retained = append(retained, p)
	}

	return strings.Join(retained, " ")
}

func cleanAsyncCommand(cmd string) string {
	cleaned := strings.TrimSpace(cmd)
	if strings.HasPrefix(cleaned, "\"") && strings.HasSuffix(cleaned, "\"") && len(cleaned) >= 2 {
		return cleaned[1 : len(cleaned)-1]
	}

	if strings.HasPrefix(cleaned, "'") && strings.HasSuffix(cleaned, "'") && len(cleaned) >= 2 {
		return cleaned[1 : len(cleaned)-1]
	}

	return cleaned
}

func executeAsyncMacroStep(
	ctx context.Context,
	step MacroStep,
	asyncOpts AsyncStepOpts,
	currentDir string,
	start time.Time,
	execOpts ExecOptions,
	idx int,
) (StepExecution, error) {
	if execOpts.DryRun {
		return executeDryRunStep(step, idx, 1, execOpts, &DirTracker{CurrentDir: currentDir}), nil
	}

	if asyncOpts.IntervalSec > 0 {
		return executeAsyncPeriodicStep(step, asyncOpts, currentDir, start, execOpts)
	}

	return executeAsyncSingleStep(step, asyncOpts, currentDir, start, execOpts)
}

func executeAsyncPeriodicStep(step MacroStep, opts AsyncStepOpts, dir string, start time.Time, execOpts ExecOptions) (StepExecution, error) {
	go runAsyncPeriodicWorker(opts, dir)
	if !isStructuredOutput(execOpts) {
		fmt.Printf("  %s✔ periodic background task registered (interval: %ds): %s%s\n",
			constants.ColorGreen, opts.IntervalSec, opts.Command, constants.ColorReset)
	}

	return createAsyncStepExecution(step, opts.Command, dir, time.Since(start), "periodic_async"), nil
}

func executeAsyncSingleStep(step MacroStep, opts AsyncStepOpts, dir string, start time.Time, execOpts ExecOptions) (StepExecution, error) {
	cmd := buildAsyncStepCmd(opts, dir)
	if err := cmd.Start(); err != nil {
		return StepExecution{CommandLine: step.CommandLine, Status: "failed", Error: err.Error()}, err
	}

	if !isStructuredOutput(execOpts) {
		fmt.Printf("  %s✔ background task spawned (PID: %d): %s%s\n",
			constants.ColorGreen, cmd.Process.Pid, opts.Command, constants.ColorReset)
	}

	return createAsyncStepExecution(step, opts.Command, dir, time.Since(start), "async"), nil
}

func createAsyncStepExecution(step MacroStep, cmdText, dir string, elapsed time.Duration, status string) StepExecution {
	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    cmdText,
		WorkingDir:     dir,
		Status:         status,
		ExitCode:       0,
		ElapsedSeconds: elapsed.Seconds(),
		Logs:           []string{"spawned in background"},
		ErrorLogs:      []string{},
	}
}

func runAsyncPeriodicWorker(opts AsyncStepOpts, dir string) {
	ticker := time.NewTicker(time.Duration(opts.IntervalSec) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cmd := buildAsyncStepCmd(opts, dir)
		_ = cmd.Run()
	}
}

func buildAsyncStepCmd(opts AsyncStepOpts, dir string) *exec.Cmd {
	cmd := selectAsyncShellCmd(opts)
	if len(dir) > 0 {
		cmd.Dir = dir
	}

	return cmd
}

func selectAsyncShellCmd(opts AsyncStepOpts) *exec.Cmd {
	if opts.ShellType == "ps" {
		return exec.Command("powershell", "-NoProfile", "-Command", opts.Command)
	}

	if opts.ShellType == "bash" {
		return exec.Command("bash", "-c", opts.Command)
	}

	if runtime.GOOS == constants.OSWindows {
		return exec.Command("powershell", "-NoProfile", "-Command", opts.Command)
	}

	return exec.Command("sh", "-c", opts.Command)
}
