# Warm-Up UI & Emailing UI Improvements: Items 31 to 60
## Modern Mail Composer Studio, Identity Grid HUD, and High-Velocity Lipgloss Terminal UX

This specification covers Improvements 31 through 60 for the GitMap Emailing & Warm-Up ecosystem across `go-email-sender`, `go-email-reader`, and `gitmap mail/warmup` CLI tooling.

---

## Cluster 4: Modern Mail Composer Studio & Live Dual-Pane Preview (Items 31–40)

### Item 031: Live Markdown-to-MIME Dual-Pane Interactive Split View
- **ID & Name**: `IMPR-031-COMPOSER-DUAL-PANE-PREVIEW`
- **Problem Statement**: Email authors cannot preview how raw markdown or templated text transforms into final multipart/alternative HTML and plain text bodies in real time, leading to broken tables and unrendered variables.
- **Concrete UI Specification**:
  - Split-pane layout using CSS Flexbox (`display: flex; height: calc(100vh - 56px);`) with responsive horizontal drag divider (`cursor: col-resize; width: 4px; background: var(--border-subtle);`).
  - Left pane houses the Monaco/CodeMirror editor; right pane houses an isolated sandboxed iframe (`sandbox="allow-same-origin"`).
  - Synchronized scroll physics: smooth debounced (16ms) scroll lock matching scroll percentages between editor and rendered DOM.
  - Zero-CDN local CSS rendering with offline typography fallback (`system-ui, -apple-system, sans-serif`).
  - Shortcut: `Ctrl+P` (or `Cmd+P`) toggles dual-pane vs single-pane editor mode.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail compose --preview --split-view --draft-id="draft_9841"
  ```
- **Data Contract**:
  ```json
  {
    "composerId": "comp_031_split",
    "isDualPaneActive": true,
    "hasSynchronizedScrolling": true,
    "isHtmlSanitized": true,
    "hasMarkdownSyntaxHighlighting": true,
    "debounceDelayMs": 16,
    "viewportRatio": 0.50
  }
  ```
- **Acceptance Criteria**:
  1. Editing markdown on the left pane updates the right iframe DOM within 30ms without full-page reloads.
  2. Pressing `Ctrl+P` switches view modes between split, editor-only, and preview-only without losing cursor focus.

---

### Item 032: Real-Time Dark/Light Theme Email Render Matrix Switcher
- **ID & Name**: `IMPR-032-COMPOSER-THEME-RENDER-SWITCHER`
- **Problem Statement**: Email clients (Outlook, Apple Mail, Gmail) automatically invert colors or apply dark mode styling, causing unreadable dark text on dark backgrounds or invisible logos.
- **Concrete UI Specification**:
  - Dual-toggle pill bar embedded in the preview pane header (`background: rgba(30, 30, 40, 0.85); border-radius: 9999px; padding: 2px 6px;`).
  - Instant CSS filter and media query emulation (`@media (prefers-color-scheme: dark)`) toggled via dataset attribute `data-client-theme="dark|light"`.
  - Contrast warning chip that turns amber if contrast ratio drops below WCAG AA (4.5:1) for any text block.
  - Shortcut: `Ctrl+Shift+D` cycles between Light, Inverted Dark, and Native Dark modes.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail render --draft-id="draft_9841" --client-theme="dark" --check-contrast
  ```
- **Data Contract**:
  ```json
  {
    "draftId": "draft_9841",
    "isDarkModeEmulated": true,
    "hasWcagContrastPassed": true,
    "hasInvertedLogoDetected": false,
    "targetColorScheme": "dark",
    "contrastRatioMin": 4.85
  }
  ```
- **Acceptance Criteria**:
  1. Switching to dark mode reflects CSS variables immediately in the preview sandbox without re-rendering network assets.
  2. If contrast falls below 4.5:1, an amber highlight border pulses on the failing element.

---

### Item 033: Context-Aware Dynamic Liquid/Mustache Merge-Tag Auto-Complete
- **ID & Name**: `IMPR-033-COMPOSER-MERGE-TAG-AUTOCOMPLETE`
- **Problem Statement**: Warm-up campaigns and transactional templates use placeholder variables (`{{recipient.name}}`, `{{account.company}}`) that users mistype, producing ugly raw curly braces in production inboxes.
- **Concrete UI Specification**:
  - Popup autocomplete dropdown triggered on `{{` typed inside the composer.
  - Positioned dynamically using caret coordinates (`transform: translate3d(x, y, 0); z-index: 1000;`).
  - Searchable list of tokens with type badge (e.g., `string`, `number`, `fallback`) and sample value preview in muted text (`var(--text-muted)`).
  - Arrow key navigation (`ArrowDown`, `ArrowUp`, `Enter`, `Tab` to insert, `Escape` to dismiss).
- **CLI Counterpart Command**:
  ```bash
  gitmap mail tags list --schema="warmup-v1" --output=json
  ```
- **Data Contract**:
  ```json
  {
    "tagPrefix": "{{",
    "isOpen": true,
    "hasFallbackConfigured": true,
    "isValidSyntax": true,
    "selectedTag": "recipient.firstName",
    "availableTagsCount": 18
  }
  ```
- **Acceptance Criteria**:
  1. Typing `{{` immediately opens the popover with zero latency (< 10ms).
  2. Selecting an entry replaces `{{` with `{{recipient.firstName | default: "there"}}`.

---

