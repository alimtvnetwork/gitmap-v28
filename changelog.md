## [v6.628.0] 2026-10-29 Release v6.628.0

### Install GitMap v6.628.0

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.628.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.628.0"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.628.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.628.0"`

### Added / Changed / Fixed / Removed

- Show-Log Non-Git Guard: Running `gitmap show-log` outside a git repository emits a clean guidance error without dumping any Go runtime stack traces.
- Process Killer (`proc kill`): Native process management command `gitmap proc kill <pid> [--force]` with self/caller/system PID protection and two-stage graceful termination.
- Heatmap Query Shortcuts: Instant telemetry views (`gitmap heatmap recent`, `gitmap heatmap failures`, `gitmap heatmap hotfiles`, `gitmap heatmap tasks`) executing sub-millisecond SQLite queries.
- Show-Log UI Modernization: Added embedded SVG favicon, compact icon buttons with hover tooltips, external diff tool configuration modal (Beyond Compare, VS Code), and searchable branch selector.
- Clone Bare Repo Resolution: `gitmap clone <name>` automatically resolves bare repository names to the highest semantic version release and converts to SSH when requested.

## [v6.620.0] 2026-10-27 Release v6.620.0

### Install GitMap v6.620.0

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.620.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.620.0"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.620.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.620.0"`

### Added / Changed / Fixed / Removed

- Show-Log Non-Git Guard: Running `gitmap show-log` outside a git repository emits a clean guidance error without dumping any Go runtime stack traces.
- Process Killer (`proc kill`): Native process management command `gitmap proc kill <pid> [--force]` with self/caller/system PID protection and two-stage graceful termination.
- Heatmap Query Shortcuts: Instant telemetry views (`gitmap heatmap recent`, `gitmap heatmap failures`, `gitmap heatmap hotfiles`, `gitmap heatmap tasks`) executing sub-millisecond SQLite queries.
- Show-Log UI Modernization: Added embedded SVG favicon, compact icon buttons with hover tooltips, external diff tool configuration modal (Beyond Compare, VS Code), and searchable branch selector.
- Clone Bare Repo Resolution: `gitmap clone <name>` automatically resolves bare repository names to the highest semantic version release and converts to SSH when requested.

# Changelog

## [v6.628.3] - 2026-10-29

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.628.2] - 2026-10-29

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.628.1] - 2026-10-29

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.628.0] - 2026-10-29

### Added
- add change-user, change-login, autobot Git identity commands and migrate machine to AI automation bot

---

## [v6.627.0] - 2026-10-28

### Added
- Context Menu Icon Isolation: Resolved duplicate child entry icon writes, ensuring the top-level `gitmap` menu item displays the crisp official application icon while child actions omit redundant collision icons.
- Non-Git Folder Show-Log Guard: When right-clicking inside non-git container folders (e.g. `D:\work`), `gitmap show-log` displays detected repositories and launches the Web UI dashboard instead of halting with an error.
- Default Browser Decoupling: Browser windows now support clean launching in default system browsers without forcing isolated Chrome profile data directories.
- Profile Tree, Single Tool Installer & Macro UI: Added context menu entries for Macro Builder UI and Profile installer, with live Profile Tree inspection and single-tool installer actions in Fleet Web UI.

---

## [v6.626.0] - 2026-10-28

### Added
- Routine release v6.626.0

---

## [v6.625.0] - 2026-10-28

### Added
- Routine release v6.625.0

---

## [v6.624.0] - 2026-10-28

### Added
- Routine release v6.624.0

---

## [v6.623.0] - 2026-10-28

### Added
- Routine release v6.623.0

---

## [v6.622.0] - 2026-10-28

### Added
- Routine release v6.622.0

---

## [v6.621.2] - 2026-10-28

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.621.1] - 2026-10-27

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.619.0] - 2026-10-27

### Added
- strip commit pseudo-module and preserve subsystem casing

---

## [v6.618.0] - 2026-10-26

### Added
- commit author information and interactive repository aum search ui

---

## [v6.617.0] - 2026-10-26

### Added
- installation right-click command open antigravity app-window browser and ui animations

---

## [v6.616.0] - 2026-10-26

### Added
- modern responsive show log single diff exit button context submenus and help

---

## [v6.615.0] - 2026-10-25

### Added
- file context menu powershell cleaner cascade flattening and ui help

---

## [v6.614.0] - 2026-10-25

### Added
- context menu show log ui and update downgrade guard

---

## [v6.613.0] - 2026-10-25

### Added
- Hardened flexible path parsing, URL decoding, and universal CLI path expansion

---

## [v6.612.0] - 2026-10-25

### Added
- Flexible path parser `CleanFlexiblePath` supporting bracketed Windows paths, markdown embeds with file URIs, numbered list prefixes, and trailing shell/code comments.
- Safe terminal image cat in `gitmap cat` rendering structured TokyoNight metadata cards (format, dimensions, size, disk verification status) instead of binary noise.
- Automated resolution for repository-relative paths and parent `../repo-cache/` cache trees.
- Documentation and concrete CLI runbooks in `.agents/skills/gitmap/skill.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`, and coding guidelines.

---

## [v6.611.0] - 2026-10-25

### Added
- Heatmap status overhaul, safe stashing, and AI confusion database

---

## [v6.610.0] - 2026-10-25

### Added
- enhance stash pop recovery and purge invalid clone fallback entries

---

## [v6.609.0] - 2026-10-24

### Added
- Plan 325 enum constants refactor and CI green release

---

## [v6.608.0] - 2026-10-24

### Added
- fix heatmap ModelName and ToolName column migration

---

## [v6.607.1] - 2026-10-24

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.607.0] - 2026-10-24

