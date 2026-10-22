---
name: next-100-ideas-for-improvements
description: Autonomously survey, discover, architect, and catalog the next 100 high-leverage innovation and improvement ideas for any target codebase across 10 ideation clusters using an 8-field specification schema, review-first gating, and pending plan enqueueing.
---

# Next 100 Ideas for Improvements

> **[/goal](slashCommand;goal)** Autonomously survey the target codebase, discover 100 high-leverage innovation and improvement ideas across 10 functional clusters, articulate each idea using the 8-field Idea Specification Schema, and enqueue them as pending plans in `.ai-memory/plans/pending/XX-next-100-ideas-for-improvements.md` with modular subtasks in `.ai-memory/plans/subtasks/XX-next-100-ideas-for-improvements/`.
> **[/plan](slashCommand;plan)** Plan first before doing work for improved work or less friction, faster work and more pure and effective work.
> **[/learn](slashCommand;learn)** Enforce the Review-First Gate (zero immediate code edits, NO-BUILD, NO-TEST during ideation discovery), GitMap search primacy, strict relative Git paths, and lowercase naming conventions.

**Version:** 1.0.0  
**Updated:** 2026-10-22  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Canonical Prompt:** `01-prompts/30-improvements/05-next-100-ideas-for-improvements.md`  
**Application Specification:** `02-spec/21-app/04-next-100-ideas-for-improvements/`  
**Execution Engine:** V6 Tree Swarm Hierarchy ($N = 400$, $A = 3$, $H = 2$, Capacity = 6)  
**Related Skills:**
- [`next-detail-out-100-improvements`](file;.agents/skills/next-detail-out-100-improvements)
- [`next-100-code-improvements`](file;.agents/skills/next-100-code-improvements)
- [`next-100-whole-system-improvements`](file;.agents/skills/next-100-whole-system-improvements)
- [`next-100-ui-improvements`](file;.agents/skills/next-100-ui-improvements)
- [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6)
- [`spec-authoring-and-validation`](file;.agents/skills/spec-authoring-and-validation)
- [`gitmap`](file;.agents/skills/gitmap)
- [`ai-verification`](file;.agents/skills/ai-verification)

---

## 1. When to Use This Skill

Activate this skill when:
- **Autonomous Ideation Flywheel:** Transforming software projects from reactive bug fixing and incremental maintenance into an autonomous innovation and superpower discovery flywheel.
- **Trigger Keywords & Aliases:** Invoked via `next-100-ideas-for-improvements`, `next-100-ideas`, `ideas-for-improvements`, `next-ideas`, `100-ideas`, or queries like `"generate next 100 great ideas for improves for the project"`.
- **Slash Command Invocation:** Triggered via `/next-100-ideas-for-improvements <project-or-module>`.
- **Pre-Refactoring Innovation Exploration:** Surveying unbuilt capabilities, developer velocity multipliers, architectural leaps, and ecosystem bridges before major engineering initiatives.
- **Review Docket Seeding:** Populating the repository's long-term backlog with rigorously specified idea cards ready for ingestion into [`next-detail-out-100-improvements`](file;.agents/skills/next-detail-out-100-improvements).

---

## 2. Core Architecture & Philosophy

The Next 100 Ideas for Improvements engine expands the canonical 100 Improvements Family beyond diagnostic audits into a **generative innovation flywheel**. It operates under four non-negotiable architectural pillars:

### 2.1 The 5-Layer Diagnostic Inspection Engine
Rather than hallucinating hypothetical ideas, the engine systematically surveys the real codebase across 5 distinct diagnostic layers using GitMap high-speed commands:
1. **Layer 1: Identity & Capabilities Discovery:** Inspects `readme.md`, `02-spec/`, CLI manifests, package dependencies, and domain boundaries to identify core mission, primary customer personas, and active feature boundaries.
2. **Layer 2: User Journeys, Entry Points & Interaction Friction:** Traces CLI command suites, HTTP endpoints, UI page routes, and terminal outputs to locate operational bottlenecks and unhandled edge cases.
3. **Layer 3: Developer Ergonomics & Operational Velocity:** Surveys runner scripts (`run.ps1`, `run.sh`), linters, test harnesses, and CI/CD pipelines to identify manual friction, slow iterations, and debugging opacity.
4. **Layer 4: Architectural Bottlenecks, Missing Abstractions & State Flows:** Evaluates data persistence (SQLite Split-DB), memory layout, caching policies, and concurrency primitives to spot unindexed queries, blocking I/O, and tight coupling.
5. **Layer 5: Ecosystem Expansion, Automation & Multi-AI Collaboration:** Surveys Model Context Protocol (MCP) tool surfaces, cross-repository hooks, export/import adapters, and subagent orchestration telemetry.

