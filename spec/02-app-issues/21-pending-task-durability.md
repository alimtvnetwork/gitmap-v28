# Pending Task Durability — Prevention Spec

## Rule

Every file or folder deletion in gitmap-v27 must be recorded as a database
task **before** the OS operation is attempted. This is a Code Red
requirement — no silent loss of delete intent is acceptable.

## Anti-Patterns (Prohibited)

1. Calling `os.RemoveAll` or `os.Remove` without a prior `PendingTask` insert.
2. Swallowing removal errors without updating `FailureReason`.
3. Deleting the `PendingTask` row without creating a `CompletedTask` row.
4. Using free-text task type strings instead of `TaskType` FK lookup.

## Required Pattern

```go
// 1. Record intent
taskID := store.InsertPendingTask(db, taskTypeID, targetPath, sourceCmd)

// 2. Attempt operation
err := os.RemoveAll(targetPath)

// 3. Record outcome
if err == nil {
    store.CompleteTask(db, taskID)  // transactional move
} else {
    store.FailTask(db, taskID, err.Error())
}
```

## Pre-Merge Checklist

- [ ] Every `os.Remove` / `os.RemoveAll` call is preceded by task insert.
- [ ] `CompleteTask` and `FailTask` are the only two exit paths.
- [ ] No duplicate pending task for same type + path.
- [ ] `golangci-lint` passes with zero errors.
- [ ] `go test ./store/...` covers insert, complete, fail, and list.

## Related

- `spec/01-app/95-pending-task-workflow.md` — Full specification.
