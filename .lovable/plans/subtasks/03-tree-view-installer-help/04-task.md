# Task 4: macro_tree.go — Macro Steps Tree

Read the plan at `.lovable/plans/pending/03-tree-view-installer-help.md` first.

1. Create `gitmap/cmd/macro_tree.go`.
2. Implement `printMacroStepsTree(m *macro.Macro)` — iterate m.Steps, print each step using `├──` (or `└──` for last).
3. Color: CYAN for glyphs, WHITE for step name, DIM for description.
4. Call `printMacroStepsTree` from `runExecuteCmd` in `macro_cmd.go` before actual execution (only when steps > 0).
5. Max 15 lines per function.
6. Compile: run `go build ./...` inside `gitmap/` to verify.
7. Track in `.lovable/temp-agents/task-04-macro-tree.md`.
