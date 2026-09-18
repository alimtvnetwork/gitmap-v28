# Subtask 04: SSH vs Cluster vs SC Comparison Table and UI Help Parity

## Status: COMPLETED

## Summary of Accomplishments
1. Implemented comparative architecture table in `cli/cmdssh/ssh_compare_table.go`:
   - Compares `gitmap ssh` (direct node-level ad-hoc management, remote exec, scan, install/update, AGY/code).
   - Compares `gitmap cluster` (multi-node cluster bootstrapping, coordination, k8s recipes, role targeting).
   - Compares `gitmap sc` (server-client daemon topology, sync, continuous monitoring, broadcast fan-out).
2. Renders table using `cli/termtable.PrintTable` with clean coloring and column alignment, plus workflow guidance.
3. Added `gitmap ssh compare` (alias `matrix`) routed through `cli/cmdssh/ssh.go`.
4. Enriched help markdown in `cli/helptext/ssh.md`, `cli/helptext/cluster.md`, and `cli/helptext/sc.md` with comparative explanations, workflow differences, and command syntax examples.
5. Added unit tests in `cli/cmdssh/ssh_compare_table_test.go`.
