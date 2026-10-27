# Specification 339: Improvements 261 to 300 — Failover, CSS Design Tokens, Split-DB WAL & AGY Event Bus

- **Sequence:** 339
- **Slug:** `339-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/04-improvements-261-to-300-failover-tokens-and-telemetry.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification authoritatively details Improvements 261 through 300 for Sequence 339 of the GitMap Warm-Up and Emailing UI Engine. Expanding upon Web Component design primitives, time-travel event debugging, automated rescue daemons, and interactive terminal studios, these 40 concrete specifications complete the platform's multi-IP network fault tolerance, declarative token styling, high-throughput SQLite engine tuning, and autonomous agent coordination:

1. **Cluster 27: Multi-IP Egress Failover & ISP Circuit Breakers (Items 261–270):** Outbound network multi-homing across dynamic IP pools, health-probed TCP socket failover, automated ASN-aware routing tables, RFC 3463 circuit state machines, and zero-downtime hot-swap egress transitions.
2. **Cluster 28: Design System CSS Custom Properties & High-Contrast Modes (Items 271–280):** Mathematical typographic scales, CSS custom property hierarchies, WCAG 2.1 AAA high-contrast themes, focus visible rings, reduced motion physics adapters, and theme token runtime injectors.
3. **Cluster 29: Split-DB Storage WAL Tuning, Memory Cache & Log Compaction (Items 281–290):** Dual SQLite split-architecture (`warmup.db` and `mail_audit.db`), WAL auto-checkpoint tuning, page cache memory allocation, Bloom filter lookups, vacuum log compaction, and deterministic Merkle tree audit logging.
4. **Cluster 30: Antigravity Autonomous Agent Event Bus & Real-Time Alerts (Items 291–300):** Pub/sub WebSocket message broker, autonomous agent execution hooks, delivery anomaly triage alerts, procedural audio engine synthesizers, terminal bell event routers, and GitMap CLI command listeners.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 339 Improvements 261 to 300 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 27 (261-270)       │ Cluster 28 (271-280)       │ Cluster 29 (281-290)         │
│ Multi-IP Egress Failover   │ Design System CSS Custom   │ Split-DB Storage WAL Tuning, │
│ & ISP Circuit Breakers     │ Properties & High Contrast │ Memory Cache & Compaction    │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 30 (291-300)                                    │ Cross-Cutting Standards      │
│ Antigravity Autonomous Agent Event Bus, WebSockets      │ Positive Booleans, Split-DB, │
│ & Real-Time Telemetry Alerts                            │ Zero Exits, Zero CDNs        │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 27: Multi-IP Egress Failover & ISP Circuit Breakers (Items 261–270)

### Improvement 261: Dynamic Outbound Multi-IP Egress Pool & Local Interface Binder
- **Description:** Multi-homed network socket orchestrator capable of binding outbound SMTP TCP connections dynamically across heterogeneous local network interfaces and dedicated IP blocks.
- **UI/Terminal Component:** IP interface matrix card displaying active outbound addresses, interface speed, subnet mask, and live connection gauge. CLI inspection via `gitmap mail egress list`.
- **Telemetry & State Fields:** `hasAssignedIp` (bool), `isInterfaceActive` (bool), `isDefaultRoute` (bool), `interfaceName` (string), `ipAddress` (string), `macAddress` (string), `activeConnectionCount` (int).
- **Backend / Engine:** `pkg/mail/egress/interface_binder.go` configuring `net.Dialer.LocalAddr` using explicit `net.TCPAddr` targets with zero overhead.
- **Failure / Guard Rail:** If specified IP is missing or unassigned at OS level, bubbles typed `*apperror.AppError` and falls back to primary default gateway interface.

### Improvement 262: Continuous Per-IP Egress Health Prober & Latency Radar
- **Description:** Background synthetic diagnostic loop issuing continuous non-intrusive TCP handshake probes (`SYN`/`SYN-ACK`) against known ISP MX gateways (Google, Microsoft, Yahoo) to determine IP-level packet loss and round-trip latency.
- **UI/Terminal Component:** Latency radar pulse chart with millisecond gauges and historical jitter sparklines.
- **Telemetry & State Fields:** `isProbeHealthy` (bool), `hasPacketLoss` (bool), `isHighLatency` (bool), `egressIp` (string), `probeTargetMx` (string), `latencyMs` (int64), `packetLossPercentage` (float64).
- **Backend / Engine:** `pkg/mail/egress/health_prober.go` running non-blocking goroutine ticks every 30 seconds with sliding window latency smoothing.
- **Failure / Guard Rail:** Automatically marks IP as degraded if 3 consecutive handshake probes exceed 1500ms or encounter TCP timeout.

### Improvement 263: Three-State ISP Circuit Breaker & Automatic Trip Engine
- **Description:** Robust finite state machine (FSM) implementing `Closed`, `Open`, and `Half-Open` circuit states per ISP domain. Detects anomalous connection drops, handshake resets, and RFC 3463 throttling codes.
- **UI/Terminal Component:** Visual circuit breaker toggle widget with animated state transitions (Emerald for Closed, Crimson for Open, Amber pulse for Half-Open).
- **Telemetry & State Fields:** `isCircuitClosed` (bool), `isCircuitOpen` (bool), `isHalfOpen` (bool), `hasTripped` (bool), `ispTarget` (string), `consecutiveErrorCount` (int), `lastStateChangeTimestamp` (int64).
- **Backend / Engine:** `pkg/mail/circuit/isp_breaker.go` coordinating state switches via atomic pointers and mutex-guarded state transitions.
- **Failure / Guard Rail:** Instantly trips circuit to `Open` when failure threshold (5 errors in 60s) is breached, rejecting outbound dispatches to that ISP.

### Improvement 264: Circuit Breaker Half-Open Canary Probe & Self-Healing Ramp
- **Description:** Controlled recovery state manager for tripped circuit breakers. Upon cooldown expiration, permits a strict quota of canary test dispatches to verify ISP gateway acceptance before fully re-opening the circuit.
- **UI/Terminal Component:** Canary status progress bar indicating permitted canary count and successful acknowledgement ratio.
- **Telemetry & State Fields:** `isCanaryActive` (bool), `hasCanaryPassed` (bool), `canarySuccessCount` (int), `canaryTargetQuota` (int), `cooldownRemainingSeconds` (int).
- **Backend / Engine:** `pkg/mail/circuit/canary_probe.go` calculating success percentages and resetting circuit state to `Closed` after 10 clean canary deliveries.
- **Failure / Guard Rail:** A single `5.X.X` rejection or connection timeout during canary trial immediately re-trips circuit to `Open` and doubles cooldown timer.

### Improvement 265: Real-Time Multi-IP Automatic Egress Failover Dispatcher
- **Description:** Outbound dispatch router that dynamically reassigns pending deliveries away from degraded or tripped IP addresses to pristine alternative IPs in the pool with zero message loss.
- **UI/Terminal Component:** Failover route redirect animation showing source IP rerouting traffic lines to target fallback IP.
- **Telemetry & State Fields:** `isFailoverTriggered` (bool), `hasAvailableFallback` (bool), `sourceIp` (string), `targetIp` (string), `reroutedMessageCount` (int), `failoverTimestamp` (int64).
- **Backend / Engine:** `pkg/mail/egress/failover_router.go` performing thread-safe queue reallocation in `warmup.db`.
- **Failure / Guard Rail:** If entire IP pool is exhausted or degraded, halts dispatches and schedules backoff pause instead of corrupting sender domain reputation.

### Improvement 266: ASN-Aware Gateway Routing Table & Geo-Proximity Optimizer
- **Description:** Dynamic routing matrix mapping destination recipient domains to their authoritative Autonomous System Numbers (ASNs) and selecting the optimal egress IP interface minimizing BGP autonomous system hops.
- **UI/Terminal Component:** Interactive ASN routing breakdown table with BGP peer indicators and hop count badges.
- **Telemetry & State Fields:** `isRouteOptimized` (bool), `hasDirectPeering` (bool), `destinationAsn` (int), `selectedEgressIp` (string), `estimatedHopCount` (int).
- **Backend / Engine:** `pkg/mail/egress/asn_router.go` reading local lightweight ASN lookup tables with zero external API latency.
- **Failure / Guard Rail:** Unknown ASNs default to standard weighted round-robin distribution across active egress pool members.

### Improvement 267: Per-IP Warming Quota Shaper & Egress Throttling Lock
- **Description:** Individual IP reputation protection governor that strictly enforces independent daily volume caps and per-minute throughput limits across every IP in the egress pool.
- **UI/Terminal Component:** Per-IP quota fill gauge with visual tick indicators and threshold safety locks.
- **Telemetry & State Fields:** `isQuotaExceeded` (bool), `hasAvailableQuota` (bool), `ipAddress` (string), `dailyDispatchedCount` (int), `dailyMaxCap` (int), `quotaUtilizationPercentage` (float64).
- **Backend / Engine:** `pkg/mail/throttle/ip_quota_shaper.go` utilizing sliding-window atomic counters and SQLite quota snapshots.
- **Failure / Guard Rail:** Blocks further dispatch attempts for the specific IP when capacity hits 100%, rerouting queued jobs to remaining pool IPs.

### Improvement 268: Egress Socket Pool Reuse & Keep-Alive Connection Manager
- **Description:** High-efficiency TCP socket pooling and keep-alive connection manager maintaining warmed persistent TLS connections to high-volume ISP MX hosts to amortize handshake overhead.
- **UI/Terminal Component:** Active connection pool monitor displaying idle, busy, and total sockets per destination host.
- **Telemetry & State Fields:** `isConnectionReused` (bool), `hasIdleConnection` (bool), `targetHost` (string), `idleSocketCount` (int), `activeSocketCount` (int), `idleTimeoutSeconds` (int).
- **Backend / Engine:** `pkg/mail/egress/socket_pool.go` implementing bounded FIFO queue with graceful idle reaper routines.
- **Failure / Guard Rail:** Prunes stale or broken TCP connections proactively before dispatching mail commands to avoid `Broken pipe` or `Connection reset` failures.

### Improvement 269: Egress Interface Hot-Plug & Dynamic IP Re-Discovery Watcher
- **Description:** Operating system network change watcher that detects added, removed, or updated IP interfaces without requiring daemon or service restart.
- **UI/Terminal Component:** System network change notification toast with detected IP delta badges (`+IP Added`, `-IP Removed`).
- **Telemetry & State Fields:** `hasInterfaceChanged` (bool), `isInterfaceAdded` (bool), `isInterfaceRemoved` (bool), `modifiedInterfaceName` (string), `discoveredIpList` ([]string).
- **Backend / Engine:** `pkg/mail/egress/interface_watcher.go` polling network interfaces and diffing against in-memory registry.
- **Failure / Guard Rail:** Safely drains active jobs from vanishing interfaces before purging them from routing tables.

### Improvement 270: CLI Egress Pool Command Suite & Failover Simulator
- **Description:** Complete command-line operational toolkit providing operators with granular controls to inspect, add, drain, remove, or simulate failover events across IP egress pools.
- **UI/Terminal Component:** Rich Lipgloss terminal table output via `gitmap mail egress status` and interactive prompt `gitmap mail egress simulate --trip`.
- **Telemetry & State Fields:** `isCommandSuccessful` (bool), `hasSimulatedTrip` (bool), `commandName` (string), `affectedIp` (string), `simulatedIsp` (string).
- **Backend / Engine:** `cli/cmdmail/cmd_egress.go` delegating to internal egress packages and formatting structured output envelopes.
- **Failure / Guard Rail:** Simulation commands require explicit `--confirm` flag to prevent accidental production circuit tripping.

---

## 3. Cluster 28: Design System CSS Custom Properties & High-Contrast Modes (Items 271–280)

### Improvement 271: Declarative Design Token System & Multi-Tier Custom Properties
- **Description:** Centralized CSS custom property token dictionary organizing colors, elevations, spacing units, border radiuses, and transitions into semantic tiers (`--sys-*`, `--comp-*`).
- **UI/Terminal Component:** Visual design token inspector previewing color swatches, elevation depths, and typography styles.
- **Telemetry & State Fields:** `hasLoadedTokens` (bool), `isCustomThemeActive` (bool), `tokenTier` (string), `activeTokenCount` (int), `themeVersion` (string).
- **Backend / Engine:** `web/static/css/tokens.css` with native CSS variables and Go-embedded string constants in `pkg/web/tokens/styles.go`.
- **Failure / Guard Rail:** Missing custom variables automatically cascade to hardcoded fallback declarations directly inside CSS declarations.

### Improvement 272: Mathematical Typographic Scale & Fluid Clamp Formulations
- **Description:** Responsive typography system governed by mathematical modular ratios (Minor Third, Major Second) combined with viewport-based `clamp()` expressions for seamless fluid scaling.
- **UI/Terminal Component:** Typography scale specimen showcase demonstrating continuous fluid scaling from mobile viewport to 4K displays.
- **Telemetry & State Fields:** `isFluidScalingActive` (bool), `hasValidRatio` (bool), `baseFontSizePx` (float64), `modularRatio` (float64), `currentScaleStep` (string).
- **Backend / Engine:** `web/static/css/typography.css` generating zero-calc CSS rules with pure relative rem/clamp units.
- **Failure / Guard Rail:** Caps upper and lower bounds strictly to prevent unreadable text truncation or excessive font expansion on ultra-wide monitors.

### Improvement 273: WCAG 2.1 AAA High-Contrast Accessibility Theme
- **Description:** Dedicated ultra-high-contrast display mode guaranteeing minimum 7:1 contrast ratio across all foreground text and interactive UI elements against background surfaces.
- **UI/Terminal Component:** High-contrast toggle pill with immediate visual contrast boost and sharp 2px solid boundary indicators.
- **Telemetry & State Fields:** `isHighContrastActive` (bool), `hasAccessibleContrast` (bool), `measuredContrastRatio` (float64), `activeForegroundHex` (string), `activeBackgroundHex` (string).
- **Backend / Engine:** `web/static/css/themes/high_contrast.css` defining distinct pure black (`#000000`) and pure white (`#ffffff`) token overrides.
- **Failure / Guard Rail:** Automatically activates if system reports `prefers-contrast: more` media feature during initial page boot.

