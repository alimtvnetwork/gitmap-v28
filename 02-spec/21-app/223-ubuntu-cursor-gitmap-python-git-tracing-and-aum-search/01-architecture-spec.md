# Architecture Specification: Cursor Ubuntu Fleet Setup, repo-secrets Tracking, and GitMap Python Runner

> **Task Slug:** `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search`  
> **Document:** `01-architecture-spec.md`  
> **Version:** `1.0.0`  
> **Status:** Canonical Architecture Specification  
> **Target Subsystems:** Cursor Fleet Provisioning (`repo-secrets/`, `cli/cmdcursor/`), GitMap Native Python Runner (`cli/cmdpy/`, `cli/cmd/`), AI & Command Telemetry (`cli/store/`)

---

## 1. User Request (Verbatim)

```text
have you improved and added cursor to the ubuntu system using gitmap or python or shell script for now and keep track in the repo-secrets folder

try to test and make sure by testing properly that these works, find issues if there is any fix it and make a minor bump and check gitmap pe -t until ci cd is green
python codes needs to run from gitmap please ensure that

make sure we have these git commands from gitmap because it can also keep a trace to split db for heat map understanding?? Clear???

make sure the agent calls are optimized using gitmap https://prnt.sc/Rk0BXJBXgkeD
https://prnt.sc/C0T16jJtXhE_
https://prnt.sc/2k7l4_Q0VFyM

Make sure we have the se AUM searches verified properly in other or current repo against other search tools to see if we have any mistakes and can be improved anything??

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Improve and add cursor to the Ubuntu system using GitMap, Python, or shell script and keep track in the repo-secrets folder.
5. Test thoroughly to ensure functionality, find and fix any issues, make a minor version bump, and check `gitmap pe -t` until CI/CD is green.
6. Ensure Python code runs from GitMap.
7. Ensure Git commands from GitMap can trace to split db for heat map understanding.
8. Optimize agent calls using GitMap.
```

---

## 2. Executive Summary & Scope

This specification establishes the architectural blueprints for two foundational pillars of Task 223:

1. **Automated Cursor IDE Ubuntu Fleet Deployment & Secrets Governance (`REQ-01`)**:
   - Programmatic upstream querying against the official Cursor release API (`https://www.cursor.com/api/download?platform=linux-x64&releaseTrack=stable`) to dynamically fetch direct AppImage download URLs, versions, and checksum hashes without brittle manual URL hardcoding.
   - Dual-language automated provisioning engine via `repo-secrets/05-scripts/setup-cursor-ubuntu.py` and `repo-secrets/05-scripts/setup-cursor-ubuntu.sh`.
   - Hermetic system installation: AppImage target at `/opt/cursor/Cursor.AppImage` (0755), universal wrapper binary at `/usr/local/bin/cursor` and `$HOME/.local/bin/cursor` with `--no-sandbox` flags, and standard XDG desktop entry at `/usr/share/applications/cursor.desktop`.
   - Fleet theme and configuration invariance: Automatic injection of Dracula Dark theme and whitespace/formatting invariants into `$HOME/.config/Cursor/User/settings.json`.
   - Fleet status audit ledger: Structured persistence of node health, version, binary paths, and verification timestamps in `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`.
   - CLI integration: Modernizing `cli/cmdcursor/cursor_install.go` to delegate `--node <alias>` commands to the automated provisioning engine.

2. **Native GitMap Python Runner (`gitmap py` / `gitmap python`) (`REQ-02`)**:
   - First-class CLI subcommand package in `cli/cmdpy/` supporting direct script execution (`gitmap py script.py [args...]`), inline evaluation (`gitmap py -c "<code>"`), and interactive invocation.
   - Root CLI routing in `cli/cmd/rootutility.go` and `cli/cmd/rootcore.go` mapping both `py` and `python` tokens.
   - Cross-platform interpreter resolution traversing virtual environments, `python3`, `python`, and Windows launcher `py`.
   - Transparent stdio streaming (`os.Stdin`, `os.Stdout`, `os.Stderr`) and faithful child process exit code propagation.
   - Dual SQLite telemetry recording: Every execution is logged to Split-DB `commands.db` (`CommandHistory` table) for developer heatmap analytics, and to AI instruction DB (`RecordAiExecution`) for autonomous agent telemetry.

