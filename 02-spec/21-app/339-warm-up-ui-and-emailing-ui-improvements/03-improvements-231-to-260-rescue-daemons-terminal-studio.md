# Specification 339: Improvements 231 to 260 (Rescue Daemons, Terminal Studio & P2P Mesh)

- **Sequence:** 339
- **Slug:** `339-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/03-improvements-231-to-260-rescue-daemons-terminal-studio.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (231 - 260)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 24** | Items 231–240 | Autonomous Lemon Service Deliverability Rescue Daemons |
| **Cluster 25** | Items 241–250 | Interactive Terminal Studio (`gitmap mail studio`, Lipgloss/Bubble Tea) |
| **Cluster 26** | Items 251–260 | P2P Multi-Tenant Reputation Mesh & Anomaly Quarantining |

---

# Cluster 24: Autonomous Lemon Service Deliverability Rescue Daemons (Items 231–240)

---

### Improvement 231: Autonomous SMTP Greylist Backoff & Rate De-escalation Daemon
- **ID:** `WARMUP-RESCUE-231`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** When major mailbox providers (e.g., Microsoft 365, Yahoo) trigger 4xx greylisting response codes (`451 4.7.500 Server busy`), sending jobs continue hammering the server, resulting in temporary IP pool blocks and reputation degradation.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Live throttling gauge widget (`.rescue-gauge-greylist`) showing current backoff state: Normal (Green `#10b981`), Cooldown (Amber `#f59e0b`), and Deep Backoff (Red `#ef4444`).
    - Real-time event feed streaming SMTP `451`/`421` responses with millisecond timestamps and destination MX hostname.
    - Manual override toggle allowing operators to force immediate rate recovery or lock safety throttle.
  - **Lipgloss Terminal UX:**
    - Styled status card with `lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#f59e0b"))`.
    - Live progress bar showing active rate de-escalation percentage and countdown timer to next tier restoration.