### Added
- **Incoming Pulled Changes Diffstat**: Automatically renders the detailed `git diff --stat` of all incoming commits integrated by `gitmap pull`, listing modified files and line deltas.
- **Git Status on Dirty Working Tree**: Renders standard Git short status (`git status --short`) with ANSI color-coding for dirty working trees, displaying staged, modified, deleted, and untracked files with full fidelity.
- **Repository Path First Tree Header**: Formats the initial command execution header with the exact filesystem repository path (`-> <repoPath>`), eliminating redundant command echoing.

---

## [v6.606.0] - 2026-10-24

### Fixed
- **Single-Repo Pull Isolation (`gitmap pull`)**: Fixed `gitmap pull` without arguments inside a repository directory (`cwd`) to pull exclusively the active repository instead of defaulting to `pull-all` across all cataloged workspaces.
- **Dynamic Invocation Target Scope**: Enhanced `printPullInvocationHeader` to dynamically resolve target scope `(repo-name)` from working directory or arguments, reserving `(all workspaces)` strictly for explicit `gitmap pull-all`, `pa`, or non-git fallback directory contexts.
- **Direct Filesystem Path Resolution**: Enhanced target slug resolver to automatically detect and pull direct filesystem directory paths matching existing Git repositories without requiring prior SQLite catalog enrollment.

---

## [v6.605.0] - 2026-10-24

### Added
- Native toolchains (cargo, rust, kill), uppercase conventional commits, and error management alignment

---

## [v6.604.0] - 2026-10-24

### Added
- Folder-specific work directory and SSH automation with bracketed commit formatter

---

## [v6.603.0] - 2026-10-24

### Added
- feat([workgroup]): folder-specific workdirs and ssh automation

---

## [v6.602.0] - 2026-10-24

### Added
- use git saved credentials instead of hardcoded fallback in conflict recovery

---

## [v6.601.0] - 2026-10-24

### Added
- fix(pipeline): resolve gosec G115 integer overflow and enhance unmerged conflict auto-recovery

---

## [v6.600.0] - 2026-10-24

### Added
- multi-machine gitmap pa interactive remediation and fleet command dispatch

---

## [v6.599.0] - 2026-10-23

### Added
- plans and memory consolidation: 98 completed micro-plans synthesized into 32 authoritative milestones (01-32)
- subtasks reduction: folded 67 completed subtask directories, preserving 10 active tracks in .ai-memory/plans/subtasks/
- redundant file elimination: pruned all dummy .gitkeep placeholders and 49 historical root plan files
- registries alignment: synchronized .ai-memory/plans/readme.md, 01-index.md, and index.md

---

## [v6.598.0] - 2026-10-23

### Added
- fix installer scripts and CI pipeline compilation

---

## [v6.597.0] - 2026-10-23

### Added
- fix installer scripts and CI legacy refs check

---

## [v6.596.0] - 2026-10-23

### Added
- installer: sync install scripts, suppress 404 discovery warnings, and align release deploy target

---

## [v6.595.0] - 2026-10-23

### Added
- fix legacy refs in tests, gofmt trailing newline, and flatten nested ifs in pipeline logs

---

## [v6.594.0] - 2026-10-23

### Added
- cmdupdate: exclude cli-updater from release asset resolution and fix self-update

---

## [v6.593.0] - 2026-10-23

### Added
- commit tracking, split-db commit history, universal sha resolution, and ai provenance integration

---

## [v6.592.1] - 2026-10-23

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.592.0] - 2026-10-23

### Added
- heatmap discovery, url cache split-db, 10-split downloader, and llm train integration

---

## [v6.591.0] - 2026-10-23

### Added
- Default pull to all workspaces and accelerate SSH fleet remote execution

---

## [v6.590.0] - 2026-10-22

### Added
- Completed the AppError migration: bare fmt.Errorf converted to typed AppError returns across all production packages (only the diag package keeps plain errors, blocked by an import cycle)
- Coding-guideline execute suite run to completion: all 39 prompts applied, partially applied with documented remainders, or closed with evidence
- Deduplicated helper libraries: filesystem and string helpers consolidated into shared packages, 51 local copies removed
- Map-based dispatch and case-insensitive matching helpers applied across command dispatchers; per-call regex compilation eliminated
- Enum guideline compliance verified by linter (9 violations fixed); heavy test suite tagged and wired into CI as a dedicated job
- Repaired the build: deduplicated constants and store types, aligned Result-envelope callers and test seams; full build and vet green
- Plan backlog drained: 24 completed plans filed away; boolean-lint residue reduced to a documented guard idiom

---

## [v6.589.10] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.9] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.8] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.7] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.6] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.5] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.4] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.3] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.2] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.1] - 2026-10-22

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.589.0] - 2026-10-21

### Added
- prompt to send commits on current machine and then fleet send all

---

## [v6.588.0] - 2026-10-21

### Added
- show dry run before dev-clean prompt and add cargo build caches

---

## [v6.587.0] - 2026-10-21

### Added
- fix legacy repo ref in spec and remove unused render function

---

## [v6.586.0] - 2026-10-21

### Added
- fix tree view structure for pa view and ssh fleet

---

## [v6.585.0] - 2026-10-21

### Added
- prompt to send commits on current machine and send all for dirty pull repositories

---

## [v6.584.1] - 2026-10-21

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.584.0] - 2026-10-20

### Added
- native mail client folder watcher and schedule email automation

---

## [v6.583.0] - 2026-10-20

### Added
- native batch screenshot prnt download and extraction engine with brain sync

---

## [v6.582.0] - 2026-10-20

### Added
- database rest command, swagger, custom ui, and record management

---

## [v6.581.1] - 2026-10-19

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.581.0] - 2026-10-19

### Added
- native sqlite database integration and exploration suite

---

## [v6.580.0] - 2026-10-19

### Added
- native sqlite database exploration and management commands

---

