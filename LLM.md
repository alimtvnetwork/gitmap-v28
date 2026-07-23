# LLM.md — GitMap AI Reference Guide

> This document is designed for large language models (LLMs) to quickly understand what GitMap is, what it can do, how it's structured, and how to help users with it. Use this as your primary context when assisting with GitMap tasks.

## What is GitMap?

GitMap is a **command-line tool** written in **Go** that scans, catalogs, and manages Git repositories across a developer's machine. Think of it as a "package manager for your Git repos" — it knows where all your repos are, their status, groups, aliases, and can perform batch operations across them.

**Key value proposition:** Developers with many Git repos (10–500+) use GitMap to scan, organize, navigate, clone, release, and maintain them from a single CLI.

## Architecture Overview

```
gitmap (Go CLI binary)
├── cmd/           — Command handlers (one file per command)
├── constants/     — All string literals, error messages, flag names
├── model/         — Data types (ScanRecord, BookmarkRecord, etc.)
├── store/         — SQLite database layer
├── release/       — Release/version parsing and tag management
├── cloner/        — Git clone operations
├── dashboard/     — HTML dashboard generation
├── verbose/       — Debug logging
├── completion/    — Shell completion scripts
├── helptext/      — Embedded Markdown help files (go:embed)
└── scripts/       — Install/uninstall scripts (PowerShell, Bash)

gitmap-updater (separate Go binary)
├── cmd/           — check, run, worker commands
└── main.go        — Entry point
```

**Database:** SQLite (single file at `~/.gitmap/gitmap.db`), managed via `store/` package.

**Platform:** Cross-platform (Windows, Linux, macOS) with amd64 and arm64 support.

**Build:** `go build` or `run.ps1` (Windows PowerShell pipeline).

## Complete Command Reference

### Scanning & Discovery

| Command | Alias | What it does |
|---------|-------|-------------|
| `scan` | `s` | Recursively scan a directory tree for Git repos. Outputs to terminal + CSV + JSON + folder-structure.md |
| `rescan` | `rs` | Re-scan previously scanned directories using cached config |
| `list` | `ls` | Show all tracked repos (filterable by type: `go`, `node`, `react`, `cpp`, `csharp`, `groups`) |
| `go-repos` | `gr` | List Go projects detected by `go.mod` |
| `node-repos` | `nr` | List Node.js projects detected by `package.json` |
| `react-repos` | `rr` | List React projects |
| `cpp-repos` | `cr` | List C++ projects |
| `csharp-repos` | `csr` | List C# projects |

**Example workflow:**
```bash
gitmap scan ~/projects                    # discover all repos
gitmap ls                                 # see what was found
gitmap ls go                              # filter to Go projects
gitmap go-repos --json                    # JSON output for scripting
```

### Navigation

| Command | Alias | What it does |
|---------|-------|-------------|
| `cd` | `go` | Navigate to a repo by slug name. Supports `--pick` for interactive selection |
| `cd repos` | — | Browse all repos interactively |
| `cd set-default` | — | Set a default path for a repo name (skip picker when multiple matches) |
| `cd clear-default` | — | Remove the default path for a repo name |

**Example workflow:**
```bash
gitmap cd my-api                          # jump to repo
gitmap cd repos --group work              # interactive pick from group
gitmap cd set-default my-api D:\repos\v2  # always use this path
gitmap cd -A api                          # navigate via alias
```

### Cloning

| Command | Alias | What it does |
|---------|-------|-------------|
| `clone` | `c` | Clone repos from a scan output file (JSON/CSV/text) |
| `clone-next` | `cn` | Clone the next versioned iteration of the current repo (e.g., v11 → v12) |
| `desktop-sync` | `ds` | Register tracked repos with GitHub Desktop |

**Example workflow:**
```bash
gitmap clone json --target-dir ./restored              # restore from scan
gitmap clone repos.json --ssh-key work                 # clone with SSH key
gitmap cn v++                                          # increment version
gitmap cn v15 --delete                                 # jump to v15, remove old
gitmap cn v++ --create-remote                          # create GitHub repo if missing
```

### Git Operations

