# Subtask: Inverted Booleans (Go)
Status: ✅ Done

## Steps
1. Edit gitmap/cmd/installctx_linux_e2e_test.go: line 213, extract !hasZenity to isZenityMissing := !hasZenity
2. Edit gitmap/cmd/prettyflag.go: line 75, extract !hasPrettyPrefix(arg) to isMissingPrefix := !hasPrettyPrefix(arg)
3. Audit and fix all remaining inverted booleans across Go codebase.
