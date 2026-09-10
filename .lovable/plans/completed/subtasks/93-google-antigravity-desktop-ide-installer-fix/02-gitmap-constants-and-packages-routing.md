# Subtask 93.02: Gitmap Constants, Tool Mappings & Binary Probes

## Goal
In `gitmap`, introduce `constants.ToolAgy = "agy"`, update `ToolAntigravity` description, and configure package aliases and probes to distinguish between the Desktop IDE (`antigravity`) and the CLI (`agy`).

## Files Impacted
- `gitmap/constants/constants_install.go`
- `gitmap/cmd/install_packages.go`
- `gitmap/cmd/installprobe.go`
- `gitmap/cmd/installverify.go`
- `gitmap/cmd/install_handlers.go`

## Acceptance Criteria
1. `ToolAgy = "agy"` added to `constants_install.go`.
2. `ToolAntigravity` points to binary `antigravity` (Desktop IDE).
3. `ToolAgy` points to binary `agy` (CLI tool).
4. `install_packages.go` maps:
   - `"antigravity"`, `"antigravity-ide"`, `"antigravity-desktop"` -> `ToolAntigravity`
   - `"agy"`, `"antigravity-cli"`, `"ag"` -> `ToolAgy`
5. Functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
