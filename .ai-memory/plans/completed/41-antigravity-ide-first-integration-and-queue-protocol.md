# Completed Plan 41: Antigravity IDE-First Integration, File-System Discovery & Queue Protocol

> **Initial Trigger:** User prompt emphasizing that Antigravity IDE (`Antigravity.exe`), not the CLI (`agy.exe`), must be the primary target for all `gitmap agy ...` commands, prompt injections, and pipeline fixes. The system must ground in the local `.gemini/antigravity` filesystem (`brain/`, `transcript.jsonl`), detect active vs idle conversations, and queue prompts when running.
> **Total Steps / Self-Loop Iterations:** 2 Phases, 4 Granular Subtasks across 2 Parallel Execution Subagents.
> **Status:** COMPLETED (100%)

---

## 1. Executive Summary & Deliverables

All deliverables requested by the user were fully implemented:
1. **IDE-First Discovery & Invocation:**
   - Inverted discovery priority so that Antigravity desktop IDE (`Antigravity.exe` at `%LOCALAPPDATA%\Programs\antigravity\Antigravity.exe` on Windows, `.app` bundle on macOS, and native/snap/flatpak binaries on Linux) is checked and prioritized first.
   - Suppressed automatic background CLI runner spawning (`agy.exe`) when an Antigravity IDE instance is detected or active.
2. **File-System Brain & Conversation Discovery:**
   - Grounded conversation and workspace matching in `~/.gemini/antigravity/` (`conversations/<conv-id>.db`, `brain/<conv-id>/.system_generated/logs/transcript.jsonl` and `transcript_full.jsonl`).
   - Added read-only mode (`mode=ro`) for SQLite database queries to avoid WAL lock contention with the running IDE.
   - Added Tier 2 fallback to inspect `transcript_full.jsonl` / `transcript.jsonl` for workspace URIs when SQLite metadata is empty or locked.
3. **Workspace Matching & Unix Path Fixes:**
   - Fixed `cleanURIStringToPath` in `cli/cmdagy/agy_conv_scanner.go` to preserve leading `/` on Unix absolute paths.
   - Fixed `isConvPathMatch` to enforce directory boundary checks using `filepath.Separator`, preventing false prefix matches across similar project names.
4. **Execution State Detection & Dynamic Queue Protocol:**
   - Implemented `CheckConversationStatus` inspecting the last line of `transcript.jsonl` (`USER_INPUT`, `MODEL` with tool calls, or `GENERIC` -> `RUNNING`; `MODEL` without tool calls -> `IDLE`).
   - If conversation is `RUNNING`: enqueues the prompt into `.ai-memory/temp/agy-prompt-queue.json` with status `queued`.
   - If conversation is `IDLE`: delivers the prompt immediately (stages to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`, copies to OS clipboard, marks `dispatched`).
5. **CLI Wiring & Commands:**
   - Implemented `gitmap agy prompt [slug/id/path] [prompt-text]` using the Queue & Injection Protocol.
   - Added `gitmap agy queue` command group (`status`, `ls`, `clear`, `pop`).
   - Updated `gitmap agy status` to display IDE process state, active conversation status (`IDLE` vs `RUNNING`), and pending prompt queue count.
   - Updated helptext in `cli/helptext/agy-fix-pipeline.md` and `cli/helptext/agy.md`.

---

## 2. Consolidated Subtasks

### Subtask 01: Antigravity IDE Path Resolution & Multi-OS Parity
- **Target Files:** `cli/cmdagy/agy_open_paths.go`, `cli/cmdagy/agy_open.go`
- **Accomplishments:**
  - Added safe fallbacks resolving from `os.UserHomeDir()` for Windows when `LOCALAPPDATA` or `ProgramFiles` environment variables are unset.
  - Added Linux candidate paths: `/usr/bin/antigravity`, `/snap/bin/antigravity`, `/var/lib/flatpak/exports/bin/antigravity`, `/opt/Antigravity/antigravity`.
  - In `launchAntigravityProcess()`, added macOS `.app` bundle execution via `open -a <path> <target>` and direct execution if inner `Contents/MacOS/Antigravity` exists.

### Subtask 02: Conversation Scanner Workspace Matching & Transcript Fallbacks
- **Target Files:** `cli/cmdagy/agy_conv_scanner.go`, `cli/cmdagy/agy_conv_ws.go`, `cli/cmdagy/agy_transcript_reader.go`, `cli/cmdagy/agy_transcript_parse.go`
- **Accomplishments:**
  - Fixed `cleanURIStringToPath`: checked if `decoded[1] == ':'`; if not and `rawURI` starts with `file:///`, prepends `"/"` to retain absolute Unix paths.
  - Fixed `isConvPathMatch`: enforced directory boundary check using `filepath.Separator`.
  - Added Tier 2 transcript-based workspace extraction fallback (`transcript_full.jsonl` / `transcript.jsonl`) when SQLite DB `trajectory_metadata_blob` is empty or locked in WAL mode.
  - Opened SQLite DB with read-only mode (`mode=ro`) to prevent lock contention.
  - Cleaned prompt text by stripping `<ADDITIONAL_METADATA>` and `<USER_SETTINGS_CHANGE>` blocks.

### Subtask 03: IDE Status Detection & Dynamic Queue Protocol
- **Target Files:** `cli/cmdagy/agy_fix_pipeline_types.go`, `cli/cmdagy/agy_fix_pipeline_queue.go`, `cli/cmdagy/agy_fix_pipeline_inject.go`
- **Accomplishments:**
  - Added `AgyInjectionModeQueued` enum and `AgyConvStatusType` (`AgyConvStatusIdle`, `AgyConvStatusRunning`, `AgyConvStatusUnknown`).
  - Implemented queue operations: `LoadPromptQueue`, `SavePromptQueue`, `EnqueuePrompt`, `GetQueueStatus`, `PopNextQueuedPrompt`, `ClearPromptQueue`.
  - Inverted priority in `InjectAgyFixTask`: IDE detection and active conversation status checked first.
  - Suppressed background `agy.exe` runner when an IDE process is active.
  - Added `InjectAgyPrompt` for general prompt injection.

### Subtask 04: CLI Wiring, agy Prompt/Queue Commands & Helptext
- **Target Files:** `cli/cmdagy/agy_cmd.go`, `cli/cmdpipeline/pipeline_fix_agy_runner.go`, `cli/cmdagy/agy_fix_pipeline_exec.go`, `cli/helptext/agy-fix-pipeline.md`, `cli/helptext/agy.md`
- **Accomplishments:**
  - Implemented `gitmap agy prompt [slug/id/path] [prompt-text]` via Queue & Injection Protocol.
  - Added `gitmap agy queue` command group with `status`, `ls`, `clear`, and `pop`.
  - Updated `gitmap agy status` to show IDE process, conversation status (`IDLE` vs `RUNNING`), and pending queue count.
  - Refactored `normalizeAgySubcommand` and `init()` into modular functions <= 15 lines.
  - Updated helptext in `cli/helptext/agy-fix-pipeline.md` and `cli/helptext/agy.md`.

---

## 3. Verification & Compliance Checklist
- [x] All functions <= 15 lines.
- [x] Affirmative booleans only (`is*`, `has*`), zero negations.
- [x] Universal `*apperror.AppError` wrappers.
- [x] Unix LF line endings preserved.
- [x] Zero test/build commands run during execution turn.
