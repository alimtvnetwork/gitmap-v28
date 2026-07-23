# Issue 03 — `gitmap code` (VS Code open) command does not exist

**Status:** open
**Created:** 2026-06-07
**Related:** `.lovable/spec/commands/05-gitmap-code-opens-vscode.md`

## Symptom
`gitmap code`, `gitmap vcode`, `gitmap vscode` all fall through to "unknown command" instead of opening the current (or argument) folder in VS Code and registering it with the Project Manager extension.

## Expected
- `gitmap code` opens cwd in VS Code AND appends the project to `projects.json` (Project Manager extension), tagged with the repo slug.
- `gitmap code <folder>` opens `<folder>` instead.
- Aliases `vcode` / `vscode` behave identically.

## Actual
- No dispatch entry, no `runCode` handler, no `CmdCode*` constants.
- Only related code is `vscode-pm-sync` (`vpm`), which walks `projects.json` but does not open the editor.
