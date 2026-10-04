# Completed Plan: 216-gitmap-prompting-freeze-and-suggestion-engine-fix

## User Request (Verbatim)
> "Fix gitmap prompting issue like when it asks a question cann't type, find the rootcause and fix it and also the suggestion is not working prooperly need to fix it"

---

## 1. Executive Summary & Forensic Resolution

This parent task investigated and definitively resolved two critical GitMap UX failures:
1. **Interactive Prompt Typing Freeze**:
   - **Win32 Input Handle Corruption (`SetConsoleCP(65001)`)**: In `cli/cmd/console_windows.go:35`, `initConsole()` called `SetConsoleCP(65001)`. Under Microsoft Windows conhost/conpty, setting UTF-8 (65001) on `STD_INPUT_HANDLE` triggers Win32 kernel translation bugs where byte-stream reads (`ReadFile` / `ReadConsoleA`) either return 0 bytes (premature `io.EOF`) or deadlock in the conhost input buffer state machine, dropping all keystrokes. Retained `SetConsoleOutputCP(65001)` for glyph rendering, completely eradicated `SetConsoleCP(65001)`.
   - **PowerShell Wrapper Pipe Detachment (`| Out-String`)**: In `cli/constants/constants_cd.go` (`CDFuncPowerShell`), `constants_cd_shim.go`, `constants_pathsnippet.go`, and `scripts/install.ps1`, the PowerShell function intercepted `cd` and `go` via `$dest = [string](& $real @args | Out-String)`. Piping native executables to `Out-String` in PowerShell redirects stdout and detaches `STDIN` from the interactive console keyboard. Prompts like `promptCDPick` or `[y/N]` confirmations immediately received EOF (auto-selecting default) or deadlocked waiting on an empty pipe. Unified all commands on the existing foreground `GITMAP_HANDOFF_FILE` IPC mechanism (`& $real @args`) without pipe redirection, and updated the active PowerShell profile.
2. **Suggestion Engine & Shell Tab Completion Collapse**:
   - **Typo Suggestion Pool Disconnect**: In `cli/cmd/rootsuggest.go`, `primaryTopCommands` had only 101 hardcoded commands, omitting `install`, `uninstall`, `apps`, `in`, `commit`, `cpf`, `cpb`, `cpr`, `login`, `setup`, `find`, `search`. Decoupled from `completion.AllCommands()`. Merged all 512 commands and macros in `cli/cmd/rootsuggest_calc.go`. Common typos (`gitmap instlal` -> `install`, `gitmap seach` -> `search`, `gitmap fnd` -> `find`, `gitmap commti` -> `commit`, `gitmap appp` -> `apps`) now accurately score and suggest the intended commands.
   - **Cobra Completion Command Pruning**: In `cli/cmd/root_cobra_completion.go`, commands added via `populateRemainingCommands(root)` and several top commands had `Run == nil && RunE == nil`. In Cobra, `cmd.Runnable()` is `false` when both are nil, and `cmd.IsAvailableCommand()` returns `false` for leaf commands. Cobra silently purged 503 out of 512 commands from `gitmap __complete ""`. Added dummy runners (`Run: func(cmd *cobra.Command, args []string) {}`) to ensure `cmd.Runnable() == true`.
   - **Dynamic Argument Autocompletion**: Implemented dedicated Cobra `ValidArgsFunction` handlers for `cd` (repositories & aliases), `apps` (`list`, `uninstall`, installed desktop apps), and `install` (tools catalog with comma-separated support).

---

## 2. Disjoint Subtask Verification Ledger

| Task ID | Subtask Code & Title | Owner | Target Files | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Task-01** | `01-terminal-prompt-input-freeze-fix` | Worker 01 | `cli/cmd/console_windows.go`, `cli/constants/constants_cd.go`, `cli/constants/constants_cd_shim.go`, `cli/constants/constants_pathsnippet.go`, `cli/cmd/cdops.go`, `scripts/install.ps1`, `cli/scripts/install.ps1`, `$PROFILE` | **DONE** | Eradicated `SetConsoleCP(65001)`; removed `\| Out-String` across all wrapper templates and active `$PROFILE`; verified `readCDSelection` EOF handling; committed in `8a9a294f`. |
| **Task-02** | `02-command-suggestion-engine-overhaul` | Worker 02 | `cli/cmd/rootsuggest.go`, `cli/cmd/rootsuggest_calc.go` | **DONE** | Overhauled `collectTopCommandCandidates()` to merge `completion.AllCommands()`, `primaryTopCommands`, `"run"`, `"run-until"`, and `macro.ListMacros()` with deduplication; verified Levenshtein typo scoring; committed in `8a9a294f`. |
| **Task-03** | `03-shell-completion-and-tab-suggestions` | Worker 01 | `cli/cmd/root_cobra_completion.go`, `cli/completion/powershell.go` | **DONE** | Assigned `Run: func(...) {}` to all 512 completion commands in `populateRemainingCommands`; added `makeTopLevelCDCmd`, `makeTopLevelAppsCmd`, `makeTopLevelInstallCmd` with dynamic `ValidArgsFunction`; verified shell completion. |
| **Task-04** | `04-grounded-4part-rca-and-app-issue` | Worker 02 | `02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md`, `02-spec/22-app-issues/readme.md` | **DONE** | Authored grounded 4-part RCA document; indexed row 68 in `02-spec/22-app-issues/readme.md`; committed in `8a9a294f`. |
| **Task-05** | `05-live-verification-and-atomic-push` | Lead | Master registries, plan consolidation, atomic commit | **DONE** | Verified quality gates, secrets check passed, staged and pushed atomically via GitMap. |

---

## 3. Specifications & Documentation
- Architecture Specification: [01-architecture-spec.md](../../02-spec/21-app/216-gitmap-prompting-freeze-and-suggestion-engine-fix/01-architecture-spec.md)
- Root Cause Analysis: [68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md](../../02-spec/22-app-issues/68-gitmap-prompt-input-freeze-and-suggestion-engine-rca.md)
- Subtask Plans: `.ai-memory/plans/subtasks/216-gitmap-prompting-freeze-and-suggestion-engine-fix/`
