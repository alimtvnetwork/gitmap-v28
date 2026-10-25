# Subtask 02: Heatmap Status Telemetry Overhaul & Visual Modernization

## Status: PLANNED
- **Subtask Slug:** `02-subtask-heatmap-status-improvement`
- **Parent Plan:** `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md`
- **Architecture Spec:** `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/01-architecture-spec.md`
- **Owner:** Telemetry & CLI Subagent (Worker 01)
- **Target Subsystem:** `cli/heatmap/`
- **Exclusive Owned Files:**
  - `cli/heatmap/heatmap_cmd_render.go`
  - `cli/heatmap/heatmap_tasks.go`
  - `cli/heatmap/heatmap_cmd.go`
  - `cli/heatmap/types.go`

---

## 1. Executive Summary & Context

The current `gitmap heatmap status` command produces basic, unpadded text that looks unpolished and lacks critical operational context:
1. It does not showcase recently completed architectural plans or currently pending tasks.
2. It lacks the professional visual styling, padded borders, and TokyoNight colors used throughout the rest of GitMap.
3. It does not provide copy-pasteable tree-view inspection command examples for users to drill into details.
4. It does not emit a machine-readable `--json` payload formatted for consumption by autonomous AI swarms executing the V6 prompt protocol.

Subtask 02 refactors the status subsystem in `cli/heatmap/` to deliver an executive-grade terminal telemetry dashboard and a compact JSON export.

---

## 2. File Modification & Structural Ownership

```text
cli/heatmap/
├── heatmap_cmd_render.go  [MODIFIED]  - Replace renderStatusText with Lipgloss TokyoNight card renderer
├── heatmap_tasks.go       [NEW]       - Task discovery from .ai-memory/plans/ & SQLite registries
├── heatmap_cmd.go         [MODIFIED]  - Wire task discovery, --json flag routing, and status dispatch
└── types.go               [MODIFIED]  - Define HeatmapStatusJSONPayload, TaskItem, and Alert types
```

---

## 3. Detailed Technical Architecture

### 3.1 Lipgloss TokyoNight Visual Palette & Padding Rules

The terminal renderer in `cli/heatmap/heatmap_cmd_render.go` will use Lipgloss components from `cli/termout/theme_lipgloss.go`:

- **Banner Header:** `termout.TokyoBannerStyle("GITMAP HEATMAP TELEMETRY STATUS")`
- **Border Box:** Rounded border container styled with `TokyoMuted` (`#565f89`) and 1-character horizontal padding.
- **Section Headers:** Bold `TokyoPrimary` (`#7aa2f7`) section labels:
  - `Telemetry Summary`
  - `Completed Tasks (Last 5)`
  - `Pending Tasks Queue (Top 5)`
  - `Drill-Down Inspection Commands`
- **Status Pills:**
  - `termout.TokyoOKPill("[✓] COMPLETED")` for finished tasks and high pass rates.
  - `termout.TokyoWarnPill("[⏳] ACTIVE")` for currently executing tasks.
  - `termout.TokyoErrorPill("[!] ALERT")` for active AI confusions or error states.
- **Card Padding & Spacing:** Exactly 1 blank line before return statements; exactly 1 blank line after closing braces `}`; strict 15-line limit per function.

### 3.2 Dual-Source Task Aggregator (`cli/heatmap/heatmap_tasks.go`)

A new dedicated module `cli/heatmap/heatmap_tasks.go` will encapsulate all task inventory discovery:

```go
package heatmap

// TaskItem represents an architectural plan or telemetry parent task.
type TaskItem struct {
	TaskSequence int    `json:"taskSequence"`
	TaskSlug     string `json:"taskSlug"`
	Title        string `json:"title"`
	TotalActions int    `json:"totalActions"`
	IsActive     bool   `json:"isActive"`
	IsCompleted  bool   `json:"isCompleted"`
}

// TaskInventory holds categorized tasks for status reporting.
type TaskInventory struct {
	Completed []TaskItem `json:"completed"`
	Pending   []TaskItem `json:"pending"`
}
```

#### Discovery Algorithm:
1. **Filesystem Scan:**
   - Scan `.ai-memory/plans/completed/` to discover the 5 most recent completed plan markdown files.
   - Scan `.ai-memory/plans/pending/` to discover the top 5 scheduled or in-progress plan markdown files.
   - Extract numeric sequence prefixes (e.g., `326`), task slugs, and titles.
2. **SQLite Registry Enrichment:**
   - Query `HeatmapParentTask` table in `root-heatmap.db` to correlate action counts and active states.
   - If SQLite records exist for the task, enrich `TotalActions` and status.
3. **Graceful Fallback:** If `.ai-memory/plans/` is absent (e.g. fresh clone), fallback seamlessly to querying SQLite `HeatmapParentTask` records without error.

### 3.3 Actionable Tree-View Drill-Down Suggestions

The status view synthesizes copy-pasteable CLI commands to guide human operators and AI agents:
- **Active Task Tree:** `gitmap heatmap tree --task=<activeSlug>`
- **Error Search:** `gitmap heatmap search --state=ERROR --task=<activeSlug>`
- **Confusion Query:** `gitmap rc ai-confusion view --task=<activeSlug>`
- **Deep Historical Status:** `gitmap heatmap status --limit=20`

