# Master Architectural Plan: Responsive Terminal-Adaptive Gitmap Pull Batch Table Layout

## 1. Overview & Root Cause Analysis

### Problem Statement
When `gitmap pull` finishes pulling repositories in batch mode, `RenderPullBatchTable` prints a summary table containing 7 columns: `REPO`, `BRANCH`, `LATEST BRANCH`, `PR/TRACK`, `STATUS`, `SHA`, `TIME`.
Currently, column widths are hardcoded to fixed values with 3-space column gaps and a hardcoded 103-dash divider line, producing an unconditional minimum line length of 107 characters.

In standard terminal environments (such as Windows PowerShell default 80-column buffers, VS Code split terminal panes, or terminals <= 104 columns as shown in user screenshot `media_1788840384813.png`), the output wraps at column 80:
- The header splits across 2 lines (`PR/TRACK` is cut in half into `P` and `R/TRACK`).
- The 103-dash divider line wraps across 2 lines.
- Every repository row wraps into 2 broken, misaligned lines (e.g. `UP_TO_DATE` splits across lines as `UP_T` and `O_DATE`, and SHA / duration wrap onto the second line).

### Root Cause
1. **Zero Terminal Width Awareness**: `PullTableLayout` never queries terminal width (`term.GetSize(int(os.Stdout.Fd()))`) or environment variables (`COLUMNS`), assuming infinite or wide terminal width.
2. **Hardcoded Over-Sized Column Budgets**: Column widths (20 + 16 + 18 + 10 + 10 + 7 + 4 = 85 chars + 18 padding + 2 indent = 105 chars) exceed standard 80-column terminal boundaries.
3. **Hardcoded Static Divider**: `PrintHeader` prints a static 103-character dash string that unconditionally wraps on narrow terminals.
4. **Redundant Dual Branch Columns on Narrow Screens**: In >90% of repos, `BRANCH` and `LATEST BRANCH` are identical (`main`/`main`), consuming 37 columns unnecessarily on narrow screens.

---

## 2. Target Architectural Design

### 2.1 Dynamic Terminal Width Detection
Query terminal width using `golang.org/x/term.GetSize(int(os.Stdout.Fd()))`:
- If valid width $W > 0$ is returned, use $W$.
- If error or non-TTY, check `os.Getenv("COLUMNS")`.
- Default to standard safe terminal width (80 columns).

### 2.2 Adaptive Multi-Tier Layout Engine
Depending on effective terminal width:
1. **Wide Mode ($W \ge 105$)**:
   - Full 7-column layout (`REPO`, `BRANCH`, `LATEST BRANCH`, `PR/TRACK`, `STATUS`, `SHA`, `TIME`).
   - Dynamic proportional sizing matching the exact terminal budget.
2. **Standard / Compact Mode ($72 \le W < 105$)**:
   - 6-column optimized layout:
     - `REPO` (18-22 chars adaptive)
     - `BRANCH` (14-16 chars; displays `active` or `active → latest` if different)
     - `PR/TRACK` (7-9 chars, e.g. `0 PRs`, `6 PRs`, `local`)
     - `STATUS` (10 chars, e.g. `UP_TO_DATE`, `DIRTY`)
     - `SHA` (7 chars)
     - `TIME` (4 chars, e.g. `1.0s`)
   - 2-space column gaps.
   - Total width strictly $\le W - 2$, fitting 80-column terminals with zero line wraps.
3. **Ultra-Compact Mode ($W < 72$)**:
   - 5 essential columns (`REPO`, `BRANCH`, `STATUS`, `SHA`, `TIME`), middle-truncated to fit $W - 2$.

### 2.3 Dynamic Divider Line
Divider length is computed dynamically as:
$$\text{dividerLen} = \text{totalTableWidth}$$
Printed with `strings.Repeat("-", dividerLen)` prefixed with `"  "`, guaranteeing the divider never wraps.

### 2.4 Profile Sign-In Scrubbing Verification
Ensure `patchImportedChromeProfilePreferencesWithOptions` and `scrubImportedPreferencesAuth` scrub stale `account_info`, `sync`, `google`, and set `signin.allowed = false` and `browser.has_seen_welcome_page = true` across all import pipelines.

---

## 3. Task-Specific Rules & Constraints

1. **Zero Line-Wrap Guarantee**: Table rows and dividers MUST NOT exceed effective terminal width $W - 2$ under any circumstances (including 72, 80, and 120 column terminals).
2. **ANSI Alignment Correctness**: Colors and text styles rendered with Lipgloss or ANSI escape codes must calculate padding using `calcAnsiPadding` / visible width so column boundaries never shift.
3. **Strict Function Sizing & Shallow Nesting**: All Go functions $\le 15$ lines; zero nested `if` statements (depth $\le 1$); positive boolean conventions.
4. **Full Test Coverage**: Unit tests simulating narrow (72-col, 80-col) and wide (120-col) terminals verifying zero line wraps and exact column alignment.
5. **Strict Relative Git Paths**: All links and file references in plans and subtasks must be relative to repository root.

---

## 4. Subtask Decomposition

- [01-terminal-width-detection-and-layout-model.md](subtasks/76-responsive-pull-batch-table-terminal-adaptive-layout/01-terminal-width-detection-and-layout-model.md): Add dynamic terminal detection and adaptive column budgeting in `pull_table_layout.go`.
- [02-adaptive-row-rendering-and-branch-collapse.md](subtasks/76-responsive-pull-batch-table-terminal-adaptive-layout/02-adaptive-row-rendering-and-branch-collapse.md): Implement adaptive row formatting, branch merging, and dynamic divider in `pull_table_row.go` and `pull_table_format.go`.
- [03-summary-and-profile-import-sign-in-integration.md](subtasks/76-responsive-pull-batch-table-terminal-adaptive-layout/03-summary-and-profile-import-sign-in-integration.md): Integrate adaptive layout into `pull_table_summary.go` and verify sign-in scrubbing across profile import handlers.
- [04-test-suite-and-quality-gates.md](subtasks/76-responsive-pull-batch-table-terminal-adaptive-layout/04-test-suite-and-quality-gates.md): Author comprehensive unit tests for narrow/standard/wide tables and run CI/CD local runner gates.
