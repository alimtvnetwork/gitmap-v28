# Subtask 05: GitMap Split-DB Repository Cache & Search Engine

## Objectives
- Implement and verify `gitmap cache` subsystem:
  - `create <path>` / `add <path>`: create root `sql.db` and `<slug>.db` per top-level folder.
  - `ls`: list cached repositories and files.
  - `rm <path>` / `remove <path>`: delete cached entries.
  - `search <query> [glob]`: text search with context lines and limits.
  - `search-multi <queries...>`: multi-term search with file patterns (`-fp`).
  - `search-multi-grep <regex>`: SQLite-native regular expression search.
  - `reconcile` / `recache` / `sync`: asynchronous background update based on filesystem `mtime`.
- Exclusion gates: skip files > 200KB, binary files, `.git/`, `.vscode/`, `.idea/`, `node_modules/`, `vendor/`.
- Never calculate SHA hashes during scan; rely on `mtime`.

## Target Files
- `cli/cmdcache/cache_cli.go`
- `cli/cmdcache/cache_create.go`
- `cli/cmdcache/cache_search.go`
- `cli/store/split_db_cache.go`
