# Specification 340: Improvements 361 to 400 — Design System Tokens, Split-DB Sharding & AGY Bus

- **Sequence:** 340
- **Slug:** `340-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/04-improvements-361-to-400-design-tokens-splitdb-agy.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/340-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/339-warm-up-ui-and-emailing-ui-improvements/04-improvements-261-to-300-failover-tokens-and-telemetry.md`
  - `02-spec/21-app/06-database-and-split-db/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification rigorously details Improvements 361 through 400 for Sequence 340 of the GitMap Warm-Up and Emailing UI Engine. Expanding upon the multi-agent collaborative deliverability swarm, cognitive sending rhythms, and self-healing circuit breakers, these 40 concrete specifications finalize the platform's multi-theme accessible design tokens, high-throughput SQLite sharding and query caching, real-time Antigravity IDE event hooks, and autonomous deliverability audit scorecards:

1. **Cluster 37: Design System Multi-Theme Palettes & Accessible Motion (Items 361–370):** Mathematical color palette ramp systems, CSS custom properties for 5 distinct operational themes, WCAG 2.1 AAA high-contrast accessibility modes, reduced-motion spring physics, and zero-CDN design token runtimes.
2. **Cluster 38: Split-DB SQLite Sharding & High-Speed Query Cache (Items 371–380):** Multi-database storage sharding (`warmup.db`, `mail_audit.db`, and `reputation_shards/`), in-memory LRU query cache, WAL pragma optimizations, Bloom filter recipient suppression lookups, and deterministic SHA-256 Merkle tree audit logging.
3. **Cluster 39: Antigravity IDE Event Bus & Dynamic Workspace Hooks (Items 381–390):** Bidirectional event broker, IDE activity hooks, real-time pub/sub WebSocket streaming, slash-command CLI triggers, audio deliverability synth feedbacks, and background daemon supervisor lifecycle.
4. **Cluster 40: Autonomous Deliverability Scorecards & Historical Audit (Items 391–400):** Machine-audited deliverability health scorecards, ISP-specific compliance radar visualizers, temporal anomaly correlation models, automated remediation playbook dispatch, and tamper-evident PDF/HTML compliance export.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│               Sequence 340 Improvements 361 to 400 Cluster Distribution               │
├────────────────────────────┬────────────────────────────┬──────────────────────────────┤
│ Cluster 37 (361-370)       │ Cluster 38 (371-380)       │ Cluster 39 (381-390)         │
│ Design System Multi-Theme  │ Split-DB SQLite Sharding   │ Antigravity IDE Event Bus,   │
│ Palettes & Accessible Motion│ & High-Speed Query Cache   │ Dynamic Hooks & Audio Synth  │
├────────────────────────────┴────────────────────────────┼──────────────────────────────┤
│ Cluster 40 (391-400)                                    │ Cross-Cutting Standards      │
│ Autonomous Deliverability Scorecards, Historical Audit   │ Positive Booleans, Split-DB, │
│ & Tamper-Evident Merkle Ledgers                          │ Zero Exits, Zero CDN         │
└─────────────────────────────────────────────────────────┴──────────────────────────────┘
```

---

## 2. Cluster 37: Design System Multi-Theme Palettes & Accessible Motion (Items 361–370)

### Improvement 361: Mathematical 10-Step Color Ramp Generator & Contrast Ratio Enforcer
- **Description:** Algorithmic color ramp synthesizer that derives a balanced 10-step tonal luminance scale (50–900) from any single brand seed hex code using OKLCH color space interpolation. Guarantees WCAG 2.1 AA (4.5:1) and AAA (7:1) contrast against dark surface tokens.
- **UI/Terminal Component:** Dynamic CSS custom properties generator injecting tokens at `:root` runtime. Terminal Lipgloss swatch preview rendered via `gitmap warmup theme inspect --contrast`.
- **Telemetry & State Fields:** `hasPassedAaaContrast` (bool), `isOklchSupported` (bool), `hasValidLuminance` (bool), `seedHex` (string), `luminanceSteps` ([]float64), `contrastRatio` (float64), `activeColorSpace` (string).
- **Backend / Engine:** `pkg/web/tokens/color_ramp.go` calculating perceptual lightness curves and generating CSS token declarations without runtime dependencies.
- **Failure / Guard Rail:** If a calculated step violates AAA contrast against the active surface background token, automatically shifts luminance value by ±8% until contrast threshold is achieved.

### Improvement 362: Five Built-In Glassmorphic Operational Theme Palettes
- **Description:** Declarative theme registry supporting 5 production glassmorphic palettes: `Obsidian Night` (default deep charcoal/cyan), `Solar Flare` (warm amber/bronze), `Cyber Emerald` (matrix green/slate), `Arctic Glacier` (frost blue/steel), and `Monochrome High-Contrast` (pure black/white).
- **UI/Terminal Component:** Theme switcher pill in top-right navigation header and terminal selector via `gitmap warmup theme set <name>`. Includes live DOM token hot-reloading with zero page refresh.
- **Telemetry & State Fields:** `isThemeActive` (bool), `hasGlassmorphism` (bool), `isHighContrastTheme` (bool), `themeId` (string), `themeName` (string), `primaryAccentHex` (string), `backdropBlurPx` (int).
- **Backend / Engine:** `pkg/web/tokens/theme_registry.go` serving cached CSS sheets and persisting user selection in local storage and Split-DB user preferences table.
- **Failure / Guard Rail:** Unknown theme identifier safely defaults to `Obsidian Night` and logs structured `*apperror.AppError` without panic or process exit.

