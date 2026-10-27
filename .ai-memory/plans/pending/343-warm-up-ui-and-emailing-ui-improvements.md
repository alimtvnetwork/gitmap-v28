# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 343 Enterprise Spec)

- **Task Sequence:** 343
- **Task Slug:** `343-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 343 of the GitMap Warm-Up and Emailing UI Engine (Improvements 601–700). Sequence 343 advances warm-up deliverability and operational infrastructure to cryptographic ARC chains, cross-tenant reputation feedback circuits, interactive Bubble Tea terminal studios, zero-trust hardware-backed credential vaults, forensic header hop visualizers, distributed SQLite sharded replicators, and autonomous Antigravity swarm mission command:
- **Cryptographic ARC & Reputation Circuits (Improvements 601–630, Clusters 61–63):** RFC 8617 ARC chain validation and seal computation, intermediate forwarder integrity defense, multi-tenant reputation score decay, automated ISP bounce surge damping, and adaptive velocity curves.
- **Interactive Bubble Tea Warm-Up Studio & Seed Analytics (Improvements 631–660, Clusters 64–66):** Full-screen interactive Bubble Tea terminal cockpit, real-time dispatch pacing visualizers, IMAP/POP3 seed response ingestion, Bayesian conversation sentiment tracking, and cross-ISP inbox placement matrix.
- **Mailbox Credential Vault & Forensic Header Trace (Improvements 661–680, Clusters 67–68):** Zero-trust AES-256-GCM / ChaCha20-Poly1305 hardware-backed credential encryption, Casbin RBAC tenant isolation, cryptographic audit logging, RFC 5322 header decomposition, hop latency waterfall timelines, ASN/GeoIP enrichment, and intermediate TLS cipher verification.
- **Distributed SQLite Shards & AGY Swarm Commander (Improvements 681–700, Clusters 69–70):** Multi-node SQLite sharding via consistent hashing, streaming WAL frame replication, page-level cryptographic Merkle-tree conflict resolution, online non-blocking vacuum compaction, Antigravity swarm peer mesh, subtask leasing arbiters, token expenditure governors, and master mission commander HUD.

## Subtasks

- [ ] `01-cryptographic-arc-and-reputation-circuits.md` — Author specifications and implement Cryptographic ARC Chain Engine, Cross-Tenant Reputation Circuits, and Velocity Shaping (Improvements 601–630, Clusters 61–63).
- [ ] `02-interactive-bubbletea-and-inbox-analytics.md` — Author specifications and implement Interactive Bubble Tea Warmup Studio, Seed Response Dynamics, and Placement Matrix (Improvements 631–660, Clusters 64–66).
- [ ] `03-credential-vault-and-header-trace.md` — Author specifications and implement Zero-Trust Mailbox Credential Vault, Enclave Security, and Forensic Header Trace Visualizer (Improvements 661–680, Clusters 67–68).
- [ ] `04-sqlite-shards-and-agy-swarm-commander.md` — Author specifications and implement Distributed SQLite Sharded Replicator, Merkle-Tree Conflict Resolver, and Autonomous Antigravity Swarm Mission Commander (Improvements 681–700, Clusters 69–70).

## Architectural Clusters Breakdown

1. **Cryptographic ARC & Reputation Circuits (Clusters 61–63 / Improvements 601–630):**
   - RFC 8617 ARC chain verification, seal generation, and header parsing in `pkg/mail/arc/`.
   - Tenant reputation scoring, decay models, and automated circuit breakers in `pkg/mail/reputation/`.
   - Non-linear warmup velocity curves and dispatch pacing in `pkg/mail/circuit/`.

2. **Interactive Bubble Tea Studio & Inbox Analytics (Clusters 64–66 / Improvements 631–660):**
   - Interactive Bubble Tea TUI, Lipgloss styling tokens, and sparkline widgets in `pkg/cli/tui/`.
   - Seed inbox response monitoring, sentiment classification, and thread depth analysis in `pkg/mail/analytics/`.
   - Cross-provider inbox placement matrix and spam folder forensics in `pkg/mail/inbox/`.

3. **Mailbox Credential Vault & Header Trace (Clusters 67–68 / Improvements 661–680):**
   - Hardware-backed credential vault, key rotation, zero-swap secret caching, and Casbin RBAC in `pkg/mail/vault/`.
   - Deep RFC 5322 header parsing, hop latency waterfall visualizer, ASN/GeoIP enrichment, and TLS cipher auditing in `pkg/mail/forensics/` and `pkg/mail/trace/`.

4. **Distributed SQLite Shards & Swarm Commander (Clusters 69–70 / Improvements 681–700):**
   - Consistent hashing shard topology, streaming WAL replication, and epoch consensus in `pkg/db/sharding/`.
   - Page-level Merkle-tree generation and conflict resolution in `pkg/db/merkle/`.
   - Antigravity swarm gossip mesh, subtask lease arbitration, token governance, deadlock detection, and mission commander HUD in `pkg/agent/swarm/`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Zero exits in leaf packages. All errors returned directly via typed `*apperror.AppError`.
- **Function Limit:** Functions <= 15 lines max (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Positive booleans with `is` or `has` prefixes exclusively (`isVaultUnlocked`, `hasValidKey`, `isReachable`, `isSyncedWithMaster`).
- **Zero CDN Dependencies:** 100% self-contained embedded web assets via Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all documentation and configurations.
