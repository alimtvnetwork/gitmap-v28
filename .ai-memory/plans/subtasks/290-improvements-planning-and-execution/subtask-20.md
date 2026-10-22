# Subtask 20: Modernize Clone List Table with Rounded Borders & Theme Roles

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdclone/clone_list_table.go`

---

## 1. Objective

Refactor `printCloneRecordsTable` in `cli/cmdclone/clone_list_table.go` to replace manual `fmt.Printf` / ASCII dash line formatting with declarative `termout.PrintTable`, using `BorderRounded` frame styling, Tokyo Night theme tokens, and structured status badges.

---

## 2. Scope & Implementation Details

1. **Table Configuration Construction**:
   - Refactor `printCloneRecordsTable(records []model.ScanRecord, targetDir string, isSSH bool)`:
     - Configure `termout.TableConfig`:
       - `BorderStyle`: `termout.BorderRounded`
       - `HasBorders`: `true`
       - `HeaderColor`: `termout.RoleSGRBold(termout.RolePrimary)`
       - `BorderColor`: `termout.SafeTokyoStyle(termout.TokyoMuted).Render("")`
       - `EllipsisText`: `"…"`
       - Columns:
         - `ID`: `MinWidth: 4`, `MaxWidth: 6`, `Align: termout.AlignRight`
         - `SLUG / REPOSITORY`: `MinWidth: 24`, `MaxWidth: 36`, `Align: termout.AlignLeft`
         - `BRANCH`: `MinWidth: 12`, `MaxWidth: 18`, `Align: termout.AlignLeft`
         - `METHOD`: `MinWidth: 8`, `MaxWidth: 8`, `Align: termout.AlignCenter`
         - `EXISTS`: `MinWidth: 8`, `MaxWidth: 8`, `Align: termout.AlignCenter`
         - `REMOTE URL`: `MinWidth: 28`, `MaxWidth: 44`, `Align: termout.AlignLeft`
       - Rows: `buildCloneTableRows(records, targetDir, isSSH)`
     - Invoke `termout.PrintTable(cfg)`.

2. **Row Assembly Helper**:
   - Decompose row generation into helper `buildCloneTableRows(records []model.ScanRecord, targetDir string, isSSH bool) []termout.Row` $\le 12$ lines:
     - For each record `r`:
       - Resolve slug, branch (default `"-"` if empty), method, exists, remote URL.
       - If `exists == "YES"`, format cell with `termout.TokyoOKPill("YES")`.
       - If `exists == "-"`, format cell with `termout.SafeTokyoStyle(termout.TokyoMuted).Render("-")`.
       - Append `termout.Row{Cells: []string{...}}`.

3. **Leaf Error Sovereignty**:
   - In `runCloneListTable(sourcePath, targetDir string, isSSH bool, onlyFilter, excludeFilter string) error`:
     - Keep returning errors wrapped in `apperror.WrapSimple`.
     - Zero calls to `cliexit.HandleError` or `os.Exit`.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Max 15 lines per function.
- **Zero Nested Ifs**: Guard clauses and early returns only.
- **Affirmative Booleans**: `isSSH`, `isGitRepo`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `printCloneRecordsTable` invokes `termout.PrintTable` instead of raw `fmt.Printf`.
- [x] Table configuration uses `termout.BorderRounded` and Tokyo Night theme tokens.
- [x] Row formatting decomposed into helper functions $\le 15$ lines.
- [x] Status indicator for existing local repos uses `termout.TokyoOKPill("YES")`.
- [x] Leaf error sovereignty maintained in `runCloneListTable`.
