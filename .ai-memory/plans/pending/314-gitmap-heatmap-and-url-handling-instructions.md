# Master Plan: GitMap Heatmap Discovery, URL Split-DB Caching, Segmented Downloader & LLM Training Instructions

## Master Plan Metadata
- **Plan File**: `.ai-memory/plans/pending/314-gitmap-heatmap-and-url-handling-instructions.md`
- **Task Identifier**: Task 314 (`gitmap-heatmap-and-url-handling-instructions`)
- **Status**: `QUEUED / READY FOR EXECUTION`
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
├── 01-subtask-heatmap-discovery-and-help.md         (Wave 1: Worker 1)
├── 02-subtask-recent-tasks-and-suggestions.md        (Wave 1: Worker 2)
├── 03-subtask-url-cache-split-db.md                  (Wave 1: Worker 3)
├── 04-subtask-segmented-downloader-and-aria2.md      (Wave 2: Worker 1)
├── 05-subtask-url-cli-command.md                     (Wave 2: Worker 2)
├── 06-subtask-llm-train-heatmap-linking.md           (Wave 2: Worker 3)
└── 07-subtask-coding-guidelines-and-skills.md        (Wave 3: Worker 1 & Lead)
```

### Detailed Subtask Disjoint Bounding Boxes

| Subtask ID | Subtask File | Scope & Deliverables | Exclusive File Ownership | Dependencies |
| :--- | :--- | :--- | :--- | :--- |
| **ST-01** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/01-subtask-heatmap-discovery-and-help.md` | Heatmap discovery scanner, multi-tier database inventory, flag routing fix, alias joins via `db.ListAliasesWithRepo()`. | `cli/heatmap/discovery.go`<br>`cli/cmd/special_repos_rc_heatmap.go` | None |
| **ST-02** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/02-subtask-recent-tasks-and-suggestions.md` | Multi-repo recent task aggregator (2–5 repos, min 5 tasks), formatted CLI table, dynamic AI suggestion pills. | `cli/heatmap/recent_tasks.go`<br>`cli/heatmap/suggestions.go` | ST-01 |
| **ST-03** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/03-subtask-url-cache-split-db.md` | Split-DB `repodb/url_cache.db` with WAL mode, single-writer mutex, URL SHA-256 hashing, ETag/mtime metadata tracking. | `cli/store/url_cache_split_db.go`<br>`cli/store/url_cache_types.go` | None |
| **ST-04** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/04-subtask-segmented-downloader-and-aria2.md` | 10-split segmented HTTP Range downloader, host `aria2c` detection and invocation, native Go 10-connection parallel Range fallback. | `cli/downloader/segmented.go`<br>`cli/downloader/aria2.go` | None |
| **ST-05** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/05-subtask-url-cli-command.md` | CLI command `gitmap url read <url>` streaming unbuffered bytes via `output.Raw()`, `--no-cache`, `--refresh`, `--json`. | `cli/cmdurl/url_read.go`<br>`cli/root.go` dispatch | ST-03, ST-04 |
| **ST-06** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/06-subtask-llm-train-heatmap-linking.md` | Heatmap telemetry linking in `gitmap llm train` / `gitmap llm train muse` via companion file or live REST URL. | `cli/cmd/llm/llm_heatmap.go`<br>`cli/cmd/llm/llm_train.go` | ST-01, ST-03 |
| **ST-07** | `.ai-memory/plans/subtasks/314-gitmap-heatmap-and-url-handling-instructions/07-subtask-coding-guidelines-and-skills.md` | Author coding guidelines 10 & 11, update `.agents/skills/gitmap/skill.md`, run relative path linter, execute minor bump release. | `02-spec/02-coding-guidelines/06-ai-optimization/`<br>`.agents/skills/gitmap/skill.md`<br>Release metadata | ST-01 through ST-06 |

---

## Multi-Phase Execution Timeline & Budgets

```mermaid
gantt
    title Task 314 Execution Timeline (400 Step Budget)
    dateFormat  X
    axisFormat %s

    section Phase 1 (Budget 100)
    Research & Discovery Subagents       :active, p1_1, 0, 35
    Modular Specifications Authoring     :active, p1_2, 35, 75
    Subtask Plans & Consistency Matrices :active, p1_3, 75, 100

    section Phase 2 (Budget 150)
    Wave 1: Heatmap & URL Split-DB Core  :crit, p2_1, 100, 160
    Wave 2: Downloader & CLI & LLM Train :p2_2, 160, 220
    Wave 3: Coding Standards & Skills     :p2_3, 220, 250

    section Phase 3 (Budget 150)
    Decoupled Unit Testing & Linter Sync :p3_1, 250, 310
    Retrospective AI Verification Audit  :p3_2, 310, 360
    Minor Bump Release & Atomic Push     :p3_3, 360, 400
