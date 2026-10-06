# App Issue 232: Ubuntu Multi-IDE & GitHub Desktop Scan Omission RCA

**Issue ID:** 232  
**Date:** 2026-10-06  
**Status:** Resolved  
**Affected Subsystems:** `cli/vscodepm/` (`path.go`, `sync.go`), `cli/desktop/` (`resolve.go`, `add.go`), `cli/cmdscan/` (`flags.go`, `scan.go`), `cli/cmdide/`, `03-ai-scripts/40-ubuntu-ide-desktop-sync.py`  
**Reference Commits:** `ide-scan-sync - add ubuntu github desktop and ide registration suite`  

---

## 1. Symptom & Incident Description

During routine repository discovery and workspace indexing on Ubuntu Linux workstations (Ubuntu 22.04 LTS and 24.04 LTS), running `gitmap scan` (or `gitmap rescan`) successfully discovers Git repositories on disk and persists them into the local SQLite tracking database (`gitmap.db`), but completely fails to register the discovered repositories into local IDEs and GUI clients:

1. **Visual Studio Code Project Manager Omission:**
   - Discovered repositories do not appear in the VS Code Project Manager extension sidebar, status bar palette, or quick pick list.
   - The extension's project storage file `projects.json` is neither created nor updated under `~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json`.
   - The CLI logs a soft skip or silent error without creating the missing directories or alerting the user.

2. **GitHub Desktop Linux Omission:**
   - Discovered repositories are not imported into GitHub Desktop on Linux.
   - The terminal summary shows zero repositories added to GitHub Desktop (`Added: 0, Failed: 0`).
   - GitHub Desktop installed via Debian packages (`/usr/bin/github`, `/usr/bin/github-desktop`, `/usr/lib/github-desktop/resources/app/static/github`) or Flatpak packages is never invoked.

3. **Cursor IDE Synchronization Blindspot:**
   - Developers using Cursor on Ubuntu (`~/.config/Cursor/User/`) observe zero synchronization.
   - `gitmap scan` possesses no registration logic for Cursor's Project Manager storage.

4. **Google Antigravity Workspace Registration Omission:**
   - Repositories are omitted from Google Antigravity project descriptors (`~/.gemini/config/projects/*.json`).
   - Developers must manually add each repository via the UI, resulting in multi-tool fragmentation.

---

## 2. Grounded Root Cause Analysis (RCA)

A comprehensive forensic inspection of the GitMap scanning pipeline, path resolvers, and flag registries revealed four discrete architectural root causes:

### 2.1 Vector 1: Omission of Auto-Mkdir in `cli/vscodepm/path.go`

In `cli/vscodepm/path.go` (lines 49–65), `ProjectsJSONPath()` resolves the storage path for the Project Manager extension:

```go
func ProjectsJSONPath() (string, error) {
	root, err := UserDataRoot()
	if err != nil {
		return root, err
	}

	extDir := filepath.Join(root,
		constants.VSCodePMUserDir,
		constants.VSCodePMGlobalStorageDir,
		constants.VSCodePMExtensionDir)

	if !dirExists(extDir) {
		return filepath.Join(extDir, constants.VSCodePMProjectsFile), ErrExtensionMissing
	}

	return filepath.Join(extDir, constants.VSCodePMProjectsFile), nil
}
```

- **Forensic Failure:**
  On a fresh Ubuntu installation, or when the Project Manager extension has been installed but has not yet saved any projects, `~/.config/Code/User` exists but `extDir` (`~/.config/Code/User/globalStorage/alefragnani.project-manager`) does not yet exist.
- Instead of proactively creating the directory hierarchy with `os.MkdirAll(extDir, 0o755)`, `ProjectsJSONPath()` strictly checks `if !dirExists(extDir)` and immediately returns `ErrExtensionMissing`.
- In `cli/cmdvscode/clonepmsync.go` (lines 143–148), `vscodepm.Sync(pairs)` intercepts `ErrExtensionMissing` via `reportVSCodePMSoftError()`. The error is treated as a benign soft skip, which suppresses fatal exits but permanently drops all records without writing `projects.json`.

