# gitmap merge-both

Two-way file-level merge between LEFT and RIGHT. Files present on
only one side are copied to the other; files present on both with
different content trigger an interactive conflict prompt. Each side
that originated from a URL is committed + pushed independently.

Spec: `spec/01-app/97-move-and-merge.md`

## Alias

mb

## Usage

    gitmap merge-both LEFT RIGHT [flags]
    gitmap mb         LEFT RIGHT [flags]

LEFT and RIGHT can each be a folder path or a remote git URL
(optionally suffixed with `:branch`).

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -y, --yes, -a, --accept-all | false | Bypass prompt; default is `--prefer-newer` |
| --prefer-left | false | LEFT always wins on conflict |
| --prefer-right | false | RIGHT always wins on conflict |
| --prefer-newer | false | Newer mtime wins on conflict |
| --prefer-skip | false | Skip every conflict; only missing files copied |
| --no-push | false | Skip git push on URL endpoints |
| --no-commit | false | Skip commit and push on URL endpoints |
| --force-folder | false | Replace folder whose origin doesn't match URL |
| --pull | false | Force `git pull --ff-only` on a folder endpoint |
| --dry-run | false | Print every action; perform none |
| --include-vcs | false | Include `.git/` in copy/diff |
| --include-node-modules | false | Include `node_modules/` in copy/diff |

## Prerequisites

None.

## Conflict Prompt

For every conflicting file:

    [L]eft  [R]ight  [S]kip  [A]ll-left  [B]all-right  [Q]uit

| Key | Action |
|-----|--------|
| L | Take LEFT's version (write into RIGHT) |
| R | Take RIGHT's version (write into LEFT) |
| S | Skip this file |
| A | Sticky: All-Left for the rest of the run |
| B | Sticky: All-Right for the rest of the run |
| Q | Quit immediately (already-applied changes are kept) |

## Examples

### Example 1: Interactive two-way merge

    gitmap merge-both ./gitmap-v27 ./gitmap-v27

**Output:**

    [merge-both] diffing trees ...
    [merge-both]   conflict: README.md
    [merge-both]     LEFT  : 4.2 KB  modified 2026-04-17 14:02
    [merge-both]     RIGHT : 4.1 KB  modified 2026-04-18 09:11
    [merge-both]   > B
    [merge-both]   conflict README.md -> took RIGHT
    [merge-both] done

### Example 2: Non-interactive (newer wins) with URL endpoint

    gitmap mb ./local https://github.com/owner/repo -y --prefer-newer

**Output:**

    [merge-both] resolving RIGHT : https://github.com/owner/repo
    [merge-both]   -> folder does not exist; cloning
    [merge-both]   -> clone OK
    [merge-both] diffing trees ...
    [merge-both] committing in https://github.com/owner/repo ...
    [merge-both]   commit 9f2e1ab "gitmap merge-both with ./local"
    [merge-both] pushing https://github.com/owner/repo ...
    [merge-both]   push OK
    [merge-both] done

### Example 3: Dry-run with prefer-left override

    gitmap merge-both ./a ./b -y --prefer-left --dry-run

**Output:**

    [merge-both] diffing trees ...
    [merge-both]   conflict config.json -> took LEFT
    [merge-both]   [dry-run] copy config.json -> RIGHT
    [merge-both] done

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Resolution, copy, commit, push failed, or user pressed Q |
| 2 | Wrong number of positional arguments |

## See Also

- [mv](mv.md) — Move LEFT into RIGHT and delete LEFT
- [merge-left](merge-left.md) — Write into LEFT only
- [merge-right](merge-right.md) — Write into RIGHT only

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter merge-both
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
