---
name: cluster-ssh-delegation
description: >-
  Autonomously manage SSH node discovery, config templating, TLS dial timeouts, and multi-node command delegation across Gitmap clusters.
---

# Cluster & SSH Delegation Skill

Autonomously implement, maintain, and audit SSH lifecycle and cluster command delegation in `gitmap/cluster/` and `gitmap/cmd/` adhering to `.lovable/spec/commands/07-cluster-command-delegation.md`, `spec/01-app/`, and `.lovable/cicd-issues/05-cluster-tls-dial-timeout-and-test-env-race.md`.

## Core Checkpoints & Invariants

1. **SSH Key Lifecycle & Config Generation:**
   - Manage multi-host SSH profiles and key generation (`ed25519` / `rsa`).
   - SSH config templating into `~/.ssh/config` must utilize idempotent marker blocks (`# BEGIN GITMAP MANAGED HOSTS` ... `# END GITMAP MANAGED HOSTS`).
   - Never mutate user-defined host entries outside the Gitmap managed marker blocks.

2. **Transport Preservation:**
   - Reclone and clone operations must honor the identified transport protocol (SSH vs HTTPS).
   - If a repository was originally cloned via SSH (`git@...`), reclone operations must strictly preserve the SSH remote URL and transport parameters.

3. **Cluster Command Delegation:**
   - Execute commands concurrently across cluster nodes with worker pools.
   - Enforce explicit TLS dial timeouts and keep-alive intervals to avoid deadlocks.
   - Guard against test environment concurrency races: ensure isolated temp directories and mock cluster targets in CI unit tests.

4. **Error & Exit Reporting:**
   - Use `cliexit.Reportf` or `cliexit.Fail` for CLI errors; wrap networking failures in `apperror.AppError`.
   - Never print bare errors directly to `os.Stderr`.
