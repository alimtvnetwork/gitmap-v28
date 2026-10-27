# Subtask Plan 03: Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 861–880)

- **Parent Task:** `345-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `345-03`
- **Spec Path:** `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/04-improvements-861-to-900-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 87: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (861–870)
  - Cluster 88: Forensic Bounce Analysis Engine & Automated Classification HUD (871–880)

## Implementation Steps

1. **DNS Over HTTPS (DoH) Sentinel & DNSSEC Chain Validator (861–865):**
   - Implement concurrent DoH queries across Cloudflare, Google, and Quad9 upstream resolvers in `pkg/mail/dnscheck/`.
   - Validate DNSSEC cryptographic signature chains (RRSIG, DNSKEY, DS) to detect cache poisoning and MitM.
   - Design automated diff analyzer identifying syntax errors, lookup limits (>10 DNS lookups in SPF), and missing alignment tags.

2. **Autonomous DNS Record Repair & Multi-Provider API Provisioning (866–870):**
   - Synthesize optimized SPF records with automated IP CIDR collapsing and `include` flattening in `pkg/mail/dnsrepair/`.
   - Generate dual-selector DKIM keys (RSA-2048 & Ed25519) and author DMARC policies (`p=reject; sp=reject; pct=100`).
   - Implement automated API patch drivers for Cloudflare, Route53, and DigitalOcean with rollback dry-run safety.

3. **Forensic SMTP Bounce Parsing & DSN RFC 3463/5248 Attribution (871–875):**
   - Parse RFC 3464 Delivery Status Notifications (DSN) and extract enhanced status codes (`5.7.1`, `5.1.1`, `4.4.2`) in `pkg/mail/bounce/`.
   - Classify bounces into granular categories: Hard Bounce, Policy Block, Reputation Throttle, Mailbox Full, Transient DNS.
   - Record forensic telemetry in Split-DB (`repodb/pipeline.db`) with strictly positive booleans (`isHardBounce`, `isPolicyBlock`).

4. **Bayesian Spam Filter Attribution & Quarantine Remediation (876–880):**
   - Extract spam filter diagnostic headers (`X-Spam-Status`, `Authentication-Results`, `X-Google-Smtp-Source`).
   - Train lightweight heuristic analyzer to identify spam trigger phrases, MIME format mismatches, and broken message encodings.
   - Trigger immediate circuit breaking and warmup mailbox quarantine when bounce bursts cross safety thresholds.

## Verification Checklist

- [ ] DoH sentinel resolves DNSSEC chains and flags simulated spoofed records without crashing.
- [ ] Record repair engine successfully generates flattened SPF records staying strictly below the 10-lookup threshold.
- [ ] Bounce parser classifies 100 sample DSN transcripts with 100% accuracy into correct categorical buckets.
- [ ] Quarantine circuit trips within 1 cycle upon detecting >5% hard bounce rate.
