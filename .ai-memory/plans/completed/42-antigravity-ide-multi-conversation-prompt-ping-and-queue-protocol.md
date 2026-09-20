# Plan 42: Antigravity IDE Multi-Conversation Prompt, Ping & Queue Protocol

## Executive Summary
This plan enforces the **IDE-first, Filesystem-First** architecture for all Antigravity interactions in GitMap:
1. **Total CLI Runner Elimination**: Completely removed `launchAgyBackgroundRunner`, `handleCLIInjectionFallback`, and all background spawning of `agy.exe`. Removed "Antigravity CLI: Detected at..." rendering from verification notices.
2. **Offline & Online Filesystem-First Dispatch**: Decoupled Antigravity IDE operations from requiring an active running process. All operations operate directly against `~/.gemini/antigravity/` and `~/.gemini/config/projects/`. If the conversation is `RUNNING`, prompts are queued to `.ai-memory/temp/agy-prompt-queue.json`. If `IDLE` or `UNKNOWN`, prompts are staged to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` and copied to the OS clipboard.
3. **Multi-Conversation Interactive Selector**:
   - `FindMatchingConversations(repoRoot)` discovers and sorts conversations by user activity descending.
   - `SelectMatchingConversation(repoRoot)` automatically selects single matches, gracefully defaults in non-interactive environments, and interactively prompts the user when multiple conversations match.
4. **Antigravity Ping Command (`gitmap agy ping`)**:
   - Diagnostic checks for executable presence (`C:\Users\Administrator\AppData\Local\Programs\antigravity\Antigravity.exe` and cross-platform paths), running IDE desktop process, filesystem health (`~/.gemini/antigravity`), workspace conversation execution state (`RUNNING` vs `IDLE`), and prompt queue status.
   - Supports `--json` (`-j`) and `--workspace` (`-w`) flags.
5. **Antigravity Prompt Inspection (`gitmap agy prompt read` & `gitmap agy prompt ls`)**:
   - `gitmap agy prompt read [conv-id]`: Reads and renders user prompts from conversation transcripts.
   - `gitmap agy prompt ls [limit]`: Lists prompts across conversations.

---

## Deliverables & Changes

### 1. Multi-Conversation Selection
- `cli/cmdagy/agy_conv_selector.go`:
  - `FindMatchingConversations(repoRoot string) ([]AgyConvInfo, error)`
  - `SelectMatchingConversation(repoRoot string) (AgyConvInfo, error)`
  - `SelectMatchingConversationWithIO(repoRoot string, r io.Reader, w io.Writer, isInteractive bool) (AgyConvInfo, error)`
  - `PromptSelectConversation(convs []AgyConvInfo, r io.Reader, w io.Writer, isInteractive bool) (AgyConvInfo, error)`
- `cli/cmdagy/agy_conv_selector_test.go`:
  - Tests for 0 matches, 1 match, multiple matches non-interactive, multiple matches interactive (indices, empty string, invalid input).

### 2. CLI Elimination & Filesystem Dispatch
- `cli/cmdagy/agy_fix_pipeline_inject.go`:
  - Removed `"os/exec"` import and all CLI runner spawning code (`launchAgyBackgroundRunner`, `handleCLIInjectionFallback`, etc.).
  - Implemented offline/online filesystem dispatch in `InjectAgyFixTask` and `InjectAgyPrompt`.
  - Updated `findMatchingActiveConvID` to use `SelectMatchingConversation`.
- `cli/cmdagy/agy_fix_pipeline_render.go`:
  - Removed `ResolveAntigravityCLI()` check from `renderQueuedVerificationNotice()`.

### 3. Antigravity Ping & Diagnostic Command
- `cli/cmdagy/agy_ping_types.go`:
  - Defined `AgyPingReport`, `AgyPingExecutableCheck`, `AgyPingProcessCheck`, `AgyFilesystemHealth`, `AgyPingWorkspaceCheck`, `AgyPingQueueCheck`.
- `cli/cmdagy/agy_ping.go`:
  - Implemented `agyPingCmd` (`gitmap agy ping`, alias `check`).
  - Implemented `ExecuteAgyPing(wsPath string) AgyPingReport` with decomposed check functions.
- `cli/cmdagy/agy_ping_render.go`:
  - Implemented `renderPingReport` (human terminal) and `renderPingJSON` (JSON).
- `cli/cmdagy/agy_ping_test.go`:
  - Comprehensive unit tests for ping diagnostics.

### 4. Antigravity Prompt Subcommands
- `cli/cmdagy/agy_prompt_subcmds.go`:
  - Registered `agyPromptReadCmd` and `agyPromptLsCmd`.
- `cli/cmdagy/agy_prompt_read.go`:
  - Implemented `RunAgyPromptRead` and `RunAgyPromptLs`.
- `cli/cmdagy/agy_prompt_read_resolve.go`:
  - Implemented `resolvePromptReadConvID`, `resolveActiveOrLatestConvID`, `resolveTargetConvID`.

### 5. CLI Wiring & Helptext
- `cli/cmdagy/agy_cmd.go`:
  - Registered `agyPingCmd` in `registerAgyBaseCommands`.
  - Added `"ping"` and `"check"` aliases to `normalizeMaintenanceSubcommands`.
  - Connected `initAgyPromptSubcommands` in `initAgyPromptAndStatusCommands`.
  - Preserved direct prompt injection (`gitmap agy prompt <slug/text>`) via `isDirectAgyPromptInjection`.
- `cli/cmdagy/agy_help.go`:
  - Added `ping (check)` to diagnostics section in help text.
- `cli/helptext/agy.md` & `cli/helptext/agy-fix-pipeline.md`:
  - Documented `gitmap agy ping`, `gitmap agy prompt read`, and `gitmap agy prompt ls`.

---

## Verification & Guidelines Compliance
- Maximum function length: <= 15 lines.
- Affirmative booleans only (`is*`, `has*`, zero `!` negation operators).
- Universal `*apperror.AppError` return wrappers.
- Blank line before every return.
- Unix LF (`\n`) line endings.
- Zero test runner or build execution during execution turns (strict adherence to ban).
