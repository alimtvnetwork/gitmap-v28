# gitmap schedule shutdown

Schedule an operating system shutdown with flexible duration syntax, state tracking, status inspection, and abort cancellation.

## Usage

```bash
gitmap schedule shutdown <duration>
gitmap schedule shutdown status
gitmap schedule shutdown cancel
```

## Description

`gitmap schedule shutdown` schedules an operating system power-off after a specified countdown duration. It automatically registers and coordinates native OS shutdown commands (`shutdown /s /t <sec>` on Windows, `shutdown -h +<min>` on Linux) while maintaining persistent state in `~/.gitmap/power_schedule.json`.

### Flexible Duration Formats
Duration parameters support human-friendly formats:
- **Colon Hours & Minutes**: `1:45hr`, `1:45h`, `2:30hr`
- **Hours**: `2h`, `4h`, `12h`
- **Minutes**: `120m`, `45m`, `90m`
- **Days**: `1day`, `1d`, `2days`, `7d`
- **Seconds**: `30s`, `600s`
- **Immediate**: `now`, `0`

### Status Inspection & Cancellation
- **`gitmap schedule shutdown status`**: Displays active countdown timer, time elapsed, time remaining, and armed status. Also surfaced automatically when running `gitmap schedule ls` or `gitmap schedule status`.
- **`gitmap schedule shutdown cancel`**: Immediately halts the native OS countdown (`shutdown /a` on Windows, `shutdown -c` on Linux) and disarms the schedule state.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `<duration>` | Schedule system shutdown after the specified duration (e.g. `1:45hr`, `2h`, `1day`) |
| `status` | View active power schedule, scheduled action, countdown remaining, and trigger time |
| `cancel` (`abort`) | Abort and cancel active operating system shutdown countdown |

## Examples

```bash
# 1. Schedule shutdown in 1 hour and 45 minutes
gitmap schedule shutdown 1:45hr

# 2. Schedule shutdown using alternate duration formats
gitmap schedule shutdown 1:45h
gitmap schedule shutdown 2h
gitmap schedule shutdown 120m
gitmap schedule shutdown 1day
gitmap schedule shutdown 1d

# 3. Inspect active shutdown countdown and remaining time
gitmap schedule shutdown status

# 4. Cancel and abort pending shutdown countdown
gitmap schedule shutdown cancel

# 5. Remote shutdown across cluster nodes via SSH
gitmap cluster exec all "gitmap schedule shutdown 1:45hr"
gitmap sc bash "gitmap schedule shutdown 2h" --except control-plane

# 6. Immediate shutdown
gitmap schedule shutdown now
```

See also: `gitmap schedule restart`, `gitmap schedule ls`, `gitmap schedule status`, `gitmap cluster exec`
