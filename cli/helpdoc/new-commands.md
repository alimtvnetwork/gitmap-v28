# gitmap new-commands

Discovers, catalogs, and filters new CLI commands and subcommands added across the last 100 commits in GitMap's git history, complete with flags and usage examples.

## Alias

nc

## Usage

```bash
gitmap new-commands [flags]
gitmap nc [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit <n>` | `-n` | `100` | Maximum number of commands to inspect from git history |
| `--filter <str>` | `-f`, `-q` | `""` | Filter commands by keyword or prefix |
| `--category <cat>` | `-c` | `""` | Filter by category (e.g. `commits`, `diagnostics`, `scanner`, `fleet`, `ai`, `os`, `spec`, `storage`, `sync`, `tooling`) |
| `--json` | `-j` | `false` | Output machine-readable JSON array of discovered commands |
| `-h`, `--help` | | `false` | Display command help menu |

## Description

The `new-commands` command analyzes commit messages, CLI constant registrations, and command router changes across recent history. It produces a formatted catalog showing:
- Command name and aliases.
- Concise summary of purpose.
- Available flags.
- Concrete, runnable examples.

This command empowers both human engineers and AI coding agents to discover recently introduced capabilities without parsing git logs manually.

## Examples

### Example 1: View recently added commands (up to 100)
```bash
gitmap nc
```

### Example 2: Filter recent commands by search term (-f / -q)
```bash
gitmap nc --filter "cache"
gitmap nc -f "cpar"
gitmap nc -q "diagnostics"
```

### Example 3: Filter by category and limit count
```bash
gitmap nc --category "commits" --limit 10
```

### Example 4: Limit output to top 10 new commands
```bash
gitmap nc --limit 10
```

### Example 5: Machine-readable JSON output for automated agents
```bash
gitmap nc --json
```
