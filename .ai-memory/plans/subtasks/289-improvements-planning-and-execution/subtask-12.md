# Subtask 12: Lipgloss Safe Tokyo & Adaptive Styling Primitives

> **Subtask ID:** Subtask-12  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/termout/theme_lipgloss.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Implement a new dedicated file `cli/termout/theme_lipgloss.go`.
2. Construct TTY-safe Lipgloss styling helper functions (`SafeTokyoStyle`, `SafeAdaptiveStyle`).
3. Guard terminal escapes against redirection/piping using `colorAllowed()` and `IsStdoutTTY()`.
4. Export predefined Tokyo Night styling instances (`TokyoHeaderStyle`, `TokyoAccentStyle`, `TokyoActiveBadgeStyle`, `TokyoMissingBadgeStyle`, `TokyoWarnBadgeStyle`, `TokyoMutedStyle`).

---

## 2. Detailed Technical Plan

### 2.1 File Creation: `cli/termout/theme_lipgloss.go`
- Package: `package termout`
- Imports:
  ```go
  import (
  	"github.com/charmbracelet/lipgloss"
  )
  ```

### 2.2 Helper Implementation
- Implement `SafeTokyoStyle`:
  ```go
  // SafeTokyoStyle returns a Lipgloss style configured with foreground colorHex,
  // falling back to an unstyled empty style when colors are disallowed.
  func SafeTokyoStyle(colorHex string) lipgloss.Style {
  	if !colorAllowed() {
  		return lipgloss.NewStyle()
  	}

  	return lipgloss.NewStyle().Foreground(lipgloss.Color(colorHex))
  }
  ```
- Implement `SafeAdaptiveStyle`:
  ```go
  // SafeAdaptiveStyle constructs an adaptive Lipgloss style supporting dark and light backgrounds.
  func SafeAdaptiveStyle(darkHex, lightHex string) lipgloss.Style {
  	if !colorAllowed() {
  		return lipgloss.NewStyle()
  	}

  	return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
  		Dark:  darkHex,
  		Light: lightHex,
  	})
  }
  ```

### 2.3 Predefined Styled Elements
- Define standard styles for cross-package use:
  ```go
  var (
  	TokyoHeaderStyle       = SafeTokyoStyle(TokyoPrimary).Bold(true)
  	TokyoAccentStyle       = SafeTokyoStyle(TokyoAccent)
  	TokyoActiveBadgeStyle  = SafeTokyoStyle(TokyoOK).Bold(true)
  	TokyoMissingBadgeStyle = SafeTokyoStyle(TokyoError).Bold(true)
  	TokyoWarnBadgeStyle    = SafeTokyoStyle(TokyoWarn).Bold(true)
  	TokyoMutedStyle        = SafeTokyoStyle(TokyoMuted)
  )
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Function Limit:** All helper functions are <= 8 lines.
- **Guard Clauses:** Early returns on `!colorAllowed()` with zero nested conditionals.
- **Positive Booleans:** Boolean check uses `colorAllowed()`.
- **Return Formatting:** Single blank line before each return statement.

---

## 4. Verification Criteria
- [x] Compiles cleanly with Lipgloss library.
- [x] Returns unstyled `lipgloss.NewStyle()` when stdout is piped (`colorAllowed() == false`).
- [x] Predefined styles initialize without nil pointer dereferences.
- [x] Truecolor rendering active when running inside an interactive terminal.