```

### Phase Breakdown

#### Phase 1: Planning, Research & Specification (Steps 1–100)
- **Goal**: Establish comprehensive architectural specifications and bounded subtasks with zero builds, zero tests, and zero git commands.
- **Milestones**:
  - Research subagents complete technical investigations across existing codebase.
  - Author specs 00, 01, 02, and 99 under `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/`.
  - Author master plan in `.ai-memory/plans/pending/` and 7 subtasks in `.ai-memory/plans/subtasks/`.
  - Validate 100% relative path hygiene and coding guideline compliance in specs.

#### Phase 2: Parallel Implementation Waves (Steps 101–250)
- **Goal**: Implement features using bounded micro-batches with strictly disjoint file boundaries.
- **Wave 1 (Concurrency = 3 workers)**:
  - Worker 1 executes ST-01: Heatmap discovery scanner & `--repo` flag routing fix.
  - Worker 2 executes ST-02: Multi-repo recent tasks aggregator & AI suggestions.
  - Worker 3 executes ST-03: Dedicated Split-DB `repodb/url_cache.db`.
- **Wave 2 (Concurrency = 3 workers)**:
  - Worker 1 executes ST-04: 10-split segmented HTTP Range & `aria2c` downloader.
  - Worker 2 executes ST-05: High-speed `gitmap url read` CLI command.
  - Worker 3 executes ST-06: LLM training heatmap telemetry linker.
- **Wave 3 (Concurrency = 2 workers)**:
  - Worker 1 executes ST-07 Part A: Author coding guidelines 10 & 11 and update `.agents/skills/gitmap/skill.md`.
  - Lead Orchestrator validates micro-batch completions and synchronizes shared registries.

#### Phase 3: Decoupled Verification, Release Ceremony & Atomic Push (Steps 251–400)
- **Goal**: Verify quality gates, execute SemVer minor bump release, and commit via GitMap.
- **Milestones**:
  - Run decoupled unit test suites (`cli/heatmap/discovery_test.go`, `cli/store/url_cache_test.go`, `cli/downloader/segmented_test.go`).
  - Run relative path and boolean linter checks (`03-ai-scripts/01-check-relative-paths.py`).
  - Execute SemVer minor version bump (`gitmap release` / `minor-bump`).
  - Perform single atomic GitMap commit: `gitmap cpf "heatmap-url-cache - implement heatmap discovery, url split-db caching, segmented downloader, and llm training integration"`.

---

## Rollback Plans & Safety Invariants

1. **Split-DB File Isolation**: `repodb/url_cache.db` is completely partitioned from core system databases. If cache corruption occurs, simply removing the single `.db` file triggers safe recreation.
2. **Network Degradation Resilience**: If `aria2c` fails or is not present, the native Go 10-worker Range downloader operates autonomously. If the server does not support HTTP Range headers (HTTP 200 returned), the engine seamlessly falls back to standard single-stream reading.
3. **Heatmap Discovery Safety**: Heatmap discovery only executes read-only SQL queries (`SELECT`) with SQLite busy timeouts set to 5000ms, preventing any locking of active agent heatmap tracking sessions.
4. **Worker Git Safety Invariant**: Subagents are strictly prohibited from invoking `git` commands, preventing `.git/index.lock` collisions.

---

## Verification Checkpoints & Quality Gates

- [ ] **Gate 1 (Spec Gate)**: All 4 specification files and 7 subtasks authored with 100% relative paths and positive booleans.
- [ ] **Gate 2 (Discovery Gate)**: `gitmap heatmap` and `gitmap rc heatmap` execute without error, inventory 3 DB tiers, list aliases, and render recent tasks and AI suggestions.
- [ ] **Gate 3 (Cache Gate)**: `gitmap url read <url>` caches payload in `repodb/url_cache.db` and satisfies repeated requests in $< 0.5\text{ ms}$.
- [ ] **Gate 4 (Downloader Gate)**: Files $\ge 1\text{--}2\text{ MB}$ download via 10 parallel segments with `aria2c` detection.
- [ ] **Gate 5 (LLM Train Gate)**: `gitmap llm train --heatmap <path|url>` successfully ingests decision traces into the training pipeline.
- [ ] **Gate 6 (Release Gate)**: Clean linter run, minor version bump incremented, and atomic push completed.
