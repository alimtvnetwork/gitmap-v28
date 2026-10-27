# Improvements 701–730: Fingerprint Morphing, Rate-Limit Mesh & SMTP Multiplexer

This document details Improvements 701 through 730 for Sequence 344 across Clusters 71, 72, and 73. All implementations strictly follow positive boolean naming (`is`/`has` only), relative repository paths, Split-DB SQLite storage (`repodb/warmup.db` and `repodb/audit.db`), and zero-external-CDN UI architectures.

---

# Cluster 71: Autonomous Envelope Fingerprint Morphing & Canonical Header Obfuscation (701–710)

---

### Improvement 701: Pseudo-Random MIME Boundary Delimiter Synthesizer
- **ID**: `IMP-701`
- **Cluster**: 71 (Envelope Fingerprint Morphing)
- **Goal**: Generates entropy-seeded, non-deterministic MIME multipart boundary strings emulating prominent MUA formats (Thunderbird, Apple Mail, Outlook, Postfix) to prevent static boundary signature matching.
- **UI/UX Terminal & Web Specification**:
  - *Terminal (Lipgloss)*: Rendered via `lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#6366F1"))`. Displays boundary syntax preview alongside entropy score meter (`[████████░░] 8.2 bits/char`).
  - *Web (Shadow DOM)*: `<mime-boundary-inspector>` custom element showing interactive boundary pattern selector and zero-CDN inline SVG entropy sparklines.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph boundary-gen [--style=<preset>] [--has-custom-seed] [--is-verbose]`
  - Flag `--style`: Target boundary template preset (`thunderbird`, `apple-mail`, `hex-epoch`, `uuid-v4`).
  - Flag `--has-custom-seed`: Enables deterministic entropy initialization for testing.
  - Flag `--is-verbose`: Prints complete boundary decomposition with character frequency counts.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MimeBoundaryTemplates (
      id TEXT PRIMARY KEY,
      presetName TEXT NOT NULL,
      prefixPattern TEXT NOT NULL,
      entropySource TEXT NOT NULL,
      delimiterLength INTEGER NOT NULL,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasTrailingHyphens INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MimeBoundaryTemplate struct {
      PresetName         string `json:"presetName"`
      PrefixPattern      string `json:"prefixPattern"`
      DelimiterLength    int    `json:"delimiterLength"`
      IsActive           bool   `json:"isActive"`
      HasTrailingHyphens bool   `json:"hasTrailingHyphens"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrates with `pkg/warmup/morph/boundary.go` and hooks into message serialization pipelines in `gitmap email warmup send`.

---

### Improvement 702: Dynamic Header Permutation Sequence Matrix
- **ID**: `IMP-702`
- **Cluster**: 71
- **Goal**: Permutes the vertical ordering of non-mandatory RFC 5322 headers within mathematically safe bounds to disrupt static machine-learning header vector classifiers.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual vertical ASCII sequence tree highlighting anchor headers (`From`, `To`, `Date`) in `#48BB78` (green) and shuffled diagnostic headers in `#ECC94B` (amber).
  - *Web*: `<header-permutation-view>` element showing interactive drag-and-drop permutation simulation with instant RFC conformity scoring.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph permute-headers --message-id=<msg-uuid> [--has-strict-anchors] [--is-dry-run]`
  - Flag `--has-strict-anchors`: Fixes `From`, `To`, `Date`, and `Subject` at top positions while permuting remaining headers.
  - Flag `--is-dry-run`: Previews the permuted header sequence without writing to the outbox database.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HeaderPermutations (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      originalSequence TEXT NOT NULL, -- JSON array of header keys
      permutedSequence TEXT NOT NULL, -- JSON array of permuted header keys
      permutationDistance INTEGER NOT NULL,
      isCompliant INTEGER NOT NULL DEFAULT 1,
      hasAnchorPreserved INTEGER NOT NULL DEFAULT 1,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HeaderPermutationRecord struct {
      MessageID           string `json:"messageId"`
      PermutationDistance int    `json:"permutationDistance"`
      IsCompliant         bool   `json:"isCompliant"`
      HasAnchorPreserved  bool   `json:"hasAnchorPreserved"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked by `pkg/warmup/morph/permutator.go` immediately prior to message signing.

---

