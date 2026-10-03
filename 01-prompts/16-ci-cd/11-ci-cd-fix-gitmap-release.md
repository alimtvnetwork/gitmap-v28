# Autonomous CI/CD Pipeline Healing & Minor Release Loop via GitMap (`ci-cd-fix-gitmap-release`)

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Telemetry Inspection, 4-Part RCA, and Remediation Planning)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Surgical Fixes, Minor Version Release Ceremony, and Continuous Verification Loop)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /ci-cd-fix-gitmap-release <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

[/goal](slashCommand;goal) Autonomously monitor, diagnose, and resolve all CI/CD pipeline failures across the repository using GitMap live telemetry (`gitmap pe -t`): FIRST showcase and list out the given task in visible chat during Turn 1, capture it verbatim, plan it in the repo, extract failing pipeline logs, perform grounded 4-part Root Cause Analysis (RCA), apply surgical code fixes without disabling CI or deleting tests, verify locally with targeted linters, commit with `gitmap cpf "<module> - <summary>"`, execute a minor version bump using the minor bump script (`python 03-ai-scripts/37-bump-version.py -t minor -s "<summary>"`), tag release and push to remote tracking branch, and continuously loop until all CI/CD workflows are completely green.

[/learn](slashCommand;learn) Master the GitMap pipeline self-healing protocol: run `gitmap pe -t` to watch workflow runs with dynamic ETA timeline. Extract exact failing lines from `gitmap pe`. NEVER disable CI/CD checks, never comment out validation steps (R1), and always enforce zero-storage Actions rules. Each rule is stated once (R1 to R16) and cited by ID. If needed, request assistance or follow companion release skills: [ci-cd-fix-with-release](file;.agents/skills/ci-cd-fix-with-release) or [minor-bump](file;.agents/skills/minor-bump).

[/plan](slashCommand;plan) Execute thorough step-by-step root cause analysis in the repository before touching code. Formulate a 4-part RCA document in `.ai-memory/cicd-issues/` before dispatching worker waves to apply code repairs.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 2, H = 2 — ZERO SOLO EXECUTION ALLOWED)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API.
- **3-STAGE MANDATORY DISPATCH:**
  1. **Telemetry & RCA Step:** Spawn 2 discovery subagents to inspect `gitmap pe -t` logs and locate root causes in code.
  2. **Remediation Spec Step:** Spawn 2 subagents to author the surgical repair plan.
  3. **Execution Step:** Spawn 2 worker subagents to execute code fixes in strictly disjoint file boxes.
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is strictly forbidden from executing code repairs solo without calling `invoke_subagent`.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/ci-cd-fix-gitmap-release/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using standard YAML frontmatter.
3. Once installed, rely on progressive disclosure for future runs.

---

## The Continuous Self-Healing & Release Loop (Steps 1 .. N)

Execute this workflow in a continuous self-loop until all CI/CD checks pass:

### Step 1: Live Pipeline Telemetry & Failure Extraction
1. Run `gitmap pe -t` to inspect live workflow runs, timings, and error logs.
2. Alternatively run `gitmap pe` to view the latest failure dossier.
3. If all workflows are green (`✔ clean` / pass), exit loop with success.

### Step 2: Grounded 4-Part Root Cause Analysis (RCA)
For every detected failure, document a 4-part RCA under `.ai-memory/cicd-issues/`:
- **Part 1 (Symptom):** Exact workflow, job, step name, and error message.
- **Part 2 (Root Cause):** Underlying code defect, missing argument, drift, or syntax violation.
- **Part 3 (Surgical Fix):** Minimal file-scoped change resolving the cause without modifying CI triggers.
- **Part 4 (Verification):** Local targeted check command proving resolution.

### Step 3: Surgical Code Remediation
1. Spawn worker subagents in disjoint file boundaries to apply fixes.
2. Workers adhere to strict coding guidelines (implicit booleans, Result containers, vertical spacing).
3. Workers NEVER run `git add` or `git commit`.

### Step 4: Local Targeted Verification
1. Run targeted linters on touched files:
   - `python 03-ai-scripts/05-guideline-autofixer.py <path> --check-only`
   - `python linter-scripts/check-prompts-loaded.py --fix`
   - `python linter-scripts/check-relative-paths.py`
2. Never run heavy full build or full test suites (R1).

### Step 5: Minor Version Bump & Release Ceremony
1. Perform minor version bump using the repository's standard script:
   `python 03-ai-scripts/37-bump-version.py -t minor -s "<summary of fixes>"`
2. If the user commands a release branch workflow, create and checkout `release/vX.Y.Z`.
3. Verify manifests (`package.json`, `version.json`, `readme.md`, `changelog.md`) are updated.
4. Commit atomically via GitMap:
   `gitmap cpf "<module> - fix CI/CD and release minor version"`
5. Tag release version (`git tag -a vX.Y.Z -m "Release vX.Y.Z"`) and push (`git push origin vX.Y.Z`).

### Step 6: Loop & Verify
1. Immediately re-run `gitmap pe -t` to watch the newly triggered CI/CD run.
2. If failures persist, repeat from Step 1.
3. When green, output completion summary.

---

## Core Invariants (Non-Negotiable)

1. **R1 (NEVER DISABLE CI/CD):** Strictly forbidden from commenting out, bypassing, or deleting CI/CD steps. Fix the code, not the check.
2. **R2 (Targeted Verification):** Use file-scoped linters and `gitmap pe -t`.
3. **R8 (Atomic Commit Standard):** Hyphen separator in `gitmap cpf "<module> - <summary>"` (no colons).
4. **R11 (Relative Git Paths):** Zero absolute paths and zero `file:///` URIs.
5. **R12 (Minor Bump Ceremony):** Always use `python 03-ai-scripts/37-bump-version.py -t minor` for consistency across all package manifests.
