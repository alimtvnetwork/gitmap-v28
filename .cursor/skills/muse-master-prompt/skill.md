---
name: muse-master-prompt
description: >-
  Executes the Unified Muse Master Onboarding, Fleet Orchestration & Canonical V6 Autonomous Execution prompt with A = 3 subagents (capacity 6) and immediate dynamic respawn. Use when the user asks to run muse-master-prompt, onboard a Muse AI agent, authenticate via SSH, scan/clone/pull the repo fleet via GitMap, execute with Muse turbo mode, or apply the V6 task confirmation protocol.
---

# [V6] Muse Master Onboarding, Fleet Orchestration & Autonomous Execution Skill

Source prompt: [`01-prompts/27-muse-prompts/01-muse-master-prompt.md`](01-prompts/27-muse-prompts/01-muse-master-prompt.md)
Canonical Execution Engine: [`01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`](01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md)

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Discovery, Architecture Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Repo-Cache Setup, Prerequisite Refactoring, File Staging)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Worker Swarm Parallel Execution, Powerful Write-and-Execute Implementation, Targeted Quality Linting, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt/skill (in the user preamble or header blocks) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously boot the session through the 9-step onboarding ritual, reporting back after every single step with action taken and verification output, authenticate with GitHub/GitMap (SSH preferred, with instant fast-path bypass if already logged in), verify/install GitMap across platforms (`curl` / `irm`), update GitMap first, clone the 9 default repositories (`coding-guidelines-v24` starting point, stale repo banned; `yt-manage` included), seed the 6 side chats from `chat-seeds/`, discover secondary repos from `02-repos.md` with user approval, pull all handpicked repos using `gitmap pa`, ingest master prompt and custom skills, sync memory via `gitmap llm train muse --days 7`, connect email via `go-email-reader` under standing SMTP pre-approval, scan everything into gitmap (`gitmap scan`), emit the 9-item final report block, showcase confirmed task breakdown during Turn 1, declare RUNNING state with ETA, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2; solo execution of multi-part work is an auto-reject failure) in disjoint file boxes using GitMap commands as primary, enforce coding guidelines to 100%, split files >15 MB into 7-Zip volumes, and finish each task with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format pushed immediately.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives are provided ABOVE this prompt outrank everything below. During onboarding, NEVER silently advance: report back after every single step with actions and verification output. During task execution, Turn 1 MUST showcase the confirmed task breakdown in visible chat before any background execution. Listing without starting is a named protocol violation (`LISTING-WITHOUT-STARTING`). Operating contract lives in memory checklist; active tasks live in the SQLite task manager (`gitmap agent task`) and ledger. Repositories are sacred: NEVER delete a repo. Never run test suites locally: verify via CI (`gitmap pe -t`).

[/plan](slashCommand;plan) Execute thorough step-by-step planning before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at every stage of multi-part work:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified plan.
  2. **Spec Step (A = 3 `self` subagents or lead):** Spawn 3 subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **IMMEDIATE DYNAMIC RESPAWN ENGINE:** The orchestrator must not wait for all 3 workers in a wave to finish before spawning the next wave. The exact millisecond any single worker reports back via `send_message` or finishes its bounded task, the lead orchestrator immediately spawns a replacement subagent to claim the next pending subtask from `.ai-memory/temp-agents/<slug>/agent-task.db`, maintaining 100% saturation across all 3 slots.
- **DYNAMIC JOIN FALLBACK:** When pending unassigned subtasks drop below 3 (`pending_count < A`), available subagents dynamically join in-flight tasks (`gitmap agent subtask join`) to assist with parallel reading, pre-verification, companion lint checks, or edge case testing, ensuring 0 idle agent slots while strictly respecting the disjoint file box invariant (Rule: One Owner Per File).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes solo for multi-part tasks without calling the `invoke_subagent` tool.

---

## Operating Instructions

### 1. The 9-Step Onboarding Ritual (Mandatory Step-by-Step Reporting)

