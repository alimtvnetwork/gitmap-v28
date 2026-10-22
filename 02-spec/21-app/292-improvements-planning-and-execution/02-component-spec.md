# Component Specification: Terminal UI Modernization, Table Framing & Output Ergonomics

**Module:** `cli/termout`, `cli/cmdstatus`, `cli/cmdhaschange`, `cli/cmdpending`, `cli/cmd`  
**Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
**Specification File:** `02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md`  
**Subtask Range:** Subtasks 11 through 20 (Improvements 11 to 20)  
**Status:** Approved for Implementation  

---

## 1. Executive Summary & Component Scope

This component specification establishes the technical requirements, Go data models, styling tokens, formatting contracts, and refactoring guidelines for Improvements 11 through 20 of the `292-improvements-planning-and-execution` milestone.

These improvements elevate GitMap's terminal user interface (TUI) and output ergonomics by standardizing Lipgloss visual primitives, adding customizable badge and header builders, delivering visual-width key-value alignment, providing robust table cell wrap modes and subtle Tokyo surface row striping, implementing rune-safe left truncation for file path displays, extracting monolithic status table rendering into structured components, formatting multi-dimensional change states with TokyoTag, modernizing pending commits visual tree structures, elevating `gitmap see` command menus, and strictly enforcing the standard error (`os.Stderr`) isolation invariant for interactive selection prompts.

### 1.1 Scope Boundaries & Subsystem Mapping

1. **`cli/termout/theme_lipgloss.go` (Improvement 11)**:
   - Implements `TokyoBadge(label, bgHex, fgHex string) string` supporting arbitrary background/foreground Tokyo Night color pairing with safe plain fallback.
   - Implements `TokyoHeader(title string) string` providing uppercase bold styled section title banners with TokyoPrimary accenting and padding.
   - Enforces the $\le 15$-line function limit and zero nested conditionals.

2. **`cli/termout/render.go` (Improvement 12)**:
   - Implements `FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string` and `PrintKeyValuePair(key, val string, keyWidth int, theme Theme)`.
   - Utilizes `VisualStringWidth` (stripping ANSI codes and accounting for East Asian and emoji double-width runes) to guarantee precision column alignment across Unicode keys.
   - Adheres to theme color semantics with automatic fallback for monochrome / non-TTY outputs.

3. **`cli/termout/termtable_types.go` (Improvement 13)**:
   - Defines strongly-typed enum `WrapModeType` (`WrapClip`, `WrapEllipsis`, `WrapWord`).
   - Expands `ColumnConfig` with `WrapMode WrapModeType`.
   - Expands `TableConfig` with `WrapMode WrapModeType` and positive boolean field `RowStripe bool`.

4. **`cli/termout/table.go` (Improvement 14)**:
   - Implements multi-mode cell wrapping logic (`WrapClip`, `WrapEllipsis`, `WrapWord`) in cell formatting pipelines.
   - Implements alternating row striping (`RowStripe`) applying subtle background styling (`TokyoSurface` or subtle tint) on alternating data rows when colors are enabled.
   - Decomposes row rendering into modular helper functions strictly $\le 15$ lines.

5. **`cli/termout/truncate.go` (Improvement 15)**:
   - Implements `TruncateLeftWithEllipsis(text string, maxWidth int) string` using rune slicing and visual width evaluation.
   - Slices trailing runes to preserve filename endings (e.g. `.../subpackage/main.go`), inserting a leading `"..."` ellipsis without Unicode character tearing.
   - Handles edge cases where `maxWidth <= 3` or `maxWidth <= 0` deterministically.

6. **`cli/cmdstatus/status_render.go` (Improvement 16)**:
   - Extracts manual, unaligned terminal table rendering out of `cli/cmdstatus/helpers_shim.go` into dedicated companion module `cli/cmdstatus/status_render.go`.
   - Wires status reporting through `termout.PrintTable` utilizing `BorderRounded`, column alignment, and `RowStripe: true`.
   - Decomposes status table construction and remediation display into discrete sub-15-line functions.

