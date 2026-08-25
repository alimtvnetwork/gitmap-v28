# gitmap chrome-profile-export

Export a Chrome profile to a JSON + CSV snapshot pair, or a ZIP archive containing SQLite databases. The JSON is the full-fidelity restore source for text-based configs; the CSV is a human-readable companion that lists extension IDs and known preferences. The ZIP output also includes cross-OS unencrypted SQLite databases (History, Web Data, etc.).

## Usage
    gitmap chrome-profile-export <name> [out] [--format=json|zip|sqlite]

- <name>: The profile directory name (e.g. Default) or its display name (e.g. Work).
- [out]: Optional output path. Defaults to .gitmap/chrome/<name>.<format>.
- --format=zip|json|sqlite: Output format. Defaults to json.

## Examples
    gitmap chrome-profile-export Default
    chrome-profile-export: wrote C:\dev\.gitmap\chrome\Default.json (28144 bytes)
    chrome-profile-export: csv  C:\dev\.gitmap\chrome\Default.csv (1812 bytes)
    chrome-profile: db synced (Default)

    gitmap chrome-profile-export Work --format=zip

## See Also
- [chrome-profile-import](chrome-profile-import.md)
- [chrome-profile-copy](chrome-profile-copy.md)
- [chrome-profile-list](chrome-profile-list.md)
