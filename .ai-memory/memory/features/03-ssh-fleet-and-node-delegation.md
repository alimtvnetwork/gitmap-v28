# Feature Catalog 03: SSH Fleet and Node Delegation

- **Domain:** Multi-Node Fleet Management, Probing, and Remote Execution
- **Status:** Authoritative Capability Catalog

## 1. Dual-Stack Machine Reachability (`gitmap nodes ping` / `gitmap ping`)
- Dual-stack ICMP and TCP probing accurately detects online nodes behind restrictive firewalls.
- Formatted ANSI box table rendering displays node alias, host, latency, and status.

## 2. Credential Security & RSA Vault
- Masked terminal input captures passwords with bullet masking (`*`) and zero echo.
- RSA-OAEP local credential vault (`credentials.vault`) stores node secrets encrypted with PBKDF2 salt.

## 3. Remote Command Streaming & Node Sync
- Streams remote stdout and stderr line-by-line with ANSI color highlighting.
- Remote clone with `--except-self` guard prevents local recursive cloning.
