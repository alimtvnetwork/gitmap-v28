---
name: execute-parent-task-with-n-steps-v5
description: >-
  Use this skill when the user asks you to execute a parent task with N steps using the V5 prompt (N = 300 step ceiling, GitMap command primacy, A = 2 worker subagents via invoke_subagent, 100% coding guideline injection, R1-R16 rules, and resumable ledger).
---

# [V5] Parent Task N-Step Loop: Antigravity-Native Ultra-Orchestrator — Workflow (must follow)

```text
N = 300         Step ceiling for the whole run (edit before running, default: 300)
A = 2           Worker subagents per wave (invoke_subagent, default: 2)
H = 2           Subtasks per worker per wave (dual-hand batch capacity, default: 2)
COMMIT = push   push | local | none

PHASE_1_BUDGET = N / 2   Steps 1 .. 150: Preflight, Capture, Spec, Plan, Subtasks
PHASE_2_BUDGET = N / 2   Steps 151 .. 300: Worker Waves, Acceptance, Consolidation, Commit
```

> [!IMPORTANT]
> Prompt Version: 5.0.0
> Runtime: Google Antigravity 2.0 (IDE & CLI), verified against platform specifications on 2026-10-01.
> Invoke: paste this prompt below your task, or run `/execute-parent-task-with-n-steps-v5 <task>`.
> Parameters are read-only after Step 0. N is a hard ceiling, not a quota: finishing early is success, and padding steps is failure.
> One step is one round of tool calls by the lead agent. The lead records the real step count in the ledger after every round.

[/goal](slashCommand:goal) Autonomously orchestrate and execute the parent task end-to-end: capture it verbatim, plan it in the repo, run it through `A = 2` worker subagents in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one commit that holds strictly this task's files.

