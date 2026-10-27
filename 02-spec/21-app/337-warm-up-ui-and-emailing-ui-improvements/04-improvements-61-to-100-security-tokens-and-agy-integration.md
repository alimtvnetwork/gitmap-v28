# Specification 337: Improvements 61 to 100 — Security Diagnostics, Design System Tokens, DB Telemetry & AGY Bridge

- **Sequence:** 337
- **Slug:** `337-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/04-improvements-61-to-100-security-tokens-and-agy-integration.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 3

---

## 1. Executive Summary & Cluster Scope

This specification defines Improvements 61 through 100 for the GitMap Warm-Up and Emailing UI Engine. These 40 improvements span four critical architecture clusters:

1. **Cluster 7: Authentication, SPF/DKIM/DMARC Diagnostics & Spam Score HUD (Improvements 61–70):** Real-time DNS record verifiers, DKIM selector visual checkers, alignment validators, SpamAssassin/Rspamd score analyzers, blackhole simulation, and RFC compliance inspectors.
2. **Cluster 8: Design System Tokens, Dark/Light Themes, Typography & WCAG AA Contrast (Improvements 71–80):** High-contrast color scales, TokyoNight dark mode, Arctic light mode, accessible interactive elements, font scaling curves, and zero-runtime CSS token systems.
3. **Cluster 9: SQLite Split-DB Telemetry, Query Speed & Audit Logs (Improvements 81–90):** Microsecond query trackers, WAL mode inspection dials, Casbin RBAC audit logs, vacuum & compaction HUDs, index hit-ratio visualizers, and tamper-evident event streaming.
4. **Cluster 10: Antigravity IDE Integration, Task Completion Email Triggers & Automation Bridge (Improvements 91–100):** AGY workspace bridge hooks, headless CLI dispatch, build/test run notification digests, slash-command trigger hooks, subagent swarm alert webhooks, and crash dump email templates.

---

## 2. Cluster 7: Authentication, SPF/DKIM/DMARC Diagnostics & Spam Score HUD (Items 61–70)

### Improvement 61: SPF Record Matrix & DNS Lookup Depth Visualizer
- **Description:** Real-time visual tree of DNS SPF mechanisms (`v=spf1 include:... ~all`). Calculates and graphs DNS lookup depth against the RFC 7208 10-lookup limit.
- **UI/Terminal Component:** Interactive SVG tree in web UI; hierarchical indented ASCII tree with color-coded depth badge in CLI (`gitmap mail auth spf --domain example.com`).
- **Telemetry & State Fields:** `hasValidSpf` (bool), `isUnderLookupLimit` (bool), `lookupCount` (int), `spfRecordRaw` (string), `spfMechanismCount` (int).
- **Backend / Engine:** `pkg/mail/auth/spf_validator.go` with Go `net.Resolver` recursive lookups; zero third-party dependencies.
- **Failure / Guard Rail:** Emits warning when lookup count exceeds 8; raises fatal deliverability block when count exceeds 10.

### Improvement 62: Multi-Selector DKIM Public/Private Key Alignment Checker
- **Description:** Live validator for active DKIM selectors (e.g. `s1`, `google`, `default`). Tests cryptographic signature generation against public DNS TXT key bit-depth (1024 vs 2048 vs 4096-bit RSA/Ed25519).
- **UI/Terminal Component:** Selector badge pill grid showing key length, algorithm badge, and DNS propagation status across 6 global DNS root servers.
- **Telemetry & State Fields:** `hasValidDkim` (bool), `isAligned` (bool), `keyBitDepth` (int), `selectorTag` (string), `hasStrictAlignment` (bool).
- **Backend / Engine:** `pkg/mail/auth/dkim_validator.go` parsing DNS TXT records, decoding base64 RSA keys, and validating test signature verifications.
- **Failure / Guard Rail:** Warns if key bit depth is below 2048-bit; rejects deprecated SHA-1 signatures.

### Improvement 63: DMARC Policy Inspector & Aggregated Reporting Visualizer
- **Description:** Analyzes domain DMARC policy (`p=none`, `p=quarantine`, `p=reject`), percentage rollout (`pct=100`), and reporting destinations (`rua`, `ruf`).
- **UI/Terminal Component:** Policy enforcement shield widget with animated ripple border; visual timeline tracking daily aggregate failure rates from incoming RUA XML feeds.
- **Telemetry & State Fields:** `hasDmarcPolicy` (bool), `isRejectMode` (bool), `isQuarantineMode` (bool), `hasReportingRua` (bool), `pctEnforcement` (int).
- **Backend / Engine:** `pkg/mail/auth/dmarc_parser.go` and `pkg/mail/auth/rua_processor.go` storing parsed aggregate reports into `mail_audit.db`.
- **Failure / Guard Rail:** Flags warning if `p=none` is used for longer than 30 warm-up days; enforces strict rejection policy for enterprise outbox.

### Improvement 64: BIMI (Brand Indicators for Message Identification) & VMC Certificate Radar
- **Description:** Scans domain DNS for BIMI assertions (`default._bimi.domain`), verifies SVG profile (`SVG-Tiny-PS` compliance), and parses Verified Mark Certificate (VMC) validity.
- **UI/Terminal Component:** Live avatar preview displaying brand mark exactly as rendered in Apple Mail, Yahoo Mail, and Gmail web interfaces.
- **Telemetry & State Fields:** `hasBimiRecord` (bool), `hasValidVmc` (bool), `isSvgTinyCompliant` (bool), `vmcExpiryDate` (string), `logoSvgUrl` (string).
- **Backend / Engine:** `pkg/mail/auth/bimi_scanner.go` with strict XML/SVG sanitizer ensuring no scripts, links, or external entities exist.
- **Failure / Guard Rail:** Blocks display and triggers alert if SVG contains malicious XML vectors or external references.

### Improvement 65: SpamAssassin & Rspamd Heuristic Score Predictor HUD
- **Description:** Offline heuristic engine scoring email content, subject, headers, and MIME boundaries against 200+ canonical spam rules before dispatch.
- **UI/Terminal Component:** Vertical score meter (0 to 10 scale). Scores `< 2.0` display emerald green; `2.0 - 5.0` amber; `> 5.0` crimson alert with rule breakdown pill list.
- **Telemetry & State Fields:** `isSpamSafe` (bool), `hasHighRiskTokens` (bool), `spamScore` (float64), `triggeredRuleCount` (int), `matchedRulesJson` (string).
- **Backend / Engine:** `pkg/mail/spam/heuristic_scorer.go` executing regex matching against subject, plain-text body, and HTML markup.
- **Failure / Guard Rail:** Sending is halted if `spamScore >= 5.0` unless forced via `--override-spam-check`.

### Improvement 66: Blackhole and Spam-Trap Pattern Detection Scanner
- **Description:** Pre-send address hygiene scanner identifying known role accounts (`abuse@`, `postmaster@`, `noreply@`), throwaway domains, and honeypot spam-trap patterns.
- **UI/Terminal Component:** Recipient risk tier table highlighting flagged addresses with quick-action bulk remove and suppress buttons.
- **Telemetry & State Fields:** `hasCleanRecipients` (bool), `hasSpamTrapRisk` (bool), `riskyRecipientCount` (int), `suppressionMatchCount` (int).
- **Backend / Engine:** `pkg/mail/hygiene/trap_detector.go` querying local SQLite suppression cache and disposable domain bloom filters.
- **Failure / Guard Rail:** Automatically sequesters flagged addresses to quarantine table; never initiates SMTP handshake to trap domains.

### Improvement 67: Reverse DNS (PTR) & IP Whitelist Consistency Matrix
- **Description:** Validates forward-confirmed reverse DNS (FCrDNS) mapping: ensures A record matches outbound IP and PTR matches HELO/EHLO hostname.
- **UI/Terminal Component:** Dual-arrow synchronization diagram showing IP -> Hostname -> IP resolution loop with round-trip latency badges.
- **Telemetry & State Fields:** `hasValidPtr` (bool), `isFcrDnsMatched` (bool), `outboundIp` (string), `ptrHostname` (string), `heloHostname` (string).
- **Backend / Engine:** `pkg/mail/auth/fcrdns_validator.go` issuing reverse PTR query followed by forward A/AAAA query.
- **Failure / Guard Rail:** Disables direct SMTP delivery mode if PTR is absent, falling back to preconfigured authenticated relay.

### Improvement 68: Mail Header Hygiene & RFC 5322 Syntax Linter
- **Description:** Deep syntax linter checking Message-ID uniqueness, Date header timezone compliance, MIME multipart boundaries, and List-Unsubscribe header formatting.
- **UI/Terminal Component:** Code-view header inspector highlighting non-compliant RFC headers with one-click auto-repair diff preview.
- **Telemetry & State Fields:** `hasRfcCompliance` (bool), `hasUnsubscribeHeader` (bool), `hasValidMessageId` (bool), `headerByteLength` (int), `linterWarningCount` (int).
- **Backend / Engine:** `pkg/mail/composer/header_linter.go` running RFC 5322 and RFC 8058 conformance test suites.
- **Failure / Guard Rail:** Rejects message generation if `Message-ID` or `Date` header is missing or violates ASCII safety boundaries.

### Improvement 69: IP/Domain Real-Time Blacklist (RBL / DNSBL) Poller
- **Description:** Parallel DNSBL query engine checking IP and sending domain across 25 tier-1 blacklists (Spamhaus ZEN, Barracuda, SORBS, SpamCop, Invaluement).
- **UI/Terminal Component:** Grid of blacklist status cards with response status code explanations (e.g. `127.0.0.2` listing reason) and removal link shortcuts.
- **Telemetry & State Fields:** `isCleanFromRbl` (bool), `hasActiveListing` (bool), `checkedRblCount` (int), `listedRblCount` (int), `listingDetailsJson` (string).
- **Backend / Engine:** `pkg/mail/auth/dnsbl_checker.go` issuing parallel Go goroutines with bounded 800ms timeouts per query.
- **Failure / Guard Rail:** If IP is listed on Spamhaus ZEN, warmup ramp is automatically paused and desktop notification triggered.

### Improvement 70: TLS Cipher & STARTTLS Handshake Diagnostic Monitor
- **Description:** Inspects outbound SMTP connection encryption, negotiating TLS 1.3/1.2, verifying certificate authority validity, and checking cipher suite strength.
- **UI/Terminal Component:** Lock icon with padlock visual strength meter; displays TLS version, cipher suite name, and SNI handshake duration.
- **Telemetry & State Fields:** `hasTlsEncryption` (bool), `isTls13` (bool), `hasValidCertificate` (bool), `tlsHandshakeMs` (int), `cipherSuiteName` (string).
- **Backend / Engine:** `pkg/mail/transport/tls_diagnostics.go` probing MX host on port 25 and 587 with opportunistic STARTTLS.
- **Failure / Guard Rail:** Flags warning if remote MX does not support STARTTLS; prevents transmission if strict TLS mode is active and handshake fails.

---

## 3. Cluster 8: Design System Tokens, Dark/Light Themes, Typography & WCAG AA Contrast (Items 71–80)

### Improvement 71: TokyoNight Zero-Runtime CSS Token Registry
- **Description:** Core design token registry defining color ramps, elevation shadows, border radii, and spacing scales inspired by TokyoNight.
- **UI/Terminal Component:** CSS custom properties (`--tn-bg-dark: #1a1b26`, `--tn-bg-card: #24283b`, `--tn-cyan: #7dcfff`, `--tn-blue: #7aa2f7`).
- **Telemetry & State Fields:** `hasDarkThemeActive` (bool), `themeTokenHash` (string), `customTokenCount` (int), `contrastRatio` (float64).
- **Backend / Engine:** `pkg/web/tokens/theme_registry.go` serving zero-dependency CSS stylesheets with inline ETags.
- **Failure / Guard Rail:** Token registry enforces immutability; unknown token lookups return fallback system variables.