### Improvement 363: WCAG 2.1 AAA High-Contrast Operational Mode
- **Description:** Dedicated accessibility toggle replacing soft glassmorphic gradients and transparency with pure, high-contrast borders (solid 2px `#FFFFFF` or `#FFFF00`), solid `#000000` surfaces, and bold sans-serif text tokens for low-vision and sunlight field visibility.
- **UI/Terminal Component:** Accessible high-contrast toggle switch with keyboard shortcut (`Alt+H`) and persistent visual badge. CLI toggle via `gitmap warmup accessibility --high-contrast=true`.
- **Telemetry & State Fields:** `isHighContrastEnabled` (bool), `hasSolidBorders` (bool), `isSunlightOptimized` (bool), `borderColorHex` (string), `contrastMultiplier` (float64), `activeA11yLevel` (string).
- **Backend / Engine:** `pkg/web/tokens/accessibility_manager.go` injecting `.theme-high-contrast` class on document root and notifying Shadow DOM components.
- **Failure / Guard Rail:** Preserves high-contrast state across browser restarts and terminal sessions via persistent preference keys.

### Improvement 364: Accessible Focus Rings & Keyboard Navigation Traversal Grid
- **Description:** High-visibility, offset focus ring design token (`--ring-focus: 3px solid #38BDF8; --ring-offset: 2px`) applying strict WCAG 2.4.7 focus visibility to all interactive buttons, cards, time scrubbers, and tabs. Supports bidirectional roving tabindex navigation.
- **UI/Terminal Component:** Animated glowing focus halo that tracks active focus position with sub-pixel alignment. Visual indicator on all interactive warm-up dials and seed cards.
- **Telemetry & State Fields:** `isFocusVisible` (bool), `hasKeyboardFocus` (bool), `isRovingTabindexActive` (bool), `focusedElementId` (string), `tabSequenceIndex` (int), `ringOffsetPx` (int).
- **Backend / Engine:** `pkg/web/components/focus_trap.go` and `web/static/js/a11y/keyboard_nav.js` managing trap groups and roving indices.
- **Failure / Guard Rail:** Traps focus within modal dialogs and deliverability rescue sheets, preventing accidental tab escapes into hidden background elements.

### Improvement 365: Spring-Physics Reduced-Motion Motion Tokens & Fluid Fallbacks
- **Description:** Responsive CSS transition and animation framework respecting the OS-level `prefers-reduced-motion: reduce` media query. In standard mode, executes smooth spring-physics cubic beziers (`cubic-bezier(0.34, 1.56, 0.64, 1)`); in reduced-motion mode, swaps to instant discrete fades (`duration: 0.05s`).
- **UI/Terminal Component:** Motion sensitivity badge in settings modal and automatic system media query listener.
- **Telemetry & State Fields:** `isReducedMotionPreferred` (bool), `hasSpringPhysics` (bool), `isMotionDisabled` (bool), `activeEasingFunction` (string), `transitionDurationMs` (int), `springStiffness` (float64).
- **Backend / Engine:** `pkg/web/tokens/motion_engine.go` compiling dynamic motion tokens into root CSS stylesheet.
- **Failure / Guard Rail:** If client system reports low GPU performance or user toggles reduced motion, automatically suppresses all CSS canvas particle effects and heavy blur filters.

### Improvement 366: Scalable Mathematical Typographic Hierarchy & Clamp Token System
- **Description:** Golden-ratio fluid typography system leveraging CSS `clamp(min, preferred, max)` functions across 8 hierarchical levels (`--font-size-xs` to `--font-size-4xl`). Adapts smoothly from mobile displays to ultra-wide 4K operations consoles without viewport jump artifacts.
- **UI/Terminal Component:** Dynamic fluid text scaling across warm-up throughput counters, reputation percentages, and temporal scrubber timelines.
- **Telemetry & State Fields:** `isFluidScalingActive` (bool), `hasClampSupport` (bool), `isWideViewport` (bool), `currentViewportWidth` (int), `baseRemSize` (float64), `activeTypographicScale` (string).
- **Backend / Engine:** `pkg/web/tokens/typography.go` defining CSS calculations and mathematical font scale curves.
- **Failure / Guard Rail:** Clamps font rendering within safe readable boundaries (`min: 12px`, `max: 36px` for headers) preventing microscopic or overflowing text on extreme display scales.

### Improvement 367: Glassmorphic Backdrop Blur Token Tiering & GPU Acceleration Guard
- **Description:** Multi-tier backdrop blur token system (`--glass-blur-sm: 4px`, `--glass-blur-md: 12px`, `--glass-blur-lg: 24px`) coupled with hardware-accelerated CSS `transform: translateZ(0)` and `will-change: backdrop-filter`. Automatically reduces blur radius on low-tier graphics chipsets.
- **UI/Terminal Component:** Translucent glassmorphic control panels, seed inbox cards, and telemetry overlays displaying depth layers.
- **Telemetry & State Fields:** `isGpuAccelerated` (bool), `hasHardwareBackdrop` (bool), `isLowPowerDevice` (bool), `activeBlurRadiusPx` (int), `frameRateFps` (float64), `surfaceLayerDepth` (int).
- **Backend / Engine:** `pkg/web/tokens/glass_surface.go` evaluating client WebGL renderer strings and adjusting CSS tokens accordingly.
- **Failure / Guard Rail:** If animation frame rate drops below 45 FPS during scrolling or scrubbing, gracefully falls back to solid dark surfaces without backdrop blur.

### Improvement 368: Procedural Web Audio API Deliverability Sonification Engine
- **Description:** Zero-audio-asset procedural synthesizer using the HTML5 Web Audio API to produce discreet micro-frequency tones for critical warm-up lifecycle events (e.g. gentle harmonic chord on successful deliverability ramp, low resonant hum on ISP rate-limit, subtle chime on seed inbox reply).
- **UI/Terminal Component:** Audio feedback toggle in top-bar with volume slider (0–100%) and mute indicator. CLI equivalent via `gitmap warmup audio --enable=true`.
- **Telemetry & State Fields:** `isAudioSynthEnabled` (bool), `hasAudioContext` (bool), `isMuted` (bool), `currentVolumeLevel` (float64), `oscillatorFrequencyHz` (float64), `lastPlayedEventTone` (string).
- **Backend / Engine:** `web/static/js/audio/audio_synthesizer.js` constructing dynamic `OscillatorNode` and `GainNode` envelopes in pure JavaScript with zero external WAV/MP3 files.
- **Failure / Guard Rail:** Strictly gates audio playback behind an explicit user gesture (click/toggle) complying with browser auto-play policies; silences gracefully if AudioContext fails.

