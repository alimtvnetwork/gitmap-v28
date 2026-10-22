# Subtask 19: Modernize See Command Pending Commits and Help Menu

## Status: Completed
## Parent Plan: .ai-memory/plans/292-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/292-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmd/cmd_see.go`

---

## 1. Objective

Modernize `printPendingCommits` in `cli/cmd/cmd_see.go` using Tokyo Night design tokens and pills, and refactor `PrintSeeHelp` from a raw unformatted multiline string into GitMap's standardized `termout.RenderMenu` architecture backed by `HelpMenu` and `HelpSection`. Enforce whitespace and coding guidelines.

---

## 2. Scope & Implementation Details

1. **Modernize `printPendingCommits`**:
   - Render section header using `termout.TokyoHeader("Pending Commits")` or `termout.PrintSectionDivider`.
   - Render repository count pill: `termout.TokyoMetricPill("Repositories", len(dirty))`.
   - For each dirty repository:
     - Format repo name in bold accent.
     - Format metrics using `termout.TokyoTag("modified", fmt.Sprint(d.Diagnosis.ModifiedCount))`, `termout.TokyoTag("untracked", fmt.Sprint(d.Diagnosis.UntrackedCount))`, and `termout.TokyoTag("staged", fmt.Sprint(d.Diagnosis.StagedCount))`.
   - Break loop rendering into a sub-15-line helper function.

2. **Standardize `PrintSeeHelp()`**:
   - Construct a `termout.HelpMenu` structure:
     - `Title: "GitMap Inspection Suite (see / c)"`
     - `UsageLines: []string{"gitmap see <subcommand> [flags]"}`
     - `Sections: []termout.HelpSection{...}` with entries:
       - `commit pending (cp)`: "List repositories with uncommitted or unpushed changes"
       - `git-ignore issues (ig)`: "List repositories with ignore or duplicate rules"
       - `errors`: "Display recent execution error logs"
       - `history`: "Display command execution history"
       - `errors ssh (ses)`: "Display execution errors across remote SSH fleet"
       - `nodes history`: "Display command history across fleet nodes"
       - `repo-manage ui`: "Launch interactive terminal dashboard"
   - Delegate rendering to `termout.RenderMenu(menu)`.

3. **Guideline & Whitespace Fixes**:
   - Ensure blank lines precede all return statements.
   - Enforce the $\le 15$-line function limit across all helpers.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Functions must be $\le 15$ lines of executable code.
- **Zero Raw Strings for Help Menus**: Help output must route through `termout.RenderMenu`.
- **Positive Booleans**: Use `isClean`, `hasHook`.
- **Return Styling**: Blank line after closing brace `}` and blank line before every `return`.
- **Zero Builds & Tests**: Strict prohibition on compiler builds or tests.

---

## 4. Verification Checklist

- [x] `printPendingCommits` formatted with Tokyo Night header and metric tags.
- [x] `PrintSeeHelp` renders standard `termout.HelpMenu` via `termout.RenderMenu`.
- [x] Raw multiline `fmt.Printf` block in `PrintSeeHelp` eliminated.
- [x] Whitespace rules satisfied (blank lines before `return`).
- [x] All functions conform to $\le 15$ lines of code.
