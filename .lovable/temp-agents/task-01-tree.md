STATUS: DONE

Task 1 completed successfully.
- Created `gitmap/cmd/installer_tree.go`.
- Defined `InstallerTreeNode` struct with `Title`, `Description string`, and `Children []InstallerTreeNode`.
- Implemented `printInstallerTree(root InstallerTreeNode, prefix string, isLast bool)` with UTF-8 box-drawing (`constants.TreeBranch` and `constants.TreeCorner`).
- Applied color styling using `constants.ColorCyan` for tree connectors, `constants.ColorWhite` for node titles, `constants.ColorDim` for descriptions, and `constants.ColorReset` for resetting ANSI escape codes.
- Implemented `printInstallSummaryHeader(slug string)` outputting `Installation Summary:` with a green checkmark (`constants.ColorGreen`).
- Strictly maintained function length <= 15 lines and semantic boolean naming conventions (`isLast`, `isChildLast`).
- Verified build and test suite with `go build ./...` and `go test ./...`.
