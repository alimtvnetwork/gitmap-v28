# Task 2: install_profile_tree.go — Profile Hierarchy Tree

Read the plan at `.lovable/plans/pending/03-tree-view-installer-help.md` first.

1. Create `gitmap/cmd/install_profile_tree.go`.
2. Define structs: `ProfileComposition{Name, Alias, Description string, Base *ProfileComposition, Tools []ToolEntry}` and `ToolEntry{Slug, Description string}`.
3. Implement `resolveProfileTree(slug string) (ProfileComposition, bool)` hardcoding the ubuntu-basic → ubuntu+vscode → ubuntu+small-dev → ubuntu+dev hierarchy.
4. Implement `printProfileTree(p ProfileComposition)` using `printInstallerTree` from Task 1.
5. Implement `printProfileInstallSummary(slug string)` — calls resolve + print.
6. Compile: run `go build ./...` inside `gitmap/` to verify.
7. Track in `.lovable/temp-agents/task-02-profile-tree.md`.
