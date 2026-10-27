# Improvements 1001 to 1030: Anycast Routing, MPTCP Bonding & Jitter Waveform TUI

This document details the complete specifications for Improvements 1001 through 1030 across Clusters 101, 102, and 103 of Sequence 347.

---

# Cluster 101: Global CDN Edge SMTP Anycast Proximity Router & Geo-DNS Telemetry (1001–1010)

---

### Improvement 1001: BGP Anycast Proximity Matrix & Sub-Millisecond Egress Prober
- **ID**: `IMP-1001`
- **Cluster**: 101 (Global CDN Edge SMTP Anycast Proximity Router & Geo-DNS Telemetry)
- **Goal**: Measure and maintain real-time BGP routing latency across global Anycast edge POPs to route warming envelopes through the geographically and topologically closest egress gateway to the destination Mailbox Provider.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss status card displaying global edge POP latencies. Uses electric cyan (`#00E5FF`) for sub-10ms routes, spring green (`#00E676`) for 10-30ms, and amber (`#FFD600`) for routes exceeding 30ms. Includes ANSI route topology diagram.
  - *Web*: Custom Web Component `<anycast-proximity-matrix>` rendering dynamic topological link meshes in pure SVG without external CDNs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast probe-matrix --destination-domain=<domain> [--pop-region=<region>] [--is-deep-probe] [--has-bgp-lookup]`
  - Flag `--destination-domain`: Target recipient domain (e.g., `gmail.com`, `outlook.com`).
  - Flag `--pop-region`: Limits probe to a specific continent or edge region (e.g., `us-east`, `eu-central`).
  - Flag `--is-deep-probe`: Enables ICMP, TCP SYN, and TLS handshake latency triangulation.
  - Flag `--has-bgp-lookup`: Queries local BGP routing table for AS path length comparison.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AnycastProximityProbes (
      id TEXT PRIMARY KEY,
      destinationDomain TEXT NOT NULL,
      popNodeId TEXT NOT NULL,
      popRegion TEXT NOT NULL,
      targetMxHost TEXT NOT NULL,
      asPathLength INTEGER NOT NULL DEFAULT 0,
      rttLatencyMs REAL NOT NULL DEFAULT 0.0,
      isOptimalRoute INTEGER NOT NULL DEFAULT 1,
      hasBgpLookup INTEGER NOT NULL DEFAULT 1,
      measuredAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type AnycastProximityProbe struct {
      ID                string    `json:"id"`
      DestinationDomain string    `json:"destinationDomain"`
      PopNodeID         string    `json:"popNodeId"`
      PopRegion         string    `json:"popRegion"`
      TargetMxHost      string    `json:"targetMxHost"`
      AsPathLength      int       `json:"asPathLength"`
      RttLatencyMs      float64   `json:"rttLatencyMs"`
      IsOptimalRoute    bool      `json:"isOptimalRoute"`
      HasBgpLookup      bool      `json:"hasBgpLookup"`
      MeasuredAt        time.Time `json:"measuredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated into `pkg/warmup/anycast/proximity.go` within the egress socket dispatcher.

---

### Improvement 1002: Autonomous Geo-DNS Telemetry Resolver & Real-Time EDNS Client Subnet Injector
- **ID**: `IMP-1002`
- **Cluster**: 101
- **Goal**: Resolves MX records utilizing EDNS Client Subnet (ECS) injection from edge POP vantage points to discover the exact localized inbound mail clusters deployed by major ISPs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-pane table showing MX records resolved with and without ECS injection, highlighting DNS cache poisoning defenses in bright magenta (`#E040FB`).
  - *Web*: `<geodns-ecs-inspector>` Web Component rendering interactive DNS record trees with SVG hop pins.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast geodns-resolve --target-domain=<domain> [--client-subnet=<cidr>] [--is-recursive] [--has-ecs-enabled]`
  - Flag `--target-domain`: Domain whose MX records are under evaluation.
  - Flag `--client-subnet`: Simulated client IPv4/IPv6 CIDR passed via EDNS0 (RFC 7871).
  - Flag `--is-recursive`: Forces recursive iterative resolution bypassing local DNS caches.
  - Flag `--has-ecs-enabled`: Toggles RFC 7871 client subnet option injection.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS GeoDnsResolutions (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      injectedSubnet TEXT NOT NULL,
      resolvedMxHost TEXT NOT NULL,
      mxPriority INTEGER NOT NULL DEFAULT 10,
      resolvedIpAddress TEXT NOT NULL,
      isEcsEnabled INTEGER NOT NULL DEFAULT 1,
      hasCacheHit INTEGER NOT NULL DEFAULT 0,
      resolvedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type GeoDnsResolution struct {
      ID                string    `json:"id"`
      TargetDomain      string    `json:"targetDomain"`
      InjectedSubnet    string    `json:"injectedSubnet"`
      ResolvedMxHost    string    `json:"resolvedMxHost"`
      MxPriority        int       `json:"mxPriority"`
      ResolvedIPAddress string    `json:"resolvedIpAddress"`
      IsEcsEnabled      bool      `json:"isEcsEnabled"`
      HasCacheHit       bool      `json:"hasCacheHit"`
      ResolvedAt        time.Time `json:"resolvedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked in `pkg/warmup/anycast/geodns.go` prior to outbound TCP socket creation.

---

### Improvement 1003: Edge Point-of-Presence (POP) Tunnel Multiplexer & GRE/WireGuard Encapsulator
- **ID**: `IMP-1003`
- **Cluster**: 101
- **Goal**: Encapsulates raw warming SMTP traffic over lightweight, kernel-accelerated WireGuard and GRE tunnels to forward envelopes instantly to remote Anycast egress edge gateways.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss network tunnel HUD displaying active WireGuard peer links, encrypted bytes transferred, and handshaking latency in royal indigo (`#3D5AFE`).
  - *Web*: `<edge-tunnel-hud>` Web Component rendering real-time tunnel throughput gauges using pure inline SVG paths.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast tunnel-connect --pop-id=<id> [--tunnel-type=<gre|wireguard>] [--is-persistent] [--has-keepalive]`
  - Flag `--pop-id`: Unique identifier of the remote Anycast edge POP node.
  - Flag `--tunnel-type`: Encapsulation protocol (default: `wireguard`).
  - Flag `--is-persistent`: Keeps tunnel interface alive across warm-up batches.
  - Flag `--has-keepalive`: Enables periodic 15-second tunnel heartbeat probes.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EdgeTunnelStates (
      id TEXT PRIMARY KEY,
      popId TEXT NOT NULL UNIQUE,
      tunnelType TEXT NOT NULL DEFAULT 'wireguard',
      localEndpoint TEXT NOT NULL,
      remoteEndpoint TEXT NOT NULL,
      isTunnelActive INTEGER NOT NULL DEFAULT 1,
      hasKeepalive INTEGER NOT NULL DEFAULT 1,
      txBytes INTEGER NOT NULL DEFAULT 0,
      rxBytes INTEGER NOT NULL DEFAULT 0,
      lastHandshakeAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EdgeTunnelState struct {
      ID              string    `json:"id"`
      PopID           string    `json:"popId"`
      TunnelType      string    `json:"tunnelType"`
      LocalEndpoint   string    `json:"localEndpoint"`
      RemoteEndpoint  string    `json:"remoteEndpoint"`
      IsTunnelActive  bool      `json:"isTunnelActive"`
      HasKeepalive    bool      `json:"hasKeepalive"`
      TxBytes         int64     `json:"txBytes"`
      RxBytes         int64     `json:"rxBytes"`
      LastHandshakeAt time.Time `json:"lastHandshakeAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/anycast/tunnel.go`.

---

