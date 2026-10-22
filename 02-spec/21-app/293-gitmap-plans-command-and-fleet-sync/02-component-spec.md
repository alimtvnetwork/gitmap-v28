# Component Specification: Plan 293 — GitMap Plans Command and Fleet Skills Synchronization

## Status: APPROVED
- **Spec ID:** 293-02
- **Task Slug:** `293-gitmap-plans-command-and-fleet-sync`
- **Architecture Spec:** [01-architecture-spec.md](01-architecture-spec.md)
- **Migration & Verification Spec:** [03-migration-and-verification-spec.md](03-migration-and-verification-spec.md)
- **Master Plan:** [.ai-memory/plans/pending/293-gitmap-plans-command-and-fleet-sync.md](../../../.ai-memory/plans/pending/293-gitmap-plans-command-and-fleet-sync.md)

---

## 1. Data Models & Type Definitions (`cli/cmdplans/plans_types.go`)

### 1.1 Status Enums
```go
type PlanStatusType string

const (
	StatusCompleted  PlanStatusType = "Completed"
	StatusPending    PlanStatusType = "Pending"
	StatusInProgress PlanStatusType = "In Progress"
)
```

### 1.2 Plan Item Model
```go
type PlanItem struct {
	ID        int            `json:"id" yaml:"id"`
	Slug      string         `json:"slug" yaml:"slug"`
	Title     string         `json:"title" yaml:"title"`
	Status    PlanStatusType `json:"status" yaml:"status"`
	FilePath  string         `json:"filePath" yaml:"filePath"`
	UpdatedAt string         `json:"updatedAt,omitempty" yaml:"updatedAt,omitempty"`
}
```

### 1.3 Summary & Report Models
```go
type PlansSummary struct {
	CompletedCount   int    `json:"completedCount" yaml:"completedCount"`
	PendingCount     int    `json:"pendingCount" yaml:"pendingCount"`
	TotalCount       int    `json:"totalCount" yaml:"totalCount"`
	DisplayedLimit   int    `json:"displayedLimit" yaml:"displayedLimit"`
	HasMorePending   bool   `json:"hasMorePending" yaml:"hasMorePending"`
	HasMoreCompleted bool   `json:"hasMoreCompleted" yaml:"hasMoreCompleted"`
	PendingOverflow  string `json:"pendingOverflow,omitempty" yaml:"pendingOverflow,omitempty"`
}

type PlansReport struct {
	Summary   PlansSummary `json:"summary" yaml:"summary"`
	Completed []PlanItem   `json:"completed" yaml:"completed"`
	Pending   []PlanItem   `json:"pending" yaml:"pending"`
}
```

### 1.4 Command Execution Options
```go
type PlansOptions struct {
	Limit           int
	IsJSON          bool
	IsYAML          bool
	IsAll           bool
	IsPendingOnly   bool
	IsCompletedOnly bool
	RepoPath        string
}
```

---

## 2. Plan Parsing & Discovery Engine (`cli/cmdplans/plans_parser.go`)

### 2.1 File Location Traversal
- Target Directories:
  1. `<repo>/.ai-memory/plans/` (root)
  2. `<repo>/.ai-memory/plans/pending/`
  3. `<repo>/.ai-memory/plans/completed/` (if present)
- Subdirectories to skip: `subtasks`, `archive`, `.gitkeep`, non-markdown assets.
- Index/Meta files to skip: `01-index.md`, `ledger.md`, `readme.md`.

### 2.2 Sequence ID Extraction
- Regex candidate: `^(\d+)-(.*)\.md$`
- When matched: `ID = strconv.Atoi(submatch[1])`, `Slug = submatch[2]`
- Fallback: inspect line 1 for `# Plan (\d+):` or `# Master Plan (\d+):`. If missing, `ID = 0`.

### 2.3 Title Extraction
- Parse first markdown heading: `^#\s+(.+)$`.
- Clean extraneous prefixes (`Master Plan 291: ` or `Plan 227: ` retained as title or sanitized).
- Fallback: kebab-case slug converted to Title Case.

### 2.4 Status Extraction & Normalization
- Scan the first 35 lines of the markdown file for:
  - `## Status:\s*([A-Za-z\s_-]+)`
  - `\*\*Status:\*\*\s*([A-Za-z\s_-]+)`
  - `>\s*\*\*Status:\*\*\s*([A-Za-z\s_-]+)`
- Normalize:
  - Case-insensitive string match on `COMPLETED` / `Completed` / `Done` -> `StatusCompleted`.
  - Case-insensitive string match on `IN PROGRESS` / `In Progress` / `Active` -> `StatusInProgress`.
  - Case-insensitive string match on `PENDING` / `Pending` / `Queued` -> `StatusPending`.
- Directory heuristics:
  - If located in `pending/` and no explicit status header -> `StatusPending`.
  - If located in `completed/` and no explicit status header -> `StatusCompleted`.
  - If located in root `plans/` and no explicit status header -> `StatusCompleted`.

### 2.5 Sorting Criteria
- Completed plans: Sorted descending by `ID` (highest/latest first).
- Pending plans: Sorted ascending by `ID` (lowest/earliest/next action first).

---

## 3. Terminal UI Lipgloss Rendering (`cli/cmdplans/plans_render.go`)

### 3.1 Tokyo Night Aesthetic
- Utilizes `termout.TokyoHeader("GitMap Plans — Task & Plan Inventory")`.
- Pill metrics:
  - `termout.TokyoMetricPill("Completed", summary.CompletedCount)`
  - `termout.TokyoMetricPill("Pending", summary.PendingCount)`
  - `termout.TokyoMetricPill("Total", summary.TotalCount)`

### 3.2 Tables
- `termout.PrintTable` configured with:
  - `HasBorders: true`
  - `BorderStyle: termout.BorderRounded`
  - `RowStripe: true`
  - `HeaderColor: "#7aa2f7"`
  - `BorderColor: "#3b4261"`
  - `StripeColor: "#1f2335"`
- Completed Table columns: `ID` (right-aligned), `Title` (word/ellipsis), `Status` (`TokyoOKPill`), `File` (relative).
- Pending Table columns: `ID` (right-aligned), `Title` (word/ellipsis), `Status` (`TokyoWarnPill` or `TokyoTag`), `File` (relative).

### 3.3 Overflow Navigation Guidance
When `summary.HasMorePending` is true:
```
  ℹ Showing 10 of 28 pending plans.
  ➜ Use 'gitmap plans 28' or 'gitmap plans --all' to view the entire backlog.
```
When `summary.HasMoreCompleted` is true:
```
  ℹ Showing 10 of 75 completed plans.
  ➜ Use 'gitmap plans 75' or 'gitmap plans --all' to view all completed plans.
```

---

## 4. Machine-Readable Formatters (`cli/cmdplans/plans_format.go`)

- **JSON Mode (`--json`, `-j`)**: Serializes `PlansReport` with `json.MarshalIndent(report, "", "  ")` to `os.Stdout`.
- **YAML Mode (`--yaml`, `--yml`)**: Serializes `PlansReport` with `yaml.Marshal(report)` to `os.Stdout`.
- Error returns: Any marshalling failure returns wrapped `*apperror.AppError` directly.

---

## 5. Command Dispatch Wiring (`cli/cmd/rootutility.go`)

- Add `CmdPlans = "plans"` and `CmdPlansAlias = "plan"` to `cli/constants/constants_commands_extra.go`.
- Register in `utilityToolEntries()`:
```go
{[]string{constants.CmdPlans, constants.CmdPlansAlias}, func() error { return cmdplans.RunPlans(argsTail()) }},
```