### Item 034: Zero-Dependency Drag-and-Drop Raw Attachment Tray with Virus Scanning Badge
- **ID & Name**: `IMPR-034-COMPOSER-ATTACHMENT-TRAY`
- **Problem Statement**: Adding attachments often lacks clear visual progress, mime-type validation, or virus-check feedback, leading to accidental large attachment failures during SMTP transit.
- **Concrete UI Specification**:
  - Dropzone strip pinned directly below subject line (`border: 1px dashed var(--border-focus); border-radius: 8px; min-height: 44px;`).
  - Drag-over state triggers smooth glow animation (`box-shadow: 0 0 0 2px var(--brand-accent); transition: 0.15s ease-in-out;`).
  - File chips display icon, truncated filename, file size in KB/MB, and positive status badge (`isScanned: true`, `hasValidMime: true`).
  - Shortcut: `Ctrl+O` opens local file browser dialog for rapid attachment.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail draft attach --draft-id="draft_9841" --file="./reports/summary.pdf"
  ```
- **Data Contract**:
  ```json
  {
    "attachmentId": "att_5521",
    "isScanned": true,
    "isMalwareFree": true,
    "hasAllowedExtension": true,
    "isSizeWithinLimit": true,
    "fileSizeBytes": 452810,
    "mimeType": "application/pdf"
  }
  ```
- **Acceptance Criteria**:
  1. Dropping a file larger than 25MB displays an instant red banner blocking the send action.
  2. Removing a chip instantly updates total email payload size calculation.

---

### Item 035: Granular Email Size & Spam Trigger Score Live Telemetry Bar
- **ID & Name**: `IMPR-035-COMPOSER-SPAM-SIZE-TELEMETRY`
- **Problem Statement**: Overweight emails (> 102KB) trigger Gmail clipping, and spam-trigger words ("100% free", "urgent wire") hurt warm-up deliverability without writer warning.
- **Concrete UI Specification**:
  - Bottom docked telemetry status bar (`height: 28px; background: var(--bg-surface-elevated); font-family: var(--font-mono); font-size: 11px;`).
  - Live byte counter with dual visual alerts: Green (< 75KB), Amber (75KB - 100KB), Red (> 102KB "Gmail Clipping Risk").
  - Real-time SpamAssassin / NLP score meter ranging from 0.0 to 10.0 with clickable chip to highlight offending phrases.
  - Shortcut: `Ctrl+Shift+S` highlights all flagged spam tokens in yellow background (`rgba(245, 158, 11, 0.2)`).
- **CLI Counterpart Command**:
  ```bash
  gitmap mail lint --draft-id="draft_9841" --strict-clipping
  ```
- **Data Contract**:
  ```json
  {
    "totalPayloadBytes": 42100,
    "isWithinClippingLimit": true,
    "isSpamScoreOptimal": true,
    "hasSpamKeywords": false,
    "spamScoreValue": 0.8,
    "maxSafeBytes": 102400
  }
  ```
- **Acceptance Criteria**:
  1. If email body exceeds 102,400 bytes, `isWithinClippingLimit` switches to `false` and a persistent red indicator activates.
  2. Clicking the spam score reveals a collapsible popover listing detected words with severity scores.

---

### Item 036: Instant Disposable Test Send & Inbox Trap Simulator Popover
- **ID & Name**: `IMPR-036-COMPOSER-DISPOSABLE-TEST-SENDER`
- **Problem Statement**: Developers and deliverability operators must manually copy-paste preview HTML into external webmail accounts to verify live headers and anti-spam headers.
- **Concrete UI Specification**:
  - "Send Test Sandbox" button positioned beside the main Send dropdown (`padding: 6px 14px; font-weight: 600;`).
  - Modal overlay with auto-generated temporary test address (`trap_abc123@warmup.internal`) or custom input with history chips.
  - Displays instantaneous dispatch log with SMTP handshake stages (EHLO, AUTH, MAIL FROM, RCPT TO, DATA).
  - Shortcut: `Ctrl+Shift+Enter` opens the test sender popover directly.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail test-send --draft-id="draft_9841" --to="sink@testmail.local" --verbose
  ```
- **Data Contract**:
  ```json
  {
    "testDispatchId": "disp_test_901",
    "isSmtpHandshakeSuccessful": true,
    "hasTlsEncrypted": true,
    "hasDkimSigned": true,
    "isSpfAligned": true,
    "latencyMs": 142
  }
  ```
- **Acceptance Criteria**:
  1. Test email dispatches through local mock SMTP daemon within 250ms.
  2. Status checklist reveals green ticks for TLS, DKIM, and SPF validation.

---

