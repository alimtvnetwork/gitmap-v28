# startup

Manage and run OS startup items, shell scripts, PowerShell tasks, binaries, and macros across Windows, Linux, and macOS.

## Aliases

`su`

## Usage

    gitmap startup <subcommand> [flags]

## Subcommands

| Subcommand | Description |
|---|---|
| `ls`, `list` | Enumerate registered startup items, targets, run frequencies, and active states |
| `add <target>` | Register a script (`.ps1`, `.sh`), executable, icon launcher, or macro to OS startup |
| `rm <name>` | Unregister and remove a startup item from the system and split database |
| `run <name>` | Execute a startup item on demand and record execution logs |
| `logs <name>` | Display execution logs and run history from `startup.db` |

## Options

    --name <name>            Custom name identifier for the startup item
    --frequency <freq>       Run frequency: `everytime` (default, every OS login) or `once-a-week`
    --icon <path>            Path to custom icon file (`.ico`, `.png`) for desktop launchers
    --type <script|macro>    Explicit target type (`script`, `macro`, `binary`)

## Description

The `startup` command manages background and user-login autostart workflows using GitMap's **Split Database** architecture. Records and run logs are persisted in `<BinaryDataDir>/startup.db` with native OS hooks:
- **Windows**: User Registry Run keys and Startup folder shortcuts.
- **Linux**: XDG Autostart `.desktop` entries in `~/.config/autostart/`.
- **macOS**: LaunchAgents in `~/Library/LaunchAgents/`.

## Examples

```bash
# List all configured startup items
gitmap startup ls

# Add a PowerShell script to run every OS restart
gitmap startup add scripts/sync-env.ps1 --frequency everytime

# Add a Bash script to run once a week with custom name
gitmap startup add scripts/cleanup.sh --name weekly-clean --frequency once-a-week

# Register a saved macro to launch on startup
gitmap startup add backup-repos --type macro

# Test run a registered startup item immediately
gitmap startup run weekly-clean

# View execution history and logs
gitmap startup logs weekly-clean

# Remove a startup item
gitmap startup rm weekly-clean
```