### 2.2 Strict Review-First Gating (Zero Immediate Execution)
> [!IMPORTANT]
> **Strict Review-First & Zero Immediate Execution Mandate:**
> Ideas discovered and formulated under this workflow MUST NEVER be implemented, edited into source files, built, or tested during the ideation turn.
> All 100 ideas are strictly enqueued into `.ai-memory/plans/pending/XX-next-100-ideas-for-improvements.md` with modular subtasks in `.ai-memory/plans/subtasks/XX-next-100-ideas-for-improvements/` and presented to the human architect for review.
> Implementation occurs only after human sign-off, proceeding through [`next-detail-out-100-improvements`](file;.agents/skills/next-detail-out-100-improvements) for deep specification and then canonical execution swarms via [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6).

### 2.3 The Synergistic 5-Pillar Improvements Flywheel
The improvements family forms a closed-loop engineering flywheel:
- `01-next-100-ui-improvements.md`: UI/UX, Design Tokens & Materiality polish.
- `02-next-100-code-improvements.md`: Code quality, guidelines compliance, monadic `Result[T]`, `*appfault.AppError`.
- `03-next-100-whole-system-improvements.md`: System friction, developer flow, and operational cause-and-effect impact audits.
- `05-next-100-ideas-for-improvements.md`: Generative innovation discovering "What Could Be" (this skill).
- `04-next-detail-out-100-improvements.md`: Deep specification of "How It Must Be Built" producing Go/No-Go Review Dockets.
- [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6): Post-approval autonomous execution.

### 2.4 Multi-Agent Concurrency & Budget Parameters
- **Step Budget ($N$):** `400` steps (100 Analyze, 150 Plan, 150 Execute/Enqueue).
- **Active Subagents ($A$):** `3` concurrent subagents.
- **Hierarchy Depth ($H$):** `2` tiers (Lead Orchestrator $\rightarrow$ Leaf Worker Subagents).
- **Pool Capacity:** `6` concurrent operations.
- **Build / Test Runs:** `DISABLED` (strict NO-BUILD and NO-TEST execution).

---

## 3. The 10 Ideation Taxonomy Clusters (100 Ideas Matrix)

The 100 improvement ideas are partitioned across 10 balanced taxonomy clusters (exactly 10 ideas per cluster, `IDEA-001` through `IDEA-100`):