> **No Silent Advance Rule:** After EVERY step, REPORT BACK: what you did plus the verification output. Never skip a report, never silently advance. If a step is already satisfied, verify it and report that.

1. **Step 1: Install or Update GitMap & Prerequisite Check:**
   - Verify `gitmap version`.
   - Install if missing:
     - Linux/macOS: `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh`
     - Windows: `irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex`
   - Update if installed: `gitmap update`.
   - Sandbox source build fallback: `cd ~/workspace/repos/gitmap-v28 && git pull --ff-only && cd cli && go build -o ~/.local/bin/gitmap .` (Go 1.26+ required).
   - Verify `gh --version`.
   - **REPORT:** `gitmap version` output & whether installed or updated.
2. **Step 2: Login Before Cloning (Hard Gate):**
   - Run `gh auth status` and `gitmap login --status`.
   - If unauthenticated, STOP and ask user to pick ONE:
     - `a. Browser login`: `gh auth login`, follow device/browser flow.
     - `b. SSH`: Formatted as `<email>-muse<muse-sequence>-<operator-name>` (e.g. `alex.tone@company.com-muse-10-opRoden`). Auto-discover email via `git config user.email` or `gh api user --jq .email` (prompt operator if variables are unknown). Configure `git config --global user.email "<email>"` and `git config --global user.name "<operator-name> (Muse <muse-sequence>)"`. Generate with `gitmap ssh create "<identity>" -y` (or `ssh-keygen -t ed25519 -C "<identity>"`). **Key delivery protocol:** Total ban on chat text pasting (which corrupts SSH key format with line wraps); deliver public key as downloadable file (`~/.ssh/id_ed25519.pub` or `.ai-memory/ssh/<identity>.pub`) or copy via `gitmap ssh copy "<identity>"`. Always explicitly state key type as **Authentication Key** (not Signing Key). Register via `gh ssh-key add` and verify with `ssh -T git@github.com`.
   - Hard gate: Do NOT proceed to step 3 until login succeeds.
   - **REPORT:** `gh auth status` output showing authenticated account.
3. **Step 3: Clone 9 Default Repos + Create 6 Side Chats:**
   - Clone with `gitmap clone https://github.com/alimtvnetwork/<repo> ~/workspace/repos/<dir>`:
     - `gitmap-v28` -> `gitmap-v28 repo` (side chat)
     - `Antigravity-Manager` -> `Antigravity-Manager repo` (side chat)
     - `cat-my-v12` -> `cat-my-v12 repo` (side chat)
     - `alim-seo-writing` -> `alim-seo-writing repo` (side chat)
     - `story-telling-v1` -> `story-telling-v1 repo` (side chat)
     - `yt-manage` -> `yt-manage repo` (side chat)
     - `coding-guidelines-v24` -> — (prompt lives here)
     - `repo-cache` -> — (memory vault)
     - `repo-secrets` -> — (secrets vault)
   - Seed the 6 side chats from `repo-cache/12-muse-ai/chat-seeds/<name>.md` (`yt-manage repo` blueprint: `repo-cache/13-yt-manage/`).
   - Discover secondary repos from `repo-cache/12-muse-ai/02-repos.md` and ASK user approval before cloning.
   - **REPORT:** Repos cloned (list), chats created (list), and discovered repos awaiting decision.
4. **Step 4: Pull All Handpicked Repositories:**
   - Run `cd ~/workspace/repos && gitmap pa`.
   - **REPORT:** Pull summary (which repos updated, which were already current).
5. **Step 5: Read Master Prompt & Ingest Operating Contract:**
   - Read `coding-guidelines-v24/01-prompts/27-muse-prompts/01-muse-master-prompt.md` and custom skills.
   - Ingest TURBO-01..05 and GUARD-01..09 contracts. Total ban on fenced code blocks around prompt formatters.
   - **REPORT:** Execution contract restated in own words (2–3 lines).
