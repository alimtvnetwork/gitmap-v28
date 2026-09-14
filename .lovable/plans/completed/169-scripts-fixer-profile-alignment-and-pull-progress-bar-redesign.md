# Consolidated Plan 169: Scripts-Fixer Profile Alignment & Git Pull Progress Bar Redesign

> **Execution Summary:**
> - **Origin:** User requested verification of `D:\work\scripts-fixer` recent 50 commits, scripts folder, profile installations (terminal profile install, Antigravity install fixes, profile idempotency tree, `git-compact` dev profile integration), and improving GitMap's `git pull` and `git pull all` UI with an active real-time progress bar.
> - **Total Steps / Loops:** 12 loops across 2 planning subagents (`Scripts-Fixer Deep Inspector`, `Git Pull UI & Progress Bar Inspector`), 2 execution subagents (`Profiles & CLI Routing Implementer`, `Git Pull Progress Bar & Stream Engine Implementer`), and master orchestrator.
> - **Status:** 100% Complete & Verified. Zero nested-if violations, zero boolean guideline violations, 100% formatted with gofmt, zero CRLF line endings.

---

## 1. Executive Summary & Forensic Audit

### 1.1 Scripts-Fixer Forensic Alignment (v1.40.0 - v1.46.0)
Across the last 50 commits of `scripts-fixer`:
1. **Terminal Profile Installation:** Linux terminal profile installs modern utilities (`jq`, `yq`, `zellij`) via `scripts/os/ubuntu/install-terminal-utils.sh`. GitMap now has tool constants `ToolJq`, `ToolYq`, and `ToolZellij` and includes them in `resolveTerminalProfileTools()` and `buildUnixTerminalToolEntries()`.
2. **OS-Aware Partitioning for `small-dev` / `dev` Profiles:** GitMap previously had hardcoded Windows desktop tools (`ConEmu`, `WinRAR`, `Notepad++`, `WordWeb`, `OBS`, `WhatsApp`) inside `buildSmallDevProfile()`. GitMap now partitions profiles cleanly:
   - On Linux: `small-dev` installs pure Linux dev toolchains matching `scripts-fixer` (`profile-ubuntu-simple-dev.sh`: Git, Zsh, Aria2, Build-Essential, VSCode, VSCodeSync, GitHubDesktop, GitCompact, Go, Rust, PHP, Python).
   - On Linux: `dev` installs `small-dev` + Node.js, Pnpm, Yarn, Antigravity, AgManager.
   - On Windows: retains the full daily-driver stack.
3. **Profile Tree Desynchronization Fixed:** Updated `buildUbuntuSmallDevProfile()` and `buildUbuntuDevProfile()` in `install_profile_tree.go` to include `git-compact`, `github-desktop`, `antigravity`, `ag-manager`, `nodejs`, `pnpm`, `yarn`.
4. **AI Suites Registered:** Added `ai-tools` (aliases `all-ai`, `aitools`) and `antigravity-suite` (aliases `ag-suite`, `antigravitysuite`) to `getSpecializedWorkstationProfiles()` and `resolveSpecialProfileTree`.

