# gitmap merge-left

One-way file-level merge that writes only into LEFT. Files missing
on LEFT are copied from RIGHT; conflicts are resolved into LEFT.
RIGHT is never modified. If LEFT originated from a URL it is
committed + pushed after the merge.

Spec: `spec/01-app/97-move-and-merge.md`

## Alias

ml

## Usage

    gitmap merge-left LEFT RIGHT [flags]
    gitmap ml         LEFT RIGHT [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -y, --yes, -a, --accept-all | false | Bypass prompt; default is `--prefer-right` |
| --prefer-left | false | LEFT always wins (skip RIGHT's version) |
| --prefer-right | false | RIGHT always wins (overwrite LEFT) |
| --prefer-newer | false | Newer mtime wins |
| --prefer-skip | false | Skip every conflict |
| --no-push | false | Skip git push on URL LEFT |
| --no-commit | false | Skip commit and push on URL LEFT |
| --force-folder | false | Replace folder whose origin doesn't match URL |
| --pull | false | Force `git pull --ff-only` on a folder endpoint |
| --dry-run | false | Print every action; perform none |
| --include-vcs | false | Include `.git/` in copy/diff |
| --include-node-modules | false | Include `node_modules/` in copy/diff |

## Prerequisites

None.

## Examples

### Example 1: Pull RIGHT's changes into LEFT (interactive)

    gitmap merge-left ./gitmap-v27 ./gitmap-v27

**Output:**

    [merge-left] diffing trees ...
    [merge-left]   conflict: docs/readme.md
    [merge-left]     LEFT  : 2.1 KB  modified 2026-04-15 10:00
    [merge-left]     RIGHT : 2.4 KB  modified 2026-04-18 11:30
    [merge-left]   > R
    [merge-left]   conflict docs/readme.md -> took RIGHT
    [merge-left] done

### Example 2: Bypass prompt (RIGHT always wins)

    gitmap ml ./local https://github.com/owner/upstream -y

**Output:**

    [merge-left] resolving RIGHT : https://github.com/owner/upstream
    [merge-left]   -> folder does not exist; cloning
    [merge-left]   -> clone OK
    [merge-left] diffing trees ...
    [merge-left]   conflict main.go -> took RIGHT
    [merge-left] done

### Example 3: Keep LEFT's version everywhere

    gitmap merge-left ./mine ./theirs -y --prefer-left

**Output:**

    [merge-left] diffing trees ...
    [merge-left]   conflict app.ts -> took LEFT
    [merge-left]   conflict app.ts -> skipped
    [merge-left] done

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Resolution, copy, commit, push failed, or user pressed Q |
| 2 | Wrong number of positional arguments |

## Notes

- RIGHT is read-only for `merge-left`; no commit or push happens
  on RIGHT even when it is a URL endpoint.
- With `-y`, the per-command default is `--prefer-right` (treat
  RIGHT as the upstream source of truth).

## See Also

- [merge-right](merge-right.md) — Mirror operation: write into RIGHT only
- [merge-both](merge-both.md) — Two-way merge
- [mv](mv.md) — Move LEFT into RIGHT and delete LEFT

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter merge-left
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
