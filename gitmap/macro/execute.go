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

// Execute runs all steps in a macro.
func Execute(ctx context.Context, m *Macro, opts ExecOptions) error {
	fmt.Printf("  %s▶ Executing Macro: %q (%d steps)%s\n\n",
		constants.ColorCyan, m.Name, len(m.Steps), constants.ColorReset)

	start := time.Now()
	for i, step := range m.Steps {
		if opts.DryRun {
			fmt.Printf("  [%2d/%d] ➜ (dry-run) %s\n", i+1, len(m.Steps), step.CommandLine)
			continue
		}
		if err := executeStep(ctx, step, i+1, len(m.Steps), opts); err != nil {
			if !step.ContinueOnError {
				fmt.Printf("  %s✖ Step %d failed: %v%s\n", constants.ColorRed, i+1, err, constants.ColorReset)
				return err
			}
			fmt.Printf("  %s▲ Step %d warning (ignored): %v%s\n", constants.ColorYellow, i+1, err, constants.ColorReset)
		}
	}

	elapsed := time.Since(start)
	fmt.Println()
	fmt.Printf("  %s✔ Macro %q completed (%d steps) · Elapsed: %.1fs%s\n\n",
		constants.ColorGreen, m.Name, len(m.Steps), elapsed.Seconds(), constants.ColorReset)
	return nil
}

func executeStep(ctx context.Context, step MacroStep, idx, total int, opts ExecOptions) error {
	fmt.Printf("  [%2d/%d] ➜ %s ... ", idx, total, step.CommandLine)
	start := time.Now()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", step.CommandLine)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", step.CommandLine)
	}

	if len(step.WorkingDir) > 0 {
		cmd.Dir = step.WorkingDir
	}
	if opts.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("%s✖ failed (%.1fs)%s\n", constants.ColorRed, elapsed.Seconds(), constants.ColorReset)
		return err
	}

	fmt.Printf("%s✔ ok (%.1fs)%s\n", constants.ColorGreen, elapsed.Seconds(), constants.ColorReset)
	return nil
}