## [v6.579.3] - 2026-10-19

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.579.2] - 2026-10-19

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.579.1] - 2026-10-19

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.579.0] - 2026-10-19

### Added
- **TortoiseGit 2.4.0.2 Installer**:
  - Added support for installing TortoiseGit 2.4.0.2 via Chocolatey (`choco install tortoisegit --version 2.4.0.2 -y`) with automatic fallback to standalone offline MSI extraction from Chocolatey `.nupkg`.
  - Added CLI tool aliases: `tortoisegit`, `tortoise-git`, `tgit`, and `tortoise`.
  - Added custom uninstall handler and registry detection probes.
  - Added unit test coverage for package manager resolution in `cmdinstall`.
  - Updated tool installer documentation in `cli/helpdoc/install.md`.

---

## [v6.578.0] - 2026-10-19

### Added
- **GitMap Background Daemon (`gitmap daemon` / `gitmap gitmap-daemon`)**:
  - Renamed `lemon` background service to `gitmap-daemon` with full backward compatibility for `lemon` and `lm` CLI aliases.
  - Added configurable polling loop (default 8s) with a non-negotiable `>= 5s` safety guard.
  - Implemented application settings toggle (`daemon.enabled`): CLI shortcuts `gitmap daemon enable` and `gitmap daemon disable`, with settings command `gitmap settings set daemon.enabled <true|false>`.
  - Added OS background service manager integration for Windows (`sc.exe`), Linux (`systemd`), and macOS (`launchd`).
- **Multi-Format Email Command Ingestion & Immediate ACK**:
  - Support for Format A (plain slash commands), Format B (JSON command envelopes), and Format C (fenced markdown and YAML headers).
  - Immediate pre-execution receipt acknowledgment (`ACK Received: Command '<cmd>' received at <timestamp> (Request ID: <id>). Executing...`).
  - Command whitelist security governance with parameter sanitization, shell injection defense, and execution outcome reporting.
- **Ent Go ORM Feasibility Analysis & Type Safety Hardening**:
  - Comprehensive feasibility study documenting empirical benchmarks (point query latency 12µs vs 24µs, memory allocations 8 vs 32 allocs) and architectural blockers for SQLite Split-DB architecture.
  - Verified compile-time type-safe parameterized SQL across all database packages with zero dynamic SQL concatenation.
- **Codebase Hygiene & File Size Compliance**:
  - Decomposed oversized settings and command files (`cmdagy`, `cmdlemon`, `cmdnodes`, `cmdsupabase`) ensuring all modified files comply strictly with the repository 300-line cap.
  - Flattened all nested `if` statements to maintain 0 violations across AST and boolean linters.

---

## [v6.577.0] - 2026-10-19

### Added
- node-prompt-fixes - fleet telemetry, prompt tree modal, telegram priority, and supabase error backoff

---

## [v6.576.0] - 2026-10-18

### Added
- ENT goal feasibility report, type-safe SQL hardening, and Lemon email command daemon service

---

## [v6.575.2] - 2026-10-18

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.575.1] - 2026-10-18

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.575.0] - 2026-10-18

### Added
- goemailsender - connection fallback delegation and email ui

---

## [v6.574.0] - 2026-10-18

### Fixed
- **GitHub Auth Test Mock Parity (`cli/usercontext`)**: Updated `TestDetectGhAuth_FallbackToken` in `gh_auth_test.go` to handle `git config --global user.name` arguments, fixing cross-platform test failures across Windows, macOS, and Ubuntu.
- **Error Management Compliance (`cli/cmddb`, `cli/heatmap`)**: Resolved 6 swallowed database call errors in `cmddb_tables.go`, `heatmap_conflict.go`, and `heatmap_view.go`. Handled all query and execution errors.
- **Legacy Version References Purged (`02-spec/01-spec-authoring-guide`)**: Replaced `gitmap-v6` references in `13-root-readme-conventions.md`, restoring 100% passing status for `check-legacy-refs.py` and `test_ci_scripts.py`. <!-- gitmap-legacy-ref-allow -->
- **Codebase Go Formatting Cleanliness (`cli/`)**: Formatted all Go files with `gofmt -w` (`rootutility_update.go`, `clone_auth_resolve.go`, `clone_auth_types.go`, `update_ui.go`, `heatmap_query.go`, `heatmap_view.go`, `heatmap_rc_store.go`).

---

## [v6.573.0] - 2026-10-18

### Fixed
- **Dead Code Elimination (`cli/cmddb`)**: Removed uncalled legacy `cmddb_help.go` functions, eliminating golangci-lint `[unused]` linter failures.
- **Static Analysis Nil Pointer Dereference Guard (`cli/cmdupdate`)**: Added immediate nil candidate guard in `finalizeUpdate` (`update_binary_swap.go`), resolving staticcheck `SA5011` possible nil pointer dereference.
- **Async Poller Interval Clamping Test Alignment (`cli/cmd`)**: Aligned `TestParseAsyncArgs` in `async_cmd_test.go` with the system-wide `>= 30s` poller floor invariant, resolving race-detector test failures.

---

## [v6.572.1] - 2026-10-18

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.572.0] - 2026-10-18

