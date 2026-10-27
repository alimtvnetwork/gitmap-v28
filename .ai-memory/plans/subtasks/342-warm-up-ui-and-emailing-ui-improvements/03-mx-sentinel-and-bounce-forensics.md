# Subtask 03: Real-Time MX DNS Sentinel & Forensic Bounce Analysis

## Status: READY
- **Subtask Slug:** `03-mx-sentinel-and-bounce-forensics`
- **Parent Plan:** `.ai-memory/plans/pending/342-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/04-improvements-561-to-600-security-audit-and-swarm-telemetry.md`
- **Target Range:** Improvements 561–580 (Clusters 57–58)
- **Target Subsystem:** `pkg/mail/mx/`, `pkg/mail/bounce/`
- **Exclusive Owned Files:**
  - `pkg/mail/mx/record_sentinel.go`
  - `pkg/mail/mx/ptr_verifier.go`
  - `pkg/mail/mx/dnsbl_scanner.go`
  - `pkg/mail/mx/spf_inspector.go`
  - `pkg/mail/mx/dkim_checker.go`
  - `pkg/mail/mx/dmarc_auditor.go`
  - `pkg/mail/mx/composite_dashboard.go`
  - `pkg/mail/bounce/dsn_parser.go`
  - `pkg/mail/bounce/semantic_categorizer.go`
  - `pkg/mail/bounce/isp_translator.go`
  - `pkg/mail/bounce/suppression_manager.go`
  - `pkg/mail/bounce/bounce_hud.go`
  - `cli/cmdwarmup/mx_cmd.go`
  - `cli/cmdwarmup/bounce_cmd.go`

---

## 1. Context & Objectives

This subtask delivers carrier-grade DNS health sentinels and forensic bounce diagnostics:
1. **Real-Time MX DNS Health & Reverse IP Blacklist Sentinel (Cluster 57, Items 561–570):** MX record latency monitoring, FCrDNS PTR verification, 50+ DNSBL blacklist queries, SPF/DKIM/DMARC/STS/BIMI verification, and composite health scoring.
2. **Forensic Bounce Analysis Engine & Automated Classification HUD (Cluster 58, Items 571–580):** RFC 3463/3464 DSN parser, Bayesian semantic bounce categorizer, ISP error translator, automated circuit-breaking suppression manager, and live diagnostic HUD.
3. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isReachable`, `isFcrdnsCompliant`, `isPermanentFailure`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: MX Sentinel & DNS Reputation Engine
- Implement MX resolution and latency probe in `pkg/mail/mx/record_sentinel.go`.
- Construct forward-confirmed reverse DNS (FCrDNS) verifier in `pkg/mail/mx/ptr_verifier.go`.
- Implement parallel DNSBL query engine in `pkg/mail/mx/dnsbl_scanner.go`.
- Build SPF, DKIM, DMARC, and BIMI auditors in `pkg/mail/mx/`.
- Aggregate DNS metrics into composite deliverability scorecard in `pkg/mail/mx/composite_dashboard.go`.

### Step 2: Forensic Bounce Parser & Semantic Categorization
- Implement high-speed RFC 3463/3464 DSN parser in `pkg/mail/bounce/dsn_parser.go`.
- Build Bayesian semantic categorizer for non-standard error strings in `pkg/mail/bounce/semantic_categorizer.go`.
- Create ISP-specific error code translator in `pkg/mail/bounce/isp_translator.go`.
- Implement automated circuit-breaking suppression list manager in `pkg/mail/bounce/suppression_manager.go`.

### Step 3: Interactive Diagnostic HUD & CLI Wiring
- Build live forensic bounce HUD in `pkg/mail/bounce/bounce_hud.go`.
- Wire CLI commands into `cli/cmdwarmup/mx_cmd.go` and `cli/cmdwarmup/bounce_cmd.go`.
