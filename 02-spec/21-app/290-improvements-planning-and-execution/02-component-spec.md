# Component Specification: Terminal UI, Table Framing, Tokyo Lipgloss & Output Standardization

**Module:** `cli/termout`, `cli/cmdcg`, `cli/cmdhistory`, `cli/cmdpending`, `cli/cmdclone`  
**Parent Plan:** `.ai-memory/plans/290-improvements-planning-and-execution.md`  
**Specification File:** `02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md`  
**Subtask Range:** Subtasks 11 through 20 (Improvements 11 to 20)  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Component Scope

This component specification establishes the technical requirements, Go interface definitions, data structures, styling tokens, formatting contracts, and refactoring guidelines for Improvements 11 through 20 of the `290-improvements-planning-and-execution` milestone.

These improvements focus on standardizing GitMap's terminal user interface (TUI), unifying table border styling, eliminating fragmented formatting routines, establishing Tokyo Night theme tokens across CLI commands, and enforcing repository-wide coding guidelines.

### 1.1 Scope Boundaries

1. **`cli/termout/theme_lipgloss.go` (Improvement 11)**:
   - Introduces `TokyoBannerStyle` and `TokyoBoxStyle` for standardized terminal frames.
   - Implements `TokyoPill`, `TokyoOKPill`, `TokyoWarnPill`, and `TokyoErrorPill` badge formatters.
   - Adds `RenderTokyoProgress` for proportional truecolor progress bars.
   - Enforces strict TTY and `colorAllowed()` safety checks.

2. **`cli/termout/render.go` (Improvement 12)**:
   - Implements `PrintKeyValPadded` and `FormatKeyValPadded` for uniform key-value pair output.
   - Integrates `termout.PadRightMax` for space-pool backed alignment.
   - Encapsulates semantic role colors (`RolePrimary`, `RoleFGSubtle`, `RoleFG`).

3. **`cli/termout/termtable_types.go` (Improvement 13)**:
   - Introduces `BorderStyleType` enum with PascalCase constants (`BorderRounded`, `BorderSharp`, `BorderMinimal`).
   - Extends `TableConfig` with `BorderStyle BorderStyleType` while preserving backward compatibility.

4. **`cli/termout/table.go` (Improvement 14)**:
   - Implements rounded (`╭`, `╮`, `╯`, `╰`, `├`, `┤`, `┬`, `┴`, `┼`), sharp (`┌`, `┐`, `┘`, `└`), and minimal frame border rendering.
   - Defaults border lines to subtle Tokyo Night palette tokens (`TokyoMuted` / `TokyoSurface`).
   - Decomposes table rendering into discrete sub-15-line functions.

5. **`cli/termout/truncate.go` (Improvement 15)**:
   - Implements `TruncateRightWithEllipsis(text string, maxWidth int) string`.
   - Guarantees rune-safe, UTF-8 aware display width truncation using `VisualStringWidth`.
   - Prevents slice out-of-bounds panics and multi-byte character corruption.

6. **`cli/cmdcg/cg_audit.go` (Improvement 16)**:
   - Purges hardcoded legacy Dracula hex color literals (`#8be9fd`, `#50fa7b`, `#ffb86c`, `#ff5555`, `#6272a4`).
   - Binds all audit styles to `termout.SafeTokyoStyle` palette tokens.
   - Standardizes compliance pill badges and status indicators.

7. **`cli/cmdcg/cg_doctor.go` (Improvement 17)**:
   - Integrates `termout.RenderTokyoProgress` for visual health ratio reporting.
   - Wraps summary diagnostics in `termout.TokyoBoxStyle` panels.
   - Refactors diagnostic renderers into clean, single-responsibility functions.

8. **`cli/cmdhistory/history.go` (Improvement 18)**:
   - Eliminates redundant local formatting routines (`padRight`, `truncateHist`, `colorize`).
   - Rebinds column alignment to `termout.PadRightMax`, truncation to `termout.TruncateRightWithEllipsis`, and styling to `termout.RoleSGR`.

