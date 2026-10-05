# Master Audit Ledger: Task 220 - Pipeline PE Unit Test Traceback Extraction, Heatmap Modernization & CI Test Remediation

**Version:** 2.0.0  
**Status:** In Progress (Phase 2 Parallel Execution)  
**Parent Task ID:** task-20261005033654-pipeline-pe-unit-test-traceback-and-heatmap  
**Database:** `.ai-memory/temp-agents/pipeline-pe-unit-test-traceback-and-heatmap/agent-task.db`  

---

## 1. User Request (Verbatim)

```text
When we do the heat map PE, PE is not showing the error message. That's one problem. I'm giving you the terminal version that is shown and what was in the error. Both of these are very important because when error happens with a unit test, you need to put the unit test information correctly with all the lines that is there. Probably not the run 18 tests, this type of line. But yes, the traceback call, which test is failed, this type of information is very important for the logger and for the AI to fix the error right away. Make sure that you respect that and you update the heat map regarding this heat map PE, and also fix this test, make a minor bump in the version, and make a release.

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Update the heat map regarding the heat map PE issue.
4. Fix the unit test to ensure error messages are correctly displayed.
5. Include the traceback call and failed test information for logging and AI error fixing.
6. Make a minor version bump.
7. Make a release.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]
```

---

## 2. Discrete Technical Deliverables (Traceable IDs)

1. **Task-01: Job-Level Log Fetcher & DB Cache Self-Healing (`Worker 01`)**
   - Query failing jobs individually via `gh run view --job <databaseId> --repo <repo> --log-failed` when run-level log fails or is in progress.
   - Normalize job log output lines to `<job_name>\t<step_name>\t<log_line>` format.
   - Prevent fallback stubs from locking SQLite Split-DB records permanently in `persistLogToRepoSplitDb`.
   - Target files: `cli/cmdpipeline/pipeline_query.go`, `cli/cmdpipeline/pipeline_persist.go`.

2. **Task-02: Traceback Parser & Heatmap Formatter (`Worker 02`)**
   - Scan backwards before `Traceback (most recent call last):` to capture preceding `FAIL: <test>` or `ERROR: <test>` line.
   - Prevent premature stack termination on empty lines in `isStackStopLine` to preserve multi-line assertion diffs.
   - Expand `isTestFailureSummary` to match `FAIL:`, `ERROR:`, `FAILED`, `--- FAIL:`.
   - Render `j.StackTrace` and `j.ErrorLines` in `printFailedJobItemToBuilder` for commit inspector / heatmap cards.
   - Author regression unit tests in `cli/cmdpipeline/pipeline_error_extract_test.go`.
   - Target files: `cli/cmdpipeline/pipeline_stacktrace.go`, `cli/cmdpipeline/pipeline_error_extract.go`, `cli/cmdpipeline/pipeline_history.go`, `cli/cmdpipeline/pipeline_error_extract_test.go`.

3. **Task-03: Targeted Verification & Quality Audit (`Lead`)**
   - Run file-scoped targeted checks: guideline autofixer, relative paths, doc path linter, forbidden strings.
   - Verify all 3,951 Go files remain gofmt-clean (`python .github/scripts/go-format-check.py --check-only`).

4. **Task-04: Minor Version Bump & Release Ceremony (`Lead`)**
   - Minor version bump from `v6.477.0` to `v6.478.0` (`version.json`, `package.json`, `changelog.md`).
   - Create release notes `.ai-memory/release/release-notes-v6.478.0.md`.
   - Execute atomic GitMap commit and release.

---

## 3. Subagent Spawning Ledger (A = 2, H = 2)

| Agent ID | Type | Role | Hands / Scope | Status |
| :--- | :--- | :--- | :--- | :--- |
| `Worker 01` | `self` | Job-Level Log Fetcher & DB Cache Engine | `pipeline_query.go`, `pipeline_persist.go` | In Progress |
| `Worker 02` | `self` | Traceback Parser & Heatmap Formatter | `pipeline_stacktrace.go`, `pipeline_error_extract.go`, `pipeline_history.go`, `pipeline_error_extract_test.go` | In Progress |
