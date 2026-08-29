# Subtask 02: CloneFrom Execution Parameter Structs

- **Status:** `PENDING`
- **Target Files:** `gitmap/clonefrom/execute.go`, `gitmap/clonefrom/execute_concurrent.go`

## Instructions

1. Encapsulate `runRowLifecycle`, `prepareAndClone`, `tryLfsAutoFix`, `applyLfsFix`, and `writeProgress` in `gitmap/clonefrom/execute.go` into dedicated value-based structs (`RowLifecycleParams`, `PrepareCloneParams`, etc.).
2. Enforce strict `is`/`has` boolean prefixes on all fields (`isSafePull`, `isConcurrent`, `hasAuth`).
3. Ensure all side-effect helper functions return `*apperror.AppError` or standard typed results.
4. Verify all tests in `gitmap/clonefrom/` pass.
