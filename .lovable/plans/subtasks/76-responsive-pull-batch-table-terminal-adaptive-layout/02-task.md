# Subtask 02: Adaptive Row Rendering & Branch Formatting

## Objective
Update row rendering in `gitmap/cmd/pull_table_row.go` and formatting helpers in `gitmap/cmd/pull_table_format.go` to adhere to the responsive layout model.

## Target Files
- `gitmap/cmd/pull_table_row.go`
- `gitmap/cmd/pull_table_format.go`

## Implementation Details
1. In `pull_table_format.go`:
   - Add `formatCombinedBranch(branch, latest string, maxLen int) string`:
     - If `latest == ""` or `strings.EqualFold(branch, latest)`: return `formatBranchName(branch, maxLen)`.
     - Otherwise: format as `<cleanBranch>→<cleanLatest>` and middle-truncate to `maxLen`.
2. In `pull_table_row.go`:
   - Refactor `PrintRow(r model.PullTableRow)`:
     - Check `l.IsWide`:
       - If wide: print full 7 columns using `l.ColGap` (3 spaces).
       - If compact: print 6 columns (`REPO`, `BRANCH [combined]`, `PR/TRACK`, `STATUS`, `SHA`, `TIME`) using `l.ColGap` (2 spaces).
     - Ensure ANSI color codes from `resolvePullStatusStyle` and `resolveRepoStatusStyle` use `calcAnsiPadding`.
     - Verify visible printed width never exceeds `l.TermWidth - 2`.
