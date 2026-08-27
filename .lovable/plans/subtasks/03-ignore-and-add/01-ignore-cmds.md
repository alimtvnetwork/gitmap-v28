# Subtask 1: Scaffold ignore & ignore-rm Commands

1. In gitmap/constants/constants_cli.go, add CmdIgnore = "ignore" and CmdIgnoreRm = "ignore-rm".
2. Add to gitmap/constants/cmd_constants_test.go 	opLevelCmds map.
3. In gitmap/cmd/roottooling.go, register them to 	oolingUtilEntries.
4. Create gitmap/cmd/ignore/ignore.go and gitmap/cmd/ignore/ignorerm.go.
5. ignore appends the pattern to .gitignore idempotently.
6. ignore-rm runs git filter-branch to wipe the pattern, then calls ignore.
7. Create entrypoints gitmap/cmd/ignore_entry.go and gitmap/cmd/ignorerm_entry.go.
