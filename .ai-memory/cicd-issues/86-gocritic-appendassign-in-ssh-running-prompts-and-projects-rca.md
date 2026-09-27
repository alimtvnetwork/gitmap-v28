# 86. CI/CD RCA: Gocritic appendAssign in SSH Running Prompts and Projects Aggregation

## 1. What Failed & Where
- **Failing CI/CD Step:** `Lint Baseline Guard: Gocritic diff (baseline-diff, full-path only)` on workflow run `#36284510621`.
- **Target Files:**
  - `cli/cmdagy/agy_running_projects_ssh.go:21:16`
  - `cli/cmdagy/agy_running_prompts_ssh.go:19:9`
  - `cli/cmdagy/agy_running_prompts_ssh.go:240:9`
- **Error Description:**
  - `gocritic` flagged `appendAssign: append result not assigned to the same slice` because local records were passed into aggregation functions and sliced with `all := append(local, remote...)` without allocating a dedicated slice or assigning back to the receiver.

## 2. Root Cause Analysis
- When concatenating slices in Go, calling `append(a, b...)` without reassigning to `a` (`a = append(a, b...)`) can cause unintended side-effects if slice `a` has capacity to spare, mutating the underlying backing array of the caller's slice.
- `gocritic` enforces that `append` results must either overwrite the same slice variable or be constructed via a freshly allocated slice with capacity (`make([]T, 0, len(a)+len(b))`).

## 3. Grounded Remediation
- Extracted helper functions `mergeProjectRecords`, `mergePromptRecords`, and `mergeBatchRecords`:
  ```go
  func mergeProjectRecords(local, remote []RunningProjectRecord) []RunningProjectRecord {
  	merged := make([]RunningProjectRecord, 0, len(local)+len(remote))
  	merged = append(merged, local...)
  	merged = append(merged, remote...)
  	return merged
  }
  ```
- Each helper allocates a new slice with precise capacity `len(local) + len(remote)` and sequentially appends without mutating caller slices.
- Satisfies function sizing guidelines (<= 8 lines per function) and gocritic strict linters.

## 4. Verification & Prevention
- Verified locally with `golangci-lint run --issues-exit-code=1 ./cmdagy/...` (0 issues).
- Verified with `python linter-scripts/check-nested-ifs.py` (0 issues).
- Verified with `python linter-scripts/check-enum-and-boolean.py` (PASS).
- Verified with `python linter-scripts/check-error-management.py` (PASS).
- Rebuilt global binary at `C:\Users\Administrator\AppData\Local\gitmap-cli\gitmap.exe`.
- Tested `gitmap running-prompts ls --ssh` locally (aggregated successfully across nodes).
- Pushed commit `3afac642b7b3fc7114ff4feadea0bff9295481e3` to `main`.
- Remote GitHub Actions CI run `#36285040528` passed 100% green across all checks (CI, Cross-Platform Build, race-detector, History Rewrite Smoke, CI Beacon).
