# Component Specification: Terminal UI, Tokyo Lipgloss Styling, Table Framing & Output Ergonomics

**Module:** `cli/termout`, `cli/cmdpipeline`, `cli/cmdcd`, `cli/cmdcg`  
**Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
**Specification File:** `02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md`  
**Subtask Range:** Subtasks 11 through 20 (Improvements 11 to 20)  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Component Scope

This component specification establishes the technical requirements, Go interface definitions, data structures, styling tokens, formatting contracts, and refactoring guidelines for Improvements 11 through 20 of the `291-improvements-planning-and-execution` milestone.

These improvements elevate GitMap's terminal user interface (TUI), standardizing Lipgloss visual primitives, introducing Tokyo Night key-value tags and metric pills, delivering subtle section dividers, adding visual-width column alignment to table rendering, ensuring rune-safe middle truncation, modernizing pipeline failure inspection cards, isolating interactive directory selection output strictly to standard error (`os.Stderr`), and enriching coding guideline doctor diagnostics.

### 1.1 Scope Boundaries

1. **`cli/termout/theme_lipgloss.go` (Improvement 11)**:
   - Implements `TokyoTag(label, value string) string` formatted with Tokyo Night styling tokens (`TokyoSurface`, `TokyoPrimary`, `TokyoFG`).
   - Ensures safe plain fallback `[label: value]` when colors are disallowed or output is not a TTY.
   - Enforces 15-line limit and zero nested conditionals.

2. **`cli/termout/theme_lipgloss.go` (Improvement 12)**:
   - Implements `TokyoMetricPill(label string, count int) string` displaying labeled counters with semantic color weighting (`TokyoOK`, `TokyoWarn`, `TokyoPrimary`).
   - Formats count metrics with bold emphasis and styled pill background.
   - Provides safe fallback for non-TTY / monochrome environments.

3. **`cli/termout/render.go` (Improvement 13)**:
   - Implements `FormatSectionDivider(title string, theme Theme) string` and `PrintSectionDivider(title string, theme Theme)`.
   - Uses Tokyo subtle line glyphs (`─`, `RoleFGSubtle`) to structure CLI help and command output sections.
   - Handles empty titles as continuous subtle divider rules.

4. **`cli/termout/termtable_types.go` (Improvement 14)**:
   - Introduces strongly-typed enum `AlignmentType` (`AlignLeft`, `AlignRight`, `AlignCenter`) adhering to the repository enum naming convention.
   - Defines canonical `ColumnConfig` struct containing `Alignment AlignmentType`, preserving `Column` and `AlignType` aliases for backward compatibility.

5. **`cli/termout/table.go` (Improvement 15)**:
   - Enhances table cell formatting with `AlignmentType` using `VisualStringWidth` to balance printable terminal widths.
   - Implements right alignment, center alignment, and left alignment across both headers and data rows.
   - Decomposes alignment logic into discrete sub-15-line helper routines.

6. **`cli/termout/truncate.go` (Improvement 16)**:
   - Implements `TruncateMiddleWithEllipsis(text string, maxWidth int) string`.
   - Replaces byte-based slicing with rune-safe, visual-width-aware middle truncation using `"..."` ellipsis.
   - Safely preserves Unicode runes and multi-byte UTF-8 sequences.

7. **`cli/cmdpipeline/pipeline_history_render.go` (Improvement 17)**:
   - Modernizes failure cards by replacing raw ASCII box borders with `termout.TokyoBoxStyle`.
   - Embeds `termout.TokyoTag` for metadata metrics (`Run`, `Workflow`, `Branch`, `SHA`, `Duration`, `Date`).
   - Refactors diagnostic rendering into clean, single-responsibility functions.

8. **`cli/cmdcd/cd_table.go` (Improvement 18)**:
   - Creates a dedicated interactive candidate selection table renderer for `gitmap cd`.
   - Routes all table framing, headers, and rows strictly to `os.Stderr` to avoid corrupting shell wrappers (e.g. `cd $(gitmap cd ...)`).
   - Formats candidate index, repository name, branch, and middle-truncated path with `termout.BorderRounded`.

