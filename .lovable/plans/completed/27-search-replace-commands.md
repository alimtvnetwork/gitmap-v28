# 08 Search and Replace Commands

## 1. Problem Context

The CLI endpoints for content-level searching (`search`, `replace`, `replace-regex`, `repo-search`, `repo-regex`, `repo-search-json`, `repo-search-regex-json`, `search-replace-all`) are currently stubbed in `search_entry.go`. We must fully wire them up to the `searcher.SearchRepoDB` and `searcher.SearchRepoDBRegex` engines. We must also update the help documentation for these commands to provide real-world regex examples (e.g., searching for `func run[A-Z]`).

## 2. Core Requirements

- **On-the-fly search:** `search`, `replace`, `replace-regex`. Will walk the filesystem or query the SQLite DB without caching (`useCache=false`).
- **DB Cached Search:** `repo-search`, `repo-regex`, `repo-search-json`. Will query SQLite DB with caching enabled (`useCache=true`).
- **Reset:** `search-replace-all reset` clears the split DBs.
- **Help Files:** Update `search.md`, `replace.md`, `repo-search.md`, etc., with examples of searching for functions.
- **UI:** Format outputs using `pterm` (similar to `find_entry.go`), printing line matches and character positions.

## 3. Subtasks Execution

1. **01-helptext.md**: Overwrite dummy help text for search commands with detailed Markdown examples showing function searching.
2. **02-search-wiring.md**: Implement `search_entry.go` fully by importing `store`, `repodb`, and `searcher`. Replace the stub functions with real `SearchRepoDB` / `SearchRepoDBRegex` invocations.
3. **03-replace-engine.md**: Implement replace stubs.
4. **04-terminal-examples.md**: Execute test runs of the CLI directly in the chat to prove functionality.
