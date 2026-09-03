# CI/CD Issue 38: Nested If Violations and Help Examples Heading

- Job: policy-check / go test
- Type: FAIL
- Detected: 2026-08-31T19:40:00Z
- Status: resolved

## Error

```
  ❌ FAILED: Found 3 violation(s):
    - gitmap/apperror/apperror.go:114: Nested 'if' detected (depth 2): 'if !more {'
    - gitmap/cmd/pipeline_ai.go:38: Nested 'if' detected (depth 2): 'if v, err := strconv.Atoi(args[i+1]); err == nil && v >= 0 {'
    - gitmap/cmd/pipeline_status.go:169: Nested 'if' detected (depth 2): 'if dur >= 10 {'

--- FAIL: TestEveryHelpFileHasExamples (0.02s)
    examples_golden_test.go:37: help files missing `## Examples` section (1):
          - llm.md
```

## Root Cause

- Nested conditionals (depth >= 2) in apperror.go, pipeline_ai.go, and pipeline_status.go during new feature additions.
- llm.md used `## 2. Step-by-Step Workflows with Real-World Examples` which did not start with `## Examples` as required by golden tests.

## Fix Applied

- Flattened all nested if blocks using guard clauses and extracted helper functions.
- Updated helptext/llm.md, llm.md, gitmap/llm.md, and gitmap/cmd/llm/llm.go to use `## Examples (Step-by-Step AI Workflows)`.
- Verified all linters and golden tests pass 100%.
