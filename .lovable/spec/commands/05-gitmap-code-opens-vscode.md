# Command: `gitmap code` opens VS Code and registers project

**Status:** open
**Created:** 2026-06-07

## Verbatim
> If we do `gitmap code`, that would open the codebase to VS Code and also add this project to the project section. The command would be `gitmap code`. We can give a folder name; if we don't give any folder name and we are in the current folder, it will take the current folder. This command can have aliases like `vcode` or `vscode` — all of those will lead to VS Code open.

## Scope
- New top-level CLI command:
  - Primary: `gitmap code`
  - Aliases: `gitmap vcode`, `gitmap vscode`
- Arguments:
  - `gitmap code` (no arg) → opens `cwd` in VS Code.
  - `gitmap code <folder>` → opens `<folder>` (relative or absolute) in VS Code.
- Side effect:
  - Append/update the folder entry in the VS Code **Project Manager** extension's `projects.json` (reuse the merge logic already in `vscode-pm-sync` / `vpm`, UNION + skip missing rootPaths).
  - Tag the project with the repo's slug (same convention as `vpm`).
- Failure modes:
  - VS Code not installed → soft-fail with stderr message; do NOT register in projects.json.
  - Target folder missing → exit non-zero with explicit error.

## When it applies
Whenever the user invokes any of the three command IDs above. Windows-first, with macOS + Linux parity (use `code` on PATH; on Windows fall back to `code.cmd` under `%LOCALAPPDATA%\Programs\Microsoft VS Code\bin\`).
