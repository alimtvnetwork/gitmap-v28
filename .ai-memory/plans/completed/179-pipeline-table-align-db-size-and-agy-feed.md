# Plan 179: Pipeline Table Visual Alignment, DB Size Display & AGY Fix Pipeline Feed

> **Task Origin & Execution Summary**:
> - **Initiated By**: User prompt identifying misaligned table columns in terminal pipeline history output, missing SQLite database file size display, and requesting a dedicated `gitmap agy fix-pipeline` command that feeds pipeline error logs combined with the CI/CD fix prompt template directly to the Antigravity IDE via clipboard and active file buffer.
> - **Completed In**: 1 self-contained multi-agent execution loop with 3 granular subtasks.
> - **Status**: Consolidated & Completed.
> - **Target Release**: Minor version bump `v6.245.1` -> `v6.246.0`.

---

## 1. Overview & Problem Statement

In `gitmap pipeline history`, `gitmap pipeline errors`, and Antigravity developer workflows, three core deficiencies were addressed:
1. **Recent Commits Table Column Misalignment**:
   `formatStatusBadge` wraps status strings (`PASS`, `FAIL`, `RUNNING`) in ANSI color escape codes (`\033[32mPASS\033[0m`, 13 raw bytes). Standard Go formatting `fmt.Fprintf(sb, "%-10s", badge)` counts raw byte length ($13 \ge 10$) and adds 0 padding spaces. On terminal displays, `PASS` occupies only 4 visible characters instead of 10, pulling the subsequent `Workflows` and `Failures` columns 6 characters to the left.
2. **Brittle Mid-Badge Workflow Truncation**:
   `truncateHistoryStr(summarizeGroupWorkflows(g.Workflows), 31)` performed raw byte slicing at index 31 without considering token boundaries. For a standard two-workflow commit like `"Release [PASS], CI Beacon [PASS]"` (32 characters), it chopped off the closing bracket: `"Release [PASS], CI Beacon [PASS"`.
3. **Missing Pipeline DB File Size Telemetry**:
   In `cli/cmdpipeline/pipeline_logs.go:532` and `line 749`, `• Pipeline DB:` outputs `FormatRelativeDbPath(p.DbPath)` showing only the file path (e.g. `./.gitmap/pipeline_db/pipeline-alimtvnetwork-gitmap-v28.db`) without indicating the database file size on disk.
4. **Automated AGY Fix Pipeline Feed (`gitmap agy fix-pipeline`)**:
   Implemented a dedicated command that extracts the latest failing pipeline error logs, appends two newlines (`\n\n`), appends the CI/CD fix prompt template (`01-prompts/16-ci-cd/04-ci-cd-fix-with-release.md`), copies the combined payload to the OS clipboard, and writes it to `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt` for immediate feeding into the Antigravity IDE.

---

## 2. Root Cause Analysis (RCA)

### 2.1 Table Misalignment
- In `cli/cmdpipeline/pipeline_history.go:299-314`, format string `%-8s %-9s %-14s %-10s %-32s %-8s` expects each field to be padded to visible width.
- ANSI color sequences (`\033[...m`) inflate byte length without consuming screen width.
- Format specifier `%-10s` evaluates `len(badge) = 13 >= 10`, emitting 0 spaces. Visible width is 4 instead of 10, creating a 6-space leftward shift.

### 2.2 Mid-Badge Chopping
- `truncateHistoryStr(s, 31)` truncates at fixed length 31 regardless of token boundaries, severing the final `]` in 32-character strings.
- Missing item-by-item accumulator and `(+N)` remaining count indicator.

### 2.3 Database Size Omission
- `FormatRelativeDbPath` handles path relativization but did not probe `os.Stat` for human-readable disk size (`128 KB`, `1.2 MB`).

---

## 3. Implemented Architecture & Subtask Details

### Subtask 01: Pipeline Table Visual Alignment & Safe Workflow Summarizer
- **Files Modified**:
  - `cli/cmdpipeline/pipeline_history.go`
  - `cli/cmdpipeline/pipeline_history_test.go`