### Improvement 274: High-Visibility Focus Indicators & Keyboard Focus Rings
- **Description:** Accessible keyboard navigation engine providing distinctive dual-ring focus outlines (`outline: 2px solid #06b6d4; outline-offset: 3px;`) conforming to Section 508.
- **UI/Terminal Component:** Animated focus ring appearing dynamically upon Tab navigation across cards, inputs, and action buttons.
- **Telemetry & State Fields:** `isKeyboardFocused` (bool), `hasFocusVisible` (bool), `focusedElementTag` (string), `focusRingColor` (string).
- **Backend / Engine:** Pure CSS `:focus-visible` pseudo-class rules in `web/static/css/accessibility.css` eliminating focus outline suppression bugs.
- **Failure / Guard Rail:** Suppresses focus ring during mouse clicks to avoid visual clutter while strictly guaranteeing keyboard focus visibility.

### Improvement 275: Reduced Motion Physics Mode & Accessible Transition Limiter
- **Description:** Motion accessibility controller that detects user OS motion preferences and completely replaces heavy transitions and physics simulations with instantaneous state switches.
- **UI/Terminal Component:** Motion sensitivity toggle in user preferences with one-click disable option for all canvas physics loops.
- **Telemetry & State Fields:** `isReducedMotionActive` (bool), `hasAnimationDisabled` (bool), `transitionDurationMs` (int), `physicsMode` (string).
- **Backend / Engine:** Media query `@media (prefers-reduced-motion: reduce)` rules and JS state toggles in `web/static/js/accessibility.js`.
- **Failure / Guard Rail:** Sets all CSS animation durations and transition times to `0.001ms` when reduced motion is detected.

