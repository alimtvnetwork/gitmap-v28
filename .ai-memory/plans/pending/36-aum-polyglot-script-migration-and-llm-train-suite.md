# Plan 36: AUM Polyglot Script Migration & LLM Train Chained Curriculum Suite

> **Origin & Execution Context:** Derived directly from user mandate and **Spec 128** (`02-spec/21-app/128-aum-automation-suite-and-roadmap.md`) and **Spec 127** (`02-spec/21-app/127-llm-train-and-chained-agent-curriculum.md`). Establishes the full architectural migration of all 44 repository automation scripts from `03-ai-scripts/` into the compiled Go `aum` (**A**utomation and **U**tility **M**anager) subsystem, and details the 2-agent parallel execution plan.
> **Status:** Pending Execution (Phases 2 through 6)
> **Active Baseline & Phase 1:** ✅ Completed & Verified in Production

---

## 1. Executive Summary & Dual Sub-Agent Division

To maximize execution throughput while eliminating concurrency collisions, the remaining 5 phases are divided strictly between **two parallel sub-agents**:

```
                              ┌────────────────────────────────────────────────────────┐
                              │                 Supervisor Controller                  │
                              └──────────────────────────┬─────────────────────────────┘
                                                         │
                        ┌────────────────────────────────┴────────────────────────────────┐
                        ▼                                                                 ▼
      ┌────────────────────────────────────┐                           ┌────────────────────────────────────┐
      │     Sub-Agent 1: Systems & DB      │                           │   Sub-Agent 2: Quality & Release   │
      ├────────────────────────────────────┤                           ├────────────────────────────────────┤
      │ • Phase 2: Database & Topology     │                           │ • Phase 4: Release & SemVer Sync   │
      │   - aum topology                   │                           │   - aum version-sync               │
      │   - aum db-generate                │                           │   - aum release-bump               │
      │   - aum db-migrate                 │                           │   - aum milestones                 │
      │   - aum schema-audit               │                           │                                    │
      │ • Phase 3: CI/CD & Multi-Core      │                           │ • Phase 5: Docs & Memory           │
      │   - aum preflight                  │                           │   - aum help-audit                 │
      │   - aum test-inventory             │                           │   - aum plan-consolidate           │
      │   - aum purge-actions              │                           │   - aum doc-links                  │
      │   - aum smoke-test                 │                           │   - aum spec-migrate               │
      │                                    │                           │ • Phase 6: Git Hygiene & Purge     │
      │                                    │                           │   - aum clean-artifacts            │
      │                                    │                           │   - aum changed-files              │
      │                                    │                           │   - aum purge-history              │
      │                                    │                           │   - aum format-go                  │
      └────────────────────────────────────┘                           └────────────────────────────────────┘
```

---

## 2. Subtask Breakdown Matrix

| Subtask ID | File Path | Scope | Assigned Agent |
|---|---|---|---|
| **Subtask 01** | `subtasks/36-aum-polyglot-script-migration-and-llm-train-suite/01-phase2-database-topology-and-migration-engine.md` | Phase 2: Topology, SQLite schemas, code generators, migrations | **Sub-Agent 1** (Systems/DB) |
| **Subtask 02** | `subtasks/36-aum-polyglot-script-migration-and-llm-train-suite/02-phase3-cicd-local-preflight-and-multi-core-checkers.md` | Phase 3: Local preflight CI runner, test inventory, actions purger | **Sub-Agent 1** (Systems/DB) |
| **Subtask 03** | `subtasks/36-aum-polyglot-script-migration-and-llm-train-suite/03-phase4-release-semver-and-version-sync.md` | Phase 4: Version synchronization, SemVer release bumper, milestones | **Sub-Agent 2** (Quality/Release) |
| **Subtask 04** | `subtasks/36-aum-polyglot-script-migration-and-llm-train-suite/04-phase5-documentation-and-memory-consolidation.md` | Phase 5: CLI help AST auditor, memory plan consolidator, link validator | **Sub-Agent 2** (Quality/Release) |
| **Subtask 05** | `subtasks/36-aum-polyglot-script-migration-and-llm-train-suite/05-phase6-git-hygiene-cleanup-and-history-purging.md` | Phase 6: Artifact purger, changed file detector, history cleaner | **Sub-Agent 2** (Quality/Release) |

---

## 3. Strict Coding & Verification Mandates
1. **Mandatory Pre-Flight Pull:** Always run `git pull` before modifying code.
2. **File & Function Sizing:** Functions <= 8–15 lines; files <= 100–200 lines.
3. **Control Flow Flattening:** Zero nested if statements (nesting depth > 1 is strictly forbidden).
4. **Error Management:** Universal `*apperror.AppError` and `result.Result[T]` return envelopes; zero swallowed errors.
5. **No Intermediate Builds:** Sub-agents must NOT run `go build` or `go test` in the middle of editing loops. Only the Supervisor Controller performs pre-commit validation.
