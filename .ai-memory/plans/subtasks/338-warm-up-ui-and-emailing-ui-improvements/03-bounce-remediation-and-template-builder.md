# Subtask 03: Autonomous Bounce Remediation & Template Builder

- **Parent Plan:** `.ai-memory/plans/pending/338-warm-up-ui-and-emailing-ui-improvements.md`
- **Spec Path:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/04-improvements-161-to-200-remediation-builder-splitdb.md`
- **Sequence:** 338
- **Subtask Index:** 03
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Objectives

Implement and verify Improvements 161 through 180 spanning Clusters 17 and 18:
1. **Cluster 17: Autonomous Bounce Remediation & Smart Throttling (Items 161–170):**
   - RFC 3463 and RFC 5248 enhanced status code parser and diagnostic classifier.
   - Adaptive per-ISP dynamic token bucket rate limiter and concurrency shaper.
   - Full-jitter exponential backoff retry queue and mailbox quarantine circuit breaker.
   - Automated warm-up cooling ramp recovery scheduler and Bloom filter suppression list.
   - SMTP greylisting detector, connection pool zombie socket pruner, and ARF abuse complaint parser.

2. **Cluster 18: Visual Email Template Studio & Live Component Builder (Items 171–180):**
   - Zero-runtime CSS grid email canvas with drag-and-drop structural blocks.
   - Fluid typography scaling and multi-device viewport emulator (Mobile, Tablet, Desktop).
   - Real-time dark mode email client preview shader with WCAG AA contrast checking.
   - Native MJML-to-HTML compilation engine and dynamic Mustache token validator.
   - Bulletproof VML CTA button generator, preheader zero-width whitespace bumper, and WebP image optimizer.
   - Visual A/B split-testing matrix and SQLite-backed template version control with visual diffing.

## Concrete Deliverables

- `pkg/mail/bounce/status_parser.go`: RFC status parser and semantic classifier.
- `pkg/mail/throttle/isp_shaper.go`: Per-ISP rate limiter and concurrency manager.
- `pkg/mail/circuit/quarantine_breaker.go`: Circuit breaker and automatic mailbox isolator.
- `pkg/web/builder/canvas_engine.go`: Visual drag-and-drop layout generator.
- `pkg/web/builder/dark_mode_shader.go`: Contrast calculator and dark mode preview injector.
- `pkg/mail/template/token_parser.go`: Mustache AST parser and schema validator.
- Static builder styles in `web/static/css/builder.css` (zero external CDNs).

## Quality Standards

- 100% positive booleans (`isThrottled`, `hasValidCanvas`, `isContrastCompliant`, `hasUndefinedVariables`).
- Functions strictly <= 15 lines.
- Complete Split-DB isolation: configuration in `warmup.db`, audit logs in `mail_audit.db`.
