# Subtask 01: Subtask Complete `--code`, Clean Error Handling & Native SQLite Introspection CLI

## Status: DONE
- **Parent Plan**: `.ai-memory/plans/284-agent-heatmap-db-introspection-and-subtask-code-completion.md`
- **Owner**: Worker 01
- **Owned Files**:
  - `cli/cliexit/handle.go`
  - `cli/cmdagent/agent_subtask.go`
  - `cli/cmdagent/agent_subtask_resolve.go`
  - `cli/cmddb/cmddb_dispatch.go`
  - `cli/cmddb/cmddb_tables.go`
  - `cli/cmddb/cmddb_schema.go`
  - `cli/cmddb/cmddb_columns.go`
  - `cli/cmddb/cmddb_query.go`
  - `cli/cmddb/cmddb_records.go`
  - `cli/cmddb/cmddb_path_resolver.go`
  - `cli/cmddb/cmddb_introspect_types.go`
  - `cli/cmd/rootutility.go`

## Objectives
1. In `cli/cliexit/handle.go`:
   - Detect CLI syntax/usage errors (`isUsageOrSyntaxError`) such as "unknown flag", "unknown shorthand flag", "flag needs an argument", "unknown command".
   - Wrap them as `apperror.NewValidation("cli", "E1000", err.Error())` with exit code 2 (`ExitCodeUsageError`).
   - Suppress callstack dumps for validation/usage errors.
2. In `cli/cmdagent/agent_subtask_resolve.go`:
   - Implement `ResolveSubtaskIdentifier(dbPath, identifier, code, slug string) (string, error)`:
     Supports lookup by ID, `TaskCode` (case-insensitive e.g. `Subtask-01`), or `TaskSlug` (case-insensitive).
3. In `cli/cmdagent/agent_subtask.go`:
   - Register `-c, --code` and `-s, --slug` flags on `subtaskCompleteCmd`, `subtaskStartCmd`, and `subtaskFailCmd`.
   - Resolve subtask dynamically before executing updates.
4. In `cli/cmddb/`:
   - Implement `gitmap db tables [db-path] [--json]`
   - Implement `gitmap db schema [table] [--db <path>] [--json]`
   - Implement `gitmap db columns <table> [--db <path>] [--json]`
   - Implement `gitmap db query "<SQL>" [--db <path>] [--json]`
   - Implement `gitmap db list <table> [--limit N] [--db <path>] [--json]`
   - Format terminal output with `cli/termout` tables.
5. In `cli/cmd/rootutility.go`:
   - Register `CmdDb` and route to `cmddb.DispatchDbSubsystem`.
6. FastGate: Ensure all modified/created Go files stay strictly <= 300 lines with nesting depth <= 1.
