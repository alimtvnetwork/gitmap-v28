# Subtask 17: Integrate Tokyo Progress & Box Framing into CG Doctor

## Status: completed
## Parent Plan: .ai-memory/plans/290-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/290-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdcg/cg_doctor.go`

---

## 1. Objective

Modernize `cli/cmdcg/cg_doctor.go` by integrating `termout.RenderTokyoProgress` for visual repository health reporting, containerizing diagnostic summaries with `termout.TokyoBoxStyle`, and standardizing issue badges with Tokyo status tokens.

---

## 2. Scope & Implementation Details

1. **Overall Health Progress Bar**:
   - In `renderDoctorOverall(totalRepos, healthyRepos int)`:
     - Compute percentage and progress bar using `termout.RenderTokyoProgress(healthyRepos, totalRepos, 24)`.
     - Output progress indicator:
       ```
       Health: [████████████████░░░░░░░░] 66% (2/3 repositories healthy)
       ```

2. **Containerized Summary Box**:
   - Wrap diagnostic summary block in `termout.TokyoBoxStyle.Render(...)` to present a cohesive terminal card.
   - Decompose box content generation into a helper function `buildDoctorSummaryContent(...) string` $\le 12$ lines.

3. **Status Badges & Issue Formatting**:
   - In `renderDoctorDiagnosis(diag CGDoctorDiagnosis)`:
     - Replace hardcoded string symbols with `termout.TokyoOKPill("HEALTHY")` and `termout.TokyoWarnPill("ACTION")`.
     - Display suggestions with `termout.TokyoOKPill("FIX")`.

4. **Leaf Error Sovereignty**:
   - Ensure `RunCGDoctor(repos []string) error` returns `nil` on success or returns `error` directly to root command dispatch without invoking exit handlers or `os.Exit`.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines.
- **Zero Nested Ifs**: Guard clauses and early returns only.
- **Affirmative Booleans**: `isHealthy`, `hasSuggestion`, `hasDrift`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] `RenderTokyoProgress` integrated into `renderDoctorOverall`.
- [x] `TokyoBoxStyle` applied to overall summary output.
- [x] Issue formatting uses Tokyo status badge pills.
- [x] Error management maintains leaf error sovereignty.
- [x] All functions $\le 15$ lines.
