# Specification: Antigravity Running Prompts Multi-Node SSH Delegation

Spec Reference: [02-spec/21-app/167-agy-running-prompts-backup-restore-and-green-automation/01-overview.md](01-overview.md)

## 1. Overview & Objectives
Enable cross-node execution and aggregation for Antigravity running prompts backup, restore, batch listing, and prompt inspection across all nodes in the cluster SSH fleet.

## 2. CLI Command Syntax & Flags
```bash
# 1. Inspect running/queued prompts across cluster nodes
gitmap agy running-prompts ls --ssh [--json] [--full] [--limit 8] [--wc 100]

# 2. Snapshot running prompts across cluster nodes
gitmap agy running-prompts backup --ssh [--json]
gitmap agy backup-running-prompts --ssh [--json]

# 3. List backup batches across cluster nodes
gitmap agy running-prompts backup ls --ssh [--json]
gitmap agy backup-running-prompts ls --ssh [--json]

# 4. Restore running prompts across cluster nodes
gitmap agy running-prompts restore --ssh [--keep] [--json]
gitmap agy restore-running-prompts --ssh [--keep] [--json]
```

## 3. Data Contracts
- For SSH aggregation, records carry:
  - `node`: Alias of the cluster SSH connection (or `"local"` for the local node).
  - `host`: IP address or hostname of the cluster node.
  - `status`: Execution state (`OK`, `OFFLINE`, `ERROR`).
- In JSON mode, results are serialized as a flat array of records with `node` and `host` populated.
- In terminal mode, tables display a `NODE` column indicating origin.

## 4. Verification Gates
1. Flags `--ssh` registered on `running-prompts ls`, `running-prompts backup`, `backup-running-prompts`, `running-prompts restore`, `restore-running-prompts`, and `backup ls`.
2. Local execution runs first; remote nodes are queried using `crypto.RunCommand` without package import cycles.
3. Node connection failures report `OFFLINE` or error gracefully without crashing.
4. Linters pass with 0 violations (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`, `golangci-lint run --issues-exit-code=1 ./...`).
