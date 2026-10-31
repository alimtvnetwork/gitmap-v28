# Subtask 02: Dedicated `gitmap agy summaries` CLI & Built-in `agy` DB Alias

- **Parent Plan:** [368-sql-query-execution-and-gitmap-skill-enhancement.md](../../pending/368-sql-query-execution-and-gitmap-skill-enhancement.md)
- **Slug:** `sql-query-execution-and-gitmap-skill-enhancement - Task 02`
- **Status:** `pending`
- **Assignee:** V7 Master Autonomous Execution Engine

## Action Items
1. Create `cli/cmdagy/agy_summaries.go` implementing `RunAgySummariesCLI(args []string) error`.
2. Connect `gitmap agy summaries`, `gitmap agy convs`, and `gitmap agy db` to automatically read `~/.gemini/antigravity/conversation_summaries.db`.
3. Support `-l, --limit`, `-j, --json`, and `-q, --query` flags.
4. Register the built-in `agy` database alias in `cli/cmddb/db_read.go` (and `cli/cmddb/` alias table) so `gitmap db-read agy "<sql>"` resolves to `conversation_summaries.db`.
5. Add unit tests in `cli/cmdagy/agy_summaries_test.go` and `cli/cmddb/db_read_test.go`.