### Item 037: Smart Plain-Text Fallback Synchronizer with One-Click Strip & Sync
- **ID & Name**: `IMPR-037-COMPOSER-PLAINTEXT-SYNCHRONIZER`
- **Problem Statement**: Multipart emails often contain broken, out-of-date, or completely missing `text/plain` parts, drastically penalizing spam deliverability scores.
- **Concrete UI Specification**:
  - Sub-tab in editor footer: `HTML Body`, `Plain Text Body`, `Diff View`.
  - Floating button: "Sync from HTML" (`border-radius: 6px; background: var(--bg-button-subtle);`).
  - Intelligently converts markdown/HTML into RFC 3676 compliant plain text (converting links to inline footnotes `[1]`, tables to markdown tables, preserving bullet points).
  - Sync lock toggle: `Keep in Sync` checkbox ensures changes in HTML automatically update plain-text.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail sync-plain --draft-id="draft_9841" --auto-wrap=78
  ```
- **Data Contract**:
  ```json
  {
    "draftId": "draft_9841",
    "isPlainTextInSync": true,
    "hasAutoWrapEnabled": true,
    "isMultipartValid": true,
    "plainTextLineLengthMax": 78
  }
  ```
- **Acceptance Criteria**:
  1. Clicking "Sync from HTML" strips tags and formats links as numbered citations at bottom of plain text.
  2. With sync enabled, typing in the HTML editor propagates text modifications to the plain-text buffer within 50ms.

---

### Item 038: Inline Email Header & Custom `X-Mailer` / `Feedback-ID` Inspector Drawer
- **ID & Name**: `IMPR-038-COMPOSER-HEADER-INSPECTOR-DRAWER`
- **Problem Statement**: Campaign operators cannot review or customize specific outbound RFC 822 email headers (`X-Campaign-ID`, `List-Unsubscribe`, `Feedback-ID`, `Message-ID`) without editing raw backend configs.
- **Concrete UI Specification**:
  - Collapsible slide-down drawer beneath the Subject line (`max-height: 0; overflow: hidden; transition: max-height 0.25s ease-out;`).
  - Key-value table with standard headers auto-populated (e.g. `From`, `Reply-To`, `List-Unsubscribe-Post: List-Unsubscribe=One-Click`).
  - Input fields with auto-validation for RFC standard compliance.
  - Shortcut: `Ctrl+H` toggles the header inspector drawer.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail headers view --draft-id="draft_9841" --format=table
  ```
- **Data Contract**:
  ```json
  {
    "draftId": "draft_9841",
    "hasCustomHeaders": true,
    "isUnsubscribeHeaderPresent": true,
    "isFeedbackIdValid": true,
    "headersCount": 8
  }
  ```
- **Acceptance Criteria**:
  1. Pressing `Ctrl+H` smoothly expands the header grid.
  2. Missing `List-Unsubscribe` displays a recommendation tooltip with auto-insert button.

---

### Item 039: Interactive Draft Version History Slider & Diff Scrubbing
- **ID & Name**: `IMPR-039-COMPOSER-DRAFT-HISTORY-SCRUBBER`
- **Problem Statement**: Accidental edits during template preparation cannot be rolled back easily, and users cannot see visual diffs between automated AI drafts and manual edits.
- **Concrete UI Specification**:
  - Top-right history clock button expanding a horizontal timeline scrubber (`height: 36px; display: flex; align-items: center; gap: 8px;`).
  - Snapshot points represented as circular nodes color-coded by author (`system`, `ai-refactor`, `operator`).
  - Side-by-side or unified green/red diff overlay inside the editor highlighting inserted/deleted lines.
  - Shortcut: `Ctrl+Z` / `Ctrl+Y` and `Ctrl+Alt+H` for timeline overview.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail draft history --draft-id="draft_9841" --limit=5
  ```
- **Data Contract**:
  ```json
  {
    "draftId": "draft_9841",
    "hasVersionHistory": true,
    "isRollbackSafe": true,
    "currentRevisionNumber": 4,
    "totalRevisionsCount": 7
  }
  ```
- **Acceptance Criteria**:
  1. Dragging the scrubber loads past revisions in < 15ms from local SQLite cache.
  2. Restoring a revision creates a new sequential snapshot without destroying forward history.

---

### Item 040: One-Click Template Snippet Injector & AI Subject Line Optimizer
- **ID & Name**: `IMPR-040-COMPOSER-SNIPPET-INJECTOR-AI-SUBJECT`
- **Problem Statement**: Cold outreach and warm-up senders struggle with low subject-line open rates and repetitively copy-paste boilerplate disclaimers or meeting links.
- **Concrete UI Specification**:
  - Quick-insert sidebar drawer with tabbed categories: `Signatures`, `Warm-Up Greetings`, `CTA Blocks`, `Legal Disclaimers`.
  - Sparkle icon inside the Subject input triggering a dropdown with 5 AI-suggested variations scored by predicted open rate (`High`, `Optimal`, `Urgent`).
  - Single click on a snippet or subject instantly writes to the cursor position with clean undo support.
  - Shortcut: `Ctrl+J` triggers snippet launcher modal.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail subject optimize --text="Quick sync regarding gitmap warmup" --style=curiosity
  ```
- **Data Contract**:
  ```json
  {
    "originalSubject": "Quick sync regarding gitmap warmup",
    "isOptimized": true,
    "hasSpamScoreCalculated": true,
    "selectedVariantIndex": 0,
    "predictedOpenRatePercent": 48.2
  }
  ```
- **Acceptance Criteria**:
  1. Clicking an AI suggestion replaces the subject line and updates the live preview immediately.
  2. Injecting a signature block preserves HTML structure and CSS inline styles.

---

## Cluster 5: Multi-Mailbox Identity Grid, Masking & Delegation HUD (Items 41–50)

