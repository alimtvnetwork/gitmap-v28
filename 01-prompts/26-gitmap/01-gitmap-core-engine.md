> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI) & Cursor IDE
> Category: 26-gitmap
> Invoke: /gitmap
>
> **Top-Instruction Priority Mandate (Above Precedence):**
> Whatever instructions, constraints, domain rules, or user prompts are given ABOVE this prompt (or passed as leading task input) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE.
> This prompt provides the autonomous execution engine; top instructions provide the mission authority.

# GitMap Core Engine & Autonomous Developer Automation — Canonical Specification (must follow)

> **Antigravity Slash Command Compatibility:**
> Use `[/goal](slashCommand;goal)` to run long-running execution without stopping until verified.
> Use `[/learn](slashCommand;learn)` to persist learned architectural conventions and rules.

You are pair programming with **GitMap**, the ultra-fast developer companion and autonomous CLI engine designed for AI coding agents and software engineers.

- **Lead Architect & Author:** MD ALIM UL KARIM (alimtvnetwork)
- **Sponsored By:** RISEUP ASIA LLC (https://riseup-asia.com)
- **Core Mission:** High-performance polyglot repository management, zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures, and AI agent pair programming.

---

## 1. Core AI Directives & Absolute Bans

1. **TOTAL BAN on Shell Search Tools:**
   - NEVER execute `Select-String`, `Get-ChildItem -Recurse`, `rg`, `ripgrep`, `grep`, `git grep`, or `findstr`.
   - All code searches, file discoveries, and symbol queries MUST be executed exclusively via GitMap commands (`gitmap aum search`, `gitmap search`, `gitmap find`, `gitmap cat`).
2. **Mandatory Scoping:**
   - Always scope `gitmap aum search` with a target directory (e.g. `cli`, `pkg`, `02-spec`) and file extension (e.g. `-e .go`, `-e .ts`, `-e .md`) to avoid scanning irrelevant assets.
3. **Temporary Script Offloading to `repo-cache` (`gitmap rc`):**
   - NEVER leave temporary diagnostic scripts, harnesses, or one-off PowerShell / Python files in repository worktrees.
   - Offload and persist all temporary scripts into `repo-cache` (`repo-storage`) via `gitmap rc file <script>` or `gitmap rc text "<content>"` so any repository can reuse them cleanly.
4. **Toolchain Location Caching (`gitmap aum locate python`):**
   - Utilize `gitmap aum locate python` to discover and cache Python runtimes in the local SQLite engine (`installation.db`) to avoid repetitive PATH resolution overhead.
5. **On-the-Fly Execution (`gitmap py` / `gitmap ps`):**
   - Execute on-the-fly Python logic with `gitmap py "<code-or-script>"` and PowerShell commands with `gitmap ps "<cmd>"` (or `gitmap pwsh "<cmd>"`).
6. **Semantic Hyphen-Separated Atomic Commits:**
   - In commit messages, GitMap automatically prefixes `Feature: ` or `Bug: `. NEVER provide a colon inside the commit message argument. Format messages with hyphens: `gitmap cpf "<module> - <summary>"` and `gitmap cpb "<module> - <summary>"`.
7. **Strict Relative Paths Mandate:**
   - All file references, links, and citations MUST be relative paths starting from the repository root (e.g. `01-prompts/26-gitmap/readme.md`). Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well (TOTAL BAN on absolute filesystem paths and `file:///` URIs).

---

## 2. The 9 Core Capabilities: Comprehensive Command Reference Table

| # | Capability | Primary CLI Invocations | Technical Mechanism | AI Operational Benefit |
| :---: | :--- | :--- | :--- | :--- |
| **01** | **Multi-Core Streaming Regex Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`<br>`gitmap aum grep` | Multi-threaded parallel disk scanner with lazy regex compilation streaming directly to stdout. | Eliminates slow shell grep; streams matching lines in milliseconds across thousands of files. |
| **02** | **Indexed Global Symbol Search** | `gitmap search "<query>" [--limit <n>]` | Sub-millisecond SQLite query against DH2D split-db symbol index (`repodb`). | Instant lookup for known functions, types, constants, and classes before invoking disk scans. |
| **03** | **Rapid File Finding** | `gitmap find "<pattern>" [-ext <ext>]`<br>`gitmap ff <name>` / `gitmap ffa <str>`<br>`gitmap ffs <prefix>` / `gitmap ffe <suffix>` | Fast directory indexer matching exact names, substrings, prefixes, or suffixes in <10ms. | Instant path resolution across repos with >10,000 files without filesystem traversal noise. |
| **04** | **File Inventory & Listing** | `gitmap lf [dir]`<br>`gitmap list-files [dir] [-ext <ext>]` | Directory inventory emitting clean relative paths matching patterns. | Clean, compact list of relative paths without PowerShell formatting clutter. |
| **05** | **Deterministic Terminal File Streaming** | `gitmap cat <filepath>` | Direct zero-disk stream of file content into process stdout. | Inspect file content immediately in resource-constrained CLI sessions without buffer bloat. |
| **06** | **Script Execution & Cache Offloading** | `gitmap py "<code-or-script>"`<br>`gitmap ps "<cmd>"` / `gitmap pwsh "<cmd>"`<br>`gitmap rc <file\|folder\|text>` | Cross-platform runners (`-NoProfile`) coupled with auto-committing `repo-cache` storage. | Create, execute, and persist test harnesses and diagnostic scripts without dirtying worktree. |
| **07** | **Python Toolchain Locator & Cache Backup** | `gitmap aum locate python`<br>`gitmap aum cache status` | Sub-15ms toolchain scanner caching absolute executable path into SQLite `installation.db`. | Ensures persistent, deterministic Python binary resolution avoiding PATH ambiguities. |
| **08** | **LLM Training Curriculum & Pipeline Diagnostics** | `gitmap llm train` / `gitmap ld`<br>`gitmap pe` / `gitmap pe history-ai` | 4-stage chained curriculum, markdown docs matrix, and pipeline error log extraction. | Instant blind-AI onboarding and automated CI/CD log extraction for 4-part RCA remediation. |
| **09** | **Semantic Hyphen-Separated Atomic Commits** | `gitmap cpf "<module> - <summary>"`<br>`gitmap cpb "<module> - <summary>"`<br>`gitmap cpr "<module> - <summary>"` | Automated git staging, semantic prefixing, atomic commit, and upstream branch push. | Standardizes commit history; prevents colon syntax corruption in automated release pipelines. |

---

## 3. Deep-Dive Specification of Core Capabilities

### Feature 1: Multi-Core Streaming Regex Search (`gitmap aum search`)
- **CLI Commands:**
  - `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
  - Alias: `gitmap aum grep "<pattern>" [dir]`
- **Rules of Engagement:**
  - `[dir]` must be passed whenever the target subsystem is known (e.g. `pkg`, `cli`, `02-spec`).
  - `-e <.ext>` must be specified to restrict scanning to relevant file types (e.g. `-e .go`, `-e .ts`, `-e .md`).
  - Flags `-i` (case-insensitive) and `-r` (recursive) can be passed as needed.

### Feature 2: Indexed Global Symbol Search (`gitmap search`)
- **CLI Commands:**
  - `gitmap search "<query>" [--limit <n>]`
- **Rules of Engagement:**
  - Query symbol name directly before initiating disk scans.
  - Returns file paths, line numbers, and symbol definitions indexed in the split-db SQLite hot cache.

### Feature 3: Rapid File Finding (`gitmap find`, `gitmap ff`, `gitmap ffa`)
- **CLI Commands:**
  - `gitmap find "<pattern>" [-ext <ext>]` — Wildcard pattern search.
  - `gitmap ff <name>` (alias: `gitmap find-files`) — Exact filename match.
  - `gitmap ffa <substring>` (alias: `gitmap find-files-any`) — Substring filename match.
  - `gitmap ffs <prefix>` (alias: `gitmap find-files-startswith`) — Prefix filename match.
  - `gitmap ffe <suffix>` (alias: `gitmap find-files-endswith`) — Suffix filename match (e.g. `_test.go`).

### Feature 4: File Inventory & Listing (`gitmap lf`)
- **CLI Commands:**
  - `gitmap lf [dir]` (alias: `gitmap list-files [dir]`)
- **Rules of Engagement:**
  - Emits clean relative file paths for directory discovery.
  - Replaces verbose directory listings and shell traversal loops.

### Feature 5: Deterministic Terminal File Streaming (`gitmap cat`)
- **CLI Commands:**
  - `gitmap cat <filepath>`
- **Rules of Engagement:**
  - Directly streams content to stdout.
  - Avoids spawning heavy editors or shell pagers (`cat`, `Get-Content`).

### Feature 6: Script Execution & Cache Offloading (`gitmap py`, `gitmap ps`, `gitmap rc`)
- **CLI Commands:**
  - Python execution: `gitmap py "<code-or-script>"`
  - PowerShell execution: `gitmap ps "<cmd>"` or `gitmap pwsh "<cmd>"`
  - Bash execution: `gitmap bash "<cmd>"` or `gitmap sh "<cmd>"`
  - Offload file to `repo-cache`: `gitmap rc file <filepath> [--repo <repo-name>]`
  - Offload folder to `repo-cache`: `gitmap rc folder <folderpath> [--repo <repo-name>]`
  - Offload inline code to `repo-cache`: `gitmap rc text "<content>" --slug <slug> --ext <.ps1|.py>`
- **Rules of Engagement:**
  - Use `repo-cache` (`gitmap rc`) for all temporary diagnostics, preventing dirty working trees while preserving reusable harnesses across all repositories.

### Feature 7: Python Toolchain Locator & Cache Backup (`gitmap aum locate python`)
- **CLI Commands:**
  - Locate and cache runtime: `gitmap aum locate python`
  - Inspect cached toolchains: `gitmap aum cache status`
- **Rules of Engagement:**
  - Runs in under 15ms and writes the resolved runtime executable into local SQLite `installation.db`.
  - Avoids ambient PATH discovery issues across different shell profiles.

### Feature 8: LLM Training Curriculum & Pipeline Diagnostics (`gitmap llm train`, `gitmap pe`)
- **CLI Commands:**
  - Run full curriculum: `gitmap llm train` (alias: `gitmap llm chain`)
  - Stdout-only curriculum: `gitmap llm train --text-only`
  - View markdown reference matrix: `gitmap ld` (alias: `gitmap llm-docs`)
  - Inspect pipeline failures: `gitmap pe` (alias: `gitmap pipeline errors`)
  - Analyze CI/CD pipeline run history: `gitmap pe history-ai`
- **Rules of Engagement:**
  - AI agents use `gitmap pe` to extract bounded error logs and run 4-part Root Cause Analysis (RCA) without manual dashboard navigation.

### Feature 9: Semantic Hyphen-Separated Atomic Commits (`gitmap cpf`, `gitmap cpb`, `gitmap cpr`)
- **CLI Commands:**
  - Feature commit & push: `gitmap cpf "<module> - <summary>"`
  - Bugfix commit & push: `gitmap cpb "<module> - <summary>"`
  - Release / chore commit & push: `gitmap cpr "<module> - <summary>"`
- **Rules of Engagement:**
  - NEVER provide a colon (`:`) inside the message argument.
  - Correct: `gitmap cpf "engine - add streaming regex search"`
  - Incorrect: `gitmap cpf "Feature: engine: add streaming search"`

---

## 4. Operational Guardrails for AI Agents

1. **Pre-Flight Pull:** Always pull before starting work: `git pull origin <branch> --no-rebase`.
2. **File Size Limit (500 KB):** Never read, generate, or stage files exceeding 500 KB without prior chunking or exclusion.
3. **Positive Boolean Hygiene:** Use affirmative booleans only (`isReady`, `hasCache`, `isValid`). Ban explicit equality against booleans (`if isReady == true`).
4. **Zero-Storage Actions:** GitHub Actions workflows must NEVER upload test or build artifacts to Actions storage.
5. **Strict Relative Git Paths:** Write relative paths from repository root everywhere in code, specs, plans, and commits.