9. **`cli/cmdcd/cdops.go` (Improvement 19)**:
   - Wires candidate selection in `promptCDPick` and `promptCDReposPick` to use `cd_table.go`.
   - Preserves stdout cleanliness for destination directory path output.
   - Replaces unaligned `fmt.Fprintf` iteration with the standardized CD selection table.

10. **`cli/cmdcg/cg_doctor.go` (Improvement 20)**:
    - Integrates `termout.PrintSectionDivider` for visual header demarcation.
    - Enhances summary reporting with `termout.TokyoMetricPill` for total, healthy, and flagged repository counts.
    - Packages overall diagnostics inside `termout.TokyoBoxStyle`.

---

## 2. Shared Data Models & Go Contracts

### 2.1 Table Alignment & Column Configuration (`cli/termout/termtable_types.go`)

```go
package termout

// AlignmentType defines column text alignment for terminal tables.
type AlignmentType string

const (
	// AlignLeft positions text flush with the left margin.
	AlignLeft AlignmentType = "left"
	// AlignRight positions text flush with the right margin (ideal for numeric metrics).
	AlignRight AlignmentType = "right"
	// AlignCenter centers text symmetrically within the column width.
	AlignCenter AlignmentType = "center"
)

// AlignType is maintained as a type alias for AlignmentType for backward compatibility.
type AlignType = AlignmentType

// ColumnConfig defines configuration options for an individual table column.
type ColumnConfig struct {
	Title     string
	MaxWidth  int
	MinWidth  int
	Alignment AlignmentType
	Align     AlignmentType // Alias for backward compatibility
}

// Column is maintained as a type alias for ColumnConfig for backward compatibility.
type Column = ColumnConfig
```

### 2.2 Tokyo Lipgloss Tags & Metric Pills (`cli/termout/theme_lipgloss.go`)

```go
package termout

// TokyoTag renders a key-value pair formatted with Tokyo Night tokens.
// In non-TTY or no-color environments, returns "[label: value]".
func TokyoTag(label, value string) string

// TokyoMetricPill renders a labeled count metric with Tokyo Night styling.
// In non-TTY or no-color environments, returns "[label: count]".
func TokyoMetricPill(label string, count int) string
```

### 2.3 Section Divider API (`cli/termout/render.go`)

```go
package termout

// FormatSectionDivider builds a horizontal subtle line divider with an optional section title.
func FormatSectionDivider(title string, theme Theme) string

// PrintSectionDivider outputs FormatSectionDivider directly to output.UI().
func PrintSectionDivider(title string, theme Theme)
```

### 2.4 Rune-Safe Visual Truncation API (`cli/termout/truncate.go`)

```go
package termout

// TruncateMiddleWithEllipsis truncates text exceeding maxWidth by placing a "..." ellipsis
// in the middle, calculated via rune visual display width to avoid character tearing.
func TruncateMiddleWithEllipsis(text string, maxWidth int) string
```

### 2.5 Dedicated CD Selection Table API (`cli/cmdcd/cd_table.go`)

```go
package cmdcd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// RenderCDSelectionTable outputs an interactive repository selection table strictly to os.Stderr.
func RenderCDSelectionTable(title string, records []model.ScanRecord)
```

---

## 3. Technical Specifications by Component

### 3.1 Subtask 11: TokyoTag Lipgloss Token Implementation

**Target File:** `cli/termout/theme_lipgloss.go`

#### Functional Requirements
- Signature: `TokyoTag(label, value string) string`.
- When `!colorAllowed()`:
  - If `value == ""`, return `fmt.Sprintf("[%s]", label)`.
  - Otherwise, return `fmt.Sprintf("[%s: %s]", label, value)`.
- When colors are enabled:
  - Format `label` with `TokyoFG` text on `TokyoSurface` background with horizontal padding (0, 1), bold text.
  - Format `value` with `TokyoPrimary` text on `TokyoBase` background with horizontal padding (0, 1).
  - Return the concatenated styled blocks.