### Improvement 703: RFC 5322 Header Name Casing Variation Engine
- **ID**: `IMP-703`
- **Cluster**: 71
- **Goal**: Applies controlled capitalization variance across case-insensitive MIME header field names (`MIME-Version`, `Mime-Version`, `mime-version`) to emulate diverse sending agents without violating RFC standards.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Side-by-side terminal diff showing standard canonical casing versus client-emulated casing with Lipgloss bold cyan highlights (`#38B2AC`).
  - *Web*: `<casing-matrix-studio>` component with client profile dropdown (e.g. "iOS Mail 16", "Outlook 365", "Postbox").
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph apply-casing --client-profile=<profile> [--is-normalized] [--has-audit-log]`
  - Flag `--client-profile`: Selects casing signature of specific MUA profiles.
  - Flag `--is-normalized`: Enforces RFC canonical title casing across all fields.
  - Flag `--has-audit-log`: Records casing mutation events in `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS HeaderCasingProfiles (
      id TEXT PRIMARY KEY,
      clientProfileName TEXT NOT NULL,
      headerKey TEXT NOT NULL,
      casingVariant TEXT NOT NULL,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasCaseSensitivityIssue INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type HeaderCasingProfile struct {
      ClientProfileName       string `json:"clientProfileName"`
      HeaderKey               string `json:"headerKey"`
      CasingVariant           string `json:"casingVariant"`
      IsActive                bool   `json:"isActive"`
      HasCaseSensitivityIssue bool   `json:"hasCaseSensitivityIssue"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Registered in `pkg/warmup/morph/casing.go` and verified during DKIM pre-sign stages.

---

### Improvement 704: Message-ID Entropy Injector & Domain Key Shuffler
- **ID**: `IMP-704`
- **Cluster**: 71
- **Goal**: Injects cryptographic noise and timestamp jitter into `Message-ID` headers to prevent statistical tracking while anchoring right-hand domain validation to authorized SPF/DKIM roots.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Displays formatted Message-ID breakdown with entropy component breakdown: `<[timestamp-hex].[entropy-base32]@[signing-domain]>`.
  - *Web*: `<message-id-visualizer>` displaying bit distribution and domain reputation alignment status.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph generate-msgid --domain=<domain> [--has-monotonic-clock] [--is-hex-encoded]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS GeneratedMessageIds (
      id TEXT PRIMARY KEY,
      domainName TEXT NOT NULL,
      localPartEntropy TEXT NOT NULL,
      fullHeaderValue TEXT NOT NULL,
      isDelivered INTEGER NOT NULL DEFAULT 0,
      hasReplied INTEGER NOT NULL DEFAULT 0,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type GeneratedMessageIdRecord struct {
      DomainName       string `json:"domainName"`
      LocalPartEntropy string `json:"localPartEntropy"`
      FullHeaderValue  string `json:"fullHeaderValue"`
      IsDelivered      bool   `json:"isDelivered"`
      HasReplied       bool   `json:"hasReplied"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Bound to `pkg/warmup/envelope/msgid.go` with metrics exported to GitMap telemetry.

---

### Improvement 705: Pre-DKIM Canonicalization Scrambler & Invariant Verifier
- **ID**: `IMP-705`
- **Cluster**: 71
- **Goal**: Verifies that morphing transformations maintain strict mathematical invariance under `relaxed/relaxed` and `simple/simple` DKIM canonicalization algorithms prior to MTA transmission.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Compact dual-hash validator display using green badges (`[OK: RELAXED]`, `[OK: SIMPLE]`) and red warning badges (`[ERR: BODY_HASH_MISMATCH]`).
  - *Web*: `<dkim-invariant-checker>` with instant cryptographic SHA-256 verification and raw canonical buffer inspection.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph test-dkim-invariant --message-id=<msg-uuid> [--is-strict] [--has-relaxed-body]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DkimInvariantChecks (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      canonicalHeaderHash TEXT NOT NULL,
      canonicalBodyHash TEXT NOT NULL,
      isInvariantSatisfied INTEGER NOT NULL DEFAULT 1,
      hasSimpleSimplePassed INTEGER NOT NULL DEFAULT 1,
      hasRelaxedRelaxedPassed INTEGER NOT NULL DEFAULT 1,
      checkedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DkimInvariantResult struct {
      MessageID               string `json:"messageId"`
      IsInvariantSatisfied    bool   `json:"isInvariantSatisfied"`
      HasSimpleSimplePassed   bool   `json:"hasSimpleSimplePassed"`
      HasRelaxedRelaxedPassed bool   `json:"hasRelaxedRelaxedPassed"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/crypto/dkim_guard.go` as a mandatory pre-flight validation gate.

---

### Improvement 706: Multi-Part Epilogue & Prologue Padding Synthesizer
- **ID**: `IMP-706`
- **Cluster**: 71
- **Goal**: Injects safe, RFC-compliant whitespace or comment padding into MIME multipart prologues and epilogues to alter total byte offsets and thwart static body size profiling.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visualizes MIME structure with byte-offset indicators and prologue/epilogue filler markers styled with Lipgloss dim tokens (`#718096`).
  - *Web*: `<mime-padding-inspector>` rendering interactive boundary delimiters and customizable padding density controls.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph inject-padding --message-id=<msg-uuid> [--has-crlf-jitter] [--is-silent]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MimePaddingInjections (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      prologueByteCount INTEGER NOT NULL DEFAULT 0,
      epilogueByteCount INTEGER NOT NULL DEFAULT 0,
      isPaddingActive INTEGER NOT NULL DEFAULT 1,
      hasCrlfJitter INTEGER NOT NULL DEFAULT 1,
      injectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MimePaddingRecord struct {
      MessageID         string `json:"messageId"`
      PrologueByteCount int    `json:"prologueByteCount"`
      EpilogueByteCount int    `json:"epilogueByteCount"`
      IsPaddingActive   bool   `json:"isPaddingActive"`
      HasCrlfJitter     bool   `json:"hasCrlfJitter"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/morph/padding.go` directly before transport packaging.

---

### Improvement 707: Content-Transfer-Encoding Alternator (7bit vs 8bit vs Quoted-Printable)
- **ID**: `IMP-707`
- **Cluster**: 71
- **Goal**: Dynamically alternates text body encoding between 7bit, 8bit, and quoted-printable while preserving byte-for-byte decoded parity across destination MTAs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Table rendering encoding choice, raw byte size delta, and target ESMTP `8BITMIME` support status.
  - *Web*: `<encoding-alternator-view>` showing real-time ASCII/hex conversion and MTA capability matrix.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph alternate-encoding --target-mx=<host> [--has-8bit-preference] [--is-dry-run]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TransferEncodingRules (
      id TEXT PRIMARY KEY,
      domainPattern TEXT NOT NULL,
      preferredEncoding TEXT NOT NULL, -- '7bit', '8bit', 'quoted-printable'
      is8BitSupported INTEGER NOT NULL DEFAULT 1,
      hasAutoFallback INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TransferEncodingRule struct {
      DomainPattern     string `json:"domainPattern"`
      PreferredEncoding string `json:"preferredEncoding"`
      Is8BitSupported   bool   `json:"is8BitSupported"`
      HasAutoFallback   bool   `json:"hasAutoFallback"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Wired into `pkg/warmup/mime/encoder.go` with automated MX capability caching.

---

### Improvement 708: Date Header Jitter & Timezone Distribution Emulator
- **ID**: `IMP-708`
- **Cluster**: 71
- **Goal**: Emulates realistic human sending distribution by jittering seconds offsets and aligning timezone tokens (`-0400`, `+0100`, `EST`, `UTC`) to simulated sender locales.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss timeline bar showing historical jitter distribution across the active sending window.
  - *Web*: `<timezone-jitter-map>` element rendering a world map with sender node markers and local timezone offsets.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph jitter-date --window-sec=<seconds> [--has-random-offset] [--is-utc]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DateJitterLogs (
      id TEXT PRIMARY KEY,
      messageId TEXT NOT NULL,
      targetTimezone TEXT NOT NULL,
      jitterOffsetSec INTEGER NOT NULL,
      formattedDateHeader TEXT NOT NULL,
      isMonotonic INTEGER NOT NULL DEFAULT 1,
      hasTimezoneAligned INTEGER NOT NULL DEFAULT 1,
      appliedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DateJitterLog struct {
      MessageID          string `json:"messageId"`
      TargetTimezone     string `json:"targetTimezone"`
      JitterOffsetSec    int    `json:"jitterOffsetSec"`
      IsMonotonic        bool   `json:"isMonotonic"`
      HasTimezoneAligned bool   `json:"hasTimezoneAligned"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/envelope/date_jitter.go` during RFC envelope assembly.

---

### Improvement 709: Optional X-Header Sanitization & Polymorphic Injection
- **ID**: `IMP-709`
- **Cluster**: 71
- **Goal**: Selectively purges or injects polymorphic diagnostic X-headers (`X-Mailer`, `X-Priority`, `X-Originating-IP`) matching selected software fingerprints (e.g. Roundcube, Thunderbird).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual checklist showing active, sanitized, and injected X-headers with color-coded confidence tags.
  - *Web*: `<xheader-polymorphic-studio>` displaying interactive policy toggle controls.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph manage-xheaders [--has-mailer-tag] [--is-strip-all] [--has-audit-log]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS PolymorphicXHeaders (
      id TEXT PRIMARY KEY,
      clientProfile TEXT NOT NULL,
      headerName TEXT NOT NULL,
      headerValueTemplate TEXT NOT NULL,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasDynamicVariables INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type PolymorphicXHeader struct {
      ClientProfile       string `json:"clientProfile"`
      HeaderName          string `json:"headerName"`
      IsActive            bool   `json:"isActive"`
      HasDynamicVariables bool   `json:"hasDynamicVariables"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Injected through `pkg/warmup/envelope/xheaders.go` with policy overrides.

---

### Improvement 710: Morph Profile Presets Manager & Regression Sentinel
- **ID**: `IMP-710`
- **Cluster**: 71
- **Goal**: Persists, evaluates, and audits morphing configuration profiles in `repodb/warmup.db` while validating deliverability regression metrics against previous baseline runs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Comprehensive profile overview table rendered with Lipgloss border styling (`#805AD5`), detailing deliverability win rates and active status.
  - *Web*: `<morph-profile-manager>` component with import/export capabilities for zero-CDN offline administration.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup morph profile list [--is-active-only] [--has-metrics]`
  - `gitmap email warmup morph profile apply --profile-name=<name> [--has-strict-validation]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MorphConfigurationPresets (
      id TEXT PRIMARY KEY,
      presetName TEXT NOT NULL UNIQUE,
      boundaryStrategy TEXT NOT NULL,
      casingStrategy TEXT NOT NULL,
      permutationLevel INTEGER NOT NULL DEFAULT 1,
      isActive INTEGER NOT NULL DEFAULT 1,
      hasRegressionGuard INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MorphConfigurationPreset struct {
      PresetName         string `json:"presetName"`
      BoundaryStrategy   string `json:"boundaryStrategy"`
      PermutationLevel   int    `json:"permutationLevel"`
      IsActive           bool   `json:"isActive"`
      HasRegressionGuard bool   `json:"hasRegressionGuard"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Central configuration registry in `pkg/warmup/morph/preset.go` linked to `gitmap email warmup config`.

---

# Cluster 72: Distributed Rate-Limiting Mesh & Cross-Instance Token Bucket Sync (711–720)

---

### Improvement 711: Peer-to-Peer UDP Gossip Rate-Limit Telemetry Synchronizer
- **ID**: `IMP-711`
- **Cluster**: 72 (Distributed Rate-Limiting Mesh)
- **Goal**: Synchronizes egress token consumption deltas across multi-node GitMap clusters using lightweight UDP gossip datagrams every 250ms with zero centralized server dependencies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Real-time peer gossip matrix displaying node IDs, round-trip latencies, and vector clock status styled with Lipgloss cyan (`#00B4D8`).
  - *Web*: `<gossip-mesh-topology>` interactive SVG node-link graph with animated UDP packet pulses.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh gossip-status [--has-peer-metrics] [--is-live]`
  - Flag `--has-peer-metrics`: Outputs packet loss and clock drift metrics per peer.
  - Flag `--is-live`: Launches a Bubble Tea continuous telemetry monitor.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MeshGossipPeers (
      id TEXT PRIMARY KEY,
      nodeId TEXT NOT NULL UNIQUE,
      ipAddress TEXT NOT NULL,
      gossipPort INTEGER NOT NULL,
      lastHeartbeatAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      isAlive INTEGER NOT NULL DEFAULT 1,
      hasClockDrift INTEGER NOT NULL DEFAULT 0,
      latencyMs REAL NOT NULL DEFAULT 0.0
  );
  ```
  ```go
  type MeshGossipPeer struct {
      NodeID        string  `json:"nodeId"`
      GossipPort    int     `json:"gossipPort"`
      IsAlive       bool    `json:"isAlive"`
      HasClockDrift bool    `json:"hasClockDrift"`
      LatencyMs     float64 `json:"latencyMs"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Core networking daemon in `pkg/warmup/mesh/gossip.go` integrated with GitMap cluster fleet.

---

### Improvement 712: Multi-Node Token Bucket Quota Partitioning Engine
- **ID**: `IMP-712`
- **Cluster**: 72
- **Goal**: Divides global target MX domain rate limits (e.g. 50 msgs/hour for `gmail.com`) dynamically among active fleet nodes based on weighted historical throughput.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual bar chart showing allocated vs consumed quotas per node with green/yellow/red threshold alerts.
  - *Web*: `<quota-partition-hud>` element showing dynamic slice rebalancing animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh partition-quota --domain=<domain> [--has-rebalance] [--is-dry-run]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DomainTokenAllocations (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      nodeId TEXT NOT NULL,
      allocatedQuota INTEGER NOT NULL,
      consumedQuota INTEGER NOT NULL DEFAULT 0,
      isPartitionBalanced INTEGER NOT NULL DEFAULT 1,
      hasExhaustedTokens INTEGER NOT NULL DEFAULT 0,
      windowExpiresAt DATETIME NOT NULL
  );
  ```
  ```go
  type DomainTokenAllocation struct {
      TargetDomain        string `json:"targetDomain"`
      NodeID              string `json:"nodeId"`
      AllocatedQuota      int    `json:"allocatedQuota"`
      ConsumedQuota       int    `json:"consumedQuota"`
      IsPartitionBalanced bool   `json:"isPartitionBalanced"`
      HasExhaustedTokens  bool   `json:"hasExhaustedTokens"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/mesh/partitioner.go` with lockless local reservations.

---

### Improvement 713: Microsecond-Grade Local Lockless Token Bucket Allocator
- **ID**: `IMP-713`
- **Cluster**: 72
- **Goal**: Employs atomic CPU operations (`sync/atomic`) for local in-memory token reservations, achieving sub-microsecond latency and zero lock contention during high-concurrency dispatch.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Displays nanosecond reservation benchmarks and current bucket saturation meters via Lipgloss formatted dials.
  - *Web*: `<lockless-bucket-meter>` rendering live SVG circular gauges for token availability.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh reserve-token --domain=<domain> [--count=<n>] [--is-blocking]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS LocalBucketReservations (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      tokensRequested INTEGER NOT NULL,
      tokensGranted INTEGER NOT NULL,
      reservationDurationMicros INTEGER NOT NULL,
      isImmediate INTEGER NOT NULL DEFAULT 1,
      hasQueuedWait INTEGER NOT NULL DEFAULT 0,
      reservedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type LocalBucketReservation struct {
      TargetDomain              string `json:"targetDomain"`
      TokensGranted             int    `json:"tokensGranted"`
      ReservationDurationMicros int    `json:"reservationDurationMicros"`
      IsImmediate               bool   `json:"isImmediate"`
      HasQueuedWait             bool   `json:"hasQueuedWait"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Located in `pkg/warmup/ratelimit/bucket.go` supporting high-speed sender routines.

---

### Improvement 714: Network Partition & Split-Brain Safe Auto-Throttling Fallback
- **ID**: `IMP-714`
- **Cluster**: 72
- **Goal**: Automatically downgrades a disconnected node's rate limit to a safe isolated floor ($1/N$ cluster capacity) whenever peer gossip heartbeats lapse beyond 1500ms.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Pulsing yellow warning indicator (`[WARN: MESH_ISOLATED - SAFE FLOOR ACTIVE]`) with fallback quota metrics.
  - *Web*: `<split-brain-alert-banner>` with automatic health recovery countdown and diagnostic probe triggers.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh test-partition --simulate-disconnect [--has-circuit-break] [--is-json]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SplitBrainFallbackEvents (
      id TEXT PRIMARY KEY,
      nodeId TEXT NOT NULL,
      isolatedSince DATETIME NOT NULL,
      fallbackQuotaCap INTEGER NOT NULL,
      isFallbackActive INTEGER NOT NULL DEFAULT 1,
      hasRecoveredQuorum INTEGER NOT NULL DEFAULT 0,
      recoveredAt DATETIME
  );
  ```
  ```go
  type SplitBrainFallbackEvent struct {
      NodeID             string `json:"nodeId"`
      FallbackQuotaCap   int    `json:"fallbackQuotaCap"`
      IsFallbackActive   bool   `json:"isFallbackActive"`
      HasRecoveredQuorum bool   `json:"hasRecoveredQuorum"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Safety watchdog in `pkg/warmup/mesh/sentinel.go` preventing accidental ISP flooding.

---

### Improvement 715: Adaptive Leaky Bucket Burst-Dampening Smoother
- **ID**: `IMP-715`
- **Cluster**: 72
- **Goal**: Smooths out traffic spikes into deterministic Poisson inter-send delays, eliminating burst artifacts that trigger ISP rate-limiting heuristics (`452 4.5.3 Too many recipients`).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: High-resolution ANSI histogram showing dispatch inter-arrival times across 50ms buckets.
  - *Web*: `<burst-dampener-histogram>` showing real-time smoothing curves alongside theoretical Poisson limits.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh smooth-burst --target-mx=<host> [--has-poisson-jitter] [--is-enforced]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BurstDampeningProfiles (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      leakRatePerSec REAL NOT NULL,
      bufferCapacity INTEGER NOT NULL,
      isSmoothingActive INTEGER NOT NULL DEFAULT 1,
      hasJitterEnabled INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BurstDampeningProfile struct {
      TargetMxHost      string  `json:"targetMxHost"`
      LeakRatePerSec    float64 `json:"leakRatePerSec"`
      IsSmoothingActive bool    `json:"isSmoothingActive"`
      HasJitterEnabled  bool    `json:"hasJitterEnabled"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Orchestrated in `pkg/warmup/ratelimit/leaky.go` prior to socket dispatch.

---

### Improvement 716: Global MX Domain Aggregation & CIDR Subnet Grouping
- **ID**: `IMP-716`
- **Cluster**: 72
- **Goal**: Automatically aggregates disparate recipient domains into common parent MX infrastructures (e.g., mapping custom business domains hosted on Google Workspace to a unified `aspmx.l.google.com` rate limit).
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Tree view grouping target domains under resolved MX parent roots with aggregate token tallies.
  - *Web*: `<mx-infrastructure-tree>` displaying interactive network topology graphs grouped by autonomous system numbers (ASN).
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh resolve-mx-group --domain=<domain> [--has-cidr-lookup] [--is-cached]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DomainMxMappings (
      id TEXT PRIMARY KEY,
      recipientDomain TEXT NOT NULL UNIQUE,
      primaryMxHost TEXT NOT NULL,
      mxAsnNumber INTEGER NOT NULL,
      isGoogleWorkspace INTEGER NOT NULL DEFAULT 0,
      hasMicrosoft365 INTEGER NOT NULL DEFAULT 0,
      resolvedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DomainMxMapping struct {
      RecipientDomain   string `json:"recipientDomain"`
      PrimaryMxHost     string `json:"primaryMxHost"`
      MxAsnNumber       int    `json:"mxAsnNumber"`
      IsGoogleWorkspace bool   `json:"isGoogleWorkspace"`
      HasMicrosoft365   bool   `json:"hasMicrosoft365"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforced in `pkg/warmup/dns/mx_group.go` for shared rate-limiting policies.

---

### Improvement 717: Dynamic Token Borrowing & Node Surplus Exchange
- **ID**: `IMP-717`
- **Cluster**: 72
- **Goal**: Permits nodes experiencing high local queue pressure to dynamically borrow unconsumed token quotas from idle peer nodes via gossip-negotiated credit grants.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual ledger table detailing active credit negotiations, borrowed amounts, and repayment deadlines.
  - *Web*: `<token-exchange-ledger>` element rendering real-time peer credit swap transactions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh borrow-tokens --from-node=<nodeId> --tokens=<n> [--has-auto-return]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TokenBorrowTransactions (
      id TEXT PRIMARY KEY,
      lenderNodeId TEXT NOT NULL,
      borrowerNodeId TEXT NOT NULL,
      tokenCount INTEGER NOT NULL,
      isRepaid INTEGER NOT NULL DEFAULT 0,
      hasExpired INTEGER NOT NULL DEFAULT 0,
      transactionTime DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TokenBorrowTransaction struct {
      LenderNodeID   string `json:"lenderNodeId"`
      BorrowerNodeID string `json:"borrowerNodeId"`
      TokenCount     int    `json:"tokenCount"`
      IsRepaid       bool   `json:"isRepaid"`
      HasExpired     bool   `json:"hasExpired"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Protocol engine in `pkg/warmup/mesh/exchange.go` maximizing cluster sending throughput.

---

### Improvement 718: Sliding Window Counter Reconciliation & Drift Detector
- **ID**: `IMP-718`
- **Cluster**: 72
- **Goal**: Reconciles distributed token counts across a 60-minute sliding window, executing cryptographic vector clock reconciliation to detect and correct token drift.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Vector clock drift summary rendered in Lipgloss table format highlighting synchronized status badges (`[SYNCED: 0ms DRIFT]`).
  - *Web*: `<vector-clock-scrubber>` element with interactive timeline drift slider and automated resync trigger.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh reconcile-drift [--is-forced] [--has-audit-trail]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SlidingWindowReconciliation (
      id TEXT PRIMARY KEY,
      windowStart DATETIME NOT NULL,
      windowEnd DATETIME NOT NULL,
      expectedTokensTotal INTEGER NOT NULL,
      actualTokensTotal INTEGER NOT NULL,
      driftDelta INTEGER NOT NULL,
      isReconciled INTEGER NOT NULL DEFAULT 1,
      hasDriftDetected INTEGER NOT NULL DEFAULT 0,
      reconciledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SlidingWindowReconciliation struct {
      ExpectedTokensTotal int  `json:"expectedTokensTotal"`
      ActualTokensTotal   int  `json:"actualTokensTotal"`
      DriftDelta          int  `json:"driftDelta"`
      IsReconciled        bool `json:"isReconciled"`
      HasDriftDetected    bool `json:"hasDriftDetected"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Audited in `pkg/warmup/mesh/reconciler.go` maintaining cluster-wide invariant guarantees.

---

### Improvement 719: Egress Backpressure Circuit Breaker & Queue Quarantining
- **ID**: `IMP-719`
- **Cluster**: 72
- **Goal**: Halts egress dispatch and moves pending outbox items to temporary quarantine when peer nodes report sudden MX reputation drops or persistent 421 deferrals.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Prominent red alert banner (`[CIRCUIT BREAKER OPEN: QUARANTINE ENGAGED]`) with affected recipient queue tallies.
  - *Web*: `<circuit-breaker-hud>` displaying visual trip triggers, cooldown timers, and manual override switches.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh circuit-break --domain=<domain> [--is-manual-trip] [--has-auto-cooldown]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS CircuitBreakerStates (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL UNIQUE,
      tripReason TEXT NOT NULL,
      trippedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      cooldownPeriodSec INTEGER NOT NULL DEFAULT 300,
      isCircuitOpen INTEGER NOT NULL DEFAULT 1,
      hasAutoRecovery INTEGER NOT NULL DEFAULT 1
  );
  ```
  ```go
  type CircuitBreakerState struct {
      TargetDomain    string `json:"targetDomain"`
      TripReason      string `json:"tripReason"`
      IsCircuitOpen   bool   `json:"isCircuitOpen"`
      HasAutoRecovery bool   `json:"hasAutoRecovery"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Core defense component in `pkg/warmup/circuit/breaker.go` protecting sender IP reputation.

---

### Improvement 720: Multi-Interface Binding & Egress IP Load Balancer
- **ID**: `IMP-720`
- **Cluster**: 72
- **Goal**: Balances outbound SMTP TCP socket connections across multiple local network interfaces and secondary IPv4/IPv6 addresses according to round-robin or least-loaded policies.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual interface card list showing active local IPs, assigned subnet masks, and outbound connection counts.
  - *Web*: `<egress-ip-balancer>` showing real-time throughput meters per network interface.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mesh bind-interfaces --ips=<ip1,ip2> [--is-round-robin] [--has-ipv6-support]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EgressInterfaces (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL,
      ipAddress TEXT NOT NULL UNIQUE,
      isIpv6 INTEGER NOT NULL DEFAULT 0,
      activeConnectionCount INTEGER NOT NULL DEFAULT 0,
      isAvailable INTEGER NOT NULL DEFAULT 1,
      hasHealthCheckPassed INTEGER NOT NULL DEFAULT 1,
      lastCheckedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EgressInterface struct {
      InterfaceName        string `json:"interfaceName"`
      IPAddress            string `json:"ipAddress"`
      IsIpv6               bool   `json:"isIpv6"`
      IsAvailable          bool   `json:"isAvailable"`
      HasHealthCheckPassed bool   `json:"hasHealthCheckPassed"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Network layer in `pkg/warmup/transport/dialer.go` supporting multi-homed warmups.

---

# Cluster 73: Real-Time Bounded SMTP Session Connection Multiplexer & Keep-Alive TUI (721–730)

---

### Improvement 721: Bounded SMTP Persistent Connection Reuse Pool
- **ID**: `IMP-721`
- **Cluster**: 73 (SMTP Multiplexer & Keep-Alive TUI)
- **Goal**: Maintains a pool of warm, TLS-established SMTP sessions to high-volume MX destinations to bypass expensive TLS handshakes while bounding maximum reuses per session.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Interactive pool monitor styled with Lipgloss borders (`#4FD1C5`), showing idle connections, active transactions, and socket ages.
  - *Web*: `<smtp-pool-dashboard>` rendering real-time connection lifecycle timelines with state progression badges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp pool-status [--target-mx=<host>] [--has-idle-detail] [--is-live]`
  - Flag `--has-idle-detail`: Prints idle timeout and total transactions completed per connection.
  - Flag `--is-live`: Launches interactive Bubble Tea session view.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpConnectionSessions (
      id TEXT PRIMARY KEY,
      poolKey TEXT NOT NULL,
      targetMxHost TEXT NOT NULL,
      localBindingIp TEXT NOT NULL,
      tlsVersion TEXT NOT NULL,
      transactionCount INTEGER NOT NULL DEFAULT 0,
      isIdle INTEGER NOT NULL DEFAULT 1,
      hasExceededMaxUses INTEGER NOT NULL DEFAULT 0,
      connectedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      lastActiveAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpConnectionSession struct {
      PoolKey            string `json:"poolKey"`
      TargetMxHost       string `json:"targetMxHost"`
      TransactionCount   int    `json:"transactionCount"`
      IsIdle             bool   `json:"isIdle"`
      HasExceededMaxUses bool   `json:"hasExceededMaxUses"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/smtp/pool.go` reducing delivery latency by >80%.

---

### Improvement 722: RFC 2920 SMTP PIPELINING Command Pipeline Bundler
- **ID**: `IMP-722`
- **Cluster**: 73
- **Goal**: Bundles `RSET`, `MAIL FROM`, and `RCPT TO` commands into a single TCP write buffer when remote MTAs advertise `PIPELINING`, cutting round-trip times over high-latency WAN links.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Hex/ASCII command stream visualizer showing pipelined command packets in cyan and remote MTA response vectors in green.
  - *Web*: `<pipelining-packet-view>` rendering waterfall diagrams of pipelined versus serialized transactions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp pipeline-test --target-mx=<host> [--is-strict-validation] [--has-chunking]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpPipeliningLogs (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      commandCount INTEGER NOT NULL,
      roundTripLatencyMs REAL NOT NULL,
      isPipeliningSuccessful INTEGER NOT NULL DEFAULT 1,
      hasOutOfOrderResponse INTEGER NOT NULL DEFAULT 0,
      executedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpPipeliningLog struct {
      ConnectionID           string  `json:"connectionId"`
      CommandCount           int     `json:"commandCount"`
      RoundTripLatencyMs     float64 `json:"roundTripLatencyMs"`
      IsPipeliningSuccessful bool    `json:"isPipeliningSuccessful"`
      HasOutOfOrderResponse  bool    `json:"hasOutOfOrderResponse"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Optimized in `pkg/warmup/smtp/pipeline.go` accelerating batch warmups.

---

### Improvement 723: Socket-Level TCP Keep-Alive & RSET Heartbeat Sentinel
- **ID**: `IMP-723`
- **Cluster**: 73
- **Goal**: Issues periodic `RSET` commands and configures TCP-level keep-alive probes (`SO_KEEPALIVE`) to prevent silent firewall drops during inter-send warming intervals.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Live pulse ticker showing seconds since last heartbeat and TCP probe acknowledgement times.
  - *Web*: `<keepalive-pulse-monitor>` rendering animated SVG heartbeat waves with health score overlays.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp keepalive-watch [--interval-sec=<n>] [--has-rset-probe] [--is-verbose]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpKeepAliveHeartbeats (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      probeType TEXT NOT NULL, -- 'tcp-probe', 'smtp-rset', 'smtp-noop'
      responseLatencyMs REAL NOT NULL,
      isConnectionHealthy INTEGER NOT NULL DEFAULT 1,
      hasTimedOut INTEGER NOT NULL DEFAULT 0,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpKeepAliveHeartbeat struct {
      ConnectionID        string  `json:"connectionId"`
      ProbeType           string  `json:"probeType"`
      ResponseLatencyMs   float64 `json:"responseLatencyMs"`
      IsConnectionHealthy bool    `json:"isConnectionHealthy"`
      HasTimedOut         bool    `json:"hasTimedOut"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Monitored by `pkg/warmup/smtp/keepalive.go` for seamless idle connection preservation.

---

### Improvement 724: Interactive Bubble Tea SMTP Session Scrubber TUI
- **ID**: `IMP-724`
- **Cluster**: 73
- **Goal**: Full-featured Bubble Tea v2 terminal interface enabling developers to inspect live SMTP socket dialogues, scrub conversation histories, and inspect raw TLS handshake states.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Split-pane TUI; left pane lists active sessions by MX target, right pane displays scrollable, color-coded protocol transcripts with keyboard navigation (`j`/`k` scroll, `Tab` switch).
  - *Web*: `<session-scrubber-tui>` web component faithfully mimicking terminal layout and keyboard shortcuts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp scrub-session [--is-interactive] [--has-color-highlighting]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpTranscriptSnippets (
      id TEXT PRIMARY KEY,
      connectionId TEXT NOT NULL,
      direction TEXT NOT NULL, -- 'outbound-client', 'inbound-server'
      commandOrResponse TEXT NOT NULL,
      timestampMicros INTEGER NOT NULL,
      isTlsEncrypted INTEGER NOT NULL DEFAULT 1,
      hasSyntaxError INTEGER NOT NULL DEFAULT 0
  );
  ```
  ```go
  type SmtpTranscriptSnippet struct {
      ConnectionID      string `json:"connectionId"`
      Direction         string `json:"direction"`
      CommandOrResponse string `json:"commandOrResponse"`
      IsTlsEncrypted    bool   `json:"isTlsEncrypted"`
      HasSyntaxError    bool   `json:"hasSyntaxError"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Standalone terminal UI launched from `pkg/tui/smtp_scrubber.go` via `gitmap email warmup session-tui`.

---

### Improvement 725: TLS Session Resumption & Zero-RTT PSK Caching Engine
- **ID**: `IMP-725`
- **Cluster**: 73
- **Goal**: Caches TLS 1.3 Pre-Shared Keys (PSK) and TLS 1.2 session tickets locally to enable 0-RTT and 1-RTT session resumption across consecutive SMTP STARTTLS negotiations.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Handshake latency comparison table displaying ticket cache hit rates and cryptographic cipher negotiation times.
  - *Web*: `<tls-resumption-matrix>` component with cipher suite breakdowns and cache validity expiration gauges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp tls-cache-status [--target-mx=<host>] [--has-psk-stats] [--is-purged]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsSessionTickets (
      id TEXT PRIMARY KEY,
      mxHost TEXT NOT NULL,
      tlsVersion TEXT NOT NULL,
      cipherSuiteId INTEGER NOT NULL,
      ticketData BLOB NOT NULL,
      isPskResumed INTEGER NOT NULL DEFAULT 0,
      hasTicketExpired INTEGER NOT NULL DEFAULT 0,
      storedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
      expiresAt DATETIME NOT NULL
  );
  ```
  ```go
  type TlsSessionTicketRecord struct {
      MxHost           string `json:"mxHost"`
      TlsVersion       string `json:"tlsVersion"`
      IsPskResumed     bool   `json:"isPskResumed"`
      HasTicketExpired bool   `json:"hasTicketExpired"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Security layer in `pkg/warmup/crypto/tls_cache.go` speeding up TLS handshakes.

---

### Improvement 726: Graceful SMTP Connection Drain & Shutdown Coordinator
- **ID**: `IMP-726`
- **Cluster**: 73
- **Goal**: Coordinates orderly teardown of open SMTP socket sessions (`QUIT` transmission and TCP `FIN` verification) during daemon restarts to avoid leaving orphan connections on remote MTAs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Step-by-step shutdown progression bar showing connection drain completion (`[Draining: 14/14 closed cleanly]`).
  - *Web*: `<shutdown-coordinator-banner>` showing drain timer countdown and graceful closure verification.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp drain-pool [--timeout-sec=<n>] [--is-force-kill] [--has-audit-log]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpDrainEvents (
      id TEXT PRIMARY KEY,
      totalSessionsToDrain INTEGER NOT NULL,
      sessionsCleanlyClosed INTEGER NOT NULL,
      sessionsForceKilled INTEGER NOT NULL,
      isCleanDrain INTEGER NOT NULL DEFAULT 1,
      hasTimeoutOccurred INTEGER NOT NULL DEFAULT 0,
      drainedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpDrainEvent struct {
      TotalSessionsToDrain  int  `json:"totalSessionsToDrain"`
      SessionsCleanlyClosed int  `json:"sessionsCleanlyClosed"`
      IsCleanDrain          bool `json:"isCleanDrain"`
      HasTimeoutOccurred    bool `json:"hasTimeoutOccurred"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed lifecycle hook in `pkg/warmup/smtp/coordinator.go` invoked during SIGTERM/SIGINT.

---

### Improvement 727: Per-MX Dynamic Socket Concurrency Limiter
- **ID**: `IMP-727`
- **Cluster**: 73
- **Goal**: Enforces strictly bounded concurrent socket counts per remote MX hostname (e.g. max 5 simultaneous TCP connections to Yahoo/AOL MTAs) to avoid connection rate throttling.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Visual gauge cards for top MX providers indicating active vs max allowed socket channels.
  - *Web*: `<mx-concurrency-gauge>` rendering circular gauges with customizable threshold limits.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp set-concurrency --mx=<host> --max-sockets=<n> [--has-auto-tune]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MxConcurrencyLimits (
      id TEXT PRIMARY KEY,
      mxHostnamePattern TEXT NOT NULL UNIQUE,
      maxAllowedSockets INTEGER NOT NULL DEFAULT 5,
      currentActiveSockets INTEGER NOT NULL DEFAULT 0,
      isAutoTuned INTEGER NOT NULL DEFAULT 1,
      hasReachedCap INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MxConcurrencyLimit struct {
      MxHostnamePattern string `json:"mxHostnamePattern"`
      MaxAllowedSockets int    `json:"maxAllowedSockets"`
      IsAutoTuned       bool   `json:"isAutoTuned"`
      HasReachedCap     bool   `json:"hasReachedCap"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Synchronized guard in `pkg/warmup/transport/limiter.go` integrated with dispatch schedulers.

---

### Improvement 728: STARTTLS Handshake Cipher Hardener & Downgrade Sentinel
- **ID**: `IMP-728`
- **Cluster**: 73
- **Goal**: Enforces modern cipher suites (TLS 1.3 ChaCha20-Poly1305 / AES-256-GCM) and aborts connections immediately if intermediate proxies attempt plaintext `STARTTLS` stripping attacks.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Cryptographic security summary table rendering cipher name, PFS key exchange curve, and security grade (`[GRADE: A+]`).
  - *Web*: `<starttls-cipher-sentinel>` element showing cryptographic negotiation details and tamper alerts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp test-starttls --target-mx=<host> [--is-strict-downgrade-guard] [--has-pfs-check]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpTlsNegotiationLogs (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      negotiatedTlsVersion TEXT NOT NULL,
      negotiatedCipherSuite TEXT NOT NULL,
      isDowngradeAttempted INTEGER NOT NULL DEFAULT 0,
      hasForwardSecrecy INTEGER NOT NULL DEFAULT 1,
      isTlsEnforced INTEGER NOT NULL DEFAULT 1,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpTlsNegotiationLog struct {
      TargetMxHost          string `json:"targetMxHost"`
      NegotiatedTlsVersion  string `json:"negotiatedTlsVersion"`
      NegotiatedCipherSuite string `json:"negotiatedCipherSuite"`
      IsDowngradeAttempted  bool   `json:"isDowngradeAttempted"`
      HasForwardSecrecy     bool   `json:"hasForwardSecrecy"`
      IsTlsEnforced         bool   `json:"isTlsEnforced"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Embedded in `pkg/warmup/crypto/tls_sentinel.go` guarding against man-in-the-middle tampering.

---

### Improvement 729: SMTP Server Greeting Banner Heuristic & Latency Scorer
- **ID**: `IMP-729`
- **Cluster**: 73
- **Goal**: Analyzes remote MTA `220` SMTP greetings and measure initial banner delivery latency to detect tarpitting, honeypot filters, or aggressive anti-spam delays.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Displays remote banner text alongside banner latency bar: `[220 banner received: 48ms] [Tarpit: No]`.
  - *Web*: `<mta-greeting-inspector>` rendering latency distribution graphs and tarpit probability indicators.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp check-greeting --target-mx=<host> [--has-tarpit-detect] [--is-verbose]`
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpGreetingAudits (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      greetingBannerRaw TEXT NOT NULL,
      bannerLatencyMs REAL NOT NULL,
      isTarpitDetected INTEGER NOT NULL DEFAULT 0,
      hasBannerFqdnMatched INTEGER NOT NULL DEFAULT 1,
      auditedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpGreetingAudit struct {
      TargetMxHost         string  `json:"targetMxHost"`
      BannerLatencyMs      float64 `json:"bannerLatencyMs"`
      IsTarpitDetected     bool    `json:"isTarpitDetected"`
      HasBannerFqdnMatched bool    `json:"hasBannerFqdnMatched"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Diagnostic scanner in `pkg/warmup/smtp/greeting.go` feeding warmup health telemetry.

---

### Improvement 730: Connection Multiplexer Diagnostics & SQLite State Dumper
- **ID**: `IMP-730`
- **Cluster**: 73
- **Goal**: Dumps real-time connection multiplexer state, pool allocations, and socket health records from `repodb/warmup.db` into formatted JSON or Markdown forensic summaries.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Formatted Lipgloss card summary detailing active pools, aggregate throughput, and recent connection errors.
  - *Web*: `<multiplexer-state-viewer>` with instant one-click sanitized JSON export and copy-to-clipboard functionality.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup smtp dump-state [--is-json] [--has-sanitized-ips] [--output-file=<path>]`
  - Flag `--is-json`: Outputs formatted JSON stream to stdout.
  - Flag `--has-sanitized-ips`: Masks public IPv4/IPv6 addresses for privacy.
  - Flag `--output-file`: Writes report to specified file path.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db` and `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MultiplexerDumpSnapshots (
      id TEXT PRIMARY KEY,
      totalActiveSockets INTEGER NOT NULL,
      totalIdleSockets INTEGER NOT NULL,
      totalTransactionsLifetime INTEGER NOT NULL,
      isHealthy INTEGER NOT NULL DEFAULT 1,
      hasSocketLeaked INTEGER NOT NULL DEFAULT 0,
      snapshotPayload TEXT NOT NULL, -- JSON blob of state
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MultiplexerDumpSnapshot struct {
      TotalActiveSockets        int    `json:"totalActiveSockets"`
      TotalIdleSockets          int    `json:"totalIdleSockets"`
      TotalTransactionsLifetime int    `json:"totalTransactionsLifetime"`
      IsHealthy                 bool   `json:"isHealthy"`
      HasSocketLeaked           bool   `json:"hasSocketLeaked"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Operational exporter in `pkg/warmup/smtp/dumper.go` integrated into `gitmap pe -t` telemetry feeds.