---

## 3. High-Level Architecture & Data Flow

```
                                  +-------------------------------------------------+
                                  |              Developer / Agent Trigger          |
                                  +-------------------------------------------------+
                                          |                                  |
                                          | [gitmap py / gitmap python]      | [gitmap cursor install --node u1]
                                          v                                  v
                    +-----------------------------+          +------------------------------------+
                    |       cli/cmdpy/py_cmd.go   |          |    cli/cmdcursor/cursor_install.go |
                    +-----------------------------+          +------------------------------------+
                                   |                                                 |
                       +-----------+-----------+                     +---------------+---------------+
                       |                       |                     |                               |
                       v                       v                     v                               v
             [-c "<code>"]            [<script.py> [args]]     [Local Linux Node]             [Remote SSH Node: u1]
                       |                       |                     |                               |
                       +-----------+-----------+                     |                               v
                                   |                                 v                     [Run via cmdssh.RunSSHExec]
                                   v                     +-------------------------+                 |
                    +-----------------------------+      | repo-secrets/05-scripts/|                 v
                    |    cli/cmdpy/py_exec.go     |      |  setup-cursor-ubuntu.py |<----------------+
                    |  resolvePythonBinary()      |      +-------------------------+
                    |  exec.Command Streaming     |                  |
                    +-----------------------------+                  | 1. Query Upstream Cursor API:
                                   |                                 |    https://www.cursor.com/api/download
                                   |                                 | 2. Download AppImage to /opt/cursor/Cursor.AppImage
                       +-----------+-----------+                     | 3. Create wrapper: /usr/local/bin/cursor
                       |                       |                     | 4. Inject settings: ~/.config/Cursor/User/settings.json
                       v                       v                     | 5. Verify headless execution (cursor --version)
          +-------------------------+  +--------------------------+  v
          |  cli/store/ai_exec_db   |  | cli/store/command_history|  +------------------------------------+
          |  RecordAiExecution()    |  | InsertCommandRecord()    |  | repo-secrets/04-ubuntu-migration/  |
          |  (.gitmap/data/ai/      |  | (.gitmap/data/history/   |  |   cursor-fleet-status.json         |
          |   instruction/sql.db)   |  |  commands.db)            |  +------------------------------------+
          +-------------------------+  +--------------------------+
                       |                       |
                       +-----------+-----------+
                                   |
                                   v
                      +--------------------------+
                      | Terminal Stdio Stream    |
                      | & Child Exit Code Exit   |
                      +--------------------------+
```

---

## 4. Cursor IDE Deployment Architecture on Ubuntu Fleet

### 4.1 Upstream Release API Contract

Cursor publishes official binary releases through dynamic routing endpoints. To eliminate stale hardcoded URLs and ensure tamper-resistant downloads, the provisioning engine queries the upstream release endpoint:

```http
GET https://www.cursor.com/api/download?platform=linux-x64&releaseTrack=stable
User-Agent: GitMap-FleetProvisioner/2.0
```

#### Protocol Semantics:
1. **HTTP Redirect Follow**: The endpoint issues an HTTP `302 Found` or `307 Temporary Redirect` pointing to the canonical CDN asset (typically hosted on `download.todesktop.com` or AWS CloudFront).
2. **Metadata Headers**: Response headers or location headers contain the explicit version slug (e.g., `Cursor-0.45.11-x86_64.AppImage`).
3. **Checksum Verification**: Where SHA256 checksums are published in accompanying headers or JSON metadata, the provisioning engine performs SHA-256 validation. When downloaded as an AppImage stream, the SHA-256 is computed and persisted in the audit ledger.

### 4.2 Provisioning Engine Architecture (`repo-secrets/05-scripts/`)

Provisioning is implemented in two synergistic scripts:
1. `repo-secrets/05-scripts/setup-cursor-ubuntu.py`: The high-level idempotent orchestrator utilizing Python 3 standard library (`urllib.request`, `hashlib`, `json`, `os`, `pathlib`, `subprocess`).
2. `repo-secrets/05-scripts/setup-cursor-ubuntu.sh`: The POSIX shell bootstrap loader suitable for one-liner curl pipes or systems where Python 3 is bootstrapping.

#### Execution Phases:

