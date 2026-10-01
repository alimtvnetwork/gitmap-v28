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
