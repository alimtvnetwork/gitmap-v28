---
name: detail-plans-and-improvements-for-review
description: >-
  Autonomously detail out all active, pending, and completed plans and improvements into exhaustive multi-step roadmaps, structural contracts, and review dockets with strict review-first gating (no immediate execution without user review).
---

# Detail Plans and Improvements for Review (`detail-plans-and-improvements-for-review`)

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Plan & Improvement Ingestion, GitMap Discovery, Review Docket Scaffold)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: 8-Field Detailed Plan Specification Card Authoring, Multi-Step Roadmap Synthesis)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Review Docket Compilation, Pending Task Backlog Enqueueing, Strict Review-First Verification Gate)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 1.0.0  
> Runtime: Google Antigravity 2.0 (IDE and CLI)  
> Invoke: `/detail-plans-and-improvements-for-review <target-plan-or-improvement-catalog>`  
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**  
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

> [/goal](slashCommand;goal) Autonomously detail out all active, pending, and completed plans, discovered improvements, and architectural backlog items into exhaustive multi-step roadmaps, structural contracts, dependency graphs, and comprehensive review dockets. Present and enqueue all detailed plans into `02-spec/21-app/XX-<slug>/` and `.ai-memory/plans/pending/XX-<slug>.md` with accompanying subtasks in `.ai-memory/plans/subtasks/XX-<slug>/` under the Strict Review-First Gate. NEVER execute or mutate production source code during this detailing phase; execution is strictly blocked pending explicit user review and sign-off.
>
> [/learn](slashCommand;learn) Enforce the Review-First Gate: improvements and detailed plans must be presented and enqueued for review; NEVER execute or implement code changes right away before explicit review. Ground all detailed plans in the 8-Field Detailed Plan Specification Schema, Zero Harm / Regression Safety Matrix, backward compatibility guarantees, strict relative Git paths, lowercase naming conventions, and GitMap command primacy.
>
> [/plan](slashCommand;plan) Plan first before doing work for improved work or less friction, faster work and more pure and effective work. Decompose each plan into clear, isolated, sequential subtasks with explicit pre-conditions, blast radius boundaries, and binary acceptance criteria before concluding.

**Version:** 1.0.0  
**Updated:** 2026-10-22  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Canonical Prompt:** `01-prompts/30-improvements/04-detail-plans-and-improvements-for-review.md`  
**Related Skills:**  
- [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6)  
- [`find-100-improvements-in-whole-system`](file;.agents/skills/find-100-improvements-in-whole-system)  
- [`find-100-improvements-in-code-base`](file;.agents/skills/find-100-improvements-in-code-base)  
- [`find-100-improvements-in-ui-level`](file;.agents/skills/find-100-improvements-in-ui-level)  
- [`ai-memory-consolidate-reduce`](file;.agents/skills/ai-memory-consolidate-reduce)  
- [`spec-authoring-and-validation`](file;.agents/skills/spec-authoring-and-validation)  
- [`gitmap`](file;.agents/skills/gitmap)  

---

## 1. When to Use This Skill

Activate this skill when:
- **Expanding Discovered Improvements:** Converting raw improvement findings (from `find-100-improvements-in-whole-system`, `find-100-improvements-in-code-base`, or `find-100-improvements-in-ui-level`) into rigorous, fully-specified execution blueprints.
- **Preparing Review Dockets for Human Sign-Off:** Generating comprehensive pre-implementation review dockets detailing architectural choices, risk assessments, dependency DAGs, and rollback strategies.
- **Decomposing Monolithic Plans into Micro-Tasks:** Breaking down oversized architectural proposals or pending backlog plans into bounded 5–8 file micro-tasks with discrete verification criteria.
- **Review-First Gating Enforcement:** Whenever a user or project policy mandates that plans, specifications, and improvements must be reviewed and approved by stakeholders before a single line of production code is written or modified.
- **Trigger Keywords & Aliases:** `detail plans for review`, `detail improvements`, `elaborate plan`, `expand pending plans`, `generate review docket`, `prepare improvements for user sign-off`, `detail out 100 improvements`, `roadmap pending plans`, `structural contracts for improvements`, `review-first improvements`, `architect review docket`.

---

## 2. Core Architecture & Invariants

The Plan Detailing & Review Docket engine operates under strict non-negotiable architectural invariants:

