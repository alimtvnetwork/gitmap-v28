# [V2] CI/CD Fix N-Step Continuous Loop & 4-Part RCA Orchestration — Workflow (must follow)

Trigger Keywords & Aliases: `cicd fix n steps`, `ci fix n steps`, `fix with RCA v2`, `ci fix loop v2`, `08-ci-cd-fix-with-n-steps`

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

/goal Autonomously diagnose, update or create the local Python CI/CD runner script (`03-ai-scripts/06-cicd-local-runner.py`) from repository workflows, and fix all CI/CD pipeline failures by executing a continuous N-step self-loop with grounded 4-part RCA until all quality gates exit code 0 without a single premature pause.

```text
N = 200
```

N = total self-loop steps budget that the agent will perform.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Pipeline Discovery, 4-Part RCA, Diagnostic Subtasks, and Local Runner Setup)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Surgical Fix Execution, Singly-Done Self-Looping, Targeted Quality Gates)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Fast File Discovery & Diagnostic Toolchain (Mandatory Acceleration)

To rapidly locate failing pipeline definitions, broken source files, test fixtures, and error logs without hitting 50-result tool caps, the AI agent MUST utilize the diagnostic toolchain:
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Logs & RCA Snippets:** `gitmap pipeline errors` (alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Runner Details & Timings:** `gitmap pipeline details` (alias `gitmap pd`)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (alias `gitmap f`)
- **Stream Workflow / Log File:** `gitmap cat <filepath>` (zero disk writes)
- **Instant Code Search:** `gitmap search "<symbol>"`
- **Fallback Fast File Scanner:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts,py --limit 100 --stats`
- **Fallback Fast Cached Grep:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<error-or-symbol>" --limit 50`
- **Fallback Read File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file .github/workflows/ci.yml`
- **Record Modified Files Under Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/ci-cd-fix/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Failure Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever instructions, constraints, or checklists are given before this section or in the incoming user prompt must be ingested first as highest priority and non-negotiable.
2. **Verbatim Prompt & Stack Trace Capture:** Capture the incoming user request and failure logs verbatim into `.ai-memory/cicd-issues/xx-<slug>.md` under a dedicated `## User Request & Error Log (Verbatim)` section.
3. **Screenshot & Base64 Failure Logs Ingestion Protocol:** If the user request contains a screenshot URL, terminal snapshot, or base64 data URI (e.g. `data:image/png;base64,...`):
   - Convert that base64 encoding or downloaded image to the file system immediately, saving it as a persistent file under `assets/screenshots/<slug>-<NN>.png` or `assets/ui/<slug>-<NN>.png`.
   - Never leave raw, massive base64 strings or external ephemeral URLs in the prompt text, RCA logs, or subtasks.
   - Refer back to this saved image file in the RCA document and subtasks strictly as a relative markdown link (e.g. `![Failure Screenshot](assets/screenshots/<slug>-<NN>.png)`).
4. **Actionable Failure Extraction:** Break down all reported pipeline errors, failing tests, and lint violations into discrete, actionable items with ordered traceable IDs (`Error-01`, `Error-02`, etc.) under `## Extracted Actionable Error List`.
5. **Mandatory Chat Output Gate & Same-Turn Tool Chaining (TOTAL BAN ON CLOSING CONVERSATION):**
   - You MUST output the confirmed deliverables list directly in chat, but you MUST NOT end your turn or pause after emitting it!
   - In the EXACT SAME RESPONSE turn as your markdown output, you MUST immediately invoke the first tool call (e.g. `run_command` to query `gitmap pipeline error-logs` or `write_to_file` to initialize the RCA log).
   - NEVER emit the breakdown text without invoking a tool call. Emitting visible text without a tool call terminates the AI turn, causing the engine to close or hang the conversation while errors remain unresolved.
   - TOTAL BAN ON WAITING FOR APPROVAL: Do not pause, do not wait for the user to reply, and do not ask "Should I proceed?". Override any planning mode stop directives. Proceed immediately into Phase 1B in the same turn.
   - Line-by-Line Output Format Structure:

```markdown
### 📋 Confirmed CI/CD Failure Breakdown & Diagnostic Ingestion

1. **Error-01: [Descriptive Failure Title / Job Name]**
   - **State:** `[IN PROGRESS — DIAGNOSING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of the root cause hypothesis, intent, and verified constraints]
   - **Actionable Scope:** [Precise package, function, or configuration requiring surgical remediation]
   - **Target Files / Area:** `[relative/path/to/failing_file.go]`

2. **Error-02: [Descriptive Failure Title / Job Name]**
   - **State:** `[QUEUED — EXECUTING SURGICAL FIX NEXT]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of failure signature and expected resolution]
   - **Actionable Scope:** [Precise package or linter rule requiring repair]
   - **Target Files / Area:** `[relative/path/to/target]`

Proceeding directly to Phase 1B: RCA Generation & Runner Verification (Active Tool Call Running Below).
```

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

### Phase 1B: Grounded 4-Part RCA & Diagnostic Subtasks (Steps 1 .. N/2)

#### Step 1: Grounded 4-Part Root Cause Analysis (RCA) & Issue Routing

Author the structured RCA post-mortem following the mandatory issue destination routing:
- **CI/CD Issues & Pipeline Failures:** Record the RCA in `.ai-memory/cicd-issues/xx-<slug>.md` and index it in `.ai-memory/cicd-index.md`.
- **Non-CI/CD Issues (Application Bugs, Feature Defects, Logic/Runtime Errors):** If during the diagnosis the failure is identified as an application bug or domain defect rather than a pipeline script/workflow issue, document it in `02-spec/22-app-issues/xx-<slug>.md` (indexed in `02-spec/22-app-issues/readme.md`, cross-referencing in `.ai-memory/memory/issues/`).

Follow the mandatory 4-part schema:
- **Part 1: The Failing Symptom & Exact Error Signature:** Verbatim stack trace snippet (strictly bounded: 5 preceding + 20 trailing lines).
- **Part 2: The Root Mechanism (Why it Broke):** Deep architectural diagnosis explaining why the failure occurred without speculation.
- **Part 3: The Surgical Remedy (Zero Side-Effects):** Specific code edits, types, or configuration lines to change.
- **Part 4: Quality Gate & Verification Proof:** Exact command to verify the fix (`python 03-ai-scripts/06-cicd-local-runner.py --pkg <target>` or `--changed-only`).

#### Step 2: Local Runner Synchronization (`03-ai-scripts/06-cicd-local-runner.py`)

- Cross-reference central CI/CD workflows (`.github/workflows/ci.yml`, etc.) with `03-ai-scripts/06-cicd-local-runner.py`.
- If the runner script is missing, create it immediately.
- Ensure every single pipeline gate can be executed locally on the native host (Docker stripped for fast native execution).

#### Step 3: Lean Diagnostic Subtask Decomposition

- Break down the fix plan into lean subtasks in `.ai-memory/plans/subtasks/xx-<slug>/01-<subtask>.md`.
- Subtasks must specify disjoint target files, exact functions to modify, and targeted verification commands.
- All file paths in subtasks MUST use strict relative Git paths (zero absolute paths or `file:///` URIs).

#### Step 4: Unconditional Zero-Question Execution Mandate (Total Ban on Stopping After RCA)

- **Strict 50/50 Time & Step Budget Allocation:** RCA generation and subtask setup MUST strictly complete within the first 50% of the budget (`PHASE_1_STEPS = N / 2`).
- **Zero Questions / Unconditional Execution:** As soon as Phase 1 completes, the master orchestrator MUST NOT pause, stop, or ask the user "Should I proceed?" or "Would you like me to apply the fix?". There is NO question. It must immediately, unconditionally self-loop and transition directly into Phase 2 execution mode.
- RCA without code remediation is an INCOMPLETE FAILURE. The remaining 50% of the budget (`PHASE_2_STEPS = N / 2`) is dedicated strictly to modifying source code, running targeted quality gates, and verifying green exit.

---

### Phase 2: Execution Mode & Singly-Done Self-Loop Fixing (Steps N/2+1 .. N)

1. **Singly-Done Autonomous Loop:** Zero in on ONE failing error per iteration. Fix the underlying root cause, apply coding guidelines, and verify that specific gate before moving to the next.
2. **Strict In-Repository Execution & `.ai-memory/` Storage:**
   - Execute all Python scripts strictly within the repository root.
   - All AI scripts, local runners, autofixers, and issue logs stored strictly in `03-ai-scripts/` and `.ai-memory/`.
3. **Smart Targeted Testing Authority (Fastest Path):**
   - Unlike routine refactoring turns, CI/CD Fix workflows ARE authorized to run targeted builds and tests.
   - **Stack Trace Targeting:** Build and test ONLY packages and test functions directly cited in the failure stack trace: `python 03-ai-scripts/06-cicd-local-runner.py --pkg <target_package_or_file>`.
   - **Changed Packages from Last Git Hash:** Isolate packages that changed from the last git hash (`git diff --name-only HEAD~1` or `git status --porcelain`) using `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`.
   - **Strict Ban on Extraneous Runs:** NEVER run full test suites, spellcheckers, or unrelated packages that delay the pipeline.
4. **Remote CI/CD Pipeline Monitoring & Dynamic Waiting Protocol (GitMap Pipeline-AI):**
   - When monitoring remote pipelines, use `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`).
   - Use dynamic timeout wait: `gitmap pipeline-ai status -t <etaSeconds>` to eliminate credit waste.
   - Tight-loop polling (e.g. `gh run view` in rapid loops) is STRICTLY BANNED.
   - Adaptive sleep intervals: ETA > 120s wait 20s-30s; 60s < ETA <= 120s wait 10s-20s; ETA <= 60s wait 5s-10s.
