# Subtask 02: Manifest Clone Error Propagation & Structured Returns

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-02.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope

1. Eliminate all leaf exit calls (`cliexit.HandleError(..., 1)`) within `cli/cmdclone/clone_manifest.go`.
2. Convert `executeClone(cf CloneFlags)` and its inner helper routines (`cloneManifestBatch`, `runCloneExecution`) into structured functions returning explicit `error` types.
3. Guarantee that task failure recording in Split-DB (`store.DB` / `failPendingTask`) is performed prior to returning the error up the stack, preventing unclosed DB connections and orphaned pending tasks.
4. Enforce strict coding guidelines:
   - Zero magic integer exit codes.
   - Max 15 lines per function.
   - Guard clauses for error handling.
   - Strictly positive boolean naming (`hasError`, `isComplete`).

---

## 2. Detailed Technical Plan

### 2.1 Refactor Concurrency Resolution in `executeClone`
- In `cli/cmdclone/clone_manifest.go`:
  - Locate `executeClone(cf CloneFlags)`.
  - Update function signature to `func executeClone(cf CloneFlags) error`.
  - When `cloneconcurrency.Resolve(cf.MaxConcurrency)` returns `!ok`:
    - Remove `cliexit.HandleError(nil, 1)`.
    - Return `fmt.Errorf(constants.ErrCloneMaxConcurrencyInvalid, cf.MaxConcurrency)`.

### 2.2 Refactor Task DB Lifecycle & Failure Propagation
- When `runCloneExecution(cf)` returns `err != nil`:
  - Execute cleanup:
    ```go
    failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneFailed, cf.Source, err))
    closeTaskDB(taskDB)
    fmt.Fprintf(os.Stderr, constants.ErrCloneFailed, cf.Source, err)

    return err
    ```
  - Remove `cliexit.HandleError(err, 1)`.
- Ensure `finalizeCloneExecution(summary, cf, taskDB, taskID)` cleanly closes resources and returns `nil` or error to the outer caller.

### 2.3 Refactor `cloneManifestBatch` and Batch Helpers
- Audit `cloneManifestBatch` to ensure repository cloning failures across workers are aggregated into `model.CloneSummary` or returned as an aggregate error if all workers fail.
- Replace any unhandled error discards (`_ = ...`) with proper error checking or structured logging.
- Decompose any function exceeding 15 lines into concise helper routines.

---

## 3. Verification Criteria

- [ ] Zero invocations of `cliexit.HandleError` remain in `cli/cmdclone/clone_manifest.go`.
- [ ] `executeClone` returns an explicit `error` rather than terminating the process directly.
- [ ] Concurrency validation failure returns a typed error to the caller.
- [ ] Pending tasks are safely marked as failed and task DB connections are closed before returning error.
- [ ] All functions in `cli/cmdclone/clone_manifest.go` comply with the $\le 15$ line rule and blank line conventions.
