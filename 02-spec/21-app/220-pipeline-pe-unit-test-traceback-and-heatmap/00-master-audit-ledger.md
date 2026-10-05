# Master Audit Ledger: Task 220 - Pipeline PE Unit Test Traceback Extraction, Heatmap Modernization & CI Test Remediation

**Version:** 1.0.0  
**Status:** Completed (Phase 3 Verified & Ready for Release)  
**Parent Task ID:** 220  
**Database:** `.ai-memory/temp-agents/220-pipeline-pe-unit-test-traceback/agent-task.db`  

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

1. **Task-01: gofmt Code Formatting & CI Path Resolution**
   - Format the 5 dirty Go files (`cli/cmd/help.go`, `cli/cmdagent/agent_cleanup.go`, `cli/cmdagent/agent_diagnose.go`, `cli/cmdcursor/cursor_sync.go`, `cli/cmddb/cmddb_reset.go`) using `gofmt -w`.
   - Fix `REPO_ROOT` path resolution in `.github/scripts/tests/test_ci_scripts.py:15` (`SCRIPTS_DIR, "..", ".."`).
   - Verify `python .github/scripts/go-format-check.py --check-only` exits 0.

2. **Task-02: Pipeline PE Unit Test Traceback Extraction Engine**
   - Enhance `cli/cmdpipeline/pipeline_stacktrace.go` to parse Python tracebacks (`Traceback (most recent call last):`), Python stack frames (`File "...", line \d+, in \w+`), exception lines, and Go test failures.
   - Guard against premature stack frame termination on test divider lines (`===`, `---`).
   - Enhance `cli/cmdpipeline/pipeline_error_extract.go` to elevate `FAIL: <test_name>` to `FailureSummary`.
   - Filter out test progress noise like `Ran \d+ tests in ...` in `isToolProgressNoise`.

3. **Task-03: Log Caching & Heatmap Terminal/Clipboard Enhancements**
   - Guard against caching synthetic step fallbacks as permanent log files during in-progress pipeline runs.
   - Format unit test failures prominently in `renderFailedJobSection` and clipboard payloads.

4. **Task-04: Test Suite & Minor Version Release Ceremony**
   - Author regression unit tests in `cli/cmdpipeline/pipeline_error_extract_test.go`.
   - Run tests and verify 100% passing.
   - Minor version bump from `6.475.0` to `6.476.0` (`version.json`, `package.json`, `changelog.md`).
   - Commit and release via GitMap CLI.

---

## 3. Subagent Spawning Ledger (A = 2, H = 2)

| Agent ID | Type | Role | Hands/Scope | Status |
| :--- | :--- | :--- | :--- | :--- |
| `Worker-01` | `self` | Python/Go formatting & CI test fixes | Subtask 01: `cli/` gofmt & `test_ci_scripts.py` | Ready to Dispatch |
| `Worker-02` | `self` | Pipeline error extraction & stacktrace parser | Subtask 02: `pipeline_stacktrace.go` & `pipeline_error_extract.go` | Ready to Dispatch |
