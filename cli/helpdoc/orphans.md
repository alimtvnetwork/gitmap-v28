# `gitmap orphans`

Find local clones whose `origin` remote returns HTTP 404 or 410 and
offer bulk delete.

## Flags

```
--root=DIR      scan root directory (default ".")
-y              delete without prompting
--dry-run       list only; never delete
```

## Examples

```
gitmap orphans
gitmap orphans --dry-run
gitmap orphans -y
```
