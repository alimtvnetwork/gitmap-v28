# gitmap chrome-profile-delete

Remove a profile **and** its stored JSON/CSV artifacts from the gitmap
database. Does NOT touch the live Chrome User Data directory.

## Alias

`cpd`

## Usage

    gitmap chrome-profile-delete <name> --yes
    gitmap cpd <name> --yes

`--yes` is required — without it the command aborts with a hint.

## Examples

### Confirmed delete

    $ gitmap cpd "Profile 2" --yes
      rm C:\dev\.gitmap\chrome\Profile 2.json
      rm C:\dev\.gitmap\chrome\Profile 2.csv
    chrome-profile-delete: removed profile "Profile 2" (2 artifacts)

### Missing confirmation

    $ gitmap cpd "Profile 2"
    chrome-profile-delete: aborted — re-run with --yes to confirm

## Exit codes

| Code | Meaning                                         |
|------|-------------------------------------------------|
| 0    | Deletion succeeded (or aborted without --yes)   |
| 6    | Usage error (missing name)                      |

## See also

- [chrome-profile-list](chrome-profile-list.md)
- [chrome-profile-export](chrome-profile-export.md)

## Examples

```bash
gitmap chrome-profile-delete
```
