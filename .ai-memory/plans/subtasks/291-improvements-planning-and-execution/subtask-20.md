# Subtask 20: Enhance Coding Guidelines Doctor with Section Dividers and TokyoMetricPill

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdcg/cg_doctor.go`

---

## 1. Objective

Enhance `cli/cmdcg/cg_doctor.go` by replacing raw ASCII divider titles with `termout.PrintSectionDivider` and upgrading the doctor summary card with `termout.TokyoMetricPill` badge counters.

---

## 2. Scope & Implementation Details

1. **Modernize Command Header in `RunCGDoctor`**:
   - Current: `fmt.Printf("%s\n\n", cgAuditTitleStyle.Render("=== Coding Guidelines Doctor Diagnostics ==="))`
   - Upgrade: Call `termout.PrintSectionDivider("Coding Guidelines Doctor Diagnostics", termout.DefaultTheme())`.

2. **Enhance Summary Content in `buildDoctorSummaryContent`**:
   - Signature: `buildDoctorSummaryContent(totalRepos, healthyRepos int) string`
   - Calculate metrics:
     - `issuesCount := totalRepos - healthyRepos`
     - `bar := termout.RenderTokyoProgress(healthyRepos, totalRepos, 24)`
   - Assemble metrics line using `TokyoMetricPill`:
     - `totalPill := termout.TokyoMetricPill("Total", totalRepos)`
     - `healthyPill := termout.TokyoMetricPill("Healthy", healthyRepos)`
     - `issuesPill := termout.TokyoMetricPill("Issues", issuesCount)`
     - Combine pills: `fmt.Sprintf("%s  %s  %s", totalPill, healthyPill, issuesPill)`
   - Format body text indicating optimal state or action required.
   - Return combined summary string.

3. **Framing in `renderDoctorOverall`**:
   - Passes assembled summary content to `termout.TokyoBoxStyle(content)`.
   - Prints styled box to stdout.

4. **Decomposition**:
   - Keep `buildDoctorSummaryContent`, `renderDoctorOverall`, and `RunCGDoctor` $\le 15$ lines each.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Invert conditions and use early returns.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `RunCGDoctor` uses `termout.PrintSectionDivider`.
- [x] `buildDoctorSummaryContent` integrates `termout.TokyoMetricPill` for Total, Healthy, and Issues metrics.
- [x] Retains `termout.RenderTokyoProgress` and `termout.TokyoBoxStyle`.
- [x] All functions $\le 15$ lines without nested `if` statements.
