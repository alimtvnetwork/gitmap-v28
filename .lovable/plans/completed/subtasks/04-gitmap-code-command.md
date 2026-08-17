# Step 4: gitmap code command

Added `gitmap code` command implementation (actually, it already existed in `cmd/code.go`).
I added `vcode` and `vscode` aliases to the dispatch table (`cmd/rootcore.go`) and the cli constants (`constants/constants_cli.go`).
Also added `CmdCodeAlias` and `CmdCodeAlias2` to the test suite in `constants/cmd_constants_test.go`.
Finally, I created the focused unit test `code_test.go` to test the `projects.json` merge paths logic (`mergeStringPaths`). Tests pass successfully.
