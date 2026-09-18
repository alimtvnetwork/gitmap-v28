# Subtask 03: Modularize cmdchrome Subpackage

## Objective
Extract ~22 chrome command files from `cli/cmd` into dedicated `cli/cmdchrome` subpackage.

## Execution Details
1. Create `cli/cmdchrome/`.
2. Move `chrome*.go` and tests from `cli/cmd` to `cli/cmdchrome`.
3. Provide `exports.go` and `helpers.go` for public types, flags, and Cobra commands.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`.
5. Verify zero circular imports via `go vet -C cli ./cmdchrome`.
