# Subtask 03: CLI Commands (db-rest, db-record-del, db-rm) & Port Fallback Engine

**Subtask ID:** `03-subtask-cli-commands-and-port-fallback`  
**Parent Task:** `database-rest-command-implementation`  
**Master Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
**Specification Reference:** `02-spec/21-app/database-rest-command-implementation/02-cli-and-rest-api-spec.md`  
**Assigned Role:** Worker 01  
**Status:** PENDING  
**Concurrency Budget:** A = 3, H = 2, Capacity = 6  

---

## 1. Subtask Overview & Scope

Worker 01 is responsible for authoring the complete CLI command parsing, terminal UI feedback, sequential port probing engine, and alias/database lifecycle operations for the `db-rest` subsystem.

### Owned Source Files (Strictly Disjoint):
- `cli/cmddb/cmd_rest.go`: Entry point for `gitmap db-rest` handling alias registration mode, target/port parsing, server startup, and terminal banner display.
- `cli/cmddb/cmd_record_del.go`: CLI command implementation for `gitmap db-record-del` and short alias `gitmap dbrd` to delete individual table records by primary key.
- `cli/cmddb/cmd_rm.go`: CLI command implementation for `gitmap db-rm` and alias `gitmap db-delete` to unregister database aliases and optionally purge physical SQLite files.
- `cli/cmddb/rest_port_picker.go`: Sequential port probing utility that verifies TCP port availability across comma-separated port candidates.
- `cli/cmddb/rest_target_resolver.go`: Comma-separated target parser that maps system aliases, custom registered aliases, and raw SQLite file paths to resolved database mount descriptors.

### Shared Integration Touchpoints (Coordinated Registration):
- `cli/cmd/rootdata275.go` & `cli/cmd/rootdata275_alias.go`: Dispatch routing for `db-rest`, `db-record-del`, `dbrd`, `db-rm`, `db-delete`.
- `cli/cmddb/cmddb_dispatch.go`: Subsystem routing within `cmddb`.

---

## 2. Invariants & Coding Standards

1. **Strict Relative Git Paths**: All documentation, comments, and log messages MUST use relative paths (e.g. `02-spec/21-app/...`, `cli/cmddb/...`, `.ai-memory/...`). Absolutely NO absolute filesystem paths or `file:///` URIs.
2. **File Size Cap**: Every Go source file MUST be strictly $\le 300$ lines. Decompose logic across cohesive helper files where necessary.
3. **Structured Errors**: Return `*appfault.AppError` using domain exit codes and structured categories (`appfault.NewInvalidInputError`, `appfault.NewNotFoundError`, `appfault.NewForbiddenError`).
4. **Single-Writer Safety**: All SQLite connections must use WAL mode and busy timeout (`_busy_timeout=5000&_journal_mode=WAL`) with `conn.SetMaxOpenConns(1)`.
5. **Standard Library Networking**: Use `net.Listen` for port probing and standard library `net/http` for HTTP serving. No external web frameworks.

---

## 3. Step-by-Step Implementation Instructions

### Step 1: Sequential Port Fallback Engine (`cli/cmddb/rest_port_picker.go`)
- Create `cli/cmddb/rest_port_picker.go` with functions:
  ```go
  // ParseCandidatePorts parses comma-separated port strings or single integer strings.
  // Defaults to []int{8080, 8081, 8082, 3000, 5000} if rawInput is empty.
  func ParseCandidatePorts(rawInput string) ([]int, error)

  // PickAvailablePort probes candidate ports in order using net.Listen("tcp", ...).
  // Immediately closes listener upon success and returns the first free port.
  func PickAvailablePort(candidates []int) (int, error)

  // CheckPortAvailable verifies whether a specific port is currently open for binding.
  func CheckPortAvailable(port int) bool
  ```
- If a candidate port is occupied, log to standard output:
  `[PORT] Port %d is currently occupied. Testing next candidate: %d...`
- If all candidates fail, return an error detailing all probed ports.

### Step 2: Target & Multi-DB Resolver (`cli/cmddb/rest_target_resolver.go`)
- Define the mount descriptor struct:
  ```go
  type DbMountTarget struct {
      Alias        string
      DatabasePath string
      DatabaseType string
      FileSizeBytes int64
      TableCount   int
      IsWritable   bool
  }
  ```
- Implement `ResolveMountTargets(rawTargets string) ([]DbMountTarget, error)`:
  1. Split `rawTargets` on commas `,`.
  2. For each trimmed token:
     - Check if token is a system alias (`master`, `pipeline`, `tasks`, `heatmap`) via `cmddb.resolveSystemDbAlias`.
     - Check if token exists in `DatabaseAlias` table via `cmddb.lookupAliasInRegistry`.
     - Check if token is an existing filesystem path. If so, derive alias from filename base without `.db` extension (e.g. `./mydata.db` -> `mydata`).
     - Check for alias collisions across the list; if duplicate alias exists, append numerical suffix (`-2`).
     - Stat file for size. Query `sqlite_master` for table count.
  3. Return slice of verified `DbMountTarget`.

