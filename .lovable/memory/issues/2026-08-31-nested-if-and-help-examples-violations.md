# Root Cause Analysis: Nested If Statements and Help Examples Heading

## 1. Why it happened

During the implementation of the pipeline-ai command suite and stack trace formatting in apperror.go, new conditional logic was added that introduced nested if statements (depth >= 2) in apperror.go, pipeline_ai.go, and pipeline_status.go. Additionally, the updated llm.md help file used ## 2. Step-by-Step Workflows with Real-World Examples instead of matching the golden test requirement for ## Examples.

## 2. How it happened

- In gitmap/apperror/apperror.go, captureStackTrace inspected runtime frames and contained if strings.Contains(...) { if !more { break } continue }.
- In gitmap/cmd/pipeline_ai.go, parsePipelineAIDelay contained if isPipelineAIDelayFlag(...) { if v, err := strconv.Atoi(...); err == nil ... }.
- In gitmap/cmd/pipeline_status.go, sumCompletedRunDurations contained if r.Status == "completed" ... { if dur >= 10 { ... } }.
- In gitmap/helptext/examples_golden_test.go, the golden test requires every help file to have a section starting with ## Examples followed by code blocks.

## 3. Root Cause

- Nested conditionals violated the coding guideline prohibiting nested if blocks (depth >= 2) enforced by check-enum-and-boolean.py.
- The section title in helptext/llm.md did not start with ## Examples, failing TestEveryHelpFileHasExamples.

## 4. Code Fix

- Refactored captureStackTrace to delegate frame appending to appendStackFrame, eliminating the nested if.
- Flattened parsePipelineAIDelay by extracting extractDelaySecondz helper and using guard checks.
- Flattened sumCompletedRunDurations using early continue guard clauses.
- Renamed the section heading in helptext/llm.md, llm.md, gitmap/llm.md, and gitmap/cmd/llm/llm.go to ## Examples (Step-by-Step AI Workflows).
- Verified all linters (check-enum-and-boolean.py, go-format-check.py) and all tests (go test ./...) pass 100%.
