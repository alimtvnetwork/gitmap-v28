# Spec 161: Antigravity Rerun Active Project Resolution, Non-Destructive Replay, and Conversation Commands

## User Request (Verbatim)

```text
PS D:\work> gitmap agy rerun 1

Antigravity IDE Rerun & Restart Suite:
  • Target Project:   #1 - wp-html-automate
  • Workspace Path:   d:\work\wp-html-automate
  • Conversation:     default

  [1/2] Terminating active Antigravity IDE process...
    ✔ Terminated Antigravity IDE (PID: 1596)
  [2/2] Launching Antigravity IDE in d:\work\wp-html-automate...

  ✔ Injected prompt into active Antigravity session (2911e8ac-d82b-4d4c-af65-4a0b4db4d63c) via agentapi!
  ✔ Completed prompt replay for project #1 (wp-html-automate)!

PS D:\work>

This is also wrong because rerun should run the last closest project which was running. In this case, it was running the WP Xampp. It couldn't figure it out. It didn't rerun. It just gives the wrong message. So I think you should be aware of it. You should fix it properly, okay, and check it thoroughly. Check it here in this machine, try to rerun without closing. Okay, then you'll understand. Rerun the running command again, send it again automatically with image and so on. Create a new conversation and do that
```

## Problem Statement

1. **Destructive Termination**: `gitmap agy rerun` defaulted `--restart` to `true`, aggressively killing the active Antigravity IDE process via `taskkill /F /PID` and `taskkill /F /IM Antigravity.exe`. This disrupted active user workflows and closed running windows.
2. **Inverted/Stale Project Sorting**: `loadSortedProjects()` only loaded static `~/.gemini/config/projects/*.json` files and sorted by the file's static `UpdatedAt` timestamp, ignoring actual running activity in `conversation_summaries.db` and pinned sidebar order in `app_storage.json`. As a result, running `gitmap agy rerun 1` selected `wp-html-automate` instead of the active project (`wp-exam`).
3. **Missing Conversation Commands**: GitMap lacked direct CLI commands exposing `agentapi new-conversation` and `agentapi send-message`, preventing users and agents from programmatically spawning fresh sessions or routing prompts to the active session.
4. **Process Name Inflexibility**: IDE process detection failed to recognize updated or copied Antigravity executables (such as `Antigravity-default-copy-4978.exe`).

## Architecture & Implementation

### 1. Non-Destructive Replay by Default
- The `--restart` (`-r`) flag on `gitmap agy rerun` defaults to `false`.
- When `--restart` is omitted or false, the IDE is never terminated; the prompt and media are injected directly into Antigravity via `agentapi`.
- Flag `--new-conversation` (`-n`) defaults to `true` on `agy rerun`, ensuring a fresh conversation session is created (`AgentAPINewConversationWithOptions`) to preserve history cleanliness.
- Supports `--model` (`-m`) for selecting model tier (`flash_lite`, `flash`, `pro`).

### 2. Activity-Driven Project Resolution
- `loadActiveSortedProjects()`:
  - Filters out projects whose workspace directories do not exist on disk (removing dummy/test folders).
  - Queries `conversation_summaries.db` for `MAX(last_modified_time)` per project and workspace.
  - Queries `app_storage.json` for `pinned_conversations_order` to align with the Antigravity IDE sidebar.
  - Sorts candidate projects: Pinned conversations first, followed by recent conversation activity, followed by project file timestamps.
- `resolveClosestActiveProject(projects, target)`:
  - If target is empty, checks if `os.Getwd()` matches a known workspace.
  - If target is a sequence number (`"1"`, `"2"`), indexes the sorted active project list.
  - If target is a name (e.g. `"wp-exam"`, `"wp-xampp"`, `"exam"`), matches using normalized slug comparison and phonetic/alias rules.

### 3. Dedicated Conversation Commands
- `gitmap agy new-conversation` (aliases: `new-conv`, `nc`):
  - Flags: `--model` (`-m`), `--title` (`-t`), `--profile`, `--file` (`-f`).
  - Creates a new session in Antigravity IDE via `agentapi new-conversation`.
- `gitmap agy send-message` (aliases: `send`, `msg`):
  - Flags: `--title` (`-t`), `--first` (`-1`), `--file` (`-f`).
  - Routes messages to explicit recipient IDs or automatically discovers the latest active conversation ID from `conversation_summaries.db`.

### 4. Resilient Process Matching
- `isAntigravityProcessName()` matches `antigravity.exe`, `language_server.exe`, and executable names with prefix `antigravity` and extension `.exe` (e.g. `Antigravity-default-copy-4978.exe`).

## Verification & Testing

- Unit tests in `cli/cmdagy/agy_rerun_test.go` and `cli/cmdagy/agy_agentapi_cmds_test.go`:
  - `TestResolveClosestActiveProject_FuzzyMatch`: verifies `"wp-xampp"`, `"wpexam"`, `"exam"`, and numeric sequences.
  - `TestResolveClosestActiveProject_CwdMatch`: verifies current working directory matching.
  - `TestSortProjectsByActivityAndPins_Ordering`: verifies pinned and activity-based ordering.
  - `TestRerunDefaultFlags_NonDestructive`: verifies `--restart` is false and `--new-conversation` is true by default.
  - `TestBuildNewConvArgsWithOptions`: verifies CLI arguments construction.
  - `TestIsAntigravityProcessName`: verifies process name resolution.
  - `TestResolveRecipientID_Tokens`: verifies `"first"` and `"latest"` conversation targeting.
- All tests pass, 0 nested-if violations, zero boolean naming infractions.
