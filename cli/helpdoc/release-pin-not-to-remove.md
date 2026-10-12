# release-pin-not-to-remove (pnr)

Pin milestone release tags so they are permanently protected from removal and auto-pruning.

## Synopsis

```
gitmap release pin-not-to-remove [tags...] [flags]
gitmap pnr                       [tags...] [flags]
```

## Description

Persists tag names in `.gitmap/release/pinned-tags.json`. All pinned tags are unconditionally preserved by `gitmap release remove-between` (`rrb`) and cleanup commands.

## Flags

- `--list, -l`       — Print all currently pinned release tags
- `--remove <tag>`   — Remove a tag from the pin registry
- `--unpin <tag>`    — Alias for `--remove`
- `--clear`          — Remove all pins from the registry

## Examples

```bash
# Pin milestone releases so they are never deleted
gitmap pnr 4.162, 4.190

# Pin using SemVer notation
gitmap pnr v4.162.0 v4.190.0

# List active pins
gitmap pnr --list

# Unpin a tag
gitmap pnr --remove 4.162

# Clear all pins
gitmap pnr --clear
```

## See Also

- `gitmap rrb` — Remove intermediate tags between two bounds
