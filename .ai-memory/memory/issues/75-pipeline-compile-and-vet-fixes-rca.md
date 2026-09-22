# RCA 75: CI Pipeline Compilation, Vet, and Unused Import Fixes

## 1. Symptom
CI run `#35732822422` failed across Cross-Platform Build, History Rewrite Smoke, and race-detector steps with multiple Go compilation and `go vet` errors:
```text
clonenext/github.go:60:2: syntax error: non-declaration statement outside function body
archive/source.go:167:10: undefined: fmt
cmdagy/agy_clean_cache_apply.go:39:2: undefined: logProcessTermination
cmdagy/agy_clean_cache_apply.go:46:2: undefined: logCacheCleaning
cmdagy/agy_clean_cache_render.go:10:2: "time" imported and not used
cmdagy/agy_clean_cache_staging.go:6:2: "fmt" imported and not used
cmdagy/agy_rm_rejoin_ops.go:11:2: "github.com/alimtvnetwork/gitmap-v28/cli/workspacesync" imported and not used
cmdagy/agy_undo_cache.go:7:2: "fmt" imported and not used
cmd/commitin/orchestrator/commit.go:65:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/commit.go:74:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/commit.go:101:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/commit.go:233:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/pipeline.go:61:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/run.go:139:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/save_profile.go:46:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/save_profile.go:69:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/setup.go:52:23: fmt.Fprintf does not support error-wrapping directive %w
cmd/commitin/orchestrator/setup.go:59:23: fmt.Fprintf does not support error-wrapping directive %w
store/scan_folder.go:102:47: fmt.Sprintf does not support error-wrapping directive %w
store/scan_folder.go:106:47: fmt.Sprintf does not support error-wrapping directive %w
```

## 2. Root Cause
The pipeline failure was caused by:
1. `cli/clonenext/github.go`: Function signature `CreateRepo` was accidentally dropped during refactoring, causing lines in the body to appear outside a function.
2. `cli/archive/source.go`: Missing `"fmt"` import used by `fmt.Errorf` on line 167.
3. `cli/cmdagy/agy_clean_cache_apply.go`: `logProcessTermination` and `logCacheCleaning` helper functions were referenced but not implemented.
4. `cli/cmdagy/`: Unused imports (`time`, `fmt`, `workspacesync`) in `agy_clean_cache_render.go`, `agy_clean_cache_staging.go`, `agy_rm_rejoin_ops.go`, and `agy_undo_cache.go`.
5. `cli/constants/constants_commitin.go` & `cli/store/scan_folder.go`: Error string constants and formatting calls used the `%w` directive with `fmt.Fprintf` and `fmt.Sprintf` (which is rejected by `go vet` because `%w` is only supported in `fmt.Errorf`).

## 3. Resolution
1. In `cli/clonenext/github.go`, restored `CreateRepo` signature and decomposed into functions <= 15 lines (`CreateRepo`, `resolveRepoToken`, `attemptRepoCreation`).
2. In `cli/archive/source.go`, added `"fmt"` to the import block.
3. In `cli/cmdagy/agy_clean_cache_apply.go`, implemented `logProcessTermination` and `logCacheCleaning`.
4. In `cli/cmdagy/` files, removed unused imports (`time`, `fmt`, `workspacesync`).
5. In `cli/constants/constants_commitin.go`, replaced `%w` with `%v` in `CommitInErr*` constants passed to `fmt.Fprintf`.
6. In `cli/store/scan_folder.go`, updated `executeDetachAndDeleteSF` to pass op names directly to `apperror.WrapSimple` without wrapping `%w` inside `fmt.Sprintf`.
7. In `cli/constants/constants_scan_folder.go`, changed `%w` to `%v` in `ErrSFAbsResolve`.

## 4. Prevention & Learnings
- **Error Formatting Directive Rules**: Never use `%w` in string constants passed to `fmt.Fprintf`, `fmt.Sprintf`, or `fmt.Printf`; `%w` is exclusively valid in `fmt.Errorf`.
- **Pre-commit Compilation Sanity**: Ensure all newly modified files compile cleanly without dangling references or unimported packages.
