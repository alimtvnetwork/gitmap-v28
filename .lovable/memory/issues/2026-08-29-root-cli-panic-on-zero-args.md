# Memory Issue: 2026-08-29 Root CLI Panic on Zero Args

## 1. Why it happened

`gitmap/cmd/root.go` had `panic("fatal error")` in `len(os.Args) < 2` branch, causing `gitmap` without subcommands to panic instead of displaying usage.

---

## 2. How it happened

A prior refactor converted `Run()` into a non-returning function and placed temporary `panic("fatal error")` calls across `gitmap/cmd/*.go`.

---

## 3. Root Cause

`panic("fatal error")` in `cmd/root.go` and various command handlers.

---

## 4. Code Fix

Replaced `panic` in `cmd/root.go` with `return`, replaced panics across `cmd/*.go` with `os.Exit(1)`, added `TestRunNoArgsDoesNotPanic`, and added build-time `CLI Zero-Args Smoke` check to the local runner.
