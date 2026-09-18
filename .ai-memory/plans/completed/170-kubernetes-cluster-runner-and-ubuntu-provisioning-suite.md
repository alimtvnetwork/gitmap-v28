# Plan 170: Kubernetes Cluster Runner & Ubuntu Provisioning Suite

## Executive Summary
Integrate the multi-node execution and cluster management from `kubernetes-training/05-server-cmds` and the node provisioning automation from `kubernetes-training/02-ubuntu-install` into GitMap. Passwords from topology JSON files (`01-config.json`) will be encrypted at rest in the SQLite database (`ssh_hosts` table) using SSH RSA-OAEP with SHA-256 via `~/.ssh/id_rsa`, with seamless in-memory decryption for remote command execution, sudo elevation (`echo $pass | sudo -S`), and remote script staging.

## Architecture Guidelines
- TOTAL BAN on running `go test`, `go build`, or `06-cicd-local-runner.py` during execution.
- Function length <= 15 lines (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal *apperror.AppError returns.
- Strict Unix LF line endings.
- Strict disjoint file assignments across subagents.

## Subtasks
1. `01-sqlite-cluster-roles-and-target-resolver.md`
2. `02-cluster-config-json-importer.md`
3. `03-cluster-remote-exec-and-sudo-elevation.md`
4. `04-ubuntu-node-provisioning-recipes-and-cli-root.md`
