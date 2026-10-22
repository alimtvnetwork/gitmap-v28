# Subtask 18: Create Dedicated CD Selection Table Renderer

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdcd/cd_table.go`

---

## 1. Objective

Create `cli/cmdcd/cd_table.go` to provide a dedicated interactive candidate selection table renderer for `gitmap cd` that prints strictly to `os.Stderr`, guaranteeing that shell wrappers like `cd $(gitmap cd ...)` receive pure target paths on standard output without UI corruption.

---

## 2. Scope & Implementation Details

1. **New File**: `cli/cmdcd/cd_table.go` in package `cmdcd`.

2. **Core Function**:
   - `RenderCDSelectionTable(title string, records []model.ScanRecord)`

3. **Table Configuration (`termout.TableConfig`)**:
   - Columns:
     - `#`: `MinWidth: 4`, `MaxWidth: 4`, `Alignment: termout.AlignRight`
     - `Repository`: `MinWidth: 16`, `MaxWidth: 24`, `Alignment: termout.AlignLeft`
     - `Branch`: `MinWidth: 10`, `MaxWidth: 16`, `Alignment: termout.AlignLeft`
     - `Path`: `MinWidth: 28`, `MaxWidth: 56`, `Alignment: termout.AlignLeft`
   - Border Style: `BorderStyle: termout.BorderRounded`.
   - Header Color: `termout.RoleSGR(termout.RolePrimary)`.
   - HasBorders: `true`.

4. **Row Assembly (`formatCDRecordRow(idx int, r model.ScanRecord) termout.Row`)**:
   - Formats `#` with index `fmt.Sprintf("%d", idx)`.
   - Formats `Repository` with `r.RepoName`.
   - Formats `Branch` with `r.Branch` (or `"-"` if empty).
   - Formats `Path` with middle-truncated absolute path: `termout.TruncateMiddleWithEllipsis(r.AbsolutePath, 56)`.

5. **CRITICAL Standard Error Invariant**:
   - `RenderCDSelectionTable` MUST output exclusively to `os.Stderr`:
     ```go
     rendered := termout.RenderTable(cfg)
     fmt.Fprintf(os.Stderr, "\n  %s\n%s\n", title, rendered)
     ```
   - Standard output (`os.Stdout`) must NEVER receive any table or prompt characters.

6. **Decomposition**:
   - `buildCDTableConfig(records []model.ScanRecord) termout.TableConfig`
   - `formatCDRecordRow(idx int, r model.ScanRecord) termout.Row`
   - `RenderCDSelectionTable(title string, records []model.ScanRecord)`
   - All functions $\le 15$ lines.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Guard clauses for empty record slice.
- **Error Sovereignty / Stream Isolation**: Strictly write to `os.Stderr`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `cli/cmdcd/cd_table.go` created in package `cmdcd`.
- [x] `RenderCDSelectionTable(title string, records []model.ScanRecord)` exported.
- [x] Table configured with rounded borders and column alignments.
- [x] Path column truncated using `termout.TruncateMiddleWithEllipsis`.
- [x] Rendered output directed strictly to `os.Stderr`.
- [x] All functions $\le 15$ lines.
