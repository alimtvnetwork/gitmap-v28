---
Slug: gitmap-code-command
Parent: 03-reclone-transport-and-vscode-open
Status: pending
Created: 2026-06-07
---

# Subtask 04 — `gitmap code` opens VS Code and registers project

## Goal
New top-level command + two aliases that open a folder in VS Code AND append/update the entry in the Project Manager extension's `projects.json`.

## CLI surface
- Primary: `gitmap code [folder]`
- Aliases: `gitmap vcode [folder]`, `gitmap vscode [folder]`
- No arg → cwd. Arg → relative or absolute folder; error if missing.

## Constants (all in `constants_cli.go` per Core rule)
```go
CmdCode       = "code"
CmdCodeAlias1 = "vcode"
CmdCodeAlias2 = "vscode"
HelpCode      = "  code (vcode, vscode)  Open folder in VS Code and register with Project Manager"
```

## Files
- `gitmap/cmd/code.go` — `runCode(args []string)` entrypoint with `// gitmap:cmd top-level` marker (per `mem://features/marker-comments`).
- Dispatch in `gitmap/cmd/roottooling.go` (or wherever `vpm` lives — co-locate).
- Reuse `gitmap/vscode/pm.go` helpers from `vpm` for the `projects.json` merge — extract a `pm.AppendOrUpdate(path, slug)` helper if it does not already exist.

## VS Code binary resolution
Mirror the `desktop.ResolveCLI` shape from v6.24.0:
1. `exec.LookPath("code")`.
2. Windows fallback: `%LOCALAPPDATA%\Programs\Microsoft VS Code\bin\code.cmd` and `%PROGRAMFILES%\Microsoft VS Code\bin\code.cmd`.
3. macOS fallback: `/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code`.
4. Linux: PATH only.
5. Empty string → soft-fail with stderr `VS Code CLI not found — skipping. (looked in PATH and known install dirs)`.

## projects.json merge
- Read existing JSON (create empty array if missing).
- If an entry with the same `rootPath` exists, leave it alone (idempotent, matches `vpm` UNION rule).
- Else append `{ "name": "<basename>", "rootPath": "<abs>", "tags": ["<slug>"], "enabled": true }`.
- Write back with the same formatting `vpm` uses.

## Tests
- Binary resolver: temp dir + fake `code.cmd`, PATH cleared → returns the fake (Windows-only test).
- projects.json merge: empty file → one entry; existing entry with same rootPath → no duplicate; existing entry with different rootPath → both present.
- Soft-fail: missing VS Code → stderr line, exit 0 (do NOT fail the command — opening editor is best-effort, but registration should still happen if `projects.json` location is known? — actually, if no editor we still register, since `vpm` works without VS Code installed).

## Definition of done
- All three command IDs dispatch to `runCode`.
- Help system shows the entry under the Tooling group.
- Completion generator picks up the marker comment (CI `generate-check` passes).