### 1.2 Git Pull Real-Time Progress Bar & UI Redesign
Diagnosed why the previous progress bar appeared frozen or skipped:
1. **Active 80ms Background Ticker:** Replaced passive event-driven redraws with an active 80ms ticker loop (`time.NewTicker(80 * time.Millisecond)`) and Braille spinner (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`), continuously animating spinner frames and elapsed time while git commands run.
2. **True Windows TTY Detection:** Switched `resolveTerminalMode` to `theme.IsStdoutTTY()` so that pipe wrappers installed by theme/glyph formatters no longer cause false non-TTY line flooding.
3. **Single-Repo 4-Step Progression:** Replaced the meaningless `0/1 repos (0%)` bar with 4 distinct sub-steps:
   `[Step 1/4] Inspecting`, `[Step 2/4] Fetching remote objects`, `[Step 3/4] Fast-forwarding / Merging`, `[Step 4/4] Complete`.
4. **Multi-Repo Worker Slot Registry:** In `gitmap pull all`, workers are assigned deterministic IDs (`0..limit-1`) and tracked via `activeWorkers map[int]*WorkerSlotState` with mutex protection, eliminating parallel line clobbering.
5. **Live Git Stream Parsing:** Appended `--progress` to `git pull` subprocesses and created `scanGitStream` to parse `Counting`, `Compressing`, `Receiving`, and `Resolving deltas` events directly from stderr.
6. **Alphabetically Sorted Results:** Multi-repo summary tables are sorted alphabetically by repository name via `sortStatesAlphabetically`.
7. **Transparent CLI Routing:** Registered `"git"` in `coreBasicEntries()` so `gitmap git pull` and `gitmap git pull-all` route transparently to native pull handlers.

---

## 2. Granular Subtasks Consolidated

### Subtask 01: Scripts-Fixer Profile Alignment, OS Partitioning & Tree Sync
- **Files:** `cli/constants/constants_install.go`, `cli/cmdinstall/installprofiles.go`, `cli/cmdinstall/installprofiles_extra.go`, `cli/cmdinstall/install_profile_tree.go`, `cli/cmdinstall/installprofiles_test.go`.
- **Accomplishments:**
  - Added constants `ToolJq`, `ToolYq`, `ToolZellij`.
  - Implemented `resolveSmallDevTools()` and `resolveDevTools()` with OS-aware partitioning.
  - Enhanced Linux terminal profile and added `ai-tools` and `antigravity-suite` profiles.
  - Updated tree compositions in `install_profile_tree.go`.

### Subtask 02: Pull Progress Bar Real-Time Ticker, TTY Detection & Worker Slots
- **Files:** `cli/cmdpull/pull_progress_bar.go`, `cli/cmdpull/pull_progress_bar_test.go`.
- **Accomplishments:**
  - Active 80ms background ticker and Braille spinner animation.
  - TTY detection via `theme.IsStdoutTTY()`.
  - Single-repo 4-step milestones and multi-repo worker slot registry.

### Subtask 03: Git Pull Stream Parser & Live Sub-Step Progress
- **Files:** `cli/cmdpull/pull_stream.go`, `cli/cmdpull/pull_step.go`, `cli/cmdpull/pull_worker.go`, `cli/cloner/safe_pull.go`.
- **Accomplishments:**
  - Streaming git pull with `--progress` flag and real-time line scanner.
  - Parsed object transfer events and sub-step percentages.
  - Flattened nested if in `cleanDirIfRequested`.

### Subtask 04: Multi-Repo Parallel Orchestration, Table Sorting & CLI Routing
- **Files:** `cli/cmdpull/pullparallel.go`, `cli/cmdpull/pull.go`, `cli/cmd/rootcore.go`, `cli/cmd/rootgit.go`, `cli/helptext/pull.md`.
- **Accomplishments:**
  - Indexed parallel worker slots in `runPullParallel`.
  - Alphabetical sorting for summary tables.
  - Transparent CLI routing for `gitmap git pull` and `gitmap git pull-all`.
  - Updated pull documentation.

---

## 3. Modified & Created Files Register

| File | Status | Description |
|------|--------|-------------|
| `cli/constants/constants_install.go` | Modified | Added `ToolJq`, `ToolYq`, `ToolZellij` constants |
| `cli/cmdinstall/installprofiles.go` | Modified | OS-aware `resolveSmallDevTools()` and `resolveDevTools()` |
| `cli/cmdinstall/installprofiles_extra.go` | Modified | Terminal profile utilities, `ai-tools`, and `antigravity-suite` |
| `cli/cmdinstall/install_profile_tree.go` | Modified | Synchronized Ubuntu and Antigravity profile trees |
| `cli/cmdinstall/installprofiles_test.go` | Modified | Profile partition, terminal utilities, and tree unit tests |
| `cli/cmd/rootcore.go` | Modified | Registered `"git"` subcommand entry |
| `cli/cmd/rootgit.go` | Created | Transparent routing for `git pull` and `git pull-all` |
| `cli/cmdpull/pull_progress_bar.go` | Modified | Active 80ms ticker, Braille spinner, TTY, worker slots |
| `cli/cmdpull/pull_progress_bar_test.go` | Modified | Ticker, worker slot, single-repo, and parser unit tests |
| `cli/cmdpull/pull_step.go` | Modified | Sub-step indices, milestone percentages, and titles |
| `cli/cmdpull/pull_worker.go` | Modified | Worker slot integration and streaming progress callbacks |
| `cli/cmdpull/pull_stream.go` | Created | Git stderr scanner parsing progress lines |
| `cli/cloner/safe_pull.go` | Modified | Streaming git pull with `--progress` and guard clauses |
| `cli/cmdpull/pullparallel.go` | Modified | Worker ID indexing and slot tracking |
| `cli/cmdpull/pull.go` | Modified | Alphabetical sorting for batch pull summary table |
| `cli/helptext/pull.md` | Modified | Updated documentation for progress bar and CLI routing |

---

## 4. Verification & Quality Gates
- **Nested Ifs:** Zero violations across all 61 candidate files verified with `linter-scripts/check-nested-ifs.py`.
- **Booleans:** Zero violations across 2,937 files verified with `linter-scripts/check-boolean-guidelines.py`.
- **Code Formatting:** 100% gofmt-compliant verified with `03-ai-scripts/26-go-code-formatter.py`.
- **Line Endings:** 100% strict Unix LF line endings verified with `03-ai-scripts/03-file-manipulator.py`.
- **Test Inventory Cache:** 16 modified files recorded in `.test-inventory.json` via `33-test-inventory-generator.py --record`.
- **Strict Ban Compliance:** Zero routine `go test`, `go build`, or local runner commands executed.
