# Subtask 04: Modularize cmdsetup Subpackage

## Objective
Extract ~21 setup command files from `cli/cmd` into dedicated `cli/cmdsetup` subpackage.

## Execution Details
1. Create `cli/cmdsetup/`.
2. Move `setup*.go` and tests from `cli/cmd` to `cli/cmdsetup`.
3. Provide `exports.go` and `helpers.go` for public types, salts, and Cobra commands.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`.
5. Verify zero circular imports via `go vet -C cli ./cmdsetup`.

## Status: COMPLETED
- Extracted 21 setup files into `cli/cmdsetup`.
- Resolved `assertContains` via local `testhelpers_test.go`.
- Decoupled `warnIfNoWrapper`, `resolveSetupConfigPath`, `isWrapperActive` via `exports.go` and bridge forwarders in `cli/cmd/clihelpers.go`.
- `go vet -C cli ./cmdsetup ./cmd` verified 100% clean.
- `gitmap setup --help` verified working.

