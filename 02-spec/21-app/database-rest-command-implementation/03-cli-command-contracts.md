# 03-cli-command-contracts: Database REST CLI Commands & Syntax

- **Spec ID:** `02-spec/21-app/database-rest-command-implementation/03-cli-command-contracts.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** CLI Commands (`db-rest`, `db-record-del`, `db-rm`), Syntax Parsing & Port Probing
- **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`
- **Dependencies:** `cli/cmddb`, `cli/store`, `cli/appfault`

---

## 1. Executive Summary & Purpose

This specification defines the exhaustive technical contract for the GitMap SQLite Database REST (`db-rest`) CLI command suite and its underlying high-performance HTTP REST API engine. 

The system enables developers to expose any local SQLite database—or multiple SQLite databases simultaneously—over a standards-compliant RESTful interface with on-the-fly schema introspection, dynamic CRUD generation, column search across tables, view authoring with joins, table alteration/DDL execution, schema migration with Split-DB caching, database cloning, and raw SQL execution.

### Architectural Invariants
1. **Universal Response Envelope**: Every HTTP response MUST return JSON with root keys `{ Status, Attributes, Results }` using PascalCase keys throughout.
2. **Standard Library HTTP Router**: Built on Go standard library `net/http` router with zero third-party web framework dependencies.
3. **CGO-Free Pure Go SQLite**: Built on `modernc.org/sqlite` with `conn.SetMaxOpenConns(1)` and WAL journal mode (`_busy_timeout=5000&_journal_mode=WAL`) to guarantee single-writer safety and eliminate database locks.
4. **Buffer-Before-Next Pattern**: All query results are fully scanned into in-memory data structures before releasing row iterators or executing subsequent queries.
5. **Multi-DB Registry Routing**: Uniformly mounts system aliases (`master`, `pipeline`, `tasks`, `heatmap`), custom registered aliases, and direct SQLite file paths under discrete namespace slugs (`/api/v1/db/{alias}/...`).

---

## 2. CLI Command Syntax & Parsing Contracts

The CLI layer is implemented in package `cli/cmddb` and registered in `cli/cmd`. The commands provide terminal access to start REST servers, bind aliases, delete individual records, and remove databases/aliases.

### 2.1 Command Grammar Summary

| Command Form | Syntax Example | Core Behavior |
|---|---|---|
| **Alias Binding** | `gitmap db-rest "db path" as alias` | Binds an arbitrary SQLite file path to a persistent alias in `DatabaseAlias` table. |
| **Single-DB Server** | `gitmap db-rest aliasName $port,$port2` | Starts REST service mounting `aliasName`, probing ports in order until one binds. |
| **Multi-DB Server (Aliases)** | `gitmap db-rest alias1,alias2 $port,$port2` | Mounts multiple database aliases simultaneously under a unified catalog. |
| **Multi-DB Server (Mixed)** | `gitmap db-rest alias1,dbPath2 $port,$port2` | Mounts registered aliases alongside direct filesystem paths with auto-generated aliases. |
| **Record Deletion** | `gitmap db-record-del <db> <table> <id>` (alias: `dbrd`) | Deletes a record from a specified table by its primary key. |
| **Alias / DB Removal** | `gitmap db-rm <alias> [--purge-file]` (alias: `db-delete`) | Unbinds an alias from registry, optionally deleting the underlying `.db` file. |

---

### 2.2 Alias Binding: `gitmap db-rest "<db-path>" as <alias>`

#### 2.2.1 Syntax & Flag Rules
```bash
gitmap db-rest "<path-to-sqlite-db>" as <alias-name>
```

- **Positional Tokens**:
  1. `<path-to-sqlite-db>`: Quoted or unquoted filesystem path to a `.db`, `.sqlite`, or `.sqlite3` file. Relative paths are cleaned and resolved relative to current working directory.
  2. `as`: Literal case-insensitive keyword token (`as` or `AS`).
  3. `<alias-name>`: Alphanumeric slug representing the alias. May contain hyphens and underscores (`^[a-zA-Z0-9_-]+$`).
- **Parsing Invariant**:
  When CLI argument list contains token `as` between two non-flag arguments, `cmddb.RunDbRest` routes directly to alias registration mode rather than starting HTTP daemon.

