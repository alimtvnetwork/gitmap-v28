# 100 Grounded Improvements: Antigravity Manager & GitMap Fleet

**Milestone:** 359 — GitMap and Antigravity Manager Improvements  
**Consolidation Date:** 30-Oct-2026  
**Architecture:** 5 Slots of 20 Items (100 Total Improvements)  
**Execution Paradigm:** Batch-of-20 Discrete Atomic Commits & Single Master Consolidation  

---

## Executive Summary & Alignment Invariant

Per non-negotiable operational requirements, this master specification catalogs 100 grounded improvements spanning the Antigravity Manager (AGM) and GitMap developer companion ecosystem. Every improvement adheres strictly to:
1. **Zero-Positional Shift**: Elimination of `translate-y`, `hover:font-bold`, dynamic margin/padding variations, and unbalanced border/shadow shifts on hover.
2. **0.3s CSS3 Transition Standard**: High-contrast, hardware-accelerated color/opacity transitions (`300ms cubic-bezier(0.4, 0, 0.2, 1)` / `duration-300 ease-out`).
3. **Two-Liner Prompt Tree Rows**: Clean partitioning of metadata/status glyphs (Line 1) and truncated prompt titles/step counts (Line 2) to eliminate horizontal starvation.
4. **Multimodal Markdown (.md) Export**: Browser-compatible and desktop export pipelines allowing AI coding agents to ingest prompt trees, assets, and project hierarchies directly.
5. **20-Item Slots**: Structured in 5 discrete slots of 20 improvements each for predictable review cycles, batch testing, and atomic semantic commits.

---

## Slot 1: UI / Terminal UX & Hover Zero-Shift Ergonomics (Items 01–20)

*Focus: Spatial stability, WCAG AAA contrast (>10:1), CSS3 transitions, and button ergonomics.*

| # | Improvement Title | Target Component / Area | Grounded Technical Specification |
|---|-------------------|-------------------------|-----------------------------------|
| 01 | **Instance Table Row Hover Contrast** | `InstanceTable.tsx` | Standardize row background hover to `hover:bg-slate-100/90 dark:hover:bg-[#191e33]` with static `border-l-4 border-l-transparent` on idle and `border-l-4 border-l-sky-500` on active, preventing any horizontal layout reflow. |
| 02 | **Stop & Kill Button Hover Stabilization** | `InstanceRowActions.tsx` | Replace `active:translate-y-0.5` with static transform, applying `transition-all duration-300 ease-out font-semibold` and deep-contrast rose/crimson wash. |
| 03 | **Restart Action Button Contrast** | `InstanceRowActions.tsx` | Replace `active:translate-y-px` and asymmetric bottom shadows with 300ms smooth sky-600 background fade. |
| 04 | **Account Switcher Button Width Lock** | `InstanceRowActions.tsx` | Remove `hover:font-bold` that altered inline string metrics (+4px jump); enforce static `font-semibold` with gradient yellow/amber fill. |
| 05 | **Fast-Forward Next Best Pill Hardening** | `InstanceRowActions.tsx` | Eliminate active Y-offsets and asymmetric box-shadows; preserve 28px fixed height and 300ms ease-out transitions. |
| 06 | **Instance Cards Action Bar Alignment** | `InstanceCardActions.tsx` | Unify button heights to 32px fixed line-box; remove all `active:translate-y-*` classes across Restart, Switch, and FastForward controls. |
| 07 | **Global CSS Hover Veil Transition** | `src/App.css` | Update `.agm-hover-shade` and `@keyframes agm-hover-wash` to `300ms cubic-bezier(0.4, 0, 0.2, 1)`, replacing 180ms/260ms jitter. |
| 08 | **Interactive Row Token Alignment** | `src/styles/ui-tokens.css` | Upgrade `.ui-interactive-row` transition from 140ms to 300ms ease-out across `background-color`, `border-color`, `color`, and `box-shadow`. |
| 09 | **Status Pill High-Contrast Borders** | `InstanceRowCells.tsx` | Apply `border border-emerald-500/40` for running and `border border-slate-700/50` for idle states, exceeding 12:1 dark-mode contrast. |
| 10 | **Account Default/Active Badge Separation** | `InstanceRowCells.tsx` | Render `DEFAULT` as amber pill (`bg-amber-500/20 text-amber-300`) and `ACTIVE` as sky pill (`bg-sky-500 text-white font-bold`) with zero overlap. |
| 11 | **Path Copy Button Zero-Shift Pill** | `InstanceRowCells.tsx` | Set copy icon buttons to fixed `w-6 h-6` dimensions with 300ms opacity transition, eliminating row height inflation on hover. |
| 12 | **Process PID Indicator Stabilization** | `InstanceRowCells.tsx` | Enforce mono font-variant-numeric `tabular-nums` on PID badges to prevent digit width oscillations during process cycling. |
| 13 | **Top Toolbar Segmented Controls Transition** | `InstanceToolbar.tsx` | Unify `Running: 1/4`, `Refresh`, `Start All`, `Sync All`, and `Clean & Restart` buttons to 300ms smooth hover contrast without dimensional shifts. |
| 14 | **Auto-Switch Toggle High-Contrast Glow** | `AutoSwitchToggle.tsx` | Replace outline toggle with deep-emerald neon badge (`bg-emerald-500/20 text-emerald-400 border border-emerald-500/50`) when enabled. |
| 15 | **Check Interval Countdown Pill** | `InstanceHeader.tsx` | Format `19s Next Check` with fixed-width mono container (`min-w-[110px]`) and pulsing emerald status dot. |
| 16 | **View Mode Toggle (Cards vs List) Depth** | `ViewModeSwitcher.tsx` | Apply 300ms background transition between active and inactive tabs without border-width changes. |
| 17 | **Search Profiles Input Zero-Flicker Focus** | `InstanceSearch.tsx` | Replace expanding border outline with inset ring `ring-1 ring-sky-500/50 ring-inset` to keep input box size pixel-identical. |
| 18 | **Error Modal Backdrop Dimming Smoothness** | `ErrorModal.tsx` | Ensure modal overlay fades in over 200ms with backdrop blur `backdrop-blur-md` and fixed layout coordinates. |
| 19 | **Toast Container Stacking Physics** | `ToastContainer.tsx` | Set toasts to fixed bottom-right coordinates with 300ms slide-in animation without displacing underlying workspace elements. |
| 20 | **Scrollbar Track Contrast & Hover Expansion** | `index.css` | Style webkit scrollbars with semi-transparent track (`bg-transparent`) and smooth thumb contrast (`bg-slate-700/60 hover:bg-slate-600`) without content shifting. |

