# Subtask 1: Central Dispatch Refactor

1. Open `gitmap/cmd/rootdispatch.go`.
2. Change `type dispatchEntry struct { names []string, handler func() error }`.
3. Update `runDispatchTable` to return `(bool, error)`.
4. Open `gitmap/cmd/root.go` and update `dispatch` to check the returned `error`.
5. If an error is returned, print it using `cliexit.Reportf` or `apperror` wrapping, and then `os.Exit(1)`.
