# Subtask 03: Installer Script Parity (`install.ps1`, `install.sh`)

## Scope & Intent
- Audit `install.ps1` and `install.sh` for profile updates and completion setup.
- If multiple profiles are configured or added, format them as a subtree under a single header or ensure consistency with the Go CLI.
- Ensure a clean trailing newline gap before returning control to the terminal prompt.
- Keep copies in sync: `install.ps1`, `cli/scripts/install.ps1`, `cli/fspath/install.ps1`, `install.sh`, `cli/scripts/install.sh`, `cli/fspath/install.sh`.

## Owned Files
- `install.ps1`
- `install.sh`
- `cli/scripts/install.ps1`
- `cli/scripts/install.sh`
- `cli/fspath/install.ps1`
- `cli/fspath/install.sh`

## Invariants
- Zero builds or full test suites.
- Strict relative paths.