7. **`cli/cmdhaschange/haschange.go` (Improvement 17)**:
   - Modernizes multi-dimensional `--all` inspection (`dirty`, `ahead`, `behind`) by formatting outputs with `termout.TokyoTag` and semantic pill colors.
   - Strictly preserves single-dimension stdout output purity (`true` or `false` without ANSI codes or decoration) for seamless shell script embedding.
   - Provides safe fallback for non-TTY / pipe output modes.

8. **`cli/cmdpending/pending_commits_render.go` (Improvement 18)**:
   - Modernizes pending commits table framing by migrating from ASCII/sharp boxes to `termout.BorderRounded` and `termout.TokyoBoxStyle`.
   - Enriches hierarchical tree remediation branches (`├──`, `└──`) using Tokyo Night subtle color tokens (`TokyoMuted`, `TokyoPrimary`).
   - Ensures visual alignment across variable-width repository names and branch labels.

9. **`cli/cmd/cmd_see.go` (Improvement 19)**:
   - Refactors `printPendingCommits` to use `termout.TokyoHeader` and `termout.TokyoTag` / `termout.TokyoMetricPill` for structured summaries.
   - Replaces raw unstructured multiline text in `PrintSeeHelp()` with standard `termout.RenderMenu` backed by `HelpMenu`, `HelpSection`, and `CommandEntry`.
   - Fixes blank line formatting preceding return statements.

10. **`cli/cmdpending/pending_commits_filter.go` (Improvement 20)**:
    - Enforces the candidate selection standard error isolation invariant by directing interactive prompts (`Press [Enter] for next, or 'q' to quit`) strictly to `os.Stderr`.
    - Guarantees standard output (`os.Stdout`) remains clean and unpolluted during interactive 1-by-1 stepping or candidate filtering.

---

## 2. Shared Data Models & Go Contracts

### 2.1 Table Wrap Modes & Row Striping Configuration (`cli/termout/termtable_types.go`)

```go
package termout

// WrapModeType defines cell overflow wrapping and truncation behavior.
type WrapModeType string

const (
	// WrapClip clips cell contents exactly at column width without appending ellipsis.
	WrapClip WrapModeType = "clip"
	// WrapEllipsis truncates overflowing cell contents by appending or inserting an ellipsis (...).
	WrapEllipsis WrapModeType = "ellipsis"
	// WrapWord breaks overflowing text into multiple wrapped lines at whitespace boundaries.
	WrapWord WrapModeType = "word"
)

// ColumnConfig defines configuration options for an individual table column.
type ColumnConfig struct {
	Title     string
	MaxWidth  int
	MinWidth  int
	Alignment AlignmentType
	Align     AlignmentType // Backward-compatible field alias
	WrapMode  WrapModeType
}

// TableConfig defines full configuration for table rendering.
type TableConfig struct {
	Columns      []ColumnConfig
	Rows         []Row
	HeaderColor  string
	BorderColor  string
	HasBorders   bool
	BorderStyle  BorderStyleType
	EllipsisText string
	WrapMode     WrapModeType
	RowStripe    bool
}
```

### 2.2 Tokyo Badge and Header Signatures (`cli/termout/theme_lipgloss.go`)

```go
package termout

// TokyoBadge renders text inside a rounded badge with customizable background and foreground hex colors.
// In non-TTY or no-color environments, returns plain unstyled "[label]".
func TokyoBadge(label, bgHex, fgHex string) string

// TokyoHeader renders a prominent section header formatted in bold with TokyoPrimary accenting.
// In non-TTY or no-color environments, returns plain uppercase title.
func TokyoHeader(title string) string
```

### 2.3 Visual-Width Key-Value Pair Contracts (`cli/termout/render.go`)

```go
package termout

// FormatKeyValuePair formats a key-value row with visual width alignment based on keyWidth.
// Color formatting adapts to the provided Theme and colorAllowed() state.
func FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string

// PrintKeyValuePair outputs FormatKeyValuePair directly to output.UI().
func PrintKeyValuePair(key, val string, keyWidth int, theme Theme)
```

### 2.4 Rune-Safe Left Truncation Contract (`cli/termout/truncate.go`)

```go
package termout

// TruncateLeftWithEllipsis truncates text exceeding maxWidth by placing a leading "..." ellipsis,
// preserving the rightmost characters (e.g. filename suffixes or path leaves) using visual width calculations.
func TruncateLeftWithEllipsis(text string, maxWidth int) string
```

