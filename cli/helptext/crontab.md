# crontab

Manage scheduled cron tasks, intervals, isolated split SQLite databases, and background execution.

## Aliases

`schedule`, `sched`

## Usage

    gitmap crontab <subcommand> [flags]

## Subcommands

| Subcommand | Description |
|---|---|
| `ls`, `list` | List all scheduled tasks and crontabs |
| `add <name> [cmds...]` | Create a new scheduled task or cron job |
| `rm <name>` | Remove a scheduled task and delete its split database |
| `run <name>` | Execute a scheduled task immediately |
| `logs <name>` | Inspect execution history from the task's split database |
| `debug <name>` | Display diagnostics, split database paths, and execution telemetry |

## Examples

```bash
# List all crontab schedules
gitmap crontab ls

# Add a cron job to run every 30 minutes
gitmap crontab add sync-task "gitmap sync" --every 30m

# Debug a crontab schedule
gitmap crontab debug sync-task
```