### Improvement 72: Arctic Light Theme & Ambient Luminosity Switcher
- **Description:** Clean, clinical light mode theme providing crisp contrast (`#f8f9fc` canvas, `#ffffff` card surface, `#2e3440` text) for high-ambient-light environments.
- **UI/Terminal Component:** Smooth CSS transition theme toggle button in navigation bar; instantaneous terminal theme sync based on terminal background detection.
- **Telemetry & State Fields:** `hasLightThemeActive` (bool), `isSystemThemeInherited` (bool), `themeTransitionMs` (int), `colorPaletteName` (string).
- **Backend / Engine:** Local client storage sync with session cookie and CLI flag `--theme=arctic|tokyo`.
- **Failure / Guard Rail:** Automatically preserves selected theme in `gitmap.db` user preferences without page reload flash.

### Improvement 73: WCAG 2.1 AA & AAA Automated Contrast Auditor HUD
- **Description:** Real-time contrast validator evaluating background vs foreground color combinations on all text labels, buttons, and status badges.
- **UI/Terminal Component:** Visual debug overlay highlighting elements with contrast ratio badge (e.g. `4.8:1 [AA]`, `7.2:1 [AAA]`, `2.9:1 [FAIL]`).
- **Telemetry & State Fields:** `hasPassedWcagAa` (bool), `hasPassedWcagAaa` (bool), `minimumContrastRatio` (float64), `failedElementCount` (int).
- **Backend / Engine:** Embedded CSS lint checker in `pkg/web/a11y/contrast_calculator.go` computing relative luminance via ITU-R BT.709 formulas.
- **Failure / Guard Rail:** CI/CD test gate fails if any production UI component drops below 4.5:1 for standard text or 3.0:1 for large text.

