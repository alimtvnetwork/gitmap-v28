# [V6] Multi-Repository Synchronization & Downstream Codebase Mirroring — Workflow (must follow)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Preflight, Discovery, Branch Validation, Target Inventory & Plan Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Backup Branch Creation, Safe Mirroring Waves, Additive Script Review, Verification & Release)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: `/sync-other-codebase` or paste below task directives
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously pull target repositories on their base branch, create and push dedicated pre-change backup branches, mirror canonical prompts (`01-prompts/`), agent skills (`.agents/skills/`, `.cursor/skills/`), coding guidelines (`02-spec/02-coding-guidelines/`), and additive AI scripts (`03-ai-scripts/`, `.agents/scripts/`) across target codebases, strictly enforcing the 5 Non-Negotiable Boundaries (Spec 21 Exclusion, Bump Script Protection, Additive-Only AI Scripts, Memory & Plans Protection, Zero Secrets Leakage), verify with targeted linters, and complete safe release ceremony per target repository with ZERO hardcoded paths.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the target repository inventory and sync parameters in visible chat before background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step target repository validation and inventory before mirroring. Ensure all target repository paths, git base branches, and backup branch names are cleanly established in the audit ledger and subtask plans before dispatching worker waves.

---

## 0. Dynamic Parameter Configuration (Zero Hardcoded Paths Mandate)

> [!CAUTION]
> **ZERO HARDCODED PATHS MANDATE:**
> The source repository and target repository paths MUST NEVER be hardcoded into this prompt or any generated instructions.
> The caller dynamically provides `SOURCE_REPO` and `TARGET_REPOS` in the user preamble, prompt invocation, or task definition.
> The orchestrator and worker subagents ingest these parameters dynamically at runtime.

```text
SOURCE_REPO  = <path-to-source-repo>   # Path to the canonical source repository (provided by caller)
TARGET_REPOS = [                       # List of target repository paths provided dynamically by caller
    "<path-to-target-repo-1>",
    "<path-to-target-repo-2>",
    ...
]
```

### Dynamic Parameter Guidelines:
1. **Dynamic Ingestion:** If `SOURCE_REPO` is unspecified, default to the current workspace root dynamically via relative path determination or current directory context.
2. **Target Repository List:** If `TARGET_REPOS` is provided in the user request, process exclusively the specified repositories. If the caller requests "all connected repos" or invokes `03-ai-scripts/38-sync-prompts-skills-scripts.py`, ingest the dynamic target list discovered by the script or CLI flags (`--repo <slug>`).
3. **No Hardcoded Absolute Paths:** NEVER write machine-specific paths (e.g. `C:\Users\...` or `/home/...`) into any file, plan, or spec. Use relative paths or dynamic shell parameter interpolation.

---

## 1. 🚨 Mandatory Subagent Spawning Gate (A = 2, H = 2 — Zero Solo Execution Allowed)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print text saying "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) across the workflow:
  1. **Planning Step (A = 2 `research` subagents):** Spawn 2 read-only discovery subagents (`TypeName: "research"`) to inspect target repositories, verify clean working trees, identify local base branches, and map existing bump scripts or customized AI scripts. Lead agent writes the unified synchronization plan.
  2. **Spec Step (A = 2 `self` subagents or lead):** Spawn 2 subagents (`TypeName: "self"`) to author modular, disjoint sync specs and target batch subtask files.
  3. **Execution Step (A = 2 `self` worker subagents):** Spawn 2 worker subagents (`TypeName: "self"`) to execute safe mirroring, diff verification, and targeted linting in strictly disjoint repository groups or file boundaries.
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning discovery, spec writing, or codebase synchronization by itself without calling `invoke_subagent`.

---

## 2. The 5 Non-Negotiable Synchronization Boundaries

Every sync operation across codebases MUST strictly enforce these five boundaries without exception:

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                   5 NON-NEGOTIABLE SYNCHRONIZATION BOUNDARIES               │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Spec 21 Exclusion      │ NEVER touch 02-spec/21-* (app/domain specs).    │
│ 2. Bump Script Protection │ NEVER overwrite repo-specific bump scripts.     │
│ 3. Additive-Only Scripts  │ Copy new AI scripts; diff/preserve modified.   │
│ 4. Memory & Plans Safe    │ NEVER touch .ai-memory/memory/ or plans/.       │
│ 5. Zero Secrets Leakage   │ NEVER sync .env or private credentials.         │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Boundary 1: Spec 21 Exclusion (TOTAL BAN)
- **Rule:** NEVER sync, copy, modify, or touch any directory matching `02-spec/21-*` (such as `21-app/`, `21-app-issues/`, `21-app-db/`, `21-app-ui-design-system/`).
- **Rationale:** Directories prefixed with `21-` contain private application specifications, user journeys, feature specs, database entity definitions, and architecture plans unique to that specific repository. Mirroring these from the source wipes out the target repository's application identity and destroys domain requirements.

