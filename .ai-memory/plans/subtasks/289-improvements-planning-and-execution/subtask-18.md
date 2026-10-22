# Subtask 18: Coding Guidelines Status Badges Modernization

> **Subtask ID:** Subtask-18  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/cmdcg/cg_version.go`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Update multi-repo status badges in `cli/cmdcg/cg_version.go` (`gitmap cg status`).
2. Replace hardcoded Dracula hex color constants with centralized Tokyo Night Lipgloss styles from `termout`.
3. Standardize active badge (`✔ active`) and missing badge (`✖ missing`) styling.
4. Ensure strict TTY safety so non-interactive output produces unescaped text.

---

## 2. Detailed Technical Plan

### 2.1 Refactor Style Declarations in `cli/cmdcg/cg_version.go`
- Replace legacy style declarations:
  ```go
  // Current legacy definitions:
  var (
  	cgStatusHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
  	cgActiveBadgeStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
  	cgMissingBadgeStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
  	cgDateStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
  )
  ```
- With centralized `termout` Tokyo styles:
  ```go
  var (
  	cgStatusHeaderStyle = termout.TokyoHeaderStyle
  	cgActiveBadgeStyle  = termout.TokyoActiveBadgeStyle
  	cgMissingBadgeStyle = termout.TokyoMissingBadgeStyle
  	cgDateStyle         = termout.TokyoMutedStyle
  	cgVersionStyle      = termout.TokyoAccentStyle
  )
  ```

### 2.2 Row Rendering Updates
- In `renderStatusActiveRow`:
  ```go
  func renderStatusActiveRow(name string, meta *CGMetadata) string {
  	badge := cgActiveBadgeStyle.Width(12).Render("✔ active")
  	ver := cgVersionStyle.Width(12).Render(meta.Version)
  	date := cgDateStyle.Render(meta.InstalledAt)

  	return fmt.Sprintf("  %-24s %s %s %s", name, badge, ver, date)
  }
  ```
- In `renderStatusMissingRow`:
  ```go
  func renderStatusMissingRow(name string) string {
  	badge := cgMissingBadgeStyle.Width(12).Render("✖ missing")
  	ver := cgDateStyle.Width(12).Render("-")
  	date := cgDateStyle.Render("-")

  	return fmt.Sprintf("  %-24s %s %s %s", name, badge, ver, date)
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **DRY Architecture:** Eliminates duplicate color definitions in command files; binds directly to `termout`.
- **Function Limit:** Functions remain <= 8 lines.
- **Return Formatting:** Single blank line before return statements.
- **Safety:** Automatically adopts `colorAllowed()` checks implemented in Subtask 12.

---

## 4. Verification Criteria
- [ ] Active repositories display bright Tokyo Green (`#9ece6a`) `✔ active` badge.
- [ ] Missing repositories display Tokyo Red (`#f7768e`) `✖ missing` badge.
- [ ] Piped output (`gitmap cg status | cat`) emits plain characters without ANSI code corruption.
- [ ] Column widths and alignment are cleanly preserved.