---

## 3. Subsystem Component Breakdown (Improvements 11–20)

| Subtask | Target Subsystem / Files | Core Technical Refactoring | Architectural Invariant |
| :--- | :--- | :--- | :--- |
| **11** | `cli/termout/theme_lipgloss.go` | Implement `TokyoBadge(label, bgHex, fgHex string)` and `TokyoHeader(title string)`. Safe plain fallback when `!colorAllowed()`. | Lipgloss Token Standard, $\le 15$ Lines |
| **12** | `cli/termout/render.go` | Implement `FormatKeyValuePair` and `PrintKeyValuePair` with visual column alignment via `VisualStringWidth`. | Visual Column Precision, Theme Decoupling |
| **13** | `cli/termout/termtable_types.go` | Add `WrapModeType` enum (`WrapClip`, `WrapEllipsis`, `WrapWord`). Add `WrapMode` and `RowStripe` to `TableConfig`. | Strongly-Typed Enums, Positive Booleans |
| **14** | `cli/termout/table.go` | Implement cell wrap formatting per `WrapModeType` and alternating row striping with subtle Tokyo surface background tint. | Visual Layout Hygiene, $\le 15$ Lines |
| **15** | `cli/termout/truncate.go` | Implement `TruncateLeftWithEllipsis(text string, maxWidth int)` using rune slicing and visual width evaluation. | Rune Display Safety, Unicode Preservation |
| **16** | `cli/cmdstatus/status_render.go`<br/>`cli/cmdstatus/helpers_shim.go` | Extract status table renderer into `status_render.go` using `termout.PrintTable`, `BorderRounded`, and `RowStripe: true`. | Component Separation, Leaf Error Sovereignty |
| **17** | `cli/cmdhaschange/haschange.go` | Modernize `--all` with `termout.TokyoTag` formatting while keeping single-mode stdout pure (`true`/`false`) for shell scripting. | Scripting Purity Invariant, TTY Adaptability |
| **18** | `cli/cmdpending/pending_commits_render.go` | Replace sharp box characters with `termout.BorderRounded` and `TokyoBoxStyle`. Modernize tree branches with Tokyo tokens. | Visual Cohesion, Standard Border Styling |
| **19** | `cli/cmd/cmd_see.go` | Modernize `printPendingCommits` with Tokyo tokens. Convert raw text in `PrintSeeHelp` to structured `termout.RenderMenu`. | Help Menu Standard, Whitespace Rules |
| **20** | `cli/cmdpending/pending_commits_filter.go` | Isolate interactive candidate selection prompts (`Press [Enter]...`) strictly to `os.Stderr`. | Stderr Prompt Isolation Invariant |

---

## 4. Detailed Component Technical Specifications

### 4.1 Component 11: TokyoBadge & TokyoHeader Lipgloss Builders (`cli/termout/theme_lipgloss.go`)

#### 4.1.1 Problem Statement
While `theme_lipgloss.go` provides standard pills (`TokyoOKPill`, `TokyoWarnPill`, `TokyoErrorPill`), callers require flexible, parameterized badge styling with custom Tokyo Night color pairings (e.g., custom repository roles, branch types, priority levels) without re-implementing Lipgloss boilerplate. Additionally, section headers across subcommands lack a uniform bold header builder.

#### 4.1.2 Technical Design
1. **`TokyoBadge(label, bgHex, fgHex string) string`**:
   - Guard clause: If `!colorAllowed()`, return `fmt.Sprintf("[%s]", label)`.
   - If `label == ""`, return `""`.
   - Style construction:
     ```go
     style := lipgloss.NewStyle().
         Bold(true).
         Foreground(lipgloss.Color(fgHex)).
         Background(lipgloss.Color(bgHex)).
         Padding(0, 1)
     return style.Render(label)
     ```
2. **`TokyoHeader(title string) string`**:
   - Guard clause: If `!colorAllowed()`, return `strings.ToUpper(title)`.
   - Style construction:
     ```go
     style := lipgloss.NewStyle().
         Bold(true).
         Foreground(lipgloss.Color(TokyoPrimary)).
         MarginBottom(1)
     return style.Render(strings.ToUpper(title))
     ```