| Command | Alias | What it does |
|---------|-------|-------------|
| `pull` | `p` | Pull a specific repo by name (or all in group) |
| `exec` | `x` | Run any git command across all tracked repos |
| `status` | `st` | Show dirty/clean status for all repos |
| `watch` | `w` | Live-refresh status dashboard |
| `has-any-updates` | `hau` | Check if remote has commits you haven't pulled |

**Example workflow:**
```bash
gitmap pull my-api                        # pull one repo
gitmap pull --group work --all            # pull all in group
gitmap exec fetch --prune                 # run across all repos
gitmap watch --interval 10 --group work   # live dashboard
gitmap hau                                # quick remote check
```

### Groups & Aliases

| Command | Alias | What it does |
|---------|-------|-------------|
| `group` | `g` | Create/manage repo groups and activate for batch ops |
| `multi-group` | `mg` | Select multiple groups for batch operations |
| `alias` | `a` | Assign short names to repos (used with `-A` flag on other commands) |

**Example workflow:**
```bash
# Groups
gitmap group create work --desc "Work repos"
gitmap group add work my-api web-app
gitmap g work                             # activate group
gitmap g pull                             # pull all in active group
gitmap g exec fetch --prune               # exec across group
gitmap g clear                            # deactivate

# Multi-group
gitmap mg backend,frontend                # select multiple groups
gitmap mg pull                            # batch pull

# Aliases
gitmap alias set api github/user/api-gateway
gitmap alias suggest --apply              # auto-create aliases
gitmap pull -A api                        # use alias anywhere
gitmap cd -A api
gitmap exec -A api status
```

### Release & Versioning

| Command | Alias | What it does |
|---------|-------|-------------|
| `release` | `r` | Create release: branch, tag, push, cross-compile binaries |
| `release-self` | `rs` | Release gitmap itself from any directory |
| `release-branch` | `rb` | Create a release branch without tagging |
| `release-pending` | `rp` | Show unreleased commits since last tag |
| `changelog` | `cl` | View changelog entries |
| `changelog-generate` | `cg` | Auto-generate changelog from commit messages |
| `list-versions` | `lv` | List all Git release tags |
| `list-releases` | `lr` | List release metadata from database |
| `revert` | — | Revert to a specific release version |
| `clear-release-json` | `crj` | Remove orphaned release metadata files |
| `temp-release` | `tr` | Create lightweight temp release branches |
| `prune` | `pr` | Delete stale release branches that have been tagged |

**Example workflow:**
```bash
gitmap release --bump patch               # auto-bump and release
gitmap release v3.0.0 --bin --compress --checksums -N "Major redesign"
gitmap release --bump minor --dry-run     # preview release
gitmap release v3.0.0 --zip-group docs-bundle   # include archives
gitmap release-pending                    # what's unreleased?
gitmap changelog-generate --write         # auto-generate changelog
gitmap revert v2.48.0                     # rollback
gitmap tr 10 v1.$$ -s 5                   # temp branches
gitmap prune                              # clean stale branches
```

### SSH Key Management

| Command | Alias | What it does |
|---------|-------|-------------|
| `ssh` | — | Generate, list, display, and delete SSH keys for Git auth |

**Example workflow:**
```bash
gitmap ssh --name work --path ~/.ssh/id_rsa_work
gitmap ssh cat --name work                # show public key
gitmap ssh list                           # list all keys
gitmap ssh config                         # regenerate ~/.ssh/config
gitmap clone repos.json --ssh-key work    # clone with specific key
```

### Zip Groups (Release Archives)

| Command | Alias | What it does |
|---------|-------|-------------|
| `zip-group` | `z` | Manage named collections of files bundled into ZIP during releases |

**Example workflow:**
```bash
gitmap z create "chrome extension" chrome-extension/dist
gitmap z add docs-bundle ./README.md ./CHANGELOG.md ./docs/
gitmap z create extras --archive extra-files.zip
gitmap z list                             # list all groups
gitmap z show docs-bundle                 # show items with expansion
gitmap release v3.0.0 --zip-group docs-bundle
```

### Data, Profiles & Bookmarks

| Command | Alias | What it does |
|---------|-------|-------------|
| `export` | `ex` | Export database to file |
| `import` | `im` | Import repos from file |
| `profile` | `pf` | Manage database profiles (create, switch, list) |
| `bookmark` | `bk` | Save commands as named bookmarks and replay them |
| `db-reset` | — | Reset the SQLite database (requires `--confirm`) |
| `history` | `hi` | Show CLI command execution history |
| `history-reset` | `hr` | Clear command history |
| `stats` | `ss` | Show aggregated usage and performance metrics |