9. **`cli/cmdpending/pending_commits_render.go` (Improvement 19)**:
   - Purges private `padRight` implementation.
   - Replaces hardcoded status strings with `termout.TokyoPill` / `termout.TokyoOKPill` / `termout.TokyoWarnPill`.
   - Unifies terminal framing with `termout.PadRightMax`.

10. **`cli/cmdclone/clone_list_table.go` (Improvement 20)**:
    - Replaces crude `fmt.Printf` / ASCII dash-line table with `termout.PrintTable`.
    - Configures structured columns with `BorderRounded`, column alignment, and Tokyo styling tokens.
    - Applies `termout.TokyoOKPill` for repository status badges.

---

## 2. Shared Data Models & Go Contracts

### 2.1 Border Style Enum & Table Configuration (`cli/termout/termtable_types.go`)

```go
package termout

// BorderStyleType defines the visual framing style for rendered terminal tables.
type BorderStyleType string

const (
	// BorderRounded uses curved UTF-8 box-drawing characters (╭, ╮, ╯, ╰).
	BorderRounded BorderStyleType = "rounded"
	// BorderSharp uses classic rectangular box-drawing characters (┌, ┐, ┘, └).
	BorderSharp BorderStyleType = "sharp"
	// BorderMinimal uses understated separator lines without outer vertical bars.
	BorderMinimal BorderStyleType = "minimal"
)

// TableConfig defines full configuration for table rendering.
type TableConfig struct {
	Columns      []Column
	Rows         []Row
	HeaderColor  string
	BorderColor  string
	HasBorders   bool
	BorderStyle  BorderStyleType
	EllipsisText string
}
```

### 2.2 Border Box Glyphs Specification

To maintain terminal consistency across operating systems (Linux, macOS, Windows Terminal), the following UTF-8 glyph sets are standardized:

| Border Style | Top Left | Top Right | Bottom Left | Bottom Right | Horiz | Vert | Cross | T-Down | T-Up | T-Right | T-Left |
|---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `BorderRounded` | `╭` | `╮` | `╰` | `╯` | `─` | `│` | `┼` | `┬` | `┴` | `├` | `┤` |
| `BorderSharp`   | `┌` | `┐` | `└` | `┘` | `─` | `│` | `┼` | `┬` | `┴` | `├` | `┤` |
| `BorderMinimal` | ` ` | ` ` | ` ` | ` ` | `─` | ` ` | `─` | `─` | `─` | ` ` | ` ` |

---

## 3. Detailed Component Specifications (Improvements 11–20)

### 3.1 Component 11: Tokyo Night Lipgloss Styling Subsystem (`cli/termout/theme_lipgloss.go`)

#### 3.1.1 Objective
Expand `cli/termout/theme_lipgloss.go` with declarative, reusable Lipgloss styles and helper routines adhering to the Tokyo Night design palette (`TokyoPrimary`, `TokyoAccent`, `TokyoOK`, `TokyoWarn`, `TokyoError`, `TokyoMuted`, `TokyoSurface`).

#### 3.1.2 Interface & Function Declarations

```go
package termout

import "github.com/charmbracelet/lipgloss"

var (
	// TokyoBannerStyle renders prominent headers with border padding.
	TokyoBannerStyle = SafeTokyoStyle(TokyoPrimary).
		Bold(true).
		Padding(0, 1)

	// TokyoBoxStyle renders subtle bordered container cards.
	TokyoBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(TokyoMuted)).
		Padding(0, 1)
)

// TokyoPill formats text inside a colored pill badge with background and foreground colors.
func TokyoPill(text, fgHex, bgHex string) string

// TokyoOKPill formats a success status badge using TokyoOK foreground.
func TokyoOKPill(text string) string

// TokyoWarnPill formats a warning status badge using TokyoWarn foreground.
func TokyoWarnPill(text string) string

// TokyoErrorPill formats an error/failure status badge using TokyoError foreground.
func TokyoErrorPill(text string) string

// RenderTokyoProgress renders an inline ASCII/UTF-8 progress bar with Tokyo truecolor styling.
func RenderTokyoProgress(current, total, barWidth int) string
```

