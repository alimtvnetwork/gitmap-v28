# Specification 340: Improvements 331 to 360 (Self-Healing SMTP Circuits, Bubble Tea Studios & Zero-Loss WAL Delivery)

- **Sequence:** 340
- **Slug:** `340-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/03-improvements-331-to-360-self-healing-circuits-tui.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 2

---

## Improvement Catalog Summary (331 - 360)

| Cluster | IDs | Scope |
| :--- | :--- | :--- |
| **Cluster 34** | Items 331–340 | Self-Healing SMTP Circuits & Auto-Remediation |
| **Cluster 35** | Items 341–350 | High-Velocity Bubble Tea / Lipgloss Terminal Studios |
| **Cluster 36** | Items 351–360 | Resilient Local Spooling & Zero-Loss WAL Delivery |

---

# Cluster 34: Self-Healing SMTP Circuits & Auto-Remediation (Items 331–340)

---

### Improvement 331: Adaptive Circuit Breaker with Dynamic Tripping & Exponential Probe Reseed
- **ID:** `WARMUP-HEAL-331`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** When downstream SMTP relays or ISP mail exchangers experience fatal connection drops, connection resets, or cascading 5xx rejections, naive deliverability runners continue blasting envelopes, triggering full IP blocklisting, exhausted connection sockets, and high failure counts without trip isolation.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Tri-state animated SVG circuit visualizer (`.circuit-state-node`) displaying states: CLOSED (emerald `#10b981`), OPEN (crimson `#ef4444`), and HALF-OPEN (amber `#f59e0b`).
    - Live error window threshold dial with drag-to-tune trip sensitivity (default: 5 consecutive failures or >25% error rate across a rolling 30-second window).
    - Half-open synthetic probe progress bar showing active sample probe health and countdown timer until circuit full closure.
  - **Lipgloss Terminal UX:**
    - High-visibility status badge with `lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#ef4444")).Foreground(lipgloss.Color("#ffffff")).Padding(0, 1).Render("CIRCUIT OPEN")`.
    - Live terminal countdown gauge showing probe reseed interval (`[=====>----] 12s to Half-Open probe`).