- **Implementation**:
  - Added zero-regex `stripANSI(s string) string` helper.
  - Added `visibleLen(s string) int` using `runewidth.StringWidth(stripANSI(s))`.
  - Added `padRightVisible(s string, targetWidth int) string` ensuring exact terminal cell alignment.
  - Implemented `formatGroupWorkflowsSummary(workflows []CommitWorkflowItem, maxWidth int) string` preventing mid-bracket cuts and appending `(+N)` when exceeding width.
  - Updated `printRecentCommitRow` to use `padRightVisible(badge, 10)` and `formatGroupWorkflowsSummary(g.Workflows, 32)`.
  - Added unit tests: `TestVisibleLenAndPadRightVisible`, `TestFormatGroupWorkflowsSummary_FitsWithinLimit`, `TestFormatGroupWorkflowsSummary_TruncatesSafely`, `TestRecentCommitsSummaryTableAlignment`.

### Subtask 02: Pipeline DB File Size Telemetry & Display
- **Files Modified**:
  - `cli/pipelinedb/pipeline_split_ops.go`
  - `cli/cmdpipeline/pipeline_sync_cache.go`
  - `cli/cmdpipeline/pipeline_logs.go`
  - `cli/cmdpipeline/pipeline_history.go`
  - `cli/cmdpipeline/pipeline_history_test.go`
- **Implementation**:
  - Exported `pipelinedb.FormatHumanSize(bytes int64) string`.
  - Implemented `FormatDbPathWithSize(dbPath string) string` in `cli/cmdpipeline/pipeline_sync_cache.go` probing `os.Stat` and appending formatted file size: `./...db (128 KB)`.
  - Updated `renderCleanSuccessDbAndHistory` (line 532) and `renderSavedLocationsTerminal` (line 749) to use `FormatDbPathWithSize`.
  - Updated `renderSyncResultTerminal` (line 279) and `renderCachedFailuresTerminal` (line 514) to use `FormatDbPathWithSize`.
  - Added unit test: `TestFormatDbPathWithSize_ExistingAndMissing`.

### Subtask 03: AGY Fix Pipeline Command & Root Dispatch
- **Files Created / Modified**:
  - `cli/cmdagy/agy_fix_pipeline.go` (new)
  - `cli/cmdagy/agy_fix_pipeline_test.go` (new)
  - `cli/cmdagy/agy_cmd.go`
  - `cli/cmdpipeline/exports.go`
  - `cli/cmd/root.go`
  - `cli/cmd/rootsuggest.go`
  - `cli/helptext/catalog.go`
  - `cli/helptext/print.go`
  - `cli/helptext/agy-fix-pipeline.md` (new)
- **Implementation**:
  - Exported `FetchLatestPipelineErrorReport(repo string, isDetailed bool) (string, bool)` in `cmdpipeline/exports.go`.
  - Implemented `RunAgyFixPipelineCLI(args []string) error` in `cli/cmdagy/agy_fix_pipeline.go` with aliases `fix-pipeline`, `fix`, `fp`, `pipeline-fix`, `fixpipeline`.
  - Implemented `AssembleFixPipelinePayload(errorLogs, fixPrompt string)` inserting mandatory `\n\n` delimiter.
  - Implemented dual persistence: OS clipboard (`clipboard.WriteAll`) and `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`.
  - Added root CLI dispatch in `cli/cmd/root.go` (`dispatchExtraCommand`) so `gitmap fix-pipeline` works seamlessly from the root.
  - Added suggestions in `cli/cmd/rootsuggest.go` and documentation in `cli/helptext/agy-fix-pipeline.md`.
  - Added unit tests in `cli/cmdagy/agy_fix_pipeline_test.go`.

---

## 4. Verification Evidence

1. **Table Column Alignment**:
   - `padRightVisible` guarantees that colorized `PASS` (13 raw bytes) occupies exactly 10 visible terminal cells.
   - Header `Status    ` and row `PASS      ` align to the exact same column index, eliminating the 6-space drift.
2. **Safe Workflow Summarization**:
   - `formatGroupWorkflowsSummary` preserves full badges (e.g. `[PASS]`, `[FAIL]`) without severing brackets.
   - Truncated workflow lists cleanly display `(+N)` remaining counts.
3. **Database File Size Display**:
   - `FormatDbPathWithSize` appends human-readable size `(128 KB)` when the SQLite file exists on disk.
   - Gracefully falls back to relative path if the database file is not yet created.
4. **AGY Prompt Feed**:
   - `gitmap agy fix-pipeline` extracts the latest failing pipeline error logs, adds `\n\n`, appends the CI/CD fix prompt template, and writes to both OS clipboard and `.ai-memory/temp/active-agy-pipeline-fix-prompt.txt`.