### Improvement 74: JetBrains Mono & Inter Font Typography Hierarchy Scale
- **Description:** Declarative typography scale coupling `Inter` for prose and navigation with `JetBrains Mono` for tabular metrics, email headers, and CLI codes.
- **UI/Terminal Component:** 8-step font ramp (`--font-size-xs: 0.75rem` up to `--font-size-4xl: 2.25rem`) with matching line heights and tabular figure alignment (`font-variant-numeric: tabular-nums`).
- **Telemetry & State Fields:** `hasTabularNumsEnabled` (bool), `fontScaleMultiplier` (float64), `activeFontFamily` (string), `isMonoRendered` (bool).
- **Backend / Engine:** Embedded font loader serving subset WOFF2 binaries locally from binary assets; zero external Google Font CDN queries.
- **Failure / Guard Rail:** Fallbacks strictly defined to standard system fonts (`system-ui, -apple-system, Segoe UI, monospace`).

### Improvement 75: High-Density Data Density Switcher (Compact / Standard / Relaxed)
- **Description:** Tri-state table and card padding modifier allowing operators to view up to 100 mailboxes per screen or switch to spacious card presentations.
- **UI/Terminal Component:** View density segmented button group (`Compact: 24px rows`, `Standard: 40px rows`, `Relaxed: 56px rows`).
- **Telemetry & State Fields:** `hasCompactDensity` (bool), `hasRelaxedDensity` (bool), `rowHeightPixels` (int), `visibleRowCount` (int).
- **Backend / Engine:** Stored in client localStorage and user session settings table in `gitmap.db`.
- **Failure / Guard Rail:** Restores default standard mode if screen height drops below 768px to prevent touch target overlap.

