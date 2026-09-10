# Subtask 05: Unit Tests, Parity Verification & CI Local Runner

## Objective
Author comprehensive unit tests and verify against quality gates:
- Write unit tests in `gitmap/store/purge_history_test.go`:
  - Test table creation, column migration (legacy `ID` to `PurgeHistoryLogId`), insert, retrieve last, mark restored.
  - Assert that `PurgeHistoryLog.Id`, `PurgeHistoryLog.PurgeHistoryLogId`, and `PurgeHistoryLog.IsRestored` are populated properly.
- Write unit tests in `gitmap/cmd/purge_test.go`:
  - Test argument parsing (`--restore`, `--confirm`, `-y`, pattern extraction).
- Run `go vet ./...` and `go test ./gitmap/store/... ./gitmap/cmd/...`.
- Run `python 03-ai-scripts/06-cicd-local-runner.py` to ensure all quality gates exit with code 0.

## Files Affected
- `gitmap/store/purge_history_test.go`
- `gitmap/cmd/purge_test.go`
