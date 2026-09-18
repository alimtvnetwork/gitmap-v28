# Subtask 03: Modularize cmdupdate Subpackage

## Objective
Extract ~23 updater, cleanup, and self-update files from `cli/cmd` into dedicated `cli/cmdupdate` subpackage.

## Status
COMPLETED

## Execution Details
1. Create `cli/cmdupdate/`. (Done)
2. Move 23 `update*.go` and tests from `cli/cmd` to `cli/cmdupdate`. (Done)
3. Provide `exports.go` and `helpers.go` for update operations. (Done)
4. Wire bridge forwarders in `cli/cmd/clihelpers.go`. (Done)
5. Verify zero circular imports via `go vet -C cli ./cmdupdate ./cmd` and all 38 CI/CD gates pass. (Done)
