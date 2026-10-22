# Find 100 Improvements in the Code Base

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, GitMap Codebase Scanning, Discovery Subagents, Up-Front Brain & Roadmap Creation)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: 10-Field Specification Card Authoring, Micro-Batching Waves, Master Catalog Assembly, Pending Plan Enqueueing)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Retrospective Quality Verification, Relative Path Hygiene, Dynamic Join Saturation, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 1.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /find-100-improvements-in-code-base <codebase-or-path>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously scan, audit, discover, and document exactly 100 high-impact, grounded code improvements across the target codebase: partition the discovery across 10 canonical Coding Guideline Rule Families (10 improvements per family = 100 total), format each improvement into an exhaustive 10-field specification card, compile a Master Catalog Summary Table, enqueue all 100 improvements as actionable, execution-ready pending plans into `.ai-memory/plans/pending/XX-<slug>.md` with modular subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`, and configure the pipeline so each successive execution seamlessly generates the next 100 improvements.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: all discoveries, analysis, and fix proposals must strictly ground themselves in repository coding guidelines (`02-spec/02-coding-guidelines/`) and error management standards (`02-spec/03-error-manage/`). Every single improvement proposal must provide copy-paste ready, refactored code satisfying all guidelines, mapped caller blast radius via `gitmap aum search`, and a concrete verification contract. Never invent hypothetical issues or use placeholder snippets.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution for improved work or less friction, faster work and more pure and effective work. Ensure all 100 improvements are systematically discovered via GitMap streaming queries, categorized into 5–8 file micro-batches, authored into modular subtasks, and registered in pending task backlogs before concluding.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print text like "Dispatched Worker..." and stop. You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the overarching audit roadmap, 100-improvement ledger, Master Catalog Summary Table, and final GitMap commit. Never executes file discoveries solo.
  2. **Branch Coordinators (Family Owners):** Own designated rule families (e.g. Worker 01 owns Families 01–03, Worker 02 owns Families 04–07, Worker 03 owns Families 08–10) and manage discovery micro-batches.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist branch coordinators with GitMap AST/pattern scanning, caller blast radius tracing via `gitmap aum search`, and specification card authoring.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Discovery Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to scan the codebase using GitMap high-speed streaming commands across the 10 rule families and return concrete candidate code locations to the lead orchestrator.
  2. **Card Authoring Step (A = 3 `self` authoring subagents):** Spawn 3 authoring subagents (`TypeName: "self"`) to author modular, disjoint 10-field improvement specification cards and subtask plans (never writing to the same file path).
  3. **Verification & Enqueueing Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to validate relative path hygiene, check against banned rules, compile pending subtasks, and register backlogs.
- **MANDATORY 3 SUBAGENTS IN SINGLE TOOL CALL:** Every call to `invoke_subagent` MUST contain an array of EXACTLY 3 subagent objects: `"Subagents": [agent1, agent2, agent3]`. Calling `invoke_subagent` with fewer than 3 subagents, or executing solo without calling `invoke_subagent`, is an immediate protocol violation and auto-reject failure.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending batch, or commands the worker to dynamically join an active subtask in collaborative assistance mode (parallel pattern searching, caller discovery, card proofreading).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from performing all 100 audits or writing all improvement cards solo without calling the `invoke_subagent` tool.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/find-100-improvements-in-code-base/skill.md` does not exist in the workspace, create it now.
2. If `.cursor/skills/find-100-improvements-in-code-base/skill.md` does not exist in the workspace, create it now with IDE-detectable Cursor skill link syntax.
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this prompt must be verified as highest priority and non-negotiable.
2. **Extreme "Understood" Turn-1 Showcase (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, you MUST output a comprehensive, in-depth understanding block directly in visible chat before any background tool execution. Any generic boilerplate or superficial 1-sentence summaries are protocol auto-reject failures. Turn 1 MUST display the top-level mandate containing all 5 Mandatory Pillars followed by the confirmed task breakdown using the 6-part verification schema:
   * **Pillar 1: Strategic Intent & Systems Reasoning:** Articulate overarching architectural goal, user intent, operational rationale, and systems integration mechanics. Explain *why* finding 100 grounded improvements elevates code health, maintainability, error resilience, and linter compliance.
   * **Pillar 2: Explicit Scope Boundaries:** Strict two-column partition defining In-Scope vs. strictly Excluded or Deferred items.
   * **Pillar 3: Non-Negotiable Constraints & Negative Guardrails Ingestion:** Explicit recitation of active constraints: raw search ban (`rg`, `ripgrep`, `Select-String`, `git grep`, `Get-ChildItem`, `dir`, `ls`), worker git ban (no git commands in subagents), zero builds/tests (never run all tests locally or full build unless explicitly asked; every time code is written, do not build the code; only do end-to-end testing if explicitly asked), strict relative Git paths, lowercase naming, zero secrets in repo (`repo-secrets` offloading), GitMap hyphen-separated commit format without colons.
   * **Pillar 4: Blast Radius & Dependency Topology:** Map impacted subsystems, target files, downstream callers, and prompt/skill indexes with explicit risk containment strategies.
   * **Pillar 5: Multi-Agent Concurrency Topology:** Detail swarm structure (`A = 3`, `H = 2`, Capacity = 6), disjoint file assignments, and immediate dynamic respawn.
   * **Task Comprehension & Model Runtime Invariants:** In the task breakdown, every single task MUST explicitly state `- **Understood / AI Comprehension:** [YES] — [1–2 concise sentences proving deep understanding]` and `- **Azure AI / Model Runtime Context:** [Document Azure AI or local model endpoint context, or state 'N/A — Local CLI/GitMap']`.
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Discrete Deliverables Extraction:** Break down requirements into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).
5. **Mandatory Same-Turn Tool Chaining:** Emit the comprehension mandate and breakdown in chat, and in the **EXACT SAME TURN**, invoke your first tool call. NEVER emit text alone.

