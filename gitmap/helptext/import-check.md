# gitmap chrome import-check

Inspect snapshot files, directories, or archives and preview profile metadata before importing.

## Usage

```bash
gitmap chrome import-check [path|ls] [flags]
gitmap chrome import ls [path] [flags]
gitmap chrome import-all ls [path] [flags]
```

## Description

Discovers all Chrome profiles contained in a folder, ZIP archive, or snapshot file without altering any local Chrome profiles. Renders profile names, display names, email addresses, bookmarks count, extensions count, and planned actions.

## Flags

- `--json`: Output discovered profile candidates as formatted JSON.
- `--file <path>`: Write inspection report or JSON output to specified file path.
- `--fnf <path>`: Save to file or fail if no candidates match.
- `--tempfile <filename>`: Write inspection report to `.lovable/temp/<filename>`.

## Examples

```bash
# Inspect profiles in current directory
gitmap chrome import-check

# Inspect a folder of profiles
gitmap chrome import ls ./chrome-ext

# Inspect with JSON output for automation
gitmap chrome import-check ./chrome-ext --json

# Export inspection to JSON file on file system
gitmap chrome import-check ./chrome-ext --json --file preview.json

# Save inspection report with --fnf flag
gitmap chrome import-check ./chrome-ext --json --fnf candidates.json
```

## See Also

- [import-all](import-all.md)
- [chrome-profile-import](chrome-profile-import.md)
