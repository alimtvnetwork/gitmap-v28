# Specification 16 - Part 6: Interactive Mode Macro Recording, Playback & Scheduled Automation

## 1. Interactive Session Recording Engine

### 1.1 Concept & Workflow

Within `gitmap interactive` (or `gitmap i`), users can start a named macro recording session:
```text
  [REC] Start Recording Session: "deploy-and-sync"
  Commands entered in this prompt execute live in your terminal while being saved to macro history.
  Type `stop-record` or press Ctrl+D to finish and save.
```

### 1.2 Live Execution & Capture

While recording:
- Every typed command is executed immediately in the user's active sub-shell / PowerShell environment.
- Captures:
  1. Step sequence index (`step_num`).
  2. Raw command string (`command_line`).
  3. Working directory context (`working_dir`).
  4. Execution outcome / exit code.
  5. Elapsed execution time.
- Commands that fail can be interactively retried, discarded, or saved with error-tolerance flags (`continue_on_error`).

## 2. Macro Storage Schema

Macros are stored in the local SQLite database across two relational tables:

```sql
CREATE TABLE IF NOT EXISTS macros (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_steps INTEGER NOT NULL DEFAULT 0,
    tags TEXT
);

CREATE TABLE IF NOT EXISTS macro_steps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    macro_id INTEGER NOT NULL,
    step_num INTEGER NOT NULL,
    command_line TEXT NOT NULL,
    working_dir TEXT,
    continue_on_error INTEGER NOT NULL DEFAULT 0,
    timeout_seconds INTEGER NOT NULL DEFAULT 300,
    FOREIGN KEY(macro_id) REFERENCES macros(id) ON DELETE CASCADE
);
```

## 3. Macro Playback & CLI Interface

### 3.1 Execution Commands

- `gitmap execute <macro_name>` (alias: `gitmap macro run <name>`): Replays the recorded steps sequentially.
  - Supports `--dry-run` to preview commands without execution.
  - Supports `--verbose` for detailed step-by-step stdout streaming.
  - Emits summary box showing passed/failed/skipped steps.
- `gitmap macro list` (alias: `gitmap macro ls`): Lists all saved macros, step counts, and creation dates.
- `gitmap macro show <name>`: Prints individual steps of a macro.
- `gitmap macro rm <name>`: Deletes a saved macro.

## 4. Scheduled Macro Automation

### 4.1 Task Integration

Saved macros can be hooked into the existing `gitmap task` scheduler:
```bash

# Run macro every weekday at 9:00 AM

gitmap task add --name "daily-sync" --macro "deploy-and-sync" --cron "0 9 * * 1-5"

# Run macro every 30 minutes

gitmap task add --name "repo-watcher" --macro "pull-and-test" --interval 30m
```
