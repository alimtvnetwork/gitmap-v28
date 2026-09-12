# gitmap power

Cross-platform OS screen timeout and sleep management framework. Inspect, configure, or reset host display and sleep timeouts with SQLite state tracking and restoration profiles.

## Usage

    gitmap power [subcommand] [flags]
    gitmap pw [subcommand]
    gitmap pwr [subcommand]

## Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| status | st | Show current OS and SQLite power configuration |
| never-sleep | never, ns | Set display and sleep timeouts to 0 (never) |
| set <mins> | - | Set custom display and sleep timeouts in minutes |
| reset | restore | Restore previous power timeout profile from SQLite |
| history | hist, logs | View audit log of power configuration changes |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --display <mins> | - | Set display screen timeout in minutes |
| --sleep <mins> | - | Set sleep/standby timeout in minutes |

## Examples

### View current status

```bash
gitmap power status
```

Output:

```text
▶ OS Power & Sleep Status (windows)
  • Display Timeout: Never
  • Sleep Timeout:   Never
  • Mode:            Never-Sleep (inhibited)
```

### Prevent screen timeout and system sleep

```bash
gitmap power never-sleep
```

Output:

```text
✓ Power setting updated: Never-Sleep mode enabled (display & standby timeouts disabled).
```

### Set custom timeouts (15m display, 30m sleep)

```bash
gitmap power set --display 15 --sleep 30
```

Output:

```text
✓ Power setting updated: Display timeout=15m, Sleep timeout=30m.
```

### Restore previous timeouts

```bash
gitmap power reset
```

Output:

```text
✓ Power settings reset: Display=15 minutes, Sleep=30 minutes.
```

### View change history

```bash
gitmap power history
```

Output:

```text
ID     ACTION       DISPLAY    SLEEP      TIMESTAMP            NOTES
--------------------------------------------------------------------------------
1      never-sleep  Never      Never      2026-09-07T12:00:00Z Configured never-sleep
```

## Supported Platforms

- **Windows**: Uses `powercfg.exe` (`/change monitor-timeout-ac`, `standby-timeout-ac`).
- **Linux/Ubuntu**: Supports GNOME `gsettings`, X11 `xset`, and headless `systemd`.
- **macOS / Future**: Extensible pluggable driver architecture.
