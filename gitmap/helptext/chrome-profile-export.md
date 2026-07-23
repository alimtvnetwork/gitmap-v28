# gitmap chrome-profile-export

Export a Chrome profile to a JSON + CSV snapshot pair. The JSON is the
full-fidelity restore source; the CSV is a human-readable companion
that lists extension IDs and known preferences.

## Alias

`cpe`

## Usage

    gitmap chrome-profile-export <name> [out.json]
    gitmap cpe <name> [out.json]

If `out.json` is omitted, both artifacts are written to
`.gitmap/chrome/<name>.json` and `.gitmap/chrome/<name>.csv`.

## Prerequisites

- The profile must exist under Chrome's User Data root. On miss, the
  error lists every real profile (v6.34.0+).
- Close Chrome before exporting for a consistent snapshot.

## Examples

### Default location

    $ gitmap cpe Default
    Artifacts:
      json  C:\dev\.gitmap\chrome\Default.json
      csv   C:\dev\.gitmap\chrome\Default.csv
    chrome-profile-export: wrote C:\dev\.gitmap\chrome\Default.json (28144 bytes)
    chrome-profile-export: csv  C:\dev\.gitmap\chrome\Default.csv (1812 bytes)
    chrome-profile: db synced (Default)

### Custom destination

    $ gitmap cpe "Profile 1" D:\backup\work.json
    Artifacts:
      json  D:\backup\work.json
      csv   D:\backup\work.csv

## Exit codes

| Code | Meaning                     |
|------|-----------------------------|
| 0    | Export succeeded            |
| 6    | Usage error (missing name)  |
| 7    | Profile not found           |

## See also

- [chrome-profile-import](chrome-profile-import.md)
- [chrome-profile-copy](chrome-profile-copy.md)
- [chrome-profile-list](chrome-profile-list.md)

## Examples

```bash
gitmap chrome-profile-export
```
