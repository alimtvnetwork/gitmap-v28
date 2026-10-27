# Improvements 601–630: Cryptographic ARC, Reputation Circuits & MIME Trees

This document specifies Improvements 601 through 630 for Sequence 343, spanning Clusters 61, 62, and 63. All specifications strictly follow positive boolean naming (`is`/`has` only), relative repository paths, Split-DB SQLite storage (`repodb/warmup.db` and `repodb/audit.db`), and zero-external-CDN UI architectures.

---

# Cluster 61: Cross-ISP DMARC / ARC Cryptographic Verification Engine & Dynamic Header Studio (601–610)

---

### Improvement 601: Recursive ARC-Seal (AS) Cryptographic Validation Chain Inspector
- **ID**: `IMP-601`
- **Cluster**: 61 (DMARC / ARC Cryptographic Verification)
- **Goal**: Traverses and cryptographically evaluates sequential ARC-Seal sets ($i=1\dots N$) across intermediate MTAs to detect signature tampering and expired RSA/Ed25519 selector keys.
- **UI/UX Terminal & Web Specification**:
  - *Terminal (Lipgloss)*: Rendered via `lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#5A67D8"))`. Shows an ASCII chain diagram with color-coded nodes: `#48BB78` (valid seal), `#E53E3E` (invalid hash/expired selector), `#ED8936` (missing intermediary).
  - *Web (Shadow DOM)*: `<arc-seal-chain>` custom component with offline SVG path connections, interactive popovers for DNS TXT responses, and zero external styling dependencies.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup arc verify-chain --message-id=<msg-uuid> [--has-strict-dns] [--is-json]`
  - Flag `--has-strict-dns`: Forces direct root nameserver queries bypassing local OS resolver caches.
  - Flag `--is-json`: Outputs machine-readable chain status to stdout.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ArcSealChainNodes (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      instanceIndex INTEGER NOT NULL,
      selectorName TEXT NOT NULL,
      signingDomain TEXT NOT NULL,
      signatureAlgorithm TEXT NOT NULL, -- 'rsa-sha256', 'ed25519-sha256'
      isSealValid INTEGER NOT NULL DEFAULT 0,
      hasKeyExpired INTEGER NOT NULL DEFAULT 0,
      rawSealBlock TEXT NOT NULL,
      verifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ArcSealNodeResult struct {
      InstanceIndex      int    `json:"instanceIndex"`
      SelectorName       string `json:"selectorName"`
      SigningDomain      string `json:"signingDomain"`
      SignatureAlgorithm string `json:"signatureAlgorithm"`
      IsSealValid        bool   `json:"isSealValid"`
      HasKeyExpired      bool   `json:"hasKeyExpired"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrates with GitMap's mail delivery diagnostics pipeline via `pkg/warmup/arc/validator.go` and triggers alert events in `gitmap pe -t` telemetry feeds.

---