### Improvement 76: Keyboard Navigation System & Focus Trap Modal Controller
- **Description:** Complete keyboard-first navigation suite: arrow keys for table navigation, `/` for search focus, `Esc` for modal dismissal, and tab focus rings.
- **UI/Terminal Component:** High-visibility double-ring focus indicator (`outline: 2px solid var(--tn-cyan); outline-offset: 2px`) on all active elements.
- **Telemetry & State Fields:** `hasFocusTrapActive` (bool), `hasKeyboardNavActive` (bool), `activeFocusedElementId` (string), `lastHotkeyKeycode` (string).
- **Backend / Engine:** Zero-dependency micro-dispatcher in `pkg/web/static/js/keyboard_navigation.js`.
- **Failure / Guard Rail:** Traps tab focus inside active modals to prevent background interaction; releases immediately upon close.

### Improvement 77: Status Pill Badges with Dual Color & Shape Coding
- **Description:** Accessible status pill system that conveys meaning via distinct geometric icon shapes in addition to color, ensuring clarity for colorblind operators.
- **UI/Terminal Component:** Success = solid check circle + green; Warning = equilateral triangle + amber; Error = octagon + red; Paused = double vertical bar + slate.
- **Telemetry & State Fields:** `hasShapeEncoding` (bool), `statusBadgeType` (string), `iconGlyphSvg` (string), `isHighContrastShape` (bool).
- **Backend / Engine:** Unified Lipgloss helper in `cli/termout/badges.go` and HTML template components.
- **Failure / Guard Rail:** Validated against Protanopia, Deuteranopia, and Tritanopia color blindness simulation filters in headless test runs.

### Improvement 78: SVG Micro-Charts & Sparkline Visual Tokens
- **Description:** Lightweight zero-dependency inline SVG sparklines rendering 14-day warmup deliverability curves inside table rows without heavy charting bundles.
- **UI/Terminal Component:** Inline 80x20px SVG vector sparklines with gradient fill and terminal Unicode block equivalents (` ▂▃▅▆▇█`).
- **Telemetry & State Fields:** `hasSparklineRendered` (bool), `sparklineDataPointCount` (int), `sparklinePeakValue` (float64), `isTrendPositive` (bool).
- **Backend / Engine:** `pkg/web/tokens/sparkline_generator.go` calculating SVG polyline points in Go templates; execution time < 10 microseconds.
- **Failure / Guard Rail:** Gracefully degrades to text percentage delta badge (`+14.2%`) when SVG rendering is unsupported.

### Improvement 79: Motion & Animation Reduced-Motion Preference Observer
- **Description:** Detects user OS accessibility settings (`prefers-reduced-motion: reduce`) and disables all progress pulses, dial rotations, and transitions.
- **UI/Terminal Component:** Media query observer replacing animations with instantaneous static state changes.
- **Telemetry & State Fields:** `hasReducedMotionEnabled` (bool), `animationDurationMs` (int), `transitionCurveName` (string).
- **Backend / Engine:** Client-side CSS media query switch with mirrored CLI option `--no-animation`.
- **Failure / Guard Rail:** Guarantees zero cognitive disturbance or vestibular disorientation for sensitive users.

