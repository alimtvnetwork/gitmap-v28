# prune

Delete stale release branches that have already been tagged.

**Alias:** `pr`

## Usage

```
gitmap prune [flags]
```

## Flags

| Flag        | Description                              |
|-------------|------------------------------------------|
| `--dry-run` | List stale branches without deleting     |
| `--confirm` | Skip interactive confirmation prompt     |
| `--remote`  | Also delete remote release branches      |

## Prerequisites

- Must be inside a Git repository.
- Release branches follow the `release/vX.Y.Z` naming convention.

## Examples

### List stale branches (dry run)

```
$ gitmap prune --dry-run

  Stale release branches (3):
    release/v2.20.0  →  tag v2.20.0 exists
    release/v2.21.0  →  tag v2.21.0 exists
    release/v2.22.0  →  tag v2.22.0 exists

  Use --confirm to delete, or run without --dry-run for interactive mode.
```

### Delete with confirmation

```
$ gitmap prune --confirm

  Pruning stale release branches...
    ✓ Deleted release/v2.20.0
    ✓ Deleted release/v2.21.0
    ✓ Deleted release/v2.22.0

  Summary: 3 deleted, 2 kept.
```

### Delete including remote branches

```
$ gitmap prune --confirm --remote

  Pruning stale release branches...
    ✓ Deleted release/v2.20.0
    ✓ Deleted remote release/v2.20.0
    ✓ Deleted release/v2.21.0
    ✓ Deleted remote release/v2.21.0

  Summary: 2 deleted, 0 kept.
```

## See Also

- `release` — Create release branches and tags
- `clear-release-json` — Remove release metadata files
- `list-releases` — Show stored releases from database

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter prune
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
