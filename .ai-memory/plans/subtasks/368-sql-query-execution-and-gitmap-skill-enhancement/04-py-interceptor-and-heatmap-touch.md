# Subtask 04: GitMap Py Interceptor Expansion & Heatmap File Touch Logging

- **Parent Plan:** [368-sql-query-execution-and-gitmap-skill-enhancement.md](../../pending/368-sql-query-execution-and-gitmap-skill-enhancement.md)
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 04`
- **Status:** `pending`
- **Assignee:** V7 Master Autonomous Execution Engine

## Action Items
1. In `cli/cmdpy/py_cmd.go`, expand `warnAdHocAntiPatterns` to detect:
   - `conversation_summaries.db` -> emit tip to use `gitmap agy summaries` or `gitmap db-read agy "<sql>"`.
   - `sqlite3.connect` on general `.db` files -> emit tip to use `gitmap db-read <db> "<sql>"`.
   - Scratch scripts in `brain/.../scratch/` trying to query databases.
2. In `cli/cmdheatmap/` (or `cli/cmdheatmap/touch.go`), ensure `gitmap heatmap touch <filepath> --reason "<why>"` registers a read/inspection touch event into `repodb/heatmap.db`.
3. Add unit tests in `cli/cmdpy/py_cmd_test.go` and `cli/cmdheatmap/touch_test.go`.
