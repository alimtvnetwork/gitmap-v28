# Subtask 01: Interactive Installer Creation Engine & Dual CLI Parity

## 1. Description
Implement interactive and scripted creation for custom installers in `gitmap install add <name> [version]` and `gitmap installer add <name> [version]`. Support macro-style interactive questions:
- Description of the installer
- Windows instructions (win)
- Unix instructions (unix)
- Specific OS instructions (ubuntu)
- Prompt: "Do you want to add or edit Unix or Ubuntu instructions later? (y/n)"
Store the record in SQLite with `model.InstallerScript` and `Scripts` map.
Also support non-interactive flags (`--desc`, `--win`, `--unix`, `--ubuntu`, `--yes`).

## 2. Files to Modify / Create
- `gitmap/cmd/install_add.go`: Implement `runInstallAdd(args []string) error` with flag parsing, terminal check, interactive prompt reader, and DB persistence.
- `gitmap/cmd/install.go`: Route `add` subcommand in `runInstall`.
- `gitmap/cmd/installer.go`: Add `add` subcommand pointing to the same creation engine as `create`.
- `gitmap/cmd/install_add_test.go`: Unit tests for interactive mock and flag-based creation.

## 3. Invariants
- Functions <= 15 lines.
- Blank line before every return.
- Zero nested if statements.
- Affirmative booleans only.
- Strict relative paths.
