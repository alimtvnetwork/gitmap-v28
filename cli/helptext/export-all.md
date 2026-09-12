# gitmap chrome export-all

Export all Chrome profiles in batch to individual JSON files, SQLite DBs, or a unified ZIP archive.

## Usage

```bash
gitmap chrome export-all [out-dir|out.zip] [flags]
gitmap chrome cpe-all [out-dir|out.zip] [flags]
```

## Description

Extracts bookmarks, preferences, and installed extension metadata across all discovered Google Chrome profiles on the local machine and packages them into the specified format (defaults to ZIP).

## Flags

- `--format <zip|json|sqlite|yaml>`: Target export format.
- `--limit <n>`: Maximum number of profiles to export.

## Examples

```bash
# Export all profiles to a zip archive
gitmap chrome export-all chrome-backup.zip

# Export all profiles to a directory
gitmap chrome export-all ./chrome-profiles --format json

# Export all profiles to SQLite database
gitmap chrome export-all chrome.sqlite --format sqlite
```

## See Also

- [chrome-profile-export](chrome-profile-export.md)
- [import-all](import-all.md)
- [copy-all](copy-all.md)
