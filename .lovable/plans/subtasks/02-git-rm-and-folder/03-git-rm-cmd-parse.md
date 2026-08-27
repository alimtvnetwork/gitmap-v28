# Subtask 3: Scaffold git-rm Command

1. In gitmap/constants/constants_cli.go, add CmdGitRm = "git-rm".
2. Create gitmap/cmd/gitrm/gitrm.go.
3. Parse input: resolve if it's a JSON file (produced by older), TXT file, CSV, or direct path.
4. Extract the flat list of relative paths to remove.
