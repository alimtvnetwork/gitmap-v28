# Subtask 06: Isolated E2E Test Suite & Minor Version Bump Release

- **Task-ID**: Task-06
- **Subtask**: 06-isolated-e2e-tests-and-release
- **Owner**: Lead Orchestrator
- **Status**: PENDING

## Owned Files (Exclusive)
- `cli/tests/e2e/goemailsender_tempe2e_test.go`
- `version.json`
- `docs/changelog/` or release notes

## Read-Only Context Files
- `02-spec/21-app/goemailsender-connection-and-testing/03-migration-and-verification-spec.md`
- `version.json`

## Objectives
1. Implement isolated end-to-end integration tests in `cli/tests/e2e/goemailsender_tempe2e_test.go` guarded by `//go:build e2e` and `GITMAP_RUN_E2E=1` checking binary invocation, repo flat commit creation, email scheduling simulation, and delegation.
2. Verify that routine `go test ./cli/...` remains unaffected and fast.
3. Perform minor version bump to `v6.575.0` in `version.json`.
4. Stage, format, and push atomic commit via `gitmap cpf "goemailsender - implement connection fallback delegation and email ui"`.
