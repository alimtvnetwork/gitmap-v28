# Subtask Plan 01: Session Pipelining and TCP Zero-Window Detection (Improvements 1201–1230)

- **Parent Task:** `349-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `349-01`
- **Spec Path:** `02-spec/21-app/349-warm-up-ui-and-emailing-ui-improvements/02-improvements-1201-to-1230-session-pipelining-and-zero-window.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 121: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (1201–1210)
  - Cluster 122: Adaptive Multi-ISP Connection Pool Telemetry & Socket Health Matrix (1211–1220)
  - Cluster 123: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1221–1230)

## Implementation Steps

1. **SMTP Session Pipelining & Command Burst Engine (1201–1210):**
   - Implement RFC 2920 command streaming pipeline packing `RSET`, `MAIL FROM`, and `RCPT TO` into single TCP payload bursts in `pkg/mail/transport/`.
   - Implement kernel TCP zero-window detection via zero-copy reflection with `getsockopt(..., TCP_INFO)`.
   - Provide sliding window backpressure when remote receive buffers fall below 4KB.
   - Enforce positive booleans (`isPipelineForced`, `hasEarlyData`, `isWindowSaturated`).

2. **Multi-ISP Connection Pool Telemetry & Socket Health (1211–1220):**
   - Implement dynamic connection pooling with provider-specific keep-alive heartbeat calibration in `pkg/mail/socket/`.
   - Compute real-time socket health matrices evaluating TCP handshake latency, TLS negotiation overhead, and RST drop ratios.
   - Enforce leaf error sovereignty with `*apperror.AppError` and zero process exit handlers in leaf functions.

3. **Bubble Tea Jitter Waveforms & Sub-Millisecond Flamegraphs (1221–1230):**
   - Build interactive Bubble Tea v2 HUD rendering sub-millisecond execution phase flamegraphs (DNS, TCP, TLS, EHLO, DATA) in Lipgloss TrueColor.
   - Generate zero-CDN SVG latency waveform visualizers for terminal and embedded web dashboards.
   - Persist socket metrics into `repodb/warmup.db` using positive schema flags (`isHealthy`, `hasZeroWindowRisk`).

## Verification Checklist

- [ ] RFC 2920 command pipeline reduces multi-turn RTT overhead by >= 60% across simulated 100ms latency links.
- [ ] Kernel TCP zero-window detection catches buffer saturation prior to packet drops.
- [ ] Bubble Tea flamegraph TUI renders sub-millisecond phases without UI lag or terminal screen flicker.
