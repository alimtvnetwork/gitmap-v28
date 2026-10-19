# GitMap v6.578.0

## What's Changed in v6.578.0

- **GitMap Background Daemon (`gitmap daemon` / `gitmap gitmap-daemon`)**:
  - Renamed `lemon` background service to `gitmap-daemon` with full backward compatibility for `lemon` and `lm` CLI aliases.
  - Added configurable polling loop (default 8s) with a non-negotiable `>= 5s` safety guard.
  - Implemented application settings toggle (`daemon.enabled`): CLI shortcuts `gitmap daemon enable` and `gitmap daemon disable`, with settings command `gitmap settings set daemon.enabled <true|false>`.
  - Added OS background service manager integration for Windows (`sc.exe`), Linux (`systemd`), and macOS (`launchd`).
- **Multi-Format Email Command Ingestion & Immediate ACK**:
  - Support for Format A (plain slash commands), Format B (JSON command envelopes), and Format C (fenced markdown and YAML headers).
  - Immediate pre-execution receipt acknowledgment (`ACK Received: Command '<cmd>' received at <timestamp> (Request ID: <id>). Executing...`).
  - Command whitelist security governance with parameter sanitization, shell injection defense, and execution outcome reporting.
- **Ent Go ORM Feasibility Analysis & Type Safety Hardening**:
  - Comprehensive feasibility study documenting empirical benchmarks (point query latency 12µs vs 24µs, memory allocations 8 vs 32 allocs) and architectural blockers for SQLite Split-DB architecture.
  - Verified compile-time type-safe parameterized SQL across all database packages with zero dynamic SQL concatenation.
- **Codebase Hygiene & File Size Compliance**:
  - Decomposed oversized settings and command files (`cmdagy`, `cmdlemon`, `cmdnodes`, `cmdsupabase`) ensuring all modified files comply strictly with the repository 300-line cap.
  - Flattened all nested `if` statements to maintain 0 violations across AST and boolean linters.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.578.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.578.0/install.sh | sh
```