### Improvement 276: Glassmorphism Blur Depth Engine & GPU Backdrop Filter Layering
- **Description:** Modern multi-plane visual depth system employing hardware-accelerated CSS `backdrop-filter: blur(16px)` and multi-stop alpha gradients for floating telemetry cards.
- **UI/Terminal Component:** Translucent glass cards with luminous edge highlights and ambient drop shadows.
- **Telemetry & State Fields:** `isGpuAccelerated` (bool), `hasBackdropSupport` (bool), `blurRadiusPx` (int), `surfaceOpacity` (float64).
- **Backend / Engine:** `web/static/css/surfaces.css` applying hardware composition layers via `transform: translateZ(0)`.
- **Failure / Guard Rail:** Automatically falls back to opaque high-contrast solid backgrounds when browser lacks `backdrop-filter` support.

### Improvement 277: Cyberpunk Emerald & Cyan Telemetry Accent System
- **Description:** Cohesive data visualization color palette utilizing vibrant Emerald (`#10b981`), Cyan (`#06b6d4`), Amber (`#f59e0b`), and Rose (`#f43f5e`) accents optimized for dark obsidian backgrounds.
- **UI/Terminal Component:** Color-coded status badges, glow borders, and pulsating telemetry LEDs across dashboard panels.
- **Telemetry & State Fields:** `hasValidPalette` (bool), `isPaletteCohesive` (bool), `accentEmerald` (string), `accentCyan` (string), `accentRose` (string).
- **Backend / Engine:** Integrated CSS palette definitions in `web/static/css/colors.css`.
- **Failure / Guard Rail:** Ensures all accent colors maintain minimum 4.5:1 contrast against dark background tokens under all viewing conditions.

