# Specification 351: Improvements 1461 to 1500 - Audit Governance, DNS Repair, Bounce Forensics & Swarm Telemetry

- **Sequence:** 351
- **Slug:** `351-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/04-improvements-1461-to-1500-audit-governance-and-swarm-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/02-improvements-1401-to-1430-session-pipelining-and-zero-window.md`
  - `02-spec/21-app/351-warm-up-ui-and-emailing-ui-improvements/03-improvements-1431-to-1460-interactive-tui-and-mailbox-diagnostics.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification defines Improvements 1461 through 1500 for Sequence 351 of the GitMap Warm-Up and Emailing UI Engine. Sequence 351 deepens high-resilience DNS over HTTPS (DoH) consensus self-healing, RFC 3463/5248 diagnostic bounce telemetry, distributed SQLite Split-DB multi-master replication, and autonomous multi-agent swarm orchestration across four major clusters:

1. **Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (Items 1461–1470):** Anycast DoH multi-resolver consensus auditing, DNSSEC NSEC3 authenticated denial-of-existence proof validation, RFC 7208 syntax validation with macro expansion, recursive SPF lookup tree flattening (<10 queries) with CIDR collapsing, dual-selector cryptographic DKIM key rotation (RSA-2048 & Ed25519), DMARC policy alignment synthesizers with stepped enforcement (`p=none` -> `p=quarantine` -> `p=reject`), multi-provider API provisioning (Cloudflare, Route53, DigitalOcean), global anycast propagation radars, and automated DNS drift rollback circuit breakers.
2. **Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 1471–1480):** High-precision RFC 3463/5248 Enhanced Mail System Status Code decoding, multipart report MIME parsers with DSN header extraction, `Authentication-Results` & `X-Spam-Status` header telemetry, offline Bayesian trigger phrase probabilistic scoring, interactive Bubble Tea terminal bounce forensics HUD, automated mailbox quarantine circuit breaking, SIEM webhook dispatchers (Splunk, Elastic, Datadog), in-memory Bloom filter suppression lists, spam trap honeypot behavioral sentinels, and longitudinal predictive bounce trend modeling.
3. **Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD (Items 1481–1490):** Vector-clock CRDT state synchronization across distributed worker replicas, encrypted WAL frame streaming transport over mutual TLS 1.3, page-level Merkle tree shard integrity auditors, distributed two-phase epoch commit coordinators, non-blocking hot VACUUM compaction engines, real-time terminal replication cluster HUDs, deterministic Last-Write-Wins (LWW) conflict arbiters with cryptographic tie-breaking, dynamic hash-partition shard rebalancers, point-in-time recovery (PITR) WAL log replayers, and zero-downtime rolling schema migration orchestrators.
4. **Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (Items 1491–1500):** Real-time multi-agent swarm cockpit HUDs, dynamic subtask auctioneers with expiring leases, token budget expenditure governors, Tarjan-based circular deadlock detection and worker self-healing, sub-250ms embedded Raft consensus master electors, closed-loop deliverability feedback governors with autonomous throttle steering, multi-turn synthetic contextual conversation generators, per-ISP reputation vector heatmaps (Google, Microsoft, Yahoo, Apple), autonomous domain warmup graduation gates, and full-screen unified mission telemetry dashboards.

```text
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 351 Improvements 1461 to 1500 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 147 (1461-1470)    │ Cluster 148 (1471-1480)    │ Cluster 149 (1481-1490)      │
│ Autonomous SPF/DKIM/DMARC  │ Forensic Bounce Analysis   │ Distributed SQLite Split-DB  │
│ Record Repair & DoH Sentinel│ Engine & Classification HUD│ Multi-Master Replication HUD │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 150 (1491-1500)                                 │ Cross-Cutting Standards      │
│ Antigravity Multi-Agent Swarm Orchestrator              │ Positive Booleans (is/has),  │
│ & Autonomous Warmup Steerer                             │ Split-DB, Zero CDN, No Exits │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel (Items 1461–1470)

### Improvement 1461: Recursive DoH Multi-Resolver Consensus Sentinel & Anycast Quorum Validator
- **Title:** Recursive DoH Multi-Resolver Consensus Sentinel & Anycast Quorum Validator
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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

### Improvement 1462: DNSSEC Cryptographic Trust Chain Traverser & RRSIG Expiration Monitor
- **Title:** DNSSEC Cryptographic Trust Chain Traverser & RRSIG Expiration Monitor
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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

### Improvement 1463: SPF Syntax Deep-Parser, Macro Expansion Engine & Multiple-Record Detector
- **Title:** SPF Syntax Deep-Parser, Macro Expansion Engine & Multiple-Record Detector
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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

