# Project Overview

## What is gitmap?

**gitmap** is a portable Go CLI tool that scans directory trees for Git repositories, extracts clone URLs and branch information, and outputs structured data in multiple formats. It can re-clone repositories from that data, preserving the original folder hierarchy. It also manages releases, SSH keys, environment variables, developer tool installations, and provides an interactive TUI.

## Current Version

**v3.1.0** (defined in `gitmap/constants/constants.go`)

## Tech Stack

| Layer | Technology |
|-------|-----------|
| CLI | Go (compiled to `gitmap` / `gitmap.exe`) |
| Database | SQLite via `modernc.org/sqlite` (CGo-free) |
| Build/Deploy | PowerShell (`run.ps1`), Makefile, `run.sh` |
| Frontend | React + Vite + Tailwind (documentation site) |
| Config | JSON (`data/config.json`) |
| CI/CD | GitHub Actions |

## Repository

`https://github.com/alimtvnetwork/gitmap-v28`

## Key Directories

| Directory | Purpose |
|-----------|---------|
| `gitmap/` | Go source code for the CLI |
| `gitmap-updater/` | Standalone updater binary |
| `spec/01-app/` | App-specific specification documents |
| `spec/02-app-issues/` | App issue post-mortems and resolutions |
| `spec/03-general/` | Reusable design patterns and guidelines |
| `spec/04-generic-cli/` | Generic CLI implementation blueprint |
| `spec/05-coding-guidelines/` | Code quality rules and conventions |
| `spec/06-design-system/` | UI design system specs |
| `spec/09-pipeline/` | CI/CD pipeline specifications |
| `src/` | React frontend (documentation site) |
| `.lovable/memory/` | AI memory and tracking |
| `.lovable/prompts/` | AI onboarding prompts |
| `.gitmap/release/` | Release metadata JSON files (DO NOT TOUCH) |
| `.gitmap/output/` | Scan output (CSV, JSON, scripts) |
| `settings/` | Editor and terminal settings sync |

## CLI Commands (60+)

The CLI supports 60+ subcommands with aliases. Key commands include: `scan`, `clone`, `clone-next`, `pull`, `release`, `release-self`, `update`, `install`, `ssh`, `env`, `interactive`, `cd`, `watch`, `task`, `doctor`, `version`, `changelog`, `stats`, `export`, `import`, `profile`, `completion`.

Full command list: `.lovable/memory/features/cli-commands.md`

## Database

SQLite with 22+ tables using strict PascalCase naming and INTEGER PRIMARY KEY AUTOINCREMENT. Connection pooling restricted to `SetMaxOpenConns(1)`. Database anchored to binary execution path via `filepath.EvalSymlinks`.

## Code Style Summary

- Files: max 200 lines
- Functions: 8-15 lines
- No negation in `if` conditions (no `!`, no `!=`)
- No `switch` statements
- No magic strings — all literals in `constants` package
- PascalCase for DB tables/columns and exported constants
- Boolean naming: always `is`/`has` prefix
- Blank line before `return`
- Zero-swallow error policy

## Version Policy

Bump on every code change. SemVer (`MAJOR.MINOR.PATCH`).