### Improvement 278: CSS Variable Runtime Token Injector & Live Preset Switcher
- **Description:** Dynamic JavaScript token manager capable of injecting and hot-swapping CSS variables at the `:root` level at runtime without reloading page assets.
- **UI/Terminal Component:** Floating theme picker menu with instant live swatch previewing and one-click preset activation.
- **Telemetry & State Fields:** `isThemeApplied` (bool), `hasCustomPreset` (bool), `themePresetName` (string), `injectedTokenCount` (int).
- **Backend / Engine:** `web/static/js/theme_manager.js` manipulating `document.documentElement.style.setProperty()` with zero DOM re-rendering.
- **Failure / Guard Rail:** Restores saved default theme from `localStorage` if an invalid or corrupted theme preset is selected.

### Improvement 279: Print Stylesheet & Clean PDF Report Exporter Layout
- **Description:** Dedicated CSS `@media print` layout formatting warm-up deliverability dashboards into pristine, high-resolution multi-page PDF documents.
- **UI/Terminal Component:** "Export PDF Report" button generating clean white-background printable reports with data tables.
- **Telemetry & State Fields:** `isPrintModeActive` (bool), `hasFormattedPageBreaks` (bool), `pageCount` (int), `reportGeneratedTimestamp` (int64).
- **Backend / Engine:** `web/static/css/print.css` stripping navigation bars, canvas animations, and dark background fills.
- **Failure / Guard Rail:** Enforces explicit `page-break-inside: avoid` on metrics cards to prevent awkward mid-card print splits.

