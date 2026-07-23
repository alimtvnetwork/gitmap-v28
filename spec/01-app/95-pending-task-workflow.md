# Pending Task Workflow

## Overview

Every file-system operation in gitmap-v27 (delete, remove, scan, clone, pull,
exec) must be recorded as a task in SQLite **before** the actual operation
is attempted. This ensures no work is silently lost when an operation fails
due to locks, permissions, or missing targets.

The system stores the full CLI command so that failed tasks can be
automatically replayed via `gitmap-v27 do-pending`.

## Database Schema

All tables use PascalCase. Primary keys are auto-incrementing integers.

### TaskType

Normalized lookup table for task categories.

```sql
CREATE TABLE IF NOT EXISTS TaskType (
    Id   INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);
```

Seed values: `Delete`, `Remove`, `Scan`, `Clone`, `Pull`, `Exec`.

### PendingTask

Holds every task that has not yet completed successfully.

```sql
CREATE TABLE IF NOT EXISTS PendingTask (
    Id               INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskTypeId       INTEGER NOT NULL REFERENCES TaskType(Id),
    TargetPath       TEXT    NOT NULL,
    WorkingDirectory TEXT    DEFAULT '',
    SourceCommand    TEXT    NOT NULL,
    CommandArgs      TEXT    DEFAULT '',
    FailureReason    TEXT    DEFAULT '',
    CreatedAt        TEXT    DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt        TEXT    DEFAULT CURRENT_TIMESTAMP
);
```

### CompletedTask

Archive of successfully executed tasks.

```sql
CREATE TABLE IF NOT EXISTS CompletedTask (
    Id               INTEGER PRIMARY KEY AUTOINCREMENT,
    OriginalTaskId   INTEGER NOT NULL,
    TaskTypeId       INTEGER NOT NULL REFERENCES TaskType(Id),
    TargetPath       TEXT    NOT NULL,
    WorkingDirectory TEXT    DEFAULT '',
    SourceCommand    TEXT    NOT NULL,
    CommandArgs      TEXT    DEFAULT '',
    CompletedAt      TEXT    DEFAULT CURRENT_TIMESTAMP,
    CreatedAt        TEXT    NOT NULL
);
```

### Column Descriptions

| Column | Purpose |
|--------|---------|
| TargetPath | File or folder being acted on |
| WorkingDirectory | CWD when the command was invoked |
| SourceCommand | Command name that initiated the task (e.g. `clone-next`) |
| CommandArgs | Full CLI arguments for replay (e.g. `scan --mode ssh ./repos`) |
| FailureReason | Human-readable reason for last failure |

## Execution Lifecycle

### 1. Task Creation

Before any file-system operation:

1. Resolve `TaskTypeId` from `TaskType`.
2. Capture current working directory.
3. Build `CommandArgs` from `os.Args`.
4. Insert row into `PendingTask`.
5. Only after successful insert, attempt the actual operation.

### 2. Success Path

1. Insert row into `CompletedTask` copying all fields + `CompletedAt`.
2. Delete row from `PendingTask`.
3. Both steps inside a single transaction.

### 3. Failure Path

1. Row stays in `PendingTask`.
2. Update `FailureReason` with human-readable context.
3. Update `UpdatedAt` to current timestamp.

### 4. Replay via do-pending

For replayable types (Scan, Clone, Pull, Exec):
1. Re-execute `gitmap-v27 <CommandArgs>` as a subprocess.
2. Set `cmd.Dir` to `WorkingDirectory`.
3. On success → complete. On failure → update reason.

For delete types (Delete, Remove):
1. Attempt `os.RemoveAll(TargetPath)`.
2. On success → complete. On failure → update reason.

## CLI Commands

### `gitmap-v27 pending`

Display all rows in `PendingTask`.

Output columns: Id, Type, TargetPath, FailureReason.

### `gitmap-v27 do-pending` (alias `dp`)

Retry all pending tasks. Each success moves to `CompletedTask`.
Each failure updates `FailureReason` and remains pending.

### `gitmap-v27 do-pending <id>`