### Boundary 2: Bump Script Protection (IMMUTABLE)
- **Rule:** NEVER overwrite version bumping scripts or version manifests in target repositories (`bump-version.mjs`, `bump_versions.py`, `37-bump-version.py`, `scripts/bump*.py`, `version.json`, etc.).
- **Rationale:** Each repository maintains unique SemVer targets, custom package manifests (Go `version.go`, Node `package.json`, Python `pyproject.toml`, Rust `Cargo.toml`, PHP `composer.json`), and custom git tag formats. Overwriting these scripts breaks release pipelines and invalidates tag history.

### Boundary 3: Additive-Only AI Scripts (NO BLIND OVERWRITES)
- **Rule:** Brand-new AI scripts from `03-ai-scripts/` and `.agents/scripts/` that do not exist in the target repository are copied cleanly. However, existing scripts that have been customized in the target repository MUST NEVER be blindly overwritten.
- **Protocol:** Before replacing an existing script, compute the diff. If the target repository has customized flags, localized paths, or repo-specific hooks, preserve those local modifications. Never clobber local script customizations.

### Boundary 4: Memory & Plans Protection (TOTAL ISOLATION)
- **Rule:** NEVER modify, sync, mirror, or touch `.ai-memory/memory/`, `.ai-memory/plans/`, `.ai-memory/temp-agents/`, `.ai-memory/cicd-issues/`, or `.ai-memory/ambiguous-questions/` in target repositories.
- **Rationale:** Target repositories maintain their own operational logs, completed plans, active execution tasks, and debugging memories. Central synchronization must never overwrite target operational memory.

### Boundary 5: Zero Secrets Leakage (STRICT PURGE)
- **Rule:** NEVER sync `.env`, `.env.*`, API keys, tokens, or credentials across repositories.
- **Location Mandate:** If credentials exist or are needed during synchronization, they must reside exclusively in the `repo-secrets` folder in the default work directory via `gitmap rs`. Standard repositories must remain 100% free of secrets.

---

## 3. Automation Engine & Script Architecture

### Primary Automation Script: `03-ai-scripts/38-sync-prompts-skills-scripts.py`

The canonical script `03-ai-scripts/38-sync-prompts-skills-scripts.py` automates the synchronization ceremony safely:

```text
SYNCHRONIZATION WORKFLOW PER TARGET REPOSITORY:
1. Detect base branch: git branch --show-current (main / master)
2. Pull latest changes: git pull origin <base_branch> --no-rebase
3. Pre-change backup branch: git checkout -b backup/pre-v3-nsteps-sync-<timestamp> && git push origin backup/...
4. Pre-change release verification: ensure pre-change release tag and release branch exist
5. Safe directory mirroring:
   - 01-prompts/                -> 01-prompts/                (mirror, excluding 06-old-prompts/)
   - .agents/skills/            -> .agents/skills/            (mirror)
   - .cursor/skills/            -> .cursor/skills/            (mirror)
   - 03-ai-scripts/             -> 03-ai-scripts/             (additive-only, bump scripts protected)
   - .agents/scripts/           -> .agents/scripts/           (additive-only)
   - 02-spec/02-coding-guidelines/ -> 02-spec/02-coding-guidelines/ (mirror)
   - 02-spec/07-design-system/  -> 02-spec/07-design-system/  (conditional mirror)
   - 02-spec/17-consolidated-guidelines/ -> 02-spec/17-...    (conditional mirror)
6. Commit & Push on base branch via GitMap or git
7. Post-change release ceremony:
   - Create release branch: release/v<next_ver>
   - Create release tag: v<next_ver>
   - Push release branch and tag to origin
   - Merge release branch back into base branch with [skip ci] and push
```

