# Completed Task: Pure Binary Updater & Modern Terminal UI Installer

**Plan Slug:** `bug-fix-and-terminal-ui-improvement`  
**Status:** Completed  
**Spec Link:** `02-spec/21-app/bug-fix-and-terminal-ui-improvement/readme.md`  
**Version Released:** `v6.557.0`  
**Completion Date:** 2026-10-16

---

## 1. Executive Summary

Resolved a critical bug in `gitmap update` where the updater previously initiated a source-discovery cascade and triggered `run.ps1 -Update` (`git pull` + `go build`) whenever a repository checkout was detected. The engine was completely redesigned to prioritize pre-compiled release binary downloads and console-safe in-place executable handoffs without requiring Go or Git. Additionally, modernized terminal UI across `install.ps1`, `install.sh`, and the in-binary `gitmap update` command using 24-bit Truecolor Dracula ANSI palettes, rounded box-drawing cards, staged progress indicators, and interactive Braille spinners. Added standalone offline local installers `local-install.ps1` and `local-install.sh`.

---

## 2. Deliverables & Subtask Audit

### Task-01: Pure Binary Download & In-Place Handoff Engine
- **Target Files:** `cli/cmdupdate/update.go`, `cli/cmdupdate/update_binary_swap.go`, `cli/cmdupdate/updateremoteinstall.go`, `cli/cmdupdate/updaterepo.go`
- **Accomplishments:**
  - Inverted update strategy in `cli/cmdupdate/update.go`: binary release update (`runBinaryReleaseUpdate`) is now the primary, default path.
  - Eliminated automatic `resolveRepoPath()` cascade and silent fallback to `run.ps1 -Update` (`git pull`, `go build`). Source rebuild strictly requires `--source-rebuild` or `--from-source`.
  - Implemented direct GitHub Releases download for platform archives (`gitmap-<tag>-windows-amd64.zip` / `gitmap-<tag>-linux-amd64.tar.gz`) into an isolated temporary staging directory in `os.TempDir()`.
  - Implemented cryptographic SHA256 integrity verification against `checksums.txt` using `subtle.ConstantTimeCompare`.
  - Unpacked binary directly using Go stdlib `archive/zip` and `archive/tar` + `compress/gzip` with directory traversal protection.
  - Implemented console-safe Windows atomic swap (`gitmap.exe` -> `.old`, new binary moved in) and spawned detached `gitmap.exe update-cleanup` with 1500ms delay (`GITMAP_UPDATE_CLEANUP_DELAY_MS=1500`).

### Task-02: Colorful Lipgloss ANSI Terminal UI in Installers
- **Target Files:** `install.ps1`, `install.sh`, `local-install.ps1`, `local-install.sh`
- **Accomplishments:**
  - Modernized `install.ps1`: enabled Win32 Virtual Terminal processing via `kernel32.dll`, removed `$PSStyle.OutputRendering = 'PlainText'` suppression, and implemented Dracula 24-bit Truecolor palette (`#8be9fd`, `#bd93f9`, `#50fa7b`, `#ffb86c`, `#ff5555`, `#6272a4`) with `NO_COLOR` / `TERM=dumb` fallback.
  - Added rounded box-drawing cards (`╭`, `╮`, `╰`, `╯`, `├`, `┤`, `─`, `│`), 5 staged execution cards (`[1/5]` to `[5/5]`), interactive Braille spinner (`⠋` to `⠏`), and completion status card.
  - Modernized `install.sh`: added Lipgloss-style card renderers, Truecolor palette, Braille spinner engine, and 5 stages.
  - Added `$GITMAP_DOWNLOAD_URL` / `$env:GITMAP_DOWNLOAD_URL` override support across both scripts.
  - Created standalone offline local installers `local-install.ps1` and `local-install.sh` supporting pre-built `./bin/` binary lookup and local `./cli` compilation fallback with identical Lipgloss ANSI UI.

### Task-03: Modernized Terminal UI for `gitmap update` Command
- **Target Files:** `cli/cmdupdate/update_ui.go`, `cli/cmdupdate/update_ui_types.go`, `cli/cmdupdate/update_binary_swap.go`
- **Accomplishments:**
  - Implemented `updateUIRendererImpl` with `pterm.DefaultMultiPrinter` and animated spinners.
  - Provided Lipgloss rounded box cards for update summary and fast-path already-updated notifications.
  - Supported 5 stages: Resolving -> Downloading (with progress reader) -> Verifying -> Swapping -> Complete.
  - Graceful degradation: `--json` emits pure structured JSON; `--quiet` runs silently; non-TTY / dumb terminals fall back to simple line-by-line output.
  - Wired into `runBinaryReleaseUpdate` in `cli/cmdupdate/update_binary_swap.go`.

### Task-04: E2E Testing, Minor Version Bump & Release Ceremony
- **Target Files:** `version.json`, `package.json`, `cli/constants/constants.go`, `changelog.md`, `readme.md`, `what-to-read.md`
- **Accomplishments:**
  - Validated dry-runs across all installer scripts (`install.ps1 -DryRun`, `local-install.ps1 -DryRun`, `install.sh --dry-run`, `local-install.sh --dry-run`).
  - Executed minor version bump from `v6.556.1` to `v6.557.0` using `37-bump-version.py --tier minor`.
  - Verified manifest synchronization (`14-version-sync-checker.py` PASS).
  - Authored release notes for `[6.557.0] - 2026-10-16` in `changelog.md` using strictly relative paths.
  - Verified zero absolute paths via `check-relative-paths.py` (PASS across 8,845 files).
  - Verified secrets gate clean via `check-forbidden-strings.py`.
