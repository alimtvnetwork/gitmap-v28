# Subtask 2: Find Engine

1. In `gitmap/searcher/finder.go` (new file):
   - Implement `FindFile(ctx, db, query, isRegex, limit)`
   - Map exact, starts-with, ends-with, contains into SQLite `LIKE` clauses on `RelativePath`.
   - Implement lazy regex evaluation on `RelativePath` for regex searches.
   - Utilize `SearchCache` table for frequency tracking exactly like `SearchRepoDB`.
- [x] Task incomplete
