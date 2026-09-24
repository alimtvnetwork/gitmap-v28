# Subtask [05]: LLM Training Error Standardization to appfault.AppError
Traceability ID: Task-06
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: cli/cmd/llm/llm.go, cli/helptext/llm.md
Action: Replace legacy patterns (`app.app`) in `gitmap llm train` dataset and help documentation with repository-standard `appfault.AppError` and structured monadic Result containers.
Acceptance Criteria:
1. All `app.app` occurrences eliminated from `cli/cmd/llm/` and help documentation.
2. Training samples adhere strictly to `02-spec/03-error-manage/`.
Targeted Verification: go test -v ./cli/cmd/llm/...
