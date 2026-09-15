---
name: cluster-sc-node-add-and-help-parity
description: >-
  Autonomously implement and audit cluster and servers-clients (sc) commands, help text parity, node add and join routing, command aliases, and descriptive examples across GitMap.
---

# Cluster and SC (Servers-Clients) Help & Node Add Parity

Autonomously implement, maintain, and audit `cluster`, `servers-clients` (`servers-client`, `sc`), `node add`, `join`, and help text parity across GitMap.

## Core Invariants

1. **Root Dispatch & Alias Parity**:
   - `gitmap servers-clients`, `gitmap servers-client`, `gitmap sc`, `gitmap clients` must all route correctly.
   - Calling root commands with `help`, `--help`, `-h`, or zero arguments must display rich, formatted help rather than "unknown sub-command token: help" or unhelpful messages.
2. **Cluster Node Add & Join Routing**:
   - `gitmap cluster node add <target> [alias] [flags]` and `gitmap cluster add <target> [alias] [flags]` must enroll nodes into the cluster / SSH host registry.
   - `gitmap cluster join` must route to `gitmap ssh-join` or cluster enrollment with clear help and copy-pasteable examples.
3. **Descriptive Leaf & Root Help Text**:
   - Every subcommand and sub-subcommand must support `-h`, `--help`, and `help`.
   - Subcommand help (e.g. `gitmap cluster nodes help`, `gitmap cluster exec help`) must show specific help for that command, not fallback to parent help blindly without context.
   - Provide concrete, copy-pasteable examples for every operation.
4. **Coding Guidelines**:
   - Functions <= 15 lines (target <= 8 lines).
   - Affirmative booleans (`is*`, `has*`).
   - Single return types with universal `*apperror.AppError` returns.
   - Unix LF line endings and UTF-8 encoding.
