# Specification: Dedicated AGY Database Commands & `agy` Alias

## 1. Overview
Antigravity stores its conversation summaries, titles, execution status, and workspace URIs in SQLite at `~/.gemini/antigravity/conversation_summaries.db`.
This specification introduces high-speed native commands in `gitmap agy` and registers the `agy` alias in `gitmap db-read`.

## 2. CLI Command Contracts

### 2.1 `gitmap agy summaries` (aliases: `gitmap agy convs`, `gitmap agy db`)
- **Usage:**
  - `gitmap agy summaries [N] [flags]`
  - `gitmap agy convs [N] [flags]`
  - `gitmap agy db [N] [flags]`
- **Arguments & Flags:**
  - `N`: Maximum number of conversations to display (default: 10).
  - `-l, --limit <N>`: Explicit row limit.
  - `-j, --json`: Return structured JSON array.
  - `-q, --query <pattern>`: Filter conversations by title, status, or conversation ID.
  - `--active-only`: Filter for `CASCADE_RUN_STATUS_RUNNING` or `not_fully_idle = 1`.
- **Output:**
  A compact, Tokyo Night styled table with columns:
  `ID | Title | Status | Idle? | Last Modified | Workspace`

### 2.2 Built-in `agy` Alias for `gitmap db-read`
- **Contract:**
  When `gitmap db-read agy "<sql>"` is invoked, GitMap automatically resolves `agy` to:
  `filepath.Join(userHome, ".gemini", "antigravity", "conversation_summaries.db")`.
- **Examples:**
  - `gitmap db-read agy "SELECT conversation_id, title FROM conversation_summaries LIMIT 5"`
  - `gitmap db-read agy ls`
  - `gitmap db-read agy table conversation_summaries --columns`

## 3. Implementation Plan
- File: `cli/cmdagy/agy_summaries.go`
- Handler: `RunAgySummariesCLI(args []string) error`
- Route: Registered in `cli/cmdagy/agy_cmd.go` and `cli/cmd/root.go`.
- Alias Resolver: Updated in `cli/cmddb/db_read.go` (or `cli/store/`) to map `"agy"` to `conversation_summaries.db`.
