# gitmap dashboard

Generate an interactive HTML dashboard for a repository.

## Alias

db

## Usage

    gitmap dashboard [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --limit | 0 (all) | Maximum number of commits to include |
| --since | (none) | Only include commits after this date (YYYY-MM-DD) |
| --no-merges | false | Exclude merge commits from the output |
| --out-dir | .gitmap/output | Output directory for dashboard files |
| --open | false | Open the generated dashboard in the default browser |

## Prerequisites

- Must be run inside a Git repository with commit history

## Examples

### Example 1: Generate a full dashboard

    gitmap dashboard

**Output:**

    Collecting repository data...
    Wrote .gitmap/output/dashboard.json (482 commits, 7 authors)
    Wrote .gitmap/output/dashboard.html
    Dashboard generated in .gitmap/output

### Example 2: Last 100 commits, open in browser

    gitmap db --limit 100 --open

**Output:**

    Collecting repository data...
    Wrote .gitmap/output/dashboard.json (100 commits, 5 authors)
    Wrote .gitmap/output/dashboard.html
    Dashboard generated in .gitmap/output
    Opening dashboard in browser...

### Example 3: Commits since a date, no merges

    gitmap dashboard --since 2025-01-01 --no-merges --out-dir ./report

**Output:**

    Collecting repository data...
    Wrote ./report/dashboard.json (63 commits, 4 authors)
    Wrote ./report/dashboard.html
    Dashboard generated in ./report

## See Also

- [stats](stats.md) — Show commit statistics for scanned repositories
- [history](history.md) — View command execution history
- [changelog](changelog.md) — Generate a changelog from Git tags

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter dashboard
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
