# Task 5: schedule_tree.go + Help Import/Export

Read the plan at `.lovable/plans/pending/03-tree-view-installer-help.md` first.

1. Create `gitmap/cmd/schedule_tree.go`.
2. Implement `printScheduleSummaryTree(taskName, interval, shellType string, steps []string)`.
3. Render a tree with header (task name + interval), each step as a tree node.
4. Call from `schedule_cmd.go` at the interactive confirmation stage.
5. Additionally, modify `gitmap/cmd/rootusage_groups.go` (or `rootusagecompact.go`) to add an `Import / Export` section to the help output listing `import-export`, `export`, and `import` commands.
6. Compile: run `go build ./...` inside `gitmap/` to verify.
7. Track in `.lovable/temp-agents/task-05-schedule-tree.md`.