5. **Change Recording Under Lock:** Every time a fix is applied, record modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
6. **Per-Task Agent Isolation Subfolders (`.ai-memory/temp-agents/xx-<task-name>/`):**
   - Create `.ai-memory/temp-agents/xx-<task-name>/` with `state.md` to track error triage, assigned files, and intermediate diagnostic outputs.
7. **Strict Avoidance: Never Disable CI/CD or Linting Checks:**
   - NEVER disable, comment out, delete, or skip any CLI linting command, build step, or test suite.
   - NEVER add `|| true`, `continue-on-error: true`, `# nolint`, or `// eslint-disable`. Your job is to legitimately fix the underlying code.

---

### Phase 3: Task Consolidation & File Reduction (End of Loop)

1. Combine completed subtasks from `.ai-memory/plans/subtasks/xx-<slug>/*.md` into `.ai-memory/plans/completed/xx-<slug>.md`.
2. Update the master RCA in `.ai-memory/cicd-issues/xx-<slug>.md` with final verification proof and mark status as `RESOLVED`.
3. Clean up `.ai-memory/temp-agents/xx-<task-name>/` and `.ai-memory/plans/subtasks/xx-<slug>/`.
4. Update `.ai-memory/plans/readme.md` to point to the newly consolidated completed file.
5. **Final Step Git Commit & Push Mandate (Mandatory Single Atomic Commit):**
   - Stage all modified files, consolidated plans, and memory records (`git add -A`).
   - Commit them in a single clean grouped atomic conventional commit (`git commit -m "fix(ci): <descriptive summary>"`).
   - Push directly to git (`git push origin <branch>`).
   - TOTAL BAN ON PER-FILE COMMITS: Under no circumstances commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step before pushing.