### Script CLI Flags & Options:
- Single repository target: `python 03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <slug>`
- Dry-run simulation: `python 03-ai-scripts/38-sync-prompts-skills-scripts.py --dry-run`
- Worker concurrency: `python 03-ai-scripts/38-sync-prompts-skills-scripts.py --workers 4`
- Skip release ceremony: `python 03-ai-scripts/38-sync-prompts-skills-scripts.py --skip-release`

---

## 4. Precedence Hierarchy & Core Operational Rules (R1–R16)

### Precedence Hierarchy (Highest First):
1. **User Instructions & Preamble:** Directives and parameters ABOVE this prompt outrank everything below.
2. **Platform Limits:** Native tools, Artifact Review Policy, permission prompts, hooks.
3. **Repo Rules:** `agents.md`, `.ai-memory/strictly-avoid.md`, and `coding-guidelines.md`.
4. **This Prompt.**

### Core Operational Rules (Cited by ID):
- **R1 Zero Builds or Test Suites (TOTAL BAN):** NEVER run `go build`, `npm run build`, `vite build`, `go test ./...`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. CI verifies full builds and suites. Routine sync turns must never waste execution steps on heavy compilation or tests.
- **R2 Targeted Checks Only:** Run only fast, file-scoped checks on specifically modified files (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`). A check scanning 0 files is a FAIL.
- **R3 Evidence or It Did Not Happen:** Every `DONE`, `PASS`, or "verified" claim MUST cite a concrete file path, git diffstat, or command exit code (`exit 0`). Vague verbal assurances are treated as hallucinations and auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths:** Before calling any script or command, verify it exists. If a tool is missing, use the documented fallback and log it in the ledger.
- **R5 Mandatory Subagents (`invoke_subagent`):** Spawning subagents is an **ABSOLUTE MUST** when `invoke_subagent` exists. Use `TypeName: "research"` for discovery and `TypeName: "self"` for execution. Solo execution is an auto-reject failure.
- **R6 One Owner Per File (Disjoint Bounding Boxes):** Within every worker wave, each file has exactly one owner. Shared indexes (`01-prompts/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/plans/readme.md`) belong exclusively to the lead orchestrator.
- **R7 Git Safety & Isolation:** Subagents are strictly banned from running `git add`, `git commit`, `git push`, or modifying git state in shared workspaces. Only the lead orchestrator executes git commands.
- **R8 Stage by Explicit Path (BAN on `git add -A`):** Staging MUST be performed using `git add -- <paths from ledger stage list>` or atomic GitMap commands. NEVER run `git add -A`, `git add --all`, or `git add .`.
- **R9 One Atomic Commit at Turn End:** Accumulate all verified changes and commit once at the final step using hyphen format (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`). TOTAL BAN on colons inside GitMap message arguments.
- **R10 Zero Unauthorized Releases:** Never bump versions or trigger release scripts unless explicitly commanded by the user or part of the authorized target repo sync release ceremony.
- **R11 Strict Relative Git Paths & Lowercase Hygiene:** TOTAL BAN on absolute filesystem paths and `file:///` URIs inside repository files. All paths must be relative from git root. All new filenames MUST use strictly lowercase naming.
- **R12 No Polling / Immediate Turn Yielding:** After dispatching subagents, print a one-line progress notification and **STOP CALLING TOOLS**. Never busy-poll `manage_task` or check the filesystem in a loop.
- **R13 Two-Strike Retry Cap & Anti-Looping:** Failing the same tool call or operation twice requires changing the approach. A subtask that fails two consecutive rounds is marked `FAILED` with an RCA.
- **R14 100% Ambiguity & Decision Boundaries:**
  - *Non-Blocking Ambiguity:* Choose the most conservative option adhering to guidelines, record it in the ledger under `Assumptions:`, and proceed immediately.
  - *Blocking Ambiguity:* Use `ask_question` once, log in `.ai-memory/ambiguous-questions/`, and continue working on unblocked tasks.
- **R15 Zero Generated Artifacts Committed:** Never commit build caches, logs, temp scripts, or newly generated binaries unless the repository already tracked them.
- **R16 Repo Secrets & Credentials Mandate:** Zero secrets in standard repositories. Offload credentials to `repo-secrets` in the default work directory via `gitmap rs`.

---

## 5. GitMap High-Speed Command Primacy

