# Subtask [06]: VM / SSH Isolated End-to-End Test Suite
Traceability ID: Task-07
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: cli/tests/e2e/ssh_fleet_e2e_test.go, cli/tests/e2e/prompt_inject_e2e_test.go
Action: Author end-to-end tests exercising prompt injection and parallel SSH delegation across nodes without storing any passwords, gated behind `//go:build e2e`.
Acceptance Criteria:
1. Tests tagged with `//go:build e2e` so standard `go test ./...` and CI/CD skips them by default.
2. Zero passwords, credentials, or secrets in test files.
3. Tests validate prompt injection and parallel dispatch logic.
Targeted Verification: go test -tags=e2e -v ./cli/tests/e2e/...
