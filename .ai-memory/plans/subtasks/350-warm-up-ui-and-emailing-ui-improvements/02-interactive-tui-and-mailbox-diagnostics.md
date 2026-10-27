# Subtask Plan 02: Interactive Bubble Tea Mailbox Health Matrix & Seed Pool Diagnostics (Improvements 1331–1360)

- **Parent Task:** `350-warm-up-ui-and-emailing-ui-improvements`
- **Subtask ID:** `350-02`
- **Spec Path:** `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/03-improvements-1331-to-1360-interactive-tui-and-mailbox-diagnostics.md`
- **Status:** PLANNED
- **Cluster Scope:**
  - Cluster 134: Interactive Bubble Tea Seed Pool Verification Matrix & Real-Time Sync (1331–1340)
  - Cluster 135: Mailbox Health Diagnostics, Spam Rescue & Reputation Auditing (1341–1350)
  - Cluster 136: Deliverability Intelligence, Forensic Reports & FBL Loop Interceptor (1351–1360)

## Implementation Steps

1. **Interactive Seed Pool Verification Matrix (1331–1340):**
   - Author multi-provider seed status matrix in Bubble Tea v2 (`pkg/cli/tui/seed_matrix.go`).
   - Implement real-time IMAP IDLE event synchronizer and live seed delivery listeners.
   - Author synthetic conversational turn dispatchers simulating natural peer engagement.
   - Record seed interaction events in Split-DB (`repodb/pipeline.db`) using positive booleans (`isSeedActive`, `hasReplied`).

2. **Mailbox Diagnostics & Spam Rescue Engine (1341–1350):**
   - Author automated IMAP spam folder rescue daemon moving misclassified warm-up emails to inbox.
   - Implement automated mark-as-important and message starring workflows to train ISP classifiers.
   - Check seed IP addresses against major DNSBL/RBL blacklists with automatic exclusion guards.
   - Implement OAuth2 token refresh rotators with AES-256-GCM encrypted credential vault.

3. **Deliverability Intelligence & FBL Abuse Interceptor (1351–1360):**
   - Parse RFC 6591 ARF (Authentication Failure Reporting Format) forensic incident reports.
   - Ingest Feedback Loop (FBL) complaint notifications across major mailbox providers (Gmail, Outlook, Yahoo).
   - Author automated recipient suppression and tenant warmup throttling circuit breakers.
   - Build Lipgloss terminal deliverability scorecard and seed health distribution charts.

## Verification Checklist

- [ ] IMAP IDLE listener detects incoming warm-up email within 2 seconds of arrival.
- [ ] Spam rescue daemon safely moves misclassified email to inbox and flags with positive boolean `isRescued`.
- [ ] RBL monitor checks 15 DNSBL lists and alerts within 500ms of listing detection.
- [ ] Bubble Tea mailbox matrix navigates smoothly via arrow keys without UI stutter.
