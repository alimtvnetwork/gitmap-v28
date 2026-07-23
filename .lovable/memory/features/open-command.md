---
name: gitmap open command
description: `gitmap open` (alias `op`) opens the current repo in BOTH GitHub Desktop and VS Code in one shot, re-injecting on every call. Resolves repo via git toplevel, falling back to cwd. v4.43.0+.
type: feature
---

# Feature: `gitmap open` / `op` (v4.43.0+)

## Rule
`gitmap open` MUST always launch BOTH GitHub Desktop AND VS Code on
the current repo, re-injecting on every call (no idempotency gate yet).

## Behavior
1. Resolve target = `git rev-parse --show-toplevel` || cwd.
2. Best-effort DB upsert (only when `origin` remote exists — same logic as `inject`).
3. `registerSingleDesktop(name, target)`.
4. `openInVSCode(target)`.
5. No shell handoff (user is already inside the folder).

## Why
User wanted a single command to surface the cwd repo in both tools
without typing the path. Inject-style flow but cwd-based and Desktop+VSCode-only.

## Files
- `gitmap/cmd/open.go` — entrypoint `runOpen`.
- `gitmap/constants/constants_open.go` — `CmdOpen`, `CmdOpenAlias`, `HelpOpen`, messages.
- `gitmap/cmd/rootcore.go` — dispatch entry.
- `gitmap/helptext/open.md` — user-facing help.

## Future
- `--force` / `-f` flag once idempotency tracking lands (LastInjectedDesktopAt / LastInjectedVSCodeAt columns).
- Spec: `spec/04-generic-cli/31-open.md` (TODO).
