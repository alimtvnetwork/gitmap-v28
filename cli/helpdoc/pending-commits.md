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
| `--no-cache` | | `false` | Bypass SQLite status cache and query git status directly |
| `--refresh` | | `false` | Force refresh SQLite status cache with fresh repository metrics |
| `--ssh` | `-ssh` | `false` | Aggregate pending commits across all registered cluster SSH nodes |
| `--json` | `-j` | `false` | Output machine-readable JSON telemetry |
| `--dirty-only` | | `true` | Display only repositories with changes |
| `--all` | | `false` | Include completely clean repositories in output |
| `-h`, `--help` | | `false` | Display command help menu |

## Table Columns & Tree View Output

The output renders a clean table consolidating untracked, modified, and staged files:
- **REPOSITORY**: Repository slug or folder name.
- **BRANCH**: Active local branch (short version).
- **UNCOMMITTED**: Consolidated count of all uncommitted working files.
- **UNPUSHED**: Number of local commits ahead of remote tracking branch.
- **STATUS**: Clean or Dirty indicator.

For dirty repositories, a structured tree view provides instant surgical fix commands:
```text
├── Option 1: gitmap cpar "wip: save changes"
└── Option 2: gitmap stash
```

A global batch remediation command is rendered in the footer to resolve all dirty repositories simultaneously:
```text
Batch fix command for all dirty repositories:
  gitmap cpar "wip: batch backup all dirty repositories"
```

## Examples

### Example 1: Standard status overview (uses 90s SQLite cache)
```bash
gitmap pc
```

### Example 2: Force refresh cache with live git scan
```bash
gitmap pc --refresh
```

### Example 3: Completely bypass cache
```bash
gitmap pc --no-cache
```

### Example 4: Output machine-readable JSON
```bash
gitmap pc --json
```
