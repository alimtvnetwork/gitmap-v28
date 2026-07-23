# gitmap rescan

Re-scan previously scanned directories using cached scan parameters.

## Alias

rs

## Usage

    gitmap rescan [--output csv|json|terminal]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --output csv\|json\|terminal | terminal | Output format |

## Prerequisites

- Run `gitmap scan` at least once to create scan cache (see scan.md)

## Examples

### Example 1: Quick rescan (picks up new repos)

    gitmap rescan

**Output:**

    Re-scanning D:\wp-work (cached parameters)...
    [1/44] github/user/my-api
    [2/44] github/user/web-app
    ...
    Found 44 repositories (+2 new since last scan)
    ✓ Database updated (44 repos)
    ✓ Output written to ./.gitmap/output/

### Example 2: Rescan with JSON output

    gitmap rs --output json

**Output:**

    Re-scanning D:\wp-work (cached parameters)...
    Found 44 repositories
    ✓ .gitmap/output/gitmap.json written
    ✓ .gitmap/output/gitmap.csv written

### Example 3: Rescan with no cache (error)

    gitmap rescan

**Output:**

    ✗ No scan cache found.
    → Run 'gitmap scan <directory>' first to create a scan cache.

## See Also

- [scan](scan.md) — Initial directory scan
- [status](status.md) — View repo statuses
- [clone](clone.md) — Clone from scan output

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter rescan
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