### Improvement 1464: Recursive SPF Lookup Counter & Tree-Graph Flattener Engine
- **Title:** Recursive SPF Lookup Counter & Tree-Graph Flattener Engine
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
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
        GeneratedAt INTEGER NOT NULL,
        PRIMARY KEY (RootDomain, GeneratedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Integrated into `pkg/mail/dnsrepair/tree_counter.go` to trigger auto-flattening pipelines.

---

### Improvement 1465: Autonomous SPF Flattening Engine with CIDR Deduplication & Range Collapsing
- **Title:** Autonomous SPF Flattening Engine with CIDR Deduplication & Range Collapsing
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Collapse nested third-party SPF include chains into direct IPv4 and IPv6 CIDR blocks, deduplicate overlapping subnet ranges using prefix tree arithmetic, and author single-line TXT records guaranteed to use <= 1 DNS lookup.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Before/after diff table highlighting query count reduction (`14 lookups -> 1 lookup`), string character lengths, and CIDR range consolidation stats.
  - *Web UI:* Real-time CIDR map visualizer with IP inclusion verification search box.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns spf-flatten --domain=<domain> [--max-chars=450] [--output-file=<path>]`
  - Flags: `--domain` (string, target domain), `--max-chars` (int, maximum characters per TXT string, default 450), `--output-file` (string, optional file to write flattened record).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpfFlattenPlan struct {
        TargetDomain     string   `json:"targetDomain"`
        OriginalLookups  int      `json:"originalLookups"`
        FlattenedLookups int      `json:"flattenedLookups"`
        ExtractedIpv4    []string `json:"extractedIpv4"`
        ExtractedIpv6    []string `json:"extractedIpv6"`
        SynthesizedTxt   string   `json:"synthesizedTxt"`
        isWithinCharCap  bool     `json:"isWithinCharCap"`
        hasCidrMerged    bool     `json:"hasCidrMerged"`
        isReadyToPublish bool     `json:"isReadyToPublish"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpfFlattenHistory (
        TargetDomain TEXT NOT NULL,
        OriginalLookups INTEGER NOT NULL,
        FlattenedLookups INTEGER NOT NULL,
        ExtractedIpv4Count INTEGER NOT NULL,
        ExtractedIpv6Count INTEGER NOT NULL,
        SynthesizedTxt TEXT NOT NULL,
        isWithinCharCap INTEGER NOT NULL CHECK (isWithinCharCap IN (0, 1)),
        hasCidrMerged INTEGER NOT NULL CHECK (hasCidrMerged IN (0, 1)),
        isReadyToPublish INTEGER NOT NULL CHECK (isReadyToPublish IN (0, 1)),
        CreatedAt INTEGER NOT NULL,
        PRIMARY KEY (TargetDomain, CreatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/dnsrepair/flattener.go` and triggered on DNS audit warnings.

---

### Improvement 1466: Dual-Selector DKIM Key Generator (RSA-2048 & Ed25519) with Selector Rotation
- **Title:** Dual-Selector DKIM Key Generator (RSA-2048 & Ed25519) with Selector Rotation
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Autonomously generate dual cryptographic DKIM key pairs (modern Ed25519 for low overhead and RFC 8463 compliance, plus RSA-2048 for universal legacy receiver fallback), formatting DNS TXT RR payloads and executing zero-downtime selector rotation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Dual-pane cryptographic key card displaying key fingerprint (SHA-256), public DNS selector name (`gm202610._domainkey`), bit length, and active rotation lifecycle status tag in Cyan (`#06B6D4`).
  - *Web UI:* One-click DNS copy button and live DKIM signature verification test panel.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dkim generate-keys --domain=<domain> [--selector-prefix=gm] [--rotate] [--json]`
  - Flags: `--domain` (string, domain name), `--selector-prefix` (string, selector naming prefix, default "gm"), `--rotate` (bool, retire old selector and mark new selector active), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DkimKeyPairMeta struct {
        DomainName     string `json:"domainName"`
        Selector       string `json:"selector"`
        AlgorithmType  string `json:"algorithmType"`
        PublicKeyDns   string `json:"publicKeyDns"`
        PrivateKeyPath string `json:"privateKeyPath"`
        RotatedAt      int64  `json:"rotatedAt"`
        isActive       bool   `json:"isActive"`
        isEd25519Dual  bool   `json:"isEd25519Dual"`
        hasDnsPublished bool  `json:"hasDnsPublished"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationDkimKeys (
        DomainName TEXT NOT NULL,
        Selector TEXT NOT NULL,
        AlgorithmType TEXT NOT NULL,
        PublicKeyDns TEXT NOT NULL,
        PrivateKeyPath TEXT NOT NULL,
        RotatedAt INTEGER NOT NULL,
        isActive INTEGER NOT NULL CHECK (isActive IN (0, 1)),
        isEd25519Dual INTEGER NOT NULL CHECK (isEd25519Dual IN (0, 1)),
        hasDnsPublished INTEGER NOT NULL CHECK (hasDnsPublished IN (0, 1)),
        PRIMARY KEY (DomainName, Selector)
    );
    ```
- **Antigravity / GitMap Integration Point:** Located in `pkg/mail/crypto/dkim_rotate.go` for cryptographic email header signing.

---

### Improvement 1467: DMARC Policy Alignment Synthesizer & Stepped Enforcement Orchestrator
- **Title:** DMARC Policy Alignment Synthesizer & Stepped Enforcement Orchestrator
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Synthesize strict RFC 7489 DMARC records with configured aggregate (`rua`) and forensic (`ruf`) report endpoints, managing automated stepped enforcement transitions (`p=none` -> `p=quarantine pct=25` -> `pct=100` -> `p=reject`) based on warm-up pass rates.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Policy progression ladder indicating current enforcement stage, alignment percentages for SPF/DKIM, and scheduled milestone dates for the next policy tightening.
  - *Web UI:* Interactive policy slider showing anticipated deliverability impact across mailbox providers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dmarc advance-policy --domain=<domain> [--target-policy=quarantine] [--pct=50] [--dry-run]`
  - Flags: `--domain` (string, domain name), `--target-policy` (string, none/quarantine/reject), `--pct` (int, enforcement percentage 1-100), `--dry-run` (bool, simulate without DNS update).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DmarcPolicyState struct {
        DomainName       string `json:"domainName"`
        CurrentPolicy    string `json:"currentPolicy"`
        EnforcementPct   int    `json:"enforcementPct"`
        RuaEndpoint      string `json:"ruaEndpoint"`
        RufEndpoint      string `json:"rufEndpoint"`
        PassRateSpfDkim  int    `json:"passRateSpfDkim"`
        LastAdvancedAt   int64  `json:"lastAdvancedAt"`
        isEnforcementMax bool   `json:"isEnforcementMax"`
        hasStrictAlign   bool   `json:"hasStrictAlign"`
        isCompliant      bool   `json:"isCompliant"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDmarcPolicy (
        DomainName TEXT NOT NULL PRIMARY KEY,
        CurrentPolicy TEXT NOT NULL,
        EnforcementPct INTEGER NOT NULL,
        RuaEndpoint TEXT NOT NULL,
        RufEndpoint TEXT NOT NULL,
        PassRateSpfDkim INTEGER NOT NULL,
        LastAdvancedAt INTEGER NOT NULL,
        isEnforcementMax INTEGER NOT NULL CHECK (isEnforcementMax IN (0, 1)),
        hasStrictAlign INTEGER NOT NULL CHECK (hasStrictAlign IN (0, 1)),
        isCompliant INTEGER NOT NULL CHECK (isCompliant IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/mail/dnsrepair/dmarc_synthesizer.go` to guarantee DMARC compliance before major sender scaling.

---

### Improvement 1468: Multi-Provider DNS Provider API Provisioner (Cloudflare, Route53, DigitalOcean)
- **Title:** Multi-Provider DNS Provider API Provisioner (Cloudflare, Route53, DigitalOcean)
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Autonomously commit repaired SPF, DKIM, and DMARC TXT records directly to authoritative upstream DNS providers via unified REST/gRPC client adapters, with pre-mutation backup snapshots and rollbacks.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Provider execution timeline rendering API request status, zone ID resolution, atomic UPSERT confirmation, and diff preview in Slate (`#64748B`) and Lime (`#84CC16`).
  - *Web UI:* Provider connection status cards with credential validation indicators and mutation audit logs.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns apply-records --domain=<domain> --provider=<cloudflare|route53|digitalocean> [--backup] [--dry-run]`
  - Flags: `--domain` (string, domain name), `--provider` (string, DNS cloud provider), `--backup` (bool, snapshot current records before mutating), `--dry-run` (bool, mock execution).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsProviderMutationLog struct {
        DomainName     string `json:"domainName"`
        ProviderType   string `json:"providerType"`
        RecordType     string `json:"recordType"`
        RecordName     string `json:"recordName"`
        PreviousValue  string `json:"previousValue"`
        NewValue       string `json:"newValue"`
        ExecutedAt     int64  `json:"executedAt"`
        isMutationDone bool   `json:"isMutationDone"`
        hasBackupSaved bool   `json:"hasBackupSaved"`
        isRolledBack   bool   `json:"isRolledBack"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsMutations (
        DomainName TEXT NOT NULL,
        ProviderType TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        RecordName TEXT NOT NULL,
        PreviousValue TEXT NOT NULL,
        NewValue TEXT NOT NULL,
        ExecutedAt INTEGER NOT NULL,
        isMutationDone INTEGER NOT NULL CHECK (isMutationDone IN (0, 1)),
        hasBackupSaved INTEGER NOT NULL CHECK (hasBackupSaved IN (0, 1)),
        isRolledBack INTEGER NOT NULL CHECK (isRolledBack IN (0, 1)),
        PRIMARY KEY (DomainName, RecordName, ExecutedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/dnsprovider/adapter.go` with strict error bubble-up to top-level CLI commands.

---

### Improvement 1469: Global DNS Anycast Propagation Radar & TTL Expiration Monitor
- **Title:** Global DNS Anycast Propagation Radar & TTL Expiration Monitor
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Track global multi-region propagation across 24 anycast edge locations (North America, Europe, Asia-Pacific, Latin America, Africa) after a record update, measuring cache TTL decay and reporting real-time propagation convergence.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* ASCII world map or regional grid with blinking status indicators (Amber for pending TTL, Green for converged) and a global percentage convergence progress meter.
  - *Web UI:* Interactive Leaflet/SVG globe displaying point-in-time regional probe responses.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns propagation-radar --domain=<domain> --record-type=TXT [--expected-val=<val>] [--watch]`
  - Flags: `--domain` (string, domain name), `--record-type` (string, RR type), `--expected-val` (string, expected substring), `--watch` (bool, poll until 100% convergence).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsPropagationStatus struct {
        DomainName       string            `json:"domainName"`
        RecordType       string            `json:"recordType"`
        ExpectedValue    string            `json:"expectedValue"`
        ConvergedPct     int               `json:"convergedPct"`
        RegionalStatus   map[string]string `json:"regionalStatus"`
        LastPolledAt     int64             `json:"lastPolledAt"`
        isFullyPropagated bool             `json:"isFullyPropagated"`
        hasStaleCaches   bool              `json:"hasStaleCaches"`
        isPollingActive  bool              `json:"isPollingActive"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsPropagation (
        DomainName TEXT NOT NULL,
        RecordType TEXT NOT NULL,
        ConvergedPct INTEGER NOT NULL,
        RegionalStatus TEXT NOT NULL,
        LastPolledAt INTEGER NOT NULL,
        isFullyPropagated INTEGER NOT NULL CHECK (isFullyPropagated IN (0, 1)),
        hasStaleCaches INTEGER NOT NULL CHECK (hasStaleCaches IN (0, 1)),
        PRIMARY KEY (DomainName, RecordType, LastPolledAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded in `pkg/mail/dnscheck/propagation.go` to gate email sending campaigns until 100% DNS convergence.

---

### Improvement 1470: Continuous DNS Drift Sentinel & Automated Rollback Circuit Breaker
- **Title:** Continuous DNS Drift Sentinel & Automated Rollback Circuit Breaker
- **Cluster:** Cluster 147: Autonomous SPF/DKIM/DMARC Record Repair & DNS Over HTTPS Sentinel
- **Goal:** Run a lightweight background sentinel daemon that continuously verifies active DNS records against known valid baseline signatures, detecting unauthorized third-party tampering, zone overrides, or unintended registrar resets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Compact status banner with lock icon, SHA-256 baseline checksum, drift delta count, and a glowing crimson alert banner when drift is identified.
  - *Web UI:* Real-time event log highlighting unexpected record mutations with an instant "Restore Baseline" action.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup dns drift-daemon --domain=<domain> [--check-interval=60s] [--auto-rollback]`
  - Flags: `--domain` (string, domain to monitor), `--check-interval` (duration, interval between checks), `--auto-rollback` (bool, automatically re-apply baseline if drift detected).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DnsDriftEvent struct {
        DomainName      string `json:"domainName"`
        BaselineHash    string `json:"baselineHash"`
        CurrentHash     string `json:"currentHash"`
        DriftDetails    string `json:"driftDetails"`
        DetectedAt      int64  `json:"detectedAt"`
        hasDriftOccurred bool  `json:"hasDriftOccurred"`
        isAutoRestored  bool   `json:"isAutoRestored"`
        isAlertDispatched bool `json:"isAlertDispatched"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDnsDriftLogs (
        DomainName TEXT NOT NULL,
        BaselineHash TEXT NOT NULL,
        CurrentHash TEXT NOT NULL,
        DriftDetails TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        hasDriftOccurred INTEGER NOT NULL CHECK (hasDriftOccurred IN (0, 1)),
        isAutoRestored INTEGER NOT NULL CHECK (isAutoRestored IN (0, 1)),
        isAlertDispatched INTEGER NOT NULL CHECK (isAlertDispatched IN (0, 1)),
        PRIMARY KEY (DomainName, DetectedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Configured in `pkg/mail/dnscheck/drift_daemon.go` with autonomous circuit breaking across all warm-up dispatches.

---

## 3. Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD (Items 1471–1480)

### Improvement 1471: Granular RFC 3463/5248 Enhanced Mail System Status Code Decoder
- **Title:** Granular RFC 3463/5248 Enhanced Mail System Status Code Decoder
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Parse and dissect standard 3-digit SMTP status codes and hierarchical 3-part enhanced status codes (e.g., `5.7.1` Security/Policy, `5.1.1` Bad Destination Mailbox, `4.4.2` Bad Connection), classifying failures into actionable taxonomy buckets.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Color-coded category tags (Permanent/Red, Transient/Yellow, Success/Green), tabular breakdown of Class, Subject, and Detail semantics, plus recommended remediation steps.
  - *Web UI:* Interactive bounce explorer with searchable RFC code dictionary and occurrence frequency sparklines.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce decode-status --code=5.7.26 [--raw-response=<text>] [--json]`
  - Flags: `--code` (string, enhanced status code e.g. 5.1.1), `--raw-response` (string, raw SMTP response string), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type EnhancedStatusCodeDetails struct {
        StatusCode        string `json:"statusCode"`
        ClassType         string `json:"classType"`
        SubjectType       string `json:"subjectType"`
        DetailType        string `json:"detailType"`
        Description       string `json:"description"`
        RecommendedAction string `json:"recommendedAction"`
        isPermanent       bool   `json:"isPermanent"`
        isPolicyRelated   bool   `json:"isPolicyRelated"`
        hasActionableFix  bool   `json:"hasActionableFix"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceClassification (
        StatusCode TEXT NOT NULL PRIMARY KEY,
        ClassType TEXT NOT NULL,
        SubjectType TEXT NOT NULL,
        DetailType TEXT NOT NULL,
        Description TEXT NOT NULL,
        RecommendedAction TEXT NOT NULL,
        isPermanent INTEGER NOT NULL CHECK (isPermanent IN (0, 1)),
        isPolicyRelated INTEGER NOT NULL CHECK (isPolicyRelated IN (0, 1)),
        hasActionableFix INTEGER NOT NULL CHECK (hasActionableFix IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Core utility in `pkg/mail/bounce/rfc_decoder.go` used by all delivery tracking pipelines.

---

### Improvement 1472: DSN Multipart Report MIME Parser & Diagnostic Header Extractor
- **Title:** DSN Multipart Report MIME Parser & Diagnostic Header Extractor
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Ingest incoming Delivery Status Notification (DSN) emails adhering to RFC 3464/6522, unpacking multipart/report MIME boundaries, extracting the `message/delivery-status` body part, and capturing original message headers (`Message-ID`, `X-Campaign-ID`).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Structured report pane displaying Reporting-MTA, Final-Recipient, Diagnostic-Code, and Action fields with collapsible raw MIME view.
  - *Web UI:* Split-pane email viewer with extracted diagnostic headers on the left and full MIME tree on the right.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce parse-dsn --file=<path> [--extract-headers] [--json]`
  - Flags: `--file` (string, path to .eml file), `--extract-headers` (bool, parse original message headers), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type DsnReportPayload struct {
        ReportingMta     string `json:"reportingMta"`
        FinalRecipient   string `json:"finalRecipient"`
        OriginalMessageId string `json:"originalMessageId"`
        DiagnosticCode   string `json:"diagnosticCode"`
        ActionTaken      string `json:"actionTaken"`
        ParsedAt         int64  `json:"parsedAt"`
        isHardBounce     bool   `json:"isHardBounce"`
        hasOriginalHeaders bool `json:"hasOriginalHeaders"`
        isSpamSignal     bool   `json:"isSpamSignal"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbDsnEvents (
        OriginalMessageId TEXT NOT NULL,
        FinalRecipient TEXT NOT NULL,
        ReportingMta TEXT NOT NULL,
        DiagnosticCode TEXT NOT NULL,
        ActionTaken TEXT NOT NULL,
        ParsedAt INTEGER NOT NULL,
        isHardBounce INTEGER NOT NULL CHECK (isHardBounce IN (0, 1)),
        hasOriginalHeaders INTEGER NOT NULL CHECK (hasOriginalHeaders IN (0, 1)),
        isSpamSignal INTEGER NOT NULL CHECK (isSpamSignal IN (0, 1)),
        PRIMARY KEY (OriginalMessageId, FinalRecipient)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/bounce/dsn_parser.go` and triggered on IMAP bounce mailbox sync.

---

### Improvement 1473: Authentication-Results & X-Spam-Status Forensic Header Analyzer
- **Title:** Authentication-Results & X-Spam-Status Forensic Header Analyzer
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Analyze complex incoming RFC 8601 `Authentication-Results` headers and vendor-specific spam headers (`X-Spam-Status`, `X-Google-Smtp-Source`, `X-Forefront-Antispam-Report`), isolating specific SPF, DKIM, DMARC, ARC, and Bayes filter sub-scores.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Badge grid rendering auth checks (`spf=pass`, `dkim=pass`, `dmarc=pass`) with individual failure reasons highlighted in High-Contrast Orange (`#F97316`).
  - *Web UI:* Radar chart showing individual authentication score dimensions across major ISPs.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce analyze-auth-headers --raw-headers=<path|string> [--json]`
  - Flags: `--raw-headers` (string, file path or raw string containing email headers), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type AuthHeaderForensics struct {
        MessageId       string  `json:"messageId"`
        SpfResult       string  `json:"spfResult"`
        DkimResult      string  `json:"dkimResult"`
        DmarcResult     string  `json:"dmarcResult"`
        ArcResult       string  `json:"arcResult"`
        SpamScore       float64 `json:"spamScore"`
        AnalyzedAt      int64   `json:"analyzedAt"`
        isFullyAligned  bool    `json:"isFullyAligned"`
        hasDkimPass     bool    `json:"hasDkimPass"`
        isSpamThreshold bool    `json:"isSpamThreshold"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbAuthHeaderLogs (
        MessageId TEXT NOT NULL PRIMARY KEY,
        SpfResult TEXT NOT NULL,
        DkimResult TEXT NOT NULL,
        DmarcResult TEXT NOT NULL,
        ArcResult TEXT NOT NULL,
        SpamScore REAL NOT NULL,
        AnalyzedAt INTEGER NOT NULL,
        isFullyAligned INTEGER NOT NULL CHECK (isFullyAligned IN (0, 1)),
        hasDkimPass INTEGER NOT NULL CHECK (hasDkimPass IN (0, 1)),
        isSpamThreshold INTEGER NOT NULL CHECK (isSpamThreshold IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Located in `pkg/mail/bounce/header_analyzer.go` for real-time seed inbox verification.

---

### Improvement 1474: Offline Naive Bayes Spam Trigger Keyword & Token Probabilistic Scorer
- **Title:** Offline Naive Bayes Spam Trigger Keyword & Token Probabilistic Scorer
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Execute a zero-dependency, local offline Naive Bayes classifier over outbound email subject lines and body copy, calculating probabilistic spam scores, identifying trigger tokens ("guaranteed", "free urgent cash"), and recommending lexical replacements.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive content linter displaying highlighted trigger words with probability badges, overall spam risk gauge (0.0 to 1.0), and suggested synonym alternatives.
  - *Web UI:* Side-by-side rich editor highlighting risky tokens in real-time with auto-rewrite suggestions.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup content lint-spam --subject=<str> --body=<path|str> [--threshold=0.3] [--json]`
  - Flags: `--subject` (string, email subject), `--body` (string, body text or file path), `--threshold` (float, maximum acceptable spam probability), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpamContentScore struct {
        SubjectScore     float64           `json:"subjectScore"`
        BodyScore        float64           `json:"bodyScore"`
        CompositeScore   float64           `json:"compositeScore"`
        FlaggedTokens    []string          `json:"flaggedTokens"`
        TokenProbMap     map[string]float64`json:"tokenProbMap"`
        isSafeToSend     bool              `json:"isSafeToSend"`
        hasTriggerWords  bool              `json:"hasTriggerWords"`
        isThresholdAlert bool              `json:"isThresholdAlert"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationSpamLexicon (
        Token TEXT NOT NULL PRIMARY KEY,
        SpamCount INTEGER NOT NULL,
        HamCount INTEGER NOT NULL,
        SpamProbability REAL NOT NULL,
        isFlaggedTrigger INTEGER NOT NULL CHECK (isFlaggedTrigger IN (0, 1)),
        hasReplacement INTEGER NOT NULL CHECK (hasReplacement IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded in `pkg/mail/content/bayes.go` as a pre-send validation step.

---

### Improvement 1475: Real-Time Interactive Terminal Bounce Forensics HUD (Bubble Tea v2)
- **Title:** Real-Time Interactive Terminal Bounce Forensics HUD (Bubble Tea v2)
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Author an interactive full-screen Bubble Tea v2 terminal HUD providing real-time visibility into incoming bounces, live stream logs, categorization breakdown charts, domain health meters, and instant manual quarantine controls.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* 3-pane layout: Left pane for live bounce stream, top right for status code distribution chart, bottom right for selected bounce forensic inspector with keyboard navigation (`j`/`k`, `q` to quarantine, `enter` for full MIME details).
  - *Web UI:* WebSocket-driven live event stream with filtering by ISP, status code, and sender domain.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce hud [--refresh-interval=1s] [--domain=<domain>]`
  - Flags: `--refresh-interval` (duration, HUD update frequency), `--domain` (string, optional domain filter).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/cli/tui/types.go`):
    ```go
    type BounceHudViewState struct {
        ActiveFilterDomain string   `json:"activeFilterDomain"`
        TotalBouncesCount  int      `json:"totalBouncesCount"`
        HardBounceCount    int      `json:"hardBounceCount"`
        SoftBounceCount    int      `json:"softBounceCount"`
        SelectedBounceId   string   `json:"selectedBounceId"`
        isLiveStreaming    bool     `json:"isLiveStreaming"`
        hasQuarantineAlert bool     `json:"hasQuarantineAlert"`
        isDetailViewOpen   bool     `json:"isDetailViewOpen"`
    }
    ```
  - SQLite DDL: Queries view over `RepodbDsnEvents` and `RepodbBounceClassification`.
- **Antigravity / GitMap Integration Point:** CLI entrypoint in `pkg/cli/tui/bounce_hud.go` utilizing standard Lipgloss theme tokens.

---

### Improvement 1476: Automated Mailbox Quarantine Circuit Breaker & Throttle Fallback Daemon
- **Title:** Automated Mailbox Quarantine Circuit Breaker & Throttle Fallback Daemon
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Continuously calculate rolling bounce and spam-complaint rates over sliding time windows (1h, 24h); immediately trip an automated circuit breaker to quarantine offending mailboxes or throttle dispatch rates when thresholds are exceeded.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* High-visibility banner in Red (`#DC2626`) showing "CIRCUIT BREAKER TRIPPED", displaying current error rate vs limit, affected mailbox ID, and cool-down timer.
  - *Web UI:* Emergency lock banner across dashboard with single-click manual override and diagnostic summary.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup quarantine status [--mailbox=<email>] [--reset] [--json]`
  - Flags: `--mailbox` (string, target mailbox address), `--reset` (bool, clear quarantine status), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type MailboxQuarantineState struct {
        MailboxEmail     string  `json:"mailboxEmail"`
        CurrentBounceRate float64 `json:"currentBounceRate"`
        ThresholdLimit   float64 `json:"thresholdLimit"`
        QuarantinedAt    int64   `json:"quarantinedAt"`
        CooldownUntil    int64   `json:"cooldownUntil"`
        TripReason       string  `json:"tripReason"`
        isQuarantined    bool    `json:"isQuarantined"`
        hasAutoReset     bool    `json:"hasAutoReset"`
        isAlertForwarded bool    `json:"isAlertForwarded"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbMailboxQuarantine (
        MailboxEmail TEXT NOT NULL PRIMARY KEY,
        CurrentBounceRate REAL NOT NULL,
        ThresholdLimit REAL NOT NULL,
        QuarantinedAt INTEGER NOT NULL,
        CooldownUntil INTEGER NOT NULL,
        TripReason TEXT NOT NULL,
        isQuarantined INTEGER NOT NULL CHECK (isQuarantined IN (0, 1)),
        hasAutoReset INTEGER NOT NULL CHECK (hasAutoReset IN (0, 1)),
        isAlertForwarded INTEGER NOT NULL CHECK (isAlertForwarded IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Enforced in `pkg/mail/throttle/circuit_breaker.go` before any mail dispatch job.

---

### Improvement 1477: SIEM Webhook Event Forwarder (Splunk, Elastic, Datadog) for Bounce Incidents
- **Title:** SIEM Webhook Event Forwarder (Splunk, Elastic, Datadog) for Bounce Incidents
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Asynchronously dispatch structured JSON telemetry events for all critical bounce events, policy violations, and quarantine trips to enterprise SIEM and observability platforms (Datadog, Splunk HEC, Elasticsearch, generic HTTP webhooks).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Forwarder status line reporting queue backlog depth, HTTP delivery success rate, average endpoint latency, and last retry attempt timestamp.
  - *Web UI:* Integration configuration drawer with endpoint URL, auth token masking, and "Test Webhook" trigger.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup siem configure --endpoint=<url> --provider=<datadog|splunk|elastic> [--token=<secret>] [--test]`
  - Flags: `--endpoint` (string, webhook URL), `--provider` (string, SIEM platform type), `--token` (string, bearer token), `--test` (bool, send synthetic test event).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SiemWebhookConfig struct {
        ProviderType     string `json:"providerType"`
        EndpointUrl      string `json:"endpointUrl"`
        AuthTokenMasked  string `json:"authTokenMasked"`
        BatchSize        int    `json:"batchSize"`
        LastForwardedAt  int64  `json:"lastForwardedAt"`
        isEnabled        bool   `json:"isEnabled"`
        hasTlsValidation bool   `json:"hasTlsValidation"`
        isQueueHealthy   bool   `json:"isQueueHealthy"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationSiemConfig (
        ProviderType TEXT NOT NULL PRIMARY KEY,
        EndpointUrl TEXT NOT NULL,
        AuthToken TEXT NOT NULL,
        BatchSize INTEGER NOT NULL,
        LastForwardedAt INTEGER NOT NULL,
        isEnabled INTEGER NOT NULL CHECK (isEnabled IN (0, 1)),
        hasTlsValidation INTEGER NOT NULL CHECK (hasTlsValidation IN (0, 1)),
        isQueueHealthy INTEGER NOT NULL CHECK (isQueueHealthy IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/bounce/siem_forwarder.go` with background worker pool and retry queue.

---

### Improvement 1478: In-Memory Bloom Filter Suppression List with Sub-Millisecond Pre-Flight Checks
- **Title:** In-Memory Bloom Filter Suppression List with Sub-Millisecond Pre-Flight Checks
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Maintain an ultra-fast, in-memory Bloom filter backed by persistent SQLite storage, enabling sub-millisecond pre-flight checks against global suppression lists (unsubscribes, hard bounces, complaint emails) prior to SMTP connection creation.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Compact filter statistics card displaying total items loaded, filter size in RAM (MB), estimated false positive probability (<0.01%), and total lookups performed.
  - *Web UI:* Searchable suppression list table with bulk CSV export and manual single-address lookup tester.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup suppression check --email=<email> [--json]`
  - `gitmap warmup suppression add --email=<email> --reason=<hardbounce|complaint|manual>`
  - Flags: `--email` (string, target address), `--reason` (string, suppression cause), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SuppressionEntry struct {
        EmailAddress string `json:"emailAddress"`
        DomainPart   string `json:"domainPart"`
        ReasonCode   string `json:"reasonCode"`
        SuppressedAt int64  `json:"suppressedAt"`
        isPermanent  bool   `json:"isPermanent"`
        hasManualAdd bool   `json:"hasManualAdd"`
        isActive     bool   `json:"isActive"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSuppressionList (
        EmailAddress TEXT NOT NULL PRIMARY KEY,
        DomainPart TEXT NOT NULL,
        ReasonCode TEXT NOT NULL,
        SuppressedAt INTEGER NOT NULL,
        isPermanent INTEGER NOT NULL CHECK (isPermanent IN (0, 1)),
        hasManualAdd INTEGER NOT NULL CHECK (hasManualAdd IN (0, 1)),
        isActive INTEGER NOT NULL CHECK (isActive IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Injected into `pkg/mail/transport/suppression_filter.go` as an unskippable safety check before recipient envelope building.

---

### Improvement 1479: Spam Trap & Honeypot Behavioral Signature Sentinel
- **Title:** Spam Trap & Honeypot Behavioral Signature Sentinel
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Identify and isolate suspected pristine and recycled spam trap addresses using behavioral telemetry (zero engagement across multiple campaigns, distinctive MX configurations, typotargeting patterns, non-standard SMTP banner responses).
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Warning indicator table with skull badge (`[TRAP RISK]`), risk confidence score (0-100%), matching signature heuristic, and recommended immediate suppression action.
  - *Web UI:* Risk distribution chart categorizing list hygiene into Clean, Suspect, and Toxic segments.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup trap-audit --list-file=<path> [--risk-threshold=80] [--auto-suppress] [--json]`
  - Flags: `--list-file` (string, path to recipient list), `--risk-threshold` (int, min risk score to flag), `--auto-suppress` (bool, automatically add flagged addresses to suppression list).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type SpamTrapRiskAssessment struct {
        EmailAddress      string  `json:"emailAddress"`
        RiskScore         int     `json:"riskScore"`
        TrapCategory      string  `json:"trapCategory"`
        MatchedHeuristic  string  `json:"matchedHeuristic"`
        AssessedAt        int64   `json:"assessedAt"`
        isHighConfidence  bool    `json:"isHighConfidence"`
        hasTypoDomain     bool    `json:"hasTypoDomain"`
        isSuppressed      bool    `json:"isSuppressed"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbSpamTrapFindings (
        EmailAddress TEXT NOT NULL PRIMARY KEY,
        RiskScore INTEGER NOT NULL,
        TrapCategory TEXT NOT NULL,
        MatchedHeuristic TEXT NOT NULL,
        AssessedAt INTEGER NOT NULL,
        isHighConfidence INTEGER NOT NULL CHECK (isHighConfidence IN (0, 1)),
        hasTypoDomain INTEGER NOT NULL CHECK (hasTypoDomain IN (0, 1)),
        isSuppressed INTEGER NOT NULL CHECK (isSuppressed IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded in `pkg/mail/bounce/trap_sentinel.go` to protect domain sender reputation.

---

### Improvement 1480: Longitudinal Predictive Bounce Trend Modeler & Pacing Advisor
- **Title:** Longitudinal Predictive Bounce Trend Modeler & Pacing Advisor
- **Cluster:** Cluster 148: Forensic Bounce Analysis Engine & Automated Classification HUD
- **Goal:** Model historical bounce data using linear regression and exponential smoothing, forecasting 7-day bounce trajectory and generating autonomous pacing recommendations to prevent mailbox blacklisting.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* ASCII sparkline graph projecting 7-day predicted bounce rate against ISP tolerance caps (0.5%), displaying recommended daily send quota adjustments.
  - *Web UI:* Trend forecast area chart with confidence intervals (95%) and historical inflection markers.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup bounce forecast --domain=<domain> [--horizon-days=7] [--json]`
  - Flags: `--domain` (string, domain name), `--horizon-days` (int, forecast horizon in days, default 7), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/types.go`):
    ```go
    type BounceTrendForecast struct {
        DomainName          string    `json:"domainName"`
        ProjectedBounceRate float64   `json:"projectedBounceRate"`
        ForecastTrendSlope  float64   `json:"forecastTrendSlope"`
        RecommendedMaxSends int       `json:"recommendedMaxSends"`
        GeneratedAt         int64     `json:"generatedAt"`
        isRateEscalating    bool      `json:"isRateEscalating"`
        hasPacingReduction  bool      `json:"hasPacingReduction"`
        isSafeTrajectory    bool      `json:"isSafeTrajectory"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbBounceForecast (
        DomainName TEXT NOT NULL,
        ProjectedBounceRate REAL NOT NULL,
        ForecastTrendSlope REAL NOT NULL,
        RecommendedMaxSends INTEGER NOT NULL,
        GeneratedAt INTEGER NOT NULL,
        isRateEscalating INTEGER NOT NULL CHECK (isRateEscalating IN (0, 1)),
        hasPacingReduction INTEGER NOT NULL CHECK (hasPacingReduction IN (0, 1)),
        isSafeTrajectory INTEGER NOT NULL CHECK (isSafeTrajectory IN (0, 1)),
        PRIMARY KEY (DomainName, GeneratedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Integrated into `pkg/mail/throttle/pacing_advisor.go` to dynamically adjust daily warmup targets.

---

## 4. Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD (Items 1481–1490)

### Improvement 1481: Vector-Clock CRDT State Synchronizer for Split-DB Replicas
- **Title:** Vector-Clock CRDT State Synchronizer for Split-DB Replicas
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Implement causality-tracking vector clocks and state-based Conflict-Free Replicated Data Types (CRDTs) to synchronize split SQLite databases (`gitmap.db`, `installation.db`, `repodb/pipeline.db`) across multi-node worker fleets without distributed deadlocks.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Replica node matrix displaying node ID, vector clock array `[node1:142, node2:139]`, replication lag in ms, and sync state in Cyan (`#06B6D4`).
  - *Web UI:* Node topology graph with animated sync pulses and causal dependency edges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db sync-status [--cluster] [--verbose] [--json]`
  - Flags: `--cluster` (bool, query all nodes in cluster), `--verbose` (bool, show raw vector clock values), `--json` (bool, structured JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ReplicaVectorClockState struct {
        NodeId         string         `json:"nodeId"`
        DatabaseName   string         `json:"databaseName"`
        ClockValues    map[string]uint64 `json:"clockValues"`
        ReplicationLagMs int64        `json:"replicationLagMs"`
        SyncedAt       int64          `json:"syncedAt"`
        isSynchronized bool           `json:"isSynchronized"`
        hasCausalDrift bool           `json:"hasCausalDrift"`
        isLeaderNode   bool           `json:"isLeaderNode"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapClusterVectorClocks (
        NodeId TEXT NOT NULL,
        DatabaseName TEXT NOT NULL,
        ClockValues TEXT NOT NULL,
        ReplicationLagMs INTEGER NOT NULL,
        SyncedAt INTEGER NOT NULL,
        isSynchronized INTEGER NOT NULL CHECK (isSynchronized IN (0, 1)),
        hasCausalDrift INTEGER NOT NULL CHECK (hasCausalDrift IN (0, 1)),
        isLeaderNode INTEGER NOT NULL CHECK (isLeaderNode IN (0, 1)),
        PRIMARY KEY (NodeId, DatabaseName)
    );
    ```
- **Antigravity / GitMap Integration Point:** Core replication layer in `pkg/db/cluster/vector_clock.go`.

---

### Improvement 1482: Encrypted WAL Frame Streaming Transport over Mutual TLS 1.3
- **Title:** Encrypted WAL Frame Streaming Transport over Mutual TLS 1.3
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Stream raw SQLite Write-Ahead Log (WAL) frames incrementally over mTLS 1.3 gRPC transport channels to peer nodes, enabling sub-second replica updates with end-to-end ChaCha20-Poly1305 frame encryption.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Live bandwidth meter rendering streaming frames/sec, compression ratio, TLS cipher suite, and socket buffer queue utilization.
  - *Web UI:* Real-time streaming chart indicating WAL frame transmission throughput and latency per peer.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db stream-wal --peer=<addr:port> [--db=repodb/pipeline.db] [--tls-cert=<path>]`
  - Flags: `--peer` (string, peer network address), `--db` (string, database file to stream), `--tls-cert` (string, client certificate path).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type WalStreamFrameMeta struct {
        FrameIndex      uint32 `json:"frameIndex"`
        PageNumber      uint32 `json:"pageNumber"`
        DatabaseFile    string `json:"databaseFile"`
        CompressedBytes int    `json:"compressedBytes"`
        ChecksumSha256  string `json:"checksumSha256"`
        StreamedAt      int64  `json:"streamedAt"`
        isCommitFrame   bool   `json:"isCommitFrame"`
        hasValidCrypto  bool   `json:"hasValidCrypto"`
        isAckReceived   bool   `json:"isAckReceived"`
    }
    ```
  - SQLite DDL: Streamed directly to disk WAL files without intermediate SQLite table overhead.
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/db/cluster/wal_streamer.go`.

---

### Improvement 1483: Page-Level Merkle Tree Shard Integrity Auditor
- **Title:** Page-Level Merkle Tree Shard Integrity Auditor
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Construct an in-memory Merkle tree of 4096-byte database pages for each SQLite shard, quickly comparing root hashes with peers to isolate corrupted or desynchronized pages without transmitting the entire database.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Visual tree structure with green/red page hash nodes, showing root hash, total page count, and desynchronized page range `[Pages 42-45 DIRTY]`.
  - *Web UI:* Interactive disk page block map with color-coded integrity status per page.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db verify-integrity --db=<path> [--peer=<node>] [--repair] [--json]`
  - Flags: `--db` (string, database path), `--peer` (string, peer node to verify against), `--repair` (bool, auto-fetch desynced pages), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type MerkleIntegrityAudit struct {
        DatabasePath     string   `json:"databasePath"`
        RootMerkleHash   string   `json:"rootMerkleHash"`
        TotalPagesCount  int      `json:"totalPagesCount"`
        MismatchedPages  []uint32 `json:"mismatchedPages"`
        AuditedAt        int64    `json:"auditedAt"`
        isHashConsistent bool     `json:"isHashConsistent"`
        hasCorruptedPages bool    `json:"hasCorruptedPages"`
        isAutoRepaired   bool     `json:"isAutoRepaired"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapDbIntegrityLogs (
        DatabasePath TEXT NOT NULL,
        RootMerkleHash TEXT NOT NULL,
        TotalPagesCount INTEGER NOT NULL,
        MismatchedPages TEXT NOT NULL,
        AuditedAt INTEGER NOT NULL,
        isHashConsistent INTEGER NOT NULL CHECK (isHashConsistent IN (0, 1)),
        hasCorruptedPages INTEGER NOT NULL CHECK (hasCorruptedPages IN (0, 1)),
        isAutoRepaired INTEGER NOT NULL CHECK (isAutoRepaired IN (0, 1)),
        PRIMARY KEY (DatabasePath, AuditedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Located in `pkg/db/cluster/merkle_audit.go`.

---

### Improvement 1484: Distributed Two-Phase Epoch Commit Coordinator for Split-DB
- **Title:** Distributed Two-Phase Epoch Commit Coordinator for Split-DB
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Coordinate atomic multi-database transactions across nodes using a lightweight Two-Phase Commit (Prepare/Commit) protocol grouped into discrete epoch windows, guaranteeing ACID guarantees across sharded SQLite instances.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Two-phase commit progress bar showing Phase 1 (Prepare vote quorum) and Phase 2 (Global Commit broadcast) with latency counters per participating node.
  - *Web UI:* Transaction sequence diagram tracing coordinator prepare requests and participant acknowledgments.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db epoch status [--current] [--json]`
  - Flags: `--current` (bool, display active epoch state), `--json` (bool, output JSON).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type EpochCommitSession struct {
        EpochId          uint64   `json:"epochId"`
        CoordinatorNode  string   `json:"coordinatorNode"`
        ParticipantNodes []string `json:"participantNodes"`
        PreparedVotes    int      `json:"preparedVotes"`
        CommitDurationMs int64    `json:"commitDurationMs"`
        CommittedAt      int64    `json:"committedAt"`
        isPrepareQuorum  bool     `json:"isPrepareQuorum"`
        hasRollbackFired bool     `json:"hasRollbackFired"`
        isFullyCommitted bool     `json:"isFullyCommitted"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapClusterEpochs (
        EpochId INTEGER NOT NULL PRIMARY KEY,
        CoordinatorNode TEXT NOT NULL,
        ParticipantNodes TEXT NOT NULL,
        PreparedVotes INTEGER NOT NULL,
        CommitDurationMs INTEGER NOT NULL,
        CommittedAt INTEGER NOT NULL,
        isPrepareQuorum INTEGER NOT NULL CHECK (isPrepareQuorum IN (0, 1)),
        hasRollbackFired INTEGER NOT NULL CHECK (hasRollbackFired IN (0, 1)),
        isFullyCommitted INTEGER NOT NULL CHECK (isFullyCommitted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/db/cluster/epoch_coordinator.go`.

---

### Improvement 1485: Non-Blocking Hot VACUUM Compaction & Storage Reclamation Engine
- **Title:** Non-Blocking Hot VACUUM Compaction & Storage Reclamation Engine
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Perform non-blocking online database compaction using SQLite `VACUUM INTO` to a temporary shadow file followed by atomic descriptor swapping, reclaiming unallocated freelist pages without interrupting ongoing read/write transactions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Storage reclamation visualizer displaying initial file size, reclaimed bytes, percentage reduction, and atomic file swap duration in milliseconds.
  - *Web UI:* Disk utilization historical chart with compaction event milestones.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db compact --db=<path> [--threshold-mb=100] [--force]`
  - Flags: `--db` (string, database file path), `--threshold-mb` (int, min unallocated space in MB before vacuuming), `--force` (bool, bypass threshold check).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type VacuumCompactionReport struct {
        DatabasePath      string `json:"databasePath"`
        OriginalSizeBytes int64  `json:"originalSizeBytes"`
        CompactedSizeBytes int64 `json:"compactedSizeBytes"`
        FreedBytes        int64  `json:"freedBytes"`
        DurationMs        int64  `json:"durationMs"`
        ExecutedAt        int64  `json:"executedAt"`
        isCompactionDone  bool   `json:"isCompactionDone"`
        hasFreedSpace     bool   `json:"hasFreedSpace"`
        isAtomicSwapSafe  bool   `json:"isAtomicSwapSafe"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationCompactionLogs (
        DatabasePath TEXT NOT NULL,
        OriginalSizeBytes INTEGER NOT NULL,
        CompactedSizeBytes INTEGER NOT NULL,
        FreedBytes INTEGER NOT NULL,
        DurationMs INTEGER NOT NULL,
        ExecutedAt INTEGER NOT NULL,
        isCompactionDone INTEGER NOT NULL CHECK (isCompactionDone IN (0, 1)),
        hasFreedSpace INTEGER NOT NULL CHECK (hasFreedSpace IN (0, 1)),
        isAtomicSwapSafe INTEGER NOT NULL CHECK (isAtomicSwapSafe IN (0, 1)),
        PRIMARY KEY (DatabasePath, ExecutedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Core compaction routine in `pkg/db/cluster/vacuum.go`.

---

### Improvement 1486: Real-Time Split-DB Replication Cluster HUD (Bubble Tea v2)
- **Title:** Real-Time Split-DB Replication Cluster HUD (Bubble Tea v2)
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Provide an interactive full-screen terminal cockpit monitoring distributed SQLite Split-DB replication health, showing live WAL replication streams, vector clock states, connected peer nodes, and epoch commit rates.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Top header displaying cluster topology status, center table rendering peer nodes with latency and frame lag, bottom console outputting real-time consensus events.
  - *Web UI:* Node topology live map with animated stream links and status badges.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db cluster-hud [--refresh=500ms]`
  - Flags: `--refresh` (duration, HUD redraw cadence).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/cli/tui/types.go`):
    ```go
    type ClusterHudViewState struct {
        ClusterLeaderNode string   `json:"clusterLeaderNode"`
        TotalPeerNodes    int      `json:"totalPeerNodes"`
        ActiveConnections int      `json:"activeConnections"`
        AverageLagMs      int64    `json:"averageLagMs"`
        isQuorumHealthy   bool     `json:"isQuorumHealthy"`
        hasDegradedPeers  bool     `json:"hasDegradedPeers"`
        isAutoFailoverOn  bool     `json:"isAutoFailoverOn"`
    }
    ```
  - SQLite DDL: Aggregates real-time state from `GitmapClusterVectorClocks`.
- **Antigravity / GitMap Integration Point:** Launched via `pkg/cli/tui/cluster_hud.go`.

---

### Improvement 1487: Deterministic Last-Write-Wins (LWW) Conflict Arbiter with Cryptographic Tie-Breaking
- **Title:** Deterministic Last-Write-Wins (LWW) Conflict Arbiter with Cryptographic Tie-Breaking
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Resolve conflicting concurrent row updates across multi-master replicas deterministically using millisecond physical timestamps combined with cryptographic Blake3 hash comparisons for identical timestamp tie-breaking.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Conflict resolution audit table displaying Table, PrimaryKey, WinnerNode, LoserNode, TimestampDelta, and deterministic tie-breaker hash.
  - *Web UI:* Side-by-side field diff highlighting discarded mutations and resolution rationale.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db conflicts list [--db=<path>] [--limit=50] [--json]`
  - Flags: `--db` (string, database file), `--limit` (int, max conflict logs to show), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type LwwConflictResolutionLog struct {
        TableName       string `json:"tableName"`
        RowIdentifier   string `json:"rowIdentifier"`
        WinnerNodeId    string `json:"winnerNodeId"`
        LoserNodeId     string `json:"loserNodeId"`
        TieBreakerHash  string `json:"tieBreakerHash"`
        ResolvedAt      int64  `json:"resolvedAt"`
        isConflictResolved bool `json:"isConflictResolved"`
        hasTieBreakerUsed  bool `json:"hasTieBreakerUsed"`
        isDataPreserved    bool `json:"isDataPreserved"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapConflictLogs (
        TableName TEXT NOT NULL,
        RowIdentifier TEXT NOT NULL,
        WinnerNodeId TEXT NOT NULL,
        LoserNodeId TEXT NOT NULL,
        TieBreakerHash TEXT NOT NULL,
        ResolvedAt INTEGER NOT NULL,
        isConflictResolved INTEGER NOT NULL CHECK (isConflictResolved IN (0, 1)),
        hasTieBreakerUsed INTEGER NOT NULL CHECK (hasTieBreakerUsed IN (0, 1)),
        isDataPreserved INTEGER NOT NULL CHECK (isDataPreserved IN (0, 1)),
        PRIMARY KEY (TableName, RowIdentifier, ResolvedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/db/cluster/lww_arbiter.go`.

---

### Improvement 1488: Dynamic Hash-Partition Shard Rebalancer & Migration Streamer
- **Title:** Dynamic Hash-Partition Shard Rebalancer & Migration Streamer
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Partition large email recipient and logging datasets across discrete SQLite shards using consistent rendezvous hashing, dynamically rebalancing partition boundaries as worker nodes join or depart the cluster.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Shard bucket distribution bar chart displaying row counts per shard, storage utilization, and active rebalancing progress percentage.
  - *Web UI:* Shard ring topology diagram showing virtual node allocations.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db rebalance [--target-shards=4] [--dry-run] [--json]`
  - Flags: `--target-shards` (int, target partition count), `--dry-run` (bool, compute partition shift without moving rows), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type ShardPartitionMeta struct {
        ShardId         int    `json:"shardId"`
        DatabaseFile    string `json:"databaseFile"`
        TotalRowsCount  int64  `json:"totalRowsCount"`
        HashRangeStart  uint32 `json:"hashRangeStart"`
        HashRangeEnd    uint32 `json:"hashRangeEnd"`
        LastRebalancedAt int64 `json:"lastRebalancedAt"`
        isPartitionActive bool `json:"isPartitionActive"`
        hasBalancedLoad   bool `json:"hasBalancedLoad"`
        isRebalancingNow  bool `json:"isRebalancingNow"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapShardPartitions (
        ShardId INTEGER NOT NULL PRIMARY KEY,
        DatabaseFile TEXT NOT NULL,
        TotalRowsCount INTEGER NOT NULL,
        HashRangeStart INTEGER NOT NULL,
        HashRangeEnd INTEGER NOT NULL,
        LastRebalancedAt INTEGER NOT NULL,
        isPartitionActive INTEGER NOT NULL CHECK (isPartitionActive IN (0, 1)),
        hasBalancedLoad INTEGER NOT NULL CHECK (hasBalancedLoad IN (0, 1)),
        isRebalancingNow INTEGER NOT NULL CHECK (isRebalancingNow IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Orchestrated by `pkg/db/cluster/rebalancer.go`.

---

### Improvement 1489: Point-in-Time Recovery (PITR) WAL Log Replayer & Time-Travel Debugger
- **Title:** Point-in-Time Recovery (PITR) WAL Log Replayer & Time-Travel Debugger
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Archive sequential WAL frame logs and replay them forward from a base SQLite snapshot to an exact target microsecond timestamp, enabling precise time-travel forensics and rollback of accidental database mutations.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Interactive timeline scrubber with frame index markers, allowing operators to step forward/backward through database commits and view snapshot state diffs.
  - *Web UI:* Visual commit history slider with instantaneous schema and data preview.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db pitr replay --db=<path> --target-timestamp=<iso8601|unix> [--out=<restored_path>]`
  - Flags: `--db` (string, source DB), `--target-timestamp` (string, restoration time), `--out` (string, output destination).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type PitrReplaySession struct {
        DatabasePath     string `json:"databasePath"`
        TargetTimestamp  int64  `json:"targetTimestamp"`
        ReplayedFrames   int    `json:"replayedFrames"`
        RestoredDbPath   string `json:"restoredDbPath"`
        CompletedAt      int64  `json:"completedAt"`
        isReplaySuccess  bool   `json:"isReplaySuccess"`
        hasTargetReached bool   `json:"hasTargetReached"`
        isVerifiedIntegrity bool `json:"isVerifiedIntegrity"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationPitrHistory (
        DatabasePath TEXT NOT NULL,
        TargetTimestamp INTEGER NOT NULL,
        ReplayedFrames INTEGER NOT NULL,
        RestoredDbPath TEXT NOT NULL,
        CompletedAt INTEGER NOT NULL,
        isReplaySuccess INTEGER NOT NULL CHECK (isReplaySuccess IN (0, 1)),
        hasTargetReached INTEGER NOT NULL CHECK (hasTargetReached IN (0, 1)),
        isVerifiedIntegrity INTEGER NOT NULL CHECK (isVerifiedIntegrity IN (0, 1)),
        PRIMARY KEY (DatabasePath, TargetTimestamp)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/db/cluster/pitr_replayer.go`.

---

### Improvement 1490: Zero-Downtime Rolling Schema Migration Orchestrator for Cluster Shards
- **Title:** Zero-Downtime Rolling Schema Migration Orchestrator for Cluster Shards
- **Cluster:** Cluster 149: Distributed SQLite Split-DB Multi-Master Replication HUD
- **Goal:** Execute zero-downtime, rolling DDL schema migrations across distributed SQLite databases using shadow table copying, trigger-based change capture, and atomic table swaps without locking read or write queries.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Step-by-step migration progress indicator: Create Shadow -> Backfill Rows -> Catchup Triggers -> Atomic Rename -> Cleanup, with elapsed time per step.
  - *Web UI:* Schema diff visualizer comparing active and shadow tables side by side.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap db migrate rolling --db=<path> --sql-file=<migration.sql> [--batch-size=1000]`
  - Flags: `--db` (string, target DB), `--sql-file` (string, DDL migration file), `--batch-size` (int, backfill rows per transaction).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/db/cluster/types.go`):
    ```go
    type RollingMigrationState struct {
        DatabasePath     string `json:"databasePath"`
        MigrationVersion string `json:"migrationVersion"`
        TargetTable      string `json:"targetTable"`
        BackfilledRows   int64  `json:"backfilledRows"`
        TotalRows        int64  `json:"totalRows"`
        StartedAt        int64  `json:"startedAt"`
        isMigrationDone  bool   `json:"isMigrationDone"`
        hasShadowSwapped bool   `json:"hasShadowSwapped"`
        isTriggersActive bool   `json:"isTriggersActive"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSchemaMigrations (
        DatabasePath TEXT NOT NULL,
        MigrationVersion TEXT NOT NULL,
        TargetTable TEXT NOT NULL,
        BackfilledRows INTEGER NOT NULL,
        TotalRows INTEGER NOT NULL,
        AppliedAt INTEGER NOT NULL,
        isMigrationDone INTEGER NOT NULL CHECK (isMigrationDone IN (0, 1)),
        hasShadowSwapped INTEGER NOT NULL CHECK (hasShadowSwapped IN (0, 1)),
        PRIMARY KEY (DatabasePath, MigrationVersion)
    );
    ```
- **Antigravity / GitMap Integration Point:** Integrated into `pkg/db/cluster/rolling_migrator.go`.

---

## 5. Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer (Items 1491–1500)

### Improvement 1491: Antigravity Multi-Agent Swarm Cockpit HUD (Bubble Tea v2)
- **Title:** Antigravity Multi-Agent Swarm Cockpit HUD (Bubble Tea v2)
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Deliver an interactive, terminal-native mission control cockpit for monitoring all concurrent subagents, inspecting live prompt executions, subtask leases, token budgets, and consensus states in real time.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Multi-column agent roster showing AgentID, Role, ActiveSubtask, ElapsedTime, and Status Badge (Idle/Working/Stalled) in vibrant colors (`#8B5CF6`, `#10B981`, `#EF4444`).
  - *Web UI:* Swarm grid dashboard with agent activity cards and real-time terminal output mirroring.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent swarm-hud [--refresh=1s] [--agent-type=all]`
  - Flags: `--refresh` (duration, redraw interval), `--agent-type` (string, filter by role/type).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/cli/tui/types.go`):
    ```go
    type SwarmCockpitState struct {
        ActiveAgentCount int      `json:"activeAgentCount"`
        CompletedTasks   int      `json:"completedTasks"`
        PendingTasks     int      `json:"pendingTasks"`
        SelectedAgentId  string   `json:"selectedAgentId"`
        isSwarmHealthy   bool     `json:"isSwarmHealthy"`
        hasStalledAgents bool     `json:"hasStalledAgents"`
        isMasterElect    bool     `json:"isMasterElect"`
    }
    ```
  - SQLite DDL: Aggregates real-time agent metrics from `GitmapAgentTasks`.
- **Antigravity / GitMap Integration Point:** CLI entrypoint in `pkg/cli/tui/swarm_hud.go`.

---

### Improvement 1492: Dynamic Subtask Auctioneer & Expiring Lease Scheduler
- **Title:** Dynamic Subtask Auctioneer & Expiring Lease Scheduler
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Implement an auction-based task scheduler that publishes subtasks to available subagents with time-bound leases (e.g., 60s); uncompleted tasks whose leases expire are automatically reclaimed and re-queued.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Auction block table showing SubtaskId, BidderAgent, LeaseExpiresIn (seconds countdown), and RetryCount.
  - *Web UI:* Kanban view of active leases with visual countdown bars.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent tasks list [--status=leased] [--json]`
  - Flags: `--status` (string, task filter: pending/leased/done), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type AgentSubtaskLease struct {
        TaskId          string `json:"taskId"`
        AssignedAgentId string `json:"assignedAgentId"`
        LeaseDurationS  int    `json:"leaseDurationS"`
        LeasedAt        int64  `json:"leasedAt"`
        ExpiresAt       int64  `json:"expiresAt"`
        RetryCount      int    `json:"retryCount"`
        isLeaseActive   bool   `json:"isLeaseActive"`
        hasTimedOut     bool   `json:"hasTimedOut"`
        isCompleted     bool   `json:"isCompleted"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapAgentTasks (
        TaskId TEXT NOT NULL PRIMARY KEY,
        AssignedAgentId TEXT NOT NULL,
        LeaseDurationS INTEGER NOT NULL,
        LeasedAt INTEGER NOT NULL,
        ExpiresAt INTEGER NOT NULL,
        RetryCount INTEGER NOT NULL,
        isLeaseActive INTEGER NOT NULL CHECK (isLeaseActive IN (0, 1)),
        hasTimedOut INTEGER NOT NULL CHECK (hasTimedOut IN (0, 1)),
        isCompleted INTEGER NOT NULL CHECK (isCompleted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Core scheduling daemon in `pkg/agent/swarm/auctioneer.go`.

---

### Improvement 1493: Token Budget Expenditure Governor & Context Window Comptroller
- **Title:** Token Budget Expenditure Governor & Context Window Comptroller
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Track cumulative token usage across LLM model invocations per agent, enforcing strict cost and quota ceilings, and preemptively compacting context windows when conversation trajectories exceed safe limits.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Token odometer gauge displaying prompt tokens, completion tokens, and dollar cost against daily budget ceiling in Amber/Green.
  - *Web UI:* Token expenditure breakdown by agent role and task sequence.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent budget status [--reset] [--json]`
  - Flags: `--reset` (bool, reset daily counters), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type TokenUsageComptroller struct {
        AgentId         string  `json:"agentId"`
        PromptTokens    int64   `json:"promptTokens"`
        OutputTokens    int64   `json:"outputTokens"`
        EstimatedCostUsd float64 `json:"estimatedCostUsd"`
        DailyBudgetLimit float64 `json:"dailyBudgetLimit"`
        LastUpdated     int64   `json:"lastUpdated"`
        isWithinBudget  bool    `json:"isWithinBudget"`
        hasQuotaExceeded bool   `json:"hasQuotaExceeded"`
        isCompacted     bool    `json:"isCompacted"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationTokenBudget (
        AgentId TEXT NOT NULL PRIMARY KEY,
        PromptTokens INTEGER NOT NULL,
        OutputTokens INTEGER NOT NULL,
        EstimatedCostUsd REAL NOT NULL,
        DailyBudgetLimit REAL NOT NULL,
        LastUpdated INTEGER NOT NULL,
        isWithinBudget INTEGER NOT NULL CHECK (isWithinBudget IN (0, 1)),
        hasQuotaExceeded INTEGER NOT NULL CHECK (hasQuotaExceeded IN (0, 1)),
        isCompacted INTEGER NOT NULL CHECK (isCompacted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Injected into `pkg/agent/swarm/token_governor.go`.

---

### Improvement 1494: Tarjan-Based Circular Deadlock Detection & Autonomous Worker Respawn
- **Title:** Tarjan-Based Circular Deadlock Detection & Autonomous Worker Respawn
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Construct a directed dependency graph of subagent blocking relationships, executing Tarjan's Strongly Connected Components algorithm to detect circular wait deadlocks, automatically killing stalled agents and re-leasing tasks.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Deadlock alert box rendering detected circular dependency cycle (`Agent-1 -> Agent-3 -> Agent-1`), termination status, and restart timeline.
  - *Web UI:* Interactive dependency graph with highlighted circular loops.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent deadlock scan [--auto-resolve] [--json]`
  - Flags: `--auto-resolve` (bool, kill cycles immediately), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type DeadlockDetectionEvent struct {
        CycleNodeIds     []string `json:"cycleNodeIds"`
        CulpritAgentId   string   `json:"culpritAgentId"`
        DetectedAt       int64    `json:"detectedAt"`
        hasCycleDetected bool     `json:"hasCycleDetected"`
        isWorkerKilled   bool     `json:"isWorkerKilled"`
        hasTaskRequeued  bool     `json:"hasTaskRequeued"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapDeadlockLogs (
        CulpritAgentId TEXT NOT NULL,
        CycleNodeIds TEXT NOT NULL,
        DetectedAt INTEGER NOT NULL,
        hasCycleDetected INTEGER NOT NULL CHECK (hasCycleDetected IN (0, 1)),
        isWorkerKilled INTEGER NOT NULL CHECK (isWorkerKilled IN (0, 1)),
        hasTaskRequeued INTEGER NOT NULL CHECK (hasTaskRequeued IN (0, 1)),
        PRIMARY KEY (CulpritAgentId, DetectedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/agent/swarm/deadlock_detector.go`.

---

### Improvement 1495: Sub-250ms Embedded Raft Consensus Module for Swarm Master Election
- **Title:** Sub-250ms Embedded Raft Consensus Module for Swarm Master Election
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Implement an embedded, zero-external-dependency Raft consensus engine across subagent processes, achieving sub-250ms leader elections and log replication to govern global orchestration decisions.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Raft cluster status indicator displaying current Term number, Leader ID, Node Roles (Follower/Candidate/Leader), and heartbeat heartbeat latency in Emerald.
  - *Web UI:* Cluster quorum state card with live heartbeat pulses.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap agent raft status [--json]`
  - Flags: `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/swarm/types.go`):
    ```go
    type RaftConsensusStatus struct {
        CurrentTerm      uint64   `json:"currentTerm"`
        LeaderNodeId     string   `json:"leaderNodeId"`
        NodeRole         string   `json:"nodeRole"`
        ActivePeers      []string `json:"activePeers"`
        LastHeartbeatMs  int64    `json:"lastHeartbeatMs"`
        isLeaderActive   bool     `json:"isLeaderActive"`
        hasQuorumReached bool     `json:"hasQuorumReached"`
        isElectionOpen   bool     `json:"isElectionOpen"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapRaftState (
        CurrentTerm INTEGER NOT NULL PRIMARY KEY,
        LeaderNodeId TEXT NOT NULL,
        NodeRole TEXT NOT NULL,
        ActivePeers TEXT NOT NULL,
        LastHeartbeatMs INTEGER NOT NULL,
        isLeaderActive INTEGER NOT NULL CHECK (isLeaderActive IN (0, 1)),
        hasQuorumReached INTEGER NOT NULL CHECK (hasQuorumReached IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Embedded in `pkg/agent/swarm/raft_node.go`.

---

### Improvement 1496: Closed-Loop Deliverability Feedback Governor & Autonomous Throttle Steerer
- **Title:** Closed-Loop Deliverability Feedback Governor & Autonomous Throttle Steerer
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Close the feedback loop between inbound deliverability signals (inbox placement, bounce codes, spam filter scores) and outbound dispatch concurrency, automatically modulating warmup volume without human intervention.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* PID controller feedback gauge displaying current velocity, error delta, adjustment factor (+15% or -30%), and target daily quota.
  - *Web UI:* Interactive throttle curve chart showing real-time feedback loop adjustments.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup throttle steer --domain=<domain> [--auto] [--json]`
  - Flags: `--domain` (string, domain name), `--auto` (bool, enable autonomous closed-loop adjustment), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/agent/steerer/types.go`):
    ```go
    type ThrottleSteeringState struct {
        DomainName       string  `json:"domainName"`
        CurrentHourlyCap int     `json:"currentHourlyCap"`
        TargetHourlyCap  int     `json:"targetHourlyCap"`
        ErrorIntegral    float64 `json:"errorIntegral"`
        AdjustmentRatio  float64 `json:"adjustmentRatio"`
        EvaluatedAt      int64   `json:"evaluatedAt"`
        isThrottleActive bool    `json:"isThrottleActive"`
        hasPacingReduced bool    `json:"hasPacingReduced"`
        isTargetReached  bool    `json:"isTargetReached"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbThrottleSteering (
        DomainName TEXT NOT NULL,
        CurrentHourlyCap INTEGER NOT NULL,
        TargetHourlyCap INTEGER NOT NULL,
        AdjustmentRatio REAL NOT NULL,
        EvaluatedAt INTEGER NOT NULL,
        isThrottleActive INTEGER NOT NULL CHECK (isThrottleActive IN (0, 1)),
        hasPacingReduced INTEGER NOT NULL CHECK (hasPacingReduced IN (0, 1)),
        isTargetReached INTEGER NOT NULL CHECK (isTargetReached IN (0, 1)),
        PRIMARY KEY (DomainName, EvaluatedAt)
    );
    ```
- **Antigravity / GitMap Integration Point:** Handled by `pkg/agent/steerer/pid_governor.go`.

---

### Improvement 1497: Multi-Turn Synthetic Contextual Conversation Thread Generator
- **Title:** Multi-Turn Synthetic Contextual Conversation Thread Generator
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Autonomously generate multi-turn, highly authentic email conversation threads between synthetic peer mailboxes, utilizing dynamic persona prompts, varied reply delays, and natural quotation headers to mimic genuine human interaction.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Conversation tree viewer displaying ThreadID, Personas, Subject, Message Depth (1 to 5 turns), and scheduled send intervals.
  - *Web UI:* Chat-style thread preview with reply scheduling controls.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup threads generate --topic=business-consulting [--depth=4] [--count=10]`
  - Flags: `--topic` (string, conversation domain), `--depth` (int, message turns per thread), `--count` (int, number of threads).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/synthetic/types.go`):
    ```go
    type SyntheticConversationThread struct {
        ThreadId        string   `json:"threadId"`
        SubjectLine     string   `json:"subjectLine"`
        TurnCount       int      `json:"turnCount"`
        ParticipantEmails []string `json:"participantEmails"`
        GeneratedAt     int64    `json:"generatedAt"`
        isThreadActive  bool     `json:"isThreadActive"`
        hasNaturalReply bool     `json:"hasNaturalReply"`
        isCompleted     bool     `json:"isCompleted"`
    }
    ```
  - SQLite DDL (`gitmap.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS GitmapSyntheticThreads (
        ThreadId TEXT NOT NULL PRIMARY KEY,
        SubjectLine TEXT NOT NULL,
        TurnCount INTEGER NOT NULL,
        ParticipantEmails TEXT NOT NULL,
        GeneratedAt INTEGER NOT NULL,
        isThreadActive INTEGER NOT NULL CHECK (isThreadActive IN (0, 1)),
        hasNaturalReply INTEGER NOT NULL CHECK (hasNaturalReply IN (0, 1)),
        isCompleted INTEGER NOT NULL CHECK (isCompleted IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Generated in `pkg/mail/synthetic/generator.go`.

---

### Improvement 1498: Per-ISP Reputation Vector Heatmap Visualizer (Google, Microsoft, Yahoo, Apple)
- **Title:** Per-ISP Reputation Vector Heatmap Visualizer (Google, Microsoft, Yahoo, Apple)
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Collect, aggregate, and render multi-dimensional reputation scores across major email service providers (Google Postmaster, Microsoft SNDS, Yahoo Sender Hub), producing unified ISP reputation vector heatmaps.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* 4-provider comparative matrix displaying SPF pass rate, DKIM pass rate, Spam Rate, User Complaint Rate, and Aggregate Health in High-Density Colored Blocks.
  - *Web UI:* 2D heatmap matrix with interactive ISP drill-down telemetry.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup reputation heatmap [--domain=<domain>] [--json]`
  - Flags: `--domain` (string, domain name), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/reputation/types.go`):
    ```go
    type IspReputationVector struct {
        DomainName     string  `json:"domainName"`
        IspProvider    string  `json:"ispProvider"`
        ReputationScore int    `json:"reputationScore"`
        SpamRate       float64 `json:"spamRate"`
        InboxPlacement float64 `json:"inboxPlacement"`
        UpdatedTimestamp int64 `json:"updatedTimestamp"`
        isHealthyScore bool    `json:"isHealthyScore"`
        hasRateWarning bool    `json:"hasRateWarning"`
        isDataFresh    bool    `json:"isDataFresh"`
    }
    ```
  - SQLite DDL (`repodb/pipeline.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS RepodbIspReputation (
        DomainName TEXT NOT NULL,
        IspProvider TEXT NOT NULL,
        ReputationScore INTEGER NOT NULL,
        SpamRate REAL NOT NULL,
        InboxPlacement REAL NOT NULL,
        UpdatedTimestamp INTEGER NOT NULL,
        isHealthyScore INTEGER NOT NULL CHECK (isHealthyScore IN (0, 1)),
        hasRateWarning INTEGER NOT NULL CHECK (hasRateWarning IN (0, 1)),
        isDataFresh INTEGER NOT NULL CHECK (isDataFresh IN (0, 1)),
        PRIMARY KEY (DomainName, IspProvider, UpdatedTimestamp)
    );
    ```
- **Antigravity / GitMap Integration Point:** Implemented in `pkg/mail/reputation/heatmap.go`.

---

### Improvement 1499: Autonomous Domain Warmup Graduation Gate & Readiness Certifier
- **Title:** Autonomous Domain Warmup Graduation Gate & Readiness Certifier
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Evaluate complete warm-up trajectory against strict graduation criteria (30+ consecutive days, >98% inbox placement, <0.1% bounce rate, 100% SPF/DKIM/DMARC alignment), automatically issuing a signed digital Readiness Certificate.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Certificate badge bordered in Gold (`#F59E0B`), displaying Domain, Graduation Date, Certified Daily Volume, and Cryptographic SHA-256 Audit Seal.
  - *Web UI:* Downloadable PDF readiness certificate and verification status badge.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup graduate verify --domain=<domain> [--issue-cert] [--json]`
  - Flags: `--domain` (string, domain name), `--issue-cert` (bool, generate graduation certificate if qualified), `--json` (bool, JSON output).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/mail/graduation/types.go`):
    ```go
    type DomainGraduationCertificate struct {
        DomainName          string `json:"domainName"`
        CertificateSerial   string `json:"certificateSerial"`
        CertifiedDailyVolume int   `json:"certifiedDailyVolume"`
        InboxPlacementRate  int    `json:"inboxPlacementRate"`
        GraduatedAt         int64  `json:"graduatedAt"`
        isGraduated         bool   `json:"isGraduated"`
        hasMetAllCriteria   bool   `json:"hasMetAllCriteria"`
        isCertIssued        bool   `json:"isCertIssued"`
    }
    ```
  - SQLite DDL (`installation.db`):
    ```sql
    CREATE TABLE IF NOT EXISTS InstallationGraduationCerts (
        DomainName TEXT NOT NULL PRIMARY KEY,
        CertificateSerial TEXT NOT NULL,
        CertifiedDailyVolume INTEGER NOT NULL,
        InboxPlacementRate INTEGER NOT NULL,
        GraduatedAt INTEGER NOT NULL,
        isGraduated INTEGER NOT NULL CHECK (isGraduated IN (0, 1)),
        hasMetAllCriteria INTEGER NOT NULL CHECK (hasMetAllCriteria IN (0, 1)),
        isCertIssued INTEGER NOT NULL CHECK (isCertIssued IN (0, 1))
    );
    ```
- **Antigravity / GitMap Integration Point:** Enforced in `pkg/mail/graduation/gate.go` prior to promoting domain to production dispatch status.

---

### Improvement 1500: Unified Executive Mission Telemetry HUD & Swarm Operations Dashboard
- **Title:** Unified Executive Mission Telemetry HUD & Swarm Operations Dashboard
- **Cluster:** Cluster 150: Antigravity Multi-Agent Swarm Orchestrator & Autonomous Warmup Steerer
- **Goal:** Synthesize the entire telemetry pipeline into a unified, full-screen executive mission dashboard displaying end-to-end swarm health, DNS consensus state, bounce metrics, distributed Split-DB replication status, and deliverability trajectories.
- **UI/UX Terminal & Web Specification:**
  - *Terminal (Lipgloss):* Four-quadrant executive dashboard: Top-Left (DNS & Auth Sentinel), Top-Right (Bounce Forensics & Pacing), Bottom-Left (SQLite Multi-Master Replication), Bottom-Right (Swarm Subagent Fleet).
  - *Web UI:* Real-time responsive control room view with configurable alert overlays.
- **CLI Command Signature & Flag Definitions:**
  - `gitmap warmup mission-hud [--fullscreen] [--refresh=1s]`
  - Flags: `--fullscreen` (bool, enter alternate terminal buffer), `--refresh` (duration, telemetry refresh rate).
- **Data Contract & Split-DB Schema:**
  - Go Struct (`pkg/cli/tui/types.go`):
    ```go
    type UnifiedMissionTelemetryState struct {
        TotalManagedDomains int      `json:"totalManagedDomains"`
        ActiveSwarmAgents   int      `json:"activeSwarmAgents"`
        ClusterSyncLagMs    int64    `json:"clusterSyncLagMs"`
        OverallInboxPlacement float64 `json:"overallInboxPlacement"`
        SystemUptimeSeconds int64    `json:"systemUptimeSeconds"`
        isMissionHealthy    bool     `json:"isMissionHealthy"`
        hasCriticalAlerts   bool     `json:"hasCriticalAlerts"`
        isTelemetryLive     bool     `json:"isTelemetryLive"`
    }
    ```
  - SQLite DDL: Queries top-level status across all subsystem tables in `gitmap.db`, `installation.db`, and `repodb/pipeline.db`.
- **Antigravity / GitMap Integration Point:** Terminal entrypoint in `pkg/cli/tui/mission_hud.go` combining all subsystem telemetry.

---

## 6. Verification and Validation Checklist

- [ ] All 40 improvements (1461 through 1500) defined with Title, ID, Cluster, and Goal.
- [ ] UI/UX Terminal (Lipgloss) & Web specifications authored for every item with zero CDN dependencies.
- [ ] Strict positive boolean convention (`is`, `has`) verified across all Go structs and SQLite DDL tables.
- [ ] Split-DB isolation maintained across `gitmap.db`, `installation.db`, and `repodb/pipeline.db`.
- [ ] 100% relative repository paths maintained with zero absolute paths or drive letters.