### Improvement 369: Terminal Lipgloss Design Token Synchronizer
- **Description:** Automated token bridge mapping CSS custom properties directly to Lipgloss terminal ANSI 256-color and TrueColor tokens. Guarantees visual identity parity between the browser web console and the command-line terminal TUI.
- **UI/Terminal Component:** Unified color and border rendering across `gitmap warmup status`, `gitmap warmup studio`, and web dashboards.
- **Telemetry & State Fields:** `isTrueColorSupported` (bool), `hasTokenParity` (bool), `isColor256Fallback` (bool), `terminalColorProfile` (string), `activeTerminalTheme` (string), `colorMapEntryCount` (int).
- **Backend / Engine:** `pkg/cli/tui/token_bridge.go` parsing web token definitions and generating Lipgloss style definitions for CLI renderers.
- **Failure / Guard Rail:** Automatically detects terminals lacking 24-bit TrueColor support and downsamples hex tokens to closest ANSI 256 or 16-color palette matches.

### Improvement 370: Zero-CDN Self-Contained Web Asset Bundler & Font Subset Embedder
- **Description:** Pure Go static asset embedder compiling all CSS tokens, SVG icons, vanilla JavaScript components, and WOFF2 variable font subsets directly into the single executable binary via Go `embed.FS`. Ensures 100% offline, air-gapped operational independence with zero external CDN requests.
- **UI/Terminal Component:** Network status indicator in footer showing `Air-Gapped: 100% Offline (Zero CDN)`.
- **Telemetry & State Fields:** `isAirGappedReady` (bool), `hasZeroCdnDependencies` (bool), `isEmbeddedFsLoaded` (bool), `embeddedAssetByteSize` (int64), `assetHashSha256` (string), `embeddedFileCount` (int).
- **Backend / Engine:** `pkg/web/assets/asset_embedder.go` serving HTTP assets with strict Content-Security-Policy headers banning remote CDN domain loading.
- **Failure / Guard Rail:** Continuous integration audit verifies that zero external URLs or network font links exist across any web template.

---

## 3. Cluster 38: Split-DB SQLite Sharding & High-Speed Query Cache (Items 371–380)

### Improvement 371: High-Throughput Split-DB WAL Pragmas & Concurrency Tuning
- **Description:** Database engine tuning module configuring dedicated Split-DB instances (`warmup.db` for operational state and `mail_audit.db` for immutable event history) with high-performance SQLite pragmas: `PRAGMA journal_mode=WAL`, `PRAGMA synchronous=NORMAL`, `PRAGMA busy_timeout=5000`, and `PRAGMA cache_size=-64000` (64MB memory cache).
- **UI/Terminal Component:** Database status widget in Admin Health view and CLI diagnostic output in `gitmap db status`.
- **Telemetry & State Fields:** `isWalModeActive` (bool), `hasSynchronousNormal` (bool), `isBusyTimeoutConfigured` (bool), `cacheSizeBytes` (int64), `walCheckpointCount` (int), `busyTimeoutMs` (int).
- **Backend / Engine:** `pkg/db/splitdb/pragmas.go` applying idempotent pragma configurations upon SQLite connection pool initialization.
- **Failure / Guard Rail:** If WAL pragma fails due to read-only filesystem or SQLite lock, wraps error into typed `*apperror.AppError` and falls back to memory mode with safe alerting.

### Improvement 372: Monthly Historical Audit Sharding & Table Partitioning
- **Description:** Time-series database partitioner that automatically shards `mail_audit.db` into monthly segmented tables (`mail_events_2026_10`, `mail_events_2026_11`) or separate physical shard files (`mail_audit_2026_10.db`). Maintains constant `O(1)` query speed regardless of multi-million row scale.
- **UI/Terminal Component:** Shard overview table showing monthly partition name, row count, file size, and query latency. CLI command `gitmap db shard list`.
- **Telemetry & State Fields:** `isPartitioned` (bool), `hasMonthlyShard` (bool), `isShardAttached` (bool), `activeShardMonth` (string), `totalPartitionCount` (int), `activeShardRowCount` (int64).
- **Backend / Engine:** `pkg/db/splitdb/shard_partitioner.go` executing automatic month-boundary table creation and unified view union query dispatch.
- **Failure / Guard Rail:** Employs SQLite `ATTACH DATABASE` dynamically; if a shard is missing or corrupt, marks that month as unavailable and completes search across remaining shards.

### Improvement 373: LRU High-Speed Memory Query Cache & Invalidation Registry
- **Description:** Thread-safe in-memory Least-Recently-Used (LRU) query cache with configurable TTL (default 15 seconds) for frequent deliverability rollups, daily volume sums, and seed engagement scores. Eliminates redundant disk reads during rapid UI refresh intervals.
- **UI/Terminal Component:** Cache hit/miss ratio gauge on diagnostic dashboard and CLI metric output.
- **Telemetry & State Fields:** `isCacheHit` (bool), `hasCapacityAvailable` (bool), `isEntryExpired` (bool), `cacheHitRatio` (float64), `cachedQueryCount` (int), `evictionCount` (int64).
- **Backend / Engine:** `pkg/db/cache/query_cache.go` wrapping database repository queries with mutex-guarded LRU bucket and automatic invalidation on write events.
- **Failure / Guard Rail:** Automatically evicts oldest entries when memory cap (default 32MB) is reached; thread-safe sync.RWMutex prevents race conditions.

