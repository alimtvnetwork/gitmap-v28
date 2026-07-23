# gitmap changelog-generate

Auto-generate changelog entries from commit messages between tags.

## Alias

cg

## Usage

    gitmap changelog-generate [--from <tag>] [--to <tag>] [--write]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --from | second-latest tag | Start tag (older boundary) |
| --to | latest tag | End tag or HEAD |
| --write | false | Prepend output to CHANGELOG.md |

## Prerequisites

- Must be inside a Git repository with at least one version tag

## Examples

### Example 1: Generate changelog between latest two tags

    gitmap changelog-generate

**Output:**

    Changelog: v2.23.0 → v2.24.0

    Preview (use --write to save):

    ## v2.24.0

    - Add TUI log viewer with detail panel
    - Add release rollback on push failure
    - Fix watch interval validation edge case

### Example 2: Generate from a specific tag to HEAD

    gitmap cg --from v2.22.0

**Output:**

    Changelog: v2.22.0 → HEAD

    Preview (use --write to save):

    ## Unreleased

    - Add changelog-generate command
    - Add TUI log viewer with detail panel
    - Fix zip-group archive naming

### Example 3: Write directly to CHANGELOG.md

    gitmap cg --from v2.23.0 --to v2.24.0 --write

**Output:**

    Changelog: v2.23.0 → v2.24.0

    ✓ Prepended changelog to CHANGELOG.md

## See Also

- [changelog](changelog.md) — View existing changelog entries
- [release](release.md) — Create a release
- [list-versions](list-versions.md) — List release tags

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter changelog-generate
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
