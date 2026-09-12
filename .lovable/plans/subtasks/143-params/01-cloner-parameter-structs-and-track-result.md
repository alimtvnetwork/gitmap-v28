# Subtask 01: Cloner Parameter Structs, TrackResult & Affirmative Options

Parent Plan: [143-argument-reduction-and-parameter-structs.md](../../pending/143-argument-reduction-and-parameter-structs.md)

## Goals
1. Refactor `cli/cloner/cloner.go`: Rename `CloneOptions` fields (`SafePull` -> `IsSafePull`, `Quiet` -> `IsQuiet`, `Clean` -> `IsClean`, `MissingOnly` -> `IsMissingOnly`). Rename `safePull` param in `CloneFromFile` and `CloneFromFileQuiet` to `isSafePull`.
2. Refactor `cli/cloner/runners.go`: Create `TrackResultParams` value struct and `SequentialRunParams` value struct. Refactor `trackResult` to `TrackResult(params TrackResultParams) *apperror.AppError` eliminating bare void return. Refactor `runSequential(params SequentialRunParams) model.CloneSummary`.
3. Refactor `cli/cloner/progress.go`: Rename `quiet` field and param to `isQuiet` in `NewProgress`. Rename `pulled` param to `isPulled` in `Done`.
4. Refactor `cli/cloner/safe_pull.go`: Update call sites of `opts.IsSafePull`.
5. Refactor `cli/cloner/batchprogress.go`: Rename `v bool` in `SetStopOnFail` to `isStopOnFail bool`.
6. Update external call sites in `cli/cmdclone/clone.go:529` and `cli/tests/heavy_test/cloner_hierarchy_e2e_test.go:112`.

## Acceptance Criteria
- [x] `TrackResultParams` and `SequentialRunParams` created and used.
- [x] `TrackResult` returns `*apperror.AppError`.
- [x] `CloneOptions` uses strict `Is*` affirmative boolean fields.
- [x] External call sites updated.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
