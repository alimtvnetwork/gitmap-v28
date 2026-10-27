# Subtask 02: Mailbox Mesh & Terminal TUI

- **Parent Plan:** `.ai-memory/plans/pending/338-warm-up-ui-and-emailing-ui-improvements.md`
- **Spec Path:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/`
- **Sequence:** 338
- **Subtask Index:** 02
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Objectives

Implement and verify Improvements 131 through 160 spanning Clusters 14, 15, and 16:
1. **Cluster 14: Multi-Mailbox Mesh Orchestration & Pool Scheduling (Items 131–140):**
   - High-concurrency mailbox mesh scheduler with weighted round-robin and priority queue dispatching.
   - Mailbox reputation tracking matrix aggregating hourly ISP placement and bounce rates.
   - Dynamic cross-mailbox quota rebalancing to prevent early daily limit exhaustion.

2. **Cluster 15: Terminal TUI Operator Consoles & Lipgloss Telemetry (Items 141–150):**
   - Interactive terminal UI dashboard (`gitmap warmup tui`) displaying real-time sender gauges, queue depths, and live event logs.
   - High-density ASCII sparklines and TokyoNight color-styled Lipgloss tables for terminal monitoring.
   - Keyboard navigation shortcuts for pausing, resuming, and quarantining sender mailboxes directly from the console.

3. **Cluster 16: Advanced DNS & Infrastructure Health Radar (Items 151–160):**
   - Automated reverse DNS (PTR) validator with forward-confirmed verification.
   - TLS handshake inspector checking TLS 1.3 cipher suite negotiation and certificate expiration dates.
   - RBL (Real-time Blackhole List) multi-server DNSBL radar scanning 50+ international blacklists.

## Concrete Deliverables

- `pkg/mail/pool/mesh_coordinator.go`: Concurrency scheduler, quota reservations, and round-robin pools.
- `pkg/cli/tui/warmup_dashboard.go`: Lipgloss and Bubble Tea terminal operator dashboard.
- `pkg/mail/dns/rdns_checker.go`: Reverse DNS, TLS certificate inspector, and DNSBL scanner.
- Terminal commands registered in `pkg/cli/warmup_cmd.go` (`gitmap warmup tui`, `gitmap mail rbl-check`).

## Quality Standards

- Leaf error sovereignty: Return `*apperror.AppError` directly.
- Strict positive boolean discipline (`isMailboxReady`, `hasRdnsMatch`, `isTlsValid`).
- Zero nested conditionals with guard clauses and early returns.
