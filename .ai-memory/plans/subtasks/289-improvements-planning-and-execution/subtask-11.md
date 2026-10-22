# Subtask 11: Tokyo Night Truecolor Palette Token Constants

> **Subtask ID:** Subtask-11  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/termout/themes.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Define Tokyo Night Storm truecolor hexadecimal constants in `cli/termout/themes.go`.
2. Provide standardized tokens (`TokyoBase`, `TokyoSurface`, `TokyoPrimary`, `TokyoAccent`, `TokyoOK`, `TokyoWarn`, `TokyoError`, `TokyoMuted`, `TokyoFG`) as the single source of truth for CLI theming.
3. Ensure exact alignment with the canonical design system specification (`02-spec/07-design-system/16-theme-catalog-and-palettes.md`).
4. Keep constants exported for consumption across `termout`, `cmdcg`, `heatmap`, and other command packages.

---

## 2. Detailed Technical Plan

### 2.1 Constant Definitions in `cli/termout/themes.go`
- Add the following constant block to `cli/termout/themes.go`:
  ```go
  // Tokyo Night truecolor palette constants (24-bit hex).
  const (
  	TokyoBase    = "#1a1b26" // Deep canvas background
  	TokyoSurface = "#24283b" // Surface card / elevated container
  	TokyoPrimary = "#7aa2f7" // Primary blue for key headings and commands
  	TokyoAccent  = "#bb9af7" // Accent purple for parameters and highlights
  	TokyoOK      = "#9ece6a" // Positive emerald green for active/passed statuses
  	TokyoWarn    = "#e0af68" // Amber warning for dirty/pending states
  	TokyoError   = "#f7768e" // Coral red for errors and missing items
  	TokyoMuted   = "#565a6e" // Muted gray for tree lines and subtitles
  	TokyoFG      = "#c0caf5" // High-contrast foreground text
  )
  ```

### 2.2 Palette Registry Synchronization
- Cross-reference `tokyoNightTheme` Palette map with the new token constants to ensure zero drift between hex values and theme registrations.

---

## 3. Strict Coding Guidelines Compliance
- **File Sizing:** `themes.go` remains clean and modular.
- **Naming Conventions:** All constants use PascalCase with `Tokyo` prefix.
- **Vertical Spacing:** Single blank lines before and after constant blocks.

---

## 4. Verification Criteria
- [x] Tokyo Night constants compile cleanly under Go 1.22+.
- [x] All 9 tokens match the canonical design system values.
- [x] Constants are accessible from external packages (`termout.TokyoPrimary`, etc.).
- [x] No regression in `ResolveTheme()` or `tokyoNightTheme` initialization.