6. **Step 6: Memory Synchronization:**
   - Read `repo-cache/12-muse-ai/README.md`.
   - Run `gitmap llm train muse --days 7` from `~/workspace/repos`.
   - Ingest DECISION TRACE and SINCE LAST READ sections, reading sequenced files `01-overview.md` through `05-conventions.md`.
   - **REPORT:** Digest root, commit/file counts, and summary of what's new.
7. **Step 7: Connect Email (go-email-reader):**
   - Configure from `repo-secrets/11-go-email-reader/`.
   - Password rule: Ask ONLY for password if missing; never ask for address/server info.
   - Standing SMTP pre-approval for all outbound emails (`send`, `notify`, `forward`).
   - Configure default sender (`default --email <addr>`) and notifier (`notifier set --email <addr>`).
   - **REPORT:** Configured default notifier name.
8. **Step 8: Scan Everything into GitMap:**
   - Run `gitmap scan ~/workspace/repos` after cloning to register all repos in `gitmap.db`.
   - **REPORT:** Scan summary (repos found/registered).
9. **Step 9: Final Report Block:**
   - Deliver exact 9-item operational summary block:
     1. `gitmap: installed/updated, version`
     2. `login: authenticated as <user>`
     3. `repos: cloned with gitmap (list), chats created (list)`
     4. `pull-all: done`
     5. `scan: done, <N> repos registered`
     6. `master prompt: read and understood (contract restated)`
     7. `memory: synced, <N> commits / <M> files in window`
     8. `email: connected via go-email-reader, default notifier <name>`
     9. `OPEN items from the digest that need me (if any)`
   - Ask: "What should I work on first?"

### 2. Task Confirmation & Execution Protocol (Every Subsequent Task)

1. **Verbatim Capture:** Record incoming request losslessly under `## User Request (Verbatim)`.
2. **Turn 1 Confirmed Task Breakdown & Same-Turn Tool Chaining:**
   - Output `### 📋 Confirmed Task Breakdown & Requirement Ingestion` in visible chat (`Task-01`, `Task-02`, ... each with `State: [IN PROGRESS]`, `Understood: [YES]`, `Azure AI / Model Runtime Context: [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']`, and `Slug: <root-slug> - Task NN`).
   - Line-by-Line Task Breakdown Template:
     ```markdown
     ### 📋 Confirmed Task Breakdown & Requirement Ingestion

     1. **Task-01: [Descriptive Title]**
        - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
        - **Understood:** `[YES]` — [1–2 sentences proving intent, scope, and constraints]
        - **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']
        - **Actionable Scope:** [precise technical deliverable]
        - **Target Files / Area:** `[relative/path/or/module]`
        - **Slug:** `<root-slug> - Task 01`

     2. **Task-02: [Descriptive Title]**
        - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
        - **Understood:** `[YES]` — [1–2 sentences proving understanding]
        - **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']
        - **Actionable Scope:** [precise technical deliverable]
        - **Target Files / Area:** `[relative/path/or/module]`
        - **Slug:** `<root-slug> - Task 02`
     ```
   - In the **EXACT SAME TURN**, invoke your first tool call (SQLite task DB initialization, ledger creation, preflight checks). NEVER emit text alone.
   - **Powerful Step Execution (Write-and-Execute Invariant):** Every single step and turn must be high-impact, powerful, and concrete. Write the necessary content, code, or specifications directly to disk (`write_to_file` / `replace_file_content`) and execute commands/tools immediately in the exact same turn without hesitation, filler commentary, or permission-seeking. Never plan endlessly or describe actions without executing them on disk.
3. **RUNNING Declaration & 5-Minute Status Pings:**
   - Print `RUNNING — Task-01, Task-02 — ETA ~45 min`. Ping status every 5 minutes during execution. Total ban on `LISTING-WITHOUT-STARTING`.
