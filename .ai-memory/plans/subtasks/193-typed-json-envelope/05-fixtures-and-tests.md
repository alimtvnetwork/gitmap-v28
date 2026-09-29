# Subtask 193.5: Automated Tests, Safe Fixtures & E2E Validation
Traceability ID: Task-05
Spec Reference: [02-spec/21-app/183-typed-json-envelope-and-format-inspection.md](../../../02-spec/21-app/183-typed-json-envelope-and-format-inspection.md)
Target Files: cli/jsonenvelope/fixtures/*.json, cli/cmd/which_format_cmd_test.go, cli/jsonenvelope/envelope_test.go
Action: Create safe mock fixtures (with fake hostnames, non-sensitive tokens). Add comprehensive test cases covering mixed folder inspections (matched vs unmatched), single-line batch command formatting, `-y` advice rendering, and backward-compatible unmarshaling.
Acceptance Criteria:
1. Mock fixtures exist for all supported types without real credentials.
2. Unit tests verify format identification, import commands, and unmatched file handling.
3. Zero-nesting (max depth 1) and positive boolean conventions verified.
Targeted Verification: `go test -v ./jsonenvelope/... ./cmd -run TestWhichFormat`
