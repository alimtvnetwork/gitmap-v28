# Subtask 93.04: Gitmap CLI Installer Decoupling & AGY Command Parity

## Goal
Decouple `agy` CLI installation into dedicated handlers in `gitmap/cmd/`, wiring `gitmap install agy`, `gitmap install antigravity-cli`, and `gitmap agy install`.

## Files Impacted
- `gitmap/cmd/install_agy.go` (NEW)
- `gitmap/cmd/install_agy_fetch.go` (NEW)
- `gitmap/cmd/agy_install.go`
- `gitmap/cmd/install_handlers.go`

## Acceptance Criteria
1. `runInstallAgyWithOpts` downloads and installs official `agy` CLI via `https://antigravity.google/cli/install.sh` / `install.ps1` or npm fallback.
2. `gitmap agy install` supports target selection: `gitmap agy install cli` installs the CLI; `gitmap agy install ide` (or `app`) installs the Desktop IDE.
3. Saves to `InstallationSplitDB` as tool `"agy"`.
4. Functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
