# `gitmap size`

Per-repo `.git` directory size report, sorted largest first. Optionally
run `git gc --aggressive --prune=now` on each listed repo.

## Flags

```
--root=DIR      scan root directory (default ".")
--top=N         show only the N largest repos (0 = all)
--prune         run `git gc --aggressive --prune=now` on each listed repo
--dry-run       with --prune: list gc invocations without running them
```

## Examples

```
gitmap size
gitmap size --top=10
gitmap size --top=5 --prune --dry-run
gitmap size --top=5 --prune
```