| Cluster | Focus Domain | Scope & High-Value Vectors |
| :--- | :--- | :--- |
| **Cluster 01** | Core Product Capabilities & Unbuilt Features (`IDEA-001` to `IDEA-010`) | Predictive batch processing, smart query filters, visual diffing, bidirectional worktree sync, and automated migration wizards. |
| **Cluster 02** | Developer Ergonomics, Tooling & Velocity (`IDEA-011` to `IDEA-020`) | CLI macro recording/playback, terminal HUDs, smart commit synthesis, fuzzy command routing, and instant shell completions. |
| **Cluster 03** | Architecture, Modularization & Boundaries (`IDEA-021` to `IDEA-030`) | Monolithic engine modularization, value semantics, memory layout optimization, monadic Result envelopes, parameter structs, and map dispatch tables. |
| **Cluster 04** | Autonomous AI Agent Optimization & Context (`IDEA-031` to `IDEA-040`) | Lockless SQLite task ring buffers (WAL mode), hierarchical context compression, sliding memory windows, subagent crash autopsies, and file-lease coordinators. |
| **Cluster 05** | API Design, Schema Contracts & Type Safety (`IDEA-041` to `IDEA-050`) | End-to-end TS/Go DTO synchronization, discriminated union RPCs, declarative OpenAPI schemas, branded primitives, and schema drift detectors. |
| **Cluster 06** | Performance, Caching & Resource Utilization (`IDEA-051` to `IDEA-060`) | Zero-allocation string normalization (`strutil.EqualFoldAny`), multi-tier SQLite caches, SIMD search acceleration, and connection multiplexing. |
| **Cluster 07** | Reliability, Fault Tolerance & Error Recovery (`IDEA-061` to `IDEA-070`) | Adaptive circuit breakers, exponential jitter backoff, transactional staging buffers, deadlock prevention, and self-healing supervisor routines. |
| **Cluster 08** | Security, Secret Isolation & Governance (`IDEA-071` to `IDEA-080`) | Pre-commit entropy-based secret scanning, ephemeral memory-only decryption, sandboxed script runners, Casbin RBAC, and SBOM generation. |
| **Cluster 09** | Test Automation, Diagnostics & CI/CD (`IDEA-081` to `IDEA-090`) | Call-graph incremental test runners, zero-storage GitHub Actions summaries, AST test fixture generators, and flaky test quarantine engines. |
| **Cluster 10** | Documentation, Onboarding & Ecosystem (`IDEA-091` to `IDEA-100`) | Interactive spec sandboxes, AST diagram generators, contributor onboarding wizards, bi-directional doc-to-code validators, and ADR auto-synthesizers. |

---

## 4. The 8-Field Idea Specification Schema

Every single one of the 100 generated ideas must strictly follow the standardized 8-field schema without placeholders (`[N]`, `TODO`, `TBD`, `...`):

```markdown
#### Card [IDEA-NNN]: [Descriptive Action Title]

1. **Idea Identifier & Domain Slug:**
   - Identifier: `IDEA-NNN` (e.g. `IDEA-001` through `IDEA-100`)
   - Domain: `[Cluster Title]`
   - Target Subsystem: `[Package, Module, or Directory Path]`
2. **Problem / Unlocked Opportunity Statement:**
   [Current structural limitation, architectural bottleneck, developer friction point, or unexploited product capability.]
3. **Concrete Innovative Solution Description:**
   [Actionable feature, algorithmic leap, architectural pattern, or tooling primitive designed to address the problem or capture the opportunity.]
4. **User & Developer Value Multiplier:**
   [Quantitative metric (e.g. latency reduction, lines saved, memory cut) or qualitative UX/DX enhancement.]
5. **Technical Feasibility & Implementation Complexity:**
   - Complexity Tier: `LOW` (<150 lines, 1-2 files) | `MEDIUM` (150-500 lines, 3-5 files) | `HIGH` (>500 lines, architectural refactor)
   - Estimated Implementation Effort: [Estimated files, lines, and micro-batches]
6. **Blast Radius, Downstream Callers & Concurrency Boundaries:**
   - Downstream Consumers / Callers: [Packages, entrypoints, and callers impacted]
   - Concurrency & Threading Boundaries: [Goroutines, locks, file leases, or state immutability]
7. **Zero Harm / Regression Safety & Backward Compatibility Guarantee:**
   [Concrete technical proof that existing public APIs, CLI flags, database schemas, and workflows remain completely unaffected.]
8. **Actionable Enqueueing Blueprint & Acceptance Criteria:**
   - Enqueue Target: `.ai-memory/plans/subtasks/XX-next-100-ideas-for-improvements/YY-subtask-[cluster].md`
   - Acceptance Criteria: [Gherkin scenario or deterministic CLI command verification contract]
```

---

## 5. Workflow & GitMap Execution

The execution follows a strict 4-phase sequence complying with canonical V6 swarming and tool primacy:

### Step 0: Preflight & Working Tree Inspection
- Confirm working tree is clean: `git status --porcelain`.
- Check cognitive heatmap telemetry: `gitmap heatmap status` and `gitmap heatmap tables`.
- Initialize task engine: `gitmap agent task init --task "next-100-ideas-for-improvements" --budget 400`.