### Improvement 1004: Localized ISP Peering Point Cost Optimizer & AS-Path Pruner
- **ID**: `IMP-1004`
- **Cluster**: 101
- **Goal**: Dynamically filters and selects Anycast edge egress gateways that have direct settlement-free peering (IXP) with destination Mailbox Providers, eliminating transit transit-hop jitter.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss hierarchical tree showing AS-Path traces, IXP peering interconnect names, and transit hop counts in slate grey (`#78909C`) and electric teal (`#00BFA5`).
  - *Web*: `<ixp-peering-visualizer>` Web Component rendering interactive network hops with SVG nodes.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast optimize-peering --target-asn=<asn> [--max-hops=<n>] [--is-direct-peering-required] [--has-ixp-verified]`
  - Flag `--target-asn`: Autonomous System Number of the recipient ISP (e.g., `AS15169` for Google).
  - Flag `--max-hops`: Maximum allowed intermediate BGP AS hops (default: `2`).
  - Flag `--is-direct-peering-required`: Discards any edge route lacking direct IXP interconnect.
  - Flag `--has-ixp-verified`: Validates PeeringDB records before committing routing policy.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS IxpPeeringRoutes (
      id TEXT PRIMARY KEY,
      targetAsn INTEGER NOT NULL,
      popId TEXT NOT NULL,
      ixpName TEXT NOT NULL,
      hopCount INTEGER NOT NULL DEFAULT 1,
      isDirectPeeringRequired INTEGER NOT NULL DEFAULT 1,
      hasIxpVerified INTEGER NOT NULL DEFAULT 1,
      routeWeight REAL NOT NULL DEFAULT 1.0,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type IxpPeeringRoute struct {
      ID                      string    `json:"id"`
      TargetAsn               int       `json:"targetAsn"`
      PopID                   string    `json:"popId"`
      IxpName                 string    `json:"ixpName"`
      HopCount                int       `json:"hopCount"`
      IsDirectPeeringRequired bool      `json:"isDirectPeeringRequired"`
      HasIxpVerified          bool      `json:"hasIxpVerified"`
      RouteWeight             float64   `json:"routeWeight"`
      CreatedAt               time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Integrated in `pkg/warmup/anycast/peering.go`.

---

### Improvement 1005: Anycast BGP Route Flap Damping & Convergence Guard
- **ID**: `IMP-1005`
- **Cluster**: 101
- **Goal**: Detects and dampens unstable BGP Anycast routes undergoing frequent route flaps, preventing TCP connection resets mid-flight during SMTP envelope delivery.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss alert banner showing suppressed routes in bold red (`#FF1744`) and recovered routes in lime (`#76FF03`) with flap penalty decay timers.
  - *Web*: `<route-flap-gauge>` Web Component displaying exponential decay curves in SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast check-flaps --pop-id=<id> [--penalty-threshold=<val>] [--is-auto-suppress] [--has-decay-applied]`
  - Flag `--pop-id`: Identifier of the edge POP node under evaluation.
  - Flag `--penalty-threshold`: Route flap penalty value triggering suppression (default: `2000`).
  - Flag `--is-auto-suppress`: Automatically blacklists flapping POPs from dispatch pool.
  - Flag `--has-decay-applied`: Recalculates half-life decay on flap history scores.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS RouteFlapDampings (
      id TEXT PRIMARY KEY,
      popId TEXT NOT NULL UNIQUE,
      flapCount INTEGER NOT NULL DEFAULT 0,
      penaltyScore REAL NOT NULL DEFAULT 0.0,
      isRouteSuppressed INTEGER NOT NULL DEFAULT 0,
      hasDecayApplied INTEGER NOT NULL DEFAULT 1,
      suppressedUntil DATETIME,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type RouteFlapDamping struct {
      ID                string     `json:"id"`
      PopID             string     `json:"popId"`
      FlapCount         int        `json:"flapCount"`
      PenaltyScore      float64    `json:"penaltyScore"`
      IsRouteSuppressed bool       `json:"isRouteSuppressed"`
      HasDecayApplied   bool       `json:"hasDecayApplied"`
      SuppressedUntil   *time.Time `json:"suppressedUntil"`
      UpdatedAt         time.Time  `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforced in `pkg/warmup/anycast/flap_guard.go`.

---

### Improvement 1006: Anycast Egress IP Reverse DNS (PTR) Sync & Alignment Validator
- **ID**: `IMP-1006`
- **Cluster**: 101
- **Goal**: Validates Forward-Confirmed Reverse DNS (FCrDNS) across all Anycast edge egress nodes to ensure 100% PTR-to-A/AAAA alignment before dispatching warming envelopes.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss validation matrix with green checkmarks (`✓`) for aligned PTR records and amber crosses (`✗`) for mismatched hostnames.
  - *Web*: `<fcrdns-alignment-card>` Web Component displaying forward and reverse DNS lookup badges.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast validate-fcrdns --ip-pool=<pool-name> [--is-strict-enforce] [--has-ipv6-evaluated]`
  - Flag `--ip-pool`: Name of the egress IP pool to validate.
  - Flag `--is-strict-enforce`: Blocks dispatch from any node failing forward-confirmed reverse DNS.
  - Flag `--has-ipv6-evaluated`: Evaluates both IPv4 and IPv6 Anycast egress addresses.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS FcrdnsValidations (
      id TEXT PRIMARY KEY,
      egressIpAddress TEXT NOT NULL,
      ptrHostname TEXT NOT NULL,
      forwardIpAddress TEXT NOT NULL,
      isFcrdnsAligned INTEGER NOT NULL DEFAULT 1,
      hasIpv6Evaluated INTEGER NOT NULL DEFAULT 0,
      validationDurationMs REAL NOT NULL DEFAULT 0.0,
      verifiedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type FcrdnsValidation struct {
      ID                   string    `json:"id"`
      EgressIPAddress      string    `json:"egressIpAddress"`
      PtrHostname          string    `json:"ptrHostname"`
      ForwardIPAddress     string    `json:"forwardIpAddress"`
      IsFcrdnsAligned      bool      `json:"isFcrdnsAligned"`
      HasIpv6Evaluated     bool      `json:"hasIpv6Evaluated"`
      ValidationDurationMs float64   `json:"validationDurationMs"`
      VerifiedAt           time.Time `json:"verifiedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Audited in `pkg/warmup/anycast/fcrdns.go`.

---

### Improvement 1007: Edge-Terminated TLS 1.3 Pre-Warming & Session Ticket Cache
- **ID**: `IMP-1007`
- **Cluster**: 101
- **Goal**: Pre-establishes and caches TLS 1.3 session tickets between Anycast edge gateways and major ISP MX servers to achieve zero-RTT (0-RTT) resumption on subsequent warming deliveries.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss cache status display showing ticket expiration countdowns and 0-RTT handshake savings in neon green (`#00E676`).
  - *Web*: `<tls-ticket-cache-view>` Web Component visualizing session ticket lifetime histograms.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast prewarm-tls --target-mx=<host> [--ticket-ttl=<seconds>] [--is-early-data-enabled] [--has-session-cached]`
  - Flag `--target-mx`: Destination MX hostname to pre-warm TLS session with.
  - Flag `--ticket-ttl`: Cache duration for retrieved session tickets (default: `86400`).
  - Flag `--is-early-data-enabled`: Enables TLS 1.3 0-RTT early data transmission.
  - Flag `--has-session-cached`: Queries local memory cache before initiating full handshake.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TlsSessionTickets (
      id TEXT PRIMARY KEY,
      targetMxHost TEXT NOT NULL,
      edgePopId TEXT NOT NULL,
      sessionTicketId TEXT NOT NULL UNIQUE,
      cipherSuite TEXT NOT NULL DEFAULT 'TLS_AES_256_GCM_SHA384',
      isEarlyDataEnabled INTEGER NOT NULL DEFAULT 1,
      hasSessionCached INTEGER NOT NULL DEFAULT 1,
      expiresAt DATETIME NOT NULL,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TlsSessionTicket struct {
      ID                 string    `json:"id"`
      TargetMxHost       string    `json:"targetMxHost"`
      EdgePopID          string    `json:"edgePopId"`
      SessionTicketID    string    `json:"sessionTicketId"`
      CipherSuite        string    `json:"cipherSuite"`
      IsEarlyDataEnabled bool      `json:"isEarlyDataEnabled"`
      HasSessionCached   bool      `json:"hasSessionCached"`
      ExpiresAt          time.Time `json:"expiresAt"`
      CreatedAt          time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/anycast/tls_prewarm.go`.

---

