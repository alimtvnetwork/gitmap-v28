# Subtask Plan 02: Interactive TUI Cockpit & Seed Intelligence Network (Improvements 931–960)

- **Parent Task:** `346-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `346-02`
- **Spec Path:** `02-spec/21-app/346-warm-up-ui-and-emailing-ui-improvements/03-improvements-931-to-960-interactive-tui-and-seed-intelligence.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 94: High-Density Bubble Tea Cluster Operations Cockpit (931–940)
  - Cluster 95: Synthetic Multi-Persona Dialog Generation & Context Injection (941–950)
  - Cluster 96: Real-Time Cross-ISP Inbox Placement Matrix & Reputation Vectors (951–960)

## Implementation Steps

1. **Bubble Tea Operations Cockpit & Visualizers (931–940):**
   - Build high-density multi-pane Bubble Tea terminal monitoring application in `pkg/cli/tui/`.
   - Incorporate Lipgloss color schemes (Dracula/Tokyo Night accents), ANSI sparkline telemetry, and live filter models.
   - Enforce zero external CSS/CDN dependencies and non-blocking event loops.

2. **Synthetic Seed Persona Dialog Generator (941–950):**
   - Implement realistic email thread generators simulating natural correspondence threads in `pkg/mail/seed/`.
   - Inject domain-specific context, simulated meeting requests, attachments, and unsubscribe links.
   - Enforce positive booleans (`isThreadedReply`, `hasAttachment`, `isPersonaActive`).

3. **Reputation Vectors & ISP Placement Matrix (951–960):**
   - Aggregate inbox placement signals across major mailbox providers (Google, Microsoft 365, Yahoo, Proton) in `pkg/mail/placement/`.
   - Generate multi-dimensional reputation vector heatmaps using ANSI true-color blocks in terminal output.
   - Persist placement snapshots into `repodb/pipeline.db` with positive booleans (`isInboxDelivered`, `hasSpamPlacement`).

## Verification Checklist

- [ ] Bubble Tea cockpit renders fluidly without screen flicker during rapid 100ms state updates.
- [ ] Synthetic dialog generator creates 50 continuous thread replies passing natural language heuristic filters.
- [ ] Cross-ISP matrix renders accurate placement ratios across 4 simulated provider endpoints.