### Improvement 280: Micro-Icon SVG Sprite Sheet & Zero-HTTP Icon Renderer
- **Description:** Inline SVG sprite sheet containing 40+ lightweight deliverability micro-icons (server, shield, clock, envelope, check, alert) rendered without external network requests.
- **UI/Terminal Component:** Crisp vector micro-icons embedded within status pills, buttons, and navigation breadcrumbs.
- **Telemetry & State Fields:** `isSvgRendered` (bool), `hasValidSprite` (bool), `iconIdentifier` (string), `svgByteSize` (int).
- **Backend / Engine:** `web/static/icons/sprites.svg` with native `<svg><use href="#icon-id" /></svg>` markup.
- **Failure / Guard Rail:** Missing icon identifier falls back to a generic diamond glyph with warning logged to console.

---

## 4. Cluster 29: Split-DB Storage WAL Tuning, Memory Cache & Log Compaction (Items 281–290)

### Improvement 281: Dual Split-DB SQLite Architecture & Domain Isolation
- **Description:** Strict architectural partitioning isolating warm-up schedule state (`warmup.db`) from immutable delivery audit trails and raw bounce logs (`mail_audit.db`).
- **UI/Terminal Component:** Storage health dashboard showing independent size, connection pool status, and transaction rates for both databases.
- **Telemetry & State Fields:** `hasConnectedWarmupDb` (bool), `hasConnectedAuditDb` (bool), `warmupDbSizeMb` (float64), `auditDbSizeMb` (float64), `totalActiveQueries` (int).
- **Backend / Engine:** `pkg/db/splitdb/manager.go` managing independent `*sql.DB` connection pools with strict schema boundaries.
- **Failure / Guard Rail:** Out-of-space condition on `mail_audit.db` buffers events to local memory ring buffer without halting `warmup.db` scheduler.

### Improvement 282: SQLite Write-Ahead Logging (WAL) Tuning & Pragmas Configurator
- **Description:** High-performance SQLite database connection initializer applying optimal performance pragmas: `journal_mode=WAL`, `synchronous=NORMAL`, `busy_timeout=5000`, and `temp_store=MEMORY`.
- **UI/Terminal Component:** DB Pragma status card verifying active journal mode, sync level, and page size in CLI and UI.
- **Telemetry & State Fields:** `isWalModeActive` (bool), `hasOptimalPragmas` (bool), `journalMode` (string), `synchronousMode` (string), `busyTimeoutMs` (int).
- **Backend / Engine:** `pkg/db/splitdb/pragmas.go` executing configuration statements upon every newly spawned database connection.
- **Failure / Guard Rail:** Verifies `PRAGMA journal_mode` returns `wal`; if unsupported, logs diagnostic warning and operates in standard rollback journal mode.

### Improvement 283: Automated WAL Checkpoint Governor & Frame Monitor
- **Description:** Proactive WAL file size governor monitoring uncheckpointed WAL frames and executing passive or full checkpoints (`PRAGMA wal_checkpoint(TRUNCATE)`) during idle cycles.
- **UI/Terminal Component:** WAL frame count gauge displaying pending frames and historical checkpoint durations.
- **Telemetry & State Fields:** `isCheckpointRunning` (bool), `hasWalTruncated` (bool), `uncheckpointedFrames` (int), `checkpointDurationMs` (int64), `lastCheckpointTimestamp` (int64).
- **Backend / Engine:** `pkg/db/splitdb/wal_governor.go` scheduling background maintenance loops every 5 minutes.
- **Failure / Guard Rail:** Prevents WAL file from growing unbounded beyond 64MB by escalating checkpoint mode to `RESTART` during off-peak intervals.

### Improvement 284: In-Memory Prepared Statement Cache & Query Plan Optimizer
- **Description:** Thread-safe prepared statement cache avoiding query compilation overhead for high-frequency warmup queries (status lookups, counter increments).
- **UI/Terminal Component:** Query cache hit-ratio gauge displaying cached statements and hit rates.
- **Telemetry & State Fields:** `hasCachedStatement` (bool), `isCacheHit` (bool), `statementKey` (string), `hitCount` (int64), `missCount` (int64), `cacheHitRatio` (float64).
- **Backend / Engine:** `pkg/db/splitdb/stmt_cache.go` wrapping `*sql.Stmt` instances with `sync.RWMutex` lifecycle guards.
- **Failure / Guard Rail:** Safely closes cached statements upon database teardown or schema migration events.

### Improvement 285: Bloom Filter Memory Index for Fast Suppressed Recipient Lookups
- **Description:** In-memory probabilistic Bloom filter pre-filtering recipient email suppression checks with sub-microsecond latency prior to hitting disk SQLite tables.
- **UI/Terminal Component:** Bloom filter memory indicator showing filter size (KB), item count, and false positive probability.
- **Telemetry & State Fields:** `isMaybeSuppressed` (bool), `hasDefinitelySuppressed` (bool), `bloomItemCount` (int), `memoryFootprintBytes` (int), `falsePositiveRate` (float64).
- **Backend / Engine:** `pkg/mail/filter/bloom_filter.go` implementing MurmurHash3 multi-hash calculations in native Go.
- **Failure / Guard Rail:** Positive match in Bloom filter triggers definitive SQLite query confirmation; zero chance of false negatives.

