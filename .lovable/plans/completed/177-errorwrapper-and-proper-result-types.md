# Plan 177: Result ErrorWrapper & Centralized Proper Result Types Suite

## 1. Overview & Problem Statement
In previous iterations, command routing and dispatch functions returning `(error, bool)` were refactored to generic `result.Result[bool]`. However, returning a generic `Result[bool]` for operations that produce no data payload introduces generic parameter clutter at call sites and obscures the purpose of the return envelope. Furthermore, common scalar Result types (`BoolResult`, `StringResult`, `IntResult`, `Int64Result`) were not centralized in `cli/result/types.go`, forcing downstream packages to redundantly re-declare type aliases.

This plan introduces:
1. `result.ErrorWrapper`: A single, dedicated non-generic envelope containing `Err *apperror.AppError` and routing match state, providing inspection methods `.IsSuccess()`, `.IsFailed()`, `.IsFailure()`, `.IsInvalid()`, `.IsMatched()`, `HasError()`, and `.AppError()`.
2. Canonical Proper Result Type Aliases in `cli/result/types.go`: `BoolResult`, `StringResult`, `IntResult`, `Int64Result`, `Uint64Result`, `ByteSliceResult`, `StringSliceResult`, `AnyResult`.
3. Repository-wide refactoring of all command dispatchers and routers from `result.Result[bool]` to `result.ErrorWrapper`.

## 2. Task-Specific Rules & Architectural Constraints
1. **Single Dedicated Type for Error / Dispatch Outcomes:**
   - Any function that performs command routing or executes side-effects without returning a value payload must return `result.ErrorWrapper` rather than generic `result.Result[bool]`.
2. **Unified Inspection Interface:**
   - Callers can inspect `.IsSuccess()`, `.IsFailed()`, `.IsFailure()`, `.IsInvalid()`, `.IsMatched()`, `HasError()`, and `.AppError()`.
3. **Proper Named Aliases Mandate:**
   - Common generic instantiations must use the canonical aliases declared in `cli/result/types.go` (`result.BoolResult`, etc.).
4. **Universal AppError Retention:**
   - `*apperror.AppError` must be preserved throughout all routing layers and only converted to standard `error` at the final Cobra `RunE` boundary.
5. **Coding Guidelines Strict Adherence:**
   - Functions strictly <= 15 lines (target <= 8 lines).
   - Affirmative booleans only (`is*`, `has*`).
   - Zero nested ifs (nesting depth <= 1).
   - Strict Unix LF line endings.
   - TOTAL BAN on running `go test`, `go build`, or runner scripts.

## 3. Disjoint Subtasks Completed
- **Subtask 01**: `01-errorwrapper-type-and-core-result-aliases.md`
  - Files: `cli/result/types.go`, `cli/result/error_wrapper.go`, `cli/result/error_wrapper_test.go`
  - Implemented `ErrorWrapper` struct, methods, constructors, and centralized type aliases in `cli/result/`.
- **Subtask 02**: `02-cmd-cluster-env-group-profile-errorwrapper.md`
  - Files: `cli/cmd/cluster.go`, `cli/cmd/cluster_test.go`, `cli/cmd/env.go`, `cli/cmd/group.go`, `cli/cmd/profile.go`
  - Refactored all routers in `cli/cmd/` to return `result.ErrorWrapper`.
- **Subtask 03**: `03-cmdssh-cmdpipeline-macro-errorwrapper.md`
  - Files: `cli/cmdpipeline/pipeline_db_dispatch.go`, `cli/cmdssh/cluster_node_cmd.go`, `cli/cmdssh/cluster_node_cmd_test.go`, `cli/cmdssh/cluster_k8s_cmd.go`, `cli/cmdssh/ssh.go`, `cli/macro/execute.go`, `cli/archive/types.go`, `cli/cmdzsh/types.go`, `cli/dbengine/result_types.go`, `cli/osfix/types.go`, `cli/osuser/types.go`
  - Refactored routers in `cli/cmdpipeline/` and `cli/cmdssh/` to return `result.ErrorWrapper`, and standardized type aliases across the repository.

## 4. Verification & Quality Gates Results
- `go vet ./...`: 0 errors.
- `check-boolean-guidelines.py`: 0 violations across 2971 files.
- `check-nested-ifs.py`: 0 violations across 2971 files.
- `check-gosec-diff.py`: 0 new findings.
- `26-go-code-formatter.py`: 2924 Go files formatted cleanly.
- `04-newline-fixer.py`: 6902 files clean LF.
- `10-encoding-normalizer.py`: 6902 files normalized UTF-8.
- Test inventory updated via `33-test-inventory-generator.py --record` (19 files recorded).