### Item 041: Multi-Tenant Identity Matrix Card Grid with Real-Time Health Gauges
- **ID & Name**: `IMPR-041-IDENTITY-MATRIX-HEALTH-GAUGES`
- **Problem Statement**: Operators running 50+ warm-up inboxes cannot easily observe which mailboxes are connected, warmed up, rate-limited, or authentication-failed at a glance.
- **Concrete UI Specification**:
  - CSS Grid container (`grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px;`).
  - Cards feature mailbox avatar, primary email, provider badge (Google Workspace, Microsoft 365, Custom IMAP/SMTP), and circular SVG health gauge (0-100%).
  - Color status rings: Green (Warm & Healthy), Amber (Gradual Ramp), Red (Auth Failure or Blacklist).
  - Quick action buttons on card footer: `Pause Warmup`, `Inspect Logs`, `Rotate Credential`.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup fleet list --format=cards --filter-status=active
  ```
- **Data Contract**:
  ```json
  {
    "mailboxId": "mbx_alpha_01",
    "emailAddress": "alex@acmecorp.dev",
    "isWarmupActive": true,
    "isAuthenticationValid": true,
    "hasBlacklistIncident": false,
    "healthScore": 96.5,
    "dailySentCount": 38,
    "dailyTargetQuota": 40
  }
  ```
- **Acceptance Criteria**:
  1. Grid renders 100+ cards smoothly without UI stutter using CSS containment (`contain: content;`).
  2. Health score gauge renders via SVG `stroke-dashoffset` calculation with smooth CSS transition.

---

### Item 042: Fast Account Switching Omnibar with In-Memory Credential Caching
- **ID & Name**: `IMPR-042-IDENTITY-OMNIBAR-ACCOUNT-SWITCHER`
- **Problem Statement**: Navigating between sender accounts requires reloading complex nested menus, breaking workflow when troubleshooting individual mailbox delivery logs.
- **Concrete UI Specification**:
  - Global omnibar activated via shortcut `Ctrl+K` or clicking top-nav identity pill.
  - Instant fuzzy search over email addresses, domain names, and tags (`sales`, `support`, `tier1`).
  - Up/Down arrow selection, `Enter` to switch active tenant context without page reload.
  - Recent accounts list shown by default with pin icon to favorite frequent mailboxes.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail switch --account="alex@acmecorp.dev"
  ```
- **Data Contract**:
  ```json
  {
    "activeMailboxId": "mbx_alpha_01",
    "isContextSwitched": true,
    "hasCachedToken": true,
    "switchDurationMs": 8,
    "cachedAccountsCount": 42
  }
  ```
- **Acceptance Criteria**:
  1. Switching identity takes less than 20ms and updates state reactively across all open UI tabs.
  2. Token validation runs from memory without triggering external HTTP auth requests on switch.

---

### Item 043: Domain Reputation & DNS Alignment Diagnostic HUD (SPF, DKIM, DMARC, MX)
- **ID & Name**: `IMPR-043-IDENTITY-DNS-ALIGNMENT-HUD`
- **Problem Statement**: Deliverability drops silently when SPF records get truncated or DKIM keys expire, with no alert presented before warm-up emails are sent.
- **Concrete UI Specification**:
  - Diagnostic HUD modal or tab with 4 primary status chips: `SPF`, `DKIM`, `DMARC`, `MX`.
  - Clicking any chip unfolds exact DNS record values, TTL, and expected vs resolved raw strings.
  - One-click "Re-verify DNS" button with rotating SVG spinner (`animation: spin 0.8s linear infinite;`).
  - Direct copy-to-clipboard button for required TXT records with toast confirmation.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup verify-domain --domain="acmecorp.dev" --json
  ```
- **Data Contract**:
  ```json
  {
    "domainName": "acmecorp.dev",
    "isSpfValid": true,
    "isDkimValid": true,
    "isDmarcValid": true,
    "isMxValid": true,
    "hasDmarcRejectPolicy": true,
    "dnsLookupDurationMs": 180
  }
  ```
- **Acceptance Criteria**:
  1. If any of SPF, DKIM, or DMARC fails, the parent identity card glows red and flags warmup as paused.
  2. Re-verify button queries local recursive DNS and updates status within 500ms.

---

### Item 044: SMTP/IMAP Connection Pool Latency Benchmark & Jitter Monitor
- **ID & Name**: `IMPR-044-IDENTITY-SMTP-LATENCY-MONITOR`
- **Problem Statement**: Transient connection drops or slow SMTP handshake responses degrade warmup schedules and lead to socket timeouts that go undetected.
- **Concrete UI Specification**:
  - Live sparkline graph showing connection handshake latency over the past 24 hours (sampled every 15 mins).
  - Metrics display: P50 latency, P95 latency, Jitter standard deviation, and TLS handshake duration.
  - Alert banner when connection latency spikes above 2500ms for 3 consecutive probes.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail ping --account="alex@acmecorp.dev" --probes=5
  ```
- **Data Contract**:
  ```json
  {
    "mailboxId": "mbx_alpha_01",
    "isHostReachable": true,
    "hasTlsHandshakeCompleted": true,
    "latencyP50Ms": 142,
    "latencyP95Ms": 310,
    "jitterStdDevMs": 18.4
  }
  ```
- **Acceptance Criteria**:
  1. Handshake latency graph renders via pure SVG path with zero external charting library.
  2. Tooltips reveal exact timestamp and ping breakdown (DNS, TCP Connect, TLS Handshake).

---

