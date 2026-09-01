# Prune — Stale Release Branch Cleanup

## Purpose

The `prune` command deletes local release branches that have already been
merged or tagged. This prevents branch clutter after releases accumulate.

## Command

```
gitmap-v28 prune [flags]
```

| Alias | Short |
|-------|-------|
| `pr`  | —     |

## Behavior

1. List all local branches matching `release/*`.
2. For each branch, check whether a corresponding tag (`v*`) exists.
3. A branch is **stale** if the tag exists (both local and remote).
4. In `--dry-run` mode, print stale branches without deleting.
5. Without `--confirm`, prompt the user before deleting.
6. Delete stale branches with `git branch -D <name>`.
7. Print a summary: N deleted, M kept.

## Flags

| Flag        | Short | Default | Description                        |
|-------------|-------|---------|------------------------------------|
| `--dry-run` | —     | `false` | List stale branches without deleting |
| `--confirm` | —     | `false` | Skip interactive confirmation prompt |
| `--remote`  | —     | `false` | Also delete remote release branches  |

## Output

### Dry-run

```
  Stale release branches (3):
    release/v2.20.0  →  tag v2.20.0 exists
    release/v2.21.0  →  tag v2.21.0 exists
    release/v2.22.0  →  tag v2.22.0 exists

  Use --confirm to delete, or run without --dry-run for interactive mode.
```

### Deletion

```
  Pruning stale release branches...
    ✓ Deleted release/v2.20.0
    ✓ Deleted release/v2.21.0
    ✓ Deleted release/v2.22.0

  Summary: 3 deleted, 2 kept.
```

### No stale branches

```
  No stale release branches found.
```

## Remote Cleanup

When `--remote` is set, after deleting the local branch:

```
git push origin --delete release/v2.20.0
```

Failures are logged as warnings but do not stop remaining deletions.

## Implementation Files

| File                            | Action | Purpose                           |
|---------------------------------|--------|-----------------------------------|
| `constants/constants_prune.go`  | CREATE | Command names, messages, errors   |
| `cmd/prune.go`                  | CREATE | Command handler + flag parsing    |
| `cmd/pruneops.go`               | CREATE | Branch listing, staleness check   |
| `helptext/prune.md`             | CREATE | Embedded help documentation       |
| `cmd/root.go`                   | MODIFY | Add dispatch entry                |
| `cmd/rootusage.go`              | MODIFY | Add help line                     |
| `constants/constants_cli.go`    | MODIFY | Add CmdPrune, CmdPruneAlias      |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No magic strings — all in `constants/constants_prune.go`.
- Positive conditionals only.
- Blank line before every `return`.
