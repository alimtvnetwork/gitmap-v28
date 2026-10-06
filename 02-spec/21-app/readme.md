# Application Specifications Index (`02-spec/21-app`)

This directory houses the authoritative application specifications for GitMap. In accordance with the deep consolidation architecture ([Spec 236](236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md)), all historical specifications and subsystems are consolidated into **Eight Canonical Domain Clusters**, representing the ratified architecture across all subsystems under the Final Version Invariant.

The top-level directory strictly maintains **exactly 10 filesystem items**: the 8 Canonical Architecture Domain Clusters, the active specification directory (`236-deep-spec-consolidation-and-canonical-reduction/`), and this master domain index (`readme.md`).

---

## 1. Canonical Architecture Domain Clusters

Each domain cluster represents the authoritative, consolidated specification for a major GitMap subsystem, embodying the Final Version Invariant (retaining ratified contracts and pruning intermediate drafts):

| Cluster | Subsystem Scope | Specifications | Status |
| :--- | :--- | :--- | :---: |
| **`01-cli-architecture/`** | Cobra CLI hierarchy, middle-ellipsized tables, shell tab completion, Win32 CP handoff, typed JSON Envelope V2 | [Architecture](01-cli-architecture/01-architecture-spec.md) • [Component](01-cli-architecture/02-component-spec.md) | `ratified` |
| **`02-scanner-and-projects/`** | High-speed filesystem discovery, project classification heuristics, `.gitmapignore` engine, AUM indexing, dedup, `EqualFoldAnyTrim` | [Architecture](02-scanner-and-projects/01-architecture-spec.md) • [Component](02-scanner-and-projects/02-component-spec.md) | `ratified` |
| **`03-git-operations-and-pull/`** | 8-worker concurrency pull pool, semantic flat commit suite (`gitmap c`), auto-staging, push self-healing, Oh-My-Zsh ignore rules | [Architecture](03-git-operations-and-pull/01-architecture-spec.md) • [Component](03-git-operations-and-pull/02-component-spec.md) | `ratified` |
| **`04-fleet-nodes-and-ssh/`** | Unified `gitmap nodes` CLI, dual-stack ICMP/TCP ping probing, RSA-OAEP salt credential vault, except-self remote clone, Ubuntu fleet | [Architecture](04-fleet-nodes-and-ssh/01-architecture-spec.md) • [Component](04-fleet-nodes-and-ssh/02-component-spec.md) | `ratified` |
| **`05-antigravity-and-ide/`** | Google Antigravity SDK workflows, multi-conversation prompt dispatch, theme parity, plugins/skills sync, Cursor IDE integration | [Architecture](05-antigravity-and-ide/01-architecture-spec.md) • [Component](05-antigravity-and-ide/02-component-spec.md) | `ratified` |
| **`06-database-and-split-db/`** | Three-tier SQLite Split-DB engine (`gitmap.db`, `installation.db`, `repodb/pipeline.db`), WAL mode, `SetMaxOpenConns(1)`, PascalCase schema, [Database ERD](06-database-and-split-db/gitmap-database-erd.mmd) | [Architecture](06-database-and-split-db/01-architecture-spec.md) • [Component](06-database-and-split-db/02-component-spec.md) | `ratified` |
| **`07-pipeline-and-diagnostics/`** | CI/CD pipeline telemetry (`pe`/`pea`), traceback extractors, heatmap failure summaries, dynamic ETA forecaster, 4-part RCA engine | [Architecture](07-pipeline-and-diagnostics/01-architecture-spec.md) • [Component](07-pipeline-and-diagnostics/02-component-spec.md) | `ratified` |
| **`08-distribution-and-release/`** | NSIS Windows installer payload auto-detection, Linux archives, SemVer release ceremony, cross-platform runners (`run.ps1`, `run.sh`) | [Architecture](08-distribution-and-release/01-architecture-spec.md) • [Component](08-distribution-and-release/02-component-spec.md) | `ratified` |

---

## 2. Active Application Specifications

The following application specification records the active feature development and consolidation milestone:

- [236-deep-spec-consolidation-and-canonical-reduction](236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md) — Deep Application Specifications Consolidation & Canonical Reduction (Specs: [Master Ledger](236-deep-spec-consolidation-and-canonical-reduction/00-master-audit-ledger.md), [Architecture](236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md), [Component](236-deep-spec-consolidation-and-canonical-reduction/02-component-and-cli-spec.md)) (Status: `in_progress`)

---

## 3. Architecture Governance & Consolidation Model

1. **The Final Version Invariant (A, B vs. X, Y):**
   - Historical specifications across 72 legacy folders and 181 standalone root files have been evaluated against the running codebase.
   - All ratified implementations (A, B) are fully synthesized into the 8 Canonical Domain Clusters.
   - Superseded exploratory drafts, speculative proposals, and intermediate spikes (X, Y) have been safely pruned.

2. **Asset Preservation:**
   - Visual diagrams and system schematics are preserved within their authoritative domain cluster (e.g., [Database ERD](06-database-and-split-db/gitmap-database-erd.mmd) in `06-database-and-split-db/`).

3. **Historical Auditability:**
   - Full historical versions of pre-consolidation specifications are preserved permanently in the Git commit history, release tag `v6.500.0`, and safety backup branch `backup/pre-deep-spec-consolidation-20261006`.
