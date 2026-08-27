# Subtask 2: Scaffold dd Command

1. In gitmap/constants/constants_cli.go, add CmdAdd = "add".
2. Add to gitmap/constants/cmd_constants_test.go 	opLevelCmds map.
3. In gitmap/cmd/roottooling.go, register to 	oolingUtilEntries.
4. Create gitmap/cmd/add/add.go.
5. Command handles common-attr (generates .gitattributes) and common-ignore (generates .gitignore).
6. Create entrypoint gitmap/cmd/add_entry.go.