Retry a single pending task by its integer Id.

## Duplicate Prevention

If a `PendingTask` already exists with the same `TaskTypeId` and
`TargetPath`, do not create a duplicate. Log a message indicating the
existing pending task Id.

For replayable task types (Scan, Clone, Pull, Exec), duplicate detection
also includes `CommandArgs`. This allows the same target path to have
multiple pending tasks if the CLI arguments differ (e.g., `scan --mode ssh`
vs `scan --mode https` on the same directory).

Delete and Remove tasks match on type+path only, since the operation
is always the same regardless of how the delete was triggered.

## Task Types Covered

| Type | Commands | Replay Method |
|------|----------|---------------|
| Delete | clone-next --delete | os.RemoveAll |
| Remove | future removal paths | os.RemoveAll |
| Scan | scan, rescan | subprocess replay |
| Clone | clone | subprocess replay |
| Pull | pull | subprocess replay |
| Exec | exec | subprocess replay |

## Integrated Commands

Commands that have been wired into the pending task system:

### scan (first integration)

The `executeScan` function enqueues a `Scan` task before starting
directory traversal. The target path is the absolute scan directory.
CLI args are captured from `os.Args` so that `do-pending` can replay
the exact scan invocation. On successful completion (all outputs
written, DB upserted, projects detected), the task is marked complete.
If `ScanDir` fails, the failure reason is recorded and the task remains
pending.

### clone (second integration)

The `executeClone` function enqueues a `Clone` task before calling
`CloneFromFile`. The target path is the absolute target directory.
On success (all repos cloned, summary printed, Desktop registration
done), the task is marked complete. If `CloneFromFile` returns an
error, the failure reason is recorded.

### pull (third integration)

The `runPull` function enqueues a `Pull` task after resolving targets
but before starting the batch pull loop. For single-repo pulls, the
target path is the repo's absolute path; for multi-repo pulls (--all
or --group), it falls back to the working directory. CLI args are
captured from `os.Args` for replay. If the batch exits with a non-zero
code, the task is failed with the exit code. On full success, the task
is marked complete.

### exec (fourth integration)

The `runExec` function enqueues an `Exec` task after parsing flags and
loading target records but before starting the batch execution loop.
The target path is the current working directory. CLI args are captured
from `os.Args` for replay. If the batch exits with a non-zero code,
the task is failed with the exit code. On full success, the task is
marked complete.

## Edge Case Handling

| Scenario | Behavior |
|----------|----------|
| Target path does not exist (delete) | Task auto-completes, logs skip |
| Working directory missing (replay) | Task fails with Code Red path diagnostic |
| Permission denied | Failure reason includes path, operation, and OS error |
| Empty CommandArgs (replay) | Task fails immediately, reason recorded |
| FailTask on non-existent ID | Returns error (RowsAffected = 0) |
| CompleteTask commit failure | Transaction rolls back, error logged |

## Constants Organization

| File | Content |
|------|---------|
| `constants_pending_task.go` | Table names, type seeds, CREATE/ALTER SQL |
| `constants_pending_task_sql.go` | CRUD query SQL |
| `constants_pending_task_msg.go` | Error/warning/message constants |

## Acceptance Criteria

1. Every file-system operation inserts into `PendingTask` before execution.
2. No operation path bypasses task creation.
3. Failed operations remain visible in `PendingTask`.
4. Successful operations appear in `CompletedTask` and are removed from `PendingTask`.
5. `gitmap-v27 pending` lists all pending tasks with Id, type, path, reason.
6. `gitmap-v27 do-pending` retries all; `gitmap-v27 dp` is an alias.
7. `gitmap-v27 do-pending <id>` retries a single task.
8. Duplicate pending tasks for the same type+path are prevented.
9. Replayable tasks use type+path+cmdArgs for duplicate detection.
10. Replayable tasks store full CLI args and working directory.
11. All commands appear in standard help, detailed help, and UI help.
12. Delete tasks auto-complete when target path no longer exists.
13. Permission errors include full Code Red path diagnostics.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect.
