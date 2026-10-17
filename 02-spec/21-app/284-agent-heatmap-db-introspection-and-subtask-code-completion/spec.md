# Specification: Agent Heatmap, DB Introspection & Subtask Code Completion

## Status
- **Task ID**: `284-agent-heatmap-db-introspection-and-subtask-code-completion`
- **Specification Version**: 1.0.0
- **Status**: APPROVED
- **Date**: 2026-10-18

---

## 1. User Request (Verbatim)

```text
gitmap agent subtask complete --code "Subtask-01"

Error: unknown flag: --code
Usage:
  agent subtask complete <subtask-id> [flags]

Flags:
  -a, --agent string      Agent role completing subtask
  -e, --evidence string   Completion evidence or summary
  -h, --help              help for complete
  -t, --task-id string    Parent task ID or slug

gitmap: [E9000:EXECUTION] cli
  origin: cliexit/handle.go:75
  cause: unknown flag: --code
  stack trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cliexit.ensureAppError (cliexit/handle.go:75)
        at github.com/alimtvnetwork/gitmap-v28/cli/cliexit.HandleError (cliexit/handle.go:41)
...

# High Priority Instruction

Let's try to fix this stack trace error. And also at the same time, you can see that sometimes the AI tries to run other Python-based commands to find the table and DB information. Why don't we have the `gitmap` commands to get these things? And also, it should be in the LLM train, and also the `gitmap` skill, and also update the `gitmap` skills in the coding guideline. So these are the things you need to do. And also try to think of other commands that the AI has run, let's say, by executing in recent few steps. Can you please check what are the commands that AI needs in most cases, and give an example with a Python example and also this one, how the AI can learn. And also, in the future, what I want is the heat map thing. What I mean by the heat map thing is that, for example, when the AI is reading the file, it writes the file line. This is analyzing, this is modifying. So every step that AI takes, the `gitmap` should have a command from the agent AI. So that would actually keep this file path, file lines, and also the text for those lines inside a SQLite DB for the heat map that we will call a speed DB for heat map. And that would reside into the temp folder for the tasks, actually. I don't think that it should be inside the temp folder inside the AI task. No. It should be besides the original DB. The DB should be heat map DB as its own root DB. But also, for every repo, it would create its own folder, and we will have the heat map DBs for that folder, and it will only contain the decision-making information...
And also at the end, do a minor bump and release. Is it clear?
```

---

## 2. Problem Statement & Root Cause Analysis

### 2.1 Subtask `--code` Failure & Scary Stack Traces
- **Problem:** AI passed `gitmap agent subtask complete --code "Subtask-01"` and received `Error: unknown flag: --code` followed by a full 15-line Go execution callstack.
- **Root Cause 1 (Missing Flag):** `subtaskCompleteCmd` in `cli/cmdagent/agent_subtask.go` only accepted positional `<subtask-id>`. It did not register `-c, --code` or `-s, --slug`.
- **Root Cause 2 (Missing Polymorphic Lookup):** Even when a positional argument was supplied, `RunSubtaskComplete` performed `WHERE SubtaskId = ?` strictly against the generated timestamp ID (`sub-153602-subtask-01`), failing to resolve user/plan codes like `Subtask-01` or slugs like `update-ui-colors`.
- **Root Cause 3 (Blind Fatal Wrapping in `cliexit`):** In `cli/cliexit/handle.go`, `ensureAppError(err)` blind-wrapped unknown Cobra flag errors using `apperror.WrapSimple(err, "cli")`. This assigned error code `E9000` (`ErrorTypeExecution`), triggering `isFatal = true` and dumping a full debug callstack for a simple command-line syntax error.

### 2.2 Lack of Native SQLite DB & Schema Introspection
- **Problem:** When AI agents need to inspect table structures, columns, or records in `gitmap.db`, `installation.db`, or `agent-task.db`, they resort to executing raw Python scripts:
  `python -c "import sqlite3; conn = sqlite3.connect('...'); print([col[1] for col in conn.execute('PRAGMA table_info(...)').fetchall()])"`
- **Root Cause:** GitMap lacked native commands for inspecting SQLite tables, schemas, columns, and data rows.

### 2.3 AI Decision Heatmap & Cognitive Step Tracking
- **Problem:** In autonomous coding sessions, AI agents analyze dozens of files and lines across a repository without persisting *why* a file was read, *what* was learned, *what* action was taken, and *whether* expected outcomes were achieved.
- **Root Cause:** GitMap had no dedicated telemetry Split-DB for logging microsecond cognitive decisions, action logs (lines analyzed, code modified, commands executed), reasoning states, and parent-child causal lineages.

---

## 3. Architecture & Requirements

### 3.1 Subtask Command Fixes & Error Handling
1. **Flag & Resolver Support:**
   - In `cli/cmdagent/agent_subtask_resolve.go` (new file, <= 300 lines):
     Implement `ResolveSubtaskIdentifier(dbPath, identifier, code, slug string) (string, error)`:
     Matches `SubtaskId`, `TaskCode` (case-insensitive), or `TaskSlug` (case-insensitive).
   - In `cli/cmdagent/agent_subtask.go`:
     Register `-c, --code` and `-s, --slug` flags on `subtaskCompleteCmd`, `subtaskStartCmd`, and `subtaskFailCmd`.
     Use `ResolveSubtaskIdentifier` before executing status updates.
