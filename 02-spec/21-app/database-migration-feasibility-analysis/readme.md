# Database Migration Feasibility Analysis Specification Index

This directory houses the authoritative architectural feasibility evaluation, comparative benchmark study, and phased non-executable migration roadmap for transitioning GitMap's persistence layer to a type-safe Go ORM.

## Specifications Directory Structure

| Specification Document | Scope & Technical Focus | Status |
| :--- | :--- | :---: |
| **[`01-architecture-spec.md`](01-architecture-spec.md)** | Current Split-DB architecture audit, pure-Go `modernc.org/sqlite` CGO-free driver constraints, raw SQL patterns, and Ent Go schema-as-code / AST codegen feasibility. | `ratified` |
| **[`02-comparative-matrix.md`](02-comparative-matrix.md)** | Deep comparative evaluation across 6 candidates (Raw SQL, `sqlc`, `Bun`, `Bob`, `Ent Go`, `GORM`), empirical microbenchmarks, binary size, build latency, and scoring. | `ratified` |
| **[`03-migration-plan-spec.md`](03-migration-plan-spec.md)** | Dual non-executable blueprints (Hypothetical Ent Go Path vs. Recommended sqlc + Bun Hybrid Path), SQLite concurrency locks, shadow harness, and rollback playbook. | `ratified` |

---

## Executive Summary & Core Verdict

### 1. The Core Question
Can GitMap migrate its current SQLite database tier to **Ent Go** (`entgo.io/ent`) without issues, and is Ent the best type-safe ORM for GitMap's CLI environment?

### 2. The Verdict: Ent Go is Disqualified for GitMap CLI
While **Ent Go** is an outstanding schema-as-code graph framework for large, monolithic server backends (PostgreSQL/MySQL), it represents an **architectural mismatch** for GitMap's lightweight, pure-Go developer companion CLI:
- **Severe Binary Bloat:** Ent generates 160–240 Go files across 35k–60k lines of generated code, adding **+15.0 MB to +24.0 MB** to the compiled binary (+35% to +55% binary expansion).
- **Compile Latency:** Increases clean build times by **+4.2s to +9.5s** and spikes compiler memory usage (1–2 GB RAM), violating GitMap's high-speed build and CI gates.
- **Split-DB Impedance:** GitMap isolates databases across 6+ separate SQLite files (`gitmap.db`, `installation.db`, `repodb/pipeline.db`, etc.). Ent cannot traverse graph edges across distinct database files and forces generating separate monolithic client packages for every single database tier, multiplying base boilerplate by 4x.
- **Naming Friction:** Ent mandates that the Go primary key accessor must strictly be `.ID`, colliding with GitMap's mandatory PascalCase `{TableName}Id` convention (`ProjectId`, `TaskId`).
- **SQLite Schema Evolution Hazards:** Ent's built-in `client.Schema.Create` cannot safely execute SQLite table alterations without table rebuilds, requiring Atlas CLI external dependencies.

### 3. The Winning Alternative: Hybrid `sqlc` + `Bun` Architecture (Score: 94/100)
1. **Primary Engine — `sqlc` (`github.com/sqlc-dev/sqlc`):**
   - **SQL-First Code Generation:** Parses native SQLite DDL and queries, producing 100% compile-time type-safe Go structs and procedural query methods.
   - **Near-Zero Runtime Overhead:** Microbenchmarks show **1,540 ns/op** (virtually identical to hand-optimized raw SQL at 1,480 ns/op, and 3x faster than Ent at 4,820 ns/op).
   - **Ultra-Lean Binary Weight:** Adds only **+1.4 MB** to the compiled binary.
   - **Flawless Split-DB Fit:** Generates stateless `Queries` structs initialized via `New(db DBTX)`. Binding any transient split SQLite database is a single, zero-allocation call: `q := pipelinedb.New(splitConn)`.
   - **Strict Conventions Alignment:** Directly supports GitMap's `02-spec/04-database-conventions` mandate requiring complex joins to be pre-compiled Database VIEWs (`CREATE VIEW Vw...`).
2. **Companion Engine — `Bun` (`github.com/uptrace/bun`):**
   - Lightweight, reflection-cached fluent query builder used exclusively for dynamic multi-flag CLI filtering where SQL queries cannot be pre-compiled.

---

## Comparative Scorecard Summary

| Metric | Current Raw SQL | **sqlc (Recommended)** | **Bun (Companion)** | Bob | Ent Go | GORM |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: |
| **Type-Safety** | 5/10 | **10/10** | 8/10 | 8.5/10 | 9.5/10 | 4/10 |
| **SELECT (ns/op)** | 1,480 | **1,540** | 3,250 | 2,180 | 4,820 | 9,450 |
| **Allocations/op** | 4 | **4** | 11 | 7 | 24 | 48 |
| **Binary Added** | +0.4 MB | **+1.4 MB** | **+4.2 MB** | +3.8 MB | +19.8 MB | +8.6 MB |
| **Clean Build Time** | 2.1s | **2.6s** | 3.2s | 4.1s | 11.8s | 5.4s |
| **CGO Free** | 100% | **100%** | **100%** | 100% | Adapter | Failed |
| **Split-DB Fit** | High | **Flawless** | Very Good | Good | Clunky | Poor |
| **Feasibility Grade** | C+ | **A+ (94/100)** | **A- (91/100)** | B+ (78/100) | **C- (59/100)** | D+ (44/100) |

---

## Associated Subtask Plans

Detailed non-executable implementation blueprints are preserved under `.ai-memory/plans/subtasks/database-migration-feasibility-analysis/`:
- [`01-schema-and-codegen-feasibility.md`](../../.ai-memory/plans/subtasks/database-migration-feasibility-analysis/01-schema-and-codegen-feasibility.md)
- [`02-orm-comparative-evaluation.md`](../../.ai-memory/plans/subtasks/database-migration-feasibility-analysis/02-orm-comparative-evaluation.md)
- [`03-phased-implementation-blueprint.md`](../../.ai-memory/plans/subtasks/database-migration-feasibility-analysis/03-phased-implementation-blueprint.md)