### 2.1 The Review-First Mandate (NO IMMEDIATE EXECUTION)
- **Zero Production Mutations:** During execution of this skill, the agent is **STRICTLY PROHIBITED** from modifying production codebase files, running builds, running tests, or initiating automated refactoring waves.
- **Complete Separation of Planning and Execution:** The output of this skill is purely analytical, architectural, and planning-centric. It authors specifications in `02-spec/21-app/XX-<slug>/`, enqueues pending plans in `.ai-memory/plans/pending/XX-<slug>.md`, creates modular subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`, and compiles the `00-review-docket.md`.
- **Review Gate Blocking:** Execution remains strictly blocked until the user explicitly reviews the presented docket and issues an execution directive (e.g. via [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6)).

### 2.2 Operational Swarm Invariants
- **Step Budget ($N = 400$):** Divided systematically across Phase 1 Ingestion & Discovery (Steps 1–100), Phase 2 Detailed Plan Card Authoring (Steps 101–250), and Phase 3 Review Docket Compilation & Enqueueing (Steps 251–400).
- **Mandatory Subagent Concurrency ($A = 3$, $H = 2$, Capacity = 6):** Solo authoring by the lead orchestrator is an auto-reject failure. The lead orchestrator dispatches $A = 3$ concurrent subagents via `invoke_subagent` in non-overlapping file and subsystem boundaries.
- **Dynamic Join & Zero Idle Workers:** When any subagent completes its analytical or authoring assignment, it immediately joins the next pending subtask to maintain 100% capacity saturation.

### 2.3 Strict Hygiene & Repository Governance
- **Strict Relative Git Paths Only:** All file paths, markdown links, citations, and subtask paths must strictly use relative paths from the repository root (e.g. `02-spec/...`, `.ai-memory/...`, `cmd/...`). Absolute filesystem paths (e.g. `C:\...`, `/home/...`) and `file:///` URIs are **TOTALLY BANNED**.
- **Strict Lowercase Naming Convention:** All files, directories, scripts, and documentation generated by this skill must strictly use lowercase filenames (e.g. `00-review-docket.md`, `01-task.md`, `readme.md`).
- **GitMap High-Speed Search Primacy:** All symbol lookups, file indexing, and codebase research must strictly use compiled GitMap commands (`gitmap aum search`, `gitmap find`, `gitmap cat`, `gitmap lf`). Shell traversal commands (`dir`, `ls`, `Get-ChildItem`) and standard grep tools (`rg`, `ripgrep`, `grep`, `Select-String`) are strictly banned.
- **Zero Secrets Storage:** Secrets, credentials, and API keys must never appear in plans or review dockets. If credentials are referenced, they must point strictly to `repo-secrets` via `gitmap rs`.

---

## 3. The 8-Field Detailed Plan Specification Schema

Every detailed plan and improvement item authored under this skill MUST strictly adhere to the standardized 8-Field Specification Schema. This ensures complete reproducibility, unambiguous review contracts, and seamless handoff to execution engines:

```markdown
### PLAN-DET-NNN: [Descriptive Plan / Improvement Title]

1. **Identifier & Scope Boundary:**
   - **ID:** `PLAN-DET-NNN` (or `IMP-DET-NNN` matching discovery catalog)
   - **Subsystem / Module:** Target package, directory, or architectural domain (e.g. `pkg/appfault`, `src/components/navigation`)
   - **Target Files:** Strictly relative Git paths of files to be modified or created
   - **Blast Radius & Caller Count:** Exact number of callers/consumers discovered via `gitmap aum search`

2. **System-Level User Experience (UX) & Architectural Value:**
   - Concrete impact on developer velocity, cognitive overhead, system responsiveness, or end-user interaction
   - Reduction in friction, eliminated boilerplate, or enhanced error clarity

3. **Current Deficiency & Friction Root Cause Analysis (RCA):**
   - Precise 4-part RCA diagnosis in 1-2 sharp sentences:
     - *Symptom:* Observable friction, defect, or inefficiency
     - *Trigger Condition:* Scenario where the deficiency manifests
     - *Root Flaw:* Underlying architectural, type, or convention violation
     - *Blast Radius:* Scope of affected consumers or subsystems

4. **Target Behavioral State & Measurable Acceptance Metrics:**
   - Unambiguous target state describing system behavior post-execution
   - Quantitative benchmarks (e.g. line count reduction, zero allocations, sub-millisecond execution, 0 lint warnings)

5. **Zero Harm / Regression Safety Matrix:**
   - Confinement boundaries: exact files and packages permitted for modification
   - Isolation guarantees: no destructive system calls, no unmocked network or OS operations
   - Build & test policy: no unauthorized builds or test runs during planning/detailing

6. **Backward Compatibility & Deprecation Shims:**
   - Preservation of existing CLI flags, public APIs, and configuration schemas
   - Forwarding shims and type aliases (e.g. `type Fault = AppError`) to prevent breaking downstream callers
   - Rollback strategy and contingency restoration procedure

7. **Concrete Implementation Blueprint & Step-by-Step Execution Roadmap:**
   - Complete, copy-paste ready code diff or blueprint (no placeholder comments or hand-waving)
   - Structured parameter DTOs, discrete handler mappings, or Result envelope signatures
   - Ordered micro-step execution breakdown (5–8 files maximum per micro-batch)

8. **Binary Verification Contract & Test Commands:**
   - Exact GitMap, micro-linter, or verification command returning exit code 0
   - Deterministic assertions and verification criteria proving 100% adherence to coding guidelines
```

