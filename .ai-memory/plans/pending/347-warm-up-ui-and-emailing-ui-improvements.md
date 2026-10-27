# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 347 Enterprise Spec)

- **Task Sequence:** 347
- **Task Slug:** `347-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 347 of the GitMap Warm-Up and Emailing UI Engine (Improvements 1001–1100). Sequence 347 elevates enterprise email delivery, post-quantum transport cryptography, multi-tenant mailbox virtualization, seed intelligence networks, autonomous DNS record self-healing, forensic bounce attribution, distributed SQLite multi-master sharding, and Antigravity multi-agent swarm orchestration to production grade:
- **Anycast Routing & Edge Telemetry (Improvements 1001–1030, Clusters 101–103):** Post-quantum Kyber/Dilithium TLS negotiation, hybrid cryptographic handshakes, dynamic multi-homed BGP egress routing, tenant socket segregation, and RFC 8461 MTA-STS autonomous enforcement.
- **Interactive TUI Cockpit & Mailbox Diagnostics (Improvements 1031–1060, Clusters 104–106):** Full-screen Lipgloss/Bubble Tea cluster status cockpit, synthetic multi-persona seed dialog generators, distributed seed inbox telemetry, real-time reputation vector heatmaps, and zero-CDN terminal dashboards.
- **Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis (Improvements 1061–1080, Clusters 107–108):** Recursive DoH multi-resolver consensus auditing, DNSSEC cryptographic trust chains, automated multi-provider DNS provisioning (Cloudflare, Route53, DigitalOcean), enhanced RFC 3463/5248 bounce diagnostic decoding, Bayesian spam filter heuristics, and automated mailbox quarantine circuit breaking.
- **Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1081–1100, Clusters 109–110):** Vector-clock CRDT multi-master database replication, encrypted WAL frame streaming, page-level Merkle tree verification, zero-downtime hot compaction, Antigravity multi-agent swarm coordination, autonomous deliverability feedback loops, token budget governors, and live warmup pacing velocity telemetry.

## Subtasks

- [ ] `01-anycast-routing-and-edge-telemetry.md` — Author specifications and implement Post-Quantum Encryption and Multi-Homed MTA Routing (Improvements 1001–1030, Clusters 101–103).
- [ ] `02-interactive-tui-and-mailbox-diagnostics.md` — Author specifications and implement Interactive TUI Cockpit and Seed Intelligence Network (Improvements 1031–1060, Clusters 104–106).
- [ ] `03-dns-repair-and-bounce-forensics.md` — Author specifications and implement Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 1061–1080, Clusters 107–108).
- [ ] `04-sqlite-shards-and-swarm-steerer.md` — Author specifications and implement Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1081–1100, Clusters 109–110).

## Architectural Clusters Breakdown

1. **Post-Quantum Encryption & MTA Routing (Clusters 101–103 / Improvements 1001–1030):**
   - Hybrid Kyber-1024 / Dilithium-5 post-quantum TLS cipher suites in `pkg/mail/crypto/`.
   - BGP multi-homed egress binding and IP pool rotation in `pkg/mail/routing/`.
   - Autonomous MTA-STS and DANE TLSA DNS verification in `pkg/mail/mtasts/`.

2. **Interactive TUI & Seed Intelligence (Clusters 104–106 / Improvements 1031–1060):**
   - High-density Bubble Tea responsive monitoring dashboards in `pkg/cli/tui/`.
   - Multi-persona contextual synthetic email conversation engines in `pkg/mail/seed/`.
   - Real-time cross-ISP inbox placement matrix and delivery scorecards in `pkg/mail/placement/`.

3. **Autonomous DNS Record Repair & Forensic Bounce Analysis (Clusters 107–108 / Improvements 1061–1080):**
   - Recursive DoH sentinel with DNSSEC cryptographic proof verification in `pkg/mail/dnscheck/`.
   - Automated SPF lookup flattening, DKIM key rotation, and DMARC synthesis in `pkg/mail/dnsrepair/`.
   - Granular RFC 3463/5248 enhanced bounce parsing and root-cause classifier in `pkg/mail/bounce/`.
   - Bayesian spam phrase tokenization and circuit-breaking mailbox quarantine in `pkg/mail/forensics/`.

4. **Distributed SQLite Split-DB Multi-Master & Swarm Steerer (Clusters 109–110 / Improvements 1081–1100):**
   - Vector-clock CRDT state synchronization and encrypted WAL replication in `pkg/db/cluster/`.
   - Split-DB page Merkle tree integrity verification and epoch commits in `pkg/db/splitdb/`.
   - Antigravity multi-agent swarm orchestrator and lease scheduler in `pkg/agent/swarm/`.
   - Deliverability feedback loop and adaptive warmup velocity governor in `pkg/agent/steerer/`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Leaf functions returning `error` or `Result[T]` must return the error directly. Zero calls to `os.Exit` or exit handlers inside leaf packages.
- **Function Limit:** Functions <= 15 lines maximum (<= 8 lines strongly preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Strictly positive booleans with `is` or `has` prefixes (`isRepaired`, `hasValidDkim`, `isBounced`, `isSpamTrap`, `isMasterNode`, `isSynced`, `isAgentActive`). NO negative or ambiguous booleans.
- **Zero CDN Dependencies:** 100% self-contained embedded web assets using Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all specifications, documentation, and configurations.