### Improvement 374: In-Memory Bloom Filter for Instant Recipient Suppression Checks
- **Description:** High-efficiency probabilistic Bloom filter (space-optimized bit array with 4 independent hash functions) held in RAM for sub-microsecond suppression checks. Instantly intercepts bounced or unsubscribed recipient addresses before initiating database queries or SMTP handshakes.
- **UI/Terminal Component:** Suppression engine status card showing Bloom filter memory footprint, item count, and false positive probability.
- **Telemetry & State Fields:** `isSuppressed` (bool), `hasBloomMatch` (bool), `isFalsePositiveChecked` (bool), `bloomFilterBits` (int), `registeredItemCount` (int), `falsePositiveRate` (float64).
- **Backend / Engine:** `pkg/mail/suppression/bloom_filter.go` maintaining bit arrays and reconciling with database on cold boot.
- **Failure / Guard Rail:** On positive Bloom filter match, performs a secondary authoritative lookup in SQLite `warmup.db` suppression table to definitively rule out false positives.

### Improvement 375: Automated Background WAL Checkpoint Governor
- **Description:** Non-blocking background governor monitoring the size of the SQLite write-ahead log (`.db-wal`). When WAL file exceeds 4MB or 1,000 uncommitted frames, executes `PRAGMA wal_checkpoint(PASSIVE)` during idle sending windows, preventing uncontrolled disk bloat.
- **UI/Terminal Component:** WAL size indicator with color-coded warning thresholds (Green <4MB, Amber 4-10MB, Red >10MB).
- **Telemetry & State Fields:** `isWalGovernorRunning` (bool), `hasPendingFrames` (bool), `isCheckpointNeeded` (bool), `walFileSize` (int64), `uncommittedFrameCount` (int), `lastCheckpointDurationMs` (int64).
- **Backend / Engine:** `pkg/db/splitdb/wal_governor.go` running 60-second ticker loop with gentle passive checkpoint sweeps.
- **Failure / Guard Rail:** Never executes blocking `TRUNCATE` checkpoints while active SMTP deliveries are writing; backs off automatically if database reports `SQLITE_BUSY`.

### Improvement 376: SHA-256 Merkle Tree Ledger for Tamper-Evident Delivery Audit
- **Description:** Cryptographic verification layer constructing an immutable Merkle tree of every outbound delivery event, bounce notification, and warm-up engagement record. Calculates periodic root hashes stored in SQLite with block chaining, guaranteeing complete audit integrity against manual database tampering.
- **UI/Terminal Component:** Cryptographic audit shield badge showing current Merkle root hash and verification timestamp. CLI audit verify command via `gitmap audit verify-integrity`.
- **Telemetry & State Fields:** `isMerkleValid` (bool), `hasTamperDetected` (bool), `isLedgerChained` (bool), `currentRootHash` (string), `leafEventCount` (int64), `verifiedBlockHeight` (int).
- **Backend / Engine:** `pkg/db/audit/merkle_tree.go` computing iterative SHA-256 hashes of serialized event records and persisting block headers.
- **Failure / Guard Rail:** If any historical row is modified or deleted outside the application layer, hash verification mismatch immediately trips a tamper alert and freezes deliverability dispatch.

### Improvement 377: Vacuum Incremental Page Compactor & Fragmentation Monitor
- **Description:** Storage maintenance routine enabling SQLite `PRAGMA auto_vacuum = INCREMENTAL`. Regularly measures B-tree page fragmentation and executes `PRAGMA incremental_vacuum(100)` during scheduled off-peak maintenance hours to reclaim unused disk pages.
- **UI/Terminal Component:** Storage fragmentation gauge showing wasted page percentage and last compactor run time.
- **Telemetry & State Fields:** `isCompactorActive` (bool), `hasHighFragmentation` (bool), `isVacuumScheduled` (bool), `freePageCount` (int), `reclaimedSpaceBytes` (int64), `fragmentationPercent` (float64).
- **Backend / Engine:** `pkg/db/splitdb/vacuum_compactor.go` calculating page statistics via SQLite `dbstat` virtual table and executing incremental page release.
- **Failure / Guard Rail:** Bounded to reclaiming a maximum of 1,000 pages per cycle to prevent prolonged I/O latency spikes on shared disks.

### Improvement 378: Prepared Statement Cache & Connection Pool Re-Use Engine
- **Description:** Enterprise SQLite connection wrapper maintaining a pool of pre-compiled `sqlite3_stmt` prepared statements across all active worker goroutines. Eliminates SQL parsing and query planning overhead for high-frequency insert and update queries.
- **UI/Terminal Component:** Statement cache metrics card showing active prepared statements, hit count, and query execution timing.
- **Telemetry & State Fields:** `isStatementCached` (bool), `hasPooledConnection` (bool), `isQueryPrepared` (bool), `cachedStatementCount` (int), `executionDurationMicros` (int64), `cacheHitCount` (int64).
- **Backend / Engine:** `pkg/db/splitdb/statement_pool.go` caching compiled queries keyed by parameterized SQL string hash.
- **Failure / Guard Rail:** Automatically clears statement handles and re-prepares them if schema migration or database re-open invalidates statement pointers.

### Improvement 379: Deadlock-Free SQLite Connection Multiplexer & Mutex Gate
- **Description:** Concurrency governor enforcing a strict single-writer, multi-reader architecture for each SQLite database file. Channels all write transactions through a dedicated FIFO channel worker, completely eliminating `SQLITE_BUSY` and database lock contention under multi-agent swarms.
- **UI/Terminal Component:** Concurrency queue depth bar showing pending write transactions and write queue latency.
- **Telemetry & State Fields:** `isWriteLocked` (bool), `hasPendingWrites` (bool), `isMultiplexerHealthy` (bool), `writeQueueDepth` (int), `averageWriteLatencyMs` (float64), `activeReaderCount` (int).
- **Backend / Engine:** `pkg/db/splitdb/write_multiplexer.go` running a dedicated background write loop with bounded buffered channels.
- **Failure / Guard Rail:** If write queue exceeds 500 items, applies gentle backpressure to caller goroutines rather than dropping transactions or deadlocking.

