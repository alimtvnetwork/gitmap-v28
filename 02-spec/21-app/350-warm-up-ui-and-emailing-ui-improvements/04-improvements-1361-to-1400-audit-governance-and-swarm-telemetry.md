# Specification 350: Improvements 1361 to 1400 - Audit Governance, DNS Repair, Bounce Forensics & Swarm Telemetry

- **Sequence:** 350
- **Slug:** `350-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/04-improvements-1361-to-1400-audit-governance-and-swarm-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/02-improvements-1301-to-1330-session-pipelining-and-zero-window.md`
  - `02-spec/21-app/350-warm-up-ui-and-emailing-ui-improvements/03-improvements-1331-to-1360-interactive-tui-and-mailbox-diagnostics.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification defines Improvements 1361 through 1400 for Sequence 350 of the GitMap Warm-Up and Emailing UI Engine. Sequence 350 establishes production-grade cryptographic identity resilience, automated DNS record self-healing, forensic SMTP bounce decomposition, distributed multi-master SQLite Split-DB sharding, and intelligent multi-agent swarm orchestration across four core clusters:

1. **Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (Items 1361–1370):** Comprehensive DNSSEC-validated DoH auditing, automated syntax validation, SPF lookup limit flattening (<10 lookups), dual-selector cryptographic DKIM key rotation (RSA-2048 & Ed25519), DMARC alignment policy synthesis, and direct DNS provider API provisioning (Cloudflare, Route53, DigitalOcean) with dry-run safety gates.
2. **Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 1371–1380):** High-precision RFC 3463/5248 Enhanced Mail System Status Code decoding, diagnostic header extraction (Authentication-Results, X-Spam-Status), Bayesian spam filter attribution, automated mailbox quarantine circuit breaking, and real-time terminal bounce forensics HUD.
3. **Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD (Items 1381–1390):** Vector-clock CRDT state synchronization, encrypted WAL frame streaming over TLS, page-level Merkle tree integrity validation, distributed two-phase epoch commits, non-blocking hot vacuum compaction, and terminal replication cluster telemetry.
4. **Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (Items 1391–1400):** Real-time multi-agent swarm cockpit HUD, dynamic subtask leasing and token budget enforcement, circular deadlock prevention, autonomous deliverability feedback steering, and live warmup pacing velocity telemetry.

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 350 Improvements 1361 to 1400 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 137 (1361-1370)    │ Cluster 138 (1371-1380)    │ Cluster 139 (1381-1390)      │
│ Autonomous SPF/DKIM/DMARC  │ Forensic Bounce Analysis   │ Distributed SQLite Split-DB  │
│ Record Repair & DoH Sentinel│ Engine & Classification HUD│ Multi-Master Replication HUD │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 140 (1391-1400)                                 │ Cross-Cutting Standards      │
│ Antigravity Multi-Agent Swarm Orchestrator              │ Positive Booleans (is/has),  │
│ & Autonomous Warmup Steerer                             │ Split-DB, Zero CDN, No Exits │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (Items 1361–1370)

### Improvement 1361: Recursive DoH Multi-Resolver Consensus Sentinel & Anycast Quorum Validator
- **Title:** Recursive DoH Multi-Resolver Consensus Sentinel & Anycast Quorum Validator
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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

### Improvement 1362: DNSSEC Cryptographic Trust Chain Traverser & RRSIG Expiration Monitor
- **Title:** DNSSEC Cryptographic Trust Chain Traverser & RRSIG Expiration Monitor
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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

### Improvement 1363: SPF Syntax Deep-Parser, Macro Expansion Engine & Multiple-Record Detector
- **Title:** SPF Syntax Deep-Parser, Macro Expansion Engine & Multiple-Record Detector
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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
        isSyntaxValid     bool     `json:"isSyntaxValid"`
        hasMultipleSpf    bool     `json:"hasMultipleSpf"`
        hasDeprecatedPtr  bool     `json:"hasDeprecatedPtr"`
        isPermErrorRaised bool     `json:"isPermErrorRaised"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpfSyntaxAudit (
        DomainName TEXT NOT NULL,
        RawRecords TEXT NOT NULL,
        EvaluationResult TEXT NOT NULL,
        RecordCount INTEGER NOT NULL,
        isSyntaxValid INTEGER NOT NULL CHECK (isSyntaxValid IN (0, 1)),
        hasMultipleSpf INTEGER NOT NULL CHECK (hasMultipleSpf IN (0, 1)),
        hasDeprecatedPtr INTEGER NOT NULL CHECK (hasDeprecatedPtr IN (0, 1)),
        isPermErrorRaised INTEGER NOT NULL CHECK (isPermErrorRaised IN (0, 1)),
        AuditedAt INTEGER NOT NULL,
        PRIMARY KEY (DomainName, AuditedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Exposed in `pkg/mail/dnscheck/spf_parser.go` for real-time domain warmup qualification checks.

---

### Improvement 1364: Recursive SPF Lookup Counter & Tree-Graph Flattener Engine
- **Title:** Recursive SPF Lookup Counter & Tree-Graph Flattener Engine
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Recursively resolve nested `include:`, `a`, `mx`, `redirect`, and `exists` SPF mechanisms, building a complete directed acyclic graph (DAG) of DNS lookups, counting total DNS queries against the strict 10-lookup RFC limit, and calculating lookup reduction targets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* ASCII hierarchical tree displaying nested include levels, query cost per node, and an aggregate query meter in Neon Amber (`#F59E0B`) if lookups >= 8 or Crimson (`#EF4444`) if > 10.
  - *Web UI:* Interactive collapsible DAG node graph showing each external domain dependency.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns spf-lookup-tree --domain=<domain> [--depth-limit=5] [--json]`
  - Flags: `--domain` (string, domain name), `--depth-limit` (int, maximum recursion traversal depth, default 5), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpfLookupTreeReport struct {
        RootDomain        string   `json:"rootDomain"`
        TotalDnsLookups   int      `json:"totalDnsLookups"`
        IncludeNodes      []string `json:"includeNodes"`
        VoidLookupsCount  int      `json:"voidLookupsCount"`
        isLookupOverLimit bool     `json:"isLookupOverLimit"`
        hasRecursiveLoop  bool     `json:"hasRecursiveLoop"`
        isSafeToDispatch  bool     `json:"isSafeToDispatch"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpfLookupTree (
        RootDomain TEXT NOT NULL,
        TotalDnsLookups INTEGER NOT NULL,
        IncludeNodes TEXT NOT NULL,
        VoidLookupsCount INTEGER NOT NULL,
        isLookupOverLimit INTEGER NOT NULL CHECK (isLookupOverLimit IN (0, 1)),
        hasRecursiveLoop INTEGER NOT NULL CHECK (hasRecursiveLoop IN (0, 1)),
        isSafeToDispatch INTEGER NOT NULL CHECK (isSafeToDispatch IN (0, 1)),
        EvaluatedAt INTEGER NOT NULL,
        PRIMARY KEY (RootDomain, EvaluatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded into `pkg/mail/dnsrepair/flattener.go` to provide the base lookup graph before triggering automated CIDR compression.

---

### Improvement 1365: Autonomous SPF Flattening Engine with CIDR Deduplication & Range Collapsing
- **Title:** Autonomous SPF Flattening Engine with CIDR Deduplication & Range Collapsing
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Automatically resolve all nested SPF `include` targets into raw IPv4 and IPv6 subnets, eliminate duplicate or overlapping IP blocks via Radix-tree subnet collapsing, and synthesize an optimized, single-record flattened SPF policy that uses 0 DNS lookups.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Side-by-side diff comparison between legacy nested SPF (e.g. 14 lookups) and newly flattened CIDR record (0 lookups), detailing collapsed subnet counts and character length optimization.
  - *Web UI:* Visual IP subnet map showing individual expanded IPs aggregated into supernets.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns spf-flatten --domain=<domain> [--output-file=<path>] [--dry-run] [--json]`
  - Flags: `--domain` (string, target domain), `--output-file` (string, optional path to write generated zone file snippet), `--dry-run` (bool, simulate without applying, default true), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpfFlatteningResult struct {
        ApexDomain            string   `json:"apexDomain"`
        OriginalLookups       int      `json:"originalLookups"`
        FlattenedLookups      int      `json:"flattenedLookups"`
        ExtractedIp4Cidrs     []string `json:"extractedIp4Cidrs"`
        ExtractedIp6Cidrs     []string `json:"extractedIp6Cidrs"`
        SynthesizedRecord     string   `json:"synthesizedRecord"`
        isRecordSplitRequired bool     `json:"isRecordSplitRequired"`
        isFlatteningApplied   bool     `json:"isFlatteningApplied"`
        hasSubnetOverlap      bool     `json:"hasSubnetOverlap"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpfFlatteningLog (
        ApexDomain TEXT NOT NULL,
        OriginalLookups INTEGER NOT NULL,
        FlattenedLookups INTEGER NOT NULL,
        ExtractedIp4Cidrs TEXT NOT NULL,
        ExtractedIp6Cidrs TEXT NOT NULL,
        SynthesizedRecord TEXT NOT NULL,
        isRecordSplitRequired INTEGER NOT NULL CHECK (isRecordSplitRequired IN (0, 1)),
        isFlatteningApplied INTEGER NOT NULL CHECK (isFlatteningApplied IN (0, 1)),
        hasSubnetOverlap INTEGER NOT NULL CHECK (hasSubnetOverlap IN (0, 1)),
        AppliedAt INTEGER NOT NULL,
        PRIMARY KEY (ApexDomain, AppliedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Core logic in `pkg/mail/dnsrepair/cidr_collapse.go`, utilized by auto-repair orchestrator.

---

### Improvement 1366: Dual-Selector DKIM Key Generator (RSA-2048 & Ed25519) with Selector Rotation
- **Title:** Dual-Selector DKIM Key Generator (RSA-2048 & Ed25519) with Selector Rotation
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Generate and rotate cryptographic DKIM keypairs supporting dual selectors: traditional RSA-2048 for legacy MTA compatibility and Ed25519 for modern lightweight cryptographic signing, storing encrypted private keys and preparing public DNS TXT payloads.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Dual-panel cryptographic manifest showing selector names (`gm2026rsa`, `gm2026ed`), algorithm details, SHA-256 public key fingerprints, and formatted DNS TXT record copy-boxes.
  - *Web UI:* Keypair management dashboard with one-click DNS snippet copy and key retirement schedule.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns dkim-rotate --domain=<domain> [--selector-prefix=gm] [--bits=2048] [--json]`
  - Flags: `--domain` (string, domain name), `--selector-prefix` (string, prefix for new selector name, default "gm"), `--bits` (int, RSA bit length, default 2048), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DkimKeypairRecord struct {
        DomainName           string `json:"domainName"`
        RsaSelector          string `json:"rsaSelector"`
        Ed25519Selector      string `json:"ed25519Selector"`
        RsaPublicKeyTxt      string `json:"rsaPublicKeyTxt"`
        EdPublicKeyTxt       string `json:"edPublicKeyTxt"`
        KeyVaultId           string `json:"keyVaultId"`
        CreatedAt            int64  `json:"createdAt"`
        isDualSelectorActive bool   `json:"isDualSelectorActive"`
        hasDnsPublished      bool   `json:"hasDnsPublished"`
        isRetired            bool   `json:"isRetired"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDkimKeypairs (
        DomainName TEXT NOT NULL,
        RsaSelector TEXT NOT NULL,
        Ed25519Selector TEXT NOT NULL,
        RsaPublicKeyTxt TEXT NOT NULL,
        EdPublicKeyTxt TEXT NOT NULL,
        KeyVaultId TEXT NOT NULL,
        CreatedAt INTEGER NOT NULL,
        isDualSelectorActive INTEGER NOT NULL CHECK (isDualSelectorActive IN (0, 1)),
        hasDnsPublished INTEGER NOT NULL CHECK (hasDnsPublished IN (0, 1)),
        isRetired INTEGER NOT NULL CHECK (isRetired IN (0, 1)),
        PRIMARY KEY (DomainName, RsaSelector, CreatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/crypto/dkim_keys.go` with key material saved directly into the encrypted GitMap vault.

---

### Improvement 1367: DMARC Policy Alignment Synthesizer & Stepped Enforcement Orchestrator
- **Title:** DMARC Policy Alignment Synthesizer & Stepped Enforcement Orchestrator
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Analyze current SPF/DKIM alignment scores and autonomously orchestrate a graduated DMARC rollout timeline (`p=none` -> `p=quarantine pct=25` -> `p=quarantine pct=100` -> `p=reject`), generating compliant `_dmarc.<domain>` TXT records with aggregated (rua) and forensic (ruf) reporting URIs.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Milestone stepper progress chart displaying the active DMARC policy stage, quarantine percentage gauge, and alignment confidence score with emerald progress indicators.
  - *Web UI:* Stepped enforcement timeline with historical deliverability vs rejection risk curves.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns dmarc-synthesize --domain=<domain> [--target-policy=quarantine] [--pct=50] [--rua-email=<email>] [--json]`
  - Flags: `--domain` (string, domain name), `--target-policy` (string, none/quarantine/reject, default quarantine), `--pct` (int, enforcement percentage 1-100, default 100), `--rua-email` (string, destination mailto for DMARC aggregate reports), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DmarcPolicyPlan struct {
        DomainName         string `json:"domainName"`
        CurrentPolicy      string `json:"currentPolicy"`
        TargetPolicy       string `json:"targetPolicy"`
        EnforcementPct     int    `json:"enforcementPct"`
        SynthesizedRecord  string `json:"synthesizedRecord"`
        RuaUri             string `json:"ruaUri"`
        RufUri             string `json:"rufUri"`
        isAlignmentStrict  bool   `json:"isAlignmentStrict"`
        hasSubdomainPolicy bool   `json:"hasSubdomainPolicy"`
        isSafeToPromote    bool   `json:"isSafeToPromote"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDmarcPolicyAudit (
        DomainName TEXT NOT NULL,
        CurrentPolicy TEXT NOT NULL,
        TargetPolicy TEXT NOT NULL,
        EnforcementPct INTEGER NOT NULL,
        SynthesizedRecord TEXT NOT NULL,
        RuaUri TEXT NOT NULL,
        RufUri TEXT NOT NULL,
        isAlignmentStrict INTEGER NOT NULL CHECK (isAlignmentStrict IN (0, 1)),
        hasSubdomainPolicy INTEGER NOT NULL CHECK (hasSubdomainPolicy IN (0, 1)),
        isSafeToPromote INTEGER NOT NULL CHECK (isSafeToPromote IN (0, 1)),
        GeneratedAt INTEGER NOT NULL,
        PRIMARY KEY (DomainName, GeneratedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed by `pkg/mail/dnsrepair/dmarc_synthesizer.go` and integrated into deliverability governance reports.

---

### Improvement 1368: Multi-Provider DNS Provider API Provisioner (Cloudflare, Route53, DigitalOcean)
- **Title:** Multi-Provider DNS Provider API Provisioner (Cloudflare, Route53, DigitalOcean)
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Direct programmatic provisioning of synthesized SPF, DKIM, and DMARC TXT records into authoritative DNS zones via provider REST APIs (Cloudflare API v4, AWS Route53 SDK, DigitalOcean DNS API), featuring automated state backup and dry-run safety gates.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive credential status banner and live API provisioning log rendering change diffs (`[+] TXT _dmarc`, `[~] TXT v=spf1`), confirmation prompts, and HTTP response codes.
  - *Web UI:* Provider connection status cards with API key health checks and rollback buttons.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns apply --provider=<cloudflare|route53|digitalocean> --domain=<domain> [--dry-run=false] [--force]`
  - Flags: `--provider` (string, target DNS provider), `--domain` (string, zone domain), `--dry-run` (bool, default true, simulate API calls without committing changes), `--force` (bool, bypass interactive confirmation prompt).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsProvisioningReceipt struct {
        ProviderName      string `json:"providerName"`
        ZoneId            string `json:"zoneId"`
        DomainName        string `json:"domainName"`
        RecordName        string `json:"recordName"`
        RecordType        string `json:"recordType"`
        PreviousValue     string `json:"previousValue"`
        NewValue          string `json:"newValue"`
        ChangeBatchId     string `json:"changeBatchId"`
        ProvisionedAt     int64  `json:"provisionedAt"`
        isDryRun          bool   `json:"isDryRun"`
        isAppliedSuccess  bool   `json:"isAppliedSuccess"`
        hasRollbackBackup bool   `json:"hasRollbackBackup"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsProvisioningLog (
        ChangeBatchId TEXT PRIMARY KEY,
        ProviderName TEXT NOT NULL,
        ZoneId TEXT NOT NULL,
        DomainName TEXT NOT NULL,
        RecordName TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        PreviousValue TEXT NOT NULL,
        NewValue TEXT NOT NULL,
        ProvisionedAt INTEGER NOT NULL,
        isDryRun INTEGER NOT NULL CHECK (isDryRun IN (0, 1)),
        isAppliedSuccess INTEGER NOT NULL CHECK (isAppliedSuccess IN (0, 1)),
        hasRollbackBackup INTEGER NOT NULL CHECK (hasRollbackBackup IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/dnsapi/driver.go` supporting pluggable provider implementations.

---

### Improvement 1369: Global DNS Anycast Propagation Radar & TTL Expiration Monitor
- **Title:** Global DNS Anycast Propagation Radar & TTL Expiration Monitor
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Continuously sample updated DNS records across 24 geographically distributed edge anycast nodes (North America, Europe, Asia-Pacific, Latin America, Africa), measuring live propagation percentage and tracking TTL cache countdowns until global convergence.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-region radar grid displaying world zones (US-East, EU-Central, AP-South) with live propagation status dots (Green = Updated, Yellow = Old TTL In-Flight, Red = Error) and an overall convergence meter.
  - *Web UI:* World map visualizer with pulsating regional nodes and propagation percentage rings.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns propagation-radar --domain=<domain> [--record-type=TXT] [--watch] [--interval=10s]`
  - Flags: `--domain` (string, apex domain), `--record-type` (string, TXT/MX/CNAME, default TXT), `--watch` (bool, stream live TUI updates), `--interval` (duration, polling frequency, default 10s).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsPropagationSnapshot struct {
        DomainName           string  `json:"domainName"`
        RecordType           string  `json:"recordType"`
        TotalNodesSampled    int     `json:"totalNodesSampled"`
        PropagatedNodesCount int     `json:"propagatedNodesCount"`
        PropagationPercent   float64 `json:"propagationPercent"`
        RemainingMaxTtlSec   int     `json:"remainingMaxTtlSec"`
        SampledAt            int64   `json:"sampledAt"`
        isFullyPropagated    bool    `json:"isFullyPropagated"`
        hasRegionalDrift     bool    `json:"hasRegionalDrift"`
        isConvergenceReached bool    `json:"isConvergenceReached"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsPropagationRadar (
        DomainName TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        TotalNodesSampled INTEGER NOT NULL,
        PropagatedNodesCount INTEGER NOT NULL,
        PropagationPercent REAL NOT NULL,
        RemainingMaxTtlSec INTEGER NOT NULL,
        SampledAt INTEGER NOT NULL,
        isFullyPropagated INTEGER NOT NULL CHECK (isFullyPropagated IN (0, 1)),
        hasRegionalDrift INTEGER NOT NULL CHECK (hasRegionalDrift IN (0, 1)),
        isConvergenceReached INTEGER NOT NULL CHECK (isConvergenceReached IN (0, 1)),
        PRIMARY KEY (DomainName, RecordType, SampledAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled in `pkg/mail/dnscheck/propagation.go` to gate high-volume dispatch until propagation reaches >= 95%.

---

### Improvement 1370: Continuous DNS Drift Sentinel & Automated Rollback Circuit Breaker
- **Title:** Continuous DNS Drift Sentinel & Automated Rollback Circuit Breaker
- **Cluster:** Cluster 137: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Periodically audit active DNS records against known valid cryptographically signed baselines; if unauthorized modifications, record deletion, or unintended TTL corruption are detected, automatically trip an emergency circuit breaker that pauses warmup dispatches and triggers automated DNS restoration.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-visibility alert banner flashing Crimson (`#DC2626`) on detected drift, detailing baseline vs detected diff, paused worker process IDs, and rollback recovery status.
  - *Web UI:* Incident timeline card with one-click manual release and baseline override controls.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns drift-watch --domain=<domain> [--check-interval=5m] [--auto-rollback] [--json]`
  - Flags: `--domain` (string, domain name), `--check-interval` (duration, polling frequency, default 5m), `--auto-rollback` (bool, enable automated restoration to last known good baseline, default false), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsDriftIncident struct {
        IncidentId          string `json:"incidentId"`
        DomainName          string `json:"domainName"`
        RecordName          string `json:"recordName"`
        ExpectedHash        string `json:"expectedHash"`
        ObservedHash        string `json:"observedHash"`
        DetectedAt          int64  `json:"detectedAt"`
        isDriftDetected     bool   `json:"isDriftDetected"`
        isWarmupSuspended   bool   `json:"isWarmupSuspended"`
        hasRollbackExecuted bool   `json:"hasRollbackExecuted"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsDriftIncidents (
        IncidentId TEXT PRIMARY KEY,
        DomainName TEXT NOT NULL,
        RecordName TEXT NOT NULL,
        ExpectedHash TEXT NOT NULL,
        ObservedHash TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        isDriftDetected INTEGER NOT NULL CHECK (isDriftDetected IN (0, 1)),
        isWarmupSuspended INTEGER NOT NULL CHECK (isWarmupSuspended IN (0, 1)),
        hasRollbackExecuted INTEGER NOT NULL CHECK (hasRollbackExecuted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded into `pkg/mail/dnscheck/drift_sentinel.go` and connected directly to the central dispatch queue pause channel.

---

## 3. Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 1371–1380)

### Improvement 1371: Granular RFC 3463/5248 Enhanced Mail System Status Code Decoder
- **Title:** Granular RFC 3463/5248 Enhanced Mail System Status Code Decoder
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Ingest Non-Delivery Reports (NDRs) and raw SMTP responses, rigorously extracting and decoding RFC 3463 and RFC 5248 enhanced status code triplets (e.g., `5.7.1` Security/Policy, `5.1.1` Bad Destination Mailbox, `4.4.2` Bad Connection, `5.4.1` DNS Routing Failure) to classify failure modes with mathematical certainty.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Color-coded status badge panel formatting status codes: Soft Bounces (`4.X.X`) in Amber (`#F59E0B`), Hard Bounces (`5.X.X`) in Rose (`#F43F5E`), and Policy Blocks (`5.7.X`) in Purple (`#A855F7`), alongside canonical RFC human descriptions.
  - *Web UI:* Expandable status code tree with historical occurrence frequency charts.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce decode --raw-code=<code_string> [--sample-log=<path>] [--json]`
  - Flags: `--raw-code` (string, e.g. "550 5.7.1 Command rejected"), `--sample-log` (string, path to SMTP delivery log file for mass decoding), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type EnhancedStatusCodeRecord struct {
        ClassCode       int    `json:"classCode"`
        SubjectCode     int    `json:"subjectCode"`
        DetailCode      int    `json:"detailCode"`
        RawEnhancedCode string `json:"rawEnhancedCode"`
        CategoryName    string `json:"categoryName"`
        RfcDescription  string `json:"rfcDescription"`
        isPermanentFail bool   `json:"isPermanentFail"`
        isTransientFail bool   `json:"isTransientFail"`
        isPolicyRelated bool   `json:"isPolicyRelated"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceStatusCodes (
        RawEnhancedCode TEXT PRIMARY KEY,
        ClassCode INTEGER NOT NULL,
        SubjectCode INTEGER NOT NULL,
        DetailCode INTEGER NOT NULL,
        CategoryName TEXT NOT NULL,
        RfcDescription TEXT NOT NULL,
        isPermanentFail INTEGER NOT NULL CHECK (isPermanentFail IN (0, 1)),
        isTransientFail INTEGER NOT NULL CHECK (isTransientFail IN (0, 1)),
        isPolicyRelated INTEGER NOT NULL CHECK (isPolicyRelated IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Used across all bounce processing pipelines in `pkg/mail/bounce/status_decoder.go`.

---

### Improvement 1372: DSN Multipart Report MIME Parser & Diagnostic Header Extractor
- **Title:** DSN Multipart Report MIME Parser & Diagnostic Header Extractor
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Parse complex RFC 3464 `multipart/report` MIME bodies (containing `message/delivery-status` and `message/rfc822` attachments), extracting original Message-ID, Remote-MTA, Final-Recipient, Diagnostic-Code, and envelope return path without truncation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Structured diagnostic summary box detailing Original-Message-ID, Reporting-MTA, Final-Recipient, Action (`failed`/`delayed`), and raw server diagnostic response wrapped in syntax-highlighted borders.
  - *Web UI:* Tabbed MIME viewer displaying human-readable summary, raw headers, and returned headers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce parse-dsn --file=<path_to_eml> [--extract-headers] [--json]`
  - Flags: `--file` (string, path to raw .eml or bounce text message), `--extract-headers` (bool, extract and output original outbound headers), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DsnParsedReport struct {
        OriginalMessageId string `json:"originalMessageId"`
        ReportingMta      string `json:"reportingMta"`
        FinalRecipient    string `json:"finalRecipient"`
        ActionTaken       string `json:"actionTaken"`
        DiagnosticCode    string `json:"diagnosticCode"`
        ReceivedTimestamp int64  `json:"receivedTimestamp"`
        hasDiagnosticText bool   `json:"hasDiagnosticText"`
        isActionFailed    bool   `json:"isActionFailed"`
        isActionDelayed   bool   `json:"isActionDelayed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbParsedDsnReports (
        OriginalMessageId TEXT PRIMARY KEY,
        ReportingMta TEXT NOT NULL,
        FinalRecipient TEXT NOT NULL,
        ActionTaken TEXT NOT NULL,
        DiagnosticCode TEXT NOT NULL,
        ReceivedTimestamp INTEGER NOT NULL,
        hasDiagnosticText INTEGER NOT NULL CHECK (hasDiagnosticText IN (0, 1)),
        isActionFailed INTEGER NOT NULL CHECK (isActionFailed IN (0, 1)),
        isActionDelayed INTEGER NOT NULL CHECK (isActionDelayed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Built in `pkg/mail/bounce/dsn_parser.go` and fed directly into the seed account IMAP inbox listener.

---

### Improvement 1373: Authentication-Results & X-Spam-Status Forensic Header Analyzer
- **Title:** Authentication-Results & X-Spam-Status Forensic Header Analyzer
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Ingest and parse mailbox receiving headers (RFC 7601 `Authentication-Results`, `X-Spam-Status`, `X-Google-Smtp-Source`, `X-Forefront-Antispam-Report`), extracting SPF, DKIM, and DMARC verification outcomes alongside spam scores and triggered SpamAssassin/Barracuda rules.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Security forensics breakdown table showing SPF/DKIM/DMARC pass/fail pills, calculated spam score meter (e.g. `Score: 3.2 / 5.0 [LOW RISK]`), and triggered rule tag badges (`URIBL_BLACK`, `BAYES_50`).
  - *Web UI:* Interactive security pill matrix with breakdown of contributing spam score factors.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce inspect-headers --header-file=<path> [--show-rules] [--json]`
  - Flags: `--header-file` (string, path to text file containing raw headers), `--show-rules` (bool, display all individual triggered spam filter rules), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type AuthHeaderForensics struct {
        MessageId         string   `json:"messageId"`
        SpfResult         string   `json:"spfResult"`
        DkimResult        string   `json:"dkimResult"`
        DmarcResult       string   `json:"dmarcResult"`
        SpamScore         float64  `json:"spamScore"`
        TriggeredRuleTags []string `json:"triggeredRuleTags"`
        AnalyzedAt        int64    `json:"analyzedAt"`
        isSpfPass         bool     `json:"isSpfPass"`
        isDkimPass        bool     `json:"isDkimPass"`
        isDmarcPass       bool     `json:"isDmarcPass"`
        hasSpamScoreRisk  bool     `json:"hasSpamScoreRisk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbAuthHeaderForensics (
        MessageId TEXT PRIMARY KEY,
        SpfResult TEXT NOT NULL,
        DkimResult TEXT NOT NULL,
        DmarcResult TEXT NOT NULL,
        SpamScore REAL NOT NULL,
        TriggeredRuleTags TEXT NOT NULL,
        AnalyzedAt INTEGER NOT NULL,
        isSpfPass INTEGER NOT NULL CHECK (isSpfPass IN (0, 1)),
        isDkimPass INTEGER NOT NULL CHECK (isDkimPass IN (0, 1)),
        isDmarcPass INTEGER NOT NULL CHECK (isDmarcPass IN (0, 1)),
        hasSpamScoreRisk INTEGER NOT NULL CHECK (hasSpamScoreRisk IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/mail/bounce/header_forensics.go` for every inbound seed confirmation message.

---

### Improvement 1374: Offline Naive Bayes Spam Trigger Keyword & Token Probabilistic Scorer
- **Title:** Offline Naive Bayes Spam Trigger Keyword & Token Probabilistic Scorer
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Analyze outbound email subject lines and body copy prior to dispatch using an embedded, zero-dependency Naive Bayes probability model trained on canonical spam corpuses, computing spam likelihood scores and flagging high-risk trigger tokens.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Live copy reviewer highlighting detected spam trigger words in Vivid Crimson (`#EF4444`) or Warning Amber (`#F59E0B`), displaying token spam probability weights and total message confidence meter.
  - *Web UI:* Rich-text drafting area with real-time underwriting linter and suggested synonym replacements.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce score-copy --subject=<text> [--body-file=<path>] [--threshold=0.75] [--json]`
  - Flags: `--subject` (string, email subject line), `--body-file` (string, path to plain text or HTML body), `--threshold` (float, maximum allowable spam score 0.0-1.0, default 0.75), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpamContentScoreReport struct {
        ContentHash        string   `json:"contentHash"`
        TotalSpamScore     float64  `json:"totalSpamScore"`
        TriggerTokensFound []string `json:"triggerTokensFound"`
        WordCount          int      `json:"wordCount"`
        ScoredAt           int64    `json:"scoredAt"`
        isPassingThreshold bool     `json:"isPassingThreshold"`
        hasUrgentTriggers  bool     `json:"hasUrgentTriggers"`
        isFinancialPattern bool     `json:"isFinancialPattern"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpamContentScores (
        ContentHash TEXT PRIMARY KEY,
        TotalSpamScore REAL NOT NULL,
        TriggerTokensFound TEXT NOT NULL,
        WordCount INTEGER NOT NULL,
        ScoredAt INTEGER NOT NULL,
        isPassingThreshold INTEGER NOT NULL CHECK (isPassingThreshold IN (0, 1)),
        hasUrgentTriggers INTEGER NOT NULL CHECK (hasUrgentTriggers IN (0, 1)),
        isFinancialPattern INTEGER NOT NULL CHECK (isFinancialPattern IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Integrated into `pkg/mail/preflight/bayes_scorer.go` as a dispatch pre-condition gate.

---

### Improvement 1375: Real-Time Interactive Terminal Bounce Forensics HUD (Bubble Tea v2)
- **Title:** Real-Time Interactive Terminal Bounce Forensics HUD (Bubble Tea v2)
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Author a full-screen, high-performance Bubble Tea v2 terminal HUD providing interactive navigation of live bounce streams, filtering by failure category (Hard Bounce, Policy Reject, Mailbox Full), and inspecting individual raw NDR messages.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Split-pane terminal layout: Left pane displays scrolling stream of recent bounce events with timestamp, recipient, and status pill; Right pane displays detailed DSN breakdown, enhanced status code analysis, and ISP root-cause diagnostics.
  - *Web UI:* Zero-CDN Web Component replica of the terminal HUD with keyboard navigation shortcuts (`j`, `k`, `enter`, `/`).
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce hud [--domain=<domain>] [--filter=<category>] [--refresh=1s]`
  - Flags: `--domain` (string, filter by sending domain), `--filter` (string, filter by hard/soft/policy/all, default all), `--refresh` (duration, UI polling frequency, default 1s).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BounceHudSnapshot struct {
        TotalBounces24h  int     `json:"totalBounces24h"`
        HardBounceCount  int     `json:"hardBounceCount"`
        SoftBounceCount  int     `json:"softBounceCount"`
        PolicyBlockCount int     `json:"policyBlockCount"`
        ActiveBounceRate float64 `json:"activeBounceRate"`
        LastUpdatedSec   int64   `json:"lastUpdatedSec"`
        isRateElevated   bool    `json:"isRateElevated"`
        hasActiveAlert   bool    `json:"hasActiveAlert"`
        isHudConnected   bool    `json:"isHudConnected"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceHudMetrics (
        DomainName TEXT NOT NULL,
        TotalBounces24h INTEGER NOT NULL,
        HardBounceCount INTEGER NOT NULL,
        SoftBounceCount INTEGER NOT NULL,
        PolicyBlockCount INTEGER NOT NULL,
        ActiveBounceRate REAL NOT NULL,
        LastUpdatedSec INTEGER NOT NULL,
        isRateElevated INTEGER NOT NULL CHECK (isRateElevated IN (0, 1)),
        hasActiveAlert INTEGER NOT NULL CHECK (hasActiveAlert IN (0, 1)),
        isHudConnected INTEGER NOT NULL CHECK (isHudConnected IN (0, 1)),
        PRIMARY KEY (DomainName, LastUpdatedSec)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/cli/tui/bounce_hud.go` as an interactive command under `gitmap warmup`.

---

### Improvement 1376: Automated Mailbox Quarantine Circuit Breaker & Throttle Fallback Daemon
- **Title:** Automated Mailbox Quarantine Circuit Breaker & Throttle Fallback Daemon
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Continually compute trailing 1-hour and 24-hour bounce ratios; if hard bounce rates exceed 4.0% or policy rejections cross 2.0%, automatically trip an emergency circuit breaker, isolating the affected sending account, purging pending queue items, and cutting dispatch pacing by 75%.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Critical alert banner displaying quarantined mailbox address, tripped threshold metric, automated cooling period countdown timer, and reset guidance instructions.
  - *Web UI:* Quarantined account status badge with manual release override switch (requiring explicit typing of "CONFIRM").
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce quarantine-status [--domain=<domain>] [--reset-mailbox=<email>] [--json]`
  - Flags: `--domain` (string, filter by domain), `--reset-mailbox` (string, manually clear quarantine state for specific mailbox), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type MailboxQuarantineRecord struct {
        MailboxAddress    string  `json:"mailboxAddress"`
        DomainName        string  `json:"domainName"`
        TriggeredReason   string  `json:"triggeredReason"`
        TrailingBouncePct float64 `json:"trailingBouncePct"`
        QuarantinedAt     int64   `json:"quarantinedAt"`
        ReleaseEligibleAt int64   `json:"releaseEligibleAt"`
        isQuarantined     bool    `json:"isQuarantined"`
        hasAdminOverride  bool    `json:"hasAdminOverride"`
        isThrottleCutDone bool    `json:"isThrottleCutDone"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbMailboxQuarantine (
        MailboxAddress TEXT PRIMARY KEY,
        DomainName TEXT NOT NULL,
        TriggeredReason TEXT NOT NULL,
        TrailingBouncePct REAL NOT NULL,
        QuarantinedAt INTEGER NOT NULL,
        ReleaseEligibleAt INTEGER NOT NULL,
        isQuarantined INTEGER NOT NULL CHECK (isQuarantined IN (0, 1)),
        hasAdminOverride INTEGER NOT NULL CHECK (hasAdminOverride IN (0, 1)),
        isThrottleCutDone INTEGER NOT NULL CHECK (isThrottleCutDone IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed in `pkg/mail/throttle/quarantine_daemon.go` with event hooks into the central queue dispatcher.

---

### Improvement 1377: SIEM Webhook Event Forwarder (Splunk, Elastic, Datadog) for Bounce Incidents
- **Title:** SIEM Webhook Event Forwarder (Splunk, Elastic, Datadog) for Bounce Incidents
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Stream structured JSON forensic bounce and policy block telemetry to enterprise SIEM platforms (Splunk HEC, Elasticsearch Bulk API, Datadog Logs API, generic HTTPS webhooks) with guaranteed at-least-once delivery, local disk spooling, and exponential backoff retry.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Real-time webhook dispatch counter showing egress throughput, successful HTTP 200 confirmations, queued buffer size, and failed transmission alerts in Cyan (`#06B6D4`).
  - *Web UI:* Outbound webhook delivery history log with payload inspection and replay buttons.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce webhook-config --endpoint=<url> [--provider=<datadog|splunk|elastic|generic>] [--auth-token=<token>] [--enable]`
  - Flags: `--endpoint` (string, HTTPS ingestion URL), `--provider` (string, telemetry format type, default generic), `--auth-token` (string, Bearer/HEC token), `--enable` (bool, toggle forwarder state).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SiemForwarderConfig struct {
        EndpointUrl       string `json:"endpointUrl"`
        ProviderType      string `json:"providerType"`
        SpoolBufferSize   int    `json:"spoolBufferSize"`
        LastForwardedAt   int64  `json:"lastForwardedAt"`
        isForwarderActive bool   `json:"isForwarderActive"`
        hasDiskSpooling   bool   `json:"hasDiskSpooling"`
        isTlsVerified     bool   `json:"isTlsVerified"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSiemWebhookConfig (
        EndpointUrl TEXT PRIMARY KEY,
        ProviderType TEXT NOT NULL,
        SpoolBufferSize INTEGER NOT NULL,
        LastForwardedAt INTEGER NOT NULL,
        isForwarderActive INTEGER NOT NULL CHECK (isForwarderActive IN (0, 1)),
        hasDiskSpooling INTEGER NOT NULL CHECK (hasDiskSpooling IN (0, 1)),
        isTlsVerified INTEGER NOT NULL CHECK (isTlsVerified IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/telemetry/siem_forwarder.go` as an asynchronous background worker.

---

### Improvement 1378: In-Memory Bloom Filter Suppression List with Sub-Millisecond Pre-Flight Checks
- **Title:** In-Memory Bloom Filter Suppression List with Sub-Millisecond Pre-Flight Checks
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Maintain an ultra-high-speed, persistent in-memory Bloom filter containing all known hard-bounced addresses, unsubscribed recipients, and honeypot traps, enabling sub-millisecond pre-flight suppression validation (<0.05ms per check) with zero database overhead.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Memory efficiency telemetry card displaying total suppressed elements count, filter memory footprint (KB), false-positive probability estimate (<0.01%), and query throughput stats.
  - *Web UI:* Suppression list management table with manual CSV import and instant search bar.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce suppression-check --email=<recipient> [--add] [--export-file=<path>] [--json]`
  - Flags: `--email` (string, recipient address to test or suppress), `--add` (bool, manually insert email into suppression filter), `--export-file` (string, export active suppression set to CSV), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BloomSuppressionMeta struct {
        FilterCapacity      int64   `json:"filterCapacity"`
        CurrentElementCount int64   `json:"currentElementCount"`
        MemorySizeBytes     int64   `json:"memorySizeBytes"`
        FalsePositiveRate   float64 `json:"falsePositiveRate"`
        LastSyncedAt        int64   `json:"lastSyncedAt"`
        isFilterLoaded      bool    `json:"isFilterLoaded"`
        hasCapacityWarning  bool    `json:"hasCapacityWarning"`
        isPersistentOnDisk  bool    `json:"isPersistentOnDisk"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSuppressionList (
        RecipientEmail TEXT PRIMARY KEY,
        SuppressedReason TEXT NOT NULL,
        SourceDsnMessageId TEXT NOT NULL,
        SuppressedAt INTEGER NOT NULL,
        isHardBounce INTEGER NOT NULL CHECK (isHardBounce IN (0, 1)),
        isManualSuppression INTEGER NOT NULL CHECK (isManualSuppression IN (0, 1)),
        hasBloomSynced INTEGER NOT NULL CHECK (hasBloomSynced IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Built in `pkg/mail/suppression/bloom_filter.go`, loaded into RAM at GitMap process initialization.

---

### Improvement 1379: Spam Trap & Honeypot Behavioral Signature Sentinel
- **Title:** Spam Trap & Honeypot Behavioral Signature Sentinel
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Detect pristine spam traps, aged honeypot domains, and typofat recipient addresses using heuristic behavioral anomaly detection (e.g. zero open/click activity over extended periods combined with immediate delivery to static MX hosts), automatically isolating high-risk targets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Warning radar panel alerting on suspect honeypot signatures with risk probability score, MX reputation tier, and recommended quarantine action in Bright Amber (`#F59E0B`).
  - *Web UI:* Interactive risk breakdown modal displaying domain registration age, MX diversity, and historical engagement metrics.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce honeypot-scan --domain=<recipient_domain> [--deep-audit] [--json]`
  - Flags: `--domain` (string, recipient domain to evaluate), `--deep-audit` (bool, perform active RDAP and MX reputation cross-checks), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpamTrapRiskAssessment struct {
        RecipientDomain     string  `json:"recipientDomain"`
        TrapRiskScore       float64 `json:"trapRiskScore"`
        DetectedAnomalyType string  `json:"detectedAnomalyType"`
        DomainAgeDays       int     `json:"domainAgeDays"`
        AssessedAt          int64   `json:"assessedAt"`
        isTrapSuspected     bool    `json:"isTrapSuspected"`
        hasPristineMxRecord bool    `json:"hasPristineMxRecord"`
        isAutoIsolated      bool    `json:"isAutoIsolated"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpamTrapAssessments (
        RecipientDomain TEXT PRIMARY KEY,
        TrapRiskScore REAL NOT NULL,
        DetectedAnomalyType TEXT NOT NULL,
        DomainAgeDays INTEGER NOT NULL,
        AssessedAt INTEGER NOT NULL,
        isTrapSuspected INTEGER NOT NULL CHECK (isTrapSuspected IN (0, 1)),
        hasPristineMxRecord INTEGER NOT NULL CHECK (hasPristineMxRecord IN (0, 1)),
        isAutoIsolated INTEGER NOT NULL CHECK (isAutoIsolated IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled in `pkg/mail/bounce/honeypot_sentinel.go` as an inbound recipient validation filter.

---

### Improvement 1380: Longitudinal Predictive Bounce Trend Modeler & Pacing Advisor
- **Title:** Longitudinal Predictive Bounce Trend Modeler & Pacing Advisor
- **Cluster:** Cluster 138: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Model 30-day longitudinal bounce trajectories using linear regression and moving average trend analysis, forecasting deliverability degradation before hard thresholds are breached and issuing prescriptive daily sending limit adjustments.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* ASCII sparkline graph projecting 7-day predicted bounce velocity vs current actual rate, alongside an advisory recommendations box suggesting exact dispatch volume ceilings.
  - *Web UI:* Predictive forecasting chart with shaded confidence intervals and interactive quota slider.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce predict-trend --domain=<domain> [--forecast-days=7] [--json]`
  - Flags: `--domain` (string, domain name), `--forecast-days` (int, days to project into future, default 7), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BounceTrendForecast struct {
        DomainName          string  `json:"domainName"`
        CurrentBounceRate   float64 `json:"currentBounceRate"`
        ForecastBounceRate  float64 `json:"forecastBounceRate"`
        ConfidenceInterval  float64 `json:"confidenceInterval"`
        RecommendedMaxDaily int     `json:"recommendedMaxDaily"`
        GeneratedAt         int64   `json:"generatedAt"`
        isDegradingTrend    bool    `json:"isDegradingTrend"`
        hasPacingReduced    bool    `json:"hasPacingReduced"`
        isSafeToExpand      bool    `json:"isSafeToExpand"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceTrendForecasts (
        DomainName TEXT NOT NULL,
        CurrentBounceRate REAL NOT NULL,
        ForecastBounceRate REAL NOT NULL,
        ConfidenceInterval REAL NOT NULL,
        RecommendedMaxDaily INTEGER NOT NULL,
        GeneratedAt INTEGER NOT NULL,
        isDegradingTrend INTEGER NOT NULL CHECK (isDegradingTrend IN (0, 1)),
        hasPacingReduced INTEGER NOT NULL CHECK (hasPacingReduced IN (0, 1)),
        isSafeToExpand INTEGER NOT NULL CHECK (isSafeToExpand IN (0, 1)),
        PRIMARY KEY (DomainName, GeneratedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed in `pkg/mail/advisor/pacing_modeler.go` to dynamically calibrate daily dispatch quotas.

---

## 4. Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD (Items 1381–1390)

### Improvement 1381: Vector-Clock CRDT State Synchronizer for Split-DB Replicas
- **Title:** Vector-Clock CRDT State Synchronizer for Split-DB Replicas
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Implement causality tracking and conflict-free replicated data types (CRDTs) using vector clocks across distributed SQLite Split-DB nodes, ensuring distributed warmup state and seed accounts converge without data loss across multi-machine clusters.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-node synchronization table detailing Node ID, current vector clock array (`[node-1:42, node-2:19]`), replication lag (ms), and active convergence state in Sky Blue (`#0284C7`).
  - *Web UI:* Node topology mesh graph displaying peer latency links and sync vectors.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster sync-status [--node-id=<id>] [--verbose] [--json]`
  - Flags: `--node-id` (string, query specific node, default local), `--verbose` (bool, display full vector clock arrays), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type VectorClockSyncState struct {
        NodeId           string `json:"nodeId"`
        ClockVectorJson  string `json:"clockVectorJson"`
        ReplicaLagMs     int64  `json:"replicaLagMs"`
        LastSyncedEpoch  int64  `json:"lastSyncedEpoch"`
        isSynchronized   bool   `json:"isSynchronized"`
        hasClockSkew     bool   `json:"hasClockSkew"`
        isConvergenceMet bool   `json:"isConvergenceMet"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbClusterVectorClocks (
        NodeId TEXT PRIMARY KEY,
        ClockVectorJson TEXT NOT NULL,
        ReplicaLagMs INTEGER NOT NULL,
        LastSyncedEpoch INTEGER NOT NULL,
        isSynchronized INTEGER NOT NULL CHECK (isSynchronized IN (0, 1)),
        hasClockSkew INTEGER NOT NULL CHECK (hasClockSkew IN (0, 1)),
        isConvergenceMet INTEGER NOT NULL CHECK (isConvergenceMet IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Built in `pkg/db/cluster/crdt_sync.go` powering multi-machine warmup deployments.

---

### Improvement 1382: Encrypted WAL Frame Streaming Transport over Mutual TLS 1.3
- **Title:** Encrypted WAL Frame Streaming Transport over Mutual TLS 1.3
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Stream raw SQLite Write-Ahead Log (WAL) frames in real-time from active writers to replica followers over mutual TLS 1.3 connections, using ChaCha20-Poly1305 frame encryption, non-blocking ring buffers, and TCP zero-copy transmission.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-speed streaming status monitor showing frames/sec throughput, cipher suite, mTLS peer certificate fingerprint, and frame verification latency.
  - *Web UI:* Frame transfer rate gauge with real-time bandwidth consumption graph.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster stream-wal --peer=<ip:port> [--tls-cert=<path>] [--tls-key=<path>] [--json]`
  - Flags: `--peer` (string, remote cluster peer address), `--tls-cert` (string, path to client certificate), `--tls-key` (string, path to private key), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type WalStreamTransportStats struct {
        PeerAddress       string  `json:"peerAddress"`
        FramesTransmitted int64   `json:"framesTransmitted"`
        BytesThroughputMb float64 `json:"bytesThroughputMb"`
        CipherSuite       string  `json:"cipherSuite"`
        ConnectedSince    int64   `json:"connectedSince"`
        isMtlsEstablished bool    `json:"isMtlsEstablished"`
        hasBufferOverflow bool    `json:"hasBufferOverflow"`
        isStreamHealthy   bool    `json:"isStreamHealthy"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbClusterWalStreamStats (
        PeerAddress TEXT PRIMARY KEY,
        FramesTransmitted INTEGER NOT NULL,
        BytesThroughputMb REAL NOT NULL,
        CipherSuite TEXT NOT NULL,
        ConnectedSince INTEGER NOT NULL,
        isMtlsEstablished INTEGER NOT NULL CHECK (isMtlsEstablished IN (0, 1)),
        hasBufferOverflow INTEGER NOT NULL CHECK (hasBufferOverflow IN (0, 1)),
        isStreamHealthy INTEGER NOT NULL CHECK (isStreamHealthy IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/db/cluster/wal_streamer.go` as the core replication transport.

---

### Improvement 1383: Page-Level Merkle Tree Shard Integrity Auditor
- **Title:** Page-Level Merkle Tree Shard Integrity Auditor
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Construct page-level Merkle hash trees across active SQLite B-tree database pages (`repodb/warmup.db`, `repodb/pipeline.db`), continuously auditing replica data pages against the master root hash to detect and isolate silent bit-rot or desynchronization.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Merkle tree audit summary rendering verified page blocks (`[PASS]`), root hash digests, and detected divergent page indices highlighted in Rose (`#F43F5E`).
  - *Web UI:* Interactive tree hierarchy visualizer allowing drilling into specific leaf page hash mismatches.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster merkle-audit --db=<path> [--check-replica=<peer>] [--json]`
  - Flags: `--db` (string, database file to audit), `--check-replica` (string, optional peer address to cross-compare root hash), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type MerkleAuditReport struct {
        DatabasePath       string `json:"databasePath"`
        TotalPagesAudited  int    `json:"totalPagesAudited"`
        RootMerkleHash     string `json:"rootMerkleHash"`
        MismatchedPageId   int    `json:"mismatchedPageId"`
        AuditedTimestamp   int64  `json:"auditedTimestamp"`
        isIntegrityPass    bool   `json:"isIntegrityPass"`
        hasBitRotDetected  bool   `json:"hasBitRotDetected"`
        isReplicaIdentical bool   `json:"isReplicaIdentical"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbMerkleAuditLog (
        DatabasePath TEXT NOT NULL,
        TotalPagesAudited INTEGER NOT NULL,
        RootMerkleHash TEXT NOT NULL,
        MismatchedPageId INTEGER NOT NULL,
        AuditedTimestamp INTEGER NOT NULL,
        isIntegrityPass INTEGER NOT NULL CHECK (isIntegrityPass IN (0, 1)),
        hasBitRotDetected INTEGER NOT NULL CHECK (hasBitRotDetected IN (0, 1)),
        isReplicaIdentical INTEGER NOT NULL CHECK (isReplicaIdentical IN (0, 1)),
        PRIMARY KEY (DatabasePath, AuditedTimestamp)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/db/cluster/merkle_audit.go` as a scheduled background consistency checker.

---

### Improvement 1384: Distributed Two-Phase Epoch Commit Coordinator for Split-DB
- **Title:** Distributed Two-Phase Epoch Commit Coordinator for Split-DB
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Orchestrate two-phase commit (2PC) epoch barriers across distributed Split-DB nodes for cross-shard operations (such as seed account migration or global quota locking), guaranteeing atomic commit semantics or synchronized rollback.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Epoch coordinator dashboard displaying active Epoch ID, participant vote status (`PREPARED`, `COMMITTED`, `ABORTED`), and phase completion latencies in Emerald (`#10B981`).
  - *Web UI:* Phase 1 / Phase 2 transaction state machine sequence chart.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster epoch-status [--epoch-id=<id>] [--timeout=5s] [--json]`
  - Flags: `--epoch-id` (string, specific epoch to query), `--timeout` (duration, transaction abort timeout, default 5s), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type TwoPhaseEpochRecord struct {
        EpochId            string `json:"epochId"`
        CoordinatorNodeId  string `json:"coordinatorNodeId"`
        ParticipatingNodes int    `json:"participatingNodes"`
        PreparedVotesCount int    `json:"preparedVotesCount"`
        EpochStatus        string `json:"epochStatus"`
        CommittedAt        int64  `json:"committedAt"`
        isCommitSuccessful bool   `json:"isCommitSuccessful"`
        hasTimedOut        bool   `json:"hasTimedOut"`
        isEpochAborted     bool   `json:"isEpochAborted"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbEpochCommitLog (
        EpochId TEXT PRIMARY KEY,
        CoordinatorNodeId TEXT NOT NULL,
        ParticipatingNodes INTEGER NOT NULL,
        PreparedVotesCount INTEGER NOT NULL,
        EpochStatus TEXT NOT NULL,
        CommittedAt INTEGER NOT NULL,
        isCommitSuccessful INTEGER NOT NULL CHECK (isCommitSuccessful IN (0, 1)),
        hasTimedOut INTEGER NOT NULL CHECK (hasTimedOut IN (0, 1)),
        isEpochAborted INTEGER NOT NULL CHECK (isEpochAborted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed in `pkg/db/cluster/epoch_coordinator.go` for multi-node cross-shard transactions.

---

### Improvement 1385: Non-Blocking Hot VACUUM Compaction & Storage Reclamation Engine
- **Title:** Non-Blocking Hot VACUUM Compaction & Storage Reclamation Engine
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Execute online, non-blocking SQLite database compaction (`VACUUM INTO` shadow database followed by atomic file swap), reclaiming disk space from deleted audit logs and WAL frame archives without blocking concurrent read/write transactions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Storage optimization progress bar detailing initial file size, compacted file size, recovered megabytes, and compaction elapsed time in Mint (`#34D399`).
  - *Web UI:* Storage efficiency sparkline showing database footprint reduction over time.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster vacuum-hot --target-db=<path> [--temp-dir=<dir>] [--max-duration=30s] [--json]`
  - Flags: `--target-db` (string, database file to compact), `--temp-dir` (string, temporary build directory for shadow DB), `--max-duration` (duration, timeout limit), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type HotVacuumReceipt struct {
        DatabasePath       string  `json:"databasePath"`
        BytesBefore        int64   `json:"bytesBefore"`
        BytesAfter         int64   `json:"bytesAfter"`
        ReclaimedBytes     int64   `json:"reclaimedBytes"`
        DurationMs         int64   `json:"durationMs"`
        ExecutedAt         int64   `json:"executedAt"`
        isCompactionPass   bool    `json:"isCompactionPass"`
        hasLockContention  bool    `json:"hasLockContention"`
        isAtomicSwapDone   bool    `json:"isAtomicSwapDone"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbHotVacuumLog (
        DatabasePath TEXT NOT NULL,
        BytesBefore INTEGER NOT NULL,
        BytesAfter INTEGER NOT NULL,
        ReclaimedBytes INTEGER NOT NULL,
        DurationMs INTEGER NOT NULL,
        ExecutedAt INTEGER NOT NULL,
        isCompactionPass INTEGER NOT NULL CHECK (isCompactionPass IN (0, 1)),
        hasLockContention INTEGER NOT NULL CHECK (hasLockContention IN (0, 1)),
        isAtomicSwapDone INTEGER NOT NULL CHECK (isAtomicSwapDone IN (0, 1)),
        PRIMARY KEY (DatabasePath, ExecutedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/db/cluster/hot_vacuum.go` during automated maintenance windows.

---

### Improvement 1386: Real-Time Split-DB Replication Cluster HUD (Bubble Tea v2)
- **Title:** Real-Time Split-DB Replication Cluster HUD (Bubble Tea v2)
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Author an interactive full-screen Bubble Tea v2 cluster telemetry HUD rendering live node health, WAL frame replication lag, quorum status, network round-trip latencies, and storage allocation across all Split-DB cluster peers.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-density multi-table cluster overview: Top table displays active cluster nodes with role (`LEADER`/`FOLLOWER`), ping latency, and frame backlog; Bottom pane displays scrolling event log of epoch commits and WAL stream checkpoints.
  - *Web UI:* Zero-CDN Web Component dashboard with live cluster topology graph.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster hud [--refresh=1s] [--sort-by=<lag|name>]`
  - Flags: `--refresh` (duration, update frequency, default 1s), `--sort-by` (string, sort column, default lag).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ClusterHudSnapshot struct {
        ClusterName      string `json:"clusterName"`
        ActiveNodesCount int    `json:"activeNodesCount"`
        QuorumThreshold  int    `json:"quorumThreshold"`
        MaxReplicaLagMs  int64  `json:"maxReplicaLagMs"`
        LastUpdatedEpoch int64  `json:"lastUpdatedEpoch"`
        isQuorumHealthy  bool   `json:"isQuorumHealthy"`
        hasDegradedNode  bool   `json:"hasDegradedNode"`
        isLeaderElected  bool   `json:"isLeaderElected"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbClusterHudTelemetry (
        ClusterName TEXT NOT NULL,
        ActiveNodesCount INTEGER NOT NULL,
        QuorumThreshold INTEGER NOT NULL,
        MaxReplicaLagMs INTEGER NOT NULL,
        LastUpdatedEpoch INTEGER NOT NULL,
        isQuorumHealthy INTEGER NOT NULL CHECK (isQuorumHealthy IN (0, 1)),
        hasDegradedNode INTEGER NOT NULL CHECK (hasDegradedNode IN (0, 1)),
        isLeaderElected INTEGER NOT NULL CHECK (isLeaderElected IN (0, 1)),
        PRIMARY KEY (ClusterName, LastUpdatedEpoch)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/cli/tui/cluster_hud.go` as an operator inspection tool.

---

### Improvement 1387: Deterministic Last-Write-Wins (LWW) Conflict Arbiter with Cryptographic Tie-Breaking
- **Title:** Deterministic Last-Write-Wins (LWW) Conflict Arbiter with Cryptographic Tie-Breaking
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Resolve concurrent multi-master write collisions deterministically using high-precision hybrid logical clocks (HLC) and cryptographic SHA-256 node signature tie-breakers, guaranteeing identical convergence across all cluster replicas.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Conflict resolution ledger table showing conflicting record keys, incoming vs local HLC timestamps, winning node ID, and deterministic arbitration reason.
  - *Web UI:* Side-by-side payload diff viewer displaying resolved conflict records.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster conflict-log [--table=<name>] [--limit=20] [--json]`
  - Flags: `--table` (string, target SQLite table name), `--limit` (int, max log entries, default 20), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ConflictArbiterRecord struct {
        ConflictId         string `json:"conflictId"`
        TargetTable        string `json:"targetTable"`
        PrimaryKeyVal      string `json:"primaryKeyVal"`
        WinningNodeId      string `json:"winningNodeId"`
        WinningTimestamp   int64  `json:"winningTimestamp"`
        DiscardedNodeId    string `json:"discardedNodeId"`
        ResolvedAt         int64  `json:"resolvedAt"`
        isConflictResolved bool   `json:"isConflictResolved"`
        hasTieBreakerUsed  bool   `json:"hasTieBreakerUsed"`
        isDataLossZero     bool   `json:"isDataLossZero"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbClusterConflictLog (
        ConflictId TEXT PRIMARY KEY,
        TargetTable TEXT NOT NULL,
        PrimaryKeyVal TEXT NOT NULL,
        WinningNodeId TEXT NOT NULL,
        WinningTimestamp INTEGER NOT NULL,
        DiscardedNodeId TEXT NOT NULL,
        ResolvedAt INTEGER NOT NULL,
        isConflictResolved INTEGER NOT NULL CHECK (isConflictResolved IN (0, 1)),
        hasTieBreakerUsed INTEGER NOT NULL CHECK (hasTieBreakerUsed IN (0, 1)),
        isDataLossZero INTEGER NOT NULL CHECK (isDataLossZero IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed in `pkg/db/cluster/conflict_arbiter.go` during replica WAL ingestion.

---

### Improvement 1388: Dynamic Hash-Partition Shard Rebalancer & Migration Streamer
- **Title:** Dynamic Hash-Partition Shard Rebalancer & Migration Streamer
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Automatically repartition domain warmup jobs and mailbox datasets across cluster nodes using consistent hashing rings (Ketama algorithm), dynamically migrating shard partitions to new nodes without stopping active email dispatches.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Consistent hashing ring visualization rendering node partition arcs, key balance variance, and migration transfer progress bars in Coral (`#FB7185`).
  - *Web UI:* Circular ring partition diagram showing dynamic rebalancing animations.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster rebalance [--nodes=<count>] [--dry-run] [--json]`
  - Flags: `--nodes` (int, target cluster size), `--dry-run` (bool, simulate shard redistribution), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ShardRebalancePlan struct {
        PlanId             string `json:"planId"`
        TotalKeysCount     int64  `json:"totalKeysCount"`
        MigratingKeysCount int64  `json:"migratingKeysCount"`
        SourceNodeId       string `json:"sourceNodeId"`
        TargetNodeId       string `json:"targetNodeId"`
        StartedAt          int64  `json:"startedAt"`
        isPlanComplete     bool   `json:"isPlanComplete"`
        hasActiveStreaming bool   `json:"hasActiveStreaming"`
        isZeroDowntimeMet  bool   `json:"isZeroDowntimeMet"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbShardRebalanceLog (
        PlanId TEXT PRIMARY KEY,
        TotalKeysCount INTEGER NOT NULL,
        MigratingKeysCount INTEGER NOT NULL,
        SourceNodeId TEXT NOT NULL,
        TargetNodeId TEXT NOT NULL,
        StartedAt INTEGER NOT NULL,
        isPlanComplete INTEGER NOT NULL CHECK (isPlanComplete IN (0, 1)),
        hasActiveStreaming INTEGER NOT NULL CHECK (hasActiveStreaming IN (0, 1)),
        isZeroDowntimeMet INTEGER NOT NULL CHECK (isZeroDowntimeMet IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/db/cluster/shard_rebalancer.go` during cluster scale-out events.

---

### Improvement 1389: Point-in-Time Recovery (PITR) WAL Log Replayer & Time-Travel Debugger
- **Title:** Point-in-Time Recovery (PITR) WAL Log Replayer & Time-Travel Debugger
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Replay archived SQLite WAL frames up to a precise microsecond timestamp or epoch ID, reconstructing previous database states into an isolated shadow database for audit forensics, disaster recovery, or debugging unexpected state transitions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Time-travel replay controller showing target timestamp, replayed frame progress meter, current shadow state snapshot, and verification checksums in Indigo (`#6366F1`).
  - *Web UI:* Interactive timeline scrubber allowing visual exploration of historical database states.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster pitr-restore --target-time=<rfc3339> --archive-dir=<path> --output-db=<path> [--json]`
  - Flags: `--target-time` (string, point-in-time recovery target timestamp), `--archive-dir` (string, path to WAL frame archive), `--output-db` (string, destination shadow SQLite DB), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type PitrReplayReceipt struct {
        TargetTimestamp    int64  `json:"targetTimestamp"`
        FramesReplayed     int64  `json:"framesReplayed"`
        OutputDbPath       string `json:"outputDbPath"`
        RestoredChecksum   string `json:"restoredChecksum"`
        ReplayDurationMs   int64  `json:"replayDurationMs"`
        isRestoreSuccess   bool   `json:"isRestoreSuccess"`
        hasTargetReached   bool   `json:"hasTargetReached"`
        isShadowDbVerified bool   `json:"isShadowDbVerified"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbPitrReplayLog (
        TargetTimestamp INTEGER NOT NULL,
        FramesReplayed INTEGER NOT NULL,
        OutputDbPath TEXT NOT NULL,
        RestoredChecksum TEXT NOT NULL,
        ReplayDurationMs INTEGER NOT NULL,
        isRestoreSuccess INTEGER NOT NULL CHECK (isRestoreSuccess IN (0, 1)),
        hasTargetReached INTEGER NOT NULL CHECK (hasTargetReached IN (0, 1)),
        isShadowDbVerified INTEGER NOT NULL CHECK (isShadowDbVerified IN (0, 1)),
        ExecutedAt INTEGER NOT NULL,
        PRIMARY KEY (TargetTimestamp, ExecutedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/db/cluster/pitr_replayer.go` as a recovery engine.

---

### Improvement 1390: Zero-Downtime Rolling Schema Migration Orchestrator for Cluster Shards
- **Title:** Zero-Downtime Rolling Schema Migration Orchestrator for Cluster Shards
- **Cluster:** Cluster 139: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Execute schema migrations across distributed multi-master Split-DB nodes using backward-compatible expand-and-contract patterns, ensuring mixed-version cluster nodes operate concurrently without query failures during version upgrades.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Rolling migration phase card rendering current migration step (`EXPAND_COLUMNS`, `BACKFILL`, `CONTRACT_LEGACY`), per-node version matrix, and health check gates.
  - *Web UI:* Version compatibility grid showing node migration progress.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap cluster migrate --migration-dir=<path> [--step=<expand|contract>] [--dry-run]`
  - Flags: `--migration-dir` (string, path to SQL migration files), `--step` (string, expand/backfill/contract), `--dry-run` (bool, validate syntax without applying).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type RollingMigrationState struct {
        MigrationId        string `json:"migrationId"`
        TargetVersion      int    `json:"targetVersion"`
        CurrentPhase       string `json:"currentPhase"`
        NodesCompleted     int    `json:"nodesCompleted"`
        TotalClusterNodes  int    `json:"totalClusterNodes"`
        AppliedAt          int64  `json:"appliedAt"`
        isExpandComplete   bool   `json:"isExpandComplete"`
        isContractComplete bool   `json:"isContractComplete"`
        isRollingSafe      bool   `json:"isRollingSafe"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbRollingMigrations (
        MigrationId TEXT PRIMARY KEY,
        TargetVersion INTEGER NOT NULL,
        CurrentPhase TEXT NOT NULL,
        NodesCompleted INTEGER NOT NULL,
        TotalClusterNodes INTEGER NOT NULL,
        AppliedAt INTEGER NOT NULL,
        isExpandComplete INTEGER NOT NULL CHECK (isExpandComplete IN (0, 1)),
        isContractComplete INTEGER NOT NULL CHECK (isContractComplete IN (0, 1)),
        isRollingSafe INTEGER NOT NULL CHECK (isRollingSafe IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed by `pkg/db/cluster/rolling_migrator.go` invoked during GitMap binary upgrades.

---

## 5. Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (Items 1391–1400)

### Improvement 1391: Antigravity Multi-Agent Swarm Cockpit HUD (Bubble Tea v2)
- **Title:** Antigravity Multi-Agent Swarm Cockpit HUD (Bubble Tea v2)
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Author a terminal mission control cockpit in Bubble Tea v2 for supervising distributed Antigravity autonomous agents, rendering live agent states (`IDLE`, `ANALYZING`, `DISPATCHING`, `PAUSED`), active task assignments, token velocity, and throughput sparklines.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* 3-column cockpit HUD: Column 1 displays Agent Swarm Roster with pulsing status glyphs; Column 2 displays current subtask execution flamegraph and active prompt context; Column 3 displays deliverability telemetry and real-time throttle gauges.
  - *Web UI:* Zero-CDN Web Component agent grid with live event feed.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent swarm-hud [--view=<full|compact>] [--refresh=500ms]`
  - Flags: `--view` (string, full cockpit or compact one-line status bar), `--refresh` (duration, UI refresh rate, default 500ms).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type SwarmCockpitSnapshot struct {
        SwarmId            string  `json:"swarmId"`
        TotalAgentsCount   int     `json:"totalAgentsCount"`
        ActiveAgentsCount  int     `json:"activeAgentsCount"`
        CompletedTasks24h  int64   `json:"completedTasks24h"`
        SwarmVelocityScore float64 `json:"swarmVelocityScore"`
        RecordedAt         int64   `json:"recordedAt"`
        isSwarmHealthy     bool    `json:"isSwarmHealthy"`
        hasIdleCapacity    bool    `json:"hasIdleCapacity"`
        isCockpitActive    bool    `json:"isCockpitActive"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSwarmCockpitMetrics (
        SwarmId TEXT NOT NULL,
        TotalAgentsCount INTEGER NOT NULL,
        ActiveAgentsCount INTEGER NOT NULL,
        CompletedTasks24h INTEGER NOT NULL,
        SwarmVelocityScore REAL NOT NULL,
        RecordedAt INTEGER NOT NULL,
        isSwarmHealthy INTEGER NOT NULL CHECK (isSwarmHealthy IN (0, 1)),
        hasIdleCapacity INTEGER NOT NULL CHECK (hasIdleCapacity IN (0, 1)),
        isCockpitActive INTEGER NOT NULL CHECK (isCockpitActive IN (0, 1)),
        PRIMARY KEY (SwarmId, RecordedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/cli/tui/swarm_hud.go` as the central agent operations interface.

---

### Improvement 1392: Dynamic Subtask Auctioneer & Expiring Lease Scheduler
- **Title:** Dynamic Subtask Auctioneer & Expiring Lease Scheduler
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Autonomously schedule and distribute decomposed warmup subtasks to worker agents using a market auction model with short-lived leases (e.g. 60-second lease with heartbeat renewals), automatically reclaiming and reassigning tasks if an agent stalls or crashes.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Active lease auction board showing Task ID, assigned Agent Conversation ID, remaining lease countdown in Amber (`#F59E0B`), and heartbeat ping indicator.
  - *Web UI:* Interactive Kanban board showing dynamic lease assignments and auto-reclaims.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent lease-list [--status=<active|expired>] [--json]`
  - Flags: `--status` (string, filter active or expired leases), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type SubtaskLeaseRecord struct {
        TaskId             string `json:"taskId"`
        AssignedAgentId    string `json:"assignedAgentId"`
        LeaseDurationSec   int    `json:"leaseDurationSec"`
        LeaseExpiresAt     int64  `json:"leaseExpiresAt"`
        LastHeartbeatAt    int64  `json:"lastHeartbeatAt"`
        isLeaseActive      bool   `json:"isLeaseActive"`
        hasHeartbeatMissed bool   `json:"hasHeartbeatMissed"`
        isReclaimed        bool   `json:"isReclaimed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSubtaskLeases (
        TaskId TEXT PRIMARY KEY,
        AssignedAgentId TEXT NOT NULL,
        LeaseDurationSec INTEGER NOT NULL,
        LeaseExpiresAt INTEGER NOT NULL,
        LastHeartbeatAt INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL CHECK (isLeaseActive IN (0, 1)),
        hasHeartbeatMissed INTEGER NOT NULL CHECK (hasHeartbeatMissed IN (0, 1)),
        isReclaimed INTEGER NOT NULL CHECK (isReclaimed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed in `pkg/agent/swarm/auctioneer.go` coordinating worker task pools.

---

### Improvement 1393: Token Budget Expenditure Governor & Context Window Comptroller
- **Title:** Token Budget Expenditure Governor & Context Window Comptroller
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Enforce per-agent token expenditure limits and context window memory ceilings, triggering automated context compaction, summarization, and task checkpointing when an agent approaches 85% of its allowed window budget.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Token utilization gauge per agent rendering prompt tokens, completion tokens, remaining context budget, and compaction warning badges in Orange (`#EA580C`).
  - *Web UI:* Context window usage graph showing memory compaction inflection points.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent token-audit [--agent-id=<id>] [--warn-threshold=0.85] [--json]`
  - Flags: `--agent-id` (string, target agent), `--warn-threshold` (float, threshold fraction, default 0.85), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type AgentTokenBudgetRecord struct {
        AgentId            string  `json:"agentId"`
        PromptTokensUsed   int64   `json:"promptTokensUsed"`
        OutputTokensUsed   int64   `json:"outputTokensUsed"`
        MaxContextLimit    int64   `json:"maxContextLimit"`
        ContextUsageRatio  float64 `json:"contextUsageRatio"`
        AuditedAt          int64   `json:"auditedAt"`
        isCompactionNeeded bool    `json:"isCompactionNeeded"`
        hasBudgetExceeded  bool    `json:"hasBudgetExceeded"`
        isThrottled        bool    `json:"isThrottled"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbAgentTokenBudgets (
        AgentId TEXT PRIMARY KEY,
        PromptTokensUsed INTEGER NOT NULL,
        OutputTokensUsed INTEGER NOT NULL,
        MaxContextLimit INTEGER NOT NULL,
        ContextUsageRatio REAL NOT NULL,
        AuditedAt INTEGER NOT NULL,
        isCompactionNeeded INTEGER NOT NULL CHECK (isCompactionNeeded IN (0, 1)),
        hasBudgetExceeded INTEGER NOT NULL CHECK (hasBudgetExceeded IN (0, 1)),
        isThrottled INTEGER NOT NULL CHECK (isThrottled IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/agent/swarm/token_comptroller.go` wrapping LLM call invocations.

---

### Improvement 1394: Tarjan-Based Circular Deadlock Detection & Autonomous Worker Respawn
- **Title:** Tarjan-Based Circular Deadlock Detection & Autonomous Worker Respawn
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Continuously construct a directed dependency graph of inter-agent messages and locked resources, applying Tarjan's strongly connected components algorithm to detect circular wait deadlocks, automatically aborting the lowest-priority lock holder and respawning unresponsive agents.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Dependency graph visualizer highlighting circular deadlock loops in Bright Crimson (`#DC2626`), with diagnostic text detailing aborted cycle links and respawn actions.
  - *Web UI:* Live graph visualizer showing lock dependency edges and cycle resolutions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent deadlock-check [--auto-break] [--json]`
  - Flags: `--auto-break` (bool, automatically terminate deadlock participant, default true), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type DeadlockIncidentReport struct {
        IncidentId         string   `json:"incidentId"`
        DeadlockCycleNodes []string `json:"deadlockCycleNodes"`
        AbortedLockId      string   `json:"abortedLockId"`
        RespawnedAgentId   string   `json:"respawnedAgentId"`
        DetectedAt         int64    `json:"detectedAt"`
        isDeadlockFound    bool     `json:"isDeadlockFound"`
        hasCycleBroken     bool     `json:"hasCycleBroken"`
        isWorkerRespawned  bool     `json:"isWorkerRespawned"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDeadlockIncidents (
        IncidentId TEXT PRIMARY KEY,
        DeadlockCycleNodes TEXT NOT NULL,
        AbortedLockId TEXT NOT NULL,
        RespawnedAgentId TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        isDeadlockFound INTEGER NOT NULL CHECK (isDeadlockFound IN (0, 1)),
        hasCycleBroken INTEGER NOT NULL CHECK (hasCycleBroken IN (0, 1)),
        isWorkerRespawned INTEGER NOT NULL CHECK (isWorkerRespawned IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/agent/swarm/deadlock_detector.go` running on a 2-second background ticker.

---

### Improvement 1395: Sub-250ms Embedded Raft Consensus Module for Swarm Master Election
- **Title:** Sub-250ms Embedded Raft Consensus Module for Swarm Master Election
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Maintain high-availability swarm coordination using an embedded, ultra-lightweight Raft consensus state machine, completing leader failover and election terms in under 250 milliseconds when the active orchestrator node disconnects.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Raft cluster term telemetry box detailing active Term ID, current Leader Node, election timeout countdown (ms), and voter quorum status in Mint (`#34D399`).
  - *Web UI:* State machine node cards showing Leader / Candidate / Follower transitions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent raft-status [--step-down] [--json]`
  - Flags: `--step-down` (bool, force current leader to yield leadership to trigger election), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type RaftConsensusStatus struct {
        ClusterTermId     int64  `json:"clusterTermId"`
        LeaderNodeId      string `json:"leaderNodeId"`
        CurrentRole       string `json:"currentRole"`
        ElectionTimeoutMs int    `json:"electionTimeoutMs"`
        VotedForNodeId    string `json:"votedForNodeId"`
        UpdatedAt         int64  `json:"updatedAt"`
        isLeaderActive    bool   `json:"isLeaderActive"`
        hasQuorumAchieved bool   `json:"hasQuorumAchieved"`
        isTermSplit       bool   `json:"isTermSplit"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbRaftConsensusMetrics (
        ClusterTermId INTEGER NOT NULL,
        LeaderNodeId TEXT NOT NULL,
        CurrentRole TEXT NOT NULL,
        ElectionTimeoutMs INTEGER NOT NULL,
        VotedForNodeId TEXT NOT NULL,
        UpdatedAt INTEGER NOT NULL,
        isLeaderActive INTEGER NOT NULL CHECK (isLeaderActive IN (0, 1)),
        hasQuorumAchieved INTEGER NOT NULL CHECK (hasQuorumAchieved IN (0, 1)),
        isTermSplit INTEGER NOT NULL CHECK (isTermSplit IN (0, 1)),
        PRIMARY KEY (ClusterTermId, UpdatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Built in `pkg/agent/swarm/raft_consensus.go` providing swarm failover resilience.

---

### Improvement 1396: Closed-Loop Deliverability Feedback Governor & Autonomous Throttle Steerer
- **Title:** Closed-Loop Deliverability Feedback Governor & Autonomous Throttle Steerer
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Dynamically couple incoming delivery signals (inbox placement rates, spam folder placements, bounce velocity) directly into the outbound dispatch rate engine, autonomously steering pacing velocity up or down without manual human intervention.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Pacing steering tachometer gauge displaying current dispatches/hour, real-time throttle scalar (`0.25x` to `2.0x`), and PID control loop feedback errors in Cyan (`#06B6D4`).
  - *Web UI:* Dynamic velocity curve chart overlaid with deliverability health signals.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent steer-status [--domain=<domain>] [--set-target-rate=<rate>] [--json]`
  - Flags: `--domain` (string, domain name), `--set-target-rate` (int, manual target override dispatches/day), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type PacingSteererSnapshot struct {
        DomainName          string  `json:"domainName"`
        CurrentRatePerHour  int     `json:"currentRatePerHour"`
        TargetRatePerHour   int     `json:"targetRatePerHour"`
        ThrottleMultiplier  float64 `json:"throttleMultiplier"`
        PidErrorScore       float64 `json:"pidErrorScore"`
        SteeredAt           int64   `json:"steeredAt"`
        isPacingAccelerated bool    `json:"isPacingAccelerated"`
        hasThrottleApplied  bool    `json:"hasThrottleApplied"`
        isTargetRateMet     bool    `json:"isTargetRateMet"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbPacingSteererLog (
        DomainName TEXT NOT NULL,
        CurrentRatePerHour INTEGER NOT NULL,
        TargetRatePerHour INTEGER NOT NULL,
        ThrottleMultiplier REAL NOT NULL,
        PidErrorScore REAL NOT NULL,
        SteeredAt INTEGER NOT NULL,
        isPacingAccelerated INTEGER NOT NULL CHECK (isPacingAccelerated IN (0, 1)),
        hasThrottleApplied INTEGER NOT NULL CHECK (hasThrottleApplied IN (0, 1)),
        isTargetRateMet INTEGER NOT NULL CHECK (isTargetRateMet IN (0, 1)),
        PRIMARY KEY (DomainName, SteeredAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed by `pkg/agent/steerer/throttle_governor.go` connected directly to the outbound mail queue.

---

### Improvement 1397: Multi-Turn Synthetic Contextual Conversation Thread Generator
- **Title:** Multi-Turn Synthetic Contextual Conversation Thread Generator
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Synthesize coherent, multi-turn threaded email dialogues across seed peer personas (handling `In-Reply-To`, `References`, subject `Re:` prefixes, and natural conversational cadence) to simulate authentic human business exchanges for mailbox reputation building.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Thread conversation visualizer displaying threaded message tree, persona avatars, subject line variations, and semantic coherence confidence scores in Violet (`#8B5CF6`).
  - *Web UI:* Email client simulation showing full threaded conversation history.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent generate-thread --topic=<business|tech|sales> [--turns=4] [--dry-run] [--json]`
  - Flags: `--topic` (string, discussion domain theme), `--turns` (int, conversation depth 2-6, default 3), `--dry-run` (bool, preview generated text without queueing), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type ConversationThreadPlan struct {
        ThreadRootId        string   `json:"threadRootId"`
        TopicCategory       string   `json:"topicCategory"`
        TotalTurnsCount     int      `json:"totalTurnsCount"`
        ParticipantSeeds    []string `json:"participantSeeds"`
        AverageTurnDelaySec int      `json:"averageTurnDelaySec"`
        GeneratedAt         int64    `json:"generatedAt"`
        isThreadCoherent    bool     `json:"isThreadCoherent"`
        hasReferencesValid  bool     `json:"hasReferencesValid"`
        isFullyDispatched   bool     `json:"isFullyDispatched"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbConversationThreads (
        ThreadRootId TEXT PRIMARY KEY,
        TopicCategory TEXT NOT NULL,
        TotalTurnsCount INTEGER NOT NULL,
        ParticipantSeeds TEXT NOT NULL,
        AverageTurnDelaySec INTEGER NOT NULL,
        GeneratedAt INTEGER NOT NULL,
        isThreadCoherent INTEGER NOT NULL CHECK (isThreadCoherent IN (0, 1)),
        hasReferencesValid INTEGER NOT NULL CHECK (hasReferencesValid IN (0, 1)),
        isFullyDispatched INTEGER NOT NULL CHECK (isFullyDispatched IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Built in `pkg/agent/steerer/thread_generator.go` utilized by peer engagement workers.

---

### Improvement 1398: Per-ISP Reputation Vector Heatmap Visualizer (Google, Microsoft, Yahoo, Apple)
- **Title:** Per-ISP Reputation Vector Heatmap Visualizer (Google, Microsoft, Yahoo, Apple)
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Compute multi-dimensional reputation vectors (inbox rate, bounce frequency, latency, spam score) segregated by destination ISP MX provider, rendering high-contrast terminal heatmaps to identify provider-specific delivery degradation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Color-coded ISP reputation matrix: Rows represent sending domains, Columns represent major ISPs (Google, Microsoft 365, Yahoo, Apple iCloud, Proton), cells colored from Dark Emerald (95-100% reputation) to Deep Crimson (<70% reputation).
  - *Web UI:* Interactive 2D heatmap grid with tooltip breakdowns of ISP-specific error rates.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent isp-heatmap [--domain=<domain>] [--window=7d] [--json]`
  - Flags: `--domain` (string, domain name or all), `--window` (duration, evaluation time window, default 7d), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type IspReputationVector struct {
        DomainName         string  `json:"domainName"`
        IspProviderName    string  `json:"ispProviderName"`
        ReputationScore    float64 `json:"reputationScore"`
        InboxPlacementRate float64 `json:"inboxPlacementRate"`
        SampleCount        int     `json:"sampleCount"`
        CalculatedAt       int64   `json:"calculatedAt"`
        isHealthyScore     bool    `json:"isHealthyScore"`
        hasThrottleRisk    bool    `json:"hasThrottleRisk"`
        isProviderTier1    bool    `json:"isProviderTier1"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbIspReputationHeatmap (
        DomainName TEXT NOT NULL,
        IspProviderName TEXT NOT NULL,
        ReputationScore REAL NOT NULL,
        InboxPlacementRate REAL NOT NULL,
        SampleCount INTEGER NOT NULL,
        CalculatedAt INTEGER NOT NULL,
        isHealthyScore INTEGER NOT NULL CHECK (isHealthyScore IN (0, 1)),
        hasThrottleRisk INTEGER NOT NULL CHECK (hasThrottleRisk IN (0, 1)),
        isProviderTier1 INTEGER NOT NULL CHECK (isProviderTier1 IN (0, 1)),
        PRIMARY KEY (DomainName, IspProviderName, CalculatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/cli/tui/isp_heatmap.go` and `pkg/agent/steerer/reputation_vector.go`.

---

### Improvement 1399: Autonomous Domain Warmup Graduation Gate & Readiness Certifier
- **Title:** Autonomous Domain Warmup Graduation Gate & Readiness Certifier
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Continuously evaluate domains against mathematical graduation criteria (e.g. minimum 21 days active warming, >98% inbox rate across top 3 ISPs, zero hard bounces in 72 hours, fully aligned DMARC), issuing cryptographically signed graduation certificates enabling production volume dispatches.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Formal graduation diploma card rendered with double-line borders (`╔═╗`), displaying domain name, warmup duration, certified daily send quota, and cryptographic certificate hash in Gold (`#F59E0B`).
  - *Web UI:* Certification badge download panel with printable PDF/SVG audit report.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent evaluate-graduation --domain=<domain> [--certify] [--json]`
  - Flags: `--domain` (string, domain name), `--certify` (bool, issue formal signed graduation certificate if criteria met), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type WarmupGraduationCertificate struct {
        CertificateId       string  `json:"certificateId"`
        DomainName          string  `json:"domainName"`
        TotalWarmupDays     int     `json:"totalWarmupDays"`
        AggregateInboxRate  float64 `json:"aggregateInboxRate"`
        ApprovedDailyQuota  int     `json:"approvedDailyQuota"`
        CertifiedTimestamp  int64   `json:"certifiedTimestamp"`
        SignatureHash       string  `json:"signatureHash"`
        isGraduationPassed  bool    `json:"isGraduationPassed"`
        hasZeroRecentBounce bool    `json:"hasZeroRecentBounce"`
        isProductionReady   bool    `json:"isProductionReady"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbGraduationCertificates (
        CertificateId TEXT PRIMARY KEY,
        DomainName TEXT NOT NULL,
        TotalWarmupDays INTEGER NOT NULL,
        AggregateInboxRate REAL NOT NULL,
        ApprovedDailyQuota INTEGER NOT NULL,
        CertifiedTimestamp INTEGER NOT NULL,
        SignatureHash TEXT NOT NULL,
        isGraduationPassed INTEGER NOT NULL CHECK (isGraduationPassed IN (0, 1)),
        hasZeroRecentBounce INTEGER NOT NULL CHECK (hasZeroRecentBounce IN (0, 1)),
        isProductionReady INTEGER NOT NULL CHECK (isProductionReady IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Governed in `pkg/agent/steerer/graduation_certifier.go` unlocking production sending limits.

---

### Improvement 1400: Unified Executive Mission Telemetry HUD & Swarm Operations Dashboard
- **Title:** Unified Executive Mission Telemetry HUD & Swarm Operations Dashboard
- **Cluster:** Cluster 140: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Consolidate all subsystem metrics (SMTP pipelining flamegraphs, kernel socket buffers, seed inbox delivery status, DoH DNS health, bounce forensics, SQLite shard replication lag, and agent swarm throughput) into a singular, high-performance Bubble Tea v2 executive mission cockpit.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Full-screen 4-quadrant executive mission display:
    - *Top-Left:* SMTP Pipelining & Socket Buffer gauge meters.
    - *Top-Right:* Autonomous DNS & DoH Sentinel status.
    - *Bottom-Left:* Forensic Bounce Classifier & Suppression list activity.
    - *Bottom-Right:* Split-DB Multi-Master Replication & Swarm Agent velocity.
  - *Web UI:* Zero-CDN responsive Shadow DOM executive dashboard with fullscreen kiosk mode.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mission-hud [--refresh=500ms] [--kiosk] [--theme=<dark|dracula|matrix>]`
  - Flags: `--refresh` (duration, UI refresh rate, default 500ms), `--kiosk` (bool, hide border navigation for dedicated wallboards), `--theme` (string, terminal color palette, default dark).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/cli/tui/types.go`):
    ```go
    type ExecutiveMissionTelemetry struct {
        SystemInstanceId   string  `json:"systemInstanceId"`
        ActiveDispatches   int     `json:"activeDispatches"`
        GlobalDeliveryRate float64 `json:"globalDeliveryRate"`
        ActiveReplicas     int     `json:"activeReplicas"`
        SwarmAgentsOnline  int     `json:"swarmAgentsOnline"`
        SnapshotEpoch      int64   `json:"snapshotEpoch"`
        isOverallGreen     bool    `json:"isOverallGreen"`
        hasActiveAlerts    bool    `json:"hasActiveAlerts"`
        isKioskMode        bool    `json:"isKioskMode"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbExecutiveMissionSnapshots (
        SystemInstanceId TEXT NOT NULL,
        ActiveDispatches INTEGER NOT NULL,
        GlobalDeliveryRate REAL NOT NULL,
        ActiveReplicas INTEGER NOT NULL,
        SwarmAgentsOnline INTEGER NOT NULL,
        SnapshotEpoch INTEGER NOT NULL,
        isOverallGreen INTEGER NOT NULL CHECK (isOverallGreen IN (0, 1)),
        hasActiveAlerts INTEGER NOT NULL CHECK (hasActiveAlerts IN (0, 1)),
        isKioskMode INTEGER NOT NULL CHECK (isKioskMode IN (0, 1)),
        PRIMARY KEY (SystemInstanceId, SnapshotEpoch)
    );
    ```
- **Antigravity / GitMap Integration Point:** Located at `pkg/cli/tui/mission_hud.go`, serving as the primary operational entry point `gitmap warmup mission-hud`.

---

## 6. Verification Checklist & Quality Gates

- [ ] All 40 improvements (1361 through 1400) specify complete CLI commands, flag definitions, and Lipgloss/Web UX specifications.
- [ ] All Go structs declare positive-only boolean fields (`is`, `has`) with zero negative booleans or banned prefixes (`can`, `should`, `was`, `will`).
- [ ] All SQLite DDL tables conform to Split-DB architecture (`repodb/pipeline.db`), PascalCase table names, and camelCase column definitions.
- [ ] All file paths in this specification use 100% relative paths (`02-spec/...`, `pkg/...`, `repodb/...`) with zero drive letters.
- [ ] Zero external CDN dependencies; all Web UI components use inline SVGs, Shadow DOM, and vanilla CSS.
