# Subtask 01: Connection Pools, TLS Multiplexers & Velocity Regulators

## Status: READY
- **Subtask Slug:** `01-connection-pools-and-velocity-regulator`
- **Parent Plan:** `.ai-memory/plans/pending/342-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/01-improvements-501-to-530-connection-pools-and-velocity.md`
- **Target Range:** Improvements 501–530 (Clusters 51–53)
- **Target Subsystem:** `pkg/mail/pool/`, `pkg/mail/velocity/`, `pkg/mail/telemetry/`
- **Exclusive Owned Files:**
  - `pkg/mail/pool/connection_pool.go`
  - `pkg/mail/pool/tls_multiplexer.go`
  - `pkg/mail/pool/socket_reuser.go`
  - `pkg/mail/velocity/rate_limiter.go`
  - `pkg/mail/velocity/leaky_bucket.go`
  - `pkg/mail/velocity/throttle_controller.go`
  - `pkg/mail/telemetry/socket_profiler.go`
  - `cli/cmdwarmup/pool_cmd.go`

---

## 1. Context & Objectives

This subtask delivers enterprise connection pooling and velocity regulation for the warm-up pipeline:
1. **SMTP Connection Pooling & Multiplexing (Cluster 51, Items 501–510):** Persistent keep-alive sockets, TLS handshake caching, connection starvation detection, and connection re-use optimization.
2. **Adaptive Velocity Regulation & Leaky Bucket Schedulers (Cluster 52, Items 511–520):** Adaptive rate limiting, dynamic windowing per sending IP/domain, feedback-driven speed throttling, and smooth delivery cadence.
3. **Socket Telemetry & Protocol Diagnostics (Cluster 53, Items 521–530):** Microsecond TCP latency monitoring, TLS cipher validation, connection pool exhaustion alarms, and socket lifecycle metrics.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isPoolActive`, `hasAvailableConnection`, `isThrottled`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies.
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Connection Pool & TLS Multiplexer
- Implement thread-safe SMTP connection manager in `pkg/mail/pool/connection_pool.go`.
- Build persistent TLS session ticket cache and handshake optimizer in `pkg/mail/pool/tls_multiplexer.go`.
- Construct connection reuse policy and idle reaper in `pkg/mail/pool/socket_reuser.go`.

### Step 2: Adaptive Velocity Regulation
- Construct token bucket and leaky bucket rate governors in `pkg/mail/velocity/`.
- Build ISP-specific throttle controller that scales send rate based on SMTP 421/450 response codes.
- Implement delivery cadence smoothing in `pkg/mail/velocity/leaky_bucket.go`.

### Step 3: Telemetry & CLI Wiring
- Build socket profiling and latency metric collector in `pkg/mail/telemetry/socket_profiler.go`.
- Wire CLI pool management commands in `cli/cmdwarmup/pool_cmd.go`.
