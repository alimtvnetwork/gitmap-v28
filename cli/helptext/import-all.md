# gitmap chrome import-all

Import all Chrome profile snapshots or archives from a directory or ZIP file.

## Usage

```bash
gitmap chrome import-all <dir|zip> [flags]
gitmap chrome cpi-all <dir|zip> [flags]
```

## Description

Recursively scans the specified folder or archive for Chrome profile snapshots (`.json`, `.zip`, `.sqlite`) and profile subdirectories (containing `Preferences`, `Bookmarks`, or `Profile.json`). Automatically resolves destination profiles using email matching or creates isolated new profiles to protect existing data.

Preflight inspection can be executed by running `gitmap chrome import-all ls <path>`.

## Flags

- `--limit <n>`: Maximum number of profiles to import.
- `--except <pattern>`: Exclude profiles matching name, email, or directory pattern.
- `--email <email>`: Filter and import only the profile matching the email address.
- `--json`: Format preview output as JSON.
- `--file <path>`: Write preview output to specified file path.
- `--fnf <path>`: Save to file or fail if no candidates match.
- `--tempfile <filename>`: Write output to `.lovable/temp/<filename>`.

## Examples

```bash
# Import all profiles from a directory
gitmap chrome import-all ./chrome-ext

# Import all profiles from a multi-profile ZIP archive
gitmap chrome import-all chrome-profiles.zip

# Preview discoverable profiles in directory before importing
gitmap chrome import-all ls ./chrome-ext

# Export preview to JSON file on file system
gitmap chrome import-all ls ./chrome-ext --json --file preview.json

# Export preview with --fnf flag
gitmap chrome import-all ls ./chrome-ext --json --fnf preview.json

# Import with limit
gitmap chrome import-all ./chrome-ext --limit 5
```

## See Also

- [chrome-profile-import](chrome-profile-import.md)
- [export-all](export-all.md)
- [copy-all](copy-all.md)
- [import-check](import-check.md)
