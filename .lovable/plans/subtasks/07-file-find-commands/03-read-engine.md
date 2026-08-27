# Subtask 3: Read Engine

1. In `gitmap/searcher/finder.go`:
   - Implement `FindAndRead(ctx, db, query, isRegex, limit) []FileReadResult`
   - If `IsBig == 0`, use `Content` from Repo DB.
   - If `IsBig == 1`, read from `AbsolutePath` on disk on-the-fly.
   - Map into a JSON-compatible struct: `{ "file_path": "...", "content": "..." }`.
- [x] Task incomplete
