---
name: ssh-join-and-recall
description: Autonomously implement and verify SSH machine joining, host alias registration, mistakes recovery, rich examples, and host recall across GitMap.
---

# SSH Join, Machine Registration & Recall Suite

## Objective
Autonomously implement first-class SSH machine joining (`gitmap ssh join <target>`, `gitmap sj <target>`), robust mistake recovery with concrete fix examples, host alias management in SQLite (`ssh_hosts`, `ssh_history`), and seamless host recall (`gitmap ssh <alias>`, `gitmap ssh join ls`) adhering to GitMap coding guidelines and zero OpenSSH hostname collision.

## Key Invariants
1. **Never Forward 'join' to OpenSSH**: `gitmap ssh join` must never execute `ssh join ...`. Join is an internal GitMap subcommand that registers machines and manages SSH host aliases.
2. **First-Class SSH Machine Enrollment**:
   - `gitmap ssh join <ip|user@ip> [alias]` registers the target host, username, and IP in SQLite.
   - Supports `--user`, `--name` / `--alias`, `--port`, and `--auth` (public key push) flags.
3. **Graceful Mistake Recovery & Concrete Examples**:
   - When arguments are missing or invalid, output formatted guidance with exact copy-pasteable commands to fix the mistake.
   - When connecting to an unknown alias, suggest matching hosts or explain how to join the machine.
4. **Rich Help Parity**:
   - `gitmap ssh join help` and `gitmap sj help` display full examples covering joining, recalling, authentication, and execution.
   - `gitmap ssh help` accurately documents all SSH subcommands including `join`, `login`, `exec`, and `as`.
5. **Strict Compliance**:
   - Functions <= 8-15 lines.
   - Affirmative booleans (`is*`, `has*`).
   - Universal `*apperror.AppError` return wrapping.
   - Strict Unix LF line endings.
   - Total ban on running `go test`, `go build`, or local runner during execution turns.
