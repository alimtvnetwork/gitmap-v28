# Subtask 02: Interactive Bubble Tea Warm-Up Studio & Seed Response Analytics

## Status: READY
- **Subtask Slug:** `02-interactive-bubbletea-and-inbox-analytics`
- **Parent Plan:** `.ai-memory/plans/pending/343-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/343-warm-up-ui-and-emailing-ui-improvements/03-improvements-631-to-660-interactive-bubbletea-and-inbox-analytics.md`
- **Target Range:** Improvements 631–660 (Clusters 64–66)
- **Target Subsystem:** `pkg/cli/tui/`, `pkg/mail/analytics/`, `pkg/mail/inbox/`
- **Exclusive Owned Files:**
  - `pkg/cli/tui/bubbletea_warmup_studio.go`
  - `pkg/cli/tui/keybindings_warmup.go`
  - `pkg/cli/tui/sparkline_renderer.go`
  - `pkg/cli/tui/tab_switcher.go`
  - `pkg/mail/analytics/seed_inbox_tracker.go`
  - `pkg/mail/analytics/sentiment_scorer.go`
  - `pkg/mail/analytics/thread_depth_monitor.go`
  - `pkg/mail/inbox/placement_matrix.go`
  - `pkg/mail/inbox/spam_folder_telemetry.go`
  - `cli/cmdtui/warmup_studio_cmd.go`

---

## 1. Context & Objectives

This subtask delivers an interactive terminal studio built on Bubble Tea, Lipgloss styling tokens, and automated seed response telemetry:
1. **Interactive Bubble Tea Warmup Studio & Real-Time Terminal HUD (Cluster 64, Items 631–640):** High-refresh interactive TUI with multi-tab navigation, dynamic sparkline renderers, real-time dispatch pacing visualizers, and interactive mailbox health inspectors.
2. **Seed Response Sentiment & Engagement Dynamics (Cluster 65, Items 641–650):** Automated seed reply ingestion, thread depth monitoring, conversational sentiment scoring, and automated engagement velocity computation.
3. **Inbox Placement & Spam Folder Forensics (Cluster 66, Items 651–660):** Cross-ISP placement rate matrix (Gmail, Microsoft 365, Yahoo, Proton, iCloud), spam folder detection, delivery delay tracking, and automated placement remediation recommendations.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isInteractiveMode`, `hasHighEngagement`, `isInboxPlaced`, `hasSpamTrigger`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Bubble Tea Terminal Warmup Studio
- Implement primary Tea model, message loops, and viewport renderers in `pkg/cli/tui/bubbletea_warmup_studio.go`.
- Construct ANSI color styles and Lipgloss tokens in `pkg/cli/tui/keybindings_warmup.go`.
- Build unicode sparkline graphs and throughput meters in `pkg/cli/tui/sparkline_renderer.go`.
- Implement zero-allocation tab switching and panel focus managers in `pkg/cli/tui/tab_switcher.go`.

### Step 2: Seed Response Analytics & Sentiment Engine
- Implement IMAP/POP3 seed response ingestion in `pkg/mail/analytics/seed_inbox_tracker.go`.
- Construct lightweight Bayesian sentiment classifier in `pkg/mail/analytics/sentiment_scorer.go`.
- Build conversation thread depth tracking and reply turnaround timers in `pkg/mail/analytics/thread_depth_monitor.go`.

### Step 3: Inbox Placement Matrix & Forensics
- Implement cross-provider inbox placement matrix in `pkg/mail/inbox/placement_matrix.go`.
- Build spam folder arrival detector and ISP penalty classification in `pkg/mail/inbox/spam_folder_telemetry.go`.
- Wire interactive studio command into `cli/cmdtui/warmup_studio_cmd.go`.