### Improvement 380: Automated Split-DB Disaster Backup & Hot Snapshot Utility
- **Description:** Online database backup engine leveraging SQLite Online Backup API (`sqlite3_backup_init`) to stream hot, point-in-time snapshot copies of `warmup.db` and `mail_audit.db` to a compressed backup directory (`.gitmap/backups/`) without stopping active sending operations.
- **UI/Terminal Component:** Backup status card showing last snapshot timestamp, backup size, integrity status, and one-click manual backup trigger. CLI command `gitmap db backup create`.
- **Telemetry & State Fields:** `isBackupSuccessful` (bool), `hasActiveBackup` (bool), `isSnapshotConsistent` (bool), `backupFilePath` (string), `backupDurationMs` (int64), `backupSizeBytes` (int64).
- **Backend / Engine:** `pkg/db/splitdb/backup_engine.go` coordinating step-by-step page copying with progress callbacks and automatic zstd compression.
- **Failure / Guard Rail:** Verifies backup file integrity using SQLite `PRAGMA quick_check` immediately after completion; discards invalid backups and triggers notification.

---

## 4. Cluster 39: Antigravity IDE Event Bus & Dynamic Workspace Hooks (Items 381–390)

### Improvement 381: Bidirectional Pub/Sub WebSocket Event Broker
- **Description:** Real-time event broker providing bidirectional pub/sub communication between the GitMap backend daemon, browser Web Components, and Antigravity IDE extensions. Broadcasts deliverability state changes, warming progress, and agent telemetry with sub-millisecond dispatch.
- **UI/Terminal Component:** Connection pulse badge in top bar showing WebSocket status (`Connected`, `Reconnecting`, `Offline`) and live event throughput counter.
- **Telemetry & State Fields:** `isWebSocketConnected` (bool), `hasActiveSubscribers` (bool), `isReconnecting` (bool), `connectedClientCount` (int), `messagesSentPerSec` (float64), `lastEventTimestamp` (string).
- **Backend / Engine:** `pkg/agent/bus/event_broker.go` managing client connection pools, subscription channels, and ping/pong keepalive intervals.
- **Failure / Guard Rail:** Automatically buffers up to 200 events in memory during client disconnects and replays them upon successful reconnection.

### Improvement 382: Antigravity IDE Workspace Activity Hook & Throttling Sensor
- **Description:** Intelligent workspace sensor detecting when the developer is actively editing code, debugging, or running test suites in the Antigravity IDE. Automatically dials back warm-up outbound volume and background CPU usage by 50% to prevent local machine slowdown.
- **UI/Terminal Component:** Work mode indicator in footer showing `Operator Active: Background Throttling 50%` or `Operator Idle: Full Ramp Mode`.
- **Telemetry & State Fields:** `isWorkspaceActive` (bool), `hasActiveIdeSession` (bool), `isThrottlingApplied` (bool), `ideActivityLevel` (string), `cpuUsageReductionPercent` (float64), `idleSeconds` (int).
- **Backend / Engine:** `pkg/agent/hooks/ide_sensor.go` listening for IDE activity events via local domain socket or file modification watchers.
- **Failure / Guard Rail:** If IDE activity ceases for 5 minutes, smoothly ramps background worker threads back to full operational velocity.

### Improvement 383: Slash-Command CLI Trigger Dispatcher (`/warmup`)
- **Description:** Universal command dispatcher enabling developers and agents to execute warm-up actions and deliverability checks directly via chat slash-commands (e.g. `/warmup status`, `/warmup pause`, `/warmup rescue --isp=google`, `/warmup report`).
- **UI/Terminal Component:** Slash-command execution panel rendering structured JSON and Markdown responses directly in chat context.
- **Telemetry & State Fields:** `isCommandRecognized` (bool), `hasValidArguments` (bool), `isExecutionSuccess` (bool), `commandTrigger` (string), `invokedByAgentId` (string), `executionTimeMs` (int64).
- **Backend / Engine:** `pkg/cli/cmdwarmup/slash_dispatcher.go` parsing command tokens and delegating to root Cobra command handlers.
- **Failure / Guard Rail:** Rejects unauthorized or malformed slash-commands with descriptive syntax help without executing underlying actions.

### Improvement 384: Statistical Deliverability Anomaly Triage & Real-Time Alert Engine
- **Description:** Real-time statistical analysis engine computing Z-scores and moving variance across deliverability metrics (bounce spikes, inbox placement drops, latency increases). When metric deviates beyond 3 standard deviations, generates priority triage alerts.
- **UI/Terminal Component:** Red floating triage alert banner with dismiss action, diagnostic breakdown, and recommended remediation button.
- **Telemetry & State Fields:** `isAnomalyDetected` (bool), `hasSevereZScore` (bool), `isAlertActive` (bool), `anomalyMetricName` (string), `currentZScore` (float64), `confidenceThreshold` (float64).
- **Backend / Engine:** `pkg/agent/anomaly/triage_engine.go` maintaining rolling windows of 100 delivery samples and evaluating deviation limits.
- **Failure / Guard Rail:** Deduplicates alerts with a 15-minute cooldown to prevent alert fatigue during sustained incident conditions.