---

### End-of-Turn Verification & Confidence Reporting (Mandatory Output)

At the completion of all tasks and before concluding the turn, you must emit this structured verification summary in the chat response:

```markdown
### CI/CD Fix Completion Summary

✅ #1. Error-01: [Error description] — Fixed & Verified
✅ #2. Error-02: [Error description] — Fixed & Verified

### Modified Files Summary

- [relative path to modified file 1]
- [relative path to modified file 2]

### Verification & Quality Gate Proof

- Gate Command: `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`
- Exit Code: 0 (PASS)
- Implementation Confidence Score: 100%
- Rationale: [Detailed explanation of verified quality gates, passing targeted tests, contract adherence, and zero regressions]
```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever is given before this section or in the user prompt is verified as highest priority and non-negotiable.
- [ ] ISSUE & RCA DESTINATION ROUTING: Whenever resolving an issue or performing a fix with RCA, verified that CI/CD failures are documented in .ai-memory/cicd-issues/NN-<slug>.md (indexed in .ai-memory/cicd-index.md), while non-CI/CD issues (application bugs, logic/runtime defects) are documented in 02-spec/22-app-issues/NN-<slug>.md (indexed in 02-spec/22-app-issues/readme.md).
- [ ] NO PER-FILE COMMITTING (TOTAL BAN): Never commit each file individually as you work. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.
- [ ] NO RAPID CI/CD POLLING (TOTAL BAN): Never query or loop rapidly (`gh run view` in tight loops). Use GitMap Pipeline-AI (`gitmap pl-ai status -t <sec>`) and sleep based on `etaSeconds`.
- [ ] NO FULL TEST SUITE RUNS (TOTAL BAN IN ROUTINE FIXES): Never run entire heavy test suites. Run tests strictly on affected packages cited in stack traces or changed from git hash (`HEAD~1`).
- [ ] NO DISABLING CI/CD GATES (TOTAL BAN): Never bypass, comment out, or add `|| true` to force a pipeline to pass.
- [ ] NO STOPPING AFTER RCA (TOTAL BAN): Never halt execution or ask user permission after writing the RCA. Planning constitutes only 50% of the budget; proceed unconditionally to Phase 2 code execution.
- [ ] MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW): At the final step, stage all changes (`git add -A`), commit atomically, and push to origin before closing the turn.