### Item 045: Automated OAuth Token Refresh & Headless Session Health Indicator
- **ID & Name**: `IMPR-045-IDENTITY-OAUTH-TOKEN-HEALTH`
- **Problem Statement**: OAuth refresh tokens expire or get revoked, causing warm-up sending to silently halt midway through a schedule.
- **Concrete UI Specification**:
  - Token health badge: `Active (Expires in 52m)`, `Refreshing...`, or `Revoked / Action Needed`.
  - One-click "Re-authenticate" button launching local headless OAuth loopback server (`localhost:18299/callback`).
  - Auto-refresh progress bar showing remaining lifetime before scheduled background token swap.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail auth refresh --account="alex@acmecorp.dev" --force
  ```
- **Data Contract**:
  ```json
  {
    "mailboxId": "mbx_alpha_01",
    "isTokenValid": true,
    "hasRefreshToken": true,
    "isExpiringSoon": false,
    "secondsUntilExpiration": 3120,
    "authProvider": "google_workspace"
  }
  ```
- **Acceptance Criteria**:
  1. Expired tokens trigger a sticky notification banner with direct action to re-auth.
  2. Silent background refresh succeeds without interrupting ongoing warmup worker threads.

---

### Item 046: Proxy & Egress IP Allocation Manager with Geo-Routing Visualizer
- **ID & Name**: `IMPR-046-IDENTITY-PROXY-ROUTING-MANAGER`
- **Problem Statement**: Sending warm-up traffic for multiple accounts from a single IP causes reputation pooling and cross-contamination across client accounts.
- **Concrete UI Specification**:
  - Visual proxy binding table mapping Mailbox -> Dedicated Egress IP -> Geo Location (Flag icon + Country).
  - Proxy health status indicator (`Active`, `Degraded`, `Rotating`, `Dead`).
  - IP leak tester button: verifies egress IP doesn't match primary host server IP.
  - Toggle switch: `Enforce Sticky Proxy Sessions` to prevent sudden IP jumping per mailbox.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup proxy test --mailbox-id="mbx_alpha_01"
  ```
- **Data Contract**:
  ```json
  {
    "mailboxId": "mbx_alpha_01",
    "egressIp": "198.51.100.42",
    "isProxyActive": true,
    "hasIpLeakDetected": false,
    "isStickySessionEnforced": true,
    "countryCode": "US"
  }
  ```
- **Acceptance Criteria**:
  1. Leaked IP alerts halt outbound warmup queues immediately for that mailbox identity.
  2. Proxy table allows rapid reassignment via inline dropdown selector.

---

### Item 047: Warm-Up Peer Pool Matching Visualizer & Exchange Network Map
- **ID & Name**: `IMPR-047-IDENTITY-PEER-POOL-NETWORK-MAP`
- **Problem Statement**: Users cannot visualize who their mailboxes are exchanging warm-up emails with, leading to uncertainty about peer pool diversity.
- **Concrete UI Specification**:
  - Interactive HTML5 Canvas/SVG network graph depicting identity nodes and peer connection edges.
  - Node clusters color-coded by provider: Google (Blue), Outlook (Teal), Yahoo (Purple), Private SMTP (Orange).
  - Hovering a node displays stats: Emails Sent, Received, Marked as Important, Saved from Spam.
  - Filter controls: View by domain, provider distribution percentage, or date range.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup pool graph --format=dot --output="./pool.svg"
  ```
- **Data Contract**:
  ```json
  {
    "totalPoolNodes": 124,
    "isPoolDiverse": true,
    "hasOptimalProviderMix": true,
    "providerDistribution": {
      "google": 0.45,
      "microsoft": 0.40,
      "other": 0.15
    }
  }
  ```
- **Acceptance Criteria**:
  1. Network map renders smoothly at 60 FPS using raw HTML5 2D Canvas without WebGL or third-party engines.
  2. If provider concentration exceeds 65% on one vendor, an alert badge flags unbalanced peer mix.

---

### Item 048: Mailbox Quota Guard & Exponential Ramp-Up Velocity Stepper
- **ID & Name**: `IMPR-048-IDENTITY-VELOCITY-RAMP-STEPPER`
- **Problem Statement**: New mailboxes burned by sending too many emails on Day 1; operators lack visual tools to plan and track automated daily volume increases.
- **Concrete UI Specification**:
  - Step-wise visual trajectory curve showing: Day 1 (5 emails) to Day 28 (50 emails).
  - Current day indicator marker with daily progress bar (`Sent 18 / 20 today`).
  - Adjustable ramp presets: `Gentle` (+2/day), `Balanced` (+4/day), `Aggressive` (+8/day).
  - Hard stop safety switch: `Emergency Brake` locks mailbox sending immediately if bounce rate > 3%.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup ramp set --mailbox-id="mbx_alpha_01" --preset="balanced" --max-daily=50
  ```
- **Data Contract**:
  ```json
  {
    "mailboxId": "mbx_alpha_01",
    "currentWarmupDay": 9,
    "dailyTarget": 25,
    "dailySent": 22,
    "isQuotaGuardActive": true,
    "hasEmergencyBrakeTriggered": false,
    "rampPreset": "balanced"
  }
  ```
- **Acceptance Criteria**:
  1. Exceeding daily limit prevents further automated dispatches and logs event to audit database.
  2. Changing presets recalculates the target curve and updates the projected completion date immediately.

---

### Item 049: Mailbox Credential Vault with GPG/Age Local File Encryption
- **ID & Name**: `IMPR-049-IDENTITY-CREDENTIAL-VAULT`
- **Problem Statement**: Mailbox passwords and OAuth secret keys stored in plaintext JSON or SQLite expose operators to serious credential leaks.
- **Concrete UI Specification**:
  - Vault management drawer showing encrypted credential stores with lock status icon.
  - "Unlock Vault" modal requiring master passphrase, keeping decrypted keys in memory only.
  - Visual security score per mailbox (e.g. `2FA Enabled`, `App Password Used`, `Raw Password Discouraged`).
  - Auto-lock countdown timer (default 15 minutes of inactivity).