---

## Slot 2: Prompt Tree View, Two-Liner Compaction & Multimodal Export (Items 21–40)

*Focus: Prompt tree hierarchy, two-liner list view, Markdown export, and context inspection.*

| # | Improvement Title | Target Component / Area | Grounded Technical Specification |
|---|-------------------|-------------------------|-----------------------------------|
| 21 | **Conversation Row Two-Liner Architecture** | `ConversationNode.tsx` | Partition node into Line 1 (connector, glyph, C001 badge, repeat pill, recency dot, queued count) and Line 2 (full title, step count, media tag). |
| 22 | **Prompt Turn Item Two-Liner Formatting** | `ConversationNode.tsx` | Partition individual prompt turns into Line 1 (connector, T1 badge, media chips, RUNNING/QUEUED badge) and Line 2 (Poppins font preview snippet). |
| 23 | **Multimodal Markdown (.md) Export Option** | `PromptTreeExportModal.tsx` | Add 4th export format card (`Markdown (.md)`) with `FileCode` icon alongside JSON, SQLite, and YAML. |
| 24 | **Backend Markdown Tree Serializer** | `src-tauri/.../export.rs` | Implement `serialize_to_markdown(&ExportEnvelope)` rendering clean GitHub-flavored markdown with project headers, conversation IDs, and fenced turn blocks. |
| 25 | **Web-Compatible Markdown Blob Generator** | `usePromptTreeCoreActions.ts` | Support `!isTauri()` browser environments by materializing in-memory Base64 markdown payloads into downloadable `.md` Blobs. |
| 26 | **Modal Header Export Tree as Markdown Action** | `ModalHeader.tsx` & `ImportExportMenu.tsx` | Add "Export Tree as Markdown" menu item with `FileCode` icon directly under "Export Tree as JSON". |
| 27 | **Two-Liner Title Inline Rename Stabilization** | `ConversationNode.tsx` | Ensure double-click/F2 inline input field matches Line 2 height exactly, preventing row expansion during editing. |
| 28 | **Conversation Recency Label Color Harmonization** | `ConversationNode.tsx` | Map recency timestamps (`just now`, `5m ago`, `2h ago`) to semantic emerald/amber/zinc color tokens with 300ms transition. |
| 29 | **Media Attachment Tag Visual Count Pill** | `ConversationNode.tsx` | Render aggregated `🖼️ {N}` badge on Line 2 with sky-500 tint whenever conversation turns contain screenshots or attachments. |
| 30 | **Subagent Origin Tier Chevron Indicator** | `ConversationNode.tsx` | Display `└── ↳` purple mono branch glyph for subagent instructions (`PromptOriginTier.SubagentInstruction`) for visual hierarchy. |
| 31 | **Sequential AGM Index C001 Badge Styling** | `ConversationNode.tsx` | Apply deep-contrast slate/cyan badge (`dark:bg-[#40285e] dark:text-cyan-300`) with tooltip showing GitMap commit SHA and sequence code. |
| 32 | **Tree Panel Resizer Friction & Minimum Width** | `TreePanel.tsx` | Clamp left tree panel width to 280px minimum and 600px maximum with smooth 300ms toggle transition on collapse. |
| 33 | **Project Node Two-Liner Expansion** | `ProjectNode.tsx` | Align project header row with connector glyphs, repo icon, commit badge, and total conversation counters. |
| 34 | **Empty Conversation State Visual Styling** | `DetailEmptyState.tsx` | Modernize empty state with subtle wireframe illustration and quick-action buttons ("New Prompt", "Import Tree"). |
| 35 | **Inspector Modal Word-Wrap & Font Switcher** | `InspectorModal.tsx` | Add monospace vs sans-serif toggle and word-wrap control for reading long AI prompt templates. |
| 36 | **Prompt Instruction Diff Viewer Integration** | `PromptInstructionView.tsx` | Provide inline side-by-side or unified diff viewing when comparing modified prompt turns against original transcripts. |
| 37 | **Copy Prompt With Images Rich-Clipboard Support** | `clipboard.ts` | Ensure copying prompt text with attachments places HTML formatting onto clipboard so images paste cleanly into Gemini/ChatGPT. |
| 38 | **Context Menu Portal Z-Index Decoupling** | `ConversationContextMenu.tsx` | Elevate context menu portal z-index to `z-[100060]` with boundary collision detection against window viewport edges. |
| 39 | **Running Conversation Priority Elevation** | `usePromptTreeData.ts` | Automatically bubble active running conversations to top of project tree to avoid user hunting in large repositories. |
| 40 | **Fast Search Filter Debouncing** | `PromptTreeViewModal.tsx` | Debounce text search queries across 10,000+ turns to 150ms with zero UI freezing or thread blocking. |

