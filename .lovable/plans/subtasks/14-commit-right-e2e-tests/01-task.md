STATUS: DONE
# Subtask 1: Audit and Fix `committransfer` Errors

1. Search `gitmap/cmd/committransfer.go` and `gitmap/committransfer/*.go` for `fmt.Errorf`.
2. Replace them with `apperror.Wrap` where appropriate.
3. Ensure no `fmt.Println` or raw `os.Exit` bypass the global error handler unless absolutely necessary for CLI entrypoints.

