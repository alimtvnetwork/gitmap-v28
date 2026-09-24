# Subtask [03]: Git Pull Efficient SSH JSON Response and Table Formatter
Traceability ID: Task-04
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: cli/cmd/pull_all_efficient.go, cli/cmdssh/ssh_pull_json.go
Action: Extend `gitmap pull all-efficient` with remote SSH JSON response mode, communicating with nodes via structured JSON and rendering the output in a clean, unified termtable.
Acceptance Criteria:
1. Remote nodes return structured JSON summaries.
2. Local receiver deserializes JSON responses and presents tables matching local UI standards.
Targeted Verification: go test -v ./cli/cmd/...
