---
name: nuclear-package-modularization
description: Autonomously modularize monolithic Go packages into smaller acyclic subpackages, isolate heavy subprocess tests into cli/tests/heavy_test, and estimate test inventory durations.
---

# Nuclear Package Modularization & Heavy Test Isolation

## Purpose & Overview

Autonomously decompose large, monolithic Go packages into small, acyclic packages following a strict Directed Acyclic Graph (DAG) architecture to eliminate compiler bloat and test latency, while isolating heavy subprocess tests into blackbox test packages (`tests/heavy_test/` under `package heavy_test`), maintaining an accurate test inventory (`.ai-memory/test-inventory.json`), and enforcing the 5-day cache freshness decision engine.

## Core Rules & Execution Flow

1. **Verbatim Prompt Recording & Deliverables Extraction (Phase 1 Step 0):**
   - Directly capture user prompt verbatim into `.ai-memory/plans/pending/xx-<slug>.md` under `## User Request (Verbatim)`.
   - Extract actionable deliverables under `## Extracted Actionable Task List`.
   - Output confirmed deliverables in chat before proceeding.

2. **5-Day Cache Freshness Decision Engine:**
   - Run `python 03-ai-scripts/33-test-inventory-generator.py --check-age --max-age-days 5`.
   - If exit code is 0 (fresh <= 5 days with profiled durations): AI skips re-running all tests and ingests cached durations directly.
   - If exit code is 1 (stale, missing, or unprofiled): AI executes baseline inventory discovery/profiling pass to populate timings.

3. **Strict DAG Architecture:**
   - Zero Go circular dependencies (`import cycle not allowed`).
   - Root packages (`main.go`, `cmd`) dispatch downwards to domain subpackages.
   - Domain subpackages only import leaf packages (`pkg/constants`, `pkg/model`, `pkg/appfault`, `pkg/fsutil`, `pkg/cliexit`).
   - Leaf packages have zero domain dependencies and never import caller packages.

4. **Heavy Test Isolation:**
   - Tests invoking `exec.Command`, external git CLI processes, sockets, or `time.Sleep` belong in `tests/heavy_test/` (or `cli/tests/heavy_test/`) under `package heavy_test`.
   - Routine package unit tests must be 100% in-memory fast unit tests (< 0.05s).

5. **Per-Task Agent Isolation & Workspace Subfolders:**
   - Create `.ai-memory/temp-agents/xx-<task-name>/` on task initialization.
   - Track progress and state in `state.md`.
   - On completion, record `STATUS: DONE` and clean up.

6. **Test Inventory Synchronization & Atomic Change Recording:**
   - Synchronize `.ai-memory/test-inventory.json` with test IDs, packages, durations, and tier tags.
   - Record modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.

7. **Coding Guidelines & Quality Gates:**
   - Functions <= 8–15 lines, blank lines before `return` statements and before `if` conditions.
   - Affirmative boolean prefixes (`is*`, `has*`), zero explicit `== true`, zero negative names.
   - Universal `*appfault.AppError` wrapping.
   - Verify with `go vet ./...` and `go build ./...` strictly at the final step (exit 0).
