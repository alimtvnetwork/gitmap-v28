# Spec: installed-dir Command & Linux Update Path Rebuild

## Overview

Adds the `installed-dir` (alias `id`) command to display the active
gitmap-v28 binary location, and integrates Linux/macOS shell-based update
flow that rebuilds the binary into the resolved install path.

---

## 1. installed-dir Command

### Behavior

- Resolves the currently running executable path via `os.Executable()`.
- Follows symlinks via `filepath.EvalSymlinks()` to find the real path.
- Prints both the full binary path and its parent directory.

### Command Registration

| Name | Alias | Dispatch |
|------|-------|----------|
| `installed-dir` | `id` | `dispatchUtility` |

### Output Format

```
  📂 Installed directory

  Binary:    /home/alim/.local/bin/gitmap-v28
  Directory: /home/alim/.local/bin
```

---

## 2. Linux/macOS Update Flow

### Problem

The `gitmap-v28 update` command's `executeUpdate` function only supported
PowerShell (`run.ps1`), making it non-functional on Linux/macOS.

### Solution

Added `executeUpdateUnix()` which:

1. Resolves the active binary's install directory using `resolveInstalledDir()`:
   - First tries `exec.LookPath("gitmap-v28")` to find the PATH binary.
   - Falls back to `os.Executable()` + `filepath.EvalSymlinks()`.
2. Runs `bash run.sh --update` from the source repository root.
3. `run.sh` already handles:
   - Git pull with conflict resolution (stash/discard/clean).
   - Go dependency resolution (`go mod tidy`).
   - Binary compilation with embedded `RepoPath` ldflags.
   - Deployment to the configured path.
   - **PATH binary sync** (lines 601-618): if the active `which gitmap-v28`
     differs from the deployed binary, it copies the new binary to the
     active PATH location.
4. Version verification and cleanup after update.

### Platform Dispatch

```
executeUpdate(repoPath)
  ├── Windows → executeUpdateWindows() → PowerShell run.ps1
  └── Linux/macOS → executeUpdateUnix() → bash run.sh --update
```

---

## 3. Install Path Resolution Strategy

| Priority | Method | Description |
|----------|--------|-------------|
| 1 | `exec.LookPath("gitmap-v28")` | Finds the binary on PATH |
| 2 | `os.Executable()` | Current process executable |
| 3 | `filepath.EvalSymlinks()` | Resolves symlinks to real path |

---

## 4. Files Changed

| File | Change |
|------|--------|
| `gitmap-v28/cmd/installeddir.go` | New: `runInstalledDir()` command |
| `gitmap-v28/cmd/updatescript.go` | Added `executeUpdateUnix()`, `resolveInstalledDir()` |
| `gitmap-v28/cmd/rootutility.go` | Added `installed-dir` / `id` dispatch |
| `gitmap-v28/constants/constants_cli.go` | Added `CmdInstalledDir`, `CmdInstalledDirAlias`, `HelpInstalledDir` |
| `gitmap-v28/constants/constants_update.go` | Added `ErrUpdateNoRunSH`, `MsgUpdateInstallDir` |
| `gitmap-v28/helptext/installed-dir.md` | New: command documentation |

---

## 5. Acceptance Criteria

1. `gitmap-v28 installed-dir` prints the binary path and directory.
2. `gitmap-v28 id` produces the same output.
3. `gitmap-v28 update` on Linux/macOS runs `run.sh --update` from the source repo.
4. After update, `gitmap-v28 version` reflects the latest version.
5. The active PATH binary is synced to the newly built version.
6. If `run.sh` is missing, a clear error is shown.
7. Symlinks are resolved to the real binary location.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
