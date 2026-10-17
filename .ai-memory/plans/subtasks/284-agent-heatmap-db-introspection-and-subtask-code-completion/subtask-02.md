# Subtask 02: AI Decision Heatmap Split-DB Architecture, Cognitive Step Tracking & Skill Updates

## Status: DONE
- **Parent Plan**: `.ai-memory/plans/284-agent-heatmap-db-introspection-and-subtask-code-completion.md`
- **Owner**: Worker 02
- **Owned Files**:
  - `cli/store/heatmap_schema.go`
  - `cli/store/heatmap_path.go`
  - `cli/store/heatmap_store.go`
  - `cli/heatmap/heatmap_cmd.go`
  - `cli/cmd/rootutility_update.go`
  - `cli/cmd/llm/llm_skill.go`
  - `.agents/skills/gitmap/SKILL.md`

## Completed Objectives
1. In `cli/store/heatmap_schema.go` (218 lines, <= 300):
   - Defined DDL constants and initialization for `HeatmapParentTask`, `HeatmapReasoningState`, and `HeatmapActionLog`.
   - Seeded reasoning states: `EXPLORATION`, `SPEC_REFINEMENT`, `CODE_EDITING`, `VERIFICATION`, `REMEDIATION`, `COMPLETED`.
   - Implemented `EnsureHeatmapTables(conn *sql.DB) error` and row scanner/tree assembler helpers.
2. In `cli/store/heatmap_path.go` (30 lines, <= 300):
   - Implemented `ResolveHeatmapRootDbPath(repoRoot string)` returning global or repo-anchored `heatmap.db`.
   - Implemented `ResolveHeatmapRepoSplitDbPath(repoSlug, repoRoot string)` via `ResolveSplitDbPath("heatmap", repoSlug, repoRoot)`.
3. In `cli/store/heatmap_store.go` (289 lines, <= 300):
   - Implemented `HeatmapDB` struct wrapping `*sql.DB`.
   - Implemented `OpenHeatmapRootDB(repoRoot string) (*HeatmapDB, error)` and `OpenHeatmapRepoDB(repoSlug, repoRoot string) (*HeatmapDB, error)`.
   - Implemented `LogAction` calculating `ExecutionDepth` and `ExecutionOrder` from `ParentLogId`.
   - Implemented `SearchActionLogs(taskSlug, query, filePath string, limit int) ([]HeatmapActionLogRecord, error)`.
   - Implemented `GetExecutionTree(taskSlug string) ([]HeatmapExecutionTreeNode, error)` using recursive SQLite CTE.
   - Implemented `GetHeatmapStatus(taskSlug string) (HeatmapStatusSummary, error)`.
4. In `cli/heatmap/heatmap_cmd.go` (280 lines, <= 300):
   - Implemented `DispatchHeatmapSubsystem(ctx context.Context, args []string, out io.Writer) error`.
   - Supported subcommands: `log`, `search`, `tree`, `status`, and `help` with ANSI table and `--json` formatting.
5. In `cli/cmd/rootutility_update.go` (195 lines, <= 300):
   - Registered `CmdHeatmap = "heatmap"` with alias `CmdHeatmapAlias = "hm"`.
   - Added `init()` hook and `RunHeatmapCommand` routing to `heatmap.DispatchHeatmapSubsystem`.
6. In `cli/cmd/llm/llm_skill.go` (272 lines, <= 300) & `.agents/skills/gitmap/SKILL.md`:
   - Documented `gitmap db` (tables, schema, columns, query, list) and `gitmap heatmap` (log, search, tree, status).
   - Added Python migration table showing replacements for `python -c "import sqlite3..."`.
   - Enforced operational guardrail 10 banning inline Python SQLite queries in favor of compiled `gitmap db` and `gitmap heatmap`.
7. FastGate Verification: All 6 Go files verified strictly <= 300 lines with nesting depth <= 1.
