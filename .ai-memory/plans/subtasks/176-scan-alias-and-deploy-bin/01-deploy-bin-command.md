# Subtask 01: Smart GitMap Binary Deploy Command (`gitmap ssh deploy-bin`)
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/176-scan-alias-migration-internal-errors-db-and-fleet-inventory-aggregation.md](../../../02-spec/21-app/176-scan-alias-migration-internal-errors-db-and-fleet-inventory-aggregation.md)
Target Files: cli/cmdssh/ssh_deploy_bin.go, cli/cmdssh/ssh_deploy_bin_help.go, cli/cmdssh/ssh.go, cli/helptext/deploy-bin.md
Status: Completed
Completed At: 2026-09-27T20:18:00+08:00

Action: Implement `gitmap ssh deploy-bin [target] [--file <path>]` (aliases: `push-bin`, `sync-bin`) to autonomously deploy the local gitmap binary to target nodes or all fleet nodes using SSH client connection, resolve remote installation path based on remote OS, upload binary, verify permissions, and test remote version.

Acceptance Criteria:
- [x] `gitmap ssh deploy-bin w1` transfers local binary to W1 and reports remote version.
- [x] `gitmap ssh deploy-bin all` broadcasts to all online fleet nodes concurrently.
- [x] `gitmap ssh deploy-bin --help` displays terminal UI help text.
- [x] Targeted Verification: python 03-ai-scripts/05-guideline-autofixer.py cli/cmdssh/ssh_deploy_bin.go (PASS)
