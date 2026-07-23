# `gitmap stale` (alias: `sta`)

List local repositories with no commits in the last N days. Optionally
move stale repos to `.gitmap/archive/<UTC-timestamp>/`.

## Flags

```
--days=N        report repos with no commits in the last N days (default 90)
--root=DIR      scan root directory (default ".")
--archive       move stale repos into .gitmap/archive/<ts>/
--dry-run       with --archive: preview moves without touching disk
```

## Examples

```
gitmap stale
gitmap sta --days=180
gitmap stale --days=365 --archive --dry-run
gitmap stale --days=365 --archive
```
