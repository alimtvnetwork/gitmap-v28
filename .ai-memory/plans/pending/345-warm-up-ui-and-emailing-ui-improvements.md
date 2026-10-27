# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 345 Enterprise Spec)

- **Task Sequence:** 345
- **Task Slug:** `345-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 345 of the GitMap Warm-Up and Emailing UI Engine (Improvements 801–900). Sequence 345 elevates automated email warmup, deliverability governance, cryptographic validation, forensic bounce categorization, distributed multi-master SQLite state synchronization, and autonomous multi-agent swarm steering to enterprise operational maturity:
- **SMTP Pipeline & Throttle Inversion (Improvements 801–830, Clusters 81–83):** Inverted adaptive SMTP pacing, dynamic backpressure governors, TLS cipher negotiation tuning, automated MX socket multiplexing, and cross-ISP priority queuing.
- **Interactive TUI & Mailbox Diagnostics (Improvements 831–860, Clusters 84–86):** Full-screen interactive Bubble Tea terminal cockpit, real-time seed inbox response telemetry, Bayesian conversation sentiment tracking, multi-ISP delivery matrix visualizer, and live dispatch pacing controls.
- **Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis (Improvements 861–880, Clusters 87–88):** DoH sentinel validation, DNSSEC chain audit, automated DNS record patch generation, multi-provider API auto-provisioning, DMARC alignment auditing, enhanced SMTP bounce classification (RFC 3463/5248), Bayesian spam filter attribution, asynchronous webhook dispatch, and automated quarantine quarantine gating.
- **Distributed SQLite Split-DB Multi-Master Replication & Swarm Steerer (Improvements 881–900, Clusters 89–90):** Vector-clock CRDT multi-master synchronization, distributed Raft-assisted leases, zero-downtime hot backup, cross-node WAL stream ingestion, self-healing corruption detection, Antigravity multi-agent swarm orchestrator, cooperative dynamic workload steering, token expenditure governors, and live warmup telemetry bridge HUD.

## Subtasks

- [ ] `01-smtp-pipeline-and-throttle-inversion.md` — Author specifications and implement SMTP Pipeline and Throttle Inversion Engine (Improvements 801–830, Clusters 81–83).
- [ ] `02-interactive-tui-and-mailbox-diagnostics.md` — Author specifications and implement Interactive TUI and Mailbox Diagnostics Cockpit (Improvements 831–860, Clusters 84–86).
- [ ] `03-dns-repair-and-bounce-forensics.md` — Author specifications and implement Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 861–880, Clusters 87–88).
- [ ] `04-sqlite-shards-and-swarm-steerer.md` — Author specifications and implement Distributed SQLite Split-DB Multi-Master Replication HUD & Swarm Steerer (Improvements 881–900, Clusters 89–90).

## Architectural Clusters Breakdown

1. **SMTP Pipeline & Throttle Inversion (Clusters 81–83 / Improvements 801–830):**
   - Inverted token-bucket pacing and dynamic backpressure in `pkg/mail/throttle/`.
   - Socket multiplexing and pipelined SMTP transmission in `pkg/mail/smtp/`.
   - Adaptive cross-ISP priority dispatch queues in `pkg/mail/queue/`.

2. **Interactive TUI & Mailbox Diagnostics (Clusters 84–86 / Improvements 831–860):**
   - Bubble Tea interactive cockpit and terminal visualizers in `pkg/cli/tui/`.
   - Seed inbox IMAP/POP3 telemetry and Bayesian sentiment tracking in `pkg/mail/diagnostics/`.
   - Cross-ISP delivery heatmaps and placement matrix HUD in `pkg/mail/placement/`.

3. **Autonomous DNS Record Repair & Forensic Bounce Analysis (Clusters 87–88 / Improvements 861–880):**
   - DNS over HTTPS (DoH) sentinel with DNSSEC cryptographic validation in `pkg/mail/dnscheck/`.
   - Autonomous SPF/DKIM/DMARC record synthesizer and provider API patcher in `pkg/mail/dnsrepair/`.
   - Forensic SMTP bounce parser with RFC 3463/5248 code attribution in `pkg/mail/bounce/`.
   - Bayesian spam trigger phrase and header heuristic classifier in `pkg/mail/forensics/`.

4. **Distributed SQLite Split-DB Multi-Master & Swarm Steerer (Clusters 89–90 / Improvements 881–900):**
   - Multi-master vector clock CRDT state synchronization and WAL frame distribution in `pkg/db/cluster/`.
   - Distributed split-db schema migration and hot-backup engine in `pkg/db/splitdb/`.
   - Antigravity multi-agent swarm orchestrator and lease coordinator in `pkg/agent/swarm/`.
   - Autonomous deliverability feedback loop and warmup pace steerer in `pkg/agent/steerer/`.

## Quality & Architectural Invariants

- **Error Return Sovereignty:** Zero exits in leaf packages. All errors returned directly via typed `*apperror.AppError`. Outer callers alone manage termination.
- **Function Limit:** Functions <= 15 lines max (<= 8 lines preferred).
- **Conditionals:** Zero nested `if` statements. Guard clauses and early returns only.
- **Boolean Naming:** Positive booleans with `is` or `has` prefixes exclusively (`isRepaired`, `hasValidDkim`, `isBounced`, `isSpamTrap`, `isMasterNode`, `isSynced`, `isAgentActive`).
- **Zero CDN Dependencies:** 100% self-contained embedded web assets via Go `embed.FS`.
- **Relative Path Standards:** 100% relative repository paths across all documentation and configurations.
