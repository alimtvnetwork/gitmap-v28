# Subtask Plan 01: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (Improvements 1401–1430)

- **Parent Task:** `351-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `351-01`
- **Spec Path:** `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/02-improvements-1401-to-1430-session-pipelining-and-zero-window.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 141: Autonomous SMTP Session Pipelining & Zero-RTT Command Bursting (1401–1410)
  - Cluster 142: Kernel Socket Buffer Diagnostics & TCP Zero-Window Sentinel (1411–1420)
  - Cluster 143: Socket Pool Multiplexing & Cross-ISP Handshake Latency Profiler (1421–1430)

## Implementation Steps

1. **SMTP Session Pipelining & Command Streaming (1401–1410):**
   - Implement RFC 2920 SMTP command pipelining engine streaming `EHLO`, `MAIL FROM`, `RCPT TO`, and `DATA` in a single network round-trip.
   - Author adaptive command burst throttlers and error short-circuiting handlers in `pkg/mail/transport/`.
   - Build real-time Lipgloss terminal pipelining flamegraph showing sub-millisecond command execution latencies.

2. **Kernel Socket Diagnostics & Zero-Window Mitigation (1411–1420):**
   - Query OS kernel socket telemetry via `getsockopt` with `TCP_INFO` to monitor send/receive buffer occupancy and congestion window (`cwnd`).
   - Detect TCP zero-window advertising from receiving MTAs, pausing outbound writes and triggering dynamic backoff.
   - Maintain positive boolean states (`hasZeroWindowOccurred`, `isBufferSaturated`) in `repodb/pipeline.db`.

3. **Socket Pool Multiplexing & ISP Latency Scoring (1421–1430):**
   - Architect managed TLS connection pool multiplexing warm SMTP sessions across target domains.
   - Record per-ISP handshake latency profiles (TLS 1.3 resumption, cipher negotiation duration) for Gmail, Outlook, Yahoo, and Apple Mail.
   - Expose interactive Bubble Tea socket pool dashboard with connection state visualizers.

## Verification Checklist

- [ ] RFC 2920 command streaming reduces SMTP handshake latency by >= 40%.
- [ ] Kernel `TCP_INFO` diagnostics capture zero-window events and trip backoff daemons.
- [ ] Positive booleans (`is`, `has`) verified across all data contracts.
- [ ] 100% relative repository paths maintained.
