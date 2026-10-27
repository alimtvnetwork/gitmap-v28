# Plan: Warm-Up UI and Emailing UI Improvements (Sequence 349 Enterprise Spec)

- **Task Sequence:** 349
- **Task Slug:** `349-warm-up-ui-and-emailing-ui-improvements`
- **Spec Path:** `02-spec/21-app/349-warm-up-ui-and-emailing-ui-improvements/`
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Executive Summary

This master plan orchestrates the complete engineering lifecycle for Sequence 349 of the GitMap Warm-Up and Emailing UI Engine (Improvements 1201–1300). Sequence 349 establishes high-throughput SMTP session pipelining, zero-copy socket buffer monitoring, interactive Bubble Tea terminal cockpits, autonomous DNS self-healing via DoH consensus, forensic bounce decomposition, distributed multi-master SQLite Split-DB sharding, and Antigravity multi-agent swarm orchestration across 10 core clusters:
- **Session Pipelining & Zero-Window Detection (Improvements 1201–1230, Clusters 121–123):** RFC 2920 command streaming pipelines, zero-copy socket buffer monitoring via `TCP_INFO`, TCP zero-window probe mitigation, connection reuse pools, keep-alive heartbeat calibration, cross-ISP TCP handshake latency scoring, Bubble Tea v2 HUD viewports, Lipgloss sub-millisecond phase flamegraphs, and zero-CDN Shadow DOM telemetry.
- **Interactive TUI Cockpit & Mailbox Diagnostics (Improvements 1231–1260, Clusters 124–126):** Interactive Bubble Tea seed pool verification matrix, real-time IMAP IDLE synchronizer, synthetic AI peer engagement workflows, deep envelope return-path & DMARC forensic report (RUF) analyzers, automated spam folder rescue daemons, seed mailbox IP reputation & RBL blacklists, SMTP MX latency and TLS handshake cipher inspectors, OAuth2 credential rotators, mailbox quota purgers, and Feedback Loop (FBL) abuse complaint interceptors.
- **Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis (Improvements 1261–1280, Clusters 127–128):** Recursive DoH multi-resolver consensus auditing, DNSSEC cryptographic trust chains, automated SPF lookup flattening (<10 lookups) with CIDR aggregation, dual-selector DKIM key rotation (RSA-2048 & Ed25519), DMARC alignment policy synthesis, dynamic Route53/Cloudflare/DigitalOcean DNS API provisioning, real-time propagation radar, continuous DNS drift circuit breakers, enhanced RFC 3463/5248 bounce status decoding, Bayesian spam trigger phrase scoring, real-time interactive terminal bounce forensics HUD, automated mailbox quarantine circuit breaking, SIEM webhook forwarding, spam trap honeypot sentinels, suppression list Bloom filters, and predictive bounce trend modeling.
- **Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1281–1300, Clusters 129–130):** Vector-clock CRDT multi-master database replication, encrypted WAL frame streaming over mutual TLS 1.3, page-level Merkle tree verification, distributed two-phase epoch commit coordinators, non-blocking hot VACUUM compaction, real-time replication cluster HUD, autonomous conflict resolution via LWW arbiters, dynamic partition rebalancing, point-in-time recovery log replay, zero-downtime rolling schema migrations, Antigravity multi-agent swarm cockpit HUD, dynamic subtask auctioneering with expiring leases, token expenditure governors, circular deadlock detection, agent consensus Raft elections, deliverability feedback loops with autonomous throttle steering, multi-persona contextual synthetic dialogues, ISP reputation vector heatmaps, autonomous warmup graduation gates, and full-screen terminal mission telemetry orchestration.

## Subtasks

- [ ] `01-session-pipelining-and-zero-window.md` — Author specifications and implement Autonomous SMTP Session Pipelining and TCP Zero-Window Detection (Improvements 1201–1230, Clusters 121–123).
- [ ] `02-interactive-tui-and-mailbox-diagnostics.md` — Author specifications and implement Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Improvements 1231–1260, Clusters 124–126).
- [ ] `03-dns-repair-and-bounce-forensics.md` — Author specifications and implement Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 1261–1280, Clusters 127–128).
- [ ] `04-sqlite-shards-and-swarm-steerer.md` — Author specifications and implement Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Improvements 1281–1300, Clusters 129–130).

## Architectural Clusters Breakdown

1. **Session Pipelining & TCP Zero-Window Telemetry (Clusters 121–123 / Improvements 1201–1230):**
   - RFC 2920 command streaming pipelines and zero-RTT command bursts in `pkg/mail/transport/`.
   - Kernel TCP socket buffer zero-copy monitoring via `getsockopt(..., TCP_INFO)` in `pkg/mail/socket/`.
   - High-definition Lipgloss sub-millisecond flamegraphs and SVG waveforms in `pkg/cli/tui/`.

2. **Interactive TUI Cockpit & Mailbox Diagnostics (Clusters 124–126 / Improvements 1231–1260):**
   - Bubble Tea interactive multi-provider seed verification matrix in `pkg/cli/tui/`.
   - Real-time IMAP IDLE synchronizer and seed arrival notification HUD in `pkg/mail/seed/`.
   - Dynamic backoff inverters and concurrency governors in `pkg/mail/throttle/`.

3. **Autonomous DNS Record Repair & Forensic Bounce Analysis (Clusters 127–128 / Improvements 1261–1280):**
   - Recursive DoH sentinel with DNSSEC cryptographic proof verification in `pkg/mail/dnscheck/`.
   - Automated SPF lookup flattening, DKIM key rotation, and DMARC synthesis in `pkg/mail/dnsrepair/`.
   - Granular RFC 3463/5248 enhanced bounce parsing and root-cause classifier in `pkg/mail/bounce/`.

4. **Distributed SQLite Split-DB Multi-Master Shards & Swarm Steerer (Clusters 129–130 / Improvements 1281–1300):**
   - Vector-clock CRDT state synchronization and encrypted WAL frame streaming in `pkg/db/cluster/`.
   - Page-level Merkle tree shard integrity auditor and two-phase epoch commits in `pkg/db/cluster/`.
   - Antigravity multi-agent swarm cockpit HUD and dynamic subtask auctioneering in `pkg/agent/swarm/`.
   - Real-time deliverability feedback loop, reputation scorecards, and warmup graduation in `pkg/agent/steerer/`.

## Quality Gates & Verification Standards

- **Strict Path Hygiene:** All file references across specifications, plans, and source code must use 100% relative repository paths (`02-spec/...`, `.ai-memory/...`). No absolute paths or Windows drive letters.
- **Positive-Only Booleans:** All Go struct fields, database column names, and CLI boolean flags must use positive prefixes (`is`, `has`). Banned prefixes (`can`, `should`, `was`, `will`) and negative polarity are strictly forbidden.
- **Split-DB Storage Separation:** Segregate operational runtime state (`repodb/warmup.db` or `repodb/pipeline.db`) from compliance ledgers (`repodb/audit.db`). All SQLite tables must use PascalCase formatting and camelCase column definitions.
- **Leaf Error Sovereignty:** All leaf functions must return structured errors (`return apperror.Wrap(...)`). Direct calls to `os.Exit` or CLI exit handlers inside leaf functions are prohibited.
- **Zero External CDN:** All TUI and web visualizers must render natively via Lipgloss ANSI formatting or inline SVGs and Shadow DOM Web Components without external asset downloads.
