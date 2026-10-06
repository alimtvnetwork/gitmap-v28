# SplitDB Error Storage & CLI Commands Specification

## 1. Overview
The current error storage implementation in cli/store/errors_split_db.go and cli/store/pull_split_db_errors.go centrally stores all error telemetry, including deep stack traces, into a single global SQLite database (gitmap-errors.db and pull databases). 

To ensure performance, isolation, and bounded database size, GitMap requires a two-tier SplitDB architecture for error storage:
1. **Root Error DB:** Acts as a central index, tracking *where* and *when* an error occurred.
2. **Per-Repo SQLite DB:** Stores the heavy diagnostic data, including deep stack traces, context JSON, and remediation commands.

## 2. Schema & Data Flow

### Root Error DB
- **Purpose:** Fast querying and indexing of error occurrences across the entire system.
- **Fields:** ErrorID, RepoSlug, RepoPath, ErrorType, ErrorCode, Timestamp, NodeID.
- **Behavior:** Dropping a repo automatically leaves an orphaned pointer, or can be cleaned up via sync.

### Per-Repo Error DB
- **Purpose:** Deep storage for telemetry.
- **Fields:** ErrorID, StackTrace, ContextJson, RemediationCmd, Message, Details.
- **Data Flow:**
  1. A failure occurs (e.g., during gitmap pull or general internal operations).
  2. The system generates a UUID ErrorID.
  3. A lean index record is inserted into the Root Error DB.
  4. The detailed trace and context are stored in the Per-Repo DB (e.g., .gitmap/repo-errors.db).

## 3. CLI Commands & Flag Handling

The new SplitDB error schema will be accessible and manageable via a fleet-aware CLI interface.

### gitmap see errors
- **Behavior:** Queries the Root Error DB to list recent errors. Can perform a join/fetch with Per-Repo DBs if the user requests detailed view.
- **Flags:** --limit <n>, --json.

### gitmap nodes errors [--json]
- **Behavior:** Uses cmdssh.RunFleetPASCommand to execute the error query across all enrolled SSH/Cluster nodes.
- **Aggregation:** Collects the Root Error DB indices from all nodes and presents a unified matrix of fleet errors.
- **Flags:** --json for machine-readable automation pipelines.

### gitmap nodes pull errors
- **Behavior:** Specifically filters for errors generated during fleet-wide pull operations.
- **Integration:** Delegates to cli/cmdpullerror/pull_error_cmd.go logic but routes it through the unified 
odes dispatcher.

### gitmap nodes errors clear
- **Behavior:** Executes a distributed clear command. 
- **Action:** Wipes the Root Error DB and cascades the clear operation to all Per-Repo Error DBs across all fleet nodes to recover disk space.

