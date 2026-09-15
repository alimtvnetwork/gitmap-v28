---
name: result-errorwrapper-and-proper-types
description: Autonomously implement and verify ErrorWrapper type and centralized proper Result aliases, eliminating raw Result[bool] returns and providing IsSuccess, IsFailed, IsInvalid, and AppError inspection across GitMap.
---

# Skill: Result ErrorWrapper & Proper Result Types Suite (`result-errorwrapper-and-proper-types`)

This skill governs the autonomous design, implementation, refactoring, and verification of the non-generic `result.ErrorWrapper` and canonical proper type aliases (`BoolResult`, `StringResult`, `IntResult`, `Int64Result`, etc.) across the GitMap codebase.

## Core Directives

1. **Dedicated Non-Generic `ErrorWrapper` Type:**
   - In `cli/result/types.go` and `cli/result/error_wrapper.go`, define `ErrorWrapper`:
     - `Err *apperror.AppError`
     - `isMatched bool`
   - Methods:
     - `IsSuccess() bool`: True if matched and `Err == nil`
     - `IsFailed() bool` / `IsFailure() bool`: True if receiver is nil or `Err != nil`
     - `IsInvalid() bool`: True if receiver is nil, unmatched, or `Err != nil`
     - `IsMatched() bool`: True if matched
     - `HasError() bool`: True if receiver is nil or `Err != nil`
     - `AppError() *apperror.AppError`: Returns underlying `*apperror.AppError`
     - `Error() string`: Returns error string or empty string
   - Constructors:
     - `SuccessWrapper()`: Returns matched `ErrorWrapper` without error.
     - `FailureWrapper(err *apperror.AppError)`: Returns matched `ErrorWrapper` with `AppError`.
     - `FailureWrapperErr(err error)`: Wraps standard `error` into `AppError` and returns matched `ErrorWrapper`.
     - `UnmatchedWrapper()`: Returns unmatched `ErrorWrapper`.

2. **Centralized Proper Type Aliases in `cli/result/types.go`:**
   - `BoolResult = Result[bool]`
   - `StringResult = Result[string]`
   - `IntResult = Result[int]`
   - `Int64Result = Result[int64]`
   - `Uint64Result = Result[uint64]`
   - `ByteSliceResult = Result[[]byte]`
   - `StringSliceResult = ResultSlice[string]`
   - `AnyResult = Result[any]`
   - `ErrorWrap = ErrorWrapper`

3. **Universal Dispatch & Router Refactoring:**
   - Replace all raw `result.Result[bool]` return signatures across command dispatchers and routers with `result.ErrorWrapper`:
     - `cli/cmdpipeline/pipeline_db_dispatch.go`: `dispatchDbMutateSubcmd`, `dispatchDbQuerySubcmd`, `dispatchPipelineDBSubcmd`
     - `cli/cmd/cluster.go`: `dispatchClusterSubcommand`, `dispatchInvertedClusterHelp`, and all `routeCluster*` functions
     - `cli/cmdssh/cluster_node_cmd.go`: `routeClusterNodeLifecycle`, `routeClusterNodeRecipes`, `routeClusterNodeCommand`
     - `cli/cmdssh/cluster_k8s_cmd.go`: `routeClusterK8sBootstrap`, `routeClusterK8sNetwork`, `routeClusterK8sAddons`, `routeClusterK8sAdmin`, `routeClusterK8sCommand`
     - `cli/cmdssh/ssh.go`: `dispatchPrimarySSH`
     - `cli/cmd/env.go`: `routeEnvVariableSub`
     - `cli/cmd/group.go`: `dispatchGroupCRUD`, `dispatchGroupScoped`
     - `cli/cmd/profile.go`: all 7 `route*Profile*` functions
     - `cli/macro/execute.go`: `executeReportStep` -> `result.BoolResult`

4. **Coding Guidelines Enforcement:**
   - Functions strictly <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`).
   - Zero nested ifs (nesting depth <= 1).
   - Strict Unix LF line endings.
   - Total ban on running tests (`go test`, runner scripts) or build commands (`go build`).
