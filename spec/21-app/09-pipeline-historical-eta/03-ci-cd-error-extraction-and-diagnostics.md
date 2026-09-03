# CI/CD Error Extraction & Diagnostics

## Overview

When checking pipeline status or wait-time, raw CI/CD logs frequently span thousands of lines containing informational output from passing steps. This document specifies the error extraction engine that isolates actionable failure lines so users and AI agents can immediately diagnose and resolve issues.

## Extraction Criteria

The diagnostic log extractor filters log streams matching the following high-priority fault markers:
1. GitHub Actions workflow errors: `##[error]`
2. Test failures: `❌ FAIL`, `--- FAIL:`, `FAIL:`
3. Compilation & linter errors: `error:`, `Error:`, `syntax error:`, `fatal error:`
4. Process terminations: `Process completed with exit code 1.`, `exit status 1`
5. Panics and stack traces: `panic:`, `goroutine \d+ \[running\]`, `Stack Trace:`

## Output Architecture

When a pipeline status or wait-time check detects that the most recent run completed with `conclusion == "failure"`, it automatically appends the extracted error block to the status card:

```text
  ● Status:           FAILED (conclusion: failure)
  ● Error Diagnostics:
    ──────────────────────────────────────────────────────────────────────
    ❌ FAIL: Found 5 absolute path / URI violation(s):
      gitmap/cmd/agy_test.go:109: Absolute file:/// URI with drive letter
    --- FAIL: TestTopLevelCmdConstantsAreUnique (0.02s)
```

If no errors occurred or the pipeline is actively running without detected segment errors, the error diagnostic block is suppressed.