### Added
- **Single-Repo Pull UX Simplification (`cli/cmdpull`)**: Replaced 0%-100% counting/compressing lines with a clean `↓ pulling remote changes...` status and compact 1-line summary showing commit delta, diffstat, clean/dirty state, and auto-merge notices.
- **Progress Bar Column Width & Wrap Fix (`cli/cmdpull`)**: Integrated `go-runewidth` to compute exact terminal cell widths for wide Unicode characters and emojis (`⚡`, `🔍`), added `termWidth - 4` safety right-margin to prevent ConHost/Terminal auto-wrap to row Y+1, and throttled non-TTY progress output.
- **Shell Wrapping Forensic RCA & Native Commands (`cli/cmdgit`)**: Documented the 5-part RCA and established strict bans on wrapping tools inside `gitmap ps`/`pwsh`. Added native `gitmap show <rev>` and `gitmap log [-n <N>]` commands to close the git inspection tooling gap.
- **Antigravity Manager (AGM) Subsystem Alignment**: Verified and documented recent AGM improvements (E9001 elimination, missing workspace fail-fast, 30s floor, JWT bridge) and aligned GitMap AGY prompt rerun/review commands.
- **Standalone Image URL Architecture (`cli/cmdagy`)**: Authored `03-image-url-architecture.md` and implemented `ExtractAndStoreBase64Images` in `agy_image_assets.go` converting inline base64 data URIs into content-addressable SHA-256 disk files and lightweight markdown URLs.
- **RC Agent & Heatmap RC Continuous Sync**: Verified the 6-stage non-blocking reconciliation protocol ensuring collision-free SQLite telemetry syncing across AI namespaces.

---

## [v6.571.2] - 2026-10-18

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.571.1] - 2026-10-17

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.571.0] - 2026-10-18

### Minor version bump
- Routine minor bump after CI-green verification: all recent pipeline runs passing at release time.
- Re-syncs the canonical version: the prior release tag v6.570.2 shipped with `version.json` still pinned at 6.570.1 (version/tag sync slip); this release aligns `version.json`, manifests, constants, and docs on 6.571.0.

---

## [v6.570.1] - 2026-10-17

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.570.0] - 2026-10-18

### Fixed
- **Direct SSH Cloning:** Bypassed PAT probe and auth prompts entirely for explicit SSH clone URLs (`git@...`, `ssh://...`) and `--ssh` flag across `clone`, `cfr`, and `clonemulti`.
- **GCM Interactive Modal Suppression:** Silenced Windows Credential Manager / GitHub Desktop sign-in popups during GitMap CLI execution via `GCM_INTERACTIVE=never`.
- **Global Git Identity Reporting:** Enforced strict fallback to `git config --global user.name` and `user.email`, eliminating `"default"` and `"github-user"` placeholders.

### Added
- **Unified Login Subsystem (`gitmap login`):** Added interactive authentication menu and commands for SSH key generation/reuse (`gitmap login ssh [alias] [-y]`), automated clipboard copying, GitLab/GitHub settings links, per-repo `core.sshCommand` binding, `login browser-flow`, and `login token`.

---

## [v6.569.0] - 2026-10-17

### Added
- pipeline-pe - parallelize log downloads, fast ANSI bypass, and SQLite PipelineShaCache speedup

---

## [v6.568.0] - 2026-10-17

### Added
- pull-efficient - fix pae post-pull hang and heartbeat self-termination

---

## [v6.567.0] - 2026-10-17

### Added
- agent-heatmap - db introspection subtask resolution and heatmap tracking

---

## [v6.566.4] - 2026-10-17

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.566.3] - 2026-10-17

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.566.2] - 2026-10-17

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.566.1] - 2026-10-17

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.566.0] - 2026-10-17

### Added
- hd-ui - improve help dashboard themes prompt finder and terminal view

---

## [v6.565.0] - 2026-10-17

### Added
- pull-summary transparency, remediation menu guidance, safe-pull conflict hardening, and alim-seo LFS fix

---

## [v6.564.0] - 2026-10-17

### Added
- merge conflict resolution and pipeline verification

---

## [v6.563.0] - 2026-10-17

### Added
- large-file-split, git-lfs-convert, repo-flatten, reclone-url-fix, and automerge-ide-guidance

---

## [v6.562.0] - 2026-10-17

### Added
- **Large Files Splitting Engine (`gitmap split large files`):** Discovers tracked and untracked files exceeding 10MB (configurable via `--size`), creates multi-part zip chunks (`<file>.zip`, `<file>.z01`, `<file>.z02`), applies STORE heuristic for pre-compressed media (`.mp3`, `.mp4`, `.png`, `.jpg`, `.wav`) and DEFLATE for data/code, removes raw files from Git history, and prompts with bold yellow irreversible warnings.
- **Git LFS Conversion Suite (`gitmap lfs convert`):** Automated scanner finding files exceeding 5MB, updates `.gitattributes` tracking patterns, migrates Git history to LFS pointers via `git lfs migrate import` with fallback staging, and previews changes via `--dry-run`.
- **Repository Flattening Engine (`gitmap repo flatten`):** Compresses bloated multi-gigabyte repositories into a clean successor (`<repo>-flat-v1`) with a single consolidated squashed snapshot commit, archives the previous repository to `<repo>-marked-for-deletion`, swaps to the canonical name, and provides GitMap `undo` restoration.
- **Flatten Undo Manager (`gitmap repo undo-flatten`):** Safely restores the previous repository archive from `<repo>-marked-for-deletion`.
- **Muse Training Priority Sync:** Updated `cli/cmd/llm/llm_train_muse.go` to prioritize `repo-secrets` checkouts for active LLM training, restored missing `decisionLog` struct/methods, and flattened nested if statements.
- **CLI Help & Troubleshooting Docs:** Authored `cli/helpdoc/split-large-files.md`, `cli/helpdoc/lfs-convert.md`, `cli/helpdoc/flatten.md`, added how-to recipes in `cli/cmdhelp/howto.go`, and added a troubleshooting issue card for 'Large Files / Git Clone Failure (No space left on device)' in `src/data/troubleshootingIssues.ts`.

---

## [v6.561.0] - 2026-10-17