| Phase | Description | Actions & Invariants |
|:---|:---|:---|
| **Phase 1: Pre-Flight Audit** | OS & Architecture Check | Verify `uname -m == x86_64`, check `/etc/os-release` for Ubuntu/Debian compatibility, verify root/sudo privileges. |
| **Phase 2: Dependencies** | System Libraries | Verify and install runtime packages: `libfuse2` (or `fuse3` AppImage hook), `libnss3`, `libasound2`, `libgbm1`, `libxss1`, `curl`, `ca-certificates`. |
| **Phase 3: Upstream Query** | API Resolution | Issue request to `https://www.cursor.com/api/download?platform=linux-x64&releaseTrack=stable`, follow redirects, resolve target URL. |
| **Phase 4: Asset Staging** | Download & Permissions | Download to `/opt/cursor/Cursor.AppImage.tmp`, compute SHA-256, rename atomically to `/opt/cursor/Cursor.AppImage`, apply `chmod 0755`. |
| **Phase 5: CLI Wrapper** | Sandbox & Shell Wrapper | Write shell wrapper to `/usr/local/bin/cursor` and `$HOME/.local/bin/cursor` passing `--no-sandbox "$@"`. |
| **Phase 6: Desktop Launcher** | XDG Integration | Write `/usr/share/applications/cursor.desktop` and extract or place icon at `/opt/cursor/cursor.png`. |
| **Phase 7: Settings & Dracula** | Editor Invariants | Ensure `$HOME/.config/Cursor/User/settings.json` has Dracula Dark theme and GitMap editor standards. |
| **Phase 8: Liveness & Status** | Verification & Persistence | Test headless execution (`cursor --version`), update `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`. |

### 4.3 Shell Wrapper & Sandboxing Specification

On modern Ubuntu releases (notably Ubuntu 24.04 LTS and hardened 22.04 LTS kernels), unprivileged user namespace restrictions (`kernel.unprivileged_userns_clone=0` or AppArmor sandboxing profiles) prevent Electron AppImages from spawning sandboxed helper processes unless granted SUID root permissions or executed with `--no-sandbox`.

To ensure seamless headless, agentic, and desktop usage across the fleet without breaking system security, GitMap installs an execution wrapper:

#### Wrapper Target Locations:
- `/usr/local/bin/cursor` (system-wide)
- `$HOME/.local/bin/cursor` (user-local fallback)

#### Wrapper Implementation:
```sh
#!/bin/sh
# GitMap Managed Cursor IDE Launcher
# Auto-generated by repo-secrets/05-scripts/setup-cursor-ubuntu.py
export ELECTRON_ENABLE_LOGGING=0
exec /opt/cursor/Cursor.AppImage --no-sandbox "$@"
```

### 4.4 Desktop Launcher Specification

Location: `/usr/share/applications/cursor.desktop`

```ini
[Desktop Entry]
Name=Cursor
GenericName=AI Code Editor
Comment=Cursor AI-first Code Editor
Exec=/usr/local/bin/cursor %F
Icon=/opt/cursor/cursor.png
Type=Application
StartupNotify=true
StartupWMClass=Cursor
Categories=Development;IDE;TextEditor;
MimeType=text/plain;inode/directory;
```

### 4.5 Editor Configuration & Dracula Theming Invariants

Location: `$HOME/.config/Cursor/User/settings.json`

The provisioning engine merges the following canonical settings without overwriting custom user keybindings:

```json
{
  "workbench.colorTheme": "Dracula Theme",
  "editor.fontFamily": "'JetBrains Mono', 'Fira Code', Consolas, monospace",
  "editor.fontSize": 14,
  "editor.lineHeight": 22,
  "editor.tabSize": 4,
  "editor.insertSpaces": true,
  "files.autoSave": "afterDelay",
  "files.autoSaveDelay": 1000,
  "files.eol": "\n",
  "files.insertFinalNewline": true,
  "files.trimTrailingWhitespace": true,
  "editor.renderWhitespace": "selection",
  "telemetry.telemetryLevel": "off"
}
```

### 4.6 Fleet Status Audit Ledger (`repo-secrets/04-ubuntu-migration/cursor-fleet-status.json`)

