# 05 — SSH Fleet Ping and Remote Execution

- **Subsystem:** Distributed Fleet Management & Remote SSH
- **Status:** Authoritative Reference

## 1. Dual-Stack Reachability Probing
- Concurrently runs OS machine ping (`ping -n` on Windows, `ping -c` on Linux) and TCP handshake on port 22.
- Accurately classifies nodes as `ONLINE`, `REACHABLE (TCP)`, `OFFLINE`, or `DEGRADED` when firewall rules drop ICMP packets.

## 2. Host Target Grammar
- Parses flexible connection targets: `user@host:port`, `host@user`, or configured machine aliases.
- Defaults to user `root` (or current login user) and port `22`.

## 3. Streaming Remote Execution
- Streams stdout and stderr line-by-line across remote SSH sessions with ANSI color highlighting.
- Propagates remote non-zero exit codes accurately to the host terminal.