### Added
- **Histogram Diff Automerge Optimization:** Upgraded `cli/cloner/safe_pull.go` to invoke automerge with `-X diff-algorithm=histogram` and `--no-edit`, preventing false merge conflicts caused by matching indentation and single-bracket lines across diverged branches.
- **Git Rerere Auto-Memory Engine:** Implemented `cli/gitutil/git_merge.go` with `EnableGitRerere` and `CheckGitRerereEnabled`, allowing Git to automatically record and re-apply conflict resolutions across past and future pulls.
- **Structured Conflicted File Extractors:** Added `ExtractConflictedFilesFromOutput` and `ExtractConflictedFilesFromRepo` in `cli/cloner/safe_pull.go`, parsing both git conflict stdout text and porcelain status codes (`UU`, `AA`, `UD`, `DU`, `UA`, `AU`, `DD`) into clean file lists.
- **Multi-IDE Launcher & Guidance:** Implemented `cli/cmdconflict/ide_launcher.go` detecting `code` (VS Code), `cursor`, `windsurf`, and `codium`. Formatted a step-by-step 3-Way Merge Editor walkthrough (`Ctrl+Shift+G` -> 'Resolve in Merge Editor' -> Accept Incoming/Current -> Complete Merge).
- **GitMap Conflict Resolution CLI:** Created `gitmap resolve <repo> [--theirs | --ours | --ide | --abort | --rerere]` and `gitmap conflicts [repo]` in `cli/cmdconflict/resolve_cmd.go`, with root command dispatch integration in `cli/cmd/rootutility.go`.
- **Expanded `gitmap fix <repo>` Actions:** Added support for `theirs` (`t`), `ours` (`o`), `resolve` (`v`, `ide`, `code`), `abort` (`a`), and `rerere` in `cli/cmdfix/fix_cmd.go`.
- **Post-Pull Merge Conflict UI Card:** Upgraded `cli/cmdpull/pull_prompt.go` to render a styled Lipgloss conflict box whenever a repository encounters conflicts during pull, providing one-key actions for IDE launch (`[v]`), remote incoming (`[t]`), local current (`[o]`), and safe abort (`[a]`).
- **Comprehensive Help & Troubleshooting Docs:** Authored `cli/helpdoc/resolve.md`, `cli/helpdoc/conflicts.md`, updated `cli/cmdhelp/howto.go`, and registered merge conflict troubleshooting in `src/data/commands.ts` and `src/data/troubleshootingIssues.ts`.

---

## [v6.560.0] - 2026-10-17

### Added
- **Parallel Fix All (`gitmap fix --all`):** Generic typed parallel worker pool (`cli/parallel/generic_pool.go`) dynamically scaling across CPU cores (up to 20 cores) with real-time per-repository progress feedback.
- **Git Repository Reclone by URL:** `gitmap reclone <url>` now directly detects repository URLs (HTTPS, SSH, git) and safely removes existing dirty targets before re-cloning without falling into manifest file parsers.
- **Redesigned Post-Pull Menu:** Reordered interactive post-pull resolution choices to prioritize `[1] Fix All (Run gitmap fix --all)`, followed by `[2] Commit All (WIP locally)`, `[3] Stash All`, `[4] Step through each repository individually`, `[5] Fleet Control`, and `[6] Skip`.
- **Repo-Cache Multi-AI Repository:** Initialized private `repo-cache` repository on GitHub and locally; migrated Muse 7-day memory digests (`12-muse-ai/`) and updated `gitmap llm train muse` to resolve `repo-cache` first.
- **Subcommand Hyphenation:** Automatic hyphen-to-space command routing across two-word verbs (`train-muse`, `llm-train-muse`, `repo-export`, `repo-import`, etc.).
- **Heatmap Skill Documentation:** Authoritative documentation of `gitmap heatmap`, `gitmap agent heatmap`, and `full-summary` in `.agents/skills/gitmap/SKILL.md` and `coding-guidelines`.

### Fixed
- **Windows GCM Credential Store (`wincredman`):** Fixed fatal GCM error `Unknown credential store 'wincred'` by setting `wincredman` across `fixcredential.go`, `pull_remediation.go`, and `git_env.go`.
- **Pull All Concurrency:** Scaled `resolveAutoProfileConcurrency` up to `runtime.NumCPU()` (20 cores) while preserving progress bar rendering and terminal ANSI color escape codes.
- **Update Probe Throttle & UI Polish:** Throttled sibling version probing in `gitmap update` to once every 24 hours via watermark cache (`probe-watermark.json`), with Title Case styling, left padding, and bold yellow transitions.

---

## [v6.559.0] - 2026-10-17

### Added
- Routine release v6.559.0

---

## [Unreleased]

### Changed
- `gitmap llm train muse`: the last-read watermark (vault HEAD SHA + timestamp) is now tracked by gitmap itself in its local state — the `12-muse-ai/last-read.md` checkpoint file is gone. Every run prints a numbered DECISION TRACE before the summary (why the vault was reused vs cloned, why pull is `--ff-only`, what is new since the last read, why the N-day window and file order), a SINCE LAST READ section (new commits since the watermark, or "unchanged"), and each AI traversal checklist step carries its rationale — so an AI learns the decision-making behind the gitmap steps, not just the commands.

### Added
- `gitmap llm train muse [--days N]`: positional sub-target of `llm train` that syncs the Muse 7-day memory digest — clone-or-pull the repo-cache vault via `gh repo clone` (fail-fast `gh auth` guidance, never prompts; an existing vault dir at the work-folder root is reused even without an `rc` record), `git pull --ff-only`, N-day `git log` diff, terminal summary (commits, files changed, sequenced digests) and a raw AI traversal checklist (`gitmap cat` / `gitmap aum search` hints); tolerates a missing `12-muse-ai/` digest with guidance, exit 0. The digest moved from repo-secrets to the dedicated `alimtvnetwork/repo-cache` vault ("GitMap distributed shared cache and AI memory vault"); repo-secrets keeps secrets only.
- `--days N` flag on `gitmap llm train` (default 7): digest window in days for `train muse`.

---