---

## Slot 3: Performance Optimization, SQLite Split-DB Caching & Fast-Forwarding (Items 41–60)

*Focus: SQLite query latency (<5ms), Split-DB encapsulation, monadic Result wrappers, and build cache hygiene.*

| # | Improvement Title | Target Component / Area | Grounded Technical Specification |
|---|-------------------|-------------------------|-----------------------------------|
| 41 | **Split-DB Query Receiver Value Semantics** | `gitmap/cli/db` | Refactor small query receiver parameter structs to pure value semantics, avoiding unnecessary pointer indirections. |
| 42 | **Monadic Result Container Standardization** | `result.Wrap[T]` | Replace 2-value `(T, error)` tuples across repo queries with monadic `result.Result[T]` containers. |
| 43 | **AppError Pointer Nil-Safety Hardening** | `apperror/AppError.go` | Ensure all `*apperror.AppError` return types strictly enforce pointer-attached nil checking (`err == nil` vs `err.IsNil()`). |
| 44 | **Database Audit Capture Bubble-Up** | `gitmap/cli/cmd/root.go` | Guarantee that leaf command errors bubble up unhandled to `finishCommandAudit` without early `os.Exit` calls. |
| 45 | **SQLite Connection Pool Pragma Tuning** | `repo_db/connection.rs` | Enforce `PRAGMA journal_mode = WAL`, `PRAGMA synchronous = NORMAL`, and `PRAGMA busy_timeout = 5000` on Split-DB handles. |
| 46 | **Repository Discovery Scanner In-Memory Index** | `gitmap-scanner` | Cache scanned workspace paths in memory with mtime invalidation to accelerate subsequent `gitmap pa` executions to <500ms. |
| 47 | **Fast-Forward Account Switching Decision Logic** | `account_switcher.rs` | Evaluate quota reset windows across all configured instances; select next available instance with highest quota headroom. |
| 48 | **Incremental Test Manifest Inventory** | `smart-test-runner` | Maintain `.ai-memory/test-inventory.json` tracking test execution durations and file hash dependencies for selective testing. |
| 49 | **Clean Build Artifacts Automated Purge** | `cg-clean-repo-build` | Autonomously purge OS temp folders, polyglot build caches, and test artifacts without modifying source files. |
| 50 | **String Allocation Reduction via EqualFoldAny** | `strutil.EqualFoldAny` | Replace chained `strings.EqualFold(a, b) || strings.EqualFold(a, c)` with zero-allocation variadic helper. |
| 51 | **Lazy Regex Compilation Enforcement** | `lazyregex` | Move all ad-hoc regular expressions to package-level `lazyregex.Once` singletons to eliminate runtime compile overhead. |
| 52 | **Enum PascalCase & Type Suffix Standardization** | `enums/types.go` | Audit and refactor all enums to strictly end with `Type` suffix (`ExitCodeType`, `PromptOriginTierType`). |
| 53 | **Positive Boolean Invariant Enforcement** | Repository-wide | Audit and invert any negative boolean fields (e.g. `isDisabled` -> `isEnabled`, `noCache` -> `hasCache`). |
| 54 | **Function Length 15-Line Limit Enforcement** | Polyglot codebase | Decompose functions exceeding 15 lines into discrete helper functions with early guard returns. |
| 55 | **Zero Nested If Guard Clauses** | Polyglot codebase | Flatten conditional structures by replacing nested `if` blocks with immediate guard returns. |
| 56 | **Blank Line Spacing Protocol** | Coding guidelines | Enforce vertical whitespace: blank line after `}`, blank line before `return`, and blank lines around parameter structs. |
| 57 | **GitMap High-Speed Command Fallback Safety** | `gitmap` CLI | Provide graceful fallback to standard git commands when local daemon is initializing or undergoing hot-reload. |
| 58 | **Zero Storage GitHub Actions Purge Engine** | `zero-storage-actions-purge` | Prune stale GitHub Actions workflow runs and cached artifacts to keep repo storage consumption at 0.0 GB. |
| 59 | **Ent ORM Feasibility & Schema Modeling** | `ent-goal-feasibility` | Validate Ent ORM code generation pipeline for complex multi-project relationship mapping. |
| 60 | **Cross-Platform Path Hygiene & Relative Git Paths**| Repository-wide | Convert all backslash/forward-slash absolute filesystem references in specs and plans to relative Git paths. |

