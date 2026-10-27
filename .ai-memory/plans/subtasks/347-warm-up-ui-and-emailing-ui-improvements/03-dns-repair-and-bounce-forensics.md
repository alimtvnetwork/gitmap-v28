# Subtask Plan 03: Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 1061–1080)

- **Parent Task:** `347-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `347-03`
- **Spec Path:** `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/04-improvements-1061-to-1100-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 107: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (961–970)
  - Cluster 108: Forensic Bounce Analysis Engine & Automated Classification HUD (971–980)

## Implementation Steps

1. **DNS Over HTTPS Sentinel & Cryptographic DNSSEC Audit (1061–1065):**
   - Concurrently query DoH resolvers (Cloudflare, Google, Quad9) for authoritative apex records in `pkg/mail/dnscheck/`.
   - Validate full DNSSEC trust anchors (DS -> DNSKEY -> RRSIG) to block DNS spoofing and cache poisoning.
   - Detect syntax errors, multiple SPF records, and recursive DNS lookup overflows (>10 lookups).

2. **Autonomous DNS Record Repair & Multi-Provider API Provisioning (1066–1070):**
   - Synthesize optimized SPF records with automated IP CIDR collapsing and `include` flattening in `pkg/mail/dnsrepair/`.
   - Generate dual-selector DKIM keys (RSA-2048 & Ed25519) and configure strict DMARC enforcement policies (`p=reject; sp=reject; pct=100`).
   - Implement API drivers for Cloudflare, Route53, and DigitalOcean with dry-run safety gates and automated rollback backups.

3. **Forensic Bounce Analysis & RFC 3473/5248 Enhanced Status Codes (1071–1075):**
   - Parse incoming RFC 3474 DSN messages and extract enhanced status codes (`5.7.1`, `5.1.1`, `4.4.2`, `5.4.1`) in `pkg/mail/bounce/`.
   - Categorize bounces into discrete failure classes: Hard Bounce, Policy Block, Reputation Throttle, Mailbox Full, Transient DNS.
   - Store bounce telemetry in Split-DB (`repodb/pipeline.db`) using positive booleans (`isHardBounce`, `isPolicyBlock`, `isReputationThrottle`).

4. **Bayesian Spam Filter Attribution & Mailbox Quarantine Circuit Breaker (1076–1080):**
   - Extract spam diagnostic headers (`X-Spam-Status`, `Authentication-Results`, `X-Google-Smtp-Source`, `X-Forefront-Antispam-Report`).
   - Implement Bayesian keyword and MIME anomaly scoring to identify trigger phrases before dispatch.
   - Automatically trip quarantine circuit breakers when bounce rate exceeds 4% or spam placement crosses 2%.

## Verification Checklist

- [ ] DoH multi-resolver consensus validates SPF/DKIM/DMARC records across 3 independent endpoints.
- [ ] SPF flattener reduces 14 recursive lookups to <= 7 lookups using optimized CIDR blocks.
- [ ] Bounce parser categorizes 150 diverse RFC 3474 test reports with 100% categorical precision.
- [ ] Quarantine circuit breaker pauses warmup pipeline within 1 cycle when hard bounce rate exceeds 4%.
