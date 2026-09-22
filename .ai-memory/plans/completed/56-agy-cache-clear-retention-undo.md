# Plan 56: AGY Cache Clear Retention, Preflight, Undo & Direct Shortcut Commands

## Task Execution Header
- **Initial Trigger**: User request to implement Antigravity cache clear commands with `--keep 10` default retention, preflight checks (`--pre`, `--precheck`, `--preflight`), `-y` confirmation, temporary undo backup (`gitmap agy undo`), and direct retention shortcuts `ccko` (keep 1) and `cckf` (keep 5).
- **Execution Lifecycle**: Completed in 1 continuous multi-phase loop (Phases 1A, 1B, 2, 3) across 15 modified/created files.
- **Verification Gates**: Nested if linter (0 violations), boolean guidelines linter (0 violations), strict file size limits (<= 100 lines per file).

---

## User Request (Verbatim)
```text
D:\work\scripts-fixer

PowerShell AST Parser: Verified 0 syntax or parsing errors across run.ps1 and all 7 extracted scripts in scripts/dispatcher/.
.\run.ps1 -Help: Verified exit 0; prints full help screen, scripts table, and version footer.
.\run.ps1 help chrome: Verified exit 0; correctly filters and prints 48 matching lines with colored term highlights.
.\run.ps1 path: Verified exit 0; reports active dev directory smart detection state.
.\run.ps1 status: Verified exit 0; outputs the Tool Status Dashboard and outdated package scan.
.\run.ps1 doctor --self-check: Verified exit 0; ran deep validation over 600 assertions.
.\run.ps1 agy check: Verified exit 0; verifies Antigravity CLI and IDE presence.
.\run.ps1 agy clear 10: Verified exit 0; safely tests retention prediction mode keeping top 10 conversations.

Try to understand these commands like AGI-related clear, keep, and other commands. So read it first very carefully, and then I want you to implement this in our CLI. Okay? So it would be like CLI space AGI. Undo, clear, cache clear. I think cache clear is already there. So inside this, you can actually add more commands, and you can have keep 10. 10 would be by default. By default, that means you're going to keep the 10 conversation. And we can also do an undo after the clear so that it keeps back the things that have been done because the clear will keep the removal in the temp directory so that it can revert back shortly. But it is also for a short period of time because temp directory can be removed. So that needs to be mentioned in the help section as well. And that should be a generic process to deal with it, the AGI cache and other stuff. I want you to understand this fully for all the OS. We will do macOS, Windows, all cases. This needs to be there. Also, it needs to be removed carefully because when we remove the conversation, we cannot remove just any conversation. We have to understand the conversation and then keep it. I think by default, we are going to keep yeah, 10, but user can have the option to change it, like keep five, keep three, keep one. There could be one direct command, like clear cache, keep one. I think that could be one. That would be like CCKO or form that is going to keep only one conversation. So that could be another, and we can have keep five. So all this needs to be in the help, UI help, doc CMD, okay, with explanation how it works, which files it removes. I could do a pre-flight as well. I could do pre-check or pre-flight, pre or pre-flight. Either one of the cases, it will tell us where and what it's going to remove. And these are destructive commands, so user needs to do a prompt Y. If they wanted to do it automatic, then they have to provide the hyphen Y, then it would do it automatically. Try to understand this and try to implement in your system properly. Is it clear?

<cli> agy cache-clear --keep 10 # by default
<cli> agy undo

<cli> agy cache-clear-keep-one (ccko)
<cli> agy cache-clear-keep-five (cckf) --precheck/--pre/-preflight
```

---

## Consolidated Subtasks & Delivered Features

### 1. Subtask 01: Enhanced AGY Cache-Clear with Configurable Retention & Preflight (Task-01)
- Added `--keep` / `-k` flag (default 10) to `gitmap agy cache-clear` and `gitmap agy clean-cache`.
- Added `--pre`, `--precheck`, `--preflight` flags simulating cleanup and conversation retention analysis without deleting files.
- Added `-y` / `--yes` confirmation bypass, with interactive `[y/N]` prompt for unconfirmed runs.

### 2. Subtask 02: Direct Retention Convenience Commands (Task-02)
- Added `gitmap agy cache-clear-keep-one` (alias: `ccko`, `cc-keep-one`, `cache-clear-keep-1`).
- Added `gitmap agy cache-clear-keep-five` (alias: `cckf`, `cc-keep-five`, `cache-clear-keep-5`).
- Wired root CLI aliases `gitmap ccko` and `gitmap cckf` directly to AGY command dispatch.

### 3. Subtask 03: Temporary Backup & Safe Undo Mechanism (Task-03)
- Staged all pruned conversations, SQLite databases, and brain artifacts into `os.TempDir()/gitmap-agy-cache-backup/<timestamp>/` with a `manifest.json`.
- Updated `gitmap agy undo` to check for and restore from the latest temporary cache backup directory.
- Added advisory notices that OS temporary directories are ephemeral.

### 4. Subtask 04: Cross-Platform OS Cache & Conversation Retention Engine (Task-04)
- Preserved pinned projects and their conversations unconditionally.
- Sorted conversations by `last_modified_time DESC` and kept top N conversations.
- Maintained cross-platform cache path discovery for Windows, macOS, and Linux.

### 5. Subtask 05: Documentation, UI Help Screens & Command Diagnostics (Task-05)
- Updated `gitmap agy help` diagnostics menu to document `cache-clear`, `ccko`, `cckf`, and `undo`.
- Updated command descriptions and preflight summaries with actionable feedback.
