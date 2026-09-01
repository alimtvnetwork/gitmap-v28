# schedule

Manage scheduled tasks, automated macro executions, interval runners, OS startup triggers, test runners, and Web UI/API integrations.

## Aliases

`sched`

---

## Subcommands

| Subcommand | Description |
|---|---|
| `add` (`create`, `new`) `<name> [cmds...]` | Create a new scheduled task (linked to a macro or on-the-fly shell commands) |
| `list` (`ls`, `status`) `[--json\|--yaml]` | List all scheduled tasks with interval, target, runs count, startup status, and last execution |
| `run` (`exec`) `<name>` | Execute a scheduled task immediately on demand |
| `test` `<name> [--delay 1s] [--times <N>]` | Test execution of a scheduled task (e.g. running 1, 2, or 3 times now) |
| `rm` (`delete`, `del`) `<name>` | Remove a scheduled task from the SQLite database |
| `startup` `<name> [--enable\|--disable]` | Register/link task execution to native OS startup (Windows Registry / Linux Autostart) |
| `restart` | Trigger native OS restart (`shutdown /r` or `systemctl reboot`) |
| `shutdown` | Trigger native OS shutdown (`shutdown /s` or `systemctl poweroff`) |

---

## Flags

| Flag | Description |
|---|---|
| `--macro <name>`, `-m <name>` | Link scheduled execution to a saved macro |
| `--every <interval>`, `-i <interval>` | Schedule interval (e.g. `1d`, `2h`, `30m`, `15s`, `1h30m`) |
| `--day <N>`, `--days <N>` | Set daily schedule interval (e.g. `--day 1` for every 1 day) |
| `--hour <N>`, `--hours <N>` | Set hourly schedule interval (e.g. `--hour 2` for every 2 hours) |
| `--minute <N>`, `--min <N>` | Set minute schedule interval (e.g. `--min 30` for every 30 minutes) |
| `--second <N>`, `--sec <N>` | Set second schedule interval (e.g. `--sec 15` for every 15 seconds) |
| `--delay <duration>`, `--sleep <duration>` | Initial delay/sleep before task execution starts (e.g. `10s`, `1m`) |
| `--startup` | Automatically enable OS startup autostart on creation |
| `--times <N>`, `-n <N>` | Number of test iterations for `schedule test` (e.g. `--times 3`) |
| `--json` | Output schedule list/status as structured JSON |
| `--yaml` | Output schedule list/status as structured YAML |

---

## Examples

```bash
# 1. Schedule a saved macro to run every 2 hours
gitmap schedule add backup-job --macro backup-db --every 2h

# 2. Schedule on-the-fly commands to run daily with a 10-second initial sleep
gitmap schedule add daily-clean "npm run clean && npm run build" --day 1 --delay 10s

# 3. Schedule a task based on minutes and link to OS startup
gitmap schedule add sync-repos "gitmap sync --all" --minute 30 --startup

# 4. List all schedules (Terminal table, JSON, YAML)
gitmap schedule list
gitmap schedule list --json
gitmap schedule list --yaml

# 5. Test run a scheduled task immediately 3 times in a row
gitmap schedule test daily-clean --times 3

# 6. Execute a scheduled task now
gitmap schedule run backup-job

# 7. Remove a scheduled task
gitmap schedule rm sync-repos
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
  { "command": "gitmap schedule list --json" }

  // Response
  {
    "success": true,
    "output": "[{\"name\":\"backup-job\",\"interval\":\"2h\"}]",
    "exitCode": 0
  }
  ```
- `GET /api/terminal/stream?session=default` — Server-Sent Events (SSE) live stream of terminal output.
- `POST /api/terminal/input?session=default` — Send keystrokes / stdin commands.
- `GET /api/terminal/autocomplete?q=gitmap+sch` — Live command autocompletion suggestions.
