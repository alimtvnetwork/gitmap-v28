# Subtask Plan 01: Autonomous SMTP Session Pipelining and TCP Zero-Window Detection (Improvements 1301–1330)

- **Parent Task:** `350-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `350-01`
- **Spec Path:** `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/02-improvements-1301-to-1330-session-pipelining-and-zero-window.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 131: High-Throughput SMTP Session Pipelining & RFC 2920 Command Streaming (1301–1310)
  - Cluster 132: Low-Level Kernel Socket Diagnostics & Zero-Window Probing (1311–1320)
  - Cluster 133: Interactive Pipelining Flamegraphs & Socket Buffer Cockpit (1321–1330)

## Implementation Steps

1. **RFC 2920 Pipelined Command Streaming (1301–1310):**
   - Implement pipelined command buffering for `RSET`, `MAIL FROM`, `RCPT TO`, and `DATA` in `pkg/mail/transport/`.
   - Implement zero-RTT command burst streaming for validated MX peers.
   - Enforce asynchronous multi-response status correlation with partial reject tolerance.
   - Author pipelined connection pooling and TLS session resumption reuse cache.

2. **Zero-Window Detection & Kernel Socket Diagnostics (1311–1320):**
   - Implement zero-copy kernel socket buffer telemetry via `getsockopt(..., TCP_INFO)` in `pkg/mail/socket/`.
   - Detect receiver TCP window exhaustion and zero-window persistent probes.
   - Author adaptive backpressure throttling and sliding window scale factor estimators.
   - Record socket buffer metrics in Split-DB (`repodb/pipeline.db`) using positive booleans (`isWindowExhausted`, `hasSndBufStall`).

3. **Bubble Tea Flamegraph HUD & Real-Time Socket Visualizer (1321–1330):**
   - Build Lipgloss ANSI terminal flamegraph visualizer rendering sub-millisecond command execution phases.
   - Author real-time TCP socket buffer gauge and window scale sparklines in Bubble Tea v2.
   - Implement zero-CDN Shadow DOM Web Component visualizer for local web dashboards.
   - Provide CLI diagnostic commands under `gitmap warmup pipelining` and `gitmap warmup socket-hud`.

## Verification Checklist

- [ ] Command streaming pipeline sends batched `MAIL FROM` and `RCPT TO` commands in a single network packet.
- [ ] `getsockopt` correctly polls `tcpi_snd_wnd` on Linux and fallback metrics on Windows.
- [ ] Bubble Tea flamegraph displays millisecond-accurate phase timings with Lipgloss theme tokens.
- [ ] All structs and SQLite schemas use strictly positive booleans (`isPipelined`, `hasWindowAlert`).
