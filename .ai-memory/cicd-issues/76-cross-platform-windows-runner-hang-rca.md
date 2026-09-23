# RCA 76: Cross-Platform Build Windows Runner Hang & Timeout Hardening

## 1. Symptom

Pipeline run `#35798625944` on workflow `Cross-Platform Build` failed on `windows-latest / go build + test` (Job ID: `106983624017`).
Duration: 46 minutes 6 seconds.
Status: Completed with conclusion `failure`.

Raw GitHub Actions error annotation:
```text
The hosted runner lost communication with the server. Anything in your workflow that terminates the runner process, starves it for CPU/Memory, or blocks its network access can cause this error.
windows-latest / go build + test: .github#1
```

Step Breakdown:
- `Set up job`: 2s (success)
- `Run actions/checkout@v6`: 7s (success)
- `Optimize Windows Runner I/O`: 4s (success)
- `Run actions/setup-go@v6`: Started at `2026-09-22T23:41:52Z`, never reached completion.
- Runner termination: `2026-09-23T00:27:39Z` (terminated by GitHub Actions runner communication reaper after 45m47s).
- Subsequent steps (`go build ./...`, `go test ./...`) remained pending.

## 2. Root Cause

1. **Missing Job Timeout Directive**:
   `.github/workflows/cross-platform.yml` lacked a `timeout-minutes` configuration on the `build-and-test` matrix job. While normal runs execute in 4–6 minutes, when a transient network disconnection, node communication drop, or socket stall occurred on the GitHub-hosted Windows runner, the runner hung silently until hitting the 45-minute communication timeout.
2. **Incomplete Windows Defender Exclusion Scope**:
   The `Optimize Windows Runner I/O` step previously only excluded `${{ github.workspace }}`. Go toolchain extraction by `actions/setup-go` operates in `C:\hostedtoolcache\windows\go`, Go module cache in `$HOME\go\pkg\mod`, and compiler build cache in `$LOCALAPPDATA\go-build`. Leaving these unexcluded allowed Windows Defender Real-Time Monitoring and Antivirus file inspection to scan, lock, and intercept Go binary extraction and socket streams during Node 24 actions execution.
3. **Codebase Lint & Guideline Regressions**:
   Targeted local audit via `03-ai-scripts/06-cicd-local-runner.py` identified:
   - Nested `if` violations in `cli/cmdpull/pull.go` (line 801), `cli/cmdpipeline/pipeline_error_extract.go` (line 1077), and `cli/cmdpipeline/pipeline_logs.go` (lines 438, 522).
   - Unused dead code in `cli/cmdpipeline/pipeline_logs.go` (`populateRunsIntoPayload`, `resolveFailedRunsForPayload`, `renderSectionStackTrace`, `renderSectionErrorLines`).
   - `linter-scripts/check-nested-ifs.py` missing `--changed-only` and `--commits` CLI argument support.

## 3. Resolution

1. **Workflow Hardening (`.github/workflows/cross-platform.yml`)**:
   - Added `timeout-minutes: 20` to the `build-and-test` job to bound runner execution and terminate hung runners within 20 minutes instead of burning 45+ minutes.
   - Expanded Windows Defender exclusions to include toolcache and Go caches:
     ```powershell
     Set-MpPreference -ExclusionPath "${{ github.workspace }}", "C:\hostedtoolcache", "$env:USERPROFILE\go", "$env:LOCALAPPDATA\go-build"
     ```
2. **Source Code Refactoring & Nested If Flattening**:
   - In `cli/cmdpull/pull.go`, extracted `resolveExplicitPullTargets` helper to eliminate nested `if`.
   - In `cli/cmdpipeline/pipeline_error_extract.go`, inverted `strings.Index(line, "\"/OUT:") == -1` guard to eliminate nested `if`.
   - In `cli/cmdpipeline/pipeline_logs.go`, extracted `resolveExplicitTargetSha` and flattened hex SHA condition in `queryWorkflowRunsForTarget`.
   - Removed unused functions `populateRunsIntoPayload`, `resolveFailedRunsForPayload`, `renderSectionStackTrace`, and `renderSectionErrorLines`.
3. **Linter Support (`linter-scripts/check-nested-ifs.py`)**:
   - Added `--changed-only` and `--commits` flags to `parse_cli_args()`.
4. **Verification**:
   - `golangci-lint run --allow-serial-runners --issues-exit-code=1 --timeout=10m -c .golangci.yml --path-prefix cli ./...` passed with exit code 0.
   - `linter-scripts/check-enum-and-boolean.py --changed-only --commits 20` passed with exit code 0.
   - `linter-scripts/check-nested-ifs.py --changed-only --commits 20` passed with exit code 0.
   - GitHub Actions rerun of job `106999283129` succeeded past `actions/setup-go@v6` and `go build ./...`.

## 4. Prevention & Learnings

- **Mandatory Timeout for All Matrix Jobs**: Every GitHub Actions workflow job that runs across matrices (`windows-latest`, `macos-latest`, `ubuntu-latest`) must specify `timeout-minutes: 15` or `20` to prevent runner hangs from exhausting runner quotas.
- **Comprehensive Defender Exclusions on Windows Runners**: Whenever running Go builds or setups on `windows-latest`, always exclude toolcache and Go paths (`C:\hostedtoolcache`, `$env:USERPROFILE\go`, `$env:LOCALAPPDATA\go-build`) in addition to `${{ github.workspace }}`.
- **Dead Code Cleanup Before Commits**: Run `golangci-lint` locally to prune replaced or deprecated functions before committing features.
