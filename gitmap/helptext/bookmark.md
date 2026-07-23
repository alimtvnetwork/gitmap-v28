# gitmap bookmark

Save and run bookmarked gitmap commands for quick re-execution.

## Alias

bk

## Usage

    gitmap bookmark <save|list|run|delete> [args]

## Flags

None.

## Prerequisites

- None

## Examples

### Example 1: Save a bookmark

    gitmap bookmark save "daily-scan" "scan D:\wp-work --quiet"

**Output:**

    ✓ Bookmark 'daily-scan' saved
    Command: scan D:\wp-work --quiet

### Example 2: List all bookmarks

    gitmap bk list

**Output:**

    NAME            COMMAND
    daily-scan      scan D:\wp-work --quiet
    work-pull       pull --group backend
    check-status    status --all
    3 bookmarks saved

### Example 3: Run a saved bookmark

    gitmap bookmark run daily-scan

**Output:**

    Running bookmark 'daily-scan'...
    → gitmap scan D:\wp-work --quiet
    Scanning D:\wp-work...
    Found 42 repositories
    ✓ Output written to ./.gitmap/output/

### Example 4: Delete a bookmark

    gitmap bookmark delete work-pull

**Output:**

    ✓ Bookmark 'work-pull' deleted

## See Also

- [history](history.md) — View command execution history
- [scan](scan.md) — Scan directories (common bookmark target)
- [pull](pull.md) — Pull repos (common bookmark target)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter bookmark
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
