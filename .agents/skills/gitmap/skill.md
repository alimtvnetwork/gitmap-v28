---
name: gitmap
description: Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation (AUM), pipeline self-healing, process management, and coding guideline enforcement.
---

# GitMap Autonomous Engineering Skill

> **Invocation:** `@[gitmap](file;.agents/skills/gitmap)`

## Overview
GitMap is an ultra-fast developer companion and autonomous CLI engine designed for AI coding agents and software engineers.

- **Lead Architect & Author:** MD ALIM UL KARIM (alimtvnetwork)
- **Sponsored By:** RISEUP ASIA LLC (https://riseup-asia.com)
- **Core Mission:** High-performance polyglot repository management, zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures, and AI agent pair programming.

---

## Non-Negotiable Command Replacement Matrix (DO NOT USE X; USE GITMAP Y)

All AI agents operating within GitMap-managed repositories MUST strictly adhere to the following command substitution matrix. Raw shell traversal and slow unbuffered tools are strictly banned:

| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr` | ✅ `gitmap aum search "<pattern>" [dir] [--ext <ext>]` or `gitmap search "<pattern>"` | Multi-core streaming live search (<15ms) with binary null-byte probe, 500 KB file cap, and DH2D SQLite hot cache. Prevents buffer bloat and terminal freezing. |
| ❌ `Get-ChildItem` (all variants: bare `Get-ChildItem`, `-Path`, `-Filter`, `-Recurse`, etc.), `dir`, `ls` (in shell), `find . -name "..."`, `dir /s /b` | ✅ `gitmap lf <dir>` or `gitmap list-files [dir]` or `gitmap find "<pattern>" [dir]` | Sub-millisecond directory listing and file glob discovery without shell traversal overhead. TOTAL BAN on PowerShell Get-ChildItem. |
| ❌ Inspecting plans manually via `Get-ChildItem -Path ".ai-memory/plans"` or `dir` | ✅ `gitmap plans [--json\|--yaml] [N]` (alias: `gitmap plan`) | Instant Tokyo Night styled inventory of completed and pending plans with sequence IDs, status pills, and backlog overflow navigation. |
| ❌ Directory File Inventory: `Get-ChildItem -Path <dir>`, `dir <dir>`, `ls <dir>` | ✅ `gitmap lf <dir>` (e.g. `gitmap lf src-tauri/src/modules/repo_db`) or `gitmap find "*.*" <dir>` | Instant directory file listing and inventory without PowerShell object pipeline or shell traversal latency. TOTAL BAN on PowerShell Get-ChildItem. |
| ❌ Raw unbuffered `cat`, `type`, `Get-Content` on source files | ✅ `gitmap cat <filepath> [--raw]` (alias: `gitmap ct`) | Direct zero-disk stream with flexible path parsing (Markdown, URIs, comments, brackets) and safe image inspection cards preventing terminal corruption. |
| ❌ Raw ambient `python script.py` or `python -c "..."` | ✅ `gitmap py "<code-or-script>"` or `gitmap py -c "<expression>"` | Auto-resolves verified Python interpreter from `installation.db` cache (<15ms), avoiding environment discrepancies. |
| ❌ `powershell -Command "..."` or `pwsh -Command "..."` | ✅ `gitmap pwsh "<cmd>"` or `gitmap ps -c "<cmd>"` | Cross-platform PowerShell execution with deterministic `-NoProfile`, UTF-8 encoding, and automatic fallback. |
| ❌ Raw `bash -c "..."` or `sh -c "..."` | ✅ `gitmap bash "<cmd>"` or `gitmap bash -c "<cmd>"` | Uniform POSIX execution with cross-platform environment isolation across Windows, macOS, and Linux. |
| ❌ Multi-repo `git status` loops or manual directory scanning | ✅ `gitmap status --dirty` / `gitmap status --json` / `gitmap st` | Single-shot multi-repository audit matrix displaying dirty counts, ahead/behind branches, and uncommitted stashes. |
| ❌ `gh auth login` interactive prompts or plaintext `.env` files | ✅ `gitmap login --web`, `gitmap login --status`, `gitmap login --token <PAT>` | Token validated via GitHub API before writing; auto-resolved for all clone, pull, and push commands. |
| ❌ Ad-hoc `git config --global user.name` / `user.email` or manual shell identity toggling for bots | ✅ `gitmap login set global ai-auto` or `gitmap login set ai-auto [agent-name]` | Locks standardized AI automation bot credentials (`gitmap automation bot - runs autonomously <info@rgitmap.com>` or named agent persona) machine-wide or per-repo with automatic original user backup. |
| ❌ Manual commit author flags (`git commit --author="..."`) or uncredited subagent commits | ✅ `gitmap cpf/cpb/cpc/cfr/cpr --ai-agent <name>` | Attaches designated subagent identity to commit provenance, Git author metadata, Split-DB telemetry, and cognitive heatmap tables. |
| ❌ Raw `git pull`, `git pull origin <branch>`, `git pull --rebase` | ✅ `gitmap pull [repo]` (or `gitmap pull` in current repo, alias `gitmap p`), `gitmap pull-all` (alias `gitmap pa`) | High-speed multi-phase pull with conflict detection, branch tracking, and auto-stash/rebase handling. TOTAL BAN on raw `git pull` without attempting `gitmap pull` first. |
| ❌ Third brackets, lowercase conventional syntax, colons in arguments, or misusing commit verbs (`gitmap cpf "chore(release): ..."` or `gitmap cpf "[module]: feat: ..."` or `git commit -m "..."`) | ✅ `gitmap cpr "<module> - <summary>"` (release), `gitmap cpc "<module> - <summary>"` (chore), `gitmap cbf "<module> - <summary>"` (bug fix), `gitmap cfr "<module> - <summary>"` (refactor), `gitmap cpf "<module> - <summary>"` (feature) | Strict verb taxonomy: NEVER use `cpf` for releases, chores, or bugs. GitMap automatically prepends the sovereign type. TOTAL BAN on square brackets `[...]`, colons in arguments, raw `git commit`, and conventional prefixes (`chore(...)`, `feat(...)`, `fix(...)`) inside arguments. |
| ❌ Raw `$env:LIBCLANG_PATH` exports or manual Cargo wrapper scripts | ✅ `gitmap cargo [args...]` or `gitmap rust status` | Automated LLVM/libclang path discovery on Windows and Linux, injecting `LIBCLANG_PATH` and `PATH` transparently. |
| ❌ Raw `taskkill /F /PID <pid>` or `kill -9 <pid>` | ✅ `gitmap kill <pid|name|:port> [--force] [--tree]` (alias: `gitmap terminate`) | Cross-platform process termination with safety guards protecting system-critical processes, IDEs, and shells. |
| ❌ Saving temporary scratch or test scripts into repo git tree | ✅ `gitmap rc text "<content>" --slug <slug> --ext .ps1` | Centralized script storage in `repo-cache` (`repo-storage`) for permanent cross-repo reuse without polluting git worktrees. |
| ❌ Committing `.env` or credentials to standard repositories | ✅ `gitmap rs text "<secret>" --slug <slug>` | Strict zero-secrets policy; stores credentials exclusively in `repo-secrets` vault. |
| ❌ `gh run watch` or tight polling loops (`while true; sleep 5`) | ✅ `gitmap pipeline-ai status -t <eta>` or `gitmap pe -t` | Dynamic timeout waiting driven by calculated workflow ETA without burning CPU or Actions API quotas. |
| ❌ Slow Python fleet sync (`python 03-ai-scripts/38-sync-prompts-skills-scripts.py`) | ✅ `gitmap sync [--workers 8] [--projects <path|json>]` | Native Go multi-repo synchronization across 43 repositories in <5s with 6-stage safe ceremony (backup branch, pre-pull, 5 boundaries, atomic commit). |
| ❌ Python SQLite task scripts or ad-hoc queries (`python -c "import sqlite3..."`) | ✅ `gitmap agent <task|subtask|log|diagnose|crashed|temp-clear>` | Native compiled Go SQLite task manager (<1ms) with WAL mode, single-writer locking, atomic in-flight telemetry, and crash autopsy. |
| ❌ Python or ad-hoc SQLite scripts (`python -c "import sqlite3..."`), `sqlite3` CLI, DB browser tools | ✅ `gitmap db-read`, `gitmap db-add`, `gitmap db-edit`, `gitmap db-find`, `gitmap db-listing`, `gitmap db-view`, `gitmap db-rest`, `gitmap db-record-del`, `gitmap db-rm` | Compiled native Go SQLite exploration, mutation, view join, REST API server with Swagger UI & visualizer dashboard, and database lifecycle management. |
| ❌ Ad-hoc Python network scrapers or image link extractors (e.g. `python -c "import urllib.request, re..."`, `fetch_urls.py`, curl pipes to scrape `prnt.sc`/lightshot/screenshot URLs) | ✅ `gitmap prnt <url1> [url2...] [--extract-only] [--json]` (aliases: `gitmap download-prnt`, `gitmap asset download-prnt`) | Native compiled parallel fetching, multi-tier CDN regex resolution, User-Agent emulation, zero external dependencies, and automatic bidirectional Antigravity Brain synchronization. |
| ❌ Raw `curl`, `wget`, `Invoke-WebRequest`, or custom Python/Go HTTP fetchers for web resources | ✅ `gitmap url read <url>` (aliases: `gitmap url <url>`, `gitmap get <url>`) | Dedicated SQLite Split-DB (`repodb/url_cache.db`), WAL mode, sub-millisecond retrieval (<0.5ms), unbuffered `output.Raw()` streaming, SHA-256 caching. TOTAL BAN on raw curl/wget. |
| ❌ Single-stream downloads for files $\ge 1\text{--}2\text{ MB}$ | ✅ GitMap 10-split segmented parallel download engine with `aria2c` fallback | Dynamic $\ge 1\text{--}2\text{ MB}$ threshold; host `aria2c -s 10 -x 10 -k 5K` or native Go 10-goroutine Range worker pool writing via `io.NewOffsetWriter` with 5KB chunking. |
| ❌ Grepping repository filesystem or guessing past agent decisions | ✅ `gitmap heatmap` or `gitmap rc heatmap --repo <alias>` | 3-tier database inventory, alias mappings via `db.ListAliasesWithRepo()`, recent tasks across 2–5 repos, dynamic suggestion pills. TOTAL BAN on filesystem grepping for agent traces. |
| ❌ Modifying shared files blindly without checking multi-agent contention | ✅ `gitmap heatmap discuss` or `gitmap agent heatmap` | Ranks write claims, active leases, and collision blast radius across agents (<5ms). Prevents file lock crashes and concurrent overwrites. |
| ❌ Raw `cargo build`, `cargo test`, `cargo check` | ✅ `gitmap cargo build [--install]`, `gitmap cargo test`, `gitmap crg check` | Automatically enriches environment with discovered toolchain paths (`CARGO_HOME`, MSVC) and offers auto-install if missing. |
| ❌ Hardcoded `$env:LIBCLANG_PATH = "..."` or manual LLVM registry searching | ✅ `gitmap aum locate llvm` + `gitmap env set LIBCLANG_PATH <path>` or `gitmap pwsh "$env:LIBCLANG_PATH=...; gitmap cargo test"` | Dynamically resolves LLVM/Clang bin directories in <15ms across machines without breaking portability. |
| ❌ `taskkill /F /PID <pid>`, `tasklist`, `kill -9`, raw `Stop-Process` | ✅ `gitmap kill <pid|name|:port> [--force] [--tree]` (alias: `gitmap terminate`), `gitmap ports [port]` | Safe cross-platform process termination with PID/port/name validation, caller self-protection guards, and graceful or forced two-stage shutdown. |
| ❌ Manual pwsh "Move-Item ...; Copy-Item ...", taskkill + copy on locked running binaries | ✅ `gitmap swap-binary [--from <path>] [--to <path>] [--kill] [--reopen]` | Safely generates an OS script (.ps1 or .sh), closes locking PIDs, rotates active binary to .old, copies new binary, and reopens/verifies version without locked file crashes. |
| ❌ Raw directory inspection of Tauri apps (`dir src-tauri`, `Get-ChildItem src-tauri`) | ✅ `gitmap lf src-tauri`, `gitmap aum search "<pattern>" src-tauri`, `gitmap pe -f tauri` | High-speed directory inventory and pipeline error filtering optimized for Tauri/Rust desktop projects. |

---

## Essential Command Cheat Sheet

### 1. Authentication & GitHub Credential Management
- `gitmap login` — Interactive authentication picker (browser login or secure token paste).
- `gitmap login --web` (alias `--browser`) — Non-interactive browser login via GitHub CLI / OAuth flow.
- `gitmap login --token <PAT>` — Non-interactive token ingestion (pre-validated against `api.github.com/user` before saving to global git config).
- `gitmap login --token <PAT> --no-verify` — Offline token ingestion without network probe.
- `gitmap login --status` — Displays current masked credential state and active user.
- `gitmap login set global ai-auto` — Locks workstation-global Git identity to `gitmap automation bot - runs autonomously <info@rgitmap.com>` while safely preserving the original developer identity in `original_user.json`.
- `gitmap login set ai-auto [agent-name]` — Locks Git identity to `gitmap automation bot - <agent-name> - runs autonomously <info@rgitmap.com>` (defaults to `gitmap automation bot - runs autonomously <info@rgitmap.com>` if agent name omitted).
- `gitmap token list` — Resolves active GitHub token and shows credential origin.
- `gitmap logout` — Purges stored GitHub credentials from Git global configuration.

### 2. Workspace & Multi-Repository Status
- `gitmap status` (alias `gitmap st`) — Comprehensive multi-repo status table across tracked workspace.
- `gitmap status --dirty` (`-d`) — Filters solely to repositories with uncommitted working changes.
- `gitmap status --ahead` — Isolates repositories with local commits waiting to push.
- `gitmap status --behind` — Isolates repositories with remote commits waiting to pull.
- `gitmap status --json` (`-j`) — Emits structured JSON array of repository states for automated scripts.
- `gitmap status --table` — Forces rich ANSI status table with colored columns.
- `gitmap status --compact` — Emits single-line compact summary per repository.
- `gitmap has-any-updates` (alias `gitmap hau`, `gitmap hac`) — Checks remote tracking branch for incoming commits.
- `gitmap latest-branch` (alias `gitmap lb`) — Discovers the most recently updated remote branch.
- `gitmap watch` (alias `gitmap w`) — Live-refresh terminal dashboard monitoring repository status changes.

### 3. Script Execution & Runner Engines
- `gitmap py "<code-or-script>"` — High-performance cross-platform Python script execution.
- `gitmap py -c "<code-or-expression>"` — Direct inline Python command evaluation.
- `gitmap pwsh "<cmd>"` / `gitmap ps "<cmd>"` — Cross-platform PowerShell execution with `-NoProfile`.
- `gitmap ps -c "<cmd>"` — Direct inline PowerShell command evaluation.
- `gitmap bash "<cmd>"` / `gitmap sh "<cmd>"` — Cross-platform Bash execution.
- `gitmap bash -c "<cmd>"` — Direct inline POSIX Bash evaluation.
- All runner commands support `--dry-run` and `--json` flags for pipeline automation.

### 4. Storage, Cache & Secrets Offloading
- `gitmap rs file <filepath> [--repo <name>]` — Copies secret file into `repo-secrets` and auto-pushes.
- `gitmap rs folder <folderpath> [--repo <name>]` — Copies secret directory into `repo-secrets` and auto-pushes.
- `gitmap rs text "<secret>" [--slug <slug>]` — Writes sensitive text into `repo-secrets/<repo>/<slug>.txt`.
- `gitmap rc file <filepath> [--repo <name>]` — Copies reusable script into `repo-cache` (`repo-storage`) and auto-pushes.
- `gitmap rc folder <folderpath> [--repo <name>]` — Copies reusable fixture directory into `repo-cache`.
- `gitmap rc text "<content>" --slug <slug> --ext <.ps1|.py>` — Writes test harnesses into `repo-cache` for permanent cross-repo reuse.
- `gitmap cd rs` / `gitmap cd rc` — Navigates directly to `repo-secrets` or `repo-cache`.

### 5. High-Performance Automation (AUM) & File Discovery
- `gitmap aum search "<pattern>" [dir] [--ext <ext>] [-r] [-i]` (alias: `gitmap aum grep`) — Multi-core streaming live search with lazy regex and binary filtering. ALWAYS scope with target `[dir]` and `--ext`. Replaces slow PowerShell `Select-String`, `Get-ChildItem -Recurse`, and `git grep`. TOTAL BAN on PowerShell `Select-String`, `rg`, `ripgrep`, and `git grep`.
  * **Named-Flag Support:** Fully supports `--query <pattern>` (or `-q <pattern>`, `--pattern <pattern>`) and `--root <dir>` (or `--dir <dir>`, `-d <dir>`, `--path <dir>`, `-p <dir>`) in any argument position (e.g. `gitmap aum search --root "cli" --query "cmdpipeline"`, `gitmap aum search --query "func Run" --dir cli -e .go`). Positional and named flags can be freely interchanged.
- `gitmap search "<query>" [--limit <n>]` — Instant SQLite cached symbol & keyword search across scanned repositories using DH2D split-db hot cache (also accepts `--query <q>`, `-q <q>`, `--pattern <p>`, `--root <r>`).
- `gitmap aum guard` — Enforces 500 KB limit, large JSON exclusion, and binary null-byte probe.
- `gitmap aum sequence` — Markdown sequence gap detector and `# XX Title` autofixer.
- `gitmap aum exclude list` — Query persistent search exclusions from SQLite.
- `gitmap aum newlines --fix` — Polyglot CRLF to LF and trailing whitespace normalizer.
- `gitmap aum cache status` — Sub-millisecond in-memory cache status.
- `gitmap aum locate [tool]` — Ultra-fast tool finder (<15ms, e.g. `vcvarsall.bat`, `msbuild`, `python`; replaces slow PowerShell traversal).
- `gitmap aum benchmark all` — Side-by-side Go vs Python execution benchmarks.
- `gitmap find "<pattern>" [-ext <ext>]` — Find files matching glob pattern in <10ms across 10,000+ files.
- `gitmap find-files <name>` (alias: `gitmap ff <name>`) — Find exact filename with optional `-ext`.
- `gitmap find-files-any <str>` (alias: `gitmap ffa <str>`) — Find files matching substring.
- `gitmap find-files-startswith <prefix>` (alias: `gitmap ffs <prefix>`) — Find by filename prefix.
- `gitmap find-files-endswith <suffix>` (alias: `gitmap ffe <suffix>`) — Find by filename suffix (e.g. `_test.go`).
- `gitmap list-files [dir]` (alias: `gitmap lf [dir]`) — List relative file paths matching pattern or directory.
- `gitmap cat <filepath> [--raw]` (alias: `gitmap ct`) — Direct zero-disk streaming of file content into stdout with flexible path parsing and terminal binary protection:
  * **Flexible Path Support:** Handles conversational strings, numbered list prefixes (`1. path`), inline comments (`path # note`), bracketed paths (`@[path]`), and Markdown file URIs (`![Image](file:///path)`).
  * **Safe Image Reader:** For image assets (`.png`, `.jpg`, `.svg`, `.webp`), automatically renders a structured metadata card (format, MIME, dimensions in px, byte size, SHA-256) instead of dumping binary control characters to stdout. Use `--raw` flag to stream raw binary bytes to pipes.
- `gitmap replace <old> <new>` — Exact literal string replacement with audit trail.
- `gitmap replace-regex <pat> <subst>` — Regex replacement across repository.

#### Directory Inventory & Path-Scoped Listing
- `gitmap lf <dir>` — Instant directory file listing (e.g. `gitmap lf src-tauri/src/modules/repo_db`).
- `gitmap find "*.*" <dir>` — Glob-filtered directory listing without disk traversal overhead.
- `gitmap aum search "<pattern>" <dir> -e <.ext>` — Scoped live code/regex search directly within target directory and extension.
- **Concrete Developer Workflows & Real-World Examples:**
  * **Module Directory Inventory:** `gitmap lf src-tauri/src/modules/repo_db` — Lists files within specific module path in <1ms without invoking PowerShell `Get-ChildItem`.
  * **Spec Folder Inspection:** `gitmap find "*.*" 02-spec/21-app/gitmap-commands-agent-spawning-heatmap-sync` — Discovers all spec files within a feature directory.
  * **Path-Scoped Code Search:** `gitmap aum search "HeatmapActionLog" src-tauri/src -e .rs` — Scoped search restricted to Rust source tree.
  * **Subtask Plan Inventory:** `gitmap lf .ai-memory/plans/subtasks` — Rapidly enumerates in-flight and pending subtask manifests.

### 6. Autonomous CI/CD Self-Healing (Pipeline AI)
- `gitmap pipeline-ai status --json` — Check workflow execution state, active branch, and ETA.
- `gitmap pipeline-ai status -t <eta>` — Wait dynamically for pipeline completion without tight polling.
- `gitmap pipeline error-logs` (alias: `gitmap pe`) — Extract failing step logs to file for 4-part RCA.
- `gitmap pe -t` — Telemetry mode extracting concise failure summaries.
- `gitmap pe history-ai` — Analyze CI/CD pipeline history across branches and recent runs.
- `gitmap pipeline purge` — Actions zero-storage purge maintaining 0.0 GB footprint.

### 7. Semantic Commit Suite, CFR & Workgroup

Disciplined atomic commit workflow enforcing hyphen-separated messages (`<Module> - <summary>`) and upstream delivery invariants.

#### Commands & Behavioral Guidelines
- `gitmap commit "<msg>"` (alias `gitmap cm "<msg>"`) — Fast 1-step auto-stage (`git add -A`) and local commit without semantic prefixing and without remote push. Ideal for quick local iterations and WIP checkpoints:
  `gitmap commit "wip - draft initial types and interfaces"`
- `gitmap cpf "<Module> - <summary>"` — Features ONLY. Stages all, formats with `Feature: <Module> - <summary>`, and pushes:
  `gitmap cpf "Skills - enrich gitmap skill with 12 real world categories"`
- `gitmap cbf "<Module> - <summary>"` (alias `gitmap cpb`) — Bug fixes ONLY. Stages all, formats with `Bug: <Module> - <summary>`, and pushes:
  `gitmap cbf "Cloner - correct concurrency flag parsing and worker limits"`
- `gitmap cfr "<Module> - <summary>"` — Refactoring ONLY. Stages all, formats with `Refactor: <Module> - <summary>`, and pushes:
  `gitmap cfr "Guidelines - decompose function argument structs"`
- `gitmap cpc "<Module> - <summary>"` — Chores, docs, sync, and maintenance ONLY. Stages all, formats with `Chore: <Module> - <summary>`, and pushes:
  `gitmap cpc "Docs - realign companion skills and remove help flags"`
- `gitmap cpr "<Module> - <summary>"` — Releases ONLY. Stages all, formats with `Release: <Module> - <summary>`, and pushes. MANDATORY verb for version bumps and releases:
  `gitmap cpr "v2.18.0 - bump minor version and update changelog"`
- `gitmap pcp "<Module> - <summary>"` — Pull latest with rebase, stage all, commit, and push with preflight verification:
  `gitmap pcp "Skills - sync cross-platform runner templates"`
- `gitmap cpf/cpb/cpc/cfr/cpr "<Module> - <summary>" --ai-agent <name>` — Attaches designated subagent identity (e.g. `--ai-agent "worker-01"`) to commit author metadata, GitMap Split-DB commit history (`repodb/commit_history.db`), and cognitive heatmap tables (`repodb/heatmap.db`). Default author becomes `gitmap automation bot - <name> - runs autonomously <info@rgitmap.com>`:
  `gitmap cpf "Module - implement feature" --ai-agent "worker-01"`
- `gitmap merge-right <left> <right> [-y]` (alias `gitmap mr`) — One-way file-level merge writing only into RIGHT.
- `gitmap merge-left <left> <right> [-y]` (alias `gitmap ml`) — One-way file-level merge writing only into LEFT.
- `gitmap merge-both <left> <right> [-y]` (alias `gitmap mb`) — Two-way bidirectional merge synchronizing endpoints.
- `gitmap merge-in <target> <inputs...> [flags]` (alias `gitmap commit-in` / `gitmap cin`) — Multi-source chronological commit replay.
- `gitmap commit-both <left> <right> [flags]` (alias `gitmap cmb`) — Bidirectional commit replay between repositories.
- `gitmap commit-left <left> <right> [flags]` (alias `gitmap cml`) — Replay commits from right repo onto left repo.
- `gitmap commit-right <left> <right> [flags]` (alias `gitmap cmr`) — Replay commits from left repo onto right repo.
- `gitmap thelp [topic]` (alias `gitmap help`) — Access terminal help screens with rich two-column menus and executable examples. Run `gitmap thelp specific` to browse topics.
- `gitmap workgroup ls` (alias `gitmap wg ls`) — Inspect folder-specific workgroups, bound SSH keys, and ancestor bindings.
- `gitmap workgroup add <name> <dir> --ssh-key <key>` — Register folder workgroup and bind dedicated SSH private key.
- `gitmap workgroup bind <dir> <key>` — Bind folder to dedicated SSH key.
- `gitmap workgroup test-ssh [name|dir]` — Test remote SSH handshake authentication against GitHub/GitLab.
- `gitmap pull [repo]` (alias `gitmap p`) — Pull targeted repository. In current directory, simply run `gitmap pull`.
- `gitmap pull-all` (alias `gitmap pa`) — Pull all repositories in workspace.
- `gitmap fix [repo] [action]` — Apply remediation to repo (aliases: `stash`, `wip`, `discard`).
- `gitmap lowercase` (alias `gitmap lcf`) — Safe 2-step `git mv` file case normalization.
- `gitmap lowercase-readme` — Safe 2-step `git mv` case normalization for root `readme.md`.
- **Subsystem Module Standard (Title Case or Lowercase Allowed):** The `<module>` token MUST represent a functional codebase subsystem (e.g. `Auth`, `Storage`, `Ui`, `Guidelines`, `Scanner`, `Cli`, `Sync`, `Prompt`). Title Case (`Auth`, `Storage`, `Ui`) or lowercase (`auth`, `storage`, `ui`) is completely valid. ALL-CAPS shout case (`AUTH`, `STORAGE`) should not be used.
- **TOTAL BAN ON 'commit' / 'commits' AS MODULE:** The word `commit` or `commits` MUST NEVER be used as a module name (e.g. NEVER run `gitmap cpc "commit - ..."` or `gitmap cpf "COMMIT - ..."`). Committing is already the version control action being performed by Git and GitMap; writing `commit` as the module produces meaningless tautological stuttering (`Chore: commit - ...`). If modifying commit-related tooling, code, or documentation, name the specific component (e.g. `Formatter`, `Provenance`, `Cli`, `Docs`) or omit the module prefix entirely (`gitmap cpc "clean constructor helpers and strip redundant tokens"`).
- **DRY-RUN BEFORE EVERY GITMAP COMMIT (MANDATORY):** run `gitmap <verb> "<module> - <summary>" --dry-run` (short flag `-n`) before the real commit, then confirm with `git status` and `git diff`. With `--dry-run` the command prints exactly what it WOULD do — the full list of files it would stage (so a sweep of unrelated files is visible before a real run), the formatted commit message, the branch and HEAD, the push target (remote/branch and URL, or local-commit-only for plain `commit`/`cm`), whether it would pull first (`pcp`), and that tag/release actions are none — then exits WITHOUT mutating anything: no commit created, no push, no pull, and no staging changes left behind. It applies to `commit`/`cm`, `cpf`, `cpb`/`cbf`, `cpc`, `cfr`, `cpr`, and `pcp`.
- **STRICT VERB TAXONOMY (NEVER MISUSE CPF):** Never use `gitmap cpf` for chore, release, or bug fix tasks. Use `cpr` for releases, `cpc` for chores, `cbf` for bug fixes, `cfr` for refactors, and `cpf` for features.
- **ZERO NESTED CONVENTIONAL PREFIXES:** Never write redundant prefixes inside GitMap arguments (e.g. BANNED: `gitmap cpf "chore(release): bump version"`, `gitmap cpr "release: v1.0"`, `gitmap cpf "feat: ..."`). GitMap automatically adds the prefix!
- **ZERO THIRD BRACKETS & COLONS IN ARGUMENTS:** Commits MUST NOT use square brackets `[...]` (colloquially called "third brackets") or colons `:` inside the argument. Use space-padded hyphen: `"<module> - summary"`. If square brackets are passed, GitMap automatically sanitizes them to round parentheses.
- **MANDATORY COMPLETION INVARIANT ('No Push = Not Done'):** Flat commits (`gitmap commit`, `gitmap cm`) are local checkpoints only. For task completion, AGENTS.md §13 strictly mandates that work must be pushed to the remote tracking branch via `gitmap cpf`, `cpc`, `cbf`, `cfr`, or `cpr`.

### 8. Native Toolchains & Process Termination (Rust, Cargo, Proc Kill)
- `gitmap cargo [args...]` (alias: `gitmap crg`) — Native Cargo runner with automatic LLVM/libclang path injection (`LIBCLANG_PATH` and `PATH`) and install suggestions.
- `gitmap rust status` (alias: `gitmap rst`, `gitmap rustc`) — Unified Rust frontend inspecting rustc, cargo, rustup, and LLVM/libclang versions and directories. Supports `gitmap rust build`, `gitmap rust test`, `gitmap rust check`, `gitmap rust clippy`.
- `gitmap proc kill <pid> [--force]` (alias: `gitmap kill <pid|name|:port>`, `gitmap terminate`) — Safe, cross-platform process termination by PID or name with safety guards protecting critical system processes (PID 0, 1, 4) and IDEs/shells (`code`, `antigravity`, `pwsh`, `bash`).

### 9. Autonomous Agent Onboarding, Curriculum & Memory Training (LLM)
- `gitmap llm train [ai-name] [flags]` (aliases: `gitmap train`, `gitmap llm-train`, `gitmap llm chain`) — Full 6-phase self-looping onboarding curriculum (Discovery, Refactoring, Verification, Semantic Commit, Telemetry, Heal & Fix), auto-generates Antigravity skill, author/sponsor attribution.
- `gitmap llm train muse [--days <N>]` (aliases: `gitmap train muse`, `gitmap train-muse`, `gitmap llm-train-muse`) — Synchronize Muse N-day memory digest from `repo-cache` (fallback `repo-secrets`), diff against watermark commit, output AI traversal checklist for sequential memory reading, and reconstitute `chat-seeds`.
- `gitmap llm train --heatmap <path|url>` (aliases: `gitmap train --heatmap`, `gitmap llm-train --heatmap`) — Ingests cognitive decision traces from local SQLite heatmap databases or remote REST endpoints into LLM training curriculum. Supports `--heatmap-repo <slug>` and `--heatmap-task <slug>`.
- `gitmap llm train --loop` / `gitmap llm train --self-loop <N>` — Autonomous self-looping flags executing continuous 6-phase AI self-looping execution cycles.
- `gitmap llm train --heal-apply <stash|wip|discard>` — Non-interactively applies git-state remediation recipes during Phase 6a.
- `gitmap llm train --skill-path <path>` — Target filesystem path for writing the generated skill markdown file (default: `.agents/skills/gitmap/skill.md`).
- `gitmap llm train --text-only` — Outputs curriculum text to stdout without writing skill file to disk.
- `gitmap llm train --urls` — Outputs the 8 authoritative public documentation registry links for LLM ingestion.
- `gitmap llm train --url` — Outputs raw public URL to `llm.md` specification.
- `gitmap llm train -j` / `gitmap llm train --json` — Emits structured machine-readable JSON curriculum and document metadata.
- `gitmap llm-docs` (alias: `gitmap ld`) — Consolidated markdown command matrix reference for LLMs.
- `gitmap llm` — Displays full LLM specification and operational guidelines.

#### Autonomous Training Flag Matrix
| Flag | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `--days <N>` | `int` | `7` | Sets memory digest inspection window in days for agent targets (e.g. `gitmap llm train muse --days 7`). |
| `--heatmap <path\|url>` | `string` | `""` | Path to local SQLite heatmap database or remote HTTP(S) URL for cognitive decision trace ingestion. |
| `--heatmap-repo <slug>` | `string` | `""` | Filters heatmap training actions by repository slug. |
| `--heatmap-task <slug>` | `string` | `""` | Filters heatmap training actions by task slug. |
| `--loop` | `bool` | `false` | Triggers continuous autonomous 6-phase AI self-looping execution cycles. |
| `--self-loop <N>` | `int` | `0` | Specifies exact number of consecutive iterations for AI self-loop. |
| `--heal-apply <recipe>` | `string` | `""` | Non-interactively applies git-state remediation recipes during Phase 6a (`stash`, `wip`, `discard`). |
| `--skill-path <path>` | `string` | `.agents/skills/gitmap/skill.md` | Target filesystem path for writing generated skill markdown file. |
| `--text-only` | `bool` | `false` | Emits complete curriculum text to stdout without writing files to disk. |
| `--urls` | `bool` | `false` | Outputs complete 8-document Authoritative Public Documentation Registry for LLM memory ingestion. |
| `-j`, `--json` | `bool` | `false` | Emits structured machine-readable JSON curriculum and document metadata. |

### 9. Multi-Repo & Toolchain Operations
- `gitmap pae --json` — Multi-repo pull with compact JSON telemetry (use only when explicitly requested; ban routine polling).
- `gitmap cargo status` — Inspect Rust and Cargo toolchain status.
- `gitmap install cargo --check` — Validate Cargo and Rust toolchain installation.
- `gitmap install cargo` — Install Rust toolchain if missing.
- `gitmap install --list` — Discover developer toolchains, profiles, and runtime packages.

### 10. Native Fleet Synchronization & SQLite Agent Task Engine
- `gitmap sync` — Synchronizes canonical prompts, skills, shared specs (`02-spec/01-20`), and additive scripts across all 43 registered fleet repositories.
- `gitmap sync --projects <path|json>` — Accepts a path to JSON file or inline JSON array of repositories (e.g. `'[{"folder": "cat-my"}]'`).
- `gitmap sync --repo <name>` — Synchronizes a single target repository by name.
- `gitmap sync --workers <N>` — Sets parallel worker pool concurrency (default: 8).
- `gitmap sync --dry-run` — Previews changes across all repositories without making git or file mutations.
- `gitmap sync --no-push` — Applies changes and commits locally without pushing to remote.
- `gitmap sync --no-release` — Disables post-sync SemVer release tagging.
- `gitmap agent task init --name "<task>" --budget <N>` — Initializes SQLite task manager in `.ai-memory/temp-agents/<slug>/agent-task.db`. If task already exists, reports parent task ID, active subagents, subtasks status, and touched files. Supports prompt MD5 hash verification to resume seamlessly without redundant reads.
- `gitmap agent task status [--json] [--ddl]` — Emits structured task status, subtask progress, and schema validation.
- `gitmap agent task ls` — Lists active and completed parent tasks across the workspace.
- `gitmap agent task pending` — Lists or counts pending subtasks across agent tasks.
- `gitmap agent task completed` — Lists completed root-level parent tasks with proof.
- [x] `gitmap agent subtask add --code "<code>" --slug "<slug>" --title "<title>" [--owned "<json>"] [--json]` — Adds a decomposed subtask bounding box to the active parent task queue. When `--json` is provided, outputs a grouped JSON response collapsing shared database metadata (`databasePath`, `parentTaskId`, `addedCount`).
- `gitmap agent subtask ls [--json|-j] [task-id]` — Lists all subtasks, their assignees, statuses, and file boxes. When `--json` (`-j`) is passed, emits a structured JSON envelope containing `status`, `databasePath`, `parentTaskId`, `summary` rollups (total, pending, in-progress, completed, failed), and `subtasks` array.
- `gitmap agent subtask claim --agent <agent-name>` — Claims the next pending subtask atomically.
- `gitmap agent subtask claim-files` — Pre-declares file leases for collision and contention tracking.
- `gitmap agent log --subtask <id> --agent <name> --action <action> --file <path> --details <desc>` — Logs in-flight agent action for crash forensics.
- `gitmap agent subtask start [id] --agent <name> [--json|-j]` — Marks subtask in-progress. Outputs compact JSON with `--json`.
- `gitmap agent subtask complete <id> --agent <name> --evidence <evidence> [--json|-j]` — Marks subtask completed with verification evidence. When `--json` is provided, outputs a compact JSON response confirming completion.
- `gitmap agent subtask fail [id] --agent <name> --reason <reason> [--json|-j]` — Marks subtask failed with diagnostic rationale. Outputs compact JSON with `--json`.
- `gitmap agent subtask join --subtask <id> --agent <name> --role "Assistant: <focus>"` — Dynamically joins an in-flight subtask as a collaborative assistant for parallel reading, caller analysis, or verification without file collisions.
- `gitmap agent task tree` — Renders hierarchical task swarm tree (Lead Orchestrator -> Branch Coordinators -> Leaf Workers & Dynamic Assistants).
- `gitmap agent diagnose` — Performs deep autopsy and diagnostic reporting on agent task runs.
- `gitmap agent crashed` — Detects crashed or abandoned agent tasks.
- `gitmap agent temp-clear` — Safely deletes ephemeral `.ai-memory/temp-agents/` directory.
- **Cognitive Heatmap & Multi-AI Attribution:** All task actions and cognitive states integrate seamlessly with `gitmap heatmap log --json '<payload>'`, recording model (`model`) and tool (`tool`) attribution for unified swarm forensics. Newly joining subagents inspect active progress via `gitmap agent subtask ls --json`, `gitmap rc heatmap tree`, and `gitmap heatmap search [slug]`.

### 10.1 AI Agent Communication & Structured Subtask Envelopes

To eliminate brittle ANSI terminal scraping and regex extraction errors in multi-agent swarms, all AI agents coordinating task execution MUST invoke GitMap commands with the `--json` / `-j` flag:

#### `gitmap agent subtask ls --json` Response Schema
Emits a structured JSON envelope with rollup statistics, relative database path, and classified file ownership:
```json
{
  "status": "success",
  "databasePath": ".ai-memory/temp-agents/task-294/agent-task.db",
  "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
  "summary": {
    "totalCount": 5,
    "pendingCount": 1,
    "inProgressCount": 2,
    "completedCount": 2,
    "failedCount": 0
  },
  "subtasks": [
    {
      "subtaskId": "sub-120000-subtask-01",
      "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
      "taskCode": "Subtask-01",
      "taskSlug": "294-code-refactoring-and-coding-guideline-update-subtask-01",
      "title": "Architecture Specification & Subtask Ledger",
      "status": "COMPLETED",
      "assignedAgentRole": "Worker 01",
      "ownedFilesJson": "[\"cli/main.go\"]",
      "files": [
        {
          "path": "cli/main.go",
          "exists": true,
          "fileType": "go"
        }
      ]
    }
  ]
}
```

#### `gitmap agent subtask add --json` Grouped Response Schema
Groups common metadata once at the top level instead of repeating across individual subtask items:
```json
{
  "status": "success",
  "databasePath": ".ai-memory/temp-agents/task-294/agent-task.db",
  "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
  "addedCount": 2,
  "subtasks": [
    {
      "subtaskId": "sub-120000-subtask-01",
      "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
      "taskCode": "Subtask-01",
      "taskSlug": "294-code-refactoring-and-coding-guideline-update-subtask-01",
      "title": "Refactor tokens",
      "status": "PENDING"
    },
    {
      "subtaskId": "sub-120001-subtask-02",
      "parentTaskId": "294-code-refactoring-and-coding-guideline-update",
      "taskCode": "Subtask-02",
      "taskSlug": "294-code-refactoring-and-coding-guideline-update-subtask-02",
      "title": "Extract repository",
      "status": "PENDING"
    }
  ]
}
```

#### `gitmap agent subtask complete --json` Response Schema
Outputs a compact, structured JSON object confirming completion with relative database path and evidence:
```json
{
  "status": "completed",
  "subtaskId": "sub-111615-task-01",
  "agentRole": "Spec Author 01",
  "databasePath": ".ai-memory/temp-agents/task-294/agent-task.db",
  "details": "PASS exit 0, 01-architecture-spec.md created"
}
```

#### `gitmap agent subtask start --json` / `fail --json` Lifecycle Schema
Compact lifecycle transitions for AI worker status synchronization:
```json
{
  "status": "in_progress",
  "subtaskId": "sub-111615-task-01",
  "agentRole": "Spec Author 01",
  "databasePath": ".ai-memory/temp-agents/task-294/agent-task.db"
}
```

#### AI Agent Swarm Protocol: JSON Parsing Invariant
Autonomous subagents (Lead Orchestrators, Branch Coordinators, and Leaf Workers) communicating with GitMap or orchestrating subagent swarms MUST parse JSON outputs directly rather than attempting to regex-scrape plain text or ANSI tables. This maximizes parsing precision, prevents regex breakages caused by column truncation, and ensures zero unhandled exceptions.

> [!CAUTION]
> **TOTAL BAN ON RAW PYTHON SQLITE (`python -c "import sqlite3..."`):**
> NEVER run inline Python sqlite3 snippets or execute raw database queries against `.ai-memory/temp-agents/<slug>/agent-task.db`.
> All agent task lifecycle operations, status inspections, claiming, logging, and diagnostics MUST use compiled `gitmap agent` CLI commands (`task init`, `task status`, `task ls`, `subtask add`, `subtask ls`, `subtask claim`, `subtask complete`, `subtask fail`, `log`, `diagnose`, `crashed`, `temp-clear`).

### 11. System Diagnostics & Performance Telemetry
- `gitmap os --json` — Structured machine and OS identity with kernel, build, and version telemetry.
- `gitmap cpu` / `gitmap processor` — Live & 5s cached CPU utilization, architecture, and core counts.
- `gitmap ram` — Live & 5s cached memory capacity, free/used MB, and visual utilization meter.
- `gitmap network [--json]` — Network adapter configuration, IPs, MAC addresses, link speeds, and throughput metrics.

---

## Autonomous AI Identity & Agent Attribution

Autonomous multi-agent swarms (e.g. V6 and V7 orchestrators) and automated pipelines require distinct identity tracking and commit provenance to differentiate automated AI agent contributions from human developer work.

### 1. Identity Locking Commands
- `gitmap login set global ai-auto`: locks global Git identity to `gitmap automation bot - runs autonomously <info@rgitmap.com>`. Safely backs up original human developer credentials to `original_user.json` before applying changes.
- `gitmap login set ai-auto [agent-name]`: locks Git identity to `gitmap automation bot - <agent-name> - runs autonomously <info@rgitmap.com>` (e.g. `gitmap login set ai-auto muse` sets `gitmap automation bot - muse - runs autonomously <info@rgitmap.com>`). If no agent name is specified, defaults to `gitmap automation bot - runs autonomously <info@rgitmap.com>`.

### 2. Semantic Commit Provenance with `--ai-agent`
- `gitmap cpf/cpb/cpc/cfr/cpr --ai-agent <name>`: attaches designated subagent identity to commit and Split-DB heatmap tables.
- Example: `gitmap cpf "Module - implement feature" --ai-agent "worker-01"` automatically sets the author to `gitmap automation bot - worker-01 - runs autonomously <info@rgitmap.com>` if `--author` was omitted, recording `meta.AIName` and `meta.IsAIMode`.

### 3. Heatmap Attribution Mechanics
GitMap cognitive heatmaps natively capture agent attribution across both local Split-DB (`repodb/heatmap.db` and `.ai-memory/heatmap.db`) and companion `repo-cache` databases:
- **Status & Summary Attribution:** `gitmap heatmap status` and `gitmap heatmap tree` display active agent models, subagent identities, and touch distributions per agent.
- **Granular Action Records:** Every action logged via `gitmap heatmap log --json '{"task":"<task-slug>","action":"EDIT","file":"<path>","state":"EXECUTION","expected":"<goal>","model":"<model>","tool":"<tool>"}'` stores the calling model, tool environment, and agent identifier in `HeatmapActionLog`.
- **Commit Telemetry Binding:** Commits executed with `--ai-agent <name>` automatically synchronize with heatmap tables, recording the exact subagent responsible for file touches, modifications, and git hashes in `repodb/heatmap.db` and `repodb/commit_history.db`.
- **Swarm Collision Forensics:** `gitmap agent heatmap` and `gitmap heatmap discuss` aggregate agent leases and write claims, displaying which specific subagent owns or recently edited a given file box.

---

### 12. File Contention Heatmap, Cognitive Telemetry & Agent Activity Observation
- `gitmap heatmap` — **Heatmap Database Discovery Catalog**: Surveys all 3 tiers of heatmap databases (local root `<EffectiveDataDir>/heatmap.db`, per-repo split DBs `<EffectiveDataDir>/heatmap/<repo>/sql.db`, and repo-cache multi-AI DBs `repo-cache/.heatmap/<seq>-<ai>/`), maps registered repository aliases via `db.ListAliasesWithRepo()`, aggregates recent tasks across 2–5 repositories (minimum 5 tasks), and provides dynamic interactive AI suggestion pills.
- `gitmap heatmap --repo <alias>` / `gitmap rc heatmap [--repo <alias>]` — Repository-scoped heatmap discovery and companion database inspection with alias enrichment and action rollups.
- `gitmap heatmap log --json '{"task":"<task-slug>","action":"READ|EDIT","file":"<path>","state":"<STATE>","expected":"<goal>","model":"<model>","tool":"<tool>"}'` — **One-shot JSON payload ingestion** logging granular file touches, cognitive reasoning states, expected outcomes, and AI provenance directly into `.ai-memory/heatmap.db`. Eliminates shell quoting failures and escapes.
- `gitmap heatmap log --task <task> --action <action> --file <file> --state <state> --expected "<goal>" [--model <name>] [--tool <name>]` — Classical CLI flag ingestion with explicit AI provenance attribution.
- `gitmap heatmap update --id <id> --outcome "<outcome>" --met <true|false>` — Updates logged action with final outcome and verified expectation status.
- `gitmap heatmap search [query] [-f <file>] [-t <task>]` — Searches logged actions, code snippets, reasoning notes, and execution history across tasks and repositories.
- `gitmap heatmap tree [-t <task-slug>] [--json]` — Renders hierarchical tree view of touched files, line ranges, action counts, and AI models.
- `gitmap heatmap status [-t <task-slug>] [--json]` — Summary of actions logged, unique files touched, expectation metrics, verification states, and active agent models.
- `gitmap heatmap discuss [--window <duration>] [--threshold <n>] [--json]` — Multi-agent file contention analysis and collision blast radius report before touching shared files.
- `gitmap heatmap live` — Real-time continuous stream of agent actions, file touches, and reasoning across concurrent agents.
- `gitmap agent log --subtask <id> --agent "<name>" --action "<read|modify|plan>" --file "<path>" --details "<reasoning>"` — In-flight subtask action logging with explicit rationale.
- `gitmap agent heatmap [--parent <slug>] [--json]` — Ranks files by write-claim count across agent tasks and active subtasks.
- `gitmap full-summary` (alias `gitmap fs`) — 48-hour active repository activity heatmap TreeView showing commit velocity and dirty working states.
- **Dual-Mode `--json` Operation:** GitMap cleanly separates one-shot JSON payload ingestion (`gitmap heatmap log --json '<payload>'`) from structured JSON output formatting (`gitmap agent heatmap --json`, `gitmap heatmap status --json`, `gitmap heatmap tree --json`). Both can be combined: `gitmap heatmap log --json '<payload>'` ingests the action and immediately returns a structured JSON confirmation envelope.
- **Mandatory AI Provenance Attribution:** Every cognitive telemetry entry records the active model (`ModelName`) and host environment (`ToolName`) via JSON payload (`"model"`, `"tool"`) or flags (`--model`, `--tool`).
  * **Google Antigravity IDE:** `model: "gemini-2.5-pro"` (or `"gemini-2.0-flash"`), `tool: "antigravity-ide"`
  * **Muse Autonomous Fleet:** `model: "muse"`, `tool: "muse-runner"`
  * **Claude Code / Standalone CLI:** `model: "claude-3-7-sonnet"`, `tool: "claude-code"`
- **MANDATORY HEATMAP & CONCURRENCY RULE:** In multi-agent autonomous sessions, agents MUST record all file read/edit decisions with reasoning via one-shot `gitmap heatmap log --json`, and inspect active file leases and contention hotspots via `gitmap agent heatmap` (or `gitmap heatmap status`) before editing shared files.
- **Roadmap-Connected Heatmap for Dynamic Join:** Newly spawned or dynamically joining subagents (`gitmap agent subtask join`) execute `gitmap rc heatmap tree` and `gitmap heatmap search [slug]` to inspect prior file reads, architectural decisions, and remaining roadmap tasks in `.ai-memory/plans/<slug>.md`, eliminating redundant reads and instantly synchronizing with existing work.

#### Mandatory 4-Gate Cognitive Heatmap Workflow
In multi-agent and autonomous workflows, every file modification or major task phase MUST strictly follow the 4-Gate Cognitive Heatmap Workflow:
- **Gate 0: Pre-Flight Discovery & Schema Verification:**
  `gitmap heatmap status` and `gitmap heatmap tables` (or `gitmap rc heatmap status` and `gitmap rc heatmap tree`)
  Inspect active task telemetry, verify SQLite schema tables, and detect any active file leases or contention hotspots before touching files.
- **Gate 1: Intent, State & Lease Declaration (One-Shot JSON Ingestion):**
  `gitmap heatmap log --json '{"task":"<task-slug>","action":"READ|EDIT|CREATE|DELETE","file":"<path>","state":"<STATE>","expected":"<goal>","model":"<model>","tool":"<tool>"}'`
  Declare intent, file path, cognitive reasoning state (`EXPLORATION`, `PLANNING`, `EXECUTION`, `VERIFICATION`), expected goal, model name, and tool environment before reading or modifying files. Returns action ID for lifecycle tracking.
- **Gate 2: Outcome Recording & Telemetry Close:**
  `gitmap heatmap update --id <id> --outcome "<outcome>" --met true`
  Record actual outcome, verify if expectation was met (`--met true` or `--met false`), and close the telemetry record immediately upon completing execution.
- **Gate 3: Verification & Hierarchy Audit:**
  `gitmap heatmap tree --task <task>` and `gitmap heatmap search <task-slug>`
  Inspect hierarchical touch tree, file churn counts, model attributions, and search logged decisions across the task.

### 13. Multi-AI Companion & repo-cache Synchronization (RC Heatmap Suite)
- **Why RC Split-DB Isolation is Mandatory:** When multiple autonomous AI engines (Antigravity/Gemini, Muse, Claude Code) collaborate simultaneously or sequentially, writing cognitive traces directly into `.ai-memory/heatmap.db` in the repository worktree causes SQLite database file lock contention (`database is locked`), dirty worktree conflicts, and merge collisions during git branch operations. The companion RC heatmap suite solves this by isolating telemetry into dedicated SQLite databases inside `repo-cache` (`repo-cache/.heatmap/<ai-name>/<repo>/heatmap.db`), keeping the main code repository 100% clean.
- **Automatic RC Vault Creation:** If `repo-cache` (`rc`) or `.ai-memory/heatmap.db` does not exist, GitMap auto-provisions the required split-db structure linked to the active mission and folder.
- `gitmap rc heatmap status` — Inspects companion RC telemetry status, active leases, and logged actions across AI engines.
- `gitmap rc heatmap tree` — Renders hierarchical roadmap and touched file tree in RC vault, showing all previous reads, edits, and decisions.
- `gitmap rc heatmap live` — Real-time activity stream in RC vault across all active AI engines.
- `gitmap rc heatmap sync` — Bidirectionally reconciles in-flight decisions with `repo-cache` upstream without creating workspace git conflicts.
- `gitmap rc heatmap tables` / `gitmap rc heatmap query "<sql>"` — Direct read-only inspection of RC companion databases with busy-timeout guards.
- `gitmap heatmap rc list-ai` — Lists all supported AI namespaces (`antigravity`, `muse`, `cloud-code`, `codex`, `grok`, `cursor`) and displays ambient detected engine.
- `gitmap heatmap rc log [flags]` — Logs cognitive actions and decisions into dedicated `repo-cache` isolated split-databases.
- `gitmap heatmap rc commit [-m "<msg>"]` — Creates a telemetry checkpoint commit in `repo-cache` and marks local changes clean.
- **Roadmap Traversal for Newly Joined Agents:** Newly joined or dynamically respawned subagents MUST run `gitmap rc heatmap tree` and `gitmap heatmap search [slug]` prior to authoring code. This enables the agent to discover:
  1. Exactly which files prior agents examined and why (`READ` actions and reasoning notes).
  2. Which architectural decisions and code edits were committed (`EDIT` actions and outcomes).
  3. Which roadmap subtasks remain pending, avoiding redundant discovery cycles.

### 14. AI Decision Heatmap & RC Admin Telemetry
- `gitmap heatmap tables [--db <root|repo|local>] [--format <table|json|yaml>]` — Inspects the SQLite master catalog and inventories all tables, views, column counts, and row counts in the active heatmap database.
- `gitmap heatmap rc tables [--ai <name>] [--repo <slug>] [--db <root|repo|local>] [--format <table|json|yaml>]` — Inventories tables across companion `repo-cache` databases (`root-heatmap.db` or `<repo-slug>/heatmap.db`).
- `gitmap heatmap view <table> [--limit <n>] [--format json|yaml|table] [--task <slug>] [--repo <slug>] [--id <id>]` — Queries records from a specific heatmap table with column-aware filtering and bounded limits.
- `gitmap heatmap rc view <table> [--limit <n>] [--format json|yaml|table] [--task <slug>] [--repo <slug>] [--id <id>] [--ai <name>]` — Views table records from the target `repo-cache` companion database.
- `gitmap heatmap query "<sql>" [--db <path|root|repo|local>] [--format table|json|yaml] [--limit <n>]` — Executes arbitrary read-only SQL queries against the local or split heatmap database with automatic `mode=ro` protection, `PRAGMA query_only = ON;`, and 5000ms busy timeout.
- `gitmap heatmap rc query "<sql>" [--ai <name>] [--repo <slug>] [--tier <repo|root>] [--format table|json|yaml] [--limit <n>]` — Direct read-only SQL execution against companion `repo-cache` split or root databases.

#### Tri-Alias Routing Parity
All companion RC commands support full tri-alias routing parity:
- `gitmap rc heatmap <status|tables|query|view|tree|live|sync>`
- `gitmap heatmap rc <status|tables|query|view|tree|live|sync>`
- `gitmap heatmap rc-heatmap <status|tables|query|view|tree|live|sync>`

#### Shorthand Commands vs Direct SQL Queries

Instead of typing multi-line SQL queries, use GitMap's dedicated shorthand commands (with direct SQL queries available as secondary alternatives):

| Operational Goal | Dedicated GitMap Command (Preferred) | Direct SQL Query Alternative |
| :--- | :--- | :--- |
| **Chronological Task Trace** | `gitmap heatmap tree --task <slug>` or `gitmap heatmap view HeatmapActionLog --task <slug> --json` | `gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, ReasoningState, IsExpectationMet FROM HeatmapActionLog WHERE TaskSlug = '<slug>' ORDER BY HeatmapActionLogId ASC" --json` |
| **Recent Activity Stream** | `gitmap heatmap recent` (alias: `latest`) or `gitmap rc heatmap recent` | `gitmap heatmap query "SELECT HeatmapActionLogId, TaskSlug, ActionType, FilePath, ReasoningState, ExpectedOutcome, ActualOutcome, CreatedAt FROM HeatmapActionLog ORDER BY HeatmapActionLogId DESC LIMIT 20"` |
| **Unmet Expectations & Failures** | `gitmap heatmap failed` (aliases: `failures`, `unmet`, `broken`) or `gitmap rc heatmap failed` | `gitmap heatmap query "SELECT HeatmapActionLogId, TaskSlug, FilePath, ExpectedOutcome, ActualOutcome FROM HeatmapActionLog WHERE IsExpectationMet = 0" --json` |
| **File Churn & Hotspot Analysis** | `gitmap heatmap hotfiles` (aliases: `hot`, `hotspots`) or `gitmap rc heatmap hotfiles` | `gitmap heatmap query "SELECT FilePath, COUNT(*) AS Touches, SUM(CASE WHEN ActionType = 'READ' THEN 1 ELSE 0 END) AS Reads, SUM(CASE WHEN ActionType = 'EDIT' THEN 1 ELSE 0 END) AS Edits FROM HeatmapActionLog GROUP BY FilePath ORDER BY Touches DESC LIMIT 20"` |
| **Parent Tasks & Spec Tracking** | `gitmap heatmap tasks` (alias: `task-list`) or `gitmap rc heatmap tasks` | `gitmap heatmap query "SELECT TaskSlug, TaskName, RootSpecFilePath, TaskFolder, Status, CreatedAt, UpdatedAt FROM HeatmapParentTask ORDER BY UpdatedAt DESC LIMIT 20"` |
| **Inspect Specific Action ID** | `gitmap heatmap view HeatmapActionLog --id <id> --yaml` or `gitmap heatmap search --id <id>` | `gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, CodeSnippet, CommandText, ReasoningNotes, ExpectedOutcome, ActualOutcome FROM HeatmapActionLog WHERE HeatmapActionLogId = <id>" --yaml` |
| **Cross-Repo Activity Rollup** | `gitmap rc heatmap status` or `gitmap heatmap rc view HeatmapRepoRegistry --tier root` | `gitmap heatmap rc query "SELECT RepoSlug, TotalActionsRecorded FROM HeatmapRepoRegistry ORDER BY TotalActionsRecorded DESC" --tier root` |

### 15. Native SQLite Database Exploration, REST API Server & Lifecycle Suite
- `gitmap db-read "<path-or-alias>" "<sql>"` — Executes read-only SQL query with request timing.
- `gitmap db-read "<path-or-alias>" ls` — Lists all tables and views with row/column counts.
- `gitmap db-read "<path>" -as <alias>` — Registers database file into persistent alias registry.
- `gitmap db-read <alias> table <t1> [t2] [--columns] [--json]` — Displays schema, data types, and primary keys.
- `gitmap db-add <alias> <table> csv <col1>, <col2>, ...` — Inserts record from CSV values.
- `gitmap db-edit <alias> -table <table> -id <id> csv <col1>, <col2>, ...` — Updates record by ID.
- `gitmap db-record-del <alias> -table <table> -id <id> [-y] [--json]` (alias `gitmap dbrd`) — Deletes single record by primary key.
- `gitmap db-find <alias> -table <t1,t2> --find "<text>" [--columns <c1,c2>] [--json]` — Full-text and wildcard search.
- `gitmap db-listing` (alias `gitmap dbl`) — Lists all registered and system databases with file sizes, table counts, and query metrics.
- `gitmap db-view <alias> -t "t1 as a, t2 as b" --viewname <name> --find "<condition>"` — Joins tables and creates persistent views on the fly.
- `gitmap db-rest "<db-path>" as <alias>` — Registers and binds database alias for REST.
- `gitmap db-rest <alias> <port1,port2>` — Launches REST API daemon probing candidate ports, auto-generates dynamic Swagger UI (`/swagger`) and custom Interactive Database Visualizer (`/ui`).
- `gitmap db-rest <alias1,alias2> <ports>` — Mounts multiple databases into a single unified REST gateway.
- `gitmap db-rm <alias> [--purge] [-y]` (alias `gitmap dbrm`, `gitmap db-delete`) — Unregisters database alias and optionally purges underlying SQLite file.

### 16. Screenshot Asset Extraction, Lightshot/prnt.sc Downloading & Brain Sync
- `gitmap prnt <url> [dest]` — Downloads Lightshot / prnt.sc screenshot, saves to `assets/screenshots/`, and auto-syncs to active Antigravity Brain.
- `gitmap prnt <url1> <url2>...` — Concurrent batch download of multiple screenshot URLs.
- `gitmap prnt <url1> <url2>... -e` (or `--extract-only`, `-p`) — Resolves and displays direct image links without downloading.
- `gitmap prnt <url1> <url2>... --json` — Emits structured Universal JSON response envelope.
- `gitmap prnt <url1> <url2>... -o <dir>` — Saves batch downloads to designated custom directory.

### 17. URL Split-DB Caching & 10-Split Segmented Downloader Suite
- `gitmap url read <url>` (aliases: `gitmap url <url>`, `gitmap read-url <url>`, `gitmap get <url>`) — Universal URL streaming and cache inspection command line interface.
  * **Dedicated Split-DB Caching:** Automatically persists remote HTTP/HTTPS responses in dedicated `repodb/url_cache.db` with WAL mode and 256MB memory mapping (`PRAGMA mmap_size = 268435456`).
  * **Sub-Millisecond Retrieval:** Cache hits stream in $<0.5\text{ms}$ without network roundtrips.
  * **Unbuffered Byte-Faithful Output:** Payloads stream directly to standard output via `output.Raw()`, preserving raw binaries, emojis, and ANSI codes without corruption. Diagnostics and transfer telemetry emit to `output.UIErr()`.
  * **Automatic 10-Split Downloader:** Payloads $\ge 1\text{--}2\text{ MB}$ automatically activate 10-connection parallel Range downloading using host `aria2c -s 10 -x 10 -k 5K` or native Go 10-goroutine Range workers with 5KB chunking.
  * **Singleflight Deduplication:** Parallel worker requests coalesce into a single fetch, eliminating duplicate network bandwidth and write lock contention.
  * **Flags:**
    - `--fresh` (alias `-f`): Bypasses existing cache and forces a fresh network download.
    - `--ttl <duration>`: Overrides expiration TTL (e.g. `1h`, `24h`, `7d`; default: `24h`).
    - `--headers`: Emits HTTP response headers to stderr before the payload.
    - `--json` (alias `-j`): Emits structured JSON envelope with payload metadata, SHA-256 hash, and cache timing.
- `gitmap url ls [--limit <n>] [--json]` — Lists cached URLs in tabular format with short hash, size, content type, access count, age, and status.
- `gitmap url info <url> [--json]` — Displays detailed metadata, ETag, headers, and access audit history for a specific URL.
- `gitmap url purge [--all]` — Evicts expired records (or all records with `--all`) and reclaims disk space.
- `gitmap url rm <url>` — Removes a specific URL entry and its cached blob.

### 18. Rust, Cargo & C Toolchain Management
- `gitmap cargo build [--install]` — Builds active Rust project with auto-toolchain resolution and installation fallback.
- `gitmap cargo test [flags]` — Executes project test suites in enriched environment with exit code normalization.
- `gitmap crg check` — Ultra-fast project check and validation (alias for `cargo check`).
- `gitmap cargo status` — Inspects Rust compiler and Cargo toolchain status.
- `gitmap aum locate llvm` (or `clang`) — Sub-millisecond locator for LLVM/Clang binaries (resolves `LIBCLANG_PATH`).
- `gitmap env set <key> <value>` — Persists environment variables across GitMap executions.
- `gitmap env get <key>` — Reads persisted environment variables.

### 19. Tauri Directory Inspection & Log Diagnostics
- `gitmap lf src-tauri` — Instant sub-millisecond directory and file inventory of Tauri application roots.
- `gitmap aum search "<pattern>" src-tauri [-e <ext>]` — Multi-core streaming live search scoped to Tauri Rust and frontend configuration.
- `gitmap pe -f tauri` — CI/CD and runtime diagnostic log extractor filtered specifically for Tauri build errors.

### 20. Process Termination & Network Port Diagnostics
- `gitmap kill <pid|name|:port> [--force] [--tree]` (alias `gitmap terminate`, `gitmap proc kill`) — Native cross-platform process termination with safety guards, protected IDE/shell checks, and graceful or forced execution:
  * Terminate by PID: `gitmap kill 12345`
  * Terminate process tree: `gitmap kill 12345 --tree` (`-t`)
  * Force immediate termination: `gitmap kill 12345 --force` (`-f`)
  * Terminate process listening on TCP port: `gitmap kill :8080 -f`
  * Terminate process by binary name: `gitmap kill node.exe -f`
- `gitmap ports [port]` — Inspects listening network ports, bound PIDs, and process names.
- `gitmap swap-binary [path] [--kill] [--reopen]` — Safely terminates locking PIDs, backs up the active binary to `.old`, copies the newly compiled binary in place, and reopens/verifies version without locked file crashes. Supports `--from <path>`, `--to <path>`, `--script-only`, and `--json`.

---

## Practical Real-World Command Examples & Workflows

This section provides authoritative, production-tested, copy-pasteable command examples with concrete flags, arguments, and expected behavior across all 18 operational categories. AI agents must refer directly to these concrete templates rather than probing `--help` flags.

### 1. Authentication & GitHub Credential Management
Managing credentials safely without interactive prompts, plaintext leakage, or browser popups during automated turns.

```bash
# Non-interactive token ingestion (pre-validated against GitHub API before saving to global git config)
gitmap login --token "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
gitmap login token "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# Offline token ingestion without network verification (for air-gapped runners or local testing)
gitmap login --token "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" --no-verify

# Non-interactive SSH key profile configuration
gitmap login ssh work -y

# Lock workstation-global Git identity to default AI automation bot
gitmap login set global ai-auto

# Lock repository Git identity to named AI automation bot (e.g. worker-01, muse)
gitmap login set ai-auto worker-01

# Lock repository Git identity to default AI automation bot (agent name omitted)
gitmap login set ai-auto

# Inspect current authenticated user and masked credential origin
gitmap login --status

# List resolved token metadata, scopes, and credential origin from git configuration
gitmap token list

# Interactive OAuth browser flow (manual developer sessions only)
gitmap login --web

# Securely purge stored credentials from global git configuration
gitmap logout
```

### 2. Workspace & Multi-Repository Status
High-speed workspace auditing across 40+ repositories in a single invocation without shell traversal loops.

```bash
# Multi-repository status matrix filtered strictly to dirty working trees
gitmap status --dirty
gitmap st -d

# Check repositories with unpushed local commits waiting for remote delivery
gitmap status --ahead

# Check repositories with incoming remote commits waiting to be pulled
gitmap status --behind

# Structured machine-readable JSON status for automated agents and scripts
gitmap status --json

# Check if active repository remote tracking branch has any incoming updates
gitmap has-any-updates
gitmap hau

# Query default remote branch name (main vs master) across repositories
gitmap branch default
gitmap b def

# Discovers the most recently committed or updated remote branch
gitmap latest-branch
gitmap lb

# Terminal live-refresh status monitor without background network fetching
gitmap watch --json --no-fetch

# Single-line compact status summary across all tracked repositories
gitmap status --compact
```

### 3. Script Execution & Runner Engines
Cross-platform polyglot execution with automatic interpreter resolution from `installation.db`, `-NoProfile` safety, and zero shell quoting breakage.

```bash
# Run file manipulator script with subcommands
gitmap py 03-ai-scripts/03-file-manipulator.py lowercase

# Run Python script with verified interpreter from installation.db
gitmap py "03-ai-scripts/06-cicd-local-runner.py"

# Inline Python expression evaluation with sub-millisecond execution (<15ms)
gitmap py -c "import sys; print(sys.version_info[:2])"

# Cross-platform PowerShell execution with deterministic -NoProfile and UTF-8
gitmap pwsh "Get-Process | Where-Object WorkingSet -gt 500MB"

# Direct inline PowerShell command evaluation
gitmap ps -c "Test-Path .ai-memory/heatmap.db"

# Cross-platform POSIX Bash execution across Windows, Linux, and macOS
gitmap bash "chmod +x scripts/*.sh && ./scripts/check.sh"

# Direct inline Bash command evaluation
gitmap sh -c "uname -sm"

# Runner dry-run preview with JSON telemetry output
gitmap py "03-ai-scripts/14-version-sync-checker.py" --dry-run --json
```

### 4. Storage, Cache & Secrets Offloading
Preventing repository pollution by routing scratch scripts, diagnostic fixtures, and secrets to dedicated vaults (`repo-secrets` / `repo-cache`).

```bash
# Offload sensitive environment file into repo-secrets and auto-push
gitmap rs file .env --repo coding-guidelines

# Offload private directory into repo-secrets
gitmap rs folder private/keys --repo coding-guidelines

# Write secret string directly to repo-secrets vault
gitmap rs text "ghp_secret_pat_value" --slug github-token

# Enumerate stored secret files and metadata in JSON format
gitmap rs ls --json

# Offload temporary PowerShell diagnostic script to repo-cache (rc)
gitmap rc file debug-runner.ps1 --repo coding-guidelines

# Offload fixture directory to repo-cache
gitmap rc folder test/fixtures --repo coding-guidelines

# Write reusable temporary test harness directly to repo-cache
gitmap rc text "$env:DEBUG=1; go test -v ./..." --slug test-debug --ext .ps1

# Fast navigation shortcuts to vault locations
gitmap cd rs
gitmap cd rc
```

### 5. High-Performance Automation (AUM) & File Discovery
Replacing slow PowerShell `Get-ChildItem`, `Select-String`, and grep with multi-core streaming, lazy regex, and SQLite hot caches.

```bash
# Multi-core streaming live code search with extension and directory scoping
gitmap aum search "RunFleetPASCommand" cli -e .go

# Named-flag search with explicit root directory and query (ideal for subagent tool calls)
gitmap aum search --root "cli" --query "cmdpipeline"
gitmap aum search --query "RunFleetPASCommand" --dir cli -e .go

# Case-insensitive regex search with boundary patterns
gitmap aum search -r "(\"pas\"|\"pa\")" cli/cmd -e .go -i

# Sub-millisecond directory listing (instant replacement for Get-ChildItem)
gitmap lf "02-spec/21-app"

# Fast glob file discovery across 10,000+ files in <10ms
gitmap find "*.md" 02-spec/21-app

# Find files matching substring across entire codebase
gitmap find-files-any "heatmap"
gitmap ffa "heatmap"

# Find files by filename prefix or suffix
gitmap find-files-startswith "01-architecture"
gitmap find-files-endswith "_test.go"

# Direct zero-disk file content streaming into stdout
gitmap cat "version.json"

# Markdown sequence gap detector and header autofixer
gitmap aum sequence 02-spec/21-app --fix

# Normalize CRLF to LF line endings and strip trailing whitespace
gitmap aum newlines . --fix

# Enforce 500 KB file limit, large JSON exclusion, and binary probes
gitmap aum guard . --max-kb 500 --json

# Exact literal string replacement across files with audit trail
gitmap replace "old_function_name" "new_function_name"
```

### 6. Autonomous CI/CD Self-Healing (Pipeline AI)
Non-polling CI/CD status inspection, calculated ETA timeout waits, bounded failure extraction, and zero-storage artifact cleanup.

```bash
# Dynamic non-polling ETA wait for running GitHub Actions workflow (in seconds)
gitmap pipeline-ai status -t 45
gitmap pipeline-ai status -t 120

# Check pipeline workflow execution state with JSON output
gitmap pipeline-ai status --json

# Extract concise failing step error summary for 4-part Root Cause Analysis (RCA)
gitmap pe -t

# Extract bounded failure frame lines (5 before, 20 after) from latest failed run
gitmap pe -1

# Clear accumulated pipeline error logs non-interactively
gitmap pe clear -y

# Analyze CI/CD pipeline historical trend across recent runs
gitmap pe history-ai

# Purge Actions workflow run artifacts to maintain strict 0.0 GB storage quota
gitmap pipeline purge
```

### 7. Semantic & Flat Commit Suite
Disciplined atomic commit workflow enforcing hyphen-separated messages and remote delivery invariants.

```bash
# Intermediate local checkpoint commit (auto-stages, no remote push, verbatim message)
gitmap commit "wip - draft initial types and interfaces"
gitmap cm "wip - draft initial types and interfaces"

# Feature commit (auto-prefixes 'Feature: ', stages all, and pushes to remote)
gitmap cpf "Skills - enrich gitmap skill with 12 real world categories"

# Bugfix commit (auto-prefixes 'Bug: ', stages all, and pushes to remote)
gitmap cbf "Cloner - correct concurrency flag parsing and worker limits"

# Chore/Maintenance commit (auto-prefixes 'Chore: ', stages all, and pushes)
gitmap cpc "Docs - realign companion skills and remove help flags"

# Release commit (auto-prefixes 'Release: ', stages, and pushes)
gitmap cpr "v2.18.0 - bump minor version and update changelog"

# Pull latest with rebase, stage all, commit, and push in one verified flow
gitmap pcp "Skills - sync cross-platform runner templates"

# Commit with explicit subagent attribution for Split-DB heatmap tracking
gitmap cpf "Module - implement autonomous subtask" --ai-agent "worker-01"

# Apply automated repository remediation (stash, discard, or auto-fix)
gitmap fix all -y

# Safe 2-step git mv file case normalization for lowercase hygiene
gitmap lowercase *.md -y
gitmap lowercase
gitmap lowercase-readme
```

### 8. Multi-Repo & Toolchain Operations
Managing multi-repository workflows, toolchain installation, and environment preparation.

```bash
# Multi-repo pull with compact JSON telemetry
gitmap pae --json

# Inspect Rust and Cargo toolchain status
gitmap cargo status

# Validate Cargo toolchain installation without modifying environment
gitmap install cargo --check

# Install missing toolchain or system runtime
gitmap install cargo

# Discover developer toolchains, profiles, and runtime packages
gitmap install --list
```

### 9. Native Fleet Synchronization & SQLite Agent Task Engine
Compiled native Go task tracking (<1ms WAL mode) and 43-repo synchronization.

```bash
# Synchronize canonical prompts, skills, and scripts across all 43 repositories
gitmap sync --workers 8

# Synchronize a single targeted repository by name
gitmap sync --repo aukgit/core

# Preview fleet synchronization without applying changes
gitmap sync --dry-run

# Initialize agent task tracking in .ai-memory/temp-agents/<slug>/agent-task.db
gitmap agent task init --name "skills-enrichment" --budget 400

# Interrogate task status and active subtask progress in JSON format
gitmap agent task status --json

# Add decomposed subtask bounding box with exclusive file ownership
gitmap agent subtask add --code "ST-01" --slug "enrich-gitmap-skill" --title "Enrich gitmap skill with examples"

# List pending subtasks in current task ledger
gitmap agent task pending

# Claim next available subtask atomically
gitmap agent subtask claim --agent "worker-01"

# Mark claimed subtask completed with verification evidence
gitmap agent subtask complete 1 --agent "worker-01" --evidence "12 categories added and verified"

# Mark subtask failed with explicit diagnostic rationale
gitmap agent subtask fail 1 --agent "worker-01" --reason "syntax error in target spec"

# Perform post-mortem autopsy and crash diagnosis
gitmap agent diagnose
gitmap agent crashed

# Clean ephemeral agent task database after verified completion
gitmap agent temp-clear
```

### 10. File Contention Heatmap, Cognitive Telemetry & RC Sync
Sub-millisecond SQLite WAL tracking of agent decisions, active leases, contention hotspots, and one-shot JSON payload ingestion with multi-AI provenance attribution.

```bash
# Heatmap database discovery catalog across 3 tiers, alias mappings, recent tasks, and suggestion pills
gitmap heatmap

# Repository-filtered heatmap discovery
gitmap heatmap --repo gitmap
gitmap rc heatmap --repo gitmap

# Step 0 Pre-Flight schema and status verification
gitmap heatmap status
gitmap heatmap tables

# One-Shot JSON Ingestion — Antigravity IDE (Gemini Pro)
gitmap heatmap log --json '{"task":"313-gitmap-execution","repo":"gitmap","action":"READ","file":"cli/heatmap/heatmap_cmd.go","state":"EXPLORATION","expected":"Audit flag parsing logic","model":"gemini-2.5-pro","tool":"antigravity-ide"}'

# One-Shot JSON Ingestion — Muse Autonomous Fleet (Muse Engine)
gitmap heatmap log --json '{"task":"313-gitmap-execution","repo":"gitmap","action":"EDIT","file":"cli/cmdssh/ssh_vault.go","state":"EXECUTION","expected":"Add transparent key vault decryption","model":"muse","tool":"muse-runner"}'

# One-Shot JSON Ingestion — Claude Code / Standalone CLI
gitmap heatmap log --json '{"task":"313-gitmap-execution","repo":"gitmap","action":"READ","file":"cli/cmdupdate/update_fallback_loop.go","state":"EXPLORATION","expected":"Verify asset matching order","model":"claude-3-7-sonnet","tool":"claude-code"}'

# Post-Action Outcome Closure (mandatory immediately after operation)
gitmap heatmap update --id 1042 --outcome "Verified one-shot JSON ingestion and model attribution" --met true

# Render hierarchical causal tree of all actions taken during the task
gitmap heatmap tree --task "313-gitmap-execution"
gitmap heatmap tree --task "313-gitmap-execution" --json

# Search prior decisions, reasoning notes, and touched files by slug or query
gitmap heatmap search "313-gitmap-execution"
gitmap heatmap search "ssh_vault" -f "cli/cmdssh/ssh_vault.go"

# Multi-agent file contention analysis and collision blast radius check
gitmap agent heatmap
gitmap agent heatmap --json

# Safe read-only SQL query against cognitive heatmap database
gitmap heatmap query "SELECT HeatmapActionLogId, ActionType, FilePath, ModelName, ToolName, IsExpectationMet FROM HeatmapActionLog ORDER BY HeatmapActionLogId DESC LIMIT 5" --json

# Companion repo-cache (rc) heatmap status, tree traversal, and synchronization
gitmap rc heatmap status
gitmap rc heatmap tree
gitmap rc heatmap live
gitmap rc heatmap sync
```

### 11. Native SQLite Database Exploration, REST API Server & Lifecycle Suite
High-speed native Go SQLite exploration, mutation, join views, and dynamic REST gateways.

```bash
# List all registered and system SQLite databases with sizes and query metrics
gitmap db-listing
gitmap dbl

# List all tables and views in database with column and row counts
gitmap db-read master ls
gitmap db-read ".ai-memory/heatmap.db" ls

# Inspect table schema, data types, and primary key constraints
gitmap db-read tasks table Subtask --columns
gitmap db-read ".ai-memory/heatmap.db" table HeatmapActionLog --columns --json

# Execute read-only SQL query with timing telemetry
gitmap db-read ".ai-memory/heatmap.db" "SELECT COUNT(*) FROM HeatmapActionLog"

# Insert new record into table from CSV values
gitmap db-add tasks Subtask csv "SUB-01","my-task","PENDING"

# Update existing record by primary key
gitmap db-edit tasks -table Subtask -id 42 csv "SUB-01","my-task","COMPLETED"

# Search database across multiple tables using full-text wildcard
gitmap db-find master -table Project,Alias --find "enrichment" --json

# Delete a single record by primary key
gitmap db-record-del tasks -table Subtask -id 42 -y

# Create persistent joined view across tables
gitmap db-view tasks -t "Subtask as s, Task as t" --viewname ViewSubtasks --find "s.TaskId = t.TaskId"

# Launch REST API daemon with auto-generated Swagger UI (/swagger) and Visualizer (/ui)
gitmap db-rest master,tasks 8080,8081
```

### 12. Screenshot Extraction & Lightshot Downloading
Native compiled parallel screenshot asset extraction, CDN resolution, and Antigravity Brain sync.

```bash
# Download single Lightshot/prnt.sc screenshot and auto-sync to Antigravity Brain
gitmap prnt "https://prnt.sc/abcdef"

# Resolve direct image URLs without downloading (dry-run probe)
gitmap prnt "https://prnt.sc/abcdef" -e

# Emits structured Universal JSON response envelope and saves to custom directory
gitmap prnt "https://prnt.sc/abcdef" -o assets/screenshots --json

# Batch download multiple screenshot URLs in parallel
gitmap prnt "https://prnt.sc/url1" "https://prnt.sc/url2" "https://prnt.sc/url3"
```

### 13. Plan & Task Discovery (`gitmap plans`)
Automated plan discovery and inventory across `.ai-memory/plans/` (root, `pending/`, and `completed/`) with Tokyo Night Lipgloss tables and machine-readable exports.

```bash
# Default view: shows last 10 completed plans and top 10 pending plans
gitmap plans
gitmap plan

# Specify custom limit per section (e.g. 20 items per table)
gitmap plans 20
gitmap plans --limit 20
gitmap plans -n 20

# View all plans without truncation
gitmap plans --all
gitmap plans -a

# Filter exclusively to pending and in-progress plans
gitmap plans --pending
gitmap plans -p

# Filter exclusively to completed plans
gitmap plans --completed
gitmap plans -c

# Emit pure machine-readable JSON to stdout (pipeable to jq)
gitmap plans --json
gitmap plans -j 20

# Emit pure machine-readable YAML to stdout
gitmap plans --yaml
gitmap plans --yml
```

### 14. URL Split-DB Caching & 10-Split Segmented Downloader (`gitmap url`)
Dedicated SQLite Split-DB (`repodb/url_cache.db`) caching with sub-millisecond retrieval, byte-faithful `output.Raw()` streaming, and automatic 10-split segmented downloading (`aria2c` or native Go Range workers).

```bash
# Read and stream remote URL with automatic sub-millisecond caching
gitmap url read https://example.com/schema.json

# Shorthand aliases for agent ergonomics
gitmap url https://example.com/schema.json
gitmap get https://example.com/schema.json

# Force fresh network download bypassing existing cache
gitmap url read https://example.com/schema.json --fresh

# Inspect cached payload metadata, ETag, and TTL
gitmap url info https://example.com/schema.json

# List cached URLs in tabular format
gitmap url ls --limit 25

# Purge expired records and reclaim disk space
gitmap url purge

# Ingest heatmap telemetry database into LLM training curriculum
gitmap llm train --heatmap repodb/heatmap.db
gitmap llm train --heatmap https://telemetry.gitmap.internal/api/heatmap/export
```

### 15. Rust, Cargo & C Toolchain Management
Native compilation, automated toolchain installation fallback, fast syntax checking, and dynamic LLVM/Clang path resolution for bindgen.

```bash
# Compile Rust project with automatic toolchain discovery and installation fallback
gitmap cargo build
gitmap cargo build --install

# Execute crate test suites in enriched environment with normalized exit codes
gitmap cargo test
gitmap cargo test -- --nocapture

# Fast syntax and type checking without producing binary artifacts (<1s)
gitmap crg check

# Inspect Rust compiler version, active toolchain targets, and Cargo state
gitmap cargo status

# Sub-millisecond locator resolving LLVM/Clang binary directories (resolves LIBCLANG_PATH)
gitmap aum locate llvm
gitmap aum locate clang

# Persist LIBCLANG_PATH in GitMap managed environment registry across sessions
gitmap env set LIBCLANG_PATH "<path-to-llvm-bin>"
gitmap env get LIBCLANG_PATH

# Scoped execution injecting LIBCLANG_PATH for a single cargo test run
gitmap pwsh "$env:LIBCLANG_PATH='<path-to-llvm-bin>'; gitmap cargo test"
```

### 16. Tauri Workspace Inspection & Build Diagnostics
Sub-millisecond directory inventories and targeted pipeline error log filtering for Tauri desktop applications.

```bash
# Instant directory file listing of Tauri workspace without PowerShell Get-ChildItem
gitmap lf src-tauri
gitmap lf src-tauri/src

# Multi-core streaming live search scoped to Tauri Rust source files
gitmap aum search "tauri::Builder" src-tauri -e .rs

# Filter pipeline error logs specifically for Tauri compilation and packaging errors
gitmap pe -f tauri
```

### 17. Process Termination, Port Unlocking & IDE Health
Cross-platform process management, network port inspection, and safe editor lock resolution without raw taskkill.

```bash
# Inspect listening network ports, bound PIDs, and process names
gitmap ports
gitmap ports 1420
gitmap ports 8080

# Native cross-platform process termination with safety guards (blocks self-kill and protected IDEs/shells unless --force)
gitmap kill 12345
gitmap kill 12345 --tree
gitmap kill 12345 --force
gitmap kill :8080 --force
gitmap kill node.exe --force

# Safely swap locked running binary on Windows (generates detached .ps1, kills locking PIDs, backs up to .old, copies, verifies version, and reopens)
gitmap swap-binary --from .\dist\gitmap.exe --to C:\Users\Administrator\go\bin\gitmap.exe --kill --reopen
gitmap swap-binary .\dist\gitmap.exe --kill

# Safely swap running binary on Linux / macOS (generates detached .sh supervisor script, halts processes, swaps executable, and verifies)
gitmap swap-binary --from ./bin/gitmap --to /usr/local/bin/gitmap --kill --reopen
gitmap swap-binary ./bin/gitmap --kill

# Inspect generated detached swap script without executing
gitmap swap-binary --from ./dist/gitmap.exe --script-only
```

---

## Skill Activation Bands & Guardrails (When to Use GitMap)

### ✅ When to Activate This Skill (Green Band)
- Operating inside any workspace or repository containing `.git` managed by or integrated with GitMap.
- Multi-agent swarm orchestration (`execute-parent-task-with-n-steps-v6`, `execute-parent-task-with-n-steps-v7`, `execute-with-e2e`).
- Performing rapid file discovery (`gitmap lf`, `gitmap find`) or streaming regex code search (`gitmap aum search`).
- Setting up AI automation bot credentials (`gitmap login set global ai-auto`, `gitmap login set ai-auto <agent-name>`).
- Staging, formatting, and pushing atomic semantic commits (`gitmap cpf`, `gitmap cbf`, `gitmap cpc`, `gitmap cfr`, `gitmap cpr`).
- Inspecting CI/CD pipeline health and error telemetry (`gitmap pipeline-ai`, `gitmap pe -t`).
- Querying, mutating, or serving SQLite Split-DB architectures (`gitmap db-read`, `gitmap db-find`, `gitmap db-view`).
- Downloading screenshot assets (`gitmap prnt`) or fetching remote URLs with caching (`gitmap url read`).
- Safe process termination (`gitmap kill`) and binary rotation (`gitmap swap-binary`).

### ⛔ When NOT to Use / Out of Bounds (Red Band)
- Subagents must NEVER execute raw `git` commands (`git commit`, `git add`, `git push`, `git checkout`); only the lead orchestrator invokes GitMap semantic commit verbs.
- Routine builds and tests: Do NOT run full test suites (`go test ./...`, `npm test`, `pytest`) or full project builds (`go build`, `npm run build`, `cargo build`) unless explicitly commanded by the user or executing `execute-with-e2e`.
- Do NOT use slow prohibited legacy tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`, `Get-ChildItem`, `dir`, `ls`).
- Do NOT run inline Python SQLite scripts (`python -c "import sqlite3..."`).
- Do NOT use absolute file paths or `file:///` URIs; use strictly relative Git paths.
- Do NOT invent CLI flags; verify commands before execution.

---

## Operational Guardrails & Non-Negotiable Invariants

1. **Learning & Skill Acquisition:** Run `gitmap llm train` (or `gitmap train --loop`) to initialize or verify GitMap skills. Run `gitmap llm train muse --days 7` when performing Muse memory IO synchronization. Never run broad keyword searches like `gitmap aum search "train"` to discover how commands work.
2. **Mandatory Pre-Flight Pull via GitMap:** Always run `gitmap pull` (fallback `git pull`) before modifying code.
3. **Scoped Search:** Always provide target directories and extensions to `gitmap aum search` (e.g. `gitmap aum search "target" cli --ext .go`).
4. **File Size & Binary Guard:** Respect 500 KB limit (Rule R19); never commit test binaries or temp artifacts.
5. **Coding Guidelines:** Max 8–15 lines per function, single return types with `*appfault.AppError`, affirmative booleans (`isReady`, `hasCache`).
6. **Script Offloading:** Save all temporary diagnostics to `repo-cache` via `gitmap rc` to prevent dirty working trees.
7. **Strict Relative Git Paths:** All paths and references must be relative to repository root (`02-spec/...`, `.ai-memory/...`); zero absolute paths and zero `file:///` URIs.
8. **Zero Storage Ban:** Zero uploads to `actions/upload-artifact`. Maintain 0.0 GB Actions storage quota across all repositories.
9. **Total Ban on Python SQLite Scripts & Ad-Hoc Database Queries:** NEVER run `python 03-ai-scripts/46-agent-sqlite-task-manager.py` or ad-hoc scripts like `python -c "import sqlite3..."`. Ad-hoc Python SQLite usage is AGGRESSIVELY PROHIBITED. You MUST use native compiled GitMap database commands exclusively: `gitmap db-find` / `gitmap dbf` for pattern-matching and regular expression searches across tables and columns; `gitmap db-read` / `gitmap dbr` for read-only table/SQL inspection; `gitmap db-listing` / `gitmap dbl` to discover registered Split-DB tiers; `gitmap db-view` / `gitmap dbv` for join analysis; `gitmap db-add` / `gitmap db-edit` / `gitmap db-rm` for safe mutations; and `gitmap agent` subcommands for task tracking.
10. **Total Ban on Ad-Hoc Python Web Scrapers:** NEVER run `python -c "import urllib.request..."` or write ad-hoc Python scripts (`fetch_urls.py`, etc.) to scrape web pages, extract images, or download screenshot URLs (e.g., `prnt.sc`, `image.prntscr.com`, `i.imgur.com`). You MUST use native compiled `gitmap prnt <url1> [url2...] [--extract-only] [--json]` (alias `gitmap download-prnt`), which automatically handles User-Agent headers, multi-tier CDN resolution, and bidirectional Antigravity Brain synchronization. See `@[gitmap](file;.agents/skills/gitmap)` for full command reference.
11. **Total Ban on Raw Web Fetchers & Single-Stream Downloads:** NEVER invoke `curl`, `wget`, `Invoke-WebRequest`, or custom HTTP scripts to read remote files. All remote URLs MUST be fetched via `gitmap url read <url>` (aliases: `gitmap url`, `gitmap get`), which provides dedicated Split-DB caching (`repodb/url_cache.db`), sub-millisecond retrieval (<0.5ms), unbuffered byte-faithful streaming via `output.Raw()`, and automatic 10-split segmented downloading (`aria2c` or native Go Range workers) for files $\ge 1\text{--}2\text{ MB}$.
12. **Total Ban on Filesystem Grepping for Past Agent Decisions:** NEVER grep source files or guess past subagent actions. You MUST run `gitmap heatmap` or `gitmap rc heatmap --repo <alias>` to discover multi-tier SQLite heatmap databases, cross-referenced repository aliases, and recent tasks across 2–5 repositories with interactive suggestion pills.
13. **Total Ban on Raw Cargo Invocations & Ad-Hoc Toolchain Variables:** NEVER run raw `cargo build`, `cargo test`, `cargo check`, or hardcode `$env:LIBCLANG_PATH` manually. You MUST use `gitmap cargo build [--install]`, `gitmap cargo test`, `gitmap crg check`, and dynamic resolution via `gitmap aum locate llvm` and `gitmap env set LIBCLANG_PATH <path>`.
14. **Total Ban on Raw Process Termination:** NEVER execute raw `taskkill /F /PID <pid>`, `Stop-Process`, or `kill -9` directly. You MUST use native compiled `gitmap kill <pid|name|:port> [--force] [--tree]` (aliases: `gitmap terminate`, `gitmap proc kill`) or `gitmap ports [port]`.
15. **Total Ban on Raw Locked Binary Operations:** NEVER attempt manual binary replacement on locked, running executables using PowerShell `Move-Item`, `Copy-Item`, raw `taskkill` + copy, or shell copy scripts that cause locked file crashes or sharing violations (`ERROR_SHARING_VIOLATION` / `0x20`). You MUST use `gitmap swap-binary [--from <path>] [--to <path>] [--kill] [--reopen]` (or alias `gitmap swap-bin`) to safely generate a platform-specific runner script, gracefully terminate locking PIDs, rotate the active binary to `.old`, install the new binary, verify version integrity, and restart.
16. **Total Ban on Raw Tauri Directory Inspection:** NEVER inspect Tauri roots with raw shell `dir` or `Get-ChildItem`. You MUST use `gitmap lf src-tauri`, `gitmap aum search "<pattern>" src-tauri`, and `gitmap pe -f tauri`.
17. **Total Ban on Raw Git Pull, Raw Git Commits & Misusing Commit Verbs:** NEVER execute raw `git pull` without trying `gitmap pull` first. NEVER execute raw `git commit`, `git commit -m`, `git add -A`, or raw `git push`. NEVER misuse commit verbs (e.g. NEVER use `cpf` for releases or chores). Use `gitmap cpr` for releases, `gitmap cpc` for chores, `gitmap cbf` for bug fixes, `gitmap cfr` for refactors, and `gitmap cpf` for features. NEVER nest conventional prefixes (like `chore(release):`, `feat(...)`, `fix(...)`) inside the GitMap argument as GitMap provides the prefix automatically. Use `<module> - summary` format with a hyphen and NO square brackets `[...]` or colons in the argument. Module must be Title Case or lowercase representing the subsystem; TOTAL BAN on `commit` / `commits` as a module name.
18. **Total Ban on Anonymous AI Autonomous Commits (Mandatory AI Name & On-The-Fly Teaching):** AI agents operating autonomously MUST provide their agent name when setting credentials (`gitmap login set ai-auto <agent-name>`) and when committing (`gitmap cpf/cpb/cpc/cfr/cpr --ai-agent <name>`). Omitting the agent name triggers on-the-fly teaching diagnostics showing valid syntax (`muse`, `gemini`, `cursor-model`, `antigravity`, `claude`). Machine-wide generic bots use `gitmap login set global ai-auto`. Default bot parameters can be configured via CLI (`gitmap settings set ai_auto.default_agent <name>`) or the Web Settings UI (`Speed & Automation` -> `ai-auto`).