3. **Coding Standards**:
   - Each function must be $\le 15$ lines.
   - Blank line after closing brace and before return.

---

### 4.2 Component 12: Visual-Aligned Key-Value Pair Output (`cli/termout/render.go`)

#### 4.2.1 Problem Statement
Existing key-value output routines rely on byte-length padding (`PadRightMax`), which misaligns columns when keys contain multi-byte Unicode characters, accented glyphs, or decorative icons.

#### 4.2.2 Technical Design
1. **`FormatKeyValuePair(key, val string, keyWidth int, theme Theme) string`**:
   - Formats the key label with a trailing colon: `keyText := key + ":"`.
   - Computes printable visual width via `VisualStringWidth(keyText)`.
   - Calculates padding: `padLen := keyWidth - visualWidth`. If `padLen < 1`, default to `1`.
   - Builds padded label string: `paddedKey := keyText + SpacePool(padLen)`.
   - When `!colorAllowed()`, returns `fmt.Sprintf("  %s %s", paddedKey, val)`.
   - When color is enabled, applies `resolveKeyColor(theme)` to `paddedKey` and appends `val`.
2. **`PrintKeyValuePair(key, val string, keyWidth int, theme Theme)`**:
   - Invokes `FormatKeyValuePair` and outputs to `output.UI()`.
3. **Decomposition**:
   - Decomposes into `resolveKeyColor(theme Theme) string` and `buildPaddedKey(key string, keyWidth int) string` to ensure functions remain under 12 lines.

---

### 4.3 Component 13: Table Wrapping Modes & Row Striping Types (`cli/termout/termtable_types.go`)

#### 4.3.1 Problem Statement
The table rendering engine lacks explicit cell wrap configuration, forcing callers to either manually pre-truncate strings or accept default middle-truncation. Furthermore, long tables lack zebra striping to assist visual scanning across wide terminals.

#### 4.3.2 Technical Design
1. **`WrapModeType` Enum**:
   - Strongly-typed string enum ending in `Type`:
     - `WrapClip`: Strict cutoff at column width.
     - `WrapEllipsis`: Truncate with `"..."` ellipsis (default).
     - `WrapWord`: Soft wrap onto subsequent lines (line wrapping within cell).
2. **Type Extensions**:
   - Add `WrapMode WrapModeType` to `ColumnConfig` (column-level override).
   - Add `WrapMode WrapModeType` to `TableConfig` (table-level default).
   - Add `RowStripe bool` to `TableConfig` (positive boolean naming convention).

---

### 4.4 Component 14: Table Cell Wrap Execution & Subtle Surface Striping (`cli/termout/table.go`)

#### 4.4.1 Problem Statement
`cli/termout/table.go` currently executes `TruncateMiddle` uniformly across all cells and lacks row background styling.

#### 4.4.2 Technical Design
1. **Cell Wrap Strategy Resolution**:
   - Function `formatCellWithWrap(text string, width int, align AlignmentType, wrap WrapModeType, ellipsis string) string`:
     - If `wrap == WrapClip`, slice string visually to `width` without ellipsis.
     - If `wrap == WrapEllipsis` or `wrap == ""`, apply `TruncateMiddle(text, width, ellipsis)`.
     - Align truncated text via existing `alignPaddedCell`.
2. **Row Striping Implementation**:
   - In `appendFramedRows` and `appendRowBlocks`, maintain row index `rowIdx`.
   - If `cfg.RowStripe && colorAllowed() && rowIdx%2 == 1`:
     - Apply subtle background tinting using `TokyoSurface` (`#1f2335` or ANSI `\033[48;2;31;35;53m`) to the formatted row content.
     - Reset ANSI background at row line boundary.
3. **Decomposition**:
   - Modularize row coloring into `resolveRowBackground(rowIdx int, cfg TableConfig) string`.

---

### 4.5 Component 15: Rune-Safe Left Truncation (`cli/termout/truncate.go`)

#### 4.5.1 Problem Statement
When displaying nested filesystem paths or long Git references in fixed-width columns, truncating from the middle or right cuts off the most distinguishing information (the filename or child directory). A left-truncating function (`.../component_spec.md`) is required that preserves Unicode runes without byte-slicing panics.

