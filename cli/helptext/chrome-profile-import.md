# gitmap chrome-profile-import

Import a Chrome profile from a JSON or ZIP export file.

## Usage

```bash
gitmap chrome-profile-import <file.json|file.zip|dir|email> [dst-profile] [flags]
gitmap chrome import <file|dir> [flags]
gitmap profile import [file|dir|email] [flags]
```

- `<file|dir|email>`: Snapshot file, multi-profile directory, ZIP archive, or email.
- `[dst-profile]`: Target profile name. If omitted, uses name embedded in snapshot.

## Flags

- `--json`: Output execution or preview report in JSON format.
- `--file <path>`: Export execution or inspection report to specified file.
- `--tempfile <name>`: Write report to `.lovable/temp/<name>`.
- `--fnf <path>`: Save to file or fail if snapshot has no valid profiles.
- `--except <rule>`: Exclude specific profiles by name or email.
- `--limit <n>`: Limit number of imported profiles.
- `--email <email>`: Filter and import only profile matching email.

## Examples

```bash
# Import a single snapshot
gitmap chrome-profile-import .gitmap/chrome/Default.json

# Import from multi-profile directory
gitmap chrome-profile-import ./chrome-ext

# Direct profile import with email filter
gitmap profile import erfan.office.n@gmail.com

# Direct profile import from current directory
gitmap profile import
```

## See Also

- [chrome-profile-export](chrome-profile-export.md)
- [chrome-profile-copy](chrome-profile-copy.md)
- [chrome-profile-list](chrome-profile-list.md)
- [import-all](import-all.md)
- [import-check](import-check.md)
