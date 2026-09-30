# Subtask 03: Fleet Clone Suggestions in cmdclone & Unit Tests

## Scope
1. In `cli/cmdclone/clone.go`:
   - Display fleet clone suggestion at the end of successful `gitmap clone` execution.
   - Suppress suggestion if `IsFleetCloneActive()` is true.
2. In `cli/cmdclone/clonefixrepo.go`:
   - Display fleet CFR suggestion at the end of successful `gitmap cfr` / `gitmap cfrp`.
   - Suppress suggestion if `IsFleetCloneActive()` is true.
3. Author `cli/cmdnodes/nodes_clone_test.go`:
   - Test command kind parsing (`clone`, `cfr`, `cfrp`).
   - Test argument decomposition and flag parsing.
   - Test manifest file detection logic.
   - Test remote command generation string formatting.
   - Test help rendering output.
4. Run Go format and repository checks.
