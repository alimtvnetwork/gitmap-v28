# Issue 47 RCA: Antigravity Rerun Process Termination and Inverted Project Selection

## 1. Reproduction

1. **Destructive Termination on Rerun:**
   - Active Antigravity IDE process running (`Antigravity-default-copy-4978.exe` PID 1596 or 9804) with multiple pinned conversations and open workspace (`wp-exam`).
   - Run `gitmap agy rerun 1` from `D:\work`.
   - Output reports:
     ```text
     [1/2] Terminating active Antigravity IDE process...
       ✔ Terminated Antigravity IDE (PID: 1596)
     [2/2] Launching Antigravity IDE in d:\work\wp-html-automate...
     ```
   - Antigravity IDE is forcefully terminated via `taskkill /F /PID` and `taskkill /F /IM Antigravity.exe`, closing the user's active window and losing work state.

2. **Inverted Project Selection:**
   - User was working in `wp-exam` (active and pinned in Antigravity IDE).
   - Rerun target `#1` resolved to `wp-html-automate` instead of `wp-exam`.
   - Replaying prompt re-targeted an unrelated project without user consent.

## 2. Root Cause Analysis

1. **`rerunRestartFlag` Defaulted to True:**
   - In `cli/cmdagy/agy_rerun.go`:
     `agyRerunCmd.Flags().BoolVarP(&rerunRestartFlag, "restart", "r", true, ...)`
   - Every execution of `gitmap agy rerun` defaulted `isRestart` to `true`, directly invoking `executeIDERestart(&plan)`.
   - `killIDEProcessByPID` executed destructive `taskkill` commands instead of using the running IDE's `agentapi` socket.

2. **Static JSON File Timestamp Sorting in `loadSortedProjects`:**
   - `loadSortedProjects` in `cli/cmdagy/agy_rerun_restart.go` loaded `~/.gemini/config/projects/*.json` and sorted strictly by `p.UpdatedAt`.
   - Project JSON files only store creation timestamps (e.g. August 2026), never updated when conversations run or user interactions occur.
   - `wp-html-automate` happened to have a newer static JSON date than `wp-exam`, causing it to be ranked at position 1.
   - The engine completely ignored `conversation_summaries.db` (the actual source of truth for runtime activity) and `app_storage.json` (the source of truth for pinned conversation order).

3. **Rigid Process Name Detection in `agy_open_paths.go`:**
   - `extractIDEProcessRecord` compared process names with `strings.EqualFold(filepath.Base(name), "antigravity.exe")`.
   - On Windows, updated or multi-instance Antigravity processes run under names like `Antigravity-default-copy-4978.exe`.
   - Process detection returned `false` ("process is not running"), failing to identify active instances.

4. **Missing Dedicated Conversation Commands:**
   - GitMap lacked high-level CLI commands for creating new conversations (`agentapi new-conversation`) and injecting messages (`agentapi send-message`), preventing clean non-destructive automation.

## 3. Corrective Implementation

1. **Non-Destructive Rerun (`cli/cmdagy/agy_rerun.go`, `cli/cmdagy/agy_rerun_restart.go`):**
   - Changed `--restart` (`-r`) to default to `false`.
   - Added `--new-conversation` (`-n`) flag defaulting to `true` to create a fresh conversation session via `AgentAPINewConversationWithOptions`.
   - When `--restart` is false, `executeIDERestart` is skipped, preserving the running IDE.
   - Formatted media attachments with both structured file references and markdown image embeds (`![Picture](...)`).

2. **Activity & Pinned Project Resolution (`cli/cmdagy/agy_rerun_project_resolve.go`):**
   - Implemented `loadActiveSortedProjects()`:
     - Excludes projects with non-existent workspace paths.
     - Retrieves `MAX(last_modified_time)` from `conversation_summaries.db`.
     - Retrieves `pinned_conversations_order` from `app_storage.json`.
     - Ranks pinned projects first in their exact sidebar order, followed by recent conversation activity, followed by project JSON timestamps.
   - Implemented `resolveClosestActiveProject`:
     - Resolves `cwd` when run within a workspace directory.
     - Supports sequence indexing (`1`, `2`, ...).
     - Supports flexible slug and phonetic matching (mapping `"wp-xampp"`, `"wpexam"`, `"exam"` to `"wp-exam"`).

3. **Flexible Process Detection (`cli/cmdagy/agy_open_paths.go`):**
   - Implemented `isAntigravityProcessName()` to match `antigravity.exe`, `language_server.exe`, and any executable starting with `antigravity` and ending with `.exe`.

4. **Dedicated Conversation Commands (`cli/cmdagy/agy_agentapi_cmds.go`):**
   - Implemented `gitmap agy new-conversation` (`new-conv`, `nc`) with `--model`, `--title`, `--profile`, and `--file`.
   - Implemented `gitmap agy send-message` (`send`, `msg`) with `--first`, `--title`, `--file`, automatically targeting the latest active conversation when recipient is omitted or specified as `"first"` / `"latest"`.

## 4. Verification & Prevention

1. **Unit & Integration Tests (`cli/cmdagy/agy_rerun_test.go`, `cli/cmdagy/agy_agentapi_cmds_test.go`):**
   - Verified that `agy rerun -d 1` selects `#1` from pinned/active projects without killing the IDE.
   - Verified that `agy rerun -d wp-xampp` resolves to `wp-exam`.
   - Verified that running `agy rerun -d` inside `d:\work\wp-exam` selects `wp-exam`.
   - Verified argument formatting and option propagation for `new-conversation` and `send-message`.
2. **Quality Gates:**
   - `check-nested-ifs.py`: PASS (0 violations).
   - `check-enum-and-boolean.py`: PASS.
   - `check-error-management.py`: PASS.
   - `check-relative-paths.py`: PASS across 7716 files.
