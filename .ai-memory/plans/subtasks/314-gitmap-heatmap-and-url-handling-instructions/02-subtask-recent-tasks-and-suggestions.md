# Subtask 02: Multi-Repository Recent Tasks Aggregator, Crisp ASCII Table/Tree Rendering, and Dynamic AI Terminal Suggestions

## Status: READY
- **Task ID:** `314-subtask-02`
- **Parent Spec:** [02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md](../../../../02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md)
- **Module:** `cli/heatmap`
- **Owned Files:**
  - `cli/heatmap/heatmap_recent_tasks.go`
  - `cli/heatmap/heatmap_suggestions.go`
  - `cli/heatmap/heatmap_cmd_render.go`
  - `cli/heatmap/types.go`

---

## 1. Objectives

1. **Multi-Repository Recent Tasks Aggregation (`cli/heatmap/heatmap_recent_tasks.go`):**
   - Consume discovered database entries from `DiscoverAllHeatmapDatabases()`.
   - Select up to 5 most recently active databases across 2 to 5 distinct repositories.
   - For each database, query the `HeatmapActionLog` table to aggregate task execution metrics:
     - `TaskSlug`: Unique identifier of the task session.
     - `TotalActions`: Total volume of logged operations (`COUNT(*)`).
     - `ExpectationsMet`: Count of actions where `IsExpectationMet == true`.
     - `ExpectationsMetPct`: Ratio of verified outcomes expressed as a percentage.
     - `ReasoningStates`: Categorical breakdown of reasoning states (`PASS`, `FAIL`, `EXPLORATION`, `RETRY`, `EXECUTION`).
     - `LastReasoningState`: Most recent reasoning state recorded.
     - `LastActiveAt`: Most recent activity timestamp.
     - `hasFailedExpectations`: Positive boolean indicating if unfulfilled expectations or failures exist.
   - Sort aggregated tasks globally by recency and return the top 5 to 10 tasks.

2. **Dynamic AI & Developer Suggestions Engine (`cli/heatmap/heatmap_suggestions.go`):**
   - Evaluate multi-repository task metrics and system database health to generate context-sensitive terminal command recommendations.
   - Generate tailored suggestions for:
     - **Failure Remediation:** If `hasFailedExpectations` is true, recommend `gitmap heatmap tree --task <taskSlug>` to review the causal hierarchy, and `gitmap heatmap search --task <taskSlug> --state FAIL` to isolate root causes.
     - **Telemetry Synchronization:** If local root database contains uncommitted actions or dirty sync state, recommend `gitmap rc heatmap sync --repo <repoSlug> --task <taskSlug>`.
     - **Multi-Agent Coordination:** If multiple AI namespaces are active, recommend `gitmap rc heatmap status --repo <repoSlug>` to inspect remote tracking.
     - **Initial Setup:** If no tasks are found, provide copyable `gitmap heatmap log` templates.

3. **Crisp ASCII Table & Visual Progress Rendering (`cli/heatmap/heatmap_cmd_render.go`):**
   - Extend terminal renderer with `renderRecentTasksTable(out io.Writer, summary MultiRepoRecentSummary)`.
   - Display cleanly aligned tabular output with repository badges, percentage meters, reasoning states, and relative time representations (e.g. `12m ago`).
   - Highlight states using `termout` color tokens: green for `PASS`, red for `FAIL`, yellow for `RETRY` / `WARN`.
   - Render dynamic suggestions as highlighted, copyable CLI command blocks with section banners.

4. **Overview Integration & JSON Serialization:**
   - Seamlessly blend recent tasks and dynamic suggestions into `RunHeatmapOverview` output.
   - When `--json` is requested, include `recentTasks` and `suggestions` arrays inside the structured JSON payload.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Author `cli/heatmap/heatmap_recent_tasks.go`
- Define data structures in `cli/heatmap/types.go` or `cli/heatmap/heatmap_recent_tasks.go`:
  - `RecentTaskSummary`
  - `MultiRepoRecentSummary`
- Implement `QueryMultiRepoRecentTasks(entries []HeatmapDbEntry, maxTasks int) (MultiRepoRecentSummary, error)`:
  - Filter `entries` to active databases sorted by `LastModifiedAt` descending.
  - Pick top 2 to 5 distinct repository databases.
  - Iterate through databases and execute `queryTasksFromDatabase(entry, maxTasks)`.
  - Combine, sort by `LastActiveAt` descending, and trim to `maxTasks` (default 5 to 10).
