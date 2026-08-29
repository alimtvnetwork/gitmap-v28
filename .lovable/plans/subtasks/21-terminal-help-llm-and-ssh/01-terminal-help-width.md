# Subtask 21.01: Compact Terminal Help Formatting & Alignment

## Target Files
- `gitmap/cmd/rootusage.go`
- `gitmap/cmd/rootusage_groups.go`
- `gitmap/constants/constants_bookmark.go`
- `gitmap/constants/constants_version_history_cmd.go`

## Actions
- [ ] Reduce `maxCmdColumnWidth` to 26–28 cols (or compute per-group max cmd width capped at 28).
- [ ] Ensure command syntax and descriptions in Bookmark, History, and Data groups are separated by double-spaces.
- [ ] Long commands (> 28 chars) print on line 1 with description on line 2 at standard 4-space indent.
- [ ] Verify `commit-left (cml)` formatting in `printGroupCommitXfer`.

## Acceptance Criteria
- [ ] Column gap between short command names and descriptions is 2–6 spaces instead of 30+ spaces.
- [ ] All rows in `bookmark`, `history`, and `data` align properly.
