# Subtask 04: Modularize cmdpull Subpackage

## Objective
Extract ~21 pull and push command files from `cli/cmd` into dedicated `cli/cmdpull` subpackage.

## Execution Details
1. Create `cli/cmdpull/`.
2. Move `pull*.go` and `push*.go` and tests from `cli/cmd` to `cli/cmdpull`.
3. Provide `exports.go` and `helpers.go` for remote branch pulling and pushing.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`.
5. Verify zero circular imports via `go vet -C cli ./cmdpull`.