```markdown
### 🧠 Deep Comprehension & Understood Mandate

#### Pillar 1: Strategic Intent & Systems Reasoning
- [Articulate architectural goal, audit depth across 10 rule families, and structured enqueueing of 100 improvements into actionable execution plans.]

#### Pillar 2: Explicit Scope Boundaries (In-Scope vs. Excluded / Deferred)
- **In-Scope:** [Strict enumeration of files, modules, rule families, and improvement deliverables generated in this run.]
- **Excluded / Deferred:** [Direct inline execution of code refactoring is deferred to downstream execution runs; this audit focuses strictly on discovery, 10-field specification card authoring, and pending plan enqueueing.]

#### Pillar 3: Non-Negotiable Constraints & Negative Guardrails Ingestion
- [Explicit recitation of active constraints: raw search ban (no rg/ripgrep/Select-String), worker git ban (no git in subagents), zero builds/tests (no go build, npm run build, full test suites), strict relative Git paths, lowercase naming, zero secrets in repo, GitMap hyphen-separated commit format without colons.]

#### Pillar 4: Blast Radius & Dependency Topology
- **Impacted Subsystems & Files:** [Files, packages, downstream callers mapped via gitmap aum search.]
- **Containment & Risk Mitigation:** [Disjoint file ownership per subagent; non-destructive audit and planning.]

#### Pillar 5: Multi-Agent Concurrency Topology & Wave Orchestration
- **Swarm Structure:** Lead Orchestrator (Root Node) -> Branch Coordinators (Family Owners) -> Leaf Workers & Dynamic Assistants
- **Capacity:** A = 3 agents × H = 2 hands = 6 concurrent operations, WAVES = ceil(subtasks / 6)
- **Disjoint Bounding Boxes:** Zero file overlap; immediate dynamic respawn on worker completion
```

---

## 1. Precedence Hierarchy & Scope (Highest First)

1. **User Instructions & Preamble:** Directives and parameters ABOVE this prompt outrank everything below.
2. **Platform Limits:** Native tools, Artifact Review Policy, permission prompts, hooks. Never claim to override them.
3. **Repo Rules:** `agents.md`, `.ai-memory/strictly-avoid.md`, and `coding-guidelines.md`.
4. **This Prompt.**

If sources conflict, follow the stricter one and record under `Conflicts:` in ledger.

### Scope Control Rules
- **Turn 1 Task Showcase:** In your very first response turn, showcase the given task in visible chat.
- **Read budget:** Read only requested paths, search hits, and Step 0 context via GitMap.
- **History read-only:** Never edit past events, changelogs, completed plans, release notes, or archive directories.
- **Out-of-scope:** Log under `Follow-ups:` in plan; never fix in this run.
- **Minimal diff:** Modify only target files required for audit documentation and plan enqueueing.
- **Indexes:** Update `01-prompts/readme.md` and `.ai-memory/prompts.md` when adding/modifying prompts. Update `.ai-memory/plans/readme.md` and `.ai-memory/what-to-read.md`.
- **Mandatory Multi-Agent Partitioning:** Work MUST be partitioned across A workers (`A = 3`). Solo execution is strictly banned.
- **Finish early:** When all 100 improvements are documented, verified, and enqueued, proceed directly to consolidation.
- **Zero releases:** Never bump versions or edit changelogs unless requested (R10).

---

## 2. Core Operational Rules (Cite by ID)

- **R1 Zero Builds, Test Suites & CI/CD Verification Ban (TOTAL BAN).** NEVER run `go build`, `npm run build`, `vite build`, `cargo build`, `go test ./...`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. Every time you write code, you do not build the code. Never run all tests locally or full build unless explicitly asked by the user. Only do end-to-end (E2E) testing based on mentioned test cases if explicitly asked. Do CI/CD verification (e.g. `gitmap pe`, `gitmap pipeline-ai`, GitHub Actions) only if asked; if not asked explicitly, don't do it.
- **R2 Targeted Checks Only.** Run only fast, file-scoped checks on specifically modified files (see Section 11). A check scanning 0 files is a **FAIL**.
- **R3 Evidence or It Did Not Happen.** Every `DONE`, `PASS`, or "verified" claim MUST cite a concrete file path, line number, git diffstat, or command exit code (`exit 0`). Vague assurances are auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths.** Verify commands with a harmless call (`gitmap lf readme.md`), not `--help`. Use documented fallbacks and log in ledger.
- **R5 Mandatory Subagents (`invoke_subagent`), Concurrency Scaling & Immediate Dynamic Respawn.** Spawning subagents via `invoke_subagent` (`A = 3`, `H = 2`, System Concurrency Capacity = 6 concurrent operations) is an **ABSOLUTE MUST** (`research` for discovery in Phase 1, `self` for card authoring and plan enqueueing in Phase 2/3). Solo execution without calling `invoke_subagent` is an auto-reject failure.
- **R6 One Owner Per File (Disjoint Bounding Boxes).** Within every worker wave, each file has exactly one owner. Shared indexes (`.ai-memory/plans/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/what-to-read.md`) belong exclusively to lead.
- **R7 Git Safety & Isolation (Worker Git Ban).** Subagents NEVER run git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces (`Workspace: "inherit"`), worker git calls create `.git/index.lock` collisions that immediately crash parallel agents.
- **R8/R9 Atomic Commit & Push via GitMap (TOTAL BAN ON RAW GIT COMMITS).** The run ends with ONE GitMap call:
  - For features: `gitmap cpf "<module> - <feature summary>"` (e.g. `gitmap cpf "audit - catalog 100 codebase improvements wave 01"`). GitMap automatically prepends `Feature: `, so format with a hyphen `-`; there is NO need to provide a colon in the GitMap argument. DO NOT provide a colon; use a hyphen `-` to separate module/scope from summary. NEVER include `feat(...)` or `feature:` in your message.
  - For bug fixes: `gitmap cpb "<module> - <fix summary>"`.
  TOTAL BAN on raw git commits (`git commit`, `git add -A`, raw `git push`), colons inside the GitMap message argument, and conventional prefixes (`docs(...)`, `feat(...)`, `fix(...)`, `chore(...)`). ZERO intermediate commits during Phase 1 or Phase 2; all files must be committed together at the final step of Phase 3.
