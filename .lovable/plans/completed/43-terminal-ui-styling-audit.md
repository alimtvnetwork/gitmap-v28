# Plan 27: Terminal UI, CLI Styling, Lipgloss & Animation Architecture Audit

## Overview

Comprehensive audit and enhancement of Terminal UI layouts, bright bold ANSI palettes (9X codes), Catppuccin pastel cycling, responsive 2-column alignment (capped at 26), Lipgloss rendering, and clean version footers.

---

## Phase 1: Terminal UI Violation Ledger

| Command / View | File Path | Line | Current UI Pattern | Defect / Limitation | Planned UI Enhancement | Status |
|---|---|:---:|---|---|---|:---:|
| `rootusage.go` | `gitmap/cmd/rootusage.go` | 13 | 2-pass width calculation | Wide column gaps | Max column cap 26, clean column 30 alignment | DONE |
| `rootusage_categories.go` | `gitmap/cmd/rootusage_categories.go` | 1 | Category Banners | Unstructured groups | Implemented `━━ INTENT BANNERS ━━` | DONE |
| `rootusage_rendering.go` | `gitmap/cmd/rootusage_rendering.go` | 20 | Lipgloss Wrapping | Long description clipping | 4-space multi-line description wrap | DONE |
| `constants_colors.go` | `gitmap/constants/constants_colors.go` | 1 | ANSI Palette | Legacy 3X dull colors | Bright bold 9X ANSI & Catppuccin pastel palette | DONE |

---

## Subtasks Breakdown

- [x] [01-terminal-palette-and-banners.md](.lovable/plans/subtasks/27-terminal-ui/01-terminal-palette-and-banners.md) — Bright bold ANSI 9X codes and super-category intent banners.
- [x] [02-lipgloss-two-pass-rendering.md](.lovable/plans/subtasks/27-terminal-ui/02-lipgloss-two-pass-rendering.md) — 2-pass column measurement with 26-char cap and responsive description wrapping.
- [x] [03-subcommand-markers-and-footers.md](.lovable/plans/subtasks/27-terminal-ui/03-subcommand-markers-and-footers.md) — Expandable subcommand markers (`▸ subcommands`) and standard version footers.
- [x] [04-verification.md](.lovable/plans/subtasks/27-terminal-ui/04-verification.md) — Verify all terminal UI test suites and linters pass cleanly.
