# Subtask Plan 02: Interactive TUI Cockpit & Mailbox Diagnostics (Improvements 1131–1160)

- **Parent Task:** `348-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `348-02`
- **Spec Path:** `02-spec/21-app/348-warm-up-ui-and-emailing-ui-improvements/03-improvements-1131-to-1160-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 114: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (1131–1140)
  - Cluster 115: Dynamic Backoff Inversion & Multi-Domain Concurrency Governor (1141–1150)
  - Cluster 116: Zero-CDN Responsive Email Content Security Sandbox & HTML Sanitizer HUD (1151–1160)

## Implementation Steps

1. **Bubble Tea Seed Health Matrix & Placement Heatmaps (1131–1140):**
   - Implement interactive terminal seed mailbox verification matrix categorizing Google, Microsoft 365, Yahoo, Proton, and Apple iCloud seeds in `pkg/cli/tui/`.
   - Implement real-time IMAP IDLE arrival synchronizer and automated spam folder rescue daemons in `pkg/mail/seed/`.
   - Track Feedback Loop (FBL) abuse complaints and RBL blacklists with instant alert toasts.
   - Enforce positive booleans (`isActive`, `isHealthy`, `hasRecentRescue`).

2. **Dynamic Backoff Inversion & Concurrency Governor (1141–1150):**
   - Implement biological leaky-integrate-and-fire (LIF) throttle models inverting dispatch rate curves upon detecting 4xx transient deferrals in `pkg/mail/throttle/`.
   - Maintain multi-domain concurrency limiters preventing IP reputation burning on destination MX hosts.
   - Store domain rate governors in `repodb/warmup.db` with positive booleans (`isThrottled`, `hasCapacity`).

3. **Content Security Sandbox & HTML Sanitizer HUD (1151–1160):**
   - Provide offline HTML content sanitization verifying CSS containment, tracking pixel integrity, and dark mode readability in `pkg/mail/sanitize/`.
   - Render terminal preview of email HTML structure using Lipgloss ASCII formatting.
   - Enforce zero external CDN asset dependencies and strict RFC 5322 compliance.

## Verification Checklist

- [ ] Bubble Tea seed matrix renders 300+ seed inboxes fluidly at 60 FPS without terminal flicker.
- [ ] IMAP IDLE daemon detects seed arrivals within 500ms and triggers automated engagement workflows.
- [ ] Concurrency governor restricts concurrent domain connections strictly under MX host limits.
