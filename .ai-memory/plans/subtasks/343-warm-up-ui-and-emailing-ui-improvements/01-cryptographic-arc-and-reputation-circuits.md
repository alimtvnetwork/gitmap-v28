# Subtask 01: Cryptographic ARC Chain Engine & Cross-Tenant Reputation Circuits

## Status: READY
- **Subtask Slug:** `01-cryptographic-arc-and-reputation-circuits`
- **Parent Plan:** `.ai-memory/plans/pending/343-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/02-improvements-601-to-630-cryptographic-arc-and-reputation-circuits.md`
- **Target Range:** Improvements 601–630 (Clusters 61–63)
- **Target Subsystem:** `pkg/mail/arc/`, `pkg/mail/reputation/`, `pkg/mail/circuit/`
- **Exclusive Owned Files:**
  - `pkg/mail/arc/chain_validator.go`
  - `pkg/mail/arc/seal_signer.go`
  - `pkg/mail/arc/header_parser.go`
  - `pkg/mail/arc/forward_handler.go`
  - `pkg/mail/reputation/scoring_engine.go`
  - `pkg/mail/reputation/circuit_breaker.go`
  - `pkg/mail/reputation/decay_accumulator.go`
  - `pkg/mail/reputation/cross_tenant_gate.go`
  - `cli/cmdmail/arc_cmd.go`
  - `cli/cmdmail/reputation_cmd.go`

---

## 1. Context & Objectives

This subtask delivers enterprise cryptographic ARC (Authenticated Received Chain) signing/validation and cross-tenant reputation feedback circuits:
1. **Authenticated Received Chain (ARC) Protocol Engine (Cluster 61, Items 601–610):** Cryptographic verification of RFC 8617 ARC chains (`ARC-Authentication-Results`, `ARC-Message-Signature`, `ARC-Seal`), forwarder integrity preservation, multi-tenant intermediate hop validation, and automated seal generation.
2. **Cross-Tenant Reputation Circuit Breaker & Penalty Gate (Cluster 62, Items 611–620):** Multi-tenant deliverability scoring, automated ISP bounce surge damping, adaptive quarantine gating, and Bayesian reputation decay algorithms.
3. **Adaptive Sending Schedule & Velocity Shaping (Cluster 63, Items 621–630):** Non-linear warmup ramp algorithms, ISP-specific rate curvature, mailbox dormancy ramp resets, and automated seed-to-target dispatch smoothing.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isValidArcChain`, `hasValidSeal`, `isReputationHealthy`, `isCircuitOpen`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: ARC Cryptographic Engine
- Implement RFC 8617 ARC chain validation in `pkg/mail/arc/chain_validator.go`.
- Implement ARC seal generation and RSA/Ed25519 signature computation in `pkg/mail/arc/seal_signer.go`.
- Construct robust ARC header tokenizers and canonicalization passes in `pkg/mail/arc/header_parser.go`.
- Build intermediate forwarder handling and validation in `pkg/mail/arc/forward_handler.go`.

### Step 2: Reputation Scoring & Circuit Breaker
- Implement real-time tenant reputation scoring models in `pkg/mail/reputation/scoring_engine.go`.
- Build automated multi-state circuit breakers (`Open`, `Half-Open`, `Closed`) in `pkg/mail/reputation/circuit_breaker.go`.
- Construct rolling Bayesian reputation decay accumulators in `pkg/mail/reputation/decay_accumulator.go`.
- Build tenant quota governor and cross-tenant penalty isolation gates in `pkg/mail/reputation/cross_tenant_gate.go`.

### Step 3: Velocity Shaping & Command Integration
- Build non-linear warmup velocity curves and dispatch smoothers in `pkg/mail/circuit/`.
- Wire CLI subcommands in `cli/cmdmail/arc_cmd.go` and `cli/cmdmail/reputation_cmd.go`.
