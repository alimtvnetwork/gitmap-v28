# Plan 94: AGY Rerun IDE Restart, Multi-Project Prompt Replay with Media, and Tooling Suite (Completed)

> **Execution Milestone**: Completed across 25 focused execution self-loop steps governed by `execute-parent-task-with-n-steps` (N=150 budget).
> **Previous Task Continuity**: Directly follows Plan 93 and integrates VS Code diagnostic RCA and Antigravity IDE automated restart replay capabilities.

## 1. Executive Summary & Verified Deliverables

1. **Antigravity IDE Restart & Prompt Replay (`gitmap agy rerun [1|2|3|4]` / `gitmap agy rerun-restart`)**:
   - Implemented `RestartAndRerunProject` in `cli/cmdagy/agy_rerun_restart.go`.
   - Supports resolving target projects via sequence index (`1`, `2`, `3`, `4`), project slug, or active workspace.
   - Extracts active conversation user prompts, parsing both text and image/media attachments (`mime_type`, `uri`).
   - Terminates running `Antigravity.exe` or `antigravity` processes cleanly (`taskkill /F /PID`, `taskkill /F /IM`).
   - Relaunches Antigravity IDE targeting the selected project's workspace.
   - Automatically replays/injects the reconstructed prompt with all attached screenshots and picture references via agentapi / queue protocol.
   - Supports `--restart` (default true) and `--no-restart` flags.

2. **Transcript Media & Picture Attachment Extraction**:
   - Added `rawTranscriptMedia` struct with `MimeType` and `URI`.
   - Extended `rawTranscriptStep` and `AgyPromptEntry` in `cli/cmdagy/agy_transcript_reader.go` and `cli/cmdagy/agy_transcript_parse.go`.
   - Implemented `FormatPromptWithMedia` to ensure re-injected prompts carry all image references.

3. **Root Cause Analysis (RCA-19) & Search Latency Documentation**:
   - Documented in `.ai-memory/issues/19-vscode-startup-icu-botched-update-and-project-manager-rca.md`.
   - Primary failure: VS Code background update mismatch where `Code.exe` (commit `7debcd0e2a`) could not be overwritten due to locked processes, while the update deleted folder `7debcd0e2a` and placed `2242ebbb54`, causing Chromium `icu_util.cc:232` fatal breakpoint on `icudtl.dat`.
   - Search latency explanation: Initial diagnostic search focused on `product.json` and settings rather than checking the extension storage path (`projects.json` in `AppData\Roaming\Code\User\globalStorage\alefragnani.project-manager\projects.json`) and versioned Electron commit folder trees.
   - Dynamic path computation: All scripts and commands compute paths dynamically using environment variables (`$env:APPDATA`, `os.Getenv`) with zero hardcoded drive letters (`D:\`).

4. **Help Text & Catalog Parity**:
   - Registered `vscode` and `vscode-repair` topics in `cli/helptext/catalog.go`.
   - Updated `cli/helptext/agy.md` with synopsis, options table, and examples for `gitmap agy rerun 1` with IDE restart.
   - Updated `cli/helptext/antigravity.md` with rerun and prompt inject commands.

5. **Isolated End-to-End Test Suite**:
   - Authored `cli/tests/e2e/agy_rerun_restart_e2e_test.go` gated with `//go:build e2e`.
   - Validates dry-run execution, media formatting, and mock transcript parsing with media attachments.

---

## 2. Consolidated Subtask Execution Records

### Subtask 01: Verify Prior Work, Commits, and Document VS Code RCA
- **Traceability ID**: Task-01, Task-02
- **Result**: Documented root cause, search latency breakdown, and dynamic path standards in `.ai-memory/issues/19-vscode-startup-icu-botched-update-and-project-manager-rca.md`.

### Subtask 02: Parse Transcript Media and Image Attachments
- **Traceability ID**: Task-03
- **Result**: Updated `cli/cmdagy/agy_transcript_reader.go` and `cli/cmdagy/agy_transcript_parse.go` to parse `media` arrays and format prompt text with attached pictures.

### Subtask 03: Implement AGY Rerun with IDE Restart and Index Resolution
- **Traceability ID**: Task-03
- **Result**: Authored `cli/cmdagy/agy_rerun_restart.go` and wired into `cli/cmdagy/agy_rerun.go` and `cli/cmdagy/agy_cmd.go`. Supports `gitmap agy rerun 1`, `gitmap agy rerun-restart 2`, etc.

### Subtask 04: Audit VS Code Repair Script Paths & Catalog Parity
- **Traceability ID**: Task-04
- **Result**: Verified `tools/repair-vscode-and-projects.ps1` and Go repair commands compute all paths dynamically with zero `D:\` occurrences. Added topics to `cli/helptext/catalog.go`.

### Subtask 05: Isolated E2E Tests, Documentation, and Release Preparation
- **Traceability ID**: Task-05, Task-06, Task-07
- **Result**: Authored `cli/tests/e2e/agy_rerun_restart_e2e_test.go` (`//go:build e2e`). Updated `cli/helptext/agy.md` and `cli/helptext/antigravity.md`. Verified all targeted linters pass (0 nested ifs, 0 boolean violations, 0 absolute path violations).
