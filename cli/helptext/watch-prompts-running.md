# gitmap watch-prompts-running

Monitor, persist, and auto-recover active Antigravity prompts and IDE processes across local and remote SSH machines.

## Synopsis

```bash
gitmap watch-prompts-running [subcommand] [flags]
gitmap wpr [subcommand] [flags]
gitmap agy wpr [subcommand] [flags]
```

## Aliases

- `watch-prompts-running`
- `wpr`
- `gitmap agy wpr`
- `gitmap agy watch-prompts-running`

## Description

`gitmap watch-prompts-running` (`wpr`) provides automated prompt persistence and process watchdog management for Antigravity (AGY) sessions. It automatically snapshots active and queued prompts to a dedicated Split-DB (`watch-prompts/<repo-slug>/sql.db`), tracks running states, handles prompt preservation during account switching, auto-recovers hung processes, and coordinates remote execution across cluster SSH fleets.

## Subcommands

- `all`: Snapshot all active and queued prompts to Split-DB and summarize.
- `ls`, `list`: Inspect watched projects, watch loop state, and intervals.
- `start [target]`: Start watch loop and auto-recovery for target project or all.
- `disable [target]`: Pause monitoring loop without deleting watched projects.
- `shutdown`, `sd`: Stop monitoring loop, terminate IDE process, and initiate shutdown.
- `restart [target]`: Restart IDE process and re-inject recent active prompts.
- `switch-account`, `sa [email]`: Switch account, backup prompts, and restart IDE.
- `fast-forward`, `ff [email]`: Fast-forward account switch with prompt preservation.
- `logs`, `log [target]`: View watcher event history and auto-recovery log entries.
- `status`, `st`: Inspect watchdog state and cached remote machine identities.
- `remove`, `rm <target|all>`: Remove specific project or all projects from watch list.
- `deploy <alias|all> <prj>`: Deploy Split-DB and enqueue WPR task via SSH/SCP stream.
- `help`: Show comprehensive command synopsis and usage guide.

## Flags & Options

| Flag | Default | Description |
|------|---------|-------------|
| `-t`, `--time <duration>` | `2m` | Polling and schedule check interval |
| `-p`, `--prefix <template>` | `default` | Prefix template prepended to re-injected prompts |
| `-s`, `--suffix <template>` | | Suffix template appended to re-injected prompts |
| `--ssh` | `false` | Execute across cluster SSH fleet nodes |
| `-j`, `--json` | `false` | Emit machine-readable JSON output |
| `-n`, `--dry-run` | `false` | Simulate operations without modifying system state |
| `-h`, `--help` | `false` | Display command usage and synopsis |

## Examples

### Example 1: View active prompts and monitoring status

```bash
gitmap wpr ls
```

### Example 2: Snapshot all active Antigravity prompts to Split-DB

```bash
gitmap wpr all
```

### Example 3: Start watchdog loop on all active projects

```bash
gitmap wpr start all --time 5m
```