- If `value == ""`, format single label pill with `TokyoPrimary` on `TokyoSurface`.

#### Invariant Governance
- Function length $\le 15$ lines.
- Zero nested `if` statements.

---

### 3.2 Subtask 12: TokyoMetricPill Component

**Target File:** `cli/termout/theme_lipgloss.go`

#### Functional Requirements
- Signature: `TokyoMetricPill(label string, count int) string`.
- When `!colorAllowed()`:
  - Return `fmt.Sprintf("[%s: %d]", label, count)`.
- When colors are enabled:
  - Resolve metric status color:
    - If `strings.EqualFold(label, "healthy") || strings.EqualFold(label, "pass") || strings.EqualFold(label, "ok")`: use `TokyoOK`.
    - If `count > 0 && (strings.EqualFold(label, "issues") || strings.EqualFold(label, "errors") || strings.EqualFold(label, "fail"))`: use `TokyoWarn` (or `TokyoError`).
    - Default: `TokyoPrimary`.
  - Styled with Lipgloss: foreground resolved color, background `TokyoSurface`, padding (0, 1).
  - Text format: `fmt.Sprintf("%s: %d", label, count)`.

#### Invariant Governance
- Function decomposed into `resolveMetricColor(label string, count int) string` and `TokyoMetricPill`.
- Both functions $\le 15$ lines.

---

### 3.3 Subtask 13: Section Dividers (`FormatSectionDivider` / `PrintSectionDivider`)

**Target File:** `cli/termout/render.go`

#### Functional Requirements
- Signatures:
  - `FormatSectionDivider(title string, theme Theme) string`
  - `PrintSectionDivider(title string, theme Theme)`
- Default divider width: 64 columns.
- Divider glyph: `─` (box drawings light horizontal `\u2500`).
- If `title == ""`:
  - Format line of 64 `─` glyphs colored with `RoleFGSubtle`.
- If `title != ""`:
  - Leading line: 3 `─` glyphs (`───`).
  - Spaced title: ` ` + `title` + ` `.
  - Trailing line: fill remaining width up to 64 columns with `─`.
- If `!colorAllowed()`:
  - Output uncolored divider line.
- `PrintSectionDivider` outputs the formatted divider to `output.UI()`.

---

### 3.4 Subtask 14: Table Alignment Enum & ColumnConfig

**Target File:** `cli/termout/termtable_types.go`

#### Functional Requirements
- Define `AlignmentType` enum with values `AlignLeft`, `AlignRight`, `AlignCenter`.
- Provide type alias `type AlignType = AlignmentType` to guarantee non-breaking compatibility.
- Define `ColumnConfig` struct with fields `Title`, `MaxWidth`, `MinWidth`, `Alignment`, and `Align`.
- Provide type alias `type Column = ColumnConfig`.

---

### 3.5 Subtask 15: Table Column Alignment Rendering

**Target File:** `cli/termout/table.go`

#### Functional Requirements
- Update `getColAlign(cols []Column, idx int) AlignmentType` to resolve `cols[idx].Alignment` or `cols[idx].Align` (defaulting to `AlignLeft`).
- In `formatCell(text string, width int, align AlignmentType, ellipsis string) string`:
  - Truncate text using `TruncateMiddle(text, width, ellipsis)`.
  - Calculate visual printable width with `VisualStringWidth(truncated)`.
  - Compute padding length `padLen := width - visWidth`.
  - If `padLen <= 0`, return `truncated`.
  - If `align == AlignRight`: return `SpacePool(padLen) + truncated`.
  - If `align == AlignCenter`: split padding into `left := padLen / 2` and `right := padLen - left`, returning `SpacePool(left) + truncated + SpacePool(right)`.
  - Default (`AlignLeft`): return `truncated + SpacePool(padLen)`.
- Update `formatHeader` to apply the column's alignment to header titles as well.

---

### 3.6 Subtask 16: Visual-Width Middle Truncation (`TruncateMiddleWithEllipsis`)

**Target File:** `cli/termout/truncate.go`

