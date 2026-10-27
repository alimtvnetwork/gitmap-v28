# Subtask Plan 03: Autonomous SPF/DKIM/DMARC Record Repair & Forensic Bounce Analysis Engine (Improvements 1561–1580)

- **Parent Task:** `352-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `352-03`
- **Spec Path:** `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/04-improvements-1561-to-1600-audit-governance-and-swarm-telemetry.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (1561–1570)
  - Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD (1571–1580)

## Implementation Steps

1. **Autonomous SPF/DKIM/DMARC DNS Over HTTPS Sentinel (1561–1565):**
   - Query DoH multi-resolvers (Cloudflare, Google, Quad9, AdGuard, OpenDNS) concurrently for authoritative apex records in `pkg/mail/dnscheck/`.
   - Validate full DNSSEC cryptographic trust chains (DS -> DNSKEY -> RRSIG) and NSEC3 authenticated denial of existence.
   - Author RFC 7208 syntax parser with macro expansion tokens and recursive SPF lookup tree flattening (<10 lookups) via CIDR aggregation in `pkg/mail/dnsrepair/`.
   - Maintain dual-selector DKIM key rotators (RSA-2048 & Ed25519) and DMARC alignment policy synthesizers (`p=none` -> `p=quarantine` -> `p=reject`).

2. **Automated DNS Multi-Provider Repair & Drift Sentinel (1566–1570):**
   - Provision DNS records automatically via Cloudflare, Route53, and DigitalOcean APIs with dry-run safety gates and automated rollbacks.
   - Build anycast propagation radars tracking 24 global nodes with real-time drift circuit breakers.
   - Persist DNS audit records in `repodb/pipeline.db` using positive booleans (`isDnssecVerified`, `isConsensusMet`, `hasValidSpf`).

3. **Forensic Bounce Analysis Engine & RFC 3463/5248 Status Codes (1571–1575):**
   - Ingest RFC 3464 Delivery Status Notifications (DSNs) and extract granular enhanced status codes (`5.7.1`, `5.1.1`, `4.4.2`, `5.4.1`) in `pkg/mail/bounce/`.
   - Categorize bounces into discrete failure classes: Hard Bounce, Policy Block, Reputation Throttle, Mailbox Full, Transient DNS.
   - Author interactive Bubble Tea terminal bounce forensics HUD with streaming NDR event logs and Lipgloss color tokens.

4. **Bayesian Spam Filter Forensics & Automated Quarantine Circuit Breaker (1576–1580):**
   - Parse `Authentication-Results`, `X-Spam-Status`, `X-Google-Smtp-Source`, and `X-Forefront-Antispam-Report` diagnostic headers.
   - Score outbound templates against offline Bayesian keyword and MIME anomaly models to detect spam trigger phrases before dispatch.
   - Automatically trip mailbox quarantine circuit breakers when bounce rates exceed 4% or spam placement crosses 2%.
   - Maintain high-speed in-memory Bloom filter suppression lists and dispatch SIEM webhooks (Splunk, Elastic, Datadog).

## Verification Checklist

- [ ] DoH multi-resolver consensus validates SPF/DKIM/DMARC records across 5 independent endpoints.
- [ ] SPF flattener reduces recursive lookups to <= 7 lookups using optimized CIDR blocks.
- [ ] Granular RFC 3463/5248 status decoder correctly maps 100+ standard enhanced error codes.
- [ ] In-memory Bloom filter yields sub-millisecond suppression lookup times.
- [ ] Positive booleans strictly followed throughout.
- [ ] 100% relative repository paths maintained with zero absolute paths.
