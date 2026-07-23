# Task Watch — File Sync Automation

## Overview

The `task` command introduces named, persistent file-sync tasks.
Each task defines a **source folder** and a **destination folder**
with one-way, timestamp-based synchronization. Tasks run in the
foreground with parallel goroutines checking file modification times
at a configurable interval (default 5 seconds).

---

## How It Works

1. User creates a named task: source path, destination path.
2. Task definition is saved to `.gitmap/tasks.json`.
3. When run, the watcher spawns one goroutine per file.
4. Every interval (default 5s), each goroutine compares `ModTime`.
5. If the source file is newer, it replaces the destination file.
6. Files matching `.gitignore` patterns in the source folder are skipped.
7. New files in source are copied; deleted files are NOT removed from dest.

---

## Commands

### `gitmap-v28 task` (alias: `tk`)

Manage and execute file-sync watch tasks.

```bash
gitmap-v28 task create <name> --src <path> --dest <path>
gitmap-v28 task list
gitmap-v28 task run <name>
gitmap-v28 task run <name> --interval 10
gitmap-v28 task delete <name>
gitmap-v28 task show <name>
```

---

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `create`   | Create a new named sync task |
| `list`     | List all saved tasks |
| `run`      | Start watching and syncing for a task |
| `show`     | Display task details (source, dest, file count) |
| `delete`   | Remove a saved task |

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--src <path>` | — | Source folder (required for `create`) |
| `--dest <path>` | — | Destination folder (required for `create`) |
| `--interval <seconds>` | `5` | Check interval in seconds (minimum 2) |
| `--verbose` | `false` | Log each file copy to stderr |
| `--dry-run` | `false` | Show what would be synced without copying |

---

## Task Storage

Tasks persist in `.gitmap/tasks.json`:

```json
{
  "tasks": [
    {
      "name": "frontend-sync",
      "source": "C:\\repos\\frontend\\src",
      "destination": "C:\\repos\\app\\vendor\\frontend",
      "created": "2026-04-04T10:00:00Z"
    }
  ]
}
```

---

## Sync Algorithm

1. Walk the source directory tree recursively.
2. For each file, check against `.gitignore` rules (if present in source root).
3. Compare `ModTime` of source file vs destination counterpart.
4. If source is newer or destination does not exist, copy the file.
5. Preserve relative directory structure in destination.
6. Deletions in source are NOT propagated (additive sync only).

### Goroutine Model

- One goroutine per top-level directory entry.
- Each goroutine handles its subtree sequentially.
- A shared `sync.WaitGroup` tracks completion of each cycle.
- File copy uses `io.Copy` with a 32 KB buffer.

---

## Cross-Platform Notes

### Windows (Primary)
- Paths use backslash; `filepath.Walk` handles this natively.
- Long paths: warn if path > 260 chars (suggest `core.longpaths`).

### Linux / macOS
- Forward-slash paths; same `filepath.Walk` logic.
- File permissions preserved via `os.Stat` mode bits.

---

## Dashboard Output

```
gitmap-v28 task run frontend-sync — checking every 5s (Ctrl+C to stop)
Last sync: 2026-04-04 10:05:32

  Watched: 142 files | Synced: 3 | Skipped (gitignore): 28
  src/App.tsx          -> synced (newer by 12s)
  src/index.css        -> synced (newer by 4s)
  src/utils/helpers.ts -> synced (new file)
```

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Source folder missing | Exit with error message |
| Destination folder missing | Create it automatically |
| File locked (Windows) | Skip with warning, retry next cycle |
| Permission denied | Log warning, continue |

---

## File Layout

| File | Purpose |
|------|---------|
| `constants/constants_task.go` | Command names, defaults, messages |
| `cmd/task.go` | Subcommand routing and flag parsing |
| `cmd/taskcreate.go` | Task creation logic |
| `cmd/taskrun.go` | Watch loop and sync engine |
| `cmd/taskformat.go` | Dashboard display formatting |
| `model/task.go` | Task struct and JSON serialization |

---

## Constraints

- Minimum interval is 2 seconds.
- Maximum concurrent goroutines capped at 64.
- `.gitignore` parsing uses simple glob matching (no nested `.gitignore`).
- All files under 200 lines, all functions 8-15 lines.
- One-way sync only (source -> destination).
- No content comparison; timestamp-only decision.

---

## Examples

```bash
# Create a task
gitmap-v28 task create ui-sync --src ./frontend/src --dest ./backend/static

# List all tasks
gitmap-v28 task list

# Run a task with default 5s interval
gitmap-v28 task run ui-sync

# Run with faster polling and verbose output
gitmap-v28 task run ui-sync --interval 2 --verbose

# Preview what would sync
gitmap-v28 task run ui-sync --dry-run

# Delete a task
gitmap-v28 task delete ui-sync
```

---

## See Also

- watch — live repo status monitor
- exec — run commands across repos
- group — manage repo collections