### Improvement 80: Touch Target Size Validator & Responsive Breakpoint Grid
- **Description:** Mobile and tablet responsiveness suite enforcing minimum 44x44px touch targets on all buttons, checkboxes, and tabs across breakpoints.
- **UI/Terminal Component:** Fluid 12-column flexbox grid adapting across Mobile (`<640px`), Tablet (`640-1024px`), Desktop (`1024-1440px`), and Ultrawide (`>1440px`).
- **Telemetry & State Fields:** `hasValidTouchTargets` (bool), `activeBreakpoint` (string), `viewportWidth` (int), `minTargetSizePixels` (int).
- **Backend / Engine:** CSS container queries and responsive utility stylesheets served locally.
- **Failure / Guard Rail:** Automatically collapses complex multi-metric data tables into expandable vertical cards on mobile viewports.

---

## 4. Cluster 9: SQLite Split-DB Telemetry, Query Speed & Audit Logs (Items 81–90)

### Improvement 81: Microsecond Query Execution Speedometer & Slow-Query Alert HUD
- **Description:** Continuous telemetry profiler measuring execution duration for every warm-up dispatch, mailbox query, and audit log write.
- **UI/Terminal Component:** Real-time query performance bar graph; queries taking `> 5.0ms` trigger an orange alert badge with EXPLAIN QUERY PLAN breakdown.
- **Telemetry & State Fields:** `isQueryFast` (bool), `hasSlowQueries` (bool), `executionDurationUs` (int64), `slowQueryThresholdMs` (int), `executedSql` (string).
- **Backend / Engine:** Wrapper around `database/sql` driver in `pkg/store/telemetry_driver.go` recording execution duration in memory ring buffer.
- **Failure / Guard Rail:** Emits warning log and initiates index optimization recommendation when query duration exceeds 10ms.

### Improvement 82: Write-Ahead Log (WAL) Checkpoint Health & Page Fragmentation Dial
- **Description:** Real-time health monitor tracking SQLite WAL file size (`gitmap.db-wal`), checkpoint frequency, and dirty page fragmentation percentage.
- **UI/Terminal Component:** Radial gauge displaying WAL file memory footprint (0 to 64MB) with status indicator (`PASSIVE`, `FULL`, `RESTART`, `TRUNCATE`).
- **Telemetry & State Fields:** `hasWalHealthy` (bool), `isWalCheckpointNeeded` (bool), `walSizeBytes` (int64), `checkpointMode` (string), `dirtyPageCount` (int).
- **Backend / Engine:** Executes `PRAGMA wal_checkpoint(PASSIVE)` and `PRAGMA freelist_count` through `pkg/store/wal_monitor.go`.
- **Failure / Guard Rail:** Triggers automated `TRUNCATE` checkpoint when WAL file size exceeds 32MB to prevent file system bloat.

### Improvement 83: Casbin RBAC Audit Trail & Security Event Streaming HUD
- **Description:** Real-time security stream logging every mailbox credential access, warmup parameter modification, and sending authorization.
- **UI/Terminal Component:** High-speed streaming event log with actor avatar, Casbin policy rule match (`p, role:admin, mailboxes, write, allow`), and timestamp.
- **Telemetry & State Fields:** `hasSecurityAuditLogged` (bool), `isAccessAllowed` (bool), `actorId` (string), `targetResource` (string), `actionVerb` (string).
- **Backend / Engine:** Interceptor in `pkg/auth/casbin_interceptor.go` recording into dedicated `installation.db` audit ledger.
- **Failure / Guard Rail:** Denied permission attempts trigger immediate security flash alert and optional email push to workspace owner.

### Improvement 84: Split-DB Architecture Sharding & Mailbox Partition Visualizer
- **Description:** Visual diagram displaying the separation of concerns across GitMap's SQLite tiers: `gitmap.db` (app state), `installation.db` (keys/config), and `mail_audit.db` (events).
- **UI/Terminal Component:** Three-cylinder database topology diagram showing active connection pool count, read/write locks, and disk file sizes.
- **Telemetry & State Fields:** `hasSplitDbHealthy` (bool), `isConnectionPoolAvailable` (bool), `gitmapDbSizeBytes` (int64), `auditDbSizeBytes` (int64), `openHandleCount` (int).
- **Backend / Engine:** `pkg/store/split_db_manager.go` managing discrete connection handles and connection pooling boundaries.
- **Failure / Guard Rail:** Prevents cross-database foreign key corruption; guarantees complete isolation during database backups or restores.

### Improvement 85: Database Auto-Vacuum & B-Tree Page Compaction Scheduler
- **Description:** Automated maintenance service analyzing database freelist space and scheduling non-blocking incremental vacuums during low sending activity windows.
- **UI/Terminal Component:** Progress bar indicating recovered disk space after compaction with next scheduled maintenance countdown timer.
- **Telemetry & State Fields:** `hasVacuumActive` (bool), `isCompactionScheduled` (bool), `freelistPageCount` (int), `reclaimedSpaceBytes` (int64).
- **Backend / Engine:** `pkg/store/compaction_worker.go` issuing `PRAGMA incremental_vacuum(100)` in background worker goroutine.
- **Failure / Guard Rail:** Compaction is strictly inhibited during active email burst dispatching to prevent disk I/O contention.

