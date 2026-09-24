# Subtask [01]: Antigravity Prompt Injection and Target Selection
Traceability ID: Task-02
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: cli/cmdagy/agy_inject.go, cli/cmdagy/agy_cmd.go
Action: Support targeting arbitrary projects via `--project` / `-p` and specific or default conversations via `--conversation` / `-c` for prompt injection, auditing list and template prompt workflows.
Acceptance Criteria:
1. `gitmap agy prompt inject` accepts `--project` and `--conversation` flags.
2. Injects prompt payload to target project directory and conversation context cleanly.
3. Help text updated to show target conversation and project injection usage.
Targeted Verification: go test -v ./cli/cmdagy/...
