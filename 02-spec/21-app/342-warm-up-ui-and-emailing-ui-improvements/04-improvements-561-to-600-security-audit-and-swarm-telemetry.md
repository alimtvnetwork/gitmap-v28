# Specification 342: Improvements 561 to 600 — MX Sentinel, Bounce Forensics, Sharded SQLite & AGY Swarm Fleet

- **Sequence:** 342
- **Slug:** `342-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/04-improvements-561-to-600-security-audit-and-swarm-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/04-improvements-461-to-500-security-audit-and-agy-bridge.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification defines Improvements 561 through 600 for Sequence 342 of the GitMap Warm-Up and Emailing UI Engine. Sequence 342 delivers carrier-grade DNS health sentinels, automated SMTP bounce diagnostics, sharded multi-replica SQLite database governance, and high-concurrency Antigravity agent swarm telemetry:

1. **Cluster 57: Real-Time MX DNS Health & Reverse IP Blacklist Sentinel (Items 561–570):** Deep MX record validation, PTR reverse DNS cross-verification, SPF/DKIM/DMARC alignment audits, and automated DNSBL propagation surveillance.
2. **Cluster 58: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 571–580):** High-speed asynchronous bounce mail processing, enhanced DSN/RFC3463 parser, ISP-specific transient error heuristics, and auto-quarantine circuit breaking.
3. **Cluster 59: Multi-Node SQLite Sharded Cluster Dashboard & Replication Sentinel (Items 581–590):** Sharded SQLite database governance across physical nodes, WAL frame streaming health, consensus state machines, and non-blocking zero-downtime compaction.
4. **Cluster 60: Autonomous Antigravity Agent Fleet Choreographer & Telemetry Bridge (Items 591–600):** Multi-agent peer fleet orchestration, task lease renewal arbitration, token cost telemetry, distributed log streaming, and bi-directional AGY telemetry bridge.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 342 Improvements 561 to 600 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 57 (561-570)       │ Cluster 58 (571-580)       │ Cluster 59 (581-590)         │
│ Real-Time MX DNS Health    │ Forensic Bounce Analysis   │ Multi-Node SQLite Sharded    │
│ & Reverse IP Sentinel      │ Engine & Classification    │ Cluster Dashboard & Sync     │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 60 (591-600)                                    │ Cross-Cutting Standards      │
│ Autonomous Antigravity Agent Fleet Choreographer        │ Positive Booleans (is/has),  │
│ & Telemetry Bridge                                      │ Split-DB, Zero CDN, No Exits │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 57: Real-Time MX DNS Health & Reverse IP Blacklist Sentinel (Items 561–570)

### Improvement 561: Real-Time MX Record Availability & DNS Latency Sentinel
- **Title:** Real-Time MX Record Availability & DNS Latency Sentinel
- **Category:** DNS Infrastructure & Sentinel Health
- **Goal:** Probe and monitor recipient domain MX host resolution latencies, TTL expirations, and failover priority hierarchies to prevent delivery timeouts.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Color-coded MX health table with Cyan host headers, Emerald status badges (`< 45ms`), Amber latency warnings (`45–150ms`), and Crimson alerts (`> 150ms` or unreachable).
  - *Web UI:* Zero-CDN responsive SVG matrix visualizing DNS resolution latency timelines across global nameservers with embedded CSS styling.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx check --domain=<domain> [--nameserver=<ip>] [--timeout=3s] [--json]`
  - Flags: `--domain` (string, target domain), `--nameserver` (string, custom DNS resolver), `--timeout` (duration, probe deadline), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/record_sentinel.go`):
    ```go
    type MxRecordHealth struct {
        DomainName         string   `json:"domainName"`
        HostName           string   `json:"hostName"`
        Priority           int      `json:"priority"`
        LatencyMillis      int64    `json:"latencyMillis"`
        IsReachable        bool     `json:"isReachable"`
        HasValidTtl        bool     `json:"hasValidTtl"`
        HasIpv6Support     bool     `json:"hasIpv6Support"`
        EvaluatedAtUnix    int64    `json:"evaluatedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MxHealthChecks (
        checkId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        hostName TEXT NOT NULL,
        priority INTEGER NOT NULL,
        latencyMillis INTEGER NOT NULL,
        isReachable INTEGER NOT NULL,
        hasValidTtl INTEGER NOT NULL,
        hasIpv6Support INTEGER NOT NULL,
        evaluatedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_mx_health_domain ON MxHealthChecks(domainName);
    ```
- **Antigravity / GitMap Integration Point:** Emits `mail.mx.latency_alert` event on the AGY event bus if latency exceeds threshold, triggering automatic warming speed de-escalation.

---

### Improvement 562: Reverse IP (PTR) DNS Alignment & Verification Sentinel
- **Title:** Reverse IP (PTR) DNS Alignment & Verification Sentinel
- **Category:** DNS Infrastructure & Deliverability Security
- **Goal:** Mechanically verify forward-confirmed reverse DNS (FCrDNS) mapping between sending IPs and their corresponding A/PTR records.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Dual-column layout displaying forward A-records vs reverse PTR records with an Emerald checkmark or Crimson mismatch banner.
  - *Web UI:* Interactive bidirectional node diagram showing IP-to-Hostname traversal with visual validation indicators.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx verify-ptr --ip=<ip> [--expected-host=<host>] [--json]`
  - Flags: `--ip` (string, sending server IP address), `--expected-host` (string, hostname for verification), `--json` (bool, structured output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/ptr_verifier.go`):
    ```go
    type PtrAlignmentResult struct {
        IpAddress          string `json:"ipAddress"`
        ResolvedHostname   string `json:"resolvedHostname"`
        HasForwardMatch    bool   `json:"hasForwardMatch"`
        HasReverseMatch    bool   `json:"hasReverseMatch"`
        IsFcrdnsCompliant  bool   `json:"isFcrdnsCompliant"`
        CheckedAtUnix      int64  `json:"checkedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PtrAlignments (
        alignmentId TEXT PRIMARY KEY,
        ipAddress TEXT NOT NULL,
        resolvedHostname TEXT NOT NULL,
        hasForwardMatch INTEGER NOT NULL,
        hasReverseMatch INTEGER NOT NULL,
        isFcrdnsCompliant INTEGER NOT NULL,
        checkedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Integrates with GitMap node validator to block email dispatches from unaligned sending IPs.

---

