# Subtask 04: Split-DB WAL Tuning, In-Memory Caches & Antigravity Event Bus

## Status: READY
- **Subtask Slug:** `04-splitdb-tuning-and-agy-event-bus`
- **Parent Plan:** `.ai-memory/plans/pending/339-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/04-improvements-261-to-300-failover-tokens-and-telemetry.md`
- **Owner:** Spec Author Subagent 3 / Database & Agent Platform Engineer
- **Target Subsystem:** `pkg/db/splitdb/`, `pkg/db/audit/`, `pkg/agent/bus/`, `pkg/agent/anomaly/`
- **Exclusive Owned Files:**
  - `pkg/db/splitdb/manager.go`
  - `pkg/db/splitdb/pragmas.go`
  - `pkg/db/splitdb/wal_governor.go`
  - `pkg/db/audit/merkle_tree.go`
  - `pkg/agent/bus/event_broker.go`
  - `pkg/agent/anomaly/triage_engine.go`
  - `web/static/js/audio_synthesizer.js`
  - `cli/cmddb/cmd_doctor.go`

---

## 1. Context & Objectives

This subtask delivers enterprise Split-DB SQLite engine optimizations and real-time autonomous agent event infrastructure:
1. **Dual Split-DB Storage WAL Tuning (Cluster 29, Items 281–290):** Physical isolation between `warmup.db` and `mail_audit.db`, optimized WAL pragmas, automated WAL checkpoint governor, prepared statement caching, in-memory Bloom filter suppression lookup, vacuum log compactor, and SHA-256 Merkle tree audit ledger.
2. **Antigravity Autonomous Agent Event Bus (Cluster 30, Items 291–300):** Real-time pub/sub WebSocket message broker, signed webhook dispatcher for AGY agents, statistical anomaly triage alerts, procedural Web Audio API deliverability synthesizer, terminal bell routers, and headless daemon supervisor.
3. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isWalModeActive`, `hasOptimalPragmas`, `isWebSocketConnected`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Split-DB Architecture & Performance Pragmas
- Configure dual database connections in `pkg/db/splitdb/manager.go` applying `PRAGMA journal_mode=WAL`, `synchronous=NORMAL`, and `busy_timeout=5000`.
- Construct background WAL governor in `pkg/db/splitdb/wal_governor.go` auto-truncating uncheckpointed WAL frames.

### Step 2: Caching, Bloom Filter & Merkle Tree Audit Ledger
- Author thread-safe prepared statement cache in `pkg/db/splitdb/stmt_cache.go`.
- Implement fast in-memory Bloom filter for suppressed recipient checks in `pkg/mail/filter/bloom_filter.go`.
- Build cryptographic SHA-256 Merkle tree calculation over audit logs in `pkg/db/audit/merkle_tree.go`.

### Step 3: WebSocket Event Bus, Webhooks & Anomaly Alerts
- Build pub/sub WebSocket broker in `pkg/agent/bus/event_broker.go` streaming typed delivery events.
- Implement Z-score anomaly detector in `pkg/agent/anomaly/triage_engine.go` triggering automated throttling.
- Implement zero-dependency Web Audio procedural chimes in `web/static/js/audio_synthesizer.js`.
- Build diagnostic CLI tool in `cli/cmddb/cmd_doctor.go`.

---

## 3. Verification & Acceptance Criteria
- [ ] SQLite connections consistently operate in WAL mode with busy timeout handling concurrent reads/writes.
- [ ] Merkle tree verification successfully detects single-bit data modifications in audit transaction logs.
- [ ] WebSocket message broker sustains 5,000 events/sec broadcast with zero dropped messages.
- [ ] Web Audio synthesizer produces clean auditory feedback without downloading external sound assets.
