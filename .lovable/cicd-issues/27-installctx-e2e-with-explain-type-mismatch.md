# 27-installctx-e2e-with-explain-type-mismatch

## Error Summary

The CI/CD pipeline failed during the compile gate and test runs with:
```
cmd/installctx_darwin_e2e_test.go:170:23: cannot use runInstallCtxMac (value of type func() error) as func() value in argument to withExplain
cmd/installctx_linux_e2e_test.go:258:23: cannot use runInstallCtxLinux (value of type func() error) as func() value in argument to withExplain
```

## Root Cause Analysis

During an earlier command refactoring to eliminate `os.Exit(1)` and return typed errors, the return signature of `runInstallCtxMac()` and `runInstallCtxLinux()` was changed from `func()` to `func() error`. However, the helper `withExplain(t *testing.T, on bool, f func())` expects a parameter of type `func()`. Passing `runInstallCtxMac` and `runInstallCtxLinux` directly caused a type signature mismatch that only surfaces when compiling for `GOOS=darwin` or `GOOS=linux`.

## Solution Applied

Wrapped `runInstallCtxMac()` and `runInstallCtxLinux()` inside anonymous `func()` closures in `cmd/installctx_darwin_e2e_test.go` and `cmd/installctx_linux_e2e_test.go`. Verified cross-compilation with `GOOS=linux` and `GOOS=darwin` locally.

## What NOT to Repeat

- When changing the signature of a command or handler to return `error`, ALWAYS verify all test call sites—especially OS-specific test files guarded by `//go:build` tags.
