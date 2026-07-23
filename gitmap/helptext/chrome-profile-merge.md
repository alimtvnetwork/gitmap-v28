# gitmap chrome-profile-merge

Merge selected slices of a source Chrome profile **into** a destination
profile **without clobbering destination values**. Unlike
`chrome-profile-copy` (which mirrors a curated subset of files src→dst),
`chrome-profile-merge` walks Preferences keys, bookmarks, and
extension folders one at a time and lets you pick the winner per
conflict.

## Alias

`cpm`

## Usage

    gitmap chrome-profile-merge <src> <dst> [--what all|settings|bookmarks|extensions] [--yes|--force] [--dry-run]
    gitmap cpm <src> <dst> [--what ...] [-y]

`<src>` and `<dst>` accept profile directory names (`Default`,
`Profile 7`) or display names from Chrome's picker (`Lovable`).

## --what

| Value         | What gets merged                                                |
|---------------|-----------------------------------------------------------------|
| `all`         | default — settings + bookmarks + extensions                     |
| `settings`    | top-level keys in `Preferences`                                 |
| `bookmarks`   | `Bookmarks` file under `roots.{bookmark_bar,other,synced}`      |
| `extensions`  | per-extension subdir under `Extensions/` (add-only)             |

## Conflict policy

Default is **interactive** — for every conflicting key/bookmark you
get:

    [k]eep destination, [o]verwrite, [a]ll-keep, [A]ll-overwrite, [q]uit

- `--yes` / `-y` — non-interactive: auto-**keep destination** on every
  conflict (safe; never loses your edits).
- `--force` — non-interactive: auto-**overwrite destination** with
  source on every conflict.
- `--dry-run` — print plan only; no files are written.

## Examples

### Pull bookmarks from a backup into your daily profile, skip duplicates

    $ gitmap cpm Backup Default --what bookmarks --yes
    ▸ chrome-profile-merge  Backup → Default  (what=bookmarks)
      source      C:\...\User Data\Backup
      destination C:\...\User Data\Default

    • bookmarks
    ✓ merge complete  added=12 skipped=3 overwrote=0

### Mirror every setting from Work into a fresh experimental profile

    gitmap cpm Work Experiments --what settings --force

### See what would change without writing

    gitmap cpm Default Backup --dry-run

## Exit codes

| Code | Meaning                          |
|------|----------------------------------|
| 0    | Merge succeeded (or aborted with `q`) |
| 6    | Usage error                      |
| 7    | Source or destination not found  |

## See also

- [chrome-profile-copy](chrome-profile-copy.md) — mirror copy (offline)
- [chrome-profile-export](chrome-profile-export.md) — snapshot to JSON+CSV
- [chrome-profile-list](chrome-profile-list.md)

## Examples

```bash
gitmap chrome-profile-merge
```
