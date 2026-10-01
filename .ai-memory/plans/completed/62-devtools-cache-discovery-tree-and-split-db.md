# [V6] Plan 62: Devtools Cache Dynamic Discovery, Tree View Rendering & Split-DB Persistence

```text
N = 300 (Total self-loop steps budget — editable top-header parameter, default: 300)
A = 2   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
PHASE_1_BUDGET = N / 2   (Steps 1 .. 150: Planning, Parallel Discovery Subagents, Detailed Spec, and Lean Subtask Generation)
PHASE_2_BUDGET = N / 2   (Steps 151 .. 300: Mandatory Parallel Subagent Execution, Self-Looping, Targeted Quality Linting)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> **Plan Slug:** `62-devtools-cache-discovery-tree-and-split-db`
> **Tracking Spec:** `02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md`
> **Runtime:** Google Antigravity 2.0 (IDE and CLI)

---

## User Request (Verbatim)

```text
C:\dev-tool\go\cache

This Git map clean or clear cache output for the dev tools is clearly stupid. First of all, you did not remove the cache folder of the Git map, which is in my machine in this location. I'm giving the location, but it will be in different ones for the different type of, let's say, setup or OS. So you need to find that. And after you find that, you actually keep these paths backed up in a SQLite DB so that you can access it very quickly. Okay. But if there is a change, the user can actually do a force, hyphen, hyphen force that would actually reset or refine the paths as well. Now, coming to the point. I think these paths can be fine. Okay, so no need to do a cache because it can change later on. So that would create another issue, okay. Now, coming to the point that the output is very terrible. As you can see, no coloring, nothing. Okay. And then you didn't do it right. It's just 0.1 megabytes, but there is so much more that you could not find those paths. And these paths needs to be listed like which ones you are removing. You did not list out, and you could list out as a tree view as well. So overall is very disappointing. I think you need to fix it. Can you please do that for me?
```

---

## Root Cause Analysis (4-Part RCA)

1. **Root Cause:**
   - Pre-purge subprocess execution (`go clean`, `npm cache clean`, etc.) wiped data before directory sweep measured sizes, leaving empty directories and recording only 0.01 MB.
   - Path resolution only checked static `go/pkg/mod` and missed `C:\dev-tool\go\cache` (752.68 MB) and `C:\dev-tool\go\pkg\mod` (711.71 MB).
   - Dynamic tools (`pnpm store path`, `npm config get cache`, `pip cache dir`) and environment variables (`GOCACHE`, `GOMODCACHE`, `GOPATH`, `CARGO_HOME`) were not queried.
2. **Blast Radius:**
   - `cli/cmdos/devclean_cmd.go` and `cli/cmdos/os_dev_clean.go`.
   - `cli/osclean/` cleaners and metric calculations.
   - All `gitmap clear devtools`, `gitmap clean devtools`, and `gitmap devtool clear` invocations.
3. **Remediation Strategy:**
   - Measure directory sizes *before* cleaning or purge.
   - Implement 3-tier dynamic discovery (CLI query, environment variables, multi-drive heuristics `C:\dev-tool`, `D:\dev-tool`, etc.).
   - Persist discovered paths in Split-DB (`.gitmap/data/devtools/cache/sql.db`) with `--force` invalidation.
   - Add rich ANSI colored cards, itemized path lists, and expandable directory tree view (`--tree`, `-t`).
4. **Prevention & Regression Gates:**
   - Strict coding guidelines (< 100 lines per file, <= 15 lines per function, affirmative booleans).
   - Zero builds / zero tests (Rule R1).
   - Verify with targeted Python linters and harmless CLI executions.

---

## Subtasks Breakdown (A = 2, H = 2)

| Task-ID | Subtask | Owner | Owned Files | Status |
| :--- | :--- | :--- | :--- | :--- |
| `Task-01` | Multi-Ecosystem Dynamic Discovery Engine | Worker 01 | `cli/cmdos/os_dev_clean_paths.go`, `cli/cmdos/os_dev_clean_discover.go` | COMPLETED |
| `Task-02` | Split-DB Cache Persistence & `--force` Invalidation | Worker 02 | `cli/store/devtools_cache_split_db.go` | COMPLETED |
| `Task-03` | ANSI Colored Rich Output & Interactive Tree View | Worker 01 | `cli/cmdos/os_dev_clean_render.go`, `cli/cmdos/os_dev_clean_tree.go` | COMPLETED |
| `Task-04` | CLI Flags Integration (`--force`, `--tree`, `--only`) & Routing | Worker 02 | `cli/cmdos/os_dev_clean.go`, `cli/cmdos/devclean_cmd.go` | COMPLETED |

---

## Verification & Acceptance Criteria
- AC-1: `gitmap clear devtools --dry-run` discovers `C:\dev-tool\go\cache`, `C:\dev-tool\go\pkg\mod`, `pnpm store`, and reports realistic sizes (> 1 GB).
- AC-2: `--tree` (`-t`) renders an ANSI colored hierarchical directory tree showing path nodes, file counts, and byte sizes.
- AC-3: Split-DB correctly persists paths in `devtools_cache_paths` and `--force` (`-f`) resets and re-probes fresh paths.
- AC-4: Itemized path list clearly shows each directory being cleaned.
- AC-5: Passes `python 03-ai-scripts/05-guideline-autofixer.py` with zero errors.
