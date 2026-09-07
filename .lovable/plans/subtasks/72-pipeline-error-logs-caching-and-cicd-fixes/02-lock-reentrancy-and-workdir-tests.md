# Subtask 02: Process Lock Reentrancy & WorkDir Tests

## Scope
- Update `gitmap/store/lock.go`:
  - Add thread-safe process-level lock tracking (`processLockMu sync.Mutex`, `processLockRefCounts map[string]int`).
  - In `handleExistingLock`:
    - If `pid == os.Getpid()`, grant reentrant access without failing with `ErrLockHeld`.
    - If `pid != os.Getpid() && processRunning(pid)`, return error.
  - In `releaseLock`:
    - Decrement reference count; remove `gitmap.lock` only when reference count reaches 0.
- Update `gitmap/cmd/workdir_test.go`:
  - Close `db` handles cleanly before calling functions that open the store.

## Files Touched
- `gitmap/store/lock.go`
- `gitmap/cmd/workdir_test.go`
- `gitmap/store/lock_test.go`
