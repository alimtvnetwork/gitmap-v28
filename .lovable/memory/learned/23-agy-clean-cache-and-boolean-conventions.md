# 23: Antigravity Clean-Cache Implementation & Boolean Linter Conventions

## Context
Implemented the `gitmap agy clean-cache` subcommand in Golang within `cli/cmdagy` to replace manual PowerShell scripts with an advanced, cross-platform cache cleaner and process terminator.

## 1. Antigravity Cache Architecture
The cleanup targets 10 cache locations across Windows and Unix:
- **Chromium Caches**: `Cache`, `Code Cache`, `GPUCache`, `blob_storage`, `logs`.
- **Dawn WebGPU Shaders**: Glob matching `Dawn*Cache` (`DawnGraphiteCache`, `DawnWebGPUCache`).
- **Updater Packages**: `%LOCALAPPDATA%\antigravity-updater` (or `~/.config/antigravity-updater` / `~/Library/Caches/antigravity-updater`).
- **Crash Dumps & Scratch**: `~/.gemini/antigravity/crashes` and `~/.gemini/antigravity/scratch`.
- **Temp**: System/user temporary files (`%TEMP%` / `$TMPDIR`).

### Data Protection Boundary
Strictly preserved without touching:
- User projects and pinned projects: `~/.gemini/config/projects/`
- User conversation history and transcripts: `~/.gemini/antigravity/brain/`
- User configurations: `Preferences`, `Local State`, `app_storage.json`
- Skills and MCP configurations: `~/.gemini/antigravity/skills/`, `~/.gemini/antigravity/mcp/`

## 2. Cross-Platform Process Discovery & Termination
- **Windows**:
  - Uses `tasklist /FO CSV /NH` for instant (<50ms) CSV parsing.
  - Terminates matching PIDs (`antigravity`, `electron`, `msedge`, `msedgewebview2`) via `taskkill /F /PID <pid>`.
  - Excludes the current runner PID (`os.Getpid()`).
- **Unix**:
  - Uses `ps -eo pid,comm` parsing.
  - Terminates matching PIDs via `syscall.Kill(pid, syscall.SIGKILL)`.

## 3. Strict Boolean Conventions (Linter Enforced)
- Banned boolean function prefixes: `can*` and `should*` are strictly banned by `linter-scripts/check-boolean-guidelines.py`.
- Allowed boolean function prefixes: `is*` and `has*` only (e.g. `isConfirmed`, `isTerminationRequired`).
- Mixed polarity (`if a && !b`) is flagged; split into explicit conditions or positive framing.