- **R10 Zero Unauthorized Releases.** Never bump versions, edit `version.json`, update changelogs, or trigger release scripts unless user explicitly requested release.
- **R11 Strict Relative Git Paths & Lowercase Hygiene.** Strict ban on absolute paths (`C:\...`, `/home/...`) and `file:///` URIs. Only add relative paths starting from git root. All filenames, documentation, and specs strictly lowercase (e.g. `readme.md`, `agents.md`, `skill.md`).
- **R12 No Polling / Immediate Turn Yielding.** When calling `invoke_subagent`, make it the sole tool action at turn end, print progress line (`Dispatched Worker 01 .. Worker <A> (wave k / WAVES); waiting for their results.`) and **STOP CALLING TOOLS**. Never poll in loop. Check `manage_subagents` once if wave runs long.
- **R13 Two-Strike Retry Cap & Anti-Looping.** Tool failing twice: worker replies `STATUS: BLOCKED` with exact error and stops. Lead takes over and logs `LEAD_FALLBACK: <reason>`.
- **R14 100% Ambiguity & Decision Boundaries.** Non-blocking: choose conservative option, log in ledger `Assumptions:`, proceed. Blocking: `ask_question` once, log in `.ai-memory/ambiguous-questions/01-new-ambiguity/`, continue unblocked tasks.
- **R15 Zero Generated Artifacts Committed.** Never commit build caches, logs, temp scripts, or newly generated binaries unless repository already tracked them.
- **R16 Zero Secrets in Standard Repos.** Never write credentials, tokens, passwords, or `.env` contents into tracked files, commits, ledger, plans, specs, or prompts (`agents.md` section 9).

---

## 3. GitMap High-Speed Command Primacy (Run Everything Faster)

GitMap is your **PRIMARY** acceleration engine. Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Execute all searches and operations through GitMap:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Live Streaming Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` | `gitmap aum grep "<pat>"` | Multi-core streaming live file search (replaces `Select-String`, `git grep`) |
| **Indexed Symbol Search** | `gitmap search "<query>" [--limit <n>]` | `gitmap search` | Instant SQLite cached keyword/symbol search across indexed repos |
| **Wildcard File Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core glob filename finder |
| **Directory Inventory** | `gitmap lf <dir>` / `gitmap find "*.*" <dir>` | `gitmap lf [pat]` | Instant indexed repository inventory (replaces `Get-ChildItem`, `dir`, `ls`) |
| **Substring File Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Direct zero-disk memory streaming to stdout (replaces `cat`, `type`, `Get-Content`) |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpf "<module> - <msg>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, formats with prefix, and pushes atomically. Format as hyphen `-` (no colons in GitMap arguments). |
| **Heatmap Action Log** | `gitmap heatmap log --task <task> --action <READ\|EDIT\|CREATE\|DELETE> --file <path> --state <STATE> --expected "<goal>"` | `gitmap heatmap log` | Concurrency-safe cognitive action & reasoning telemetry into `.ai-memory/heatmap.db` |
| **Contention & Leases** | `gitmap agent heatmap [--json]` / `gitmap heatmap status` | `gitmap agent heatmap` | Pre-edit file lease conflict & contention hotspot analysis across active agents |
| **RC Heatmap Suite** | `gitmap rc heatmap <status\|tables\|query\|view\|tree\|live\|sync>` | `gitmap rc heatmap` | Unified RC heatmap commands with automatic vault auto-provisioning |

### 🔍 Code & Symbol Search Protocol (TOTAL BAN ON `rg`, `ripgrep`, `Select-String`, `git grep`, `Get-ChildItem`, `dir`, `ls`)
- **Directory file inventory:** `gitmap lf <dir>` or `gitmap find "*.*" <dir>`
- **Live streaming code/symbol search:** `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
- **Finding files:** `gitmap find "<pattern>" [dir]`
- **Stream file:** `gitmap cat <filepath>`
- **TOTAL BAN (AUTO-REJECT FAILURE):** Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Running raw grep, ripgrep, or slow shell search pipelines wastes execution steps, causes Windows process hangs, and violates GitMap primacy.

---

## 4. The 10 Coding Guideline Rule Families (10 Improvements Each = 100 Total)

To achieve complete, balanced repository-wide audit coverage, exactly 10 improvements MUST be identified and specified under each of the 10 rule families below:

### Family 01: Error Management & Structured `*appfault.AppError` Returns (Items 001–010)
- **Canonical Specifications:** `02-spec/03-error-manage/02-error-architecture/06-apperror-package/02-apperror-reference.md`, `02-spec/03-error-manage/01-error-resolution/03-retrospectives/06-appfault-package-migration.md`.
- **Core Principles:**
  - Functions returning structured error metadata MUST use `*appfault.AppError` as return type (`func doWork() *appfault.AppError`).
  - Total ban on returning bare untyped `error` for internal domain operations.
  - Package naming convention: `pkg/appfault` with struct `AppError` eliminates package stutter (`appfault.AppError` instead of `apperror.AppError`).
  - Zero swallowed errors: every error branch must wrap, log, or return structured metadata with operational fault codes.
  - Standardized universal error response envelope serialization.
