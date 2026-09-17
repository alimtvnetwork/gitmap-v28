# gitmap schedule restart

Schedule an operating system reboot with flexible duration syntax, state tracking, status inspection, and abort cancellation.

## Usage

```bash
gitmap schedule restart <duration>
gitmap schedule restart status
gitmap schedule restart cancel
```

## Description

`gitmap schedule restart` schedules an operating system reboot after a specified countdown duration. It coordinates native OS reboot commands (`shutdown /r /t <sec>` on Windows, `shutdown -r +<min>` on Linux) while maintaining persistent state in `~/.gitmap/power_schedule.json`.

### Flexible Duration Formats
Duration parameters support human-friendly formats:
- **Colon Hours & Minutes**: `1:45hr`, `1:45h`, `2:30hr`
- **Hours**: `2h`, `4h`, `12h`
- **Minutes**: `120m`, `45m`, `90m`
- **Days**: `1day`, `1d`, `2days`, `7d`
- **Seconds**: `30s`, `600s`
- **Immediate**: `now`, `0`

### Status Inspection & Cancellation
- **`gitmap schedule restart status`**: Displays active countdown timer, time elapsed, time remaining, and armed status. Also surfaced automatically when running `gitmap schedule ls` or `gitmap schedule status`.
- **`gitmap schedule restart cancel`**: Immediately halts the native OS countdown (`shutdown /a` on Windows, `shutdown -c` on Linux) and disarms the schedule state.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `<duration>` | Schedule system reboot after the specified duration (e.g. `1:45hr`, `2h`, `1day`) |
| `status` | View active reboot schedule, scheduled action, countdown remaining, and trigger time |
| `cancel` (`abort`) | Abort and cancel active operating system reboot countdown |

## Examples

```bash
# 1. Schedule restart in 1 hour and 45 minutes
gitmap schedule restart 1:45hr

# 2. Schedule restart using alternate duration formats
gitmap schedule restart 1:45h
gitmap schedule restart 2h
gitmap schedule restart 120m
gitmap schedule restart 1day
gitmap schedule restart 1d

# 3. Inspect active reboot countdown and remaining time
gitmap schedule restart status

# 4. Cancel and abort pending reboot countdown
gitmap schedule restart cancel

# 5. Remote restart across cluster nodes via SSH
gitmap cluster exec all "gitmap schedule restart 1:45hr"
gitmap sc bash "gitmap schedule restart 2h" --except control-plane

# 6. Immediate reboot
gitmap schedule restart now
```

See also: `gitmap schedule shutdown`, `gitmap schedule ls`, `gitmap schedule status`, `gitmap cluster exec`