### Improvement 563: Real-Time Multi-Provider DNSBL IP & Domain Reputation Sentinel
- **Title:** Real-Time Multi-Provider DNSBL IP & Domain Reputation Sentinel
- **Category:** DNSBL Sentinel & Blacklist Protection
- **Goal:** Concurrently query 50+ tier-1 DNSBL providers (Spamhaus, Barracuda, SORBS) to detect immediate blacklist listings.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Real-time scrolling progress HUD displaying provider query results, response codes, and overall threat percentage.
  - *Web UI:* Status card grid displaying listed providers in bold Crimson, clean providers in soft Emerald, and query latency gauges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx dnsbl-scan --target=<ip-or-domain> [--concurrency=20] [--json]`
  - Flags: `--target` (string, IP or domain to audit), `--concurrency` (int, simultaneous DNS queries), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/dnsbl_scanner.go`):
    ```go
    type DnsblAuditSummary struct {
        TargetIdentifier   string   `json:"targetIdentifier"`
        TotalQueried       int      `json:"totalQueried"`
        ListedCount        int      `json:"listedCount"`
        IsListed           bool     `json:"isListed"`
        HasCriticalListing bool     `json:"hasCriticalListing"`
        ListingProviders   []string `json:"listingProviders"`
        ScannedAtUnix      int64    `json:"scannedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DnsblAudits (
        auditId TEXT PRIMARY KEY,
        targetIdentifier TEXT NOT NULL,
        totalQueried INTEGER NOT NULL,
        listedCount INTEGER NOT NULL,
        isListed INTEGER NOT NULL,
        hasCriticalListing INTEGER NOT NULL,
        listingProvidersJson TEXT NOT NULL,
        scannedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_dnsbl_target ON DnsblAudits(targetIdentifier);
    ```
- **Antigravity / GitMap Integration Point:** Dispatches a critical warning to AGY multi-agent swarm to trigger autonomous unlisting or IP pool rotation.

---

### Improvement 564: Autonomous SPF Record Syntax & Flattened Include Depth Inspector
- **Title:** Autonomous SPF Record Syntax & Flattened Include Depth Inspector
- **Category:** DNS Authentication & Syntax Auditing
- **Goal:** Audit SPF TXT records for syntax validity, duplicate directives, and RFC7208 10-DNS-lookup limit exceedances.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Hierarchical ASCII tree depicting nested SPF `include:` terms with a lookup counter badge (`[8/10 Lookups Used]`).
  - *Web UI:* Collapsible visual tree graph highlighting lookups contributing toward the 10-query limit with inline syntax fix suggestions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx spf-inspect --domain=<domain> [--strict] [--json]`
  - Flags: `--domain` (string, domain name), `--strict` (bool, treat `~all` as warning vs `-all`), `--json` (bool, machine-readable output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/spf_inspector.go`):
    ```go
    type SpfInspectionReport struct {
        DomainName         string   `json:"domainName"`
        RawRecord          string   `json:"rawRecord"`
        LookupCount        int      `json:"lookupCount"`
        HasLookupExceeded  bool     `json:"hasLookupExceeded"`
        IsValidSyntax      bool     `json:"isValidSyntax"`
        HasAllMechanism    bool     `json:"hasAllMechanism"`
        InspectedAtUnix    int64    `json:"inspectedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SpfInspections (
        inspectionId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        rawRecord TEXT NOT NULL,
        lookupCount INTEGER NOT NULL,
        hasLookupExceeded INTEGER NOT NULL,
        isValidSyntax INTEGER NOT NULL,
        hasAllMechanism INTEGER NOT NULL,
        inspectedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Injects synthetic SPF flattening proposals directly into GitMap DNS management plan queues.

---

### Improvement 565: DKIM Selector Health & Cryptographic Key Length Sentinel
- **Title:** DKIM Selector Health & Cryptographic Key Length Sentinel
- **Category:** DNS Authentication & Cryptographic Integrity
- **Goal:** Validate active DKIM DNS selectors, public key availability, and verify RSA modulus length (minimum 2048-bit) or Ed25519 signature validity.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Cryptographic key summary badge indicating cipher algorithm, bit strength (Emerald `2048/4096-bit`, Crimson `< 1024-bit`), and selector state.
  - *Web UI:* Security certificate card showing key fingerprint, algorithm, bit security rating, and selector expiration timeline.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx dkim-check --domain=<domain> --selector=<selector> [--json]`
  - Flags: `--domain` (string, signing domain), `--selector` (string, DKIM selector key), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/dkim_checker.go`):
    ```go
    type DkimKeyVerification struct {
        DomainName         string `json:"domainName"`
        Selector           string `json:"selector"`
        KeyAlgorithm       string `json:"keyAlgorithm"`
        KeyBitLength       int    `json:"keyBitLength"`
        IsValidPublicKey   bool   `json:"isValidPublicKey"`
        HasAdequateLength  bool   `json:"hasAdequateLength"`
        VerifiedAtUnix     int64  `json:"verifiedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DkimVerifications (
        verificationId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        selector TEXT NOT NULL,
        keyAlgorithm TEXT NOT NULL,
        keyBitLength INTEGER NOT NULL,
        isValidPublicKey INTEGER NOT NULL,
        hasAdequateLength INTEGER NOT NULL,
        verifiedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Prevents warm-up worker threads from sending email if selector DNS keys are revoked or weak.

---

### Improvement 566: DMARC Policy Enforcement & Alignment Telemetry HUD
- **Title:** DMARC Policy Enforcement & Alignment Telemetry HUD
- **Category:** DNS Authentication & DMARC Surveillance
- **Goal:** Monitor published DMARC policy tags (`p=none`, `p=quarantine`, `p=reject`), percentage coverage (`pct`), and aggregate reporting endpoints (`rua`/`ruf`).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-contrast shield card detailing policy status, reporting endpoints, and strictness rating.
  - *Web UI:* DMARC policy status pill accompanied by an alignment preview dial and recipient aggregate email destination validator.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx dmarc-audit --domain=<domain> [--json]`
  - Flags: `--domain` (string, sender domain to verify), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/dmarc_auditor.go`):
    ```go
    type DmarcPolicyAudit struct {
        DomainName         string `json:"domainName"`
        PolicyAction       string `json:"policyAction"`
        SubdomainPolicy    string `json:"subdomainPolicy"`
        PercentageEnforced int    `json:"percentageEnforced"`
        HasReportingUri    bool   `json:"hasReportingUri"`
        IsStrictAlignment  bool   `json:"isStrictAlignment"`
        AuditedAtUnix      int64  `json:"auditedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS DmarcAudits (
        auditId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        policyAction TEXT NOT NULL,
        subdomainPolicy TEXT NOT NULL,
        percentageEnforced INTEGER NOT NULL,
        hasReportingUri INTEGER NOT NULL,
        isStrictAlignment INTEGER NOT NULL,
        auditedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Supplies DMARC policy state to warmup planner so velocity ramp curves match domain authentication readiness.

---

### Improvement 567: MTA-STS & TLS-RPT DNS Sentinel & Policy Fetcher
- **Title:** MTA-STS & TLS-RPT DNS Sentinel & Policy Fetcher
- **Category:** Mail Transport Security & Policy Compliance
- **Goal:** Verify presence of `_mta-sts` and `_smtp._tls` DNS TXT records, download the HTTPS STS policy, and validate TLS enforcement mode.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Security padlock HUD displaying STS mode (`enforce` vs `testing`), HTTPS policy TTL, and MX host pattern match list.
  - *Web UI:* TLS transport verification badge showing TLS version compatibility, certificate status, and policy cache age.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx mta-sts-audit --domain=<domain> [--json]`
  - Flags: `--domain` (string, domain name), `--json` (bool, JSON output format).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/mta_sts_auditor.go`):
    ```go
    type MtaStsPolicyReport struct {
        DomainName         string   `json:"domainName"`
        PolicyMode         string   `json:"policyMode"`
        MaxAgeSeconds      int64    `json:"maxAgeSeconds"`
        HasValidRecord     bool     `json:"hasValidRecord"`
        HasHttpsPolicy     bool     `json:"hasHttpsPolicy"`
        IsEnforced         bool     `json:"isEnforced"`
        FetchedAtUnix      int64    `json:"fetchedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MtaStsReports (
        reportId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        policyMode TEXT NOT NULL,
        maxAgeSeconds INTEGER NOT NULL,
        hasValidRecord INTEGER NOT NULL,
        hasHttpsPolicy INTEGER NOT NULL,
        isEnforced INTEGER NOT NULL,
        fetchedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Automatically logs transport security compliance to `gitmap.db` audit ledger for enterprise deliverability audits.

---

### Improvement 568: DNS Propagation Speed & Anycast Global Resolver Monitor
- **Title:** DNS Propagation Speed & Anycast Global Resolver Monitor
- **Category:** DNS Infrastructure & Propagation Telemetry
- **Goal:** Benchmark DNS propagation time across geographically distributed Anycast resolvers (Cloudflare, Google, OpenDNS, Quad9, Local).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Geographic multi-node terminal gauge with resolver name, resolved value, TTL, and status pills.
  - *Web UI:* Global Anycast world map visualizer with regional latency markers and propagation percentage ring.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx propagation-check --domain=<domain> --record-type=TXT [--expected=<val>] [--json]`
  - Flags: `--domain` (string, target domain), `--record-type` (string, A, MX, TXT), `--expected` (string, expected record content), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/propagation_monitor.go`):
    ```go
    type PropagationCheckResult struct {
        DomainName         string `json:"domainName"`
        RecordType         string `json:"recordType"`
        ResolverIp         string `json:"resolverIp"`
        ResolverName       string `json:"resolverName"`
        IsPropagated       bool   `json:"isPropagated"`
        HasMatchingValue   bool   `json:"hasMatchingValue"`
        LatencyMillis      int64  `json:"latencyMillis"`
        CheckedAtUnix      int64  `json:"checkedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS PropagationChecks (
        checkId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        recordType TEXT NOT NULL,
        resolverIp TEXT NOT NULL,
        resolverName TEXT NOT NULL,
        isPropagated INTEGER NOT NULL,
        hasMatchingValue INTEGER NOT NULL,
        latencyMillis INTEGER NOT NULL,
        checkedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Suspends newly provisioned sending identities until propagation achieves 100% across all standard Anycast resolvers.

---

### Improvement 569: BIMI (Brand Indicators) SVG & VMC DNS Readiness Checker
- **Title:** BIMI (Brand Indicators) SVG & VMC DNS Readiness Checker
- **Category:** Brand Deliverability & Visual Identity
- **Goal:** Validate `default._bimi` DNS record, verify Tiny-SVG-1.2 format compliance of brand avatar, and audit Verified Mark Certificate (VMC) chains.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Brand indicator verification card displaying SVG dimensions, XML namespace validation, and VMC certificate authority.
  - *Web UI:* Brand avatar live visualizer rendering the validated SVG alongside VMC issuer trust badges and DMARC enforcement prerequisites.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx bimi-check --domain=<domain> [--json]`
  - Flags: `--domain` (string, domain name), `--json` (bool, output as JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/bimi_checker.go`):
    ```go
    type BimiReadinessReport struct {
        DomainName         string `json:"domainName"`
        SvgLocationUrl     string `json:"svgLocationUrl"`
        VmcCertificateUrl  string `json:"vmcCertificateUrl"`
        IsValidSvgProfile  bool   `json:"isValidSvgProfile"`
        HasValidCertificate bool  `json:"hasValidCertificate"`
        IsBimiEligible     bool   `json:"isBimiEligible"`
        AuditedAtUnix      int64  `json:"auditedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BimiReadiness (
        reportId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        svgLocationUrl TEXT NOT NULL,
        vmcCertificateUrl TEXT NOT NULL,
        isValidSvgProfile INTEGER NOT NULL,
        hasValidCertificate INTEGER NOT NULL,
        isBimiEligible INTEGER NOT NULL,
        auditedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Supplies visual brand metadata to the warmup email template renderer for verified inbox display simulations.

---

### Improvement 570: MX Cluster Health Composite Dashboard & Alerting Bridge
- **Title:** MX Cluster Health Composite Dashboard & Alerting Bridge
- **Category:** DNS Aggregation & Sentinel Dashboard
- **Goal:** Aggregate all Cluster 57 DNS health metrics (MX, PTR, DNSBL, SPF, DKIM, DMARC, STS, BIMI) into a unified composite deliverability score.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Master DNS terminal scorecard with Lipgloss double-borders, overall health grade (`A+` to `F`), and sub-system status indicators.
  - *Web UI:* Unified zero-CDN executive dashboard with radar chart representing 8 DNS vectors and instant remediation action buttons.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mx dashboard --domain=<domain> [--watch] [--interval=10s] [--json]`
  - Flags: `--domain` (string, domain), `--watch` (bool, continuous refresh), `--interval` (duration, polling cadence), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/mx/composite_dashboard.go`):
    ```go
    type CompositeDnsHealth struct {
        DomainName         string  `json:"domainName"`
        CompositeScore     float64 `json:"compositeScore"`
        GradeLetter        string  `json:"gradeLetter"`
        IsPassingAllChecks bool    `json:"isPassingAllChecks"`
        HasSecurityFault   bool    `json:"hasSecurityFault"`
        HasBlacklistFault  bool    `json:"hasBlacklistFault"`
        GeneratedAtUnix    int64   `json:"generatedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_dns.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CompositeDnsReports (
        reportId TEXT PRIMARY KEY,
        domainName TEXT NOT NULL,
        compositeScore REAL NOT NULL,
        gradeLetter TEXT NOT NULL,
        isPassingAllChecks INTEGER NOT NULL,
        hasSecurityFault INTEGER NOT NULL,
        hasBlacklistFault INTEGER NOT NULL,
        generatedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_composite_domain ON CompositeDnsReports(domainName);
    ```
- **Antigravity / GitMap Integration Point:** Serves as the primary DNS health gate for the Antigravity swarm choreographer before launching warm-up phases.

---

## 3. Cluster 58: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 571–580)

### Improvement 571: High-Speed RFC 3463 / RFC 3464 DSN Bounce Message Parser
- **Title:** High-Speed RFC 3463 / RFC 3464 DSN Bounce Message Parser
- **Category:** Bounce Forensics & RFC Compliance
- **Goal:** Parse raw delivery status notifications (DSN) and non-delivery reports (NDR) into structured status codes (`5.X.X` permanent, `4.X.X` persistent transient).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Terminal forensic diagnostic card highlighting action, status code, diagnostic code, and final recipient.
  - *Web UI:* Split-pane email parser view displaying raw MIME multipart/report on left and structured parsed AST on right.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce parse --file=<path-to-eml> [--json]`
  - Flags: `--file` (string, raw bounce email file), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/dsn_parser.go`):
    ```go
    type ParsedBounceDsn struct {
        OriginalMessageId  string `json:"originalMessageId"`
        RecipientAddress   string `json:"recipientAddress"`
        ReportingMta       string `json:"reportingMta"`
        EnhancedStatusCode string `json:"enhancedStatusCode"`
        DiagnosticCode     string `json:"diagnosticCode"`
        IsPermanentFailure bool   `json:"isPermanentFailure"`
        HasEnhancedCode    bool   `json:"hasEnhancedCode"`
        ParsedAtUnix       int64  `json:"parsedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ParsedBounces (
        bounceId TEXT PRIMARY KEY,
        originalMessageId TEXT NOT NULL,
        recipientAddress TEXT NOT NULL,
        reportingMta TEXT NOT NULL,
        enhancedStatusCode TEXT NOT NULL,
        diagnosticCode TEXT NOT NULL,
        isPermanentFailure INTEGER NOT NULL,
        hasEnhancedCode INTEGER NOT NULL,
        parsedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_bounce_recipient ON ParsedBounces(recipientAddress);
    ```
- **Antigravity / GitMap Integration Point:** Directly feeds structured bounce records to the reputation decay calculator to adjust warm-up curves.

---

### Improvement 572: Bayesian Semantic Bounce Reason Categorizer
- **Title:** Bayesian Semantic Bounce Reason Categorizer
- **Category:** Bounce Intelligence & NLP Classification
- **Goal:** Classify unstructured and non-standard bounce error strings (e.g., custom ISP rejections, challenge-response systems) into discrete failure categories.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Diagnostic category tag (`[SPAM_BLOCK]`, `[MAILBOX_FULL]`, `[DOMAIN_UNKNOWN]`) with classification confidence percentage.
  - *Web UI:* Color-coded category pill with hover tooltip explaining semantic classification rationale and matched token weights.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce classify --text="<diagnostic-text>" [--json]`
  - Flags: `--text` (string, raw error string from MTA), `--json` (bool, JSON response).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/semantic_categorizer.go`):
    ```go
    type BounceClassification struct {
        CategoryName       string  `json:"categoryName"`
        ConfidenceScore    float64 `json:"confidenceScore"`
        MatchedKeywords    []string`json:"matchedKeywords"`
        IsSpamRelated      bool    `json:"isSpamRelated"`
        IsMailboxExistent  bool    `json:"isMailboxExistent"`
        HasTransientNature bool    `json:"hasTransientNature"`
        ClassifiedAtUnix   int64   `json:"classifiedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BounceClassifications (
        classificationId TEXT PRIMARY KEY,
        categoryName TEXT NOT NULL,
        confidenceScore REAL NOT NULL,
        matchedKeywordsJson TEXT NOT NULL,
        isSpamRelated INTEGER NOT NULL,
        isMailboxExistent INTEGER NOT NULL,
        hasTransientNature INTEGER NOT NULL,
        classifiedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Notifies Antigravity agent when `isSpamRelated` is true to revise spintax generation and subject line formulations.

---

### Improvement 573: ISP-Specific Error Code Translator & Knowledge Matrix
- **Title:** ISP-Specific Error Code Translator & Knowledge Matrix
- **Category:** ISP Forensics & Deliverability Troubleshooting
- **Goal:** Translate esoteric error codes from Google Workspace, Microsoft 365, Yahoo, and Comcast into human-actionable remediation instructions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Two-tone remediation prompt with ISP brand icon, literal error code explanation, and concrete next-step checklist.
  - *Web UI:* ISP troubleshooting modal with direct links to Google Postmaster Tools, Microsoft SNDS, and automated resolution playbooks.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce translate --isp=<google|microsoft|yahoo> --code=<code> [--json]`
  - Flags: `--isp` (string, target provider), `--code` (string, diagnostic error code), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/isp_translator.go`):
    ```go
    type IspErrorTranslation struct {
        IspIdentifier      string `json:"ispIdentifier"`
        RawErrorCode       string `json:"rawErrorCode"`
        FriendlyMessage    string `json:"friendlyMessage"`
        RemediationStep    string `json:"remediationStep"`
        IsRateLimit        bool   `json:"isRateLimit"`
        HasReputationFault bool   `json:"hasReputationFault"`
        TranslatedAtUnix   int64  `json:"translatedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS IspTranslations (
        translationId TEXT PRIMARY KEY,
        IspIdentifier TEXT NOT NULL,
        rawErrorCode TEXT NOT NULL,
        friendlyMessage TEXT NOT NULL,
        remediationStep TEXT NOT NULL,
        isRateLimit INTEGER NOT NULL,
        hasReputationFault INTEGER NOT NULL,
        translatedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Automatically modulates sender retry queues and velocity schedules based on whether the translation indicates rate limiting.

---

### Improvement 574: Automated Circuit-Breaking Suppression List Manager
- **Title:** Automated Circuit-Breaking Suppression List Manager
- **Category:** Deliverability Protection & Suppression Governance
- **Goal:** Automatically isolate and quarantine addresses receiving hard bounces (5.X.X) or repeat soft bounces across all nodes without manual intervention.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Safety shield banner displaying quarantined address, quarantine reason, and time-to-live expiration.
  - *Web UI:* Searchable suppression table with one-click manual release, audit log history, and bulk export capabilities.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce suppress add --email=<email> --reason=<reason> [--duration=720h] [--json]`
  - Flags: `--email` (string, target address), `--reason` (string, reason), `--duration` (duration, suppression TTL), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/suppression_manager.go`):
    ```go
    type QuarantinedAddress struct {
        EmailAddress       string `json:"emailAddress"`
        SuppressionReason  string `json:"suppressionReason"`
        ExpiresAtUnix      int64  `json:"expiresAtUnix"`
        IsPermanent        bool   `json:"isPermanent"`
        HasActiveOverride  bool   `json:"hasActiveOverride"`
        CreatedAtUnix      int64  `json:"createdAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SuppressedAddresses (
        suppressionId TEXT PRIMARY KEY,
        emailAddress TEXT NOT NULL UNIQUE,
        suppressionReason TEXT NOT NULL,
        expiresAtUnix INTEGER NOT NULL,
        isPermanent INTEGER NOT NULL,
        hasActiveOverride INTEGER NOT NULL,
        createdAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_suppress_email ON SuppressedAddresses(emailAddress);
    ```
- **Antigravity / GitMap Integration Point:** Synchronizes global suppression lists across all connected GitMap nodes during pre-flight checks.

---

### Improvement 575: Transient Retry Backoff Calculator with Jitter Telemetry
- **Title:** Transient Retry Backoff Calculator with Jitter Telemetry
- **Category:** Queue Optimization & Algorithmic Resilience
- **Goal:** Compute exponential backoff delay schedules with randomized decorrelated jitter for transient bounces (4.X.X) to prevent SMTP thundering herds.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Step-wise timeline showing retry attempt numbers, computed backoff delay seconds, and jitter offsets.
  - *Web UI:* Interactive retry simulation graph demonstrating queued retry distribution over a 24-hour window.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce backoff-calc --attempt=3 [--base-delay=30s] [--max-delay=2h] [--json]`
  - Flags: `--attempt` (int, current attempt count), `--base-delay` (duration), `--max-delay` (duration), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/backoff_calculator.go`):
    ```go
    type BackoffScheduleResult struct {
        AttemptNumber      int    `json:"attemptNumber"`
        DelaySeconds       int64  `json:"delaySeconds"`
        JitterOffsetMillis int64  `json:"jitterOffsetMillis"`
        NextRunAtUnix      int64  `json:"nextRunAtUnix"`
        IsMaxDelayReached  bool   `json:"isMaxDelayReached"`
        HasRetryRemaining  bool   `json:"hasRetryRemaining"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RetrySchedules (
        scheduleId TEXT PRIMARY KEY,
        attemptNumber INTEGER NOT NULL,
        delaySeconds INTEGER NOT NULL,
        jitterOffsetMillis INTEGER NOT NULL,
        nextRunAtUnix INTEGER NOT NULL,
        isMaxDelayReached INTEGER NOT NULL,
        hasRetryRemaining INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Dispatches scheduled task wakeups into the Antigravity timer queue when retry windows mature.

---

### Improvement 576: Feedback Loop (FBL) ARF Complaint Forensics Analyzer
- **Title:** Feedback Loop (FBL) ARF Complaint Forensics Analyzer
- **Category:** Compliance & Complaint Analysis
- **Goal:** Ingest and parse Abuse Reporting Format (ARF) complaints from Yahoo, Microsoft, and other FBL providers to isolate offending headers and campaigns.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Warning incident card highlighting complaining ISP, reported sender IP, original user agent, and feedback type (`abuse`).
  - *Web UI:* Incident forensic view with masked recipient identifiers, campaign source attribution, and one-click seed revocation.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce parse-fbl --file=<path-to-arf> [--json]`
  - Flags: `--file` (string, path to ARF complaint), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/fbl_analyzer.go`):
    ```go
    type ArfComplaintReport struct {
        ComplaintId        string `json:"complaintId"`
        ReportingProvider  string `json:"reportingProvider"`
        FeedbackType       string `json:"feedbackType"`
        ReportedSenderIp   string `json:"reportedSenderIp"`
        OriginalMessageId  string `json:"originalMessageId"`
        IsActionable       bool   `json:"isActionable"`
        HasSenderAuthMatch bool   `json:"hasSenderAuthMatch"`
        LoggedAtUnix       int64  `json:"loggedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FblComplaints (
        complaintId TEXT PRIMARY KEY,
        reportingProvider TEXT NOT NULL,
        feedbackType TEXT NOT NULL,
        reportedSenderIp TEXT NOT NULL,
        originalMessageId TEXT NOT NULL,
        isActionable INTEGER NOT NULL,
        hasSenderAuthMatch INTEGER NOT NULL,
        loggedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_fbl_sender ON FblComplaints(reportedSenderIp);
    ```
- **Antigravity / GitMap Integration Point:** Emits high-priority `mail.complaint.received` event, causing swarm coordinator to immediately pause the offending campaign.

---

### Improvement 577: Challenge-Response & Auto-Reply Disposable Filter
- **Title:** Challenge-Response & Auto-Reply Disposable Filter
- **Category:** Bounce Forensics & Noise Reduction
- **Goal:** Distinguish automated Out-of-Office (OOO), vacation auto-replies, and challenge-response filters (e.g. Boxbe) from genuine delivery failures.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Subdued gray info pill indicating auto-reply suppression, sender address, and vacation return date without incrementing bounce counters.
  - *Web UI:* Filtered auto-reply activity feed clearly separated from genuine failure metrics to keep reputation calculations accurate.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce filter-auto --subject="<subject>" --body="<body>" [--json]`
  - Flags: `--subject` (string, email subject), `--body` (string, message body), `--json` (bool, JSON response).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/auto_reply_filter.go`):
    ```go
    type AutoReplyDetection struct {
        HeaderPrecedence   string `json:"headerPrecedence"`
        IsAutoSubmitted    bool   `json:"isAutoSubmitted"`
        IsVacationNotice   bool   `json:"isVacationNotice"`
        HasChallengePrompt bool   `json:"hasChallengePrompt"`
        ShouldIgnoreBounce bool   `json:"shouldIgnoreBounce"`
        DetectedAtUnix     int64  `json:"detectedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AutoReplyDetections (
        detectionId TEXT PRIMARY KEY,
        headerPrecedence TEXT NOT NULL,
        isAutoSubmitted INTEGER NOT NULL,
        isVacationNotice INTEGER NOT NULL,
        hasChallengePrompt INTEGER NOT NULL,
        shouldIgnoreBounce INTEGER NOT NULL,
        detectedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Prevents false-positive reputation penalties in GitMap deliverability analytics.

---

### Improvement 578: Bounce Cluster Pattern Matcher & ISP Trend Telemetry
- **Title:** Bounce Cluster Pattern Matcher & ISP Trend Telemetry
- **Category:** Bounce Intelligence & Predictive Clustering
- **Goal:** Detect sudden spikes of similar bounce messages across multiple recipient domains sharing common MX infrastructure (e.g., Proofpoint, Mimecast).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Clustered anomaly report displaying MX cluster name, error frequency, and recommended cooldown period.
  - *Web UI:* Heatmap matrix representing bounce frequency indexed by gateway provider vs time buckets.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce clusters --window=1h [--min-cluster=5] [--json]`
  - Flags: `--window` (duration, time window), `--min-cluster` (int, minimum incident threshold), `--json` (bool, JSON format).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/cluster_matcher.go`):
    ```go
    type BounceClusterSummary struct {
        ClusterKey         string  `json:"clusterKey"`
        GatewayProvider    string  `json:"gatewayProvider"`
        OccurrenceCount    int     `json:"occurrenceCount"`
        AverageVelocity    float64 `json:"averageVelocity"`
        IsAnomalyDetected  bool    `json:"isAnomalyDetected"`
        HasMitigationPlan  bool    `json:"hasMitigationPlan"`
        AnalyzedAtUnix     int64   `json:"analyzedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BounceClusters (
        clusterId TEXT PRIMARY KEY,
        clusterKey TEXT NOT NULL,
        gatewayProvider TEXT NOT NULL,
        occurrenceCount INTEGER NOT NULL,
        averageVelocity REAL NOT NULL,
        isAnomalyDetected INTEGER NOT NULL,
        hasMitigationPlan INTEGER NOT NULL,
        analyzedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Autonomous warm-up coordinator executes automated throttle rules specifically for that gateway infrastructure.

---

### Improvement 579: Seed Account Engagement Feedback Loop Connector
- **Title:** Seed Account Engagement Feedback Loop Connector
- **Category:** Warmup Telemetry & Inbox Engagement
- **Goal:** Synchronize with IMAP/API seed account testboxes to record opens, reply interactions, and marking messages as "not spam" (spam recovery).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Engagement tracker table with Green upward arrows for inbox placement, Yellow for promotion tabs, and Red for spam rescues.
  - *Web UI:* Real-time placement breakdown bar chart showing Inbox vs Spam vs Promotion percentages across major providers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce engagement-sync --seed-pool=<pool-id> [--json]`
  - Flags: `--seed-pool` (string, identifier of the seed testing pool), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/engagement_connector.go`):
    ```go
    type SeedEngagementMetrics struct {
        SeedAddress        string  `json:"seedAddress"`
        FolderLocation     string  `json:"folderLocation"`
        IsOpened           bool    `json:"isOpened"`
        IsRescuedFromSpam  bool    `json:"isRescuedFromSpam"`
        HasReplied         bool    `json:"hasReplied"`
        LatencyToOpenSecs  int64   `json:"latencyToOpenSecs"`
        SynchronizedAtUnix int64   `json:"synchronizedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SeedEngagements (
        engagementId TEXT PRIMARY KEY,
        seedAddress TEXT NOT NULL,
        folderLocation TEXT NOT NULL,
        isOpened INTEGER NOT NULL,
        isRescuedFromSpam INTEGER NOT NULL,
        hasReplied INTEGER NOT NULL,
        latencyToOpenSecs INTEGER NOT NULL,
        synchronizedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_seed_address ON SeedEngagements(seedAddress);
    ```
- **Antigravity / GitMap Integration Point:** Feeds real-time reward signals directly into the Antigravity warmup reinforcement learning loop.

---

### Improvement 580: Live Forensic Bounce Diagnostic HUD & Action Center
- **Title:** Live Forensic Bounce Diagnostic HUD & Action Center
- **Category:** Terminal TUI & Diagnostic Operator Studio
- **Goal:** Deliver an interactive Lipgloss terminal workspace showing real-time incoming bounce streams, active quarantines, and manual release options.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Split-view interactive Bubble Tea TUI with live log stream on top, bounce categorization analytics on bottom, and keyboard shortcuts (`q` quit, `r` release, `s` suppress).
  - *Web UI:* Action center console with real-time websocket updates, audio alerts on critical spikes, and batch remediation controls.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce hud [--follow] [--refresh=2s] [--json]`
  - Flags: `--follow` (bool, live stream mode), `--refresh` (duration, UI update interval), `--json` (bool, emit initial state as JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/bounce/bounce_hud.go`):
    ```go
    type BounceHudSnapshot struct {
        ActiveQuarantines  int     `json:"activeQuarantines"`
        BouncesLastHour    int     `json:"bouncesLastHour"`
        HardBounceRate     float64 `json:"hardBounceRate"`
        IsUnderThreshold   bool    `json:"isUnderThreshold"`
        HasCriticalAlarm   bool    `json:"hasCriticalAlarm"`
        GeneratedAtUnix    int64   `json:"generatedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/mail_bounce.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS BounceHudSnapshots (
        snapshotId TEXT PRIMARY KEY,
        activeQuarantines INTEGER NOT NULL,
        bouncesLastHour INTEGER NOT NULL,
        hardBounceRate REAL NOT NULL,
        isUnderThreshold INTEGER NOT NULL,
        hasCriticalAlarm INTEGER NOT NULL,
        generatedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Provides human-in-the-loop diagnostic controls and real-time oversight for autonomous mail agent swarms.

---

## 4. Cluster 59: Multi-Node SQLite Sharded Cluster Dashboard & Replication Sentinel (Items 581–590)

### Improvement 581: Sharded SQLite Cluster Discovery & Topology Mesh Mapper
- **Title:** Sharded SQLite Cluster Discovery & Topology Mesh Mapper
- **Category:** Distributed Storage & Cluster Topology
- **Goal:** Autonomously discover and register SQLite database shards across cluster nodes, mapping replica topologies and shard keys.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Tree-style node mesh visualization showing node IDs, assigned shard ranges, storage paths, and health indicators.
  - *Web UI:* Interactive zero-CDN network graph with dynamic nodes showing master vs replica shards and connection latency links.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard discover [--config=<cluster.json>] [--json]`
  - Flags: `--config` (string, cluster configuration path), `--json` (bool, output topology as JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/topology_mapper.go`):
    ```go
    type ShardTopologyNode struct {
        NodeIdentifier     string `json:"nodeIdentifier"`
        ShardRangeStart    int64  `json:"shardRangeStart"`
        ShardRangeEnd      int64  `json:"shardRangeEnd"`
        IsMasterNode       bool   `json:"isMasterNode"`
        HasActiveReplicas  bool   `json:"hasActiveReplicas"`
        DiscoveredAtUnix   int64  `json:"discoveredAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShardTopologyNodes (
        nodeId TEXT PRIMARY KEY,
        nodeIdentifier TEXT NOT NULL,
        shardRangeStart INTEGER NOT NULL,
        shardRangeEnd INTEGER NOT NULL,
        isMasterNode INTEGER NOT NULL,
        hasActiveReplicas INTEGER NOT NULL,
        discoveredAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Coordinates with GitMap cluster fleet command (`gitmap cluster`) to route database read/write queries to correct local shards.

---

### Improvement 582: WAL Frame Streaming Replication Lag Sentinel
- **Title:** WAL Frame Streaming Replication Lag Sentinel
- **Category:** Distributed Storage & Replication Lag
- **Goal:** Measure replication lag across SQLite WAL streaming connections in terms of unapplied frame counts and microsecond latency.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Linear horizontal progress bar indicating frame synchronization state (Green `< 5 frames`, Yellow `5–50 frames`, Red `> 50 frames`).
  - *Web UI:* Multi-series time-series line chart tracking per-node WAL frame lag and network throughput.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard lag-check --shard-id=<id> [--threshold=50] [--json]`
  - Flags: `--shard-id` (string, database shard identifier), `--threshold` (int, maximum allowed uncommitted frames), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/wal_lag_sentinel.go`):
    ```go
    type WalReplicationLag struct {
        ShardIdentifier    string `json:"shardIdentifier"`
        ReplicaNodeId      string `json:"replicaNodeId"`
        UnappliedFrames    int64  `json:"unappliedFrames"`
        LagDurationMillis  int64  `json:"lagDurationMillis"`
        IsSynchronized     bool   `json:"isSynchronized"`
        HasExceededLimit   bool   `json:"hasExceededLimit"`
        SampledAtUnix      int64  `json:"sampledAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS WalReplicationLags (
        lagId TEXT PRIMARY KEY,
        shardIdentifier TEXT NOT NULL,
        replicaNodeId TEXT NOT NULL,
        unappliedFrames INTEGER NOT NULL,
        lagDurationMillis INTEGER NOT NULL,
        isSynchronized INTEGER NOT NULL,
        hasExceededLimit INTEGER NOT NULL,
        sampledAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_wal_lag_shard ON WalReplicationLags(shardIdentifier);
    ```
- **Antigravity / GitMap Integration Point:** Halts write promotions on replica nodes if `hasExceededLimit` is true to avoid split-brain replication divergences.

---

### Improvement 583: Distributed Consensus & Epoch Lease Coordinator
- **Title:** Distributed Consensus & Epoch Lease Coordinator
- **Category:** Distributed Systems & Master Election
- **Goal:** Maintain leader leases and epoch counters across SQLite shards using time-bounded atomic lease renewals to avoid dual-master conflicts.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Lease tenure badge showing current epoch number, leader node ID, remaining lease seconds, and renewal heartbeats.
  - *Web UI:* Consensus state diagram showing live heartbeat pulses and smooth failover election sequences.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard lease-status --shard-id=<id> [--json]`
  - Flags: `--shard-id` (string, shard ID), `--json` (bool, output status as JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/epoch_coordinator.go`):
    ```go
    type ShardLeaderLease struct {
        ShardIdentifier    string `json:"shardIdentifier"`
        CurrentLeaderNode  string `json:"currentLeaderNode"`
        EpochCounter       int64  `json:"epochCounter"`
        LeaseExpiresAtUnix int64  `json:"leaseExpiresAtUnix"`
        IsLeaseActive      bool   `json:"isLeaseActive"`
        HasQuorumAgreement bool   `json:"hasQuorumAgreement"`
        RenewedAtUnix      int64  `json:"renewedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShardLeaderLeases (
        leaseId TEXT PRIMARY KEY,
        shardIdentifier TEXT NOT NULL UNIQUE,
        currentLeaderNode TEXT NOT NULL,
        epochCounter INTEGER NOT NULL,
        leaseExpiresAtUnix INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL,
        hasQuorumAgreement INTEGER NOT NULL,
        renewedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Arbitrates which Antigravity agent node possesses write authority for global warmup campaign schedules.

---

### Improvement 584: Cross-Node SQLite Checkpoint & Incremental Auto-Vacuum Governor
- **Title:** Cross-Node SQLite Checkpoint & Incremental Auto-Vacuum Governor
- **Category:** Database Maintenance & Performance Optimization
- **Goal:** Coordinate `PRAGMA wal_checkpoint(TRUNCATE)` and incremental auto-vacuums during low-traffic windows across all cluster shards.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Progress spinner showing checkpoint execution, reclaimed disk space in MB, and locked page metrics.
  - *Web UI:* Cluster storage defragmentation card showing reclaimed storage trends and scheduling controls.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard vacuum --shard-id=<id> [--mode=PASSIVE|FULL|RESTART|TRUNCATE] [--json]`
  - Flags: `--shard-id` (string, shard ID), `--mode` (string, checkpoint mode), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/vacuum_governor.go`):
    ```go
    type CheckpointExecutionReport struct {
        ShardIdentifier    string `json:"shardIdentifier"`
        CheckpointMode     string `json:"checkpointMode"`
        PagesCheckpointed  int    `json:"pagesCheckpointed"`
        BytesReclaimed     int64  `json:"bytesReclaimed"`
        IsCompleted        bool   `json:"isCompleted"`
        HasBlockedReaders  bool   `json:"hasBlockedReaders"`
        ExecutedAtUnix     int64  `json:"executedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS CheckpointReports (
        reportId TEXT PRIMARY KEY,
        shardIdentifier TEXT NOT NULL,
        checkpointMode TEXT NOT NULL,
        pagesCheckpointed INTEGER NOT NULL,
        bytesReclaimed INTEGER NOT NULL,
        isCompleted INTEGER NOT NULL,
        hasBlockedReaders INTEGER NOT NULL,
        executedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Ensures local SQLite instances never experience locking or unbounded WAL growth during intensive warm-up runs.

---

### Improvement 585: Sharded Query Router & Monadic Result Aggregator
- **Title:** Sharded Query Router & Monadic Result Aggregator
- **Category:** Distributed Querying & Data Federation
- **Goal:** Scatter-gather read queries across multiple SQLite shard databases and combine individual monadic `Result[T]` responses into unified datasets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Execution breakdown table showing per-shard query latency, returned row count, and federation assembly duration.
  - *Web UI:* Query profiling tree detailing scatter execution times, parallel worker performance, and consolidated table views.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard query --query="<sql>" [--parallel=4] [--json]`
  - Flags: `--query` (string, read-only SQL query), `--parallel` (int, worker concurrency), `--json` (bool, JSON result).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/query_router.go`):
    ```go
    type ShardedQuerySummary struct {
        QueryFingerprint   string `json:"queryFingerprint"`
        ShardsTargeted     int    `json:"shardsTargeted"`
        TotalRowsReturned  int64  `json:"totalRowsReturned"`
        ExecutionTimeMillis int64 `json:"executionTimeMillis"`
        IsAllShardsSuccess bool   `json:"isAllShardsSuccess"`
        HasPartialResults  bool   `json:"hasPartialResults"`
        ExecutedAtUnix     int64  `json:"executedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShardedQueryLogs (
        logId TEXT PRIMARY KEY,
        queryFingerprint TEXT NOT NULL,
        shardsTargeted INTEGER NOT NULL,
        totalRowsReturned INTEGER NOT NULL,
        executionTimeMillis INTEGER NOT NULL,
        isAllShardsSuccess INTEGER NOT NULL,
        hasPartialResults INTEGER NOT NULL,
        executedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Exposes federated telemetry metrics to Antigravity planning scripts without requiring manual shard navigation.

---

### Improvement 586: Shard Split & Rebalance Migration Orchestrator
- **Title:** Shard Split & Rebalance Migration Orchestrator
- **Category:** Cluster Resizing & Data Migration
- **Goal:** Autonomously partition oversized shards (> 5 GB) into child shards with zero downtime using non-blocking online row migration.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive migration progress bar showing copy phase, dual-write synchronization phase, and cutover point.
  - *Web UI:* Storage rebalance diagram illustrating old shard boundaries versus new split partitions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard split --source-id=<id> --target-a=<id1> --target-b=<id2> [--json]`
  - Flags: `--source-id` (string, source shard ID), `--target-a` (string, child 1), `--target-b` (string, child 2), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/shard_rebalancer.go`):
    ```go
    type ShardSplitOperation struct {
        SourceShardId      string `json:"sourceShardId"`
        TargetShardA       string `json:"targetShardA"`
        TargetShardB       string `json:"targetShardB"`
        MigratedRowCount   int64  `json:"migratedRowCount"`
        IsCutoverComplete  bool   `json:"isCutoverComplete"`
        HasDualWriteActive bool   `json:"hasDualWriteActive"`
        CompletedAtUnix    int64  `json:"completedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShardSplitLogs (
        splitId TEXT PRIMARY KEY,
        sourceShardId TEXT NOT NULL,
        targetShardA TEXT NOT NULL,
        targetShardB TEXT NOT NULL,
        migratedRowCount INTEGER NOT NULL,
        isCutoverComplete INTEGER NOT NULL,
        hasDualWriteActive INTEGER NOT NULL,
        completedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Emits migration events that update GitMap's Split-DB routing table without requiring process restarts.

---

### Improvement 587: SQLite B-Tree Page Corruption & Integrity Check Sentinel
- **Title:** SQLite B-Tree Page Corruption & Integrity Check Sentinel
- **Category:** Data Durability & Integrity Sentinel
- **Goal:** Periodically execute `PRAGMA integrity_check` and `PRAGMA quick_check` across all local shards, alerting on bitrot or damaged B-tree nodes.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Hex-style integrity verification summary with Green `INTEGRITY_OK` badge or Bright Red `PAGE_CORRUPTION` alert.
  - *Web UI:* Database health dashboard displaying B-tree depth, page allocation tables, and free list status.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard integrity --shard-id=<id> [--full] [--json]`
  - Flags: `--shard-id` (string, target shard), `--full` (bool, execute comprehensive integrity_check vs quick_check), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/integrity_sentinel.go`):
    ```go
    type ShardIntegrityReport struct {
        ShardIdentifier    string   `json:"shardIdentifier"`
        PageCount          int64    `json:"pageCount"`
        CorruptedPageCount int      `json:"corruptedPageCount"`
        ErrorDetails       []string `json:"errorDetails"`
        IsPristine         bool     `json:"isPristine"`
        HasPageBitrot      bool     `json:"hasPageBitrot"`
        VerifiedAtUnix     int64    `json:"verifiedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShardIntegrityReports (
        reportId TEXT PRIMARY KEY,
        shardIdentifier TEXT NOT NULL,
        pageCount INTEGER NOT NULL,
        corruptedPageCount INTEGER NOT NULL,
        errorDetailsJson TEXT NOT NULL,
        isPristine INTEGER NOT NULL,
        hasPageBitrot INTEGER NOT NULL,
        verifiedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Automatically isolates corrupted shard from read pools and triggers automated restoration from hot standby WAL replicas.

---

### Improvement 588: Connection Pool Starvation & Mutex Contention Profiler
- **Title:** Connection Pool Starvation & Mutex Contention Profiler
- **Category:** Concurrency Diagnostics & Pool Telemetry
- **Goal:** Profile Go `sql.DB` connection pool statistics (`MaxOpenConnections`, `InUse`, `WaitCount`, `WaitDuration`) to detect thread lock contention.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Dual bar gauge showing active pool utilization percentage and average connection acquisition latency.
  - *Web UI:* Microsecond latency histogram of DB handle acquisition with threshold warning markers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard pool-stats --shard-id=<id> [--json]`
  - Flags: `--shard-id` (string, shard ID), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/pool_profiler.go`):
    ```go
    type ConnectionPoolMetrics struct {
        ShardIdentifier    string `json:"shardIdentifier"`
        OpenConnections    int    `json:"openConnections"`
        InUseConnections   int    `json:"inUseConnections"`
        IdleConnections    int    `json:"idleConnections"`
        WaitCount          int64  `json:"waitCount"`
        WaitDurationMillis int64  `json:"waitDurationMillis"`
        IsStarved          bool   `json:"isStarved"`
        HasContention      bool   `json:"hasContention"`
        CapturedAtUnix     int64  `json:"capturedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ConnectionPoolLogs (
        logId TEXT PRIMARY KEY,
        shardIdentifier TEXT NOT NULL,
        openConnections INTEGER NOT NULL,
        inUseConnections INTEGER NOT NULL,
        idleConnections INTEGER NOT NULL,
        waitCount INTEGER NOT NULL,
        waitDurationMillis INTEGER NOT NULL,
        isStarved INTEGER NOT NULL,
        hasContention INTEGER NOT NULL,
        capturedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Automatically throttles concurrent warmup workers when mutex contention exceeds 15ms.

---

### Improvement 589: Zero-Downtime Backup & Snapshot Archival Coordinator
- **Title:** Zero-Downtime Backup & Snapshot Archival Coordinator
- **Category:** Disaster Recovery & Backup Orchestration
- **Goal:** Leverage the SQLite Online Backup API (`sqlite3_backup_init`) to stream point-in-time snapshots of sharded databases to local or remote storage.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Compact backup status row showing backup progress percentage, throughput in MB/s, and SHA256 checksum.
  - *Web UI:* Snapshot history table with retention timelines, size indicators, and one-click test-restore verification.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard backup --shard-id=<id> --out-dir=<dir> [--compress] [--json]`
  - Flags: `--shard-id` (string, shard identifier), `--out-dir` (string, destination directory), `--compress` (bool, enable zstd compression), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/backup_coordinator.go`):
    ```go
    type SnapshotBackupResult struct {
        ShardIdentifier    string `json:"shardIdentifier"`
        DestinationPath    string `json:"destinationPath"`
        SnapshotSizeBytes  int64  `json:"snapshotSizeBytes"`
        Sha256Checksum     string `json:"sha256Checksum"`
        DurationMillis     int64  `json:"durationMillis"`
        IsConsistent       bool   `json:"isConsistent"`
        HasCompleted       bool   `json:"hasCompleted"`
        FinishedAtUnix     int64  `json:"finishedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ShardSnapshots (
        snapshotId TEXT PRIMARY KEY,
        shardIdentifier TEXT NOT NULL,
        destinationPath TEXT NOT NULL,
        snapshotSizeBytes INTEGER NOT NULL,
        sha256Checksum TEXT NOT NULL,
        durationMillis INTEGER NOT NULL,
        isConsistent INTEGER NOT NULL,
        hasCompleted INTEGER NOT NULL,
        finishedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_snapshot_shard ON ShardSnapshots(shardIdentifier);
    ```
- **Antigravity / GitMap Integration Point:** Enforces cold-storage retention policies across GitMap repository backup hierarchies.

---

### Improvement 590: Sharded Cluster Live Management & Telemetry Cockpit
- **Title:** Sharded Cluster Live Management & Telemetry Cockpit
- **Category:** Distributed Storage Cockpit & Operator HUD
- **Goal:** Provide an administrative terminal and web console aggregating health, WAL replication status, and shard storage capacities across the fleet.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Full-screen Bubble Tea TUI with multi-shard table view, live IOPS counters, storage gauges, and keyboard navigation.
  - *Web UI:* Responsive dashboard with multi-node shard topology grid, live query throughput gauges, and zero external JS/CSS dependencies.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster shard cockpit [--watch] [--interval=2s] [--json]`
  - Flags: `--watch` (bool, continuous interactive mode), `--interval` (duration, polling frequency), `--json` (bool, output status snapshot as JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/cluster_cockpit.go`):
    ```go
    type ShardClusterStatus struct {
        TotalShardsCount   int     `json:"totalShardsCount"`
        HealthyShardsCount int     `json:"healthyShardsCount"`
        AggregateSizeBytes int64   `json:"aggregateSizeBytes"`
        IsClusterHealthy   bool    `json:"isClusterHealthy"`
        HasReplicationLag  bool    `json:"hasReplicationLag"`
        HasCorruptedNodes  bool    `json:"hasCorruptedNodes"`
        GeneratedAtUnix    int64   `json:"generatedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS ClusterCockpitSnapshots (
        snapshotId TEXT PRIMARY KEY,
        totalShardsCount INTEGER NOT NULL,
        healthyShardsCount INTEGER NOT NULL,
        aggregateSizeBytes INTEGER NOT NULL,
        isClusterHealthy INTEGER NOT NULL,
        hasReplicationLag INTEGER NOT NULL,
        hasCorruptedNodes INTEGER NOT NULL,
        generatedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Functions as the operational monitoring plane for GitMap distributed multi-node fleet deployments.

---

## 5. Cluster 60: Autonomous Antigravity Agent Fleet Choreographer & Telemetry Bridge (Items 591–600)

### Improvement 591: Agent Peer Discovery & Heartbeat Mesh Telemetry
- **Title:** Agent Peer Discovery & Heartbeat Mesh Telemetry
- **Category:** Multi-Agent Swarm & Peer Topology
- **Goal:** Autonomously track active Antigravity subagents, heartbeat liveness, assigned subtasks, and compute capacity across local and remote processes.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Live agent roster table with role names, active conversation IDs, heartbeat timestamps, and CPU/memory footprint badges.
  - *Web UI:* Real-time agent peer grid showing active agent cards, status pulses (Emerald active, Amber idle, Gray offline), and subtask allocations.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm roster [--filter=<role>] [--json]`
  - Flags: `--filter` (string, filter by agent role), `--json` (bool, JSON output format).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/peer_discovery.go`):
    ```go
    type AgentPeerInfo struct {
        ConversationId     string `json:"conversationId"`
        AgentRole          string `json:"agentRole"`
        ProcessId          int    `json:"processId"`
        IsHealthy          bool   `json:"isHealthy"`
        IsIdle             bool   `json:"isIdle"`
        HasAssignedSubtask bool   `json:"hasAssignedSubtask"`
        LastHeartbeatUnix  int64  `json:"lastHeartbeatUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AgentPeers (
        peerId TEXT PRIMARY KEY,
        conversationId TEXT NOT NULL UNIQUE,
        agentRole TEXT NOT NULL,
        processId INTEGER NOT NULL,
        isHealthy INTEGER NOT NULL,
        isIdle INTEGER NOT NULL,
        hasAssignedSubtask INTEGER NOT NULL,
        lastHeartbeatUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_agent_role ON AgentPeers(agentRole);
    ```
- **Antigravity / GitMap Integration Point:** Interacts with AGY subagent runtime (`define_subagent` / `invoke_subagent`) to dynamically rebalance work loads.

---

### Improvement 592: Subtask Lease Arbitration & Distributed Task Steerer
- **Title:** Subtask Lease Arbitration & Distributed Task Steerer
- **Category:** Multi-Agent Swarm & Task Scheduling
- **Goal:** Provide mutual exclusion and time-bounded task leases on pending plan items, preventing multiple agents from editing the same spec or code files.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Task allocation card showing subtask ID, assigned agent, remaining lease duration, and release key.
  - *Web UI:* Kanban-style task board with active agent avatar tokens anchored to locked task cards.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm lease acquire --task-id=<id> --agent-id=<cid> [--ttl=300s] [--json]`
  - Flags: `--task-id` (string, task slug), `--agent-id` (string, conversation ID), `--ttl` (duration, lease tenure), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/lease_arbitrator.go`):
    ```go
    type SubtaskLeaseRecord struct {
        TaskIdentifier     string `json:"taskIdentifier"`
        HolderAgentId      string `json:"holderAgentId"`
        LeaseExpiresAtUnix int64  `json:"leaseExpiresAtUnix"`
        IsLeaseActive      bool   `json:"isLeaseActive"`
        HasLockConflict    bool   `json:"hasLockConflict"`
        AcquiredAtUnix     int64  `json:"acquiredAtUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SubtaskLeases (
        leaseId TEXT PRIMARY KEY,
        taskIdentifier TEXT NOT NULL UNIQUE,
        holderAgentId TEXT NOT NULL,
        leaseExpiresAtUnix INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL,
        hasLockConflict INTEGER NOT NULL,
        acquiredAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Directly referenced by `.ai-memory/plans/` orchestration scripts to enforce file exclusivity.

---

### Improvement 593: Token Budget Consumption & LLM Cost Governance Telemetry
- **Title:** Token Budget Consumption & LLM Cost Governance Telemetry
- **Category:** Resource Governance & Model Cost Telemetry
- **Goal:** Real-time tracking of input tokens, output tokens, tool call invocations, and estimated API expenses per subagent conversation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Financial telemetry widget with token gauges, USD dollar estimations, and budget exhaustion warning bars.
  - *Web UI:* Stacked bar chart showing token consumption across agents over time with configurable budget limits.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm tokens [--agent-id=<cid>] [--budget=50.0] [--json]`
  - Flags: `--agent-id` (string, target agent ID), `--budget` (float64, spending cap in USD), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/token_governor.go`):
    ```go
    type AgentTokenTelemetry struct {
        ConversationId     string  `json:"conversationId"`
        InputTokens        int64   `json:"inputTokens"`
        OutputTokens       int64   `json:"outputTokens"`
        TotalEstimatedCost float64 `json:"totalEstimatedCost"`
        IsWithinBudget     bool    `json:"isWithinBudget"`
        HasExhaustedTokens bool    `json:"hasExhaustedTokens"`
        SampledAtUnix      int64   `json:"sampledAtUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS AgentTokenUsage (
        usageId TEXT PRIMARY KEY,
        conversationId TEXT NOT NULL,
        inputTokens INTEGER NOT NULL,
        outputTokens INTEGER NOT NULL,
        totalEstimatedCost REAL NOT NULL,
        isWithinBudget INTEGER NOT NULL,
        hasExhaustedTokens INTEGER NOT NULL,
        sampledAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_token_conv ON AgentTokenUsage(conversationId);
    ```
- **Antigravity / GitMap Integration Point:** Emits warning when subagent spending reaches 85% of budget, allowing graceful compaction or task wrap-up.

---

### Improvement 594: Inter-Agent Message Transit & Event Bus Telemetry
- **Title:** Inter-Agent Message Transit & Event Bus Telemetry
- **Category:** Swarm Communication & Telemetry Bus
- **Goal:** Capture, validate, and index all asynchronous inter-agent communications (`send_message`), tracking delivery latency and receipt acknowledgments.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Communication flow graph depicting messages passing between parent and subagents with millisecond latencies.
  - *Web UI:* Message transit timeline with expandable message bodies, payload schemas, and recipient status badges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm bus-log [--limit=20] [--follow] [--json]`
  - Flags: `--limit` (int, message history depth), `--follow` (bool, continuous tail), `--json` (bool, JSON stream).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/message_transit.go`):
    ```go
    type InterAgentMessageLog struct {
        MessageIdentifier  string `json:"messageIdentifier"`
        SenderAgentId      string `json:"senderAgentId"`
        RecipientAgentId   string `json:"recipientAgentId"`
        PayloadSizeChars   int    `json:"payloadSizeChars"`
        IsDelivered        bool   `json:"isDelivered"`
        HasPayloadErrors   bool   `json:"hasPayloadErrors"`
        TransitedAtUnix    int64  `json:"transitedAtUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InterAgentMessages (
        messageId TEXT PRIMARY KEY,
        messageIdentifier TEXT NOT NULL,
        senderAgentId TEXT NOT NULL,
        recipientAgentId TEXT NOT NULL,
        payloadSizeChars INTEGER NOT NULL,
        isDelivered INTEGER NOT NULL,
        hasPayloadErrors INTEGER NOT NULL,
        transitedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_msg_sender ON InterAgentMessages(senderAgentId);
    CREATE INDEX IF NOT EXISTS idx_msg_recipient ON InterAgentMessages(recipientAgentId);
    ```
- **Antigravity / GitMap Integration Point:** Plugs into the Antigravity conversation message hook to provide centralized observability without modifying agent prompts.

---

### Improvement 595: Autonomous Deadlock & Idle Swarm Detector
- **Title:** Autonomous Deadlock & Idle Swarm Detector
- **Category:** Swarm Health & Deadlock Remediation
- **Goal:** Detect cyclic task wait states, blocked tool invocations, or subagents that have remained idle with unreleased file locks.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Bright Amber deadlock alert box showing blocked subagent IDs, contending resource files, and auto-release countdown.
  - *Web UI:* Circular dependency graph highlighting contending resource nodes in pulsing Red with an immediate "Break Deadlock" action button.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm deadlock-check [--auto-recover] [--json]`
  - Flags: `--auto-recover` (bool, forcefully terminate expired locks), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/deadlock_detector.go`):
    ```go
    type SwarmDeadlockReport struct {
        ReportIdentifier   string   `json:"reportIdentifier"`
        DeadlockedAgentIds []string `json:"deadlockedAgentIds"`
        ContendedResources []string `json:"contendedResources"`
        IsDeadlocked       bool     `json:"isDeadlocked"`
        HasAutonomousFix   bool     `json:"hasAutonomousFix"`
        DetectedAtUnix     int64    `json:"detectedAtUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SwarmDeadlocks (
        reportId TEXT PRIMARY KEY,
        reportIdentifier TEXT NOT NULL,
        deadlockedAgentsJson TEXT NOT NULL,
        contendedResourcesJson TEXT NOT NULL,
        isDeadlocked INTEGER NOT NULL,
        hasAutonomousFix INTEGER NOT NULL,
        detectedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Invokes `manage_subagents(kill)` or cancels stuck leases when an irrecoverable cycle is confirmed.

---

### Improvement 596: Subagent Work Verification & Output Quality Sentinel
- **Title:** Subagent Work Verification & Output Quality Sentinel
- **Category:** Quality Assurance & Spec Compliance
- **Goal:** Mechanically audit subagent outputs against coding guidelines (positive booleans, function length <= 15 lines, zero nested ifs, relative paths) before parent merging.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Quality scorecard card displaying guideline check passes/fails, violating line numbers, and compliance percentages.
  - *Web UI:* Review docket visualizer highlighting code diffs side-by-side with rule compliance annotations.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm verify-output --path=<file-or-dir> [--strict] [--json]`
  - Flags: `--path` (string, target file or directory), `--strict` (bool, reject on minor warnings), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/output_verifier.go`):
    ```go
    type SubagentWorkAudit struct {
        TargetFilePath     string  `json:"targetFilePath"`
        ViolationsCount    int     `json:"violationsCount"`
        ComplianceScore    float64 `json:"complianceScore"`
        IsPassed           bool    `json:"isPassed"`
        HasNegativeBoolean bool    `json:"hasNegativeBoolean"`
        HasOversizedFunc   bool    `json:"hasOversizedFunc"`
        VerifiedAtUnix     int64   `json:"verifiedAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/audit_ledger.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS SubagentAudits (
        auditId TEXT PRIMARY KEY,
        targetFilePath TEXT NOT NULL,
        violationsCount INTEGER NOT NULL,
        complianceScore REAL NOT NULL,
        isPassed INTEGER NOT NULL,
        hasNegativeBoolean INTEGER NOT NULL,
        hasOversizedFunc INTEGER NOT NULL,
        verifiedAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_audit_file ON SubagentAudits(targetFilePath);
    ```
- **Antigravity / GitMap Integration Point:** Rejects subagent completion reports if code violates GitMap repository standards, sending correction instructions back to the worker.

---

### Improvement 597: Dynamic Subagent Respawn & Fleet Capacity Scaler
- **Title:** Dynamic Subagent Respawn & Fleet Capacity Scaler
- **Category:** Swarm Orchestration & Autoscaling
- **Goal:** Dynamically scale worker subagents between configured minimum and maximum limits (e.g. A=3, capacity 6) based on backlog task counts.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Capacity scaling banner displaying current worker pool size, pending queue length, and scale-up/scale-down actions.
  - *Web UI:* Autoscaling status dial with concurrent agent capacity slots and active task fill percentages.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm autoscale --min=2 --max=6 [--backlog-threshold=4] [--json]`
  - Flags: `--min` (int, minimum workers), `--max` (int, maximum worker capacity), `--backlog-threshold` (int, tasks per agent threshold), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/fleet_scaler.go`):
    ```go
    type FleetScalingDecision struct {
        ActiveAgentCount   int    `json:"activeAgentCount"`
        QueuedTaskCount    int    `json:"queuedTaskCount"`
        TargetAgentCount   int    `json:"targetAgentCount"`
        ScalingAction      string `json:"scalingAction"`
        IsScaleTriggered   bool   `json:"isScaleTriggered"`
        HasMaxCapacityHit  bool   `json:"hasMaxCapacityHit"`
        EvaluatedAtUnix    int64  `json:"evaluatedAtUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS FleetScalingLogs (
        logId TEXT PRIMARY KEY,
        activeAgentCount INTEGER NOT NULL,
        queuedTaskCount INTEGER NOT NULL,
        targetAgentCount INTEGER NOT NULL,
        scalingAction TEXT NOT NULL,
        isScaleTriggered INTEGER NOT NULL,
        hasMaxCapacityHit INTEGER NOT NULL,
        evaluatedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Triggers `invoke_subagent` when backlog increases, ensuring rapid completion of large plan sequences.

---

### Improvement 598: Conversation Transcript Indexer & AST Knowledge Miner
- **Title:** Conversation Transcript Indexer & AST Knowledge Miner
- **Category:** Swarm Memory & Knowledge Persistence
- **Goal:** Ingest and index JSONL agent conversation logs (`transcript.jsonl`) into local SQLite full-text search (FTS5) to enable semantic recall of past solutions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Search result snippet highlighting matched tool calls, error resolutions, and file paths with relevance rankings.
  - *Web UI:* Interactive transcript search console with tool invocation filters, timestamp navigation, and full markdown rendering.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm search-logs --query="<search-text>" [--role=<role>] [--json]`
  - Flags: `--query` (string, search term), `--role` (string, filter by subagent role), `--json` (bool, JSON format).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/transcript_miner.go`):
    ```go
    type TranscriptSearchResult struct {
        ConversationId     string `json:"conversationId"`
        StepIndex          int    `json:"stepIndex"`
        ToolName           string `json:"toolName"`
        MatchedSnippet     string `json:"matchedSnippet"`
        RelevanceScore     float64`json:"relevanceScore"`
        HasErrorContext    bool   `json:"hasErrorContext"`
        IndexedAtUnix      int64  `json:"indexedAtUnix"`
    }
    ```
  - SQLite Schema (`installation.db`):
    ```sql
    CREATE VIRTUAL TABLE IF NOT EXISTS TranscriptFts USING fts5(
        conversationId,
        stepIndex,
        toolName,
        content,
        thinking
    );
    ```
- **Antigravity / GitMap Integration Point:** Supplies rich context to `read-memory-enhanced` skill by querying past resolutions for identical error traces.

---

### Improvement 599: Unified Antigravity Telemetry Bridge & Event Stream Broadcaster
- **Title:** Unified Antigravity Telemetry Bridge & Event Stream Broadcaster
- **Category:** Telemetry Streaming & Event Distribution
- **Goal:** Broadcast real-time agent lifecycle events, tool calls, and warm-up milestones across local IPC sockets and zero-CDN websockets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Streaming event log viewer displaying timestamp, emitting agent, event type, and compact JSON payload summary.
  - *Web UI:* Live event water-fall monitor showing streaming microsecond events with interactive category toggles.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agy swarm stream-events [--types=mail,agent,db] [--json]`
  - Flags: `--types` (comma-separated event categories), `--json` (bool, output stream as newline-delimited JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/event_broadcaster.go`):
    ```go
    type SwarmTelemetryEvent struct {
        EventIdentifier    string `json:"eventIdentifier"`
        EventType          string `json:"eventType"`
        SourceAgentId      string `json:"sourceAgentId"`
        PayloadJson        string `json:"payloadJson"`
        IsPriorityEvent    bool   `json:"isPriorityEvent"`
        HasSubscriberAck   bool   `json:"hasSubscriberAck"`
        BroadcastAtUnix    int64  `json:"broadcastAtUnix"`
    }
    ```
  - SQLite Schema (`repodb/event_stream.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS TelemetryEvents (
        eventId TEXT PRIMARY KEY,
        eventIdentifier TEXT NOT NULL,
        eventType TEXT NOT NULL,
        sourceAgentId TEXT NOT NULL,
        payloadJson TEXT NOT NULL,
        isPriorityEvent INTEGER NOT NULL,
        hasSubscriberAck INTEGER NOT NULL,
        broadcastAtUnix INTEGER NOT NULL
    );
    CREATE INDEX IF NOT EXISTS idx_telemetry_type ON TelemetryEvents(eventType);
    ```
- **Antigravity / GitMap Integration Point:** Powers live dashboards and webhook listeners for enterprise automation pipelines.

---

### Improvement 600: Master Swarm Orchestrator Cockpit & Fleet Command Deck
- **Title:** Master Swarm Orchestrator Cockpit & Fleet Command Deck
- **Category:** Unified Operator Cockpit & Executive HUD
- **Goal:** Unify all 100 improvements (501–600) into a single high-performance terminal and web mission control cockpit for warm-up and emailing operations.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-pane interactive terminal command deck integrating DNS sentinels, bounce forensics, sharded databases, and agent swarms into an executive HUD.
  - *Web UI:* Full-screen responsive operations dashboard featuring live health gauges, swarm controls, shard maps, and audio alert triggers with zero external CDNs.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup cockpit [--port=8080] [--open-browser] [--json]`
  - Flags: `--port` (int, local web dashboard port), `--open-browser` (bool, launch local browser window), `--json` (bool, output full system state snapshot).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/master_cockpit.go`):
    ```go
    type MasterSwarmCockpitState struct {
        ActiveCampaignCount int     `json:"activeCampaignCount"`
        TotalNodesOnline    int     `json:"totalNodesOnline"`
        OverallSystemHealth float64 `json:"overallSystemHealth"`
        IsFleetOperational  bool    `json:"isFleetOperational"`
        HasActiveAlarms     bool    `json:"hasActiveAlarms"`
        HasPendingTasks     bool    `json:"hasPendingTasks"`
        CapturedAtUnix      int64   `json:"capturedAtUnix"`
    }
    ```
  - SQLite Schema (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS MasterCockpitStates (
        stateId TEXT PRIMARY KEY,
        activeCampaignCount INTEGER NOT NULL,
        totalNodesOnline INTEGER NOT NULL,
        overallSystemHealth REAL NOT NULL,
        isFleetOperational INTEGER NOT NULL,
        hasActiveAlarms INTEGER NOT NULL,
        hasPendingTasks INTEGER NOT NULL,
        capturedAtUnix INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Acts as the apex command center for the entire GitMap Warm-Up and Emailing UI Engine.

---

## 6. Architectural Invariants & Quality Standards

1. **Relative Path Hygiene:**
   - 100% of all path references in documentation, code imports, and configurations must use relative repository paths (`02-spec/...`, `pkg/...`, `.ai-memory/...`). Absolute paths with drive letters (`D:\work\...`, `C:\Users\...`) are strictly prohibited.
2. **Boolean Naming Conventions:**
   - Every boolean field and parameter must strictly begin with `is` or `has` prefixes exclusively (`isReachable`, `hasValidMxRecords`, `isPermanentFailure`).
   - No negative booleans (`isNotValid`, `hasNoErrors`) or prohibited prefixes (`can`, `should`, `was`, `will`).
3. **Split-DB Architecture Adherence:**
   - Global cluster topology, leases, and cockpit states persist in `gitmap.db`.
   - Node-specific peer metadata, token usage, and subtask leases persist in `installation.db`.
   - DNS checks, bounce forensics, and telemetry streams persist in dedicated domain databases (`repodb/mail_dns.db`, `repodb/mail_bounce.db`, `repodb/event_stream.db`).
4. **Error Handling & Return Sovereignty:**
   - Leaf functions return `*apperror.AppError` directly.
   - Zero calls to `os.Exit` or CLI exit handlers inside leaf libraries.
5. **Function & File Size Guardrails:**
   - Functions <= 15 lines max (<= 8 lines preferred).
   - Zero nested `if` statements. Conditionals must be flattened using guard clauses and early returns.
