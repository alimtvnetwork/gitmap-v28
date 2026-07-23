# gitmap merge-right

One-way file-level merge that writes only into RIGHT. Files missing
on RIGHT are copied from LEFT; conflicts are resolved into RIGHT.
LEFT is never modified. If RIGHT originated from a URL it is
committed + pushed after the merge.

Spec: `spec/01-app/97-move-and-merge.md`

## Alias

mr

## Usage

    gitmap merge-right LEFT RIGHT [flags]
    gitmap mr          LEFT RIGHT [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -y, --yes, -a, --accept-all | false | Bypass prompt; default is `--prefer-left` |
| --prefer-left | false | LEFT always wins (overwrite RIGHT) |
| --prefer-right | false | RIGHT always wins (skip LEFT's version) |
| --prefer-newer | false | Newer mtime wins |
| --prefer-skip | false | Skip every conflict |
| --no-push | false | Skip git push on URL RIGHT |
| --no-commit | false | Skip commit and push on URL RIGHT |
| --force-folder | false | Replace folder whose origin doesn't match URL |
| --pull | false | Force `git pull --ff-only` on a folder endpoint |
| --dry-run | false | Print every action; perform none |
| --include-vcs | false | Include `.git/` in copy/diff |
| --include-node-modules | false | Include `node_modules/` in copy/diff |

## Prerequisites

None.

## Examples

### Example 1: Push LEFT's changes into a remote repo

    gitmap merge-right ./local https://github.com/owner/repo -y

**Output:**

    [merge-right] resolving RIGHT : https://github.com/owner/repo
    [merge-right]   -> folder does not exist; cloning
    [merge-right]   -> clone OK
    [merge-right] diffing trees ...
    [merge-right]   conflict src/app.ts -> took LEFT
    [merge-right] committing in https://github.com/owner/repo ...
    [merge-right]   commit 7a3f9c2 "gitmap merge-right from ./local"
    [merge-right] pushing https://github.com/owner/repo ...
    [merge-right]   push OK
    [merge-right] done

### Example 2: Push to a specific branch

    gitmap mr ./local https://github.com/owner/repo:develop -y

**Output:**

    [merge-right] resolving RIGHT : https://github.com/owner/repo:develop
    [merge-right]   -> folder does not exist; cloning
    [merge-right]   -> clone OK
    [merge-right] pushing https://github.com/owner/repo ...
    [merge-right]   push OK
    [merge-right] done

### Example 3: Stage changes locally without pushing

    gitmap merge-right ./local ./mirror --no-push

**Output:**

    [merge-right] diffing trees ...
    [merge-right]   conflict README.md -> took LEFT
    [merge-right] done

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Resolution, copy, commit, push failed, or user pressed Q |
| 2 | Wrong number of positional arguments |

## Notes

- LEFT is read-only for `merge-right`; no commit or push happens
  on LEFT even when it is a URL endpoint.
- With `-y`, the per-command default is `--prefer-left` (treat
  LEFT as the source of truth being published).

## See Also

- [merge-left](merge-left.md) — Mirror operation: write into LEFT only
- [merge-both](merge-both.md) — Two-way merge
- [mv](mv.md) — Move LEFT into RIGHT and delete LEFT

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter merge-right
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
