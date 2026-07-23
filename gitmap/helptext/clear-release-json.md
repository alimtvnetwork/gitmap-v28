# gitmap clear-release-json

Remove a specific release metadata JSON file from the `.gitmap/release/` directory.

## Alias

crj

## Usage

    gitmap clear-release-json <version> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --dry-run | false | Preview which file would be removed without deleting it |

## Prerequisites

- A `.gitmap/release/vX.Y.Z.json` file must exist for the given version
- Run `gitmap release` first to generate metadata (see [release](release.md))

## Examples

### Example 1: Remove a release JSON file

    gitmap clear-release-json v2.20.0

**Output:**

    Found .gitmap/release/v2.20.0.json (1.2 KB)
    ✓ Removed .gitmap/release/v2.20.0.json

### Example 2: Dry-run preview

    gitmap clear-release-json v2.20.0 --dry-run

**Output:**

    [dry-run] Found .gitmap/release/v2.20.0.json (1.2 KB)
    [dry-run] Would remove .gitmap/release/v2.20.0.json
    No changes made.

### Example 3: Version not found

    gitmap clear-release-json v9.9.9

**Output:**

    ✗ Error: no release file found for v9.9.9
    Available versions in .gitmap/release/:
      v2.22.0, v2.21.0, v2.20.0, v2.19.0
    → Use 'gitmap list-releases' to see all stored releases

### Example 4: Version with zero-padding normalization

    gitmap crj v2

**Output:**

    Found .gitmap/release/v2.0.0.json (0.8 KB)
    ✓ Removed .gitmap/release/v2.0.0.json

### Example 5: Clean up after orphaned metadata prompt

    gitmap release --bump patch
    # ⚠ Release metadata exists for v2.20.0 but no tag found
    # User decides to clean up manually:
    gitmap crj v2.20.0

**Output:**

    Found .gitmap/release/v2.20.0.json (1.2 KB)
    ✓ Removed .gitmap/release/v2.20.0.json
    → You can now re-run 'gitmap release --bump patch'

## See Also

- [release](release.md) — Create a new versioned release
- [list-releases](list-releases.md) — Show all stored releases
- [list-versions](list-versions.md) — List version tags in a repository
- [revert](revert.md) — Roll back a release
- [db-reset](db-reset.md) — Clear the entire database

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter clear-release-json
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
