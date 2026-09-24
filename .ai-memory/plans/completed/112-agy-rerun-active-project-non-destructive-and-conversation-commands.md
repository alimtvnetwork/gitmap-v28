# Plan 112: Antigravity Rerun Active Project Resolution, Non-Destructive Replay, and Conversation Commands

## Metadata
- Spec Reference: [02-spec/21-app/161-agy-rerun-active-project-non-destructive-and-conversation-commands.md](../../../02-spec/21-app/161-agy-rerun-active-project-non-destructive-and-conversation-commands.md)
- RCA Reference: [02-spec/22-app-issues/47-agy-rerun-killed-ide-and-selected-wrong-project-rca.md](../../../02-spec/22-app-issues/47-agy-rerun-killed-ide-and-selected-wrong-project-rca.md)
- Status: Completed
- Duration / Cycles: 1 cycle (Parent Task N-Step Loop v2.5.0)

## Overview & Scope
Resolved `gitmap agy rerun` process termination and inverted project selection:
1. `gitmap agy rerun 1` forcefully terminated the active Antigravity IDE process (`taskkill /F /PID` and `taskkill /F /IM Antigravity.exe`) because `--restart` was enabled by default.
2. `loadSortedProjects` sorted exclusively by static `~/.gemini/config/projects/*.json` creation timestamps (`UpdatedAt`), failing to identify active or pinned projects (selecting `wp-html-automate` instead of `wp-exam`).
3. Process detection failed to recognize updated or copied Antigravity binaries (`Antigravity-default-copy-4978.exe`).
4. Added dedicated Antigravity CLI commands for programmatic session management:
   - `gitmap agy new-conversation` (aliases: `new-conv`, `nc`)
   - `gitmap agy send-message` (aliases: `send`, `msg`)

## Outcomes & Verification
- Updated `cli/cmdagy/agy_rerun.go` so `--restart` defaults to `false` and added `--new-conversation` defaulting to `true` with `--model`.
- Updated `cli/cmdagy/agy_rerun_restart.go` to invoke `loadActiveSortedProjects()` and replay without terminating the running IDE process.
- Created `cli/cmdagy/agy_rerun_project_resolve.go`:
  - Dynamically extracts latest conversation activity from `conversation_summaries.db` (`MAX(last_modified_time)`).
  - Retrieves `pinned_conversations_order` from `app_storage.json` to mirror the Antigravity IDE sidebar.
  - Implements flexible target matching supporting sequence numbers (`1`, `2`), cwd detection, and phonetic/alias resolution (`wp-xampp` -> `wp-exam`).
- Created `cli/cmdagy/agy_agentapi_cmds.go`:
  - Implemented `gitmap agy new-conversation` with `--model`, `--title`, `--profile`, and `--file`.
  - Implemented `gitmap agy send-message` with `--first`, `--title`, and `--file`, automatically resolving the latest active conversation ID.
- Updated `cli/cmdagy/agy_open_paths.go` with `isAntigravityProcessName()` to detect process variations like `Antigravity-default-copy-4978.exe`.
- Formatted prompt payload with markdown image embeds (`![Picture](...)`) alongside structured file references.
- Verified test suite: all 8 new tests pass (`cli/cmdagy/agy_rerun_test.go` and `cli/cmdagy/agy_agentapi_cmds_test.go`).
- Verified linters: 0 violations across nested ifs, booleans, error management, and relative paths.
- Verified live CLI dry-runs:
  - `gitmap agy rerun -d 1` selects `#1` without restart prompts.
  - `gitmap agy rerun -d wp-xampp` resolves to `wp-exam`.
  - Inside `d:\work\wp-exam`, `gitmap agy rerun -d` automatically selects `wp-exam`.
