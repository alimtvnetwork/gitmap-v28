# Subtask Plan 02: Interactive TUI Cockpit & Seed Intelligence Network (Improvements 1031–1060)

- **Parent Task:** `347-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `347-02`
- **Spec Path:** `02-spec/21-app/347-warm-up-ui-and-emailing-ui-improvements/03-improvements-1031-to-1060-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 104: High-Density Bubble Tea Cluster Operations Cockpit (1031–1040)
  - Cluster 105: Synthetic Multi-Persona Dialog Generation & Context Injection (1041–1050)
  - Cluster 106: Real-Time Cross-ISP Inbox Placement Matrix & Reputation Vectors (1051–1060)

## Implementation Steps

1. **Bubble Tea Operations Cockpit & Visualizers (1031–1040):**
   - Build high-density multi-pane Bubble Tea terminal monitoring application in `pkg/cli/tui/`.
   - Incorporate Lipgloss color schemes (Dracula/Tokyo Night accents), ANSI sparkline telemetry, and live filter models.
   - Enforce zero external CSS/CDN dependencies and non-blocking event loops.

2. **Synthetic Seed Persona Dialog Generator (1041–1050):**
   - Implement realistic email thread generators simulating natural correspondence threads in `pkg/mail/seed/`.
   - Inject domain-specific context, simulated meeting requests, attachments, and unsubscribe links.
   - Enforce positive booleans (`isThreadedReply`, `hasAttachment`, `isPersonaActive`).

3. **Reputation Vectors & ISP Placement Matrix (1051–1060):**
   - Aggregate inbox placement signals across major mailbox providers (Google, Microsoft 365, Yahoo, Proton) in `pkg/mail/placement/`.
   - Generate multi-dimensional reputation vector heatmaps using ANSI true-color blocks in terminal output.
   - Persist placement snapshots into `repodb/pipeline.db` with positive booleans (`isInboxDelivered`, `hasSpamPlacement`).

## Verification Checklist

- [ ] Bubble Tea cockpit renders fluidly without screen flicker during rapid 100ms state updates.
- [ ] Synthetic dialog generator creates 50 continuous thread replies passing natural language heuristic filters.
- [ ] Cross-ISP matrix renders accurate placement ratios across 4 simulated provider endpoints.