---

## 4. Structure of the Review Docket (`00-review-docket.md`)

When detailing a batch of plans or improvements, the subagent fleet compiles a centralized **Review Docket** placed at `02-spec/21-app/XX-<slug>/00-review-docket.md`. The docket acts as the single source of truth for stakeholder review:

```markdown
# Review Docket: [Batch / Task Title]

> [!IMPORTANT]
> **Status:** PENDING USER REVIEW (Execution Blocked)
> **Authoring Model:** Autonomous Detailing Fleet (A = 3, H = 2)
> **Review Gate Invariant:** No code changes will be applied until explicit user approval is granted.

## 1. Executive Summary & Impact Scorecard
- High-level overview of proposed changes across the target domain
- Aggregate metrics: Total items detailed, files affected, estimated lines changed, cognitive friction eliminated

## 2. Risk & Impact Matrix
| Plan ID | Title | Blast Radius | Risk Tier (Low/Med/High) | Rollback Complexity | Backward Compatible? |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `PLAN-DET-001` | ... | 4 files | Low | Trivial (git revert) | Yes (`type Fault = AppError`) |

## 3. Multi-Stage Dependency DAG (Directed Acyclic Graph)
- Sequence of execution waves respecting dependencies (e.g. Data Layer -> Service Layer -> Presentation Layer)
- Explicit callout of independent tasks suitable for concurrent execution

## 4. Master Specification Cards
[Includes all detailed items formatted strictly in the 8-Field Schema above]

## 5. User Review Gate & Sign-Off Checklist
- [ ] Architectural boundary containment verified
- [ ] Backward compatibility guarantees confirmed
- [ ] Verification contracts validated
- [ ] User Sign-Off Decision: [PENDING / APPROVED / REVISE]
```

---

## 5. Workflow & Tool Execution

The detailing workflow strictly follows canonical multi-agent coordination across three distinct phases:

### Phase 1: Ingestion & GitMap Discovery (Steps 1–100)
1. **Preflight Handshake:**
   - Inspect cognitive heatmap: `gitmap heatmap status` and `gitmap heatmap tables`.
   - Log pre-action intent:
     ```bash
     gitmap heatmap log --task <task-slug> --action READ --file "01-prompts/30-improvements/" --state EXPLORATION --expected "Ingest improvements catalog"
     ```
2. **Catalog & Plan Ingestion:**
   - Ingest pending plans from `.ai-memory/plans/pending/` or discovered improvements from `01-prompts/30-improvements/`.
3. **High-Speed Codebase Discovery (A = 3 Discovery Subagents):**
   - Dispatch `A = 3` read-only subagents (`TypeName: "research"`) via `invoke_subagent`.
   - Map exact caller blast radius and type definitions using GitMap:
     - `gitmap lf <dir>`
     - `gitmap aum search "<symbol>" [dir] [-e <.ext>]`
     - `gitmap cat <filepath>`
4. **Scaffold Docket & Plan Trees:**
   - Initialize target directory: `02-spec/21-app/XX-<slug>/`.
   - Establish `.ai-memory/plans/subtasks/XX-<slug>/` directory structure.

