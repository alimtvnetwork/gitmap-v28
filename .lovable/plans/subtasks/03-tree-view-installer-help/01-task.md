# Task 1: installer_tree.go — Core Tree Renderer

Read the plan at `.lovable/plans/pending/03-tree-view-installer-help.md` first.

1. Create `gitmap/cmd/installer_tree.go`.
2. Define `InstallerTreeNode` struct: `Title`, `Description string`, `Children []InstallerTreeNode`.
3. Implement `printInstallerTree(root InstallerTreeNode, prefix string, isLast bool)` with UTF-8 box-drawing.
4. Color: `constants.ColorCyan` for `├──`/`└──`, `constants.ColorWhite` for title, `constants.ColorDim` for description, `constants.ColorReset` to reset.
5. Implement `printInstallSummaryHeader(slug string)` that prints `  Installation Summary:` with a green checkmark.
6. Max 15 lines per function. Boolean naming rules: `is`/`has` prefix only.
7. Compile: run `go build ./...` inside `gitmap/` to verify.
8. Track in `.lovable/temp-agents/task-01-tree.md`.