### Improvement 385: Signed Webhook Dispatcher for External Autonomous Agents
- **Description:** Outbound webhook notification service delivering cryptographically signed HMAC-SHA256 payloads to external autonomous agents, alerting them to deliverability milestone completions, spam trap incidents, or required human review.
- **UI/Terminal Component:** Webhook configuration card showing endpoint URL, secret key rotation button, and historical delivery delivery logs.
- **Telemetry & State Fields:** `isWebhookEnabled` (bool), `hasValidSignature` (bool), `isDelivered` (bool), `targetEndpointUrl` (string), `httpStatusCode` (int), `retryAttemptCount` (int).
- **Backend / Engine:** `pkg/agent/bus/webhook_dispatcher.go` generating `X-GitMap-Signature` headers and managing exponential backoff retries.
- **Failure / Guard Rail:** Automatically disables webhook endpoint if 5 consecutive retry attempts fail with 5xx or connection timeout errors.

### Improvement 386: Terminal Bell Audio/Visual Event Router
- **Description:** Cross-platform notification router sending ASCII visual flash bell (`\a`) or desktop notifications for critical warm-up milestones (e.g. 100% daily warming target achieved, fatal SPF/DKIM authentication failure).
- **UI/Terminal Component:** Visual screen flash effect in terminal TUI and desktop OS notification toast.
- **Telemetry & State Fields:** `isBellEnabled` (bool), `hasVisualFlash` (bool), `isDesktopNotificationSent` (bool), `notificationTitle` (string), `severityLevel` (string), `bellTriggerEvent` (string).
- **Backend / Engine:** `pkg/cli/tui/bell_router.go` checking terminal capabilities and dispatching appropriate OS notification hooks.
- **Failure / Guard Rail:** Respects user's "Do Not Disturb" system settings and CLI `--quiet` flags; never sounds audible bell in non-interactive CI environments.

### Improvement 387: Headless Background Daemon Supervisor & Process Health Check
- **Description:** Process supervisor managing the GitMap Warm-Up background worker as a detached headless daemon. Handles graceful shutdown on `SIGINT`/`SIGTERM`, PID file locking, automatic crash restart, and heartbeat health checks.
- **UI/Terminal Component:** Daemon status widget showing PID, uptime, memory consumption, and restart count. CLI commands `gitmap warmup daemon [start|stop|status]`.
- **Telemetry & State Fields:** `isDaemonRunning` (bool), `hasValidPidLock` (bool), `isHealthy` (bool), `processId` (int), `uptimeSeconds` (int64), `memoryUsageBytes` (int64).
- **Backend / Engine:** `pkg/daemon/supervisor.go` acquiring exclusive file locks on `.gitmap/warmup.pid` and orchestrating clean component termination.
- **Failure / Guard Rail:** Detects stale PID files from dirty system reboots and safely clears locks after verifying process existence.

### Improvement 388: Structured JSON-Lines Audit Log Streamer
- **Description:** Real-time log streaming pipeline outputting structured, JSON-Lines formatted event logs to `stdout` or rotating log files (`.gitmap/logs/warmup.log`). Enables seamless ingestion into Datadog, Grafana Loki, or local developer grep pipelines.
- **UI/Terminal Component:** Interactive streaming log viewer tab with live search filter, log level selector (`DEBUG`, `INFO`, `WARN`, `ERROR`), and pause button.
- **Telemetry & State Fields:** `isLogStreamingActive` (bool), `hasLogFileRotated` (bool), `isFormattedJson` (bool), `currentLogLevel` (string), `logEntryCount` (int64), `activeLogFilePath` (string).
- **Backend / Engine:** `pkg/logger/json_streamer.go` writing buffered JSON structures with timestamp, subsystem, level, and contextual message attributes.
- **Failure / Guard Rail:** Implements size-based log rotation (max 10MB per file, retaining 5 compressed backups) to prevent disk exhaustion.

### Improvement 389: Live Metric Exporter for Prometheus & OpenTelemetry
- **Description:** Embedded HTTP metrics endpoint (`/metrics`) exposing Prometheus-formatted time-series gauges and counters for outbound SMTP rates, deliverability percentages, bounce counts, and queue depth.
- **UI/Terminal Component:** Prometheus integration toggle in settings with copyable endpoint URL (`http://127.0.0.1:9095/metrics`).
- **Telemetry & State Fields:** `isMetricsServerActive` (bool), `hasOpenTelemetryEnabled` (bool), `isScraped` (bool), `metricsPort` (int), `scrapeCount` (int64), `lastScrapeTimestamp` (string).
- **Backend / Engine:** `pkg/telemetry/prometheus_exporter.go` registering and updating standard Prometheus collector gauges without external dependencies.
- **Failure / Guard Rail:** Binds exclusively to `127.0.0.1` by default to prevent unauthorized external network exposure; supports token authentication.

### Improvement 390: Agent Multi-Session State Synchronizer
- **Description:** State synchronization coordinator reconciling concurrent changes made across multiple open browser tabs, terminal TUI sessions, and autonomous agent processes. Utilizes optimistic concurrency control with monotonic version vectors.
- **UI/Terminal Component:** Sync conflict resolver modal prompting operator only when irreconcilable simultaneous edits occur.
- **Telemetry & State Fields:** `isStateSynchronized` (bool), `hasVersionConflict` (bool), `isOptimisticLockValid` (bool), `currentVersionVector` (int64), `activeSessionCount` (int), `lastSyncDurationMs` (int64).
- **Backend / Engine:** `pkg/agent/bus/session_sync.go` checking revision numbers before committing state changes to SQLite `warmup.db`.
- **Failure / Guard Rail:** On version collision, cleanly reloads latest database state and attempts automated three-way merge before surfacing UI conflict.

---

## 5. Cluster 40: Autonomous Deliverability Scorecards & Historical Audit (Items 391–400)

