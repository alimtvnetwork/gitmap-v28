# Subtask 01: SSH Node Liveness Probing, Availability Caching, and Default Key Discovery

## Status: COMPLETED

## Summary of Accomplishments
1. Implemented in-memory TTL liveness cache in `cli/cmdssh/ssh_liveness.go` (89 lines) with thread-safe `sync.RWMutex`, 45-second TTL, and `CheckConnLiveness` / `CheckNodeLiveness` / `InvalidateLivenessCache`.
2. Probed TCP reachability before attempting command execution in `cli/cmdssh/sshexec.go`: unreachable/offline nodes are skipped immediately with `[alias|ip] OFFLINE (skipped: reason)`.
3. Added `gitmap ssh scan` in `cli/cmdssh/ssh_scan_cmd.go` and `cli/cmdssh/ssh_scan_render.go` to probe fleet nodes in parallel and render a styled table with `termtable.PrintTable`.
4. Implemented automatic user key discovery in `cli/cmdssh/ssh_keys.go`: probes `~/.ssh/id_ed25519`, `~/.ssh/id_rsa`, and `~/.ssh/id_ecdsa` when no password or key was explicitly supplied during join.
5. Formatted structured guidance (`formatMissingAuthAdvice`) explaining how to authenticate (`gitmap ssh join-auth` or `--key` / `add-with-pass`).
6. Added comprehensive unit tests in `cli/cmdssh/ssh_liveness_test.go` and `cli/cmdssh/ssh_keys_test.go`.
