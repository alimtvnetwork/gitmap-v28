# AI Code Execution & Behavioral Restrictions

This rule is mandatory across all AI agents (Antigravity, Cursor, Claude, Muse, Gemini) operating in GitMap-managed repositories. Bypassing these restrictions breaks audit trails, causes multi-agent write collisions, and violates repository standards.

---

## 1. Strictly Prohibited Code Patterns & Mandatory Replacements

### 1.1 Total Ban on Ad-Hoc Python Database Scripts
- **❌ BANNED:** Writing scratch Python scripts (e.g. `check_summaries.py`, `query.py`, `brain/.../scratch/*.py`) or executing inline Python one-liners (`python -c "import sqlite3..."`, `gitmap py -c "import sqlite3..."`) to query, inspect, or mutate SQLite databases.
- **✅ MANDATORY REPLACEMENTS:**
  | Target Database / Operation | Mandatory GitMap Command |
  | :--- | :--- |
  | Antigravity Conversation States | `gitmap agy summaries [N] [-q query] [-a] [-j]` (aliases: `gitmap agy convs`, `gitmap convs`, `gitmap agy db`) |
  | SQL Query on Antigravity DB | `gitmap db-read agy "SELECT conversation_id, title, status FROM conversation_summaries LIMIT 10"` |
  | Read-only SQL on Any Database | `gitmap db-read <alias\|path> "<sql>" [-l limit] [--json]` (alias: `gitmap dbr`) |
  | List Tables & Row Counts | `gitmap db-read <alias\|path> ls` |
  | Inspect Table Schema & Columns | `gitmap db-read <alias\|path> table <table> --columns` |
  | Pattern Search across Columns | `gitmap db-find <alias\|path> "<pattern>"` (alias: `gitmap dbf`) |
  | Task Swarm Tracking | `gitmap agent task ls`, `gitmap agent subtask ls <task>` |

---

### 1.2 Mandatory Heatmap Touch Notification During File Analysis
- **❌ BANNED:** Reading, inspecting, or analyzing source code files and specifications during troubleshooting or planning without registering intent with the GitMap Heatmap system.
- **✅ MANDATORY REQUIREMENT:**
  Whenever an AI agent begins analyzing a target file, it MUST execute:
  ```bash
  gitmap heatmap touch <filepath> --reason "<why the file is being analyzed>"
  ```
  *(Alternatively, run `gitmap cat <filepath> --touch` to read and register in a single stream).*
  This logs active file leases and inspection intent to `heatmap.db`, eliminating concurrency collisions across agent swarms.

---

### 1.3 Total Ban on Ad-Hoc Python File Reading & In-Memory Slicing
- **❌ BANNED:** `gitmap py -c "open('...').read().splitlines()"` or `lines=open(...).readlines()`.
- **✅ MANDATORY REPLACEMENT:** `gitmap cat <filepath> [-L start:end] [-q query] [-t tail] [--head N] [--num] [--raw]` (alias: `gitmap ct`).

---

### 1.4 Total Ban on Ad-Hoc Python Filesystem Traversal
- **❌ BANNED:** `gitmap py -c "import os; [print(...) for ... in os.walk(...) if ...]"`.
- **✅ MANDATORY REPLACEMENT:** `gitmap find "*<pattern>*" [dir]` or `gitmap aum search "<pattern>" [dir]`.

---

### 1.5 Total Ban on Ad-Hoc Python Sockets for Chrome DevTools / CDP
- **❌ BANNED:** Writing raw socket scripts (`socket.socket`, `send_frame`, `recv_frame`) to probe DevTools pages or execute JavaScript.
- **✅ MANDATORY REPLACEMENT:** `gitmap cdp probe`, `gitmap cdp eval "<expr>"`, `gitmap cdp target`.

---

## 2. Enforcement & Audit
1. `gitmap py` actively intercepts prohibited code signatures and displays educational command tips on `stderr`.
2. All file touches and database operations are captured in Split-DB audit logs (`gitmap.db`, `heatmap.db`, `agent-task.db`).
