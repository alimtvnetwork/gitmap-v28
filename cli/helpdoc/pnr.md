# pnr

Alias for `gitmap release pin-not-to-remove`.

## Synopsis

```
gitmap pnr [tags...] [flags]
```

## Description

Persists tag names in `.gitmap/release/pinned-tags.json`. All pinned tags are unconditionally preserved by `gitmap release remove-between` (`rrb`) and cleanup commands.

## Examples

```bash
gitmap pnr 4.162, 4.190
gitmap pnr --list
gitmap pnr --remove 4.162
gitmap pnr --clear
```