- **CLI Counterpart Command**:
  ```bash
  gitmap mail vault status --unlocked
  ```
- **Data Contract**:
  ```json
  {
    "isVaultLocked": false,
    "hasMasterKeyLoaded": true,
    "isEncryptedAtRest": true,
    "cipherAlgorithm": "age-x25519",
    "unlockedCredentialsCount": 12
  }
  ```
- **Acceptance Criteria**:
  1. Plaintext secrets are never written to disk or exposed in client DOM inspect tools.
  2. Auto-lock timeout purges master keys from memory and flips `isVaultLocked` to `true`.

---

### Item 050: Multi-Mailbox Bulk Tagging, Filtering & Batch Action Toolbar
- **ID & Name**: `IMPR-050-IDENTITY-BULK-ACTIONS-TOOLBAR`
- **Problem Statement**: Managing 100+ accounts one-by-one is tedious; users need fast bulk selection to apply tags, pause warmup, or update ramp schedules simultaneously.
- **Concrete UI Specification**:
  - Floating bottom action toolbar that slides into view when one or more cards/rows are selected (`transform: translateY(0); transition: 0.2s cubic-bezier(0.16, 1, 0.3, 1);`).
  - Selection counter (`"18 mailboxes selected"`), Select All / Invert Selection buttons.
  - Action pills: `Pause Warmup`, `Resume Warmup`, `Assign Tag...`, `Set Ramp Curve...`, `Export Diagnostics`.
  - Shortcut: `Ctrl+A` (when grid focused) selects all filtered accounts.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup fleet batch --filter-tag="tier2" --action="pause"
  ```
- **Data Contract**:
  ```json
  {
    "selectedMailboxIds": ["mbx_alpha_01", "mbx_alpha_02"],
    "hasBulkActionRunning": false,
    "isAllSelected": false,
    "selectedCount": 2,
    "appliedTag": "staging"
  }
  ```
- **Acceptance Criteria**:
  1. Selecting items updates the floating bar count instantaneously without re-rendering unaffected cards.
  2. Batch action execution provides a live progress bar and per-account success/failure summary.

---

## Cluster 6: High-Velocity CLI Terminal UI, Lipgloss Tables & Interactive ANSI Prompts (Items 51–60)

### Item 051: Bubbletea/Lipgloss Interactive Full-Terminal Warmup Monitor Dashboard
- **ID & Name**: `IMPR-051-CLI-LIPGLOSS-DASHBOARD`
- **Problem Statement**: CLI operators lack a rich, live terminal dashboard to monitor multi-mailbox warm-up progress without cluttering terminal scrollback buffers.
- **Concrete UI Specification**:
  - Full-screen alternate screen buffer application (`tea.EnterAltScreen`) using Charmbracelet Bubbletea & Lipgloss.
  - Three-column layout: Left (Mailbox list with mini health gauges), Center (Real-time live send/receive event log), Right (Key metrics, delivery rates, queue backlog).
  - High-contrast ANSI colors adhering to user terminal theme with fallback 16-color mode.
  - Keybindings: `q` or `Ctrl+C` to quit, `tab` to switch panes, `j`/`k` or arrow keys for navigation.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup tui --refresh-rate=1s
  ```
- **Data Contract**:
  ```json
  {
    "isAltScreenActive": true,
    "hasTermSizeDetected": true,
    "termWidth": 140,
    "termHeight": 42,
    "activePaneIndex": 0,
    "isUtf8Supported": true
  }
  ```
- **Acceptance Criteria**:
  1. Launching `gitmap warmup tui` cleanly switches to alt screen without flashing or residual artifacts.
  2. Exiting completely restores user's prior terminal screen buffer and cursor state.

---

### Item 052: Lipgloss Colored Unicode Table for Mailbox Diagnostics & Deliverability
- **ID & Name**: `IMPR-052-CLI-LIPGLOSS-UNICODE-TABLE`
- **Problem Statement**: Standard CLI tabular output is often misaligned, truncates email addresses without warning, and lacks color-coded severity indicators.
- **Concrete UI Specification**:
  - Border styles using Lipgloss rounded unicode borders (`lipgloss.RoundedBorder()`).
  - Columns: `STATUS`, `MAILBOX`, `PROVIDER`, `SENT / QUOTA`, `HEALTH`, `DNS (S/D/M)`, `LAST EVENT`.
  - Color styling: Green bullet (`●`) for active, yellow triangle (`▲`) for ramp, red cross (`✖`) for error.
  - Responsive column collapsing when terminal width is < 100 columns.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup status --format=table --color=always
  ```
- **Data Contract**:
  ```json
  {
    "hasUnicodeBorders": true,
    "isColorEnabled": true,
    "rowCount": 15,
    "isTruncationHandled": true,
    "tableRenderTimeMs": 2.4
  }
  ```
- **Acceptance Criteria**:
  1. Terminal columns adjust automatically to terminal dimensions without line wrapping errors.
  2. Piping output (`gitmap warmup status | cat`) automatically disables ANSI escape codes and borders.

---

### Item 053: Interactive Fuzzy Mailbox Selector & Action Launcher
- **ID & Name**: `IMPR-053-CLI-INTERACTIVE-FUZZY-SELECTOR`
- **Problem Statement**: Typing long email addresses into CLI command flags (`--mailbox="john.doe.engineering@company.com"`) is error-prone and slow.
- **Concrete UI Specification**:
  - Interactive prompt using Bubbletea textinput and fuzzy matching engine (matching on email, domain, or alias).
  - Instant highlighted matching characters in bold cyan (`lipgloss.Color("#00D7FF")`).
  - Preview panel on the right displaying mailbox configuration and current token health.
  - Keys: `Enter` selects, `Esc` cancels, `Ctrl+N`/`Ctrl+P` moves selection.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail pick
  ```
