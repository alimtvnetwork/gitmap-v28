# Root Cause Analysis: Exhaustive Switch Missing Cases in CI Diff Gate

## 1. Context and Problem Statement

During CI/CD execution of the Lint Baseline Guard (`check-single-linter-diff.sh gitmap` with `LINTER=exhaustive`), the baseline diff gate failed with:
```text
Error: [exhaustive] missing cases in switch of type db.ResultStatusType: db.ResultStatusPending, db.ResultStatusDeferred, db.ResultStatusRequiresAuth (NEW vs baseline)
Error: [exhaustive] missing cases in switch of type cmd.replaceMode: cmd.replaceModeUnknown (NEW vs baseline)
FAIL: 2 new exhaustive finding(s). Fix the issues above.
```

## 2. Root Cause

1. `cmd/clustercommand.go` switched on `db.ResultStatusType` without covering pending, deferred, and auth-required status values.
2. `cmd/replace.go` switched on `replaceMode` without explicitly covering `replaceModeUnknown`.
3. Other unhandled enum switch cases existed across the repository in `cmd/dedupe.go`, `cmd/orphans.go`, `cmd/size.go`, `cmd/stale.go`, `movemerge/merge.go`, `render/pretty_emit.go`, `theme/theme.go`, and `vscodepm/mergemode.go`.

## 3. Corrective and Preventive Actions

- Updated `cmd/clustercommand.go` to explicitly branch on all `db.ResultStatusType` values.
- Updated `cmd/replace.go` to explicitly handle `case replaceModeUnknown:`.
- Updated all remaining switches in `cmd/dedupe.go`, `cmd/orphans.go`, `cmd/size.go`, `cmd/stale.go`, `movemerge/merge.go`, `render/pretty_emit.go`, `theme/theme.go`, and `vscodepm/mergemode.go` to explicitly handle all declared enum members.

## 4. Verification

- Ran `golangci-lint run --no-config --disable-all --enable=exhaustive ./...` and confirmed 0 findings remain across the entire codebase.
- Ran `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` and confirmed all 19 gates passed (exit code 0).
