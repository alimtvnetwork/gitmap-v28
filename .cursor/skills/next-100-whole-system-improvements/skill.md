---
name: next-100-whole-system-improvements
description: "Autonomously audit, discover, plan, and orchestrate 100 high-impact improvements across the whole system focusing strictly on user and developer experience, Zero Harm regression safety, and automated plan enqueuing into .ai-memory/plans/pending/."
---

# Next 100 Whole System Improvements (`next-100-whole-system-improvements`)

> **[/goal](slashCommand;goal)** Autonomously audit the entire system architecture, discover 100 high-leverage improvements focused strictly on user experience, developer ergonomics, agent friction reduction, performance, reliability, error clarity, and spec governance. Enqueue all 100 improvements as structured pending plans in `.ai-memory/plans/pending/` with subtasks in `.ai-memory/plans/subtasks/` using the Continuous Improvement Flywheel.
> **[/learn](slashCommand;learn)** Enforce the 8-Part Specification Schema, Zero Harm / Regression Safety Matrix, backward compatibility guarantees (e.g. `type Fault = AppError`), strict relative Git paths, lowercase naming, and GitMap command primacy.
> **[/plan](slashCommand;plan)** Plan first before doing work for improved work, less friction, faster execution, and pure, effective outcomes.

**Version:** 1.0.0  
**Updated:** 2026-10-22  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Canonical Prompt:** `01-prompts/30-improvements/03-next-100-whole-system-improvements.md`  
**Related Skills:**
- [`unified-ai-prompt-v4`](file;.cursor/skills/unified-ai-prompt-v4)
- [`execute-parent-task-with-n-steps-v6`](file;.cursor/skills/execute-parent-task-with-n-steps-v6)
- [`ai-memory-consolidate-reduce`](file;.cursor/skills/ai-memory-consolidate-reduce)
- [`spec-authoring-and-validation`](file;.cursor/skills/spec-authoring-and-validation)
- [`gitmap`](file;.cursor/skills/gitmap)

---

## 1. When to Use This Skill

Activate this skill when:
- Conducting comprehensive system-wide UX and DX architectural reviews.
- Seeking to identify concrete, non-breaking quality improvements across developer ergonomics, agent swarm orchestration, UI materiality, multi-repo fleet synchronization, performance, and spec governance.
- Populating the backlog with actionable, bite-sized tasks structured in the canonical 8-part schema.
- Running the Continuous Improvement Flywheel to continuously generate the next 100 system improvements on recurring development cycles.

---

## 2. Core Architecture & Philosophy

The Whole-System Improvements engine is anchored on four foundational architectural principles:

1. **User Experience (UX) Primacy:**
   - Every improvement must directly reduce friction, eliminate cognitive load, accelerate iteration speed, or improve error clarity for human developers or autonomous AI agents.
2. **Zero Harm / Regression Safety Matrix:**
   - Every improvement must be non-destructive, strictly confined within target boundaries, free of breaking changes, and executable without unrequested full builds or test suites.
3. **Backward Compatibility Guarantees:**
   - Preservation of existing CLI flags, API signatures, deprecation type shims (e.g. `type Fault = AppError`), and rollback safety.
4. **Continuous Improvement Flywheel:**
   - Discovered improvements are persisted directly into `.ai-memory/plans/pending/XX-<slug>.md` and `.ai-memory/plans/subtasks/XX-<slug>/`, and each execution dynamically synthesizes the next 100 improvement candidates.

---

## 3. The 10 System UX Clusters

The system evaluates enhancements across 10 distinct, non-overlapping operational clusters (10 items each = 100 total):

| Cluster | Focus Subsystem | Target Outcomes |
| :--- | :--- | :--- |
| **Cluster 01** | Developer Experience (DX) & CLI Ergonomics | Typo-tolerant fuzzy commands, terminal auto-scaling, one-step completions, live spinners, color formatting, environment health doctor. |
| **Cluster 02** | UI/UX, Design Tokens & Materiality | Unified design tokens, WCAG AAA contrast, fluid typography clamp scales, virtualized tree lists, glassmorphic depth, zero-FOUC hydration. |
| **Cluster 03** | Autonomous Agent Orchestration & Heatmap | SQLite WAL task queues, file lease contention locks, deadlock autopsy, live terminal tree matrix, dynamic hand capacity, pre/post-action telemetry. |
| **Cluster 04** | Prompt Engineering & Instruction Architecture | Canonical V6 parameterized templates, top-instruction priority verification, IDE skill link synthesis, lossless verbatim capture, negative guardrails. |
| **Cluster 05** | Multi-Repository Fleet Synchronization | Parallel multi-repo sync, 5 non-negotiable boundaries, pre-flight backup branches, zero secrets offloading (`gitmap rs`), script vault distribution (`gitmap rc`). |
| **Cluster 06** | Cross-Language Code Quality & Immutability | Monadic Result envelopes, domain `*appfault.AppError`, positive booleans (`is*`/`has*`), parameter structs, map dispatch tables, allocation-free strings. |
| **Cluster 07** | Performance, Caching & Resource Efficiency | Zero-disk memory streaming (`gitmap cat`), SQLite WAL tuning, SIMD AUM search, shared engine cache, lazy regex compilation, stack value semantics. |
| **Cluster 08** | Diagnostics, Observability & Error Management | 4-part RCA engine, pipeline error extractor (`gitmap pe -t`), contextual stack trace truncation, non-polling ETA waiting, zero-storage summaries. |
| **Cluster 09** | Testing, Verification & Regression Isolation | Destructive OS call mocking, temporary E2E test isolation, surgical micro-linters, in-memory fixtures, behavioral AAA naming, installer smoke tests. |
| **Cluster 10** | Documentation, Memory & Spec Governance | Strict relative path linters, lowercase filename governance, 3-tier application specs, lean subtask queues, memory consolidation, continuous flywheel. |

