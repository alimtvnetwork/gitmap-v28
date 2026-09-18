# Plan 176: Wrapped Result Dispatch Architecture & Universal AppError Retention

## 1. Executive Summary
Eliminated multi-value `(error, bool)` return tuples repository-wide across all CLI command dispatchers and internal routers. Replaced them with strongly-typed `result.Result[bool]` wrapped result envelopes. Preserved `*apperror.AppError` throughout all internal routing layers, avoiding premature conversions to standard library `error` until the final CLI entrypoint stage.

## 2. Completed Subtasks

### Subtask 01: Cluster Command Wrapped Result & AppError Retention
- Modified `cli/cmd/cluster.go`:
  - Updated all 11 routing functions (`routeClusterBootstrapOrExec`, `routeClusterScriptOrNode`, `routeClusterSSH`, `routeClusterJoinOps`, `routeClusterLegacyOps`, `routeClusterPasswordOps`, `routeClusterNodeOps`, `routeClusterK8s`, `routeClusterCore`, `routeClusterExt`, `dispatchClusterSubcommand`, `dispatchInvertedClusterHelp`) to return `result.Result[bool]`.
  - Refactored `dispatchInvertedClusterHelp` to propagate `result.Result[bool]` directly and return `result.RouteMatchedAppErr` with `*apperror.AppError`.
  - In `runCluster`, converted to standard `error` strictly at the final stage via `res.AppError()`.
- Updated `cli/cmd/cluster_test.go`:
  - Updated tests for `dispatchClusterSubcommand` and `dispatchInvertedClusterHelp` to assert on `res.Data` and `res.AppError()`.

### Subtask 02: Cmdssh Node, K8s, and SSH Wrapped Result & AppError Retention
- Modified `cli/cmdssh/cluster_node_cmd.go`:
  - Updated `routeClusterNodeLifecycle`, `routeClusterNodeRecipes`, and `routeClusterNodeCommand` to return `result.Result[bool]`.
  - In `RunClusterNodeCLI`, routed using `res.Data` and returned `res.AppError()`.
- Updated `cli/cmdssh/cluster_node_cmd_test.go`:
  - Updated `routeClusterNodeLifecycle` test to assert on `res.Data`.
- Modified `cli/cmdssh/cluster_k8s_cmd.go`:
  - Updated `routeClusterK8sBootstrap`, `routeClusterK8sNetwork`, `routeClusterK8sAddons`, `routeClusterK8sAdmin`, and `routeClusterK8sCommand` to return `result.Result[bool]`.
  - In `RunClusterK8sCLI`, routed using `res.Data` and returned `res.AppError()`.
- Modified `cli/cmdssh/ssh.go`:
  - Updated `dispatchPrimarySSH` to return `result.Result[bool]`.
  - In `dispatchSSH`, inspected `resPrimary.Data` and returned `resPrimary.AppError()`.

### Subtask 03: Env, Group, Profile, Pipeline DB, and Macro Wrapped Result
- Modified `cli/cmd/env.go`:
  - Updated `routeEnvVariableSub` to return `result.Result[bool]`.
  - In `routeEnvSub`, checked `resVar.Data` and returned `resVar.AppError()`.
- Modified `cli/cmd/group.go`:
  - Updated `dispatchGroupCRUD` and `dispatchGroupScoped` to return `result.Result[bool]`.
  - In `dispatchGroup`, inspected `res.Data` and returned `res.AppError()`.
- Modified `cli/cmd/profile.go`:
  - Updated all 7 profile routing functions (`routeInstallProfileSub`, `routeGitProfileSub`, `routeDBProfileSub`, `routeDBProfileExtra`, `routeChromeProfileSub`, `routeChromeImportSub`, `routeChromeExportSub`) to return `result.Result[bool]`.
  - In `routeProfileSub`, propagated `res.Data` and returned `res.AppError()`.
- Modified `cli/cmdpipeline/pipeline_db_dispatch.go`:
  - Updated `dispatchDbMutateSubcmd`, `dispatchDbQuerySubcmd`, and `dispatchPipelineDBSubcmd` to return `result.Result[bool]`.
  - In `handlePipelineDB`, routed with `res.Data` and returned `res.AppError()`.
- Modified `cli/macro/execute.go`:
  - Updated `executeReportStep` to return `result.Result[bool]` (`res.Data = shouldBreak`, `res.Err = apperror.WrapSimple(err, ...)`).
  - Cleanly eliminated the last remaining `(error, bool)` tuple in the repository.
- Modified `cli/result/result.go`:
  - Added constructors `RouteMatched(err error)`, `RouteMatchedAppErr(appErr *apperror.AppError)`, and `RouteUnmatched()` for clean, idiomatic result instantiation.
  - Added tests in `cli/result/result_test.go`.

## 3. Quality & Verification Gates
- `go vet ./...`: PASS (0 warnings, 0 errors).
- `check-boolean-guidelines.py`: PASS (0 violations across 2969 files).
- `check-nested-ifs.py`: PASS (0 violations across 2969 files).
- `check-gosec-diff.py`: PASS (0 new findings).
- `26-go-code-formatter.py`: Formatted 2922 Go files.
- `04-newline-fixer.py`: All 6898 files have clean newlines.
- `10-encoding-normalizer.py`: All 6898 text files verified.
- `33-test-inventory-generator.py`: Recorded 13 modified files, 349 associated test files.
