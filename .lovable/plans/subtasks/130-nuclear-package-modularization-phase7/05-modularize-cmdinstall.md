# Subtask 05: Modularize cmdinstall Subpackage

## Objective
Extract ~81 tooling and shell context installation files from `cli/cmd` into dedicated `cli/cmdinstall` subpackage.

## Execution Details
1. Create `cli/cmdinstall/`.
2. Move `install*.go` and tests from `cli/cmd` to `cli/cmdinstall`.
3. Provide `exports.go` and `helpers.go` for public types, profile trees, and Cobra commands.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`.
5. Verify zero circular imports via `go vet -C cli ./cmdinstall`.

## Status: COMPLETED
- Extracted 83 files (install*.go, osdetect*.go, agy_install.go) into `cli/cmdinstall`.
- Provided `exports.go`, `helpers.go`, and `testhelpers_test.go`.
- Decoupled installer dependencies with `cli/cmdinstaller` DAG structure.
- Added bridge forwarders and type aliases in `cli/cmd/clihelpers.go`.
- `go vet -C cli ./cmdinstall ./cmd` and `go vet -C cli ./...` verified 100% clean.
- `gitmap install --help` and `gitmap installer --help` verified working.