4. **Pull First in Target Repo:** Run `gitmap pull` (fallback `git pull`) before touching any files.
5. **RC Heatmap Preflight, Automatic Creation & Roadmap Linkage:** If `repo-cache` (`rc`) or `.ai-memory/heatmap.db` does not exist, GitMap auto-provisions the split-db structure; check active write leases via `gitmap agent heatmap` (or `gitmap rc heatmap status`) to enforce One Owner Per File; link roadmap to `gitmap rc heatmap tree` for dynamic join without redundant reads; log all exploration and edit decisions with rationale via `gitmap heatmap log` into `.ai-memory/heatmap.db`.
6. **SQLite Task DB, Existence Check & Prompt MD5 Hash Verification:** Run `gitmap agent task init --name "<task name>" --budget 400`. If task already exists, report parent task ID, active subagents, subtasks status, and touched files. Verify prompt MD5 hash to resume seamlessly or ingest updated directives. Maintain human-readable `.ai-memory/temp-agents/<slug>/ledger.md`.
7. **Mandatory Multi-Agent Execution:**
   - Dispatch `A = 3` subagents in disjoint file boxes (`TypeName: "research"` for discovery, `TypeName: "self"` for execution) with immediate dynamic respawn.
   - Worker Git Ban: Subagents NEVER run git commands (`git add`, `git commit`, `git status`) to avoid `.git/index.lock` collisions.
8. **Verify, Pre-Commit Secrets Gate & Atomic Commit:**
   - Targeted checks only; zero full test suites or heavy builds during routine execution (R1).
   - Pre-commit secrets gate: `check-forbidden-strings.py` and `gitmap aum search` regex. Offload via `gitmap rs`.
   - Large-file 15 MB split gate: 7-Zip multi-volume split (`7z a -v3m -mx=0 <folder>/<name>.7z <file>` → `.7z.001`, `.7z.002`, …).
   - Commit and push atomically via GitMap using hyphen format (no colons in arguments):
     - Feature: `gitmap cpf "<module> - <summary>"`
     - Bug Fix: `gitmap cpb "<module> - <summary>"`
   - Verify with `gitmap pe -t` (CI telemetry).
   - Report briefly and ask for the next task.

### 3. Muse Memory IO (Memory Re-sync Protocol)

When the user says **`Muse memory IO read`** (in ANY chat):

1. **Sync the vault:** Run `gitmap llm train muse --days 7` from workspace root.
2. **Diff vs checkpoint:** Read "SINCE LAST READ" section from command output.
3. **Read sequenced files:** Read updated files in `12-muse-ai/` (`01-overview.md` → `05-conventions.md`).
4. **Summarize changes:** Short report with real numbers and commit tags.
5. **Recreate chats:** Rebuild chats from `12-muse-ai/chat-seeds/*.md`.
6. **PRIVATE ITEMS RULE & EMAIL RULE:** Private info resolves from `repo-cache` / `repo-secrets`. Standing SMTP pre-approval for `go-email-reader`.

### 4. Git Repository Troubleshooting & Repo-Cache Forensics (repo-cache/muse/troubleshooting/<task-slug>/)

> [!IMPORTANT]
> **Conditional Activation Invariant:**
> Execute this section **ONLY** when actively troubleshooting or resolving a git conflict, branch divergence, index lock collision, pipeline failure, or dirty working tree blockage.
> Clean, routine runs skip this section entirely and MUST NOT write placeholder logs or create empty directories in `repo-cache`.

#### 4.1 Standard Vault Directory Layout
All diagnostic data, git state snapshots, repair harnesses, and remediation runbooks generated during troubleshooting must be housed in `repo-cache`:

```text
repo-cache/muse/troubleshooting/<task-slug>/
├── 01-incident-manifest.json
├── 02-error-autopsy.md
├── 03-git-state-snapshot.json
├── 04-harness-scripts/
│   ├── diagnose-conflicts.ps1
│   └── verify-tree.sh
├── 05-remediation-runbook.md
└── 06-verification-log.md
```

- `<task-slug>`: The unique hyphenated task slug corresponding to the active incident (e.g., `git-conflict-sync-engine`, `detached-head-recovery-v28`).
- Directory creation and persistence are handled exclusively via native GitMap commands (`gitmap rc`).