- **CLI Counterpart Command:**
  - `gitmap warmup rescue daemon status --service "lemon-greylist"`
  - `gitmap warmup rescue throttle apply --pool "us-east-pool-1" --mx "protection.outlook.com" --rate-limit 15`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["daemonId", "poolId", "mxHost", "backoffLevel", "isActive", "hasTriggeredBackoff", "messagesPerMinute", "cooldownRemainingSeconds"],
    "properties": {
      "daemonId": { "type": "string" },
      "poolId": { "type": "string" },
      "mxHost": { "type": "string" },
      "backoffLevel": { "type": "string", "enum": ["nominal", "cooldown", "deep_backoff", "emergency_stop"] },
      "isActive": { "type": "boolean" },
      "hasTriggeredBackoff": { "type": "boolean" },
      "messagesPerMinute": { "type": "integer", "minimum": 0 },
      "cooldownRemainingSeconds": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - When 3 consecutive 4xx greylist responses occur within 10 seconds, the daemon automatically drops pool outbound rate by 50%.
  - Zero-CDN Web UI updates without reloading using native Server-Sent Events (SSE).
  - CLI command exits with code 0 and outputs structured status table.

---

### Improvement 232: DKIM Signature Drift & Key Rotation Rescuer
- **ID:** `WARMUP-RESCUE-232`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** DNS propagation delays or corrupted selector records during automated 90-day DKIM key rotations lead to widespread DKIM authentication failure (`dkim=fail (body hash did not verify)`), silently dropping inbox placement.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Split-screen verification viewer showing local active private key selector vs public DNS TXT record query result.
    - Visual status badges for selector parity: `DKIM_ALIGNED` (Green), `DKIM_DRIFT_DETECTED` (Red), `DNS_QUERY_PENDING` (Gray).
    - One-click rollback button to restore previous verified selector from local Lemon SQLite vault.
  - **Lipgloss Terminal UX:**
    - High-contrast diff view rendered with Lipgloss: Left column shows expected canonical key record, right column shows live DNS resolution with highlighting on mismatching characters.
- **CLI Counterpart Command:**
  - `gitmap mail dkim rescue verify --domain "send.example.com" --selector "k1"`
  - `gitmap mail dkim rescue rollback --domain "send.example.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "selector", "isAligned", "hasDriftDetected", "publicKeyBits", "lastCheckedTimestamp"],
    "properties": {
      "domain": { "type": "string" },
      "selector": { "type": "string" },
      "isAligned": { "type": "boolean" },
      "hasDriftDetected": { "type": "boolean" },
      "publicKeyBits": { "type": "integer", "enum": [1024, 2048, 4096] },
      "lastCheckedTimestamp": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Rescuer daemon polls DNS every 120 seconds during key transition periods.
  - In case of signature drift, sending is paused instantly on affected selector and diverted to healthy backup selector.
  - CLI `verify` displays exact DNS discrepancy without throwing unhandled exceptions.

---

### Improvement 233: Spam Trap Proximity & Seed List Quarantining Daemon
- **ID:** `WARMUP-RESCUE-233`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** Pristine spam traps or recycled honeypot addresses accidentally ingested into warmup seed pools cause immediate IP blacklisting on Spamhaus or Barracuda.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Trap probability risk matrix displaying recipient seed addresses scored from 0.0 to 1.0 based on engagement patterns and bounce telemetry.
    - Red alert banner displaying flagged suspicious seed addresses quarantined by the daemon.
    - Visual timeline of seed interaction history with granular opt-out toggle.
  - **Lipgloss Terminal UX:**
    - Table rendered with `lipgloss.NewStyle()` showing quarantined seeds, risk score colored with ANSI color ramp, and detection reason.
- **CLI Counterpart Command:**
  - `gitmap warmup trap-rescue scan --pool-file "./data/seeds/pool-alpha.json"`
  - `gitmap warmup trap-rescue quarantine list`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["seedEmail", "riskScore", "isQuarantined", "hasTrapCharacteristics", "quarantinedAt"],
    "properties": {
      "seedEmail": { "type": "string", "format": "email" },
      "riskScore": { "type": "number", "minimum": 0, "maximum": 1 },
      "isQuarantined": { "type": "boolean" },
      "hasTrapCharacteristics": { "type": "boolean" },
      "quarantinedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Seeds with risk score > 0.85 are automatically moved to quarantine state and excluded from warmup dispatch within 50ms.
  - Quarantined records write to Lemon SQLite database with complete audit trace.
  - CLI returns non-zero exit code only if an unexpected operational error occurs, adhering to AppError guidelines.

---

### Improvement 234: Autonomous TLS Cipher & Cert Expiry Health Daemon
- **ID:** `WARMUP-RESCUE-234`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** Expired Let's Encrypt TLS certificates or deprecated cipher suites (TLS 1.0/1.1) on custom SMTP relay daemons cause target MXs (like Gmail) to refuse STARTTLS and silently classify traffic as unencrypted or suspicious.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Certificate health card showing certificate issuer, expiration countdown in days/hours, and supported cipher suite badges.
    - Yellow warning pill when cert expiry is under 14 days; red flashing pill when under 48 hours.
    - Autonomous renewal trigger button with real-time ACME challenge log stream.
  - **Lipgloss Terminal UX:**
    - Sleek TUI card displaying TLS handshake negotiation test result (`TLS_AES_256_GCM_SHA384` vs legacy ciphers) with green checkmarks.
- **CLI Counterpart Command:**
  - `gitmap mail tls inspect --host "smtp.outbound-node.internal" --port 587`
  - `gitmap mail tls renew-cert --host "smtp.outbound-node.internal" --force`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["host", "port", "certExpiryDays", "isValidCert", "hasTls13Support", "isRenewalRequired"],
    "properties": {
      "host": { "type": "string" },
      "port": { "type": "integer" },
      "certExpiryDays": { "type": "integer" },
      "isValidCert": { "type": "boolean" },
      "hasTls13Support": { "type": "boolean" },
      "isRenewalRequired": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Daemon performs synthetic STARTTLS probe every 360 minutes to verify certificate validity.
  - Certs with <= 7 days remaining automatically trigger ACME DNS-01 or HTTP-01 renewal without downtime.
  - UI renders cipher suites without relying on external CDNs or external style sheets.

---

### Improvement 235: MX Fallback & Emergency Smart-Host Failover Daemon
- **ID:** `WARMUP-RESCUE-235`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** Upstream IP blockages or route flapping at a hosting provider completely stall email warmup queues, dropping overall sender reputation due to burst re-attempts upon recovery.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Dual-circuit routing topology diagram displaying primary smart-host connection status and live backup pool readiness.
    - Automatic failover trigger rules configuration: failure threshold (e.g., >10% bounces in 5 min) and recovery stabilization window.
    - Simulated drill button to safely test failover routing in dry-run mode.
  - **Lipgloss Terminal UX:**
    - Dual column table: Primary Route vs Standby Route with active throughput and ping latency metrics.
- **CLI Counterpart Command:**
  - `gitmap mail route failover status`
  - `gitmap mail route failover trigger --target-pool "standby-aws-ses" --reason "ip-reputation-drop"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["activeRouteId", "standbyRouteId", "isFailoverActive", "hasCircuitTripped", "tripReason", "switchedAt"],
    "properties": {
      "activeRouteId": { "type": "string" },
      "standbyRouteId": { "type": "string" },
      "isFailoverActive": { "type": "boolean" },
      "hasCircuitTripped": { "type": "boolean" },
      "tripReason": { "type": "string" },
      "switchedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Circuit breaker trips automatically upon 5 consecutive connection timeouts or 10% hard bounce rate.
  - Queue reroutes to standby smart-host with zero dropped messages in the SQLite spool.
  - CLI provides immediate feedback and structured status report.

---

### Improvement 236: Autonomous Bounce Log Telemetry Parser & Classification Daemon
- **ID:** `WARMUP-RESCUE-236`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** Standard SMTP bounce logs combine hard bounces (unknown mailbox `550 5.1.1`), policy blocks (`554 5.7.1`), and temporary rate limits into generic error strings, preventing targeted remediation.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Real-time classification donut chart rendered via SVG showing Hard Bounce vs Policy Block vs Reputation Reject vs Greylist.
    - Regex pattern discovery table identifying new or unclassified ISP error codes.
    - Drill-down log viewer with syntax highlighted SMTP response chunks and suggested corrective action.
  - **Lipgloss Terminal UX:**
    - Colorized log stream using Lipgloss styles with color-coded tags: `[HARD-BOUNCE]` (Red), `[POLICY]` (Purple), `[TEMP-LIMIT]` (Yellow).
- **CLI Counterpart Command:**
  - `gitmap mail bounce parse --log-file "./logs/lemon-smtp.log" --classify`
  - `gitmap mail bounce summary --since "24h"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["logId", "smtpCode", "enhancedCode", "category", "isHardBounce", "hasPermanentFailure", "recommendation"],
    "properties": {
      "logId": { "type": "string" },
      "smtpCode": { "type": "integer" },
      "enhancedCode": { "type": "string" },
      "category": { "type": "string", "enum": ["hard_bounce", "soft_bounce", "policy_block", "rate_limit", "greylist"] },
      "isHardBounce": { "type": "boolean" },
      "hasPermanentFailure": { "type": "boolean" },
      "recommendation": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Daemon parses 5,000 bounce events per second with pure Go regular expression parser.
  - RFC 3463 enhanced status codes correctly map to internal classification taxonomy.
  - CLI bounce summary computes accurate percentages for all classified categories.

---

### Improvement 237: DMARC Forensic RUF Report Parsing & Anomaly Daemon
- **ID:** `WARMUP-RESCUE-237`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** DMARC forensic failure reports (RUF / AFRF format) generated when emails fail DKIM or SPF checks contain crucial evidence of spoofing or configuration bugs, but are typically dropped due to raw XML/MIME parsing complexity.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Interactive forensic report table displaying failing sender IP, authenticating domain, and auth-result breakdown.
    - Inspect modal presenting anonymized raw AFRF message body and parsed MIME headers.
    - IP geolocation map showing geographical origin of unauthorized senders.
  - **Lipgloss Terminal UX:**
    - Formatted terminal report showing failing auth mechanisms (`SPF: FAIL`, `DKIM: NONE`) and recipient ISP.
- **CLI Counterpart Command:**
  - `gitmap mail dmarc parse-ruf --input-dir "./data/dmarc/incoming/"`
  - `gitmap mail dmarc anomalies list --severity "high"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["reportId", "sourceIp", "spfResult", "dkimResult", "isAligned", "hasSpoofingSignature", "receivedAt"],
    "properties": {
      "reportId": { "type": "string" },
      "sourceIp": { "type": "string" },
      "spfResult": { "type": "string", "enum": ["pass", "fail", "softfail", "neutral", "none"] },
      "dkimResult": { "type": "string", "enum": ["pass", "fail", "none"] },
      "isAligned": { "type": "boolean" },
      "hasSpoofingSignature": { "type": "boolean" },
      "receivedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Parser ingests RFC 6591 format forensic emails and stores parsed records in SQLite.
  - Anomaly detection flags source IPs sending unaligned emails exceeding 5 events per hour.
  - No external JavaScript or CSS CDNs required for table rendering.

---

### Improvement 238: Reverse DNS (PTR) Drift & Forward Alignment Daemon
- **ID:** `WARMUP-RESCUE-238`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** Cloud providers (e.g., Hetzner, DigitalOcean) occasionally reset or wipe custom Reverse DNS (rDNS/PTR) records during VPS reboots or network migrations, causing immediate delivery failures on major consumer ISPs (Gmail, Outlook).
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - FCrDNS (Forward Confirmed Reverse DNS) status table for each active IP in pool.
    - Three-way verification indicators: IP -> PTR Hostname, Hostname -> A Record, Parity Match.
    - One-click resync button that interacts with cloud provider API to re-assert PTR record.
  - **Lipgloss Terminal UX:**
    - Terminal check list showing IP, resolved PTR, forward A record, and alignment status with green/red checkmarks.
- **CLI Counterpart Command:**
  - `gitmap mail ptr verify --ip "198.51.100.25"`
  - `gitmap mail ptr audit-fleet --pool "production-pool"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["ipAddress", "ptrHostname", "forwardIp", "isAligned", "hasPtrRecord", "checkedAt"],
    "properties": {
      "ipAddress": { "type": "string" },
      "ptrHostname": { "type": "string" },
      "forwardIp": { "type": "string" },
      "isAligned": { "type": "boolean" },
      "hasPtrRecord": { "type": "boolean" },
      "checkedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Daemon verifies FCrDNS every 60 minutes across all registered sending pool IPs.
  - If a drift is detected, sending from the affected IP is halted immediately with warning emitted.
  - CLI `audit-fleet` returns tabular verification summary with positive booleans.

---

### Improvement 239: ISP Feedback Loop (FBL) Abuse Rate Spike Throttler
- **ID:** `WARMUP-RESCUE-239`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** If spam complaints on Yahoo/AOL or Outlook Smart Network Data Services (SNDS) exceed 0.1% during warmup, domain reputation enters a spiral where recovery takes weeks without instant emergency rate cuts.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Real-time abuse rate gauge with visual safety threshold markers at 0.05% (caution) and 0.1% (critical).
    - Auto-pause countdown banner that triggers when complaint rate breaches 0.08%.
    - Recipient unsubscriber list populated instantaneously upon FBL ARF (Abuse Reporting Format) receipt.
  - **Lipgloss Terminal UX:**
    - Lipgloss bar chart displaying complaint percentage over the last 7 days compared to sending volume.
- **CLI Counterpart Command:**
  - `gitmap warmup fbl status --domain "outbound.company.com"`
  - `gitmap warmup fbl configure --threshold 0.0008 --action "pause-and-alert"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "complaintRatio", "thresholdRatio", "isThrottled", "hasBreachedThreshold", "suppressedCount"],
    "properties": {
      "domain": { "type": "string" },
      "complaintRatio": { "type": "number", "minimum": 0 },
      "thresholdRatio": { "type": "number", "minimum": 0 },
      "isThrottled": { "type": "boolean" },
      "hasBreachedThreshold": { "type": "boolean" },
      "suppressedCount": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Ingestion of an ARF complaint instantly adds the complainer to the suppression table in Lemon DB.
  - If complaints exceed configured threshold, the warmup daemon automatically throttles campaign velocity by 75%.
  - CLI command displays exact ratio and current suppression count.

---

### Improvement 240: IP Reputation Blacklist Continuous Monitor & Delisting Automator
- **ID:** `WARMUP-RESCUE-240`
- **Cluster:** 24 (Autonomous Lemon Service Deliverability Rescue Daemons)
- **Problem Statement:** Delayed awareness of DNSBL listings (Spamhaus ZEN, SORBS, Barracuda, Invaluement) means warm-up traffic continues broadcasting from listed IPs, multiplying delivery penalties across other non-listed networks.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Live DNSBL grid displaying 40+ DNSBL providers with real-time status pills: `CLEAN` (Green), `LISTED` (Red), `QUERY_TIMEOUT` (Orange).
    - Delisting action drawer providing direct automated delisting links and generated pre-filled delisting request templates.
    - Historical listing timeline tracking listing duration and delisting confirmation.
  - **Lipgloss Terminal UX:**
    - Multi-column grid output rendered with Lipgloss showing DNSBL provider names and colored status indicators.
- **CLI Counterpart Command:**
  - `gitmap mail dnsbl check --ip "198.51.100.25"`
  - `gitmap mail dnsbl generate-delist-request --ip "198.51.100.25" --dnsbl "zen.spamhaus.org"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["ipAddress", "dnsblHost", "isClean", "hasActiveListing", "responseCode", "checkedAt"],
    "properties": {
      "ipAddress": { "type": "string" },
      "dnsblHost": { "type": "string" },
      "isClean": { "type": "boolean" },
      "hasActiveListing": { "type": "boolean" },
      "responseCode": { "type": "string" },
      "checkedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Daemon queries top 25 DNSBLs concurrently via Go goroutines with a 3-second timeout per provider.
  - Any positive listing triggers an immediate status badge update and optional auto-quarantine of the IP.
  - Generated delisting request includes valid reverse DNS, contact address, and remediation explanation.

---

# Cluster 25: Interactive Terminal Studio (`gitmap mail studio`, Lipgloss/Bubble Tea) (Items 241–250)

---

### Improvement 241: Interactive Bubble Tea Mail Studio Dashboard
- **ID:** `WARMUP-STUDIO-241`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Operators managing email infrastructure from SSH terminals must juggle disparate shell commands (`mailq`, `postqueue`, custom scripts) without a unified, responsive curses-style control panel.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Full-terminal dashboard split into 3 responsive panes: Left Sidebar (Navigation & Accounts), Main Center (Live Warmup Stage, In-Flight Queue, Deliverability Score), Bottom Footer (System Telemetry, Log Stream, Hotkey Shortcuts: `[q] Quit`, `[p] Pause`, `[r] Refresh`, `[tab] Switch Pane`).
    - Dynamic resizing support handling terminal sizes from 80x24 up to full 4K terminal emulators.
    - Color tokens matching GitMap design: Lavender `#a78bfa`, Forest `#10b981`, Amber `#f59e0b`, Coral `#ef4444`.
  - **Zero-CDN Web UI Mirror:**
    - Web terminal emulator view rendering the exact same ASCII/ANSI stream via HTML `<pre>` with live WebSocket mirroring.
- **CLI Counterpart Command:**
  - `gitmap mail studio`
  - `gitmap mail studio --account "corp-primary" --view "overview"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["activeAccountId", "paneFocus", "terminalWidth", "terminalHeight", "isActive", "hasUnreadAlerts"],
    "properties": {
      "activeAccountId": { "type": "string" },
      "paneFocus": { "type": "string", "enum": ["sidebar", "viewport", "footer", "dialog"] },
      "terminalWidth": { "type": "integer", "minimum": 80 },
      "terminalHeight": { "type": "integer", "minimum": 24 },
      "isActive": { "type": "boolean" },
      "hasUnreadAlerts": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Launches in under 100ms with smooth 60fps keyboard-driven cursor navigation.
  - Handles `tea.WindowSizeMsg` dynamically without UI shearing or broken border lines.
  - Clean exit restoring terminal cursor and alternate screen buffer on `q` or `Ctrl+C`.

---

### Improvement 242: Terminal Interactive MIME Email Composer & Template Previewer
- **ID:** `WARMUP-STUDIO-242`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Drafting and testing rich HTML/multipart warmup templates requires external tools, making quick terminal-based testing of variable substitution and headers tedious.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Dual-pane editor: Left pane contains raw Markdown/HTML template source with syntax highlighting; Right pane renders live formatted terminal preview using ANSI styles and terminal image protocol (Sixel/Kitty) where supported.
    - Header config bar at top: `From:`, `To:`, `Subject:`, `DKIM-Selector:`, `Precedence: bulk`.
    - Variable substitution evaluator showing real-time resolution of `{{.SeedName}}`, `{{.Company}}`, `{{.RandomEntropy}}`.
  - **Zero-CDN Web UI:**
    - Pure HTML/CSS preview card with toggle between Raw MIME, HTML rendered, and Plaintext fallback views.
- **CLI Counterpart Command:**
  - `gitmap mail studio compose --template "./templates/warmup-intro.html"`
  - `gitmap mail template preview --file "./templates/warmup-intro.html" --data '{"SeedName":"Alice"}'`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["templateId", "subject", "htmlByteLength", "textByteLength", "isValidMime", "hasVariablesResolved"],
    "properties": {
      "templateId": { "type": "string" },
      "subject": { "type": "string" },
      "htmlByteLength": { "type": "integer", "minimum": 0 },
      "textByteLength": { "type": "integer", "minimum": 0 },
      "isValidMime": { "type": "boolean" },
      "hasVariablesResolved": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Edits in template source pane reflect in preview pane within 16ms.
  - Emits valid multipart/alternative MIME message structure with boundary strings.
  - Supports standard vim keybindings (`h`, `j`, `k`, `l`, `i`, `Esc`) in editor mode.

---

### Improvement 243: Bubble Tea Interactive Spool Queue Manager & Inspector
- **ID:** `WARMUP-STUDIO-243`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** When thousands of warmup messages are queued, operators have no easy terminal-based way to inspect individual message headers, flush specific MX queues, or purge stuck items.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Interactive scrollable table displaying Spool ID, Recipient, Target MX, Status (`QUEUED`, `DEFERRED`, `IN_TRANSIT`), Retry Count, and Next Attempt ETA.
    - Quick action bar: `[f] Flush Queue`, `[p] Purge Selected`, `[v] View Raw Header`, `[/] Search Filter`.
    - Modal popup presenting raw RFC 5322 email headers and full delivery trace with syntax highlighting.
  - **Zero-CDN Web UI:**
    - Virtualized table rendering 10,000+ items without DOM lag, styled with pure native CSS variables.
- **CLI Counterpart Command:**
  - `gitmap mail queue list --status "deferred"`
  - `gitmap mail queue flush --target-domain "gmail.com"`
  - `gitmap mail queue inspect --id "msg-98374-af3"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["messageId", "sender", "recipient", "status", "retryCount", "isQueued", "hasDeferredError"],
    "properties": {
      "messageId": { "type": "string" },
      "sender": { "type": "string", "format": "email" },
      "recipient": { "type": "string", "format": "email" },
      "status": { "type": "string", "enum": ["queued", "in_transit", "deferred", "delivered", "bounced"] },
      "retryCount": { "type": "integer", "minimum": 0 },
      "isQueued": { "type": "boolean" },
      "hasDeferredError": { "type": "boolean" },
      "deferredReason": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Keyboard filter updates the list instantly as user types.
  - Selecting `[f]` immediately dispatches background flush worker and animates spinner indicator.
  - Purging confirms via dialog to avoid accidental bulk queue deletion.

---

### Improvement 244: Terminal Real-Time Sparkline Heatmap of Delivery Latency
- **ID:** `WARMUP-STUDIO-244`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Deliverability degradation happens gradually; without immediate visual feedback on SMTP handshake and transmission latency, operators miss subtle server slowdowns.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Unicode sparkline graph (` ▂▃▄▅▆▇█`) rendering rolling 60-second delivery latency per major ISP (Google, Microsoft, Yahoo, Fastmail).
    - Latency color thresholds: `<500ms` (Green), `500ms-2000ms` (Yellow), `>2000ms` (Red).
    - Summary telemetry block showing p50, p95, and p99 delivery latency percentiles.
  - **Zero-CDN Web UI:**
    - SVG sparkline with zero external chart libraries, rendered via pure `<polyline>` points.
- **CLI Counterpart Command:**
  - `gitmap mail studio latency --isp "google" --interval "1s"`
  - `gitmap mail latency report --since "1h" --format json`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["ispName", "p50Ms", "p95Ms", "p99Ms", "isHealthy", "hasHighLatencyAlert", "sampleCount"],
    "properties": {
      "ispName": { "type": "string" },
      "p50Ms": { "type": "number", "minimum": 0 },
      "p95Ms": { "type": "number", "minimum": 0 },
      "p99Ms": { "type": "number", "minimum": 0 },
      "isHealthy": { "type": "boolean" },
      "hasHighLatencyAlert": { "type": "boolean" },
      "sampleCount": { "type": "integer", "minimum": 1 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Sparklines update every 1,000ms without flickering using Bubble Tea ticks.
  - Correct percentile calculations verified with automated test fixtures.
  - TUI consumes less than 1% CPU when idle in background terminal.

---

### Improvement 245: Interactive SMTP Handshake Terminal Debugger
- **ID:** `WARMUP-STUDIO-245`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Debugging raw SMTP dialogs (`EHLO`, `STARTTLS`, `AUTH`, `MAIL FROM`, `RCPT TO`, `DATA`) usually requires cumbersome manual `openssl s_client` or `telnet` commands without structured analysis.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Step-by-step interactive stepper: Step 1 DNS MX resolution, Step 2 TCP Connect, Step 3 TLS Negotiation, Step 4 SMTP Greeting, Step 5 Authentication, Step 6 Envelope Transaction.
    - Protocol stream pane showing client commands in Cyan and server responses in Yellow/Green.
    - Timing column displaying exact round-trip latency for each individual command/response pair.
  - **Zero-CDN Web UI:**
    - Interactive stepper card with expandable JSON inspector for each handshake stage.
- **CLI Counterpart Command:**
  - `gitmap mail test handshake --to "test@gmail.com" --interactive`
  - `gitmap mail test handshake --to "test@outlook.com" --dump-tls`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["targetMx", "port", "isHandshakeSuccess", "hasTlsEstablished", "cipherSuite", "totalDurationMs"],
    "properties": {
      "targetMx": { "type": "string" },
      "port": { "type": "integer" },
      "isHandshakeSuccess": { "type": "boolean" },
      "hasTlsEstablished": { "type": "boolean" },
      "cipherSuite": { "type": "string" },
      "totalDurationMs": { "type": "number", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Interactive mode allows pausing between steps to inspect server headers.
  - Captures and displays any SSL/TLS certificate warnings or cipher mismatches clearly.
  - Exits gracefully on connection abort without hanging the terminal.

---

### Improvement 246: Terminal Account Switching & Multi-Profile HUD
- **ID:** `WARMUP-STUDIO-246`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Operators managing warmup across multiple client accounts, sending domains, and brand profiles must constantly exit and pass different `--config` flags to CLI tools.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Quick-switcher modal overlay triggered by `[Ctrl+A]` or `[:]`: fuzzy searchable list of configured client profiles, active domains, and warmup tiers.
    - Top HUD status bar always visible showing: Active Profile Name, Associated Sending Domain, Current Warmup Tier (e.g. Day 14 / Tier 3), Daily Send Target (e.g. 150/200).
    - Status pill indicating health status (`OPTIMAL`, `WARMING`, `PAUSED`).
  - **Zero-CDN Web UI:**
    - Top-right dropdown selector with profile avatar and quick switch shortcuts.
- **CLI Counterpart Command:**
  - `gitmap mail studio --profile "fintech-corp"`
  - `gitmap mail profiles switch --name "marketing-agency"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["profileId", "profileName", "sendingDomain", "tierLevel", "isActive", "hasWarningStatus"],
    "properties": {
      "profileId": { "type": "string" },
      "profileName": { "type": "string" },
      "sendingDomain": { "type": "string" },
      "tierLevel": { "type": "integer", "minimum": 1 },
      "isActive": { "type": "boolean" },
      "hasWarningStatus": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Switching profiles updates all active TUI panes and metrics in under 30ms.
  - Active profile state persists in local SQLite configuration table.
  - Fuzzy matcher accurately sorts profiles by typing keystrokes.

---

### Improvement 247: Terminal Seed Network Topology & Ping Matrix
- **ID:** `WARMUP-STUDIO-247`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** When seeds fail to respond, operators have difficulty determining whether the issue is ISP-wide, specific to an IMAP provider, or an authentication lock.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - 2D grid matrix: Columns represent Seed Email Providers (Gmail, Outlook, Yahoo, Proton, iCloud, Fastmail); Rows represent Sender Nodes.
    - Cells contain visual latency / health indicators: `●` (Green = healthy, <200ms), `▲` (Yellow = degraded, >1000ms), `✖` (Red = unreachable/locked).
    - Highlighted cursor box allowing directional arrow navigation to inspect any individual cell.
  - **Zero-CDN Web UI:**
    - CSS Grid matrix with zero JavaScript dependencies, styled with hover tooltips and accessible aria labels.
- **CLI Counterpart Command:**
  - `gitmap warmup seed matrix --pool "production"`
  - `gitmap warmup seed probe --seed-id "seed-gmail-42"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["seedId", "provider", "status", "latencyMs", "isReachable", "hasAuthError"],
    "properties": {
      "seedId": { "type": "string" },
      "provider": { "type": "string" },
      "status": { "type": "string", "enum": ["healthy", "degraded", "unreachable", "auth_locked"] },
      "latencyMs": { "type": "number", "minimum": 0 },
      "isReachable": { "type": "boolean" },
      "hasAuthError": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Interactive grid navigation updates detail panel instantly on cell selection.
  - Probes execute asynchronously without freezing the Bubble Tea rendering loop.
  - CLI output formats cleanly into tabular terminal output.

---

### Improvement 248: Terminal Live Tail Log Viewer with Regex Highlight Filters
- **ID:** `WARMUP-STUDIO-248`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Standard `tail -f` outputs walls of unstructured text, making it painful to spot critical email events like rate-limiting, authentication rejections, or spam classifications in real time.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - Real-time scrolling log buffer with configurable backscroll (default 1,000 lines).
    - Auto-colorized log levels: `INFO` (Cyan), `WARN` (Yellow), `ERROR` (Red), `DEBUG` (Dark Gray).
    - Bottom command bar supporting interactive regex filtering (`/regex/`) and inverse filtering (`!pattern`).
    - Pause/Resume scrolling toggle with `[Space]` key, showing unread counter while paused.
  - **Zero-CDN Web UI:**
    - Virtualized terminal log viewer with auto-scroll lock and search highlighting.
- **CLI Counterpart Command:**
  - `gitmap mail studio logs --follow`
  - `gitmap mail studio logs --filter "status=deferred" --level "warn"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["logId", "timestamp", "level", "message", "isFiltered", "hasMatchedPattern"],
    "properties": {
      "logId": { "type": "string" },
      "timestamp": { "type": "string", "format": "date-time" },
      "level": { "type": "string", "enum": ["debug", "info", "warn", "error"] },
      "message": { "type": "string" },
      "isFiltered": { "type": "boolean" },
      "hasMatchedPattern": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Log tailing consumes less than 5MB RAM and zero frame drops during burst logging (1,000 lines/sec).
  - Regex filter matches instantly without locking up the terminal UI.
  - Spacebar toggles pause state reliably, resuming at original tail position.

---

### Improvement 249: Terminal Warmup Schedule Builder & Ramp Curve Visualizer
- **ID:** `WARMUP-STUDIO-249`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** Configuring warmup curves (linear vs exponential vs sigmoid) across 30 days is difficult via flat config files without seeing the projected daily volume progression.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - ASCII bar chart rendering daily projected email volumes across a 30-day timeline.
    - Interactive parameter adjustment sliders: Starting Volume (e.g. 5/day), Daily Increment (e.g. 20%), Max Daily Cap (e.g. 500/day), Curve Type (`linear`, `exponential`, `sigmoid`).
    - Instant recalculation of total projected warmup emails and projected completion date.
  - **Zero-CDN Web UI:**
    - Interactive SVG ramp curve graph with draggable control points.
- **CLI Counterpart Command:**
  - `gitmap warmup curve view --curve "sigmoid" --start 5 --cap 500 --days 30`
  - `gitmap warmup curve apply --domain "send.example.com" --curve "sigmoid"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["curveType", "startVolume", "maxCap", "totalDays", "isApproved", "hasCompletedCalculation", "dailyProjections"],
    "properties": {
      "curveType": { "type": "string", "enum": ["linear", "exponential", "sigmoid", "custom"] },
      "startVolume": { "type": "integer", "minimum": 1 },
      "maxCap": { "type": "integer", "minimum": 1 },
      "totalDays": { "type": "integer", "minimum": 7 },
      "isApproved": { "type": "boolean" },
      "hasCompletedCalculation": { "type": "boolean" },
      "dailyProjections": {
        "type": "array",
        "items": { "type": "integer" }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Parameter changes update the ASCII bar chart in real time.
  - Sigmoid curve calculation adheres to standard mathematical logistic formula.
  - Saving the schedule commits the configuration to the SQLite store.

---

### Improvement 250: Terminal Panic Button & Immediate Kill-Switch TUI Modal
- **ID:** `WARMUP-STUDIO-250`
- **Cluster:** 25 (Interactive Terminal Studio)
- **Problem Statement:** In an emergency (e.g., accidental broadcast to unverified list or sudden spam trap breach), operators need an instant, unmistakable panic kill-switch in the TUI to halt all outbound SMTP traffic.
- **Concrete UI Specification:**
  - **Lipgloss / Bubble Tea TUI:**
    - High-visibility modal triggered by `[Ctrl+K]` or `[!]`: Bright flashing Red Border (`#ef4444`), bold warning prompt requiring confirmation text (`KILL ALL`).
    - Action options: 1. Pause active queues (retain messages), 2. Purge outbound spool (discard unsent), 3. Sever SMTP relay connections immediately.
    - Post-abort summary card showing total halted in-flight messages and preserved spool state.
  - **Zero-CDN Web UI:**
    - Floating red panic button with confirmation dialog and audible/visual pulse.
- **CLI Counterpart Command:**
  - `gitmap mail panic --force --reason "accidental-send"`
  - `gitmap mail panic resume --pool "primary"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["isTriggered", "haltedMessageCount", "abortedAt", "triggerSource", "hasPreservedSpool"],
    "properties": {
      "isTriggered": { "type": "boolean" },
      "haltedMessageCount": { "type": "integer", "minimum": 0 },
      "abortedAt": { "type": "string", "format": "date-time" },
      "triggerSource": { "type": "string" },
      "hasPreservedSpool": { "type": "boolean" },
      "abortReason": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Triggering the kill-switch terminates all active SMTP goroutines in under 15ms.
  - Active queues are placed into `PAUSED` state without data corruption in SQLite.
  - Modal prompts for explicit confirmation string to prevent accidental keystroke triggers.

---

# Cluster 26: P2P Multi-Tenant Reputation Mesh & Anomaly Quarantining (Items 251–260)

---

### Improvement 251: P2P Node Discovery & Reputation Gossip Protocol
- **ID:** `WARMUP-MESH-251`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Individual warmup nodes operate in isolation, unable to share real-time ISP deliverability intelligence (such as temporary Microsoft IP greylisting or Gmail throttling windows) across the fleet.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Interactive P2P mesh topology visualizer showing peer nodes, connected links, latency, and shared reputation gossip packets.
    - Status pill for each node: `SYNCED` (Green), `DISCOVERING` (Cyan), `OUT_OF_SYNC` (Yellow), `QUARANTINED` (Red).
    - Gossip packet log stream showing peer-to-peer reputation score exchange.
  - **Lipgloss Terminal UX:**
    - Styled node list showing Node ID, Peer Address, Ping RTT, and Number of Shared Reputation Vectors.
- **CLI Counterpart Command:**
  - `gitmap mesh peers list`
  - `gitmap mesh gossip status`
  - `gitmap mesh connect --peer "node-02.fleet.internal:9443"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["nodeId", "peerAddress", "isPeerConnected", "hasValidGossipState", "latencyMs", "reputationVectorsShared"],
    "properties": {
      "nodeId": { "type": "string" },
      "peerAddress": { "type": "string" },
      "isPeerConnected": { "type": "boolean" },
      "hasValidGossipState": { "type": "boolean" },
      "latencyMs": { "type": "number", "minimum": 0 },
      "reputationVectorsShared": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Nodes discover peers via lightweight mDNS / local bootstrap registry within 5 seconds.
  - Gossip protocol exchanges ISP health vectors every 30 seconds using minimal binary payloads (<2KB).
  - CLI command lists active peers with positive booleans and RTT metrics.

---

### Improvement 252: Byzantine Fault-Tolerant Reputation Consensus Engine
- **ID:** `WARMUP-MESH-252`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Compromised, misconfigured, or malfunctioning sender nodes could broadcast false deliverability metrics, poisoning the shared reputation pool and triggering unnecessary throttling across healthy nodes.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Consensus round telemetry panel displaying current epoch number, validator node votes, and quorum status (e.g. 7/10 nodes agreed).
    - Visual timeline of consensus rounds with green blocks for reached quorum and red blocks for rejected outlier proposals.
    - Outlier node warning card highlighting nodes whose submitted metrics deviate >3 standard deviations from the fleet mean.
  - **Lipgloss Terminal UX:**
    - Table displaying active epoch, proposing leader node, quorum percentage, and consensus status.
- **CLI Counterpart Command:**
  - `gitmap mesh consensus status`
  - `gitmap mesh consensus verify-epoch --epoch 1042`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["epochId", "quorumRatio", "isConsensusReached", "hasOutlierDetected", "votingNodeCount"],
    "properties": {
      "epochId": { "type": "integer", "minimum": 1 },
      "quorumRatio": { "type": "number", "minimum": 0, "maximum": 1 },
      "isConsensusReached": { "type": "boolean" },
      "hasOutlierDetected": { "type": "boolean" },
      "votingNodeCount": { "type": "integer", "minimum": 1 },
      "detectedOutlierNode": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Requires 2/3 + 1 supermajority of peer votes to ratify global ISP status changes.
  - Malicious or corrupted metrics from up to 33% of nodes are automatically rejected without affecting consensus.
  - Verified epoch records write directly to local SQLite ledger.

---

### Improvement 253: Automated Peer Anomaly Quarantining & Isolation Daemon
- **ID:** `WARMUP-MESH-253`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** When a specific node experiences a critical compromise (e.g. private key leak or runaway bounce rate), it must be quarantined immediately from the mesh to prevent cross-tenant contamination.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Quarantined nodes drawer with red hazard perimeter styling.
    - Isolation trigger log detailing the exact anomaly event (e.g., sudden 40% bounce rate breach on Outlook).
    - Secure manual re-admission workflow with cryptographic identity re-verification.
  - **Lipgloss Terminal UX:**
    - High-visibility banner `[QUARANTINE ENFORCED]` with node identifier, reason code, and duration of isolation.
- **CLI Counterpart Command:**
  - `gitmap mesh quarantine list`
  - `gitmap mesh quarantine isolate --node-id "node-chicago-04" --reason "sudden-bounce-spike"`
  - `gitmap mesh quarantine release --node-id "node-chicago-04"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["nodeId", "isolationReason", "isQuarantined", "hasSafetyLock", "isolatedAt"],
    "properties": {
      "nodeId": { "type": "string" },
      "isolationReason": { "type": "string" },
      "isQuarantined": { "type": "boolean" },
      "hasSafetyLock": { "type": "boolean" },
      "isolatedAt": { "type": "string", "format": "date-time" },
      "quarantineDurationSeconds": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Quarantining an anomaly node severs all peer gossip and relay routing connections within 100ms.
  - Quarantined nodes cannot participate in consensus voting until formal re-admission.
  - CLI command requires explicit `--reason` flag for administrative audit logging.

---

### Improvement 254: Cryptographic Multi-Tenant Reputation Vector Signer
- **ID:** `WARMUP-MESH-254`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Metrics shared across multi-tenant clusters can be tampered with or spoofed in transit if transmitted without rigorous asymmetric signatures.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Key management dashboard showing Ed25519 signing public key fingerprint, creation date, and signed payload counter.
    - Signature verification inspector allowing operators to paste a hex payload and verify cryptographic validity.
    - Visual badge for payload authenticity: `VALID_ED25519_SIG` (Green) vs `INVALID_SIGNATURE` (Red).
  - **Lipgloss Terminal UX:**
    - Security summary card displaying node public key, active cipher suite, and count of verified vs rejected peer signatures.
- **CLI Counterpart Command:**
  - `gitmap mesh crypto keys show`
  - `gitmap mesh crypto verify-payload --input "./data/gossip/packet-204.bin"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["keyId", "publicKeyFingerprint", "isValidSignature", "hasCryptoVerificationPassed", "signedAt"],
    "properties": {
      "keyId": { "type": "string" },
      "publicKeyFingerprint": { "type": "string" },
      "isValidSignature": { "type": "boolean" },
      "hasCryptoVerificationPassed": { "type": "boolean" },
      "signedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - All gossip messages are signed using Ed25519 keys stored in the secure local key store.
  - Payloads with missing or invalid signatures are discarded immediately with a security log entry.
  - Verification benchmarks complete in under 50 microseconds per payload.

---

### Improvement 255: Global ISP Rate-Limit Shared Knowledge Base
- **ID:** `WARMUP-MESH-255`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Mailbox providers change their connection rate limits dynamically depending on global server load; individual nodes take hours to rediscover these limits through trial and error.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Global ISP rate-limit heatmap displaying current safe connection rates (messages per minute, max connections) per destination MX domain.
    - Real-time confidence score based on the volume of observations across all participating mesh nodes.
    - Toggle to sync local sending throttles automatically with the global mesh recommendation.
  - **Lipgloss Terminal UX:**
    - Tabular summary of top destination domains (`gmail.com`, `outlook.com`, `yahoo.com`) with recommended concurrent connections and hourly limits.
- **CLI Counterpart Command:**
  - `gitmap mesh limits list`
  - `gitmap mesh limits sync --provider "all"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["destinationDomain", "recommendedMaxConnections", "messagesPerHourLimit", "confidenceScore", "isGloballySynced", "hasHighConfidence"],
    "properties": {
      "destinationDomain": { "type": "string" },
      "recommendedMaxConnections": { "type": "integer", "minimum": 1 },
      "messagesPerHourLimit": { "type": "integer", "minimum": 1 },
      "confidenceScore": { "type": "number", "minimum": 0, "maximum": 1 },
      "isGloballySynced": { "type": "boolean" },
      "hasHighConfidence": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Shared knowledge base dynamically aggregates observations from all peer nodes.
  - Auto-sync updates local dispatcher throttle rates within 1 second of consensus update.
  - CLI `sync` outputs confirmation of updated destination rate bounds.

---

### Improvement 256: Multi-Tenant Tenant Separation & Fair-Share Queue Allocator
- **ID:** `WARMUP-MESH-256`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** In a shared infrastructure node, one aggressive tenant running a high-volume warmup can monopolize outbound SMTP connections, starving lower-volume tenants.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Deficit Round Robin (DRR) fair-share queue allocation visualizer showing live bandwidth and connection slot distribution per tenant.
    - Tenant allocation slider allowing administrators to assign weight or minimum guaranteed throughput shares.
    - Starvation warning alert triggered if any tenant's queued messages experience wait times > 180 seconds.
  - **Lipgloss Terminal UX:**
    - Progress bars for each tenant showing used vs allocated connection quotas with Lipgloss gradient coloring.
- **CLI Counterpart Command:**
  - `gitmap mail multi-tenant status`
  - `gitmap mail multi-tenant set-quota --tenant-id "tenant-alpha" --weight 2`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["tenantId", "allocatedWeight", "activeConnections", "isThrottled", "hasSlaStarvation", "queuedMessageCount"],
    "properties": {
      "tenantId": { "type": "string" },
      "allocatedWeight": { "type": "integer", "minimum": 1 },
      "activeConnections": { "type": "integer", "minimum": 0 },
      "isThrottled": { "type": "boolean" },
      "hasSlaStarvation": { "type": "boolean" },
      "queuedMessageCount": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - DRR scheduler strictly enforces fair connection dispatch across all active tenants.
  - No tenant can consume more than their assigned max share while other tenants have pending messages.
  - Changes to tenant weights take effect immediately on next queue round.

---

### Improvement 257: Mesh Zero-Knowledge Proofs for Anti-Abuse Compliance
- **ID:** `WARMUP-MESH-257`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Tenants in a multi-tenant warmup mesh want proof that other participating tenants adhere to bounce and complaint limits without revealing proprietary email recipient lists or campaign contents.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Zero-Knowledge Compliance verification panel displaying cryptographic validity proofs (e.g., zk-SNARK / bulletproofs verifying bounce rate < 1.0%).
    - Status pill: `ZK_VERIFIED_COMPLIANT` (Green) vs `PROOF_EXPIRED` (Yellow) vs `VERIFICATION_FAILED` (Red).
    - Audit drawer displaying public proof artifacts and verification timestamp.
  - **Lipgloss Terminal UX:**
    - High-tech terminal status view showing tenant anonymized hash, proof commitment ID, and verification result.
- **CLI Counterpart Command:**
  - `gitmap mesh zk verify --tenant-hash "a8f3b...9e"`
  - `gitmap mesh zk generate-proof --tenant-id "tenant-primary"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["tenantHash", "proofCommitment", "isZkVerified", "hasValidCompliance", "verifiedAt"],
    "properties": {
      "tenantHash": { "type": "string" },
      "proofCommitment": { "type": "string" },
      "isZkVerified": { "type": "boolean" },
      "hasValidCompliance": { "type": "boolean" },
      "verifiedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Proof confirms bounce rate is below threshold without exposing individual email addresses or totals.
  - Proof verification executes in under 20ms using local pure Go cryptography routines.
  - CLI prints clean verification summary with zero data leakage.

---

### Improvement 258: P2P Distributed Seed Pool Rotation & Cross-Account Reply Mesh
- **ID:** `WARMUP-MESH-258`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Static seed networks create predictable interaction patterns that anti-spam machine learning models (like Gmail's Spam Brain) easily fingerprint and devalue.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Cross-account conversational mesh visualization displaying randomized seed interaction pairs across different accounts.
    - Conversation depth distribution chart (Thread depth: 1 reply, 2 replies, 3+ back-and-forth exchanges).
    - Entropy randomized scheduler setting realistic business-hour reply delays (e.g. 15 to 140 minutes).
  - **Lipgloss Terminal UX:**
    - Lipgloss table showing active conversational threads, participating seeds, current turn, and next scheduled reply timestamp.
- **CLI Counterpart Command:**
  - `gitmap warmup mesh-seeds status`
  - `gitmap warmup mesh-seeds schedule-thread --account-a "tenant-1" --account-b "tenant-2" --depth 3`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["threadId", "initiatorSeed", "responderSeed", "threadDepth", "isActive", "hasNextReplyScheduled"],
    "properties": {
      "threadId": { "type": "string" },
      "initiatorSeed": { "type": "string", "format": "email" },
      "responderSeed": { "type": "string", "format": "email" },
      "threadDepth": { "type": "integer", "minimum": 1 },
      "isActive": { "type": "boolean" },
      "hasNextReplyScheduled": { "type": "boolean" },
      "nextScheduledReplyAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Scheduled replies simulate realistic human response intervals with gaussian distribution.
  - Reply headers preserve valid `In-Reply-To` and `References` RFC 5322 headers.
  - CLI lists active mesh threads with accurate status booleans.

---

### Improvement 259: Automated IP Pool Shuffling & Coordinated Re-Balancing
- **ID:** `WARMUP-MESH-259`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** Sending IP pools experience uneven reputation accumulation; fresh IPs receive too little traffic while mature IPs bear excessive load, slowing fleet-wide warmup progression.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - IP Pool re-balancing dashboard displaying load distribution across Tier 1 (Cold), Tier 2 (Warming), and Tier 3 (Production) IP sets.
    - Re-balancing slider allowing automated traffic weight shifting based on current sender score.
    - One-click trigger for coordinated fleet re-shuffle with gradual weight migration.
  - **Lipgloss Terminal UX:**
    - Columnar table of IP addresses with assigned warmup tier, send weight percentage, and current daily utilization bar.
- **CLI Counterpart Command:**
  - `gitmap mail pool rebalance --dry-run`
  - `gitmap mail pool rebalance apply --strategy "proportional-reputation"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["poolId", "ipAddress", "currentTier", "assignedWeight", "isBalanced", "hasSufficientTraffic"],
    "properties": {
      "poolId": { "type": "string" },
      "ipAddress": { "type": "string" },
      "currentTier": { "type": "string", "enum": ["cold", "warming", "production", "cooldown"] },
      "assignedWeight": { "type": "number", "minimum": 0, "maximum": 1 },
      "isBalanced": { "type": "boolean" },
      "hasSufficientTraffic": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Re-balancing algorithm adjusts outbound traffic routing smoothly without dropping in-flight connections.
  - Dry-run mode previews proposed weight adjustments without applying modifications to routing tables.
  - All status properties strictly use positive `is`/`has` booleans.

---

### Improvement 260: Mesh Global Incident Broadcast & Post-Mortem Archiver
- **ID:** `WARMUP-MESH-260`
- **Cluster:** 26 (P2P Multi-Tenant Reputation Mesh)
- **Problem Statement:** After a major deliverability incident (e.g. mass blocklist event or network outage), incident details and remediation steps remain siloed on the affected node rather than preserved for collective fleet learning.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Incident broadcast timeline displaying root cause summary, affected IP pools, duration, and post-mortem resolution notes.
    - Cryptographically signed post-mortem document viewer with export options (JSON, Markdown).
    - Status badge: `INCIDENT_ACTIVE` (Red), `REMEDIATING` (Amber), `RESOLVED_AND_ARCHIVED` (Green).
  - **Lipgloss Terminal UX:**
    - Lipgloss-bordered post-mortem report card with formatted Markdown text and key metrics summary table.
- **CLI Counterpart Command:**
  - `gitmap mesh incident broadcast --title "Yahoo 421 Rate Limit Wave" --severity "critical"`
  - `gitmap mesh incident list --status "resolved"`
  - `gitmap mesh incident view --id "inc-2026-10-01"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["incidentId", "title", "severity", "isResolved", "hasPostMortemArchived", "createdAt", "resolvedAt"],
    "properties": {
      "incidentId": { "type": "string" },
      "title": { "type": "string" },
      "severity": { "type": "string", "enum": ["low", "medium", "high", "critical"] },
      "isResolved": { "type": "boolean" },
      "hasPostMortemArchived": { "type": "boolean" },
      "createdAt": { "type": "string", "format": "date-time" },
      "resolvedAt": { "type": "string", "format": "date-time" },
      "rootCause": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Broadcast distributes incident notification across all connected peer nodes in under 2 seconds.
  - Incident post-mortem documents persist in local and replicated Split-DB SQLite storage.
  - CLI `view` prints full structured report with zero unhandled exceptions.
