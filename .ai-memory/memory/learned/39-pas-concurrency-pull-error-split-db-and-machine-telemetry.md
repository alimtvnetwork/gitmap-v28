# Learned Memory 39: PAS Worker Concurrency, Pull-Error Split-DB Subsystem & Machine Telemetry

Date: 2026-10-01
Related Specs:
- [02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)
- [.ai-memory/plans/pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../plans/pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)
- [.ai-memory/spec/commands/11-pas-worker-concurrency-pull-error-and-telemetry.md](../../spec/commands/11-pas-worker-concurrency-pull-error-and-telemetry.md)

## 1. Architectural Lessons & Principles

### 1.1 Machine Telemetry Clarity (`gitmap machine`)
- Bare command invocations that represent identity inspection must never dump generic help text.
- `gitmap machine` defaults to rendering the local node's `MachineIdentity` (Alias, Hostname, IP, OS, Architecture, CPU count, User, GitMap Version).
- `--json` emits clean, single-root JSON objects for automated pipelines and AI consumption.

### 1.2 SSH Fleet Startup Version Observability
- Distributed fleet operations (`gitmap pas`) must provide visual verification of the remote binary version alongside liveness (`Online`/`Offline`).
- Probing remote versions asynchronously during liveness checks prevents version divergence across nodes.
- Local execution always prints the host machine version without SSH network hops.

### 1.3 Concurrency Resolver & `paswh` Command
- Throttling pull operations across low-resource machines requires explicit worker pool controls (`--w <N>`) and hands (`--hand <Y>`).
- The `paswh` command provides an ergonomic shorthand (`gitmap paswh [N] [Y]`, defaulting to $1 \times 1$) for single-worker/single-hand execution to protect constrained remote fleet nodes.

### 1.4 Heartbeat Ticker for Long-Running Operations (>30s)
- Batch multi-repo pulls exceeding 30 seconds must emit non-destructive background heartbeat progress ticks (e.g. `[⏳ 35s elapsed] 14/48 completed...`) to prevent terminal freeze perceptions.

### 1.5 Pull-Error Split-DB Subsystem (`gitmap pull-error`)
- Isolating pull errors into dedicated Split-DB tables (`data/pull-errors/sql.db`) provides durable, queryable error history across sessions.
- Context-aware target resolution (current repo when inside `.git`, all repos when outside) enables instant developer and AI self-healing.