## [v6.558.0] - 2026-10-16

### Added
- Routine release v6.558.0

---

## [6.557.0] - 2026-10-16

### Added
- **Pure Binary Self-Updater (`gitmap update`):** Inverted `gitmap update` strategy so pre-compiled binary download and in-place handoff is the primary default path, permanently eliminating the automatic fallback to `run.ps1 -Update`, `git pull`, and `go build`. Pre-compiled release binaries are downloaded directly from GitHub Releases with SHA-256 integrity verification and console-safe Windows process handoff (`gitmap-update-<pid>.exe`).
- **Terminal UI Modernization:** Modernized `install.ps1` and `install.sh` with 24-bit Truecolor Dracula ANSI palettes (`#8be9fd`, `#bd93f9`, `#50fa7b`, `#ffb86c`, `#ff5555`, `#6272a4`), rounded box-drawing cards (`╭`, `╮`, `╰`, `╯`, `│`, `─`), 5 staged execution steps, interactive Braille spinners (`⠋` to `⠏`), and completion status cards.
- **Offline Local Installers:** Created standalone `local-install.ps1` and `local-install.sh` for offline installations capable of deploying directly from pre-built `./bin/` artifacts or compiling from local `./cli` source with identical Truecolor terminal UI.
- **CLI Update UI:** Upgraded `gitmap update` with rich Lipgloss-style terminal UI, multi-stage progress indicators, and fast-path notifications.
- **Installer Download URL Override:** Added `$GITMAP_DOWNLOAD_URL` and `$env:GITMAP_DOWNLOAD_URL` override support across `install.ps1` and `install.sh` for isolated testing and enterprise staging environments.

### Fixed
- **Bug Fix**: Inverted `gitmap update` strategy so pre-compiled binary download and in-place handoff is the primary default path, permanently eliminating the automatic fallback to `run.ps1 -Update`, `git pull`, and `go build`.
- **PowerShell ANSI PlainText Suppression:** Removed `$PSStyle.OutputRendering = 'PlainText'` in `install.ps1` and enabled Windows Virtual Terminal Processing (`ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004`) via Win32 `SetConsoleMode`.
- **Source Build Priority Inversion:** Resolved updater bug where `gitmap update` erroneously prioritized local repository rebuilds over pre-compiled release binaries.

---

## [v6.556.1] - 2026-10-16

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.556.0] - 2026-10-16

### Added
- `gitmap os --json` / `gitmap os json`: Fixed subcommand dispatch to output structured `MachineIdentity` JSON without throwing `E_INVALID_OS_SUBCMD`.
- `gitmap os`: Enhanced 3-part layout displaying Modern OS menu, Active Host Machine Identity, and sample real-world commands at the bottom.
- `gitmap cpu`, `processor`, `ram`, `disk`: Instant hardware metrics across Windows (PowerShell CIM average query) and Linux (`/proc/stat`, `/proc/meminfo`, `syscall.Statfs`) with 5-second cross-process caching (`gitmap_hw_metrics_cache.json`) to eliminate CPU spikes.
- `gitmap os network` / `gitmap network`: Network adapter discovery, MAC addresses, IPv4/IPv6 addresses, link status (`isUp`), link speed, and byte/packet counters with `--json` support.
- `gitmap nodes tasks`: Distributed task queue interrogation across cluster nodes, reporting CPU, RAM, active/pending task counts, heavy-duty task classification (`isHeavyDuty`), and role-based scheduling placeholder (`?`).
- Fixed node enumeration in `nodes_resources.go` by integrating unified node discovery (`collectUnifiedNodes`).
- Terminal help documentation (`cli/helpdoc/` for `cpu`, `ram`, `disk`, `network`, `nodes-tasks`) and command palette definitions in `src/data/commands.ts`.

---

## [v6.555.0] - 2026-10-16

### Added
- `gitmap nodes-running-only` (alias `gitmap nro`, `nodes-running`, `nodes-online`, `nodes-open`, `nodes-active`): Discovers and displays strictly active, online fleet machines (`is_online: true`), filtering out offline/timed-out machines.
- `gitmap nro --json` / `gitmap nodes-running-only --json`: Machine-readable JSON array containing strictly verified reachable nodes.
- 10-Second Node Liveness Caching: In-memory and cross-process JSON cache (`gitmap_nodes_liveness_cache.json`) in OS temp directory with 10s TTL, eliminating repeated 1200ms TCP probe timeouts during sequential CLI loops.
- Cache Invalidation Flags: Added `--refresh` (`-r`) and `--no-cache` to force immediate live reachability probing.
- Running-Only Delegation: Seamless forwarding for `nro exec "<cmd>"`, `nro tasks`, `nro resources`, `nro ping`, `nro pull`, `nro cp all`, and `nro scan`, executing solely on active nodes without offline delay or warnings.
- Pre-Dispatch Offline Filtering: Updated `gitmap nodes tasks` and `gitmap nodes resources` to accept `--running-only` / `-ro` / `--open-only` flags, skipping unreachable nodes before spawning goroutines.
- Help & UI Documentation: Authored `cli/helpdoc/nodes-running-only.md`, registered completion in `cli/completion/completion.go`, and added Web UI catalog entry in `src/data/commands.ts`.

---

## [v6.554.0] - 2026-10-16

### Added
- `gitmap os network` / `gitmap network`: Cross-platform network adapter telemetry, status, MAC, IPv4/IPv6, link speed, and byte traffic with `--json` output.
- `gitmap cpu` / `gitmap processor`: Real-time processor utilization and specs with 5-second cross-process cache protection.
- `gitmap ram`: Real-time physical memory and swap utilization metrics with 5-second cache protection.
- `gitmap nodes tasks` / `gitmap nodes queue`: Fleet-wide AI agent task queue inspection from `agent-task.db` across cluster nodes.
- `gitmap os`: Added active host identity summary footer to default bare OS command.

