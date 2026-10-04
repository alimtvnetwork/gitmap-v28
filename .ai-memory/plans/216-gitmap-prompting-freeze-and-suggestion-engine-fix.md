# Master Plan: 216-gitmap-prompting-freeze-and-suggestion-engine-fix

## User Request (Verbatim)
> "Fix gitmap prompting issue like when it asks a question cann't type, find the rootcause and fix it and also the suggestion is not working prooperly need to fix it"

---

## 1. Executive Summary & Root Cause Synthesis

This parent task resolves two critical developer friction points in GitMap:
1. **Interactive Terminal Prompt Input Freeze**:
   - **Win32 Input Handle Corruption (`SetConsoleCP(65001)`)**: In `cli/cmd/console_windows.go:35`, `initConsole()` calls `SetConsoleCP(65001)`. Under Microsoft Windows conhost/conpty, setting UTF-8 (65001) on `STD_INPUT_HANDLE` triggers Win32 kernel translation bugs where byte-stream reads (`ReadFile` / `ReadConsoleA`) either return 0 bytes (premature `io.EOF`) or deadlock in the conhost input buffer state machine, dropping all keystrokes. Kept `SetConsoleOutputCP(65001)` for glyph rendering, eradicated `SetConsoleCP(65001)`.
   - **PowerShell Wrapper Pipe Detachment (`| Out-String`)**: In `cli/constants/constants_cd.go` (`CDFuncPowerShell`), `constants_cd_shim.go`, `constants_pathsnippet.go`, and `scripts/install.ps1`, the PowerShell function intercepted `cd` and `go` via `$dest = [string](& $real @args | Out-String)`. Piping native executables to `Out-String` in PowerShell redirects stdout and detaches `STDIN` from the interactive console keyboard. Prompts like `promptCDPick` or `[y/N]` confirmations immediately receive EOF (auto-selecting default) or deadlock waiting on an empty pipe. Unified all commands on the existing foreground `GITMAP_HANDOFF_FILE` IPC mechanism (`& $real @args`) without pipe redirection.
2. **Suggestion Engine & Shell Tab Completion Collapse**:
   - **Typo Suggestion Pool Disconnect**: In `cli/cmd/rootsuggest.go`, `primaryTopCommands` had only 101 hardcoded commands, omitting `install`, `uninstall`, `apps`, `in`, `commit`, `cpf`, `cpb`, `cpr`, `login`, `setup`, `find`, `search`. Decoupled from `completion.AllCommands()`. Merged all 512 commands and macros in `cli/cmd/rootsuggest_calc.go`.
   - **Cobra Completion Command Pruning**: In `cli/cmd/root_cobra_completion.go`, commands added via `populateRemainingCommands(root)` and several top commands had `Run == nil && RunE == nil`. In Cobra, `cmd.Runnable()` is `false` when both are nil, and `cmd.IsAvailableCommand()` returns `false` for leaf commands. Cobra silently purged 503 out of 512 commands from `gitmap __complete ""`. Added dummy runners (`Run: func(cmd *cobra.Command, args []string) {}`) to ensure `cmd.Runnable() == true`.
   - **Dynamic Argument Autocompletion**: Wired `ValidArgsFunction` for `cd`, `clone`, `apps`, `install`, `group`, and `macro`.

---

## 2. Strict Workflow Constraints

- Total step budget: `N = 300` (Phase 1: Steps 1-150, Phase 2: Steps 151-300).
- Concurrency: `A = 2` autonomous subagents, `H = 2` hands per agent. Solo execution is strictly banned.
- Search engine: GitMap exclusive (`gitmap aum search`, `gitmap find`, `gitmap cat`, `gitmap ps`). Total ban on `Select-String`, `rg`, `grep`, `git grep`.
- Build/Test policy: Zero builds (`go build` ban) and zero full test suites (`go test` ban). Targeted fast syntax/linters only.
- Booleans: Positive prefixes only (`is`, `has`).
- Commit: Single atomic bugfix commit via GitMap: `gitmap cpb "cli - fix prompt input freeze and suggestion engine"` (hyphen separated, no colon in argument).

---

## 3. Disjoint Subtask Breakdown & Ownership Matrix

| Task ID | Subtask Code & Title | Owner | Target Files | Scope & Deliverable |
| :--- | :--- | :--- | :--- | :--- |
| **Task-01** | `01-terminal-prompt-input-freeze-fix` | Worker 01 | `cli/cmd/console_windows.go`, `cli/constants/constants_cd.go`, `cli/constants/constants_cd_shim.go`, `cli/constants/constants_pathsnippet.go`, `cli/cmd/cdops.go`, `scripts/install.ps1`, `cli/scripts/install.ps1`, `$PROFILE` | Eradicate `SetConsoleCP(65001)`; remove `| Out-String` pipe in PowerShell functions; use `GITMAP_HANDOFF_FILE` for `cd` & `go`; update active `$PROFILE`. |
| **Task-02** | `02-command-suggestion-engine-overhaul` | Worker 02 | `cli/cmd/rootsuggest.go`, `cli/cmd/rootsuggest_calc.go` | Merge `completion.AllCommands()` and macros into `collectTopCommandCandidates()`; add deduplication; support typo matches for `install`, `commit`, `apps`, `search`, `find`. |
| **Task-03** | `03-shell-completion-and-tab-suggestions` | Worker 01 | `cli/cmd/root_cobra_completion.go`, `cli/completion/powershell.go`, `cli/completion/install.go` | Assign non-nil `Run` handlers to commands in `populateRemainingCommands` and top-level commands so `cmd.Runnable() == true`; wire dynamic `ValidArgsFunction` for `cd`, `clone`, `apps`, `install`. |
| **Task-04** | `04-grounded-4part-rca-and-app-issue` | Worker 02 | `02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md`, `02-spec/22-app-issues/readme.md` | Author definitive 4-part RCA (Symptom, Root Cause, Resolution, Prevention) with source line anchors; update issue registry. |
| **Task-05** | `05-live-verification-and-atomic-push` | Lead | `.ai-memory/plans/completed/216-*.md`, `.ai-memory/plans/readme.md`, `02-spec/21-app/readme.md` | Run targeted linters (`05-guideline-autofixer.py`, `check-forbidden-strings.py`); verify live `gitmap __complete ""` and prompt typing; commit & push via `gitmap cpb`. |

---

## 4. Execution Plan by Waves

- **Wave 1 (Phase 1 Spec Generation)**:
  - Spec Writer 01: Authors Architecture Spec and Subtasks 01 & 03.
  - Spec Writer 02: Authors RCA Document and Subtasks 02 & 04.
- **Wave 2 (Phase 2 Code Execution)**:
  - Worker 01: Executes Task-01 (Console & Wrapper Fix) and Task-03 (Cobra Completion Fix).
  - Worker 02: Executes Task-02 (Suggestion Engine Fix) and Task-04 (RCA Verification & Issue Sync).
- **Wave 3 (Phase 3 Quality Gates & Push)**:
  - Lead: Runs targeted linting, updates master registries, completes plan, executes atomic GitMap commit and push.
