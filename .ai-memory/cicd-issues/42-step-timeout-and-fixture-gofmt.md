# CI/CD Issue 42: Step Timeout Flakiness and Fixture Gofmt Backup Dirtiness

- Job: Full Suite Guard, macos-latest / go build + test, windows-latest / go build + test, ubuntu-latest / go build + test
- Type: FAIL
- Detected: 2026-09-12T05:16:11Z
- Status: resolved

## Error
1. macro_test.go:161: expected timeout error under 2.8s, got err: signal: killed, took: 3.003s
2. gofmt_e2e_test.go:75: gofmt -l . reported dirty files after fix-repo: .gitmap/backup/repo-v12/v12/fix-repo/.../files/aligned_map.go

## Root Cause
1. macro.buildStepCmd lacked cmd.WaitDelay, so lingering pipes held by child processes delayed cmd.Run() return on Unix.
2. alignedMapSource included trailing blank lines making the initial fixture dirty, and runGofmtList checked the .gitmap/backup/ directory instead of only working tree files.

## Fix Applied
1. Added cmd.WaitDelay = 250 * time.Millisecond in buildStepCmd and exec sleep %d in macro_test.go.
2. Cleaned alignedMapSource in fixture_helpers_test.go and filtered .gitmap/ and .git/ in runGofmtList.