### Fixed
- `cli/cmdos/os.go`: Resolved `E_INVALID_OS_SUBCMD` regression on `gitmap os --json` and `gitmap os -j`.
- `cli/cmdnodes/nodes_resources.go`: Replaced defective `cmdssh.FetchAllSSHConnections()` with `collectUnifiedNodes(ctx, true)` to properly include local node and cluster database nodes.

### Documented
- Added `cli/helpdoc/cpu.md` and `cli/helpdoc/ram.md`.
- Updated `cli/helpdoc/os.md` and `cli/helpdoc/nodes.md`.
- Registered `network`, `cpu`, `ram`, and `nodes tasks` in `src/data/commands.ts`.

---

## [v6.553.0] - 2026-10-15

### Added
- Routine release v6.553.0

---

## [v6.552.0] - 2026-10-14

### Added
- Routine release v6.552.0

---

## [v6.551.0] - 2026-10-14

### Added
- Routine release v6.551.0

---

## [v6.550.0] - 2026-10-14

### Added
- Routine release v6.550.0

---

## [v6.549.0] - 2026-10-12

### Added
- Routine release v6.549.0

---

## [v6.548.0] - 2026-10-12

### Added
- update-antigravity-ide - Antigravity IDE Update Integration

---

## [v6.547.0] - 2026-10-12

### Added
- feat(release): implement release-remove-between (rrb) with auto dry-run and pin-not-to-remove (pnr)

---

## [v6.546.0] - 2026-10-11

### Added
- ci - decompose root_cobra_completion and completion/install for 300-line size gate

---

## [v6.545.0] - 2026-10-11

### Added
- ci - format schemaregistry_assert_test.go with gofmt to pass lint-script tests

---

## [v6.544.0] - 2026-10-11

### Added
- ci - remove unused schemaregistry assert helpers for lint-baseline-diff

---

## [v6.543.0] - 2026-10-11

### Added
- ci - remove unused legacy statusformat, statusmissing, and statusprint for lint-baseline-diff

---

## [v6.542.0] - 2026-10-11

### Added
- decompose pipeline_db_ops into pipeline_db_ops and pipeline_db_purge for file size gate

---

## [v6.541.0] - 2026-10-11

### Added
- ci - remove unused mv_relocate.go for lint baseline diff

---

## [v6.540.0] - 2026-10-11

### Added
- sync companion repo-cache gitmap json on merge-ai and clone

---

## [v6.539.0] - 2026-10-11

### Added
- ci - remove unused pe all cache and error symbols for lint baseline guard

---

## [v6.538.0] - 2026-10-11

### Added
- ci - fix redundant newline in fmt.Print within rootutility_update_args.go

---

## [v6.537.0] - 2026-10-11

### Added
- ci - decompose rootutility into rootutility_update and rootutility_update_args for 300-line size gate

---

## [v6.536.0] - 2026-10-11

### Added
- ci - fix unused import cmdsummary in rootutility breaking release binaries

---

## [v6.535.0] - 2026-10-11

### Added
- repo-management - repo cache json injection subtree solutions and native path acceleration

---

## [v6.534.0] - 2026-10-11

### Added
- reseed lint baseline cache to v4 and make merge-ai test hermetic

---

## [v6.533.0] - 2026-10-11

### Added
- ci lint baseline cache v3, unused symbols cleanup, and helper functions refactor

---

## [v6.532.0] - 2026-10-11

### Added
- fix file size gate 300-line cap and unused symbols

---

## [v6.531.0] - 2026-10-10

### Added
- GitMap Repo-Cache Scan, Clone, and LS Ecosystem:
  - `gitmap scan . --rc` / `gitmap scan --rc`: Export and auto-merge scanned repositories into `repo-cache/01-gitmap/gitmap.json` with URL normalization and deduplication.
  - `gitmap scan . --rc --separate` (aliases `--seprate`, `--sep`, `--s`): Auto-allocate next sequential manifest (`02-gitmap.json`, `03-gitmap.json`).
  - `gitmap clone/cfr/cfrp --rc` / `rc`: Split-DB cached (<5ms) manifest discovery and interactive/auto (`-y`) cloning with concise `owner/repo` TreeView.
  - `gitmap ls rc` / `gitmap list rc`: Manifest inspector with `[SSH]` and `[Public HTTPS]` protocol badges and copy-pasteable import hints.

---

## [v6.530.0] - 2026-10-10

### Added
- Implement gitmap fix release tags command to audit and purge orphan and broken release tags

---

## [v6.529.0] - 2026-10-10

### Added
- Pull failure subtree yellow solutions, non-repo smart clone and unified batch fix

---

## [v6.528.0] - 2026-10-10

### Added
- Tree export flags, tree-search engine family, pending-commits backup DB and relative paths

---

## [v6.527.0] - 2026-10-10

### Added
- cd  default workdir autonomous resolution and persistence

---

## [v6.526.0] - 2026-10-10

### Added
- pending commits table improvement, status cache backup db and new commands

---

## [v6.525.0] - 2026-10-10

### Added
- summary, full summary treeview, nodes delegation, pe all, merge-ai and universal repo feature

---

## [v6.524.0] - 2026-10-10

### Added
- pending commits table improvement status cache and new commands

---

## [v6.523.1] - 2026-10-09

### Added
- Routine release v6.523.1

---

## [v6.523.0] - 2026-10-09

### Added
- Routine release v6.523.0

---

## [v6.522.0] - 2026-10-09

### Added
- pipeline errors all aggregator plus te shortcut

---

## [v6.521.0] - 2026-10-09

### Added
- Routine release v6.521.0

---

## [v6.520.0] - 2026-10-09

### Added
- DB-backed spec number issuance (gitmap spec next)

