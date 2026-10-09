# gitmap shutdown-until

Monitor designated projects and trigger automated operating system shutdown when all CI/CD pipelines turn green.

## Synopsis

```bash
gitmap shutdown-until <command> [flags]
gitmap sug <command> [flags]
gitmap agy shutdown-until <command> [flags]
gitmap agy sug <command> [flags]
```

## Aliases

- `shutdown-until`
- `shutdown-until-green`
- `sug`
- `gitmap agy shutdown-until`
- `gitmap agy sug`

## Description

`gitmap shutdown-until` (concise alias `sug`) automates overnight agent pipeline workflows. It tracks designated projects and periodically evaluates their CI/CD status using cached pipeline error data and remote pipeline status checks. When every registered project's CI/CD pipeline turns green with no failures or active runs, GitMap gracefully triggers the operating system's shutdown command.

To test safely without shutting down the system, use the `--dry-run` (`-n`) flag.

## Subcommands

- `ls` / `list`: Display the registered project targets in the shutdown watch list.
- `add-projects <targets...>` / `add`: Register one or more project names or directory paths into the watch list.
- `rm <targets...>` / `remove`: Remove designated projects from the watch list.
- `agy-running-projects` / `running-projects`: Automatically discover all currently active Antigravity workspaces (with queued prompts or active fix prompts) and set them as the watch list.
- `run [flags]`: Start the continuous evaluation loop. When all pipelines turn green, executes OS shutdown.
- `help`: Render the interactive two-column terminal help menu.

## Flags & Options

| Flag | Default | Description |
|------|---------|-------------|
| `-t`, `--time <duration>` | `5m` | Polling check interval (minimum `2m`, or `10s` during `--dry-run`) |
| `-n`, `--dry-run` | `false` | Simulate OS shutdown execution without powering off system |
| `-1`, `--once` | `false` | Execute a single status evaluation cycle and exit immediately |
| `-h`, `--help` | `false` | Display comprehensive two-column terminal help menu |

## Examples

### Example 1: View shutdown-until interactive help

```bash
gitmap sug help
```

### Example 2: Automatically watch all running Antigravity projects

```bash
gitmap sug agy-running-projects
```

### Example 3: Inspect currently monitored projects

```bash
gitmap sug ls
```

### Example 4: Add specific project targets

```bash
gitmap sug add-projects gitmap alim-status-sample
```

### Example 5: Test the watch loop safely with dry-run

```bash
gitmap sug run --dry-run -t 1m
```

### Example 6: Single-pass dry-run check

```bash
gitmap shutdown-until run -n -1
```

### Example 7: Launch live overnight shutdown monitor

```bash
gitmap sug run -t 5m
```