To enable fleet-wide visibility and tracking, the provisioning engine reads, updates, and writes a persistent JSON ledger:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "CursorFleetStatus",
  "version": "1.0.0",
  "lastUpdated": "2026-10-05T12:00:00Z",
  "nodes": {
    "u1": {
      "nodeAlias": "u1",
      "hostname": "ubuntu-u1",
      "os": "Ubuntu 24.04 LTS",
      "arch": "x86_64",
      "cursorInstalled": true,
      "cursorVersion": "0.45.11",
      "appImagePath": "/opt/cursor/Cursor.AppImage",
      "wrapperPath": "/usr/local/bin/cursor",
      "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "wrapperVerified": true,
      "settingsApplied": true,
      "theme": "Dracula Theme",
      "status": "healthy",
      "exitCode": 0,
      "lastVerified": "2026-10-05T12:00:00Z"
    }
  }
}
```

### 4.7 Integration in `cli/cmdcursor/cursor_install.go`

`cli/cmdcursor/cursor_install.go` is refactored so that:
1. `installOnLinux`: Invokes `python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py` (or falls back to `bash repo-secrets/05-scripts/setup-cursor-ubuntu.sh`).
2. `delegateRemoteInstall`: SSHs into `--node <alias>` and executes:
   ```sh
   python3 -c "import urllib.request; ..." || bash -c "curl -fsSL ... | bash"
   ```
   Or executes the synced `repo-secrets/05-scripts/setup-cursor-ubuntu.py` on the target remote node, capturing stdout/stderr and verifying JSON exit status.

---

## 5. Native GitMap Python Runner (`gitmap py` / `gitmap python`)

### 5.1 CLI Command Syntax & Semantics

The Python runner is invoked via two primary command tokens: `gitmap py` and `gitmap python`.

```bash
# 1. Inline Code Evaluation
gitmap py -c "import sys; print(sys.version)"
gitmap python -c "print('Hello from GitMap Python runner')"

# 2. Script File Execution with Positional and Flag Arguments
gitmap py 03-ai-scripts/37-bump-version.py --minor
gitmap python repo-secrets/05-scripts/setup-cursor-ubuntu.py --node u1

