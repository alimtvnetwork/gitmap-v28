# CI/CD Issue 13: `undefined: swapIP` Cross-Platform Build & Vet Failure

## Context

Running `GOOS=windows GOARCH=amd64 go vet ./...` (or compiling on macOS / `GOOS=darwin` matrix jobs) failed with:
```

# github.com/alimtvnetwork/gitmap-v28/gitmap/cmd

Error: cmd/ipchange_cmd.go:47:12: undefined: swapIP
Error: cmd/ipchange_cmd.go:57:8: undefined: swapIP

# github.com/alimtvnetwork/gitmap-v28/gitmap/cmd

# [github.com/alimtvnetwork/gitmap-v28/gitmap/cmd]

Error: vet: cmd/ipchange_cmd.go:47:12: undefined: swapIP
Error: Process completed with exit code 1.
```

## Root Cause Analysis

1. **Implicit Go Filename Suffix Constraints**: Go toolchain conventions (`go/build`) automatically treat file suffixes `*_$GOOS.go` (such as `*_linux.go`, `*_windows.go`) as hard compiler build constraints regardless of internal build tags.
2. **The Missing POSIX Declaration Gap**:
   - `ipchange_linux.go` contained the POSIX implementation of `swapIP` with `//go:build !windows`.
   - However, because the filename ended in `_linux.go`, the Go toolchain restricted compilation strictly to `GOOS=linux`.
   - On macOS (`darwin`), FreeBSD, or during cross-compile matrix checking, `ipchange_windows.go` was ignored (`GOOS != windows`) and `ipchange_linux.go` was also ignored (`GOOS != linux`).
   - Consequently, neither file was compiled, leaving `swapIP` completely undeclared in package `cmd` and triggering compilation/vet failures.

## Resolution

1. **Renamed File**: Renamed `gitmap/cmd/ipchange_linux.go` to `gitmap/cmd/ipchange_posix.go` with `//go:build !windows` header.
2. **Verified Multi-OS Declarations**: Confirmed `ipchange_windows.go` provides `swapIP` on Windows, and `ipchange_posix.go` provides `swapIP` across all non-Windows systems (`linux`, `darwin`, `freebsd`, `openbsd`).
3. **Normalized Installer OS Suffixes**: Renamed `prompt_runner_unix.go` and `prompt_runner_windows.go` to `prompt_runner_unix_impl.go` and `prompt_runner_win_impl.go` to ensure clean multi-platform compilation.

## What NOT to Repeat

- **Never rely on `//go:build !windows` inside a file named `*_linux.go`**: Go ignores `//go:build !windows` when the filename has `_linux.go` because the filename suffix filter takes precedence and rejects `darwin`, `freebsd`, etc.
- Always use `_posix.go` or `_unix.go` with `//go:build !windows` when implementing fallback logic for all non-Windows operating systems.