### Improvement 86: SQLite Index Hit-Ratio & Full-Table-Scan Sentinel
- **Description:** Telemetry analyzer detecting missing indexes by auditing `sqlite_stat1` and intercepting table scan events during warm-up state queries.
- **UI/Terminal Component:** Radar indicator displaying cache hit ratio (target `> 99.0%`) and a list of table queries executing without covering indexes.
- **Telemetry & State Fields:** `hasOptimalIndexes` (bool), `hasTableScanDetected` (bool), `indexHitRatio` (float64), `scannedTableCount` (int).
- **Backend / Engine:** `pkg/store/index_sentinel.go` parsing query execution plans from SQLite query planner.
- **Failure / Guard Rail:** Generates suggested `CREATE INDEX` migration statement whenever a table scan touches more than 500 rows.

### Improvement 87: Tamper-Evident SHA-256 Audit Log Chaining
- **Description:** Cryptographic hash-chaining for outbound email audit records: every log entry contains the SHA-256 hash of the preceding entry, forming a tamper-evident audit ledger.
- **UI/Terminal Component:** Integrity shield icon showing verification status (`CHAIN_VALID`, `BLOCKS_VERIFIED: 14,290`) with one-click audit verification scan.
- **Telemetry & State Fields:** `hasValidAuditChain` (bool), `isTamperEvident` (bool), `verifiedBlockCount` (int64), `lastBlockHash` (string), `chainBreakIndex` (int64).
- **Backend / Engine:** `pkg/store/audit_chain.go` computing `SHA256(prev_hash + current_record_bytes)` prior to SQLite commit.
- **Failure / Guard Rail:** Immediate hard alert if any historical record has been modified, deleted, or backdated.

### Improvement 88: In-Memory Prepared Statement Cache & Prepared Handle Tracker
- **Description:** Prepared statement cache manager maximizing SQLite throughput by eliminating repetitive statement compilation and binding.
- **UI/Terminal Component:** Statement cache statistics card displaying hit count, eviction count, and memory allocation size.
- **Telemetry & State Fields:** `hasStatementCached` (bool), `statementHitCount` (int64), `statementMissCount` (int64), `cacheCapacity` (int).
- **Backend / Engine:** `pkg/store/statement_cache.go` managing LRU cache of `*sql.Stmt` pointers with thread-safe mutex guards.
- **Failure / Guard Rail:** Automatically closes evicted statements to prevent SQLite statement handle exhaustion.

### Improvement 89: Live Database Backup & Atomic Snapshot HUD
- **Description:** Zero-downtime hot backup coordinator using SQLite online backup API to export snapshot archives without blocking active warm-up threads.
- **UI/Terminal Component:** Backup drawer with snapshot history table, archive file sizes, SHA-256 checksums, and one-click manual backup trigger.
- **Telemetry & State Fields:** `hasBackupCompleted` (bool), `isBackupInProgress` (bool), `backupFileSizeBytes` (int64), `backupDurationMs` (int), `lastBackupTimestamp` (string).
- **Backend / Engine:** `pkg/store/backup_service.go` invoking SQLite backup API to replicate source database into compressed snapshot file.
- **Failure / Guard Rail:** Verifies backup integrity with `PRAGMA integrity_check` before marking snapshot valid.

### Improvement 90: Multi-Process SQLite Busy Handler & Lock Contention Telemetry
- **Description:** Real-time lock contention profiler capturing database lock wait durations and retry counts across CLI, background daemons, and web dashboards.
- **UI/Terminal Component:** Contention warning badge displaying max lock wait time in milliseconds and retry loop count.
- **Telemetry & State Fields:** `hasLockContention` (bool), `isBusyHandlerTriggered` (bool), `busyRetryCount` (int), `maxLockWaitMs` (int).
- **Backend / Engine:** Configured via `PRAGMA busy_timeout = 5000` with telemetry callbacks in `pkg/store/busy_interceptor.go`.
- **Failure / Guard Rail:** Fails with typed `ErrDatabaseBusy` if lock wait exceeds 5000ms, triggering caller retry with exponential jitter.

---

## 5. Cluster 10: Antigravity IDE Integration, Task Completion Email Triggers & Automation Bridge (Items 91–100)

