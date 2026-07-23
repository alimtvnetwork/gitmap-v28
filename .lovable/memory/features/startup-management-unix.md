---
name: startup-management-unix
description: Linux/Unix XDG autostart list/remove with marker-based scoping to gitmap-managed entries
type: feature
---

# Startup Management (Linux/Unix)

`gitmap startup-list` (alias `sl`) and `gitmap startup-remove` (alias
`sr`) enumerate and delete XDG autostart entries created by gitmap.

## Scoping rule (CRITICAL — never widen)

An entry is "gitmap-managed" if and only if:
1. Filename ends in `.desktop` AND starts with `gitmap-`
2. Body contains `X-Gitmap-Managed=true` (constant: `StartupMarkerKey`)

Both gates required. List and Remove BOTH re-check the marker — never
trust filename alone. A future `startup-add` MUST emit the marker.

## Directory resolution

`AutostartDir()` returns:
- `$XDG_CONFIG_HOME/autostart` if `XDG_CONFIG_HOME` is set
- otherwise `$HOME/.config/autostart`
- error on Windows / macOS (use platform-specific commands there)

## Safety contract

- Missing autostart dir → empty list, exit 0 (NOT an error)
- Unknown name in remove → `RemoveNoOp`, exit 0
- Existing third-party file → `RemoveRefused`, exit 0, file untouched
- Path-separator / NUL in name → `RemoveBadName`, exit 0, no I/O
- Real I/O failure → exit 1

All "soft" outcomes exit 0 so the commands are idempotent and
script-safe.

## File layout

- `gitmap/startup/startup.go` — public API: `Entry`, `AutostartDir`, `List`
- `gitmap/startup/desktop.go` — .desktop parser + manage filter
- `gitmap/startup/remove.go` — `Remove` + `RemoveStatus` enum
- `gitmap/startup/scanner.go` — bufio shim
- `gitmap/cmd/startup.go` — CLI runners
- `gitmap/constants/constants_startup.go` — Cmd*, marker, messages

Shipped in v3.133.0.