- **GitMap Audit Queries:**
  - `gitmap aum search "func.*error\)" -e .go`
  - `gitmap aum search "apperror\.AppError" -e .go`
  - `gitmap aum search "fmt\.Errorf\(" -e .go`
  - `gitmap aum search "errors\.New\(" -e .go`

### Family 02: Monadic `Result[T]` Envelopes & Types Extraction (Items 011–020)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/25-generic-return-types.md`, `02-spec/03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/04-result-types.md`.
- **Core Principles:**
  - Replace multi-value tuple returns `(T, error)` or `(map[K]V, error)` with monadic containers: `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`.
  - Result containers provide `.AppError()`, `.Fault()`, `.Value()`, and core affirmative predicates (`IsDefined()`, `IsEmpty()`, `HasRecord()`).
  - Pointer null-safety: functions returning `*Result[T]` must include defensive line-1 guard `if r == nil`.
  - Concrete Types Centralization (`types.go`): extract domain structs, generic instantiations, and Result wrappers into dedicated `types.go` files as single reusable named types (never leak bare inline generic instantiations like `Result[*CustomConfig]`).
- **GitMap Audit Queries:**
  - `gitmap aum search "\([a-zA-Z0-9_\*]+,\s*error\)" -e .go`
  - `gitmap aum search "\([a-zA-Z0-9_\*]+,\s*\*appfault\.AppError\)" -e .go`
  - `gitmap aum search "result\.Result\[" -e .go`

### Family 03: Positive Booleans (Rules P1–P4: Affirmative Prefixes & Negation Elimination) (Items 021–030)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-naming-prefixes.md`, `02-spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`.
- **Core Principles:**
  - `is` and `has` prefixes ONLY. All other prefixes (`can`, `should`, `was`, `will`, `did`, `must`) are strictly banned.
  - Ban single-letter boolean parameters (`v bool`, `b bool`, `flag bool`) in setters and methods.
  - Ban bare unprefixed names (`stop bool` -> `isStopped bool`, `stopOnFail bool` -> `isStopOnFail bool`, `defined` -> `isDefined`).
  - Total ban on negative boolean identifiers (`isNotValid`, `isUndefined`, `hasNoData`, `isUnset`).
  - Inverted negative replacement: replace `!isEmpty` or `!res.IsEmpty()` with affirmative `isDefined` or `res.IsDefined()`.
- **GitMap Audit Queries:**
  - `gitmap aum search "([a-z]+)\s+bool" -e .go`
  - `gitmap aum search "([a-zA-Z]+):\s*boolean" -e .ts`
  - `gitmap aum search "!.*IsEmpty\(\)"`
  - `gitmap aum search "isNot[A-Z]"`

### Family 04: Condition Simplification (Rules P5–P9: No Mixed Polarity & Guard Clauses) (Items 031–040)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md`, `02-spec/02-coding-guidelines/01-cross-language/04-code-style/03-conditions-and-extraction.md`.
- **Core Principles:**
  - NO MIXED POLARITY: never combine a positive check and a negative check in the same `if` condition (`if isA && !isB` is banned; evaluate separately or extract named affirmative composite predicate).
  - Implicit boolean evaluation only: NEVER write `if isReady == true` or `if (isValid === false)`.
  - Flatten nested branching: maximum nesting depth 1 with early returns and inverted guard clauses.
  - No inline assignment in `if`: avoid `if err := doWork(); err != nil` -> assign variable, blank line, then `if err != nil`.
- **GitMap Audit Queries:**
  - `gitmap aum search "==\s*true"`
  - `gitmap aum search "===\s*true"`
  - `gitmap aum search "==\s*false"`
  - `gitmap aum search "if.*&&.*!"`
  - `gitmap aum search "if.*;.*!= nil" -e .go`

### Family 05: Function Sizing (Hard Ceiling <=15 Lines, Target 8–10 Lines) (Items 041–050)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`.
- **Core Principles:**
  - Functions must adhere to a hard ceiling of <=15 lines (target 8–10 lines).
  - Decompose oversized functions into discrete, single-responsibility helper functions.
  - Map dispatch tables: replace repetitive chained if-else ladders and switch blocks matching keys/tokens with declarative map dispatch tables.
  - Clean constructor helpers: keep factory constructors compact and declarative.
- **GitMap Audit Queries:**
  - `gitmap aum search "func " -e .go`
  - `gitmap aum search "function " -e .ts`
  - `python 03-ai-scripts/13-file-size-guard.py` (or targeted line inspection via `gitmap cat`)

### Family 06: Parameter Structs (The 3-Parameter Cap & DTOs) (Items 051–060)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`, `02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/04-code-style-and-errors.md`.
- **Core Principles:**
  - Functions with >2-3 parameters MUST be refactored into a dedicated `*Params` / `*Options` struct or DTO.
  - Eliminate loose string/int/bool sequences that lead to argument transposition bugs.
  - Multiline formatting: when calling functions with >2 arguments, format one argument per line with trailing commas.
  - Variadic and spread parameters (`...T`, `...items`, `&[T]`): use variadic signatures where appropriate to eliminate artificial slice wrappers.
- **GitMap Audit Queries:**
  - `gitmap aum search "func.*\(.*,.*,.*,.*\)" -e .go`
  - `gitmap aum search "function.*\(.*,.*,.*,.*\)" -e .ts`

