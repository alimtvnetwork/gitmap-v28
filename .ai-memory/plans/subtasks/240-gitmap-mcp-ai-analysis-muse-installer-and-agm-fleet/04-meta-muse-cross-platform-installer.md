# Subtask 04: Meta Muse Multi-Platform Native Installer

> **Subtask ID:** Subtask-04  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmdinstall/`, `cli/cmd/muse_cmd.go`, `cli/constants/`  
> **Owned Files:**  
> - `cli/cmdinstall/install_muse.go`  
> - `cli/cmdinstall/install_muse_test.go`  
> - `cli/cmd/muse_cmd.go`  
> - `cli/constants/constants_install.go`  

---

## 1. Concrete Objectives

1. **Native Unified Installer Command:**
   - Address the installation flow illustrated in `assets/screenshots/240-muse-meta-installer.png`.
   - Provide a native, single-command installation experience: `gitmap muse install` (with alias `gitmap install muse`).
   - Automatically detect the host operating system and invoke the official Meta Muse installation pipeline.

2. **Cross-Platform Strategy Matrix:**
   - **Windows:**
     - Execute via PowerShell: `powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://dev.meta.ai/install.ps1 | iex"`.
     - Stream download progress in real-time, matching telemetry seen in `240-muse-meta-installer.png` (`Downloading muse 1.4.3-R5018.1 (429 MB)`).
     - Verify executable binary `muse.exe` is registered in system or user PATH.
   - **Linux / Ubuntu:**
     - Execute via bash: `bash -c "curl -fsSL https://dev.meta.ai/install.sh | bash"`.
     - Handle non-root installation paths (`~/.local/bin/muse`) and verify execution permissions.
   - **macOS (Darwin):**
     - Execute via bash: `bash -c "curl -fsSL https://dev.meta.ai/install.sh | bash"`.
     - Support platform architecture detection (Apple Silicon `arm64` vs Intel `x86_64`).

3. **Pre-Flight Guards & Flag Controls:**
   - **Connectivity Probe:** Verify network connectivity to `https://dev.meta.ai` prior to launching child processes.
   - **Existing Binary Guard:** Probe PATH for existing `muse` binary. If present, report detected version and abort unless `--force` is supplied.
   - **Dry-Run Inspection:** `--dry-run` displays target platform, script URL, and command string without altering system state.
   - **Platform Override:** `--platform <windows|linux|darwin>` allows testing or cross-environment script generation.
   - **Post-Install Verification:** `--verify` runs `muse --version` after installation and outputs structured JSON status.

4. **Centralized Constants & Audit Logging:**
   - In `cli/constants/constants_install.go`, define official distribution endpoints:
     - `MetaMuseInstallUrlWindows = "https://dev.meta.ai/install.ps1"`
     - `MetaMuseInstallUrlUnix = "https://dev.meta.ai/install.sh"`
   - Record installation status, timestamp, duration, and exit code in `installation.db` Split-DB.

---

## 2. Core Domain Types & Structs (`cli/cmdinstall/install_muse.go`)

```go
package cmdinstall

import "time"

// MusePlatformType represents the target operating system platform.
type MusePlatformType string

const (
	MusePlatformWindows MusePlatformType = "windows"
	MusePlatformLinux   MusePlatformType = "linux"
	MusePlatformDarwin  MusePlatformType = "darwin"
)

// MuseInstallOptions configures the installation process.
type MuseInstallOptions struct {
	Platform MusePlatformType
	Force    bool
	DryRun   bool
	Verify   bool
	Timeout  time.Duration
}

// MuseInstallResult captures execution telemetry and post-install status.
type MuseInstallResult struct {
	Platform        MusePlatformType `json:"platform"`
	CommandExecuted string           `json:"commandExecuted"`
	Success         bool             `json:"success"`
	ExitCode        int              `json:"exitCode"`
	DurationMs      int64            `json:"durationMs"`
	InstalledPath   string           `json:"installedPath,omitempty"`
	VersionDetected string           `json:"versionDetected,omitempty"`
	ErrorMessage    string           `json:"errorMessage,omitempty"`
}
```

