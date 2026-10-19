# CLI Command Contracts: gitmap db-read, db-add, db-edit

- **Spec ID:** `02-spec/21-app/command-and-database-integration-instructions/02-cli-command-contracts.md`
- **Module:** `cli/cmddb`, `cli/cmd`, `cli/store`
- **Status:** APPROVED
- **Version:** 1.0.0
- **Target Audience:** Core CLI Engineers, Autonomous AI Agents, Multi-Agent Orchestrators
- **Target Release:** Minor Bump

---

## 1. Command Overview & Aliases

GitMap exposes three primary database commands with corresponding short aliases:

| Canonical Command | Short Alias | Primary Function | Safety Guarantee |
|---|---|---|---|
| `gitmap db-read` | `gitmap dbr` | Query tables, list schema, execute read-only SQL | Strict read-only (`PRAGMA query_only = ON`) |
| `gitmap db-add` | `gitmap dba` | Insert records via CSV values or interactive TUI | Atomic transaction (`BEGIN IMMEDIATE`) |
| `gitmap db-edit` | `gitmap dbe` | Update existing records by primary key ID | Atomic update (`BEGIN IMMEDIATE`) |

All subcommands inherit global database resolution (`-as, --alias`), output formatting (`--json`), and verbose logging flags.

---

## 2. gitmap db-read / gitmap dbr Contract

The `db-read` command inspects and queries database tables across any Split-DB tier.

### 2.1 Execution Modes

1. **SQL Query Mode:** Executes a raw read-only SQL statement.
   ```bash
   gitmap db-read "SELECT id, repoName, remoteUrl FROM Repo WHERE isClean = 0"
   gitmap dbr "SELECT count(*) AS total FROM Node" -as install
   ```
2. **Table Listing Mode (`ls`):** Lists all user tables, row counts, and column counts.
   ```bash
   gitmap db-read ls
   gitmap dbr ls -as pipeline --json
   ```
3. **Table Query Mode (`table`):** Selects records from a named table with optional limits and column projections.
   ```bash
   gitmap db-read table Repo -c "id,repoName" -l 10
   gitmap dbr table PullError -as pipeline -l 5 --json
   ```

### 2.2 db-read Argument Parsing Table

| Argument / Flag | Shorthand | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `<query>` / `<mode>` | N/A | string | Yes | N/A | SQL query string, `ls` keyword, or `table <name>` |
| `--alias` | `-as` | string | No | `master` | Target DB alias (`master`, `install`, `pipeline`, `heatmap`, `tasks`) or file path |
| `--limit` | `-l` | int | No | `50` | Maximum rows to return (`0` = all rows) |
| `--columns` | `-c` | string | No | `*` | Comma-delimited list of columns to project |
| `--where` | `-w` | string | No | `""` | Optional WHERE filter clause (table mode only) |
| `--order` | `-o` | string | No | `""` | Optional ORDER BY column and direction (e.g. `id DESC`) |
| `--json` | N/A | bool | No | `false` | Output structured machine-readable JSON |

---

## 3. gitmap db-add / gitmap dba Contract

The `db-add` command inserts new records into any user table, supporting scripted CSV ingestion or interactive prompts.

### 3.1 Execution Modes

1. **CSV Value Insertion:** Ingests values directly from CLI flags.
   ```bash
   gitmap db-add -t Repo -c "repoName,absolutePath,remoteUrl" -v "gitmap,/work/gitmap,https://github.com/org/gitmap"
   gitmap dba -t Node -v "alpha-worker,192.168.1.50,22" -as install
   ```
2. **Interactive TUI Prompt Mode:** Introspects table schema via `PRAGMA table_info` and prompts for non-autoincrement fields.
   ```bash
   gitmap db-add -t Bookmark -i -as install
   ```

### 3.2 db-add Argument Parsing Table

| Argument / Flag | Shorthand | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `--table` | `-t` | string | Yes | N/A | Target table name |
| `--values` | `-v` | string | No* | `""` | Comma-delimited column values (*required unless `--interactive`) |
| `--columns` | `-c` | string | No | `""` | Comma-delimited target column names (inferred if omitted) |
| `--interactive` | `-i` | bool | No | `false` | Prompt interactively for each column value |
| `--alias` | `-as` | string | No | `master` | Target DB alias or explicit database file path |
| `--json` | N/A | bool | No | `false` | Output insertion metadata in JSON |

---

## 4. gitmap db-edit / gitmap dbe Contract

The `db-edit` command updates existing records identified by their primary key `id`.

### 4.1 Execution Syntax

```bash
gitmap db-edit -t Repo -id 12 -v "isClean=1,remoteUrl=https://github.com/new/url"
gitmap dbe -t Node -id 3 -v "isOnline=1,lastPingMs=14" -as install --json
```

### 4.2 db-edit Argument Parsing Table

| Argument / Flag | Shorthand | Type | Required | Default | Description |
|---|---|---|---|---|---|
| `--table` | `-t` | string | Yes | N/A | Target table name |
| `--id` | `-id` | int64 | Yes | N/A | Primary key identifier of target row |
| `--values` | `-v` | string | Yes | N/A | Comma-delimited `col=val` pairs or ordered values |
| `--alias` | `-as` | string | No | `master` | Target DB alias or explicit database file path |
| `--json` | N/A | bool | No | `false` | Output update confirmation in JSON |

