# CI/CD Issue 35: Exhaustive Switch Missing Cases in CI Diff Gate

- **Job**: Lint Baseline Guard (`check-single-linter-diff.sh gitmap` with `LINTER=exhaustive`)
- **Type**: FAIL
- **Detected**: 2026-08-29
- **Status**: resolved

## Error

```text
Error: [exhaustive] missing cases in switch of type db.ResultStatusType: db.ResultStatusPending, db.ResultStatusDeferred, db.ResultStatusRequiresAuth (NEW vs baseline)
Error: [exhaustive] missing cases in switch of type cmd.replaceMode: cmd.replaceModeUnknown (NEW vs baseline)
FAIL: 2 new exhaustive finding(s). Fix the issues above.
Error: Process completed with exit code 1.
```

## Root Cause

1. In `gitmap/cmd/clustercommand.go`, the `switch res.ResultStatus` on enum `db.ResultStatusType` only handled `ResultStatusSucceeded`, `ResultStatusFailed`, and `ResultStatusSkipped`, omitting `ResultStatusPending`, `ResultStatusDeferred`, and `ResultStatusRequiresAuth`.
2. In `gitmap/cmd/replace.go`, `dispatchReplaceMode` switched on `replaceMode` enum without an explicit `case replaceModeUnknown:`.
3. Other unhandled enum switch statements existed across the repository (`cmd/dedupe.go`, `cmd/orphans.go`, `cmd/size.go`, `cmd/stale.go`, `movemerge/merge.go`, `render/pretty_emit.go`, `theme/theme.go`, `vscodepm/mergemode.go`).

## Fix Applied

1. Updated `cmd/clustercommand.go` to explicitly handle all 6 enum cases (`ResultStatusSucceeded`, `ResultStatusFailed`, `ResultStatusRequiresAuth`, `ResultStatusSkipped`, `ResultStatusPending`, `ResultStatusDeferred`).
2. Updated `cmd/replace.go` to explicitly handle `case replaceModeUnknown:`.
3. Updated `cmd/dedupe.go`, `cmd/orphans.go`, `cmd/size.go`, and `cmd/stale.go` to explicitly include `case hygieneFormatTable:`.
4. Updated `movemerge/merge.go` to explicitly include `case DiffConflict:`.
5. Updated `render/pretty_emit.go` to explicitly handle remaining `blockKind` enum cases.
6. Updated `theme/theme.go` to explicitly include `case ModeBright:`.
7. Updated `vscodepm/mergemode.go` to explicitly include `case MergeModeUnion:`.
8. Verified with `golangci-lint run --no-config --disable-all --enable=exhaustive ./...` (0 warnings remain).