- **Data Contract**:
  ```json
  {
    "searchQuery": "eng",
    "hasMatchedResults": true,
    "matchedCount": 4,
    "selectedMailboxId": "mbx_eng_04",
    "isInteractivePrompt": true
  }
  ```
- **Acceptance Criteria**:
  1. Typing filters through 500 mailboxes with sub-millisecond response latency (< 1ms).
  2. Selected email address is passed cleanly to the subsequent CLI sub-command.

---

### Item 054: Lipgloss Dual-Stream Progress Bar for Concurrent Email Dispatch & Read
- **ID & Name**: `IMPR-054-CLI-DUAL-STREAM-PROGRESS-BAR`
- **Problem Statement**: Warm-up cycles run both outbound SMTP dispatches and inbound IMAP retrieval checks concurrently; standard single progress bars fail to show both operations.
- **Concrete UI Specification**:
  - Stacked Lipgloss progress bars: Top bar (Outbound Send Progress: Cyan ramp), Bottom bar (Inbound Read & Interaction Progress: Magenta ramp).
  - Numeric stats appended on the right: `[42/50 Sent | 84%]` and `[38/40 Read | 95%]`.
  - Smooth animation using tick events (`tea.Tick`) with percentage decimal precision.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup run-cycle --parallel=4
  ```
- **Data Contract**:
  ```json
  {
    "isOutboundActive": true,
    "isInboundActive": true,
    "hasDualProgressRendered": true,
    "outboundPercent": 84.0,
    "inboundPercent": 95.0,
    "activeThreads": 4
  }
  ```
- **Acceptance Criteria**:
  1. Bars advance independently without visual tearing or terminal line skipping.
  2. Upon 100% completion, bars freeze into a solid green summary line.

---

### Item 055: ANSI Live Streaming Log Viewer with Granular Level Filtering & Pause
- **ID & Name**: `IMPR-055-CLI-STREAMING-LOG-VIEWER`
- **Problem Statement**: Verbose logs during warm-up execution stream too rapidly to read, while normal logs hide important SMTP exchange details.
- **Concrete UI Specification**:
  - Live log streaming viewport widget with auto-scroll lock.
  - Interactive hotkeys: `space` toggles pause/resume stream, `f` opens live regex filter prompt, `1-4` toggles log levels (`DEBUG`, `INFO`, `WARN`, `ERROR`).
  - Syntax highlighted RFC protocol exchanges: `>>> MAIL FROM:...` in green, `<<< 250 OK` in cyan, `<<< 550 Blocked` in bold red.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup logs --tail --level=debug
  ```
- **Data Contract**:
  ```json
  {
    "isStreamPaused": false,
    "hasActiveFilter": true,
    "activeLogLevel": "DEBUG",
    "bufferedLinesCount": 1500,
    "isAutoScrollLocked": true
  }
  ```
- **Acceptance Criteria**:
  1. Hitting `space` freezes the buffer instantly while retaining incoming entries in a ring buffer.
  2. Filtering by regex highlights matches in yellow and filters non-matching lines within 5ms.

---