### Family 07: File Modularity & Sizing (Ceiling <=100 Lines Components, <=300 Lines Packages) (Items 061–070)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/04-code-style/05-function-and-type-size.md`.
- **Core Principles:**
  - Hard ceiling <=100 lines for components, views, and focused utility files.
  - Hard ceiling <=300 lines for complex package implementations.
  - Split oversized monolithic files into cohesive sibling modules: `types.go`, `handlers.go`, `helpers.go`, `adapters.go`.
  - React frontend components: decompose oversized JSX templates into extracted sub-components and custom hooks (`use*.ts`).
- **GitMap Audit Queries:**
  - `gitmap find "*.go"`
  - `gitmap find "*.tsx"`
  - `gitmap find "*.ts"`

### Family 08: Vertical Line Gaps & Whitespace Spacing (Blank Padding Rules) (Items 071–080)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md`, `02-spec/02-coding-guidelines/01-cross-language/21-newline-styling-examples.md`.
- **Core Principles:**
  - Mandatory blank line before `if` statements (unless at function start).
  - Mandatory blank line after closing brace `}` (unless followed by `else` or closing parent block).
  - Mandatory blank line before `return` statements (unless function body is only 1-2 lines).
  - Blank line padding around multiline struct instantiations and multiline function calls.
  - Strictly Unix LF (`\n`) line endings and UTF-8 without BOM.
- **GitMap Audit Queries:**
  - `gitmap aum search -r "[^\n]\n\s*if\s+" -e .go`
  - `gitmap aum search -r "\}\n\s*[a-zA-Z]" -e .go`
  - `gitmap aum search -r "[^\n]\n\s*return\s+" -e .go`

### Family 09: String Efficiency & Enums (Zero-Allocation & PascalCase `Type` Suffix) (Items 081–090)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/15-master-coding-guidelines/03-boolean-and-enum.md`, `02-spec/02-coding-guidelines/01-cross-language/16-lazy-evaluation-patterns.md`.
- **Core Principles:**
  - Zero-allocation string comparison: replace chained equality checks and `strings.ToLower(x) == "y"` with `strutil.EqualFoldAny(target, "a", "b")` or `strutil.EqualFoldAnyTrim`.
  - Lazy regex compilation: ban raw inline `regexp.MustCompile` inside loops or hot paths; use package `lazyregex` with thread-safe cached compilation.
  - Cross-language enum standards: all enums MUST be PascalCase and carry mandatory `Type` suffix (e.g. `StatusType`, `LogLevelType`, `RoleType`).
  - String-backed enums with exhaustive validator methods (`IsValid() bool`).
- **GitMap Audit Queries:**
  - `gitmap aum search "strings\.ToLower\(" -e .go`
  - `gitmap aum search "\.toLowerCase\(\)" -e .ts`
  - `gitmap aum search "regexp\.MustCompile\(" -e .go`
  - `gitmap aum search "type .* string" -e .go`

### Family 10: Testing Hygiene & Pathutil Abstraction (No Destructive OS Calls) (Items 091–100)
- **Canonical Specifications:** `02-spec/02-coding-guidelines/01-cross-language/14-test-naming-and-structure.md`, `02-spec/02-coding-guidelines/04-golang/06-testing-standards.md`.
- **Core Principles:**
  - Unit tests MUST NEVER trigger real OS shutdown, reboot, power-off, system modifications, or unmocked process termination.
  - Injectable command executors and filesystem interfaces for all system operations.
  - Mock duration verification: tests must complete within milliseconds without sleeping real clocks.
  - Cross-platform path handling: strict pathutil abstraction, never hardcoding `\` or `/` filesystem separators.
  - Bounded test execution: test suites must isolate side effects to temporary sandboxes.
- **GitMap Audit Queries:**
  - `gitmap aum search "exec\.Command\(" -e .go`
  - `gitmap aum search "os\.Exit\(" -e .go`
  - `gitmap aum search "time\.Sleep\(" -e .go`
  - `gitmap find "*test*.go"`

---

## 5. Structured 10-Field Specification Card Schema

Every single one of the 100 improvements MUST be documented using the strict, structured 10-field specification card schema below. No fields may be omitted, truncated, or replaced with placeholders:

```markdown
### Improvement-NNN: [Concise Descriptive Improvement Title]

1. **Target File Path:** `[strictly/relative/path/to/file.ext]`
2. **File Name:** `file.ext`
3. **Target Function / Symbol / Scope:** `PackageName.FunctionName` or `ClassName.methodName`
4. **Line Numbers:** Lines L12–L34
5. **Rule Violated & Canonical Citation:** `[Rule Code: Rule Name]` — Canonical Citation: `02-spec/02-coding-guidelines/.../rule.md`
6. **Violation Severity:** `[Critical | Major | Minor]`
7. **Current Code Context:**
```<lang>
// Verbatim offending snippet from the codebase
```
8. **Violation Analysis & Architectural Rationale:**
[In-depth technical explanation of why the current implementation violates repository standards, the operational risks or technical debt introduced, and how resolving it improves codebase resilience, readability, and performance.]

9. **Exact Fix Proposal:**
```<lang>
// Copy-paste ready refactored code satisfying all guidelines
```

10. **Blast Radius, Downstream Callers & Verification Contract:**
- **Downstream Callers (discovered via `gitmap aum search`):**
  * `[relative/path/to/caller1.ext:Line]`: CallerFunction()
  * `[relative/path/to/caller2.ext:Line]`: AnotherCaller()
