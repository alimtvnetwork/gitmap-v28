# Plan: 284-agent-heatmap-db-introspection-and-subtask-code-completion

## Metadata
- **Task ID**: `284-agent-heatmap-db-introspection-and-subtask-code-completion`
- **Specification**: `02-spec/21-app/284-agent-heatmap-db-introspection-and-subtask-code-completion/spec.md`
- **Total Steps**: 2 Subtasks across 1 Wave (A=2, H=2)
- **Status**: IN_PROGRESS

---

## User Request (Verbatim)

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

## Subtask Breakdown

### Subtask-01: Subtask Complete `--code` & Resolution, Clean Syntax Error Handling & Native SQLite Introspection CLI
- **Owner**: Worker 01
- **Files Owned**:
  - `cli/cliexit/handle.go`
  - `cli/cmdagent/agent_subtask.go`
  - `cli/cmdagent/agent_subtask_resolve.go` (new file)
  - `cli/cmddb/` (new package: `cmddb_dispatch.go`, `cmddb_tables.go`, `cmddb_schema.go`, `cmddb_columns.go`, `cmddb_query.go`, `cmddb_records.go`, `cmddb_path_resolver.go`, `cmddb_introspect_types.go`)
  - `cli/cmd/rootutility.go`
- **Actions**:
  1. In `cliexit/handle.go`, classify usage/syntax errors and suppress raw callstack dumps.
  2. In `agent_subtask.go` and `agent_subtask_resolve.go`, bind `-c, --code` and `-s, --slug` flags on `complete`, `start`, and `fail` commands, resolving subtasks by ID, Code (e.g. `Subtask-01`), or Slug.
  3. Create `cli/cmddb/` implementing `gitmap db tables`, `gitmap db schema`, `gitmap db columns`, `gitmap db query`, and `gitmap db list` with clean `termout` tables and `--json` support.
  4. Register `gitmap db` command in `cli/cmd/rootutility.go`.
  5. Enforce line counts <= 300 across all modified/created Go files.

### Subtask-02: AI Decision Heatmap Split-DB Architecture, Cognitive Step Tracking & Skill Updates
- **Owner**: Worker 02
- **Files Owned**:
  - `cli/store/heatmap_schema.go` (new file)
  - `cli/store/heatmap_path.go` (new file)
  - `cli/store/heatmap_store.go` (new file)
  - `cli/heatmap/heatmap_cmd.go` (new file)
  - `cli/cmd/rootutility_update.go`
  - `cli/cmd/llm/llm_skill.go`
  - `.agents/skills/gitmap/SKILL.md`
- **Actions**:
  1. In `cli/store/heatmap_*.go`, create dual-layer SQLite heatmap engine (`heatmap.db` root store and `<baseDir>/heatmap/<repoSlug>/sql.db` per-repo split store).
  2. Implement `HeatmapParentTask`, `HeatmapReasoningState`, and `HeatmapActionLog` schemas with recursive parent-child execution lineages.
  3. In `cli/heatmap/heatmap_cmd.go`, implement `gitmap heatmap log`, `gitmap heatmap search`, `gitmap heatmap tree`, and `gitmap heatmap status`.
  4. Wire `gitmap heatmap` command into `cli/cmd/rootutility_update.go`.
  5. Update `cli/cmd/llm/llm_skill.go` and `.agents/skills/gitmap/SKILL.md` with DB introspection and heatmap decision tracking protocols.
  6. Enforce line counts <= 300 across all modified/created Go files.
