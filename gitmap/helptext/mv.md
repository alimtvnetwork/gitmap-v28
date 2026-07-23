# gitmap mv

Move every file from LEFT into RIGHT, then delete LEFT entirely.
Either endpoint can be a local folder OR a remote git URL with an
optional `:branch` suffix. URL endpoints are auto-cloned (or
re-pulled if the working folder already matches origin), and a
commit + push is made on the URL side after the file copy.

Spec: `spec/01-app/97-move-and-merge.md`

## Alias

move

## Usage

    gitmap mv   LEFT RIGHT [flags]
    gitmap move LEFT RIGHT [flags]

LEFT and RIGHT can each be:

- a local folder path (relative or absolute)
- a remote git URL with optional `:branch` suffix
  (e.g. `https://github.com/owner/repo:develop`)

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --no-push | false | Skip git push on URL endpoints (still commits) |
| --no-commit | false | Skip both commit and push on URL endpoints |
| --force-folder | false | Replace folder whose origin doesn't match the URL |
| --pull | false | Force `git pull --ff-only` on a folder endpoint |
| --init | false | When RIGHT is auto-created, also `git init` it |
| --dry-run | false | Print every action; perform none |
| --include-vcs | false | Include `.git/` in the copy (default: excluded) |
| --include-node-modules | false | Include `node_modules/` in the copy |

## Prerequisites

None. Both endpoints are resolved on demand: folders are validated
in place and URLs are cloned into a folder named after the repo.

## Examples

### Example 1: Move one local folder into another

    gitmap mv ./gitmap-v28 ./gitmap-v28

**Output:**

    [mv] resolving LEFT : ./gitmap-v28 (folder)
    [mv] resolving RIGHT : ./gitmap-v28 (folder)
    [mv] copying files LEFT -> RIGHT (excluding .git/) ...
    [mv]   copied 142 files
    [mv] deleting LEFT (./gitmap-v28) ...
    [mv]   deleted
    [mv] done

### Example 2: Move a local folder into a remote repo (clone + push)

    gitmap mv ./gitmap-v28 https://github.com/owner/gitmap-v28

**Output:**

    [mv] resolving RIGHT : https://github.com/owner/gitmap-v28
    [mv]   -> mapped to working folder: /work/gitmap-v28
    [mv]   -> folder does not exist; cloning
    [mv]   -> clone OK
    [mv] copying files LEFT -> RIGHT (excluding .git/) ...
    [mv]   copied 142 files
    [mv] committing in https://github.com/owner/gitmap-v28 ...
    [mv]   commit a1b2c3d "gitmap mv from ./gitmap-v28"
    [mv] pushing https://github.com/owner/gitmap-v28 ...
    [mv]   push OK
    [mv] done

### Example 3: Preview without writing anything

    gitmap mv ./gitmap-v28 ./gitmap-v28 --dry-run

**Output:**

    [mv] copying files LEFT -> RIGHT (excluding .git/) ...
    [mv]   copied 142 files
    [mv] deleting LEFT (./gitmap-v28) ...
    [mv]   deleted
    [mv] done

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Resolution, copy, commit, or push failed (message on stderr) |
| 2 | Wrong number of positional arguments (need exactly LEFT RIGHT) |

## Notes

- The `.git/` folder is never copied; LEFT's `.git/` is removed
  along with the rest of LEFT after the copy.
- LEFT and RIGHT must not resolve to the same folder, and one must
  not be nested inside the other — checked before any write.
- On a URL endpoint, the commit message is `gitmap mv from <LEFT>`.

## See Also

- [merge-both](merge-both.md) — Two-way file-level merge
- [merge-left](merge-left.md) — Merge into LEFT only
- [merge-right](merge-right.md) — Merge into RIGHT only
- [clone](clone.md) — Clone repositories from URL or scan output

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter mv
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
