# Subtask Plan 01: SMTP Pipeline & Throttle Inversion Engine (Improvements 801–830)

- **Parent Task:** `345-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `345-01`
- **Spec Path:** `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/02-improvements-801-to-830-smtp-pipeline-and-throttle-inversion.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 81: High-Throughput SMTP Socket Multiplexing & Connection Pooling (801–810)
  - Cluster 82: Inverted Adaptive Backpressure & Rate Limit Token Mesh (811–820)
  - Cluster 83: Cross-ISP Dynamic Pacing & Envelope Header Morphing (821–830)

## Implementation Steps

1. **SMTP Connection Pool & Multiplexer Architecture (801–810):**
   - Implement asynchronous keep-alive socket reuse for high-volume warm-up endpoints in `pkg/mail/smtp/`.
   - Implement pipelined command batching (`RSET`, `MAIL FROM`, `RCPT TO`) complying with RFC 2920.
   - Enforce positive booleans (`isSocketOpen`, `hasTlsSession`, `isReusedConnection`).

2. **Inverted Adaptive Throttle & Backpressure Governor (811–820):**
   - Implement inverted token bucket algorithm dynamically responding to ISP deferral codes (421/450/451) in `pkg/mail/throttle/`.
   - Implement memory-bounded sliding window telemetry buffers for rate estimation.
   - Guard against thread starvation and clamp burst spikes via backpressure feedback.

3. **Cross-ISP Envelope Customization & Morphing (821–830):**
   - Build ISP-specific header ordering matrices (Gmail, Outlook, Yahoo/AOL, Proton) in `pkg/mail/envelope/`.
   - Implement deterministic Message-ID entropy randomization and MIME boundary morphing.
   - Enforce strict error wrapping with `*apperror.AppError` and zero direct exit calls.

## Verification Checklist

- [ ] Socket pool benchmark confirms zero connection leaks under 5,000 parallel test dispatches.
- [ ] Rate-limit governor adjusts token allocation within 50ms upon receiving simulated 421 deferral signals.
- [ ] Envelope synthesizer generates compliant RFC 5322 MIME structures with valid positive booleans.
