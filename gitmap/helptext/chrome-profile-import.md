# gitmap chrome-profile-import

Import a Chrome profile from a JSON or ZIP export file.

## Usage
    gitmap chrome-profile-import <file.json|file.zip|file.sqlite> [dst-profile]

- <file>: The snapshot or ZIP archive to import.
- [dst-profile]: Target profile name. If omitted, uses the name embedded in the snapshot/zip.

## Examples
    gitmap chrome-profile-import .gitmap\chrome\Default.json
    chrome-profile-import: imported .gitmap\chrome\Default.json into profile "Default"
    chrome-profile: db synced (Default)

    gitmap chrome-profile-import work.zip WorkRestored

## See Also
- [chrome-profile-export](chrome-profile-export.md)
- [chrome-profile-copy](chrome-profile-copy.md)
- [chrome-profile-list](chrome-profile-list.md)