GitMap is your **PRIMARY** acceleration engine across all operations:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Live Streaming Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` | `gitmap aum grep "<pat>"` | Multi-core streaming live file search (replaces `Select-String`, `git grep`) |
| **Indexed Symbol Search** | `gitmap search "<query>" [--limit <n>]` | `gitmap search` | Instant SQLite cached keyword/symbol search across indexed repos |
| **Wildcard Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core file finding |
| **File Listing** | `gitmap list-files [pattern] [-ext <ext>]` | `gitmap lf [pat]` | Instant indexed repository inventory |
| **Substring Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Zero-disk memory streaming to stdout |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpf "<module> - <summary>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, formats prefix, and pushes atomically (hyphen `-` format) |
| **Pipeline Waiting** | `gitmap pipeline-ai status --json` | `gitmap pl-ai` | Non-polling dynamic ETA CI/CD monitor |

### 🔍 Code & Symbol Search Protocol (TOTAL BAN ON `Select-String` & `git grep`)
- **Default:** Always use `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r]` for text and regex discovery.
- **TOTAL BAN:** NEVER run PowerShell `Select-String`, `Get-ChildItem -Recurse`, `git grep`, `grep`, or `findstr`.

---

## 6. The Unified Master Pipeline (Atomic Numbered Steps)

Execute multi-repository synchronization via a strict 3-Phase pipeline.

### Phase 1A: Verbatim Capture, Target Parameter Extraction & Chat Showcase (Turn 1)

1. **Top-Instruction Priority Verification:** Verify whatever directives or repository parameters are supplied above this prompt.
2. **Showcase Given Task First (Turn 1 Action):** In your VERY FIRST response turn, output the confirmed task breakdown and target repository list directly in visible chat. Never execute tools silently without displaying the task breakdown first!
3. **Lossless Verbatim Capture:** Store incoming prompt under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Target Inventory Extraction:** Extract `SOURCE_REPO` and list of `TARGET_REPOS` into discrete actionable items (`Task-01`, `Task-02`, ...).
5. **Mandatory Same-Turn Tool Chaining:** Emit the breakdown in chat with clean vertical formatting, and in the **EXACT SAME TURN**, invoke your first tool call. Never emit text alone and stop.

```markdown
### 📋 Confirmed Task Breakdown & Target Ingestion

1. **Task-01: Target Codebase Preflight & Safety Audit**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Validate base branch, clean working tree, and bump script inventory for target repos]
   - **Actionable Scope:** [Target repositories inventory and branch verification]
   - **Target Areas:** `[TARGET_REPOS]`

2. **Task-02: Multi-Repository Synchronization & Safe Mirroring**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Create backup branches, mirror canonical assets enforcing 5 non-negotiable boundaries]
   - **Actionable Scope:** [Safe mirroring of prompts, skills, guidelines, and additive scripts]
   - **Target Areas:** `[TARGET_REPOS]`

