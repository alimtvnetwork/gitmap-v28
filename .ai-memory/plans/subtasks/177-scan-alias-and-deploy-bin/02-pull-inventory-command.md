# Subtask 02: Remote File Fetching & Fleet Inventory Automation (`gitmap ssh pull-inventory`)
Traceability ID: Task-03
Spec Reference: [02-spec/21-app/177-scan-alias-migration-internal-errors-db-and-fleet-inventory-aggregation.md](../../../02-spec/21-app/177-scan-alias-migration-internal-errors-db-and-fleet-inventory-aggregation.md)
Target Files: cli/cmdssh/ssh_pull_inventory.go, cli/cmdssh/ssh.go
Status: Completed
Completed At: 2026-09-27T20:18:00+08:00

Action: Implement `gitmap ssh pull-inventory [target]` (alias: `fetch-inventory`, `sync-inventory`) to execute remote scans on W1, W2, W3 via SSH, read the resulting `gitmap.json` and node configs, and save them directly into `./repo-secrets\04-w1-machine\`, `05-w2-machine\`, and `06-w3-machine\`.

Acceptance Criteria:
- [x] Executes remote scan via SSH without stalling.
- [x] Fetches remote `gitmap.json`, `gitmap-ssh-nodes.json`, `gitmap-ssh.json`, `ooshutup10.cfg`.
- [x] Populates `./repo-secrets\04-w1-machine\`, `05-w2-machine\`, and `06-w3-machine\`.
- [x] Committed and pushed to `alimtvnetwork/repo-secrets` on GitHub.
- [x] Targeted Verification: python 03-ai-scripts/05-guideline-autofixer.py cli/cmdssh/ssh_pull_inventory.go (PASS)