### Step 3: CLI Entry Point & Mode Switcher (`cli/cmddb/cmd_rest.go`)
- Implement `RunDbRest(args []string) *appfault.AppError`:
  1. Print help text if args are empty or match `-h`, `--help`, `help`.
  2. Detect Alias Binding Mode:
     - Check if arguments contain the token `as` (case-insensitive), e.g. `gitmap db-rest "<path>" as <alias>`.
     - Extract path and alias name.
     - Call `cmddb.RegisterAlias(alias, path, "Registered via gitmap db-rest")`.
     - Print success confirmation and return `nil`.
  3. Detect REST Server Startup Mode:
     - Argument 0: target databases string (e.g. `analytics`, `master,pipeline`, `master,./data/store.db`).
     - Argument 1 (optional): port candidates string (e.g. `8080,8081,3000`).
     - Resolve candidate ports via `ParseCandidatePorts`.
     - Pick free port via `PickAvailablePort`.
     - Resolve target databases via `ResolveMountTargets`.
     - Instantiate and start `MultiDbRestServer` (from Subtask 04).
     - Print the terminal startup ASCII banner with all mounted database paths, catalog URL, Swagger URL, and UI URL.
     - Block on `server.ListenAndServe()` and handle graceful SIGINT / SIGTERM shutdown.

### Step 4: Record Deletion CLI (`cli/cmddb/cmd_record_del.go`)
- Implement `RunDbRecordDel(args []string) *appfault.AppError`:
  1. Support both positional syntax (`<db> <table> <id>`) and flagged syntax (`<db> -table <table> -id <id> [-y]`).
  2. Verify database target, resolve path, open single connection via `cmddb.OpenResolvedDb`.
  3. Retrieve table schema via `fetchColumnInfos(conn, table)`.
  4. Detect primary key column using `findPrimaryKeyCol(cols)`. If no explicit PK, fallback to `rowid`.
  5. Check if record exists: `SELECT count(*) FROM <table> WHERE <pkCol> = ?`. If count == 0, return `appfault.NewNotFoundError`.
  6. Check interactive prompt:
     - If not `-y` / `--yes` / `--force` and stdin is a TTY, prompt `[CONFIRM] Delete record #%v from '%s'? [y/N]: `.
     - If user declines, print `[ABORTED] Record deletion cancelled.` and return `nil`.
  7. Execute `DELETE FROM <table> WHERE <pkCol> = ?`.
  8. Print confirmation: `[SUCCESS] Deleted 1 record (%s=%v) from table '%s' in database '%s'.`.

### Step 5: Database / Alias Removal CLI (`cli/cmddb/cmd_rm.go`)
- Implement `RunDbRm(args []string) *appfault.AppError`:
  1. Syntax: `gitmap db-rm <alias-name> [--purge-file] [-y]`.
  2. Guard against system aliases:
     - If target alias is `master`, `gitmap`, `pipeline`, `tasks`, `heatmap`, return `appfault.NewForbiddenError("system database alias cannot be removed")`.
  3. Look up alias in `DatabaseAlias` table to find underlying file path.
  4. If `--purge-file` flag is provided:
     - If `-y` not passed, prompt user for confirmation: `[WARNING] Permanently delete database file '%s'? [y/N]: `.
     - Remove physical files: `<path>`, `<path>-wal`, `<path>-shm`.
  5. Delete alias entry from `DatabaseAlias` table:
     `DELETE FROM DatabaseAlias WHERE LOWER(AliasName) = LOWER(?)`.
  6. Print summary: `[SUCCESS] Removed alias '%s' from registry.`.

### Step 6: Dispatch Wiring & Route Registration
- Wire `db-rest`, `db-record-del`, `dbrd`, `db-rm`, `db-delete` into `cli/cmd/rootdata275.go` and `cli/cmd/rootdata275_alias.go`.
- Ensure hyphenated and spaced variants work uniformly (`gitmap db-rest` and `gitmap db rest`).

---

## 4. Function Signatures & Data Contracts

```go
package cmddb

import (
    "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// RunDbRest executes the REST API server daemon or alias binding command.
func RunDbRest(args []string) *appfault.AppError

// RunDbRecordDel executes primary-key record deletion from CLI.
func RunDbRecordDel(args []string) *appfault.AppError

// RunDbRm executes database alias unbinding and optional file purge.
func RunDbRm(args []string) *appfault.AppError

// ParseCandidatePorts parses comma-separated port numbers.
func ParseCandidatePorts(rawInput string) ([]int, error)

// PickAvailablePort probes ports in order until one is available.
func PickAvailablePort(candidates []int) (int, error)

// ResolveMountTargets resolves a comma-separated list of database aliases or file paths.
func ResolveMountTargets(rawTargets string) ([]DbMountTarget, error)
```

---

## 5. Verification Checklist & Acceptance Criteria

- [ ] **Port Fallback**: `PickAvailablePort` successfully detects an in-use port and binds to the subsequent candidate.
- [ ] **Alias Binding**: `gitmap db-rest "path/to/test.db" as testalias` successfully inserts into `DatabaseAlias` table.
- [ ] **Multi-DB Parsing**: `ResolveMountTargets("master,tasks,./data/custom.db")` returns 3 valid `DbMountTarget` entries.
- [ ] **Record Deletion**: `gitmap db-record-del tasks Subtask 101 -y` deletes the row and confirms deletion.
- [ ] **Safe Alias Removal**: `gitmap db-rm testalias` removes registration; attempting `gitmap db-rm master` returns `ForbiddenError`.
- [ ] **Line Count Cap**: All created `.go` files have $\le 300$ lines.
