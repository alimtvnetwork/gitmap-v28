# release-remove-between (rrb)

Prune intermediate release tags between two bounds, strictly preserving both endpoint milestones and pinned tags.

## Synopsis

```
gitmap release remove-between <start> <end> [flags]
gitmap rrb                    <start> <end> [flags]
```

## Description

Iterates through all Git tags in the repository, evaluates versions in semantic order between `<start>` and `<end>` (exclusive), and safely deletes intermediate tags and release objects.

Milestone anchors `<start>` and `<end>` are **always kept**. Any tag registered in `.gitmap/release/pinned-tags.json` or passed via `--pin` is also **strictly kept**.

## What it does

1. Validates SemVer bounds (`<start>` and `<end>`)
2. Evaluates repository tags and outputs an **automatic dry-run preview**
3. Excludes milestone anchors (`start`, `end`) and pinned tags
4. Prompts the user for confirmation (bypassable **only** via `--confirm` or `-y`)
5. Deletes remote GitHub Release entities (`gh release delete`)
6. Deletes remote Git tags in batch chunks (`git push origin --delete <tags...>`)
7. Deletes local Git tags in batch (`git tag -d <tags...>`)
8. Removes local release metadata sidecars (`.gitmap/release/<tag>.json`)
9. Prunes remote tracking refs and runs garbage collection (`git gc --prune=now`)

## Flags

- `--confirm`           — Bypass confirmation prompt and proceed immediately with removal
- `-y, --yes`           — Alias for `--confirm`
- `--dry-run`           — Preview candidate tags without prompting or making any modifications
- `--pin <tag1,tag2...>` — Ad-hoc additional tags to protect from removal
- `--keep-remote`       — Only delete local tags; preserve origin tags

## Examples

```bash
# Evaluates and prints auto dry-run preview, then prompts for confirmation:
gitmap rrb 4.162, 4.190

# Bypass confirmation prompt directly via --confirm:
gitmap rrb 4.162, 4.190 --confirm

# Bypass confirmation prompt with shorthand -y:
gitmap rrb 4.162 4.190 -y

# Pure preview mode (no prompt, zero side-effects):
gitmap rrb 4.162, 4.190 --dry-run

# Protect an additional specific tag during range deletion:
gitmap rrb 4.162, 4.190 --pin 4.175.0 --confirm
```

## See Also

- `gitmap pnr` — Pin critical tags so they are immune to range deletion
- `gitmap ru`  — Roll back / undo a specific release tag
