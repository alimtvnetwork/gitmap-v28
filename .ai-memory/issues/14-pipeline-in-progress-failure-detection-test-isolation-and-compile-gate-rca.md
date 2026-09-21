# RCA-14: In-Progress Pipeline Failure Immediate Detection, Test Isolation, and Go Compile Gate Fix

## 1. Root Cause Analysis

### Symptoms
1. **GitHub Actions CI (Run 35576139471)**: Failed on `Relative Path Check` due to 14 absolute paths and `file:///` URIs hardcoded in `.ai-memory/issues/13-pipeline-cache-ttl-and-table-order-inversion-rca.md`.
2. **GitHub Actions Cross-Platform Build (Run 35576139366)**: Failed on `ubuntu-latest / go build + test` in `TestPersistErrorReport` (`cli/cmdpipeline/pipeline_errorlogs_test.go:264`).
3. **Internal Diagnostic Gate (`runInternalCICDChecks`)**: Failed with:
   `✖ FAIL Go Compile Gate: pattern ./...: directory prefix . does not contain main module or its selected dependencies`
4. **User Experience (`gitmap pe` / `gitmap p`)**:
   When a multi-job workflow run was actively `in_progress` on GitHub and had already suffered a step or job failure (such as `Relative Path Check` failing 13s into `CI`), running `gitmap pe` reported only `● Active Pipeline is RUNNING` and exited early without showing the failed section or error logs. Furthermore, the history summary table was suppressed during running states.

### Root Causes
1. **Hardcoded Drive Letters in RCA-13**:
   RCA-13 document contained literal Windows drive paths (`d:/work/gitmap/...`) and `file:///` URIs which were flagged by `linter-scripts/check-relative-paths.py`.
2. **Missing In-Progress Job Inspection in `collectFailedRuns`**:
   In GitHub Actions, multi-job workflows remain `status: "in_progress"` and `conclusion: ""` until all jobs complete. `collectFailedRuns` in `cli/cmdpipeline/pipeline_logs.go` only checked `isFailingConclusion(r.Conclusion)`. Because `r.Conclusion` is empty while running, failed jobs inside in-progress workflows were ignored, leading `len(p.FailedRuns) == 0`. In `renderErrorLogsTerminal`, `if p.IsRunning && len(p.FailedRuns) == 0` returned immediately, hiding the failure.
3. **Test State Pollution in `TestPersistErrorReport`**:
   `TestPersistErrorReport` called `writeCombinedErrorReport` without repo qualification, writing to the default repo directory `.gitmap/data/pipeline/alimtvnetwork-gitmap-v28/pipeline_errors.log`. During parallel test execution across packages, concurrent test runs (e.g. `clearLocalErrorLogs` or `pipelinedb` cache purges) deleted or replaced this file before `os.ReadFile` executed.
4. **Root Module Directory Misidentification in `resolveGitmapDirFromRoot`**:
   `resolveGitmapDirFromRoot` in `cli/cmdpipeline/pipeline_cicd_checker.go` checked for a subdirectory `gitmap/` instead of `cli/` where `go.mod` lives. As a result, `probeCompileGate` ran `go build ./...` in the repository root directory rather than `cli/`, triggering module not found errors.
5. **Database WAL Timestamp Invalidation**:
   In SQLite WAL mode, database writes modify `sql.db-wal` while `sql.db` modification time may lag. Without explicit timestamp touching upon recording fresh runs, TTL evaluation can return stale cache evaluations.

---

## 2. Key Changes

1. **Clean Relative Paths**:
   Replaced all absolute drive letters and `file:///` URIs in `.ai-memory/issues/13-pipeline-cache-ttl-and-table-order-inversion-rca.md` with relative repository paths.
2. **Immediate In-Progress Failure Detection**:
   - In `cli/cmdpipeline/pipeline_logs.go`, updated failed run collection to inspect active runs via `hasActiveRunFailedJobs(repo, r)` and query jobs with `queryRunJobs(repo, r.DatabaseId)`. When any job or step fails in an in-progress workflow, the run is flagged as failing immediately.
   - In `renderErrorLogsTerminal`, removed the early return that suppressed error logs when running, ensuring failure breakdown, failed sections, logs, and `RenderHistorySummaryTable` are always displayed immediately.
3. **Prevent Premature Job Caching**:
   In `cli/cmdpipeline/pipeline_segments.go`, ensured `writeCachedPipelineJobs` only writes to persistent disk cache when `isAllJobsCompleted(jobs)` is true, allowing subsequent calls to receive fresh step progress while in progress.
4. **Test Isolation in `TestPersistErrorReport`**:
   Updated `TestPersistErrorReport` in `cli/cmdpipeline/pipeline_errorlogs_test.go` to use `writeCombinedErrorReportForRepo` with an isolated test repository name and `t.Cleanup`.
5. **Module Directory Correction**:
   In `cli/cmdpipeline/pipeline_cicd_checker.go`, updated `resolveGitmapDirFromRoot` to detect `cli/go.mod` and target `cli/` for Go toolchain commands.
6. **Explicit DB Timestamp Updates**:
   In `cli/cmdpipeline/pipeline_recorder.go` and `cli/cmdpipeline/pipeline_cache_eval.go`, touch database modification time on write and inspect WAL timestamp on cache checks.

---

## 3. Verification Plan

1. **Relative Path Linter**:
   Run `python linter-scripts/check-relative-paths.py` to confirm zero violations.
2. **Code Style & Guidelines**:
   Run `python linter-scripts/check-nested-ifs.py`, `python linter-scripts/check-boolean-guidelines.py`, and `python linter-scripts/check-newline-styling.py`.
3. **Compilation & Unit Tests**:
   Run `go test -v ./cmdpipeline -run "TestPersistErrorReport|TestRunInternalCICDChecks"` and `go test ./... -count=1`.
4. **Binary Build & Release**:
   Build local binary and orchestrate minor release `v6.289.0`.

---

## 4. Prevention & Lessons Learned

- Always check in-progress runs for individual failing jobs rather than relying exclusively on top-level workflow conclusion.
- Unit tests must never write to production default directories; always isolate using unique repo slugs or `t.TempDir()`.
- Diagnostic suites probing the Go compiler must inspect `cli/go.mod` when the repository uses a subfolder module layout.
