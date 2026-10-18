# ENT Goal Feasibility and Email Command Integration (`02-spec/21-app/ent-goal-feasibility-and-email-command-integration`)

Authoritative specification suite for the ENT goal feasibility analysis, type-safe database architecture, and cross-platform Lemon background daemon / email command execution engine.

---

## Specification Documents

1. [Architectural Specification (`spec.md`)](spec.md)
   - Comprehensive systems design for the Lemon background service, multi-format email parser, immediate receipt acknowledgment (`ACK Received`), sandboxed command execution, and cross-platform OS service managers (`systemd`, `launchd`, `sc.exe`).
2. [Ent Go Feasibility & Type-Safe Database Evaluation (`ent-feasibility.md`)](ent-feasibility.md)
   - Benchmark-grounded architectural analysis addressing Ent Go performance implications, Split-DB multi-SQLite concurrency deadlocks, and dynamic SQL elimination via parameterized constant queries.

---

## Key Invariants & Architectural Contracts

- **Ent Go Verdict:** Disqualified due to 25–60% query latency degradation, 300–400% memory allocation explosion, 30–50ms cold-start latency penalty, and fundamental SQLite connection pool deadlocks. Replaced by compile-time type safety with zero dynamic SQL.
- **Daemon Interval Safety Guard:** Default 8 seconds (`8s`). Strictly clamped to minimum 5 seconds (`5s`) if configured lower.
- **Immediate ACK Protocol:** Dispatches receipt acknowledgment before command execution begins.
- **Sandboxed Command Execution:** Strict whitelist (`status`, `pull`, `pe`, `version`, `nodes`, `agm`, `agy`, `space`, `scan`) with shell injection protection and timeout guards.
