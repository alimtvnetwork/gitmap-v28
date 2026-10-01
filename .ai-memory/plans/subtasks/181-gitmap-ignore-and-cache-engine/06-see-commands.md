---
status: PENDING
---

# Subtask: See Commands

## Objective
Implement read-only queries via the `gitmap see` (alias `gitmap c`) command for better visibility into the cache, ssh nodes, and pending tasks.

## Requirements
- `gitmap see all fast` and `gitmap see all fast <group>`: Display fast paths.
- `gitmap [c] see all repo-cache` and `gitmap [c] see all repo-cache <group>`: Query `repodb.db` to show cached repository data.
- `gitmap [c] see ssh-node all`: Query `nodes.db` and list all SSH nodes registered globally or locally.
- `gitmap [c] see ssh-node <id>`: Query `nodes.db` and view detailed information for a specific node (id, os, type, ip, identityfile, alive status, load/cpu/ram, latency).
- `gitmap [c] see pending-tasks all <group>`: List all pending-tasks for a group in ai-memory.

## Next Steps
1. Add new subcommands for `see` in `cli/cmdsee/` for each query type.
2. Connect `see all repo-cache` to `repodb` cache queries.
3. Connect `see ssh-node` commands to `nodes.db` to fetch SSH node history and metrics.
4. Verify the output formats meet the UX standards outlined in the spec.
