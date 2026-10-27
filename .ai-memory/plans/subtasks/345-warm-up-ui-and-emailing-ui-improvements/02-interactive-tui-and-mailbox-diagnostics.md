# Subtask Plan 02: Interactive TUI & Mailbox Diagnostics Cockpit (Improvements 831–860)

- **Parent Task:** `345-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `345-02`
- **Spec Path:** `02-spec/21-app/345-warm-up-ui-and-emailing-ui-improvements/03-improvements-831-to-860-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 84: Interactive Bubble Tea Warm-Up Operations Cockpit (831–840)
  - Cluster 85: Multi-Protocol Seed Mailbox Telemetry & Reply Harvester (841–850)
  - Cluster 86: Bayesian Sentiment Analysis & Cross-ISP Placement Heatmaps (851–860)

## Implementation Steps

1. **Bubble Tea Full-Screen TUI Cockpit (831–840):**
   - Implement responsive multi-panel Bubble Tea terminal interface in `pkg/cli/tui/`.
   - Incorporate Lipgloss color palettes (Nord/Tokyo Night inspired), sparkline velocity graphs, and status badges.
   - Enforce zero external CSS/web dependencies and non-blocking keymap event loops.

2. **Seed Mailbox Harvester & Protocol Adapter (841–850):**
   - Implement asynchronous IMAP4rev1 and POP3 seed account pollers in `pkg/mail/diagnostics/`.
   - Extract raw RFC 822 email responses, delivery delays, and folder destinations (Inbox, Spam, Promotions).
   - Ensure clean goroutine pooling with cancellation context propagation.

3. **Bayesian Sentiment & Placement Matrix HUD (851–860):**
   - Build lightweight naive Bayesian classifier for automated warmup reply tone evaluation.
   - Render cross-ISP placement heatmaps in CLI via ANSI true-color blocks and ASCII grids.
   - Store placement snapshots in `repodb/pipeline.db` using positive booleans (`isInboxDelivered`, `hasSpamFlag`).

## Verification Checklist

- [ ] Interactive TUI resizes dynamically without terminal visual artifacting across standard terminal dimensions.
- [ ] IMAP harvester parses multi-part MIME responses and correctly tags folder destinations without memory bloat.
- [ ] Bayesian classifier calculates placement scores with >95% accuracy on test corpus.
