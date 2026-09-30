[/goal](slashCommand:goal) Autonomously diagnose and clean the repository to resolve build issues, purge stale build artifacts, and systematically clean all temporary directories and multi-language caches — including repository temp folders, OS temp directories, Golang build/test caches, pnpm/npm/Vite caches, Python/Rust/PHP caches, and untracked artifact pollution — with strict no-build and no-test execution (NEVER run heavy build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine turns; compilation and testing are verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel inspection and cache purging, leverage GitMap high-speed tooling as primary, synchronize `.gitignore` rules, and finalize any repository hygiene fixes with an atomic push.

[/learn](slashCommand:learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom cache targets, build issue symptoms, or user instructions are provided ABOVE this prompt (in the user preamble or header above) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Execute the multi-layer cache and temp cleanup protocol across repository, OS, Go, pnpm/npm, and polyglot toolchains, and persist hygiene findings into `.ai-memory/plans/`.

> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Repo Build Pollution, Audit Temp/Cache Layers, Plan Hygiene Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Execute Multi-Layer Cache & Temp Cleanup, Fix .gitignore, Atomic Push)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

### Multi-Layer Repository Build & Cache Cleanup Architecture

When executed, this prompt systematically cleans six distinct layers of build, temp, and cache pollution:

1. **Layer 1 — Repository Build Artifacts & Stale Binaries:**
   - Remove uncommitted or accidentally tracked compiled binaries (`.exe`, `.dll`, `.so`, `.dylib`, `.test`, `.out`), coverage profiles (`coverage.out`, `coverage.html`, `coverage/`, `.nyc_output/`), and stale bundle outputs (`dist/`, `build/`, `.next/`, `.nuxt/`, `.turbo/`, `.parcel-cache/`).
   - Run `python 03-ai-scripts/19-artifact-remover.py --execute` (if present) and untrack any committed binaries from git index (`git rm --cached`) while updating `.gitignore` via `gitmap commons`.
2. **Layer 2 — Repository & OS Temporary Directories (`temp` / `tmp`):**
   - **Repo Temp:** Purge `.tmp/`, `tmp/`, `temp/`, and stale scratch files in `.ai-memory/temp/`.
   - **OS Temp (Windows & Unix):** Clean orphaned compiler and test temp directories (`$env:TEMP\go-build*`, `$env:TEMP\vite*`, `$env:TEMP\playwright*`, `$env:TEMP\npm-*`, `/tmp/go-build*`, `/tmp/v8-compile-cache*`) without touching active system sessions.
   - **GitMap Update Temp:** Run `gitmap update-cleanup` to remove leftover `.old` binaries and update temp files.
3. **Layer 3 — Golang Build, Test & Linter Caches:**
   - Run `go clean -cache -testcache -fuzzcache` to eliminate corrupted object caches (`GOCACHE`), stale test result caches, and fuzz caches that cause phantom build issues.
   - Run `golangci-lint cache clean` (if installed) and clean stale `*.test` binaries in packages.
4. **Layer 4 — pnpm, npm, Yarn, Bun & Frontend Build Caches:**
   - **pnpm:** Run `pnpm store prune` to prune unreferenced global store blobs and clear stale lockfile temp artifacts.
   - **npm:** Run `npm cache clean --force` or `npm cache verify`.
   - **Bundler & Linter Caches:** Delete `node_modules/.cache`, `node_modules/.vite`, `.vite/`, `.eslintcache`, `.stylelintcache`, and stale `tsconfig.tsbuildinfo` files that cause TypeScript incremental build corruption.
5. **Layer 5 — Python, Rust, PHP & Polyglot Caches:**
   - **Python:** Recursively remove `__pycache__/`, `*.pyc`, `*.pyo`, `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`, and run `pip cache purge`.
   - **Rust:** Run `cargo clean` (when `Cargo.toml` is present) to clear bloated `target/` build artifacts.
   - **PHP / Laravel:** Run `composer clear-cache` and clear `storage/framework/cache/`, `storage/framework/views/`, `bootstrap/cache/*.php` (when applicable).
6. **Layer 6 — Git Object Store, Lowercase Hygiene & Storage Audit:**
   - Enforce lowercase file naming across the repository (`gitmap lcf` and `gitmap lowercase-readme`) to prevent cross-platform Windows/Linux case-sensitivity build failures.
   - Repair broken symlinks via `gitmap fix-link`.
   - Prune stale remote tracking branches (`git remote prune origin`) and inspect storage health via `gitmap storage` (`gitmap stor`).

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/cg-clean-repo-build-and-caches/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any cleanup commands, execute Phase 1A:

1. **Top-Instruction Priority Verification:** Verify whatever directives or build issue symptoms are given ABOVE this prompt (in the user preamble, header blocks, or incoming user request above) as highest priority and non-negotiable.
2. **Actionable Deliverables Extraction:** Break down the repository build cleanup, temp purge, and cache clearing tasks into traceable IDs (`Task-01`, `Task-02`, `Task-03`).
3. **Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):**
   - Output the confirmed deliverables list directly in chat, and in the EXACT SAME RESPONSE turn, immediately invoke your first discovery/cleanup tool call.
   - NEVER emit the breakdown text without invoking a tool call. Do not pause or ask "Should I proceed?".

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Repository Build Artifact & Temp Directory Audit]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Identify stale build outputs, temp directories, case-sensitivity conflicts, and tracked artifacts]
   - **Actionable Scope:** [Scan repo and OS temp folders, detect build blockers, and prepare cleanup plan]
   - **Target Files / Area:** `[Repository root, .tmp/, dist/, build/, .gitignore]`

