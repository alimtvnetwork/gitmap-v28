# Subtask Plan 03: Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 1261–1280)

- **Parent Task:** `349-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `349-03`
- **Spec Path:** `02-spec/21-app/349-warm-up-ui-and-emailing-ui-improvements/04-improvements-1261-to-1300-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 127: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (1261–1270)
  - Cluster 128: Forensic Bounce Analysis Engine & Automated Classification HUD (1271–1280)

## Implementation Steps

1. **DNS Over HTTPS Sentinel & Cryptographic DNSSEC Audit (1261–1265):**
   - Concurrently query DoH resolvers (Cloudflare, Google, Quad9, AdGuard, OpenDNS) for authoritative apex records in `pkg/mail/dnscheck/`.
   - Validate full DNSSEC trust anchors (DS -> DNSKEY -> RRSIG) to block DNS spoofing and cache poisoning.
   - Detect syntax errors, multiple SPF records, and recursive DNS lookup overflows (>10 lookups).
   - Generate dual-selector DKIM keys (RSA-2048 & Ed25519) in `pkg/mail/crypto/` with positive booleans (`isActive`, `isEd25519Dual`, `hasDnsPublished`).

2. **Autonomous DNS Record Repair & Multi-Provider API Provisioning (1266–1270):**
   - Synthesize optimized SPF records with automated IP CIDR collapsing and `include` flattening in `pkg/mail/dnsrepair/`.
   - Author DMARC alignment policy synthesis and stepped enforcement schedules (`p=none` -> `p=quarantine` -> `p=reject`).
   - Implement API drivers for Cloudflare, Route53, and DigitalOcean with dry-run safety gates and automated rollback backups.
   - Track global DNS propagation across 24 anycast nodes with drift detection circuit breakers.

3. **Forensic Bounce Analysis & RFC 3463/5248 Enhanced Status Codes (1271–1275):**
   - Parse incoming RFC 3464 DSN messages and extract enhanced status codes (`5.7.1`, `5.1.1`, `4.4.2`, `5.4.1`) in `pkg/mail/bounce/`.
   - Categorize bounces into discrete failure classes: Hard Bounce, Policy Block, Reputation Throttle, Mailbox Full, Transient DNS.
   - Author interactive Bubble Tea terminal bounce forensics HUD with real-time streaming NDR logs.
   - Store bounce telemetry in Split-DB (`repodb/pipeline.db`) using positive booleans (`isHardBounce`, `isPolicyBlock`, `isReputationThrottle`).

4. **Bayesian Spam Filter Attribution & Mailbox Quarantine Circuit Breaker (1276–1280):**
   - Extract spam diagnostic headers (`X-Spam-Status`, `Authentication-Results`, `X-Google-Smtp-Source`, `X-Forefront-Antispam-Report`).
   - Implement offline Bayesian keyword and MIME anomaly scoring to identify trigger phrases before dispatch.
   - Automatically trip quarantine circuit breakers when bounce rate exceeds 4% or spam placement crosses 2%.
   - Sync global suppression lists using high-speed in-memory Bloom filters and forward SIEM webhooks.

## Verification Checklist

- [ ] DoH multi-resolver consensus validates SPF/DKIM/DMARC records across 5 independent endpoints.
- [ ] SPF flattener reduces 14 recursive lookups to <= 7 lookups using optimized CIDR blocks.
- [ ] Bounce parser categorizes 150 diverse RFC 3464 test reports with 100% categorical precision.
- [ ] Quarantine circuit breaker pauses warmup pipeline within 1 cycle when hard bounce rate exceeds 4%.
