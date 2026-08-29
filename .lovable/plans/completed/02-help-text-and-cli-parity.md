# Execution Plan: Help Text & CLI Parity

## Goal

Implement missing CLI commands for `gitmap agy` (Antigravity) and `gitmap vscode`, fix their aliases, and massively improve the `gitmap help` formatting as specified by the user.

## Detailed Requirements

### 1. Help Text Formatting & Alignment

- Add empty lines (line gaps) underneath major headers in the help text (e.g., `Quick start:`, `Scanning & Discovery:`).
- Add the `installer` / `install` and `macro` sections into the main help text (under an `install` side section).
- Ensure the `schedule` command is properly documented in the main help text.
- Add `vscode` and `antigravity (ag)` to the main help text so users know they exist.
- Add explicit notes for commands that have multi-step/subcommands (e.g., `sj`, `ag`, `vscode`) so users know to expand them with `--help`.

### 2. Antigravity (`ag`) Commands

Add the following to `agy_cmd.go` (Cobra) or create dedicated handlers:
- `add-project` / `add` (supports multiple args / paths / IDs)
- `rm` / `remove` / `del` (supports commas / multiple arguments)
- `clear` (clears AG cache/projects)
- `open` (opens antigravity / specific project)
- `prompt` (sends a prompt to AG)
- `rw` (enables rewrite both for project)
- `sync` (loads/syncs all projects)
- `prompt-all-project` / `pap` (sends prompt to all)
- `export-projects` / `ep` (zip backup of AG projects)
- `import-projects` / `ip` (import zip backup)
- `stat` / `status`
- `plugins ls` (lists AG plugins)
- `plugin install` (installs an AG plugin)
- Add alias `ag` for `antigravity` command.

### 3. VS Code (`vscode`) Commands

Add the following to `vscode_cmd.go` (switch-based dispatch):
- `pap` (prompt all project - VS Code equivalent? No, user says "now similar needs to be available for vscode: gitmap vscode pap, plugins, add-project (ap), rm, ls")
- `plugins`
- `add-project` (ap)
- `rm` / `remove` / `del`
- `ls` / `list`

## Execution Strategy

1. **Phase 1 (Help Text Audit & Gap Fixes)**: Update `rootusage.go`, `rootusage_groups.go`, `rootusagecompact.go` with newline gaps, missing `ag`, `vscode`, `schedule`, `macro`, `installer`, and `sj` expand hints.
2. **Phase 2 (Antigravity Command Parity)**: Stub and implement the requested `ag` Cobra commands in `agy_cmd.go`. Use `apperror` correctly.
3. **Phase 3 (VS Code Command Parity)**: Extend `dispatchVSCodeAction` in `vscode_cmd.go` to handle the new subcommands and their aliases.
4. **Phase 4 (Validation)**: Run unit tests, verify boolean naming (`is`, `has`), ensure no generic variable names, limit function sizes.
5. **Phase 5 (Release Ceremony)**: Bump minor version in `version.json`, update `changelog.md`, commit with `feat(cli)`.

## Coding Guidelines Checklist (To Enforce)

- [x] No `temp`, `data`, `obj` variables.
- [x] Boolean prefixes: `is`, `has`, `can`, `should`. No inverted success variables.
- [x] Max 15 lines per function.
- [x] Strict error wrapping with `apperror`.
- [x] Fenced code blocks in markdown.
