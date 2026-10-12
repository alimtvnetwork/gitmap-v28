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
2. Excludes milestone anchors (`start`, `end`) and pinned tags
3. Deletes remote GitHub Release entities (`gh release delete`)
4. Deletes remote Git tags (`git push origin --delete <tags...>`)
5. Deletes local Git tags (`git tag -d <tags...>`)
6. Removes local release metadata sidecars (`.gitmap/release/<tag>.json`)
7. Prunes remote tracking refs and runs garbage collection (`git gc --prune=now`)

## Flags

- `--pin <tag1,tag2...>` — Ad-hoc additional tags to protect from removal
- `--keep-remote`       — Only delete local tags; preserve origin tags
- `--dry-run`           — Preview candidate tags without modifying anything
- `-y, --yes`           — Skip the confirmation prompt

## Examples

```bash
# Keep only 4.162 and 4.190, removing intermediate tags 4.162.1..4.189.0
gitmap rrb 4.162, 4.190

# Same as above with space-separated syntax
gitmap rrb 4.162 4.190 -y

# Preview which tags would be deleted before running
gitmap rrb 4.162, 4.190 --dry-run

# Protect an additional specific tag during range deletion
gitmap rrb 4.162, 4.190 --pin 4.175.0 -y
```

## See Also

- `gitmap pnr` — Pin critical tags so they are immune to range deletion
- `gitmap ru`  — Roll back / undo a specific release tag
