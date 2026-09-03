# Chrome Browser Profile Suite (`gitmap chrome`)

Manage, backup, clone, and merge Chrome browser profiles.

## Commands

| Subcommand | Alias | Description |
|---|---|---|
| `gitmap chrome open <profile>` | `launch` | Open Chrome with a specific user profile |
| `gitmap chrome observe` | `watch` | Monitor running Chrome processes, locks, and profile memory |
| `gitmap chrome list` | `cpl` | List all discovered Chrome profile directories |
| `gitmap chrome copy <src> <dst>` | `cpc` | Copy live profile into an offline backup profile |
| `gitmap chrome export <name>` | `cpe` | Export profile snapshot (`.json`, `.zip`, `.sqlite`) |
| `gitmap chrome import <file>` | `cpi` | Import a Chrome profile from archive |
| `gitmap chrome delete <name>` | `cpd` | Delete profile and stored artifacts |
| `gitmap chrome merge <src> <dst>` | `cpm` | Merge bookmarks, extensions, and preferences |
| `gitmap chrome optimize-projects` | — | Remove duplicate profile references |
| `gitmap chrome reset <profile>` | — | Purge caches and temporary session data |
| `gitmap chrome extensions <profile>` | `ext` | List installed extensions and permissions |
| `gitmap chrome flags <profile>` | `switches` | Inspect configured feature flags |
