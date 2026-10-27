# Subtask Plan 01: Autonomous SMTP Session Pipelining & TCP Zero-Window Detection (Improvements 1501–1530)

- **Parent Task:** `352-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `352-01`
- **Spec Path:** `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/02-improvements-1501-to-1530-session-pipelining-and-zero-window.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 151: Autonomous SMTP Session Pipelining & Zero-RTT Command Bursting (1501–1510)
  - Cluster 152: Kernel Socket Buffer Diagnostics & TCP Zero-Window Sentinel (1511–1520)
  - Cluster 153: Socket Pool Multiplexing & Cross-ISP Handshake Latency Profiler (1521–1530)

## Implementation Steps

1. **Autonomous SMTP Session Pipelining & Command Streaming (1501–1510):**
   - Implement RFC 2920 pipelined command multiplexing streaming `RSET`, `MAIL FROM`, `RCPT TO`, and `DATA` in a single coalesced TCP segment burst in `pkg/mail/transport/`.
   - Author adaptive command burst throttlers and error short-circuiting handlers in `pkg/mail/transport/pipeline.go`.
   - Build real-time Lipgloss terminal pipelining flamegraph rendering sub-millisecond command execution latencies.

2. **Kernel Socket Diagnostics & Zero-Window Mitigation (1511–1520):**
   - Query OS kernel socket telemetry via `getsockopt` with `TCP_INFO` to monitor send/receive buffer occupancy, unacknowledged packets, and congestion window (`cwnd`).
   - Detect TCP zero-window advertising from receiving MTAs, pausing outbound writes and triggering dynamic backoff probes.
   - Maintain positive boolean states (`hasZeroWindowOccurred`, `isBufferSaturated`, `isEgressThrottled`) in `repodb/warmup.db`.

3. **Socket Pool Multiplexing & Cross-ISP Latency Scoring (1521–1530):**
   - Architect managed TLS connection pool multiplexing warm SMTP sessions across target domains.
   - Record per-ISP handshake latency profiles (TLS 1.3 0-RTT resumption, cipher negotiation duration) for Gmail, Microsoft 365, Yahoo, and Apple Mail.
   - Expose interactive Bubble Tea socket pool dashboard with connection state visualizers and SVG waveforms in `pkg/cli/tui/`.

## Verification Checklist

- [ ] RFC 2920 command streaming reduces SMTP handshake latency by >= 40%.
- [ ] Kernel `TCP_INFO` diagnostics capture zero-window events and trip backoff daemons.
- [ ] Positive booleans (`is`, `has`) verified across all data contracts.
- [ ] 100% relative repository paths maintained with zero absolute paths.