- **Blast Radius Assessment:** [Low | Medium | High] — [Explanation of affected boundaries and backward-compatibility preservation]
- **Verification Command Contract:** `[exact non-build targeted linter or syntax check]` -> Expected: `exit 0`
```

---

## 6. Step 0: Preflight, Platform Handshake & Discovery Setup (Phase 1 Budget)

1. **Platform Handshake:** Confirm tools (`invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, `write_file`, `edit_file`, `run_command`).
2. **Commands & Directory:** Confirm `gitmap --version` and `python --version` exit 0. Verify GitMap with a harmless call (`gitmap lf readme.md`), not `--help`. `run_command` uses `Cwd` in workspace root, paths relative.
3. **Working Tree Cleanliness:** Run `git status --porcelain`. Record modified files in ledger; never touch them. Confirm root `readme.md` is lowercase. Read `.ai-memory/what-to-read.md`, `strictly-avoid.md`, `coding-guidelines.md`.
4. **Step 0 Heatmap & RC Preflight Checks:**
   - Preflight status & tables: `gitmap heatmap status` and `gitmap heatmap tables`
   - Pre-action heatmap logging: `gitmap heatmap log --task <task> --action READ --file <path> --state EXPLORATION --expected "<goal>"`
5. **SQLite Task DB & Deterministic Slug Initialization:**
   `gitmap agent task init --name "100-improvements-<codebase>" --budget 400`
6. **Ledger Creation:** Create `.ai-memory/temp-agents/<nn>-find-100-improvements/ledger.md` tracking all 100 items across the 10 families.

---

## 7. Phase 1: Codebase Scanning & Discovery (Steps 1 .. PHASE_1_BUDGET)

The lead orchestrator partitions the 10 rule families across 3 concurrent `research` discovery subagents:

### Wave 1 Discovery Dispatch (`TypeName: "research"`):
- **Worker 01 (Research):** Discovers Families 01–03 (Error Management, Monadic Results, Positive Booleans).
- **Worker 02 (Research):** Discovers Families 04–07 (Condition Simplification, Function Sizing, Parameter Structs, File Modularity).
- **Worker 03 (Research):** Discovers Families 08–10 (Vertical Spacing, String Efficiency/Enums, Testing Hygiene).

Each discovery subagent executes targeted GitMap searches:
- `gitmap lf <target-dir>` to inventory active files.
- `gitmap aum search "<pattern>" <target-dir> -e <ext>` for pattern violations.
- `gitmap cat <filepath>` to extract verbatim line numbers and context.
- Trace caller blast radius using `gitmap aum search "<symbol>"`.

Subagents report their findings back to the lead orchestrator via `send_message` with verbatim file paths, line ranges, offending snippets, and caller dependencies.

---

## 8. Phase 2: Card Authoring, Verification & Enqueueing (Steps 101 .. 250)

Once candidate improvements are received, the lead orchestrator coordinates 3 concurrent `self` authoring subagents to author the 10-field specification cards and assemble the Master Catalog Summary Table.

### Micro-Batching Waves (5–8 Files Per Wave)
All 100 improvements are grouped into micro-batches of 5–8 files to preserve bounded context and prevent tool payload truncation:
- **Wave 01:** Improvements 001–010 (Family 01: Error Management)
- **Wave 02:** Improvements 011–020 (Family 02: Monadic Results & Types Extraction)
- **Wave 03:** Improvements 021–030 (Family 03: Positive Booleans)
- **Wave 04:** Improvements 031–040 (Family 04: Condition Simplification)
- **Wave 05:** Improvements 041–050 (Family 05: Function Sizing)
- **Wave 06:** Improvements 051–060 (Family 06: Parameter Structs)
- **Wave 07:** Improvements 061–070 (Family 07: File Modularity & Sizing)
- **Wave 08:** Improvements 071–080 (Family 08: Vertical Line Gaps & Whitespace)
- **Wave 09:** Improvements 081–090 (Family 09: String Efficiency & Enums)
- **Wave 10:** Improvements 091–100 (Family 10: Testing Hygiene & Pathutil)

### Enqueueing Plans to `.ai-memory/plans/pending/` & Subtasks

> [!IMPORTANT]
> **Mandatory Plan Enqueueing & Sequential Next-Wave Generation:**
> All 100 improvements MUST be systematically preserved and enqueued as pending execution plans in the repository so they can be executed by subsequent refactoring sessions:
> 1. Write the Master Plan to `.ai-memory/plans/pending/XX-find-100-improvements-wave-YY.md` (where `XX` is the zero-padded sequence number in `pending/`, and `YY` is the wave counter).
> 2. Decompose the 100 improvements into modular subtasks stored in `.ai-memory/plans/subtasks/XX-find-100-improvements-wave-YY/`:
>    - `01-family-01-error-management.md` (Items 001–010)
>    - `02-family-02-monadic-results.md` (Items 011–020)
>    - `03-family-03-positive-booleans.md` (Items 021–030)
>    - `04-family-04-condition-simplification.md` (Items 031–040)
>    - `05-family-05-function-sizing.md` (Items 041–050)
>    - `06-family-06-parameter-structs.md` (Items 051–060)
>    - `07-family-07-file-modularity.md` (Items 061–070)
>    - `08-family-08-vertical-spacing.md` (Items 071–080)
>    - `09-family-09-string-efficiency-enums.md` (Items 081–090)
>    - `10-family-10-testing-hygiene.md` (Items 091–100)
> 3. Register the new pending plan in `.ai-memory/plans/readme.md` under `## Pending Plans`.
> 4. Enqueue the roadmap in `.ai-memory/what-to-read.md`.
> 5. **Sequential Generation Rule:** Each subsequent run of `/find-100-improvements-in-code-base` MUST inspect existing pending and completed improvement plans, increment the wave counter (`wave-02`, `wave-03`, ...), scan unaddressed files or deeper architectural layers, and generate the NEXT 100 unique improvements without duplicate cards.