### Improvement 91: Antigravity IDE Workspace Context & Status Bar Indicator
- **Description:** Extension bridge providing an AGY status bar item displaying active warm-up progress, health score, and mailbox counts directly inside the IDE.
- **UI/Terminal Component:** AGY status bar item: `$(mail) GitMap: 12 Mailboxes (98% Health, Day 7/14 Ramp)` with click-to-open dashboard shortcut.
- **Telemetry & State Fields:** `hasAgyConnected` (bool), `isWorkspaceActive` (bool), `statusMessage` (string), `activeMailboxCount` (int), `aggregateHealthPct` (float64).
- **Backend / Engine:** Local Unix socket / named pipe RPC server in `pkg/ide/agy_bridge.go` exchanging JSON messages with IDE extension.
- **Failure / Guard Rail:** Disconnects silently and re-polls every 5 seconds if IDE process terminates.

### Improvement 92: Autonomous Task Completion Email Dispatcher
- **Description:** Automated email dispatcher that sends formatted executive summaries when long-running background tasks (e.g. CI/CD runs, test suites, major refactors) complete.
- **UI/Terminal Component:** Notification settings drawer allowing developers to configure destination addresses, success/failure filters, and template previews.
- **Telemetry & State Fields:** `hasTaskCompletionTriggered` (bool), `hasSentNotification` (bool), `parentTaskId` (string), `taskDurationSeconds` (int), `taskExitCode` (int).
- **Backend / Engine:** `pkg/bridge/task_notifier.go` triggered by GitMap task runner completion hooks; routes email through warmed mailbox cluster.
- **Failure / Guard Rail:** Rate-limited to max 1 email per 10 minutes to prevent notification storms; falls back to local terminal alert if mailboxes are cold.

### Improvement 93: Antigravity Slash-Command Bridge (`/mail status`, `/mail ramp`)
- **Description:** Native slash-command handlers inside Antigravity AI chat allowing agents to inspect deliverability and adjust ramp parameters mid-session.
- **UI/Terminal Component:** Formatted Markdown card responses rendered inside AGY chat stream with interactive action buttons.
- **Telemetry & State Fields:** `hasCommandHandled` (bool), `isCommandValid` (bool), `executedSlashCommand` (string), `responsePayloadJson` (string).
- **Backend / Engine:** `pkg/bridge/slash_handler.go` dispatching requests to existing CLI command executors.
- **Failure / Guard Rail:** Rejects unauthorized command execution without operator confirmation if modifying active ramp rates.

### Improvement 94: Subagent Swarm Deliverability Alert Webhook
- **Description:** Real-time event publisher broadcasting deliverability alerts (e.g. sudden spam surge, blacklist listing) to active AI subagent swarms.
- **UI/Terminal Component:** Swarm notification banner in AGY UI indicating agent intervention requested to analyze bounce logs.
- **Telemetry & State Fields:** `hasAlertDispatched` (bool), `hasSubagentSubscribed` (bool), `alertSeverity` (string), `affectedDomain` (string), `activeAgentCount` (int).
- **Backend / Engine:** In-process event bus in `pkg/bridge/swarm_bus.go` broadcasting over AGY IPC channels.
- **Failure / Guard Rail:** Drops duplicate alerts within a 15-minute sliding window to avoid agent distraction.

### Improvement 95: Crash Dump & Panic Diagnostic Email Generator
- **Description:** Automatic panic-recovery handler formatting application crash stack traces into sanitized, privacy-filtered diagnostic emails sent to the lead developer.
- **UI/Terminal Component:** Pre-send modal with redacted environment variables, sanitized paths, and one-click dispatch button.
- **Telemetry & State Fields:** `hasPanicCaptured` (bool), `hasCrashDumpSent` (bool), `sanitizedStackLength` (int), `crashTimestamp` (string).
- **Backend / Engine:** `pkg/bridge/crash_reporter.go` wrapping top-level command dispatchers with `recover()` blocks.
- **Failure / Guard Rail:** Redacts all passwords, API tokens, and private SSH keys using regex filters before email generation.

### Improvement 96: CI/CD Pipeline Failure Digest & Diff Breakdown Template
- **Description:** Rich HTML/Plain-text email template summarizing failed pipeline steps, failing test assertions, and git commit diffs.
- **UI/Terminal Component:** Dual-pane email preview in dashboard showing HTML and fallback ASCII terminal rendering.
- **Telemetry & State Fields:** `hasPipelineFailed` (bool), `hasDigestRendered` (bool), `failedStepName` (string), `commitSha` (string), `authorEmail` (string).
- **Backend / Engine:** `pkg/bridge/pipeline_digest.go` rendering Go templates in `templates/emails/pipeline_failure.html`.
- **Failure / Guard Rail:** Compiles to both multipart/alternative HTML and strict 80-column plain-text for terminal email clients.