#### Functional Requirements
- Signature: `TruncateMiddleWithEllipsis(text string, maxWidth int) string`.
- Guard conditions:
  - If `maxWidth <= 0`, return `""`.
  - If `VisualStringWidth(text) <= maxWidth`, return `text`.
  - If `maxWidth <= 3`, return `"..."[:maxWidth]`.
- Calculation:
  - Ellipsis: `"..."` (width 3).
  - Target available width: `avail := maxWidth - 3`.
  - Prefix visual budget: `prefixBudget := (avail + 1) / 2`.
  - Suffix visual budget: `suffixBudget := avail - prefixBudget`.
- Rune slicing:
  - Convert text to `[]rune`.
  - Traverse runes from left accumulating display width up to `prefixBudget`.
  - Traverse runes from right accumulating display width up to `suffixBudget`.
  - Join: `prefix + "..." + suffix`.
- Guarantees zero byte slicing panics and zero UTF-8 rune corruption.

---

### 3.7 Subtask 17: Pipeline Failure Cards Modernization

**Target File:** `cli/cmdpipeline/pipeline_history_render.go`

#### Functional Requirements
- Refactor `renderSingleCachedFailureCard(db *pipelinedb.PipelineSplitDb, r pipelinedb.PipelineRunRecord, idx, total int, isDetailed bool)`.
- Replace legacy raw ASCII box drawings (`┌─`, `│`, `└─`) with `termout.TokyoBoxStyle`.
- Header section:
  - `termout.TokyoTag("FAIL", fmt.Sprintf("[%d/%d] #%d", idx, total, r.RunId))`
  - Workflow title bolded with Tokyo styling.
- Metadata metrics section:
  - Branch tag: `termout.TokyoTag("Branch", r.Branch)`
  - SHA tag: `termout.TokyoTag("SHA", r.Sha)`
  - Duration tag: `termout.TokyoTag("Duration", formatDurationSeconds(r.DurationSeconds))`
  - Date tag: `termout.TokyoTag("Date", formatRunTimestamp(r.CreatedAt))`
- Diagnostics section:
  - Steps and errors formatted cleanly with `termout.TokyoErrorPill("STEP")`.
- Assemble card string and print via `fmt.Println(termout.TokyoBoxStyle(content))`.

---

### 3.8 Subtask 18: Dedicated CD Selection Table (`cd_table.go`)

**Target File:** `cli/cmdcd/cd_table.go`

#### Functional Requirements
- Create new package file `cli/cmdcd/cd_table.go`.
- Implements `RenderCDSelectionTable(title string, records []model.ScanRecord)`.
- CRITICAL REQUIREMENT: Output must be written strictly to `os.Stderr`.
- Construct `termout.TableConfig`:
  - Columns:
    - `#`: `MinWidth: 4`, `MaxWidth: 4`, `Alignment: termout.AlignRight`
    - `Repository`: `MinWidth: 16`, `MaxWidth: 24`, `Alignment: termout.AlignLeft`
    - `Branch`: `MinWidth: 10`, `MaxWidth: 16`, `Alignment: termout.AlignLeft`
    - `Path`: `MinWidth: 30`, `MaxWidth: 56`, `Alignment: termout.AlignLeft`
  - Style: `BorderStyle: termout.BorderRounded`.
  - HeaderColor: `termout.RoleSGR(termout.RolePrimary)`.
- Row construction:
  - Truncate long absolute paths using `termout.TruncateMiddleWithEllipsis(r.AbsolutePath, 56)`.
- Render table string via `termout.RenderTable(cfg)` and write directly: `fmt.Fprint(os.Stderr, rendered)`.

---

### 3.9 Subtask 19: Wire CD Selection Table in `cdops.go`

**Target File:** `cli/cmdcd/cdops.go`

#### Functional Requirements
- Update `promptCDPick(name string, records []model.ScanRecord) (string, error)`:
  - Replace raw loop with `RenderCDSelectionTable(fmt.Sprintf("Multiple repositories matching '%s'", name), records)`.
  - Output pick prompt to `os.Stderr`: `fmt.Fprintf(os.Stderr, "\nPick [1-%d] (default 1): ", len(records))`.
