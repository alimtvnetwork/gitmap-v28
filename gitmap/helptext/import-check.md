# gitmap chrome import-check

Inspect snapshot files, directories, or archives and preview profile metadata before importing.

## Usage

```bash
gitmap chrome import-check [path|ls] [flags]
gitmap chrome import ls [path]
gitmap chrome import-all ls [path]
```

## Description

Discovers all Chrome profiles contained in a folder, ZIP archive, or snapshot file without altering any local Chrome profiles. Renders profile names, display names, email addresses, bookmarks count, extensions count, and planned actions.

## Flags

- `--json`: Output discovered profile candidates as formatted JSON.

## Examples

```bash
# Inspect profiles in current directory
gitmap chrome import-check

# Inspect a folder of profiles
gitmap chrome import ls ./chrome-ext

# Inspect a multi-profile ZIP archive
gitmap chrome import-check backup.zip

# Inspect with JSON output for automation
gitmap chrome import-check ./chrome-ext --json
```

## See Also

- [import-all](import-all.md)
- [chrome-profile-import](chrome-profile-import.md)
