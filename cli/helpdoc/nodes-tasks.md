# gitmap nodes tasks

Inspect active, pending, and heavy-duty AI agent tasks enqueued in `agent-task.db` across all cluster nodes with real-time CPU and RAM telemetry.

## Synopsis

```bash
gitmap nodes tasks [node-alias] [flags]
gitmap nodes queue [node-alias] [flags]
gitmap nodes jobs [node-alias] [flags]
```

## Description

Interrogates cluster nodes in parallel, reading task state from run-scoped Tier 2 SQLite task databases (`agent-task.db`), and aggregates workload distribution across the fleet.

A task is classified as **Heavy-Duty** (`isHeavyDuty: true`) if:
1. `ParentTask.TotalStepsBudget >= 100`, OR
2. Task slug or notes contain heavy-load keywords (`heavy`, `full-audit`, `deep-scan`, `migration`, `build-all`), OR
3. Number of owned files in `OwnedFilesJson` exceeds 8.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` (`-j`) | `false` | Output aggregated task queues and node health as JSON |
| `--heavy-only` | `false` | Filter only tasks tagged with `isHeavyDuty: true` |
| `--open-only` (`-ro`) | `false` | Probe strictly reachable online nodes, avoiding TCP connection timeouts on offline machines |

## Role-Based Delegation

Role-based scheduling is currently marked with `?` as documented in `.ai-memory/ambiguous-questions/01-new-ambiguity/04-role-based-task-scheduling.md`.

## Examples

```bash
# Display aggregated task queues across all fleet nodes
gitmap nodes tasks

# Query task queue on specific remote node
gitmap nodes tasks vm-02

# Filter only heavy-duty tasks across reachable nodes
gitmap nodes tasks --heavy-only --open-only

# Output raw JSON telemetry for orchestrator agents
gitmap nodes tasks --json
```

## See Also

- [nodes](nodes.md) — Unified fleet infrastructure command and node status
- [nodes-running-only](nodes-running-only.md) — Discover and execute commands strictly on online fleet nodes
- [nodes-resources](nodes.md) — Probe CPU and RAM telemetry across fleet
