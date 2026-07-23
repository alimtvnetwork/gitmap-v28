# chrome

Umbrella for Chrome profile snapshot/diff utilities.

## Subcommands

- `backup` — snapshot all profiles to a tar.gz under `.gitmap/chrome/backup/`.
- `restore <tarball>` — restore a snapshot into the User Data dir.
- `diff <A> <B>` — list extensions/bookmarks only-in-A vs only-in-B.
- `export-bookmarks <profile> [--format md|html|json] [--out <file>] [--root <bookmark_bar|other|synced>] [--folder <path/to/folder>] [--match <substr>] [--title <exact>]` — export bookmarks tree.
- `which` — print directory + display name of the currently-active profile.

## Examples

```bash
gitmap chrome backup
gitmap chrome backup --out ~/chrome-2026.tar.gz
gitmap chrome restore ~/chrome-2026.tar.gz
gitmap chrome diff Default "Profile 1"
gitmap chrome which

# export-bookmarks — md/html/json with --root and --folder
# Whole tree as Markdown (default):
gitmap chrome export-bookmarks Lovable --out bm.md

# Only the bookmarks bar, as HTML:
gitmap chrome export-bookmarks Lovable --root bookmark_bar --format html --out bar.html

# A nested subtree from the bookmarks bar, as JSON:
gitmap chrome export-bookmarks "Profile 1" --root bookmark_bar --folder "Work/Docs" --format json --out work-docs.json

# "Other bookmarks" → Markdown to stdout:
gitmap chrome export-bookmarks Default --root other --format md

# Synced bookmarks → HTML file:
gitmap chrome export-bookmarks Default --root synced --format html --out synced.html

# Combine --folder with --match to grab only items containing "spec":
gitmap chrome export-bookmarks Lovable --root bookmark_bar --folder "Work" --match spec --format json
```

Notes:
- `--root` accepts `bookmark_bar`, `other`, or `synced` (case-insensitive).
- `--folder` is slash-delimited and case-insensitive (e.g. `--folder "Work/Docs"`).
- Omit `--out` to print to stdout; combine with shell redirection if preferred.

