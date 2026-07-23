# gitmap chrome-profile-import

Restore a Chrome profile from a previously exported snapshot. Accepts
either a `.json` file (full fidelity) or a `.csv` file (extension IDs
and known preferences — bookmarks are NOT restored from CSV).

## Alias

`cpi`

## Usage

    gitmap chrome-profile-import <file.json|file.csv> [dst-profile]
    gitmap cpi <file.json|file.csv> [dst-profile]

If `dst-profile` is omitted, the destination name is inferred from the
file basename (e.g. `Default.json` → `Default`).

## Examples

### Full restore from JSON

    $ gitmap cpi .gitmap\chrome\Default.json
    chrome-profile-import: imported .gitmap\chrome\Default.json into profile "Default"
    chrome-profile: db synced (Default)

### Partial restore from CSV

    $ gitmap cpi work.csv WorkRestored
    chrome-profile-import: csv source detected — restoring extension IDs + known preferences (bookmarks omitted)
    chrome-profile-import: imported work.csv into profile "WorkRestored"

## Exit codes

| Code | Meaning                     |
|------|-----------------------------|
| 0    | Import succeeded            |
| 6    | Usage error (missing file)  |

## See also

- [chrome-profile-export](chrome-profile-export.md)
- [chrome-profile-copy](chrome-profile-copy.md)
- [chrome-profile-list](chrome-profile-list.md)

## Examples

```bash
gitmap chrome-profile-import
```