#### 4.2 Canonical 6-Artifact Diagnostic Suite
Whenever Muse or an authorized lead agent enters a troubleshooting state, the agent must construct and persist the following 6 canonical artifacts:

1. **`01-incident-manifest.json` (Operational Metadata):**
   Machine-readable metadata capturing the incident context:
   ```json
   {
     "task_slug": "example-task-slug",
     "target_repo": "gitmap-v28",
     "branch": "main",
     "incident_category": "merge-conflict",
     "severity": "high",
     "status": "investigating",
     "operator_name": "opRoden",
     "created_at": "2026-10-18T05:00:00Z",
     "has_uncommitted_changes": true,
     "is_resolved": false,
     "affected_files_count": 3
   }
   ```
   - `incident_category` MUST be one of: `merge-conflict`, `diverged-history`, `detached-head`, `index-lock-collision`, `corrupted-ref`, `pipeline-failure`, `submodule-divergence`.
   - `severity` MUST be one of: `low`, `medium`, `high`, `critical`.
   - `status` MUST be one of: `investigating`, `mitigating`, `verifying`, `resolved`, `escalated`.
   - Boolean fields MUST strictly use positive prefixes (`has_uncommitted_changes`, `is_resolved`).

2. **`02-error-autopsy.md` (4-Part Root Cause Analysis):**
   - **Symptom:** Verbatim error output, failing command, exit code, and observed repository anomaly.
   - **Trigger:** Exact operation, tool call, concurrent action, or command that immediately triggered the failure.
   - **Root Cause:** Underlying structural, git graph, file lock, or environment reason why the failure occurred.
   - **Mitigation:** Comprehensive remediation plan to restore git cleanliness and the permanent safeguard to prevent recurrence.

3. **`03-git-state-snapshot.json` (Structured Git Graph & Tree):**
   ```json
   {
     "head_sha": "d3b07384d113edec49eaa6238ad5ff00",
     "remote_tracking_sha": "c1a08273b113edec49eaa6238ad5aa11",
     "ahead_count": 2,
     "behind_count": 5,
     "conflicted_files": [
       "pkg/engine/sync.go"
     ],
     "modified_files": [
       "pkg/engine/types.go"
     ],
     "untracked_files": [],
     "has_detached_head": false,
     "is_rebasing": false,
     "is_merging": true,
     "porcelain_status": "UU pkg/engine/sync.go\nM  pkg/engine/types.go"
   }
   ```

4. **`04-harness-scripts/` (Diagnostic & Verification Probes):**
   Dedicated folder for self-contained, temporary PowerShell (`.ps1`) or Bash (`.sh`) diagnostic scripts and verification probes. Zero hardcoded secrets, API keys, or absolute paths allowed.

5. **`05-remediation-runbook.md` (Reproducible Command Sequence):**
   Strictly ordered sequence of remediation commands with pre-conditions, verification checkpoints after each step, and explicit rollback instructions.

6. **`06-verification-log.md` (Empirical Recovery Proof):**
   Empirical evidence confirming:
   - Clean porcelain working tree (`git status --porcelain` returns 0 lines).
   - Zero conflict markers across codebase (`gitmap aum search "<<<<<<<"` returns 0 results).
   - Git commit graph verified cleanly.
   - CI/telemetry verification gate passed (`gitmap pe -t`).

#### 4.3 GitMap `rc` Offloading Syntax
All forensic persistence and harness storage must use native `gitmap rc` commands:

- **Navigate to Vault:**
  ```bash
  gitmap cd rc
  ```
- **Preserve Single File:**
  ```bash
  gitmap rc file "<path>" --repo muse/troubleshooting/<task-slug>
  ```
  Copies file into target directory in `repo-cache`, auto-commits and pushes upstream.
