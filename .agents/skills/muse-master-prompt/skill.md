---
name: muse-master-prompt
description: >-
  Executes the Unified Muse Master Onboarding, Fleet Orchestration & Canonical V6 Autonomous Execution prompt. Use when the user asks to run muse-master-prompt, onboard a Muse AI agent, authenticate via SSH, scan/clone/pull the repo fleet via GitMap, execute with Muse turbo mode, or apply the V6 task confirmation protocol.
---

# [V6] Muse Master Onboarding, Fleet Orchestration & Autonomous Execution Skill

Source prompt: [`01-prompts/27-muse-prompts/01-muse-master-prompt.md`](01-prompts/27-muse-prompts/01-muse-master-prompt.md)
Canonical Execution Engine: [`01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`](01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Discovery, Architecture Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt/skill (in the user preamble or header blocks) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously boot the session through the 9-step onboarding ritual, reporting back after every single step with action taken and verification output, authenticate with GitHub/GitMap (SSH preferred, with instant fast-path bypass if already logged in), verify/install GitMap across platforms (`curl` / `irm`), update GitMap first, clone the 8 default repositories (`coding-guidelines-v24` starting point, stale repo banned), seed the 5 side chats from `chat-seeds/`, discover secondary repos from `02-repos.md` with user approval, pull all handpicked repos using `gitmap pa`, ingest master prompt and custom skills, sync memory via `gitmap llm train muse --days 7`, connect email via `go-email-reader` under standing SMTP pre-approval, scan everything into gitmap (`gitmap scan`), emit the 9-item final report block, showcase confirmed task breakdown during Turn 1, declare RUNNING state with ETA, spawn autonomous subagents via `invoke_subagent` (A = 2, H = 2; solo execution of multi-part work is an auto-reject failure) in disjoint file boxes using GitMap commands as primary, enforce coding guidelines to 100%, split files >15 MB into 7-Zip volumes, and finish each task with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format pushed immediately.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives are provided ABOVE this prompt outrank everything below. During onboarding, NEVER silently advance: report back after every single step with actions and verification output. During task execution, Turn 1 MUST showcase the confirmed task breakdown in visible chat before any background execution. Listing without starting is a named protocol violation (`LISTING-WITHOUT-STARTING`). Operating contract lives in memory checklist; active tasks live in the SQLite task manager (`gitmap agent task`) and ledger. Repositories are sacred: NEVER delete a repo. Never run test suites locally: verify via CI (`gitmap pe -t`).

[/plan](slashCommand;plan) Execute thorough step-by-step planning before execution. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 2` agents (`invoke_subagent`) at every stage of multi-part work:
  1. **Planning Step (A = 2 `research` subagents):** Spawn 2 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified plan.
  2. **Spec Step (A = 2 `self` subagents or lead):** Spawn 2 subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks.
  3. **Execution Step (A = 2 `self` worker subagents):** Spawn 2 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
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
     - `b. SSH`: `ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N ""`, display public key in chat for user to add, then retry.
   - Hard gate: Do NOT proceed to step 3 until login succeeds.
   - **REPORT:** `gh auth status` output showing authenticated account.
3. **Step 3: Clone 8 Default Repos + Create 5 Side Chats:**
   - Clone with `gitmap clone https://github.com/alimtvnetwork/<repo> ~/workspace/repos/<dir>`:
     - `gitmap-v28` -> `gitmap-v28 repo` (side chat)
     - `Antigravity-Manager` -> `Antigravity-Manager repo` (side chat)
     - `cat-my-v12` -> `cat-my-v12 repo` (side chat)
     - `alim-seo-writing` -> `alim-seo-writing repo` (side chat)
     - `story-telling-v1` -> `story-telling-v1 repo` (side chat)
     - `coding-guidelines-v24` -> — (prompt lives here)
     - `repo-cache` -> — (memory vault)
     - `repo-secrets` -> — (secrets vault)
   - Seed the 5 side chats from `repo-cache/12-muse-ai/chat-seeds/<name>.md`.
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
   - Output `### 📋 Confirmed Task Breakdown & Requirement Ingestion` in visible chat (`Task-01`, `Task-02`, ... each with `State: [IN PROGRESS]` and `Understood: [YES]`, `Slug: <root-slug> - Task NN`).
   - In the **EXACT SAME TURN**, invoke your first tool call (SQLite task DB initialization, ledger creation, preflight checks). NEVER emit text alone.
3. **RUNNING Declaration & 5-Minute Status Pings:**
   - Print `RUNNING — Task-01, Task-02 — ETA ~45 min`. Ping status every 5 minutes during execution. Total ban on `LISTING-WITHOUT-STARTING`.
4. **Pull First in Target Repo:** Run `gitmap pull` (fallback `git pull`) before touching any files.
5. **SQLite Task DB & Ledger Preflight:**
   - Run `gitmap agent task init --name "<task name>" --budget 300`.
   - Maintain human-readable `.ai-memory/temp-agents/<slug>/ledger.md`.
6. **Mandatory Multi-Agent Execution:**
   - Dispatch `A = 2` subagents in disjoint file boxes (`TypeName: "research"` for discovery, `TypeName: "self"` for execution).
   - Worker Git Ban: Subagents NEVER run git commands (`git add`, `git commit`, `git status`) to avoid `.git/index.lock` collisions.
7. **Verify, Pre-Commit Secrets Gate & Atomic Commit:**
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