#### 3.1.3 Implementation Invariants
1. **TTY Safety**: Every pill and banner helper MUST query `colorAllowed()`. When color output is disabled (e.g. piped stdout or `NO_COLOR=1`), badges must degrade gracefully to plain bracketed strings: `fmt.Sprintf("[%s]", text)`.
2. **Progress Bar Algorithm**:
   - Guard against invalid arguments: If `total <= 0` or `barWidth <= 0`, return an empty or zero-percent bar `[░░░░] 0%`.
   - Clamp `current`: If `current > total`, clamp `current = total`; if `current < 0`, clamp `current = 0`.
   - Calculate filled cells: `filled = (current * barWidth) / total`.
   - Calculate percentage: `pct = (current * 100) / total`.
   - Build bar using filled block `█` (styled with `TokyoOK` or `TokyoPrimary`) and unfilled block `░` (styled with `TokyoMuted`).
3. **Function Sizing**: All functions must be $\le 15$ lines. Helper routines must be extracted for percentage formatting and segment construction.

---

### 3.2 Component 12: Unified Key-Value Output Engine (`cli/termout/render.go`)

#### 3.2.1 Objective
Eliminate scattered `fmt.Printf("  %-20s: %s\n", key, val)` statements throughout CLI commands by introducing a centralized, theme-aware key-value formatter.

#### 3.2.2 Interface & Function Declarations

```go
package termout

// PrintKeyValPadded outputs a padded key-value pair directly to output.UI().
func PrintKeyValPadded(key string, val any, keyWidth int)

// FormatKeyValPadded builds a padded key-value string with semantic role colors.
func FormatKeyValPadded(key string, val any, keyWidth int) string
```

#### 3.2.3 Implementation Details
- The formatted key must include the colon suffix (`key + ":"`).
- Uses `termout.PadRightMax(key+":", keyWidth)` to ensure deterministic alignment backed by the zero-allocation space pool.
- Key styling: `RoleSGRBold(RolePrimary)`.
- Value styling: `RoleSGR(RoleFG)` or neutral text.
- If colors are disabled, renders clean unstyled text: `fmt.Sprintf("  %-*s %v", keyWidth, key+":", val)`.
- Line length limit: maximum 12 lines per function.

---

### 3.3 Component 13: Border Style Enum Specification (`cli/termout/termtable_types.go`)

#### 3.3.1 Objective
Provide strongly typed border styling options for `termout.TableConfig` to replace ad-hoc boolean flags and magic strings.

#### 3.3.2 Requirements
- Define `BorderStyleType` type ending in `Type` as required by GitMap coding guidelines.
- Define constants: `BorderRounded`, `BorderSharp`, `BorderMinimal`.
- Add field `BorderStyle BorderStyleType` to `TableConfig`.
- Maintain backward compatibility: When `BorderStyle == ""`, table framing falls back to standard line rendering if `HasBorders == true`.

---

### 3.4 Component 14: Structured Table Framing Engine (`cli/termout/table.go`)

#### 3.4.1 Objective
Enhance `RenderTable` and `PrintTable` in `cli/termout/table.go` to support complete outer and inner grid framing based on `TableConfig.BorderStyle`.

#### 3.4.2 Framing Architecture
1. **Top Border**: `╭──────┬──────┬──────╮` (Rounded) or `┌──────┬──────┬──────┐` (Sharp).
2. **Header Row**: `│ Name │ Type │ Note │`.
3. **Separator Row**: `├──────┼──────┼──────┤` (Inner T-junctions).
4. **Data Rows**: `│ item │ data │ text │`.
5. **Bottom Border**: `╰──────┴──────┴──────╯` (Rounded) or `└──────┴──────┴──────┘` (Sharp).
6. **Minimal Mode**: Omits outer vertical pipes (`│`); renders top and bottom lines as subtle separators.

