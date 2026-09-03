# `gitmap db reset`

Reset database records and remove split search databases while preserving core configuration.

## Usage

```bash
gitmap db reset [flags]
gitmap db clear [flags]
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--all` | `false` | Resets master repo records and deletes all split DBs |
| `--dry-run` | `false` | Preview files to be cleared without deleting |
| `--yes`, `-y` | `false` | Skip interactive confirmation prompt |

## Examples

```bash
gitmap db reset
gitmap db reset --yes
```
