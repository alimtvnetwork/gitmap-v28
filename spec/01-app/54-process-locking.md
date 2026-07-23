# Concurrent Process Locking

## Purpose

Prevent multiple gitmap-v27 instances from writing to the same SQLite
database simultaneously, which can cause corruption or lock errors.

## Mechanism

A **file-based lock** (`gitmap.lock`) is created in the database
directory when the database is opened. The lock file contains the
PID of the owning process.

## Behavior

1. On `openDBAt()`, attempt to create `gitmap.lock` exclusively.
2. If the lock file exists:
   - Read the PID from the file.
   - Check if that PID is still running.
   - If running → print warning and exit.
   - If not running (stale) → remove lock and proceed.
3. Write current PID to the lock file.
4. On `Close()`, remove the lock file.

## Output

```
  ⚠ Another gitmap-v27 process is running (PID 12345).
  If this is incorrect, delete: /path/to/data/gitmap.lock
```

## Implementation Files

| File | Action | Purpose |
|------|--------|---------|
| `store/lock.go` | CREATE | Lock acquisition, release, stale detection |
| `store/store.go` | MODIFY | Acquire lock in openDBAt, release in Close |
| `constants/constants_store.go` | MODIFY | Add lock file name and messages |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No magic strings — all in constants.
- Lock is advisory (file-based), not OS-level.