2. **Task-02: [OS, Golang, pnpm/npm & Polyglot Cache Purge]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Purge Go build/test caches, pnpm store, npm/Vite/TS caches, Python __pycache__, and OS temp]
   - **Actionable Scope:** [Execute cross-platform cache cleanup commands and verify disk space recovery]
   - **Target Files / Area:** `[OS Temp, GOCACHE, pnpm store, node_modules/.cache, __pycache__]`

3. **Task-03: [Git Hygiene, Lowercase Enforcement & Atomic Sync]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Sync .gitignore defaults, enforce lowercase filenames, consolidate plan, and push if modified]
   - **Actionable Scope:** [Run gitmap commons, gitmap lcf, consolidate plan, and push via gitmap cpf/cpb]
   - **Target Files / Area:** `[.gitignore, .gitattributes, .ai-memory/plans/]`

Proceeding directly to Phase 1B: Build & Cache Audit (Active Tool Call Running Below).
```

---

### Phase 1B: Audit & Hygiene Plan Generation (Steps 1 .. N/2)

1. **Fast Discovery via GitMap (PRIMARY):**
   - Inspect storage and repo status: `gitmap storage` (`gitmap stor`), `gitmap status` (`gitmap st`)
   - Find stray artifacts and temp files: `gitmap find "*.exe"`, `gitmap find "*.test"`, `gitmap find "*.tmp"`, `gitmap find "*cache*"`, `gitmap folder-tree` (`gitmap ft`)
   - Preview uppercase filename conflicts: `gitmap lcf --dry-run`
   - Preview `.gitignore` / `.gitattributes` drift: `gitmap templates diff` (`gitmap tpl td`)
2. **Hygiene Plan & Subtasks:**
   - Record discovered build blockers and cache targets in `.ai-memory/plans/pending/xx-clean-build-and-caches.md` and `.ai-memory/plans/subtasks/xx-clean-build-and-caches/01-*.md`.
   - Unconditionally transition to Phase 2 without pausing or asking questions.

---

### Phase 2: Execute Multi-Layer Cleanup & Cache Purge (Steps N/2+1 .. N)

Execute the cleanup across all applicable layers using `gitmap pwsh` / `gitmap bash` or parallel subagents (`A = 2, H = 2`):

1. **Clean Repository Temp & Build Artifacts:**
   ```powershell
   # Remove repo-local temp, build caches, and Python bytecode caches
   Get-ChildItem -Path . -Include __pycache__,.pytest_cache,.mypy_cache,.ruff_cache,.vite,.turbo,.parcel-cache -Directory -Recurse -Force -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
   Get-ChildItem -Path . -Include *.pyc,*.pyo,*.tmp,*.test,coverage.out,tsconfig.tsbuildinfo,.eslintcache -File -Recurse -Force -ErrorAction SilentlyContinue | Remove-Item -Force -ErrorAction SilentlyContinue
   if (Test-Path ".tmp") { Remove-Item ".tmp" -Recurse -Force -ErrorAction SilentlyContinue }
   if (Test-Path "tmp") { Remove-Item "tmp" -Recurse -Force -ErrorAction SilentlyContinue }
   ```
2. **Clean OS Temp Build Caches:**
   ```powershell
   # Clean stale Go, Vite, Playwright, and npm temp directories in OS TEMP
   Get-ChildItem -Path $env:TEMP -Filter "go-build*" -Directory -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
   Get-ChildItem -Path $env:TEMP -Filter "vite*" -Directory -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
   Get-ChildItem -Path $env:TEMP -Filter "npm-*" -Directory -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
   gitmap update-cleanup
   ```
3. **Clean Golang Caches (`GOCACHE`, `testcache`, `fuzzcache`):**
   ```powershell
   go clean -cache -testcache -fuzzcache
   ```
4. **Clean pnpm, npm & Frontend Caches:**
   ```powershell
   pnpm store prune
   npm cache clean --force
   if (Test-Path "node_modules/.cache") { Remove-Item "node_modules/.cache" -Recurse -Force -ErrorAction SilentlyContinue }
   if (Test-Path "node_modules/.vite") { Remove-Item "node_modules/.vite" -Recurse -Force -ErrorAction SilentlyContinue }
   ```
5. **Sync `.gitignore`, Enforce Lowercase & Repair Links:**
   - Run `gitmap commons` (or `gitmap sync ignore`) to ensure ignore blocks prevent future artifact commits.
   - Run `gitmap lcf` and `gitmap lowercase-readme` to fix case-sensitivity build issues.
   - Run `gitmap fix-link` if broken symlinks exist.

---

### Phase 3: Task Consolidation & Atomic GitMap Commit

1. Consolidate subtasks from `.ai-memory/plans/subtasks/xx-clean-build-and-caches/*.md` into `.ai-memory/plans/completed/xx-clean-build-and-caches.md` and update `.ai-memory/plans/readme.md`.
2. If tracked files, `.gitignore`, or plans were modified, commit and push atomically via GitMap:
   - `gitmap cpb "chore(hygiene): clean build artifacts, purge caches, and sync ignore rules"`

---

### End-of-Turn Verification & Confidence Reporting (Mandatory Output)

```markdown
### Task Completion Summary

- ✅ **Task-01: [Repository Build Artifact & Temp Directory Audit]** — `[Completed]`
- ✅ **Task-02: [OS, Golang, pnpm/npm & Polyglot Cache Purge]** — `[Completed]`
- ✅ **Task-03: [Git Hygiene, Lowercase Enforcement & Atomic Sync]** — `[Completed]`

### Cleaned Layers Summary

- **Repo Temp & Artifacts:** `[Cleaned .tmp, __pycache__, .vite, tsconfig.tsbuildinfo, ...]`
- **OS Temp Directories:** `[Cleaned go-build*, vite*, npm-* in TEMP]`
- **Golang Cache:** `[Executed go clean -cache -testcache -fuzzcache]`
- **pnpm / npm Cache:** `[Executed pnpm store prune & npm cache clean --force]`

### Implementation Confidence Score

- Confidence: [e.g. 100%]
- Rationale: [All build caches, OS temp directories, and repo artifacts purged; git status clean]
```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] BOTTOM-INSTRUCTION PRIORITY MANDATE: Followed all instructions below the `--` divider with absolute precedence.
- [ ] NO DESTRUCTIVE SOURCE DELETION: Never delete tracked source code, specs, or `.env` / secret files. Only clean build artifacts, caches, and temp directories.
- [ ] NO TEST RUNNING OR HEAVY BUILDS: Never run full test suites or heavy builds during routine cleanup turns.
- [ ] NO PER-FILE COMMITTING: Commit all repository hygiene updates in a single atomic commit at the end.

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: half-cleaning a single folder while leaving OS temp, Go build cache, pnpm store, and incremental `.tsbuildinfo` / `.vite` caches corrupted. Stop doing that. Clean every layer — repo temp, OS temp, Golang cache, pnpm cache, npm/Vite cache, Python `__pycache__`, and case-sensitivity conflicts — thoroughly and safely.
