---
name: powershell-cross-platform-scripting
description: >-
  Autonomously author, debug, and validate cross-platform PowerShell and Bash scripts with UTF-8 BOM encoding,
  exit code contracts, and execution wrappers.
---

# Cross-Platform PowerShell & Shell Scripting Skill

Autonomously author, review, and debug PowerShell (`.ps1`) and shell (`.sh`) scripts adhering to `spec/11-powershell-integration/`, `spec/17-consolidated-guidelines/10-powershell-integration.md`, and `.lovable/strictly-avoid.md`.

## Core Checkpoints & Mandatory Invariants

1. **Character Encoding & Windows Compatibility:**
   - PowerShell scripts that output emojis, box-drawing characters, or ANSI escape sequences MUST be saved with **UTF-8 with BOM** (`utf-8-sig`) so Windows PowerShell 5.1 and pwsh parse Unicode string literals accurately.
   - Cross-platform scripts must use universal line endings (`\n` or `\r\n` handled gracefully). Never use PowerShell literal newlines (`` `n ``) in non-PowerShell target files.

2. **No Shell Wrappers Around Python Scripts:**
   - Cross-platform Python tools must be invoked directly via `python script.py`.
   - Strictly avoid creating `.sh` or `.ps1` wrappers that merely forward arguments to Python, as argument escaping breaks in cross-platform GitHub Actions runners.

3. **Standardized Exit Codes:**
   - `0`: Success.
   - `1`: Operational or runtime error.
   - `2`: Invalid command-line arguments or syntax error.
   - Any script encountering an unrecoverable state must exit with a non-zero code to fail CI pipelines.

4. **Stdout vs. Stderr Separation:**
   - Diagnostic messages, progress indicators, banners, and errors must be written to `stderr` (e.g. `[Console]::Error.WriteLine(...)` or `write-error` / `echo >&2`).
   - Pure machine-readable output (JSON, paths, hashes) must go to `stdout`.

5. **Shell Directory Handoff Sentinel:**
   - Never mutate shell environments directly using `os.Setenv("GITMAP_SHELL_HANDOFF", ...)`.
   - Use the `GITMAP_HANDOFF_FILE` temporary file pattern (`$env:GITMAP_HANDOFF_FILE`) for shell function directory transitions (e.g. `gcd`, `cn`).