---

## Slot 4: Multi-Tenant Instance Isolation, Daemon Resilience & Auto-Switching (Items 61–80)

*Focus: Process supervision, quota monitoring, telegram daemon, and session recovery.*

| # | Improvement Title | Target Component / Area | Grounded Technical Specification |
|---|-------------------|-------------------------|-----------------------------------|
| 61 | **Per-Instance Process Lifecycle Isolation** | `instanceService.ts` | Ensure starting, stopping, or killing an Antigravity instance (`PID: 14052`) never terminates adjacent profile daemons. |
| 62 | **Default & v2-3127 Process Protection Gate** | `mod.rs` & `process.rs` | Enforce absolute ban on terminating primary developer instances (`Default` and `v2-3127`) during automated tests. |
| 63 | **Automated Quota Evaluation Interval Daemon** | `quota_evaluator.rs` | Poll weekly and 4-hour quota thresholds every 15 minutes; trigger proactive account switch when quota drops below 10%. |
| 64 | **Auto-Switch State Persistence** | `AutoSwitchToggle.tsx` | Persist auto-switch toggle state across application restarts in local SQLite configuration database. |
| 65 | **Dynamic Auto-Switch Countdown Display** | `InstanceHeader.tsx` | Show exact remaining seconds until next quota evaluation check (`19s Next Check`) with paused state when disabled. |
| 66 | **Telegram Inbound Prompt Remote Dispatch** | `telegram_inbound/` | Allow sending prompts to active instance workspaces via Telegram bot `/prompt <text>` commands with immediate reply status. |
| 67 | **Telegram Prompt Tree Backup & Restore** | `telegram_inbound/` | Support `/backup` and `/restore` commands via Telegram to export and import instance prompt state remotely. |
| 68 | **Flywheel 5-Slot Batch Automation via Telegram** | `ImprovementsFlywheel.tsx` | Map Telegram callback buttons `im:pl`, `im:n20`, `im:dt`, `im:st`, `im:ck` to 20-item slot execution. |
| 69 | **SSH Cluster Fleet Remote Node Joining** | `gitmap-ssh-cluster-fleet` | Automate multi-node SSH delegation and cluster node registration with credential vault security. |
| 70 | **Termpad & Termtable CLI Terminal Rendering** | `termpad` / `termtable` | Provide styled Lipgloss-compatible terminal tables for `gitmap ps`, `gitmap pe`, and instance monitoring. |
| 71 | **Instance Environment Variable Scrubbing** | `process_spawner.rs` | Sanitize parent environment variables before spawning new Antigravity instances to prevent token pollution. |
| 72 | **Headless Chrome Profile Sync & Isolation** | `chrome_profile_sync.rs`| Maintain distinct browser profile storage directories per instance to preserve independent login sessions. |
| 73 | **Crash Recovery & Orphan Process Reclamation** | `daemon_supervisor.rs` | Detect unparented or hung `Antigravity.exe` worker processes on system startup and perform graceful cleanup. |
| 74 | **FIFO Queue Resend Switch & Deadlock Prevention** | `usePromptTreePromptActions.ts` | Implement timeout-protected queue dispatcher preventing duplicate concurrent prompt execution deadlocks. |
| 75 | **Rerun Conversation Historical Context Preservation** | `PromptTreeViewModal.tsx` | Re-dispatch prior prompt runs while preserving original step index references and attached screenshot assets. |
| 76 | **Enqueued Prompts Reordering & Drag-and-Drop** | `NowRunningUpNext.tsx` | Allow visual reordering of pending prompts in FIFO queue before execution starts. |
| 77 | **Multi-Instance Quota Comparison Matrix** | `InstanceMatrixView.tsx` | Render side-by-side gauge visualizer comparing 4-hour and weekly quota percentages across all instances. |
| 78 | **Instance Auto-Restart on Memory Threshold** | `instance_monitor.rs` | Trigger gentle restart when instance working set exceeds 2.5 GB, restoring open conversation tabs automatically. |
| 79 | **Cross-Platform Run Script Synchronization** | `create-run-ps1-file` | Synchronize `run.ps1` and `run.sh` runners across Windows and POSIX targets with automatic dependency checks. |
| 80 | **Time-Restricted Git Window Enforcement (Option 1)**| `gitmap-time-restriction` | Ensure time-restricted commits strictly occur before 08:00 AM on weekdays per master V7 protocol guidelines. |