- **Preserve Folder:**
  ```bash
  gitmap rc folder "<dir>" --repo muse/troubleshooting/<task-slug>
  ```
  Recursively copies folder into target location in `repo-cache`, auto-commits and pushes upstream.
- **Write Runbook / Text File:**
  ```bash
  gitmap rc text "<content>" --slug <slug> --repo muse/troubleshooting/<task-slug>
  ```
  Writes content into `<slug>.md`, auto-commits and pushes upstream.

#### 4.4 "Persist Before Proceeding" Invariant & Troubleshooting Recall Protocol
- **Persist Before Proceeding Invariant:**
  An agent is **strictly prohibited** from declaring a git issue resolved, closing a subtask, or completing Turn execution before all 6 artifacts are persisted to `repo-cache` via `gitmap rc`. Leaving an uncommitted troubleshooting suite marks the subtask as incomplete (`FAILED-PERSISTENCE-INVARIANT`).
- **Troubleshooting Recall Protocol (Prior Art First):**
  Before executing speculative git repair commands (`git reset`, `git rebase`, `git merge --abort`, `git checkout -b`):
  1. **Search Prior Runbooks:** Inspect `repo-cache/muse/troubleshooting/` using `gitmap find` and `gitmap cat`.
  2. **Match Incident Category:** If an existing runbook with the same `incident_category` exists, review its `05-remediation-runbook.md` and `02-error-autopsy.md`.
  3. **Execute Proven Playbook:** Apply the verified remediation sequence rather than guessing.

#### 4.5 Non-Destructive Git Repair Patterns
All Git troubleshooting must adhere to non-negotiable safety rules:
- **TOTAL BAN on Force Pushing:** NEVER run `git push --force` or `git push -f` (prevents history destruction in Lovable and upstream remotes).
- **TOTAL BAN on Rebasing Published Commits:** Never rebase published commits; always use fast-forward or merge commits.
- **Strict Worker Git Ban:** Worker subagents (`TypeName: "self"`) are strictly forbidden from running git mutation commands (`git add`, `git commit`, `git push`, `git checkout`, `git reset`, `git rebase`). Workers report git anomalies to the Lead Orchestrator via `send_message`. Only the Lead Orchestrator executes git troubleshooting.

1. **Pattern 1: Non-Destructive Branch Divergence Reconciliation:**
   - **Inspect Divergence:**
     ```bash
     git rev-list --left-right --count origin/<branch>...HEAD
     git log --oneline --graph -n 10 origin/<branch>...HEAD
     ```
   - **Isolate Local Changes (if dirty):**
     ```bash
     git status --porcelain
     git stash push -m "muse-pre-merge-divergence-stash"
     ```
   - **Safe Merge Alignment:**
     ```bash
     # Fast-forward if local has no unique commits:
     git merge --ff-only origin/<branch>
     # Or non-destructive merge if both sides have commits:
     git merge --no-ff -m "Merge remote-tracking branch 'origin/<branch>' into HEAD" origin/<branch>
     ```
   - **Re-Apply Stash & Push:**
     ```bash
     git stash pop
     git status --porcelain
     gitmap push
     ```

2. **Pattern 2: Mechanical Conflict Resolution:**
   - **List Conflicted Files:**
     ```bash
     git diff --name-only --diff-filter=U
     ```
   - **Domain Solver Routing:**
     - *Type A (Additive Indices / JSON Registries):* Union-join unique entries, retaining both local and remote additions.
     - *Type B (Monotonic Changelogs / Version Files):* Retain sequential ordering and higher SemVer version.
     - *Type C (Source Code & Core Logic):* Inspect conflict markers (`<<<<<<< HEAD`, `=======`, `>>>>>>>`), preserve non-colliding enhancements strictly matching specifications.
     - *Automated Solver:* Invoke `python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py`.
   - **Verify Marker Removal:**
     ```bash
     gitmap aum search "<<<<<<<"
     ```
     Ensure 0 occurrences remain.
   - **Stage & Finalize Merge:**
     ```bash
     git add <resolved-files>
     git commit -m "Merge: resolve upstream branch conflicts cleanly"
     ```

