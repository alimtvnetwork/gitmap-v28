# Subtask 02: Multi-Mailbox Management, Outbound Composer Studio & Bounce Monitors

## Status: PLANNED
- **Subtask Slug:** `02-composer-studio-and-multimailbox-hud`
- **Parent Plan:** `.ai-memory/plans/pending/337-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Owner:** Frontend & Studio Subagent
- **Target Subsystem:** `pkg/mail/mailbox/`, `pkg/web/composer/`, `pkg/mail/telemetry/`
- **Improvements Covered:** 31 to 50 (Clusters 4 and 5)
- **Exclusive Owned Files:**
  - `pkg/mail/mailbox/pool_manager.go`
  - `pkg/mail/mailbox/rotation_scheduler.go`
  - `pkg/mail/mailbox/types.go`
  - `pkg/web/composer/studio.go`
  - `pkg/web/composer/preview_renderer.go`
  - `pkg/mail/telemetry/bounce_analyzer.go`
  - `pkg/mail/telemetry/event_stream.go`

---

## 1. Executive Summary & Context

Subtask 02 addresses the operations and user-facing authoring capabilities of the GitMap email engine:
1. **Multi-Mailbox Management & Routing (Improvements 31–40):** Health card decks, round-robin rotation algorithms, per-mailbox sending limits, credential status indicators, and automated mailbox failover.
2. **Outbound Composer Studio & Bounce Telemetry (Improvements 41–50):** Markdown-to-HTML email composition with variable interpolation (`{{firstName}}`, `{{company}}`), dark/light split preview, spam word highlighting, and classified bounce telemetry (Hard, Soft, Policy, Challenge).

---

## 2. File Modification & Structural Ownership

```text
pkg/mail/
├── mailbox/
│   ├── pool_manager.go         [NEW] - Multi-account pool coordinator and health registry
│   ├── rotation_scheduler.go   [NEW] - Round-robin, weighted, and lowest-latency dispatch
│   └── types.go                [NEW] - MailboxConfig, AccountStatus, and Credentials types
├── telemetry/
│   ├── bounce_analyzer.go      [NEW] - DSN (Delivery Status Notification) parser & classifier
│   └── event_stream.go         [NEW] - Real-time SSE / WebSocket event dispatcher
pkg/web/composer/
├── studio.go                   [NEW] - Email composition template engine and variable parser
└── preview_renderer.go         [NEW] - Split-pane dark/light email rendering generator
```

---

## 3. Technical Requirements & Architecture

### 3.1 Multi-Mailbox Pool & Rotation
- Supports multi-tenant account pooling with independent SMTP/IMAP credentials.
- Rotation strategies:
  - `RotationRoundRobin`: Sequential even distribution across healthy accounts.
  - `RotationWeighted`: Proportionate sending based on mailbox warm-up age.
  - `RotationLeastLoaded`: Routes to the mailbox with highest remaining daily capacity.
- Automatic account quarantine if consecutive SMTP connection failures exceed 3.

### 3.2 Composer Studio & Variable Interpolation
- Clean, zero-CDN web composer supporting Markdown input with live HTML output.
- Variable syntax: `{{variableName}}` with strict sanitization and missing key warnings.
- Dark and light mode email preview rendering simulating Outlook, Apple Mail, and Gmail.
- Inline spam token highlighter detecting high-risk phrases (e.g. "guarantee", "free gift", "act now").

### 3.3 Bounce Categorization Engine
- Parses RFC 3463 and RFC 1893 status codes (e.g. `5.1.1` User Unknown, `4.2.2` Mailbox Full, `5.7.1` Policy Rejection).
- Auto-suppresses hard bounce addresses in the SQLite suppression database.
- Emits real-time event packets for UI notification toasts.

---

## 4. Verification & Testing Checklist

- [ ] Unit tests in `pkg/mail/mailbox/rotation_scheduler_test.go` verifying balanced account rotation.
- [ ] Composer tests verifying markdown compilation and variable interpolation safety.
- [ ] Bounce analyzer tests validating categorization of 20+ standard SMTP error strings.
- [ ] Code limit checks: <= 15 line functions, positive booleans (`isQuarantined`, `hasCapacity`).