### 2.2 Vector 2: GitHub Desktop Linux CLI Discovery Blindspot in `cli/desktop/resolve.go`

In `cli/desktop/resolve.go` (lines 43–56), `knownInstallCandidates()` defines platform-specific fallback search paths for the GitHub Desktop CLI shim:

```go
func knownInstallCandidates() []string {
	if runtime.GOOS == constants.OSWindows {
		return windowsCandidates()
	}

	if runtime.GOOS == "darwin" {
		return darwinCandidates()
	}

	return nil
}
```

- **Forensic Failure:**
  `knownInstallCandidates()` returns `nil` for Linux (`runtime.GOOS == "linux"`).
- `ResolveCLI()` in `cli/desktop/resolve.go` (lines 21–28) first attempts `exec.LookPath(constants.GitHubDesktopBin)`, where `constants.GitHubDesktopBin` is hardcoded to `"github"`.
- On Linux distributions:
  1. Many package distributions install the binary under `/usr/bin/github-desktop` or `/usr/local/bin/github-desktop`, rather than `"github"`.
  2. The canonical CLI shim resides at `/usr/lib/github-desktop/resources/app/static/github`.
  3. If `/usr/bin/github` is missing from `PATH` or only `github-desktop` is installed, `exec.LookPath("github")` fails with `exec.ErrNotFound`.
  4. The resolver falls back to `resolveFromKnownInstalls()`, which iterates over `knownInstallCandidates()`. Because Linux returns `nil`, the function returns an empty string `""`, completely aborting Desktop registration.
- In addition, GitHub Desktop on Linux does not implement a dedicated `add` CLI verb; its command interface uses `open <path>`.

### 2.3 Vector 3: Default False Scan Flag in `cli/cmdscan/flags.go`

In `cli/cmdscan/flags.go` (lines 62–71):

```go
func registerScanToggles(fs *flag.FlagSet, flagPtrs *scanFlagPointers) {
	flagPtrs.ghDesktopFlag = fs.Bool("github-desktop", false, constants.FlagDescGHDesktop)
	flagPtrs.openFlag = fs.Bool("open", false, constants.FlagDescOpen)
	flagPtrs.quietFlag = fs.Bool("quiet", false, constants.FlagDescQuiet)
	...
}
```

- **Forensic Failure:**
  The `ghDesktopFlag` flag defaults to `false`.
- In `cli/cmdscan/scan.go` (lines 183–185 and 342–348):
  ```go
  bench.Phase("scan.addToDesktop", func() {
      addToDesktop(records, ghDesktop)
  })
  ```
  Because `ghDesktop` defaults to `false`, routine runs of `gitmap scan` completely bypass GitHub Desktop registration unless the user explicitly supplies `--github-desktop`. Users expect `gitmap scan` to automatically register detected repositories across local developer tooling.

### 2.4 Vector 4: Absence of Multi-IDE Registration Pipeline

In `cli/cmdscan/scan.go` (lines 183–188):

```go
bench.Phase("scan.addToDesktop", func() {
    addToDesktop(records, ghDesktop)
})
bench.Phase("scan.vscodePMSync", func() {
    syncRecordsToVSCodePM(records, noVSCodeSync, noAutoTags)
})
```

- **Forensic Failure:**
  The post-scan hook only delegates to `addToDesktop` and `syncRecordsToVSCodePM`.
- No registration hooks exist for:
  1. **Cursor IDE:** `~/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json`.
  2. **Google Antigravity:** `~/.gemini/config/projects/<uuid>.json`.
- Developers working across modern polyglot editors experience fractured workspace states where repositories are indexed in GitMap but completely missing from their daily editors.

---

## 3. Blast Radius & Impact Assessment

