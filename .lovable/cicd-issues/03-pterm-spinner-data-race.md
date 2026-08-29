# CI/CD Issue: pterm Spinner Data Race

## Error

```
WARNING: DATA RACE
Read at 0x00c0003fc098 by goroutine 349:
  github.com/pterm/pterm.SpinnerPrinter.Start.func1()

Previous write at 0x00c0003fc098 by goroutine 348:
  github.com/pterm/pterm.(*SpinnerPrinter).Stop()
  github.com/pterm/pterm.(*SpinnerPrinter).Success()
  github.com/alimtvnetwork/gitmap-v28/gitmap/cluster.RunPool.func3()
```

## Root Cause

The `github.com/pterm/pterm@v0.12.83` package has a known internal data race inside the `SpinnerPrinter`. When `Start()` is called, it unconditionally spawns a goroutine that continuously reads the `IsActive` boolean field. When the main goroutine calls `Stop()`, `Success()`, or `Fail()`, it writes `IsActive = false` without a mutex. 

In `gitmap/cluster/pool.go`, the `pterm.DefaultSpinner` was unconditionally initialized and started regardless of whether the output was active (`pterm.Output`), meaning even during `go test -race` runs (which disable output via `pterm.DisableOutput()`), the background goroutine was spawned and the data race was triggered when `spinner.Success()` was called.

## Solution

Modified `gitmap/cluster/pool.go` to conditionally instantiate and start the `SpinnerPrinter` only if `isMultiActive` is `true`. When output is disabled (e.g. during tests), the application now falls back to simple `pterm.Info.Printf`, `pterm.Success.Printf`, etc., entirely bypassing the `SpinnerPrinter` allocation and its racing goroutines.

## What NOT to Repeat

- **Unconditional UI Component Initialization:** Do not unconditionally start `pterm` background UI components (like Spinners) without wrapping them in an active-output check (e.g. `pterm.Output` or equivalent). This will leak goroutines and trigger data races in test environments where output is discarded or disabled.
