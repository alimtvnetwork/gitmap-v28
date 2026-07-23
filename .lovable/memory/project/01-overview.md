# Project Overview

## What is gitmap?

gitmap is a portable Go CLI tool that scans directory trees for Git repositories, extracts clone URLs and branch information, and outputs structured data in multiple formats. It can also re-clone repositories from that data, preserving the original folder hierarchy.

## Current Version

**v2.36.3** (defined in `gitmap/constants/constants.go`)

## Tech Stack

- **CLI**: Go (compiled to `gitmap.exe`)
- **Database**: SQLite via `modernc.org/sqlite` (CGo-free)
- **Build/Deploy**: PowerShell (`run.ps1`)
- **Frontend**: React + Vite + Tailwind (documentation site, currently placeholder)
- **Config**: JSON (`powershell.json`, `data/config.json`)

## Key Directories

| Directory | Purpose |
|-----------|---------|
| `gitmap/` | Go source code for the CLI |
| `spec/01-app/` | App-specific specification documents |
| `spec/02-app-issues/` | App issue post-mortems and resolutions |
| `spec/03-general/` | Reusable design patterns & guidelines (generic, shareable) |
| `spec/04-generic-cli/` | Generic CLI implementation blueprint |
| `src/` | React frontend (documentation site) |
| `.lovable/memory/` | AI memory and tracking |
| `.gitmap/release/` | Release metadata JSON files |

## CLI Commands (44 total)

| Command | Alias | Description | Status |
|---------|-------|-------------|--------|
| `scan [dir]` | `s` | Scan directory for Git repos, output all formats, auto-import releases | ✅ Done |
| `clone <source>` | `c` | Re-clone from CSV/JSON/text preserving hierarchy | ✅ Done |
| `pull <name>` | `p` | Pull a specific repo by name | ✅ Done |
| `rescan` | `rs` | Re-run last scan with cached flags | ✅ Done |
| `desktop-sync` | `ds` | Sync repos to GitHub Desktop from scan output | ✅ Done |
| `setup` | — | Configure Git global settings from JSON | ✅ Done |
| `status` | `st` | Show dirty/clean, ahead/behind for all repos | ✅ Done |
| `exec <args>` | `x` | Run any git command across all repos | ✅ Done |
| `release [ver]` | `r` | Create release branch, tag, push, persist to DB, cross-compile Go assets | ✅ Done |
| `release-branch` | `rb` | Complete release from existing branch | ✅ Done |
| `release-pending` | `rp` | Release all pending branches without tags | ✅ Done |
| `changelog [ver]` | `cl` | Show concise release notes, filterable by `--source` | ✅ Done |
| `latest-branch` | `lb` | Find most recently updated remote branch | ✅ Done |
| `list` | `ls` | Show all tracked repos with slugs, labeled fields, inline cd hints | ✅ Done |
| `group <sub>` | `g` | Manage repo groups (create/add/remove/show/list/delete) | ✅ Done |
| `multi-group` | `mg` | Cross-group operations | ✅ Done |
| `list-versions` | `lv` | Show all release tags with changelog, filterable by `--source` | ✅ Done |
| `list-releases` | `lr` | Show stored releases from database, filterable by `--source` | ✅ Done |
| `revert <ver>` | — | Revert to a specific release version | ✅ Done |
| `doctor` | — | Diagnose PATH, deploy, and version issues | ✅ Done |
| `update` | — | Self-update via copy-and-handoff + auto-cleanup | ✅ Done |
| `update-cleanup` | — | Remove update temp files and .old backups | ✅ Done |
| `db-reset` | — | Clear all repos, groups, releases from database | ✅ Done |
| `version` | `v` | Print version string and exit | ✅ Done |
| `help` | — | Show usage information | ✅ Done |
| `seo-write` | `sw` | Automated commit scheduler with SEO-rich messages | ✅ Done |
| `amend` | `am` | Amend commit author metadata | ✅ Done |
| `amend-list` | `al` | List amended commits | ✅ Done |
| `history` | `hi` | Show command execution history | ✅ Done |
| `history-reset` | `hr` | Clear command history | ✅ Done |
| `stats` | — | Show repository statistics | ✅ Done |
| `export` | — | Export repo data | ✅ Done |
| `import` | — | Import repo data | ✅ Done |
| `profile` | — | Manage scan profiles | ✅ Done |
| `diff-profiles` | — | Compare two profiles | ✅ Done |
| `cd` | — | Navigate to repo directory | ✅ Done |
| `watch` | `w` | Watch repos for changes | ✅ Done |
| `bookmark` | `bm` | Save/run bookmarked commands | ✅ Done |
| `gomod` | — | Go module rename operations | ✅ Done |
| `completion` | `comp` | Shell completion for bash/zsh/powershell | ✅ Done |
| `scan-projects` | `sp` | Detect project types (Go, C#, Node, React, C++) | ✅ Done |
| `go-repos` | — | List Go repositories | ✅ Done |
| `csharp-repos` | — | List C# repositories | ✅ Done |
| `audit` | — | Compliance audit | ✅ Done |

## Database Tables (PascalCase)

| Table | Purpose |
|-------|---------|
| `Repos` | Discovered Git repositories |
| `Groups` | Named collections of repos |
| `GroupRepos` | Join table linking repos to groups |
| `Releases` | Release metadata with changelog and source tracking (`release` or `import`) |
| `CommitTemplates` | SEO commit message templates (title/description) with placeholders |

## Release Assets Pipeline (v2.14.0+)

The `release` command includes a full Go cross-compilation pipeline:

| Feature | Description |
|---------|-------------|
| Auto-detect | Detects `go.mod` + `cmd/` entries automatically |
| Cross-compile | Builds for 6 targets: windows/linux/darwin × amd64/arm64 |
| `--compress` | Archives assets (.zip for Windows, .tar.gz for others) |
| `--checksums` | Generates SHA256 `checksums.txt` |
| `--no-assets` | Skips binary compilation |
| `--targets` | Custom OS/arch matrix (e.g., `windows/amd64,linux/arm64`) |
| GitHub Upload | Native HTTP API upload with retry (no `gh` CLI needed) |

## Code Style Rules

- No negation in `if` conditions (no `!`, no `!=`)
- No `switch` statements — use `if`/`else if` chains
- Functions: 8–15 lines
- Files: 100–200 lines max
- One responsibility per package
- Blank line before `return` (unless sole line in `if` block)
- All string literals in `constants` package (no magic strings)
- All DB table/column names in PascalCase

## Version Policy

- **Bump on every code change** that alters behavior or output
- Follows SemVer (`MAJOR.MINOR.PATCH`)
- Displayed in terminal banner, `help`, and `version` command
- `run.ps1` prints version after each build

## File Naming Convention

- All `.md` files use **lowercase-hyphen** naming (e.g. `01-overview.md`, `19-list-versions.md`)
- Go files use lowercase (e.g. `listversions.go`, `revertscript.go`)
