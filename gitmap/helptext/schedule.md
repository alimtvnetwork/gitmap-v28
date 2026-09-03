# schedule

Manage scheduled tasks, automated macro executions, interval runners, OS startup triggers, isolated split SQLite databases, execution logs, multi-format import/export, and Web UI/API integrations.

## Aliases

`sched`

---

## Subcommands

| Subcommand | Description |
|---|---|
| `add` (`create`, `new`) `<name> [cmds...]` | Create a new scheduled task (linked to a macro or on-the-fly shell commands with isolated split DB) |
| `list` (`ls`) `[--json\|--yaml] [-f <file>]` | List all scheduled tasks with status (enabled/disabled), interval, target, runs count, startup, and split DB slug |
| `status` `[name\|*] [--json\|--yaml] [-f <file>]` | View detailed status of a specific schedule or all schedules (when it ran, enabled status, split DB path) |
| `enable` `<name>` | Enable a disabled scheduled task |
| `disable` `<name>` | Disable a scheduled task (preserves configuration and logs) |
| `logs` (`log`, `history`) `<name> [--limit <N>] [--json\|--yaml] [-f <file>]` | Inspect execution logs and run history from the schedule's dedicated split SQLite database |
| `export` `[name\|*] [-f <file>] [--format=json\|yaml\|sqlite\|zip]` | Export single schedule or all schedules with run logs to JSON, YAML, SQLite, or ZIP |
| `export-all` `[-f <file>] [--json\|--yaml\|--sqlite\|--zip] [-except "name1, name2"]` | Export all schedules with exclusion filters to JSON, YAML, SQLite DB, or ZIP |
| `import` `<file> [-except "name1, name2"]` | Import single or batch schedules from JSON, YAML, SQLite DB, or ZIP into root and split DBs |
| `import-all` `-f <file> [-except "name1, name2"]` | Batch import all schedules with exclusion filter |
| `reset` `<name>` | Reset and clear execution run logs in the schedule's split database |
| `reset-all` | Reset and clear execution logs across all schedule split databases |
| `run` (`exec`) `<name>` | Execute a scheduled task immediately on demand and record execution history |
| `test` `<name> [--delay 1s] [--times <N>]` | Test execution of a scheduled task (e.g. running 1, 2, or 3 times now) |
| `rm` (`delete`, `del`) `<name>` | Remove a scheduled task from the root database and delete its split database file |
| `startup` `<name> [--enable\|--disable]` | Register/link task execution to native OS startup (Windows Registry / Linux Autostart) |
| `restart` | Trigger native OS restart (`shutdown /r` or `systemctl reboot`) |
| `shutdown` | Trigger native OS shutdown (`shutdown /s` or `systemctl poweroff`) |

---

## Flags

| Flag | Description |
|---|---|
| `-f <path>`, `--file <path>`, `-o <path>` | File path for export, import, or saving status/log outputs |
| `-except <list>`, `--except <list>` | Comma-separated list of schedule names to exclude during export-all or import-all |
| `--format <format>` | Explicit export format: `json`, `yaml`, `sqlite` (`db`), `zip` (inferred automatically from `-f` extension) |
| `--macro <name>`, `-m <name>` | Link scheduled execution to a saved macro |
| `--every <interval>`, `-i <interval>` | Schedule interval (e.g. `1d`, `2h`, `30m`, `15s`, `1h30m`) |
| `--day <N>`, `--days <N>` | Set daily schedule interval (e.g. `--day 1` for every 1 day) |
| `--hour <N>`, `--hours <N>` | Set hourly schedule interval (e.g. `--hour 2` for every 2 hours) |
| `--minute <N>`, `--min <N>` | Set minute schedule interval (e.g. `--min 30` for every 30 minutes) |
| `--second <N>`, `--sec <N>` | Set second schedule interval (e.g. `--sec 15` for every 15 seconds) |
| `--delay <duration>`, `--sleep <duration>` | Initial delay/sleep before task execution starts (e.g. `10s`, `1m`) |
| `--startup` | Automatically enable OS startup autostart on creation |
| `--limit <N>`, `-n <N>` | Number of log records to fetch for `schedule logs` (default: 20) |
| `--times <N>` | Number of test iterations for `schedule test` (e.g. `--times 3`) |
| `--json` | Output schedule list, status, or run logs as structured JSON |
| `--yaml` | Output schedule list, status, or run logs as structured YAML |

---

## Split Database Architecture

To prevent large output logs from bloating the main database, GitMap uses a **Split Database** model:
- **Root DB (`scheduler_tasks`)**: Lightweight index containing task name, slug, interval, enabled status, and last run timestamp.
- **Split DB (`data/schedules/<slug>.db`)**: Dedicated SQLite database per schedule storing complete metadata (`schedule_config`) and massive run histories (`schedule_logs`) with timestamps, durations, exit codes, usernames, and stdout/stderr outputs.

---

## Examples

```bash

# 1. Schedule a saved macro to run every 2 hours

gitmap schedule add backup-job --macro backup-db --every 2h

# 2. Schedule on-the-fly commands to run daily with a 10-second initial sleep

gitmap schedule add daily-clean "npm run clean && npm run build" --day 1 --delay 10s

# 3. View execution history and logs for a schedule

gitmap schedule log "daily-clean"
gitmap schedule log "daily-clean" --json -f "clean_log.json"
gitmap schedule log "daily-clean" --yaml -f "clean_log.yaml"

# 4. View single task status

gitmap schedule status "daily-clean"
gitmap schedule status "daily-clean" --json

# 5. Export single or all schedules (auto-infers format from file extension)

gitmap schedule export "daily-clean" -f "daily-clean.json"
gitmap schedule export-all -f "schedules_backup.yaml"
gitmap schedule export-all -f "schedules_archive.zip"
gitmap schedule export-all -f "schedules_backup.db"
gitmap schedule export-all --json -except "backup-job, test-task" -f "all_except.json"

# 6. Import schedules into root and split databases

gitmap schedule import "daily-clean.json"
gitmap schedule import-all -f "schedules_archive.zip" -except "deprecated-task"

# 7. Disable and re-enable a schedule

gitmap schedule disable daily-clean
gitmap schedule enable daily-clean

# 8. Reset execution logs for a schedule

gitmap schedule reset daily-clean

# 9. Test run a scheduled task immediately 3 times

gitmap schedule test daily-clean --times 3

# 10. Execute on demand

gitmap schedule run backup-job
```

---

## Web UI & REST API Integration

When running `gitmap ui` or `gitmap help-dashboard`, the built-in HTTP server provides real-time terminal and API execution endpoints:

- `http://localhost:5173/terminal` — Interactive web console for running GitMap commands.
- `POST /api/command/exec` — Execute any GitMap command via REST API:
  ```json
  // Request
  POST http://localhost:5173/api/command/exec
  Content-Type: application/json
  { "command": "gitmap schedule logs daily-clean --json" }

  // Response
  {
    "success": true,
    "output": "[{\"runNumber\":1,\"isSuccess\":true,\"durationMs\":120}]",
    "exitCode": 0
  }
  ```
