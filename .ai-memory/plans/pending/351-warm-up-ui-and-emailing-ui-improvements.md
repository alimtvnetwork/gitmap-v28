# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 351 Enterprise Spec)

- **Task Sequence:** 351
- **Task Slug:** `351-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 351 of the GitMap Warm-Up and Emailing UI Engine (Improvements 1401–1500). Sequence 351 extends and hardens high-throughput SMTP session pipelining, socket diagnostics, interactive Bubble Tea terminal cockpits, automated DNS self-healing via DoH consensus, forensic bounce decomposition, distributed multi-master SQLite Split-DB sharding, and Antigravity multi-agent swarm orchestration across 10 core clusters:
- **Session Pipelining & Zero-Window Detection (Improvements 1401–1430, Clusters 141–143):** High-throughput command streaming pipelining, kernel zero-copy socket buffer monitoring via `TCP_INFO`, TCP zero-window probe mitigations, pooled socket multiplexing, cross-ISP handshake latency scoring, Bubble Tea v2 viewport renderers, Lipgloss sub-millisecond flamegraphs, and zero-CDN Shadow DOM telemetry.
- **Interactive TUI Cockpit & Mailbox Diagnostics (Improvements 1431–1460, Clusters 144–146):** Bubble Tea interactive seed pool verification matrix, real-time IMAP IDLE event listeners, synthetic AI peer conversational workflows, deep envelope return-path & DMARC forensic report (RUF) analyzers, automated spam folder rescue daemons, seed mailbox IP reputation & RBL blacklists, SMTP MX latency and TLS handshake cipher inspectors, OAuth2 token credential rotators, mailbox quota purgers, and Feedback Loop (FBL) abuse complaint interceptors.
- **Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis (Improvements 1461–1480, Clusters 147–148):** Recursive DoH multi-resolver consensus auditing, DNSSEC cryptographic trust chains, automated SPF lookup flattening (<10 lookups) with CIDR aggregation, dual-selector DKIM key rotation (RSA-2048 & Ed25519), DMARC alignment policy synthesis, dynamic Route53/Cloudflare/DigitalOcean DNS API provisioning, real-time propagation radar, continuous DNS drift circuit breakers, enhanced RFC 3463/5248 bounce status decoding, Bayesian spam trigger phrase scoring, real-time interactive terminal bounce forensics HUD, automated mailbox quarantine circuit breaking, SIEM webhook forwarding, spam trap honeypot sentinels, suppression list Bloom filters, and predictive bounce trend modeling.
- **Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1481–1500, Clusters 149–150):** Vector-clock CRDT multi-master database replication, encrypted WAL frame streaming over mutual TLS 1.3, page-level Merkle tree verification, distributed two-phase epoch commit coordinators, non-blocking hot VACUUM compaction, real-time replication cluster HUD, autonomous conflict resolution via LWW arbiters, dynamic partition rebalancing, point-in-time recovery log replay, zero-downtime rolling schema migrations, Antigravity multi-agent swarm cockpit HUD, dynamic subtask auctioneering with expiring leases, token expenditure governors, circular deadlock detection, agent consensus Raft elections, deliverability feedback loops with autonomous throttle steering, multi-persona contextual synthetic dialogues, ISP reputation vector heatmaps, autonomous warmup graduation gates, and full-screen terminal mission telemetry orchestration.

## Subtasks

- [ ] `01-session-pipelining-and-zero-window.md` — Author specifications and implement Autonomous SMTP Session Pipelining and TCP Zero-Window Detection (Improvements 1401–1430, Clusters 141–143).
- [ ] `02-interactive-tui-and-mailbox-diagnostics.md` — Author specifications and implement Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Improvements 1431–1460, Clusters 144–146).
- [ ] `03-dns-repair-and-bounce-forensics.md` — Author specifications and implement Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 1461–1480, Clusters 147–148).
- [ ] `04-sqlite-shards-and-swarm-steerer.md` — Author specifications and implement Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1481–1500, Clusters 149–150).

## Architectural Clusters Breakdown

1. **Session Pipelining & TCP Zero-Window Telemetry (Clusters 141–143 / Improvements 1401–1430):**
   - RFC 2920 command streaming pipelines and zero-RTT command bursts in `pkg/mail/transport/`.
   - Kernel TCP socket buffer zero-copy monitoring via `getsockopt(..., TCP_INFO)` in `pkg/mail/socket/`.
   - High-definition Lipgloss sub-millisecond flamegraphs and SVG waveforms in `pkg/cli/tui/`.

2. **Interactive TUI Cockpit & Mailbox Diagnostics (Clusters 144–146 / Improvements 1431–1460):**
   - Bubble Tea interactive multi-provider seed verification matrix in `pkg/cli/tui/`.
   - Real-time IMAP IDLE synchronizer and seed arrival notification HUD in `pkg/mail/seed/`.
   - Dynamic backoff inverters and concurrency governors in `pkg/mail/throttle/`.

3. **Autonomous DNS Record Repair & Forensic Bounce Analysis (Clusters 147–148 / Improvements 1461–1480):**
   - Recursive DoH sentinel with DNSSEC cryptographic proof verification in `pkg/mail/dnscheck/`.
   - Automated SPF lookup flattening, DKIM key rotation, and DMARC synthesis in `pkg/mail/dnsrepair/`.
   - Granular RFC 3463/5248 enhanced bounce parsing and root-cause classifier in `pkg/mail/bounce/`.

4. **Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Clusters 149–150 / Improvements 1481–1500):**
   - Vector-clock CRDT state synchronization and encrypted WAL frame streaming in `pkg/db/cluster/`.
   - Page-level Merkle tree shard integrity auditor and two-phase epoch commits in `pkg/db/cluster/`.
   - Antigravity multi-agent swarm cockpit HUD and dynamic subtask auctioneering in `pkg/agent/swarm/`.
   - Real-time deliverability feedback loop, reputation scorecards, and warmup graduation in `pkg/agent/steerer/`.

## Quality Gates & Verification Standards

- **Strict Path Hygiene:** All file references across specifications, plans, and source code must use 100% relative repository paths (`02-spec/...`, `.ai-memory/...`). No absolute paths or Windows drive letters.
- **Positive Booleans Only:** Enforce positive booleans (`is`, `has`) exclusively across Go types and SQLite DDL.
- **Split-DB Isolation:** Respect boundaries between `gitmap.db`, `installation.db`, and `repodb/pipeline.db`.
- **Zero CDN & Self-Contained UI:** All Lipgloss styling, SVG visualizations, and web templates operate without external dependencies.
