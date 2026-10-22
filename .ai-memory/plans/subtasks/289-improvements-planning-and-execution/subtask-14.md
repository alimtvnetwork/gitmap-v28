# Subtask 14: Help Menu Sectional Badge Rendering

> **Subtask ID:** Subtask-14  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/termout/render.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Update `PrintSection` in `cli/termout/render.go` to support sectional badges.
2. Render badges in clean bracketed format: `  <SectionTitle> [<BadgeText>]:`.
3. Support custom badge colors (`sec.BadgeColor`), with graceful fallback to `RoleFGSubtle` / muted color.
4. Ensure omission of brackets when `sec.BadgeText` is empty.
5. Adhere to strict 15-line function limit by extracting header formatting into a dedicated helper.

---

## 2. Detailed Technical Plan

### 2.1 Refactor `PrintSection` in `cli/termout/render.go`
- Extract section header line formatting to `formatSectionHeader`:
  ```go
  // PrintSection renders a section with its header and entries.
  func PrintSection(sec HelpSection, colWidth int, theme Theme) {
  	headerColor := sec.Color
  	if len(headerColor) == 0 {
  		headerColor = theme.HeaderColor
  	}

  	fmt.Fprint(output.UI(), formatSectionHeader(sec, headerColor))
  	for _, entry := range sec.Entries {
  		printEntryLine(entry, colWidth, theme)
  	}
  	fmt.Fprintln(output.UI())
  }
  ```

### 2.2 Implement `formatSectionHeader` Helper
- Build clean bracketed header string:
  ```go
  func formatSectionHeader(sec HelpSection, headerColor string) string {
  	if sec.BadgeText == "" {
  		return fmt.Sprintf("  %s%s:%s\n", headerColor, sec.Title, Reset())
  	}

  	badgeColor := sec.BadgeColor
  	if badgeColor == "" {
  		badgeColor = RoleSGR(RoleFGSubtle)
  	}

  	return fmt.Sprintf(" %s[%s]%s", badgeColor, sec.BadgeText, Reset())
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Function Limit:** `PrintSection` is 11 lines; `formatSectionHeader` is 12 lines (both <= 15 lines).
- **Conditional Flattening:** Guard clause on `sec.BadgeText == ""` returns early.
- **Return Formatting:** Single blank line before return statements and after `}` blocks.
- **Color Safety:** Emits `Reset()` after styled segments to avoid color bleed.

---

## 4. Verification Criteria
- [x] Sections without badges render identically to legacy output (`  Title:`).
- [x] Sections with `BadgeText` render formatted brackets (`  Title [BADGE]:`).
- [x] Badge colors correctly reflect `BadgeColor` when specified.
- [x] Output suppresses ANSI escapes when redirected or piped.