#### 4.5.2 Technical Design
1. **`TruncateLeftWithEllipsis(text string, maxWidth int) string`**:
   - Guard clauses:
     - `if maxWidth <= 0`: return `""`.
     - `if VisualStringWidth(text) <= maxWidth`: return `text`.
     - `if maxWidth <= 3`: return `"..."[:maxWidth]`.
   - Execution:
     - Target budget for trailing content: `budget := maxWidth - 3`.
     - Slices suffix runes using existing helper `sliceSuffixRunes([]rune(text), budget)`.
     - Returns `"..." + suffix`.
2. **Guidelines**:
   - Function must be under 12 lines.
   - Blank line before return.

---

### 4.6 Component 16: Status Table Extraction & Rounded Framing (`cli/cmdstatus/status_render.go`)

#### 4.6.1 Problem Statement
`cli/cmdstatus/helpers_shim.go` is over 440 lines long and contains manual, imperative formatting loops (`printStatusTableHeader`, `printStatusTableRow`, `printStatusTableWithContext`) that construct tables with ad-hoc padding and dashes instead of using the repository's centralized `termout.PrintTable` engine.

#### 4.6.2 Technical Design
1. **Module Creation**:
   - Create `cli/cmdstatus/status_render.go` in package `cmdstatus`.
2. **Table Construction**:
   - Define columns:
     - `Repository` (MaxWidth 28, AlignLeft)
     - `Branch` (MaxWidth 20, AlignLeft)
     - `Status` (MaxWidth 12, AlignCenter)
     - `Sync` (MaxWidth 14, AlignLeft)
     - `Stash` (MaxWidth 10, AlignRight)
     - `Files` (MaxWidth 16, AlignLeft)
   - Configure `TableConfig`:
     - `BorderStyle: termout.BorderRounded`
     - `RowStripe: true`
     - `HeaderColor: constants.ColorYellow`
3. **Shim Refactoring**:
   - Delegate `printStatusTableWithContext` to `renderStatusTableEngine(c *statusTableContext)` in `status_render.go`.
   - Keep `helpers_shim.go` clean and focused on data accumulation.

---

### 4.7 Component 17: HasChange TokyoTag Formatting for Multi-Dimension Output (`cli/cmdhaschange/haschange.go`)

#### 4.7.1 Problem Statement
`gitmap has-change <repo> --all` outputs raw, unstyled text `dirty=true ahead=false behind=true`. In terminal sessions, this is visually bland and harder to read at a glance. However, scripts frequently rely on single-mode `gitmap has-change <repo>` returning raw `"true"` or `"false"`.

#### 4.7.2 Technical Design
1. **Single-Mode Invariant**:
   - `printHasChangeOne(target, mode string)` MUST strictly retain its pure stdout output (`fmt.Println("true")` / `fmt.Println("false")`). Zero ANSI codes or styling tags are permitted in single mode.
2. **Multi-Dimension Modernization (`--all`)**:
   - In `printHasChangeAll(target string)`:
     - Check `termout.ColorAllowed()`.
     - If color is allowed, format the three dimensions using `termout.TokyoTag`:
       - `termout.TokyoTag("dirty", boolStr(dirty))`
       - `termout.TokyoTag("ahead", boolStr(ahead > 0))`
       - `termout.TokyoTag("behind", boolStr(behind > 0))`
       - Print the space-separated tags on a single line.
     - If color is disallowed (or piped to a file/script), retain the standard output format: `dirty=true ahead=false behind=true`.
3. **Guidelines**:
   - Functions must remain $\le 15$ lines with zero nested conditionals.

---

### 4.8 Component 18: Pending Commits Rounded Framing & Tree Modernization (`cli/cmdpending/pending_commits_render.go`)

#### 4.8.1 Problem Statement
`cli/cmdpending/pending_commits_render.go` constructs tables and summaries using hardcoded ASCII box-drawing characters (`┌───`, `│`, `└───`), which clash with modern rounded Tokyo Night framing across the rest of the CLI.

#### 4.8.2 Technical Design
1. **Banner & Header Modernization**:
   - Replace manual `┌───...───┐` in `renderTableBanner` with `termout.BorderRounded` or `termout.TokyoBoxStyle`.
   - Use `termout.TokyoHeader("Gitmap Pending Commits Summary")`.
