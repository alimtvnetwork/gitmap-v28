# Subtask 3: Command Implementation

1. Create gitmap/cmd/ag/ag.go and gitmap/cmd/vscode/vscode.go.
2. Parse args: if install, invoke gitmap install antigravity or gitmap install vscode.
3. If no args, run g or code respectively via xec.Command.
4. Create entrypoints g_entry.go and scode_entry.go.