- **CLI Counterpart Command:**
  - `gitmap mail circuit status --pool "us-west-relay"`
  - `gitmap mail circuit trip --relay "smtp.us-west-1.relay.net" --reason "forced-maintenance"`
  - `gitmap mail circuit reset --relay "smtp.us-west-1.relay.net"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["relayHost", "poolId", "circuitState", "isOpen", "isHalfOpen", "hasProbeSucceeded", "consecutiveFailures", "probeIntervalSeconds"],
    "properties": {
      "relayHost": { "type": "string" },
      "poolId": { "type": "string" },
      "circuitState": { "type": "string", "enum": ["closed", "half_open", "open"] },
      "isOpen": { "type": "boolean" },
      "isHalfOpen": { "type": "boolean" },
      "hasProbeSucceeded": { "type": "boolean" },
      "consecutiveFailures": { "type": "integer", "minimum": 0 },
      "probeIntervalSeconds": { "type": "integer", "minimum": 1 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Upon 5 consecutive connection timeout or 5xx response faults, circuit transitions to `OPEN` state in <10ms and reroutes traffic.
  - While `OPEN`, queued messages bypass the failed relay to fallback pools or disk-spool without network socket creation.
  - Half-open probe automatically fires after configured backoff; on success, circuit closes and normal traffic resumes.

---

### Improvement 332: Real-time Reverse DNS (rDNS/PTR) Drift Auto-Correction & Registrar Sync
- **ID:** `WARMUP-HEAL-332`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Outbound warm-up nodes assigned floating IPs frequently suffer from mismatched rDNS/PTR records when cloud providers remap interfaces or DNS lease TTLs expire, causing severe Gmail and Yahoo policy rejections (`550 5.7.1 The IP address sending this message does not have a PTR record`).
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Side-by-side PTR matrix displaying IP Address, Detected PTR FQDN, Forward A Record, and Verification Badge (`isPTRValid`).
    - One-click "Auto-Fix via Cloud Provider API" modal for AWS Route53, Cloudflare, and Hetzner DNS.
    - Historical PTR drift timeline recording lease renegotiations and latency dips.
  - **Lipgloss Terminal UX:**
    - Dual-column tabular view using Lipgloss borders: Left shows IP and current PTR, right displays forward resolution comparison with green checkmark or red cross indicators.
- **CLI Counterpart Command:**
  - `gitmap warmup rdns verify --ip "198.51.100.24"`
  - `gitmap warmup rdns remediate --ip "198.51.100.24" --expected-fqdn "mail01.sender.org"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["ipAddress", "ptrRecord", "forwardARecord", "isPtrMatching", "hasForwardMatch", "lastVerifiedAt"],
    "properties": {
      "ipAddress": { "type": "string" },
      "ptrRecord": { "type": "string" },
      "forwardARecord": { "type": "string" },
      "isPtrMatching": { "type": "boolean" },
      "hasForwardMatch": { "type": "boolean" },
      "lastVerifiedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - The checker verifies both PTR lookup (`dig -x IP`) and forward lookup (`dig A PTR_HOST`) to enforce strict Forward-Confirmed Reverse DNS (FCrDNS).
  - Flags mismatches within 500ms of IP allocation and alerts the Lemon rescue coordinator before dispatching outbound warm-up batches.

---

### Improvement 333: Autonomous SPF Macro Alignment & DNS Lookup Limit Compressor
- **ID:** `WARMUP-HEAL-333`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Rapidly scaling multi-tenant sending nodes causes customer SPF records to exceed the RFC 7208 strict 10-DNS-lookup limit (`PermError: SPF permanent error: too many DNS lookups`), leading receiving servers to drop emails straight to spam.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Visual SPF lookup odometer meter displaying current lookup count out of 10 (`7/10 Lookups Used`).
    - Tree diagram of nested `include:` directives with expandable lookup depths and resolution latencies.
    - Interactive "Flatten SPF" action button that compresses complex nested chains into unified CIDR `ip4:` and `ip6:` blocks.
  - **Lipgloss Terminal UX:**
    - Lipgloss bar chart displaying lookup count: Yellow bar when 7-9 lookups, flashing red border when >= 10 lookups.
- **CLI Counterpart Command:**
  - `gitmap mail spf audit --domain "example.com"`
  - `gitmap mail spf flatten --domain "example.com" --out-txt`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "rawSpfRecord", "dnsLookupCount", "isUnderLimit", "hasNestedIncludes", "suggestedFlattenedSpf"],
    "properties": {
      "domain": { "type": "string" },
      "rawSpfRecord": { "type": "string" },
      "dnsLookupCount": { "type": "integer", "minimum": 0, "maximum": 25 },
      "isUnderLimit": { "type": "boolean" },
      "hasNestedIncludes": { "type": "boolean" },
      "suggestedFlattenedSpf": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - SPF parser recurses through all `include`, `a`, `mx`, `ptr`, and `exists` mechanisms, counting lookups according to RFC 7208.
  - Fails validation and generates auto-compressed CIDR string if lookup count is >= 10.

---

### Improvement 334: Real-Time DMARC Aggregate Feedback (RUA/RUF) Anomaly Sentinel
- **ID:** `WARMUP-HEAL-334`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** DMARC aggregate XML reports delivered to mailboxes are processed manually or via delayed batch jobs, allowing domain spoofing, alignment failures, and sender reputation crashes to go unnoticed for days before mitigation occurs.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Live XML ingestion dashboard (`.dmarc-rua-feed`) auto-parsing incoming reports from Google, Microsoft, and Yahoo.
    - Domain compliance donut chart showing Percentage Aligned vs DKIM Fail vs SPF Fail.
    - Policy enforcement slider to smoothly upgrade policies from `p=none` to `p=quarantine` to `p=reject` with projected traffic pass rates.
  - **Lipgloss Terminal UX:**
    - Lipgloss sparkline table displaying 7-day DMARC compliance trends across primary destination mailbox providers.
- **CLI Counterpart Command:**
  - `gitmap mail dmarc ingest --report-file "./dmarc-rua-google.xml"`
  - `gitmap mail dmarc report --domain "example.com" --window "24h"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "policyApplied", "totalMessages", "alignedMessages", "isCompliant", "hasSpoofingDetected", "topFailingIps"],
    "properties": {
      "domain": { "type": "string" },
      "policyApplied": { "type": "string", "enum": ["none", "quarantine", "reject"] },
      "totalMessages": { "type": "integer", "minimum": 0 },
      "alignedMessages": { "type": "integer", "minimum": 0 },
      "isCompliant": { "type": "boolean" },
      "hasSpoofingDetected": { "type": "boolean" },
      "topFailingIps": { "type": "array", "items": { "type": "string" } }
    }
  }
  ```
- **Acceptance Criteria:**
  - Ingests raw `.xml` and `.xml.gz` RUA reports directly into SQLite database in <250ms.
  - Flags spoofing spikes (>5% unaligned volume from unknown IPs) and sends instant remediation alerts to operator terminals.

---

### Improvement 335: Automated TLS Cipher Suite Negotiation Downgrade Guard & MX Probe
- **ID:** `WARMUP-HEAL-335`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Strict outbound STARTTLS configurations fail when sending to legacy governmental or educational recipient MX servers that only support TLS 1.2 or obsolete ciphers, causing unhandled handshake termination without safe fallback negotiation.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - TLS compatibility matrix for target domains listing supported protocols (TLS 1.3, TLS 1.2), supported ciphers, and certificate authority validity.
    - Policy mode selector: Strict MTA-STS (Require TLS 1.3/1.2 with valid CA) vs Opportunistic TLS with cipher negotiation fallback.
    - Warning banner for weak ciphers or expired remote certificates.
  - **Lipgloss Terminal UX:**
    - Colored handshake cipher details: Cyan for `TLS_AES_256_GCM_SHA384`, Yellow for `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`, Red for plain unencrypted fallback.
- **CLI Counterpart Command:**
  - `gitmap mail tls probe --mx "mail.edu-domain.org"`
  - `gitmap mail tls policy --domain "edu-domain.org" --mode "opportunistic"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["mxHost", "tlsVersion", "cipherSuite", "isTlsActive", "hasCertExpired", "isMtaStsCompliant"],
    "properties": {
      "mxHost": { "type": "string" },
      "tlsVersion": { "type": "string" },
      "cipherSuite": { "type": "string" },
      "isTlsActive": { "type": "boolean" },
      "hasCertExpired": { "type": "boolean" },
      "isMtaStsCompliant": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Probes remote MX server with STARTTLS; accurately identifies TLS 1.3 and 1.2 capabilities without dropping connection.
  - Enforces MTA-STS policy when published by recipient; prevents downgrade attacks to unencrypted plaintext if STS mode is `enforce`.

---

### Improvement 336: SMTP Response Classifier with RFC 3463 Enhanced Status Code Resolver
- **ID:** `WARMUP-HEAL-336`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Standard SMTP client libraries only inspect 3-digit basic response codes (e.g., `550`), ignoring the RFC 3463 enhanced status codes (e.g., `5.7.26`, `5.1.1`, `5.4.1`), causing mailbox provider-specific nuances (bad mailbox vs authentication reject vs IP reputation block) to be grouped indiscriminately into generic delivery errors.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Interactive response code decoder tool: paste raw SMTP banner line and view instant breakdown of Class (Permanent/Transient), Subject, and Detail meanings.
    - Provider-specific badge tagger (`Gmail DMARC Policy`, `Outlook SNDS Spam Trap`, `Yahoo Rate Block`).
    - Real-time stream of classified rejections with targeted one-click mitigation suggestions.
  - **Lipgloss Terminal UX:**
    - Formatted terminal breakdown showing: `[Raw Code: 550] [Enhanced: 5.7.26] [Category: AUTH_DKIM_FAIL]`.
- **CLI Counterpart Command:**
  - `gitmap mail classify "550 5.7.26 This message does not have authentication info"`
  - `gitmap mail errors list --category "authentication"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["rawBanner", "statusCode", "enhancedCode", "category", "isPermanentFailure", "hasReputationImpact"],
    "properties": {
      "rawBanner": { "type": "string" },
      "statusCode": { "type": "integer" },
      "enhancedCode": { "type": "string" },
      "category": { "type": "string", "enum": ["mailbox_unavailable", "authentication_fail", "reputation_block", "content_filtered", "rate_limit", "system_fault"] },
      "isPermanentFailure": { "type": "boolean" },
      "hasReputationImpact": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Regex parser extracts both basic code (2xx, 4xx, 5xx) and enhanced dotted code (X.X.X) with 100% precision.
  - Maps code to standardized internal enum (`ErrCategoryAuth`, `ErrCategoryRateLimit`, `ErrCategoryMailboxMissing`) to trigger targeted recovery workflows.

---

### Improvement 337: Automated Feedback Loop (FBL) ARF Parser & Auto-Unsubscribe Dispatcher
- **ID:** `WARMUP-HEAL-337`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** When recipients click "Report Spam", Feedback Loop (FBL) emails formatted as Abuse Reporting Format (ARF, RFC 5965) bounce back to abuse mailboxes; failure to immediately suppress these recipients results in exponential reputation penalties and account suspension.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Incoming ARF incident viewer showing Reporter MX, Offending Message ID, Target Recipient, Feedback Type (`abuse`), and Original Headers.
    - Global suppression ledger table with instant search and filter by recipient email, domain, or timestamp.
    - Auto-unsubscribe confirmation modal indicating synchronous deletion/suppression from all active warm-up campaigns.
  - **Lipgloss Terminal UX:**
    - Red alert box rendered with Lipgloss on ARF detection: `[SPAM COMPLAINT DETECTED] recipient: user@yahoo.com -> Auto-suppressed`.
- **CLI Counterpart Command:**
  - `gitmap mail fbl ingest --eml-path "./complaint-arf.eml"`
  - `gitmap mail suppression add --email "user@yahoo.com" --reason "fbl_abuse"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["feedbackType", "reportedEmail", "originalMessageId", "reportingMta", "isSuppressed", "hasUnsubscribedImmediately"],
    "properties": {
      "feedbackType": { "type": "string", "enum": ["abuse", "fraud", "virus", "other"] },
      "reportedEmail": { "type": "string" },
      "originalMessageId": { "type": "string" },
      "reportingMta": { "type": "string" },
      "isSuppressed": { "type": "boolean" },
      "hasUnsubscribedImmediately": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Parses multipart/report MIME messages conforming to RFC 5965; extracts original recipient and message ID.
  - Writes suppression entry into `gitmap.db` within 5ms; guarantees zero future messages will be dispatched to that address.

---

### Improvement 338: Real-Time DNSBL / RBL Blacklist Monitor with Auto-Egress IP Drain
- **ID:** `WARMUP-HEAL-338`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Sending warm-up traffic from an IP that gets listed on major blacklists (Spamhaus ZEN, Barracuda, SORBS, SpamCop) contaminates the entire sending pool and burns sender domain reputation before admins can manually react.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - IP Pool Blacklist Radar grid: matrix of all egress IP addresses across 25+ reputable DNSBL providers.
    - Status indicators: Clean (Green `#10b981`), Listed (Red `#ef4444`), Query Timeout (Yellow `#f59e0b`).
    - Automated "Emergency Drain" toggle: instantly removes any listed IP from the active round-robin rotation and queues warm-up traffic through healthy standby IPs.
  - **Lipgloss Terminal UX:**
    - Tabular Lipgloss terminal grid displaying IP listings with query response codes (e.g., `127.0.0.2` for Spamhaus SBL).
- **CLI Counterpart Command:**
  - `gitmap warmup rbl check --ip "198.51.100.50"`
  - `gitmap warmup pool drain --ip "198.51.100.50" --reason "spamhaus_listed"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["ipAddress", "dnsblHost", "isListed", "hasActiveDelistingRequest", "listingCode", "responseTimeMs"],
    "properties": {
      "ipAddress": { "type": "string" },
      "dnsblHost": { "type": "string" },
      "isListed": { "type": "boolean" },
      "hasActiveDelistingRequest": { "type": "boolean" },
      "listingCode": { "type": "string" },
      "responseTimeMs": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Performs asynchronous parallel DNS lookups against top 20 DNSBL providers in under 2 seconds.
  - Automatically invokes IP drain in <100ms when Spamhaus or Barracuda returns positive 127.0.0.x response.

---

### Improvement 339: MX Connection Pool Throttler with Fair-Queuing Destination Leaky Bucket
- **ID:** `WARMUP-HEAL-339`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Sending threads burst hundreds of parallel connections to a single destination MX cluster (e.g., `gmail-smtp-in.l.google.com`), triggering anti-DDoS connection rate limits (`421 4.7.0 Try again later, closing connection`), crashing delivery speed for all messages to that provider.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Per-destination leaky bucket rate monitor: shows active concurrent connections, tokens available, and queue backlog per MX domain.
    - Interactive rate slider allowing operators to adjust max concurrent connections (e.g., max 10 to Yahoo, max 20 to Gmail) with live visual pressure gauge.
    - Priority queue table showing VIP transactional vs background warm-up traffic slots.
  - **Lipgloss Terminal UX:**
    - ASCII connection pool meter rendering real-time capacity: `[||||||....] 6/10 active connections to google.com`.
- **CLI Counterpart Command:**
  - `gitmap mail pool config --mx-domain "google.com" --max-connections 15 --rate-per-min 120`
  - `gitmap mail pool metrics`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["mxDomain", "activeConnections", "maxConnections", "tokensRemaining", "isThrottled", "hasAvailableCapacity"],
    "properties": {
      "mxDomain": { "type": "string" },
      "activeConnections": { "type": "integer", "minimum": 0 },
      "maxConnections": { "type": "integer", "minimum": 1 },
      "tokensRemaining": { "type": "integer", "minimum": 0 },
      "isThrottled": { "type": "boolean" },
      "hasAvailableCapacity": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Guarantees strict per-MX-domain connection concurrency caps using Go token bucket or semaphore primitives.
  - Queued envelopes wait without dropping until tokens refresh, eliminating 421 concurrency rejections.

---

### Improvement 340: SMTP Server Capability Caching & ESMTP Extension Fingerprinter
- **ID:** `WARMUP-HEAL-340`
- **Cluster:** 34 (Self-Healing SMTP Circuits & Auto-Remediation)
- **Problem Statement:** Every outbound SMTP connection redundantly issues `EHLO` and negotiates capabilities (`PIPELINING`, `SIZE`, `CHUNKING`, `8BITMIME`, `SMTPUTF8`, `DSN`) from scratch, wasting round-trip latencies on every single message delivery.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Cached MX Capability Viewer: table of remote MX hostnames, supported ESMTP extensions, advertised max message size, and cache expiration TTL.
    - Capability matrix checkboxes showing supported protocol features across top destination providers.
    - Manual cache purge button to force immediate re-fingerprinting on next connection.
  - **Lipgloss Terminal UX:**
    - Structured capability tag list: `[PIPELINING: YES] [CHUNKING/BDAT: NO] [8BITMIME: YES] [SIZE: 50MB]`.
- **CLI Counterpart Command:**
  - `gitmap mail ehlo inspect --host "mx1.mail.icloud.com"`
  - `gitmap mail ehlo cache-clear --host "mx1.mail.icloud.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["mxHost", "supportsPipelining", "supportsChunking", "supports8BitMime", "supportsSmtpUtf8", "isCached", "hasSizeLimit", "maxSizeBytes"],
    "properties": {
      "mxHost": { "type": "string" },
      "supportsPipelining": { "type": "boolean" },
      "supportsChunking": { "type": "boolean" },
      "supports8BitMime": { "type": "boolean" },
      "supportsSmtpUtf8": { "type": "boolean" },
      "isCached": { "type": "boolean" },
      "hasSizeLimit": { "type": "boolean" },
      "maxSizeBytes": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Caches `EHLO` extension response for remote MX hosts in local memory/SQLite for 24 hours.
  - If remote MX supports `PIPELINING`, pipeline `MAIL FROM` and `RCPT TO` commands into a single round-trip packet, reducing handshake latency by >= 40%.

---

# Cluster 35: High-Velocity Bubble Tea / Lipgloss Terminal Studios (Items 341–350)

---

### Improvement 341: Interactive Bubble Tea Mail Studio (`gitmap mail studio`)
- **ID:** `WARMUP-TUI-341`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Operators navigating complex email warm-up schedules, DNS validations, and live spool health are forced to context-switch across dozens of disjoint CLI subcommands, slowing operational debugging and incident response.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - N/A (Dedicated Fullscreen Terminal TUI Studio application).
  - **Lipgloss Terminal UX:**
    - Fullscreen 3-column layout built with Charmbracelet Bubble Tea & Lipgloss:
      - Left column (width 25%): Navigational Tree (Campaigns, Mailbox Pools, Live Spool, DNS Radar, Blacklist Watcher, WAL Journal).
      - Center column (width 50%): Dynamic Workspace View showing selected campaign schedule curves, throughput meters, and message composer.
      - Right column (width 25%): Live Inspector displaying real-time SMTP event stream, peer nodes, and circuit breaker health gauges.
    - Hotkey bar footer (`[q] Quit | [tab] Switch Pane | [r] Refresh | [space] Pause/Resume | [/] Filter | [?] Help`).
- **CLI Counterpart Command:**
  - `gitmap mail studio`
  - `gitmap mail studio --theme "monokai" --fullscreen`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["activePane", "selectedEntityId", "isStudioActive", "hasFilterApplied", "themeName", "terminalWidth", "terminalHeight"],
    "properties": {
      "activePane": { "type": "string", "enum": ["nav", "workspace", "inspector"] },
      "selectedEntityId": { "type": "string" },
      "isStudioActive": { "type": "boolean" },
      "hasFilterApplied": { "type": "boolean" },
      "themeName": { "type": "string" },
      "terminalWidth": { "type": "integer", "minimum": 80 },
      "terminalHeight": { "type": "integer", "minimum": 24 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Runs inside modern ANSI terminals via Bubble Tea event loop at steady 60 FPS without flickering or CPU spinning.
  - Responds immediately to terminal resize events (`tea.WindowSizeMsg`) reflowing all Lipgloss pane boundaries cleanly.

---

### Improvement 342: Real-Time Terminal Warmup Curve & Velocity Visualizer (`gitmap warmup curve-tui`)
- **ID:** `WARMUP-TUI-342`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Reviewing progressive warm-up ramping curves (linear, quadratic, sigmoid, logarithmic) currently requires launching web dashboards, preventing terminal-bound engineers and remote SSH operators from monitoring day-by-day sending schedules.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Mirroring Web SVG chart with identical data points for dual-surface operational parity.
  - **Lipgloss Terminal UX:**
    - Braille / ASCII sparkline graph rendered in terminal showing day 1 to day 30 target volume curve vs actual delivered volume.
    - Active day indicator with glowing neon cursor (`● Day 14: Target 180 / Sent 178`).
    - Color-coded velocity trend line: Cyan for planned curve, Green for delivered, Red dots for bounced volumes.
- **CLI Counterpart Command:**
  - `gitmap warmup curve-tui --campaign "primary-saas-warmup"`
  - `gitmap warmup curve-tui --model "sigmoid" --days 30`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["campaignId", "modelType", "currentDay", "targetVolume", "actualVolume", "isRampingNominal", "hasDeviatedFromCurve"],
    "properties": {
      "campaignId": { "type": "string" },
      "modelType": { "type": "string", "enum": ["linear", "quadratic", "sigmoid", "exponential"] },
      "currentDay": { "type": "integer", "minimum": 1 },
      "targetVolume": { "type": "integer", "minimum": 0 },
      "actualVolume": { "type": "integer", "minimum": 0 },
      "isRampingNominal": { "type": "boolean" },
      "hasDeviatedFromCurve": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Accurately plots 30-day projection curves using terminal character grid rendering.
  - Dynamically recalculates volume targets if deliverability score dips below safe threshold.

---

### Improvement 343: Interactive SMTP Dialogue Debugger & Raw Protocol Hex Inspector
- **ID:** `WARMUP-TUI-343`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Diagnosing rejected SMTP handshakes requires wading through dense raw log files without protocol syntax highlighting, missing handshake timing breakdowns and TLS certificate inspection.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Web counterpart viewer with collapsible command/response pairs and hex dump side pane.
  - **Lipgloss Terminal UX:**
    - Colored protocol trace viewer:
      - Client commands rendered in Cyan (`> EHLO mail.sender.com`).
      - Server responses rendered in Emerald (`< 250-STARTTLS`).
      - TLS negotiation handshake block highlighted in Purple (`[TLS 1.3 / X25519 / AES_256_GCM]`).
      - Hex & ASCII side-by-side dump view for investigating byte-level encoding bugs or CRLF injection.
- **CLI Counterpart Command:**
  - `gitmap mail debug-dialogue --host "gmail-smtp-in.l.google.com" --port 25`
  - `gitmap mail debug-dialogue --session-id "sess_98234ab8"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["sessionId", "remoteHost", "totalDurationMs", "dialogueSteps", "isTlsEstablished", "hasProtocolViolation"],
    "properties": {
      "sessionId": { "type": "string" },
      "remoteHost": { "type": "string" },
      "totalDurationMs": { "type": "integer", "minimum": 0 },
      "dialogueSteps": {
        "type": "array",
        "items": {
          "type": "object",
          "required": ["direction", "command", "responseCode", "payload"],
          "properties": {
            "direction": { "type": "string", "enum": ["client", "server"] },
            "command": { "type": "string" },
            "responseCode": { "type": "integer" },
            "payload": { "type": "string" }
          }
        }
      },
      "isTlsEstablished": { "type": "boolean" },
      "hasProtocolViolation": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Captures and displays every step of the SMTP protocol handshake with millisecond latency resolution.
  - Highlights protocol violations (missing CRLF line endings, oversized lines > 998 octets, invalid character bytes).

---

### Improvement 344: Live Spool Memory & Disk Queue TUI Monitor (`gitmap mail queue-top`)
- **ID:** `WARMUP-TUI-344`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Operators have no interactive `top`-like command to monitor live outbound spool queues, active delivery workers, message retry intervals, and memory consumption in real time.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Web dashboard view with sortable tables by queue latency, size, and destination provider.
  - **Lipgloss Terminal UX:**
    - `htop`-inspired terminal dashboard:
      - Header metrics: Active Spool Count, Processing Rate (msg/s), Disk Spool Size (MB), Failed / Retry Count.
      - Sortable interactive message list: ID, Recipient, MX Domain, State (`PENDING`, `SENDING`, `BACKOFF`), Age, Attempt #.
      - Action hotkeys: `[k] Kill Message | [r] Retry Now | [f] Freeze Queue | [p] Purge Dead Letters`.
- **CLI Counterpart Command:**
  - `gitmap mail queue-top`
  - `gitmap mail queue-top --refresh-interval 500ms`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["totalQueued", "activeWorkers", "messagesPerSecond", "spoolSizeBytes", "isQueueFrozen", "hasDeadLetters"],
    "properties": {
      "totalQueued": { "type": "integer", "minimum": 0 },
      "activeWorkers": { "type": "integer", "minimum": 0 },
      "messagesPerSecond": { "type": "number", "minimum": 0.0 },
      "spoolSizeBytes": { "type": "integer", "minimum": 0 },
      "isQueueFrozen": { "type": "boolean" },
      "hasDeadLetters": { "type": "boolean" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Refreshes spool queue telemetry at user-defined interval (default: 1000ms) with zero UI flicker.
  - Supports keyboard-driven sorting by Age, Recipient, and Attempt count.

---

### Improvement 345: Interactive DNS & Authentication Wizard TUI (`gitmap mail auth-wizard`)
- **ID:** `WARMUP-TUI-345`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Onboarding a new sending domain requires manually constructing SPF records, generating 2048-bit DKIM keys, crafting DMARC policies, and verifying rDNS across external tools, leading to frequent typos and configuration errors.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Step-by-step guided stepper in Web UI reflecting identical domain verification inputs.
  - **Lipgloss Terminal UX:**
    - Bubble Tea multi-step wizard:
      - Step 1: Input domain name and selector with input field validation.
      - Step 2: Auto-generate 2048-bit RSA and Ed25519 DKIM keys with one-click copy to clipboard.
      - Step 3: Interactive DMARC policy builder with radio select for `none`, `quarantine`, or `reject`.
      - Step 4: Real-time DNS polling status checker with live spinner waiting for TXT record propagation.
- **CLI Counterpart Command:**
  - `gitmap mail auth-wizard --domain "outbound.company.com"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["domain", "selector", "isDkimConfigured", "isSpfConfigured", "isDmarcConfigured", "hasDnsPropagated", "generatedTxtRecords"],
    "properties": {
      "domain": { "type": "string" },
      "selector": { "type": "string" },
      "isDkimConfigured": { "type": "boolean" },
      "isSpfConfigured": { "type": "boolean" },
      "isDmarcConfigured": { "type": "boolean" },
      "hasDnsPropagated": { "type": "boolean" },
      "generatedTxtRecords": {
        "type": "array",
        "items": { "type": "string" }
      }
    }
  }
  ```
- **Acceptance Criteria:**
  - Generates cryptographically secure 2048-bit RSA keypair in <100ms.
  - Polls authoritative DNS nameservers directly with timeout and displays live confirmation badges once TXT records resolve.

---

### Improvement 346: Bubble Tea Responsive Email Template Previewer & Dark Mode Tester
- **ID:** `WARMUP-TUI-346`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Previewing HTML and plaintext warm-up email templates across desktop and mobile screen viewports requires building separate preview servers or sending test emails to external devices.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Web viewport switcher (Desktop 1200px, Tablet 768px, Mobile 375px) with dark/light mode toggle.
  - **Lipgloss Terminal UX:**
    - Terminal-based dual viewport renderer:
      - Plaintext Markdown / ANSI conversion viewer on the left.
      - HTML structural tree and image-alt tags inspector on the right.
      - Viewport toggle (`[m] Mobile 40-col | [d] Desktop 80-col | [t] Dark/Light Theme Simulator`).
- **CLI Counterpart Command:**
  - `gitmap mail preview --template "./templates/warmup-intro.html"`
  - `gitmap mail preview --template "./templates/warmup-intro.html" --viewport "mobile"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["templatePath", "viewportWidth", "isDarkModeSimulated", "hasPlaintextAlternative", "htmlSizeBytes", "spamScorePrediction"],
    "properties": {
      "templatePath": { "type": "string" },
      "viewportWidth": { "type": "integer", "enum": [40, 80, 120] },
      "isDarkModeSimulated": { "type": "boolean" },
      "hasPlaintextAlternative": { "type": "boolean" },
      "htmlSizeBytes": { "type": "integer", "minimum": 0 },
      "spamScorePrediction": { "type": "number", "minimum": 0.0, "maximum": 10.0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Renders HTML and text alternatives directly in terminal without third-party external dependencies.
  - Flags potential spam triggers (oversized HTML payload > 102KB leading to Gmail clipping, missing alt tags, bad link ratios).

---

### Improvement 347: Interactive Peer Cluster Mesh & Remote Node Monitor (`gitmap mail nodes-tui`)
- **ID:** `WARMUP-TUI-347`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** In distributed multi-node Lemon service deployments, administrators struggle to inspect cross-node pool balancing, SSH peer connectivity, and local reputation telemetry from a single synchronized terminal view.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Topology map showing interconnected master and worker nodes with latency links and load percentages.
  - **Lipgloss Terminal UX:**
    - Lipgloss interactive table of cluster nodes:
      - Node Name, IP, Region, SSH Status, Assigned Campaigns, Outbound Rate (msg/s), Health Status (`ONLINE`, `DEGRADED`, `OFFLINE`).
      - Selecting a node opens detailed bottom drawer with CPU, Memory, SQLite WAL size, and active socket counts.
- **CLI Counterpart Command:**
  - `gitmap mail nodes-tui`
  - `gitmap mail nodes ping --all`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["nodeId", "nodeFqdn", "region", "isOnline", "hasSshReachable", "activeCampaignsCount", "latencyMs"],
    "properties": {
      "nodeId": { "type": "string" },
      "nodeFqdn": { "type": "string" },
      "region": { "type": "string" },
      "isOnline": { "type": "boolean" },
      "hasSshReachable": { "type": "boolean" },
      "activeCampaignsCount": { "type": "integer", "minimum": 0 },
      "latencyMs": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Pings and queries all peer nodes in parallel via Go goroutines with 1.5s timeout.
  - Renders responsive interactive Bubble Tea table supporting Arrow keys, Enter for details, and Esc to return.

---

### Improvement 348: Real-Time Bounded Log Tailer with Lipgloss Syntax Highlighting
- **ID:** `WARMUP-TUI-348`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Standard `tail -f` outputs unformatted raw JSON or mixed log output that overwhelms operators with thousands of routine debug lines during active email sending runs.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Web log viewer with level filtering (DEBUG, INFO, WARN, ERROR) and full-text search.
  - **Lipgloss Terminal UX:**
    - High-performance streaming log viewer in Bubble Tea:
      - Color-coded severity chips (`INFO` Blue, `WARN` Yellow, `ERROR` Red, `FATAL` Magenta).
      - Interactive fuzzy search (`/`) filtering entries dynamically in real time without pausing the stream.
      - Pause/Resume stream toggle (`[space]`) with buffer scrollback memory.
- **CLI Counterpart Command:**
  - `gitmap mail logs tail --level "warn" --follow`
  - `gitmap mail logs tail --filter "550 5.1.1"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["timestamp", "logLevel", "subsystem", "message", "isErrorLevel", "hasStackTrace", "contextData"],
    "properties": {
      "timestamp": { "type": "string", "format": "date-time" },
      "logLevel": { "type": "string", "enum": ["debug", "info", "warn", "error", "fatal"] },
      "subsystem": { "type": "string" },
      "message": { "type": "string" },
      "isErrorLevel": { "type": "boolean" },
      "hasStackTrace": { "type": "boolean" },
      "contextData": { "type": "object" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Handles continuous bursts of >5,000 log lines per second with zero UI lockup using ring buffer memory capping.
  - Renders syntax-highlighted key-value metadata tags using Lipgloss styling.

---

### Improvement 349: Interactive Terminal Mail Composer & One-Click Test Sender (`gitmap mail send-tui`)
- **ID:** `WARMUP-TUI-349`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Testing delivery across various email relays or specific IP pools requires writing temporary scripts or executing complex raw CLI commands with multiple parameters.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Quick test email composer modal with preset test accounts (Gmail, Yahoo, Outlook, Mailosaur).
  - **Lipgloss Terminal UX:**
    - Bubble Tea form interface with focused text inputs:
      - From Address, To Address, Subject, Body (multi-line textarea), Relay Pool Selector (dropdown).
      - Live DKIM selector and SPF preview check indicators before dispatching.
      - "Send Envelope" button triggering immediate delivery with real-time response modal.
- **CLI Counterpart Command:**
  - `gitmap mail send-tui`
  - `gitmap mail send-test --to "test@gmail.com" --pool "primary"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["fromAddress", "toAddress", "subject", "poolId", "isDispatched", "hasDeliveredSuccessfully", "smtpResponseBanner"],
    "properties": {
      "fromAddress": { "type": "string" },
      "toAddress": { "type": "string" },
      "subject": { "type": "string" },
      "poolId": { "type": "string" },
      "isDispatched": { "type": "boolean" },
      "hasDeliveredSuccessfully": { "type": "boolean" },
      "smtpResponseBanner": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Provides frictionless keyboard navigation across input fields using Tab and Shift-Tab.
  - Displays instant delivery confirmation with actual remote SMTP response code and latency within 2 seconds.

---

### Improvement 350: Configurable Lipgloss Theme & Color Palette Studio (`gitmap mail theme`)
- **ID:** `WARMUP-TUI-350`
- **Cluster:** 35 (High-Velocity Bubble Tea / Lipgloss Terminal Studios)
- **Problem Statement:** Operators working in various terminal environments (standard ANSI 16-color, 256-color, or 24-bit TrueColor) frequently encounter illegible contrast ratios or unreadable text when fixed colors are hardcoded into TUI components.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Web theme sync selector allowing user to mirror terminal palette (Catppuccin Mocha, Nord, Dracula, Monokai, High Contrast).
  - **Lipgloss Terminal UX:**
    - Theme switcher interactive previewer:
      - Cycles through pre-built palettes in real-time.
      - Tests WCAG contrast compliance for every status chip against current terminal background color.
      - Saves preferred theme to `~/.gitmap/tui-theme.json`.
- **CLI Counterpart Command:**
  - `gitmap mail theme list`
  - `gitmap mail theme set "catppuccin-mocha"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["themeName", "isTrueColorSupported", "hasHighContrastMode", "primaryColor", "successColor", "warningColor", "errorColor"],
    "properties": {
      "themeName": { "type": "string" },
      "isTrueColorSupported": { "type": "boolean" },
      "hasHighContrastMode": { "type": "boolean" },
      "primaryColor": { "type": "string" },
      "successColor": { "type": "string" },
      "warningColor": { "type": "string" },
      "errorColor": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically detects terminal color profile (`termenv.ColorProfile()`) and gracefully degrades from TrueColor to ANSI 256 or ANSI 16 without throwing errors.
  - Persists theme preferences immediately and applies universally across all `gitmap mail` and `gitmap warmup` TUI commands.

---

# Cluster 36: Resilient Local Spooling & Zero-Loss WAL Delivery (Items 351–360)

---

### Improvement 351: SQLite WAL-Backed Zero-Loss Outbound Spool Engine
- **ID:** `WARMUP-SPOOL-351`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Outbound warm-up engines that store queued envelopes purely in system RAM lose unsent batches upon system crashes, abrupt machine reboots, or process termination, causing non-deterministic delivery states and lost customer emails.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Spool storage health widget displaying SQLite database size, WAL log file size, uncheckpointed pages, and write throughput.
    - Manual "Checkpoint WAL" button triggering synchronous truncate checkpoint (`PRAGMA wal_checkpoint(TRUNCATE)`).
    - Queue integrity audit indicator showing zero corruption and 100% durable page sync.
  - **Lipgloss Terminal UX:**
    - Terminal gauge showing WAL status: `[WAL Size: 2.4 MB] [Checkpointed: YES] [Durable Envelopes: 1,420]`.
- **CLI Counterpart Command:**
  - `gitmap mail spool status`
  - `gitmap mail spool checkpoint --truncate`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["spoolPath", "totalMessagesDurable", "walSizeBytes", "isWalModeActive", "hasPendingCheckpoint", "uncommittedTransactions"],
    "properties": {
      "spoolPath": { "type": "string" },
      "totalMessagesDurable": { "type": "integer", "minimum": 0 },
      "walSizeBytes": { "type": "integer", "minimum": 0 },
      "isWalModeActive": { "type": "boolean" },
      "hasPendingCheckpoint": { "type": "boolean" },
      "uncommittedTransactions": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Every outbound envelope is committed to SQLite with `PRAGMA synchronous = NORMAL` and WAL journal mode before returning success to caller.
  - Guarantees zero message loss across unexpected SIGKILL or power outage recovery tests.

---

### Improvement 352: Idempotent Delivery Token & De-Duplication Fingerprint Ledger
- **ID:** `WARMUP-SPOOL-352`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Network blips during SMTP `DATA` acknowledgment can cause worker retry loops to re-transmit identical envelopes, causing embarrassing duplicate emails to recipients and triggering spam filters.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Deduplication collision inspector showing prevented duplicate send attempts, recipient address, and fingerprint hash.
    - Idempotency TTL configuration input (default: 24 hours).
    - Ledger search tool to look up delivery receipt tokens and original transmission timestamps.
  - **Lipgloss Terminal UX:**
    - Cyan status badge on deduplication event: `[DEDUP PREVENTED] Duplicate hash 8f14a9... blocked for recipient: user@domain.com`.
- **CLI Counterpart Command:**
  - `gitmap mail dedup check --hash "8f14a9c2e0b7"`
  - `gitmap mail dedup purge --older-than "48h"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["messageFingerprint", "recipientEmail", "firstSeenAt", "isDuplicateBlocked", "hasDeliveredOnce", "tokenTtlSeconds"],
    "properties": {
      "messageFingerprint": { "type": "string" },
      "recipientEmail": { "type": "string" },
      "firstSeenAt": { "type": "string", "format": "date-time" },
      "isDuplicateBlocked": { "type": "boolean" },
      "hasDeliveredOnce": { "type": "boolean" },
      "tokenTtlSeconds": { "type": "integer", "minimum": 60 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Computes deterministic SHA-256 fingerprint from `Message-ID`, canonicalized recipient email, and body digest.
  - Drops duplicate attempts within TTL window in <1ms without opening outbound network socket.

---

### Improvement 353: Segmented Priority Queuing with Fair Starvation Prevention
- **ID:** `WARMUP-SPOOL-353`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** In shared spool databases, massive batches of background warm-up emails clog worker threads, starving critical transactional emails (password resets, OTP verification codes) for minutes or hours.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Dual-lane spool throughput visualizer: High Priority Lane (Transactional) vs Standard Lane (Warm-up Ramping).
    - Worker reservation ratio slider (e.g., 70% workers reserved for High Priority, 30% for Warm-up).
    - Starvation alarm indicator showing maximum wait time of background items to prevent complete stall.
  - **Lipgloss Terminal UX:**
    - Split progress bar: Top lane shows VIP Transactional queue depth (Emerald), bottom lane shows Warm-up queue depth (Purple).
- **CLI Counterpart Command:**
  - `gitmap mail queue priority-stats`
  - `gitmap mail queue set-allocation --priority-ratio 0.70`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["priorityTier", "queuedCount", "processedCount", "isStarvationPrevented", "hasCapacityReserved", "avgLatencyMs"],
    "properties": {
      "priorityTier": { "type": "string", "enum": ["transactional_critical", "warmup_normal", "bulk_low"] },
      "queuedCount": { "type": "integer", "minimum": 0 },
      "processedCount": { "type": "integer", "minimum": 0 },
      "isStarvationPrevented": { "type": "boolean" },
      "hasCapacityReserved": { "type": "boolean" },
      "avgLatencyMs": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Transactional messages achieve <500ms dispatch latency regardless of background warm-up queue volume.
  - Implements fair-share deficit round-robin so low-priority items still make progress without indefinite starvation.

---

### Improvement 354: Exponential Jitter Backoff & Greylist Retry Scheduler
- **ID:** `WARMUP-SPOOL-354`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Fixed retry intervals (e.g., retry every 5 minutes) cause thundering herd problems when destination servers recover, causing immediate re-choking and continued delivery rejections.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Retry timeline curve showing computed exponential backoff with randomized jitter spreads.
    - Configuration sliders for Base Interval (default: 30s), Max Interval (default: 4 hours), and Max Attempts (default: 8).
    - Next-scheduled retry bucket histogram grouping pending retries across upcoming time windows.
  - **Lipgloss Terminal UX:**
    - Visual timeline of retry attempts rendered with Lipgloss: `Attempt 1: +45s -> Attempt 2: +180s -> Attempt 3: +720s`.
- **CLI Counterpart Command:**
  - `gitmap mail retry config --base-seconds 30 --max-attempts 8 --jitter-factor 0.25`
  - `gitmap mail retry inspect --message-id "msg_48723910"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["messageId", "attemptCount", "maxAttempts", "nextRetryTimestamp", "isJitterApplied", "hasReachedMaxRetries", "backoffSeconds"],
    "properties": {
      "messageId": { "type": "string" },
      "attemptCount": { "type": "integer", "minimum": 0 },
      "maxAttempts": { "type": "integer", "minimum": 1 },
      "nextRetryTimestamp": { "type": "string", "format": "date-time" },
      "isJitterApplied": { "type": "boolean" },
      "hasReachedMaxRetries": { "type": "boolean" },
      "backoffSeconds": { "type": "integer", "minimum": 1 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Applies Full Jitter algorithm: `Sleep = rand(0, min(MaxInterval, BaseInterval * 2^attempt))`.
  - Distributes retry attempts evenly across time to eliminate synchronized thundering herds.

---

### Improvement 355: Dead Letter Queue (DLQ) Quarantine & One-Click Inspection Vault
- **ID:** `WARMUP-SPOOL-355`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Envelopes that exhaust maximum retry limits or fail permanently with 5xx fatal errors are either silently deleted or dumped into obscure log files where debugging is impossible.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Dead Letter Queue repository view with paginated table: Message ID, Recipient, Terminal Error Reason, Creation Date, Failure Count.
    - Deep message inspector modal revealing raw headers, MIME structure, and exact remote SMTP failure transcript.
    - Batch actions: "Re-enqueue to Alternative Pool", "Export to EML File", "Purge Permanently".
  - **Lipgloss Terminal UX:**
    - Red bordered DLQ list card with keyboard navigation: `[Enter] View Full Transcript | [r] Re-queue | [d] Delete`.
- **CLI Counterpart Command:**
  - `gitmap mail dlq list`
  - `gitmap mail dlq requeue --id "dlq_839120"`
  - `gitmap mail dlq export --id "dlq_839120" --out-file "./failed-msg.eml"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["dlqId", "originalMessageId", "recipient", "fatalErrorCode", "isQuarantined", "hasRequeuedSuccessfully", "quarantinedAt"],
    "properties": {
      "dlqId": { "type": "string" },
      "originalMessageId": { "type": "string" },
      "recipient": { "type": "string" },
      "fatalErrorCode": { "type": "string" },
      "isQuarantined": { "type": "boolean" },
      "hasRequeuedSuccessfully": { "type": "boolean" },
      "quarantinedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Permanently moves exhausted messages to `dead_letter_spool` table with full diagnostic metadata intact.
  - Re-enqueue action creates fresh message envelope with reset retry counters and assigns selected fallback pool.

---

### Improvement 356: Local Spool Cryptographic Encryption at Rest (AES-256-GCM)
- **ID:** `WARMUP-SPOOL-356`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Sensitive emails containing user verification links or customer communications stored unencrypted in local SQLite databases are vulnerable to exposure if laptop disks or cloud instances are compromised.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Vault security status indicator: Encryption Algorithm (`AES-256-GCM`), Key Derivation (`Argon2id`), Key ID, and Re-keying status.
    - Master key unlock / status verification badge.
    - Audit log recording all decryption events and key access operations.
  - **Lipgloss Terminal UX:**
    - Green shield icon in header: `[ENCRYPTED AT REST: AES-256-GCM]`.
- **CLI Counterpart Command:**
  - `gitmap mail vault status`
  - `gitmap mail vault init --key-file "./vault.key"`
  - `gitmap mail vault rotate-key --new-key-file "./vault-new.key"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["keyId", "algorithm", "isVaultUnlocked", "hasValidMasterKey", "encryptedEnvelopesCount", "lastRotatedAt"],
    "properties": {
      "keyId": { "type": "string" },
      "algorithm": { "type": "string", "enum": ["AES-256-GCM", "ChaCha20-Poly1305"] },
      "isVaultUnlocked": { "type": "boolean" },
      "hasValidMasterKey": { "type": "boolean" },
      "encryptedEnvelopesCount": { "type": "integer", "minimum": 0 },
      "lastRotatedAt": { "type": "string", "format": "date-time" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Envelopes stored in SQLite `payload` column are encrypted with unique 96-bit nonce and verified with 128-bit authentication tag.
  - Unlocked purely in-memory during transmission; plaintext never touches persistent disk storage.

---

### Improvement 357: Automatic SQLite Spool Vacuuming & Dynamic WAL Auto-Checkpointing
- **ID:** `WARMUP-SPOOL-357`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** High-volume warm-up dispatching causes continuous database insertions and deletions, resulting in bloated SQLite WAL files exceeding gigabytes of disk space and fragmented page caches that degrade I/O performance.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Disk storage optimization card: Database Size, WAL Size, Free Space Reclamation Opportunity, and Fragmentation %.
    - Automated threshold settings: Trigger auto-checkpoint at 1,000 pages (~4MB) and incremental vacuum when fragmentation exceeds 20%.
    - Live I/O write graph showing storage reclaimed over time.
  - **Lipgloss Terminal UX:**
    - Terminal progress spinner during optimization: `Optimizing spool storage... Reclaimed 142 MB [DONE]`.
- **CLI Counterpart Command:**
  - `gitmap mail spool optimize`
  - `gitmap mail spool vacuum --threshold-mb 50`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["databaseSizeBytes", "reclaimedBytes", "isVacuumCompleted", "hasWalTruncated", "fragmentationPercent", "durationMs"],
    "properties": {
      "databaseSizeBytes": { "type": "integer", "minimum": 0 },
      "reclaimedBytes": { "type": "integer", "minimum": 0 },
      "isVacuumCompleted": { "type": "boolean" },
      "hasWalTruncated": { "type": "boolean" },
      "fragmentationPercent": { "type": "number", "minimum": 0.0, "maximum": 100.0 },
      "durationMs": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Automatically executes incremental vacuum in background without locking active reader or writer transactions.
  - Keeps WAL log file size under 10MB during continuous sustained sending workloads.

---

### Improvement 358: EML File-System Drop-Spooler with Inotify / ReadDirectoryChanges Watcher
- **ID:** `WARMUP-SPOOL-358`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Legacy applications or offline background scripts cannot easily invoke Go APIs or CLI binaries to send emails, needing a simple file-drop directory where dropping `.eml` files immediately triggers transmission.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Drop-spool folder manager: Active Watch Path, Ingested Files Count, Processing Status, Error Folder.
    - Drag-and-drop web dropzone to ingest `.eml` files directly into active outbound queue.
  - **Lipgloss Terminal UX:**
    - Status line: `[DROP-SPOOL WATCHER ACTIVE] path: /var/spool/gitmap/outgoing (watching with fsnotify)`.
- **CLI Counterpart Command:**
  - `gitmap mail spool watch --dir "./spool/outgoing"`
  - `gitmap mail spool watch --dir "./spool/outgoing" --daemon`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["watchDir", "filesIngested", "isWatcherActive", "hasIngestError", "activeFileCount", "lastIngestedFile"],
    "properties": {
      "watchDir": { "type": "string" },
      "filesIngested": { "type": "integer", "minimum": 0 },
      "isWatcherActive": { "type": "boolean" },
      "hasIngestError": { "type": "boolean" },
      "activeFileCount": { "type": "integer", "minimum": 0 },
      "lastIngestedFile": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Uses native OS file notification events (`inotify` on Linux, `kqueue` on macOS, `ReadDirectoryChangesW` on Windows) via `fsnotify`.
  - Atomically moves parsed files from `outgoing/` to `processing/` to `completed/` (or `failed/`), guaranteeing exactly-once ingestion.

---

### Improvement 359: Outbound Transmission Receipt Journal with Merkle Proof Verification
- **ID:** `WARMUP-SPOOL-359`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** Proving that an email was strictly accepted by remote MX servers (non-repudiation) is impossible without cryptographically verified delivery receipts containing remote server banners, timestamps, and message digests.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Cryptographic receipt viewer displaying: Message Hash, Delivery Timestamp, Remote Server Host & IP, Remote 250 OK Queue ID, and Merkle Proof Root.
    - "Verify Receipt Authenticity" validator widget.
    - PDF / JSON cryptographic receipt export button.
  - **Lipgloss Terminal UX:**
    - Clean terminal receipt printout with verified green seal: `[RECEIPT VERIFIED: MERKLE ROOT 4a82c1...]`.
- **CLI Counterpart Command:**
  - `gitmap mail receipt get --message-id "msg_123456"`
  - `gitmap mail receipt verify --proof-file "./receipt.json"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["messageId", "recipient", "remoteMxBanner", "smtpQueueId", "isReceiptVerified", "hasValidMerkleProof", "merkleRootHash"],
    "properties": {
      "messageId": { "type": "string" },
      "recipient": { "type": "string" },
      "remoteMxBanner": { "type": "string" },
      "smtpQueueId": { "type": "string" },
      "isReceiptVerified": { "type": "boolean" },
      "hasValidMerkleProof": { "type": "boolean" },
      "merkleRootHash": { "type": "string" }
    }
  }
  ```
- **Acceptance Criteria:**
  - Generates append-only cryptographic journal entry containing raw 250 response banner, timestamp, and message body hash.
  - Verifies proofs mathematically against daily root Merkle tree hash in <5ms.

---

### Improvement 360: Offline Resilient Local Relay Daemon with P2P Mesh Spool Sync
- **ID:** `WARMUP-SPOOL-360`
- **Cluster:** 36 (Resilient Local Spooling & Zero-Loss WAL Delivery)
- **Problem Statement:** When laptops or edge worker nodes lose internet connectivity, outbound email generation fails completely unless a local lightweight SMTP daemon accepts envelopes on `127.0.0.1:2525` and syncs them automatically across peer nodes when connectivity recovers.
- **Concrete UI Specification:**
  - **Zero-CDN Web UI:**
    - Local daemon status tile: Listening Port (`2525`), Local Queue Count, Mesh Peer Sync Status, Connectivity State (Online / Offline Mesh Mode).
    - Peer discovery list showing local LAN/Tailscale nodes available for mesh delivery relay.
    - "Sync Spool to Peer Node" manual trigger button.
  - **Lipgloss Terminal UX:**
    - Amber pulsing indicator when offline: `[OFFLINE SPOOL MODE: 12 Envelopes Buffered] -> Will auto-drain via mesh`.
- **CLI Counterpart Command:**
  - `gitmap mail daemon start --port 2525 --mesh-sync`
  - `gitmap mail daemon status`
  - `gitmap mail daemon flush --peer "node-02.fleet.internal"`
- **Data Contract:**
  ```json
  {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "required": ["listenAddress", "bufferedCount", "isDaemonRunning", "isMeshSyncActive", "hasInternetConnectivity", "peerCount"],
    "properties": {
      "listenAddress": { "type": "string" },
      "bufferedCount": { "type": "integer", "minimum": 0 },
      "isDaemonRunning": { "type": "boolean" },
      "isMeshSyncActive": { "type": "boolean" },
      "hasInternetConnectivity": { "type": "boolean" },
      "peerCount": { "type": "integer", "minimum": 0 }
    }
  }
  ```
- **Acceptance Criteria:**
  - Runs local lightweight SMTP listener on localhost; accepts standard RFC 5321 envelopes and queues them into SQLite WAL spool without network dependency.
  - Detects network restoration or peer node reachability and automatically flushes buffered queue without operator intervention.