#### 3.4.3 Decomposition Plan
To adhere to the 15-line function limit:
- `renderTableFramed(cfg TableConfig, widths []int) string` ($\le 12$ lines)
- `buildTopBorder(cfg TableConfig, widths []int) string` ($\le 10$ lines)
- `buildHeaderBorder(cfg TableConfig, widths []int) string` ($\le 10$ lines)
- `buildRowBorder(cfg TableConfig, widths []int) string` ($\le 10$ lines)
- `buildBottomBorder(cfg TableConfig, widths []int) string` ($\le 10$ lines)
- `resolveBorderGlyphs(style BorderStyleType) borderGlyphSet` ($\le 12$ lines)

---

### 3.5 Component 15: UTF-8 Aware Right-Truncation (`cli/termout/truncate.go`)

#### 3.5.1 Objective
Replace naive byte slicing with a UTF-8 safe, display-width aware right-truncation function in `cli/termout/truncate.go`.

#### 3.5.2 Function Signature & Contract

```go
package termout

// TruncateRightWithEllipsis shortens text exceeding maxWidth by trimming characters from the right
// and appending an ellipsis. It accounts for multi-byte UTF-8 runes and double-width terminal characters.
func TruncateRightWithEllipsis(text string, maxWidth int) string
```

#### 3.5.3 Algorithm Specification
1. If `maxWidth <= 0`, return `""`.
2. Compute visual width: `currentWidth := VisualStringWidth(text)`.
3. If `currentWidth <= maxWidth`, return `text` untouched.
4. Ellipsis token is `…` (visual width 1) or `...` (visual width 3). Standardize on `…` when visual width allows, or `...` if width is sufficient.
5. If `maxWidth <= 1`, return `"…"` (or `"."` if single width).
6. Target rune display width: `availWidth := maxWidth - 1` (reserving 1 column for `…`).
7. Iterate through runes in `text`, accumulating runes while `accumWidth + runeWidth <= availWidth`.
8. Return accumulated rune string + `"…"`.
9. Enforce zero nested loops and max 15 lines per function.

---

### 3.6 Component 16: Coding Guidelines Audit Theme Modernization (`cli/cmdcg/cg_audit.go`)

#### 3.6.1 Objective
Purge obsolete Dracula hex codes from `cli/cmdcg/cg_audit.go` and bind all audit diagnostic renderers to `termout.SafeTokyoStyle` tokens.

#### 3.6.2 Migration Mapping

| Old Dracula Style | Old Hex | New Tokyo Night Replacement | Token Constant |
|---|---|---|---|
| `cgAuditTitleStyle` | `#8be9fd` | `termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true)` | `termout.TokyoPrimary` |
| `cgPassStyle` | `#50fa7b` | `termout.SafeTokyoStyle(termout.TokyoOK).Bold(true)` | `termout.TokyoOK` |
| `cgWarnStyle` | `#ffb86c` | `termout.SafeTokyoStyle(termout.TokyoWarn).Bold(true)` | `termout.TokyoWarn` |
| `cgFailStyle` | `#ff5555` | `termout.SafeTokyoStyle(termout.TokyoError).Bold(true)` | `termout.TokyoError` |
| `cgMutedStyle` | `#6272a4` | `termout.SafeTokyoStyle(termout.TokyoMuted)` | `termout.TokyoMuted` |

#### 3.6.3 Audit Score Rendering
- Pass checks: Render with `termout.TokyoOKPill("PASS")` or `cgPassStyle.Render("✔")`.
- Warning checks: Render with `termout.TokyoWarnPill("WARN")` or `cgWarnStyle.Render("!")`.
- Fail checks: Render with `termout.TokyoErrorPill("FAIL")` or `cgFailStyle.Render("✖")`.
- Audit summary score: Display compliance percentage using `termout.RenderTokyoProgress(passed, total, 20)`.

---

### 3.7 Component 17: Coding Guidelines Doctor Diagnostics Enhancements (`cli/cmdcg/cg_doctor.go`)

#### 3.7.1 Objective
Enhance `cli/cmdcg/cg_doctor.go` with structured visual progress metrics and Tokyo-themed diagnostic box containers.

