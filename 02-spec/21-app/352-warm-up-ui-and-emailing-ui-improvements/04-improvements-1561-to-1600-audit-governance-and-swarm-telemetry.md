# Specification 352: Improvements 1561 to 1600 - Audit Governance, DNS Repair, Bounce Forensics & Swarm Telemetry

- **Sequence:** 352
- **Slug:** `352-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/04-improvements-1561-to-1600-audit-governance-and-swarm-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/02-improvements-1501-to-1530-session-pipelining-and-zero-window.md`
  - `02-spec/21-app/352-warm-up-ui-and-emailing-ui-improvements/03-improvements-1531-to-1560-interactive-tui-and-mailbox-diagnostics.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification defines Improvements 1561 through 1600 for Sequence 352 of the GitMap Warm-Up and Emailing UI Engine. Sequence 352 hardens DNS over HTTPS (DoH) consensus self-healing, RFC 3463/5248 diagnostic bounce telemetry, distributed SQLite Split-DB multi-master replication, and autonomous multi-agent swarm orchestration across four major clusters:

1. **Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (Items 1561–1570):** Anycast DoH multi-resolver consensus auditing, DNSSEC NSEC3 authenticated denial-of-existence proof validation, RFC 7208 syntax validation with macro expansion, recursive SPF lookup tree flattening (<10 queries) with CIDR collapsing, dual-selector cryptographic DKIM key rotation (RSA-2048 & Ed25519), DMARC policy alignment synthesizers with stepped enforcement (`p=none` -> `p=quarantine` -> `p=reject`), multi-provider API provisioning (Cloudflare, Route53, DigitalOcean), global anycast propagation radars, and automated DNS drift rollback circuit breakers.
2. **Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 1571–1580):** High-precision RFC 3463/5248 Enhanced Mail System Status Code decoding, multipart report MIME parsers with DSN header extraction, `Authentication-Results` & `X-Spam-Status` header telemetry, offline Bayesian trigger phrase probabilistic scoring, interactive Bubble Tea terminal bounce forensics HUD, automated mailbox quarantine circuit breaking, SIEM webhook dispatchers (Splunk, Elastic, Datadog), in-memory Bloom filter suppression lists, spam trap honeypot behavioral sentinels, and longitudinal predictive bounce trend modeling.
3. **Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD (Items 1581–1590):** Vector-clock CRDT state synchronization across distributed worker replicas, encrypted WAL frame streaming transport over mutual TLS 1.3, page-level Merkle tree shard integrity auditors, distributed two-phase epoch commit coordinators, non-blocking hot VACUUM compaction engines, real-time terminal replication cluster HUDs, deterministic Last-Write-Wins (LWW) conflict arbiters with cryptographic tie-breaking, dynamic hash-partition shard rebalancers, point-in-time recovery (PITR) WAL log replayers, and zero-downtime rolling schema migration orchestrators.
4. **Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (Items 1591–1600):** Real-time multi-agent swarm cockpit HUDs, dynamic subtask auctioneers with expiring leases, token budget expenditure governors, Tarjan-based circular deadlock detection and worker self-healing, sub-250ms embedded Raft consensus master electors, closed-loop deliverability feedback governors with autonomous throttle steering, multi-turn synthetic contextual conversation generators, per-ISP reputation vector heatmaps (Google, Microsoft, Yahoo, Apple), autonomous domain warmup graduation gates, and full-screen unified mission telemetry dashboards.

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 352 Improvements 1561 to 1600 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 157 (1561-1570)    │ Cluster 158 (1571-1580)    │ Cluster 159 (1581-1590)      │
│ Autonomous SPF/DKIM/DMARC  │ Forensic Bounce Analysis   │ Distributed SQLite Split-DB  │
│ Record Repair & DoH Sentinel│ Engine & Classification HUD│ Multi-Master Replication HUD │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 160 (1591-1600)                                 │ Cross-Cutting Standards      │
│ Antigravity Multi-Agent Swarm Orchestrator              │ Positive Booleans (is/has),  │
│ & Autonomous Warmup Steerer                             │ Split-DB, Zero CDN, No Exits │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (Items 1561–1570)

### Improvement 1561: Recursive DoH Multi-Resolver Consensus Sentinel & Anycast Quorum Validator
- **Title:** Recursive DoH Multi-Resolver Consensus Sentinel & Anycast Quorum Validator
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Concurrently query authoritative DNS records across multiple independent DoH endpoints (Cloudflare, Google, Quad9, AdGuard, OpenDNS) with recursive CNAME unrolling and cryptographic quorum verification to preempt DNS cache poisoning and ISP interception.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-density resolver matrix panel rendering per-resolver latency (ms), HTTP/3 status, unrolled CNAME hops, and RDATA payload consensus state in Emerald (`#10B981`) for consensus or Crimson (`#EF4444`) for divergence.
  - *Web UI:* Zero-CDN CSS grid comparing resolver response trees side-by-side with diff highlighting.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns doh-sentinel --domain=<domain> [--record-type=TXT] [--min-quorum=3] [--json]`
  - Flags: `--domain` (string, apex domain to inspect), `--record-type` (string, target RR type: TXT, MX, A, AAAA), `--min-quorum` (int, minimum consensus matching threshold, default 3), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DohSentinelConsensusAudit struct {
        ApexDomain       string   `json:"apexDomain"`
        RecordType       string   `json:"recordType"`
        ActiveResolvers  []string `json:"activeResolvers"`
        ConsensusRecords []string `json:"consensusRecords"`
        AverageLatencyMs int64    `json:"averageLatencyMs"`
        EvaluatedAt      int64    `json:"evaluatedAt"`
        isConsensusMet   bool     `json:"isConsensusMet"`
        hasCnameLoop     bool     `json:"hasCnameLoop"`
        isDnssecVerified bool     `json:"isDnssecVerified"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDohSentinelAudit (
        ApexDomain TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        ActiveResolvers TEXT NOT NULL,
        ConsensusRecords TEXT NOT NULL,
        AverageLatencyMs INTEGER NOT NULL,
        EvaluatedAt INTEGER NOT NULL,
        isConsensusMet INTEGER NOT NULL CHECK (isConsensusMet IN (0, 1)),
        hasCnameLoop INTEGER NOT NULL CHECK (hasCnameLoop IN (0, 1)),
        isDnssecVerified INTEGER NOT NULL CHECK (isDnssecVerified IN (0, 1)),
        PRIMARY KEY (ApexDomain, RecordType, EvaluatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Hooked into `pkg/mail/dnscheck/sentinel.go` as a mandatory pre-flight validation gate prior to starting warm-up batch dispatch.

---

### Improvement 1562: DNSSEC Cryptographic Trust Chain Traverser & RRSIG Expiration Monitor
- **Title:** DNSSEC Cryptographic Trust Chain Traverser & RRSIG Expiration Monitor
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Traverse and cryptographically verify the complete DNSSEC chain of trust from root KSK down through TLD DS records and zone ZSK/RRSIG signatures, alerting on impending signature expiration or broken crypto anchors.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Visual hierarchical tree representation of DNSSEC trust nodes with lock icons (`[VALID]`), key tag identifiers, digest types (SHA-256), and time-to-expiration progress bars color-coded by remaining days.
  - *Web UI:* SVG trust-chain ladder diagram displaying cryptographic links and expiration countdown gauges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns dnssec-verify --domain=<domain> [--warn-days=7] [--json]`
  - Flags: `--domain` (string, domain to verify), `--warn-days` (int, days threshold before RRSIG expiry triggers warning, default 7), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnssecTrustChainReport struct {
        DomainName         string `json:"domainName"`
        RootTrustAnchor    string `json:"rootTrustAnchor"`
        ZoneKeyTag         int    `json:"zoneKeyTag"`
        AlgorithmType      string `json:"algorithmType"`
        RrsigExpiresAt     int64  `json:"rrsigExpiresAt"`
        DaysUntilExpiry    int    `json:"daysUntilExpiry"`
        isChainValid       bool   `json:"isChainValid"`
        hasBogusProof      bool   `json:"hasBogusProof"`
        isExpiringSoon     bool   `json:"isExpiringSoon"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnssecAudit (
        DomainName TEXT NOT NULL,
        ZoneKeyTag INTEGER NOT NULL,
        AlgorithmType TEXT NOT NULL,
        RrsigExpiresAt INTEGER NOT NULL,
        DaysUntilExpiry INTEGER NOT NULL,
        isChainValid INTEGER NOT NULL CHECK (isChainValid IN (0, 1)),
        hasBogusProof INTEGER NOT NULL CHECK (hasBogusProof IN (0, 1)),
        isExpiringSoon INTEGER NOT NULL CHECK (isExpiringSoon IN (0, 1)),
        RecordedAt INTEGER NOT NULL,
        PRIMARY KEY (DomainName, ZoneKeyTag, RecordedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Invoked in `pkg/mail/dnscheck/dnssec.go` before initiating cryptographic DKIM / DMARC verification cycles.

---

### Improvement 1563: SPF Syntax Deep-Parser, Macro Expansion Engine & Multiple-Record Detector
- **Title:** SPF Syntax Deep-Parser, Macro Expansion Engine & Multiple-Record Detector
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Perform deep RFC 7208 syntax parsing on domain SPF records, evaluating macro expansion tokens (`%{s}`, `%{l}`, `%{o}`), detecting multiple duplicate TXT `v=spf1` records (which cause RFC fatal `permerror`), and validating mechanism order.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Tokenized syntax highlighter displaying each mechanism (`+ip4`, `include`, `redirect`, `~all`) with syntax status badges and inline warning annotations for duplicate records.
  - *Web UI:* Interactive SPF token visualizer with simulated sender IP validation simulator.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns spf-audit --domain=<domain> [--sender-ip=<ip>] [--json]`
  - Flags: `--domain` (string, domain to inspect), `--sender-ip` (string, optional test IP to evaluate match against parsed SPF mechanisms), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpfSyntaxAuditResult struct {
        DomainName        string   `json:"domainName"`
        RawRecords        []string `json:"rawRecords"`
        ExtractedTokens   []string `json:"extractedTokens"`
        EvaluationResult  string   `json:"evaluationResult"`
        RecordCount       int      `json:"recordCount"`
        EvaluatedAt       int64    `json:"evaluatedAt"`
        isValidSyntax     bool     `json:"isValidSyntax"`
        hasMultipleSpf    bool     `json:"hasMultipleSpf"`
        hasDeprecatedPtr  bool     `json:"hasDeprecatedPtr"`
        hasAllQualifier   bool     `json:"hasAllQualifier"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpfSyntaxAudit (
        DomainName TEXT NOT NULL,
        RawRecords TEXT NOT NULL,
        ExtractedTokens TEXT NOT NULL,
        EvaluationResult TEXT NOT NULL,
        RecordCount INTEGER NOT NULL,
        EvaluatedAt INTEGER NOT NULL,
        isValidSyntax INTEGER NOT NULL CHECK (isValidSyntax IN (0, 1)),
        hasMultipleSpf INTEGER NOT NULL CHECK (hasMultipleSpf IN (0, 1)),
        hasDeprecatedPtr INTEGER NOT NULL CHECK (hasDeprecatedPtr IN (0, 1)),
        hasAllQualifier INTEGER NOT NULL CHECK (hasAllQualifier IN (0, 1)),
        PRIMARY KEY (DomainName, EvaluatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Core module in `pkg/mail/dnsrepair/spf_parser.go` consumed by pre-dispatch warming health engines.

---

### Improvement 1564: Recursive SPF Lookup Tree Flattening (<10 Lookups) & CIDR Collapsing Engine
- **Title:** Recursive SPF Lookup Tree Flattening (<10 Lookups) & CIDR Collapsing Engine
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Traverse recursive SPF `include:`, `a:`, `mx:`, and `redirect:` directives, count DNS query consumption, and autonomously synthesize flattened `ip4:` and `ip6:` CIDR blocks collapsed to strictly stay under the RFC 7208 10-lookup limit.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive tree visualization showing lookup depth per branch, current total count with a colored gauge (Green: <=7, Amber: 8-9, Red: >=10), and the proposed flattened record with diff display.
  - *Web UI:* Expandable DAG graph displaying the full external DNS resolution tree and generated CIDR aggregates.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns spf-flatten --domain=<domain> [--max-cidrs=50] [--dry-run] [--json]`
  - Flags: `--domain` (string, domain whose SPF to flatten), `--max-cidrs` (int, max CIDR entries per flattened TXT block, default 50), `--dry-run` (bool, preview synthesis without writing changes), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpfFlatteningPlan struct {
        ApexDomain          string   `json:"apexDomain"`
        OriginalLookups     int      `json:"originalLookups"`
        FlattenedLookups    int      `json:"flattenedLookups"`
        ExtractedIp4Cidrs   []string `json:"extractedIp4Cidrs"`
        ExtractedIp6Cidrs   []string `json:"extractedIp6Cidrs"`
        SynthesizedRecord   string   `json:"synthesizedRecord"`
        GeneratedAt         int64    `json:"generatedAt"`
        isUnderLimit        bool     `json:"isUnderLimit"`
        hasCidrCollapsed    bool     `json:"hasCidrCollapsed"`
        isReadyForPublish   bool     `json:"isReadyForPublish"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpfFlattenPlan (
        ApexDomain TEXT NOT NULL,
        OriginalLookups INTEGER NOT NULL,
        FlattenedLookups INTEGER NOT NULL,
        ExtractedIp4Cidrs TEXT NOT NULL,
        ExtractedIp6Cidrs TEXT NOT NULL,
        SynthesizedRecord TEXT NOT NULL,
        GeneratedAt INTEGER NOT NULL,
        isUnderLimit INTEGER NOT NULL CHECK (isUnderLimit IN (0, 1)),
        hasCidrCollapsed INTEGER NOT NULL CHECK (hasCidrCollapsed IN (0, 1)),
        isReadyForPublish INTEGER NOT NULL CHECK (isReadyForPublish IN (0, 1)),
        PRIMARY KEY (ApexDomain, GeneratedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/dnsrepair/flattener.go` and directly executable as part of `gitmap warmup repair-spf`.

---

### Improvement 1565: Dual-Selector Cryptographic DKIM Key Rotator (RSA-2048 & Ed25519) & DNS Publisher
- **Title:** Dual-Selector Cryptographic DKIM Key Rotator (RSA-2048 & Ed25519) & DNS Publisher
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Generate and rotate modern dual-selector DKIM keys (RSA-2048 for legacy MTA compatibility and RFC 8463 Ed25519 for high-speed elliptic curve verification) and manage zero-downtime overlapping key transitions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Split-view key manager card displaying Selector A (Active) and Selector B (Staged/Retiring), cryptographic fingerprint (SHA-256), public key preview, and retirement countdown timer.
  - *Web UI:* One-click key pair generation dashboard with instant public key DNS copy triggers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dkim rotate --domain=<domain> [--selector-prefix=gm] [--algo=dual] [--stage-only]`
  - Flags: `--domain` (string, sending domain), `--selector-prefix` (string, selector naming prefix, default "gm"), `--algo` (string, "rsa", "ed25519", "dual", default "dual"), `--stage-only` (bool, generate without activating).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DkimKeyRotationState struct {
        DomainName      string `json:"domainName"`
        ActiveSelector  string `json:"activeSelector"`
        StagedSelector  string `json:"stagedSelector"`
        AlgorithmType   string `json:"algorithmType"`
        PublicKeyDns    string `json:"publicKeyDns"`
        PrivateKeyPath  string `json:"privateKeyPath"`
        RotatedAt       int64  `json:"rotatedAt"`
        RetireAfter     int64  `json:"retireAfter"`
        isActive        bool   `json:"isActive"`
        isEd25519Dual   bool   `json:"isEd25519Dual"`
        hasDnsPublished bool   `json:"hasDnsPublished"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationDkimKeys (
        DomainName TEXT NOT NULL,
        ActiveSelector TEXT NOT NULL,
        StagedSelector TEXT NOT NULL,
        AlgorithmType TEXT NOT NULL,
        PublicKeyDns TEXT NOT NULL,
        PrivateKeyPath TEXT NOT NULL,
        RotatedAt INTEGER NOT NULL,
        RetireAfter INTEGER NOT NULL,
        isActive INTEGER NOT NULL CHECK (isActive IN (0, 1)),
        isEd25519Dual INTEGER NOT NULL CHECK (isEd25519Dual IN (0, 1)),
        hasDnsPublished INTEGER NOT NULL CHECK (hasDnsPublished IN (0, 1)),
        PRIMARY KEY (DomainName, ActiveSelector)
    );
    ```
- **Antigravity / GitMap Integration Point:** Managed by `pkg/mail/crypto/dkim_rotator.go` with secure local private key persistence in `installation.db`.

---

### Improvement 1566: DMARC Policy Alignment Synthesizer & Stepped Enforcement Scheduler
- **Title:** DMARC Policy Alignment Synthesizer & Stepped Enforcement Scheduler
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Analyze domain authentication telemetry, verify SPF/DKIM identifier alignment (strict vs relaxed), and generate automated stepped enforcement schedules transitioning domains progressively from `p=none` to `p=quarantine; pct=25` up to `p=reject; pct=100`.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Stepped progression ladder visualization indicating current enforcement milestone, pass/fail alignment ratios, RUA/RUF destination endpoints, and next phase trigger date.
  - *Web UI:* Alignment audit matrix showing From header vs Return-Path and DKIM d= alignment scores.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dmarc synthesize --domain=<domain> [--rua-email=<rua>] [--target-policy=reject] [--step-days=14]`
  - Flags: `--domain` (string, domain to configure), `--rua-email` (string, aggregate reporting email), `--target-policy` (string, final desired policy: "none", "quarantine", "reject"), `--step-days` (int, days between progression steps, default 14).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DmarcPolicyProgression struct {
        DomainName          string `json:"domainName"`
        CurrentPolicy       string `json:"currentPolicy"`
        TargetPolicy        string `json:"targetPolicy"`
        CurrentPercentage   int    `json:"currentPercentage"`
        RuaDestination      string `json:"ruaDestination"`
        RufDestination      string `json:"rufDestination"`
        PhaseStartedAt      int64  `json:"phaseStartedAt"`
        NextPhaseAt         int64  `json:"nextPhaseAt"`
        isStrictAlignment   bool   `json:"isStrictAlignment"`
        hasActiveViolations bool   `json:"hasActiveViolations"`
        isGraduatedToReject bool   `json:"isGraduatedToReject"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDmarcProgression (
        DomainName TEXT NOT NULL,
        CurrentPolicy TEXT NOT NULL,
        TargetPolicy TEXT NOT NULL,
        CurrentPercentage INTEGER NOT NULL,
        RuaDestination TEXT NOT NULL,
        RufDestination TEXT NOT NULL,
        PhaseStartedAt INTEGER NOT NULL,
        NextPhaseAt INTEGER NOT NULL,
        isStrictAlignment INTEGER NOT NULL CHECK (isStrictAlignment IN (0, 1)),
        hasActiveViolations INTEGER NOT NULL CHECK (hasActiveViolations IN (0, 1)),
        isGraduatedToReject INTEGER NOT NULL CHECK (isGraduatedToReject IN (0, 1)),
        PRIMARY KEY (DomainName, PhaseStartedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Orchestrated in `pkg/mail/dnsrepair/dmarc_synthesizer.go` and surfaced on the CLI warmup status console.

---

### Improvement 1567: Multi-Provider DNS API Provisioner (Cloudflare, Route53, DigitalOcean) with Rollback
- **Title:** Multi-Provider DNS API Provisioner (Cloudflare, Route53, DigitalOcean) with Rollback
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Interface directly with major authoritative DNS provider APIs to automatically publish synthesized SPF, DKIM, and DMARC records with dry-run verification and automatic snapshot rollback upon detection of transient outages.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive confirmation modal with side-by-side RDATA diff (Original vs New), provider API ping latency meter, and transactional commit progress bar.
  - *Web UI:* Provider connection status cards with API token health verification badges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns apply --domain=<domain> --provider=<cloudflare|route53|digitalocean> [--commit] [--dry-run]`
  - Flags: `--domain` (string, zone domain), `--provider` (string, DNS provider adapter), `--commit` (bool, execute real remote updates), `--dry-run` (bool, simulate without API calls, default true).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsProvisionTransaction struct {
        TransactionId     string `json:"transactionId"`
        DomainName        string `json:"domainName"`
        ProviderType      string `json:"providerType"`
        PreviousRecord    string `json:"previousRecord"`
        PublishedRecord   string `json:"publishedRecord"`
        ExecutedAt        int64  `json:"executedAt"`
        isCommitted       bool   `json:"isCommitted"`
        hasRollbackOccurred bool `json:"hasRollbackOccurred"`
        isDryRun          bool   `json:"isDryRun"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationDnsTransactions (
        TransactionId TEXT PRIMARY KEY,
        DomainName TEXT NOT NULL,
        ProviderType TEXT NOT NULL,
        PreviousRecord TEXT NOT NULL,
        PublishedRecord TEXT NOT NULL,
        ExecutedAt INTEGER NOT NULL,
        isCommitted INTEGER NOT NULL CHECK (isCommitted IN (0, 1)),
        hasRollbackOccurred INTEGER NOT NULL CHECK (hasRollbackOccurred IN (0, 1)),
        isDryRun INTEGER NOT NULL CHECK (isDryRun IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Adapter layer in `pkg/mail/dnsrepair/providers/` implementing unified `DnsProviderClient` interface.

---

### Improvement 1568: Global Anycast Propagation Radar & Bounded TTL Polling Sentinel
- **Title:** Global Anycast Propagation Radar & Bounded TTL Polling Sentinel
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Continuously sample DNS resolution across 24 global edge points (North America, Europe, Asia-Pacific, Latin America) following a record update, rendering real-time propagation percentage and pausing warming until worldwide propagation completes.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* ASCII world map highlighting edge nodes in Green (Updated) or Yellow (Pending/Stale), accompanied by a global propagation percentage bar (`[████████░░] 83%`).
  - *Web UI:* Interactive geo-coordinate map pinpointing each probe location and local cache TTL remaining.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns propagation --domain=<domain> [--record-type=TXT] [--target-propagation=95] [--watch]`
  - Flags: `--domain` (string, domain to monitor), `--record-type` (string, DNS record type), `--target-propagation` (int, target propagation percentage threshold, default 95), `--watch` (bool, stream live terminal radar).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type AnycastPropagationStatus struct {
        DomainName         string  `json:"domainName"`
        RecordType         string  `json:"recordType"`
        TotalProbes        int     `json:"totalProbes"`
        PropagatedProbes   int     `json:"propagatedProbes"`
        PropagationPercent float64 `json:"propagationPercent"`
        MinTtlRemainingSec int64   `json:"minTtlRemainingSec"`
        SampledAt          int64   `json:"sampledAt"`
        isTargetMet        bool    `json:"isTargetMet"`
        hasStaleCaches     bool    `json:"hasStaleCaches"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsPropagation (
        DomainName TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        TotalProbes INTEGER NOT NULL,
        PropagatedProbes INTEGER NOT NULL,
        PropagationPercent REAL NOT NULL,
        MinTtlRemainingSec INTEGER NOT NULL,
        SampledAt INTEGER NOT NULL,
        isTargetMet INTEGER NOT NULL CHECK (isTargetMet IN (0, 1)),
        hasStaleCaches INTEGER NOT NULL CHECK (hasStaleCaches IN (0, 1)),
        PRIMARY KEY (DomainName, RecordType, SampledAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Runs in `pkg/mail/dnscheck/propagation_radar.go` as a non-blocking background watcher.

---

### Improvement 1569: Continuous DNS Drift Detection Engine & Automated Rollback Circuit Breaker
- **Title:** Continuous DNS Drift Detection Engine & Automated Rollback Circuit Breaker
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Periodically scan active domain DNS configuration to detect unexpected upstream changes, unauthorized record overwrites, or registrar DNS outages, tripping an immediate warm-up circuit breaker to protect domain reputation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Flashing Amber alert banner when drift is detected, showing expected vs observed RDATA diff and circuit breaker trip state.
  - *Web UI:* Audit timeline with historical snapshots and side-by-side drift visualizer.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns drift-check --domain=<domain> [--auto-trip-breaker] [--json]`
  - Flags: `--domain` (string, domain to inspect), `--auto-trip-breaker` (bool, automatically pause warming queue on drift, default true), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsDriftEventRecord struct {
        DomainName       string `json:"domainName"`
        RecordType       string `json:"recordType"`
        ExpectedDigest   string `json:"expectedDigest"`
        ObservedDigest   string `json:"observedDigest"`
        ObservedRdata    string `json:"observedRdata"`
        DetectedAt       int64  `json:"detectedAt"`
        isDriftDetected  bool   `json:"isDriftDetected"`
        hasBreakerTripped bool  `json:"hasBreakerTripped"`
        isAutoRestored   bool   `json:"isAutoRestored"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsDriftEvents (
        DomainName TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        ExpectedDigest TEXT NOT NULL,
        ObservedDigest TEXT NOT NULL,
        ObservedRdata TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        isDriftDetected INTEGER NOT NULL CHECK (isDriftDetected IN (0, 1)),
        hasBreakerTripped INTEGER NOT NULL CHECK (hasBreakerTripped IN (0, 1)),
        isAutoRestored INTEGER NOT NULL CHECK (isAutoRestored IN (0, 1)),
        PRIMARY KEY (DomainName, RecordType, DetectedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded sentinel in `pkg/mail/dnscheck/drift_sentinel.go` feeding warming throttle state.

---

### Improvement 1570: DNS Over HTTPS Health Dashboard & TUI Configuration Inspector
- **Title:** DNS Over HTTPS Health Dashboard & TUI Configuration Inspector
- **Cluster:** Cluster 157: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Deliver a unified interactive terminal dashboard summarizing SPF lookup consumption, DKIM selector states, DMARC enforcement phase, and upstream DoH health metrics across all configured warming domains.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Comprehensive multi-column table displaying Domain, SPF Lookups (e.g., `4/10`), DKIM Algo (`Ed25519+RSA`), DMARC Policy (`p=quarantine 50%`), DoH Quorum (`5/5`), and Health Badge (`[PRIME]`).
  - *Web UI:* Responsive HTML/SVG scorecard viewable without external fonts or scripts.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns status [--filter-domain=<domain>] [--export-markdown=<path>]`
  - Flags: `--filter-domain` (string, filter specific domain), `--export-markdown` (string, file path to render markdown report).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsOverallHealthSummary struct {
        DomainName          string `json:"domainName"`
        SpfLookupCount      int    `json:"spfLookupCount"`
        DkimSelectorState   string `json:"dkimSelectorState"`
        DmarcPolicySummary  string `json:"dmarcPolicySummary"`
        ConsensusScore      int    `json:"consensusScore"`
        OverallGrade        string `json:"overallGrade"`
        isReadyForWarmup    bool   `json:"isReadyForWarmup"`
        hasActionRequired   bool   `json:"hasActionRequired"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsOverallHealth (
        DomainName TEXT PRIMARY KEY,
        SpfLookupCount INTEGER NOT NULL,
        DkimSelectorState TEXT NOT NULL,
        DmarcPolicySummary TEXT NOT NULL,
        ConsensusScore INTEGER NOT NULL,
        OverallGrade TEXT NOT NULL,
        isReadyForWarmup INTEGER NOT NULL CHECK (isReadyForWarmup IN (0, 1)),
        hasActionRequired INTEGER NOT NULL CHECK (hasActionRequired IN (0, 1)),
        UpdatedAt INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** CLI presentation layer in `pkg/cli/tui/dns_dashboard.go`.

---

## 3. Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 1571–1580)

### Improvement 1571: RFC 3463/5248 Enhanced Status Code Deep-Decoder & Root-Cause Classifier
- **Title:** RFC 3463/5248 Enhanced Status Code Deep-Decoder & Root-Cause Classifier
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Ingest raw SMTP error strings and Non-Delivery Reports (NDRs), extract standard 3-digit and enhanced 5-digit error codes (e.g., `5.7.1`, `4.4.2`, `5.1.1`), and classify them into deterministic operational root causes.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Colorized error decomposition view displaying the raw message, matched RFC definition token, categorized cause (e.g., `[REPUTATION_BLOCK]`, `[MAILBOX_DOES_NOT_EXIST]`), and suggested remediation action.
  - *Web UI:* Interactive status code lookup table with historical occurrence frequencies.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce decode --raw-error="<smtp-error-string>" [--json]`
  - Flags: `--raw-error` (string, raw SMTP bounce string to parse), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type EnhancedBounceClassification struct {
        RawErrorString      string `json:"rawErrorString"`
        BasicStatusCode     int    `json:"basicStatusCode"`
        EnhancedStatusCode  string `json:"enhancedStatusCode"`
        RootCauseCategory   string `json:"rootCauseCategory"`
        RemediationAction   string `json:"remediationAction"`
        isPermanentFailure  bool   `json:"isPermanentFailure"`
        isTransientFailure  bool   `json:"isTransientFailure"`
        hasReputationImpact bool   `json:"hasReputationImpact"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceClassification (
        EnhancedStatusCode TEXT PRIMARY KEY,
        RootCauseCategory TEXT NOT NULL,
        RemediationAction TEXT NOT NULL,
        isPermanentFailure INTEGER NOT NULL CHECK (isPermanentFailure IN (0, 1)),
        isTransientFailure INTEGER NOT NULL CHECK (isTransientFailure IN (0, 1)),
        hasReputationImpact INTEGER NOT NULL CHECK (hasReputationImpact IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/bounce/decoder.go` as the foundational parser for all inbound delivery failure telemetry.

---

### Improvement 1572: Multipart Report MIME Parser & DSN Header Extractor
- **Title:** Multipart Report MIME Parser & DSN Header Extractor
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Deconstruct RFC 3464 `multipart/report` MIME structures, extracting human-readable explanations, machine-readable `message/delivery-status` headers (`Final-Recipient`, `Action`, `Remote-MTA`, `Diagnostic-Code`), and original message headers.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Structured inspection card cleanly separating DSN Headers, Reporting MTA, and original Message-ID with key-value alignment.
  - *Web UI:* Collapsible raw MIME tree with instant header search.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce parse-dsn --file=<path-to-eml> [--json]`
  - Flags: `--file` (string, path to raw EML file), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DsnReportMetadata struct {
        OriginalMessageId   string `json:"originalMessageId"`
        FinalRecipient      string `json:"finalRecipient"`
        ReportingMta        string `json:"reportingMta"`
        DsnAction           string `json:"dsnAction"`
        DiagnosticCode      string `json:"diagnosticCode"`
        ReceivedAt          int64  `json:"receivedAt"`
        isFinalBounce       bool   `json:"isFinalBounce"`
        hasOriginalEnvelope bool   `json:"hasOriginalEnvelope"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDsnReports (
        OriginalMessageId TEXT NOT NULL,
        FinalRecipient TEXT NOT NULL,
        ReportingMta TEXT NOT NULL,
        DsnAction TEXT NOT NULL,
        DiagnosticCode TEXT NOT NULL,
        ReceivedAt INTEGER NOT NULL,
        isFinalBounce INTEGER NOT NULL CHECK (isFinalBounce IN (0, 1)),
        hasOriginalEnvelope INTEGER NOT NULL CHECK (hasOriginalEnvelope IN (0, 1)),
        PRIMARY KEY (OriginalMessageId, FinalRecipient)
    );
    ```
- **Antigravity / GitMap Integration Point:** Parser utility in `pkg/mail/bounce/dsn_parser.go`.

---

### Improvement 1573: Authentication-Results & X-Spam-Status Forensic Header Telemetry Extractor
- **Title:** Authentication-Results & X-Spam-Status Forensic Header Telemetry Extractor
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Parse diagnostic anti-spam headers (`Authentication-Results`, `X-Spam-Status`, `X-Google-Smtp-Source`, `X-Forefront-Antispam-Report`) from seed mailboxes and test NDRs to extract granular provider spam scores and authentication verdict tags.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Diagnostic badges showing DKIM (`pass`), SPF (`pass`), DMARC (`pass`), Spam Score (`0.12/5.0`), and Provider Verdict (`Clean` vs `Junked`).
  - *Web UI:* Visual meter comparing inbound spam scores across Google Workspace, Microsoft 365, and Yahoo Mail.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce headers --message-id=<msg-id> [--json]`
  - Flags: `--message-id` (string, message identifier to query), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpamHeaderTelemetry struct {
        MessageId          string  `json:"messageId"`
        TargetProvider     string  `json:"targetProvider"`
        SpfVerdict         string  `json:"spfVerdict"`
        DkimVerdict        string  `json:"dkimVerdict"`
        DmarcVerdict       string  `json:"dmarcVerdict"`
        SpamScore          float64 `json:"spamScore"`
        EvaluatedAt        int64   `json:"evaluatedAt"`
        isSpamFlagged      bool    `json:"isSpamFlagged"`
        hasAuthPassed      bool    `json:"hasAuthPassed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpamHeaderAudit (
        MessageId TEXT PRIMARY KEY,
        TargetProvider TEXT NOT NULL,
        SpfVerdict TEXT NOT NULL,
        DkimVerdict TEXT NOT NULL,
        DmarcVerdict TEXT NOT NULL,
        SpamScore REAL NOT NULL,
        EvaluatedAt INTEGER NOT NULL,
        isSpamFlagged INTEGER NOT NULL CHECK (isSpamFlagged IN (0, 1)),
        hasAuthPassed INTEGER NOT NULL CHECK (hasAuthPassed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Extractor module in `pkg/mail/bounce/header_telemetry.go`.

---

### Improvement 1574: Offline Bayesian Spam Trigger Phrase Probabilistic Scorer & Content Linter
- **Title:** Offline Bayesian Spam Trigger Phrase Probabilistic Scorer & Content Linter
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Evaluate outbound warmup email content against an embedded Bayesian bag-of-words model and known trigger phrase database, scoring the probability of spam placement and highlighting problematic phrases prior to dispatch.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* In-line content linter with highlighted red spans over suspicious trigger words (e.g., "100% free", "guaranteed", "act immediately"), probability meter (`[P(Spam): 0.04]`), and safe substitute suggestions.
  - *Web UI:* WYSIWYG editor with live Bayesian scoring sidebar.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup content lint --subject="<subject>" --body-file=<path> [--threshold=0.15]`
  - Flags: `--subject` (string, email subject line), `--body-file` (string, path to body text), `--threshold` (float64, max acceptable spam probability, default 0.15).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type ContentSpamScoreReport struct {
        ContentHash         string   `json:"contentHash"`
        SpamProbability     float64  `json:"spamProbability"`
        IdentifiedTriggers  []string `json:"identifiedTriggers"`
        SuggestedFixes      []string `json:"suggestedFixes"`
        ScoredAt            int64    `json:"scoredAt"`
        isContentSafe       bool     `json:"isContentSafe"`
        hasTriggerPhrases   bool     `json:"hasTriggerPhrases"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbContentSpamAudit (
        ContentHash TEXT PRIMARY KEY,
        SpamProbability REAL NOT NULL,
        IdentifiedTriggers TEXT NOT NULL,
        SuggestedFixes TEXT NOT NULL,
        ScoredAt INTEGER NOT NULL,
        isContentSafe INTEGER NOT NULL CHECK (isContentSafe IN (0, 1)),
        hasTriggerPhrases INTEGER NOT NULL CHECK (hasTriggerPhrases IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded scorer in `pkg/mail/bounce/bayesian_scorer.go`.

---

### Improvement 1575: Interactive Bubble Tea Terminal Bounce Forensics HUD & Live Streamer
- **Title:** Interactive Bubble Tea Terminal Bounce Forensics HUD & Live Streamer
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Provide a live terminal dashboard displaying real-time incoming delivery failures, categorizing them by ISP (Google, Microsoft, Yahoo, Other) and failure type, with interactive log expansion and keyboard navigation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Dynamic split-pane TUI; top pane charts bounce rate by provider over time; bottom pane streams incoming bounce events with timestamps, recipients, and colored status chips.
  - *Web UI:* Streaming event log with real-time SSE WebSocket push updates.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce hud [--follow] [--filter-provider=<isp>]`
  - Flags: `--follow` (bool, stream live events, default true), `--filter-provider` (string, filter specific ISP).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BounceStreamEvent struct {
        EventId         string `json:"eventId"`
        RecipientEmail  string `json:"recipientEmail"`
        DomainName      string `json:"domainName"`
        ProviderName    string `json:"providerName"`
        EnhancedCode    string `json:"enhancedCode"`
        EventTimestamp  int64  `json:"eventTimestamp"`
        isHardBounce    bool   `json:"isHardBounce"`
        isPolicyBlock   bool   `json:"isPolicyBlock"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceEvents (
        EventId TEXT PRIMARY KEY,
        RecipientEmail TEXT NOT NULL,
        DomainName TEXT NOT NULL,
        ProviderName TEXT NOT NULL,
        EnhancedCode TEXT NOT NULL,
        EventTimestamp INTEGER NOT NULL,
        isHardBounce INTEGER NOT NULL CHECK (isHardBounce IN (0, 1)),
        isPolicyBlock INTEGER NOT NULL CHECK (isPolicyBlock IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Real-time TUI module in `pkg/cli/tui/bounce_hud.go`.

---

### Improvement 1576: Automated Mailbox Quarantine Circuit Breaker & Reputation Protector
- **Title:** Automated Mailbox Quarantine Circuit Breaker & Reputation Protector
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Automatically pause warming or permanently quarantine individual sending mailboxes or recipient pools when hard bounces exceed 4% or spam complaint thresholds cross 0.1%, protecting apex domain IP reputation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Bold red warning panel indicating quarantined mailbox, trigger metric, threshold exceeded, and time quarantined.
  - *Web UI:* Quarantined accounts management table with manual override release buttons.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup quarantine list [--show-all]`
  - `gitmap warmup quarantine release --mailbox=<email> [--reason=<notes>]`
  - Flags: `--mailbox` (string, mailbox to release from quarantine), `--reason` (string, audit log reason).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type MailboxQuarantineRecord struct {
        MailboxEmail     string  `json:"mailboxEmail"`
        QuarantinedAt    int64   `json:"quarantinedAt"`
        TriggerMetric    string  `json:"triggerMetric"`
        ObservedValue    float64 `json:"observedValue"`
        ThresholdValue   float64 `json:"thresholdValue"`
        ReleasedAt       int64   `json:"releasedAt"`
        isQuarantined    bool    `json:"isQuarantined"`
        hasAutoQuarantined bool  `json:"hasAutoQuarantined"`
        isManualOverride bool    `json:"isManualOverride"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbMailboxQuarantine (
        MailboxEmail TEXT PRIMARY KEY,
        QuarantinedAt INTEGER NOT NULL,
        TriggerMetric TEXT NOT NULL,
        ObservedValue REAL NOT NULL,
        ThresholdValue REAL NOT NULL,
        ReleasedAt INTEGER NOT NULL,
        isQuarantined INTEGER NOT NULL CHECK (isQuarantined IN (0, 1)),
        hasAutoQuarantined INTEGER NOT NULL CHECK (hasAutoQuarantined IN (0, 1)),
        isManualOverride INTEGER NOT NULL CHECK (isManualOverride IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded circuit breaker in `pkg/mail/throttle/circuit_breaker.go`.

---

### Improvement 1577: SIEM Webhook Dispatcher (Splunk, Elastic, Datadog) for Critical Deliverability Events
- **Title:** SIEM Webhook Dispatcher (Splunk, Elastic, Datadog) for Critical Deliverability Events
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Stream high-severity delivery alerts (circuit breaker trips, DMARC alignment failures, major ISP reputation blocks) to external enterprise SIEM systems via authenticated HTTP webhooks with backoff and retry guarantees.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Configuration summary card displaying registered SIEM endpoints, auth schemes, and dispatch success rates.
  - *Web UI:* Webhook delivery test button and execution log.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup siem dispatch-test --endpoint=<url> [--payload-type=bounce] [--dry-run]`
  - Flags: `--endpoint` (string, destination webhook URL), `--payload-type` (string, "bounce", "circuit_breaker", "dns_drift"), `--dry-run` (bool, validate payload without sending).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SiemDispatchLog struct {
        DispatchId       string `json:"dispatchId"`
        TargetEndpoint   string `json:"targetEndpoint"`
        EventType        string `json:"eventType"`
        PayloadJson      string `json:"payloadJson"`
        ResponseStatus   int    `json:"responseStatus"`
        DispatchedAt     int64  `json:"dispatchedAt"`
        isSuccessful     bool   `json:"isSuccessful"`
        hasRetryExhausted bool  `json:"hasRetryExhausted"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationSiemLogs (
        DispatchId TEXT PRIMARY KEY,
        TargetEndpoint TEXT NOT NULL,
        EventType TEXT NOT NULL,
        PayloadJson TEXT NOT NULL,
        ResponseStatus INTEGER NOT NULL,
        DispatchedAt INTEGER NOT NULL,
        isSuccessful INTEGER NOT NULL CHECK (isSuccessful IN (0, 1)),
        hasRetryExhausted INTEGER NOT NULL CHECK (hasRetryExhausted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Webhook dispatcher in `pkg/mail/bounce/siem_dispatcher.go`.

---

### Improvement 1578: In-Memory Bloom Filter Suppression List with Sub-Millisecond Membership Tests
- **Title:** In-Memory Bloom Filter Suppression List with Sub-Millisecond Membership Tests
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Maintain millions of bounced, unsubscribed, and complained email addresses in an optimal in-memory Bloom filter backed by local SQLite storage, providing O(1) sub-millisecond suppression checks with zero disk I/O during high-volume warming bursts.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Telemetry box displaying filter capacity, number of elements, estimated false-positive rate (`0.0001%`), and bitset memory utilization (e.g., `1.8 MB`).
  - *Web UI:* Filter rebuild progress bar and test query box.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup suppression test --email=<email>`
  - `gitmap warmup suppression sync [--rebuild]`
  - Flags: `--email` (string, email to check), `--rebuild` (bool, rehydrate bitset from database).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BloomFilterMetadata struct {
        FilterName          string  `json:"filterName"`
        TotalEntries        int64   `json:"totalEntries"`
        BitsetSizeBytes     int64   `json:"bitsetSizeBytes"`
        FalsePositiveRate   float64 `json:"falsePositiveRate"`
        LastRebuiltAt       int64   `json:"lastRebuiltAt"`
        isHydrated          bool    `json:"isHydrated"`
        hasPendingFlushes   bool    `json:"hasPendingFlushes"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSuppressionList (
        EmailHash TEXT PRIMARY KEY,
        OriginalEmail TEXT NOT NULL,
        ReasonCode TEXT NOT NULL,
        AddedAt INTEGER NOT NULL,
        isActive INTEGER NOT NULL CHECK (isActive IN (0, 1)),
        hasPermanentBlock INTEGER NOT NULL CHECK (hasPermanentBlock IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** High-speed filtering service in `pkg/mail/bounce/bloom_suppression.go`.

---

### Improvement 1579: Spam Trap & Honeypot Behavioral Sentinel with Early Warning Signals
- **Title:** Spam Trap & Honeypot Behavioral Sentinel with Early Warning Signals
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Detect subtle signs of pristine or recycled spam trap hits (e.g., zero engagement over 180 days followed by sudden delivery without MX verification, or specific RFC 3464 trap tags) and alert operators before major RBL listings occur.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Critical red radar warning for potential trap contact with domain risk assessment score.
  - *Web UI:* Honeypot mitigation workflow wizard with recipient purge suggestions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trap-audit --domain=<domain> [--scan-recipients] [--json]`
  - Flags: `--domain` (string, domain to analyze), `--scan-recipients` (bool, scan pending queue for honeypot patterns).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpamTrapIncidentReport struct {
        IncidentId        string `json:"incidentId"`
        SuspectEmail      string `json:"suspectEmail"`
        SendingMailbox    string `json:"sendingMailbox"`
        TrapCategory      string `json:"trapCategory"`
        RiskScore         int    `json:"riskScore"`
        DetectedAt        int64  `json:"detectedAt"`
        isQuarantined     bool   `json:"isQuarantined"`
        hasRblListed      bool   `json:"hasRblListed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpamTrapIncidents (
        IncidentId TEXT PRIMARY KEY,
        SuspectEmail TEXT NOT NULL,
        SendingMailbox TEXT NOT NULL,
        TrapCategory TEXT NOT NULL,
        RiskScore INTEGER NOT NULL,
        DetectedAt INTEGER NOT NULL,
        isQuarantined INTEGER NOT NULL CHECK (isQuarantined IN (0, 1)),
        hasRblListed INTEGER NOT NULL CHECK (hasRblListed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Audit engine in `pkg/mail/bounce/spamtrap_sentinel.go`.

---

### Improvement 1580: Longitudinal Predictive Bounce Trend Modeler & Capacity Planner
- **Title:** Longitudinal Predictive Bounce Trend Modeler & Capacity Planner
- **Cluster:** Cluster 158: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Apply exponential smoothing and linear regression over historical warming telemetry to forecast bounce rates for scheduled volume increases, recommending optimal ramp velocities per sending domain.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Projected delivery success curve plotted as an ASCII trend graph, highlighting recommended safe daily volumes and projected warning zones.
  - *Web UI:* Interactive volume vs bounce projection simulator with slider controls.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce forecast --domain=<domain> [--days-ahead=7] [--target-volume=500]`
  - Flags: `--domain` (string, domain to forecast), `--days-ahead` (int, forecast horizon in days, default 7), `--target-volume` (int, target daily volume).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BounceForecastProjection struct {
        DomainName          string  `json:"domainName"`
        ForecastHorizonDays int     `json:"forecastHorizonDays"`
        ProjectedDailyVol   int     `json:"projectedDailyVol"`
        PredictedBounceRate float64 `json:"predictedBounceRate"`
        RecommendedMaxVol   int     `json:"recommendedMaxVol"`
        CalculatedAt        int64   `json:"calculatedAt"`
        isSafeToScale       bool    `json:"isSafeToScale"`
        hasThrottleRisk     bool    `json:"hasThrottleRisk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceForecast (
        DomainName TEXT NOT NULL,
        ForecastHorizonDays INTEGER NOT NULL,
        ProjectedDailyVol INTEGER NOT NULL,
        PredictedBounceRate REAL NOT NULL,
        RecommendedMaxVol INTEGER NOT NULL,
        CalculatedAt INTEGER NOT NULL,
        isSafeToScale INTEGER NOT NULL CHECK (isSafeToScale IN (0, 1)),
        hasThrottleRisk INTEGER NOT NULL CHECK (hasThrottleRisk IN (0, 1)),
        PRIMARY KEY (DomainName, CalculatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Forecasting model in `pkg/mail/bounce/forecast.go`.

---

## 4. Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD (Items 1581–1590)

### Improvement 1581: Vector-Clock CRDT State Synchronizer for Distributed SQLite Replicas
- **Title:** Vector-Clock CRDT State Synchronizer for Distributed SQLite Replicas
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Synchronize warmup state and mailbox dispatch counters across multiple distributed worker nodes using state-based CRDTs and vector clocks, ensuring causal consistency and eventual convergence without central locks.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Cluster topology grid displaying Node IDs, IP endpoints, current Vector Clock epochs, and clock skew in milliseconds (`Δ: +4ms`).
  - *Web UI:* Node sync graph showing causal message paths and convergence latency.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster sync-status [--json]`
  - Flags: `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type VectorClockSyncState struct {
        NodeId           string            `json:"nodeId"`
        VectorEpochs     map[string]uint64 `json:"vectorEpochs"`
        LastSyncEpoch    int64             `json:"lastSyncEpoch"`
        PendingChanges   int               `json:"pendingChanges"`
        isConverged      bool              `json:"isConverged"`
        hasConflictState bool              `json:"hasConflictState"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapClusterNodes (
        NodeId TEXT PRIMARY KEY,
        EndpointAddress TEXT NOT NULL,
        VectorClockJson TEXT NOT NULL,
        LastSyncEpoch INTEGER NOT NULL,
        PendingChanges INTEGER NOT NULL,
        isConverged INTEGER NOT NULL CHECK (isConverged IN (0, 1)),
        hasConflictState INTEGER NOT NULL CHECK (hasConflictState IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Core synchronization daemon in `pkg/db/cluster/crdt_sync.go`.

---

### Improvement 1582: Encrypted WAL Frame Streaming Transport over Mutual TLS 1.3
- **Title:** Encrypted WAL Frame Streaming Transport over Mutual TLS 1.3
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Stream raw SQLite Write-Ahead Log (WAL) frames in real-time from active writers to replica nodes over high-performance streaming mTLS 1.3 sockets, achieving sub-second replication latency with complete transport security.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Bandwidth throughput meter showing frames/sec, compression ratio (zstd), and mTLS certificate validity countdown.
  - *Web UI:* Live WAL streaming pipeline diagram with backpressure indicators.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster stream-wal --target-node=<id> [--compression=zstd] [--port=9443]`
  - Flags: `--target-node` (string, replica node ID), `--compression` (string, "none", "gzip", "zstd", default "zstd"), `--port` (int, replication port, default 9443).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type WalStreamMetricRecord struct {
        StreamId        string `json:"streamId"`
        SourceNodeId    string `json:"sourceNodeId"`
        ReplicaNodeId   string `json:"replicaNodeId"`
        FramesStreamed  int64  `json:"framesStreamed"`
        ThroughputBytes int64  `json:"throughputBytes"`
        RecordedAt      int64  `json:"recordedAt"`
        isTlsActive     bool   `json:"isTlsActive"`
        hasDroppedFrames bool  `json:"hasDroppedFrames"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbWalStreamMetrics (
        StreamId TEXT PRIMARY KEY,
        SourceNodeId TEXT NOT NULL,
        ReplicaNodeId TEXT NOT NULL,
        FramesStreamed INTEGER NOT NULL,
        ThroughputBytes INTEGER NOT NULL,
        RecordedAt INTEGER NOT NULL,
        isTlsActive INTEGER NOT NULL CHECK (isTlsActive IN (0, 1)),
        hasDroppedFrames INTEGER NOT NULL CHECK (hasDroppedFrames IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Transport engine in `pkg/db/cluster/wal_streamer.go`.

---

### Improvement 1583: Page-Level Merkle Tree Shard Integrity Auditor & Diff Localizer
- **Title:** Page-Level Merkle Tree Shard Integrity Auditor & Diff Localizer
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Construct cryptographic Merkle trees over 4096-byte database pages to rapidly compare shard state across nodes, localizing corrupted or divergent pages in logarithmic time without reading entire database files.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Hierarchical Merkle tree visualizer showing Root Hash, Level 1 branches, and highlighting mismatched leaf pages in red with byte offset information.
  - *Web UI:* Interactive tree explorer with page hex inspection modal.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster audit-integrity --db=<path> [--compare-node=<id>] [--repair]`
  - Flags: `--db` (string, database file path), `--compare-node` (string, node to compare against), `--repair` (bool, auto-fetch and overwrite corrupted pages).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type MerkleTreeAuditSummary struct {
        DatabasePath     string   `json:"databasePath"`
        RootMerkleHash   string   `json:"rootMerkleHash"`
        TotalPages       int64    `json:"totalPages"`
        MismatchedPages  []int64  `json:"mismatchedPages"`
        AuditedAt        int64    `json:"auditedAt"`
        isIntegrityValid bool     `json:"isIntegrityValid"`
        hasRepairedPages bool     `json:"hasRepairedPages"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationMerkleAudits (
        DatabasePath TEXT NOT NULL,
        RootMerkleHash TEXT NOT NULL,
        TotalPages INTEGER NOT NULL,
        MismatchedPages TEXT NOT NULL,
        AuditedAt INTEGER NOT NULL,
        isIntegrityValid INTEGER NOT NULL CHECK (isIntegrityValid IN (0, 1)),
        hasRepairedPages INTEGER NOT NULL CHECK (hasRepairedPages IN (0, 1)),
        PRIMARY KEY (DatabasePath, AuditedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Module in `pkg/db/cluster/merkle_audit.go`.

---

### Improvement 1584: Distributed Two-Phase Epoch Commit Coordinator for Split-DB Shards
- **Title:** Distributed Two-Phase Epoch Commit Coordinator for Split-DB Shards
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Coordinate atomic transactions spanning multiple Split-DB shards (`gitmap.db`, `installation.db`, `repodb/pipeline.db`) across distributed workers via an epoch-based two-phase commit (2PC) protocol with automatic timeout rollback.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* 2PC phase indicator (`Prepare` -> `Pre-Commit` -> `Commit`), showing voting consensus status per participating shard node.
  - *Web UI:* Transaction sequence diagram tracing coordinator messages and shard votes.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster 2pc-status [--epoch=<id>] [--json]`
  - Flags: `--epoch` (string, specific commit epoch ID to inspect), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type TwoPhaseCommitEpoch struct {
        EpochId         string   `json:"epochId"`
        CoordinatorId   string   `json:"coordinatorId"`
        ParticipatingShards []string `json:"participatingShards"`
        CommitState     string   `json:"commitState"`
        StartedAt       int64    `json:"startedAt"`
        CommittedAt     int64    `json:"committedAt"`
        isCommitted     bool     `json:"isCommitted"`
        hasAborted      bool     `json:"hasAborted"`
        isTimedOut      bool     `json:"isTimedOut"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapEpochCommits (
        EpochId TEXT PRIMARY KEY,
        CoordinatorId TEXT NOT NULL,
        ParticipatingShards TEXT NOT NULL,
        CommitState TEXT NOT NULL,
        StartedAt INTEGER NOT NULL,
        CommittedAt INTEGER NOT NULL,
        isCommitted INTEGER NOT NULL CHECK (isCommitted IN (0, 1)),
        hasAborted INTEGER NOT NULL CHECK (hasAborted IN (0, 1)),
        isTimedOut INTEGER NOT NULL CHECK (isTimedOut IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Transaction coordinator in `pkg/db/cluster/coordinator.go`.

---

### Improvement 1585: Non-Blocking Hot VACUUM Compaction Engine & Free-List Defragmenter
- **Title:** Non-Blocking Hot VACUUM Compaction Engine & Free-List Defragmenter
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Execute online incremental SQLite page defragmentation (`PRAGMA incremental_vacuum`) and background compaction without obtaining exclusive database locks or interrupting active warm-up message dispatching.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Reclaimed disk space gauge, before/after file size indicators, and freelist page countdown.
  - *Web UI:* Real-time database size and compaction history sparkline.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db compact [--db=<path>] [--max-pages=5000] [--auto-schedule]`
  - Flags: `--db` (string, target SQLite database file), `--max-pages` (int, max freelist pages to reclaim per run, default 5000), `--auto-schedule` (bool, enable periodic background maintenance).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type CompactionEventRecord struct {
        DatabaseName       string `json:"databaseName"`
        InitialSizeBytes   int64  `json:"initialSizeBytes"`
        ReclaimedBytes     int64  `json:"reclaimedBytes"`
        FreelistPagesFreed int    `json:"freelistPagesFreed"`
        DurationMs         int64  `json:"durationMs"`
        ExecutedAt         int64  `json:"executedAt"`
        isSuccessful       bool   `json:"isSuccessful"`
        hasLockContention  bool   `json:"hasLockContention"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationDbCompaction (
        DatabaseName TEXT NOT NULL,
        InitialSizeBytes INTEGER NOT NULL,
        ReclaimedBytes INTEGER NOT NULL,
        FreelistPagesFreed INTEGER NOT NULL,
        DurationMs INTEGER NOT NULL,
        ExecutedAt INTEGER NOT NULL,
        isSuccessful INTEGER NOT NULL CHECK (isSuccessful IN (0, 1)),
        hasLockContention INTEGER NOT NULL CHECK (hasLockContention IN (0, 1)),
        PRIMARY KEY (DatabaseName, ExecutedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Maintenance service in `pkg/db/cluster/compactor.go`.

---

### Improvement 1586: Interactive Bubble Tea Replication Cluster HUD & Topology Navigator
- **Title:** Interactive Bubble Tea Replication Cluster HUD & Topology Navigator
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Author an interactive full-screen Bubble Tea TUI cockpit showing multi-master node status, replication lag per shard, WAL throughput, active lock holders, and leader election states with keyboard selection.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-density cluster grid displaying Node ID, Role (Leader/Follower), Shard Allocations, Lag (ms), and Health Status with Lipgloss color styling.
  - *Web UI:* Full web dashboard with cluster topology graph and node metrics.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster hud [--refresh-interval=1s]`
  - Flags: `--refresh-interval` (duration, dashboard refresh frequency, default 1s).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ClusterNodeTelemetry struct {
        NodeId           string  `json:"nodeId"`
        RoleName         string  `json:"roleName"`
        ReplicationLagMs int64   `json:"replicationLagMs"`
        WalBufferSize    int64   `json:"walBufferSize"`
        ActiveQueries    int     `json:"activeQueries"`
        CpuUsagePercent  float64 `json:"cpuUsagePercent"`
        isHealthy        bool    `json:"isHealthy"`
        isLeader         bool    `json:"isLeader"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbClusterTelemetry (
        NodeId TEXT PRIMARY KEY,
        RoleName TEXT NOT NULL,
        ReplicationLagMs INTEGER NOT NULL,
        WalBufferSize INTEGER NOT NULL,
        ActiveQueries INTEGER NOT NULL,
        CpuUsagePercent REAL NOT NULL,
        isHealthy INTEGER NOT NULL CHECK (isHealthy IN (0, 1)),
        isLeader INTEGER NOT NULL CHECK (isLeader IN (0, 1)),
        UpdatedAt INTEGER NOT NULL
    );
    ```
- **Antigravity / GitMap Integration Point:** Presentation HUD in `pkg/cli/tui/cluster_hud.go`.

---

### Improvement 1587: Deterministic Last-Write-Wins (LWW) Conflict Arbiter with Cryptographic Tie-Breaking
- **Title:** Deterministic Last-Write-Wins (LWW) Conflict Arbiter with Cryptographic Tie-Breaking
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Resolve conflicting concurrent writes to replicated SQLite tables using hybrid logical clocks (HLC) and deterministic Blake3 cryptographic tie-breaking, ensuring all cluster replicas achieve identical state without human intervention.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Conflict resolution log panel indicating Table, Primary Key, Conflicting Timestamps, Winning Node, and Blake3 hash tie-breaker proof.
  - *Web UI:* Conflict audit timeline with side-by-side payload diffing.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster conflicts [--table=<table>] [--tail=50] [--json]`
  - Flags: `--table` (string, filter specific table), `--tail` (int, number of past events to display, default 50), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ConflictResolutionRecord struct {
        ConflictId      string `json:"conflictId"`
        TargetTable     string `json:"targetTable"`
        RecordPrimaryKey string `json:"recordPrimaryKey"`
        WinningNodeId   string `json:"winningNodeId"`
        WinningHlcTime  int64  `json:"winningHlcTime"`
        ResolutionProof string `json:"resolutionProof"`
        ResolvedAt      int64  `json:"resolvedAt"`
        isAutoResolved  bool   `json:"isAutoResolved"`
        hasDataLoss     bool   `json:"hasDataLoss"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbConflictLog (
        ConflictId TEXT PRIMARY KEY,
        TargetTable TEXT NOT NULL,
        RecordPrimaryKey TEXT NOT NULL,
        WinningNodeId TEXT NOT NULL,
        WinningHlcTime INTEGER NOT NULL,
        ResolutionProof TEXT NOT NULL,
        ResolvedAt INTEGER NOT NULL,
        isAutoResolved INTEGER NOT NULL CHECK (isAutoResolved IN (0, 1)),
        hasDataLoss INTEGER NOT NULL CHECK (hasDataLoss IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Arbitration module in `pkg/db/cluster/arbiter.go`.

---

### Improvement 1588: Dynamic Hash-Partition Shard Rebalancer & Rendezvous Router
- **Title:** Dynamic Hash-Partition Shard Rebalancer & Rendezvous Router
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Partition high-volume email recipient and event logs across multiple SQLite shard files using highest-random-weight (rendezvous) hashing, automatically migrating shard assignments as worker nodes join or depart.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Partition distribution bar chart displaying shard balance percentages, active partition counts per node, and in-flight partition migration progress.
  - *Web UI:* Shard assignment visualizer with dynamic ring mapping.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster rebalance [--partitions=64] [--dry-run]`
  - Flags: `--partitions` (int, total virtual partition count, default 64), `--dry-run` (bool, simulate assignment changes).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type PartitionShardAssignment struct {
        PartitionId       int    `json:"partitionId"`
        AssignedNodeId    string `json:"assignedNodeId"`
        AssignedShardFile string `json:"assignedShardFile"`
        TotalRecords      int64  `json:"totalRecords"`
        RebalancedAt      int64  `json:"rebalancedAt"`
        isAssigned        bool   `json:"isAssigned"`
        isMigrating       bool   `json:"isMigrating"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapShardPartitions (
        PartitionId INTEGER PRIMARY KEY,
        AssignedNodeId TEXT NOT NULL,
        AssignedShardFile TEXT NOT NULL,
        TotalRecords INTEGER NOT NULL,
        RebalancedAt INTEGER NOT NULL,
        isAssigned INTEGER NOT NULL CHECK (isAssigned IN (0, 1)),
        isMigrating INTEGER NOT NULL CHECK (isMigrating IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Routing router in `pkg/db/cluster/partition_router.go`.

---

### Improvement 1589: Point-in-Time Recovery (PITR) WAL Log Replayer & Disaster Restoration Engine
- **Title:** Point-in-Time Recovery (PITR) WAL Log Replayer & Disaster Restoration Engine
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Enable exact point-in-time database restoration by replaying sequential WAL frame archives up to a specific millisecond timestamp or transaction epoch, preventing data loss following operator error or disk failure.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive timeline scrubber allowing operators to select a target recovery timestamp, displaying replayed frame count and estimated time to completion.
  - *Web UI:* Visual snapshot selector with replay progress meter.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db restore-pitr --db=<path> --target-timestamp=<epoch-ms> [--dry-run]`
  - Flags: `--db` (string, database to restore), `--target-timestamp` (int64, millisecond recovery target), `--dry-run` (bool, verify WAL frames without writing).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type PitrRecoveryPlan struct {
        TargetDatabase   string `json:"targetDatabase"`
        TargetTimestamp  int64  `json:"targetTimestamp"`
        FramesReplayed   int64  `json:"framesReplayed"`
        RestorationEpoch int64  `json:"restorationEpoch"`
        CompletedAt      int64  `json:"completedAt"`
        isRestored       bool   `json:"isRestored"`
        hasVerificationPassed bool `json:"hasVerificationPassed"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationPitrRestorations (
        RestorationId TEXT PRIMARY KEY,
        TargetDatabase TEXT NOT NULL,
        TargetTimestamp INTEGER NOT NULL,
        FramesReplayed INTEGER NOT NULL,
        RestorationEpoch INTEGER NOT NULL,
        CompletedAt INTEGER NOT NULL,
        isRestored INTEGER NOT NULL CHECK (isRestored IN (0, 1)),
        hasVerificationPassed INTEGER NOT NULL CHECK (hasVerificationPassed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Disaster recovery tool in `pkg/db/cluster/pitr_replayer.go`.

---

### Improvement 1590: Zero-Downtime Rolling Schema Migration Orchestrator with Shadow Tables
- **Title:** Zero-Downtime Rolling Schema Migration Orchestrator with Shadow Tables
- **Cluster:** Cluster 159: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Perform live schema migrations across active distributed SQLite shards using shadow tables, dual-writing triggers, and atomic table swaps (`ALTER TABLE ... RENAME`), ensuring zero downtime and lock-free evolution.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Step-by-step migration progress indicator (`Shadow Created` -> `Backfill 100%` -> `Dual-Write Active` -> `Atomic Swap Completed`).
  - *Web UI:* Database schema diff viewer with historical migration versions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db migrate --migration-file=<path> [--step-by-step] [--rollback]`
  - Flags: `--migration-file` (string, SQL migration script), `--step-by-step` (bool, pause between migration phases), `--rollback` (bool, undo active shadow migration).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type RollingMigrationState struct {
        MigrationId      string `json:"migrationId"`
        TargetDatabase   string `json:"targetDatabase"`
        CurrentPhase     string `json:"currentPhase"`
        RecordsBackfilled int64 `json:"recordsBackfilled"`
        StartedAt        int64  `json:"startedAt"`
        CompletedAt      int64  `json:"completedAt"`
        isCompleted      bool   `json:"isCompleted"`
        hasRollbackOccurred bool `json:"hasRollbackOccurred"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSchemaMigrations (
        MigrationId TEXT PRIMARY KEY,
        TargetDatabase TEXT NOT NULL,
        CurrentPhase TEXT NOT NULL,
        RecordsBackfilled INTEGER NOT NULL,
        StartedAt INTEGER NOT NULL,
        CompletedAt INTEGER NOT NULL,
        isCompleted INTEGER NOT NULL CHECK (isCompleted IN (0, 1)),
        hasRollbackOccurred INTEGER NOT NULL CHECK (hasRollbackOccurred IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Migration engine in `pkg/db/cluster/migrator.go`.

---

## 5. Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (Items 1591–1600)

### Improvement 1591: Real-Time Antigravity Multi-Agent Swarm Cockpit HUD & Telemetry Viewer
- **Title:** Real-Time Antigravity Multi-Agent Swarm Cockpit HUD & Telemetry Viewer
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Author an interactive terminal cockpit dashboard monitoring all autonomous Antigravity subagents, displaying active subagent roles, assigned subtasks, execution durations, context window sizes, and state transitions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-agent grid displaying Agent ID, Role, Status (`IDLE`, `WORKING`, `REVIEWING`), Current Subtask, and Uptime with vibrant status badges.
  - *Web UI:* Real-time agent activity feed and communication graph.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm cockpit [--watch] [--filter-role=<role>]`
  - Flags: `--watch` (bool, live updating terminal view, default true), `--filter-role` (string, filter specific agent role).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type AgentSwarmTelemetry struct {
        AgentId           string `json:"agentId"`
        RoleTitle         string `json:"roleTitle"`
        CurrentTaskId     string `json:"currentTaskId"`
        AgentStatus       string `json:"agentStatus"`
        ContextTokensUsed int    `json:"contextTokensUsed"`
        StartedAt         int64  `json:"startedAt"`
        LastHeartbeatAt   int64  `json:"lastHeartbeatAt"`
        isActive          bool   `json:"isActive"`
        hasActiveTask     bool   `json:"hasActiveTask"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSwarmTelemetry (
        AgentId TEXT PRIMARY KEY,
        RoleTitle TEXT NOT NULL,
        CurrentTaskId TEXT NOT NULL,
        AgentStatus TEXT NOT NULL,
        ContextTokensUsed INTEGER NOT NULL,
        StartedAt INTEGER NOT NULL,
        LastHeartbeatAt INTEGER NOT NULL,
        isActive INTEGER NOT NULL CHECK (isActive IN (0, 1)),
        hasActiveTask INTEGER NOT NULL CHECK (hasActiveTask IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Swarm cockpit in `pkg/cli/tui/swarm_hud.go`.

---

### Improvement 1592: Dynamic Subtask Auctioneer with Expiring Leases & Automatic Requeueing
- **Title:** Dynamic Subtask Auctioneer with Expiring Leases & Automatic Requeueing
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Distribute pending warming subtasks across autonomous agents using an auction mechanism with expiring heartbeat leases (e.g., 60 seconds), automatically revoking and requeueing tasks from crashed or stalled subagents.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Task auction board showing Task ID, Bidding Agents, Winning Agent, Remaining Lease (`42s`), and Priority Weight.
  - *Web UI:* Drag-and-drop task queue with lease countdown gauges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm auction [--lease-sec=60] [--poll-interval=5s]`
  - Flags: `--lease-sec` (int, heartbeat lease duration in seconds, default 60), `--poll-interval` (duration, task poll rate, default 5s).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type SubtaskLeaseRecord struct {
        TaskId          string `json:"taskId"`
        AssignedAgentId string `json:"assignedAgentId"`
        LeaseAcquiredAt int64  `json:"leaseAcquiredAt"`
        LeaseExpiresAt  int64  `json:"leaseExpiresAt"`
        AttemptCount    int    `json:"attemptCount"`
        isLeaseActive   bool   `json:"isLeaseActive"`
        hasReassigned   bool   `json:"hasReassigned"`
        isCompleted     bool   `json:"isCompleted"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSubtaskLeases (
        TaskId TEXT PRIMARY KEY,
        AssignedAgentId TEXT NOT NULL,
        LeaseAcquiredAt INTEGER NOT NULL,
        LeaseExpiresAt INTEGER NOT NULL,
        AttemptCount INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL CHECK (isLeaseActive IN (0, 1)),
        hasReassigned INTEGER NOT NULL CHECK (hasReassigned IN (0, 1)),
        isCompleted INTEGER NOT NULL CHECK (isCompleted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Auction manager in `pkg/agent/swarm/auctioneer.go`.

---

### Improvement 1593: Token Budget Expenditure Governor & Auto-Compaction Gate
- **Title:** Token Budget Expenditure Governor & Auto-Compaction Gate
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Monitor and govern LLM token expenditures across all swarm agents, enforcing hard budget caps per subtask and triggering automated context summarization / compaction before context exhaustion occurs.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Token meter showing cumulative prompt/completion tokens, cost tracking ($), and compaction trigger alerts.
  - *Web UI:* Agent spending breakdown pie chart and token rate over time.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm budget [--set-max-tokens=100000] [--status]`
  - Flags: `--set-max-tokens` (int, budget cap per session), `--status` (bool, display current spending).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type TokenBudgetAccounting struct {
        SessionId          string  `json:"sessionId"`
        PromptTokensTotal  int64   `json:"promptTokensTotal"`
        OutputTokensTotal  int64   `json:"outputTokensTotal"`
        EstimatedCostUsd   float64 `json:"estimatedCostUsd"`
        LastCompactionAt   int64   `json:"lastCompactionAt"`
        isBudgetExceeded   bool    `json:"isBudgetExceeded"`
        hasCompactedRecent bool    `json:"hasCompactedRecent"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationTokenBudgets (
        SessionId TEXT PRIMARY KEY,
        PromptTokensTotal INTEGER NOT NULL,
        OutputTokensTotal INTEGER NOT NULL,
        EstimatedCostUsd REAL NOT NULL,
        LastCompactionAt INTEGER NOT NULL,
        isBudgetExceeded INTEGER NOT NULL CHECK (isBudgetExceeded IN (0, 1)),
        hasCompactedRecent INTEGER NOT NULL CHECK (hasCompactedRecent IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Accounting module in `pkg/agent/swarm/budget_governor.go`.

---

### Improvement 1594: Circular Deadlock Detection Engine (Tarjan SCC) & Agent Self-Healing
- **Title:** Circular Deadlock Detection Engine (Tarjan SCC) & Agent Self-Healing
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Periodically construct the directed dependency graph of active agent tasks and execute Tarjan's Strongly Connected Components (SCC) algorithm to detect circular wait-states, terminating deadlocks by preemptively aborting and re-ordering the lowest priority task.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Visual cycle graph diagram pinpointing the circular dependency nodes with red arrows, accompanied by automated recovery notification.
  - *Web UI:* Live dependency DAG with detected cycle highlights.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm deadlock-check [--auto-heal] [--json]`
  - Flags: `--auto-heal` (bool, automatically break detected cycles, default true), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type DeadlockDetectionIncident struct {
        IncidentId        string   `json:"incidentId"`
        CycleNodes        []string `json:"cycleNodes"`
        PreemptedTaskId   string   `json:"preemptedTaskId"`
        DetectedAt        int64    `json:"detectedAt"`
        isDeadlockPresent bool     `json:"isDeadlockPresent"`
        hasSelfHealed     bool     `json:"hasSelfHealed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDeadlockIncidents (
        IncidentId TEXT PRIMARY KEY,
        CycleNodes TEXT NOT NULL,
        PreemptedTaskId TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        isDeadlockPresent INTEGER NOT NULL CHECK (isDeadlockPresent IN (0, 1)),
        hasSelfHealed INTEGER NOT NULL CHECK (hasSelfHealed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Sentinel in `pkg/agent/swarm/deadlock_detector.go`.

---

### Improvement 1595: Sub-250ms Embedded Raft Consensus Master Elector for Swarm Nodes
- **Title:** Sub-250ms Embedded Raft Consensus Master Elector for Swarm Nodes
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Embed an ultra-lightweight Raft consensus state machine allowing swarm agents across distributed machines to elect a master orchestrator in under 250ms, guaranteeing single-leader task assignment even during network partitions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Raft term counter, leader node badge, heartbeats/sec meter, and election vote distribution matrix.
  - *Web UI:* Real-time Raft log visualization showing committed vs uncommitted state entries.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap swarm raft-status [--step-down]`
  - Flags: `--step-down` (bool, force current leader to abdicate).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type RaftConsensusStatus struct {
        CurrentTerm      uint64 `json:"currentTerm"`
        LeaderNodeId     string `json:"leaderNodeId"`
        VotedFor         string `json:"votedFor"`
        CommitIndex      uint64 `json:"commitIndex"`
        LastHeartbeatAt  int64  `json:"lastHeartbeatAt"`
        isLeader         bool   `json:"isLeader"`
        hasCandidateState bool  `json:"hasCandidateState"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapRaftState (
        CurrentTerm INTEGER PRIMARY KEY,
        LeaderNodeId TEXT NOT NULL,
        VotedFor TEXT NOT NULL,
        CommitIndex INTEGER NOT NULL,
        LastHeartbeatAt INTEGER NOT NULL,
        isLeader INTEGER NOT NULL CHECK (isLeader IN (0, 1)),
        hasCandidateState INTEGER NOT NULL CHECK (hasCandidateState IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Consensus engine in `pkg/agent/swarm/raft_elector.go`.

---

### Improvement 1596: Closed-Loop Deliverability Feedback Governor & Autonomous Throttle Steerer
- **Title:** Closed-Loop Deliverability Feedback Governor & Autonomous Throttle Steerer
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Ingest real-time bounce and open-rate feedback into a Proportional-Integral-Derivative (PID) control loop, autonomously adjusting warmup email dispatch rates and concurrency limits to maintain deliverability within targeted thresholds.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* PID tuning gauge displaying Error Margin, Proportional Term, Integral Drift, Derivative Rate, and current Adjusted Concurrency (`12 msgs/min`).
  - *Web UI:* Closed-loop response curve graph charting deliverability setpoint vs measured performance.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup steer --domain=<domain> [--target-inbox-rate=98.0] [--auto-tune]`
  - Flags: `--domain` (string, domain to govern), `--target-inbox-rate` (float64, desired inbox placement percent, default 98.0), `--auto-tune` (bool, enable autonomous PID adjustments).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type DeliverabilityPidMetrics struct {
        DomainName          string  `json:"domainName"`
        TargetPlacementRate float64 `json:"targetPlacementRate"`
        ActualPlacementRate float64 `json:"actualPlacementRate"`
        PidControlOutput    float64 `json:"pidControlOutput"`
        CalculatedAt        int64   `json:"calculatedAt"`
        isEgressThrottled   bool    `json:"isEgressThrottled"`
        hasDampedOvershoot  bool    `json:"hasDampedOvershoot"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbPidControlMetrics (
        DomainName TEXT NOT NULL,
        TargetPlacementRate REAL NOT NULL,
        ActualPlacementRate REAL NOT NULL,
        PidControlOutput REAL NOT NULL,
        CalculatedAt INTEGER NOT NULL,
        isEgressThrottled INTEGER NOT NULL CHECK (isEgressThrottled IN (0, 1)),
        hasDampedOvershoot INTEGER NOT NULL CHECK (hasDampedOvershoot IN (0, 1)),
        PRIMARY KEY (DomainName, CalculatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Control engine in `pkg/agent/steerer/pid_steerer.go`.

---

### Improvement 1597: Multi-Turn Contextual Synthetic Conversation Generator for Seed Inboxes
- **Title:** Multi-Turn Contextual Synthetic Conversation Generator for Seed Inboxes
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Generate authentic, multi-turn threaded email conversations between diverse synthetic AI personas, simulating realistic human exchange patterns, varied reply intervals, and natural subject threading to build high mailbox engagement scores.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Thread view displaying Persona A and Persona B exchanges, message preview snippets, simulated reading delay (`14m`), and sentiment authenticity score.
  - *Web UI:* Conversation tree explorer showing complete message threads.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup generate-thread --persona-a=<id> --persona-b=<id> [--turns=3] [--topic="Project Update"]`
  - Flags: `--persona-a` (string, sender persona ID), `--persona-b` (string, receiver persona ID), `--turns` (int, conversation depth, default 3), `--topic` (string, seed topic).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type SyntheticConversationThread struct {
        ThreadId         string `json:"threadId"`
        SenderPersonaId  string `json:"senderPersonaId"`
        TargetPersonaId  string `json:"targetPersonaId"`
        SubjectLine      string `json:"subjectLine"`
        TurnCount        int    `json:"turnCount"`
        LastMessageSentAt int64 `json:"lastMessageSentAt"`
        isThreadActive   bool   `json:"isThreadActive"`
        hasRepliesGenerated bool `json:"hasRepliesGenerated"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSyntheticThreads (
        ThreadId TEXT PRIMARY KEY,
        SenderPersonaId TEXT NOT NULL,
        TargetPersonaId TEXT NOT NULL,
        SubjectLine TEXT NOT NULL,
        TurnCount INTEGER NOT NULL,
        LastMessageSentAt INTEGER NOT NULL,
        isThreadActive INTEGER NOT NULL CHECK (isThreadActive IN (0, 1)),
        hasRepliesGenerated INTEGER NOT NULL CHECK (hasRepliesGenerated IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Content generator in `pkg/agent/steerer/conversation_generator.go`.

---

### Improvement 1598: Per-ISP Reputation Vector Heatmap (Google, Microsoft, Yahoo, Apple)
- **Title:** Per-ISP Reputation Vector Heatmap (Google, Microsoft, Yahoo, Apple)
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Compute multi-dimensional reputation scores across distinct major mailbox providers based on inbox placement, spam folder count, bounce rates, and deferral frequency, visualizing provider-specific deliverability vectors.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* 4-provider color-coded matrix showing Google, Microsoft, Yahoo, and Apple reputation scores (0–100), with status tags (`[PRIME]`, `[GOOD]`, `[WARMING]`, `[DEGRADED]`).
  - *Web UI:* Spider radar chart comparing domain standing across all 4 major ecosystems.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup reputation --domain=<domain> [--export-json=<path>]`
  - Flags: `--domain` (string, target sending domain), `--export-json` (string, export raw metrics to JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type IspReputationScorecard struct {
        DomainName       string  `json:"domainName"`
        GoogleScore      float64 `json:"googleScore"`
        MicrosoftScore   float64 `json:"microsoftScore"`
        YahooScore       float64 `json:"yahooScore"`
        AppleScore       float64 `json:"appleScore"`
        CalculatedAt     int64   `json:"calculatedAt"`
        isAllGreen       bool    `json:"isAllGreen"`
        hasDegradedProvider bool `json:"hasDegradedProvider"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbIspReputation (
        DomainName TEXT NOT NULL,
        GoogleScore REAL NOT NULL,
        MicrosoftScore REAL NOT NULL,
        YahooScore REAL NOT NULL,
        AppleScore REAL NOT NULL,
        CalculatedAt INTEGER NOT NULL,
        isAllGreen INTEGER NOT NULL CHECK (isAllGreen IN (0, 1)),
        hasDegradedProvider INTEGER NOT NULL CHECK (hasDegradedProvider IN (0, 1)),
        PRIMARY KEY (DomainName, CalculatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Evaluation engine in `pkg/agent/steerer/reputation_heatmap.go`.

---

### Improvement 1599: Autonomous Domain Warmup Graduation Gate & Production Hand-Off Engine
- **Title:** Autonomous Domain Warmup Graduation Gate & Production Hand-Off Engine
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Autonomously evaluate graduated domains against strict production readiness criteria (e.g., 21 continuous days of warming, >=98% inbox placement across all ISPs, zero spam trap hits, valid DMARC enforcement), promoting ready domains to active production dispatch pools.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Ceremonial graduation banner with checklist items (Inbox Rate, DMARC Status, Volume Quota, Age), green checkmarks, and formal promotion certificate badge.
  - *Web UI:* Domain graduation certificate with printable compliance report.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup graduation-check --domain=<domain> [--promote]`
  - Flags: `--domain` (string, domain to evaluate), `--promote` (bool, execute promotion into production pool if eligible).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type WarmupGraduationRecord struct {
        DomainName         string `json:"domainName"`
        ContinuousDaysWarm int    `json:"continuousDaysWarm"`
        OverallInboxRate   float64 `json:"overallInboxRate"`
        GraduationStatus   string `json:"graduationStatus"`
        EvaluatedAt        int64  `json:"evaluatedAt"`
        isGraduated        bool   `json:"isGraduated"`
        hasFailedCriteria  bool   `json:"hasFailedCriteria"`
        isProductionReady  bool   `json:"isProductionReady"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationDomainGraduation (
        DomainName TEXT PRIMARY KEY,
        ContinuousDaysWarm INTEGER NOT NULL,
        OverallInboxRate REAL NOT NULL,
        GraduationStatus TEXT NOT NULL,
        EvaluatedAt INTEGER NOT NULL,
        isGraduated INTEGER NOT NULL CHECK (isGraduated IN (0, 1)),
        hasFailedCriteria INTEGER NOT NULL CHECK (hasFailedCriteria IN (0, 1)),
        isProductionReady INTEGER NOT NULL CHECK (isProductionReady IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Graduation controller in `pkg/agent/steerer/graduation_gate.go`.

---

### Improvement 1600: Full-Screen Terminal Mission Telemetry Dashboard & Swarm Command Center
- **Title:** Full-Screen Terminal Mission Telemetry Dashboard & Swarm Command Center
- **Cluster:** Cluster 160: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Unify all operational metrics from Innovations 1501 through 1600 into an all-in-one full-screen terminal command center, rendering live session pipelining, socket buffers, DNS DoH consensus, bounce classification, Split-DB replication, and agent swarm status in an integrated high-authority interface.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* 4-quadrant mission control display: Q1 (SMTP Socket Pipelining & TCP Zero-Window), Q2 (DoH DNS Consensus & DNSSEC), Q3 (Bounce Forensics & Bayesian Scorer), Q4 (Split-DB Multi-Master & Antigravity Swarm).
  - *Web UI:* Comprehensive operational cockpit with dark-mode aesthetic and zero external CDN requirements.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mission-control [--fullscreen] [--refresh-rate=500ms]`
  - Flags: `--fullscreen` (bool, utilize full alternate terminal buffer, default true), `--refresh-rate` (duration, update frequency, default 500ms).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/cli/tui/types.go`):
    ```go
    type MissionControlSummaryState struct {
        ActiveWarmupDomains int    `json:"activeWarmupDomains"`
        TotalMessagesToday  int64  `json:"totalMessagesToday"`
        OverallInboxPercent float64 `json:"overallInboxPercent"`
        ActiveAgentsCount   int    `json:"activeAgentsCount"`
        ClusterReplicas     int    `json:"clusterReplicas"`
        SystemUptimeSec     int64  `json:"systemUptimeSec"`
        isAllSystemsNominal bool   `json:"isAllSystemsNominal"`
        hasActiveAlerts     bool   `json:"hasActiveAlerts"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbMissionControl (
        RecordedAt INTEGER PRIMARY KEY,
        ActiveWarmupDomains INTEGER NOT NULL,
        TotalMessagesToday INTEGER NOT NULL,
        OverallInboxPercent REAL NOT NULL,
        ActiveAgentsCount INTEGER NOT NULL,
        ClusterReplicas INTEGER NOT NULL,
        SystemUptimeSec INTEGER NOT NULL,
        isAllSystemsNominal INTEGER NOT NULL CHECK (isAllSystemsNominal IN (0, 1)),
        hasActiveAlerts INTEGER NOT NULL CHECK (hasActiveAlerts IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Flagship command in `pkg/cli/tui/mission_control.go` exposed as `gitmap warmup mission-control`.