#### 2.2.2 Execution Flow
1. Verify target file exists on disk. If file does not exist, return `appfault.NewNotFoundError("database file not found: <path>")`.
2. Call `cmddb.RegisterAlias(aliasName, cleanPath, "User registered via gitmap db-rest")`.
3. Read file metadata (size in bytes, table count via `SELECT count(*) FROM sqlite_master WHERE type='table'`).
4. Output confirmation banner:
   ```text
   [SUCCESS] Registered database alias 'analytics' -> ./data/warehouse.db (2.4 MB, 14 tables)
   You can now start the REST server with: gitmap db-rest analytics 8080
   ```

---

### 2.3 Starting REST Server: `gitmap db-rest <targets> [ports]`

#### 2.3.1 Syntax Variations
```bash
# 1. Single alias with comma-separated port fallback
gitmap db-rest analytics 8080,8081,3000

# 2. Direct path with default ports (8080,8081,8082,3000)
gitmap db-rest ./data/cache.db

# 3. Multiple registered aliases
gitmap db-rest master,pipeline,tasks 8080,8081

# 4. Mixed aliases and raw database paths
gitmap db-rest master,./custom/store.db 9000,9001
```

#### 2.3.2 Target Resolver Rules
The first positional argument `<targets>` is split on commas `,` into individual target identifiers:
- For each token $T$:
  1. If $T$ matches a known system alias (`master`, `pipeline`, `tasks`, `heatmap`), resolve to its standard path.
  2. If $T$ exists in `DatabaseAlias` table in `gitmap.db`, resolve to registered `DatabasePath`.
  3. If $T$ is an existing file path, synthesize an alias slug from base filename.
  4. If $T$ cannot be found on disk or in registry, return `appfault.NewNotFoundError("database target not found: " + T)`.

#### 2.3.3 Port Probing & Sequential Fallback Engine
The port argument can be a single port or a comma-separated list of candidate ports (e.g. `8080,8081,3000`). If omitted, defaults to `8080,8081,8082,3000,5000`.

- **Fallback Logging**:
  ```text
  [PORT] Port 8080 is already occupied. Falling back to next candidate: 8081...
  [PORT] Successfully bound to port 8081.
  ```

#### 2.3.4 Terminal Startup Banner
```text
================================================================================
                    GITMAP DATABASE REST SERVER (v28)
================================================================================
 Server Listening: http://localhost:8081
 Mounted Databases (3):
   * master    -> store/gitmap.db (1.2 MB)
   * pipeline  -> repodb/pipeline.db (450 KB)
   * analytics -> data/warehouse.db (2.4 MB)

 Endpoints:
   * Catalog API:        http://localhost:8081/api/v1/databases
   * Master Tables:      http://localhost:8081/api/v1/db/master/tables
   * Interactive UI:     http://localhost:8081/ui
   * Swagger / OpenAPI:  http://localhost:8081/swagger
================================================================================
 Press Ctrl+C to stop the server.
```

---

### 2.4 Record Deletion CLI: `gitmap db-record-del` / `dbrd`

#### 2.4.1 Syntax Rules
```bash
gitmap db-record-del <alias-or-path> <table> <id>
gitmap db-record-del <alias-or-path> -table <table> -id <id> [-y]
gitmap dbrd tasks Subtask 42 -y
```

#### 2.4.2 Parameter Specification & Execution Logic
1. Resolve database path and open single connection.
2. Query `PRAGMA table_info(<table>)` to discover primary key column name.
3. Fetch existing record. If not found, return `appfault.NewNotFoundError("record not found")`.
4. If stdin is a TTY and `-y` flag is absent, prompt confirmation.
5. Execute `DELETE FROM <table> WHERE <pkCol> = ?`.
6. Output success summary: `[SUCCESS] Deleted 1 record (SubtaskId=42) from table 'Subtask' in database 'tasks'.`

---

### 2.5 Database / Alias Removal CLI: `gitmap db-rm` / `db-delete`

#### 2.5.1 Syntax Rules
```bash
gitmap db-rm <alias-name> [-y]
gitmap db-rm <alias-name> --purge-file [-y]
gitmap db-delete custom-db --purge-file
```

#### 2.5.2 Safeguards & Protection Rules
- **System Alias Immutability**:
  Aliases `master`, `gitmap`, `pipeline`, `pipelines`, `tasks`, `agents`, `heatmap` cannot be removed.
- **Physical File Purge Guard**:
  When `--purge-file` is passed, requires explicit confirmation unless `-y` is supplied.
- **Execution Flow**:
  1. Remove row from `DatabaseAlias` in `gitmap.db`.
  2. If `--purge-file` requested, call `os.Remove(resolvedPath)` along with WAL/SHM files.