#### 3.7.2 Refactoring Plan
1. In `RunCGDoctor(repos []string) error`:
   - Compute aggregate health metrics: `healthyCount`, `totalRepos`.
   - Render diagnostic summary block wrapped in `termout.TokyoBoxStyle`.
   - Display overall repository health progress:
     ```
     Overall Health: [██████████████░░░░░░] 70% (7/10 Repositories Healthy)
     ```
2. In `renderDoctorDiagnosis(diag CGDoctorDiagnosis)`:
   - Replace raw indentation with `termout.PrintKeyValPadded`.
   - Format issue labels with `termout.TokyoWarnPill("ISSUE")` and fix recommendations with `termout.TokyoOKPill("FIX")`.
3. Preserve leaf error sovereignty: `RunCGDoctor` returns `nil` or wrapped `error` without calling exit handlers directly.

---

### 3.8 Component 18: History CLI Output Consolidation (`cli/cmdhistory/history.go`)

#### 3.8.1 Objective
Eliminate duplicate formatting routines in `cli/cmdhistory/history.go` and bind column formatting to `termout` core utilities.

#### 3.8.2 Obsolete Routines to Remove
- `padRight(s string, w int) string` (lines 350–356) $\rightarrow$ delete; use `termout.PadRightMax(s, w)`.
- `truncateHist(s string, w int) string` (lines 358–368) $\rightarrow$ delete; use `termout.TruncateRightWithEllipsis(s, w)`.
- `colorize(color, s string) string` (lines 370–376) $\rightarrow$ delete; use `termout.RoleSGR` / `termout.SafeTokyoStyle`.

#### 3.8.3 Status Pill Migration
In `colorizedStatus(code int) string`:
```go
func colorizedStatus(code int) string {
	if code == 0 {
		return termout.TokyoOKPill("✓ OK")
	}

	return termout.TokyoErrorPill("✗ FAIL")
}
```

---

### 3.9 Component 19: Pending Commits Output Standardization (`cli/cmdpending/pending_commits_render.go`)

#### 3.9.1 Objective
Purge local padding helper from `cli/cmdpending/pending_commits_render.go` and standardize table status badges on `termout.TokyoPill`.

#### 3.9.2 Obsolete Routines to Remove
- `padRight(s string, width int) string` (lines 141–147) $\rightarrow$ delete; use `termout.PadRightMax(s, width)`.

#### 3.9.3 Status Pill Integration
Replace ad-hoc status color strings:
- Old pending status: `constants.ColorYellow + "● PEND   " + reset` $\rightarrow$ `termout.TokyoWarnPill("PEND")`.
- Old clean status: `constants.ColorGreen + "○ CLEAN  " + reset` $\rightarrow$ `termout.TokyoOKPill("CLEAN")`.
- Box boundary alignment: Replace `padRight(summaryText, 78)` with `termout.PadRightMax(summaryText, 78)`.

---

### 3.10 Component 20: Clone List Table Modernization (`cli/cmdclone/clone_list_table.go`)

#### 3.10.1 Objective
Migrate crude fixed-width ASCII printing in `printCloneRecordsTable` to declarative `termout.PrintTable` with `BorderRounded`.

#### 3.10.2 Table Schema Definition
```go
func printCloneRecordsTable(records []model.ScanRecord, targetDir string, isSSH bool) {
	cfg := termout.TableConfig{
		BorderStyle:  termout.BorderRounded,
		HasBorders:   true,
		HeaderColor:  termout.RoleSGRBold(termout.RolePrimary),
		BorderColor:  termout.SafeTokyoStyle(termout.TokyoMuted).Render(""),
		EllipsisText: "…",
		Columns: []termout.Column{
			{Title: "ID", MinWidth: 4, MaxWidth: 6, Align: termout.AlignRight},
			{Title: "SLUG / REPOSITORY", MinWidth: 24, MaxWidth: 36, Align: termout.AlignLeft},
			{Title: "BRANCH", MinWidth: 12, MaxWidth: 18, Align: termout.AlignLeft},
			{Title: "METHOD", MinWidth: 8, MaxWidth: 8, Align: termout.AlignCenter},
			{Title: "EXISTS", MinWidth: 8, MaxWidth: 8, Align: termout.AlignCenter},
			{Title: "REMOTE URL", MinWidth: 28, MaxWidth: 44, Align: termout.AlignLeft},
		},
		Rows: buildCloneTableRows(records, targetDir, isSSH),
	}

	termout.PrintTable(cfg)
}
```

