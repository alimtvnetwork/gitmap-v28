# Task 3: installer_history_tree.go — History with Tree

Read the plan at `.lovable/plans/pending/03-tree-view-installer-help.md` first.

1. Create `gitmap/cmd/installer_history_tree.go`.
2. Implement `printInstallerHistoryTree(db *store.DB)`.
3. Read the install history from `db.ListInstallHistory()` (or the equivalent method — look in `gitmap/store/` to find the correct method name).
4. For each unique entry (by slug, latest timestamp): if profile slug → `resolveProfileTree` + render tree. Else: single `└──` description.
5. Separate entries with a DIM `------` divider.
6. Wire into `installer ls` with `--tree` flag in `installer_ls.go`.
7. Compile: run `go build ./...` inside `gitmap/` to verify.
8. Track in `.lovable/temp-agents/task-03-history-tree.md`.