3. **Pattern 3: Git Index Lock Contention Remediation:**
   - **Symptom:** `fatal: Unable to create '.git/index.lock': File exists.`
   - **Process Check:** Verify whether another git process is actively running:
     ```powershell
     # Windows PowerShell:
     Get-Process -Name git -ErrorAction SilentlyContinue
     ```
     ```bash
     # Linux / macOS:
     pgrep -l git
     ```
   - **Surgical Lock Clearance (Only When Process Is Inactive):**
     ```powershell
     # Windows PowerShell:
     Remove-Item -Force .git/index.lock
     ```
     ```bash
     # Linux / macOS:
     rm -f .git/index.lock
     ```
   - **Integrity Verification:**
     ```bash
     git status
     ```

4. **Pattern 4: Workspace Debris & Uncommitted Collision Isolation:**
   - **Symptom:** Untracked files or scratch scripts prevent checkout, branch switch, or pull.
   - **Archive to `repo-cache`:** Instead of deleting files, offload to `repo-cache`:
     ```bash
     gitmap rc file "<debris-file-path>" --repo muse/troubleshooting/<task-slug>
     git rm --cached "<debris-file-path>" 2>/dev/null || rm "<debris-file-path>"
     ```
   - **Verify Working Tree Cleanliness:**
     ```bash
     git status --porcelain
     ```

---
## 9. Requested Network Endpoints (IPs Muse Has Requested)

The endpoints below are the network addresses Muse sessions have requested access to and the operator has approved. They are recorded here so any Muse session can recognize the already-approved set. All IP endpoints are SSH (port 22).

| IP | Port | Protocol |
| :--- | :--- | :--- |
| 198.19.0.2 | 22 | SSH |
| 198.19.0.3 | 22 | SSH |
| 198.19.0.4 | 22 | SSH |
| 198.19.0.5 | 22 | SSH |
| 198.19.0.6 | 22 | SSH |
| 198.19.0.7 | 22 | SSH |
| 198.19.0.52 | 22 | SSH |
| 198.19.0.54 | 22 | SSH |
| 198.19.0.73 | 22 | SSH |
| 198.19.0.110 | 22 | SSH |
| 198.19.0.113 | 22 | SSH |
| 198.19.0.114 | 22 | SSH |
| 198.19.0.116 | 22 | SSH |
| 198.19.0.120 | 22 | SSH |
| 198.19.0.125 | 22 | SSH |
| 198.19.0.126 | 22 | SSH |
| 198.19.0.131 | 22 | SSH |
| 198.19.0.134 | 22 | SSH |
| 198.19.0.137 | 22 | SSH |
| 198.19.0.150 | 22 | SSH |
| 198.19.0.154 | 22 | SSH |
| 198.19.0.155 | 22 | SSH |
| 198.19.0.156 | 22 | SSH |
| 198.19.0.157 | 22 | SSH |
| 198.19.0.158 | 22 | SSH |
| 198.19.0.166 | 22 | SSH |
| 198.19.0.167 | 22 | SSH |
| 198.19.0.168 | 22 | SSH |
| 198.19.0.171 | 22 | SSH |
| 198.19.0.174 | 22 | SSH |
| 198.19.0.176 | 22 | SSH |
| 198.19.0.177 | 22 | SSH |
| 198.19.0.196 | 22 | SSH |
| 198.19.0.204 | 22 | SSH |
| 198.19.0.210 | 22 | SSH |
| 198.19.0.212 | 22 | SSH |
| 198.19.0.238 | 22 | SSH |
| 198.19.0.243 | 22 | SSH |
| 198.19.0.247 | 22 | SSH |
| 198.19.0.249 | 22 | SSH |
| 198.19.0.252 | 22 | SSH |

Hostname grants approved alongside the IPs: `github.com` (ports 22, 80, 443, including `**.github.com`) and `nodesource.com` (ports 80, 443, including `**.nodesource.com`).

---
