# Master Plan: GitMap Heatmap Discovery, URL Split-DB Caching, Segmented Downloader & LLM Training Instructions

## Master Plan Metadata
- **Plan File**: `.ai-memory/plans/completed/314-gitmap-heatmap-and-url-handling-instructions.md`
- **Task Identifier**: Task 314 (`gitmap-heatmap-and-url-handling-instructions`)
- **Status**: `COMPLETED`
- **Specification References**:
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/00-overview.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/02-component-spec.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/99-consistency-report.md`
- **Execution Swarm Concurrency**: Lead Orchestrator -> 3 Parallel Subagents ($A=3, H=2$, Concurrency Capacity = 6)
- **Total Step Budget**: $N=400$ (Phase 1: 100, Phase 2: 150, Phase 3: 150)

---

## Executive Summary & Strategic Objectives

Task 314 delivers four integrated capabilities to the GitMap developer companion:
1. **Heatmap Discovery & AI Help**: Unifies `gitmap heatmap` and `gitmap rc heatmap` into a self-documenting CLI suite that catalogs multi-tier heatmap databases, maps repository aliases via `db.ListAliasesWithRepo()`, aggregates recent execution tasks across 2–5 repositories (minimum 5 tasks), and provides interactive AI terminal suggestions.
2. **Dedicated URL Split-DB Caching**: Implements `repodb/url_cache.db` operating in WAL mode with single-writer mutex synchronization and memory mapping (`PRAGMA mmap_size = 268435456`). Parallel worker URL requests stream in $< 0.5\text{ ms}$ directly to stdout via `output.Raw()` without duplicate network downloads.
3. **10-Split Segmented & Aria2 Downloader**: Automatically partitions downloads exceeding 1–2 MB into $\ge 10$ parallel Range segments (or 5KB chunks), detects host `aria2c` (`aria2c -s 10 -x 10 -k 5K`), and provides seamless fallback to native Go parallel Range workers.
4. **LLM Training Integration & Guidelines**: Connects heatmap companion telemetry (`repo-cache/.heatmap/<seq>-<ai>/<repo>/heatmap.db` and live REST endpoints) into `gitmap llm train` / `gitmap llm train muse`, expands coding guidelines with anti-patterns and caching standards, updates `.agents/skills/gitmap/skill.md`, and executes a SemVer minor release.

---

## Subtask Topology & Disjoint Ownership

The master plan is partitioned into 7 strictly disjoint subtasks across three execution waves:

```text
.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/
├── 01-subtask-heatmap-discovery-and-help.md         (COMPLETED: Worker 1)
├── 02-subtask-recent-tasks-and-suggestions.md        (COMPLETED: Worker 1)
├── 03-subtask-url-cache-split-db.md                  (COMPLETED: Worker 2)
├── 04-subtask-segmented-downloader-and-aria2.md      (COMPLETED: Worker 3)
├── 05-subtask-url-cli-command.md                     (COMPLETED: Worker 2)
├── 06-subtask-llm-train-heatmap-linking.md           (COMPLETED: Worker 4)
└── 07-subtask-coding-guidelines-and-skills.md        (COMPLETED: Worker 5 & Lead)
```

### Detailed Subtask Disjoint Bounding Boxes

| Subtask ID | Subtask File | Scope & Deliverables | Exclusive File Ownership | Status |
| :--- | :--- | :--- | :--- | :--- |
| **ST-01** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/01-subtask-heatmap-discovery-and-help.md` | Heatmap discovery scanner, multi-tier database inventory, flag routing fix, alias joins via `db.ListAliasesWithRepo()`. | `cli/heatmap/heatmap_discovery.go`<br>`cli/cmd/special_repos_rc_heatmap.go` | COMPLETED |
| **ST-02** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/02-subtask-recent-tasks-and-suggestions.md` | Multi-repo recent task aggregator (2–5 repos, min 5 tasks), formatted CLI table, dynamic AI suggestion pills. | `cli/heatmap/heatmap_overview.go`<br>`cli/heatmap/heatmap_cmd.go` | COMPLETED |
| **ST-03** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/03-subtask-url-cache-split-db.md` | Split-DB `repodb/url_cache.db` with WAL mode, single-writer mutex, URL SHA-256 hashing, ETag/mtime metadata tracking. | `cli/store/url_cache_split_db.go`<br>`cli/store/url_cache_types.go`<br>`cli/store/split_db_path.go` | COMPLETED |
| **ST-04** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/04-subtask-segmented-downloader-and-aria2.md` | 10-split segmented HTTP Range downloader, host `aria2c` detection and invocation, native Go 10-connection parallel Range fallback. | `cli/downloader/segmented.go`<br>`cli/downloader/aria2.go`<br>`cli/downloader/download.go` | COMPLETED |
| **ST-05** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/05-subtask-url-cli-command.md` | CLI command `gitmap url read <url>` streaming unbuffered bytes via `output.Raw()`, `--no-cache`, `--refresh`, `--json`. | `cli/cmdurl/cmd_url.go`<br>`cli/cmd/rootutility.go` dispatch | COMPLETED |
| **ST-06** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/06-subtask-llm-train-heatmap-linking.md` | Heatmap telemetry linking in `gitmap llm train` / `gitmap llm train muse` via companion file or live REST URL. | `cli/cmd/llm/llm_heatmap.go`<br>`cli/cmd/llm/llm_train.go`<br>`cli/cmdai/ai_train.go` | COMPLETED |
| **ST-07** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/07-subtask-coding-guidelines-and-skills.md` | Author coding guidelines 10 & 11, update `.agents/skills/gitmap/skill.md`, run relative path linter, execute minor bump release. | `02-spec/02-coding-guidelines/06-ai-optimization/`<br>`.agents/skills/gitmap/skill.md`<br>`cli/cmd/llm/llm_urls.go`<br>Release metadata | COMPLETED |