| Subsystem / Tool | Impact Level | Observed Manifestation |
| :--- | :--- | :--- |
| **VS Code Project Manager** | High | Zero scanned repositories appear in the PM sidebar on Ubuntu machines where the extension directory was not pre-created. |
| **GitHub Desktop (Linux)** | High | Discovered repositories are never registered; CLI discovery fails if `github-desktop` is the primary binary name. |
| **Cursor IDE** | Medium | No synchronization hook; manual project importing is required. |
| **Google Antigravity** | Medium | No workspace descriptor generation during scan; agent sessions cannot discover repositories without manual onboarding. |
| **Developer Experience** | High | Users observe broken expectations: `gitmap scan` reports successful scans, yet no editors show the newly discovered repositories. |

---

## 4. Comprehensive Remediation Matrix

To permanently resolve all four failure vectors, a structured multi-phase remediation is implemented across standalone tooling, core CLI libraries, and scanning hooks:

| Item | Failure Vector | Target File(s) | Remediation Architecture | Verification Gate |
| :---: | :--- | :--- | :--- | :--- |
| **R-01** | Vector 1: Omission of Auto-Mkdir | `cli/vscodepm/path.go` | Implement `EnsureProjectsJSONPath()` that executes `os.MkdirAll(extDir, 0o755)` prior to returning path for write operations. | Invoking write-mode resolution when `extDir` is missing creates directory and returns valid path. |
| **R-02** | Vector 2: Linux CLI Resolution | `cli/desktop/resolve.go` | Add `linuxCandidates()` probing `/usr/bin/github`, `/usr/bin/github-desktop`, `/usr/local/bin/github`, and `/usr/lib/github-desktop/resources/app/static/github`. Probe `github-desktop` in `exec.LookPath`. | `ResolveCLI()` discovers GitHub Desktop binary on Ubuntu Linux with 100% reliability. |
| **R-03** | Vector 2: CLI Invocation Grammar | `cli/desktop/add.go` | Ensure registration dispatches `open <path>` when executing GitHub Desktop CLI on Linux. | Verified execution of `github open <path>` registers repository without errors. |
| **R-04** | Vector 3 & 4: Multi-IDE Scan Hooks | `cli/cmdscan/flags.go`, `cli/cmdscan/scan.go` | Introduce affirmative `--sync-ide`, `--skip-sync`, and `--exclude-sync` flags. Enable IDE synchronization by default with granular exclusion filters. | `gitmap scan` registers repositories across all active IDEs; `--skip-sync desktop` cleanly skips GitHub Desktop. |
| **R-05** | Vector 4: Cursor & Antigravity Support | `cli/cmdide/` | Architect `gitmap ide` command group (`add`, `sync`, `remove`, `ls`, `status`) supporting VS Code, Cursor, Antigravity, and GitHub Desktop. | `gitmap ide status` reports health for all 4 targets; `gitmap ide sync` populates all 4 targets. |
| **R-06** | All Vectors: Standalone Script | `03-ai-scripts/40-ubuntu-ide-desktop-sync.py` | Create standalone Python 3 utility discovering repos from disk or `gitmap.db`, auto-creating directories, and syncing all 4 targets with `--dry-run`, `--json`, and filters. | Standalone script executes with `--dry-run` and `--json` emitting complete sync telemetry without errors. |

---

## 5. Prevention & Learnings

1. **Write-Intent Directory Preconditions:**
   Path resolution functions intended for mutating operations must proactively ensure parent directory existence (`os.MkdirAll`) rather than failing when an extension's storage folder has not yet been initialized.
2. **Cross-Platform Binary Discovery Parity:**
   When defining platform-specific candidate paths, Linux must receive equal priority with Windows and macOS. Shims and package-manager aliases (`github`, `github-desktop`) must both be probed.
3. **Explicit Affirmative Configuration:**
   Scan synchronization features must be governed by affirmative boolean flags (`isVSCodeTargeted`, `isGitHubDesktopTargeted`) with sensible default enablement and explicit skip controls (`--skip-sync`).
4. **Standalone Automation First:**
   Providing a standalone script (`40-ubuntu-ide-desktop-sync.py`) alongside compiled Go CLI integration ensures rapid recovery, debugging transparency, and headless execution in automation pipelines.
