# Completed Task: 197-pas-fix

Spec Reference: [02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md](../../../02-spec/21-app/197-gitmap-pas-fix-and-repo-cache-commands.md)

Completed in a 3-wave continuous execution loop using Subagents.

# Subtask 01: Nodes Clone Table Alignment & W3 Resilience

**Objective**: Fix fleet nodes clone table alignment and W3 reachability diagnostics as specified in Spec 194 and Spec 196.

## Requirements
- Fix table column gutter padding in `cli/cmdnodes/nodes_clone_table.go`.
- Ensure W3 offline vs unreachable diagnostics are properly distinguished.
- Prevent ANSI escape code leakage into column measurement calculations.

## Constraints
- Bounding Box: `cli/cmdnodes/*.go`
- Coding Rules: Positive booleans only, AppError wrappers, functions <= 15 lines.
# Subtask 02: PA Async Refactor

**Objective**: Decouple ignore duplicate and tracking checks from synchronous git pull workflow in `gitmap pa`.

## Requirements
- Move `.gitignore` inspection into background async workers in `cli/cmdpull/pull.go`.
- Ensure pull begins immediately without blocking on ignore parsing.
- Aggregate ignore duplicate summaries at the completion of all pulls.

## Constraints
- Bounding Box: `cli/cmdpull/*.go`
- Coding Rules: Positive booleans only, AppError wrappers, functions <= 15 lines.
# Subtask 03: Ignore Group Logic

**Objective**: Implement the new `gitmap ignore` group management commands specified in Spec 197.

## Requirements
- Create commands in `cli/cmdignore/`:
  - `add <string>`: Adds ignore rule to default group.
  - `ls`: List groups, rule count, application rules.
  - `add-group` / `remove-group (rm-grp)`: Manage ignore groups.
  - `set-default-group` / `add-grp-to-default (agtd)`: Set default execution groups.
  - `connect-group-with-repo (cgwp)`: Maps a group to a repo alias/path.
  - `apply`: Applies ignore rules to the current folder/repo.
  - `export` / `import`: JSON-based backup/restore of settings.
- Implement storage in `cli/store/split_db_ignore.go` or similar.

## Constraints
- Bounding Box: `cli/cmdignore/*.go`, `cli/store/split_db_ignore*.go`
- Coding Rules: Positive booleans only, `*appfault.AppError`, functions <= 15 lines. No builds/tests.
# Subtask 04: Cache Split-DB

**Objective**: Implement the high-performance `gitmap cache` SQLite split-db architecture from Spec 197.

## Requirements
- Create commands in `cli/cmdcache/` and storage in `cli/store/split_db_cache.go`.
- `gitmap cache create <path>`:
  - Root DB `sql.db`: repo URL, root files, folder paths, last modified times.
  - Subfolder DB `[slug].db`: relative path and files.
  - Exclude files > 200KB, images, and binaries.
- `gitmap cache search / search-multi / search-multi-grep`:
  - `gitmap cache search "text search" -file-pattern (fp) "a*.md", "b*.md" [--lines 10] [--limit 20]`
  - `gitmap cache search-multi "text search", "multi *" -file-pattern (fp) "a*.md" [--lines 10] [--limit 20]`
  - `gitmap cache search-multi-grep "regex search", "multi *" -file-pattern (fp) "a*.md" [--lines 10] [--limit 20]`
- Reconcile logic: update cache automatically if filesystem `last_modified` is newer.

## Constraints
- Bounding Box: `cli/cmdcache/*.go`, `cli/store/split_db_cache*.go`, `cli/utils/cache*.go`
- Coding Rules: Positive booleans only, `*appfault.AppError`, functions <= 15 lines. No builds/tests.
# Subtask 05: CPAR and FIA Commands

**Objective**: Implement `gitmap commit-push-all-repos (cpar)` and `gitmap fix-ignore-all (fia)` commands per Spec 197.

## Requirements
- `gitmap cpar [-y] [--review (r)] [--commit-only (co)]`: Scans all repos for pending changes, prompts for review/batch commit.
- `gitmap fix-ignore-all (fia) [-y]` and `gitmap fix-ignores-all-ssh (fias) [-y]`: Sweeps all repos for ignore consistency, dedupes ignores, and applies fixes. fias uses the PAS formula (SSH delegation).
- `gitmap see (c)` subcommands: `c commit pending`, `c gitignore issues (ig)`, `c errors ssh (ses)`
- Add commands to `cli/cmdcpar/`, `cli/cmdfixgit/` or `cli/cmdignore/`.

## Constraints
- Bounding Box: `cli/cmdcpar/*.go`, `cli/cmdignore/*fix*.go`, `cli/cmdsee/*.go`
- Coding Rules: Positive booleans only, `*appfault.AppError`, functions <= 15 lines. No builds/tests.