---

## Multi-Phase Execution Timeline & Budgets

```mermaid
gantt
    title Task 314 Execution Timeline (400 Step Budget)
    dateFormat  X
    axisFormat %s

    section Phase 1 (Budget 100)
    Research & Discovery Subagents       :done, p1_1, 0, 35
    Modular Specifications Authoring     :done, p1_2, 35, 75
    Subtask Plans & Consistency Matrices :done, p1_3, 75, 100

    section Phase 2 (Budget 150)
    Wave 1: Heatmap & URL Split-DB Core  :done, p2_1, 100, 160
    Wave 2: Downloader & CLI & LLM Train :done, p2_2, 160, 220
    Wave 3: Coding Standards & Skills     :done, p2_3, 220, 250

    section Phase 3 (Budget 150)
    Decoupled Unit Testing & Linter Sync :done, p3_1, 250, 310
    Retrospective AI Verification Audit  :done, p3_2, 310, 360
    Minor Bump Release & Atomic Push     :done, p3_3, 360, 400
```

---

## Verification Checkpoints & Quality Gates

- [x] **Gate 1 (Spec Gate)**: All 4 specification files and 7 subtasks authored with 100% relative paths and positive booleans.
- [x] **Gate 2 (Discovery Gate)**: `gitmap heatmap` and `gitmap rc heatmap` execute without error, inventory 3 DB tiers, list aliases, and render recent tasks and AI suggestions.
- [x] **Gate 3 (Cache Gate)**: `gitmap url read <url>` caches payload in `repodb/url_cache.db` and satisfies repeated requests in $< 0.5\text{ ms}$.
- [x] **Gate 4 (Downloader Gate)**: Files $\ge 1\text{--}2\text{ MB}$ download via 10 parallel segments with `aria2c` detection.
- [x] **Gate 5 (LLM Train Gate)**: `gitmap llm train --heatmap <path|url>` successfully ingests decision traces into the training pipeline.
- [x] **Gate 6 (Release Gate)**: Clean linter run, minor version bump incremented, and atomic push completed.