2. **Clean Syntax Error Handling:**
   - In `cli/cliexit/handle.go`:
     Implement `isUsageOrSyntaxError(err error) bool` detecting "unknown flag:", "unknown command", "flag needs an argument", etc.
     If true, return `apperror.NewValidation("cli", "E1000", err.Error())` with exit code `2` (`ExitCodeUsageError`).
     Suppress debug stack traces on validation/usage errors.

### 3.2 Native SQLite DB Introspection CLI (`gitmap db ...`)
Create modular package `cli/cmddb/` (all files <= 300 lines):
- `gitmap db tables [db-path] [--json]`: Lists tables, views, column counts, and record counts.
- `gitmap db schema [table] [--db <path>] [--json]`: Dumps formatted `CREATE TABLE` and index DDL.
- `gitmap db columns <table> [--db <path>] [--json]`: Displays column names, types, nullability, defaults, and primary key badges.
- `gitmap db query "<SQL>" [--db <path>] [--json]`: Executes read-only queries with tabular and JSON output.
- `gitmap db list <table> [--limit N] [--db <path>] [--json]`: Displays formatted table rows.
- Semantic `--db` shortcuts: `master`, `gitmap`, `heatmap`, `pipeline`, `tasks`, `agents`, or file path.

### 3.3 AI Decision Heatmap Split-DB System
Implement package `cli/heatmap/` and `cli/store/`:
- **Dual-Layer Topology:**
  - Root Global DB: `store.EffectiveDataDir()/heatmap.db`
  - Per-Repo Split DBs: `<baseDir>/heatmap/<repoSlug>/sql.db`
- **Schema:**
  - `HeatmapParentTask`: Registers parent tasks, slugs, and status.
  - `HeatmapReasoningState`: Standardized state machine (`EXPLORATION`, `SPEC_REFINEMENT`, `CODE_EDITING`, `VERIFICATION`, `REMEDIATION`, `COMPLETED`).
  - `HeatmapActionLog`: Records `LogId`, `ParentLogId`, `TaskSlug`, `RepoName`, `ActionType` (`READ`, `ANALYZE`, `MODIFY`, `EXEC_CMD`), `FilePath`, `StartLine`, `EndLine`, `CodeSnippet`, `CommandText`, `ReasoningState`, `ReasoningNotes`, `ExpectedOutcome`, `ActualOutcome`, `IsExpectationMet`, `ExecutionDepth`, `ExecutionOrder`.
- **CLI Commands:**
  - `gitmap heatmap log [--task <slug>] [--action <type>] [--file <path>] [--lines <start-end>] [--reason <notes>] [--expected <outcome>] [--parent <id>]`
  - `gitmap heatmap search "<query>" [--task <slug>] [--file <path>] [--json]`
  - `gitmap heatmap tree [--task <slug>] [--json]`
  - `gitmap heatmap status [--task <slug>] [--json]`

### 3.4 Skill & Guidelines Updates
- Update `.agents/skills/gitmap/SKILL.md` to document `gitmap db` and `gitmap heatmap` commands.
- Update `cli/cmd/llm/llm_skill.go` so AI self-training learns native `gitmap db` commands.

### 3.5 Minor Bump & Release
- Execute minor version bump (`v6.567.0`), update `changelog.md` and `version.json`, tag, and publish.

---

## 4. File Boundaries & Ownership

- **Worker 01**:
  - `cli/cliexit/handle.go` (syntax error classification, stack trace suppression)
  - `cli/cmdagent/agent_subtask.go` (register `--code` and `--slug` flags)
  - `cli/cmdagent/agent_subtask_resolve.go` (new file: polymorphic subtask resolver)
  - `cli/cmddb/` (new package: `cmddb_dispatch.go`, `cmddb_tables.go`, `cmddb_schema.go`, `cmddb_columns.go`, `cmddb_query.go`, `cmddb_records.go`, `cmddb_path_resolver.go`, `cmddb_introspect_types.go`)
  - `cli/cmd/rootutility.go` (register `CmdDb`)
- **Worker 02**:
  - `cli/store/heatmap_schema.go` (new file: DDL, indexes, seed states)
  - `cli/store/heatmap_path.go` (new file: dual-layer path resolver)
  - `cli/store/heatmap_store.go` (new file: CRUD operations, action logging, CTE tree)
  - `cli/heatmap/heatmap_cmd.go` (new file: CLI commands for heatmap log, search, tree, status)
  - `cli/cmd/rootutility_update.go` or `cli/cmd/rootutility.go` (register `CmdHeatmap`)
  - `cli/cmd/llm/llm_skill.go` (update LLM train skills)
  - `.agents/skills/gitmap/SKILL.md` (document db and heatmap commands)