### Improvement 97: Antigravity AI Memory Synchronization Hook
- **Description:** Automatic synchronization hook archiving daily warm-up telemetry, reputation metrics, and deliverability insights into `.ai-memory/warmup/`.
- **UI/Terminal Component:** Synchronization status indicator showing memory file size, last sync date, and active token count.
- **Telemetry & State Fields:** `hasMemorySynced` (bool), `isMemoryUpToDate` (bool), `syncedRecordCount` (int), `memoryFilePath` (string).
- **Backend / Engine:** `pkg/bridge/memory_sync.go` writing markdown and JSON summary ledgers into `.ai-memory/warmup/daily_reputation.md`.
- **Failure / Guard Rail:** Overwrites memory atomically using temporary file rename to prevent file corruption.

### Improvement 98: Local Headless Automation Daemon (`gitmap lemon warmup`)
- **Description:** Cross-platform background worker service executing scheduled warm-up ramps, seed interactions, and health checks as an OS daemon.
- **UI/Terminal Component:** Service controller widget showing daemon PID, uptime, memory consumption (RSS), and start/stop/restart toggles.
- **Telemetry & State Fields:** `isDaemonRunning` (bool), `hasAutoRestart` (bool), `daemonPid` (int), `uptimeSeconds` (int64), `rssMemoryBytes` (int64).
- **Backend / Engine:** `pkg/daemon/lemon_service.go` integrating with Windows Service Control Manager, macOS launchd, and Linux systemd.
- **Failure / Guard Rail:** Heartbeat monitor restarts worker thread automatically if unhandled network exception occurs.

### Improvement 99: Dynamic SMTP Relay Load Balancer & Fallback Switch
- **Description:** Automated transport failover engine switching between local postfix, AWS SES, SendGrid, and Mailgun based on real-time rate limits and latency.
- **UI/Terminal Component:** Active transport routing diagram with green active pipes and latency counters per relay provider.
- **Telemetry & State Fields:** `hasRelayHealthy` (bool), `isFallbackActive` (bool), `activeRelayProvider` (string), `relayLatencyMs` (int), `failoverCount` (int).
- **Backend / Engine:** `pkg/mail/transport/relay_balancer.go` evaluating circuit breakers and routing emails across healthy pools.
- **Failure / Guard Rail:** Fallback occurs in `< 250ms` upon connection timeout; alerts operator if all available relays enter failure state.

### Improvement 100: Unified Outbound Emailing Executive KPI Dashboard
- **Description:** Master summary view uniting all 100 improvements: warm-up ramp curves, deliverability health, seed interactions, security auth matrices, and DB performance into a single executive command center.
- **UI/Terminal Component:** Full-width panoramic dashboard featuring TokyoNight glassmorphic cards, sparkline grids, live activity tickers, and instant action buttons.
- **Telemetry & State Fields:** `hasDashboardInitialized` (bool), `isSystemHealthy` (bool), `activeMailboxTotal` (int), `aggregateReputationScore` (float64), `totalSentCount` (int64), `inboxPlacementPct` (float64).
- **Backend / Engine:** Single-query aggregation endpoint `pkg/web/handlers/executive_dashboard.go` reading from unified Split-DB cache.
- **Failure / Guard Rail:** Renders complete dashboard within 35ms; supports zero-JavaScript offline viewing mode via pure HTML/CSS.

---

## 6. Implementation Summary & Verification Checklist

| Cluster | Improvement IDs | Subsystem / Package | Key Metrics & Telemetry | Verification Invariant |
| :--- | :--- | :--- | :--- | :--- |
| **Cluster 7** | 61–70 | `pkg/mail/auth/`, `pkg/mail/spam/` | SPF depth, DKIM keys, DMARC RUA, Spam score, PTR/FCrDNS | Zero external DNS blocking; timeouts < 800ms |
| **Cluster 8** | 71–80 | `pkg/web/tokens/`, `cli/termout/` | TokyoNight/Arctic tokens, WCAG AA 4.5:1, JetBrains Mono | 100% zero-CDN; local WOFF2; reduced-motion |
| **Cluster 9** | 81–90 | `pkg/store/`, `mail_audit.db` | Execution < 5ms, WAL checkpoint, Casbin logs, SHA-256 chain | Tamper-evident hash chain; no cross-DB FKs |
| **Cluster 10** | 91–100 | `pkg/ide/`, `pkg/bridge/`, `pkg/daemon/` | AGY status bar, Task triggers, Slash commands, Executive HUD | Instant fallback; zero data leaks; sanitized logs |

---

## 7. Relative Path & Positive Boolean Conformance

- All file references across this specification use strict repository-relative paths (`02-spec/21-app/...`, `pkg/...`, `.ai-memory/...`).
- All state fields, flags, and entity attributes use strictly positive boolean naming (`hasValidSpf`, `isAligned`, `hasPassedWcagAa`, `isSystemHealthy`), avoiding negative prefixes (`no`, `not`, `non`).
- Functions adhere to the 15-line limit and return errors directly without leaf exit handlers.