Proceeding directly to Preflight & Phase 1B Discovery (Active Tool Call Running Below).
```

### Phase 1B: Discovery & Safety Audit (A = 2 `research` Discovery Subagents)

1. **Platform Handshake & SQLite Task DB:**
   - Confirm active tools.
   - Initialize task DB: `python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "<sync task name>" --budget 300`
   - Initialize `.ai-memory/temp-agents/<nn>-<slug>/ledger.md`.
2. **Spawn Discovery Subagents (A = 2 `research` subagents):**
   - Subagent 1: Inspect target repositories, check git working trees (`git status`), detect base branches (`main`/`master`), verify existing tags.
   - Subagent 2: Scan target repositories for existing bump scripts (`bump-version.mjs`, `bump_versions.py`, `37-bump-version.py`, `version.json`) and existing custom scripts in `03-ai-scripts/` to ensure Boundary 2 and Boundary 3 compliance.
3. **Yield Turn:** Lead stops calling tools after `invoke_subagent` to allow reactive wakeup.
4. **Master Sync Plan Generation (Lead Agent):** Lead synthesizes discovery reports and writes `.ai-memory/plans/pending/nn-<slug>.md`.

### Phase 1C: Spec Step & Modular Subtasks (A = 2 `self` Subagents or Lead)

1. Author canonical sync spec in `02-spec/21-app/nn-<slug>/01-sync-architecture-spec.md` (in source repo only; NEVER in target repos).
2. Author modular subtasks in `.ai-memory/plans/subtasks/nn-<slug>/`.
3. Populate subtasks into SQLite database:
   `python 03-ai-scripts/46-agent-sqlite-task-manager.py add-subtasks --db <databasePath> --tasks-json '[...]'`
4. Unconditionally proceed to Phase 2. ZERO intermediate commits during Phase 1!

### Phase 2: Execution Step (Worker Waves) (Steps 151 .. 300)

1. **Mandatory Subagent Dispatch (`invoke_subagent`, A = 2 `self` workers):**
   - Subagent 1: Executes pre-change backup branch creation (`backup/pre-v3-nsteps-sync-<timestamp>`) and pushes backup branches.
   - Subagent 2: Executes safe asset mirroring (prompts, skills, guidelines, additive scripts) adhering strictly to the 5 Non-Negotiable Boundaries.
   - Alternatively, workers process disjoint batches of `TARGET_REPOS` via `03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <slug>`.
2. **Self-Contained Worker Brief Envelope:**
   Inject complete instructions into the worker brief including SQLite logging commands, git ban for workers, and 5 Non-Negotiable Boundaries.
3. **Turn-Yielding & Crash Forensics:**
   - Lead agent yields turn after `invoke_subagent`.
   - On worker completion or failure, lead inspects SQLite task manager (`diagnose` / `status`).
4. **Targeted Verification Linters:**
   - Verify all synced files in target repositories have valid LF line endings, strictly relative paths, and zero syntax errors.
   - Run targeted linting: `python 03-ai-scripts/05-guideline-autofixer.py <synced-dir> --check-only`.

### Phase 3: Consolidation, Evidence Verification & Release Ceremony

1. **Consolidate Subtasks:**
   Merge completed subtasks into `.ai-memory/plans/completed/nn-<slug>.md`. Remove pending subtask files.
2. **Release Ceremony per Target Repository (if applicable):**
   - If full release ceremony is enabled: create post-change release branch (`release/v<next>`), create tag (`v<next>`), push to origin, merge back to base branch with `[skip ci]`, and push.
3. **Update Central Registers:**
   - Update `01-prompts/readme.md`, `.ai-memory/prompts.md`, and `.ai-memory/plans/readme.md`.
4. **Final Atomic Commit in Source Repo:**
   - Call `gitmap cpf "SYNC - synchronize prompts skills and guidelines to connected repositories"`.
   - Use hyphen `-` format (no colons in GitMap message argument).

---

## 7. Self-Contained Worker Brief Template

```text
You are Worker <NN> for task nn-<slug>. You have no prior chat context; this brief is your complete specification.

### Boundaries & Crash Prevention:
- Read any file in the workspace; edit ONLY your Owned Files / Repositories: <relative paths or assigned target repos>.
- TOTAL BAN ON GIT COMMANDS IN SHARED WORKSPACES: Workers NEVER run raw git commands (git add, git commit, git push, git checkout) directly in the source workspace. In target repositories, execute exclusively through 03-ai-scripts/38-sync-prompts-skills-scripts.py or lead-delegated operations.
- TOTAL BAN ON COMMITS IN SOURCE REPO: Workers NEVER commit or push in source repo. Committing is exclusively reserved for the Lead Agent at Phase 3 via GitMap.
- Code & Symbol Search: Use GitMap exclusively: `gitmap aum search "<pattern>" [dir]` or `gitmap search`. TOTAL BAN on PowerShell Select-String or git grep.
- After C tool calls, stop and report what you have.
- A tool failing twice: reply "STATUS: BLOCKED" with exact error and stop.
- Adhere to R1, R2, and R11 by ID.

### The 5 Non-Negotiable Boundaries:
1. Spec 21 Exclusion: NEVER sync or touch 02-spec/21-*.
2. Bump Script Protection: NEVER overwrite version bump scripts (bump-version.mjs, bump_versions.py, 37-bump-version.py, version.json).
3. Additive-Only AI Scripts: Copy new scripts; diff and preserve existing modified scripts.
4. Memory & Plans Protection: NEVER touch .ai-memory/memory/ or .ai-memory/plans/ in target repos.
5. Zero Secrets Leakage: NEVER sync .env or credentials.

### Assigned Subtasks:
- Subtask 1: .ai-memory/plans/subtasks/nn-<slug>/01-<name>.md

