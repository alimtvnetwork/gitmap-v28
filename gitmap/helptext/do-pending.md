# Retry Pending Tasks

Retry all pending tasks or a specific task by its ID.

## Alias

dp

## Usage

    gitmap do-pending [task-id]

## Description

Re-executes tasks that remain in the PendingTask table. For delete
and remove tasks, the system retries the file removal. For replayable
tasks (scan, clone, pull, exec), the system re-runs the original CLI
command using the stored arguments and working directory.

Successfully completed tasks are moved to the CompletedTask table.
Failed tasks remain pending with an updated failure reason.

## Arguments

| Argument | Default | Description                          |
|----------|---------|--------------------------------------|
| task-id  | (all)   | Retry only this specific task by ID  |

## Examples

### Retry all pending tasks

    $ gitmap do-pending
      Retrying 2 pending task(s)...
      Task #1 completed: D:\wp-work\riseup-asia\scripts-fixer-v5
      Task #2 failed: command replay failed: network timeout

### Retry using alias

    $ gitmap dp
      Retrying 1 pending task(s)...
      Task #3 completed: D:\projects\my-repo

### Retry a specific task by ID

    $ gitmap do-pending 2
      Retrying task #2...
      Replaying: gitmap clone json
      Task #2 completed: clone json

### No pending tasks

    $ gitmap dp
      No pending tasks.

## See Also

- [pending](pending.md) — List all pending tasks
- [clone-next](clone-next.md) — Clone next versioned iteration
- [task](task.md) — Manage file-sync watch tasks

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter do-pending
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