### Improvement 391: Autonomous Multi-Factor Deliverability Health Scorecard
- **Description:** Algorithmic scorecard aggregator evaluating 10 weighted deliverability dimensions (SPF/DKIM/DMARC health, bounce rate, spam complaint ratio, warming ramp fidelity, seed inbox engagement, ISP queue friction, IP reputation, domain age, content entropy, and DNS consistency). Produces a unified 0–100 Deliverability Health Score.
- **UI/Terminal Component:** Circular SVG gauge in main dashboard with color-coded score tiers (90-100 Elite Green, 75-89 Good Amber, <75 Critical Red) and breakdown radar chart.
- **Telemetry & State Fields:** `hasPassedAudit` (bool), `isEliteDeliverability` (bool), `hasCriticalFailure` (bool), `compositeScore` (float64), `dimensionScores` (map[string]float64), `auditTimestamp` (string).
- **Backend / Engine:** `pkg/mail/audit/scorecard_engine.go` aggregating metric weights and calculating normalized scorecards.
- **Failure / Guard Rail:** If any single critical dimension (e.g. DMARC reject or bounce rate >5%) fails catastrophically, caps overall score at maximum 60.

### Improvement 392: ISP-Specific Compliance Radar & Requirement Matrix
- **Description:** Specialized compliance verification matrix checking outbound practices against published ISP policies (Google Sender Guidelines 2024, Yahoo Postmaster Guidelines, Microsoft SNDS requirements). Detects missing one-click unsubscribe headers, TLS 1.3 requirements, and spam complaint thresholds.
- **UI/Terminal Component:** Interactive radar chart comparing compliance percentages across Google, Microsoft, Yahoo, and Apple Mail.
- **Telemetry & State Fields:** `isGoogleCompliant` (bool), `isYahooCompliant` (bool), `isMicrosoftCompliant` (bool), `googleComplianceScore` (float64), `yahooComplianceScore` (float64), `nonCompliantRuleCount` (int).
- **Backend / Engine:** `pkg/mail/audit/compliance_matrix.go` evaluating message headers and DNS records against codified ISP rule sets.
- **Failure / Guard Rail:** Highlights specific missing headers (e.g. `List-Unsubscribe-Post: List-Unsubscribe=One-Click`) with actionable one-click fix buttons.

### Improvement 393: Historical Deliverability Trend & Regression Analyzer
- **Description:** Longitudinal analytics engine comparing deliverability scores, inbox placement ratios, and bounce statistics over 7-day, 30-day, and 90-day rolling windows. Automatically detects statistical trends and alerts to subtle reputation regressions before major blacklisting occurs.
- **UI/Terminal Component:** Multi-series interactive SVG trend chart with toggleable metric layers and baseline regression markers.
- **Telemetry & State Fields:** `isRegressing` (bool), `hasPositiveTrend` (bool), `isTrendStatisticallySignificant` (bool), `trendSlope` (float64), `changePercentage30d` (float64), `historicalWindowDays` (int).
- **Backend / Engine:** `pkg/mail/audit/trend_analyzer.go` calculating linear regression slopes and Pearson correlation coefficients over time-series shards.
- **Failure / Guard Rail:** Ignores minor statistical noise by requiring at least 5 days of consistent downward slope before triggering a regression alert.

### Improvement 394: Seed Inbox Response Matrix & Engagement Audit Inspector
- **Description:** Deep inspection table auditing the status of all seed inboxes across the multi-agent network. Displays arrival folder (`Inbox`, `Spam`, `Promotions`, `Missing`), open state, automated reply timestamp, and thread persistence.
- **UI/Terminal Component:** Filterable data table with search, ISP badge filters, folder placement indicators, and raw MIME message viewer drawer.
- **Telemetry & State Fields:** `isArrivedInInbox` (bool), `hasReplied` (bool), `isThreaded` (bool), `seedAddress` (string), `ispProvider` (string), `arrivalFolder` (string), `responseTimeSeconds` (int).
- **Backend / Engine:** `pkg/mail/seed/audit_inspector.go` querying IMAP/POP3 seed mailboxes and summarizing folder placement statistics.
- **Failure / Guard Rail:** Flags unresponsive seed inboxes that fail to authenticate or sync for over 2 hours, automatically replacing them with reserve seeds.

### Improvement 395: Automated Remediation Playbook Dispatcher
- **Description:** Autonomous playbook execution engine triggered when deliverability score drops below acceptable thresholds. Offers one-click or fully automated remediation actions (e.g. "Trigger 48h Volume Throttling", "Exclude Microsoft IP Pool", "Activate Seed Reply Boost", "Regenerate Content Entropy").
- **UI/Terminal Component:** Remediation action drawer displaying prioritized playbooks, expected recovery time, and "Execute Plan" button.
- **Telemetry & State Fields:** `isPlaybookActive` (bool), `hasAutoRemediated` (bool), `isActionApproved` (bool), `playbookName` (string), `recommendedActionId` (string), `estimatedRecoveryHours` (int).
- **Backend / Engine:** `pkg/mail/remediation/playbook_dispatcher.go` orchestrating system parameter adjustments and scheduling recovery verification.
- **Failure / Guard Rail:** Requires explicit operator confirmation if a remediation playbook involves permanent recipient suppression or major sending halts.

### Improvement 396: SPF/DKIM/DMARC Continuous DNS Health Verifier
- **Description:** Background diagnostic loop verifying live DNS records for all sending domains every 6 hours. Detects DNS propagation issues, accidental TXT record deletions, syntax errors, and misaligned SPF `include:` mechanisms.
- **UI/Terminal Component:** DNS health status grid with green checkmarks for SPF, DKIM, and DMARC, and diagnostic record viewer. CLI command `gitmap warmup dns verify`.
- **Telemetry & State Fields:** `isSpfValid` (bool), `isDkimValid` (bool), `isDmarcValid` (bool), `sendingDomain` (string), `dnsLookupLatencyMs` (int64), `dnsRecordError` (string).
- **Backend / Engine:** `pkg/mail/dns/health_verifier.go` performing direct DNS queries against Google (8.8.8.8) and Cloudflare (1.1.1.1) resolvers to bypass local cache poisoning.
- **Failure / Guard Rail:** Instantly pauses active warming schedule if SPF or DKIM DNS records fail to resolve, preventing massive ISP reputation destruction.

