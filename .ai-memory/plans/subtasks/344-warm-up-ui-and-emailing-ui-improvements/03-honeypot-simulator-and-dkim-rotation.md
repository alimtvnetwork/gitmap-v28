# Subtask 03: Honeypot Interaction Simulator & Cryptographic DKIM Rotation Sentinel

## Status: READY
- **Subtask Slug:** `03-honeypot-simulator-and-dkim-rotation`
- **Parent Plan:** `.ai-memory/plans/pending/344-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/04-improvements-761-to-800-audit-governance-and-swarm-telemetry.md`
- **Target Range:** Improvements 761–780 (Clusters 77–78)
- **Target Subsystem:** `pkg/mail/honeypot/`, `pkg/mail/weaver/`, `pkg/mail/dkim/`, `pkg/mail/dnscheck/`
- **Exclusive Owned Files:**
  - `pkg/mail/honeypot/trap_detector.go`
  - `pkg/mail/honeypot/heuristic_scorer.go`
  - `pkg/mail/honeypot/quarantine_shield.go`
  - `pkg/mail/weaver/thread_synthesizer.go`
  - `pkg/mail/weaver/context_memory.go`
  - `pkg/mail/weaver/turn_pacer.go`
  - `pkg/mail/dkim/key_rotator.go`
  - `pkg/mail/dkim/dual_selector.go`
  - `pkg/mail/dkim/revocation_circuit.go`
  - `pkg/mail/dnscheck/propagation_hud.go`
  - `pkg/mail/dnscheck/resolver_pool.go`
  - `cli/cmdmail/honeypot_cmd.go`
  - `cli/cmdmail/dkim_cmd.go`

---

## 1. Context & Objectives

This subtask delivers an intelligent synthetic honeypot detection engine, contextual AI thread weaving, and automated cryptographic DKIM key lifecycle management:
1. **Autonomous Honeypot Interaction Simulator & Conversational Thread Weaver (Cluster 77, Items 761–770):** Identifying spam-trap honeypot heuristics, simulated natural engagement responses, multi-turn AI email thread continuity, and adaptive delay pacing.
2. **Cryptographic DKIM Key Rotation Sentinel & DNS Propagation HUD (Cluster 78, Items 771–780):** Zero-downtime dual-selector DKIM key rotation (2048-bit RSA and Ed25519), automated global DNS TXT/CNAME propagation verification HUD, compromised key revocation, and DNSSEC validation.
3. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isHoneypotDetected`, `hasNaturalEngagement`, `isDkimKeyRotated`, `isDnsPropagated`, `hasValidDnssec`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Honeypot Detection & Safe Simulation
- Build spam-trap address heuristic analyzer in `pkg/mail/honeypot/trap_detector.go`.
- Implement honeypot risk scoring accumulator in `pkg/mail/honeypot/heuristic_scorer.go`.
- Construct safe disengagement quarantine shield in `pkg/mail/honeypot/quarantine_shield.go`.

### Step 2: Multi-Turn Conversational Thread Weaving
- Build multi-turn contextual email generator in `pkg/mail/weaver/thread_synthesizer.go`.
- Implement long-context thread memory cache in `pkg/mail/weaver/context_memory.go`.
- Construct realistic inter-turn dispatch delay pacer in `pkg/mail/weaver/turn_pacer.go`.

### Step 3: DKIM Key Rotation Sentinel
- Build dual-selector zero-downtime DKIM rotator in `pkg/mail/dkim/key_rotator.go`.
- Implement active/staging selector coordinator in `pkg/mail/dkim/dual_selector.go`.
- Construct compromised key emergency revocation circuit in `pkg/mail/dkim/revocation_circuit.go`.

### Step 4: DNS Propagation HUD
- Implement terminal DNS propagation HUD in `pkg/mail/dnscheck/propagation_hud.go`.
- Build multi-region authoritative DNS resolver pool in `pkg/mail/dnscheck/resolver_pool.go`.

### Step 5: CLI Command Integration
- Wire honeypot simulation commands in `cli/cmdmail/honeypot_cmd.go`.
- Wire DKIM rotation and DNS check commands in `cli/cmdmail/dkim_cmd.go`.
