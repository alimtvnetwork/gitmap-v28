# Migration and Verification Specification: GoEmailSender, Multi-Profile Fallback & Fleet Delegation

## 1. Migration Strategy

1. **Zero Breaking Changes**:
   - Existing invocations of `gitmap space send-email` and `gitmap space schedule-email` retain exact syntax and parameter contracts.
   - `gitmap clone <target>` continues to work unchanged when target is accessible by default credentials.
   - Database operations continue to support existing SQLite tables; path metadata additions are purely additive in JSON envelopes.
2. **Backward Compatibility**:
   - Unix and macOS binary resolution remain functional.
   - In environments without `go-email-reader`, fallback to built-in SMTP sender remains operational.
   - In environments with a single profile, fallback immediately exits on error without unnecessary probing overhead.

## 2. Verification Protocol

### Test Isolation Rule
All end-to-end integration tests touching external binaries (`go-email-reader`), live network connections, or mock SSH servers MUST be isolated:
- Build tag: `//go:build e2e` OR environment variable guard:
  ```go
  if os.Getenv("GITMAP_RUN_E2E") != "1" {
      t.Skip("skipping e2e integration test; run with GITMAP_RUN_E2E=1")
  }
  ```
- This guarantees that routine local tests and CI/CD pipelines (`go test ./...`) remain fast, deterministic, and green.

### Acceptance Verification Matrix

| ID | Verification Target | Verification Command / Check | Expected Result |
| :--- | :--- | :--- | :--- |
| V-01 | Binary resolution on Windows | `go test -run TestResolveGoEmailReaderWindows ./cli/cmdspace/...` | PASS (detects `.exe` and GOPATH) |
| V-02 | Multi-login fallback unit test | `go test -run TestMultiProfileCloneFallback ./cli/cmdclone/...` | PASS (probes secondary profile on 401/404) |
| V-03 | Email UI help and dry-run | `gitmap space email ui --help` | Exit 0 with interactive UI options |
| V-04 | Relative DB path in Heatmap | `gitmap heatmap --json` | JSON output contains relative `database` path |
| V-05 | Relative DB path in Agent | `gitmap agent task status --json` | JSON output contains `relativeDbPath` & `repoSecretsUrl` |
| V-06 | Fleet Nodes command routing | `go test -run TestNodesRemoteCommandResolution ./cli/cmdnodes/...` | PASS (recognizes `space`, `heatmap`, `agent`) |
| V-07 | E2E Isolated Test Suite | `go test -v -tags=e2e ./cli/tests/e2e/...` (with `GITMAP_RUN_E2E=1`) | PASS (all E2E assertions verified) |