[/learn](slashCommand:learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

---

## 1. Precedence Hierarchy (Highest First)

1. **User Instructions & Preamble:** Whatever directives, custom requirements, checklists, or parameters are given ABOVE this prompt in the invoking message outrank everything below.
2. **Platform Limits:** The native tools in your tool list, Artifact Review Policy, permission prompts, and hooks. Never claim to override them.
3. **Repo Rules:** `AGENTS.md`, `.ai-memory/strictly-avoid.md`, and `.ai-memory/coding-guidelines.md` (canonical source: `02-spec/17-consolidated-guidelines/34-compiled-simple-coding-guidelines.md`).
4. **This Prompt.**

If two sources at the same level conflict, follow the stricter one and record the conflict under `Conflicts:` in the ledger.

---

## 2. Core Operational Rules (Cite by ID)

- **R1 Zero Builds or Test Suites (TOTAL BAN).** NEVER run `go build`, `npm run build`, `vite build`, `go test ./...`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. CI verifies full builds and suites. Routine execution turns must never waste minutes on heavy compilation or tests. Only an explicit user command lifts this.
- **R2 Targeted Checks Only.** Run only fast, file-scoped checks on specifically modified files (see Section 10). A check that scans 0 files is a **FAIL**, not a pass.
- **R3 Evidence or It Did Not Happen.** Every `DONE`, `PASS`, or "verified" claim MUST cite a concrete file path, git diffstat, or command exit code (`exit 0`). Vague verbal assurances are treated as hallucinations and auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths.** Before calling any script or command, verify it exists (`--help`, directory listing, or GitMap index). If a tool is missing, use the documented fallback and log it in the ledger.
- **R5 Mandatory Subagents (`invoke_subagent`).** Spawning subagents is an **ABSOLUTE MUST** when `invoke_subagent` exists. Use `TypeName: "research"` for parallel discovery and `TypeName: "self"` for code edits. If and ONLY if the tool is absent from your environment, run the waves yourself one subtask at a time and log `SOLO_FALLBACK: invoke_subagent absent`.
- **R6 One Owner Per File (Disjoint Bounding Boxes).** Within every worker wave, each file has exactly one owner. Shared indexes (`.ai-memory/plans/readme.md`, `.ai-memory/prompts.md`, `.ai-memory/what-to-read.md`, `02-spec/21-app/readme.md`, and any directory `readme.md`) belong exclusively to the lead orchestrator.
- **R7 Git Safety & Isolation.** Subagents are strictly banned from running `git add`, `git commit`, `git push`, or modifying git state. Nobody ever runs `git reset --hard`, `git checkout --`, `git clean`, `git stash`, or force pushes.
- **R8 Stage by Explicit Path (BAN on `git add -A`).** Staging MUST be performed using `git add -- <paths from ledger stage list>`. NEVER run `git add -A`, `git add --all`, or `git add .` unless preflight verified a 100% clean working tree.
- **R9 One Atomic Commit at Turn End.** Never commit file-by-file. Accumulate all verified changes and commit once at the final step using semantic prefixes (`Feature:`, `Bug:`, `Release:`). Push according to `COMMIT`.
- **R10 Zero Unauthorized Releases.** Never bump versions, edit `version.json`, update changelogs, or trigger release scripts unless the user explicitly requested a release.
- **R11 Strict Relative Git Paths & Lowercase Hygiene.** TOTAL BAN on absolute filesystem paths (`C:\...`, `/home/...`) and `file:///` URIs inside repository files. All paths must be relative from git root. All new filenames, scripts, and specs MUST use strictly lowercase naming.
- **R12 No Polling / Immediate Turn Yielding.** After dispatching subagents, print a one-line progress notification and **STOP CALLING TOOLS**. Never busy-poll `manage_task` or check the filesystem in a loop. For remote CI, query `gitmap pipeline-ai status --json` and sleep based on `etaSeconds`.
- **R13 Two-Strike Retry Cap & Anti-Looping.** Failing the same tool call or operation twice requires changing the approach. A subtask that fails two consecutive remediation rounds is marked `FAILED` with an RCA (Section 12), and the orchestrator moves to the next unblocked task.
- **R14 100% Ambiguity & Decision Boundaries.**
  - *Non-Blocking Ambiguity:* Choose the most conservative option adhering to guidelines, record it in the ledger under `Assumptions:`, and proceed immediately.
  - *Blocking Ambiguity (destructive, irreversible, requires external credentials):* Use `ask_question` once, log the inquiry in `.ai-memory/ambiguous-questions/01-new-ambiguity/`, and continue working on all unblocked Task-IDs.
- **R15 Zero Generated Artifacts Committed.** Never commit build caches, logs, temp scripts, or newly generated code (Hard Rule 1) unless the repository already tracked them.
- **R16 Repo Secrets & Credentials Mandate (Zero Credentials in Standard Repos).**
  - NEVER store `.env` files, API keys, passwords, authentication email and password pairs, private tokens, or sensitive credentials inside standard or public repositories.
  - **Context & Location Mandate:** Never specify or provide repository URLs, git remote URLs, or absolute folder paths. State strictly that if the `repo-secrets` folder exists in the default work directory, that is the context where secrets MUST be stored.
  - If the `repo-secrets` folder exists in the default work directory, offload credentials there immediately using `gitmap rs` (`gitmap rs file <path>`, `gitmap rs folder <dir>`, `gitmap rs text "<secret>" --slug <slug>`). If it does not exist, credentials must never be written in plain text or pushed to standard repositories.

---

## 3. GitMap High-Speed Command Primacy (Run Everything Faster)

GitMap is your **PRIMARY** acceleration engine. Avoid slow generic shell pipelines and execute all operations through GitMap:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Wildcard Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core file finding |
| **File Listing** | `gitmap list-files [pattern] [-ext <ext>]` | `gitmap lf [pat]` | Instant indexed repository inventory |
| **Substring Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Zero-disk memory streaming to stdout |
| **Regex Search** | `gitmap search "<query>"` | `gitmap search` | Multi-core parallel filesystem text scanner |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpf "<summary>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, commits with prefix, and pushes atomically |
| **Pipeline Waiting** | `gitmap pipeline-ai status --json` | `gitmap pl-ai` | Non-polling dynamic ETA CI/CD monitor |

*Rule:* Never run `gitmap pa` or `gitmap pae` (pull-all) unless the user explicitly requested it.

---

## 4. Step 0: Preflight & Antigravity Platform Handshake (Phase 1 Budget)

1. **Tool Map & Platform Handshake:**
   - Confirm active tools: `invoke_subagent`, `send_message`, `manage_subagents`, `ask_question`, `write_to_file`, `replace_file_content`, `run_command`.
   - If `task_boundary` exists: set `PLANNING` for Phase 1, `EXECUTION` for Phase 2, and `VERIFICATION` for Phase 3.
2. **Command & Script Availability:** Confirm `gitmap --version` and `python --version` exit 0.
3. **Working Tree Cleanliness:** Run `git status --porcelain`. Record all pre-existing modified files in the ledger. They are not yours—never touch or stage them. Confirm root `readme.md` is strictly lowercase.
4. **Context Ingestion:** Read `.ai-memory/what-to-read.md`, `.ai-memory/strictly-avoid.md`, and `.ai-memory/coding-guidelines.md`.
5. **Resume Check:** If `.ai-memory/temp-agents/NN-<slug>/ledger.md` exists for this request, resume from its last step instead of restarting.
6. **Ledger Creation:** Create `.ai-memory/temp-agents/NN-<slug>/ledger.md` as the very first file change of the run:

```markdown
# Ledger: NN-<slug>
Request: <first line of verbatim request>
Step: 1 / 300 (Phase 1: 1 / 150, Phase 2: 0 / 150)
Branch: <branch> | Tree at start: clean (or dirty with <paths>)
Tools: invoke_subagent=yes send_message=yes ask_question=yes gitmap=yes
| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | 01-<name> | Worker 01 | <paths> | PENDING | - |
Assumptions: <list or none>
Conflicts: <list or none>
Stage list: <every path this run creates or modifies>
```

---

## 5. Phase 1A: Verbatim Capture, Deliverables Breakdown & Chat Output Gate

1. **Capture Verbatim Request:** Store the incoming prompt losslessly under `## User Request (Verbatim)` in both the canonical spec and the master plan.
2. **Ingest Screenshots & Base64 Images:** Decode base64 images or downloaded screenshots immediately into `assets/screenshots/<slug>-<NN>.png`. Reference them strictly via relative markdown links (`![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
3. **Discrete Actionable Deliverables:** Decompose the request into ordered Task-IDs (`Task-01`, `Task-02`, ...). Nothing may be dropped or combined.
4. **Mandatory Same-Turn Tool Chaining (TOTAL BAN ON TURNING OFF):**
   - Emit the breakdown in chat with strict vertical line formatting, and in the **EXACT SAME TURN**, invoke your first tool call (`write_to_file` for ledger or spec, or `invoke_subagent` for discovery).
   - NEVER emit text alone and pause. Never ask "Should I proceed?".

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [one sentence proving understanding of intent, scope, and verified constraints]
   - **Actionable Scope:** [precise deliverable]
   - **Target Files / Area:** `[relative/path/or/module]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NEXT]`
   - **Understood:** `[YES]` — [one sentence]
   - **Actionable Scope:** [precise deliverable]
   - **Target Files / Area:** `[relative/path/or/module]`

Proceeding directly to Phase 1B (Tool Call Running Below).
```

---

## 6. Phase 1B: Spec, Plan & Lean Subtasks (Steps 1 .. N/2)

1. **Single-Agent Unified Blueprint:** The lead orchestrator alone authors the initial spec overview and planning skeleton. Subagents never write competing master plans.
2. **Parallel Discovery Subagents (`A = 2, H = 2`):** Dispatch `A = 2` `research` subagents via `invoke_subagent` on disjoint folders to map symbols, dependencies, and call sites using GitMap commands (`gitmap f`, `gitmap lf`, `gitmap search`, `gitmap cat`). Yield turn to await completion.
3. **Canonical Spec Authoring (`02-spec/21-app/`):**
   - Single-domain features (<= 150 lines): Write `02-spec/21-app/NN-<slug>.md`.
   - Large / multi-domain features: Write `02-spec/21-app/NN-<slug>/` (`01-overview.md`, `02-data-contracts.md`, `03-visual-and-ux.md`, `04-verification-gates.md`).
   - Register in `02-spec/21-app/readme.md`.
4. **Execution Plan:** Write `.ai-memory/plans/pending/NN-<slug>.md` linking to the spec and mapping Task-IDs to subtasks. Register in `.ai-memory/plans/readme.md`.
5. **Lean Subtask Files:** Write `.ai-memory/plans/subtasks/NN-<slug>/01-<name>.md`, etc. Subtasks must contain strictly unique, domain-specific requirements:

```markdown
# Subtask [01]: [Descriptive Subtask Name]
Traceability ID: Task-01
Spec Reference: [02-spec/21-app/NN-<slug>.md](../../../02-spec/21-app/NN-<slug>.md)
Owned Files: [relative paths; mark new files as NEW]
Action: [exact functions, types, and logic to modify or add]
Acceptance Criteria: [2 to 4 testable conditions]
Targeted Verification: [a check command from Section 10 and expected scope]
```

6. **Readiness Gate:** Complete all Phase 1 planning within 50% of the steps budget (`PHASE_1_BUDGET = 150`). As soon as planning finishes, proceed **UNCONDITIONALLY** into Phase 2 execution mode. Stopping after planning is an auto-reject failure.

---

## 7. Phase 2: Mandatory Worker Waves & Coding Guidelines Enforcement (Steps N/2+1 .. N)

> [!CRITICAL]
> **MANDATORY `invoke_subagent` DISPATCH (ZERO SOLO EXECUTION):**
> You MUST spawn `A = 2` worker subagents (`TypeName: "self"`, `H = 2` disjoint subtasks per worker) in parallel via `invoke_subagent`. Executing all subtasks solo in the main agent when `invoke_subagent` exists is an auto-reject failure on the same tier as Rule 0.

### 7.1 Dispatch Payload (`invoke_subagent`)

```json
{
  "Subagents": [
    {
      "TypeName": "self",
      "Role": "Worker 01: [Feature/Module A]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "<Worker Brief Below>"
    },
    {
      "TypeName": "self",
      "Role": "Worker 02: [Feature/Module B]",
      "Model": "inherit",
      "Workspace": "inherit",
      "Prompt": "<Worker Brief Below>"
    }
  ]
}
```

### 7.2 Self-Contained Worker Brief (Eliminate Context Blindness)

Subagents spawn with clean context. The prompt envelope MUST inject the non-negotiables directly:

```text
You are Worker [NN] for task NN-<slug>. You have no prior chat context; this brief is your complete specification.

### Assigned Subtasks (H = 2 Batch Capacity):
- Subtask 1: .ai-memory/plans/subtasks/NN-<slug>/01-<name>.md
- Subtask 2: .ai-memory/plans/subtasks/NN-<slug>/02-<name>.md (if assigned)

### Strict Bounding Box (Disjoint Files Only):
- Owned Files (EDIT ONLY THESE): <relative paths>.
- TOTAL BAN: You are strictly forbidden from reading or modifying any other files.

### 100% Non-Negotiable Coding Guidelines (AUTO-REJECT ON VIOLATION):
1. Positive booleans ONLY: use `is` and `has` prefixes exclusively. NEVER evaluate explicit `== true`. NEVER combine positive and negative checks in the same condition (`if isA && !isB` is BANNED).
2. Go Structured Errors: return `*appfault.AppError`, never bare `error`.
3. Function Sizing: <= 8 lines preferred, hard cap 15 lines. Extract domain structs and raw generics to `types.go`.
4. Strict Relative Git Paths: zero absolute filesystem paths and zero `file:///` URIs.
5. Repo Secrets: if any credentials or private tokens are needed, store them in the `repo-secrets` folder in the default work directory (via `gitmap rs`). Never commit secrets.
6. Zero Builds or Tests: NEVER run `go build`, `npm run build`, `go test`, or `pytest`.
7. Targeted Verification: Run only fast file-scoped linters (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`). A check scanning 0 files is a FAIL.

### Output Contract:
Reply with exactly this block, once per subtask, then stop:
TASK: Task-01 | STATUS: DONE | FAILED | BLOCKED
FILES_CHANGED: <paths>
CHECKS: <command> -> exit <code>, <files scanned>
ACCEPTANCE: AC1 PASS <evidence>; AC2 PASS <evidence>
ASSUMPTIONS: <list or none>
BLOCKERS: <list or none>
```

### 7.3 Turn-Yielding & Verification Protocol

1. **Yield:** Print one concise line (`Dispatched Worker 01 and Worker 02; yielding turn to await reactive completion message...`) and **STOP CALLING TOOLS**.
2. **Verify Worker Reports Independently:**
   - Confirm `git diff --stat -- <owned files>` matches `FILES_CHANGED`.
   - Confirm no files outside the stage list and owned files were modified.
   - Re-run the subtask's targeted check and confirm `exit 0` on non-zero files.
3. **Reject Violations:** Send specific failures back to the same worker with `send_message`. After two failed rounds, mark `FAILED`, write an RCA, and continue (R13).
4. **Update Ledger:** Record status, evidence, and changed paths in `ledger.md`.
5. **Loop:** Dispatch subsequent waves until all subtasks are `DONE` or `FAILED`.

---

## 8. Phase 3: Consolidation, Explicit Staging & Atomic Push

1. **Consolidate Subtasks:** Merge completed subtasks into `.ai-memory/plans/completed/NN-<slug>.md`. Document real steps used from the ledger and link to the canonical spec. Delete `.ai-memory/plans/subtasks/NN-<slug>/` and the pending plan file. Keep the canonical spec in `02-spec/21-app/` permanently intact.
2. **Update Index Registers:** Update `.ai-memory/plans/readme.md`, `01-prompts/readme.md`, and `.ai-memory/prompts.md`.
3. **Explicit Staging (R8):** Run `git add -- <every path in the ledger stage list>`. Verify with `git diff --cached --name-only` that ONLY intended files are staged.
4. **Atomic Commit & Push (R9):**
   - If tree was completely clean at preflight, use GitMap semantic commit: `gitmap cpf "<summary>"` (features) or `gitmap cpb "<summary>"` (fixes).
   - If preflight was dirty, use explicit commit: `git commit -m "<summary>"` and push according to `COMMIT`.

---

## 9. Final Report Format (Strict Vertical Lines)

```markdown
### Task Completion Summary

- ✅ **Task-01: [Descriptive Task Title]** — `[Completed]` — [diff/check evidence]
- ❌ **Task-02: [Descriptive Task Title]** — `[Failed]` — [RCA link]
- ⏳ **Task-03: [Descriptive Task Title]** — `[Deferred]` — [blocker]

### Modified Files Summary

- [relative/path/to/modified/file1.ext]
- [relative/path/to/modified/file2.ext]

### Steps Used

- Phase 1: [x] of [PHASE_1_BUDGET], Phase 2: [y] of [PHASE_2_BUDGET] (Total: [x+y] / 300)

### Implementation Confidence Score

- Confidence: [100%]
- Rationale: [Verified evidence across all criteria, passing targeted linters, zero regressions]

### 🤖 Independent AI Verification & Audit Prompt

(Emit the self-contained audit prompt linking to the spec, plan, and modified files)
```

---

## 10. Targeted Verification Checks (R2)

Confirm scripts exist before invoking (R4). Run on changed files/folders only:

- **Coding Guidelines & Boolean Linter:** `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only --ext <.ext>`
- **Relative Path Linter:** `python linter-scripts/check-relative-paths.py`
- **Prompts & Spec Index Linter:** `python linter-scripts/check-prompts-loaded.py`
- **Markdown Link & Doc Path Linter:** `python 03-ai-scripts/22-doc-path-linter.py <folder>`
- **Sequence Integrity Linter:** `python 03-ai-scripts/21-sequence-integrity-linter.py`
- **Forbidden Strings Check:** `python linter-scripts/check-forbidden-strings.py`

---

## 11. Discovery Toolchain (GitMap Primary)

- **GitMap:** `gitmap f "<pattern>"`, `gitmap lf [pat]`, `gitmap ffa "<str>"`, `gitmap cat <file>`, `gitmap search "<term>"`.
- **Python Fallback:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --stats`, `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<text>"`, `python 03-ai-scripts/17-fast-file-reader.py --read-file <file>`.

---

## 12. Issue Destination & RCA Routing

- **CI/CD, Workflow, Runner Failures:** `.ai-memory/cicd-issues/NN-<slug>.md`, indexed in `.ai-memory/cicd-index.md`.
- **Application Bugs (Logic, UI, CLI, API):** `02-spec/22-app-issues/NN-<slug>.md` with 4-part RCA (Reproduction, Cause, Fix, Prevention), indexed in `02-spec/22-app-issues/readme.md`.
- **Failed Subtasks (R13):** Log RCA in `.ai-memory/memory/issues/` and link from `ledger.md`.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run the targeted checks (builds and full unit tests stay in CI per R1), group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.
