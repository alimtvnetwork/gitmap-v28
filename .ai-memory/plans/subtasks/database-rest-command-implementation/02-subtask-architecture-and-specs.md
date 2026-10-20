# Subtask 02: Architecture & Specification Authoring

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/database-rest-command-implementation.md`  
> **Status:** DONE
> **Target Subsystems:** `02-spec/21-app/database-rest-command-implementation/`  
> **Owned Files:**  
> - `02-spec/21-app/database-rest-command-implementation/readme.md`  
> - `02-spec/21-app/database-rest-command-implementation/01-architecture-spec.md`  
> - `02-spec/21-app/database-rest-command-implementation/02-cli-and-rest-api-spec.md`  
> - `02-spec/21-app/database-rest-command-implementation/03-swagger-ui-and-verification-spec.md`  

---

## 1. Concrete Objectives

1. **Specification Suite Index & Overview (`readme.md`):**
   - Provide executive summary of the `db-rest` subsystem, supported CLI syntax, specification suite directory, and mandatory compliance invariants.

2. **System Architecture Specification (`01-architecture-spec.md` - Spec Author 01):**
   - Author comprehensive architectural blueprint covering:
     - End-to-end subsystem topology and Mermaid flowchart.
     - Multi-database target resolution for single, multiple, and mixed aliases/paths.
     - Sequential port probing and fallback engine over comma-separated candidate ports.
     - Split-DB migration and alteration caching architecture (`.gitmap/repodb/dbrest_cache.db`).
     - Concurrency governance, single-writer SQLite connection safety (`MaxOpenConns(1)`), in-memory row buffering, and WAL mode.
     - Universal Response Envelope (`{ Status, Attributes, Results }`) with PascalCase keys and positive booleans (`IsActive`, `IsSuccess`, `HasCompleted`).
     - REST routing topology and CLI companion commands architecture (`gitmap db-record-del`, `gitmap db-rm`).

3. **CLI Commands & REST API Contract Specification (`02-cli-and-rest-api-spec.md` - Spec Author 02):**
   - Author detailed interface contracts covering:
     - CLI syntax and argument parsing for `gitmap db-rest`, `gitmap db-record-del`, and `gitmap db-rm`.
     - REST API request and response payload schemas for all endpoints:
       - Catalog & Meta: `GET /api/v1/databases`
       - Schema & Tables: `GET /api/v1/db/{alias}/tables`, `GET /api/v1/db/{alias}/schema`
       - Column Discovery: `GET /api/v1/db/{alias}/find-columns`
       - Record CRUD: `GET/POST /api/v1/db/{alias}/{table}`, `GET/PUT/DELETE /api/v1/db/{alias}/{table}/{id}`
       - Direct SQL: `GET/POST /api/v1/db/{alias}/query`
       - Dynamic Views: `GET/POST/DELETE /api/v1/db/{alias}/views`
       - DDL Alterations: `POST /api/v1/db/{alias}/alter`, `POST /api/v1/db/{alias}/tables/create`
       - Migrations: `GET/POST /api/v1/db/{alias}/migrations`, `POST /api/v1/db/{alias}/migrate`
       - Database Clone: `POST /api/v1/db/{alias}/clone`

4. **Dynamic OpenAPI, Swagger & Visualizer UI Specification (`03-swagger-ui-and-verification-spec.md` - Spec Author 03):**
   - Author technical specifications covering:
     - On-the-fly OpenAPI 3.0 document generation from live SQLite schemas (`/swagger/doc.json`).
     - Embedded Swagger UI web asset serving (`/swagger`).
     - Custom interactive Database Visualizer dashboard (`/ui`) supporting table navigation, column introspection, row CRUD modals, dynamic view designer, and direct SQL console.
     - Verification gates, acceptance test cases, and error edge scenarios.

5. **Cross-Specification Invariants & Guardrails:**
   - Strict Relative Git Paths ONLY (`02-spec/...`, `.ai-memory/...`, `cli/...`). TOTAL BAN on absolute filesystem paths and `file:///` URIs.
   - Strict lowercase filenames across all files.
   - PascalCase table and column names across all schema definitions and JSON envelopes.
   - Positive boolean conventions (`IsActive`, `IsSuccess`, `HasCompleted`, `IsReadOnly`).
   - Single-writer SQLite connection safety (`MaxOpenConns(1)`).
   - Zero external web framework dependencies (standard library `net/http` only).
   - Pure-Go CGO-free SQLite driver (`modernc.org/sqlite`).

---

## 2. Work Breakdown Matrix (Spec Authors A=3)

| Role | Assigned Document | Subsystem Focus | Status |
|---|---|---|---|
| **Spec Author 01** | `01-architecture-spec.md` | System topology, multi-DB engine, port probe fallback, Split-DB caching, SQLite invariants | DONE |
| **Spec Author 01** | `01-subtask-git-sync-and-preflight.md` | Task preflight, ledger registration, lifecycle budget | DONE |
| **Spec Author 01** | `02-subtask-architecture-and-specs.md` | Architecture and spec authoring subtask decomposition | DONE |
| **Spec Author 02** | `02-cli-and-rest-api-spec.md` | CLI parsing, route handlers, Universal Response Envelope, CRUD, DDL, views, clone | IN_PROGRESS |
| **Spec Author 03** | `03-swagger-ui-and-verification-spec.md` | Dynamic OpenAPI 3.0 generator, Swagger UI, visualizer dashboard UI, acceptance criteria | IN_PROGRESS |

---

## 3. Action Items Checklist

- [x] 1. Author specification suite index in `02-spec/21-app/database-rest-command-implementation/readme.md`.
- [x] 2. Author master architecture specification in `02-spec/21-app/database-rest-command-implementation/01-architecture-spec.md` (Spec Author 01).
- [x] 3. Author Subtask 01 plan in `.ai-memory/plans/subtasks/database-rest-command-implementation/01-subtask-git-sync-and-preflight.md` (Spec Author 01).
- [x] 4. Author Subtask 02 plan in `.ai-memory/plans/subtasks/database-rest-command-implementation/02-subtask-architecture-and-specs.md` (Spec Author 01).
- [ ] 5. Author CLI and REST API specification in `02-spec/21-app/database-rest-command-implementation/02-cli-and-rest-api-spec.md` (Spec Author 02).
- [ ] 6. Author Swagger, Visualizer UI, and verification specification in `02-spec/21-app/database-rest-command-implementation/03-swagger-ui-and-verification-spec.md` (Spec Author 03).
- [ ] 7. Cross-specification compliance audit across all 3 specifications against quality gates G1–G6.

---

## 4. Verification Evidence

- `02-spec/21-app/database-rest-command-implementation/01-architecture-spec.md` authoring complete and validated.
- `.ai-memory/plans/subtasks/database-rest-command-implementation/01-subtask-git-sync-and-preflight.md` authoring complete and validated.
- `.ai-memory/plans/subtasks/database-rest-command-implementation/02-subtask-architecture-and-specs.md` authoring complete and validated.
- All documents enforce strict relative Git paths, lowercase filenames, PascalCase schema/keys, positive booleans, and single-writer SQLite connection safety.
