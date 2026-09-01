# chrome

Unified CLI suite for managing Google Chrome installation, profile replication, backup, export/import, diff, and bookmark extraction.

## Aliases

`gitmap cprof <subcommand>`  
`gitmap chrome-profile <subcommand>`

---

## Subcommands

### Installation & System Setup
- `install` (`in`) — Install Google Chrome via system package manager (Winget, Chocolatey, Apt, Homebrew).

### Single Profile Operations
- `copy` (`cpc`, `profile-copy`) `<src> <dst>` — Copy a Chrome profile (bookmarks, extensions, preferences, flags) into an offline destination profile.
- `export` (`cpe`, `profile-export`) `<name> [out]` — Export a profile snapshot (`--format=json|zip|sqlite`).
- `import` (`cpi`, `profile-import`) `<file> [name]` — Import a profile from a JSON, ZIP, or SQLite export file.
- `list` (`ls`, `cpl`, `profiles`) — List all Chrome profiles discovered on the local machine and tracked in GitMap database.
- `delete` (`rm`, `del`, `cpd`) `<name> [--yes]` — Remove a profile and its stored artifacts from the database.
- `merge` (`cpm`) `<src> <dst> [--what=all|settings|bookmarks|extensions]` — Merge pieces of one profile into another.

### Batch Profile Operations
- `copy-all` (`all-profile-copy`, `cpc-all`) `[dst-dir]` — Replicate ALL discovered Chrome profiles to a destination directory.
- `export-all` (`all-profile-export`, `cpe-all`) `[out-dir] [--format=json|zip|sqlite]` — Batch export all Chrome profiles.
- `import-all` (`all-profile-import`, `cpi-all`) `<dir>` — Batch import all profile snapshots from a directory.

### Backup, Diff & Inspection
- `backup` `[--out <tarball>]` — Snapshot all Chrome profiles into a compressed `tar.gz` archive.
- `restore` `<tarball>` — Restore all Chrome profiles from a `tar.gz` archive.
- `diff` `<A> <B>` — Compare extensions, bookmarks, and settings between two profiles.
- `export-bookmarks` (`bookmarks`) `<profile> [--format md|html|json] [--out <file>]` — Export Chrome bookmarks tree.
- `which` — Print Chrome User Data root and executable installation path.

---

## Examples

```bash
# Install Chrome
gitmap chrome install

# List local profiles
gitmap chrome list
gitmap cp ls

# Replicate a single profile
gitmap chrome copy Default "Profile Work"
gitmap cp cpc Default "Profile Work"

# Batch copy all profiles
gitmap chrome copy-all ~/chrome-backups/all-profiles

# Batch export all profiles to JSON snapshots
gitmap chrome export-all .gitmap/chrome/

# Batch import profile snapshots
gitmap chrome import-all ~/chrome-backups/snapshots/

# Snapshot all profiles into a tarball
gitmap chrome backup --out ~/chrome-2026.tar.gz

# Restore from tarball
gitmap chrome restore ~/chrome-2026.tar.gz

# Compare two profiles
gitmap chrome diff Default "Profile 1"

# Export bookmarks
gitmap chrome export-bookmarks Default --format html --out bookmarks.html
```