#### 3.10.3 Row Formatting Helper
- Decomposed helper `buildCloneTableRows(records []model.ScanRecord, targetDir string, isSSH bool) []termout.Row` ($\le 12$ lines).
- Local status badge:
  - If repo exists: `termout.TokyoOKPill("YES")`.
  - If missing: `termout.SafeTokyoStyle(termout.TokyoMuted).Render("-")`.

---

## 4. Coding Guidelines & Architecture Invariants

### 4.1 Function Size Limit
- Every newly created or refactored function MUST NOT exceed 15 lines of executable body.
- Target function size: 6 to 10 lines.
- Deep or branchy operations must be broken down into discrete named helper functions.

### 4.2 Boolean Polarity & Naming
- All booleans must begin with `is` or `has` ONLY (`can`, `should`, `was`, `will` are forbidden).
- Negative boolean names (`isNotDirty`, `hasNoUpstream`) are strictly prohibited. Always use positive semantics (`isClean`, `hasUpstream`).

### 4.3 Control Flow & Nesting
- Zero nested `if` statements. All conditional logic must be flattened with guard clauses and early returns.
- Example:
  ```go
  // FORBIDDEN
  if isValid {
      if hasData {
          doWork()
      }
  }

  // REQUIRED
  if !isValid {
      return
  }
  if !hasData {
      return
  }
  doWork()
  ```

### 4.4 Spacing & Formatting Rules
- Blank line after closing curly brace `}` when followed by further statements.
- Blank line before every `return` statement, unless it is the only statement in a block.

### 4.5 Error Management
- Leaf functions returning `error` must return the error directly (`return apperror.Wrap(...)` or `return err`).
- Leaf functions must NEVER invoke exit handlers (`cliexit.HandleError`) or `os.Exit`.

---

## 5. Subtask Traceability Matrix

| Subtask ID | File Path | Focus Area | Verification Criteria |
|---|---|---|---|
| **Subtask 11** | `cli/termout/theme_lipgloss.go` | Tokyo Lipgloss helpers & progress bar | `TokyoBannerStyle`, `TokyoPill`, `RenderTokyoProgress` defined; $\le 15$ lines per func. |
| **Subtask 12** | `cli/termout/render.go` | Key-value formatting with `PadRightMax` | `PrintKeyValPadded` & `FormatKeyValPadded` implemented using role tokens. |
| **Subtask 13** | `cli/termout/termtable_types.go` | Border style enum & TableConfig extension | `BorderStyleType` enum with PascalCase constants; `TableConfig` updated. |
| **Subtask 14** | `cli/termout/table.go` | Frame border rendering (`BorderRounded`, etc.) | Full top, separator, and bottom border drawing with subtle Tokyo border colors. |
| **Subtask 15** | `cli/termout/truncate.go` | UTF-8 safe right truncation | `TruncateRightWithEllipsis` handles multi-byte runes and visual width accurately. |
| **Subtask 16** | `cli/cmdcg/cg_audit.go` | Dracula color purge & Tokyo binding | Dracula hex literals purged; bound to `termout.SafeTokyoStyle`. |
| **Subtask 17** | `cli/cmdcg/cg_doctor.go` | Progress bar & box frame integration | Doctor diagnostics integrate `RenderTokyoProgress` and `TokyoBoxStyle`. |
| **Subtask 18** | `cli/cmdhistory/history.go` | Purge local pad/trunc/color routines | Replaced local helpers with `termout` centralized functions. |
| **Subtask 19** | `cli/cmdpending/pending_commits_render.go` | Status pill & pad replacement | Replaced local `padRight` and status strings with `TokyoPill` and `PadRightMax`. |
| **Subtask 20** | `cli/cmdclone/clone_list_table.go` | Table migration to `termout.PrintTable` | Migrated `printCloneRecordsTable` to structured `termout.PrintTable` with `BorderRounded`. |
