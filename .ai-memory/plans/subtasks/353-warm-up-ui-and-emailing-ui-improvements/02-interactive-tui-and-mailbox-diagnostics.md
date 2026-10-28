# Subtask Plan 02: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Improvements 1631–1660)

- **Parent Task:** `353-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `353-02`
- **Spec Path:** `02-spec/21-app/353-warm-up-ui-and-emailing-ui-improvements/03-improvements-1631-to-1660-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 164: Neuromorphic Feedback Calibration & Real-time ISP Throttle Inversion (1631–1640)
  - Cluster 165: Dynamic Autonomous IP Warmup Ramp-Up Schedule & Cohort Balancer (1641–1650)
  - Cluster 166: Zero-Loss Envelope Spool Journaling & POSIX Shared-Memory Ring Buffers (1651–1660)

## Implementation Steps

1. **Neuromorphic Feedback Calibration & Throttle Inverter (1631–1640):**
   - Author leaky-integrate-and-fire (LIF) membrane potential engine accumulating microsecond TCP ACK jitter and banner latency variances in `pkg/mail/throttle/`.
   - Invert dispatch rates dynamically into micro-trickle modes before defensive gateways issue hard blocks.
   - Render real-time membrane potential gauges and stress levels via Bubble Tea interactive components.

2. **Autonomous Warmup Ramp-Up & Cohort Balancing (1641–1650):**
   - Implement mathematical sigmoid and exponential ramp curves adapting to mailbox ISP feedback and reputation thresholds.
   - Balance recipient cohort loads between seed inboxes and cold target inboxes with dynamic priority weighting.
   - Store warmup schedules and cohort distributions in `repodb/warmup.db` using positive booleans (`isRampActive`, `hasReachedCeiling`).

3. **Zero-Loss Envelope Spool & Shared-Memory Ring Buffers (1651–1660):**
   - Author lockless IPC shared-memory ring buffers for high-throughput zero-copy envelope staging in `pkg/mail/spool/`.
   - Provide crash-resilient WAL journaling guaranteeing zero envelope loss across sudden worker termination.
   - Author interactive Bubble Tea mailbox health matrix and spool utilization HUD in `pkg/cli/tui/`.

## Verification Checklist

- [ ] Neuromorphic throttle engine dampens egress spikes before 4xx deferral errors trigger.
- [ ] Ramp-up cohort balancer adjusts volume dynamically without violating ISP daily quotas.
- [ ] Ring buffer achieves microsecond spool staging and recovery on process restart.
- [ ] Positive booleans (`is`, `has`) verified across all data contracts.
- [ ] 100% relative repository paths maintained with zero absolute paths.
