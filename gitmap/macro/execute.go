package macro

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// Execute runs all steps in a macro maintaining dynamic directory state.
func Execute(ctx context.Context, m *Macro, opts ExecOptions) error {
	fmt.Printf("  %s▶ Executing Macro: %q (%d steps)%s\n\n",
		constants.ColorCyan, m.Name, len(m.Steps), constants.ColorReset)
	start := time.Now()
	initialDir, _ := os.Getwd()
	dt := NewDirTracker(initialDir)
	return runExecuteSteps(ctx, m, opts, dt, start)
}

func runExecuteSteps(ctx context.Context, m *Macro, opts ExecOptions, dt *DirTracker, start time.Time) error {
	for i, step := range m.Steps {
		if opts.DryRun {
			fmt.Printf("  [%2d/%d] ➜ (dry-run) %s\n", i+1, len(m.Steps), step.CommandLine)
			continue
		}
		err := executeStep(ctx, step, i+1, len(m.Steps), opts, dt)
		if err == nil {
			continue
		}
		if !step.ContinueOnError {
			fmt.Printf("  %s✖ Step %d failed: %v%s\n", constants.ColorRed, i+1, err, constants.ColorReset)
			return err
		}
		fmt.Printf("  %s▲ Step %d warning (ignored): %v%s\n", constants.ColorYellow, i+1, err, constants.ColorReset)
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

func executeStep(ctx context.Context, step MacroStep, idx, total int, opts ExecOptions, dt *DirTracker) error {
	expandedCmd := ExpandPathAndEnv(step.CommandLine)
	fmt.Printf("  [%2d/%d] ➜ %s ... ", idx, total, expandedCmd)
	start := time.Now()
	if isDirChange := dt.ProcessCd(expandedCmd); isDirChange {
		elapsed := time.Since(start)
		fmt.Printf("%s✔ ok (%.1fs)%s\n", constants.ColorGreen, elapsed.Seconds(), constants.ColorReset)
		return nil
	}
	return runStepProcess(ctx, expandedCmd, step, opts, dt, start)
}

func runStepProcess(ctx context.Context, cmdText string, step MacroStep, opts ExecOptions, dt *DirTracker, start time.Time) error {
	targetDir := dt.CurrentDir
	if len(targetDir) == 0 {
		targetDir = ExpandPathAndEnv(step.WorkingDir)
	}
	cmd := buildStepCmd(ctx, cmdText, targetDir, opts)
	err := cmd.Run()
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("%s✖ failed (%.1fs)%s\n", constants.ColorRed, elapsed.Seconds(), constants.ColorReset)
		return err
	}
	fmt.Printf("%s✔ ok (%.1fs)%s\n", constants.ColorGreen, elapsed.Seconds(), constants.ColorReset)
	return nil
}

func buildStepCmd(ctx context.Context, cmdText, dir string, opts ExecOptions) *exec.Cmd {
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
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
