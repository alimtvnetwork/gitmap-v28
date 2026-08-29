# CI/CD Issue 34: Gocritic ifElseChain New Findings in CI Diff Gate

- **Job**: Lint Baseline Guard (`check-single-linter-diff.sh gitmap` with `LINTER=gocritic`)
- **Type**: FAIL
- **Detected**: 2026-08-29
- **Status**: resolved

## Error
```text
Error: [gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)
Error: [gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)
FAIL: 2 new gocritic finding(s). Fix the issues above.
Error: Process completed with exit code 1.
```

## Root Cause
Several newly touched and refactored files used multi-branch `if ... else if ... else ...` constructs instead of Go `switch` statements, triggering gocritic's `ifElseChain` rule in the diff against the baseline.

## Fix Applied
1. Refactored `cluster/exec_install.go` to use `switch mgr` for package manager selection.
2. Refactored `cmd/cg.go` to use `switch` in `resolveCGRepos`.
3. Refactored `cmd/cg_worker.go` to use `switch` in `printCGUpdateSummary`.
4. Refactored `cmd/clustersubcmd.go` to use `switch cmdToken` for command dispatch.
5. Refactored `cmd/filemanipulator.go` to use `switch` in `assignPositional` and `handleFixSeqArg`.
6. Refactored `cmd/schedule_cmd.go` and `cmd/schedule_os.go` to use `switch` for subcommand and OS dispatch.
7. Refactored `cmd/ssh_parser.go` and `cmd/workdir_flags.go` to use `switch` for argument parsing.
8. Refactored `movemerge/conflict.go`, `render/pretty_parse.go`, `tui/releases.go`, and `osutil/detector.go` to use `switch`.
9. Verified with `golangci-lint run --no-config --disable-all --enable=gocritic ./...` (0 `ifElseChain` occurrences remain in entire codebase).
