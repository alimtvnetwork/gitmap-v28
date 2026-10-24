# Master Plan: 324 — Package Constants & Enums DRY Restructuring (Phase 6: Repository Footprint, LFS & Tracking CLI Packages)

## Executive Summary
This milestone executes the strict `01-prompts/15-cg-execute/06-constants-and-enums.md` coding guideline protocol across 5 central repository footprint, LFS, and tracking packages in `cli/`:
1. `cli/cmdflatten` (Repository Flattening & Hierarchy Normalization — 2 files)
2. `cli/cmdlfs` (Large File Storage Conversion & Execution — 2 files)
3. `cli/cmdgitrm` (Git Cached & Tracked Object Removal — 2 files)
4. `cli/cmdsize` (Repository Disk Usage & Storage Sizing — 3 files)
5. `cli/cmdurl` (Remote Git URL Parsing & Browser Navigation — 2 files)

Each package will receive:
- A dedicated leaf package `enums/` with strongly-typed enums ending in `Type`, deterministic methods (`IsValid()`, `String()`, `Parse*()`, `All*()`), positive booleans (`is*`, `has*`), and comprehensive unit tests (`enums_test.go`).
- A dedicated leaf package `constants/` with zero magic strings, centralized command identifiers, flags, prompt tokens, exit codes, error messages, and unit tests (`constants_test.go`).
- Refactored caller files in each parent package eliminating magic string comparisons, raw literals, and repetitive token definitions.

## Target Packages & Boundaries
- `cli/cmdflatten`: `enums/flatten_enums.go`, `constants/constants.go`, refactored callers in `flatten_cmd.go`, `flatten_copy.go`.
- `cli/cmdlfs`: `enums/lfs_enums.go`, `constants/constants.go`, refactored callers in `lfs_convert.go`, `lfs_exec.go`.
- `cli/cmdgitrm`: `enums/gitrm_enums.go`, `constants/constants.go`, refactored callers in `gitrm_entry.go`, `hook_shim.go`.
- `cli/cmdsize`: `enums/size_enums.go`, `constants/constants.go`, refactored callers in `size.go`, `helpers_shim.go`, `hook_shim.go`.
- `cli/cmdurl`: `enums/url_enums.go`, `constants/constants.go`, refactored callers in `cmd_url.go`, `cmd_url_test.go`.

## Non-Negotiable Invariants
1. Enums must end with `Type` suffix.
2. Single-responsibility leaf packages: `enums/` and `constants/` must have zero cyclic imports back to the parent package.
3. Positive boolean predicates only (`isValid`, `hasValue`, `isDefined`).
4. Functions <= 15 lines (<= 8 preferred). Blank line after `}` and before `return`.
5. Error sovereignty: leaf helpers return `error` or `*appfault.AppError` directly. Outer callers handle reporting.
6. Zero raw git commits or `git pull`. Final atomic commit via GitMap (`gitmap cfr`).
