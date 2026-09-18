# Subtask 01: Banned Function Prefixes & Negative Variables Refactoring

Parent Plan: [142-boolean-principles-negatives-and-complex-conditions.md](../../pending/142-boolean-principles-negatives-and-complex-conditions.md)

## Goals
1. Refactor `cli/cmdscan/flags.go`: Rename `wasFlagPassed` to `isFlagPassed` and update call sites.
2. Refactor `cli/cluster/exec_install.go`: Rename `noCmdErr := cmdErr == nil` to `isCmdSuccess := cmdErr == nil`.
3. Refactor `cli/cmd/create_ops.go`: Rename `NoRemote bool` and `noRemote := ...` to `IsSkipRemote bool` and `isSkipRemote := ...`.
4. Refactor `cli/cmd/mv_flags.go` and `cli/cmd/mv_relocate.go`: Rename `noVSCode` and `noDesktop` to `isSkipVSCode` and `isSkipDesktop`.
5. Refactor `cli/cmdvscode/vscodepmscan.go`, `cli/cmdvscode/exports.go`, `cli/cmdscan/helpers.go`, `cli/cmdscan/exports.go`, and `cli/cmd/clihelpers.go`: Rename `noVSCodeSync` and `noAutoTags` parameters to `isSkipVSCodeSync` and `isSkipAutoTags`.

## Acceptance Criteria
- [x] Zero `wasFlagPassed` occurrences remaining.
- [x] Zero `noCmdErr`, `noRemote`, `noVSCode`, `noDesktop`, `noVSCodeSync`, and `noAutoTags` variable/field names remaining in target files.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