### Improvement 286: Autonomous Vacuum Log Compaction & Segment Archiver
- **Description:** Background database maintenance worker that compacts aged audit logs older than 90 days, exports compressed gzip archives, and executes incremental vacuuming.
- **UI/Terminal Component:** Archival progress card showing archived row counts, reclaimed disk space (MB), and next scheduled run.
- **Telemetry & State Fields:** `isArchivingActive` (bool), `hasReclaimedSpace` (bool), `archivedRowCount` (int64), `reclaimedSpaceBytes` (int64), `lastVacuumTimestamp` (int64).
- **Backend / Engine:** `pkg/db/splitdb/compactor.go` streaming chunks of rows to gzip files before executing deletion queries.
- **Failure / Guard Rail:** Executes operations in batches of 1000 rows to prevent long-running table write-locks.

### Improvement 287: Cryptographic Merkle Tree Audit Ledger & Tamper Detection
- **Description:** SHA-256 Merkle tree calculation engine over `mail_audit.db` transaction logs, guaranteeing tamper-evident cryptographic verification of all warm-up dispatches.
- **UI/Terminal Component:** Cryptographic audit verification badge with live root hash display and one-click "Verify Audit Integrity" check.
- **Telemetry & State Fields:** `isLedgerValid` (bool), `hasTamperDetected` (bool), `merkleRootHash` (string), `verifiedBlockCount` (int64), `lastVerifiedTimestamp` (int64).
- **Backend / Engine:** `pkg/db/audit/merkle_tree.go` building binary leaf trees from sequential log entry hashes.
- **Failure / Guard Rail:** Any modified, injected, or removed audit log entry immediately invalidates Merkle root hash, raising high-priority security alert.

### Improvement 288: LRU Cache for Mailbox Deliverability Health Scores
- **Description:** Least-Recently-Used (LRU) memory cache storing computed mailbox health scores, reducing CPU overhead from repeated reputation recalculations.
- **UI/Terminal Component:** Health cache status pill with live capacity count and memory utilization metrics.
- **Telemetry & State Fields:** `isCachedScore` (bool), `hasHealthScore` (bool), `mailboxAddress` (string), `healthScore` (float64), `cacheTtlRemainingSeconds` (int).
- **Backend / Engine:** `pkg/mail/pool/health_cache.go` maintaining doubly linked list and hash map with 5-minute eviction policy.
- **Failure / Guard Rail:** Invalidation hooks clear cached mailbox score instantly upon receipt of any new bounce or spam complaint.

### Improvement 289: SQLite Database Schema Migration Engine & Version Ledger
- **Description:** Zero-dependency sequential migration runner applying schema DDL updates to both `warmup.db` and `mail_audit.db` inside transactional boundaries.
- **UI/Terminal Component:** Schema migration history table in terminal CLI via `gitmap db migrate status`.
- **Telemetry & State Fields:** `isSchemaCurrent` (bool), `hasPendingMigrations` (bool), `currentSchemaVersion` (int), `appliedMigrationCount` (int), `migrationDurationMs` (int64).
- **Backend / Engine:** `pkg/db/splitdb/migrator.go` reading embedded SQL scripts and logging executed versions in `schema_migrations` table.
- **Failure / Guard Rail:** Any error during migration execution automatically triggers `ROLLBACK` and leaves database at prior healthy version.

### Improvement 290: DB Diagnostic Doctor CLI & Integrity Check Runner
- **Description:** Comprehensive terminal diagnostic utility running `PRAGMA integrity_check`, `PRAGMA foreign_key_check`, index fragmentation analysis, and page allocation audits.
- **UI/Terminal Component:** Formatted Lipgloss report card printed to console upon running `gitmap db doctor`.
- **Telemetry & State Fields:** `isIntegrityValid` (bool), `hasCorruptedPages` (bool), `hasForeignKeyViolations` (bool), `fragmentationRatio` (float64), `checkedPageCount` (int64).
- **Backend / Engine:** `cli/cmddb/cmd_doctor.go` inspecting database internals and outputting remediation suggestions.
- **Failure / Guard Rail:** Non-zero corruption immediately triggers recommendation to restore from automated backup snapshots.

---

## 5. Cluster 30: Antigravity Autonomous Agent Event Bus & Real-Time Alerts (Items 291–300)

### Improvement 291: WebSocket Telemetry Event Bus & Topic Message Broker
- **Description:** High-throughput in-memory pub/sub message broker broadcasting typed warm-up events (`warmup.dispatch`, `warmup.bounce`, `warmup.circuit`) over WebSockets to UI clients and subagents.
- **UI/Terminal Component:** WebSocket connection status LED in web status bar with live message rate counter (msg/sec).
- **Telemetry & State Fields:** `isWebSocketConnected` (bool), `hasActiveSubscribers` (bool), `connectedClientCount` (int), `broadcastMessageRate` (float64), `activeTopicCount` (int).
- **Backend / Engine:** `pkg/agent/bus/event_broker.go` utilizing Go channels and client subscription maps with zero message dropping.
- **Failure / Guard Rail:** Disconnected or slow-reading WebSocket clients are automatically purged after channel buffer exceeds 1024 messages.

