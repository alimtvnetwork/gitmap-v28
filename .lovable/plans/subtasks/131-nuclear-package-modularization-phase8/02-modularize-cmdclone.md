# Subtask 02: Modularize cmdclone Subpackage

## Objective
Extract ~74 repository cloning, checkout, and batch sync files from `cli/cmd` into dedicated `cli/cmdclone` subpackage.

## Execution Details
1. Create `cli/cmdclone/`.
2. Move `clone*.go` and tests from `cli/cmd` to `cli/cmdclone`.
3. Provide `exports.go` and `helpers.go` for public types and Cobra/runner entrypoints.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`.
5. Verify zero circular imports via `go vet -C cli ./cmdclone`.