**Example workflow:**
```bash
gitmap export                             # export DB
gitmap import gitmap-export.json          # import repos
gitmap profile create work                # new profile
gitmap profile switch work                # activate profile
gitmap bookmark save daily scan ~/projects
gitmap bookmark run daily                 # replay saved command
gitmap history --limit 10                 # recent commands
gitmap stats --json                       # usage stats
gitmap db-reset --confirm                 # nuclear option
```

### Author Amendment

| Command | Alias | What it does |
|---------|-------|-------------|
| `amend` | `am` | Rewrite commit author name/email across commits |
| `amend-list` | `al` | List previous amendments |

**Example workflow:**
```bash
gitmap amend --name "John" --email "john@co.com"           # all commits
gitmap amend abc123 --name "John" --email "john@co.com"    # specific commit
gitmap amend --name "John" --email "john@co.com" --dry-run # preview
gitmap amend --name "John" --email "john@co.com" --force-push
gitmap amend-list --json --limit 5
```

### Environment & Tools

| Command | Alias | What it does |
|---------|-------|-------------|
| `env` | `ev` | Manage persistent environment variables and PATH entries |
| `install` | `in` | Install developer tools via platform package manager |

**Supported tools:** Node.js, Yarn, Bun, pnpm, Go, Python, VS Code, Git, Git LFS, GitHub CLI, GitHub Desktop, PHP, C++ (MinGW), PowerShell.

**Example workflow:**
```bash
gitmap env set GOPATH "/home/user/go"
gitmap env path add /usr/local/go/bin
gitmap install node                       # install Node.js
gitmap install go --check                 # check if installed
gitmap install python --dry-run           # preview command
gitmap install --list                     # show all tools
```

### File-Sync Tasks

| Command | Alias | What it does |
|---------|-------|-------------|
| `task` | `tk` | Create and run one-way folder synchronization tasks |

**Example workflow:**
```bash
gitmap task create my-sync --src ./src --dest ./backup
gitmap tk run my-sync --interval 10 --verbose
gitmap task list
gitmap task delete my-sync
```

### Utilities

| Command | Alias | What it does |
|---------|-------|-------------|
| `setup` | — | Interactive first-time configuration wizard |
| `doctor` | — | Diagnose PATH, deploy, and version issues |
| `update` | — | Self-update (from source or via gitmap-updater fallback) |
| `version` | `v` | Show version number |
| `completion` | `cmp` | Generate shell tab-completion (PowerShell, Bash, Zsh) |
| `interactive` | `i` | Full-screen interactive TUI |
| `docs` | `d` | Open documentation website in browser |
| `dashboard` | `db` | Generate interactive HTML dashboard for a repo |
| `gomod` | `gm` | Rename Go module path across repo |
| `seo-write` | `sw` | Auto-commit SEO messages from CSV |
| `diff-profiles` | `dp` | Compare repos across two scan profiles |

### Global Flags

These flags work with most commands:

| Flag | Description |
|------|-------------|
| `--help`, `-h` | Show help text for any command |
| `--verbose` | Enable debug logging |
| `--json` | JSON output (where supported) |
| `-A`, `--alias` | Use a repo alias instead of slug |

## Coding Conventions

When modifying GitMap code, follow these rules:

1. **No magic strings** — All literals go in `constants/` package
2. **Functions ≤ 15 lines** — Extract helpers liberally
3. **Files ≤ 200 lines** — Split when approaching limit
4. **PascalCase** for exported constants
5. **`is`/`has` prefix** for boolean variables and functions
6. **Blank line before `return`** statements
7. **Chained `if` + `return`** for dispatch (not switch)
8. **Use `fmt.Fprint`** (not `Fprintln`) when the constant already ends with `\n`
9. **Group same-type parameters** — `func(a, b bool)` not `func(a bool, b bool)`
10. **Positive logic** in `if` conditions

## Project Structure

