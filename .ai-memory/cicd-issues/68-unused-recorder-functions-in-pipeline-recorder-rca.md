# RCA: Unused Functions `recordSingleFailedRun` and `isRunFailureAlreadyRecorded` in Pipeline Recorder

- **Date:** 2026-09-21
- **Affected Run:** [CI #35554604260](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35554604260)
- **Target Commit:** `76db95c5c9bf05959ab1ab697f2ce9fbdb4c19c8`
- **Scope:** CI/CD Quality Gates (`golangci-lint` unused analyzer)

---

## 1. Symptom

Pipeline run [#35554604260](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35554604260) failed across 2 jobs:
1. **Lint** (`golangci-lint` strict, fail on any error):
   ```text
   cmdpipeline/pipeline_recorder.go:92:6: func `recordSingleFailedRun` is unused (unused)
       func recordSingleFailedRun(pipeDb *pipelinedb.PipelineSplitDb, repo string, r ghRunItem) {
       ^
   cmdpipeline/pipeline_recorder.go:112:6: func `isRunFailureAlreadyRecorded` is unused (unused)
       func isRunFailureAlreadyRecorded(pipeDb *pipelinedb.PipelineSplitDb, r ghRunItem) bool {
       ^
   Process completed with exit code 1.
   ```
2. **Full Suite Guard** (`golangci-lint` strict, full suite):
   ```text
   cmdpipeline/pipeline_recorder.go:92:6: func `recordSingleFailedRun` is unused (unused)
       func recordSingleFailedRun(pipeDb *pipelinedb.PipelineSplitDb, repo string, r ghRunItem) {
       ^
   cmdpipeline/pipeline_recorder.go:112:6: func `isRunFailureAlreadyRecorded` is unused (unused)
       func isRunFailureAlreadyRecorded(pipeDb *pipelinedb.PipelineSplitDb, r ghRunItem) bool {
       ^
   Process completed with exit code 1.
   ```

---

## 2. Root Cause

During the pipeline error recording refactoring in `cli/cmdpipeline/pipeline_recorder.go`, `recordSingleFailedRun` and `isRunFailureAlreadyRecorded` became orphaned dead code when run failure recording was restructured to use `persistSingleFailedRunLog` and caller batching directly.

---

## 3. Resolution

1. **Dead Code Elimination**: Removed the orphaned private functions `recordSingleFailedRun` and `isRunFailureAlreadyRecorded` from `cli/cmdpipeline/pipeline_recorder.go`.
2. **Local Lint Gate Verification**: Executed `golangci-lint run ./...` inside `cli/`, confirming 0 errors and a clean exit code 0.
3. **Quality Gates Verification**: Executed `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` across all 35 static analysis and compile gates, confirming 100% green pass in 294.95s.
4. **Smart Test Verification**: Executed `python 03-ai-scripts/06-cicd-local-runner.py run-smart`, verifying all 148 Go packages pass.

---

## 4. Prevention & Learnings

- **Pre-Commit Lint Sweep**: Always run `golangci-lint run ./...` in `cli/` whenever refactoring internal helper functions to guarantee that no orphaned functions remain before pushing commits.
- **AST & Dead Code Audits**: Ensure unused symbol linters (`unused` analyzer in `golangci-lint`) run as part of pre-release verification gates.