### Phase 2: Detailed Plan Card Authoring (Steps 101–250)
1. **Dispatch Authoring Wave (A = 3 Authoring Subagents):**
   - Dispatch `A = 3` authoring subagents (`TypeName: "self"`) via `invoke_subagent` across partitioned clusters.
   - Author exhaustive 8-Field Detailed Plan Specification Cards into `02-spec/21-app/XX-<slug>/`.
2. **Subtask Plan Generation:**
   - Author modular, execution-ready subtask files in `.ai-memory/plans/subtasks/XX-<slug>/01-task.md` through `NN-task.md`.
   - Ensure each subtask is strictly bounded to 5–8 files.
3. **Register Tasks in GitMap Task Engine:**
   - Use compiled GitMap commands to initialize task structures:
     ```bash
     gitmap agent task init --id <slug> --title "<Task Title>"
     gitmap agent subtask add --task <slug> --id 01-subtask --title "<Subtask Title>"
     ```

### Phase 3: Review Docket Assembly & Gating (Steps 251–400)
1. **Compile Master Review Docket:**
   - Assemble `02-spec/21-app/XX-<slug>/00-review-docket.md` with executive scorecard, dependency DAG, and sign-off checklist.
2. **Enqueue Pending Parent Plan:**
   - Create or update `.ai-memory/plans/pending/XX-<slug>.md` referencing the review docket and all modular subtasks.
   - Update `.ai-memory/plans/readme.md` and `.ai-memory/what-to-read.md` index files.
3. **Enforce the Review-First Gate:**
   - Update heatmap telemetry: `gitmap heatmap tree --task <task-slug>` and verify 100% action closure.
   - **HALT BEFORE CODE EXECUTION:** Present the complete Review Docket, risk assessment, and plan summary to the user in visible chat.
   - Highlight key architectural decisions, proposed blueprints, and await explicit user confirmation before any implementation proceeds.

---

## 6. IDE Skill Link Syntax Standard

All skill references, command triggers, and tool documentation links MUST strictly follow IDE-detectable syntax:

- **Antigravity Agents:**
  - [`detail-plans-and-improvements-for-review`](file;.agents/skills/detail-plans-and-improvements-for-review)
  - [`find-100-improvements-in-whole-system`](file;.agents/skills/find-100-improvements-in-whole-system)
  - [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6)
- **Cursor IDE:**
  - [`detail-plans-and-improvements-for-review`](file;.cursor/skills/detail-plans-and-improvements-for-review)
  - [`find-100-improvements-in-whole-system`](file;.cursor/skills/find-100-improvements-in-whole-system)
  - [`execute-parent-task-with-n-steps-v6`](file;.cursor/skills/execute-parent-task-with-n-steps-v6)
- **Interactive Slash Commands:**
  - [`[/plan](slashCommand;plan)`](slashCommand;plan)
  - [`[/goal](slashCommand;goal)`](slashCommand;goal)
  - [`[/learn](slashCommand;learn)`](slashCommand;learn)
- **Prohibited Syntaxes:** Never append `/skill.md` or use raw filesystem paths (e.g. `d:/...` or `file:///...`) inside links, as they break IDE interactivity and violate path hygiene mandates.

---

## 7. Verification & Acceptance Criteria

Before concluding this skill workflow, verify all following acceptance criteria:

- [ ] YAML frontmatter is valid and specifies `name: detail-plans-and-improvements-for-review` with descriptive text.
- [ ] Source canonical prompt `01-prompts/30-improvements/04-detail-plans-and-improvements-for-review.md` is referenced.
- [ ] Review-First Gate is strictly enforced: zero code mutations occurred during detailing turns.
- [ ] All detailed plan items strictly implement the complete 8-Field Specification Schema.
- [ ] Centralized Review Docket is authored at `02-spec/21-app/XX-<slug>/00-review-docket.md` with risk matrix and dependency DAG.
- [ ] Modular subtasks are enqueued in `.ai-memory/plans/subtasks/XX-<slug>/` (bounded to 5–8 files each).
- [ ] Parent plan is registered in `.ai-memory/plans/pending/XX-<slug>.md`, `.ai-memory/plans/readme.md`, and `.ai-memory/what-to-read.md`.
- [ ] All file paths use strictly relative Git paths (TOTAL BAN on absolute paths and `file:///` URIs).
- [ ] All generated files and directories use strictly lowercase naming.
- [ ] Codebase discovery strictly used GitMap commands (`gitmap aum search`, `gitmap find`, `gitmap cat`, `gitmap lf`).
- [ ] Heatmap actions logged and verified with 100% closure.
