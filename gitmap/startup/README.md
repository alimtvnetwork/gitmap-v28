# `gitmap/startup` — build tags & cross-OS strategy

This package manages OS-level autostart entries for the `gitmap`
binary. It supports three backends across three operating systems
from a **single import path** that the rest of the codebase calls
without ever switching on `runtime.GOOS` itself.

| OS      | Backend                          | Primary files                              |
|---------|----------------------------------|--------------------------------------------|
| Linux   | XDG `.desktop` autostart entry   | `startup.go`, `desktop.go`                 |
| macOS   | `~/Library/LaunchAgents/*.plist` | `addplist.go`, `plist.go`, `plistxml.go`   |
| Windows | Registry `Run` keys + `.lnk`     | `winbackend.go`, `winregistry_*.go`, `winshortcut*.go` |

---

## File-suffix convention

Go applies build tags **automatically** based on filename suffix:

| Suffix          | Compiled when                  |
|-----------------|--------------------------------|
| `*_windows.go`  | `GOOS=windows` only            |
| `*_darwin.go`   | `GOOS=darwin` only             |
| `*_linux.go`    | `GOOS=linux` only              |
| `*_other.go`    | **Has no implicit meaning** — must declare `//go:build !<os>` explicitly |
| (no suffix)     | Always compiled, every OS      |

The implicit-suffix rule is why `winregistry_windows.go` does not need
a `//go:build windows` directive — but **we add one anyway** for human
readability. The `_other.go` suffix is a project convention only; the
`//go:build !windows` line is what actually excludes it from Windows
builds.

---

## The three file classes in this package

### 1. Untagged dispatchers — `winbackend.go`, `winshortcut.go`

These files compile on **every OS** and contain code that *decides
which backend to use*. They reference Windows-only symbols (e.g.
`addWindowsRegistry`, `trackingSubkeyExists`) but guard the calls
with `runtime.GOOS == "windows"` so the non-Windows branch never
executes them.

For the Go compiler to accept these references on Linux/macOS, the
symbols must still **exist** at compile time — which is what the
stub files in class 3 provide.

### 2. Per-OS implementations — `*_windows.go`, `*_darwin.go`

Real implementations gated by filename suffix + explicit `//go:build`
tag. These hold the OS-specific syscalls, registry access, plist
marshaling, etc.

### 3. Stub files — `winregistry_other.go`

`//go:build !windows` files that provide **no-op or error-returning
implementations** of every Windows-only symbol referenced from the
untagged dispatchers in class 1.

Without these stubs, `go build` on Linux or macOS fails with
`undefined: addWindowsRegistry`. The stubs are not dead code — they
are the load-bearing piece that lets a cross-OS dispatcher reference
platform-gated symbols.

---

## Why not just tag the dispatchers `//go:build windows`?

Two reasons:

1. **Shared logic.** `winbackend.go` contains the *selection* logic
   (registry vs. shortcut, HKCU vs. HKLM, scope parsing) that the
   non-Windows test suite exercises via table-driven tests. Tagging
   the file Windows-only would make those tests un-runnable on the
   Linux CI matrix.

2. **Single dispatch surface.** Callers do `startup.Add(...)` once;
   the package internally routes to the right backend. Splitting the
   dispatcher per-OS would push that branching into every caller.

---

## Adding a new Windows-only symbol — checklist

If you add a new function that lives in a `*_windows.go` file and is
called from an untagged file, you **must** also add a stub in
`winregistry_other.go` (or a sibling `_other.go` file). The stub
signature must match exactly.

CI will catch a missing stub via `.github/workflows/startup-build-tags.yml`,
which cross-compiles this package for `GOOS=linux`, `darwin`, and
`windows` on every PR that touches `gitmap/startup/**`. A missing
stub fails the Linux/macOS leg with `undefined: <symbol>`.

Local pre-flight:

```bash
cd gitmap
GOOS=linux   go build ./startup/...
GOOS=darwin  go build ./startup/...
GOOS=windows go build ./startup/...
```

All three must succeed.

---

## Test helpers

OS-gated test helpers (e.g. `withFakeLaunchAgentsDir` in
`launchagents_testhelper_test.go`) follow the same rules — the helper
lives in **one** `*_test.go` file, gated by build tag if it uses
OS-specific paths, and is consumed by sibling test files without
redeclaration. See the per-package test files for the established
pattern.

---

## Related

- CI guard: [`.github/workflows/startup-build-tags.yml`](../../.github/workflows/startup-build-tags.yml)
- Cross-platform spec: [`spec/01-app/42-cross-platform.md`](../../spec/01-app/42-cross-platform.md)
- Constants strategy: [`spec/04-generic-cli/02-project-structure.md`](../../spec/04-generic-cli/02-project-structure.md)