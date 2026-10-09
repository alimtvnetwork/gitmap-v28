# gitmap pending-commits

Scans local workspace repositories or fleet cluster nodes over SSH to identify uncommitted working tree modifications and unpushed commits ahead of origin.

## Alias

pc

## Usage

```bash
gitmap pending-commits [repoName|all] [flags]
gitmap pc [repoName|all] [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--sort <mode>` | `-s` | `priority` | Sort order: `priority`, `name`, or `count` |
| `--detail <mode>` | `-d` | `summary` | Detail view: `summary`, `all`, `1`, or `none` |
| `--ssh` | `-ssh` | `false` | Aggregate pending commits across all registered cluster SSH nodes |
| `--json` | `-j` | `false` | Output machine-readable JSON telemetry |
| `--dirty-only` | | `true` | Display only repositories with changes |
| `--all` | | `false` | Include completely clean repositories in output |
| `-h`, `--help` | | `false` | Display command help menu |

## Examples

### Example 1: Summary overview of dirty repositories

```bash
gitmap pc
```

### Example 2: Detailed file breakdown sorted by modification count

```bash
gitmap pc --sort=count --detail=all
```

### Example 3: Inspect pending commits across remote SSH cluster nodes

```bash
gitmap pc --ssh --json
```
