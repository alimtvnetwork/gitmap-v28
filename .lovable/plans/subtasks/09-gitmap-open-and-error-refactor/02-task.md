# Subtask 2: AST Refactor Script for Handlers

1. Create a script in `.lovable/temp-scripts/refactor_cmds.go`.
2. The script must iterate over all `gitmap/cmd/*.go` files.
3. Replace all `func run[A-Z][a-zA-Z0-9_]*\(args \[\]string\)` with `func run...(args []string) error`.
4. Replace all `func run[A-Z][a-zA-Z0-9_]*\(\)` with `func run...() error`.
5. For all commands that call `os.Exit(0)` or `os.Exit(1)`, attempt to replace them with `return err` or `return nil`.
6. Add `return nil` at the end of every `runXxx` function.
7. Update the dispatch tables (e.g. `roottooling.go`, `rootcore.go`) to match the new signature:
   `func() error { return runXxx(...) }` or similar.
8. Execute the script, fix any compilation errors manually via `go build ./...`.
