# Subtask 01: TLS Fingerprint Morphing & Distributed Rate Limit Mesh

## Status: READY
- **Subtask Slug:** `01-fingerprint-morphing-and-rate-limit-mesh`
- **Parent Plan:** `.ai-memory/plans/pending/344-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/01-improvements-701-to-730-fingerprint-morphing-and-rate-limit-mesh.md`
- **Target Range:** Improvements 701–730 (Clusters 71–73)
- **Target Subsystem:** `pkg/mail/fingerprint/`, `pkg/mail/emulation/`, `pkg/mail/ratelimit/`
- **Exclusive Owned Files:**
  - `pkg/mail/fingerprint/ja3_morpher.go`
  - `pkg/mail/fingerprint/ja4_signature.go`
  - `pkg/mail/fingerprint/tcp_window_scaler.go`
  - `pkg/mail/fingerprint/cipher_permutator.go`
  - `pkg/mail/emulation/mouse_entropy.go`
  - `pkg/mail/emulation/scroll_cadence.go`
  - `pkg/mail/emulation/webmail_session.go`
  - `pkg/mail/ratelimit/token_bucket_mesh.go`
  - `pkg/mail/ratelimit/isp_adaptive_throttle.go`
  - `pkg/mail/ratelimit/distributed_sync.go`
  - `cli/cmdmail/fingerprint_cmd.go`
  - `cli/cmdmail/ratelimit_cmd.go`

---

## 1. Context & Objectives

This subtask delivers carrier-grade TLS client fingerprint virtualization, behavioral webmail humanization, and cross-node distributed rate limiting:
1. **Dynamic TLS JA3/JA4 Fingerprint Morphing (Cluster 71, Items 701–710):** Cryptographic permutation of client hello cipher suites, TLS extensions, elliptic curves, point formats, and TCP MSS/window scale emulation to prevent network-level mail agent fingerprinting.
2. **Behavioral Webmail Session Humanization (Cluster 72, Items 711–720):** Realistic human mouse trajectory generation via Bézier curves, micro-pause keystroke jitter, viewport scroll cadence simulation, and cookie lifecycle persistence.
3. **Distributed Rate Limiting Token Bucket Mesh (Cluster 73, Items 721–730):** Peer-to-peer token synchronization, ISP-specific dynamic throttle gates, burst allowance damping, and backpressure telemetry.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isFingerprintMorphed`, `hasValidJa4Signature`, `isWithinRateLimit`, `hasAvailableToken`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: JA3/JA4 TLS Fingerprint Engine
- Implement client hello cipher permutation in `pkg/mail/fingerprint/ja3_morpher.go`.
- Build JA4 signature computation and validation in `pkg/mail/fingerprint/ja4_signature.go`.
- Implement TCP window scaling and MSS emulator in `pkg/mail/fingerprint/tcp_window_scaler.go`.
- Build cipher extension randomize passes in `pkg/mail/fingerprint/cipher_permutator.go`.

### Step 2: Behavioral Humanization
- Build cubic Bézier mouse movement engine in `pkg/mail/emulation/mouse_entropy.go`.
- Implement variable keystroke pause generator in `pkg/mail/emulation/scroll_cadence.go`.
- Construct authenticated webmail session manager in `pkg/mail/emulation/webmail_session.go`.

### Step 3: Distributed Rate Limiting Mesh
- Implement distributed token bucket algorithm in `pkg/mail/ratelimit/token_bucket_mesh.go`.
- Build adaptive ISP throttle controllers in `pkg/mail/ratelimit/isp_adaptive_throttle.go`.
- Implement gossip-based token balance synchronizer in `pkg/mail/ratelimit/distributed_sync.go`.

### Step 4: CLI Command Dispatch
- Wire fingerprint inspection commands into `cli/cmdmail/fingerprint_cmd.go`.
- Wire mesh rate limiting controls into `cli/cmdmail/ratelimit_cmd.go`.