### Improvement 292: Autonomous Agent Task Webhook Dispatcher & Payload Signer
- **Description:** Secure webhook delivery engine dispatching cryptographically signed JSON payloads (`X-GitMap-Signature: sha256=...`) to autonomous Antigravity agents upon critical warm-up triggers.
- **UI/Terminal Component:** Webhook delivery logs panel showing endpoint URL, HTTP response status, retry attempts, and latency.
- **Telemetry & State Fields:** `isWebhookDelivered` (bool), `hasSignatureValid` (bool), `targetEndpoint` (string), `httpStatusCode` (int), `responseDurationMs` (int64).
- **Backend / Engine:** `pkg/agent/hooks/webhook_dispatcher.go` utilizing HMAC-SHA256 signing and exponential backoff retries.
- **Failure / Guard Rail:** Automatically disables misbehaving external webhook endpoints after 10 consecutive connection failures.

### Improvement 293: Anomaly Triage Engine & Real-Time Reputation Drop Alert
- **Description:** Statistical anomaly detection engine evaluating rolling 15-minute dispatch metrics. Instantly sounds alerts if bounce rate spikes >2.5x above the historical baseline.
- **UI/Terminal Component:** Urgent crimson alert banner at top of dashboard with one-click "Halt Senders & Investigate" action.
- **Telemetry & State Fields:** `hasAnomalyDetected` (bool), `isReputationDegraded` (bool), `anomalySeverity` (string), `currentMetricValue` (float64), `baselineMetricValue` (float64).
- **Backend / Engine:** `pkg/agent/anomaly/triage_engine.go` computing standard deviations and Z-scores over sliding metrics buffers.
- **Failure / Guard Rail:** Automatically throttles warm-up dispatches across all active senders by 50% when a critical anomaly is confirmed.

### Improvement 294: Procedural Web Audio API Deliverability Sound Synthesizer
- **Description:** Client-side zero-dependency audio feedback engine using the native Web Audio API to synthesize procedural chime tones for delivery success, warnings, and error alerts.
- **UI/Terminal Component:** Sound effects volume slider and mute toggle in navigation header (`web/static/js/audio_controls.js`).
- **Telemetry & State Fields:** `isAudioMuted` (bool), `hasAudioSynthesizer` (bool), `audioVolumePercent` (int), `activeOscillatorType` (string).
- **Backend / Engine:** `web/static/js/audio_synthesizer.js` constructing dynamic `OscillatorNode` and `GainNode` sound envelopes with zero audio file downloads.
- **Failure / Guard Rail:** Gracefully respects browser autoplay policies; audio remains dormant until user initiates first interactive click.

### Improvement 295: Terminal Bell Alert Router & CLI Audio Notification
- **Description:** CLI notification adapter emitting audible terminal bell frequencies (`\a`) or ASCII escape alerts for critical events when running long-running terminal monitors.
- **UI/Terminal Component:** Terminal bell indicator icon and CLI flag `--bell` on `gitmap warmup monitor`.
- **Telemetry & State Fields:** `isBellEnabled` (bool), `hasBellRung` (bool), `alertEventType` (string), `bellTriggerCount` (int).
- **Backend / Engine:** `pkg/cli/tui/bell_router.go` checking terminal TTY capabilities before writing raw bell bytes.
- **Failure / Guard Rail:** Suppresses terminal bell automatically if stdout is redirected to a pipe, file, or non-interactive subshell.

### Improvement 296: Antigravity IDE Workspace Status Bar Bridge
- **Description:** JSON-RPC local socket bridge feeding live warm-up metrics directly into the Antigravity IDE bottom status bar (showing active mailboxes, warming speed, and deliverability %).
- **UI/Terminal Component:** Compact status bar segment in Antigravity editor: `[Warmup: 99.2% | 12 Mailboxes | 45 msg/h]`.
- **Telemetry & State Fields:** `isIdeConnected` (bool), `hasStatusUpdated` (bool), `currentDeliverabilityScore` (float64), `activeSenderCount` (int), `hourlyRate` (int).
- **Backend / Engine:** `pkg/ide/bridge/status_bar.go` exposing lightweight local IPC socket for editor extension consumption.
- **Failure / Guard Rail:** Socket operates entirely asynchronously; editor UI is never blocked by warm-up background processing.

### Improvement 297: Autonomous Subagent Warming Ramp Coordinator
- **Description:** Agentic decision engine that autonomously calculates optimal warming acceleration schedules, dynamically adjusting daily dispatch targets based on deliverability signals.
- **UI/Terminal Component:** "AI Auto-Pilot" badge on warm-up control deck indicating autonomous scheduling is active.
- **Telemetry & State Fields:** `isAutonomousRampActive` (bool), `hasOptimalRampTarget` (bool), `recommendedDailyIncrease` (int), `confidenceScore` (float64).
- **Backend / Engine:** `pkg/agent/orchestrator/ramp_coordinator.go` executing heuristic rule matrices to determine next-day quota targets.
- **Failure / Guard Rail:** Agent suggestions are bounded by strict hard upper limits (e.g. maximum +20% day-over-day volume expansion).

