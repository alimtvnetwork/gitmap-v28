# Subtask 94.04: Go Fast Cache Optimization, Upfront Batch Indexing & Binary Sniffer

## Goal
Optimize Go's file indexing and caching in `gitmap/indexer/walker.go` and `gitmap/searcher/db_search.go` by adopting Python's fast cache patterns: batch upfront timestamp preloading, 8KB null-byte binary probe, and comprehensive directory pruning.

## Files Impacted
- `gitmap/indexer/walker.go`
- `gitmap/searcher/db_search.go`
- `gitmap/fsutil/discovery_cache.go`

## Acceptance Criteria
1. `indexer.Walker`: replaces per-file `db.QueryRowContext` with an upfront batch load of `RelativePath` -> `WriteTime` into an in-memory `map[string]int64`, eliminating thousands of redundant SQLite queries during traversal.
2. `indexer.Walker`: adds an 8KB null-byte sniffer (`bytes.IndexByte(head, 0) != -1`) before storing text into `RepoFile`, marking binary files with `IsBig = true` (or skipping content) to prevent database corruption.
3. Directory pruning: updates directory ignore rules to prune `.venv`, `dist`, `build`, `bin`, `vendor`, `.gemini`, `coverage`, `tmp`, `__pycache__`, and `.turbo`.
4. All functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
