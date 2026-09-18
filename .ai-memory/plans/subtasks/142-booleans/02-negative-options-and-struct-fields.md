# Subtask 02: Negative Options & Struct Fields Refactoring

Parent Plan: [142-boolean-principles-negatives-and-complex-conditions.md](../../pending/142-boolean-principles-negatives-and-complex-conditions.md)

## Goals
1. Refactor `cli/cmdclone/helpers.go`, `cli/cmd/codingguidelines_commit.go`, and `cli/cmd/clihelpers.go`: Rename `NoCommit` and `NoPush` in `CGCommitOpts` to `IsSkipCommit` and `IsSkipPush`, and update `emitCGSkipNotes` parameters.
2. Refactor `cli/cmd/gomod.go` and `cli/cmd/gomod_test.go`: Rename `noMerge` and `noTidy` in `goModOpts` to `isSkipMerge` and `isSkipTidy`, affirmative flag pointers, and `runGoModTidy(isSkipTidy bool)`.
3. Refactor `cli/cmd/latestbranch.go`: Rename `noFetch` to `isSkipFetch`, `shouldFetch` to `isFetchEnabled`, `shouldSwitch` to `isSwitchEnabled`, and `filterByRemote` to `isRemoteFiltered`.
4. Refactor `cli/movemerge/types.go`, `cli/cmd/movemergeflags.go`, `cli/movemerge/finalize.go`, and `cli/movemerge/integration_test.go`: Rename `IsNoPush` and `IsNoCommit` to `IsSkipPush` and `IsSkipCommit`.
5. Refactor `cli/cmd/historyrewrite_flags.go` and `cli/cmd/historyrewrite_push.go`: Rename `noPush` to `isSkipPush`.

## Acceptance Criteria
- [x] All negative options and struct fields across targeted packages refactored to affirmative `IsSkip*` / `isSkip*` / `is*Enabled`.
- [x] Call sites and tests updated and passing.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
