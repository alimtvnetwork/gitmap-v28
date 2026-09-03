# chrome

Unified CLI suite for managing Google Chrome installation, profile replication, URL & tab launching, extensions/plugins, experimental flags, cache reset, observation, and backup.

## Aliases

`gitmap cprof <subcommand>`  
`gitmap chrome-profile <subcommand>`

---

## Subcommands

### Launching & Tabs

- `open` (`launch`, `tab`) `[urls...] [--profile=<name>] [--incognito] [--new-window] [--app=<url>]` — Open URL(s) in specified or active profile, or launch multiple profiles in one command (e.g. `"Profile 1=https://a.com,Profile 2=https://b.com"`).
- `observe` (`tabs`, `status`, `ps`) `[profile] [--json|--yaml|--all]` — Inspect active Chrome processes, open tabs, page titles, and loading status.

### Extensions & Plugins

- `extensions` (`ext`, `plugins`) `[profile] [--json|--yaml|--all]` — List installed extensions, versions, IDs, and enabled/disabled state.
- `extension-install` (`ext-in`, `plugin-install`) `<path> [--profile=<name>]` — Inject an unpacked extension directory or `.crx` file into a profile.
- `extension-enable` (`ext-on`, `plugin-enable`) `<pattern|id> [--profile=<name>]` — Enable extension(s) matching an ID or pattern.
- `extension-disable` (`ext-off`, `plugin-disable`) `<pattern|id> [--profile=<name>]` — Disable extension(s) matching an ID or pattern.
- `extension-disable-all` (`ext-off-all`) `[--profile=<name>]` — Disable ALL extensions in the target profile.

### Feature Flags & Profile Reset

- `flags` (`experiments`) `[ls|enable|disable|reset] [flag-name]` — Inspect, toggle, or reset Chrome experimental feature flags in `Local State`.
- `reset` (`clean`, `clear`) `[profile] [--cache|--cookies|--history|--extensions|--all]` — Purge caches, cookies, history, or restore clean default preferences.

### Single Profile Operations

- `copy` (`cpc`, `profile-copy`) `<src> <dst>` — Copy a Chrome profile into an offline destination profile.
- `export` (`cpe`, `profile-export`) `<name> [out] [--format=json|sqlite|yaml|zip]` — Export profile snapshot (auto-infers format from extension).
- `import` (`cpi`, `profile-import`) `<file> [name]` — Import a profile from a JSON, SQLite DB, YAML, or ZIP snapshot.
- `undo` `[profile]` — Revert the last profile import or restore from an automatic snapshot.
- `redo` `[profile]` — Reapply the undone Chrome profile mutation.
- `group` `[ls|add <group> <profiles...>|rm <group> [profile]]` — Manage named groups of Chrome profiles.
- `list` (`ls`, `cpl`, `profiles`) — List all Chrome profiles discovered on the local machine.
- `delete` (`rm`, `del`, `cpd`) `<name> [--yes]` — Remove a profile and its stored artifacts.
- `merge` (`cpm`) `<src> <dst> [--what=all|settings|bookmarks|extensions]` — Merge pieces of one profile into another.

### Batch Profile Operations

- `copy-all` (`all-profile-copy`, `cpc-all`) `[dst-dir]` — Replicate ALL discovered Chrome profiles to a directory.
- `export-all` (`all-profile-export`, `cpe-all`) `[out-dir]` — Batch export all Chrome profiles to JSON, SQLite, YAML, or ZIP.
- `import-all` (`all-profile-import`, `cpi-all`) `<dir>` — Batch import all profile snapshots from a directory.

### Backup, Diff & Maintenance

- `install` (`in`) — Install Google Chrome via system package manager (Winget, Chocolatey, Apt, Homebrew).
- `backup` `[--out <tarball>]` — Snapshot all Chrome profiles into a compressed `tar.gz` archive.
- `restore` `<tarball>` — Restore all Chrome profiles from a `tar.gz` archive.
- `diff` `<A> <B>` — Compare extensions, bookmarks, and settings between two profiles.
- `export-bookmarks` (`bookmarks`) `<profile> [--format md|html|json] [--out <file>]` — Export Chrome bookmarks tree.
- `which` — Print Chrome User Data root and executable installation path.

---

## Examples

```bash

# Launch URL in default profile

gitmap chrome open "https://github.com"

# Launch in specific profile

gitmap chrome open "https://github.com" --profile="Default"
gitmap chrome open "Profile 1" "https://google.com"

# Launch multiple profiles with different URLs in one command

gitmap chrome open "Profile 1=https://github.com,Profile 2=https://google.com"

# Observe open tabs and active processes

gitmap chrome observe
gitmap chrome observe --json
gitmap chrome observe --yaml

# List extensions

gitmap chrome extensions
gitmap chrome extensions "Default" --json

# Inject / install local extension

gitmap chrome extension-install ./my-extension --profile="Default"

# Enable / disable extensions

gitmap chrome extension-disable "Dark Reader" --profile="Default"
gitmap chrome extension-enable "Dark Reader" --profile="Default"
gitmap chrome extension-disable-all --profile="Default"

# Inspect and toggle experimental flags

gitmap chrome flags
gitmap chrome flags enable enable-gpu-rasterization
gitmap chrome flags disable enable-gpu-rasterization
gitmap chrome flags reset

# Reset cache or full profile

gitmap chrome reset Default --cache
gitmap chrome reset Default --all
```