---

## 9. Master Catalog Summary Table

The generated audit specification MUST contain the Master Catalog Summary Table indexing all 100 improvements:

| Id | Target File | Target Function / Symbol | Rule Violated & Code | Severity | Summary of Fix Proposal |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `IMP-001` | `[pkg/service/auth.go]` | `ValidateSession` | `CG-ERR-001: *appfault.AppError` | Critical | Return `*appfault.AppError` instead of bare `error` |
| `IMP-002` | `[pkg/store/user.go]` | `FindByID` | `CG-ERR-002: Swallowed Error` | Critical | Wrap SQL error in `appfault.New()` with code `ErrDatabaseQuery` |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-010` | `[pkg/api/handler.go]` | `SendResponse` | `CG-ERR-005: Envelope Parity` | Major | Serialize error using standard response envelope schema |
| `IMP-011` | `[pkg/repo/query.go]` | `FetchConfig` | `CG-RES-001: Monadic Result` | Major | Replace `(*Config, error)` with `Result[*Config]` |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-020` | `[pkg/repo/types.go]` | `Types Extraction` | `CG-TYPE-001: types.go Centralization` | Minor | Extract inline Result aliases into dedicated `types.go` |
| `IMP-021` | `[pkg/worker/progress.go]` | `SetStopOnFail` | `P1: Affirmative Prefixing` | Major | Rename parameter `v bool` to `isStopOnFail bool` |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-030` | `[pkg/model/status.go]` | `CheckPresence` | `P3: Inverted Empty Elimination` | Minor | Replace `!res.IsEmpty()` with affirmative `res.IsDefined()` |
| `IMP-031` | `[pkg/auth/guard.go]` | `CheckAccess` | `P6: Mixed Polarity Ban` | Critical | Separate `if isValid && !isExpired` into discrete guard clauses |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-040` | `[pkg/parser/lexer.go]` | `ScanToken` | `P7: Inline Assignment Ban` | Minor | Extract `if tok := next(); tok != nil` to standalone line |
| `IMP-041` | `[pkg/engine/eval.go]` | `EvaluateAST` | `CG-SIZE-001: 15-Line Function Ceiling` | Major | Decompose 45-line evaluator into 4 sub-expression handlers |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-050` | `[pkg/router/route.go]` | `DispatchRoute` | `CG-DISP-001: Map Dispatch Table` | Major | Replace 35-line switch statement with map dispatch lookup |
| `IMP-051` | `[pkg/client/connect.go]` | `DialCluster` | `CG-PARAM-001: 3-Parameter Cap` | Major | Refactor 6 string/int arguments into `*DialParams` struct |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-060` | `[pkg/format/render.go]` | `RenderTable` | `CG-PARAM-002: Variadic Spread` | Minor | Convert `items []string` to variadic `...string` |
| `IMP-061` | `[pkg/ui/component.tsx]` | `DashboardView` | `CG-FILE-001: 100-Line Component Ceiling`| Major | Extract metrics cards and chart sections into child components |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-070` | `[pkg/service/monolith.go]` | `ServiceEngine` | `CG-FILE-002: 300-Line Package Ceiling` | Critical | Split 620-line file into `engine.go`, `types.go`, `helpers.go` |
| `IMP-071` | `[pkg/util/filter.go]` | `ApplyFilters` | `CG-STYLE-001: Blank Line Before If` | Minor | Insert required blank newline before `if` statement |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-080` | `[pkg/auth/token.go]` | `GenerateClaims` | `CG-STYLE-004: Blank Line Before Return`| Minor | Insert required blank newline before final `return` |
| `IMP-081` | `[pkg/search/query.go]` | `MatchCommand` | `CG-STR-001: Zero-Alloc Case Folding` | Major | Replace `strings.ToLower(x) == "y"` with `EqualFoldAny` |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-090` | `[pkg/config/types.go]` | `LogLevel` | `CG-ENUM-001: PascalCase Type Suffix` | Major | Rename `loglevel` to `LogLevelType` with validator method |
| `IMP-091` | `[pkg/sys/process_test.go]` | `TestRebootGuard` | `CG-TEST-001: OS Destructive Call Ban` | Critical | Inject mock executor to eliminate unmocked system reboot call |
| `...` | `...` | `...` | `...` | `...` | `...` |
| `IMP-100` | `[pkg/fs/path_test.go]` | `TestResolvePaths` | `CG-TEST-003: Pathutil Abstraction` | Minor | Replace hardcoded `\` separators with `filepath.Join` |

---

## 10. Phase 3: Consolidation, Evidence Verification & Atomic GitMap Push

1. **Verify Deliverables:**
   - Confirm all 100 specification cards are present and complete (zero placeholders, zero missing fields).
   - Confirm Master Catalog Summary Table indexes exactly 100 rows across all 10 rule families.
   - Confirm pending plan file exists at `.ai-memory/plans/pending/XX-find-100-improvements-wave-YY.md`.
   - Confirm 10 subtask files exist under `.ai-memory/plans/subtasks/XX-find-100-improvements-wave-YY/`.
   - Confirm `.ai-memory/plans/readme.md` registers the new plan.
2. **Push Gate Verification:**
   - Targeted checks exit 0 (>0 files scanned; every time code is written, do not build code).
   - Secrets gate clean: check changed files via `git status --porcelain`, run `python linter-scripts/check-forbidden-strings.py`.
   - Strict relative Git paths verified across all authored markdown documents.
   - Phase 3 cognitive heatmap verification: `gitmap heatmap tree --task <task>` and `gitmap heatmap status --task <task>`.
3. **Atomic Commit & Push:**
   Call `gitmap cpf "<module> - <feature summary>"` (e.g. `gitmap cpf "audit - catalog 100 codebase improvements wave 01"`).
   Format as hyphen `-` because GitMap automatically prepends `Feature: ` (providing the colon). TOTAL BAN on colons inside the GitMap argument. NEVER include conventional prefixes (`feat(...)`, `docs(...)`). Conclude turn immediately upon successful GitMap push.

---

## 11. Targeted Verification Checks (R2)

Run only fast, file-scoped checks on modified files:
- **Coding Guidelines & Boolean Linter:** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`
- **Relative Path Linter:** `python linter-scripts/check-relative-paths.py`
- **Markdown Link & Doc Path Linter:** `python 03-ai-scripts/22-doc-path-linter.py <folder>`
- **Sequence Integrity Linter:** `python linter-scripts/check-sequence-integrity.py`
- **Forbidden Strings Check:** `python linter-scripts/check-forbidden-strings.py`

