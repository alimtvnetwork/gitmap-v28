# Subtask [04]: Supabase Integration Scaffold and Tool Registry Placeholder
Traceability ID: Task-05
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: cli/constants/constants_install.go, cli/cmdinstall/installlist.go
Action: Add Supabase to GitMap installation categories and tools registry flagged with `[?]` (Planned / To-Do) awaiting future architectural details.
Acceptance Criteria:
1. Supabase appears under database/backend tools with placeholder indicator `[?]`.
2. Clean integration into `gitmap install ls` without breaking profile resolution.
Targeted Verification: go test -v ./cli/cmdinstall/...
