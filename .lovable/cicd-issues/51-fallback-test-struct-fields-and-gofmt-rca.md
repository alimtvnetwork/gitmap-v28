# CI/CD Issue 51: Fallback Test Struct Fields and gofmt Drift RCA

## 1. Symptom
In GitHub Actions CI run `#35263234345` on commit `75215ea`, two jobs failed:
1. **Lint** (`Go vet`):
```text
vet: cmdpipeline/pipeline_fallback_test.go:41:4: unknown field CommitSha in struct literal of type CommitPipelineGroup
```
2. **Lint Script Unit Tests** (`Run lint-script unit tests`):
```text
FAIL: test_gofmt_check_clean_repo (__main__.TestGoFormatCheck.test_gofmt_check_clean_repo)
AssertionError: 1 != 0
```

---

## 2. Root Cause
1. **`CommitPipelineGroup` Field Name Drift**: In `cli/cmdpipeline/pipeline_fallback_test.go`, the test mock instantiated `CommitPipelineGroup` with `{CommitSha: "...", ShortSha: "..."}`. The actual definition in `pipeline_commit_groups.go` defines `HeadSha string` with no `CommitSha` or `ShortSha` fields.
2. **Trailing Line / Formatting Inconsistency in New Files**: `cli/cmdpipeline/pipeline_fallback.go` contained a trailing empty line after the closing brace, triggering `test_gofmt_check_clean_repo` to fail during strict repository format checks.

---

## 3. Resolution
1. Corrected `cli/cmdpipeline/pipeline_fallback_test.go` to use `{HeadSha: "..."}` and asserted against `fallback.HeadSha`.
2. Cleaned and formatted `cli/cmdpipeline/pipeline_fallback.go` and `cli/cmdpipeline/pipeline_fallback_test.go` with `gofmt -w`.
3. Verified zero unformatted files across all 2,853 Go files using `python .github/scripts/go-format-check.py --check-only`.

---

## 4. Prevention & Learnings
- **Inspect Struct Definitions**: Always verify exact struct field names from domain declarations before creating test fixtures.
- **Pre-commit `gofmt` Validation**: Always run repository gofmt checker before committing Go source files.
