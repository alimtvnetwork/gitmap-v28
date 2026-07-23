# gitmap self-update

Probe the GitHub release feed for the newest tag and, when a newer
version is available, re-run `gitmap self-install` non-interactively.

## Usage

```
gitmap self-update [--dry-run] [--force]
```

## Flags

| Flag        | Effect                                                                |
|-------------|-----------------------------------------------------------------------|
| `--dry-run` | Print the would-be install command and exit 0 without touching disk.  |
| `--force`   | Re-install even when the local version already matches latest.        |

## Exit codes

| Code | Meaning                                            |
|------|----------------------------------------------------|
| 0    | Already on latest, or successful update.           |
| 1    | self-install failed (network, permissions, etc.).  |
| 2    | Could not reach the release API.                   |

## Examples

```
gitmap self-update
gitmap self-update --dry-run
gitmap self-update --force
```
