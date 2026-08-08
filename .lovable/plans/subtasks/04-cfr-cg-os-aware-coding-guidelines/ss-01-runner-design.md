---
Slug: runner-design
Parent: 04-cfr-cg-os-aware-coding-guidelines
Status: pending
Created: 2026-07-16
---

# SS-01 — Coding Guidelines runner design

## Goal

A single `RunCodingGuidelinesInstall(opts CodingGuidelinesOpts) error` entry point used by `cfr cg`, `cfrp cg`, and any future integration.

## Shape

```go
type CodingGuidelinesOpts struct {
    WorkingDir string          // repo root after clone
    Runner     func(name string, args ...string) *exec.Cmd // injectable for tests
    Stdout     io.Writer
    Stderr     io.Writer
    Stdin      io.Reader
}
```

## Dispatch rules

- `runtime.GOOS == "windows"` → resolve `powershell` then `pwsh` (reuse `resolvePowerShellBinary` from `installcleancode.go`; extract into a shared helper file `gitmap/cmd/pwshresolve.go` to avoid duplication). Command:
  `powershell -NoProfile -ExecutionPolicy Bypass -Command "irm <WindowsURL> | iex"`
- Otherwise → require `bash` and `curl` on PATH; error early with actionable hint if missing. Command:
  `bash -c "curl -fsSL <UnixURL> | bash"`

## URLs

Sourced from `constants_codingguidelines.go`:

- `DefaultCodingGuidelinesURLWindows = "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/error-manage-install.ps1"`
- `DefaultCodingGuidelinesURLUnix    = "https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh"`

## Errors

- Missing shell → `constants.ErrCGShellNotFound` with the exact one-liner fallback the user can copy.
- Installer non-zero exit → wrap with `%w` and surface OS + URL used for diagnosis.
- Zero-swallow rule (per project core memory): every failure goes to `os.Stderr` in the standardized format.

## Function-size cap

Split into `dispatchWindows`, `dispatchUnix`, `runInstaller` so each stays under the 15-line cap.