---

## 5. Standardized Error Codes

All database subcommands return deterministic `AppError` exit codes:

| Error Code | Exit Code | Trigger Condition | Remediation Guidance |
|---|---|---|---|
| `ERR_DB_NOT_FOUND` | 20 | Resolved SQLite file does not exist on disk. | Verify database alias or run `gitmap scan` to initialize. |
| `ERR_DB_LOCKED` | 21 | Busy timeout exceeded (5000ms lock contention). | Check for dangling writer processes or hung background agents. |
| `ERR_TABLE_NOT_FOUND` | 22 | Requested table does not exist in schema. | Execute `gitmap db-read ls` to inspect existing tables. |
| `ERR_INVALID_SQL` | 23 | SQL statement contains syntax errors. | Verify SQL syntax against standard SQLite grammar. |
| `ERR_READONLY_VIOLATION` | 24 | Mutation attempted during `db-read` operation. | Use `gitmap db-add` or `gitmap db-edit` for write operations. |
| `ERR_COLUMN_MISMATCH` | 25 | Values count does not match columns count. | Align values count with table schema or supply explicit `-c`. |
| `ERR_RECORD_NOT_FOUND` | 26 | Specified primary key ID not found during edit. | Query table via `gitmap db-read table <T>` to confirm row ID. |
| `ERR_CONSTRAINT_VIOLATION` | 27 | UNIQUE, FOREIGN KEY, or NOT NULL failure. | Ensure input values satisfy schema integrity constraints. |
| `ERR_PARSE_FAILURE` | 28 | Malformed CSV or assignment string. | Format values using `key=value` or valid comma delimiters. |

---

## 6. JSON Schema Envelopes

When invoked with `--json`, commands suppress terminal decorative chrome and output a standardized envelope.

### 6.1 Base Envelope Structure

```json
{
  "success": true,
  "command": "db-read",
  "database": "gitmap.db",
  "alias": "master",
  "timestamp": "2026-10-19T16:15:00Z",
  "data": {},
  "error": null
}
```

### 6.2 gitmap db-read Output Envelope

```json
{
  "success": true,
  "command": "db-read",
  "database": "gitmap.db",
  "alias": "master",
  "timestamp": "2026-10-19T16:15:00Z",
  "data": {
    "columns": ["id", "repoName", "isClean"],
    "rows": [
      [1, "gitmap", 1],
      [2, "repo-cache", 1]
    ],
    "count": 2,
    "executionTimeMs": 1.42
  },
  "error": null
}
```

### 6.3 gitmap db-read ls Output Envelope

```json
{
  "success": true,
  "command": "db-read:ls",
  "database": "installation.db",
  "alias": "install",
  "timestamp": "2026-10-19T16:15:00Z",
  "data": {
    "tables": [
      {"name": "Node", "type": "table", "columnCount": 7, "rowCount": 4},
      {"name": "Bookmark", "type": "table", "columnCount": 5, "rowCount": 12}
    ],
    "totalTables": 2
  },
  "error": null
}
```

### 6.4 gitmap db-add Output Envelope

```json
{
  "success": true,
  "command": "db-add",
  "database": "gitmap.db",
  "alias": "master",
  "timestamp": "2026-10-19T16:15:00Z",
  "data": {
    "table": "Repo",
    "insertedId": 48,
    "rowsAffected": 1
  },
  "error": null
}
```

### 6.5 gitmap db-edit Output Envelope

```json
{
  "success": true,
  "command": "db-edit",
  "database": "gitmap.db",
  "alias": "master",
  "timestamp": "2026-10-19T16:15:00Z",
  "data": {
    "table": "Repo",
    "id": 48,
    "rowsAffected": 1,
    "updatedFields": ["isClean", "remoteUrl"]
  },
  "error": null
}
```

### 6.6 Failure Envelope

```json
{
  "success": false,
  "command": "db-read",
  "database": "gitmap.db",
  "alias": "master",
  "timestamp": "2026-10-19T16:15:00Z",
  "data": null,
  "error": {
    "code": "ERR_READONLY_VIOLATION",
    "exitCode": 24,
    "message": "mutation attempted during db-read: attempt to write a readonly database",
    "remediation": "Use 'gitmap db-add' or 'gitmap db-edit' for write operations."
  }
}
```

---

## 7. Acceptance Criteria & Compliance Gates

```yaml
verificationGates:
  isDbReadSupported: true
  isDbAddSupported: true
  isDbEditSupported: true
  isAliasesParityGuaranteed: true
  isErrorCodesDocumented: true
  isJsonSchemaEnvelopesSpecified: true
```

- [x] Full argument parsing tables defined for `db-read`, `db-add`, and `db-edit`.
- [x] Short aliases `dbr`, `dba`, and `dbe` provide 100% flag and execution parity.
- [x] Exit codes 20 through 28 mapped to specific failure states.
- [x] Standard JSON envelope schemas defined for all success and error responses.
- [x] Total line count of this specification does not exceed 300 lines.
