# Subtask 02: Modularize cmdinstaller Subpackage

## Objective
Extract ~34 installer files from `cli/cmd` into dedicated `cli/cmdinstaller` subpackage.

## Execution Details
1. Create `cli/cmdinstaller/`.
2. Move `installer*.go` and tests from `cli/cmd` to `cli/cmdinstaller`.
3. Provide `exports.go` and `helpers.go` for public types and Cobra commands.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go` or root command dispatcher.
5. Verify zero circular imports via `go vet -C cli ./cmdinstaller`.