### Improvement 602: ARC-Message-Signature (AMS) Body Canonicalization Hash Verifier
- **ID**: `IMP-602`
- **Cluster**: 61
- **Goal**: Validates `bh=` (body hash) parameters against relaxed/simple canonicalization streams to pinpoint whitespace alterations or injected HTML disclaimers by intermediate relays.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual side-by-side diff showing exact byte offsets where canonicalization broke, styled with Lipgloss background highlights (`#742A2A` for mutated bytes, `#22543D` for matching bytes).
  - *Web*: `<ams-hash-diff>` visualizer with dual pane synchronized scrolling and inline hex inspector.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup arc check-ams --message-id=<msg-uuid> [--has-relaxed-body] [--is-verbose]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AmsBodyHashResults (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      instanceIndex INTEGER NOT NULL,
      canonicalizationHeader TEXT NOT NULL, -- 'relaxed', 'simple'
      canonicalizationBody TEXT NOT NULL,
      expectedBodyHash TEXT NOT NULL,
      computedBodyHash TEXT NOT NULL,
      isHashMatched INTEGER NOT NULL DEFAULT 0,
      hasTrailingWhitespaceMutation INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AmsBodyHashRecord struct {
      MessageID                     string `json:"messageId"`
      InstanceIndex                 int    `json:"instanceIndex"`
      CanonicalizationBody          string `json:"canonicalizationBody"`
      IsHashMatched                 bool   `json:"isHashMatched"`
      HasTrailingWhitespaceMutation bool   `json:"hasTrailingWhitespaceMutation"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Linked directly into `gitmap email warmup test-deliverability` to flag transport-level body alterations.

---

### Improvement 603: ARC-Authentication-Results (AAR) Header Synthesizer & Policy Inspector
- **ID**: `IMP-603`
- **Cluster**: 61
- **Goal**: Parses, normalizes, and validates intermediate `ARC-Authentication-Results` headers across heterogeneous MTA formats (Postfix, Exim, Sendmail, Google Workspace, Exchange).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Formatted tabular summary using `lipgloss.NewStyle().Foreground(lipgloss.Color("#ECC94B"))` listing each hop's SPF, DKIM, and DMARC evaluated verdict.
  - *Web*: `<aar-policy-matrix>` element rendering pass/neutral/fail tags with zero external CSS frameworks.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup arc inspect-aar --input-eml=<path> [--is-summary] [--has-dkim-detail]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AarPolicyEvaluations (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      authservId TEXT NOT NULL,
      instanceIndex INTEGER NOT NULL,
      spfVerdict TEXT NOT NULL, -- 'pass', 'neutral', 'softfail', 'fail'
      dkimVerdict TEXT NOT NULL,
      dmarcVerdict TEXT NOT NULL,
      isAuthentic INTEGER NOT NULL DEFAULT 0,
      hasSpfAligned INTEGER NOT NULL DEFAULT 0,
      evaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AarPolicyResult struct {
      AuthservID    string `json:"authservId"`
      SpfVerdict    string `json:"spfVerdict"`
      DkimVerdict   string `json:"dkimVerdict"`
      IsAuthentic   bool   `json:"isAuthentic"`
      HasSpfAligned bool   `json:"hasSpfAligned"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Used by `pkg/warmup/header/aar.go` to construct verifiable incoming hop graphs.

---

### Improvement 604: DMARC Alignment Matrix Visualizer (Strict vs Relaxed Modes)
- **ID**: `IMP-604`
- **Cluster**: 61
- **Goal**: Evaluates RFC 5322 From domain against RFC 5321 Return-Path (SPF alignment) and DKIM `d=` signatures (DKIM alignment) across both strict (`aspf=s`, `adkim=s`) and relaxed (`aspf=r`, `adkim=r`) modes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual 2x2 Lipgloss matrix showing alignment axes. Green `#38A169` tick marks for compliant domains; red `#E53E3E` cross marks for sub-domain misalignment in strict mode.
  - *Web*: `<dmarc-alignment-grid>` standalone web component rendering SVG vectors representing organizational root domain boundaries.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup dmarc align --from-domain=<dom> --return-path=<path> --dkim-domain=<dom> [--has-strict-mode]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DmarcAlignmentSnapshots (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      fromHeaderDomain TEXT NOT NULL,
      returnPathDomain TEXT NOT NULL,
      dkimSignDomain TEXT NOT NULL,
      isSpfAligned INTEGER NOT NULL DEFAULT 0,
      isDkimAligned INTEGER NOT NULL DEFAULT 0,
      hasStrictPolicy INTEGER NOT NULL DEFAULT 0,
      evaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DmarcAlignmentRecord struct {
      DomainName      string `json:"domainName"`
      IsSpfAligned    bool   `json:"isSpfAligned"`
      IsDkimAligned   bool   `json:"isDkimAligned"`
      HasStrictPolicy bool   `json:"hasStrictPolicy"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforces pre-flight warmup checks in `gitmap email warmup start` before triggering warm-up send queues.

---

### Improvement 605: Automated DKIM Selector Key-Pair Rotation & DNS Prober
- **ID**: `IMP-605`
- **Cluster**: 61
- **Goal**: Automates generation, rotation, and active DNS propagation polling for DKIM 2048-bit RSA and Ed25519 key-pairs, ensuring zero message drops during key rollover.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Progress spinner with step checkpoints: `[✓] Keygen -> [✓] DNS Poll (8/8 Nameservers) -> [✓] Active Selector Switch`. Styled with Lipgloss Cyan `#4FD1C5`.
  - *Web*: Real-time countdown timer element `<dkim-rotation-timeline>` tracking TTL expiry.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup dkim rotate --domain=<domain> --selector=<name> [--has-ed25519] [--is-dry-run]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DkimKeySelectors (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      selectorName TEXT NOT NULL,
      publicKeyRecord TEXT NOT NULL,
      privateKeyRef TEXT NOT NULL,
      algorithmType TEXT NOT NULL, -- 'rsa2048', 'ed25519'
      isCurrentlyActive INTEGER NOT NULL DEFAULT 0,
      hasDnsPropagated INTEGER NOT NULL DEFAULT 0,
      rotatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DkimSelectorRecord struct {
      DomainName       string `json:"domainName"`
      SelectorName     string `json:"selectorName"`
      IsCurrentlyActive bool   `json:"isCurrentlyActive"`
      HasDnsPropagated bool   `json:"hasDnsPropagated"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Secure key storage managed via local GitMap credential vault (`pkg/vault/credentials.go`).

---

### Improvement 606: Dynamic RFC 5322 Raw Header Folding & ASCII Syntax Highlighting Studio
- **ID**: `IMP-606`
- **Cluster**: 61
- **Goal**: Provides syntax-highlighted, interactive terminal navigation of RFC 5322 header blocks, displaying folding anomalies, illegal ASCII control characters, and line-length violations (>998 chars).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Full-screen Bubble Tea viewport. Headers colored by function: routing (`#63B3ED`), authentication (`#9AE6B4`), MIME structure (`#F6AD55`). Whitespace folding indicators marked in dim gray `#4A5568`.
  - *Web*: Encapsulated code block with copyable individual folded header tokens.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup header inspect --file=<eml-file> [--has-folding-lint] [--is-interactive]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HeaderSyntaxViolations (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      headerKey TEXT NOT NULL,
      lineIndex INTEGER NOT NULL,
      violationType TEXT NOT NULL, -- 'excessive_length', 'illegal_control_char', 'bad_folding'
      isConformant INTEGER NOT NULL DEFAULT 0,
      hasOverlongLine INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HeaderSyntaxViolationRecord struct {
      MessageID       string `json:"messageId"`
      HeaderKey       string `json:"headerKey"`
      IsConformant    bool   `json:"isConformant"`
      HasOverlongLine bool   `json:"hasOverlongLine"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Executed as part of the `gitmap email lint` command suite.

---

### Improvement 607: Cross-ISP DMARC Aggregate (RUA) XML Ingestion & Anomaly Detector
- **ID**: `IMP-607`
- **Cluster**: 61
- **Goal**: Ingests, uncompresses (`.gz`/`.zip`), parses, and aggregates multi-ISP DMARC RUA XML feedback files (Google, Yahoo, Fastmail, Comcast) to highlight spoofing outbreaks.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-density Lipgloss table summarizing ISP reporting orgs, total volume, SPF fail percentages, and DKIM alignment rates.
  - *Web*: `<dmarc-rua-viewer>` custom element displaying bar charts rendered purely through SVG path elements.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup dmarc ingest-rua --path=<xml-dir> [--has-auto-purge] [--is-recursive]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DmarcRuaReports (
      id TEXT PRIMARY KEY,
      orgName TEXT NOT NULL,
      reportId TEXT NOT NULL,
      dateRangeBegin DATETIME NOT NULL,
      dateRangeEnd DATETIME NOT NULL,
      totalCount INTEGER NOT NULL DEFAULT 0,
      failCount INTEGER NOT NULL DEFAULT 0,
      isSpoofingDetected INTEGER NOT NULL DEFAULT 0,
      hasPolicyOverride INTEGER NOT NULL DEFAULT 0,
      importedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DmarcRuaSummary struct {
      OrgName            string `json:"orgName"`
      ReportID           string `json:"reportId"`
      TotalCount         int    `json:"totalCount"`
      IsSpoofingDetected bool   `json:"isSpoofingDetected"`
      HasPolicyOverride  bool   `json:"hasPolicyOverride"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Scheduled background processing hook via `gitmap email warmup sync-rua`.

---

### Improvement 608: BIMI (Brand Indicators for Message Identification) SVG Tiny PS Sanitizer
- **ID**: `IMP-608`
- **Cluster**: 61
- **Goal**: Validates sender logos against RFC-compliant SVG Tiny Portable/Secure profile specifications, stripping interactive scripts, external entities, and non-conforming vector tags.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Terminal-rendered ASCII art approximation of the sanitized SVG with green/red status checks for `baseProfile="tiny-ps"`, missing root viewBox, and XML script tag absence.
  - *Web*: `<bimi-preview-box>` demonstrating mock inbox appearance in Gmail and Apple Mail headers.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup bimi validate --svg=<path> [--has-vmc-cert] [--is-strict]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BimiValidationLedger (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      svgFileHash TEXT NOT NULL,
      vmcCertificateUrl TEXT,
      isSvgTinyPsCompliant INTEGER NOT NULL DEFAULT 0,
      hasScriptTagsPurged INTEGER NOT NULL DEFAULT 0,
      hasRootViewBox INTEGER NOT NULL DEFAULT 1,
      validatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BimiValidationRecord struct {
      DomainName           string `json:"domainName"`
      IsSvgTinyPsCompliant bool   `json:"isSvgTinyPsCompliant"`
      HasScriptTagsPurged  bool   `json:"hasScriptTagsPurged"`
      HasRootViewBox       bool   `json:"hasRootViewBox"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked during brand asset setup in `pkg/warmup/bimi/sanitizer.go`.

---

### Improvement 609: SPF Macro Evaluation & DNS Lookup Limit (>10) Early Warning Guard
- **ID**: `IMP-609`
- **Cluster**: 61
- **Goal**: Recursively evaluates SPF `include:`, `a`, `mx`, `ptr`, and `redirect` mechanisms to prevent RFC 7208 PermError failures triggered by exceeding the 10 DNS lookup limit.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss tree view of DNS query hops. Shows a lookup counter `[7/10 Lookups]` in yellow; turns flashing red `#E53E3E` if count $\ge 10$.
  - *Web*: Collapsible SVG tree visualizer displaying each DNS query branch and resolution latency.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup spf inspect --domain=<domain> [--has-macro-expand] [--is-flatten-recommended]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SpfEvaluationLogs (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      dnsLookupCount INTEGER NOT NULL DEFAULT 0,
      voidLookupCount INTEGER NOT NULL DEFAULT 0,
      isLimitExceeded INTEGER NOT NULL DEFAULT 0,
      hasMacroSyntax INTEGER NOT NULL DEFAULT 0,
      rawSpfRecord TEXT NOT NULL,
      evaluatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SpfEvaluationResult struct {
      DomainName       string `json:"domainName"`
      DnsLookupCount   int    `json:"dnsLookupCount"`
      IsLimitExceeded  bool   `json:"isLimitExceeded"`
      HasMacroSyntax   bool   `json:"hasMacroSyntax"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `gitmap email warmup verify` before kicking off any outbound warmup thread.

---

### Improvement 610: MTA-STS & TLS-RPT Policy Conformance Terminal Auditor
- **ID**: `IMP-610`
- **Cluster**: 61
- **Goal**: Audits domain MTA-STS (`_mta-sts.domain.tld`) HTTPS policy files and TLSRPT (`_smtp._tls.domain.tld`) records, enforcing TLS 1.3/1.2 cipher constraints on incoming and outbound warmup relays.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Formatted Lipgloss card showing policy mode (`enforce` vs `testing`), `max_age` countdown, and MX hostname matches.
  - *Web*: `<mta-sts-inspector>` component displaying HTTPS certificate validation chain and policy fetch diagnostics.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup tls-sts audit --domain=<domain> [--has-strict-cert] [--is-enforced]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MtaStsAuditRecords (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      policyMode TEXT NOT NULL, -- 'enforce', 'testing', 'none'
      maxAgeSeconds INTEGER NOT NULL DEFAULT 604800,
      isPolicyValid INTEGER NOT NULL DEFAULT 0,
      hasTlsRptConfigured INTEGER NOT NULL DEFAULT 0,
      checkedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MtaStsAuditRecord struct {
      DomainName          string `json:"domainName"`
      PolicyMode          string `json:"policyMode"`
      IsPolicyValid       bool   `json:"isPolicyValid"`
      HasTlsRptConfigured bool   `json:"hasTlsRptConfigured"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Feeds directly into security posture evaluations reported in GitMap developer hygiene diagnostics.

---

# Cluster 62: Autonomous Reputation Feedback Loops & Micro-Bounce Dampening Circuit Breakers (611–620)

---

### Improvement 611: Multi-Tier Micro-Bounce Classification Engine (RFC 3463 Enhanced Codes)
- **ID**: `IMP-611`
- **Cluster**: 62 (Reputation & Circuit Breakers)
- **Goal**: Parses SMTP delivery responses and RFC 3463 status codes (`5.X.X` vs `4.X.X`) to distinguish transient mailbox-full or rate-limit events from permanent mailbox-does-not-exist rejections.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Real-time stream of classified bounces using Lipgloss tags: `[TRANSIENT-RATE-LIMIT]` in Orange `#DD6B20`, `[HARD-USER-UNKNOWN]` in Dark Red `#9B2C2C`.
  - *Web*: Live updating event stream component with zero external JavaScript libraries.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup bounce classify --smtp-response="<raw>" [--is-transient-check] [--has-enhanced-code]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BounceClassifications (
      id TEXT PRIMARY KEY,
      campaignId TEXT NOT NULL,
      recipientAddress TEXT NOT NULL,
      rawSmtpCode INTEGER NOT NULL,
      enhancedCode TEXT NOT NULL,
      bounceCategory TEXT NOT NULL, -- 'rate_limit', 'quota_exceeded', 'mailbox_missing', 'policy_rejection'
      isTransient INTEGER NOT NULL DEFAULT 1,
      hasReputationPenalty INTEGER NOT NULL DEFAULT 0,
      classifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BounceClassificationRecord struct {
      CampaignID           string `json:"campaignId"`
      RecipientAddress     string `json:"recipientAddress"`
      IsTransient          bool   `json:"isTransient"`
      HasReputationPenalty bool   `json:"hasReputationPenalty"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Consumed by `pkg/warmup/reputation/classifier.go` and logged to `repodb/audit.db`.

---

### Improvement 612: Probabilistic Circuit-Breaker State Machine (Closed/Half-Open/Open)
- **ID**: `IMP-612`
- **Cluster**: 62
- **Goal**: Automatically trips warmup velocity into dampening mode (`StateOpen`) when bounce rates cross $1.8\%$ within a rolling window, transitioning to `StateHalfOpen` for synthetic health probes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Prominent ANSI status badge in Bubble Tea header: `[CIRCUIT: CLOSED]` in Bright Green `#38A169`, `[CIRCUIT: OPEN - THROTTLED]` in Flashing Red `#E53E3E`.
  - *Web*: `<circuit-breaker-gauge>` displaying animated SVG state transitions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup circuit status --domain=<domain> [--has-auto-reset] [--is-forced-reset]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS WarmupCircuitBreakers (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      targetIsp TEXT NOT NULL,
      currentState TEXT NOT NULL DEFAULT 'closed', -- 'closed', 'open', 'half_open'
      rollingBounceRatio REAL NOT NULL DEFAULT 0.0,
      tripThreshold REAL NOT NULL DEFAULT 0.018,
      isOpenTripped INTEGER NOT NULL DEFAULT 0,
      hasHalfOpenProbed INTEGER NOT NULL DEFAULT 0,
      lastTransitionTime DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CircuitBreakerStateRecord struct {
      DomainName        string  `json:"domainName"`
      CurrentState      string  `json:"currentState"`
      RollingBounceRatio float64 `json:"rollingBounceRatio"`
      IsOpenTripped     bool    `json:"isOpenTripped"`
      HasHalfOpenProbed bool    `json:"hasHalfOpenProbed"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Governs rate governor concurrency in `pkg/warmup/scheduler/governor.go`.

---

### Improvement 613: Rolling Micro-Window Bounce Dampener & Backoff Governor
- **ID**: `IMP-613`
- **Cluster**: 62
- **Goal**: Calculates an exponentially weighted moving average (EWMA) of rejection rates over a sliding 100-dispatch window to execute gradual velocity deceleration before hard ISP blocks occur.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Horizontal ASCII sparkline rendered via Lipgloss block characters (` ▂▃▄▅▆▇█`) indicating rolling rejection pressure.
  - *Web*: Canvas-free SVG area chart charting volume vs bounce probability.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup governor throttle --domain=<domain> [--has-ewma] [--is-aggressive-backoff]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RollingGovernorMetrics (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      slidingWindowSize INTEGER NOT NULL DEFAULT 100,
      currentEwma REAL NOT NULL DEFAULT 0.0,
      backoffFactor REAL NOT NULL DEFAULT 1.0,
      isThrottlingActive INTEGER NOT NULL DEFAULT 0,
      hasBackoffTriggered INTEGER NOT NULL DEFAULT 0,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RollingGovernorRecord struct {
      DomainName          string  `json:"domainName"`
      CurrentEwma         float64 `json:"currentEwma"`
      IsThrottlingActive  bool    `json:"isThrottlingActive"`
      HasBackoffTriggered bool    `json:"hasBackoffTriggered"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Directly modulates dispatch intervals in `pkg/warmup/scheduler/throttle.go`.

---

### Improvement 614: ISP-Specific Reputation Scoring Aggregator (Google, M365, Yahoo)
- **ID**: `IMP-614`
- **Cluster**: 62
- **Goal**: Segregates delivery signals per recipient MX cluster (Gmail, Microsoft 365, Yahoo/AOL, ProtonMail) to calculate independent 0–100 health metrics per ISP.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: 4-quadrant Lipgloss overview displaying ISP logos in ASCII with scores: Google (98/100, Green), M365 (84/100, Yellow), Yahoo (92/100, Green).
  - *Web*: `<isp-reputation-quadrant>` component showing per-provider latency and inbox placement estimates.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup reputation scores --domain=<domain> [--has-per-isp-breakdown] [--is-live]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IspReputationScores (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      ispProviderKey TEXT NOT NULL, -- 'google', 'm365', 'yahoo', 'other'
      scoreValue REAL NOT NULL DEFAULT 100.0,
      sampleSize INTEGER NOT NULL DEFAULT 0,
      isHealthy INTEGER NOT NULL DEFAULT 1,
      hasWarningTriggered INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IspScoreSnapshot struct {
      DomainName          string  `json:"domainName"`
      IspProviderKey      string  `json:"ispProviderKey"`
      ScoreValue          float64 `json:"scoreValue"`
      IsHealthy           bool    `json:"isHealthy"`
      HasWarningTriggered bool    `json:"hasWarningTriggered"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Telemetry synchronized with GitMap cluster status commands (`gitmap pe -t`).

---

### Improvement 615: Automated Feedback Loop (FBL) ARF Complaint Processing Engine
- **ID**: `IMP-615`
- **Cluster**: 62
- **Goal**: Ingests Abuse Reporting Format (ARF) complaints from participating ISP feedback loops, instantly blacklisting complaining recipients and isolating sender sub-identities.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-priority Red `#E53E3E` alert banner in terminal showing: `[FBL COMPLAINT RECEIVED] Domain: example.com | Recipient: user@yahoo.com | Sender Sub-ID Suppressed`.
  - *Web*: `<fbl-complaint-log>` element showing raw headers and redacted recipient identifiers.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup fbl ingest --arf-file=<path> [--has-auto-suppress] [--is-dry-run]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS FblComplaintRecords (
      id TEXT PRIMARY KEY,
      campaignId TEXT NOT NULL,
      senderDomain TEXT NOT NULL,
      suppressedRecipientHash TEXT NOT NULL,
      abuseType TEXT NOT NULL DEFAULT 'spam',
      isSuppressedPermanently INTEGER NOT NULL DEFAULT 1,
      hasSenderPenalized INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type FblComplaintRecord struct {
      CampaignID              string `json:"campaignId"`
      SenderDomain            string `json:"senderDomain"`
      IsSuppressedPermanently bool   `json:"isSuppressedPermanently"`
      HasSenderPenalized      bool   `json:"hasSenderPenalized"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Auto-appends hashes to suppression lists in `repodb/warmup.db`.

---

### Improvement 616: Synthetic Probe Dispatcher for Half-Open Verification
- **ID**: `IMP-616`
- **Cluster**: 62
- **Goal**: Dispatches low-frequency canary messages with RFC 5322 tracking headers to validated internal seed addresses to verify whether an ISP has lifted rate-limiting.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Pulse indicator in Bubble Tea dashboard: `[PROBING] Canary ID: canary_8492 -> MX: alt1.aspmx.l.google.com -> Status: Accepted (250 OK)`. Styled with Magenta `#B794F4`.
  - *Web*: `<canary-probe-status>` component detailing probe latency distributions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup probe dispatch --domain=<domain> --isp=<target> [--has-canary-seed] [--is-synchronous]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SyntheticProbeRuns (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      targetIsp TEXT NOT NULL,
      canaryAddress TEXT NOT NULL,
      smtpResponseCode INTEGER NOT NULL,
      latencyMs INTEGER NOT NULL,
      isProbeSuccessful INTEGER NOT NULL DEFAULT 0,
      hasClearedRateLimit INTEGER NOT NULL DEFAULT 0,
      executedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SyntheticProbeResult struct {
      DomainName          string `json:"domainName"`
      TargetIsp           string `json:"targetIsp"`
      IsProbeSuccessful   bool   `json:"isProbeSuccessful"`
      HasClearedRateLimit bool   `json:"hasClearedRateLimit"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Automates recovery testing in `pkg/warmup/circuit/probe.go`.

---

### Improvement 617: DNSBL / RBL Real-Time Blacklist Sentinel Prober
- **ID**: `IMP-617`
- **Cluster**: 62
- **Goal**: Queries 40+ major DNS-based Blackhole Lists (Spamhaus ZEN, Barracuda, SpamCop, SORBS) every 15 minutes to flag sender IP or root domain blacklisting instantly.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Status grid listing queried RBLs with green checkmarks or bold flashing red warnings: `Spamhaus ZEN: [CLEAN]`, `SpamCop: [CLEAN]`.
  - *Web*: `<dnsbl-sentinel-grid>` web component rendering an interactive list with direct links to removal request workflows.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup dnsbl check --ip=<ip> --domain=<domain> [--has-full-audit] [--is-json]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DnsblCheckHistory (
      id TEXT PRIMARY KEY,
      senderIp TEXT NOT NULL,
      senderDomain TEXT NOT NULL,
      rblZoneName TEXT NOT NULL,
      returnCode TEXT, -- e.g., '127.0.0.2'
      isBlacklisted INTEGER NOT NULL DEFAULT 0,
      hasNotificationSent INTEGER NOT NULL DEFAULT 0,
      queriedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DnsblQueryResult struct {
      SenderIP            string `json:"senderIp"`
      RblZoneName         string `json:"rblZoneName"`
      IsBlacklisted       bool   `json:"isBlacklisted"`
      HasNotificationSent bool   `json:"hasNotificationSent"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Triggers emergency pauses in `gitmap email warmup halt-if-listed`.

---

### Improvement 618: Deferred Queue Depth Balancer & Exponential Re-Queue Engine
- **ID**: `IMP-618`
- **Cluster**: 62
- **Goal**: Schedules deferred messages (SMTP 421/450 responses) into exponential backoff buckets (5m, 15m, 45m, 2h) to avoid secondary bursts to already saturated destination MTAs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual ASCII queue depth bar chart using Lipgloss, categorized by destination domain and retry bucket.
  - *Web*: `<queue-depth-histogram>` element rendering bucket occupancy bars.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup queue retry-plan --domain=<domain> [--has-jitter] [--is-dry-run]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DeferredRetryQueue (
      id TEXT PRIMARY KEY,
      campaignId TEXT NOT NULL,
      recipientAddress TEXT NOT NULL,
      retryCount INTEGER NOT NULL DEFAULT 0,
      nextScheduledAttempt DATETIME NOT NULL,
      isPendingDispatch INTEGER NOT NULL DEFAULT 1,
      hasJitterApplied INTEGER NOT NULL DEFAULT 1,
      enqueuedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DeferredQueueItem struct {
      CampaignID        string `json:"campaignId"`
      RecipientAddress  string `json:"recipientAddress"`
      RetryCount        int    `json:"retryCount"`
      IsPendingDispatch bool   `json:"isPendingDispatch"`
      HasJitterApplied  bool   `json:"hasJitterApplied"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed by `pkg/warmup/queue/deferred.go`.

---

### Improvement 619: Sender Identity Reputation Decay & Recovery Model
- **ID**: `IMP-619`
- **Cluster**: 62
- **Goal**: Computes mathematical reputation recovery curves following incident resolutions, calculating safe ramp-up velocity limits over a 14-day post-trip recovery interval.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss line chart visualizing the projected 14-day recovery curve vs actual warmup volume.
  - *Web*: `<recovery-curve-model>` component with interactive slider for simulated incident severity.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup recovery project --domain=<domain> [--has-penalty-decay] [--is-conservative]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ReputationRecoveryModels (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      incidentSeverityScore REAL NOT NULL,
      recoveryDayIndex INTEGER NOT NULL DEFAULT 0,
      maxAllowedDailyVolume INTEGER NOT NULL,
      isRecovered INTEGER NOT NULL DEFAULT 0,
      hasStrictCapActive INTEGER NOT NULL DEFAULT 1,
      calculatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ReputationRecoveryRecord struct {
      DomainName         string  `json:"domainName"`
      RecoveryDayIndex   int     `json:"recoveryDayIndex"`
      IsRecovered        bool    `json:"isRecovered"`
      HasStrictCapActive bool    `json:"hasStrictCapActive"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Restricts target volume caps in `WarmupCampaigns`.

---

### Improvement 620: Real-Time Deliverability Telemetry Broadcast Event Hub
- **ID**: `IMP-620`
- **Cluster**: 62
- **Goal**: Broadcasts real-time warmup delivery outcomes, bounce classifications, and circuit transitions to local Unix sockets, named pipes, and SQLite event tables.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Live streaming log window with color-coded event icons: `⚡ TRIP`, `✔ DELIVERED`, `⏳ DEFERRED`, `⛊ PROBE`.
  - *Web*: Encapsulated WebSocket-free custom element polling local SQLite state at 100ms intervals.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup events stream [--has-circuit-events] [--is-follow]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DeliverabilityBroadcastEvents (
      id TEXT PRIMARY KEY,
      campaignId TEXT NOT NULL,
      eventType TEXT NOT NULL, -- 'delivery', 'bounce', 'circuit_change', 'throttle'
      eventPayloadJson TEXT NOT NULL,
      isDispatched INTEGER NOT NULL DEFAULT 1,
      hasConsumerAcknowledged INTEGER NOT NULL DEFAULT 0,
      emittedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DeliverabilityEventRecord struct {
      CampaignID              string `json:"campaignId"`
      EventType               string `json:"eventType"`
      IsDispatched            bool   `json:"isDispatched"`
      HasConsumerAcknowledged bool   `json:"hasConsumerAcknowledged"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Ingested by GitMap's global telemetry aggregator (`gitmap pe -t`).

---

# Cluster 63: Live MIME Composition Tree Diffing & Attachment Payload Streamer (621–630)

---

### Improvement 621: RFC 2046 Nested MIME Multipart AST Parser & Tree Visualizer
- **ID**: `IMP-621`
- **Cluster**: 63 (MIME Tree Diffing & Payload Streamer)
- **Goal**: Constructs an abstract syntax tree (AST) of nested MIME structures (`multipart/alternative`, `multipart/related`, `multipart/mixed`) to validate structural compliance.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Hierarchical Lipgloss tree with ASCII connectors:
    ```
    multipart/mixed (boundary="==b1==")
    ├── multipart/related (boundary="==b2==")
    │   ├── text/html (2.4 KB)
    │   └── image/png (cid:logo@company.com, 14.2 KB)
    └── application/pdf (contract.pdf, 184.1 KB)
    ```
  - *Web*: `<mime-ast-tree>` component with collapsible nodes and inline size indicators.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime tree --input=<eml-path> [--has-size-lint] [--is-verbose]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MimeTreeStructures (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      partIndex INTEGER NOT NULL,
      parentIndex INTEGER,
      contentType TEXT NOT NULL,
      contentDisposition TEXT,
      byteLength INTEGER NOT NULL,
      isMultipartNode INTEGER NOT NULL DEFAULT 0,
      hasValidBoundary INTEGER NOT NULL DEFAULT 1,
      parsedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MimePartNode struct {
      PartIndex        int    `json:"partIndex"`
      ParentIndex      int    `json:"parentIndex"`
      ContentType      string `json:"contentType"`
      IsMultipartNode  bool   `json:"isMultipartNode"`
      HasValidBoundary bool   `json:"hasValidBoundary"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Core parsing utility in `pkg/warmup/mime/parser.go`.

---

### Improvement 622: Visual MIME Tree Structural Diffing Engine
- **ID**: `IMP-622`
- **Cluster**: 63
- **Goal**: Diffs two versions of a warmup message template to highlight accidental removal of plain-text alternatives, missing Content-ID tags, or added high-risk attachments.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Unified diff view showing node additions in Green `+`, deletions in Red `-`, and attribute modifications in Cyan `~`.
  - *Web*: `<mime-diff-slider>` component rendering side-by-side tree comparisons.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime diff --file-a=<eml> --file-b=<eml> [--has-ignore-boundary] [--is-tree-only]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MimeDiffHistory (
      id TEXT PRIMARY KEY,
      templateId TEXT NOT NULL,
      versionA TEXT NOT NULL,
      versionB TEXT NOT NULL,
      structuralDeltaSummary TEXT NOT NULL,
      isTextAlternativePreserved INTEGER NOT NULL DEFAULT 1,
      hasAttachmentAdded INTEGER NOT NULL DEFAULT 0,
      diffedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MimeDiffReport struct {
      TemplateID                 string `json:"templateId"`
      IsTextAlternativePreserved bool   `json:"isTextAlternativePreserved"`
      HasAttachmentAdded         bool   `json:"hasAttachmentAdded"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Triggered on template edits in `gitmap email warmup template update`.

---

### Improvement 623: Zero-Copy Attachment Chunked Streamer & SHA-256 Checksummer
- **ID**: `IMP-623`
- **Cluster**: 63
- **Goal**: Streams large attachments using bounded 64KB memory buffers directly to the SMTP socket without reading entire binary payloads into RAM, calculating running SHA-256 checksums.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Real-time progress bar in Lipgloss (`[████████░░░░] 64% | 1.2 MB/s | SHA: 7f8a...`) with current buffer memory utilization.
  - *Web*: `<attachment-stream-gauge>` custom element monitoring active memory footprint.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime stream-test --file=<path> [--has-checksum-verify] [--is-bench]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS StreamedAttachmentPayloads (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      filename TEXT NOT NULL,
      mimeType TEXT NOT NULL,
      totalBytes INTEGER NOT NULL,
      sha256Hash TEXT NOT NULL,
      isZeroCopyUsed INTEGER NOT NULL DEFAULT 1,
      hasChecksumMatched INTEGER NOT NULL DEFAULT 1,
      streamedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AttachmentStreamMeta struct {
      Filename           string `json:"filename"`
      TotalBytes         int64  `json:"totalBytes"`
      IsZeroCopyUsed     bool   `json:"isZeroCopyUsed"`
      HasChecksumMatched bool   `json:"hasChecksumMatched"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Powers outbound payload transmission in `pkg/warmup/transport/streamer.go`.

---

### Improvement 624: Base64 / Quoted-Printable Transfer Encoding Integrity Linter
- **ID**: `IMP-624`
- **Cluster**: 63
- **Goal**: Verifies line-folding standards (76 chars for base64, 72 chars for QP) and flags illegal unencoded high-bit characters in 7bit ASCII transports.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual indicator of line length violations with exact row/col markers in Red `#FEB2B2`.
  - *Web*: `<encoding-linter-view>` rendering raw text with toggleable whitespace character visualization.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime lint-encoding --input=<eml> [--has-strict-7bit] [--is-fix-mode]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EncodingLintResults (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      partIndex INTEGER NOT NULL,
      encodingType TEXT NOT NULL, -- 'base64', 'quoted-printable', '7bit', '8bit'
      illegalCharCount INTEGER NOT NULL DEFAULT 0,
      isEncodingValid INTEGER NOT NULL DEFAULT 1,
      hasLineLengthViolation INTEGER NOT NULL DEFAULT 0,
      lintedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EncodingLintRecord struct {
      MessageID              string `json:"messageId"`
      EncodingType           string `json:"encodingType"`
      IsEncodingValid        bool   `json:"isEncodingValid"`
      HasLineLengthViolation bool   `json:"hasLineLengthViolation"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into pre-dispatch validation rules.

---

### Improvement 625: Content-ID (CID) Inline Asset Referencer & Orphan Detector
- **ID**: `IMP-625`
- **Cluster**: 63
- **Goal**: Cross-references HTML `<img src="cid:...">` tags with MIME `Content-ID` headers to detect broken inline images or orphaned embedded attachments bloating email payloads.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Two-column Lipgloss table matching HTML references to MIME parts. Orphan parts marked with `[ORPHAN: ATTACHMENT WASTING 42KB]` in Orange `#F6AD55`.
  - *Web*: `<cid-orphan-detector>` showing thumbnail previews of all referenced and orphaned assets.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime check-cid --input=<eml> [--has-auto-strip-orphans] [--is-json]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CidReferenceAudits (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      contentId TEXT NOT NULL,
      isReferencedInHtml INTEGER NOT NULL DEFAULT 1,
      isPartPresent INTEGER NOT NULL DEFAULT 1,
      hasOrphanPayload INTEGER NOT NULL DEFAULT 0,
      auditedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type CidAuditRecord struct {
      ContentID          string `json:"contentId"`
      IsReferencedInHtml bool   `json:"isReferencedInHtml"`
      IsPartPresent      bool   `json:"isPartPresent"`
      HasOrphanPayload   bool   `json:"hasOrphanPayload"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Diagnostic filter in `pkg/warmup/mime/cid.go`.

---

### Improvement 626: Plain-Text Alternative Parity Scorer & Spam Trigger Scanner
- **ID**: `IMP-626`
- **Cluster**: 63
- **Goal**: Evaluates lexical similarity between `text/html` and `text/plain` alternatives to prevent high spam-score penalties assigned by SpamAssassin and Bayesian filters for disparate copies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Similarity gauge `[HTML <-> Plain Parity: 96% Match]` in Green `#48BB78`. Highlights missing copy or disparate links in Red.
  - *Web*: `<multipart-parity-view>` element displaying synchronized side-by-side text diffing.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime score-parity --input=<eml> [--has-link-parity-check] [--is-strict]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PlainTextParityScores (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      similarityRatio REAL NOT NULL,
      htmlTokenCount INTEGER NOT NULL,
      plainTokenCount INTEGER NOT NULL,
      isParitySufficient INTEGER NOT NULL DEFAULT 1,
      hasDisparateLinks INTEGER NOT NULL DEFAULT 0,
      scoredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ParityScoreRecord struct {
      MessageID          string  `json:"messageId"`
      SimilarityRatio    float64 `json:"similarityRatio"`
      IsParitySufficient bool    `json:"isParitySufficient"`
      HasDisparateLinks  bool    `json:"hasDisparateLinks"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforced in warmup content generator routines.

---

### Improvement 627: Anti-Spam Zero-Font & Hidden Text CSS Heuristic Analyzer
- **ID**: `IMP-627`
- **Cluster**: 63
- **Goal**: Parses inline CSS and HTML attributes to identify deceptive techniques such as `font-size: 0px`, `display: none`, matching foreground/background colors, or `opacity: 0`.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual alert card listing detected hidden text snippets and exact CSS rules triggered. Styled with Warning Yellow `#D69E2E`.
  - *Web*: `<css-heuristic-inspector>` element showing rendered HTML with hidden elements revealed in bright outlines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime scan-hidden-text --input=<eml> [--has-color-contrast-check] [--is-blocking]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HiddenTextDetections (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      detectedSnippet TEXT NOT NULL,
      cssPropertyViolated TEXT NOT NULL,
      isDeceptiveHiddenText INTEGER NOT NULL DEFAULT 1,
      hasZeroFontSize INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HiddenTextDetectionRecord struct {
      MessageID             string `json:"messageId"`
      CssPropertyViolated   string `json:"cssPropertyViolated"`
      IsDeceptiveHiddenText bool   `json:"isDeceptiveHiddenText"`
      HasZeroFontSize       bool   `json:"hasZeroFontSize"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Security verification gate in `pkg/warmup/mime/heuristics.go`.

---

### Improvement 628: MIME Boundary Collisions & Entropy Robustness Validator
- **ID**: `IMP-628`
- **Cluster**: 63
- **Goal**: Scans message body payloads to confirm that generated multipart boundary strings never appear verbatim in attachment streams, preventing payload truncation.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Formatted string validator showing boundary entropy bits: `Boundary: ==_Part_849102 [Entropy: 128-bit] -> Collision Check: 0 Occurrences [PASS]`.
  - *Web*: `<boundary-entropy-meter>` element detailing boundary generation parameters.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime verify-boundary --input=<eml> [--has-entropy-audit] [--is-strict]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BoundaryEntropyAudits (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      boundaryString TEXT NOT NULL,
      entropyBits INTEGER NOT NULL DEFAULT 128,
      isCollisionFree INTEGER NOT NULL DEFAULT 1,
      hasStandardDelimiters INTEGER NOT NULL DEFAULT 1,
      verifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BoundaryEntropyRecord struct {
      MessageID             string `json:"messageId"`
      BoundaryString        string `json:"boundaryString"`
      IsCollisionFree       bool   `json:"isCollisionFree"`
      HasStandardDelimiters bool   `json:"hasStandardDelimiters"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Pre-dispatch invariant check in message builders.

---

### Improvement 629: Dynamic Apple Mail & Outlook MSO Conditional Comment Stripper
- **ID**: `IMP-629`
- **Cluster**: 63
- **Goal**: Evaluates Microsoft Office `<!--[if mso]>` and WebKit conditional comments in warmup message bodies to verify rendering parity across diverse email clients.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Dual tree view displaying conditional rendering blocks active for MSO vs standard standards-compliant clients.
  - *Web*: `<mso-conditional-toggle>` element toggling live client preview engines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mime test-conditionals --input=<eml> [--has-mso-audit] [--is-strip-preview]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ConditionalCommentAudits (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      msoBlockCount INTEGER NOT NULL DEFAULT 0,
      isMsoSupported INTEGER NOT NULL DEFAULT 1,
      hasUnclosedConditionals INTEGER NOT NULL DEFAULT 0,
      auditedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ConditionalCommentRecord struct {
      MessageID               string `json:"messageId"`
      IsMsoSupported          bool   `json:"isMsoSupported"`
      HasUnclosedConditionals bool   `json:"hasUnclosedConditionals"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Injected into `pkg/warmup/mime/conditional.go`.

---

### Improvement 630: Streaming RFC 822 EML File Packager & Unpackager CLI Studio
- **ID**: `IMP-630`
- **Cluster**: 63
- **Goal**: High-speed CLI utility to stream, dissect, pack, and repack `.eml` raw files and Maildir folders with instant cryptographic re-signing and terminal inspection.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-performance Bubble Tea progress bar unpacking 10,000 `.eml` files/second with thread status indicators and live IOPS metrics.
  - *Web*: `<eml-studio-workspace>` component offering offline drag-and-drop file dissection.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup eml pack --input-dir=<dir> --output=<bundle.eml> [--has-auto-sign] [--is-stream]`
  - `gitmap email warmup eml unpack --input=<bundle.eml> --output-dir=<dir> [--has-preserve-times] [--is-quiet]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EmlPackagingOperations (
      id TEXT PRIMARY KEY,
      operationType TEXT NOT NULL, -- 'pack', 'unpack'
      fileCount INTEGER NOT NULL,
      totalBytesProcessed INTEGER NOT NULL,
      durationMs INTEGER NOT NULL,
      isCompletedSuccessfully INTEGER NOT NULL DEFAULT 1,
      hasValidationErrors INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EmlPackagingRecord struct {
      OperationType           string `json:"operationType"`
      FileCount               int    `json:"fileCount"`
      TotalBytesProcessed     int64  `json:"totalBytesProcessed"`
      IsCompletedSuccessfully bool   `json:"isCompletedSuccessfully"`
      HasValidationErrors     bool   `json:"hasValidationErrors"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Foundational IO tooling in `pkg/warmup/eml/studio.go` supporting GitMap automated diagnostic exports.
