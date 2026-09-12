# Subtask 03: Modularize cmdupdate Subpackage

## Objective
Extract ~23 updater, cleanup, and self-update files from `cli/cmd` into dedicated `cli/cmdupdate` subpackage.

## Execution Details
1. Create `cli/cmdupdate/`.
2. Move `update*.go` and tests from `cli/cmd` to `cli/cmdupdate`.
3. Provide `exports.go` and `helpers.go` for update operations.
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`.
5. Verify zero circular imports via `go vet -C cli ./cmdupdate`.
