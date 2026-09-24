# Subtask [02]: SSH Fleet Parallel Execution, Node Exclusion, and Live Streaming
Traceability ID: Task-03
Spec Reference: [02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md](../../../02-spec/21-app/143-antigravity-commands-ssh-fleet-and-e2e-benchmarking.md)
Target Files: cli/cmdssh/fleet_parallel.go, cli/cmdssh/ssh_cmd.go, cli/cmdagm/agm_cmd.go, cli/cmdagy/agy_cmd.go
Action: Implement parallel SSH fleet execution routing for `gitmap agm update-all ssh` and `gitmap agy ssh`, supporting `--except` node exclusion, real-time node alias and IP emission upon start and finish, and aggregated summary.
Acceptance Criteria:
1. `gitmap agm update-all ssh` and `gitmap agy ssh <cmd>` dispatch concurrently across registered nodes.
2. `--except` / `-e` excludes specified nodes from execution.
3. Terminal displays starting node alias and IP, streaming results asynchronously as they complete.
4. Final execution summary is printed with pass/fail counts.
Targeted Verification: go test -v ./cli/cmdssh/...
