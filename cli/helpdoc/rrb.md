# rrb

Alias for `gitmap release remove-between`.

## Synopsis

```
gitmap rrb <start> <end> [flags]
```

## Description

Prune intermediate release tags between two bounds, strictly preserving both endpoint milestones and pinned tags.

See full documentation at `gitmap release-remove-between --help`.

## Examples

```bash
gitmap rrb 4.162, 4.190
gitmap rrb 4.162, 4.190 --confirm
gitmap rrb 4.162 4.190 -y
gitmap rrb 4.162, 4.190 --dry-run
gitmap rrb 4.162, 4.190 --pin 4.175.0 --confirm
```