These commands are printed in a dedicated terminal sub-panel with syntax highlighting (`TokyoAccent`).

### 3.4 Machine JSON Output (`--json`)

When `gitmap heatmap status --json` or `-j` is supplied:
1. `runHeatmapStatus()` resolves the task inventory and status summary.
2. Constructs `HeatmapStatusJSONPayload` (defined in `cli/heatmap/types.go`).
3. Emits indented JSON to `io.Writer` via `writeIndentedJSON()`.
4. Bypasses all Lipgloss ANSI codes to guarantee pure machine ingestion.

```go
type HeatmapStatusJSONPayload struct {
	IsSuccess         bool                      `json:"isSuccess"`
	WorkspaceSlug     string                    `json:"workspaceSlug"`
	ActiveTaskSlug    string                    `json:"activeTaskSlug"`
	Telemetry         store.HeatmapStatusSummary `json:"telemetry"`
	CompletedTasks    []TaskItem                `json:"completedTasks"`
	PendingTasks      []TaskItem                `json:"pendingTasks"`
	ConfusionAlerts   []ConfusionAlertItem      `json:"confusionAlerts"`
	SuggestedCommands map[string]string         `json:"suggestedCommands"`
}
```

---

## 4. Execution Plan & Implementation Steps

### Step 1: Type Definitions (`cli/heatmap/types.go`)
- Add `TaskItem`, `TaskInventory`, `ConfusionAlertItem`, and `HeatmapStatusJSONPayload` types.
- Ensure all boolean fields use positive `is*` / `has*` naming (`isCompleted`, `isActive`, `hasPendingSubtasks`).
- Keep function lengths under 15 lines.

### Step 2: Task Aggregator (`cli/heatmap/heatmap_tasks.go`)
- Implement `DiscoverHeatmapTasks(repoRoot string, limit int) (TaskInventory, error)`.
- Implement `scanPlanDirectory(dirPath string, isCompleted bool, limit int) ([]TaskItem, error)`.
- Implement `enrichTaskWithDbMetrics(db *store.HeatmapDB, tasks []TaskItem) []TaskItem`.
- Guard all file system reads with existence checks and relative path resolution.

### Step 3: Terminal Card Renderer (`cli/heatmap/heatmap_cmd_render.go`)
- Refactor `renderStatusText(out io.Writer, s store.HeatmapStatusSummary)` to accept `TaskInventory` and `map[string]string` command suggestions.
- Implement modular helper renderers (under 15 lines each):
  - `renderStatusHeader(out io.Writer, taskSlug string)`
  - `renderTelemetrySummaryBox(out io.Writer, s store.HeatmapStatusSummary)`
  - `renderTaskQueue(out io.Writer, title string, tasks []TaskItem)`
  - `renderDrillDownSuggestions(out io.Writer, suggestions map[string]string)`
- Apply TokyoNight color tokens (`TokyoHeaderStyle`, `TokyoBoxStyle`, `TokyoOKPill`, `TokyoWarnPill`, `TokyoMutedStyle`).

### Step 4: CLI Dispatcher & JSON Wiring (`cli/heatmap/heatmap_cmd.go`)
- Update `runHeatmapStatus(flags map[string]string, out io.Writer) error`:
  - Open heatmap root database via `store.OpenHeatmapRootDB("")`.
  - Fetch telemetry summary via `db.GetHeatmapStatus(flags["task"])`.
  - Discover tasks via `DiscoverHeatmapTasks(".", 5)`.
  - Check for `--json` flag: if true, marshal and emit `HeatmapStatusJSONPayload`.
  - If false, call `renderStyledStatusDashboard(out, summary, inventory, suggestions)`.

---

## 5. Coding Guideline & Architectural Invariants

1. **Leaf Function Sovereignty:** `DiscoverHeatmapTasks` and all rendering helpers must return errors directly (`return err` or `return apperror.WrapSimple(...)`). Never call `os.Exit` or exit handlers.
2. **15-Line Function Limit:** No function in `cli/heatmap/` may exceed 15 lines. Decompose large rendering logic into discrete, single-responsibility helpers.
3. **Zero Nested Ifs:** All conditionals must be flattened with guard clauses and early returns.
4. **Return Styling:** Always insert a blank line before `return` statements and a blank line after closing braces `}`.
5. **Positive Booleans Only:** Only `is*` or `has*` prefixes permitted. Zero negative booleans.
6. **No-Build / No-Test Rule During Spec Turn:** Zero compilation or test commands executed during this turn.

---

## 6. Acceptance Criteria

- [ ] `cli/heatmap/types.go` contains `HeatmapStatusJSONPayload` with full schema support.
- [ ] `cli/heatmap/heatmap_tasks.go` discovers completed and pending tasks from `.ai-memory/plans/` and SQLite.
- [ ] `cli/heatmap/heatmap_cmd_render.go` renders Lipgloss TokyoNight rounded cards, padded boxes, and pill badges.
- [ ] `gitmap heatmap status` displays the last 5 completed tasks and top 5 pending tasks.
- [ ] `gitmap heatmap status` outputs actionable drill-down commands (e.g., `gitmap heatmap tree --task=<slug>`).
- [ ] `gitmap heatmap status --json` emits clean, valid JSON matching the architectural specification.
- [ ] All functions adhere to the 15-line limit and positive boolean naming conventions.
