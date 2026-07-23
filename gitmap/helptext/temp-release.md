# Temp Release

Create lightweight, temporary release branches from recent commits for quick experimentation. No tags, no metadata — just branches.

## Usage

```
gitmap temp-release <count> <version-pattern> [-s <start>]
gitmap tr <count> <version-pattern> [-s <start>]
```

## Subcommands

- `list` — Show all temp-release branches
- `remove <version>` — Remove a single temp-release branch
- `remove <v1> to <v2>` — Remove a range of temp-release branches
- `remove all` — Remove all temp-release branches

## Flags

| Flag | Description |
|------|-------------|
| `-s`, `--start` | Starting sequence number (default: auto-increment) |
| `--dry-run` | Preview branch names without creating |
| `--json` | JSON output for `list` subcommand |

## Version Pattern

Use `$` placeholders for zero-padded sequence numbers:
- `$$` → 2 digits (05, 12, 99)
- `$$$` → 3 digits (005, 012, 099)
- `$$$$` → 4 digits (0005, 0012, 0099)

## Examples

```bash
# Create 10 branches from last 10 commits, starting at sequence 5
gitmap tr 10 v1.$$ -s 5

# Create 1 branch, auto-increment from last temp-release
gitmap tr 1 v1.$$

# Preview without creating
gitmap tr 5 v2.$$$ --dry-run

# List all temp-release branches
gitmap tr list

# Remove a single branch
gitmap tr remove v1.05

# Remove a range
gitmap tr remove v1.05 to v1.10

# Remove all
gitmap tr remove all
```

## See Also

- [release](release.md) — Full release with tags and metadata
- [release-branch](release-branch.md) — Complete release from existing branch
- [prune](prune.md) — Clean up stale branches

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter temp-release
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