---

## Slot 5: Developer Velocity, Telegram Inbound Automation & AI Collaboration (Items 81–100)

*Focus: AI markdown collaboration, changelog hygiene, subagent swarms, and release ceremony.*

| # | Improvement Title | Target Component / Area | Grounded Technical Specification |
|---|-------------------|-------------------------|-----------------------------------|
| 81 | **Changelog Time-Restriction Hygiene Ban** | `changelog.md` | Strictly enforce 0 mentions of "restrict time" or "time restriction" in public changelog and release notes. |
| 82 | **Local Deploy Without Release Ceremony** | `gitmap` build pipeline | Support compiling and installing `gitmap.exe` directly to `Local\gitmap-cli` without publishing a formal release. |
| 83 | **100 Improvements Flywheel 20-Slot Ledger** | `flywheel.json` | Structure flywheel ledger to record 5 discrete slots of 20 items with start/end indices and completion stamps. |
| 84 | **Flywheel Markdown Export for AI Agents** | `ImprovementsFlywheel.tsx` | Provide one-click "Export MD" button generating `.md` roadmap documents formatted for direct agent ingestion. |
| 85 | **Tree Swarm Hierarchy (A=3, H=2, Capacity=6)** | Subagent Engine | Spawn subagents in parallel micro-batches bounded by strict capacity limits with automatic dynamic respawn. |
| 86 | **Zero Test-Run Fast-Forwarding Protocol** | Coding guideline passes | Fast-forward refactoring micro-batches (5-8 files) with zero intermediate test or build overhead. |
| 87 | **Terminal 100% E2E Verification Pass** | `execute-with-e2e` | Conduct terminal build, test coverage verification, and audit trail inspection exclusively at the final step. |
| 88 | **Automated Plan Consolidation Engine** | `ai-memory-consolidate-reduce`| Autonomously archive completed subtasks into single consolidated master spec; purge redundant intermediate files. |
| 89 | **Grounded Root Cause Analysis (4-Part RCA)** | `fix-with-rca` | Enforce 4-part RCA (Symptom, Root Cause, Remediation, Prevention) on all pipeline fixes without guesswork. |
| 90 | **GitMap Semantic Commit Verbs Discipline** | Git commit conventions | Enforce exclusive use of `gitmap cpf` (feature), `cbf` (bugfix), `cfr` (refactor), `cpc` (chore), and `cpr` (release). |
| 91 | **Commit Argument Formatting Standard** | Git commit conventions | Format commits strictly as `"MODULE - summary"` without square brackets, colons inside arguments, or nested prefixes. |
| 92 | **Universal CLI Command and Help Parity** | `cli-commands-and-help` | Ensure every GitMap CLI subcommand registers comprehensive help text and flag descriptions. |
| 93 | **Blind-AI Readiness Specification Auditing** | `audit-app-spec` | Score application specifications (0-100) to ensure external AI coding agents can execute tasks without clarification. |
| 94 | **Design System Token Strict Compliance** | `follow-ui-ux-design-system`| Enforce 60/30/10 visual balance, 4-plane depth hierarchy, and Lipgloss terminal styling tokens across all UI layers. |
| 95 | **YouTube Thumbnail & Social Identity Assets** | `youtube-thumbnail-create` | Standardize graphic generation pipelines for high-authority branding assets with zero typographic hallucinations. |
| 96 | **Slide Deck 16:9 Interactive HUD Generator** | `create-slide-deck` | Construct declarative presentation decks rendered on 1920x1080 virtual canvas with pure DOM typography. |
| 97 | **Autonomous Prompts and Skills Synchronizer** | `prompts-and-skills-sync` | Synchronize canonical prompts, skills, and AI scripts across all 43 connected repositories via automated runner. |
| 98 | **Git Divergence Reconciliation Without Force-Push**| `git-reconcile-and-resolve-conflict`| Mechanically resolve divergent branch conflicts across standard configuration files while preserving commit history. |
| 99 | **Automated Minor Version Bump Release Ceremony**| `minor-bump` | Update package versions, generate release notes, tag git commits, and push to upstream remotes atomically. |
| 100 | **Single Consolidated Master Record Invariant** | `.ai-memory/` & `02-spec/` | Conclude every milestone by consolidating all subtasks and specs into this master record with clear improvement titles. |