---

## [v6.519.0] - 2026-10-09

### Added
- Routine release v6.519.0

---

## [v6.518.0] - 2026-10-09

### Added
- Routine release v6.518.0

---

## [v6.517.0] - 2026-10-09

### Fixed
- `gitmap agm update` terminal UI: coordinated pterm progress UI (spinner + live installer log tail + stage checkmarks) fed by piped child output — no more minutes of silence, no raw installer escape sequences corrupting the screen.
- `gitmap agm update` no longer misrouted to `gitignore`: the `agm` alias collision in the core dispatcher is removed, so the documented command reaches the real updater.
- Failure path: structured error panel (failed stage, last installer lines, next-step hint) instead of raw output dumps + Go stack traces; verification failures point at the diagnostic log file.
- `agm update --help` now lists the `--yes` flag.

### Added
- `GITMAP_AGM_INSTALL_SCRIPT` env override for the AGM installer source (test seam; default behavior unchanged).

---

## [v6.515.0] - 2026-10-08

### Added
- multi-user project switching, github cli auth status, and git config manager

---

## [v6.514.0] - 2026-10-08

### Added
- user - multi-user project switching, github cli auth status, and git config manager

---

## [v6.513.0] - 2026-10-08

### Added
- Concise bare `gitmap` CLI display (< 15 lines) with structured suggestions:
  - Sub-point 1: Catalog exploration via `gitmap help` and `gitmap -h`.
  - Sub-point 2: Mandatory AI model self-training via `gitmap llm train` (and `gitmap ai llm-train`) before starting workspace tasks.
- Elimination of redundant blank line between binary location triplet and version identity block.
- Authoritative `gitmap llm train` subsystem:
  - Emits full copy-pasteable Antigravity skill with YAML frontmatter directly to stdout.
  - Automatically writes and updates `.agents/skills/gitmap/SKILL.md`.
  - Added `--urls` flag emitting public raw GitHub documentation URLs for LLM memory ingestion.
  - Directs autonomous AI models to generate and persist their own skill from training output.
  - Embeds recursive Git network learning directives (multi-repo topology, prior mistake prevention via `gitmap pe history-ai`, dynamic ETA wait, memory ledgers).
- Modernized skills documentation in `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md`.

---

## [v6.512.0] - 2026-10-08

### Added
- suggestion engine, fastgate pre-commit runner, and muse help clusters

---

## [v6.511.0] - 2026-10-08

### Added
- fix pipeline pe -1 nested ifs boolean guidelines relative paths and misspell

---

## [v6.510.0] - 2026-10-08

### Added
- chrome profile import-export e2e verification and test logging in repo-secrets

---

## [v6.509.0] - 2026-10-08

### Added
- cli build verification, cpar documentation, and chrome profile automation

---

## [v6.508.0] - 2026-10-08

### Added
- universal file runner, supabase vault, cpc chore commits, pe-ai and muse skills v6.508.0

---

## [v6.507.2] - 2026-10-08

### Added
- sync - add native fleet sync and sqlite agent task engine

---

## [v6.507.1] - 2026-10-07

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.507.0] - 2026-10-07

### Added
- pipeline pe cache invalidation relative path sanitization and error extraction fix v6.507.0

---

## [v6.506.2] - 2026-10-07

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.506.1] - 2026-10-07

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.506.0] - 2026-10-07

### Added
- **Cursor Settings Guard Validation**: Added `requireCursorSettingsFile` in `cli/cmdcursor/cursor_settings.go` with descriptive `E1037` error envelopes prompting users to run `gitmap cursor settings apply`.
- **Views AI Settings UI Modernization**: Refactored `src/pages/Settings.tsx` to invert `/api/instances` fetch validation into clean guard clauses, eliminating nested `if` statements.
- **Process Discovery Spacing Hygiene**: Enforced vertical blank line spacing rules in `cli/cmdagy/agy_instance_discovery.go` and `cli/cmdcursor/cursor_settings_test.go`.
- **Quality Gates Verification**: Verified zero violations across `check-nested-ifs.py`, `check-enum-and-boolean.py`, and `check-relative-paths.py`.

---

## [v6.505.0] - 2026-10-06

### Added
- history purge dispatch fix, recycle bin export, and relative path sanitation

---

## [v6.504.0] - 2026-10-06

### Added
- `gitmap login`: unified GitHub authentication — `--token <PAT>` (validated via api.github.com), `--web` browser login (gh device flow or guided token page), interactive chooser, `--status`
- `gitmap logout`: removes the stored GitHub credential
- Logged-in token is reused automatically by clone/pull/push via token resolution
- Secure no-echo token prompt; fail-fast guidance when stdin is not a TTY

---

## [v6.503.0] - 2026-10-06

### Added
- chrome profile import export e2e test verification and test suite fixes

---

## [v6.502.0] - 2026-10-06

### Added
- resolve CI/CD sends push assertion and ERD parity cluster path

---

## [v6.501.0] - 2026-10-06

### Added
- Deep spec consolidation, canonical 8-cluster reduction, and memory unification

---

## [v6.500.0] - 2026-10-06

### Added
- Pre-deep-consolidation baseline release

---

## [v6.499.0] - 2026-10-06

### Added
- App spec and completed plans consolidation and reduction

---

## [v6.498.0] - 2026-10-06

### Added
- Pre-consolidation baseline release

---

## [v6.497.0] - 2026-10-06

### Added
- **Ubuntu IDE and GitHub Desktop Scan Sync Suite**: Resolved omission where `gitmap scan` on Ubuntu failed to add repositories into VS Code and GitHub Desktop.
- **Root Cause Analysis (Issue 71)**: Grounded 4-part RCA in `02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md` identifying missing directory creation in VS Code Project Manager extension path and missing Linux CLI resolution candidates for GitHub Desktop.