- Update `promptCDReposPick(records []model.ScanRecord) (string, error)`:
  - Replace raw loop with `RenderCDSelectionTable("Available Repositories", records)`.
  - Output pick prompt to `os.Stderr`: `fmt.Fprintf(os.Stderr, constants.MsgCDPickPrompt, len(records))`.
- Guarantees standard output captures only the destination path.

---

### 3.10 Subtask 20: Coding Guideline Doctor Summary Enhancement

**Target File:** `cli/cmdcg/cg_doctor.go`

#### Functional Requirements
- Update `RunCGDoctor(repos []string) error`:
  - Replace raw `=== Coding Guidelines Doctor Diagnostics ===` with:
    `termout.PrintSectionDivider("Coding Guidelines Doctor Diagnostics", termout.DefaultTheme())`
- Update `buildDoctorSummaryContent(totalRepos, healthyRepos int) string`:
  - Build metrics summary bar using:
    - `termout.TokyoMetricPill("Total", totalRepos)`
    - `termout.TokyoMetricPill("Healthy", healthyRepos)`
    - `termout.TokyoMetricPill("Issues", totalRepos - healthyRepos)`
  - Retain `termout.RenderTokyoProgress(healthyRepos, totalRepos, 24)` inline.
  - Return composite message wrapped in `termout.TokyoBoxStyle`.

---

## 4. Coding Guidelines & Invariant Governance

All implementations authored under this specification must strictly comply with the following invariants:

1. **Error Sovereignty**: Leaf functions returning an `error` must bubble it up. No `os.Exit` or swallowing errors.
2. **15-Line Function Limit**: Functions must not exceed 15 lines of executable code (excluding comments/blank lines). Functions $\le 8$ lines are preferred.
3. **Zero Nested Conditionals**: All nested `if` statements are forbidden. Flatten all branches using guard clauses, inverted conditions, and early returns.
4. **Positive Boolean Naming**: Boolean variables, functions, and parameters must strictly begin with `is` or `has` (e.g. `isDetailed`, `hasBorders`). Negative booleans (`isNotValid`, `hasNoMatch`) are banned.
5. **Vertical Spacing & Formatting**:
   - Exactly one blank line before `if` statements (unless directly at block start).
   - Exactly one blank line after closing brace `}`.
   - Exactly one blank line before `return` statements.
6. **Zero Compiler Builds and Tests**: No `go build` or `go test` invocations during planning.

---

## 5. Implementation Dependencies & File Boundaries

| Subtask | Target File | Package | Primary Symbol / Feature |
| :--- | :--- | :--- | :--- |
| **11** | `cli/termout/theme_lipgloss.go` | `termout` | `TokyoTag(label, value string) string` |
| **12** | `cli/termout/theme_lipgloss.go` | `termout` | `TokyoMetricPill(label string, count int) string` |
| **13** | `cli/termout/render.go` | `termout` | `FormatSectionDivider`, `PrintSectionDivider` |
| **14** | `cli/termout/termtable_types.go` | `termout` | `AlignmentType`, `ColumnConfig` |
| **15** | `cli/termout/table.go` | `termout` | Visual column alignment rendering in `formatCell` |
| **16** | `cli/termout/truncate.go` | `termout` | `TruncateMiddleWithEllipsis(text string, maxWidth int) string` |
| **17** | `cli/cmdpipeline/pipeline_history_render.go` | `cmdpipeline` | `TokyoBoxStyle` + `TokyoTag` in failure cards |
| **18** | `cli/cmdcd/cd_table.go` | `cmdcd` | `RenderCDSelectionTable` printing to `os.Stderr` |
| **19** | `cli/cmdcd/cdops.go` | `cmdcd` | Wiring `promptCDPick` / `promptCDReposPick` to `cd_table.go` |
| **20** | `cli/cmdcg/cg_doctor.go` | `cmdcg` | `TokyoMetricPill` + section dividers in doctor diagnostics |