---

## 4. 8-Part Structured Specification Schema

Every improvement item MUST strictly follow the 8-part specification schema:

```markdown
### SYS-IMP-NNN: [Descriptive Title]

1. **Identifier & Subsystem:** `SYS-IMP-NNN` | Subsystem | Target Boundaries
2. **System-Level User Experience (UX) Impact:** Cognitive load reduction, workflow velocity, error clarity
3. **Cause (Current Deficiency / Root Friction):** Root cause analysis of current friction in 1-2 sentences
4. **Effect (Target State & Behavioral Outcome):** Precise target state, quantitative/qualitative metrics
5. **Zero Harm / Regression Safety Matrix:** Confinement rules, no-build/no-test rules, zero secrets invariants
6. **Backward Compatibility Guarantee:** CLI flags preservation, non-breaking signatures, deprecation shims
7. **Concrete Implementation Blueprint:** Concise code snippet, schema diff, or configuration enhancement
8. **Binary Verification & Acceptance Criteria:** Exact GitMap command or test check returning exit code 0
```

---

## 5. Workflow & GitMap High-Speed Execution

Execute this workflow adhering to canonical V6 multi-agent orchestration:

### Step 0: Preflight & Tool Handshake
- Verify GitMap availability: `gitmap --version`
- Check working tree cleanliness: `git status --porcelain` (never touch uncommitted files)
- Initialize or inspect cognitive heatmap: `gitmap heatmap status` and `gitmap heatmap tables`

### Phase 1: Planning & Specification
- Dispatch `A = 3` read-only discovery subagents (`TypeName: "research"`) via `invoke_subagent` using GitMap high-speed commands:
  - Inventory files: `gitmap lf <dir>` or `gitmap find "*.*" <dir>`
  - Streaming code search: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
  - Zero-disk streaming: `gitmap cat <filepath>`
- Dispatch `A = 3` authoring subagents (`TypeName: "self"`) to author modular specifications in `02-spec/21-app/` and subtask plans in `.ai-memory/plans/subtasks/`.

### Phase 2: Plan Enqueueing & Flywheel Activation
- Enqueue all 100 improvements as pending plans:
  - Parent plan: `.ai-memory/plans/pending/XX-next-100-whole-system-improvements.md`
  - Modular subtasks: `.ai-memory/plans/subtasks/XX-next-100-whole-system-improvements/01-task.md` through `NN-task.md`
- Register the newly enqueued tasks in `.ai-memory/plans/readme.md` and `.ai-memory/what-to-read.md`.
- Synthesize next 100 improvement candidates for subsequent flywheel cycles.

### Phase 3: Push Gate & Atomic GitMap Commit
- Run surgical micro-linters:
  - `python linter-scripts/check-relative-paths.py`
  - `python linter-scripts/check-prompts-loaded.py`
- Verify cognitive heatmap: `gitmap heatmap tree --task <task>`
- Stage and commit atomically using GitMap semantic commit command:
  ```bash
  gitmap cpf "Improvements - next 100 whole system improvements and enqueue plans"
  ```
- Format strictly with hyphen `-` (no colons in GitMap arguments).
- Ensure `git status --porcelain` is 100% clean and upstream remote reflects the commit ('No Push = Not Done').

---

## 6. Verification & Acceptance Criteria

- [ ] Source prompt `01-prompts/30-improvements/03-next-100-whole-system-improvements.md` exists and contains all 100 improvements in the 8-part schema.
- [ ] Category index `01-prompts/30-improvements/readme.md` lists all 4 improvement prompts.
- [ ] All file references strictly use relative paths from repository root (TOTAL BAN on absolute paths and `file:///` URIs).
- [ ] All filenames are strictly lowercase.
- [ ] Working tree is clean and changes pushed upstream to GitHub.
