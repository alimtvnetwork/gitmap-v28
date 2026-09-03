# 41 - Terminal UI, CLI Styling, Lipgloss & Animations Audit Specification

## 1. Verbatim Acceptance Criteria Echo (from spec/02-coding-guidelines/01-cross-language/97-acceptance-criteria.md)

### AC-01: Guideline Coverage
- [ ] Boolean principles define naming, evaluation, and composition patterns
- [ ] Casting elimination patterns cover type-safe alternatives to type assertions
- [ ] Code style defines formatting, naming, and structural conventions

### AC-02: Enforcement
- [ ] All guidelines include ❌ (forbidden) and ✅ (compliant) code examples
- [ ] ESLint/linter rules are documented for automated enforcement
- [ ] Master guidelines document consolidates all standards for AI reference

---

## 2. Task-Specific Rule Set (Domain Rules)

1. **Rule TUI-1 (Bright & Bold ANSI Palette):** Use bright bold 9X ANSI sequences (`\033[1;9Xm`) with mandatory `ColorReset` (`\033[0m`) suffixes.
2. **Rule TUI-2 (Catppuccin Pastel Palette):** Utilize TrueColor pastel colors (`ColorPastelGreen`, `ColorPastelCyan`, `ColorPastelYellow`, `ColorPastelMagenta`) and `ColorCycle` for multi-item list rotations.
3. **Rule TUI-3 (Super-Category Intent Banners):** Organize commands under bold intent banners (`━━ SECTION ━━━━━━━━━━`).
4. **Rule TUI-4 (Responsive 2-Column Alignment):** Implement 2-pass width calculation (measurement pass + render/wrap pass) with `maxCmdColumnWidth = 26`.
5. **Rule TUI-5 (Expandable Markers & Footers):** Mark nested subcommands with `▸ subcommands` and provide standard version/repository footers.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/constants/constants_terminal.go | 11-40 | ANSI & Catppuccin color definitions | Verified bright bold 9X palette and ColorCycle rotation | VERIFIED |
| V-02 | gitmap/cmd/rootusage.go | Various | Responsive 2-pass column layout and intent banners | Verified 2-pass calculation with 26 col max width | VERIFIED |
| V-03 | gitmap/cloner/progress.go | 66-120 | Terminal progress formatters | Verified Lipgloss spinner and ColorGreen status formatters | VERIFIED |
| V-04 | 03-ai-scripts/06-cicd-local-runner.py | N/A | Full CI quality runner verification | Verified all quality gates pass (exit 0) | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/41-terminal-ui/01-ansi-palette-and-box-drawing.md`)**:
   Verify bright bold ANSI palette and Catppuccin pastel cycling in terminal constants.
2. **Subtask 02 (`.lovable/plans/subtasks/41-terminal-ui/02-responsive-alignment-and-banners.md`)**:
   Verify responsive 2-column help layout and super-category intent banners.
3. **Subtask 03 (`.lovable/plans/subtasks/41-terminal-ui/03-ci-runner-verification.md`)**:
   Execute full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.
