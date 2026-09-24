# Spec 144: AGY Rerun IDE Restart, Multi-Project Prompt Replay with Media, and VS Code Repair Integration

> **/goal** Provide automated Antigravity IDE process restart, multi-project index-based prompt replay with full media/image attachments, end-to-end testing, and complete VS Code repair integration with zero hardcoded paths.
> **/learn** Enforce coding guidelines, PascalCase types, AppError wrapping, relative paths, and isolated `//go:build e2e` test suites.

## 1. Overview & Objective

When multiple projects are active or running commands in Google Antigravity, developers need the ability to immediately rerun the last prompt of a specific project (e.g. project 1, 2, 3, or 4) while restarting the Antigravity IDE process. Rerunning must extract the exact active prompt payload—including text and attached media/pictures—and dispatch it back to that project's active conversation automatically.

Furthermore, this specification details the architectural integration and root cause analysis (RCA) of the VS Code startup crash and Project Manager JSON repairs, ensuring dynamic path resolution across all scripts and platforms with zero hardcoded drive letters.

---

## 2. AGY Rerun with IDE Restart Architecture

### 2.1 Command Interface & Syntax

```bash
gitmap agy rerun [1|2|3|4|project-name] [flags]
gitmap agy rerun-restart [1|2|3|4] [flags]
gitmap agy rr [1|2|3|4] [flags]
```

### 2.2 Flags & Options

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--restart` | `-r` | boolean | `true` | Restart the Antigravity IDE process during prompt rerun |
| `--no-restart` | | boolean | `false` | Replay prompt directly without restarting the IDE |
| `--project` | `-P` | string | `""` | Target project by index (`1`, `2`, ...) or slug/name |
| `--conversation` | `-c` | string | `""` | Target conversation ID (auto-resolved if empty) |
| `--prompt` | `-p` | string | `""` | Optional prompt template prefix (e.g. `is-done`) |
| `--dry-run` | `-d` | boolean | `false` | Preview restart actions and extracted payload without executing |

### 2.3 Execution Flow

```text
User executes: gitmap agy rerun 1
       │
       ▼
1. Resolve Target Project
   - Load projects from ~/.gemini/config/projects/
   - Sort by recent / pinned
   - Index 1 resolves to projects[0]
       │
       ▼
2. Extract Last User Prompt & Media
   - Select matching conversation via SelectMatchingConversation(projectPath)
   - Open conversation transcript (transcript_full.jsonl / transcript.jsonl)
   - Read last USER_INPUT step:
     * Extract prompt text (clean of metadata tags)
     * Extract media array: mime_type, uri (pictures / screenshots)
   - Reconstruct complete payload preserving image links and references
       │
       ▼
3. Restart Antigravity IDE Process
   - Detect running IDE: DetectRunningAntigravityIDE()
   - Terminate IDE process gracefully (taskkill / pkill)
   - Launch Antigravity IDE with target project workspace
   - Wait brief stabilization window (500ms - 1s)
       │
       ▼
4. Dispatch / Replay Prompt
   - Stage active prompt in project (.ai-memory/prompts/active-agy-prompt.md)
   - Inject into target conversation via AgentAPI / Queue protocol
   - Copy payload to system clipboard as safety fallback
   - Output rich summary (Project, ConvID, Media Count, Restart Status)
```

---

## 3. Media & Picture Attachment Extraction

### 3.1 Transcript Media Data Model

Antigravity transcripts store image and media attachments in the `media` JSON array of `USER_INPUT` steps:

```json
{
  "step_index": 383,
  "source": "USER_EXPLICIT",
  "type": "USER_INPUT",
  "status": "DONE",
  "created_at": "2026-09-24T08:32:01Z",
  "content": "<USER_REQUEST>...\n</USER_REQUEST>",
  "media": [
    {
      "mime_type": "image/png",
      "uri": "C:/Users/Administrator/.gemini/antigravity/brain/.../media_1.png"
    }
  ]
}
```

### 3.2 Payload Reconstruction

When replaying a prompt that contains attached pictures:
1. All media URIs are parsed and verified on the local filesystem.
2. Relative or local markdown image references (`![Attached Image](<path>)`) are preserved or appended to the prompt payload so that the re-invoked AI agent receives the exact image context.
3. Media attachment metadata is reported in terminal output and preserved in the staged prompt file.

---

## 4. Root Cause Analysis (RCA): VS Code Startup Failure & Search Latency

### 4.1 Technical Root Cause (ICU Breakpoint 0x80000003)
- An interrupted background update attempted to upgrade VS Code from commit `7debcd0e2a` to commit `2242ebbb54`.
- Because running VS Code processes (`Code.exe`, `code-tunnel`) were holding file locks, the installer updated subdirectories but could not overwrite `Code.exe`.
- When launched, `Code.exe` (built for `7debcd0e2a`) looked for `icudtl.dat` in folder `7debcd0e2a`. Because only `2242ebbb54` existed, Chromium's `icu_util.cc:232` threw a fatal breakpoint.

### 4.2 Latency Explanation (Why Search Took 1 Hour)
1. **Initial Misdirection:** Exploration initially searched `product.json` and general settings rather than inspecting Chromium crash dumps and the extension storage directory (`AppData\Roaming\Code\User\globalStorage\alefragnani.project-manager\projects.json`).
2. **Versioned Update Discovery:** Windows VS Code utilizes a dual-path structure (System in `Program Files` and User in `%LOCALAPPDATA%\Programs`). Discovering the mismatch between the root `Code.exe` commit SHA and the versioned subfolder required binary disassembly and directory diffing.

### 4.3 Zero Hardcoded Paths Rule
All repair scripts (`tools/repair-vscode-and-projects.ps1`, Go utilities) must calculate paths dynamically:
- `$env:APPDATA`, `$env:LOCALAPPDATA`, `$env:ProgramFiles`
- `filepath.Join`, `os.UserHomeDir()`, `os.Getenv()`
- Strict avoidance of hardcoded drive letters (`D:\`).

---

## 5. End-to-End Test Suite (`//go:build e2e`)

All E2E tests are strictly isolated under `cli/tests/e2e/agy_rerun_restart_e2e_test.go` with the build tag:
```go
//go:build e2e
```
Key test cases:
1. `TestAgyRerun_ProjectIndexResolution`: Verifies numeric indices `1`, `2`, `3` resolve to correct projects.
2. `TestAgyRerun_MediaAttachmentExtraction`: Verifies transcript parsing extracts both prompt text and attached media URIs.
3. `TestAgyRerun_PayloadPreservation`: Verifies reconstructed prompt preserves all pictures and text.
4. `TestAgyRerun_RestartSimulation`: Verifies mock IDE process termination, relaunch command synthesis, and queue staging.

---

## 6. Acceptance Criteria

- **AC-144-01**: `gitmap agy rerun 1` resolves project #1, extracts prompt + media, and triggers IDE restart.
- **AC-144-02**: All attached pictures/media in transcript steps are extracted and formatted into the rerun payload.
- **AC-144-03**: VS Code repair utility and PowerShell helpers calculate paths dynamically with zero `D:\` occurrences.
- **AC-144-04**: E2E tests are gated with `//go:build e2e` and pass cleanly.
- **AC-144-05**: Terminal help text updated in `cli/helptext/agy.md`, `cli/helptext/vscode-repair.md`, and `cli/helptext/catalog.go`.
