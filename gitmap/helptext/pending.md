# Pending Tasks

List all pending tasks that have not yet completed successfully.

## Usage

    gitmap pending

## Description

Displays all tasks currently in the PendingTask table. A task is
pending if it was enqueued but either failed or has not yet been
attempted. Each task shows its ID, type, target path, and the
last known failure reason (if any).

## Output Columns

| Column  | Description                              |
|---------|------------------------------------------|
| ID      | Auto-incrementing task identifier        |
| Type    | Task category (Delete, Scan, Clone, etc) |
| Path    | Target file or folder path               |
| Reason  | Last failure reason (empty if untried)   |

## Examples

### List pending tasks

    $ gitmap pending
      Pending Tasks:
        #1      Delete   D:\wp-work\riseup-asia\scripts-fixer-v5   retry removal failed: access denied
        #2      Clone    D:\projects\my-repo                        command replay failed: network timeout

### No pending tasks

    $ gitmap pending
      No pending tasks.

## See Also

- [do-pending](do-pending.md) — Retry pending tasks
- [clone-next](clone-next.md) — Clone next versioned iteration
- [task](task.md) — Manage file-sync watch tasks

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pending
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
