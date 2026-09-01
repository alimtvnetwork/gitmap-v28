# 19 — list-versions Command

## Purpose

`gitmap list-versions` (`lv`) lists all Git release tags (matching `v*`) in the current repository, sorted from highest to lowest semantic version.

## Command Signature

```
gitmap list-versions [flags]
gitmap lv [flags]
```

## Flags

| Flag       | Short | Default | Description                        |
|------------|-------|---------|------------------------------------|
| `--json`   |       | false   | Output as JSON array               |
| `--limit`  |       | 0       | Show only the top N versions (0 = all) |
| `--source` |       | (all)   | Filter by source: `release` or `import` |

## Behavior

1. Run `git tag --list "v*"` to collect all version tags.
2. Parse each tag with `release.Parse()`. Skip unparseable tags silently.
3. Sort descending by semantic version (highest first).
4. Cross-reference tags with the `Releases` database table to attach `source` metadata.
5. If `--source` is set, keep only versions matching that source value.
6. Print each version on its own line, v-prefixed (e.g. `v2.11.0`). Append `[source]` when available.
7. If `--json` is set, output a JSON array of version objects with `source` field.
8. If no tags are found (after filtering), print an error and exit 1.

## Terminal Output Example

```
v2.11.0
v2.10.0
v2.9.0
v2.8.0
```

## JSON Output Example

```json
["v2.11.0","v2.10.0","v2.9.0","v2.8.0"]
```

## Implementation Files

| File                       | Responsibility                              |
|----------------------------|---------------------------------------------|
| `cmd/listversions.go`      | Command handler, flag parsing, output       |
| `constants/constants_cli.go` | `CmdListVersions`, `CmdListVersionsAlias` |

## Code Style

All functions ≤ 15 lines. Positive logic. Blank line before every return. No magic strings. No switch statements.
