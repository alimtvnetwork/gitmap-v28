# gitmap list-versions

List all available Git release tags in the repository.

## Alias

lv

## Usage

    gitmap list-versions [--json] [--limit N]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --json | false | Output as structured JSON |
| --limit \<N\> | 0 | Show only the top N versions (0 = all) |

## Prerequisites

- Must be inside a Git repository with tags

## Examples

### Example 1: List all versions

    gitmap list-versions

**Output:**

    VERSION     DATE          COMMITS  LATEST
    v2.22.0     2025-03-10    5        ✓
    v2.21.0     2025-03-08    3
    v2.20.0     2025-02-28    4
    v2.19.0     2025-02-20    6
    v2.18.0     2025-02-15    2
    ...
    12 versions found

### Example 2: Show top 3 versions

    gitmap lv --limit 3

**Output:**

    VERSION     DATE          COMMITS  LATEST
    v2.22.0     2025-03-10    5        ✓
    v2.21.0     2025-03-08    3
    v2.20.0     2025-02-28    4
    3 versions shown (12 total)

### Example 3: JSON output for scripting

    gitmap lv --json --limit 3

**Output:**

    [
      {"version":"v2.22.0","date":"2025-03-10","commits":5,"latest":true},
      {"version":"v2.21.0","date":"2025-03-08","commits":3,"latest":false},
      {"version":"v2.20.0","date":"2025-02-28","commits":4,"latest":false}
    ]

## See Also

- [list-releases](list-releases.md) — List stored release metadata
- [changelog](changelog.md) — View release notes
- [release](release.md) — Create a release
- [revert](revert.md) — Revert to a specific version

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter list-versions
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