### Concurrency-Safe SQLite Action Logging:
- Database: <databasePath>
- Claim: python 03-ai-scripts/46-agent-sqlite-task-manager.py claim --db <databasePath> --agent "Worker <NN>"
- Before action: python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <databasePath> --subtask-id <id> --agent "Worker <NN>" --action "write_to_file" --file "<path>" --details "<action>"
- On complete: python 03-ai-scripts/46-agent-sqlite-task-manager.py complete --db <databasePath> --subtask-id <id> --agent "Worker <NN>" --evidence "PASS <evidence>"
- On failure: python 03-ai-scripts/46-agent-sqlite-task-manager.py fail --db <databasePath> --subtask-id <id> --agent "Worker <NN>" --reason "<reason>"

### Output Contract:
Write subtask output to .ai-memory/plans/subtasks/nn-<slug>/01-<name>.json and reply with this JSON block, then stop:
{
  "task": "Task-01",
  "status": "DONE",
  "filesChanged": ["<path1>", "<path2>"],
  "checks": "<command> -> exit <code>, <files scanned>",
  "acceptance": { "ac1": "PASS <evidence>" },
  "assumptions": [],
  "blockers": []
}
```

---

## 8. Final Report Format (Strict Vertical Lines)

```markdown
### Synchronization Completion Summary

- ✅ **Task-01: Target Codebase Preflight & Safety Audit** — `[Completed]` — [preflight exit 0, backup branches pushed]
- ✅ **Task-02: Multi-Repository Synchronization & Safe Mirroring** — `[Completed]` — [prompts, skills, guidelines mirrored; 5 boundaries enforced]

### Modified / Synchronized Target Repositories

- [repo-1-slug] (base branch: main, backup branch: backup/pre-v3-nsteps-sync-..., release: vX.Y.Z)
- [repo-2-slug] (base branch: main, backup branch: backup/pre-v3-nsteps-sync-..., release: vX.Y.Z)

### Steps Used

- Step x / N (Phase 1: y / PHASE_1_BUDGET, Phase 2: z / PHASE_2_BUDGET), Wave k / WAVES

### Implementation Confidence Score

- Confidence: [passed checks / total checks]
- Rationale: [Verified evidence across all 5 non-negotiable boundaries, targeted linters passing, zero regressions]
```

---

## 9. Targeted Verification Checks (R2)

Run on changed files/folders only:

- **Coding Guidelines & Boolean Linter:** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`
- **Relative Path Linter:** `python linter-scripts/check-relative-paths.py`
- **Prompts & Spec Index Linter:** `python linter-scripts/check-prompts-loaded.py`
- **Doc Path Linter:** `python 03-ai-scripts/22-doc-path-linter.py <folder>`
- **Sequence Integrity Linter:** `python linter-scripts/check-sequence-integrity.py`
- **Forbidden Strings Check:** `python linter-scripts/check-forbidden-strings.py`

---

## 10. AI Fix Scripts Memory (Reusable Tooling)

- [ ] [/goal](slashCommand;goal) Reuse First: Scanned and learned `03-ai-scripts/readme.md` before writing temporary code.
- [ ] Strict In-Repository Execution: All Python scripts executed strictly within the repository context.
- [ ] Multi-Repo Sync Automation: Leverage `python 03-ai-scripts/38-sync-prompts-skills-scripts.py` as primary synchronization automation engine.
- [ ] Strict .ai-memory/ Folder Storage: All helper scripts, local runners, and linters stored in `03-ai-scripts/`.

---

