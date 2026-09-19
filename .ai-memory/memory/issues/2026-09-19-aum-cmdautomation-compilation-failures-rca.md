# 4-Part Root Cause Analysis (RCA): AUM Polyglot Worker Compilation Failures in CI (Commit cc7e36f)

## 1. Why It Happened (High-Level Business & Architectural Context)
In commit `cc7e36f` (Plan 38), we implemented the Go supervisor and polyglot worker orchestrator for `gitmap aum run`, `runtimes`, and local CI/CD testing tooling per Spec 124. During rapid multi-file scaffolding across `cli/cmdautomation`, several function identifier collisions, struct field name mismatches, undefined signatures, and Go build exclusion by filename convention (`smoke_test.go` ending in `_test.go`) caused `go build ./...` compilation failures in GitHub Actions CI runs (`#35448965867`, `#35448965838`, `#35448965806`).

---

## 2. How It Happened (Technical Execution Flow)
1. **Identifier Collisions within Package `cmdautomation`**:
   - `collectSearchFiles` was declared in both `search.go:47` and `spec_migrate.go:145`.
   - `resolveThreshold` was declared in both `plan_consolidate.go:39` (taking `int`) and `test_inventory.go:42` (taking `float64`).
2. **Undefined Identifiers in `runtimes_cmd.go`**:
   - `runtimes_cmd.go` referenced `ListRuntimes()`, `RefreshRuntimes()`, and `RuntimeInfo`, but the struct in `worker_types.go` was named `RuntimeRecord`, and the helper functions had not been declared in `runtime_probe.go`.
3. **Typo in Function Call**:
   - `preflight.go:121` called `RunRelPathsAudit` while the actual function name in `relpaths.go:18` was `RunRelPathAudit`.
4. **Type Mismatch and Missing Fields in `run_cmd.go`**:
   - `RunWorkerPool(opts)` returned `WorkerRunResultMonad` (`result.Result[WorkerRunResult]`), but was passed into `toError(*apperror.AppError)`.
   - `WorkerRunOptions` lacked `IsJson`, and `opts.Target` was accessed instead of `opts.Script`.
5. **Go Build Exclusion via `_test.go` Suffix**:
   - `smoke_test.go` contained implementation code (`RunSmokeTest`, `aggregateSmokeResult`, etc.) rather than unit tests. In Go, any file matching `*_test.go` is excluded from regular package builds (`go build ./...`), causing `smoke_test_cmd.go` to fail with `undefined: RunSmokeTest`.
6. **Result Wrapper Alignment & Spec Directory Resolution**:
   - `BuildFileContext`, `EncodeToStream`, and `DecodeFromStream` signatures in `stream_encoding.go` and `worker_pool.go` required alignment with `cg-result-wrapper` and `worker_test.go`.
   - `resolveSpecDir` in `spec_migrate.go` did not inspect `02-spec/21-app` within custom directory paths, causing `TestSpecMigrate_ResequencesAndUpdates` to fail.

---

## 3. Root Cause (Exact File, Line, and Mechanism)
- `cli/cmdautomation/spec_migrate.go:145:6`: Redeclaration of `collectSearchFiles` (collision with `search.go:47:6`).
- `cli/cmdautomation/test_inventory.go:42:6`: Redeclaration of `resolveThreshold` (collision with `plan_consolidate.go:39:6`).
- `cli/cmdautomation/runtimes_cmd.go:34, 43, 52, 70`: Referenced undefined `ListRuntimes`, `RefreshRuntimes`, and `RuntimeInfo`.
- `cli/cmdautomation/preflight.go:121:11`: Called `RunRelPathsAudit` instead of `RunRelPathAudit`.
- `cli/cmdautomation/run_cmd.go:26, 42, 51`: Type mismatch with `toError`, undefined `opts.Target`, undefined `workerOpts.AsJson`.
- `cli/cmdautomation/smoke_test.go`: File named with `_test.go` suffix despite containing non-test implementation code.
- `cli/cmdautomation/stream_encoding.go` & `worker_pool.go`: Missing monad wrapper returns and signature desync with `worker_test.go`.
- `cli/cmdautomation/spec_migrate.go:36`: `resolveSpecDir` returned `customDir` directly without checking child spec directories.

---

## 4. Code Fix (Exact Changes)

### Fix 1: Rename `collectSearchFiles` in `cli/cmdautomation/spec_migrate.go`
Renamed to `collectSpecMigrateFiles(root string) []string`.

### Fix 2: Rename `resolveThreshold` in `cli/cmdautomation/test_inventory.go`
Renamed to `resolveSlowThreshold(val float64) float64`.

### Fix 3: Add `ListRuntimes` & `RefreshRuntimes` in `cli/cmdautomation/runtime_probe.go` and Update `runtimes_cmd.go`
Implemented `ListRuntimes() RuntimeListResultMonad`, `RefreshRuntimes() RuntimeListResultMonad`, and `reprobeRuntime()`. Updated `runtimes_cmd.go` to use `RuntimeRecord`.

### Fix 4: Fix Call in `cli/cmdautomation/preflight.go`
Updated `RunRelPathsAudit` to `RunRelPathAudit`.

### Fix 5: Update `WorkerRunOptions` and `cli/cmdautomation/run_cmd.go`
Added `IsJson bool` to `WorkerRunOptions`. Mapped `opts.Script = args[1]`. Checked `res.IsFailure()` and added `renderWorkerResult`.

### Fix 6: Rename `cli/cmdautomation/smoke_test.go` to `smoke_runner.go`
Renamed via `git mv cli/cmdautomation/smoke_test.go cli/cmdautomation/smoke_runner.go` so it is compiled into the main package.

### Fix 7: Align `stream_encoding.go` and `worker_pool.go` with `result.Result` Wrappers
Updated `EncodeToStream` and `DecodeFromStream` to return `result.Result[T]`. Refactored `BuildFileContext` into modular helpers adhering to 8–15 line limits.

### Fix 8: Fix `resolveSpecDir` in `cli/cmdautomation/spec_migrate.go`
Updated `resolveSpecDir` to inspect `02-spec/21-app` and `02-spec` in the base directory.