# 3. Arbitrary Python Subcommands & Flags
gitmap py -m venv .venv
gitmap py --version
```

### 5.2 Package Architecture (`cli/cmdpy/`)

Following GitMap's 100-line modular design standard, `cli/cmdpy` is decomposed into three focused files:

```
cli/cmdpy/
├── py_cmd.go      # CLI entrypoint, argument parsing, routing, help text (<= 95 lines)
├── py_exec.go     # Binary resolution, exec.Command streaming, exit code extraction (<= 90 lines)
└── py_record.go   # Telemetry recording to RecordAiExecution and CommandHistory (<= 85 lines)
```

#### Detailed File Specifications:

#### 1. `cli/cmdpy/py_cmd.go`
- Declares `RunPy(args []string) error`.
- Handles `--help`, `-h`, and `help`.
- Dispatches inline `-c` execution vs script file execution.
- Returns idiomatic errors or child process exit codes.

#### 2. `cli/cmdpy/py_exec.go`
- Implements `resolvePythonBinary() string`:
  1. Checks active virtualenv: `$VIRTUAL_ENV/bin/python` (Linux/macOS) or `$VIRTUAL_ENV/Scripts/python.exe` (Windows).
  2. Traverses `exec.LookPath("python3")`.
  3. Traverses `exec.LookPath("python")`.
  4. On Windows, falls back to `exec.LookPath("py")`.
- Implements `executePythonStreaming(binary string, args []string, dir string) (int, error)`:
  - Hooks `cmd.Stdin = os.Stdin`, `cmd.Stdout = os.Stdout`, `cmd.Stderr = os.Stderr`.
  - Captures start timestamp and duration.
  - Extracts process exit code from `*exec.ExitError`.

#### 3. `cli/cmdpy/py_record.go`
- Dual telemetry recorder:
  - `recordAiExecution(cmdLine string, args []string, dir string, durationMs int, exitCode int, err error)`:
    Calls `store.RecordAiExecution("python_runner", cmdLine, argsJSON, dir, "127.0.0.1", durationMs, exitCode, "", errMsg, isSuccess)`.
  - `recordCommandHistory(cmdLine string, exitCode int, durationMs int64)`:
    Opens `store.OpenCommandHistorySplitDB("")`, calls `histDB.InsertCommandRecord(cmdLine, "py", exitCode, durationMs)`.

### 5.3 CLI Routing Architecture

In `cli/cmd/rootutility.go`:
```go
func utilityToolEntries() []dispatchEntry {
    return []dispatchEntry{
        {[]string{"py", "python"}, func() error { return cmdpy.RunPy(argsTail()) }},
        // ... existing entries
    }
}
```

In `cli/cmd/rootcore.go`:
Ensures `py` and `python` are recognized across all dispatch layers and help dashboard listings.

---

## 6. Dual SQLite Telemetry & Heatmap Integration

### 6.1 Developer Activity Heatmap (`commands.db`)

GitMap maintains a Split-DB SQLite database at `.gitmap/data/history/commands.db` with table:

```sql
CREATE TABLE IF NOT EXISTS CommandHistory (
    CommandId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommandLine TEXT NOT NULL,
    CommandName TEXT NOT NULL,
    ExitCode INTEGER NOT NULL DEFAULT 0,
    ExecutedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    DurationMs INTEGER NOT NULL DEFAULT 0
);
```

When `gitmap py <script>` runs:
- `CommandName` is recorded as `"py"`.
- `CommandLine` records the exact relative command line: e.g. `py 03-ai-scripts/37-bump-version.py`.
- `DurationMs` records elapsed wall-clock milliseconds.
- `ExitCode` records process status.
- This table feeds the GitMap developer productivity heatmap and command frequency analyzer.

### 6.2 Autonomous AI Instruction Telemetry (`sql.db`)

GitMap maintains an AI execution database at `.gitmap/data/ai-instruction/sql.db` via `store.RecordAiExecution`:

- Category: `"python_runner"`.
- Enables agents to trace whether automated fix scripts or test scripts succeeded or failed, with exact duration and error details.

---

## 7. Sandboxing, Security, and Error Handling

| Scenario | Risk / Failure Mode | Remediation / Architectural Invariant |
|:---|:---|:---|
| **Missing Python Binary** | `exec: "python3": executable file not found in $PATH` | `py_exec.go` returns structured `apperror.NewWithDetails("cmd.py.exec", "E1032", "no Python interpreter found on PATH; please install python3", ...)` with actionable guidance. |
| **Ubuntu AppImage FUSE Error** | AppImage fails with `dlopen(): error loading libfuse.so.2` | Provisioning script checks `dpkg -s libfuse2` and installs `libfuse2` via `apt-get` automatically. |
| **Ubuntu Electron UserNS Error** | AppImage fails with `The SUID sandbox helper binary was not found or is not root SUID` | Launch wrapper explicitly passes `--no-sandbox` as the default execution argument for all invocations. |
| **Corrupted Download Stream** | Network interruption leaves half-written AppImage | Script downloads to temporary `.tmp` file, validates content length and executable headers, then renames atomically. |
| **Non-Zero Python Child Exit** | Script fails with exit code 1 or 2 | `gitmap py` flushes stderr, records exit code to SQLite, and returns exit code faithfully to shell so CI/CD and scripts detect failures. |

---

## 8. Verification & Acceptance Criteria

1. **Cursor Ubuntu Provisioning**:
   - `python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py` completes without errors on Ubuntu nodes.
   - `/usr/local/bin/cursor` and `$HOME/.local/bin/cursor` exist and point to the `--no-sandbox` wrapper.
   - `cursor --version` runs headlessly without throwing sandbox or X11 errors.
   - `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json` reflects `"status": "healthy"` and `"cursorInstalled": true`.
2. **GitMap Python Runner**:
   - `gitmap py -c "print(1+1)"` outputs `2` with exit code 0.
   - `gitmap py -c "import sys; sys.exit(42)"` exits with status code 42.
   - `gitmap python` aliases `gitmap py` identically.
   - Each execution inserts a corresponding row in `.gitmap/data/history/commands.db` and `.gitmap/data/ai-instruction/sql.db`.
3. **Repository Cleanliness & Code Standards**:
   - All newly created Go files in `cli/cmdpy/` strictly obey the `<= 100` lines file length cap and `<= 15` lines function length rule.
   - Zero hardcoded absolute paths; strictly relative paths (`repo-secrets/...`, `cli/...`).
