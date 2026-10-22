# Subtask 16: Extract Status Table Renderer to status_render.go

## Status: Pending
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdstatus/status_render.go`, `cli/cmdstatus/helpers_shim.go`

---

## 1. Objective

Extract monolithic status table formatting out of `cli/cmdstatus/helpers_shim.go` into a dedicated companion module `cli/cmdstatus/status_render.go`. Migrate the status display engine to utilize GitMap's centralized `termout.PrintTable` with `BorderRounded`, column alignment, and `RowStripe: true`.

---

## 2. Scope & Implementation Details

1. **Create `cli/cmdstatus/status_render.go`**:
   - Belongs to `package cmdstatus`.
   - Function `renderStatusTableEngine(c *statusTableContext)`:
     - Constructs `termout.TableConfig`:
       - `BorderStyle: termout.BorderRounded`
       - `HeaderColor: constants.ColorYellow`
       - `RowStripe: true`
       - `Columns: buildStatusColumns(c)`
       - `Rows: buildStatusRows(c)`
     - Calls `termout.PrintTable(cfg)`.
     - Calls `printMissingRepoRemediation(c)` if missing repositories exist.

2. **Column & Row Construction**:
   - `buildStatusColumns(c *statusTableContext) []termout.ColumnConfig`:
     - Configures columns: `Repository`, `Branch`, `Status`, `Sync`, `Stash`, `Files`.
     - Sets alignments: `AlignLeft` for repo and branch, `AlignCenter` for status icon, `AlignRight` for stash count.
   - `buildStatusRows(c *statusTableContext) []termout.Row`:
     - Iterates through `c.Rows`, converting each `statusRow` into a `termout.Row` with cells: `[r.RepoName, r.Branch, r.StateIcon, r.SyncText, r.StashText, r.FilesText]`.
     - Colorizes repository name with rotating pastel colors.

3. **Refactor `cli/cmdstatus/helpers_shim.go`**:
   - Update `printStatusTableWithContext(c *statusTableContext)` to forward to `renderStatusTableEngine(c)`.
   - Remove obsolete manual padding and dash divider routines (`printStatusTableHeader`, `printStatusTableRow`, `calculateStatusDividerLen`).

---

## 3. Strict Coding Guidelines & Invariants

- **Component Sizing**: Break construction logic into small helpers $\le 15$ lines.
- **Leaf Error Sovereignty**: Preserve error propagation without introducing void exit handlers.
- **Zero Nested Ifs**: Invert conditions to eliminate nested conditional branching.
- **Whitespace Hygiene**: Blank line after closing brace `}` and blank line before every `return`.
- **Zero Builds & Tests**: Strict prohibition on running `go test` or `go build`.

---

## 4. Verification Checklist

- [ ] Dedicated companion module `cli/cmdstatus/status_render.go` created.
- [ ] Status table rendered via `termout.PrintTable` with `BorderRounded` and `RowStripe: true`.
- [ ] Manual formatting code in `helpers_shim.go` cleanly refactored.
- [ ] Missing repository remediation guidance preserved and displayed.
- [ ] All functions $\le 15$ lines of executable code.
