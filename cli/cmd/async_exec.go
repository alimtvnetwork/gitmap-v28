package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runAsyncJob(args []string) error {
	opts, err := parseAsyncArgs(args)
	if err != nil {
		return err
	}

	fmt.Printf("%s Starting async monitor: %q", constants.ColorCyan+"ℹ"+constants.ColorReset, opts.command)
	if opts.intervalS > 0 {
		fmt.Printf(" (interval: %ds)\n", opts.intervalS)

		return runAsyncPeriodicLoop(opts)
	}

	fmt.Println(" (single execution)")

	return runAsyncSingleJob(opts)
}

func runAsyncPeriodicLoop(opts *asyncTaskOpts) error {
	ticker := time.NewTicker(time.Duration(opts.intervalS) * time.Second)
	defer ticker.Stop()

	iteration := 1
	executeAsyncIteration(opts.command, iteration)

	for range ticker.C {
		iteration++
		executeAsyncIteration(opts.command, iteration)
		if opts.maxCount > 0 && iteration >= opts.maxCount {
			break
		}
	}

	return nil
}

func executeAsyncIteration(command string, iter int) {
	tStr := time.Now().Format("15:04:05")
	fmt.Printf("\n[%s] #%d ➜ %s\n", tStr, iter, command)

	cmd := buildAsyncExecCmd(context.Background(), command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	printIterationResult(err, time.Since(start))
}

func printIterationResult(err error, duration time.Duration) {
	if err != nil {
		fmt.Printf("  %s✗ exit error: %v (%v)%s\n", constants.ColorRed, err, duration, constants.ColorReset)

		return
	}

	fmt.Printf("  %s✓ ok (%v)%s\n", constants.ColorGreen, duration, constants.ColorReset)
}

func runAsyncSingleJob(opts *asyncTaskOpts) error {
	cmd := buildAsyncExecCmd(context.Background(), opts.command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, "async.Start")
	}

	fmt.Printf("%s Process spawned in background (PID: %d)\n",
		constants.ColorGreen+"✓"+constants.ColorReset, cmd.Process.Pid)

	return nil
}

func buildAsyncExecCmd(ctx context.Context, command string) *exec.Cmd {
	if constants.OSWindows == "windows" {
		return exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", command)
	}

	return exec.CommandContext(ctx, "sh", "-c", command)
}