### Improvement 298: GitMap Slash-Command CLI Dispatch Listener
- **Description:** Interactive slash-command parser enabling operators and AI agents to issue conversational management commands (`/warmup status`, `/warmup pause`, `/warmup boost`) via terminal and chat consoles.
- **UI/Terminal Component:** Autocomplete command palette and interactive prompt displaying recognized slash commands.
- **Telemetry & State Fields:** `isCommandRecognized` (bool), `hasExecutedCommand` (bool), `commandString` (string), `commandExecutionDurationMs` (int64).
- **Backend / Engine:** `pkg/agent/slash/command_listener.go` mapping slash tokens to internal CLI handler functions.
- **Failure / Guard Rail:** Invalid slash commands output formatted help card listing available syntax without throwing runtime errors.

### Improvement 299: Headless Background Daemon Manager & PID Lockfile Engine
- **Description:** Process supervisor managing background daemon execution of the warm-up and emailing engines with POSIX/Windows PID lockfile enforcement and clean SIGTERM handling.
- **UI/Terminal Component:** Daemon status indicator in CLI via `gitmap warmup daemon status` (running, stopped, pid).
- **Telemetry & State Fields:** `isDaemonRunning` (bool), `hasPidLock` (bool), `processId` (int), `uptimeSeconds` (int64), `memoryUsageBytes` (int64).
- **Backend / Engine:** `pkg/daemon/supervisor/pid_manager.go` managing file locks in standard OS temp folders.
- **Failure / Guard Rail:** Stale lockfiles from previous system crashes are detected, logged, and gracefully reclaimed during daemon startup.

### Improvement 300: Unified System Health Telemetry Aggregator & Executive Dashboard
- **Description:** Comprehensive top-level telemetry aggregator combining egress health, circuit breaker statuses, design token states, Split-DB metrics, and agent alerts into a unified real-time dashboard.
- **UI/Terminal Component:** Full-screen executive command center dashboard with synchronized metric tiles, live event waterfall, and health scores.
- **Telemetry & State Fields:** `isSystemHealthy` (bool), `hasActiveAlerts` (bool), `aggregateDeliverabilityScore` (float64), `activeEgressPoolSize` (int), `systemUptimeSeconds` (int64).
- **Backend / Engine:** `pkg/telemetry/aggregator/system_health.go` pulling state snapshots across all subsystems every 1000ms.
- **Failure / Guard Rail:** If any single subsystem reports failure, aggregate health indicator transitions to warning state without crashing telemetry feed.

---

## 6. Implementation Architecture & Data Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       GitMap Telemetry & Failover Mesh                      │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
        ┌──────────────────────────────┼──────────────────────────────┐
        ▼                              ▼                              ▼
┌───────────────────────┐      ┌───────────────────────┐      ┌───────────────────────┐
│     Egress Mesh       │      │   Split-DB Storage    │      │    Agent Event Bus    │
│  - Multi-IP Binding   │      │  - warmup.db (WAL)    │      │  - WebSocket Broker   │
│  - Circuit Breakers   │◄────►│  - mail_audit.db      │◄────►│  - Webhook Dispatch   │
│  - Latency Radar      │      │  - Merkle Tree Audit  │      │  - Anomaly Alerts     │
└───────────┬───────────┘      └───────────┬───────────┘      └───────────┬───────────┘
            │                              │                              │
            └──────────────────────────────┼──────────────────────────────┘
                                           │
                                           ▼
                       ┌───────────────────────────────────────┐
                       │   Visual View & Terminal Consoles     │
                       │   - CSS Custom Property Tokens        │
                       │   - WCAG 2.1 AAA High Contrast        │
                       │   - Web Audio Procedural Chimes       │
                       │   - `gitmap mail egress` / `db` CLI   │
                       └───────────────────────────────────────┘
```

---

## 7. Verification & Quality Invariants

Every improvement in this specification strictly adheres to GitMap core invariants:
- **Error Sovereignty:** All internal Go functions return typed `Result[T]` or `*apperror.AppError`. Zero direct `os.Exit` or `cliexit.HandleError` calls.
- **Positive Booleans Only:** Every boolean flag across all structs, APIs, and components strictly uses positive prefixes (`is*`, `has*`).
- **Function Limit:** Functions are strictly <= 15 lines max (<= 8 lines preferred). Conditionals are flattened with guard clauses.
- **Relative Paths:** All file links, asset URLs, and documentation cross-references use strictly relative repository paths.
- **Zero External CDNs:** All stylesheets, icons, fonts, and scripts are embedded directly into repository packages.
