# Learned: Nuclear Package Modularization (Phase 6) & Heavy Test Isolation

## 1. Domain Clusters Extracted in Phase 6

- `cli/cmdfixgit` (11 domain files + helpers + exports):
  - Self-contained auto-repair pipelines: permission fixing, index corruption remediation, safe directory configuration, and index lock clearing.
  - Zero upward dependencies on `cmd`.
- `cli/cmdcg` (19 domain files + helpers + exports):
  - Coding guideline tooling: workspace auditor, prompt installer, multi-repo version manifest synchronizer.
  - Relies only on leaves (`model`, `store`, `apperror`, `fsutil`, `cmdprompt`).
- `cli/cmdssh` (43 domain files + helpers + exports):
  - SSH node manager: SSH client, authentication key generator/distributor, cluster join, alias commands, multi-IP broadcaster, and execution engine.
  - Decoupled from `join.go` and `profile.go` using callback injection (`JoinRunner`, `ProfileRunner`).

## 2. Heavy Test Isolation Guarantees

- `TestCloneAllPreservesNestedHierarchy` (in `cloner`) and `TestAgentIntegration` (in `cluster`) were relocated to `cli/tests/heavy_test/` (`package heavy_test`).
- Core package unit tests across all non-test packages (`cloner`, `cluster`, `gitutil`, `clonefrom`, etc.) now execute purely in memory without spawning subprocesses or opening network sockets.
- Heavy tests are categorized as tier: "slow" (5.0s) in `.lovable/test-inventory.json` (50 slow tests total).

## 3. Strict Acyclic DAG Architecture

```text
[constants, model, store, apperror, cliexit, db, dbengine, gitutil, cluster]
                                 ▲
            ┌────────────────────┼────────────────────┐
            │                    │                    │
     [cli/cmdfixgit]        [cli/cmdcg]          [cli/cmdssh]
            │                    │                    │
            └────────────────────┼────────────────────┘
                                 │
                           [cli/cmd] (Orchestrator)
                                 │
                            [main.go]
```
