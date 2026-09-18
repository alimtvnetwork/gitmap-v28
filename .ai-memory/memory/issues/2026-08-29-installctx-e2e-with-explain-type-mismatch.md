# Issue: installctx E2E withExplain Type Signature Mismatch Under Linux/Darwin

- date: 2026-08-29
- status: resolved
- scope: test/ci

## 1. Why it happened

When refactoring CLI functions to return errors instead of calling `os.Exit(1)`, the signatures of `runInstallCtxMac()` and `runInstallCtxLinux()` were updated from `func()` to `func() error`.

## 2. How it happened

In `cmd/installctx_darwin_e2e_test.go` and `cmd/installctx_linux_e2e_test.go`, the tests invoked `withExplain(t, true, runInstallCtxMac)` and `withExplain(t, true, runInstallCtxLinux)`. The helper function `withExplain` expects an argument of type `func()`. In Linux and Darwin test runs (including CI's compile gate `go test -run=^$ ./...`), the compiler rejected the function argument as a type mismatch.

## 3. Root Cause

- Passing `func() error` directly to a parameter requiring `func()` at `cmd/installctx_darwin_e2e_test.go:170` and `cmd/installctx_linux_e2e_test.go:258`.
- These tests are conditionally built using `//go:build darwin || linux`, which prevented native Windows test commands from catching the regression locally without cross-compilation checks.

## 4. Code Fix

- Wrapped calls in anonymous parameterless closures:
  ```go
  withExplain(t, true, func() {
      _ = runInstallCtxMac()
  })
  ```
  and
  ```go
  withExplain(t, true, func() {
      _ = runInstallCtxLinux()
  })
  ```
- Verified cross-compilation locally using `GOOS=linux go test -c -o NUL ./cmd` and `GOOS=darwin go test -c -o NUL ./cmd`.
