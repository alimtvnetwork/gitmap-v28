# gitmap agy running-prompts

Snapshot, inspect, restore, and govern active and enqueued Antigravity AI prompts across all workspaces.

## Synopsis

    gitmap agy running-prompts <command> [flags]
    gitmap backup-running-prompts [flags]
    gitmap restore-running-prompts [flags]
    gitmap agy running-prompts ls [--limit <N>] [--wc <N>] [--full] [--json] [--ssh]
    gitmap agy running-prompts backup [--file <path>] [--json] [--ssh]
    gitmap agy running-prompts restore [--keep] [--file <path>] [--json] [--ssh]
    gitmap agy running-prompts clean [--force]

## Description

`gitmap agy running-prompts` provides unified prompt lifecycle management for Google Antigravity (AGY) workflows.

### 1. Active & Queued Prompts Inspection (`ls`)
Scans all registered Antigravity projects and current workspace queues (`agy-prompt-queue.json` and active pipeline fix text files). Displays project name, execution status (`RUNNING` or `ENQUEUED`), word count, and prompt preview.

### 2. Split-DB Backup Snapshot (`backup` / `backup-running-prompts`)
Safely captures all active running and enqueued prompts across projects into an isolated SQLite Split-DB (`data/backup-prompts/sql.db`). Each backup creates an immutable batch record (`b-<hash>`) with project path, sequence ID, prompt text, and word counts. Enables developers and automated runners to safely power off or reboot machines without losing queued prompts.

### 3. Queue Restoration & Injection (`restore` / `restore-running-prompts`)
Restores backed-up prompt items from the latest or specified batch and re-injects them directly back into each target project's queue file (`.ai-memory/temp/agy-prompt-queue.json` or root `agy-prompt-queue.json`). Injected entries are marked with status `queued` and type `restored_prompt` so agent loops pick them up seamlessly.

### 4. Retention & Cleanup (`clean`)
By default, restored entries older than 24 hours that are not flagged with `--keep` are pruned to keep the Split-DB compact. Passing `--force` purges all historical backup records.

## Subcommands

| Subcommand | Description |
| :--- | :--- |
| `ls` | Inspect and list running and queued prompts across projects. |
| `backup` | Snapshot all active and queued prompts to SQLite Split-DB. |
| `restore` | Restore backed-up prompts from database into project prompt queues. |
| `clean` | Prune expired restored batches (or force purge all batches). |
| `export` | Export running prompts to an external SQLite `.db` or JSON `.json` file. |
| `import` | Import prompts from an external SQLite or JSON file into queues. |
| `help` | Render interactive two-column help menu. |

## Flags

| Flag | Shorthand | Default | Description |
| :--- | :--- | :--- | :--- |
| `--limit <N>` | `-l` | 8 | Maximum prompts to display in table output. |
| `--wc <N>` | *(none)* | 100 | Truncate prompt text preview to N words. |
| `--full` | *(none)* | false | Display full prompt text without word truncation. |
| `--json` | `-j` | false | Emit structured JSON output. |
| `--ssh` | *(none)* | false | Query or backup prompts across multi-node cluster fleet. |
| `--file <path>` | `-f` | "" | Target database or import/export file path. |
| `--keep` | `-k` | false | Exclude restored batch from 24-hour automatic retention pruning. |
| `--force` | *(none)* | false | Force delete all backup batches during clean. |
| `--help` | `-h` | false | Show command help menu. |

## Examples

### Example 1: Inspect Running & Queued Prompts
```bash
gitmap agy running-prompts ls
```

**Output:**
```text
  ● Antigravity Running & Queued Prompts (4 items)
  PROJECT              STATUS     WORDS    PROMPT
  ────────────────────────────────────────────────────────────────────────────────
  gitmap               RUNNING    45       Audit and refactor fleet update JSON...
  Antigravity-Manager  ENQUEUED   32       Update electron window layout tokens...
  scripts-fixer        ENQUEUED   28       Run test-matrix and verify linter...
  D:\test-gitmap       RUNNING    52       Verify clone cache direct tree replay...
```

### Example 2: Snapshot All Prompts Before Machine Reboot
```bash
gitmap backup-running-prompts
```

### Example 3: Re-inject Backed-up Prompts After Reboot
```bash
gitmap restore-running-prompts
```

### Example 4: List Historical Backup Batches
```bash
gitmap agy backup-running-prompts ls
```

## See Also

- `gitmap agy fpug` — finish prompts until green
- `gitmap agy sug` — shutdown system once all project pipelines are green
- `gitmap rerun` — rerun latest prompts by sequence number
