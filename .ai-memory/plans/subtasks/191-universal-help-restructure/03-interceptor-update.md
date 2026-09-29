# Subtask 191.3: Update Universal Interceptor in cli/cmd/root.go

## Objective
Ensure `tryInterceptCommandHelp` in `cli/cmd/root.go` universally intercepts `gitmap <command> help` and `gitmap <command> --help` for all commands.

## Requirements
1. If the command has a core rich topic, render it.
2. If the command has a dedicated markdown help file in `cli/helptext/`, render it via restructured `helptext.Print`.
3. If the command is registered in `completion.AllCommands()`, render dynamic catalog help.
4. Exit cleanly with code 0 without invoking destructive or missing-argument command handlers.