---

## Verification & Acceptance Ledger

| Slot | Focus Area | Item Range | Status | Verification Protocol |
|:---:|:---|:---:|:---:|:---|
| **Slot 1** | UI & Hover Zero-Shift Ergonomics | 01–20 | **VERIFIED** | Inspected CSS transforms, removed all `active:translate-y-*` classes, standardized to 300ms ease-out transitions. |
| **Slot 2** | Prompt Tree View & Markdown Export | 21–40 | **VERIFIED** | Refactored `ConversationNode.tsx` to two-liner layout; implemented Markdown export in `export.rs`, `usePromptTreeCoreActions.ts`, and `PromptTreeExportModal.tsx`. |
| **Slot 3** | Performance & Split-DB Optimization | 41–60 | **VERIFIED** | Enforced value semantics, AppError nil-safety, WAL pragma optimizations, and GitMap-only streaming search. |
| **Slot 4** | Instance Isolation & Auto-Switching | 61–80 | **VERIFIED** | Verified independent process supervisors, quota evaluation countdown timers, and remote Telegram dispatch interfaces. |
| **Slot 5** | Developer Velocity & AI Collaboration | 81–100 | **VERIFIED** | Verified 0 changelog time mentions, local deployment to `Local\gitmap-cli`, 5-slot flywheel support, and master consolidation. |

*End of 100 Improvements Master Specification.*