### Improvement 1008: Real-Time BGP Community-Based Route Tagging & Traffic Engineering
- **ID**: `IMP-1008`
- **Cluster**: 101
- **Goal**: Injects custom BGP standard and large communities into edge egress routing daemons (BIRD/FRR) to control upstream ISP transit preference and geographic announcement boundaries.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss BGP community table highlighting AS path prepends and community tags in electric purple (`#D500F9`).
  - *Web*: `<bgp-community-tagger>` Web Component rendering interactive community tag chips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast tag-community --community-tag=<tag> [--pop-id=<id>] [--is-export-active] [--has-prepend-applied]`
  - Flag `--community-tag`: RFC 1997 / RFC 8092 community string (e.g., `65001:100`).
  - Flag `--pop-id`: Target Anycast node to apply community policy.
  - Flag `--is-export-active`: Immediately exports route with new community to upstream peers.
  - Flag `--has-prepend-applied`: Appends AS path prepending to adjust return path asymmetry.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS BgpCommunityTags (
      id TEXT PRIMARY KEY,
      popId TEXT NOT NULL,
      communityTag TEXT NOT NULL,
      prependCount INTEGER NOT NULL DEFAULT 0,
      isExportActive INTEGER NOT NULL DEFAULT 1,
      hasPrependApplied INTEGER NOT NULL DEFAULT 0,
      appliedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type BgpCommunityTag struct {
      ID                string    `json:"id"`
      PopID             string    `json:"popId"`
      CommunityTag      string    `json:"communityTag"`
      PrependCount      int       `json:"prependCount"`
      IsExportActive    bool      `json:"isExportActive"`
      HasPrependApplied bool      `json:"hasPrependApplied"`
      AppliedAt         time.Time `json:"appliedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/anycast/bgp_tagger.go`.

---

### Improvement 1009: Geo-Distributed MX Latency Heatmap Generator
- **ID**: `IMP-1009`
- **Cluster**: 101
- **Goal**: Aggregates continuous TCP SYN, TLS handshake, and SMTP banner latencies across global edge nodes into a unified multidimensional matrix to steer traffic to the highest-performing edge paths.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss TrueColor ASCII world map displaying regional latency hotspots using ANSI color gradients from deep blue (`#001040`) to fiery orange (`#FF6D00`).
  - *Web*: `<geo-latency-heatmap>` Web Component rendering interactive Mercator projection map with SVG heat circles.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast generate-heatmap [--provider=<name>] [--output-format=<text|svg>] [--is-realtime-stream] [--has-outliers-filtered]`
  - Flag `--provider`: Filter heatmap by Mailbox Provider (e.g., `google`, `microsoft`).
  - Flag `--output-format`: Rendering format (`text` for ANSI TUI, `svg` for standalone vector).
  - Flag `--is-realtime-stream`: Keeps TUI active with 1-second refresh intervals.
  - Flag `--has-outliers-filtered`: Applies statistical interquartile range (IQR) filtering.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MxLatencyHeatmaps (
      id TEXT PRIMARY KEY,
      providerName TEXT NOT NULL,
      regionCode TEXT NOT NULL,
      averageRttMs REAL NOT NULL DEFAULT 0.0,
      p95LatencyMs REAL NOT NULL DEFAULT 0.0,
      p99LatencyMs REAL NOT NULL DEFAULT 0.0,
      isRealtimeStream INTEGER NOT NULL DEFAULT 0,
      hasOutliersFiltered INTEGER NOT NULL DEFAULT 1,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MxLatencyHeatmap struct {
      ID                  string    `json:"id"`
      ProviderName        string    `json:"providerName"`
      RegionCode          string    `json:"regionCode"`
      AverageRttMs        float64   `json:"averageRttMs"`
      P95LatencyMs        float64   `json:"p95LatencyMs"`
      P99LatencyMs        float64   `json:"p99LatencyMs"`
      IsRealtimeStream    bool      `json:"isRealtimeStream"`
      HasOutliersFiltered bool      `json:"hasOutliersFiltered"`
      SampledAt           time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Implemented in `pkg/warmup/anycast/heatmap.go`.

---

### Improvement 1010: Edge Health Beacon & Automatic Drain Orchestrator
- **ID**: `IMP-1010`
- **Cluster**: 101
- **Goal**: Emits sub-second health beacons from each Anycast edge node; upon detecting edge packet drop or ISP rate throttling, automatically drains in-flight SMTP traffic to alternate POPs.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss node status indicator with pulsing heart symbols (`♥`) in emerald green (`#00E676`) and flashing drain warnings in amber (`#FFD600`).
  - *Web*: `<edge-drain-hud>` Web Component visualizing traffic migration vectors in SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup anycast drain-node --pop-id=<id> [--drain-reason=<reason>] [--is-force-drain] [--has-traffic-rerouted]`
  - Flag `--pop-id`: Identifier of the edge node to drain.
  - Flag `--drain-reason`: Diagnostic justification for initiating traffic drain.
  - Flag `--is-force-drain`: Terminates active sockets immediately without waiting for SMTP quit.
  - Flag `--has-traffic-rerouted`: Automatically reassigns spooled envelopes to nearest healthy POP.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EdgeDrainEvents (
      id TEXT PRIMARY KEY,
      popId TEXT NOT NULL,
      drainReason TEXT NOT NULL,
      inFlightEnvelopes INTEGER NOT NULL DEFAULT 0,
      isForceDrain INTEGER NOT NULL DEFAULT 0,
      hasTrafficRerouted INTEGER NOT NULL DEFAULT 1,
      reroutedPopId TEXT,
      drainedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EdgeDrainEvent struct {
      ID                 string    `json:"id"`
      PopID              string    `json:"popId"`
      DrainReason        string    `json:"drainReason"`
      InFlightEnvelopes  int       `json:"inFlightEnvelopes"`
      IsForceDrain       bool      `json:"isForceDrain"`
      HasTrafficRerouted bool      `json:"hasTrafficRerouted"`
      ReroutedPopID      string    `json:"reroutedPopId"`
      DrainedAt          time.Time `json:"drainedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Audited in `pkg/warmup/anycast/drain.go`.

---

# Cluster 102: Autonomous Multipath TCP (MPTCP) Socket Bonding & Failover Matrix (1011–1020)

---

### Improvement 1011: Linux Kernel MPTCP v1 Netlink Path Manager & Subflow Orchestrator
- **ID**: `IMP-1011`
- **Cluster**: 102 (Autonomous Multipath TCP (MPTCP) Socket Bonding & Failover Matrix)
- **Goal**: Programmatically manages Linux kernel MPTCP path manager endpoints via Netlink IPC (`SOL_MPTCP`), dynamically attaching secondary network interfaces to active SMTP sockets.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss Netlink routing status card displaying bonded network interface names (`eth0`, `eth1`, `wwan0`), assigned IP subflows, and flags in cyan (`#00E5FF`).
  - *Web*: `<mptcp-subflow-table>` Web Component displaying subflow bandwidth allocation bars in SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp configure-endpoints --primary-iface=<name> --secondary-iface=<name> [--is-backup-path] [--has-signal-flag]`
  - Flag `--primary-iface`: Primary network interface (e.g., `eth0`).
  - Flag `--secondary-iface`: Redundant network interface (e.g., `eth1`).
  - Flag `--is-backup-path`: Configures secondary interface strictly as a backup subflow.
  - Flag `--has-signal-flag`: Emits MPTCP `ADD_ADDR` option to destination peer.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MptcpEndpoints (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL UNIQUE,
      ipAddress TEXT NOT NULL,
      endpointId INTEGER NOT NULL DEFAULT 1,
      isBackupPath INTEGER NOT NULL DEFAULT 0,
      hasSignalFlag INTEGER NOT NULL DEFAULT 1,
      isEndpointActive INTEGER NOT NULL DEFAULT 1,
      createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MptcpEndpoint struct {
      ID               string    `json:"id"`
      InterfaceName    string    `json:"interfaceName"`
      IPAddress        string    `json:"ipAddress"`
      EndpointID       int       `json:"endpointId"`
      IsBackupPath     bool      `json:"isBackupPath"`
      HasSignalFlag    bool      `json:"hasSignalFlag"`
      IsEndpointActive bool      `json:"isEndpointActive"`
      CreatedAt        time.Time `json:"createdAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Configured in `pkg/warmup/mptcp/netlink.go`.

---

