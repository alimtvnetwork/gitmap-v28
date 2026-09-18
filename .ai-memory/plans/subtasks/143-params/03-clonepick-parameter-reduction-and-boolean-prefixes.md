# Subtask 03: ClonePick Parameter Reduction & Affirmative Booleans

Parent Plan: [143-argument-reduction-and-parameter-structs.md](../../pending/143-argument-reduction-and-parameter-structs.md)

## Goals
1. Refactor `cli/clonepick/parse.go`: Encapsulate `shorthandToURL` 4 loose parameters into `ShorthandURLParams`. Rename `askMode` to `isAskMode` in `normalisePaths`.
2. Refactor `cli/clonepick/replay.go`: Rename `dryRun` to `isDryRun` in `TouchAfterReplay`.
3. Refactor `cli/clonepick/sparse.go`: Rename `keepGit` to `isKeepGit` in `removeDotGitIfRequested`.
4. Refactor `cli/clonepick/picker_nav.go`: Encapsulate `clampScroll` 4 loose parameters into `ScrollBoundsParams`.
5. Refactor `cli/clonepick/picker_view.go`: Encapsulate `formatRow` 4 loose parameters into `FormatRowParams` (`IsSelected`, `IsCursor`).

## Acceptance Criteria
- [x] Parameter structs defined for all functions with >3 parameters in `cli/clonepick/`.
- [x] All boolean parameters renamed to affirmative `is*` prefixes.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