---

## 12. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or user instructions are given before this prompt are verified as highest priority and non-negotiable.
- [ ] MANDATORY SUBAGENT SPAWNING GATE (ZERO SOLO EXECUTION): Verified that `invoke_subagent` (`A = 3, H = 2, Capacity = 6`) was explicitly called in Phase 1 discovery and Phase 2 card authoring.
- [ ] TREE SWARM & DYNAMIC JOIN MANDATE: Verified Tree Swarm Hierarchy (Root Orchestrator -> Branch Coordinators -> Leaf Workers & Dynamic Assistants).
- [ ] NO SUBAGENT FILE WRITE COLLISIONS: Subagents must have strictly disjoint, non-overlapping file boundaries.
- [ ] NO READ-ONLY SUBAGENT FILE WRITES: Subagents of `TypeName: "research"` have read-only tools and cannot write files.
- [ ] NO INVALID SUBAGENT TYPENAMES: `TypeName` MUST be strictly `"research"` or `"self"`.
- [ ] NO SUBAGENT GIT EXECUTION: Subagents must NEVER run any git commands.
- [ ] NO LOCAL TESTS (TOTAL BAN): Never run all tests locally or full test suites (`06-cicd-local-runner.py`, `pytest`, `npm test`, `go test ./...`) unless explicitly asked.
- [ ] NO BUILD CHECKING (TOTAL BAN): Every time you write code, you do not build the code. Never run build commands (`go build`, `npm run build`, `vite build`, `cargo build`).
- [ ] NO AUTOMATIC RELEASES: Never bump versions, edit `version.json`, or update changelogs unless explicitly commanded.
- [ ] NO RAW GIT COMMITS: Never execute raw `git commit`, `git add -A`, or raw `git push`. All staging and pushing must be via GitMap semantic commands (`gitmap cpf "<module> - <summary>"`).
- [ ] NO INTERMEDIATE COMMITS: All files across the run must be accumulated and committed together in a single atomic GitMap commit at the final step of Phase 3.
- [ ] MANDATORY COMPLETION INVARIANT (NO PUSH = NOT DONE): If the code is not committed to Git and pushed upstream to GitHub, the task is strictly INCOMPLETE and NOT DONE.
- [ ] NO RG, RIPGREP, OR RAW SHELL SEARCHES (TOTAL BAN): Banned ALL variants of `Get-ChildItem`, `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Search exclusively through `gitmap aum search`, `gitmap find`, `gitmap lf`, or `gitmap cat`.

---

## 13. Non-Negotiable Coding Guidelines Checklist

- [ ] Master Guidelines: Grounded strictly in `02-spec/02-coding-guidelines/` and `02-spec/03-error-manage/`.
- [ ] Error Management: Structured `*appfault.AppError` returns, zero swallowed errors.
- [ ] Monadic Results: `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]` with `types.go` centralization.
- [ ] Positive Booleans: `is` and `has` prefixes only. No negative flags or bare parameter names.
- [ ] Condition Simplification: No mixed polarity, implicit checks, guard clauses.
- [ ] Function Sizing: Ceiling <=15 lines, target 8–10 lines.
- [ ] Parameter Structs: The 3-parameter cap, dedicated DTO structs for >2-3 parameters.
- [ ] File Modularity: Ceiling <=100 lines for components, <=300 lines for packages.
- [ ] Vertical Spacing: Blank line before `if`, blank line after `}`, blank line before `return`.
- [ ] String Efficiency & Enums: Zero-allocation `EqualFoldAny`, PascalCase enums with `Type` suffix.
- [ ] Testing Hygiene: Zero OS shutdown/reboot calls, mock duration verification.
- [ ] Strict Relative Git Paths: All markdown links strictly relative; total ban on absolute paths and `file:///` URIs.

---

## 14. Anti-Hallucination & Blast Radius Checklist

- [ ] Verbatim Code Snippets: All current code contexts must be exact, verbatim extractions from the repository (via `gitmap cat`).
- [ ] Caller Tracing: Downstream callers must be traced using `gitmap aum search "<symbol>"`.
- [ ] Exact Fix Proposals: Fix proposals must be valid, copy-paste ready syntax satisfying all repository coding guidelines.
- [ ] Concrete Verification Contracts: Verification commands must be actionable, non-build syntax or linter checks.
- [ ] Zero Placeholders: Confirmed zero `TODO`, `TBD`, or `[...]` placeholders remain in specification cards.
- [ ] Plan & Subtask Enqueueing: Confirmed that all 100 improvements are enqueued as pending plans in `.ai-memory/plans/pending/XX-<slug>.md` and subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`.
- [ ] Successive Wave Continuity: Confirmed that re-running this prompt automatically increments wave count and targets the next 100 improvements.
