# Plan: Pipeline PE Unit Test Traceback Extraction & Heatmap Modernization

> **Plan ID:** 220-pipeline-pe-unit-test-traceback-and-heatmap  
> **Status:** Completed  
> **Created:** 2026-10-05  
> **Associated Spec:** `02-spec/21-app/220-pipeline-pe-unit-test-traceback-and-heatmap/`  

---

## User Request (Verbatim)

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

## Subtasks Sequence

| Subtask ID | Title | Status | Target Area |
| :--- | :--- | :--- | :--- |
| `01` | gofmt Formatting & CI Test Fix | DONE | `cli/`, `.github/scripts/tests/test_ci_scripts.py` |
| `02` | Pipeline PE Unit Test Traceback Extraction | DONE | `cli/cmdpipeline/pipeline_stacktrace.go`, `pipeline_error_extract.go` |
| `03` | Pipeline Heatmap & Log Cache Enhancement | DONE | `cli/cmdpipeline/pipeline_query.go`, `pipeline_details.go` |
| `04` | Parser Unit Tests, Verification & Release Ceremony | DONE | `cli/cmdpipeline/pipeline_error_extract_test.go`, `version.json`, release |
