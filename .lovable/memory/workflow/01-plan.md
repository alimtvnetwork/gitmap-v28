# Development Plan

## Completed Work

### v1.1.0 → v1.1.3
- ✅ Self-update handoff, direct SSH clone output, deploy retry logic
- ✅ Desktop-sync command, enhanced terminal clone hints

### v2.0.0 → v2.1.0
- ✅ Removed GitHub Release integration (Git-only + local metadata)
- ✅ Nested deploy structure, update enhancements, update-cleanup command
- ✅ Generic spec files in `spec/03-general/`

### v2.2.0 → v2.9.0
- ✅ Release-pending, changelog, doctor, latest-branch commands
- ✅ Date formatting, sort/filter flags, CSV/JSON output formats
- ✅ Database with repos, groups, group management commands
- ✅ Self-update hardening (rename-first, stale-process fallback)

### v2.10.0 — Compliance Audit
- ✅ Full compliance audit (Wave 1 + Wave 2): all 75+ source files pass code style rules
- ✅ Trimmed oversized files, fixed negation/switch violations, extracted constants

### v2.11.0 — Constants Inventory
- ✅ Added constants inventory audit section documenting ~280 constants

### v2.12.0 — New Commands
- ✅ `list-versions` (`lv`): show all release tags sorted highest-first with changelog
- ✅ `revert <version>`: checkout tag + handoff rebuild (same mechanism as update)

### v2.13.0 — Changelog Enhancements
- ✅ Release metadata JSON includes changelog field from changelog.md
- ✅ `list-versions` shows changelog notes as sub-points (terminal + JSON)

### v2.14.0 — Go Release Assets, Compression & Checksums
- ✅ `--compress`, `--checksums`, `--no-assets`, `--targets` flags
- ✅ Go cross-compilation pipeline (6 targets, auto-detect, GitHub upload)
- ✅ Config-driven release targets, checksums, and compress booleans

### v2.15.0 — Cross-Platform & CI/CD
- ✅ Full documentation site, `run.sh`, Makefile, GitHub Actions CI/Release

### v2.15.1 — Database Path Fix
- ✅ Fixed DB path resolution: database now at `<binary-dir>/data/`

### v2.16.0 — Interactive TUI
- ✅ Bubble Tea TUI with 6 views, fuzzy search, multi-select

### v2.17.0 → v2.23.0
- ✅ Enhanced group management, gomod, diff-profiles, watch, zip-group, alias commands
- ✅ Shell completion and cross-platform build parity

### v2.24.0 — Release Workflow Restructure
- ✅ Metadata committed on original branch, `--notes`/`--no-commit`/`--skip-meta` flags

### v2.35.0 — Directory Consolidation & ID Migration
- ✅ Consolidated under `.gitmap/`, migrated UUID to INTEGER PK

### v2.36.0 → v2.36.7 — Refactoring & Integration Tests
- ✅ File splits (Wave 1-3), migration hardening, output path fix
- ✅ Integration tests: SkipMeta, rollback, E2E release, edge cases

### v2.49.0 — Polish & Test Coverage
- ✅ Wire `--shell` flag in env commands

### v2.72.0 — VS Code Admin Mode Bypass
- ✅ 3-tier launch strategy, isolated user-data-dir, multi-path discovery

### v2.74.0 — Setup Config & Doctor Checks
- ✅ `gitmap doctor` setup config resolution + shell wrapper detection
- ✅ Shell wrapper scripts export `GITMAP_WRAPPER=1`
- ✅ `gitmap setup` resolves config relative to binary path
- ✅ Post-setup verification step + `gitmap cd` wrapper warnings

### v2.75.0 — Auto-Flatten Clone + Version History DB
- ✅ `gitmap clone-next` flattens by default (base name folder, no `-vN` suffix)
- ✅ `gitmap clone <url>` auto-flattens versioned URLs when no custom folder given
- ✅ New `RepoVersionHistory` SQLite table tracking version transitions
- ✅ `Repos` table gains `CurrentVersionTag` and `CurrentVersionNum` columns
- ✅ Auto-remove existing flattened folder before re-clone (no prompt)
- ✅ `GITMAP_SHELL_HANDOFF` set to flattened path

### v2.76.0 — Version History Command + Specs + ERD
- ✅ New `gitmap version-history` (`vh`) command with `--limit`/`--json`
- ✅ Tab completion for `version-history`/`vh` (Bash, Zsh, PowerShell)
- ✅ Specs `59-clone-next.md` and `87-clone-next-flatten.md` updated for flatten-by-default
- ✅ Full database ERD (Mermaid) covering all 22 tables
- ✅ Docs site page `src/pages/VersionHistory.tsx` with terminal previews
- ✅ Help text `helptext/version-history.md`

## Pending Work

### Unit Tests (from v2.49.0)
- ⬜ Unit tests for `task` commands: create, list, show, delete, validation
- ⬜ Unit tests for `env` commands: set, get, delete, list, path operations
- ⬜ Unit tests for `install` commands: tool validation, manager detection
- ⬜ Unit tests for platform-specific env persistence
- ⬜ Fix `install --check` to print distinct "not found" message
- ⬜ Update docs site command entries with `--shell` flag for env commands
- ⬜ Update `helptext/env.md` examples with `--shell` usage

### Docs Site Navigation
- ⬜ Add `version-history` to sidebar/commands navigation on docs site
- ⬜ Add `clone` page to docs site (currently has clone-next but not clone)

### Install System (from plan.md Parts B–F)
- ⬜ Expand supported tools (databases, package managers)
- ⬜ Multi-platform package manager resolution
- ⬜ `gitmap uninstall` command
- ⬜ Enhanced `--list`, `--status`, `--upgrade` flags
