# 221 - Component & CLI Spec: Pipeline Error Handling & Test Parsing

## 1. Component Boundaries

### 1.1 `cli/cmdpipeline/pipeline_persist.go`
- **Component:** SQLite Pipeline Error Log Persistence
- **Package:** `cmdpipeline`
- **Target Function:** `persistLogToRepoSplitDb`
- **Helper Added:** `isLogAlreadyPersisted(pipeDb *pipelinedb.PipelineSplitDb, runId uint64) bool`
- **Requirements:**
  - Function length <= 8 lines.
  - Zero nested `if` statements.
  - Guard clauses with early returns.

### 1.2 `cli/cmdpipeline/pipeline_error_extract.go`
- **Component:** Pipeline Error & Traceback Extractor
- **Package:** `cmdpipeline`
- **Target Functions:** `isStrongerSummary`, `isTestFailureSummary`, `isTestStatusBanner`
- **Requirements:**
  - `isTestStatusBanner(s string) bool`: detects generic runners like `--- FAIL: ` and `FAIL\t` and `FAIL: `.
  - `isStrongerSummary(candidate, current string) bool`:
    - Candidate with location information (`isLocationSummary`) always overrides a banner.
    - Candidate containing `Expected ` or `AssertionError` always overrides a banner.
    - Preserves Python traceback and unittest / pytest error prioritization.

---

## 2. Test Plan
- Run `TestParseFailedLogLines`: verifies `sample_test.go:12: Expected true to be false` is chosen over `--- FAIL: TestSample (0.00s)`.
- Run `TestParseFailedLogLines_PythonUnitTestE2E`: verifies Python tracebacks and unittest failures are captured.
- Run `TestParseFailedLogLinesWithWarnings`: verifies compiler warnings and failure context.