## 11. Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Verified whatever directives or parameters are provided above this prompt are treated as highest priority.
- [ ] ZERO HARDCODED PATHS MANDATE: Verified source and target repository paths are supplied dynamically from caller and never hardcoded in prompt or scripts.
- [ ] BOUNDARY 1 (SPEC 21 EXCLUSION): Verified `02-spec/21-*` directories were NEVER touched, synced, or overwritten in target repositories.
- [ ] BOUNDARY 2 (BUMP SCRIPT PROTECTION): Verified version bumping scripts (`bump-version.mjs`, `bump_versions.py`, `37-bump-version.py`, `version.json`) in target repos were NEVER overwritten.
- [ ] BOUNDARY 3 (ADDITIVE-ONLY AI SCRIPTS): Verified new AI scripts were added cleanly, and existing customized scripts were diffed and preserved, never overwritten blindly.
- [ ] BOUNDARY 4 (MEMORY & PLANS SAFE): Verified `.ai-memory/memory/` and `.ai-memory/plans/` in target repositories were never modified or overwritten.
- [ ] BOUNDARY 5 (ZERO SECRETS LEAKAGE): Verified no `.env` or credential files were synchronized across repositories.
- [ ] MANDATORY SUBAGENT SPAWNING GATE: Verified `invoke_subagent` (`A = 2, H = 2`) was explicitly called in Phase 1 and Phase 2. Zero solo execution allowed.
- [ ] NO SUBAGENT FILE WRITE COLLISIONS: Never assign two subagents to create or modify the same file path simultaneously.
- [ ] NO READ-ONLY SUBAGENT FILE WRITES: Subagents of `TypeName: "research"` have read-only tools and cannot write files.
- [ ] NO INVALID SUBAGENT TYPENAMES: `TypeName` MUST strictly be `"research"` or `"self"`.
- [ ] NO SUBAGENT GIT EXECUTION: Subagents never run git commands in shared workspaces.
- [ ] NO TEST RUNNING (TOTAL BAN): Never run test suites (`go test`, `pytest`, runner scripts) during routine turns.
- [ ] NO BUILD CHECKING (TOTAL BAN): Never run build commands (`go build`, `npm run build`) to verify compilation.
- [ ] NO RAW GIT COMMITS OR CONVENTIONAL COMMIT PREFIXES: All staging and committing executed via GitMap with hyphen format (`gitmap cpf "<module> - <summary>"`). TOTAL BAN on colons inside message arguments.
- [ ] NO INTERMEDIATE COMMITS: Zero commits during Phase 1 or mid-Phase 2. Commit once at Phase 3 final step.
- [ ] NO RAPID CI/CD POLLING: Query `gitmap pipeline-ai status --json` with adaptive sleep based on `etaSeconds`.
- [ ] NO POWERSHELL OR SHELL SEARCHES (TOTAL BAN): Never run `Select-String`, `Get-ChildItem -Recurse`, `git grep`, or `grep`. Use GitMap high-speed search tools.

---

## 12. Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

- [ ] Master Guidelines: Fully enforced every file in `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] Error Management: Enforced `02-spec/03-error-manage/` using domain-specific `*appfault.AppError`, never generic error.
- [ ] Boolean Conventions: All booleans begin with `is` or `has` only. Implicit checks only; no `== true`. No mixed polarity (`if isA && !isB` is banned).
- [ ] Semantic Naming: Zero generic garbage names (`temp`, `data`, `obj`).
- [ ] Line Endings & Encoding: Strictly Unix LF (`\n`) and UTF-8 without BOM.
- [ ] Function Sizing: Functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Strict Relative Git Paths & Lowercase: Zero absolute paths or `file:///` URIs. All new filenames strictly lowercase.

---

## 13. Anti-Hallucination & Blast Radius Checklist

- [ ] Echo Back the Spec: Verified Acceptance Criteria from the Spec file verbatim.
- [ ] Pre-Commit Diff Proof: Verified `git status` shows actual modified files before committing.
- [ ] No Placeholder Search: Confirmed zero `TODO` or `\[.*\]` placeholders remain in modified files.
- [ ] Index Sync Deadman Switch: Registered `23-sync/` in `01-prompts/readme.md` and `.ai-memory/prompts.md`.
- [ ] Blast Radius Acknowledgment: Global search across codebase performed via `gitmap aum search "<symbol>" [dir]`.
- [ ] Final Step Commit & Push Verified: Staged and committed all changes atomically via GitMap using hyphen format (`gitmap cpf "<module> - <summary>"`).

---

## 14. Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] MANDATORY FINAL COMMIT & PUSH VIA GITMAP (ANYHOW): At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans consolidated, use GitMap semantic commit commands exclusively: `gitmap cpf "<module> - <summary>"` (features) or `gitmap cpb "<module> - <summary>"` (bug fixes). Mention and format as a hyphen `-` to separate module and summary; no need to provide a colon `:` in the GitMap `cpf`/`cpb` message argument as the colon is already provided automatically by GitMap (`Feature: ` or `Bug: `). TOTAL BAN on raw `git commit`, `git add -A`, colons `:` in the message argument, or conventional prefixes. ZERO intermediate commits during Phase 1 or Phase 2. Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] TOTAL BAN ON PER-FILE COMMITS: Never commit file-by-file. All modified files across the turn must be accumulated in the working tree and committed together in a single grouped atomic commit at the final step before pushing.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase readme.md files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
