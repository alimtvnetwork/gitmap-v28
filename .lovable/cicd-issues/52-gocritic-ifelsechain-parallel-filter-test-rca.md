# CI/CD Issue 52: gocritic ifElseChain Violation in Parallel Filter Test RCA

## 1. Symptom
In GitHub Actions CI run `#35263724019` on commit `e9cb7384`, the `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)` job failed:

```text
[CI/CD ERROR REPORT] New gocritic Violations Detected (1 finding(s)):
##[error][gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)
  ┌─ Error [1/1]: [gocritic]
  │ File:     cmdpipeline/pipeline_parallel_filter_test.go:11:3
  │ Location: cmdpipeline/pipeline_parallel_filter_test.go at line 11, col 3
  │ Linter:   gocritic
  │ Message:  ifElseChain: rewrite if-else to switch statement
  │ Context:  NEW finding vs baseline
  │ Script:   .github/scripts/check-single-linter-diff.py
  └──────────────────────────────────────────────────────────
FAIL: 1 new gocritic finding(s) detected!
```

---

## 2. Root Cause
In `cli/cmdpipeline/pipeline_parallel_filter_test.go:11-17`, test line generation used an `if i%3 == 0 { ... } else if i%3 == 1 { ... } else { ... }` chain. The repository's strict `gocritic` linter mandates replacing multi-branch conditionals with a switch statement (`switch i % 3`).

---

## 3. Resolution
1. Refactored lines 11–17 of `cli/cmdpipeline/pipeline_parallel_filter_test.go` from the `if-else` chain to:
```go
switch i % 3 {
case 0:
    lines[i] = fmt.Sprintf("PASS: step #%d passed successfully", i)
case 1:
    lines[i] = fmt.Sprintf("✔ ok: task %d finished", i)
default:
    lines[i] = fmt.Sprintf("Error at line %d: syntax error", i)
}
```
2. Ran `gofmt -w cli/cmdpipeline/pipeline_parallel_filter_test.go` and verified repository gofmt cleanliness.
3. Verified zero guideline violations via boolean, nested-if, and error-management linters.

---

## 4. Prevention & Learnings
- **Conditionals with Modulo/Enum Targets**: Always prefer `switch` constructs over sequential `if-else if-else` chains in Go to comply with `gocritic` conventions.
