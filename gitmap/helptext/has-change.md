# gitmap has-change

Print `true` or `false` indicating whether a tracked repo has uncommitted changes, unpushed commits, or unpulled commits.

## Alias

hc

## Usage

    gitmap has-change <repo> [--mode dirty|ahead|behind] [--all] [--fetch=false]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --mode \<m\> | dirty | Dimension to check: `dirty`, `ahead`, or `behind` |
| --all | false | Print all three dimensions as `dirty=X ahead=Y behind=Z` |
| --fetch | true | Run `git fetch` before checking ahead/behind |

## Examples

### Example 1: Check if working tree is dirty

    gitmap hc gitmap

**Output:**

    true

### Example 2: Check if local has unpushed commits

    gitmap hc gitmap --mode ahead

**Output:**

    false

### Example 3: Get all three dimensions at once

    gitmap hc gitmap --all

**Output:**

    dirty=true ahead=false behind=true

## See Also

- [status](status.md) — Full per-repo status dashboard
- [has-any-updates](has-any-updates.md) — Check current repo for remote updates

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter has-change
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