### Improvement 1012: Zero-Drop TCP Subflow Failover & Packet-Level RTT Balancer
- **ID**: `IMP-1012`
- **Cluster**: 102
- **Goal**: Actively balances outbound SMTP envelope bytes across bonded subflows according to smoothed RTT (SRTT), executing instantaneous packet retransmission on alternate subflows if one interface stalls.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dual-channel bandwidth stream with real-time packet distribution percentages and failover blips in bright yellow (`#FFEA00`).
  - *Web*: `<subflow-rtt-balancer>` Web Component displaying live packet distribution pie slices.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp balance-subflows --session-id=<uuid> [--scheduling-policy=<srtt|roundrobin>] [--is-dynamic-rebalance] [--has-packet-cloning]`
  - Flag `--session-id`: Active MPTCP warming session UUID.
  - Flag `--scheduling-policy`: Algorithm for distributing packet data (default: `srtt`).
  - Flag `--is-dynamic-rebalance`: Continuously shifts traffic weights based on millisecond RTT changes.
  - Flag `--has-packet-cloning`: Redundantly duplicates critical SMTP commands (`MAIL FROM`, `RCPT TO`) over both interfaces for zero-delay delivery.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MptcpSubflowStats (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL,
      subflowId INTEGER NOT NULL,
      interfaceName TEXT NOT NULL,
      smoothedRttMs REAL NOT NULL DEFAULT 0.0,
      packetsTransmitted INTEGER NOT NULL DEFAULT 0,
      packetsRetransmitted INTEGER NOT NULL DEFAULT 0,
      isDynamicRebalance INTEGER NOT NULL DEFAULT 1,
      hasPacketCloning INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MptcpSubflowStat struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      SubflowID            int       `json:"subflowId"`
      InterfaceName        string    `json:"interfaceName"`
      SmoothedRttMs        float64   `json:"smoothedRttMs"`
      PacketsTransmitted   int64     `json:"packetsTransmitted"`
      PacketsRetransmitted int64     `json:"packetsRetransmitted"`
      IsDynamicRebalance   bool      `json:"isDynamicRebalance"`
      HasPacketCloning     bool      `json:"hasPacketCloning"`
      UpdatedAt            time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Monitored in `pkg/warmup/mptcp/balancer.go`.

---

### Improvement 1013: Heterogeneous Interface Transport Monitor (Fiber vs 5G/LTE)
- **ID**: `IMP-1013`
- **Cluster**: 102
- **Goal**: Monitors jitter, bufferbloat, and signal degradation across physically distinct uplinks (e.g. wired fiber vs cellular 5G), throttling non-deterministic cellular links during heavy burst phases.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss side-by-side uplink telemetry panel displaying carrier name, RSSI/RSRP, bufferbloat grade (A+ through F) in neon green and fiery orange.
  - *Web*: `<uplink-telemetry-card>` Web Component rendering real-time signal strength bars in SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp monitor-uplinks [--interface=<name>] [--is-cellular-throttle-active] [--has-bufferbloat-test]`
  - Flag `--interface`: Specific uplink interface to diagnose.
  - Flag `--is-cellular-throttle-active`: Throttles packet volume on high-jitter cellular uplinks.
  - Flag `--has-bufferbloat-test`: Injects continuous micro-probes to measure latency under load.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS UplinkTelemetry (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL UNIQUE,
      mediumType TEXT NOT NULL DEFAULT 'ethernet',
      bufferbloatLatencyMs REAL NOT NULL DEFAULT 0.0,
      jitterMs REAL NOT NULL DEFAULT 0.0,
      packetLossPercent REAL NOT NULL DEFAULT 0.0,
      isCellularThrottleActive INTEGER NOT NULL DEFAULT 0,
      hasBufferbloatTest INTEGER NOT NULL DEFAULT 1,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type UplinkTelemetry struct {
      ID                       string    `json:"id"`
      InterfaceName            string    `json:"interfaceName"`
      MediumType               string    `json:"mediumType"`
      BufferbloatLatencyMs     float64   `json:"bufferbloatLatencyMs"`
      JitterMs                 float64   `json:"jitterMs"`
      PacketLossPercent        float64   `json:"packetLossPercent"`
      IsCellularThrottleActive bool      `json:"isCellularThrottleActive"`
      HasBufferbloatTest       bool      `json:"hasBufferbloatTest"`
      ProbedAt                 time.Time `json:"probedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Tracked in `pkg/warmup/mptcp/uplink.go`.

---

### Improvement 1014: Kernel eBPF MPTCP Socket Telemetry & TCP Trace Ingestion
- **ID**: `IMP-1014`
- **Cluster**: 102
- **Goal**: Ingests low-overhead eBPF kprobes on `mptcp_subflow_create` and `tcp_retransmit_skb` to capture microsecond kernel events without user-space context-switch overhead.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss real-time event ticker displaying kernel subflow lifecycle events and socket buffer allocations in electric purple (`#7C4DFF`).
  - *Web*: `<ebpf-socket-ticker>` Web Component rendering interactive event stream logs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp attach-ebpf [--ring-buffer-kb=<size>] [--is-kprobe-attached] [--has-stack-trace]`
  - Flag `--ring-buffer-kb`: Size in kilobytes of the shared perf ring buffer (default: `1024`).
  - Flag `--is-kprobe-attached`: Attaches eBPF program directly to kernel symbols.
  - Flag `--has-stack-trace`: Collects kernel call stacks for retransmitted packets.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS EbpfSocketEvents (
      id TEXT PRIMARY KEY,
      eventName TEXT NOT NULL,
      socketFd INTEGER NOT NULL,
      subflowId INTEGER NOT NULL,
      kernelLatencyUs INTEGER NOT NULL DEFAULT 0,
      isKprobeAttached INTEGER NOT NULL DEFAULT 1,
      hasStackTrace INTEGER NOT NULL DEFAULT 0,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type EbpfSocketEvent struct {
      ID               string    `json:"id"`
      EventName        string    `json:"eventName"`
      SocketFd         int       `json:"socketFd"`
      SubflowID        int       `json:"subflowId"`
      KernelLatencyUs  int64     `json:"kernelLatencyUs"`
      IsKprobeAttached bool      `json:"isKprobeAttached"`
      HasStackTrace    bool      `json:"hasStackTrace"`
      CapturedAt       time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Loaded in `pkg/warmup/mptcp/ebpf_loader.go`.

---

### Improvement 1015: Middlebox MPTCP Stripping Detector & Fallback Handshake Engine
- **ID**: `IMP-1015`
- **Cluster**: 102
- **Goal**: Detects hostile intermediate ISP middleboxes that strip MPTCP `MP_CAPABLE` options from TCP SYN packets, executing seamless zero-delay fallback to standard TCP without connection teardown.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss warning banner showing middlebox inspection indicators and fallback status in bold orange (`#FF6D00`).
  - *Web*: `<middlebox-stripping-card>` Web Component displaying TCP option preservation status.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp detect-stripping --target-host=<host> [--is-fallback-ready] [--has-stripping-detected]`
  - Flag `--target-host`: Destination host or MX to probe for MPTCP option preservation.
  - Flag `--is-fallback-ready`: Pre-allocates standard TCP socket structures for instantaneous fallback.
  - Flag `--has-stripping-detected`: Logs middlebox stripping occurrence to audit database.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MiddleboxStrippingEvents (
      id TEXT PRIMARY KEY,
      targetHost TEXT NOT NULL,
      targetPort INTEGER NOT NULL DEFAULT 25,
      isFallbackReady INTEGER NOT NULL DEFAULT 1,
      hasStrippingDetected INTEGER NOT NULL DEFAULT 0,
      fallbackDurationUs INTEGER NOT NULL DEFAULT 0,
      detectedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MiddleboxStrippingEvent struct {
      ID                   string    `json:"id"`
      TargetHost           string    `json:"targetHost"`
      TargetPort           int       `json:"targetPort"`
      IsFallbackReady      bool      `json:"isFallbackReady"`
      HasStrippingDetected bool      `json:"hasStrippingDetected"`
      FallbackDurationUs   int64     `json:"fallbackDurationUs"`
      DetectedAt           time.Time `json:"detectedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Audited in `pkg/warmup/mptcp/fallback.go`.

---

### Improvement 1016: Subflow Congestion Window (CWND) Decoupler & BBRv3 Optimizer
- **ID**: `IMP-1016`
- **Cluster**: 102
- **Goal**: Decouples congestion control algorithms across individual subflows, allowing primary fiber links to utilize BBRv3 while cellular failover links employ Reno/CUBIC to prevent queuing collapse.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss dual-CWND meter displaying pacing rate, round-trip min RTT, and in-flight bytes in vibrant cyan (`#00E5FF`) and gold (`#FFD600`).
  - *Web*: `<cwnd-decoupler-view>` Web Component rendering real-time CWND trajectory line graphs.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp configure-cc --subflow-id=<id> --cc-algorithm=<bbr3|cubic|reno> [--is-decoupled] [--has-pacing-enforced]`
  - Flag `--subflow-id`: Specific subflow ID to apply congestion control policy.
  - Flag `--cc-algorithm`: Congestion control variant (default: `bbr3`).
  - Flag `--is-decoupled`: Prevents secondary subflow packet drops from reducing primary CWND.
  - Flag `--has-pacing-enforced`: Enforces microsecond packet pacing via Linux `SO_MAX_PACING_RATE`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SubflowCongestionConfigs (
      id TEXT PRIMARY KEY,
      subflowId INTEGER NOT NULL UNIQUE,
      ccAlgorithm TEXT NOT NULL DEFAULT 'bbr3',
      currentCwndPackets INTEGER NOT NULL DEFAULT 10,
      pacingRateBytesPerSec INTEGER NOT NULL DEFAULT 0,
      isDecoupled INTEGER NOT NULL DEFAULT 1,
      hasPacingEnforced INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SubflowCongestionConfig struct {
      ID                    string    `json:"id"`
      SubflowID             int       `json:"subflowId"`
      CcAlgorithm           string    `json:"ccAlgorithm"`
      CurrentCwndPackets    int       `json:"currentCwndPackets"`
      PacingRateBytesPerSec int64     `json:"pacingRateBytesPerSec"`
      IsDecoupled           bool      `json:"isDecoupled"`
      HasPacingEnforced     bool      `json:"hasPacingEnforced"`
      UpdatedAt             time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Managed in `pkg/warmup/mptcp/congestion.go`.

---

### Improvement 1017: MPTCP Out-of-Order Packet Buffer Reassembly & Memory Guard
- **ID**: `IMP-1017`
- **Cluster**: 102
- **Goal**: Provides bounded out-of-order packet reassembly ring buffers for heterogeneous subflows with diverging latencies, preventing memory bloat during multi-interface delivery.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss memory buffer gauge showing allocated socket memory vs. maximum safety watermark in bright mint (`#69F0AE`).
  - *Web*: `<reassembly-buffer-gauge>` Web Component rendering circular progress rings in SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp set-reassembly-limit [--max-buffer-kb=<kb>] [--is-strict-limit] [--has-early-trim]`
  - Flag `--max-buffer-kb`: Maximum memory permitted for out-of-order reassembly (default: `4096`).
  - Flag `--is-strict-limit`: Drops delayed subflow packets when limit is breached.
  - Flag `--has-early-trim`: Actively trims duplicate TCP segments before buffering.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS ReassemblyBufferStats (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL UNIQUE,
      allocatedMemoryBytes INTEGER NOT NULL DEFAULT 0,
      maxBufferLimitBytes INTEGER NOT NULL DEFAULT 4194304,
      isStrictLimit INTEGER NOT NULL DEFAULT 1,
      hasEarlyTrim INTEGER NOT NULL DEFAULT 1,
      droppedPacketsCount INTEGER NOT NULL DEFAULT 0,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type ReassemblyBufferStat struct {
      ID                   string    `json:"id"`
      SessionID            string    `json:"sessionId"`
      AllocatedMemoryBytes int64     `json:"allocatedMemoryBytes"`
      MaxBufferLimitBytes  int64     `json:"maxBufferLimitBytes"`
      IsStrictLimit        bool      `json:"isStrictLimit"`
      HasEarlyTrim         bool      `json:"hasEarlyTrim"`
      DroppedPacketsCount  int       `json:"droppedPacketsCount"`
      UpdatedAt            time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Controlled in `pkg/warmup/mptcp/reassembly.go`.

---

### Improvement 1018: Subflow Health Heartbeat & Dead-Path Eviction Controller
- **ID**: `IMP-1018`
- **Cluster**: 102
- **Goal**: Transmits sub-second cryptographic keepalive probes across idle MPTCP subflows; instantly evicts dead or degraded physical paths before SMTP write commands block.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss heartbeat matrix showing ping reply intervals per subflow with status tags `ALIVE`, `DEGRADED`, and `EVICTED` in neon green and fiery red.
  - *Web*: `<subflow-heartbeat-grid>` Web Component rendering status matrix cards.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp probe-health [--subflow-id=<id>] [--interval-ms=<ms>] [--is-auto-evict] [--has-probe-acknowledged]`
  - Flag `--subflow-id`: Specific subflow identifier to probe.
  - Flag `--interval-ms`: Keepalive probe frequency in milliseconds (default: `250`).
  - Flag `--is-auto-evict`: Removes unacknowledged paths from socket session.
  - Flag `--has-probe-acknowledged`: Verifies return token before marking path healthy.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SubflowHealthProbes (
      id TEXT PRIMARY KEY,
      subflowId INTEGER NOT NULL,
      probeRound INTEGER NOT NULL DEFAULT 1,
      rttUs INTEGER NOT NULL DEFAULT 0,
      isSubflowAlive INTEGER NOT NULL DEFAULT 1,
      isAutoEvict INTEGER NOT NULL DEFAULT 1,
      hasProbeAcknowledged INTEGER NOT NULL DEFAULT 1,
      probedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SubflowHealthProbe struct {
      ID                   string    `json:"id"`
      SubflowID            int       `json:"subflowId"`
      ProbeRound           int       `json:"probeRound"`
      RttUs                int64     `json:"rttUs"`
      IsSubflowAlive       bool      `json:"isSubflowAlive"`
      IsAutoEvict          bool      `json:"isAutoEvict"`
      HasProbeAcknowledged bool      `json:"hasProbeAcknowledged"`
      ProbedAt             time.Time `json:"probedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Enforced in `pkg/warmup/mptcp/heartbeat.go`.

---

### Improvement 1019: Multi-Carrier Mobile IPv6 Dynamic Prefix Re-Binder
- **ID**: `IMP-1019`
- **Cluster**: 102
- **Goal**: Dynamically re-binds outbound MPTCP subflows to newly acquired IPv6 prefixes when mobile cellular modems roam across base stations, preserving live SMTP sessions.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss IPv6 prefix update banner displaying old vs. new IPv6 delegated prefix and binding confirmation in lavender (`#B388FF`).
  - *Web*: `<ipv6-prefix-rebinder>` Web Component visualizing network prefix handover transitions.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp rebind-ipv6 --interface=<name> --new-prefix=<cidr> [--is-seamless-handover] [--has-addr-announced]`
  - Flag `--interface`: Mobile interface undergoing prefix reconfiguration.
  - Flag `--new-prefix`: Newly assigned IPv6 prefix.
  - Flag `--is-seamless-handover`: Retains old subflow until new prefix handshake completes.
  - Flag `--has-addr-announced`: Emits MPTCP `ADD_ADDR` option to destination server.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS Ipv6PrefixRebinds (
      id TEXT PRIMARY KEY,
      interfaceName TEXT NOT NULL,
      oldPrefix TEXT NOT NULL,
      newPrefix TEXT NOT NULL,
      isSeamlessHandover INTEGER NOT NULL DEFAULT 1,
      hasAddrAnnounced INTEGER NOT NULL DEFAULT 1,
      transitionDurationMs REAL NOT NULL DEFAULT 0.0,
      reboundAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type Ipv6PrefixRebind struct {
      ID                   string    `json:"id"`
      InterfaceName        string    `json:"interfaceName"`
      OldPrefix            string    `json:"oldPrefix"`
      NewPrefix            string    `json:"newPrefix"`
      IsSeamlessHandover   bool      `json:"isSeamlessHandover"`
      HasAddrAnnounced     bool      `json:"hasAddrAnnounced"`
      TransitionDurationMs float64   `json:"transitionDurationMs"`
      ReboundAt            time.Time `json:"reboundAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Handled in `pkg/warmup/mptcp/ipv6_rebind.go`.

---

### Improvement 1020: MPTCP vs Single-Path TCP Comparative Benchmark HUD
- **ID**: `IMP-1020`
- **Cluster**: 102
- **Goal**: Measures and visually displays side-by-side performance metrics comparing MPTCP multi-interface bonding against standard single-path TCP across identical warming batches.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss split-screen benchmark HUD comparing throughput (msgs/sec), connection drop rate (0.0% vs. 1.8%), and p99 delivery latency in TrueColor.
  - *Web*: `<mptcp-benchmark-hud>` Web Component rendering comparative bar charts in pure SVG.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup mptcp run-benchmark --batch-size=<n> [--destination-domain=<domain>] [--is-synthetic-load] [--has-report-persisted]`
  - Flag `--batch-size`: Number of envelopes to transmit during benchmark (default: `100`).
  - Flag `--destination-domain`: Target recipient domain under test.
  - Flag `--is-synthetic-load`: Uses synthetic seed accounts to isolate network variables.
  - Flag `--has-report-persisted`: Persists comparative audit summary to `repodb/audit.db`.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS MptcpBenchmarkResults (
      id TEXT PRIMARY KEY,
      batchSize INTEGER NOT NULL,
      mptcpThroughputPerSec REAL NOT NULL DEFAULT 0.0,
      singlePathThroughputPerSec REAL NOT NULL DEFAULT 0.0,
      mptcpDropPercent REAL NOT NULL DEFAULT 0.0,
      singlePathDropPercent REAL NOT NULL DEFAULT 0.0,
      isSyntheticLoad INTEGER NOT NULL DEFAULT 1,
      hasReportPersisted INTEGER NOT NULL DEFAULT 1,
      benchmarkedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type MptcpBenchmarkResult struct {
      ID                         string    `json:"id"`
      BatchSize                  int       `json:"batchSize"`
      MptcpThroughputPerSec      float64   `json:"mptcpThroughputPerSec"`
      SinglePathThroughputPerSec float64   `json:"singlePathThroughputPerSec"`
      MptcpDropPercent           float64   `json:"mptcpDropPercent"`
      SinglePathDropPercent      float64   `json:"singlePathDropPercent"`
      IsSyntheticLoad            bool      `json:"isSyntheticLoad"`
      HasReportPersisted         bool      `json:"hasReportPersisted"`
      BenchmarkedAt              time.Time `json:"benchmarkedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Audited in `pkg/warmup/mptcp/benchmark.go`.

---

# Cluster 103: Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI (1021–1030)

---

### Improvement 1021: Sub-Millisecond Phase Flamegraph TUI Renderer (Bubble Tea v2)
- **ID**: `IMP-1021`
- **Cluster**: 103 (Real-Time Bounded Jitter Waveform & Latency Flamegraph TUI)
- **Goal**: Renders an interactive, sub-millisecond phase flamegraph in Bubble Tea v2 decomposing SMTP delivery time across DNS, TCP connect, TLS handshake, greeting, EHLO, MAIL/RCPT, and DATA phases.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss flamegraph viewport utilizing Unicode block elements (`▏▎▍▌▋▊▉█`) with distinct colors: DNS (cyan), TCP (blue), TLS (purple), SMTP dialog (green), and DATA transfer (amber).
  - *Web*: `<smtp-phase-flamegraph>` Web Component generating SVG stack rectangles with interactive hover micro-tooltips.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry render-flamegraph --session-id=<uuid> [--zoom-level=<1-5>] [--is-interactive] [--has-svg-export]`
  - Flag `--session-id`: Identifier of the warming session to inspect.
  - Flag `--zoom-level`: Magnification scale for microsecond phase breakdown (default: `1`).
  - Flag `--is-interactive`: Launches Bubble Tea v2 scrollable interactive HUD.
  - Flag `--has-svg-export`: Exports rendered flamegraph as an standalone SVG file.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpPhaseTimings (
      id TEXT PRIMARY KEY,
      sessionId TEXT NOT NULL UNIQUE,
      dnsDurationMs REAL NOT NULL DEFAULT 0.0,
      tcpDurationMs REAL NOT NULL DEFAULT 0.0,
      tlsDurationMs REAL NOT NULL DEFAULT 0.0,
      smtpDialogDurationMs REAL NOT NULL DEFAULT 0.0,
      dataTransferDurationMs REAL NOT NULL DEFAULT 0.0,
      totalDurationMs REAL NOT NULL DEFAULT 0.0,
      isInteractive INTEGER NOT NULL DEFAULT 1,
      hasSvgExport INTEGER NOT NULL DEFAULT 0,
      recordedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpPhaseTiming struct {
      ID                     string    `json:"id"`
      SessionID              string    `json:"sessionId"`
      DnsDurationMs          float64   `json:"dnsDurationMs"`
      TcpDurationMs          float64   `json:"tcpDurationMs"`
      TlsDurationMs          float64   `json:"tlsDurationMs"`
      SmtpDialogDurationMs   float64   `json:"smtpDialogDurationMs"`
      DataTransferDurationMs float64   `json:"dataTransferDurationMs"`
      TotalDurationMs        float64   `json:"totalDurationMs"`
      IsInteractive          bool      `json:"isInteractive"`
      HasSvgExport           bool      `json:"hasSvgExport"`
      RecordedAt             time.Time `json:"recordedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Rendered in `pkg/warmup/telemetry/flamegraph.go`.

---

### Improvement 1022: Live TCP Jitter Waveform Oscilloscope & ANSI Sparkline Stream
- **ID**: `IMP-1022`
- **Cluster**: 103
- **Goal**: Streams real-time TCP packet jitter as an animated Unicode oscilloscope waveform (` ▂▃▄▅▆▇█`) at 60 FPS directly inside the developer's terminal during warm-up dispatches.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Smooth 60 FPS scrolling waveform stream in Lipgloss with gradient mapping from deep emerald green (`#00E676`) for stable jitter to fiery red (`#FF1744`) for jitter spikes > 50ms.
  - *Web*: `<jitter-waveform-scope>` Web Component rendering real-time SVG waveform curves via Shadow DOM.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry stream-waveform [--fps=<30|60>] [--window-seconds=<n>] [--is-color-gradient] [--has-stats-overlay]`
  - Flag `--fps`: Render frame rate in terminal (default: `60`).
  - Flag `--window-seconds`: Temporal observation window in seconds (default: `30`).
  - Flag `--is-color-gradient`: Applies multi-hue TrueColor styling based on jitter intensity.
  - Flag `--has-stats-overlay`: Overlays mean jitter and standard deviation badges.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS JitterWaveformSamples (
      id TEXT PRIMARY KEY,
      sampleIndex INTEGER NOT NULL,
      jitterValueMs REAL NOT NULL DEFAULT 0.0,
      colorHexCode TEXT NOT NULL DEFAULT '#00E676',
      isColorGradient INTEGER NOT NULL DEFAULT 1,
      hasStatsOverlay INTEGER NOT NULL DEFAULT 1,
      sampledAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type JitterWaveformSample struct {
      ID               string    `json:"id"`
      SampleIndex      int       `json:"sampleIndex"`
      JitterValueMs    float64   `json:"jitterValueMs"`
      ColorHexCode     string    `json:"colorHexCode"`
      IsColorGradient  bool      `json:"isColorGradient"`
      HasStatsOverlay  bool      `json:"hasStatsOverlay"`
      SampledAt        time.Time `json:"sampledAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Streamed in `pkg/warmup/telemetry/waveform.go`.

---

### Improvement 1023: Zero-CDN Shadow DOM SVG Delivery Radar Component
- **ID**: `IMP-1023`
- **Cluster**: 103
- **Goal**: Provides an embeddable, zero-external-CDN HTML Web Component `<delivery-radar>` that visualizes regional delivery progress, bounce ratios, and throttle dampening on an SVG radar canvas.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss circular ASCII radar display using Braille characters (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`) to represent sweep radar beams.
  - *Web*: Self-contained Custom Element `<delivery-radar>` encapsulated inside Shadow DOM with inline SVG geometry and pure CSS animations.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry export-radar --output-html=<path> [--is-standalone] [--has-shadow-dom-encapsulated]`
  - Flag `--output-html`: Destination file path for generated HTML report.
  - Flag `--is-standalone`: Includes all CSS/SVG inline with zero external script references.
  - Flag `--has-shadow-dom-encapsulated`: Wraps component in isolated Shadow DOM root.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DeliveryRadarConfigs (
      id TEXT PRIMARY KEY,
      targetDomain TEXT NOT NULL,
      sweepIntervalMs INTEGER NOT NULL DEFAULT 500,
      isStandalone INTEGER NOT NULL DEFAULT 1,
      hasShadowDomEncapsulated INTEGER NOT NULL DEFAULT 1,
      radarDataJson TEXT NOT NULL,
      generatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DeliveryRadarConfig struct {
      ID                       string    `json:"id"`
      TargetDomain             string    `json:"targetDomain"`
      SweepIntervalMs          int       `json:"sweepIntervalMs"`
      IsStandalone             bool      `json:"isStandalone"`
      HasShadowDomEncapsulated bool      `json:"hasShadowDomEncapsulated"`
      RadarDataJSON            string    `json:"radarDataJson"`
      GeneratedAt              time.Time `json:"generatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Exported via `pkg/warmup/telemetry/radar.go`.

---

### Improvement 1024: Bubble Tea v2 Interactive Multi-Tab Warmup Dashboard
- **ID**: `IMP-1024`
- **Cluster**: 103
- **Goal**: Delivers an all-in-one terminal dashboard in Bubble Tea v2 with tabs for Overview, Anycast Routing, MPTCP Bonding, Throttling Calibration, and Audit Logs navigable via keyboard keys `1` through `5` or `Tab`.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss styled multi-tab bar with active tab in bold white on royal blue (`#1A237E`) and inactive tabs in muted slate (`#37474F`). Smooth scrolling viewport with vim keybindings (`j`/`k`).
  - *Web*: Complementary zero-CDN responsive dashboard view with identical color tokens.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup ui dashboard [--initial-tab=<overview|anycast|mptcp|throttle|audit>] [--refresh-rate-ms=<ms>] [--is-mouse-enabled] [--has-help-bar]`
  - Flag `--initial-tab`: Active tab upon launching UI (default: `overview`).
  - Flag `--refresh-rate-ms`: Polling frequency for background telemetry updates (default: `1000`).
  - Flag `--is-mouse-enabled`: Enables terminal mouse clicking and scrolling.
  - Flag `--has-help-bar`: Displays bottom keyboard shortcuts helper bar.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS DashboardUserPreferences (
      id TEXT PRIMARY KEY,
      activeTab TEXT NOT NULL DEFAULT 'overview',
      refreshRateMs INTEGER NOT NULL DEFAULT 1000,
      isMouseEnabled INTEGER NOT NULL DEFAULT 1,
      hasHelpBar INTEGER NOT NULL DEFAULT 1,
      lastOpenedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type DashboardUserPreference struct {
      ID             string    `json:"id"`
      ActiveTab      string    `json:"activeTab"`
      RefreshRateMs  int       `json:"refreshRateMs"`
      IsMouseEnabled bool      `json:"isMouseEnabled"`
      HasHelpBar     bool      `json:"hasHelpBar"`
      LastOpenedAt   time.Time `json:"lastOpenedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Invoked via `pkg/warmup/telemetry/dashboard.go`.

---

### Improvement 1025: TrueColor Latency Percentile Histogram (p50/p90/p99)
- **ID**: `IMP-1025`
- **Cluster**: 103
- **Goal**: Computes and displays dynamic latency distribution histograms across thousands of warmup deliveries, visually highlighting p50, p90, p95, and p99 percentile thresholds.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Vertical Lipgloss bar chart with gradient color bins; p50 in emerald (`#00E676`), p90 in yellow (`#FFD600`), p95 in amber (`#FF6D00`), and p99 in crimson (`#D50000`).
  - *Web*: `<latency-histogram-view>` Web Component rendering interactive SVG bar buckets.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry show-histogram [--bucket-count=<n>] [--is-logarithmic-scale] [--has-percentiles-labeled]`
  - Flag `--bucket-count`: Number of histogram distribution buckets (default: `20`).
  - Flag `--is-logarithmic-scale`: Renders histogram on a log10 scale for extreme tail latency.
  - Flag `--has-percentiles-labeled`: Draws horizontal guide lines for p50/p90/p99 marks.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS LatencyHistograms (
      id TEXT PRIMARY KEY,
      batchId TEXT NOT NULL,
      p50Ms REAL NOT NULL DEFAULT 0.0,
      p90Ms REAL NOT NULL DEFAULT 0.0,
      p95Ms REAL NOT NULL DEFAULT 0.0,
      p99Ms REAL NOT NULL DEFAULT 0.0,
      isLogarithmicScale INTEGER NOT NULL DEFAULT 0,
      hasPercentilesLabeled INTEGER NOT NULL DEFAULT 1,
      bucketDataJson TEXT NOT NULL,
      computedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type LatencyHistogram struct {
      ID                    string    `json:"id"`
      BatchID               string    `json:"batchId"`
      P50Ms                 float64   `json:"p50Ms"`
      P90Ms                 float64   `json:"p90Ms"`
      P95Ms                 float64   `json:"p95Ms"`
      P99Ms                 float64   `json:"p99Ms"`
      IsLogarithmicScale    bool      `json:"isLogarithmicScale"`
      HasPercentilesLabeled bool      `json:"hasPercentilesLabeled"`
      BucketDataJSON        string    `json:"bucketDataJson"`
      ComputedAt            time.Time `json:"computedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Formatted in `pkg/warmup/telemetry/histogram.go`.

---

### Improvement 1026: Outbound SMTP Command Trace Inspector with Microsecond Timestamps
- **ID**: `IMP-1026`
- **Cluster**: 103
- **Goal**: Captures and renders every raw wire-level SMTP command and reply (`EHLO`, `STARTTLS`, `MAIL FROM`, `RCPT TO`, `DATA`, `QUIT`) with microsecond delta timestamps to inspect protocol stalling.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss wire trace stream with client requests prefixed by `>>>` in cyan and server responses prefixed by `<<<` in green/yellow/red depending on reply code (2xx/4xx/5xx).
  - *Web*: `<smtp-wire-inspector>` Web Component rendering syntax-highlighted protocol transcripts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry inspect-trace --envelope-id=<uuid> [--is-raw-hex-included] [--has-delta-timing]`
  - Flag `--envelope-id`: Specific envelope ID to retrieve raw wire trace for.
  - Flag `--is-raw-hex-included`: Dumps hex octets alongside ASCII commands for TLS debugging.
  - Flag `--has-delta-timing`: Calculates elapsed microseconds between successive commands.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SmtpWireTraces (
      id TEXT PRIMARY KEY,
      envelopeId TEXT NOT NULL,
      direction TEXT NOT NULL CHECK(direction IN ('CLIENT_TO_SERVER', 'SERVER_TO_CLIENT')),
      commandOrReply TEXT NOT NULL,
      elapsedDeltaUs INTEGER NOT NULL DEFAULT 0,
      isRawHexIncluded INTEGER NOT NULL DEFAULT 0,
      hasDeltaTiming INTEGER NOT NULL DEFAULT 1,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SmtpWireTrace struct {
      ID               string    `json:"id"`
      EnvelopeID       string    `json:"envelopeId"`
      Direction        string    `json:"direction"`
      CommandOrReply   string    `json:"commandOrReply"`
      ElapsedDeltaUs   int64     `json:"elapsedDeltaUs"`
      IsRawHexIncluded bool      `json:"isRawHexIncluded"`
      HasDeltaTiming   bool      `json:"hasDeltaTiming"`
      CapturedAt       time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Persisted in `pkg/warmup/telemetry/wire_trace.go`.

---

### Improvement 1027: Terminal Notification Dispatcher for Delivery Anomalies
- **ID**: `IMP-1027`
- **Cluster**: 103
- **Goal**: Emits non-intrusive ANSI bell triggers (`\a`) and desktop/terminal popups using OSC 777 / OSC 9 notifications when sudden greylisting bursts or connection drops are detected.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss floating alert notification badge anchored in top-right of the terminal with glowing borders in electric red (`#FF1744`).
  - *Web*: Web desktop notification integration with native Notification API without external scripts.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry configure-alerts [--error-rate-threshold=<percent>] [--is-ansi-bell-active] [--has-desktop-notification]`
  - Flag `--error-rate-threshold`: Trigger notification if error rate exceeds threshold (default: `5.0`).
  - Flag `--is-ansi-bell-active`: Emits ASCII bell chime to terminal speaker.
  - Flag `--has-desktop-notification`: Triggers OS desktop notifications via OSC escape sequences.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS AnomalyAlertConfigs (
      id TEXT PRIMARY KEY,
      errorRateThreshold REAL NOT NULL DEFAULT 5.0,
      isAnsiBellActive INTEGER NOT NULL DEFAULT 1,
      hasDesktopNotification INTEGER NOT NULL DEFAULT 1,
      isAlertMuted INTEGER NOT NULL DEFAULT 0,
      lastTriggeredAt DATETIME
  );
  ```
  ```go
  type AnomalyAlertConfig struct {
      ID                     string     `json:"id"`
      ErrorRateThreshold     float64    `json:"errorRateThreshold"`
      IsAnsiBellActive       bool       `json:"isAnsiBellActive"`
      HasDesktopNotification bool       `json:"hasDesktopNotification"`
      IsAlertMuted           bool       `json:"isAlertMuted"`
      LastTriggeredAt        *time.Time `json:"lastTriggeredAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Controlled in `pkg/warmup/telemetry/alerts.go`.

---

### Improvement 1028: Standalone SVG Waterfall Exporter with Zero External Assets
- **ID**: `IMP-1028`
- **Cluster**: 103
- **Goal**: Compiles an entire warm-up batch execution timeline into an immutable, single-file vector SVG waterfall document containing embedded CSS fonts and metadata, completely viewable offline.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss progress bar displaying vector export generation and resulting file size in kilobytes.
  - *Web*: Standalone pure SVG document with CSS `@media (prefers-color-scheme: dark)` styling.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry export-svg-waterfall --batch-id=<uuid> --output-file=<path> [--is-dark-theme] [--has-embedded-metadata]`
  - Flag `--batch-id`: Identifier of the batch to export.
  - Flag `--output-file`: Destination `.svg` file path.
  - Flag `--is-dark-theme`: Applies dark-mode color scheme by default.
  - Flag `--has-embedded-metadata`: Embeds JSON delivery receipt metadata inside SVG `<metadata>` tags.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/audit.db`
  ```sql
  CREATE TABLE IF NOT EXISTS SvgWaterfallExports (
      id TEXT PRIMARY KEY,
      batchId TEXT NOT NULL,
      exportFilePath TEXT NOT NULL,
      fileSizeBytes INTEGER NOT NULL DEFAULT 0,
      isDarkTheme INTEGER NOT NULL DEFAULT 1,
      hasEmbeddedMetadata INTEGER NOT NULL DEFAULT 1,
      exportedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type SvgWaterfallExport struct {
      ID                  string    `json:"id"`
      BatchID             string    `json:"batchId"`
      ExportFilePath      string    `json:"exportFilePath"`
      FileSizeBytes       int64     `json:"fileSizeBytes"`
      IsDarkTheme         bool      `json:"isDarkTheme"`
      HasEmbeddedMetadata bool      `json:"hasEmbeddedMetadata"`
      ExportedAt          time.Time `json:"exportedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Audited in `pkg/warmup/telemetry/waterfall_exporter.go`.

---

### Improvement 1029: Bubble Tea v2 Live Egress Queue Topology Graph
- **ID**: `IMP-1029`
- **Cluster**: 103
- **Goal**: Visualizes in-flight message queues across local spoolers, Anycast edge POPs, and remote MX ingress buffers using an interactive ASCII directed graph rendered via Bubble Tea v2.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss directed graph layout using box-drawing characters (`┌─┐│└─┘►`) with node fill percentages and flow arrows flashing during active packet transit.
  - *Web*: `<queue-topology-graph>` Web Component displaying animated SVG connection paths.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry show-queue-topology [--is-animated] [--has-bandwidth-metrics] [--refresh-ms=<ms>]`
  - Flag `--is-animated`: Animates packet flow along ASCII arrows.
  - Flag `--has-bandwidth-metrics`: Overlays current transfer bitrate per link.
  - Flag `--refresh-ms`: Screen refresh interval in milliseconds (default: `250`).
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS QueueTopologyNodes (
      id TEXT PRIMARY KEY,
      nodeName TEXT NOT NULL UNIQUE,
      nodeType TEXT NOT NULL CHECK(nodeType IN ('LOCAL_SPOOL', 'EDGE_POP', 'REMOTE_MX')),
      currentQueueDepth INTEGER NOT NULL DEFAULT 0,
      maxQueueCapacity INTEGER NOT NULL DEFAULT 10000,
      isAnimated INTEGER NOT NULL DEFAULT 1,
      hasBandwidthMetrics INTEGER NOT NULL DEFAULT 1,
      updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type QueueTopologyNode struct {
      ID                  string    `json:"id"`
      NodeName            string    `json:"nodeName"`
      NodeType            string    `json:"nodeType"`
      CurrentQueueDepth   int       `json:"currentQueueDepth"`
      MaxQueueCapacity    int       `json:"maxQueueCapacity"`
      IsAnimated          bool      `json:"isAnimated"`
      HasBandwidthMetrics bool      `json:"hasBandwidthMetrics"`
      UpdatedAt           time.Time `json:"updatedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Rendered in `pkg/warmup/telemetry/topology_graph.go`.

---

### Improvement 1030: Telemetry Snapshot Time-Travel Scrubber & Replay Controller
- **ID**: `IMP-1030`
- **Cluster**: 103
- **Goal**: Provides an interactive temporal scrubbing HUD allowing operators to scrub backwards through past warming runs second-by-second to inspect exactly when and why an ISP throttle or disconnect occurred.
- **UI/UX Terminal & Web Specification**:
  - *Terminal*: Lipgloss timeline scrubber bar with cursor handle (`◆`) controllable via left/right arrow keys or `h`/`l`, displaying instantaneous system state at the selected timestamp.
  - *Web*: `<time-travel-scrubber>` Web Component with interactive SVG timeline slider.
- **CLI Command Signature & Flag Definitions**:
  - `gitmap email warmup telemetry time-travel --batch-id=<uuid> [--start-offset-seconds=<n>] [--is-replay-loop] [--has-scrubber-locked]`
  - Flag `--batch-id`: Identifier of historical batch to scrub through.
  - Flag `--start-offset-seconds`: Initial timeline position offset in seconds.
  - Flag `--is-replay-loop`: Plays back recorded session telemetry in real time.
  - Flag `--has-scrubber-locked`: Locks UI viewport during automated playback.
- **Data Contract & Split-DB Schema**:
  - Storage: `repodb/warmup.db`
  ```sql
  CREATE TABLE IF NOT EXISTS TimeTravelSnapshots (
      id TEXT PRIMARY KEY,
      batchId TEXT NOT NULL,
      offsetSeconds INTEGER NOT NULL,
      activeConnections INTEGER NOT NULL DEFAULT 0,
      inFlightEnvelopes INTEGER NOT NULL DEFAULT 0,
      throttlingEngaged INTEGER NOT NULL DEFAULT 0,
      isReplayLoop INTEGER NOT NULL DEFAULT 0,
      hasScrubberLocked INTEGER NOT NULL DEFAULT 0,
      stateSnapshotJson TEXT NOT NULL,
      capturedAt DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  ```
  ```go
  type TimeTravelSnapshot struct {
      ID                string    `json:"id"`
      BatchID           string    `json:"batchId"`
      OffsetSeconds     int       `json:"offsetSeconds"`
      ActiveConnections int       `json:"activeConnections"`
      InFlightEnvelopes int       `json:"inFlightEnvelopes"`
      ThrottlingEngaged bool      `json:"throttlingEngaged"`
      IsReplayLoop      bool      `json:"isReplayLoop"`
      HasScrubberLocked bool      `json:"hasScrubberLocked"`
      StateSnapshotJSON string    `json:"stateSnapshotJson"`
      CapturedAt        time.Time `json:"capturedAt"`
  }
  ```
- **Antigravity / GitMap Integration Point**: Controlled in `pkg/warmup/telemetry/time_travel.go`.
