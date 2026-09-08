# Subtask 01: Dynamic Terminal Detection & Adaptive Layout Model

## Objective
Detect terminal width using `golang.org/x/term` and compute adaptive column budgets in `gitmap/cmd/pull_table_layout.go`.

## Target Files
- `gitmap/cmd/pull_table_layout.go`

## Implementation Details
1. Implement `detectTerminalWidth() int`:
   - Calls `term.GetSize(int(os.Stdout.Fd()))`.
   - If width $\le 0$, fallback to `os.Getenv("COLUMNS")`.
   - If missing/invalid, default to 80 columns.
   - Clamp minimum to 60 columns.
2. Extend `PullTableLayout` struct:
   - `TermWidth int`
   - `IsWide bool` (true if `TermWidth >= 105`)
   - `IsCompact bool` (true if `TermWidth < 105`)
   - `ColGap int` (2 spaces in compact, 3 spaces in wide)
   - `DividerLen int`
3. Adaptive sizing algorithm:
   - When `IsWide == true`:
     - Show 7 columns: `REPO`, `BRANCH`, `LATEST BRANCH`, `PR/TRACK`, `STATUS`, `SHA`, `TIME`.
   - When `IsWide == false` (Compact):
     - Calculate remaining width for `REPO` and `BRANCH`:
       Fixed widths: `PR/TRACK` (8), `STATUS` (10), `SHA` (7), `TIME` (4), indent (2), gaps (5 * 2 = 10). Fixed sum = 41.
       Remaining budget = `TermWidth - 2 - 41`.
       Allocate: `MaxRepo = remaining * 55 / 100`, `MaxBranch = remaining * 45 / 100`.
       If `Branch` != `LatestBranch`, format branch as `branch→latest` (middle-truncated) to retain full visibility without line wrapping.
4. Dynamic `PrintHeader()`:
   - Formats columns using computed widths and gap.
   - Prints dynamic divider: `strings.Repeat("-", l.DividerLen)` prefixed with 2 spaces.
   - Total width strictly $\le TermWidth - 2$.
