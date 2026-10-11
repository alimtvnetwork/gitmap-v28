# GitMap v6.537.0 Release Notes

## Summary of Changes
GitMap v6.537.0 delivers architectural modularization to satisfy the 300-line Go file size gate (spec 243.4), alongside completed enhancements for companion repo-cache discovery, two-tier subtree remediation rendering, antivirus false-positive mitigations, and native Golang path query acceleration.

---

### 1. Codebase Architecture & File Size Compliance
- **Decomposed `cli/cmd/rootutility.go`:**
  - `cli/cmd/rootutility.go` (269 lines): Retains utility command dispatch tables, categorization metadata, and CLI help routing.
  - `cli/cmd/rootutility_update.go` (138 lines): Core update orchestration, fleet update dispatch, remote SSH update routing, and AGM updates.
  - `cli/cmd/rootutility_update_args.go` (189 lines): Argument parsing, token validation, command suggestions, and update help UI.
- All modified `.go` files strictly adhere to the repository $\le 300$ line cap.
- Verified with `03-ai-scripts/52-file-size-check.py` with zero violations.

---

### 2. Feature Highlights (Cumulative in v6.537.0)
- **Repo-Cache JSON Injection:**
  - `cli/cmdscan/scanoutput.go` and `cli/cmdscan/scan.go`: `ProbeWorkspaceRepoCacheRecord()` checks the default work directory (`store.OpenDefault().GetDefaultWorkDir()`) and `SpecialReposSplitDB`. Synthesizes a companion `model.ScanRecord` with notes `companion:repo-cache` for automatic inclusion in `gitmap.json` and `.gitmap/output/gitmap.json`.
- **Two-Tier Subtree Solutions & Native Remediation:**
  - `cli/cmdpull/pull_efficient_render.go`: Refactored remediation solutions from flat inline text into structured two-tier subtrees (`├── Solutions:`, `│   ├── Option 1:`, `│   │   └── command`).
  - `cli/cmdpull/pull_remediation_fallback.go` & `cli/cmdpull/pull_remediation_hint.go`: Replaced raw `git -C ...` strings with native GitMap commands (`gitmap sends cp`, `gitmap stash`, `gitmap fix`, `gitmap pull`, and fleet option `gitmap sends cp all "wip: save changes"`).
- **Antivirus & EDR False-Positive Mitigations:**
  - `cli/cmdos/os_autologin_windows.go`: Removed writing plaintext passwords to `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon` (`Trojan:Win32/AutoLogon` mitigation).
  - `cli/cmdos/os_tweak_windows_privacy.go`: Removed `sc.exe DiagTrack` tampering commands (`Behavior:Win32/TamperWinDefService` mitigation).
  - `cli/cmdagy/agy_window_focus_windows.go`: Hardened user32 window focus callbacks with affirmative booleans and early exits.
- **Golang Native Path Acceleration (<5ms):**
  - `cli/cmd/special_repos_cli.go`: Added `--path` and `path` handling to `gitmap rc` and `gitmap rs`.
  - `cli/cmd/path_cli.go`: Added `gitmap path <repo|rc|rs>` and `--json` flag.
  - `cli/cmdcd/cdops.go`: Added `--path` flag to `gitmap cd <repo> --path`.
  - Replaces slow AI Python subprocess calls (`gitmap py -c "import subprocess; print(subprocess.check_output(['gitmap', 'cd', 'rc'], text=True))"`).
- **Skills & Coding Guidelines Synchronized:**
  - `.agents/skills/gitmap/SKILL.md`: Added Command Replacement Matrix and Section 5 Cheat Sheet.
  - `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md`: Added Mistake #16.
  - `02-spec/02-coding-guidelines/12-python/02-standards.md`: Added Section 5 rule banning Python subprocesses for repository path discovery.

---

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.537.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.537.0/install.sh | sh
```