### Phase 1: 5-Layer Codebase Survey ($A = 3$ `research` Subagents)
- Dispatch parallel discovery subagents using GitMap high-speed streaming commands:
  - Inventory files: `gitmap lf <dir>` or `gitmap find "*.*" <dir>`
  - Live symbol and pattern search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
  - Substring file finder: `gitmap find "<pattern>" [dir]`
  - Zero-disk file streaming: `gitmap cat <filepath>`
- **TOTAL BAN on Shell Search Tools:** NEVER execute `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`, `Get-ChildItem`, `dir`, or `ls`.

### Phase 2: Idea Card Authoring & Cluster Assembly ($A = 3$ `self` Subagents)
- Partition authoring across subagents with disjoint cluster boundaries:
  - Worker 01: Clusters 01, 02, 05, and 08 (Core Features, Automation, Data, Performance).
  - Worker 02: Clusters 03, 04, 06, and 07 (DX, UX, Extensibility, Observability).
  - Worker 03: Clusters 09 and 10 (Security, Zero Harm, AI Orchestration, Detailing Alignment).
- Author exactly 10 fully articulated idea cards per cluster adhering strictly to the 8-field schema.

### Phase 3: Plan Enqueueing into `.ai-memory/plans/pending/`
- Author the Master Plan: `.ai-memory/plans/pending/XX-next-100-ideas-for-improvements.md` (where `XX` is the zero-padded sequential plan index).
- Author Modular Subtasks: `.ai-memory/plans/subtasks/XX-next-100-ideas-for-improvements/` (bounded files covering clusters).
- Update central indices:
  - Register under `## Pending Plans` in `.ai-memory/plans/readme.md`.
  - Enqueue in `.ai-memory/what-to-read.md`.

### Phase 4: Verification Audit & Atomic Push Gate
- Validate relative path hygiene: zero absolute paths or `file:///` URIs.
- Verify zero secrets: verify absence of tokens or private credentials.
- Final commit and push executed by the Lead Orchestrator:
  ```bash
  gitmap cpf "Improvements - next 100 ideas for improvements catalog and pending plans"
  ```
- Subagents never run git commands; all git mutations are reserved strictly for the Lead Orchestrator.

---

## 6. Binary Verification Criteria

An ideation run is complete and verified only when all criteria pass:

- [ ] **100 Ideas Cataloged:** Exactly 100 ideas are authored across all 10 clusters (10 per cluster, `IDEA-001` through `IDEA-100`).
- [ ] **8-Field Schema Compliance:** Every idea card provides all 8 fields with concrete technical mechanisms and zero placeholders (`[N]`, `TODO`, `TBD`).
- [ ] **Review-First Gate Respected:** Exactly 0 source code files modified, 0 build commands run, and 0 test suites executed during the discovery turn.
- [ ] **Master Plan Enqueued:** Master plan authored at `.ai-memory/plans/pending/XX-next-100-ideas-for-improvements.md`.
- [ ] **Subtask Files Created:** Modular subtasks authored in `.ai-memory/plans/subtasks/XX-next-100-ideas-for-improvements/`.
- [ ] **Registry Synchronization:** Indexed in `.ai-memory/plans/readme.md` and `.ai-memory/what-to-read.md`.
- [ ] **Strict Relative Git Paths:** Zero absolute filesystem paths or `file:///` URIs anywhere in the plan files, specs, or skills.
- [ ] **Strict Lowercase File Naming:** All paths and files strictly lowercase.
- [ ] **Upstream Push Completed:** Changes committed and pushed upstream ('No Push = Not Done').

---

## 7. Related Skills

- [`next-detail-out-100-improvements`](file;.agents/skills/next-detail-out-100-improvements)
- [`next-100-code-improvements`](file;.agents/skills/next-100-code-improvements)
- [`next-100-whole-system-improvements`](file;.agents/skills/next-100-whole-system-improvements)
- [`next-100-ui-improvements`](file;.agents/skills/next-100-ui-improvements)
- [`execute-parent-task-with-n-steps-v6`](file;.agents/skills/execute-parent-task-with-n-steps-v6)
- [`spec-authoring-and-validation`](file;.agents/skills/spec-authoring-and-validation)
- [`gitmap`](file;.agents/skills/gitmap)
- [`ai-verification`](file;.agents/skills/ai-verification)
