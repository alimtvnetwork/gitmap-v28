# Master Plan: GitMap Commit Tracking and Database Integration (Task 315)

## Master Plan Metadata
- **Plan File**: `.ai-memory/plans/completed/315-gitmap-commit-tracking-and-database-integration.md`
- **Task Identifier**: Task 315 (`gitmap-commit-tracking-and-database-integration`)
- **Status**: `COMPLETED`
- **Specification References**:
  - `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/00-overview.md`
  - `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/01-architecture-spec.md`
  - `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/02-component-spec.md`
  - `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/99-consistency-report.md`
- **Execution Swarm Concurrency**: Lead Orchestrator -> 3 Parallel Subagents ($A=3, H=2$, Concurrency Capacity = 6)
- **Total Step Budget**: $N=400$ (Phase 1: 100, Phase 2: 150, Phase 3: 150)
- **Target Version**: `v6.593.0` (SemVer Minor Bump)

---

## 1. Executive Summary & Strategic Intent

Task 315 delivers an end-to-end **Commit Provenance Tracking & Database Integration System** for GitMap (`v6.593.0`). The system resolves the critical opacity in distributed multi-agent workflows by capturing, indexing, and visually reporting the complete cognitive provenance of every git commit across four interconnected layers: the local Git repository, an external high-performance Split-DB (`repodb/commit_history.db`), the cognitive heatmap databases, and agent task manager stores.

Key capabilities delivered:
1. **Universal Commit Info Command (`gitmap commit info <sha>`, `ci`, `commit-info`)**: Instant multi-stage resolution of any 40-char SHA, 7-char short SHA, prefix, suffix, or arbitrary substring, rendering rich ASCII tree visualizers detailing author metadata, AI model, subagent ID, heatmap session, task slug, subtasks completed, spec link, and modified files.
2. **External Root Commit Split-DB (`<EffectiveDataDir>/repodb/commit_history.db`)**: High-throughput SQLite split-db outside the repository operating in WAL mode with 256MB memory mapping (`PRAGMA mmap_size = 268435456`) and single-writer mutex synchronization, guaranteeing $< 0.5\text{ ms}$ indexed lookups.
3. **Commit Provenance Flags & Auto-Ingestion**: Parameter parsing across all commit entry points (`commit`, `cm`, `cpf`, `cpb`, `cpr`, `cpc`, `pcp`) supporting `--agent-id`, `--heatmap-id`, `--ai`, `--spec`, and `--task`, ambient environment fallbacks, automatic flag stripping from messages, and automatic persistence on every commit.
4. **Origin Disambiguation**: Deterministic machine attribution reporting whether a commit originated from the local machine's AI or an external worker/node/human.
5. **Coding Guidelines, Skills & V6 Prompt Synchronization**: Authors guideline `12-commit-provenance-and-tracking.md`, updates GitMap skills, updates the V6 execution prompt and skill, and completes the minor release ceremony to `v6.593.0`.

---

## 2. The 5 Mandatory V6 Pillars

### Pillar 1: Strategic Intent & Systems Reasoning
In modern agentic swarms, code generation is delegated across multiple LLM engines and worker nodes. Standard Git logs capture only local author name and email, creating total attribution blindness when regressions occur. Task 315 equips GitMap with deep cognitive provenance tracking, connecting each commit SHA to its originating task, AI model, heatmap action, and specification. By storing this cache in an external Split-DB, GitMap ensures zero repository file pollution, zero git merge conflicts, and instant cross-repo commit discovery.

### Pillar 2: Explicit Scope Boundaries
- **In-Scope**: CLI command `gitmap commit info <sha>` (with aliases `ci`, `commit-info`); 4-stage universal SHA resolution; external Split-DB `repodb/commit_history.db`; commit flags (`--agent-id`, `--heatmap-id`, `--ai`, `--spec`, `--task`); ambient environment fallbacks; auto-ingestion commit hooks; origin disambiguation; ASCII tree rendering; coding guideline 12; GitMap skill updates; V6 prompt and skill updates; minor bump release.
- **Excluded**: Interactive git rebasing; remote API network calls during local SHA queries; in-repository SQLite databases; mandatory flags for human developers; blocking commit failures on DB logging errors.

### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
- **Search Ban**: Absolute ban on `rg`, `ripgrep`, `grep`, `git grep`, and `Select-String`. All discovery executed via GitMap toolchains (`gitmap aum search`, `gitmap find`, `gitmap cat`).
- **Worker Git Ban**: Autonomous subagents must never execute git commands. Staging and committing is performed exclusively by the Lead Orchestrator via `gitmap cpf`.
- **Zero Builds & Tests Mandate**: Zero full builds (`go build ./...`) or test sweeps (`go test ./...`) during routine implementation.
- **Relative Path Invariant**: 100% relative repository paths (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero drive letters (`C:`, `D:`) or file URIs.
- **Positive Booleans**: All booleans must begin with `is` or `has` (`isLocalAI`, `hasHeatmapMatch`).
- **Function Caps**: Functions $\le 15$ lines max ($\le 8$ lines preferred). Zero nested ifs.

### Pillar 4: Blast Radius & Dependency Topology
The primary blast radius is contained within `cli/cmdcommit/` and `cli/store/`. Hooking into `commitAndOptionalPush` and `performCommitPush` is wrapped in safe exception guards (`defer/recover` and non-fatal logging) so that database failures never compromise git commit execution. Upstream prompt changes in `01-prompts/14-execute/` and `.agents/skills/` are additive.

### Pillar 5: Multi-Agent Concurrency Topology
- Concurrency Parameters: $A = 3$ concurrent subagents, $H = 2$ hands per agent, Concurrency Capacity = 6.
- Phase 1: 3 Authoring Subagents author modular specs (00, 01, 02, 99) and 6 subtask plans.
- Phase 2: Parallel implementation across disjoint file sets (Commit Info, Split-DB, Flags/Ingestion).
- Phase 3: Guideline authoring, skill updates, V6 prompt sync, targeted linting, and atomic release.

---

## 3. Subtask Topology & Disjoint Ownership

```text
.ai-memory/plans/subtasks/315-gitmap-commit-tracking-and-database-integration/
├── 01-subtask-commit-info-command.md          (Phase 2: Wave 1 - CLI Routing & Entry Point) [COMPLETED]
├── 02-subtask-sha-resolver-and-tree-view.md   (Phase 2: Wave 1 - Resolution Pipeline & Tree UI) [COMPLETED]
├── 03-subtask-commit-history-split-db.md      (Phase 2: Wave 2 - External Split-DB & Schema) [COMPLETED]
├── 04-subtask-commit-flags-and-auto-ingestion.md (Phase 2: Wave 2 - Flags Parser & Auto-Ingestion) [COMPLETED]
├── 05-subtask-guidelines-and-skills-sync.md   (Phase 3: Wave 3 - Guideline 12 & GitMap Skill) [COMPLETED]
└── 06-subtask-v6-prompt-sync-and-release.md   (Phase 3: Wave 3 - V6 Prompt Sync & Release) [COMPLETED]
```

### Detailed Subtask Disjoint Bounding Boxes

| Subtask ID | Subtask Plan File | Scope & Deliverables | Exclusive File Ownership | Assigned Worker | AI Comprehension | Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **ST-01** | `.../01-subtask-commit-info-command.md` | CLI command entry point, alias registration (`ci`, `commit-info`), argument parsing, JSON flag, error sovereignty. | `cli/cmdcommit/commit_info.go`<br>`cli/cmd/rootcore.go` | Worker 3 | **[YES]** Understood: command parses SHA argument, routes to resolver, and delegates to tree renderer without exit calls. | `COMPLETED` |
| **ST-02** | `.../02-subtask-sha-resolver-and-tree-view.md` | 4-stage SHA resolver (rev-parse, Split-DB, Heatmap, rev-list stream), cross-db metadata query, origin attribution, ASCII tree rendering. | `cli/cmdcommit/commit_info_resolve.go`<br>`cli/cmdcommit/commit_info_query.go`<br>`cli/cmdcommit/commit_info_tree.go` | Worker 3 | **[YES]** Understood: resolves partial SHAs, checks local vs external machine origin, and prints structured tree. | `COMPLETED` |
| **ST-03** | `.../03-subtask-commit-history-split-db.md` | External Split-DB `<EffectiveDataDir>/repodb/commit_history.db`, WAL mode, 256MB mmap, models (`CommitHistoryRecord`, `CommitFileTouch`), registry sync. | `cli/store/split_db_path.go`<br>`cli/store/commit_history_*.go`<br>`cli/store/split_database_registry_sync.go` | Worker 1 | **[YES]** Understood: external database path outside repo, single-writer mutex, sub-millisecond WAL queries. | `COMPLETED` |
| **ST-04** | `.../04-subtask-commit-flags-and-auto-ingestion.md` | Commit flags parser (`--agent-id`, `--heatmap-id`, `--ai`, `--spec`, `--task`), ambient fallbacks, token stripping, auto-ingestion commit hooks. | `cli/cmdcommit/commit_types.go`<br>`cli/cmdcommit/commit_flags.go`<br>`cli/cmdcommit/commit_cmd.go`<br>`cli/cmdcommit/commit_push.go`<br>`cli/cmdcommit/commit_batch.go`<br>`cli/cmdcommit/commit_history_hook.go` | Worker 2 | **[YES]** Understood: parses flags across all commit commands, strips them from message, and records to DB on commit success. | `COMPLETED` |
| **ST-05** | `.../05-subtask-guidelines-and-skills-sync.md` | Author guideline `12-commit-provenance-and-tracking.md`, update guidelines catalog, update `.agents/skills/gitmap/skill.md`. | `02-spec/02-coding-guidelines/06-ai-optimization/`<br>`.agents/skills/gitmap/skill.md` | Worker 4 | **[YES]** Understood: author canonical guideline for commit tracking and document all new commands/flags in GitMap skills. | `COMPLETED` |
| **ST-06** | `.../06-subtask-v6-prompt-sync-and-release.md` | Update V6 prompt and skill to instruct AI agents on commit flags, sync connected repos, targeted linters, SemVer minor release `v6.593.0`. | `01-prompts/14-execute/`<br>`.agents/skills/execute-parent-task-with-n-steps-v6/`<br>`cli/constants/constants.go`<br>Release metadata | Worker 5 / Lead | **[YES]** Understood: update V6 prompts, synchronize repos, verify clean linters, and commit atomic minor release via `gitmap cpf`. | `COMPLETED` |

---

## 4. Verification Checkpoints & Definition of Done

All verification checkpoints evaluated to 100% green:
- [x] **CP-1 (Specifications & Plans)**: All modular specifications (00, 01, 02, 99) and 6 subtask plans are complete, consistent, and committed.
- [x] **CP-2 (Split-DB Sub-Millisecond Queries)**: `repodb/commit_history.db` operates outside the repo in WAL mode, verified with query latency $< 0.5\text{ ms}$.
- [x] **CP-3 (Universal SHA Resolution)**: `gitmap commit info <sha>` resolves full SHAs, 7-char short SHAs, prefixes, suffixes, and substrings accurately.
- [x] **CP-4 (Origin Attribution)**: Local AI commits display model and agent details; external commits state `"The current machine's AI has not committed this (likely another worker/node/human)"`.
- [x] **CP-5 (Provenance Flags & Stripping)**: All commit commands parse `--agent-id`, `--heatmap-id`, `--task` and ambient environment variables, cleanly stripping tokens before committing.
- [x] **CP-6 (Guidelines & Skills Sync)**: Guideline 12 authored, GitMap skill updated, and V6 execution prompt updated to mandate AI commit attribution.
- [x] **CP-7 (Targeted Linters)**: `check-relative-paths.py`, `check-enum-and-boolean.py`, and `check-forbidden-strings.py` pass with 100% compliance.
- [x] **CP-8 (Minor Release Ceremony)**: Version bumped to `v6.593.0`, release notes published, changelogs updated, and atomic commit pushed via `gitmap cpf`.