### Item 056: Interactive Email Template Preview in Terminal with Inline ANSI Renderer
- **ID & Name**: `IMPR-056-CLI-INLINE-TERMINAL-TEMPLATE-PREVIEW`
- **Problem Statement**: Developers tweaking email templates in headless SSH sessions cannot verify visual layout or link placements without opening a browser.
- **Concrete UI Specification**:
  - Terminal-based HTML/Markdown renderer converting styles to ANSI 256-color / true-color equivalents.
  - Headings displayed in bold underline, links formatted with OSC 8 clickable terminal hyperlink escapes (`\x1b]8;;url\x1b\text\x1b]8;;\x1b\`).
  - Tables converted to ASCII/Unicode grid lines.
  - Fallback to colored plain-text if terminal does not support truecolor.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail template preview --template="welcome-warmup.md" --terminal
  ```
- **Data Contract**:
  ```json
  {
    "templatePath": "./templates/welcome-warmup.md",
    "isTrueColorSupported": true,
    "hasClickableOsc8Links": true,
    "isRenderSuccessful": true,
    "renderDurationMs": 4.1
  }
  ```
- **Acceptance Criteria**:
  1. Links in supported terminals (iTerm2, Alacritty, WezTerm, Windows Terminal) are directly clickable.
  2. Terminal width constraints are respected with clean word-wrapping without splitting words.

---

### Item 057: Quick DNS Diagnosis CLI Wizard with Interactive Fix Generator
- **ID & Name**: `IMPR-057-CLI-DNS-WIZARD-FIX-GENERATOR`
- **Problem Statement**: Diagnosing DNS misconfigurations in the terminal usually requires chaining `dig`, `nslookup`, and manual string comparisons.
- **Concrete UI Specification**:
  - Guided step-by-step wizard prompting for domain name and target providers.
  - Executes live DNS lookups for MX, TXT, CNAME records and highlights mismatches side-by-side.
  - Outputs copy-pasteable BIND / Cloudflare / Route53 record formats with exact quotes and escapes.
  - Confirmation prompt: `[y/N]` to export generated zone file snippet to disk.
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup doctor --domain="acmecorp.dev" --wizard
  ```
- **Data Contract**:
  ```json
  {
    "domain": "acmecorp.dev",
    "isDoctorCompleted": true,
    "hasDnsDiscrepancy": true,
    "suggestedFixCount": 2,
    "isExportSnippetAvailable": true
  }
  ```
- **Acceptance Criteria**:
  1. Wizard pinpoints exact missing or conflicting SPF tokens (e.g. `~all` vs `-all`).
  2. Cloudflare formatted zone lines can be exported directly with a single confirmation.

---

### Item 058: Interactive Outbound Queue Purge & Reschedule TUI Modal
- **ID & Name**: `IMPR-058-CLI-QUEUE-PURGE-RESCHEDULE-MODAL`
- **Problem Statement**: When an error occurs (such as an incorrect subject line), stopping and purging hundreds of queued warm-up messages via standard kill commands can orphan background threads.
- **Concrete UI Specification**:
  - Modal prompt overlay in terminal (`border: double red; padding: 1 2;`).
  - Displays queue summary: Total Pending, In-Flight, Scheduled, Failed.
  - Action choices via arrow keys: `Cancel Pending Only`, `Purge Failed`, `Pause Entire Fleet`, `Abort & Reschedule`.
  - Safety confirmation: Requires typing uppercase `PURGE` to prevent accidental mass deletion.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail queue purge --interactive
  ```
- **Data Contract**:
  ```json
  {
    "isPurgeConfirmed": true,
    "hasActiveJobsCanceled": true,
    "purgedJobsCount": 142,
    "isFleetPaused": true
  }
  ```
- **Acceptance Criteria**:
  1. Typing anything other than `PURGE` cancels the action cleanly without touching the database.
  2. Active in-flight SMTP connections are gracefully terminated with proper TCP FIN handshakes.

---

### Item 059: Lipgloss ASCII Health Sparklines & Fleet Deliverability Summary
- **ID & Name**: `IMPR-059-CLI-LIPGLOSS-ASCII-SPARKLINES`
- **Problem Statement**: Viewing historical warmup trends across the last 7 to 30 days traditionally requires launching a web browser or generating heavy image files.
- **Concrete UI Specification**:
  - Compact terminal summary card utilizing Unicode block characters (` ▂▃▄▅▆▇█`) to draw 14-day trend sparklines directly in stdout.
  - Side-by-side fleet metrics: Open Rate Sparkline, Spam-Folder Recovery Sparkline, Inbox Placement Rate.
  - Color gradient applied across bars: Low (Red/Yellow) to High (Emerald Green).
- **CLI Counterpart Command**:
  ```bash
  gitmap warmup summary --days=14 --sparkline
  ```
- **Data Contract**:
  ```json
  {
    "timeframeDays": 14,
    "hasSparklineGenerated": true,
    "isPlacementImproving": true,
    "currentInboxPlacementPercent": 94.2,
    "sparklineGlyphs": " ▂▄▅▆▇██▇█████"
  }
  ```
- **Acceptance Criteria**:
  1. Sparklines render accurately across all UTF-8 compliant terminal emulators.
  2. Non-interactive or piped output outputs clean numeric CSV/JSON if `--no-ansi` is detected.

---

### Item 060: Zero-Config CLI Configuration Wizard with Encrypted Profile Creation
- **ID & Name**: `IMPR-060-CLI-ZERO-CONFIG-SETUP-WIZARD`
- **Problem Statement**: Onboarding new SMTP/IMAP accounts into GitMap requires manual editing of YAML configuration files, resulting in syntax errors and unvalidated credentials.
- **Concrete UI Specification**:
  - Interactive multi-step form wizard powered by Charmbracelet Huh? / Bubbletea.
  - Step 1: Mailbox Email & Display Name.
  - Step 2: Auto-detection of SMTP/IMAP servers based on MX lookup (Google, Microsoft, Fastmail, Postmark, Custom).
  - Step 3: Password / App Password input with masked characters (`****`).
  - Step 4: Instant connection verification test before saving profile.
  - Step 5: Master passphrase prompt to encrypt configuration into `~/.gitmap/credentials.age`.
- **CLI Counterpart Command**:
  ```bash
  gitmap mail setup --wizard
  ```
- **Data Contract**:
  ```json
  {
    "isWizardCompleted": true,
    "hasValidatedConnection": true,
    "isProfileEncrypted": true,
    "configuredProvider": "google_workspace",
    "savedProfileName": "primary-outreach"
  }
  ```
- **Acceptance Criteria**:
  1. Auto-detection resolves standard IMAP/SMTP endpoints correctly for 95% of top email providers.
  2. Invalid credentials fail fast at Step 4, highlighting connection error before writing anything to disk.

---

## Traceability Matrix & File Reference
- Document Location: `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/03-improvements-31-to-60-composer-studio-and-cli-ux.md`
- Relative Paths Verified: All references to files, templates, and specs use relative paths.
- Coding Guidelines: Booleans strictly adhere to positive `is` / `has` naming conventions with zero negative prefixes.
