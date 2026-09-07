# gitmap chrome copy-all

Copy all discovered Chrome profiles to a destination directory.

## Usage

```bash
gitmap chrome copy-all <dst-dir> [flags]
gitmap chrome cpc-all <dst-dir> [flags]
```

## Description

Copies each active Chrome profile directory from Chrome User Data into offline standalone profile directories under the specified target folder.

## Flags

- `--limit <n>`: Maximum number of profiles to copy.

## Examples

```bash
# Copy all profiles to backup folder
gitmap chrome copy-all /backup/chrome-profiles

# Copy first 3 profiles
gitmap chrome copy-all /backup/chrome-profiles --limit 3
```

## See Also

- [chrome-profile-copy](chrome-profile-copy.md)
- [export-all](export-all.md)
- [import-all](import-all.md)
