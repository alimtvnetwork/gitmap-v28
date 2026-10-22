# Subtask 20: Pull Renderer Theme Token Standardization

> **Subtask ID:** Subtask-20  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/cmdpull/pull_efficient_render_helpers.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Standardize status coloring in `cli/cmdpull/pull_efficient_render_helpers.go` using centralized `termout.RoleSGR` theme tokens.
2. Replace hardcoded `constants.ColorGreen`, `constants.ColorYellow`, `constants.ColorRed`, `constants.ColorDim` in `StyleRepoStatusLabel` and `formatDiffStatsColor`.
3. Support seamless palette switching (`--theme`) and automatic plain-text fallback when stdout is redirected or piped.
4. Keep functions within the 15-line limit with zero nested if statements.

---

## 2. Detailed Technical Plan

### 2.1 Refactor `StyleRepoStatusLabel`
- Replace raw constants with semantic role lookups:
  ```go
  // StyleRepoStatusLabel decorates status labels with theme-aware ANSI colors for fast visual triage.
  func StyleRepoStatusLabel(statusLabel string) string {
  	switch statusLabel {
  	case "up-to-date":
  		return termout.RoleSGR(termout.RoleOK) + "up-to-date" + termout.Reset()
  	case "dirty":
  		return termout.RoleSGR(termout.RoleWarn) + "dirty" + termout.Reset()
  	case "failed":
  		return termout.RoleSGR(termout.RoleError) + "failed" + termout.Reset()
  	}

  	if strings.HasPrefix(statusLabel, "+") {
  		return formatDiffStatsColor(statusLabel)
  	}

  	return statusLabel
  }
  ```

### 2.2 Refactor `formatDiffStatsColor`
- Modernize diff stat coloring (+ins / -del):
  ```go
  func formatDiffStatsColor(statusLabel string) string {
  	parts := strings.SplitN(statusLabel, " ", 2)
  	diffPart := parts[0]
  	slashIdx := strings.Index(diffPart, "/")
  	if slashIdx == -1 {
  		return termout.RoleSGR(termout.RoleOK) + statusLabel + termout.Reset()
  	}

  	ins := termout.RoleSGR(termout.RoleOK) + diffPart[:slashIdx] + termout.Reset()
  	del := termout.RoleSGR(termout.RoleError) + diffPart[slashIdx:] + termout.Reset()
  	coloredDiff := ins + del
  	if len(parts) > 1 {
  		return coloredDiff + " " + termout.RoleSGR(termout.RoleFGSubtle) + parts[1] + termout.Reset()
  	}

  	return coloredDiff
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Error & Theme Sovereignty:** Delegating escape generation to `termout.RoleSGR` respects TTY and theme rules.
- **Function Limit:** `StyleRepoStatusLabel` is 14 lines; `formatDiffStatsColor` is 15 lines.
- **Zero Nested Ifs:** Linear flow with guard clauses and early returns.
- **Return Formatting:** Single blank line before return statements.

---

## 4. Verification Criteria
- [x] Pull progress and results show theme-consistent status indicators.
- [x] Positive insertions show green, deletions show red, and branch counters show subtle foreground.
- [x] Running `gitmap pull ... | cat` outputs clean status text without escape artifacts.
- [x] All functions adhere strictly to coding guidelines and 15-line limit.
