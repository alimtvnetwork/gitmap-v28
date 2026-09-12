# Subtask 01: Inverted Success Elimination & Method Hardening

**Parent Plan:** `.lovable/plans/pending/96-booleans-and-complex-conditions-audit.md`  
**Status:** In Progress  
**Target Files:**
- `gitmap/model/record.go`
- `gitmap/model/prompt_status.go`
- `gitmap/pipelinedb/pipeline_repository_test.go`
- `gitmap/lazyregex/compile_result.go`
- `gitmap/lazyregex/lazyregex_test.go`
- `gitmap/result/result_test.go`
- `gitmap/cmd/install_audit_runner_test.go`
- `gitmap/cmd/prompt_failure_reporter.go`
- `gitmap/cmd/schedule_cmd.go`
- `gitmap/cmd/pull.go`
- `gitmap/cmd/pullparallel.go`
- `gitmap/cmd/push.go`
- `gitmap/cmd/pushparallel.go`
- `gitmap/cmd/install_logs.go`
- `gitmap/tests/cmd_test/prompt_install_e2e_test.go`

---

## 1. Objectives & Grounded Rules

1. **Add Affirmative Failure Methods:**
   - In `gitmap/model/record.go`: Add `func (r CloneResult) IsFailed() bool { return !r.IsSuccess }` and `func (r CloneResult) IsFail() bool { return !r.IsSuccess }`.
   - In `gitmap/model/prompt_status.go`: Add `func (r PromptRunStatusRecord) IsFailed() bool { return !r.IsSuccess }` and `func (r PromptRunStatusRecord) IsFail() bool { return !r.IsSuccess }`.
   - In `gitmap/lazyregex/compile_result.go`: Refactor `IsFailed()` to `return it.appError != nil || it.re == nil` rather than calling `!it.IsSuccess()`.

2. **Replace All 17 Inverted Success Checks:**
   - `gitmap/cmd/install_audit_runner_test.go:13`: `res.IsFailed()`
   - `gitmap/cmd/install_logs.go:75`: `r.IsFailed() && len(filtered) < limit`
   - `gitmap/cmd/prompt_failure_reporter.go:13`: `r.IsFailed()`
   - `gitmap/cmd/pull.go:517`: `result.IsFailed()`
   - `gitmap/cmd/pullparallel.go:108`: `result.IsFailed()`
   - `gitmap/cmd/pullparallel.go:112`: `isStopRequested := result.IsFailed() && stopOnFail` (also eliminates `shouldStop` banned prefix and inverted success)
   - `gitmap/cmd/push.go:223`: `result.IsFailed()`
   - `gitmap/cmd/pushparallel.go:99`: `result.IsFailed()`
   - `gitmap/cmd/pushparallel.go:102`: `result.IsFailed() && stopOnFail`
   - `gitmap/cmd/pushparallel.go:105`: `result.IsFailed()`
   - `gitmap/cmd/schedule_cmd.go:408`: `r.IsFailed()`
   - `gitmap/lazyregex/lazyregex_test.go:212`: `validRes.IsFailed()`
   - `gitmap/lazyregex/lazyregex_test.go:359`: `res.IsFailed() || res.Value != re1`
   - `gitmap/pipelinedb/pipeline_repository_test.go:90`: `runRecord.IsFailed()`
   - `gitmap/result/result_test.go:13`: `res.IsFailed()`
   - `gitmap/tests/cmd_test/prompt_install_e2e_test.go:15`: `res.IsFailed()`

3. **Coding Standards:**
   - UNIX LF line endings, UTF-8 without BOM.
   - Blank lines before `return` and after `}` blocks.
   - Functions <= 8 lines preferred (hard cap 15 lines).
   - Verify with `go test -v -short -run "Test.*(Pull|Push|Clone|Install|Regex)" ./...`.
