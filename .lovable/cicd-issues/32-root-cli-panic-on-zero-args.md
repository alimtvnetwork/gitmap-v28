# 32-root-cli-panic-on-zero-args

## Error Summary

Running `gitmap` with no arguments or unknown subcommands panicked with `panic: fatal error`:
```text
panic: fatal error

goroutine 1 [running]:
github.com/alimtvnetwork/gitmap-v28/gitmap/cmd.Run()
        /home/runner/work/gitmap-v28/gitmap-v28/gitmap/cmd/root.go:26 +0x892
main.main()
        /home/runner/work/gitmap-v28/gitmap-v28/gitmap/main.go:7 +0xf
```

---

## 4-Part Root Cause Analysis

### 1. Why it happened

In `gitmap/cmd/root.go`, lines 26 and 62 contained an explicit `panic("fatal error")` directly after printing binary locations and CLI usage:
```go
if len(os.Args) < 2 {
    PrintBinaryLocations()
    printUsage()
    panic("fatal error")
}
```
When `gitmap` was executed with zero arguments, `cmd.Run()` hit this branch and panicked instead of gracefully returning or exiting 0. Additionally, multiple command files had `panic("fatal error")` calls in error handling paths.

### 2. How it happened

During earlier refactoring that changed `Run()` from returning an error to `void`, placeholder `panic("fatal error")` calls were inserted instead of clean returns and proper `os.Exit(1)` or `cliexit.Exit(1)` error terminations.

### 3. Root Cause

Inadvertent `panic("fatal error")` calls in `gitmap/cmd/root.go:26`, `gitmap/cmd/root.go:62`, and throughout various command handlers instead of graceful `return` or `os.Exit(1)`.

### 4. Code Fix

- In `gitmap/cmd/root.go`, replaced `panic("fatal error")` with `return` when `len(os.Args) < 2` so running `gitmap` with no arguments displays usage and exits cleanly with code 0.
- Replaced all other occurrences of `panic("fatal error")` in `gitmap/cmd/` with graceful `os.Exit(1)` / `cliexit.Exit(1)`.
- Added unit test `TestRunNoArgsDoesNotPanic` in `gitmap/cmd/root_no_args_test.go`.
- Added `CLI Zero-Args Smoke` (`go run .`) to `.lovable/ai-fix-scripts/03-cicd-local-runner.py` to ensure zero-argument binary execution is verified during build-time CI gates.
