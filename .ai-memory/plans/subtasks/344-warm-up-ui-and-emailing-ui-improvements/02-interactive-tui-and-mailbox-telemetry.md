# Subtask 02: Interactive TUI Cockpit & Mailbox Seed Telemetry

## Status: READY
- **Subtask Slug:** `02-interactive-tui-and-mailbox-telemetry`
- **Parent Plan:** `.ai-memory/plans/pending/344-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/344-warm-up-ui-and-emailing-ui-improvements/02-improvements-731-to-760-interactive-tui-and-mailbox-telemetry.md`
- **Target Range:** Improvements 731–760 (Clusters 74–76)
- **Target Subsystem:** `pkg/cli/tui/`, `pkg/mail/telemetry/`, `pkg/mail/placement/`
- **Exclusive Owned Files:**
  - `pkg/cli/tui/warmup_cockpit.go`
  - `pkg/cli/tui/sparkline_widget.go`
  - `pkg/cli/tui/mailbox_grid.go`
  - `pkg/cli/tui/color_theme.go`
  - `pkg/mail/telemetry/imap_ingester.go`
  - `pkg/mail/telemetry/pop3_poller.go`
  - `pkg/mail/telemetry/sentiment_classifier.go`
  - `pkg/mail/placement/inbox_matrix.go`
  - `pkg/mail/placement/spam_classifier.go`
  - `pkg/mail/placement/reputation_tracker.go`
  - `cli/cmdtui/cockpit_cmd.go`
  - `cli/cmdmail/telemetry_cmd.go`

---

## 1. Context & Objectives

This subtask delivers a high-performance interactive terminal studio and automated seed mailbox intelligence system:
1. **Interactive Bubble Tea Warm-Up Cockpit (Cluster 74, Items 731–740):** Full-screen terminal dashboard with real-time dispatch pacing charts, health dials, animated progress gauges, and Lipgloss HD theme styling.
2. **Automated Seed Mailbox Telemetry & Sentiment Classification (Cluster 75, Items 741–750):** Concurrent IMAP IDLE and POP3 polling, automated multi-folder message ingestion, thread association, and Bayesian sentiment scoring.
3. **Cross-ISP Inbox Placement Matrix & Forensic Folder Tracking (Cluster 76, Items 751–760):** Matrix tracking of Inbox vs Spam vs Promotions folder placement across major providers (Google Workspace, Microsoft 365, Yahoo, Fastmail), tab-separated export, and automated spam rescue workflows.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isInteractiveMode`, `hasSeedResponse`, `isInboxPlaced`, `hasSpamAlert`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Interactive Bubble Tea Cockpit
- Build main Bubble Tea model and message loop in `pkg/cli/tui/warmup_cockpit.go`.
- Implement high-density ANSI sparkline widget in `pkg/cli/tui/sparkline_widget.go`.
- Implement responsive multi-column mailbox grid in `pkg/cli/tui/mailbox_grid.go`.
- Configure Lipgloss color palettes and status tokens in `pkg/cli/tui/color_theme.go`.

### Step 2: Seed Ingestion & Sentiment Analysis
- Build concurrent IMAP IDLE listener in `pkg/mail/telemetry/imap_ingester.go`.
- Implement fallback POP3 poller with backoff in `pkg/mail/telemetry/pop3_poller.go`.
- Construct Bayesian sentiment classification model in `pkg/mail/telemetry/sentiment_classifier.go`.

### Step 3: Inbox Placement Matrix
- Build cross-provider placement matrix aggregator in `pkg/mail/placement/inbox_matrix.go`.
- Implement automated spam folder classification and rescue trigger in `pkg/mail/placement/spam_classifier.go`.
- Build provider-level reputation scoring tracker in `pkg/mail/placement/reputation_tracker.go`.

### Step 4: CLI Integration
- Wire terminal dashboard launch commands in `cli/cmdtui/cockpit_cmd.go`.
- Wire mailbox telemetry commands in `cli/cmdmail/telemetry_cmd.go`.
