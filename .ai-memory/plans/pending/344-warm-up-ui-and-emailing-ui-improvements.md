# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 344 Enterprise Spec)

- **Task Sequence:** 344
- **Task Slug:** `344-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 344 of the GitMap Warm-Up and Emailing UI Engine (Improvements 701–800). Sequence 344 pushes warm-up deliverability, client emulation, forensic reputation tracking, cryptographic identity, and swarm mission governance to industry-leading standards:
- **Fingerprint Morphing & Rate Limit Mesh (Improvements 701–730, Clusters 71–73):** Dynamic TLS client Hello fingerprint morphing (JA3/JA4), TCP MSS and window scaling virtualization, behavioral mouse/scroll humanization in webmail emulators, cross-ISP distributed rate-limiting token bucket mesh, and adaptive throttling.
- **Interactive TUI Cockpit & Mailbox Telemetry (Improvements 731–760, Clusters 74–76):** Full-screen interactive Bubble Tea terminal cockpit, ANSI sparkline spark charts, IMAP/POP3 seed response ingestion, Bayesian conversation sentiment tracking, cross-ISP inbox placement matrix, and real-time dispatch pacing visualizers.
- **Honeypot Interaction Simulator & Cryptographic DKIM Rotation Sentinel (Improvements 761–780, Clusters 77–78):** Autonomous synthetic honeypot traps detection and conversational simulation, multi-turn AI email thread weaving, cryptographic DKIM key rotation sentinel, automated DNS TXT/CNAME propagation verification HUD, and public key revocation circuits.
- **SQLite Split-DB Multi-Cluster Replicator & Swarm Telemetry Bridge (Improvements 781–800, Clusters 79–80):** Multi-cluster SQLite split-db synchronization, WAL streaming replication, page Merkle-tree conflict resolution, online non-blocking vacuum compaction, Antigravity multi-agent swarm orchestrator, subtask leasing arbiters, token expenditure governors, and live warmup telemetry bridge HUD.

## Subtasks

- [ ] `01-fingerprint-morphing-and-rate-limit-mesh.md` — Author specifications and implement TLS JA3/JA4 Fingerprint Morphing, Emulated Behavioral Dynamics, and Distributed Rate Limit Mesh (Improvements 701–730, Clusters 71–73).
- [ ] `02-interactive-tui-and-mailbox-telemetry.md` — Author specifications and implement Interactive Bubble Tea TUI Cockpit, IMAP/POP3 Seed Ingestion, Sentiment Classification, and Placement Analytics (Improvements 731–760, Clusters 74–76).
- [ ] `03-honeypot-simulator-and-dkim-rotation.md` — Author specifications and implement Autonomous Honeypot Interaction Simulator, Multi-Turn Conversational Thread Weaver, Cryptographic DKIM Rotation Sentinel, and DNS Propagation HUD (Improvements 761–780, Clusters 77–78).
- [ ] `04-sqlite-shards-and-swarm-telemetry.md` — Author specifications and implement SQLite Split-DB Multi-Cluster Replicator, Cross-Node Ledger Synchronization, Antigravity Multi-Agent Swarm Orchestrator, and Live Warmup Telemetry Bridge (Improvements 781–800, Clusters 79–80).

## Architectural Clusters Breakdown

1. **Fingerprint Morphing & Rate Limit Mesh (Clusters 71–73 / Improvements 701–730):**
   - JA3/JA4 client hello randomized cipher suite permutation in `pkg/mail/fingerprint/`.
   - Webmail browser session humanization and mouse entropy simulation in `pkg/mail/emulation/`.
   - Distributed cross-node token bucket rate limiting in `pkg/mail/ratelimit/`.

2. **Interactive TUI Cockpit & Mailbox Telemetry (Clusters 74–76 / Improvements 731–760):**
   - Terminal Bubble Tea dashboard, Lipgloss tokens, and ANSI sparklines in `pkg/cli/tui/`.
   - Automated seed inbox synchronization and reply classification in `pkg/mail/telemetry/`.
   - Cross-provider inbox vs spam folder delivery matrix in `pkg/mail/placement/`.

3. **Honeypot Simulator & Cryptographic DKIM Rotation (Clusters 77–78 / Improvements 761–780):**
   - Synthetic honeypot tripwire detection and simulated natural reply generation in `pkg/mail/honeypot/`.
   - Multi-turn AI conversational thread weaving and contextual relevance in `pkg/mail/weaver/`.
   - Automated DKIM 2048-bit RSA/Ed25519 key generation and dual-selector staging in `pkg/mail/dkim/`.
   - Real-time DNS TXT/CNAME propagation verification HUD across global nameservers in `pkg/mail/dnscheck/`.

4. **SQLite Split-DB Multi-Cluster & Swarm Telemetry Bridge (Clusters 79–80 / Improvements 781–800):**
   - Multi-node SQLite sharding via consistent hashing and WAL streaming replication in `pkg/db/cluster/`.
   - Cryptographic Merkle-tree page-level reconciliation and last-write-wins (LWW) resolution in `pkg/db/merkle/`.
   - Antigravity multi-agent swarm coordination, subtask lease arbitration, and token budget governance in `pkg/agent/swarm/`.
   - Live warmup telemetry bridge streaming metrics, error bursts, and health indices to the mission HUD in `pkg/agent/telemetry/`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Zero exits in leaf packages. All errors returned directly via typed `*apperror.AppError`.
- **Function Limit:** Functions <= 15 lines max (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Positive booleans with `is` or `has` prefixes exclusively (`isHoneypotDetected`, `hasValidDkimKey`, `isDnsPropagated`, `isLeaderNode`, `isSyncedWithMaster`, `isAgentAlive`).
- **Zero CDN Dependencies:** 100% self-contained embedded web assets via Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all documentation and configurations.