- Implement `queryTasksFromDatabase(entry HeatmapDbEntry, limit int) ([]RecentTaskSummary, error)`:
  - Open database in read-only mode (`mode=ro`).
  - Execute grouped query on `HeatmapActionLog` by `TaskSlug`.
  - Query reasoning state breakdown for each task:
    `SELECT ReasoningState, COUNT(*) FROM HeatmapActionLog WHERE TaskSlug = ? GROUP BY ReasoningState`.
  - Calculate `ExpectationsMetPct` and set positive boolean `hasFailedExpectations`.
- Ensure all helper functions stay strictly $\le 15$ lines.

### Step 2: Implement `cli/heatmap/heatmap_suggestions.go`
- Implement `GenerateDynamicSuggestions(summary MultiRepoRecentSummary, entries []HeatmapDbEntry) []string`:
  - Initialize empty slice `suggestions := make([]string, 0, 8)`.
  - Check for failed expectations across recent tasks via `hasFailedExpectations`.
  - If failures detected, append failure inspection commands via `buildFailureSuggestions(task)`.
  - Check for local unsynchronized telemetry across `entries`. If dirty, append sync commands via `buildSyncSuggestions(task)`.
  - If no tasks exist, append initialization command.
  - Return bounded list of top 3 to 5 highest-priority suggestions.
- Implement specialized helper functions:
  - `buildFailureSuggestions(task RecentTaskSummary) []string`
  - `buildSyncSuggestions(task RecentTaskSummary) []string`
- Guarantee zero nested conditionals and functions $\le 15$ lines.

### Step 3: Enhance `cli/heatmap/heatmap_cmd_render.go`
- Implement `renderRecentTasksTable(out io.Writer, summary MultiRepoRecentSummary)`:
  - Header: `Recent Multi-Repository Tasks:`
  - Columns: `TASK SLUG`, `REPO`, `ACTIONS`, `MET %`, `LAST STATE`, `UPDATED`.
  - Print formatted rows with aligned whitespace.
- Implement `renderDynamicSuggestions(out io.Writer, suggestions []string)`:
  - Header: `Dynamic Suggestions for AI & Developers:`
  - Print each suggestion preceded by accent pointer `▶` and indented command block.
- Implement `formatRelativeTime(t time.Time) string`:
  - Compute `time.Since(t)` and format cleanly as `<N>s ago`, `<N>m ago`, `<N>h ago`, or `<N>d ago`.

### Step 4: Wire Unified Terminal Output in `cli/heatmap/heatmap_overview.go`
- In `RunHeatmapOverview`:
  - Run discovery: `entries, err := DiscoverAllHeatmapDatabases()`.
  - Enrich aliases: `enriched := enrichEntriesWithAliases(entries)`.
  - Run recent tasks query: `recent, err := QueryMultiRepoRecentTasks(enriched, 10)`.
  - Generate suggestions: `suggestions := GenerateDynamicSuggestions(recent, enriched)`.
  - If JSON: serialize composite struct to stdout.
  - If Text: sequentially call `renderDiscoveryTable`, `renderRecentTasksTable`, and `renderDynamicSuggestions`.

---

## 3. Strict Coding Guidelines & Constraints

- **Function Line Limit:** Every function must be $\le 15$ lines (preferred $\le 8$ lines).
- **Positive Booleans Only:** Use `hasFailedExpectations`, `hasActiveTasks`, `isExpectationMet`, `isRecent`. Never use negative booleans (`isNotFailed`, `unverified`).
- **Zero Nested Conditionals:** Guard clauses and early returns for all conditional checks.
- **Error Wrapping:** Explicitly wrap all errors with `apperror.WrapSimple(err, "context")`.
- **Relative Path Hygiene:** Zero hardcoded drive letters (`C:`, `D:`) or URI schemes (`file:///`).
- **Zero Builds / Zero Tests:** Do not execute `go build` or `go test` during this authoring phase.

---

## 4. Verification Evidence & Acceptance Criteria

1. Multi-repository task aggregation correctly surveys top 5 to 10 tasks across 2 to 5 repositories.
2. Percentage metrics accurately calculate `ExpectationsMet / TotalActions * 100.0`.
3. Tasks with failed actions trigger dynamic tree and search recommendations in the terminal suggestions section.
4. Output formatting aligns neatly across 80-column and wider terminals with `termout` color styling.
5. All functions pass strict 15-line limit and positive boolean naming conventions.
