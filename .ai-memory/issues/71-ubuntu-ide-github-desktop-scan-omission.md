# App Issue 232: Ubuntu Multi-IDE & GitHub Desktop Scan Omission RCA

**Issue ID:** 232  
**Date:** 2026-10-06  
**Status:** Resolved  
**Parent Architecture Spec:** `02-spec/21-app/232-ubuntu-ide-and-github-desktop-scan-sync/01-architecture-spec.md`  
**Detailed Specification:** `02-spec/22-app-issues/232-ubuntu-ide-github-desktop-scan-omission.md`  
**Affected Subsystems:** `cli/vscodepm/`, `cli/desktop/`, `cli/cmdscan/`, `cli/cmdide/`, `03-ai-scripts/40-ubuntu-ide-desktop-sync.py`  
**Reference Commits:** `ide-scan-sync - add ubuntu github desktop and ide registration suite`  

---

## 1. Symptom & Incident Description

When executing `gitmap scan` or `gitmap rescan` on Ubuntu Linux workstations:
- Repositories are properly detected on disk and registered into `gitmap.db`.
- However, zero repositories are registered into Visual Studio Code's Project Manager extension (`~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json`).
- Repositories are not registered into GitHub Desktop on Linux (`Added: 0, Failed: 0`).
- Cursor IDE Project Manager storage (`~/.config/Cursor/User/...`) receives zero updates.
- Google Antigravity project descriptors (`~/.gemini/config/projects/`) receive zero updates.

---

## 2. Root Cause Analysis (4 Vectors)

1. **Vector 1: Omission of Auto-Mkdir in `cli/vscodepm/path.go` (lines 49–65):**
   - `ProjectsJSONPath()` checks `if !dirExists(extDir)` and immediately returns `ErrExtensionMissing`.
   - On freshly installed systems or before the Project Manager extension writes data, `extDir` is missing.
   - `syncClonedReposToVSCodePM()` swallows this error as a soft skip, dropping all records without creating the folder hierarchy or writing `projects.json`.
2. **Vector 2: Linux CLI Discovery Blindspot in `cli/desktop/resolve.go` (lines 43–56):**
   - `knownInstallCandidates()` returns `nil` when `runtime.GOOS` is Linux.
   - `ResolveCLI()` probes `exec.LookPath("github")`. If `/usr/bin/github` is missing from `PATH` or only `/usr/bin/github-desktop` is installed, candidate discovery returns empty string.
   - GitHub Desktop on Linux requires `github open <path>` invocation rather than non-existent `add` subcommands.
3. **Vector 3: Default False Scan Flag in `cli/cmdscan/flags.go` (line 63):**
   - `flagPtrs.ghDesktopFlag = fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)`.
   - GitHub Desktop sync defaults to disabled during routine scans unless `--github-desktop` is explicitly provided.
4. **Vector 4: Absence of Multi-IDE Registration Pipeline in `cli/cmdscan/scan.go` (lines 183–188):**
   - Scan loop only dispatches to `addToDesktop` and `syncRecordsToVSCodePM`, completely omitting Cursor and Google Antigravity.

---

## 3. Blast Radius & Affected Files

- `cli/vscodepm/path.go`: Path resolution failing with `ErrExtensionMissing` instead of auto-creating parent directories.
- `cli/desktop/resolve.go`: Missing Linux candidate paths (`/usr/bin/github`, `/usr/bin/github-desktop`, `/usr/lib/github-desktop/resources/app/static/github`).
- `cli/desktop/add.go`: Linux CLI dispatch mechanism requiring `open <path>`.
- `cli/cmdscan/flags.go`: Scan flags lacking affirmative multi-IDE synchronization and skip toggles.
- `cli/cmdscan/scan.go`: Scan execution lacking unified multi-IDE hooks.
- `03-ai-scripts/40-ubuntu-ide-desktop-sync.py`: Required standalone automation script.

---

## 4. Comprehensive Remediation Matrix

| Item | Problem Vector | Target File | Action Taken | Verification |
| :---: | :--- | :--- | :--- | :--- |
| **R-01** | Vector 1: Auto-Mkdir | `cli/vscodepm/path.go` | Implement `EnsureProjectsJSONPath()` with `os.MkdirAll` (0o755). | Write path succeeds even when parent directory is absent. |
| **R-02** | Vector 2: Linux CLI | `cli/desktop/resolve.go` | Add `linuxCandidates()` and probe both `github` and `github-desktop`. | `ResolveCLI()` finds `/usr/bin/github` or `/usr/bin/github-desktop` on Ubuntu. |
| **R-03** | Vector 3: Flag Defaults | `cli/cmdscan/flags.go` | Add affirmative `--sync-ide`, `--skip-sync`, `--exclude-sync`. | Defaults enable IDE sync with flexible skip support. |
| **R-04** | Vector 4: Multi-IDE | `cli/cmdide/` | Implement dedicated `gitmap ide` command group (`add`, `sync`, `remove`, `ls`, `status`). | All 4 targets (VS Code, Cursor, Antigravity, Desktop) manageable via CLI. |
| **R-05** | Standalone Script | `03-ai-scripts/40-ubuntu-ide-desktop-sync.py` | Create standalone Python tool discovering repos, auto-creating directories, and syncing all 4 IDE targets. | Script runs with `--dry-run` and `--json` emitting complete telemetry. |

---

## 5. Prevention & Verification Guidelines

- Strictly evaluate positive booleans (`isVSCodeEnabled`, `isGitHubDesktopEnabled`, `isAutoMkdirEnabled`).
- Ensure all file references use relative repository paths.
- Standalone synchronization script in `03-ai-scripts/40-ubuntu-ide-desktop-sync.py` provides immediate self-contained remediation.