```
/
├── gitmap/                    # Main CLI (Go module)
│   ├── cmd/                   # Command handlers
│   ├── constants/             # All string constants
│   ├── model/                 # Data types
│   ├── store/                 # SQLite database
│   ├── release/               # Version/tag management
│   ├── cloner/                # Clone operations
│   ├── dashboard/             # HTML dashboard
│   ├── verbose/               # Debug logging
│   ├── completion/            # Shell completions
│   ├── helptext/              # Embedded help (go:embed)
│   └── scripts/               # Install/uninstall scripts
├── gitmap-updater/            # Standalone updater (Go module)
├── spec/                      # Specifications & design docs
│   ├── 01-app/                # Feature specs
│   ├── 02-app-issues/         # Post-mortems
│   ├── 03-general/            # Architecture patterns
│   ├── 04-generic-cli/        # CLI blueprints
│   ├── 05-coding-guidelines/  # Quality standards
│   └── 06-design-system/      # UI standards
├── src/                       # React docs site
├── hooks/                     # Git hooks (pre-commit lint)
├── .github/workflows/         # CI/CD pipelines
├── CHANGELOG.md
├── README.md
└── LLM.md                    # This file
```

## Database Schema (Conceptual)

GitMap uses SQLite with tables for:

- **repos** — scanned repository records (path, slug, URLs, branch, type)
- **groups** — named groups of repos
- **group_members** — repo-to-group mappings
- **aliases** — short names pointing to repo slugs
- **bookmarks** — saved command configurations
- **amendments** — author amendment audit trail
- **releases** — release metadata
- **ssh_keys** — SSH key records
- **zip_groups** / **zip_group_items** — release archive configurations
- **history** — command execution history
- **tasks** — file-sync task definitions

## Installation Methods

| Method | Command |
|--------|---------|
| One-liner (Windows) | `irm https://raw.githubusercontent.com/.../install.ps1 \| iex` |
| One-liner (Unix) | `curl -fsSL https://raw.githubusercontent.com/.../install.sh \| sh` |
| Custom directory | `.\install.ps1 -InstallDir "D:\tools\gitmap"` |
| Pinned version | `.\install.ps1 -Version v2.49.1` |
| From source | `cd gitmap && go build -o ../gitmap .` |
| Future: Chocolatey | `choco install gitmap` (planned) |
| Future: Winget | `winget install AliMTVNetwork.GitMap` (planned) |

## Self-Update Flow

1. `gitmap update` checks for embedded repo path
2. If found: pulls latest source and rebuilds via PowerShell
3. If not found: looks for `gitmap-updater` on PATH
4. `gitmap-updater` queries GitHub Releases API, downloads `install.ps1`, runs it
5. Uses handoff pattern (temp binary copy) to avoid Windows file locks
6. Manual override: `gitmap update --repo-path C:\gitmap-src`

## CI/CD Pipeline

- **CI:** golangci-lint (28 linters) + go test + go vet on every push/PR
- **Vulncheck:** Weekly govulncheck scan (Mondays 09:00 UTC)
- **Release:** Push `v*` tag → cross-compile 6 targets → checksums → GitHub Release with assets
- **Targets:** windows/amd64, windows/arm64, linux/amd64, linux/arm64, darwin/amd64, darwin/arm64

## Common Patterns for LLM Assistance

### "I want to find/navigate to a repo"
```bash
gitmap cd <repo-name>
gitmap cd repos                    # interactive picker
gitmap cd -A <alias>               # via alias
```

### "I want to update all my repos"
```bash
gitmap g work                      # activate group
gitmap g pull                      # pull all
# or
gitmap exec pull                   # pull across ALL repos
```

### "I want to organize my repos"
```bash
gitmap scan ~/projects             # discover repos
gitmap group create work           # create group
gitmap group add work api web      # add repos
gitmap alias suggest --apply       # auto-name repos
```

### "I want to release my project"
```bash
gitmap release-pending             # what's unreleased?
gitmap changelog-generate --write  # generate changelog
gitmap release --bump patch --bin  # release with binaries
```

### "I want to clone a project iteration"
```bash
gitmap cn v++                      # next version
gitmap cn v++ --delete             # and remove old
```

### "I want to check repo health"
```bash
gitmap doctor                      # diagnose issues
gitmap hau                         # check for updates
gitmap status                      # dirty/clean status
gitmap watch                       # live dashboard
```
