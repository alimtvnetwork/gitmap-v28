# Plan 180: Pipeline DB Repo Location Resolution, Next-Line Size Display & Rust Test Log Filtering

> **Task Origin & Execution Summary**:
> - **Initiated By**: User prompt identifying that pipeline SQLite databases incorrectly defaulted to AppData (`C:/Users/.../AppData/Local/gitmap-cli/data/pipeline_db/...`) when invoked from outside a repo (e.g. `D:\work`), rather than resolving to `.gitmap/data/pipeline.db` relative to the repository installation. Additionally, user required database file size to be displayed on the **next line** immediately below `• Pipeline DB:`. A subsequent prompt instructed that Rust test logs (`test <name> ... ok`) should be recognized as passing lines and filtered out so only failing tests (`test <name> ... FAILED`) and error logs are retained in compact error logs.
> - **Completed In**: 1 self-contained multi-agent execution loop with 2 granular subtasks.
> - **Status**: Consolidated & Completed.
> - **Target Release**: Minor version bump `v6.246.0` -> `v6.247.0`.

---

## 1. Overview & Problem Statement

1. **Pipeline Database Path Relocation**:
   - When executing `gitmap pipeline` from outside a repo root (e.g. `D:\work`), `gitutil.RepoRoot(".")` failed.
   - It previously defaulted immediately to `fallbackBinaryPipelineDbPath(repoSlug)`, dumping databases in `AppData/Local/gitmap-cli/data/pipeline_db/`.
   - The user required database storage to be close to the repo installation (`.gitmap/data/pipeline.db`), never in AppData.
2. **Next-Line Database File Size Display**:
   - Previously, database file size was either omitted or appended in parentheses on the same line.
   - The user required the file size to be clearly presented on the next line:
     ```text
       • Pipeline DB:     .gitmap/data/pipeline.db
       • DB Size:         128 KB
     ```
3. **Rust Test Log Noise Filtering**:
   - In Cargo / Rust test outputs, dozens or hundreds of lines follow `test <path> ... ok`.
   - Previously, `isOkLogLine` only recognized Go and test-runner tokens, causing passing Rust test lines to remain in error reports.
   - `isRustOkLine` now identifies `test <name> ... ok`, `test <name> ... ignored`, and `test result: ok.`, filtering them out so only `test <name> ... FAILED` and error details are visible.

---

## 2. Root Cause Analysis (RCA)

### 2.1 Database Path Resolution
- In `cli/pipelinedb/pipeline_split_db.go`: `ResolvePipelineDbPath` called `gitutil.RepoRoot(".")` and immediately fell back to `fallbackBinaryPipelineDbPath(repoSlug)` if `.` was not a git repo.
- Fixed by introducing candidate folder probing (`filepath.Base(repoSlug)` and stripped version suffixes), database store registry lookup (`store.OpenDefault()` -> `FindBySlug`), and falling back to `./.gitmap/data/pipeline.db`.

### 2.2 Next-Line Size Display
- In `cli/cmdpipeline/pipeline_logs.go`, `pipeline_sync_cache.go`, and `pipeline_history.go`, output lines used `FormatDbPathWithSize` which placed the size in parentheses on the same line or failed silently if unstated.
- Fixed by adding `ResolveDbFileSize` and formatting `• Pipeline DB:` and `• DB Size:` on separate lines with matching indentations.

### 2.3 Rust Test Log Filtering
- In `cli/cmdpipeline/pipeline_error_extract.go`: `isOkLogLine` did not recognize `test ... ok` or `test result: ok.`, and `appendContextLine` allowed passing test lines to be recorded as error context.
- Fixed by adding `isRustOkLine` to `isOkLogLine` and checking `isOkLogLine` in `appendContextLine`.

---

## 3. Implemented Subtasks

### Subtask 01: Pipeline DB Repo Location Resolution (`cli/pipelinedb/`)
- **Files Modified**:
  - `cli/pipelinedb/pipeline_split_db.go`
  - `cli/pipelinedb/pipeline_split_db_test.go`
- **Implementation**:
  - Refactored `ResolvePipelineDbPath` to resolve repo root through:
    1. Current git repo match (`isRepoMatchingSlug`).
    2. Subdirectory candidates (`findCandidateRepoRoot`).
    3. Stored repository registry (`findRepoRootInStore`).
    4. Any active git repo in current directory.
    5. Fallback `./.gitmap/data/pipeline.db` (never AppData).
  - Added unit test `TestResolvePipelineDbPath`.

### Subtask 02: Next-Line DB Size Display & Rust Test Log Filtering (`cli/cmdpipeline/`)
- **Files Modified**:
  - `cli/cmdpipeline/pipeline_sync_cache.go`
  - `cli/cmdpipeline/pipeline_logs.go`
  - `cli/cmdpipeline/pipeline_history.go`
  - `cli/cmdpipeline/pipeline_error_extract.go`
  - `cli/cmdpipeline/pipeline_compact_test.go`
  - `cli/cmdpipeline/pipeline_history_test.go`
- **Implementation**:
  - Implemented `ResolveDbFileSize` and normalized `FormatRelativeDbPath`.
  - Updated `renderCleanSuccessDbAndHistory` and `renderSavedLocationsTerminal` to display DB size on next line.
  - Updated `renderSyncResultTerminal` and `printEmptyCachedFailures`.
  - Added `isRustOkLine` to `isOkLogLine` and prevented context pollution in `appendContextLine`.
  - Added test cases in `pipeline_compact_test.go` and `pipeline_history_test.go`.

---

## 4. Verification Evidence
1. `TestResolvePipelineDbPath`: Verified repo slug resolves to `.gitmap` while `test-` slugs stay isolated.
2. `TestResolveDbFileSize` & `TestFormatRelativeDbPath_RepoScoped`: Verified human file size formatting and relative path normalization.
3. `TestIsOkLogLine_WithPassingInputs_ReturnsTrue` & `TestIsOkLogLine_WithFailingInputs_ReturnsFalse`: Verified Rust test lines (`test ... ok`, `test result: ok.`) are properly identified as passing, and `test ... FAILED` as failing.
4. Go code formatting verified across all files via `26-go-code-formatter.py`.
