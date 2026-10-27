# Subtask 03: Multi-IP Egress Failover, ISP Circuit Breakers & Design System Tokens

## Status: READY
- **Subtask Slug:** `03-failover-and-design-tokens`
- **Parent Plan:** `.ai-memory/plans/pending/339-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/04-improvements-261-to-300-failover-tokens-and-telemetry.md`
- **Owner:** Spec Author Subagent 3 / Network & Design Tokens Engineer
- **Target Subsystem:** `pkg/mail/egress/`, `pkg/mail/circuit/`, `web/static/css/`
- **Exclusive Owned Files:**
  - `pkg/mail/egress/interface_binder.go`
  - `pkg/mail/egress/failover_router.go`
  - `pkg/mail/circuit/isp_breaker.go`
  - `pkg/mail/egress/health_prober.go`
  - `web/static/css/tokens.css`
  - `web/static/css/themes/high_contrast.css`
  - `cli/cmdmail/cmd_egress.go`

---

## 1. Context & Objectives

This subtask implements robust multi-IP network fault tolerance and declarative design token hierarchies:
1. **Multi-IP Egress Pool & Socket Binding (Cluster 27, Items 261–270):** Outbound multi-homing across local network interfaces, continuous latency radar probing, 3-state ISP circuit breakers (`Closed`, `Open`, `Half-Open`), canary probe ramps, and automatic queue failover.
2. **Design System Tokens & Accessibility (Cluster 28, Items 271–280):** Multi-tier CSS custom property dictionary (`--sys-*`, `--comp-*`), fluid typography scales via `clamp()`, WCAG 2.1 AAA high-contrast modes, high-visibility focus indicators, and reduced motion physics adapters.
3. **Code Quality & Invariants:**
   - Every function <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isCircuitClosed`, `hasAvailableFallback`, `isHighContrastActive`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Multi-IP Egress Socket Binder & Health Prober
- Implement `pkg/mail/egress/interface_binder.go` binding TCP dialers to explicit local IP addresses.
- Construct `pkg/mail/egress/health_prober.go` running non-intrusive TCP handshake probes against ISP MX hosts.

### Step 2: 3-State ISP Circuit Breaker & Automatic Failover Router
- Architect `pkg/mail/circuit/isp_breaker.go` implementing `Closed`, `Open`, and `Half-Open` state transitions.
- Build `pkg/mail/egress/failover_router.go` reallocating queued dispatches away from degraded IPs.

### Step 3: Design System Tokens & WCAG 2.1 AAA Theme
- Author `web/static/css/tokens.css` defining glassmorphic obsidian palettes and semantic color tokens.
- Implement `web/static/css/themes/high_contrast.css` with 7:1 contrast ratios and solid boundaries.
- Build CLI command suite in `cli/cmdmail/cmd_egress.go` supporting `gitmap mail egress list` and failover simulation.

---

## 3. Verification & Acceptance Criteria
- [ ] Outbound TCP dialer binds cleanly to assigned secondary IP interface with zero socket leaks.
- [ ] Tripped circuit breaker shifts to `Half-Open` after cooldown and recovers upon 10 successful canary deliveries.
- [ ] High-contrast mode passes automated WCAG 2.1 AAA contrast audit (>7:1 ratio on all text).
