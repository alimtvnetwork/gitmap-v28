# Root Cause Analysis: Gocritic ifElseChain Diff Failure in CI

## 1. Context and Problem Statement
During CI/CD execution of the Lint Baseline Guard (`check-single-linter-diff.sh gitmap` with `LINTER=gocritic`), the baseline diff gate failed with:
`Error: [gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)`.

## 2. Root Cause
Refactored code paths and new command dispatch routines used repeated `if ... else if ... else ...` conditional chains where Go `switch` statements are more idiomatic, triggering gocritic's `ifElseChain` rule in the diff against the baseline.

## 3. Corrective and Preventive Actions
- Refactored `cluster/exec_install.go` to use `switch mgr`.
- Refactored `cmd/cg.go`, `cmd/cg_worker.go`, and `cmd/clustersubcmd.go` to use `switch` constructs.
- Refactored `cmd/filemanipulator.go`, `cmd/schedule_cmd.go`, `cmd/schedule_os.go`, `cmd/ssh_parser.go`, and `cmd/workdir_flags.go` to use `switch` constructs.
- Refactored all other `ifElseChain` occurrences repository-wide in `movemerge/conflict.go`, `render/pretty_parse.go`, `tui/releases.go`, and `osutil/detector.go`.

## 4. Verification
- Ran `golangci-lint run --no-config --disable-all --enable=gocritic ./...` locally and verified 0 `ifElseChain` findings across the entire codebase.
- Ran `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` and verified all 13 checks pass cleanly (exit code 0).
