# Subtask 1: Scaffold older Command

1. In gitmap/constants/constants_cli.go, add CmdFolder = "folder" and descriptions.
2. In gitmap/cmd_constants_test.go, add it to the topLevelCmds map (or skip).
3. Create gitmap/cmd/folder/folder.go.
4. Parse args: <dir>, <output_file>, -exclude <pattern> <0|1>.
5. Implement recursive directory walking (s.WalkDir) avoiding .git.
6. Implement logic to write a text tree (.txt, .md).