---

## 3. Step-by-Step Implementation Plan

### Step 1: Distribution Constants Declaration

- In `cli/constants/constants_install.go`, declare:
  - `MetaMuseInstallUrlWindows = "https://dev.meta.ai/install.ps1"`
  - `MetaMuseInstallUrlUnix = "https://dev.meta.ai/install.sh"`
  - `MetaMuseBinaryNameWindows = "muse.exe"`
  - `MetaMuseBinaryNameUnix = "muse"`

### Step 2: Platform Strategy Resolver

- In `cli/cmdinstall/install_muse.go`, implement `ResolveMuseInstallCommand(platform MusePlatformType) (string, []string, error)`:
  - For Windows: `powershell`, `[]string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "irm https://dev.meta.ai/install.ps1 | iex"}`.
  - For Linux: `bash`, `[]string{"-c", "curl -fsSL https://dev.meta.ai/install.sh | bash"}`.
  - For Darwin: `bash`, `[]string{"-c", "curl -fsSL https://dev.meta.ai/install.sh | bash"}`.

### Step 3: Installation Runner with Progress Streaming

- In `cli/cmdinstall/install_muse.go`, implement `RunMuseInstaller(opts MuseInstallOptions) (*MuseInstallResult, error)`:
  - Probe connectivity using `http.Head("https://dev.meta.ai")`.
  - Check existing binary via `exec.LookPath`.
  - If dry-run: return planned commands immediately.
  - Construct child process with streaming stdout and stderr.
  - Capture process duration and exit code.
  - If verify: execute `muse --version` and parse output.

### Step 4: Installation DB Audit Logging

- In `cli/cmdinstall/install_muse.go`, implement `logMuseInstallAudit(result *MuseInstallResult)`:
  - Connect to `installation.db` via `store.ResolveSplitDbPath(store.SectionInstallation, "global", "")`.
  - Insert record into `PackageInstallHistory`.

### Step 5: Cobra CLI Command Integration

- Create `cli/cmd/muse_cmd.go`:
  - Define `MuseCmd` (`gitmap muse`) and subcommand `install`.
  - Add alias `gitmap install muse` in `cli/cmdinstall/install_handlers.go`.
  - Bind flags `--platform`, `--force`, `--dry-run`, `--verify`, `--timeout`.

### Step 6: Automated Testing Suite

- Create `cli/cmdinstall/install_muse_test.go`:
  - Test command string resolution across Windows, Linux, and macOS.
  - Test pre-flight existing binary check logic.
  - Test `--dry-run` flag behavior asserting zero child process execution.
  - Test error handling when network endpoint is unreachable.

---

## 4. Acceptance Criteria

- [x] `gitmap muse install` resolves the appropriate platform command for the host OS.
- [x] Windows executes PowerShell one-liner `irm https://dev.meta.ai/install.ps1 | iex`.
- [x] Linux / Ubuntu and macOS execute `curl -fsSL https://dev.meta.ai/install.sh | bash`.
- [x] `--dry-run` outputs the planned execution command without invoking the network or shell.
- [x] `--force` allows overriding an existing `muse` installation.
- [x] `--platform <name>` allows explicit platform target selection.
- [x] Successful installation records an audit entry in `installation.db`.
- [x] Zero absolute paths or `file:///` URIs exist in any owned file.
- [x] All unit tests in `cli/cmdinstall/install_muse_test.go` pass with 100% success rate.

---

## 5. Verification Commands

```bash
  # 1. Run unit tests for Meta Muse installer
  go test -v ./cli/cmdinstall -run "TestMuse.*"

  # 2. Test dry-run execution on Windows
  gitmap muse install --dry-run --platform windows

  # 3. Test dry-run execution on Linux
  gitmap muse install --dry-run --platform linux

  # 4. Test dry-run execution on macOS
  gitmap muse install --dry-run --platform darwin

  # 5. Test alias command dispatch
  gitmap install muse --dry-run

  # 6. Run repository relative paths linter
  python linter-scripts/check-relative-paths.py
```
