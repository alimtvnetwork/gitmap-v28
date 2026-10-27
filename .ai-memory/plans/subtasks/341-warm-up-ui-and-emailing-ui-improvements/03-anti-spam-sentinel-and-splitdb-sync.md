# Subtask 03: Autonomous Anti-Spam Sentinel, Template AST Diffing & Split-DB Sync

## Status: READY
- **Subtask Slug:** `03-anti-spam-sentinel-and-splitdb-sync`
- **Parent Plan:** `.ai-memory/plans/pending/341-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/04-improvements-461-to-500-security-audit-and-agy-bridge.md`
- **Target Range:** Improvements 461–490 (Clusters 47–49)
- **Target Subsystem:** `pkg/mail/trap/`, `pkg/mail/dnsbl/`, `pkg/mail/template/`, `pkg/db/splitdb/`
- **Exclusive Owned Files:**
  - `pkg/mail/trap/bayesian_classifier.go`
  - `pkg/mail/trap/lexical_scanner.go`
  - `pkg/mail/dnsbl/sentinel.go`
  - `pkg/mail/template/ast_parser.go`
  - `pkg/mail/template/ast_differ.go`
  - `pkg/mail/template/spintax_verifier.go`
  - `pkg/db/splitdb/wal_replicator.go`
  - `pkg/db/splitdb/vacuum_scheduler.go`
  - `pkg/db/splitdb/fragmentation_analyzer.go`
  - `cli/cmdwarmup/sentinel_cmd.go`

---

## 1. Context & Objectives

This subtask delivers the anti-spam trap sentinel, polyglot template AST diffing engine, and Split-DB distributed replication layer:
1. **Autonomous Anti-Spam Trap Detection & Honeypot Avoidance (Cluster 47, Items 461–470):** Bayesian probability scoring, pristine honeypot lexical scanners, DNSBL blocklist sentinel, MX clustering, and automated quarantine isolation.
2. **Polyglot Campaign Orchestration & Template AST Diffing (Cluster 48, Items 471–480):** Multi-engine template AST parsing, semantic AST diffing, Shannon text entropy verification, i18n translation validation, and Thompson-sampling bandit schedules.
3. **Split-DB Cross-Replica Sync HUD & SQLite Health (Cluster 49, Items 481–490):** Streaming SQLite WAL replication via vector clocks, non-blocking incremental auto-vacuum scheduler, B-tree page fragmentation analyzer, and prepared statement LRU pools.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isTrapSuspected`, `isSynchronized`, `isSemanticallyEqual`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Anti-Spam Sentinel & Trap Avoidance
- Build Bayesian trap classifier and lexical pattern scanner in `pkg/mail/trap/`.
- Implement concurrent DNSBL multi-provider query engine in `pkg/mail/dnsbl/sentinel.go`.
- Construct dynamic quarantine manager and threat HUD in `pkg/mail/trap/`.

### Step 2: Polyglot Template AST Parser & Semantic Diffing
- Implement multi-engine template tokenizer and AST builder in `pkg/mail/template/ast_parser.go`.
- Build structural AST diffing inspector in `pkg/mail/template/ast_differ.go`.
- Verify spintax combinatorial entropy and i18n coherence in `pkg/mail/template/`.

### Step 3: Split-DB Distributed WAL Replication & Vacuum Scheduler
- Implement streaming WAL frame replicator with vector clocks in `pkg/db/splitdb/wal_replicator.go`.
- Build non-blocking incremental auto-vacuum scheduler in `pkg/db/splitdb/vacuum_scheduler.go`.
- Implement SQLite B-tree page fragmentation inspector and statement pool in `pkg/db/splitdb/`.
- Wire CLI commands into `cli/cmdwarmup/sentinel_cmd.go`.
