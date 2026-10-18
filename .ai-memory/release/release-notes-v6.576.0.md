# GitMap v6.576.0

## What's Changed in v6.576.0

- **ENT Goal Feasibility Analysis & Performance Report:**
  - Evaluated query latency (+25% to +60%), heap allocations (+300% to +400%), cold CLI startup latency (+15ms to +40ms), binary size bloat (+15MB to +35MB), and Split-DB SQLite architectural impedance.
  - Documented complete technical findings in `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/ent-feasibility.md` and app spec `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/spec.md`.
- **Compile-Time Type Safety & SQL Hardening:**
  - Hardened SQLite compiler identifier quoting (`cli/dbengine/sqlite/compiler.go`) to escape embedded quotes.
  - Replaced dynamic SQL formatting with static typed SQL statements in `cli/store/inject_idempotency.go` and `cli/constants/constants_inject_idempotency.go`.
- **Lemon CLI Daemon & Email Command Integration Subsystem (`cli/cmdlemon/`):**
  - Multi-format email parser (`cli/cmdlemon/email_parser.go`) supporting Format A (Plain text / slash commands), Format B (JSON envelopes), and Format C (Markdown / YAML blocks).
  - Immediate receipt acknowledgment replier (`cli/cmdlemon/email_ack.go`) transmitting `ACK Received` before command execution.
  - Sandboxed command executor (`cli/cmdlemon/command_executor.go`) with whitelist validation (`cli/cmdlemon/command_whitelist.go`) for Antigravity Manager / Tools commands (`gitmap agm`, `gitmap agy`, `gitmap status`, `gitmap pull`, `gitmap nodes`).
  - Background daemon polling runner (`cli/cmdlemon/lemon_daemon.go`) with non-negotiable $\ge 5$s safety interval clamp (`cli/cmdlemon/interval_guard.go`) and OS signal handling.
  - Cross-platform service hooks (`cli/cmdlemon/service_linux.go`, `cli/cmdlemon/service_darwin.go`, `cli/cmdlemon/service_windows.go`).
  - Cobra CLI integration (`cli/cmdlemon/lemon_cmd.go`, `cli/cmd/roottooling.go`).
- **Comprehensive End-to-End Tests:**
  - Authored test suite `cli/tests/e2e/lemon_email_e2e_test.go` covering multi-format parsing, interval clamping, immediate ACK delivery, command execution, and compiler quote escaping.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.576.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.576.0/install.sh | sh
```