### Improvement 397: DMARC Aggregate Report (RUA) Ingestion & Forensic Parser
- **Description:** Forensic XML parser capable of ingesting and analyzing DMARC aggregate XML reports (RFC 7489) received from Google, Yahoo, and Microsoft postmasters. Identifies legitimate sending IP pass rates and pinpoints unauthorized spoofing attempts.
- **UI/Terminal Component:** Forensic report dashboard showing passing vs failing IP addresses, geographic origin map, and spoofing alert cards.
- **Telemetry & State Fields:** `isReportParsed` (bool), `hasSpoofingDetected` (bool), `isAlignmentPassed` (bool), `reportOrgName` (string), `totalReportedMessages` (int64), `passingPercentage` (float64).
- **Backend / Engine:** `pkg/mail/dmarc/rua_parser.go` decompressing `.gz` or `.zip` XML attachments and storing structured statistics in `mail_audit.db`.
- **Failure / Guard Rail:** Safely discards malformed or oversized XML payloads with a structured `*apperror.AppError` and moves to next report.

### Improvement 398: Tamper-Evident Deliverability Audit Certificate Generator
- **Description:** Report generator synthesizing complete deliverability scorecards, Merkle root verification proofs, and ISP compliance matrices into an authoritative, digitally signed HTML/PDF audit certificate for executive review and compliance archiving.
- **UI/Terminal Component:** "Export Compliance Certificate" button in header, producing a clean, printable vector report with embedded QR code verification link.
- **Telemetry & State Fields:** `isCertificateGenerated` (bool), `hasCryptographicSignature` (bool), `isVerificationValid` (bool), `certificateId` (string), `signerIdentity` (string), `reportGenerationTimestamp` (string).
- **Backend / Engine:** `pkg/mail/audit/certificate_generator.go` compiling self-contained HTML templates with inline SVG graphics and cryptographic signatures.
- **Failure / Guard Rail:** Verifies current Merkle tree root hash before generating certificate; aborts with error if database integrity cannot be validated.

### Improvement 399: Deliverability Incident Post-Mortem & Timeline Reconstructor
- **Description:** Incident reconstruction tool that automatically compiles a step-by-step chronology of events surrounding deliverability dips (e.g. bounce spike at 14:02, ISP throttle response at 14:05, circuit breaker trip at 14:07, rescue daemon activation at 14:10).
- **UI/Terminal Component:** Chronological incident timeline card with expandable event nodes, raw SMTP transcript snippets, and root-cause analysis tags.
- **Telemetry & State Fields:** `isPostMortemAvailable` (bool), `hasRootCauseIdentified` (bool), `isResolved` (bool), `incidentId` (string), `incidentDurationMinutes` (int), `rootCauseCategory` (string).
- **Backend / Engine:** `pkg/mail/audit/post_mortem.go` querying historical audit shards across the incident time window and correlating event patterns.
- **Failure / Guard Rail:** Automatically archives incident chronologies in SQLite for retrospective training of autonomous deliverability agents.

### Improvement 400: Autonomous Deliverability Self-Healing Master Coordinator
- **Description:** The pinnacle coordinator integrating all Sequence 340 subsystems: monitors swarm agents, sending rhythms, P2P mesh seeds, Split-DB health, and real-time scorecards. When systemic deliverability risks arise, orchestrates holistic self-healing adjustments across the entire platform without manual operator action.
- **UI/Terminal Component:** Global Autonomous Self-Healing master card showing real-time coordinator status (`Vigilant`, `Optimizing`, `Self-Healing`, `Stable`), active mitigation workflows, and historical healing score improvements.
- **Telemetry & State Fields:** `isSelfHealingActive` (bool), `hasAutonomousMitigation` (bool), `isSystemStable` (bool), `activeMitigationWorkflow` (string), `healingInterventionCount` (int), `reputationGainDelta` (float64).
- **Backend / Engine:** `pkg/mail/coordinator/self_healing_master.go` binding agent swarm consensus, circuit breakers, warming curves, and database governors into a unified closed-loop feedback controller.
- **Failure / Guard Rail:** Maintains a safety lock threshold: if autonomous interventions fail to stabilize deliverability score within 60 minutes, gracefully transitions to failsafe pause mode and summons human operator via high-priority notification.

---

## 6. Coding Guidelines & Architecture Verification Checklist

All code authored in support of Improvements 361 through 400 must strictly conform to:
- [x] **Error Return Sovereignty**: Functions returning `error` or `(*Result, error)` MUST return the error directly (`return err` or `return apperror.Wrap(...)`). Never call `cliexit.HandleError` or `os.Exit` inside leaf functions.
- [x] **Function Line Limit**: Every function <= 15 lines max (<= 8 lines preferred). Decompose large operations into focused helpers.
- [x] **Zero Nested Ifs**: Flatten all conditionals with guard clauses and early returns.
- [x] **Positive Booleans**: All boolean variables, fields, and parameters MUST begin with `is` or `has` (`isHealthy`, `hasPassedAudit`, `isThemeActive`). Zero negative booleans (`isNotValid` forbidden).
- [x] **Zero External CDN**: Pure vanilla JavaScript, CSS custom properties, and SVG icons embedded directly into Go binary via `embed.FS`.
- [x] **Split-DB Isolation**: Strict separation between operational state (`warmup.db`) and historical audit ledger (`mail_audit.db`).
- [x] **Relative Path Hygiene**: All documentation, code references, and configuration paths must strictly use relative repository paths.