2. **Table Summary Pill Formatting**:
   - Format summary metrics in `renderTableSummary` using `termout.TokyoMetricPill` for `Scanned`, `Dirty`, `Uncommitted`, and `Unpushed`.
3. **Tree Remediation Hints**:
   - In `renderRemediationTitleBranch` and `renderRemediationChildBranch`, style tree branch glyphs (`├──`, `└──`, `│`) using `termout.SafeTokyoStyle(termout.TokyoMuted)`.
   - Style option commands using `termout.SafeTokyoStyle(termout.TokyoAccent)`.
4. **Footer Modernization**:
   - Replace manual footer borders with `termout.BorderRounded` bottom glyphs.

---

### 4.9 Component 19: See Suite Pending Commits & Modernized Help Menu (`cli/cmd/cmd_see.go`)

#### 4.9.1 Problem Statement
1. `printPendingCommits` in `cli/cmd/cmd_see.go` prints raw text with `=== PENDING COMMITS ===`.
2. `PrintSeeHelp` prints a raw multiline string with `fmt.Printf`, bypassing the repository's unified `termout.RenderMenu` architecture.
3. Code style violations: blank lines missing preceding return statements.

#### 4.9.2 Technical Design
1. **Modernize `printPendingCommits`**:
   - Render header with `termout.TokyoHeader("Pending Commits")`.
   - Format repository count with `termout.TokyoMetricPill("Repositories", len(dirty))`.
   - Format each dirty repo line with `termout.TokyoTag` for modified, untracked, and staged counts.
2. **Standardize `PrintSeeHelp`**:
   - Construct a `termout.HelpMenu`:
     - `Title: "GitMap Inspection Suite (see / c)"`
     - `UsageLines: []string{"gitmap see <subcommand> [flags]"}`
     - `Sections: []termout.HelpSection{...}` detailing `commit pending`, `git-ignore issues`, `errors`, `history`, `errors ssh`, `nodes history`, and `repo-manage ui`.
   - Invoke `termout.RenderMenu(menu)`.
3. **Style Compliance**:
   - Enforce $\le 15$-line functions and mandatory blank line before every return.

---

### 4.10 Component 20: Candidate Selection Stderr Isolation Invariant (`cli/cmdpending/pending_commits_filter.go`)

#### 4.10.1 Problem Statement
In interactive modes (such as 1-by-1 stepping through pending commits), prompt queries (`Press [Enter] for next, or 'q' to quit: `) are output to standard output (`os.Stdout`). If a user or automated agent pipes stdout to a file or JSON parser, these interactive prompt strings corrupt the data stream.

#### 4.10.2 Technical Design
1. **Stderr Isolation Invariant**:
   - All interactive selection prompts, user confirmation queries, and pagination controls MUST be written strictly to `os.Stderr`.
2. **Refactoring `isDetailPromptTerminated`**:
   - Replace `fmt.Printf("  Press [Enter] for next, or 'q' to quit: ")` with:
     ```go
     fmt.Fprintf(os.Stderr, "  Press [Enter] for next, or 'q' to quit: ")
     ```
3. **Refactoring `renderPendingCommitsInteractive`**:
   - Write stepping pagination header to `os.Stderr` while preserving stdout for structured payload streaming when redirected.

---

## 5. Architectural Invariants & Non-Negotiables

1. **Strictly Relative Git Paths Only**:
   - All references to files across specs, code comments, subtasks, and documentation MUST use relative forward-slash paths (e.g. `cli/termout/theme_lipgloss.go`). No absolute paths or URI schemes.
2. **Zero Nested Ifs**:
   - All branching logic must use inverted guard clauses with immediate early returns.
3. **15-Line Limit**:
   - Every function must be $\le 15$ lines of executable code ($\le 8$ lines preferred).
4. **Positive Booleans**:
   - All boolean variables and parameters must begin with `is` or `has` (e.g., `isColorSupported`, `hasBorder`, `hasChange`). Negative booleans (`isNotDirty`, `hasNoWrap`) are strictly banned.
5. **Zero Builds and Tests**:
   - Spec and plan authoring strictly prohibits executing `go test` or `go build`.
