# Subtask Plan 02: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Improvements 1431–1460)

- **Parent Task:** `351-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `351-02`
- **Spec Path:** `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/03-improvements-1431-to-1460-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 144: Interactive Seed Mailbox Health Matrix & Real-Time IMAP IDLE Synchronizer (1431–1440)
  - Cluster 145: Autonomous Spam Folder Rescue Daemon & Deliverability Forensics HUD (1441–1450)
  - Cluster 146: OAuth2 Credential Rotation Engine & Mailbox Storage Quota Purger (1451–1460)

## Implementation Steps

1. **Seed Mailbox Health Matrix & IMAP IDLE Sync (1431–1440):**
   - Author Bubble Tea interactive table displaying seed mailboxes across Gmail, Outlook, Yahoo, and custom IMAP nodes in `pkg/cli/tui/`.
   - Maintain persistent background IMAP IDLE socket connections to detect inbound email arrivals in sub-second time.
   - Synchronize seed arrival signals with central Split-DB (`installation.db` and `repodb/pipeline.db`).

2. **Automated Spam Rescue Daemon & Forensic Header Audit (1441–1450):**
   - Implement autonomous spam rescue daemon that checks `Junk`/`Spam` folders, moves messages to `Inbox`, and marks them as important.
   - Parse `Authentication-Results` and `X-Spam-Status` headers from seed inboxes, storing per-ISP scores.
   - Author deliverability forensics HUD with visual inbox placement percentages per domain.

3. **OAuth2 Rotation & Mailbox Maintenance Daemon (1451–1460):**
   - Build automated OAuth2 refresh token rotation daemon for Google Workspace and Microsoft Graph accounts.
   - Schedule non-destructive mailbox storage quota purgers to delete test warmup messages older than 14 days.
   - Intercept Feedback Loop (FBL) complaints and trigger immediate suppression.

## Verification Checklist

- [ ] Bubble Tea mailbox matrix reflects real-time status with zero flickering.
- [ ] IMAP IDLE notifications wake worker loops in under 500ms.
- [ ] Spam rescue daemon successfully un-marks spam messages and marks them important.
- [ ] Positive booleans strictly followed throughout.
